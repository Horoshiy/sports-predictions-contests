package contest_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/contest-service/internal/models"
	"github.com/sports-prediction-contests/contest-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&models.Contest{}, &models.Participant{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestParticipantCountConsistency(t *testing.T) {
	db := setupTestDB(t)
	contestRepo := repository.NewContestRepository(db)
	participantRepo := repository.NewParticipantRepository(db)

	// Create a contest
	contest := &models.Contest{
		Title:           "Test Contest",
		Description:     "Test Description",
		SportType:       "football",
		Status:          "active",
		StartDate:       time.Now().Add(24 * time.Hour),
		EndDate:         time.Now().Add(48 * time.Hour),
		MaxParticipants: 100,
		CreatorID:       1,
	}

	if err := contestRepo.Create(contest); err != nil {
		t.Fatalf("Failed to create contest: %v", err)
	}

	// Add participants
	for i := 1; i <= 3; i++ {
		participant := &models.Participant{
			ContestID: contest.ID,
			UserID:    uint(i),
			Role:      "participant",
			Status:    "active",
			JoinedAt:  time.Now(),
		}
		if err := participantRepo.Create(participant); err != nil {
			t.Fatalf("Failed to create participant %d: %v", i, err)
		}
	}

	// Count participants using repository method
	count, err := participantRepo.CountByContest(contest.ID)
	if err != nil {
		t.Fatalf("Failed to count participants: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 participants, got %d", count)
	}

	// Remove one participant
	if err := participantRepo.DeleteByContestAndUser(contest.ID, 1); err != nil {
		t.Fatalf("Failed to delete participant: %v", err)
	}

	// Verify count is updated
	count, err = participantRepo.CountByContest(contest.ID)
	if err != nil {
		t.Fatalf("Failed to count participants after deletion: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 participants after deletion, got %d", count)
	}
}

func TestTimezoneHandling(t *testing.T) {
	// Test UTC timezone handling
	contest := &models.Contest{
		Title:           "Test Contest",
		Description:     "Test Description",
		SportType:       "football",
		Status:          "draft",
		StartDate:       time.Now().UTC().Add(1 * time.Hour),
		EndDate:         time.Now().UTC().Add(25 * time.Hour),
		MaxParticipants: 100,
		CreatorID:       1,
	}

	if err := contest.ValidateDates(); err != nil {
		t.Errorf("Expected valid dates with UTC handling, got error: %v", err)
	}

	// Test date too far in past
	contest.StartDate = time.Now().UTC().Add(-25 * time.Hour)
	if err := contest.ValidateDates(); err == nil {
		t.Error("Expected error for start date too far in past")
	}
}

func TestSportTypeExtensibility(t *testing.T) {
	contest := &models.Contest{
		Title:           "Test Contest",
		Description:     "Test Description",
		SportType:       "esports", // New sport type not in hardcoded list
		Status:          "draft",
		StartDate:       time.Now().Add(24 * time.Hour),
		EndDate:         time.Now().Add(48 * time.Hour),
		MaxParticipants: 100,
		CreatorID:       1,
	}

	if err := contest.ValidateSportType(); err != nil {
		t.Errorf("Expected new sport type to be valid, got error: %v", err)
	}

	// Test empty sport type still fails
	contest.SportType = ""
	if err := contest.ValidateSportType(); err == nil {
		t.Error("Expected error for empty sport type")
	}
}

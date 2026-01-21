package repository

import (
	"testing"

	"github.com/sports-prediction-contests/challenge-service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Challenge{}, &models.ChallengeParticipant{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCreateWithParticipantsTransaction(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChallengeRepository(db)

	t.Run("Successful atomic creation", func(t *testing.T) {
		challenge := &models.Challenge{
			ChallengerID: 1,
			OpponentID:   2,
			EventID:      1,
			Message:      "Test challenge",
			Status:       "pending",
		}

		participants := []*models.ChallengeParticipant{
			{
				UserID: 1,
				Role:   "challenger",
				Status: "active",
			},
			{
				UserID: 2,
				Role:   "opponent",
				Status: "active",
			},
		}

		err := repo.CreateWithParticipants(challenge, participants)
		if err != nil {
			t.Errorf("CreateWithParticipants failed: %v", err)
		}

		// Verify challenge was created
		if challenge.ID == 0 {
			t.Error("Challenge ID should be set after creation")
		}

		// Verify participants were created with correct challenge ID
		for _, participant := range participants {
			if participant.ChallengeID != challenge.ID {
				t.Errorf("Participant challenge ID (%d) should match challenge ID (%d)", participant.ChallengeID, challenge.ID)
			}
		}

		// Verify data in database
		var count int64
		db.Model(&models.Challenge{}).Count(&count)
		if count != 1 {
			t.Errorf("Expected 1 challenge in database, got %d", count)
		}

		db.Model(&models.ChallengeParticipant{}).Count(&count)
		if count != 2 {
			t.Errorf("Expected 2 participants in database, got %d", count)
		}
	})

	t.Run("Validation errors", func(t *testing.T) {
		// Test nil challenge
		err := repo.CreateWithParticipants(nil, []*models.ChallengeParticipant{})
		if err == nil {
			t.Error("Expected error for nil challenge")
		}

		// Test empty participants
		challenge := &models.Challenge{
			ChallengerID: 1,
			OpponentID:   2,
			EventID:      1,
		}
		err = repo.CreateWithParticipants(challenge, []*models.ChallengeParticipant{})
		if err == nil {
			t.Error("Expected error for empty participants")
		}
	})
}

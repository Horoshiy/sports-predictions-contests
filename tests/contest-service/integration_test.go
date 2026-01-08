//go:build integration

package contest_test

import (
	"context"
	"testing"
	"time"

	"github.com/sports-prediction-contests/contest-service/internal/models"
	"github.com/sports-prediction-contests/contest-service/internal/repository"
	"github.com/sports-prediction-contests/contest-service/internal/service"
	pb "github.com/sports-prediction-contests/shared/proto/contest"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&models.Contest{}, &models.Participant{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestContestServiceIntegration(t *testing.T) {
	db := setupTestDB(t)
	
	contestRepo := repository.NewContestRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	contestService := service.NewContestService(contestRepo, participantRepo)

	// Test contest creation
	req := &pb.CreateContestRequest{
		Title:           "Integration Test Contest",
		Description:     "Test Description",
		SportType:       "football",
		Rules:           `{"scoring": {"exact_score": 3}}`,
		StartDate:       timestamppb.New(time.Now().Add(24 * time.Hour)),
		EndDate:         timestamppb.New(time.Now().Add(48 * time.Hour)),
		MaxParticipants: 100,
	}

	// Note: This test would need proper JWT context setup
	// For now, we'll test the repository layer directly
	contest := &models.Contest{
		Title:           req.Title,
		Description:     req.Description,
		SportType:       req.SportType,
		Rules:           req.Rules,
		StartDate:       req.StartDate.AsTime(),
		EndDate:         req.EndDate.AsTime(),
		MaxParticipants: uint(req.MaxParticipants),
		CreatorID:       1,
		Status:          "draft",
	}

	// Test repository create
	if err := contestRepo.Create(contest); err != nil {
		t.Errorf("Failed to create contest: %v", err)
	}

	// Test repository get
	retrieved, err := contestRepo.GetByID(contest.ID)
	if err != nil {
		t.Errorf("Failed to retrieve contest: %v", err)
	}

	if retrieved.Title != contest.Title {
		t.Errorf("Expected title %s, got %s", contest.Title, retrieved.Title)
	}

	// Test participant creation
	participant := &models.Participant{
		ContestID: contest.ID,
		UserID:    1,
		Role:      "admin",
		Status:    "active",
		JoinedAt:  time.Now(),
	}

	if err := participantRepo.Create(participant); err != nil {
		t.Errorf("Failed to create participant: %v", err)
	}

	// Test participant listing
	participants, total, err := participantRepo.ListByContest(contest.ID, 10, 0)
	if err != nil {
		t.Errorf("Failed to list participants: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected 1 participant, got %d", total)
	}

	if len(participants) != 1 {
		t.Errorf("Expected 1 participant in list, got %d", len(participants))
	}
}

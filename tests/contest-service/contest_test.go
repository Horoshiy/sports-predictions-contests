package contest_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/contest-service/internal/models"
)

func TestContestValidation(t *testing.T) {
	contest := &models.Contest{
		Title:           "Test Contest",
		Description:     "Test Description",
		SportType:       "football",
		Status:          "draft",
		StartDate:       time.Now().Add(24 * time.Hour),
		EndDate:         time.Now().Add(48 * time.Hour),
		MaxParticipants: 100,
		CreatorID:       1,
	}

	// Test title validation
	if err := contest.ValidateTitle(); err != nil {
		t.Errorf("Expected valid title, got error: %v", err)
	}

	// Test empty title
	contest.Title = ""
	if err := contest.ValidateTitle(); err == nil {
		t.Error("Expected error for empty title")
	}

	// Test sport type validation
	contest.Title = "Test Contest"
	contest.SportType = "invalid_sport"
	if err := contest.ValidateSportType(); err == nil {
		t.Error("Expected error for invalid sport type")
	}

	// Test valid sport type
	contest.SportType = "football"
	if err := contest.ValidateSportType(); err != nil {
		t.Errorf("Expected valid sport type, got error: %v", err)
	}
}

func TestParticipantValidation(t *testing.T) {
	participant := &models.Participant{
		ContestID: 1,
		UserID:    1,
		Role:      "participant",
		Status:    "active",
		JoinedAt:  time.Now(),
	}

	// Test role validation
	if err := participant.ValidateRole(); err != nil {
		t.Errorf("Expected valid role, got error: %v", err)
	}

	// Test invalid role
	participant.Role = "invalid_role"
	if err := participant.ValidateRole(); err == nil {
		t.Error("Expected error for invalid role")
	}

	// Test status validation
	participant.Role = "participant"
	participant.Status = "invalid_status"
	if err := participant.ValidateStatus(); err == nil {
		t.Error("Expected error for invalid status")
	}
}

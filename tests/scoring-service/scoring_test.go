package scoring_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
)

func TestScoreValidation(t *testing.T) {
	score := &models.Score{
		UserID:       1,
		ContestID:    1,
		PredictionID: 1,
		Points:       10.5,
	}

	// Test valid score
	if err := score.ValidateUserID(); err != nil {
		t.Errorf("Expected valid user ID, got error: %v", err)
	}

	if err := score.ValidateContestID(); err != nil {
		t.Errorf("Expected valid contest ID, got error: %v", err)
	}

	if err := score.ValidatePoints(); err != nil {
		t.Errorf("Expected valid points, got error: %v", err)
	}

	// Test invalid user ID
	score.UserID = 0
	if err := score.ValidateUserID(); err == nil {
		t.Error("Expected error for invalid user ID")
	}

	// Test negative points
	score.UserID = 1
	score.Points = -5.0
	if err := score.ValidatePoints(); err == nil {
		t.Error("Expected error for negative points")
	}
}

func TestScoreBeforeUpdate(t *testing.T) {
	score := &models.Score{
		ID:           1,
		UserID:       1,
		ContestID:    1,
		PredictionID: 1,
		Points:       10.5,
		ScoredAt:     time.Now(),
	}

	// Test that points validation works
	score.Points = -5.0
	if err := score.BeforeUpdate(nil); err == nil {
		t.Error("Expected error for negative points in BeforeUpdate")
	}

	// Test valid points
	score.Points = 15.0
	if err := score.ValidatePoints(); err != nil {
		t.Errorf("Expected valid points, got error: %v", err)
	}
}

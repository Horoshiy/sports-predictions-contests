package scoring_test

import (
	"testing"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
)

func TestLeaderboardValidation(t *testing.T) {
	leaderboard := &models.Leaderboard{
		ContestID:   1,
		UserID:      1,
		TotalPoints: 25.5,
		Rank:        1,
	}

	// Test valid leaderboard entry
	if err := leaderboard.ValidateContestID(); err != nil {
		t.Errorf("Expected valid contest ID, got error: %v", err)
	}

	if err := leaderboard.ValidateUserID(); err != nil {
		t.Errorf("Expected valid user ID, got error: %v", err)
	}

	if err := leaderboard.ValidateTotalPoints(); err != nil {
		t.Errorf("Expected valid total points, got error: %v", err)
	}

	// Test invalid contest ID
	leaderboard.ContestID = 0
	if err := leaderboard.ValidateContestID(); err == nil {
		t.Error("Expected error for invalid contest ID")
	}

	// Test negative total points
	leaderboard.ContestID = 1
	leaderboard.TotalPoints = -10.0
	if err := leaderboard.ValidateTotalPoints(); err == nil {
		t.Error("Expected error for negative total points")
	}
}

func TestLeaderboardHelpers(t *testing.T) {
	leaderboard := &models.Leaderboard{
		TotalPoints: 15.5,
		Rank:        3,
	}

	// Test HasPoints
	if !leaderboard.HasPoints() {
		t.Error("Expected leaderboard with points to return true")
	}

	leaderboard.TotalPoints = 0
	if leaderboard.HasPoints() {
		t.Error("Expected leaderboard with zero points to return false")
	}

	// Test IsRanked
	leaderboard.Rank = 5
	if !leaderboard.IsRanked() {
		t.Error("Expected leaderboard with rank to return true")
	}

	leaderboard.Rank = 0
	if leaderboard.IsRanked() {
		t.Error("Expected leaderboard with zero rank to return false")
	}
}

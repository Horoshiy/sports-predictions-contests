package scoring_test

import (
	"testing"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
)

func TestStreakValidation(t *testing.T) {
	streak := &models.UserStreak{
		UserID:    1,
		ContestID: 1,
	}

	// Test valid streak
	if err := streak.ValidateUserID(); err != nil {
		t.Errorf("Expected valid user ID, got error: %v", err)
	}

	if err := streak.ValidateContestID(); err != nil {
		t.Errorf("Expected valid contest ID, got error: %v", err)
	}

	// Test invalid user ID
	streak.UserID = 0
	if err := streak.ValidateUserID(); err == nil {
		t.Error("Expected error for invalid user ID")
	}

	// Test invalid contest ID
	streak.UserID = 1
	streak.ContestID = 0
	if err := streak.ValidateContestID(); err == nil {
		t.Error("Expected error for invalid contest ID")
	}
}

func TestGetMultiplier(t *testing.T) {
	tests := []struct {
		streak     uint
		multiplier float64
	}{
		{0, 1.0},
		{1, 1.0},
		{2, 1.0},
		{3, 1.25},
		{4, 1.25},
		{5, 1.5},
		{6, 1.5},
		{7, 1.75},
		{8, 1.75},
		{9, 1.75},
		{10, 2.0},
		{15, 2.0},
		{100, 2.0},
	}

	for _, tt := range tests {
		streak := &models.UserStreak{CurrentStreak: tt.streak}
		got := streak.GetMultiplier()
		if got != tt.multiplier {
			t.Errorf("GetMultiplier() for streak %d = %v, want %v", tt.streak, got, tt.multiplier)
		}
	}
}

func TestIncrementStreak(t *testing.T) {
	streak := &models.UserStreak{
		UserID:        1,
		ContestID:     1,
		CurrentStreak: 0,
		MaxStreak:     0,
	}

	// First increment
	streak.IncrementStreak(100)
	if streak.CurrentStreak != 1 {
		t.Errorf("Expected current streak 1, got %d", streak.CurrentStreak)
	}
	if streak.MaxStreak != 1 {
		t.Errorf("Expected max streak 1, got %d", streak.MaxStreak)
	}
	if streak.LastPredictionID == nil || *streak.LastPredictionID != 100 {
		t.Error("Expected last prediction ID to be 100")
	}
	if streak.LastPredictionCorrect == nil || !*streak.LastPredictionCorrect {
		t.Error("Expected last prediction correct to be true")
	}

	// Second increment
	streak.IncrementStreak(101)
	if streak.CurrentStreak != 2 {
		t.Errorf("Expected current streak 2, got %d", streak.CurrentStreak)
	}
	if streak.MaxStreak != 2 {
		t.Errorf("Expected max streak 2, got %d", streak.MaxStreak)
	}

	// Increment to 5 for multiplier change
	streak.IncrementStreak(102)
	streak.IncrementStreak(103)
	streak.IncrementStreak(104)
	if streak.CurrentStreak != 5 {
		t.Errorf("Expected current streak 5, got %d", streak.CurrentStreak)
	}
	if streak.GetMultiplier() != 1.5 {
		t.Errorf("Expected multiplier 1.5 at streak 5, got %v", streak.GetMultiplier())
	}
}

func TestResetStreak(t *testing.T) {
	streak := &models.UserStreak{
		UserID:        1,
		ContestID:     1,
		CurrentStreak: 5,
		MaxStreak:     5,
	}

	// Reset streak
	streak.ResetStreak(200)
	if streak.CurrentStreak != 0 {
		t.Errorf("Expected current streak 0 after reset, got %d", streak.CurrentStreak)
	}
	if streak.MaxStreak != 5 {
		t.Errorf("Expected max streak to remain 5, got %d", streak.MaxStreak)
	}
	if streak.LastPredictionID == nil || *streak.LastPredictionID != 200 {
		t.Error("Expected last prediction ID to be 200")
	}
	if streak.LastPredictionCorrect == nil || *streak.LastPredictionCorrect {
		t.Error("Expected last prediction correct to be false")
	}

	// Verify multiplier is back to 1.0
	if streak.GetMultiplier() != 1.0 {
		t.Errorf("Expected multiplier 1.0 after reset, got %v", streak.GetMultiplier())
	}
}

func TestMaxStreakPreservation(t *testing.T) {
	streak := &models.UserStreak{
		UserID:        1,
		ContestID:     1,
		CurrentStreak: 0,
		MaxStreak:     0,
	}

	// Build up to streak of 10
	for i := 1; i <= 10; i++ {
		streak.IncrementStreak(uint(i))
	}
	if streak.MaxStreak != 10 {
		t.Errorf("Expected max streak 10, got %d", streak.MaxStreak)
	}

	// Reset
	streak.ResetStreak(100)
	if streak.MaxStreak != 10 {
		t.Errorf("Expected max streak to remain 10 after reset, got %d", streak.MaxStreak)
	}

	// Build new streak of 5
	for i := 1; i <= 5; i++ {
		streak.IncrementStreak(uint(100 + i))
	}
	if streak.CurrentStreak != 5 {
		t.Errorf("Expected current streak 5, got %d", streak.CurrentStreak)
	}
	if streak.MaxStreak != 10 {
		t.Errorf("Expected max streak to remain 10, got %d", streak.MaxStreak)
	}

	// Build beyond previous max
	for i := 6; i <= 12; i++ {
		streak.IncrementStreak(uint(100 + i))
	}
	if streak.MaxStreak != 12 {
		t.Errorf("Expected max streak to update to 12, got %d", streak.MaxStreak)
	}
}

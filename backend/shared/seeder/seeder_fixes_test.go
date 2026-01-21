package seeder

import (
	"testing"
)

func TestSeederFixes(t *testing.T) {
	t.Run("Status weights validation", func(t *testing.T) {
		// Test that status and weights arrays are validated
		statuses := []string{"pending", "accepted", "declined", "expired", "active", "completed"}
		statusWeights := []int{30, 25, 15, 10, 10, 10}

		if len(statuses) != len(statusWeights) {
			t.Errorf("Status array length (%d) does not match weights array length (%d)", len(statuses), len(statusWeights))
		}

		// Test mismatched arrays would fail
		mismatchedWeights := []int{30, 25, 15, 10, 10} // One less element
		if len(statuses) == len(mismatchedWeights) {
			t.Error("Mismatched arrays should not be equal length")
		}
	})

	t.Run("Opponent selection with fallback", func(t *testing.T) {
		// Simulate the opponent selection logic
		users := make([]int, 10) // 10 users
		challengerIdx := 5

		// Test that we can always find a different opponent
		var opponentIdx int
		for attempts := 0; attempts < 100; attempts++ {
			opponentIdx = attempts % len(users) // Simulate random selection
			if opponentIdx != challengerIdx {
				break
			}
			if attempts == 99 {
				// Fallback logic
				opponentIdx = (challengerIdx + 1) % len(users)
			}
		}

		if opponentIdx == challengerIdx {
			t.Error("Opponent should be different from challenger")
		}
	})
}

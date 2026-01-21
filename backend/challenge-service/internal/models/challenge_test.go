package models

import (
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestChallengeValidationFixes(t *testing.T) {
	t.Run("Status validation uses efficient lookup", func(t *testing.T) {
		challenge := &Challenge{Status: "pending"}
		if err := challenge.ValidateStatus(); err != nil {
			t.Errorf("Expected valid status, got error: %v", err)
		}

		challenge.Status = "invalid_status"
		if err := challenge.ValidateStatus(); err == nil {
			t.Error("Expected error for invalid status")
		}
	})

	t.Run("Time zone consistency with UTC", func(t *testing.T) {
		challenge := &Challenge{
			ChallengerID: 1,
			OpponentID:   2,
			EventID:      1,
			Status:       "pending",
		}

		// Test BeforeCreate sets UTC time
		err := challenge.BeforeCreate(&gorm.DB{})
		if err != nil {
			t.Errorf("BeforeCreate failed: %v", err)
		}

		// Check that ExpiresAt is set and in the future
		if challenge.ExpiresAt.IsZero() {
			t.Error("ExpiresAt should be set")
		}

		if !challenge.ExpiresAt.After(time.Now().UTC()) {
			t.Error("ExpiresAt should be in the future")
		}

		// Test Accept method uses UTC
		challenge.Accept()
		if challenge.AcceptedAt == nil {
			t.Error("AcceptedAt should be set")
		}

		// Test Complete method uses UTC
		challenge.Complete(10.0, 8.0)
		if challenge.CompletedAt == nil {
			t.Error("CompletedAt should be set")
		}
	})

	t.Run("CanAccept and IsExpired use UTC", func(t *testing.T) {
		challenge := &Challenge{
			Status:    "pending",
			ExpiresAt: time.Now().UTC().Add(1 * time.Hour),
		}

		if !challenge.CanAccept() {
			t.Error("Challenge should be acceptable")
		}

		if challenge.IsExpired() {
			t.Error("Challenge should not be expired")
		}

		// Test expired challenge
		challenge.ExpiresAt = time.Now().UTC().Add(-1 * time.Hour)
		if challenge.CanAccept() {
			t.Error("Expired challenge should not be acceptable")
		}

		if !challenge.IsExpired() {
			t.Error("Challenge should be expired")
		}
	})
}

func TestChallengeParticipantValidation(t *testing.T) {
	t.Run("Validation methods work correctly", func(t *testing.T) {
		participant := &ChallengeParticipant{
			ChallengeID: 1,
			UserID:      1,
			Role:        "challenger",
		}

		if err := participant.ValidateUserID(); err != nil {
			t.Errorf("Expected valid user ID, got error: %v", err)
		}

		if err := participant.ValidateChallengeID(); err != nil {
			t.Errorf("Expected valid challenge ID, got error: %v", err)
		}

		if err := participant.ValidateRole(); err != nil {
			t.Errorf("Expected valid role, got error: %v", err)
		}
	})
}

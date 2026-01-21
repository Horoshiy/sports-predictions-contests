package challenge_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/challenge-service/internal/models"
)

func TestChallengeValidation(t *testing.T) {
	challenge := &models.Challenge{
		ChallengerID: 1,
		OpponentID:   2,
		EventID:      1,
		Message:      "Test challenge message",
		Status:       "pending",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Test user IDs validation
	if err := challenge.ValidateUserIDs(); err != nil {
		t.Errorf("Expected valid user IDs, got error: %v", err)
	}

	// Test same user IDs
	challenge.OpponentID = 1
	if err := challenge.ValidateUserIDs(); err == nil {
		t.Error("Expected error for same challenger and opponent")
	}

	// Test empty challenger ID
	challenge.ChallengerID = 0
	challenge.OpponentID = 2
	if err := challenge.ValidateUserIDs(); err == nil {
		t.Error("Expected error for empty challenger ID")
	}

	// Test empty opponent ID
	challenge.ChallengerID = 1
	challenge.OpponentID = 0
	if err := challenge.ValidateUserIDs(); err == nil {
		t.Error("Expected error for empty opponent ID")
	}

	// Test event ID validation
	challenge.OpponentID = 2
	challenge.EventID = 0
	if err := challenge.ValidateEventID(); err == nil {
		t.Error("Expected error for empty event ID")
	}

	// Test valid event ID
	challenge.EventID = 1
	if err := challenge.ValidateEventID(); err != nil {
		t.Errorf("Expected valid event ID, got error: %v", err)
	}

	// Test status validation
	challenge.Status = "invalid_status"
	if err := challenge.ValidateStatus(); err == nil {
		t.Error("Expected error for invalid status")
	}

	// Test valid status
	challenge.Status = "pending"
	if err := challenge.ValidateStatus(); err != nil {
		t.Errorf("Expected valid status, got error: %v", err)
	}

	// Test message validation
	challenge.Message = string(make([]byte, 501)) // 501 characters
	if err := challenge.ValidateMessage(); err == nil {
		t.Error("Expected error for message too long")
	}

	// Test valid message
	challenge.Message = "Valid message"
	if err := challenge.ValidateMessage(); err != nil {
		t.Errorf("Expected valid message, got error: %v", err)
	}
}

func TestChallengeParticipantValidation(t *testing.T) {
	participant := &models.ChallengeParticipant{
		ChallengeID: 1,
		UserID:      1,
		Role:        "challenger",
		Status:      "active",
		JoinedAt:    time.Now(),
	}

	// Test user ID validation
	if err := participant.ValidateUserID(); err != nil {
		t.Errorf("Expected valid user ID, got error: %v", err)
	}

	// Test empty user ID
	participant.UserID = 0
	if err := participant.ValidateUserID(); err == nil {
		t.Error("Expected error for empty user ID")
	}

	// Test challenge ID validation
	participant.UserID = 1
	participant.ChallengeID = 0
	if err := participant.ValidateChallengeID(); err == nil {
		t.Error("Expected error for empty challenge ID")
	}

	// Test valid challenge ID
	participant.ChallengeID = 1
	if err := participant.ValidateChallengeID(); err != nil {
		t.Errorf("Expected valid challenge ID, got error: %v", err)
	}

	// Test role validation
	participant.Role = "invalid_role"
	if err := participant.ValidateRole(); err == nil {
		t.Error("Expected error for invalid role")
	}

	// Test valid role
	participant.Role = "challenger"
	if err := participant.ValidateRole(); err != nil {
		t.Errorf("Expected valid role, got error: %v", err)
	}

	participant.Role = "opponent"
	if err := participant.ValidateRole(); err != nil {
		t.Errorf("Expected valid role, got error: %v", err)
	}
}

func TestChallengeStateMethods(t *testing.T) {
	challenge := &models.Challenge{
		ChallengerID: 1,
		OpponentID:   2,
		EventID:      1,
		Status:       "pending",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Test CanAccept
	if !challenge.CanAccept() {
		t.Error("Expected challenge to be acceptable")
	}

	// Test expired challenge
	challenge.ExpiresAt = time.Now().Add(-1 * time.Hour)
	if challenge.CanAccept() {
		t.Error("Expected expired challenge to not be acceptable")
	}

	// Test IsExpired
	if !challenge.IsExpired() {
		t.Error("Expected challenge to be expired")
	}

	// Test Accept method
	challenge.ExpiresAt = time.Now().Add(24 * time.Hour)
	challenge.Accept()
	if challenge.Status != "accepted" {
		t.Errorf("Expected status to be 'accepted', got %s", challenge.Status)
	}
	if challenge.AcceptedAt == nil {
		t.Error("Expected AcceptedAt to be set")
	}

	// Test IsActive
	if !challenge.IsActive() {
		t.Error("Expected accepted challenge to be active")
	}

	// Test Complete method
	challenge.Complete(10.5, 8.0)
	if challenge.Status != "completed" {
		t.Errorf("Expected status to be 'completed', got %s", challenge.Status)
	}
	if challenge.ChallengerScore != 10.5 {
		t.Errorf("Expected challenger score to be 10.5, got %f", challenge.ChallengerScore)
	}
	if challenge.OpponentScore != 8.0 {
		t.Errorf("Expected opponent score to be 8.0, got %f", challenge.OpponentScore)
	}
	if challenge.WinnerID == nil || *challenge.WinnerID != 1 {
		t.Error("Expected challenger to be winner")
	}
	if challenge.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}

	// Test IsCompleted
	if !challenge.IsCompleted() {
		t.Error("Expected challenge to be completed")
	}

	// Test tie scenario
	challenge2 := &models.Challenge{
		ChallengerID: 1,
		OpponentID:   2,
		EventID:      1,
		Status:       "accepted",
	}
	challenge2.Complete(10.0, 10.0)
	if challenge2.WinnerID != nil {
		t.Error("Expected no winner for tie")
	}
}

func TestChallengeBusinessLogic(t *testing.T) {
	// Test challenge creation with default values
	challenge := &models.Challenge{
		ChallengerID: 1,
		OpponentID:   2,
		EventID:      1,
	}

	// Simulate BeforeCreate hook
	if challenge.Status == "" {
		challenge.Status = "pending"
	}
	if challenge.ExpiresAt.IsZero() {
		challenge.ExpiresAt = time.Now().Add(24 * time.Hour)
	}

	if challenge.Status != "pending" {
		t.Errorf("Expected default status to be 'pending', got %s", challenge.Status)
	}

	if challenge.ExpiresAt.Before(time.Now()) {
		t.Error("Expected expiration to be in the future")
	}

	// Test challenge workflow
	if !challenge.CanAccept() {
		t.Error("New challenge should be acceptable")
	}

	challenge.Accept()
	if challenge.Status != "accepted" {
		t.Error("Challenge should be accepted")
	}

	if !challenge.IsActive() {
		t.Error("Accepted challenge should be active")
	}

	challenge.Complete(15.0, 12.0)
	if !challenge.IsCompleted() {
		t.Error("Challenge should be completed")
	}

	if challenge.WinnerID == nil || *challenge.WinnerID != 1 {
		t.Error("Challenger should be winner")
	}
}

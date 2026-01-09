package prediction_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
)

func TestPredictionValidation(t *testing.T) {
	prediction := &models.Prediction{
		ContestID:      1,
		UserID:         1,
		EventID:        1,
		PredictionData: `{"outcome": "home_win"}`,
		Status:         "pending",
		SubmittedAt:    time.Now().UTC(),
	}

	// Test contest ID validation
	if err := prediction.ValidateContestID(); err != nil {
		t.Errorf("Expected valid contest ID, got error: %v", err)
	}

	// Test empty contest ID
	prediction.ContestID = 0
	if err := prediction.ValidateContestID(); err == nil {
		t.Error("Expected error for empty contest ID")
	}

	// Test user ID validation
	prediction.ContestID = 1
	if err := prediction.ValidateUserID(); err != nil {
		t.Errorf("Expected valid user ID, got error: %v", err)
	}

	// Test empty user ID
	prediction.UserID = 0
	if err := prediction.ValidateUserID(); err == nil {
		t.Error("Expected error for empty user ID")
	}

	// Test event ID validation
	prediction.UserID = 1
	if err := prediction.ValidateEventID(); err != nil {
		t.Errorf("Expected valid event ID, got error: %v", err)
	}

	// Test empty event ID
	prediction.EventID = 0
	if err := prediction.ValidateEventID(); err == nil {
		t.Error("Expected error for empty event ID")
	}

	// Test prediction data validation
	prediction.EventID = 1
	if err := prediction.ValidatePredictionData(); err != nil {
		t.Errorf("Expected valid prediction data, got error: %v", err)
	}

	// Test empty prediction data
	prediction.PredictionData = ""
	if err := prediction.ValidatePredictionData(); err == nil {
		t.Error("Expected error for empty prediction data")
	}

	// Test status validation
	prediction.PredictionData = `{"outcome": "home_win"}`
	if err := prediction.ValidateStatus(); err != nil {
		t.Errorf("Expected valid status, got error: %v", err)
	}

	// Test invalid status
	prediction.Status = "invalid_status"
	if err := prediction.ValidateStatus(); err == nil {
		t.Error("Expected error for invalid status")
	}
}

func TestEventValidation(t *testing.T) {
	event := &models.Event{
		Title:     "Test Match",
		SportType: "football",
		HomeTeam:  "Team A",
		AwayTeam:  "Team B",
		EventDate: time.Now().Add(24 * time.Hour),
		Status:    "scheduled",
	}

	// Test title validation
	if err := event.ValidateTitle(); err != nil {
		t.Errorf("Expected valid title, got error: %v", err)
	}

	// Test empty title
	event.Title = ""
	if err := event.ValidateTitle(); err == nil {
		t.Error("Expected error for empty title")
	}

	// Test sport type validation
	event.Title = "Test Match"
	if err := event.ValidateSportType(); err != nil {
		t.Errorf("Expected valid sport type, got error: %v", err)
	}

	// Test empty sport type
	event.SportType = ""
	if err := event.ValidateSportType(); err == nil {
		t.Error("Expected error for empty sport type")
	}

	// Test team validation
	event.SportType = "football"
	if err := event.ValidateTeams(); err != nil {
		t.Errorf("Expected valid teams, got error: %v", err)
	}

	// Test empty home team
	event.HomeTeam = ""
	if err := event.ValidateTeams(); err == nil {
		t.Error("Expected error for empty home team")
	}

	// Test empty away team
	event.HomeTeam = "Team A"
	event.AwayTeam = ""
	if err := event.ValidateTeams(); err == nil {
		t.Error("Expected error for empty away team")
	}

	// Test event date validation
	event.AwayTeam = "Team B"
	if err := event.ValidateEventDate(); err != nil {
		t.Errorf("Expected valid event date, got error: %v", err)
	}

	// Test past event date
	event.EventDate = time.Now().Add(-48 * time.Hour)
	if err := event.ValidateEventDate(); err == nil {
		t.Error("Expected error for past event date")
	}

	// Test status validation
	event.EventDate = time.Now().Add(24 * time.Hour)
	if err := event.ValidateStatus(); err != nil {
		t.Errorf("Expected valid status, got error: %v", err)
	}

	// Test invalid status
	event.Status = "invalid_status"
	if err := event.ValidateStatus(); err == nil {
		t.Error("Expected error for invalid status")
	}
}

func TestPredictionCanUpdate(t *testing.T) {
	prediction := &models.Prediction{
		Status: "pending",
	}

	if !prediction.CanUpdate() {
		t.Error("Expected pending prediction to be updatable")
	}

	prediction.Status = "scored"
	if prediction.CanUpdate() {
		t.Error("Expected scored prediction to not be updatable")
	}
}

func TestEventCanAcceptPredictions(t *testing.T) {
	event := &models.Event{
		Status:    "scheduled",
		EventDate: time.Now().Add(24 * time.Hour),
	}

	if !event.CanAcceptPredictions() {
		t.Error("Expected scheduled future event to accept predictions")
	}

	event.Status = "completed"
	if event.CanAcceptPredictions() {
		t.Error("Expected completed event to not accept predictions")
	}

	event.Status = "scheduled"
	event.EventDate = time.Now().Add(-1 * time.Hour)
	if event.CanAcceptPredictions() {
		t.Error("Expected past event to not accept predictions")
	}
}

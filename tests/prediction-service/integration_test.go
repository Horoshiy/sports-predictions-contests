// +build integration

package prediction_test

import (
	"context"
	"testing"
	"time"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"github.com/sports-prediction-contests/prediction-service/internal/repository"
	"github.com/sports-prediction-contests/shared/database"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database
	db, err := database.NewConnectionFromEnv()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate test schema
	if err := db.AutoMigrate(&models.Prediction{}, &models.Event{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	// Clean up test data
	db.Exec("DELETE FROM predictions")
	db.Exec("DELETE FROM events")
}

func TestPredictionRepositoryIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := repository.NewPredictionRepository(db)

	// Create test event first
	event := &models.Event{
		Title:     "Test Match",
		SportType: "football",
		HomeTeam:  "Team A",
		AwayTeam:  "Team B",
		EventDate: time.Now().Add(24 * time.Hour),
		Status:    "scheduled",
	}
	if err := db.Create(event).Error; err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}

	// Test Create
	prediction := &models.Prediction{
		ContestID:      1,
		UserID:         1,
		EventID:        event.ID,
		PredictionData: `{"outcome": "home_win"}`,
		Status:         "pending",
		SubmittedAt:    time.Now().UTC(),
	}

	err := repo.Create(prediction)
	if err != nil {
		t.Fatalf("Failed to create prediction: %v", err)
	}

	if prediction.ID == 0 {
		t.Error("Expected prediction ID to be set after creation")
	}

	// Test GetByID
	retrieved, err := repo.GetByID(prediction.ID)
	if err != nil {
		t.Fatalf("Failed to get prediction by ID: %v", err)
	}

	if retrieved.ContestID != prediction.ContestID {
		t.Errorf("Expected contest ID %d, got %d", prediction.ContestID, retrieved.ContestID)
	}

	// Test GetByUserAndContest
	predictions, err := repo.GetByUserAndContest(1, 1)
	if err != nil {
		t.Fatalf("Failed to get predictions by user and contest: %v", err)
	}

	if len(predictions) != 1 {
		t.Errorf("Expected 1 prediction, got %d", len(predictions))
	}

	// Test GetByUserContestAndEvent
	specific, err := repo.GetByUserContestAndEvent(1, 1, event.ID)
	if err != nil {
		t.Fatalf("Failed to get specific prediction: %v", err)
	}

	if specific == nil {
		t.Error("Expected to find specific prediction")
	}

	// Test Update
	prediction.PredictionData = `{"outcome": "away_win"}`
	err = repo.Update(prediction)
	if err != nil {
		t.Fatalf("Failed to update prediction: %v", err)
	}

	updated, err := repo.GetByID(prediction.ID)
	if err != nil {
		t.Fatalf("Failed to get updated prediction: %v", err)
	}

	if updated.PredictionData != `{"outcome": "away_win"}` {
		t.Errorf("Expected updated prediction data, got %s", updated.PredictionData)
	}

	// Test List
	allPredictions, total, err := repo.List(10, 0, 1, 0)
	if err != nil {
		t.Fatalf("Failed to list predictions: %v", err)
	}

	if len(allPredictions) != 1 {
		t.Errorf("Expected 1 prediction in list, got %d", len(allPredictions))
	}

	if total != 1 {
		t.Errorf("Expected total count 1, got %d", total)
	}

	// Test CountByContest
	count, err := repo.CountByContest(1)
	if err != nil {
		t.Fatalf("Failed to count predictions by contest: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Test Delete
	err = repo.Delete(prediction.ID)
	if err != nil {
		t.Fatalf("Failed to delete prediction: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(prediction.ID)
	if err == nil {
		t.Error("Expected error when getting deleted prediction")
	}
}

func TestEventRepositoryIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := repository.NewEventRepository(db)

	// Test Create
	event := &models.Event{
		Title:     "Test Match",
		SportType: "football",
		HomeTeam:  "Team A",
		AwayTeam:  "Team B",
		EventDate: time.Now().Add(24 * time.Hour),
		Status:    "scheduled",
	}

	err := repo.Create(event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	if event.ID == 0 {
		t.Error("Expected event ID to be set after creation")
	}

	// Test GetByID
	retrieved, err := repo.GetByID(event.ID)
	if err != nil {
		t.Fatalf("Failed to get event by ID: %v", err)
	}

	if retrieved.Title != event.Title {
		t.Errorf("Expected title %s, got %s", event.Title, retrieved.Title)
	}

	// Test GetBySportType
	events, err := repo.GetBySportType("football")
	if err != nil {
		t.Fatalf("Failed to get events by sport type: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	// Test GetUpcoming
	upcoming, err := repo.GetUpcoming(10)
	if err != nil {
		t.Fatalf("Failed to get upcoming events: %v", err)
	}

	if len(upcoming) != 1 {
		t.Errorf("Expected 1 upcoming event, got %d", len(upcoming))
	}

	// Test List
	allEvents, total, err := repo.List(10, 0, "football", "scheduled")
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	if len(allEvents) != 1 {
		t.Errorf("Expected 1 event in list, got %d", len(allEvents))
	}

	if total != 1 {
		t.Errorf("Expected total count 1, got %d", total)
	}

	// Test Update
	event.Status = "completed"
	event.ResultData = `{"home_score": 2, "away_score": 1}`
	err = repo.Update(event)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	updated, err := repo.GetByID(event.ID)
	if err != nil {
		t.Fatalf("Failed to get updated event: %v", err)
	}

	if updated.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", updated.Status)
	}

	// Test Delete
	err = repo.Delete(event.ID)
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(event.ID)
	if err == nil {
		t.Error("Expected error when getting deleted event")
	}
}

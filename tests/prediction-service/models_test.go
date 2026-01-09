package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Prediction{}, &Event{})
	return db
}

func TestPredictionUniqueConstraint(t *testing.T) {
	db := setupTestDB()

	// Create test event
	event := &Event{
		Title:     "Test Event",
		SportType: "football",
		HomeTeam:  "Team A",
		AwayTeam:  "Team B",
		EventDate: time.Now().Add(24 * time.Hour),
		Status:    "scheduled",
	}
	db.Create(event)

	// Create first prediction
	prediction1 := &Prediction{
		ContestID:      1,
		UserID:         1,
		EventID:        event.ID,
		PredictionData: `{"score": "2-1"}`,
		Status:         "pending",
		SubmittedAt:    time.Now(),
	}
	err := db.Create(prediction1).Error
	assert.NoError(t, err)

	// Try to create duplicate prediction
	prediction2 := &Prediction{
		ContestID:      1,
		UserID:         1,
		EventID:        event.ID,
		PredictionData: `{"score": "1-0"}`,
		Status:         "pending",
		SubmittedAt:    time.Now(),
	}
	err = db.Create(prediction2).Error
	assert.Error(t, err, "Should fail due to unique constraint")
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestPredictionBeforeUpdateDoesNotChangeSubmittedAt(t *testing.T) {
	db := setupTestDB()

	// Create test event
	event := &Event{
		Title:     "Test Event",
		SportType: "football",
		HomeTeam:  "Team A",
		AwayTeam:  "Team B",
		EventDate: time.Now().Add(24 * time.Hour),
		Status:    "scheduled",
	}
	db.Create(event)

	// Create prediction
	originalTime := time.Now().Add(-1 * time.Hour)
	prediction := &Prediction{
		ContestID:      1,
		UserID:         1,
		EventID:        event.ID,
		PredictionData: `{"score": "2-1"}`,
		Status:         "pending",
		SubmittedAt:    originalTime,
	}
	db.Create(prediction)

	// Update prediction
	prediction.PredictionData = `{"score": "3-1"}`
	db.Save(prediction)

	// Verify SubmittedAt didn't change
	var updatedPrediction Prediction
	db.First(&updatedPrediction, prediction.ID)
	assert.Equal(t, originalTime.Unix(), updatedPrediction.SubmittedAt.Unix())
}

func TestEventDateValidation(t *testing.T) {
	tests := []struct {
		name        string
		eventDate   time.Time
		shouldError bool
	}{
		{
			name:        "Future event should be valid",
			eventDate:   time.Now().Add(24 * time.Hour),
			shouldError: false,
		},
		{
			name:        "Event 30 minutes ago should be valid",
			eventDate:   time.Now().Add(-30 * time.Minute),
			shouldError: false,
		},
		{
			name:        "Event 2 hours ago should be invalid",
			eventDate:   time.Now().Add(-2 * time.Hour),
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{
				Title:     "Test Event",
				SportType: "football",
				HomeTeam:  "Team A",
				AwayTeam:  "Team B",
				EventDate: tt.eventDate,
				Status:    "scheduled",
			}

			err := event.ValidateEventDate()
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCanAcceptPredictions(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		date     time.Time
		expected bool
	}{
		{
			name:     "Scheduled future event should accept predictions",
			status:   "scheduled",
			date:     time.Now().Add(1 * time.Hour),
			expected: true,
		},
		{
			name:     "Scheduled past event should not accept predictions",
			status:   "scheduled",
			date:     time.Now().Add(-1 * time.Hour),
			expected: false,
		},
		{
			name:     "Live event should not accept predictions",
			status:   "live",
			date:     time.Now().Add(1 * time.Hour),
			expected: false,
		},
		{
			name:     "Completed event should not accept predictions",
			status:   "completed",
			date:     time.Now().Add(1 * time.Hour),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{
				Status:    tt.status,
				EventDate: tt.date,
			}
			assert.Equal(t, tt.expected, event.CanAcceptPredictions())
		})
	}
}

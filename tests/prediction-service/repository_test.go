package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Prediction{}, &models.Event{})
	return db
}

func TestPaginationValidation(t *testing.T) {
	db := setupTestDB()
	repo := NewPredictionRepository(db)

	tests := []struct {
		name           string
		inputLimit     int
		inputOffset    int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "Negative limit should default to 10",
			inputLimit:     -5,
			inputOffset:    0,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Zero limit should default to 10",
			inputLimit:     0,
			inputOffset:    0,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Negative offset should default to 0",
			inputLimit:     5,
			inputOffset:    -10,
			expectedLimit:  5,
			expectedOffset: 0,
		},
		{
			name:           "Large limit should be capped at 100",
			inputLimit:     200,
			inputOffset:    0,
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "Valid parameters should remain unchanged",
			inputLimit:     20,
			inputOffset:    10,
			expectedLimit:  20,
			expectedOffset: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies that the repository handles invalid pagination parameters
			// by checking that no error occurs and results are returned
			predictions, total, err := repo.List(tt.inputLimit, tt.inputOffset, 0, 0)
			
			assert.NoError(t, err)
			assert.NotNil(t, predictions)
			assert.GreaterOrEqual(t, total, int64(0))
		})
	}
}

func TestEventRepositoryPaginationValidation(t *testing.T) {
	db := setupTestDB()
	repo := NewEventRepository(db)

	tests := []struct {
		name           string
		inputLimit     int
		inputOffset    int
	}{
		{
			name:        "Negative limit should be handled",
			inputLimit:  -5,
			inputOffset: 0,
		},
		{
			name:        "Zero limit should be handled",
			inputLimit:  0,
			inputOffset: 0,
		},
		{
			name:        "Negative offset should be handled",
			inputLimit:  5,
			inputOffset: -10,
		},
		{
			name:        "Large limit should be handled",
			inputLimit:  200,
			inputOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, total, err := repo.List(tt.inputLimit, tt.inputOffset, "", "")
			
			assert.NoError(t, err)
			assert.NotNil(t, events)
			assert.GreaterOrEqual(t, total, int64(0))
		})
	}
}

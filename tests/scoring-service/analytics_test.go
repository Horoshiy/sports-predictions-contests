package scoring_service_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
)

func TestTimeRangeToDate(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name      string
		timeRange string
		expected  func() time.Time
	}{
		{
			name:      "7 days",
			timeRange: "7d",
			expected:  func() time.Time { return now.AddDate(0, 0, -7) },
		},
		{
			name:      "30 days",
			timeRange: "30d",
			expected:  func() time.Time { return now.AddDate(0, 0, -30) },
		},
		{
			name:      "90 days",
			timeRange: "90d",
			expected:  func() time.Time { return now.AddDate(0, 0, -90) },
		},
		{
			name:      "all time returns zero",
			timeRange: "all",
			expected:  func() time.Time { return time.Time{} },
		},
		{
			name:      "unknown returns zero",
			timeRange: "unknown",
			expected:  func() time.Time { return time.Time{} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.TimeRangeToDate(tt.timeRange)
			expected := tt.expected()

			if tt.timeRange == "all" || tt.timeRange == "unknown" {
				if !result.IsZero() {
					t.Errorf("TimeRangeToDate(%s) = %v, want zero time", tt.timeRange, result)
				}
			} else {
				// Allow 1 second tolerance for test execution time
				diff := result.Sub(expected)
				if diff < -time.Second || diff > time.Second {
					t.Errorf("TimeRangeToDate(%s) = %v, want approximately %v", tt.timeRange, result, expected)
				}
			}
		})
	}
}

func TestUserAnalyticsStruct(t *testing.T) {
	analytics := &models.UserAnalytics{
		UserID:             1,
		TotalPredictions:   100,
		CorrectPredictions: 75,
		OverallAccuracy:    75.0,
		TotalPoints:        150.5,
		TimeRange:          "30d",
	}

	if analytics.UserID != 1 {
		t.Errorf("UserID = %d, want 1", analytics.UserID)
	}
	if analytics.TotalPredictions != 100 {
		t.Errorf("TotalPredictions = %d, want 100", analytics.TotalPredictions)
	}
	if analytics.OverallAccuracy != 75.0 {
		t.Errorf("OverallAccuracy = %f, want 75.0", analytics.OverallAccuracy)
	}
}

func TestSportAccuracyStruct(t *testing.T) {
	accuracy := models.SportAccuracy{
		SportType:          "football",
		TotalPredictions:   50,
		CorrectPredictions: 40,
		AccuracyPercentage: 80.0,
		TotalPoints:        100.0,
	}

	if accuracy.SportType != "football" {
		t.Errorf("SportType = %s, want football", accuracy.SportType)
	}
	if accuracy.AccuracyPercentage != 80.0 {
		t.Errorf("AccuracyPercentage = %f, want 80.0", accuracy.AccuracyPercentage)
	}
}

func TestPlatformStatsStruct(t *testing.T) {
	stats := &models.PlatformStats{
		AverageAccuracy:            65.5,
		AveragePointsPerPrediction: 1.5,
		TotalUsers:                 1000,
		TotalPredictions:           50000,
	}

	if stats.TotalUsers != 1000 {
		t.Errorf("TotalUsers = %d, want 1000", stats.TotalUsers)
	}
	if stats.AverageAccuracy != 65.5 {
		t.Errorf("AverageAccuracy = %f, want 65.5", stats.AverageAccuracy)
	}
}

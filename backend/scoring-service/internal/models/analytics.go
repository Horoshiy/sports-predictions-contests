package models

import "time"

// SportAccuracy represents accuracy statistics for a specific sport
type SportAccuracy struct {
	SportType          string  `json:"sport_type"`
	TotalPredictions   int     `json:"total_predictions"`
	CorrectPredictions int     `json:"correct_predictions"`
	AccuracyPercentage float64 `json:"accuracy_percentage"`
	TotalPoints        float64 `json:"total_points"`
}

// LeagueAccuracy represents accuracy statistics for a specific league
type LeagueAccuracy struct {
	LeagueID           uint    `json:"league_id"`
	LeagueName         string  `json:"league_name"`
	SportType          string  `json:"sport_type"`
	TotalPredictions   int     `json:"total_predictions"`
	CorrectPredictions int     `json:"correct_predictions"`
	AccuracyPercentage float64 `json:"accuracy_percentage"`
}

// PredictionTypeAccuracy represents accuracy by prediction type
type PredictionTypeAccuracy struct {
	PredictionType     string  `json:"prediction_type"`
	TotalPredictions   int     `json:"total_predictions"`
	CorrectPredictions int     `json:"correct_predictions"`
	AccuracyPercentage float64 `json:"accuracy_percentage"`
	AveragePoints      float64 `json:"average_points"`
}

// AccuracyTrend represents accuracy over a time period
type AccuracyTrend struct {
	Period             string  `json:"period"`
	TotalPredictions   int     `json:"total_predictions"`
	CorrectPredictions int     `json:"correct_predictions"`
	AccuracyPercentage float64 `json:"accuracy_percentage"`
	TotalPoints        float64 `json:"total_points"`
}

// PlatformStats represents platform-wide statistics for comparison
type PlatformStats struct {
	AverageAccuracy            float64 `json:"average_accuracy"`
	AveragePointsPerPrediction float64 `json:"average_points_per_prediction"`
	TotalUsers                 int     `json:"total_users"`
	TotalPredictions           int     `json:"total_predictions"`
}

// UserAnalytics aggregates all analytics for a user
type UserAnalytics struct {
	UserID             uint                     `json:"user_id"`
	TotalPredictions   int                      `json:"total_predictions"`
	CorrectPredictions int                      `json:"correct_predictions"`
	OverallAccuracy    float64                  `json:"overall_accuracy"`
	TotalPoints        float64                  `json:"total_points"`
	BySport            []SportAccuracy          `json:"by_sport"`
	ByLeague           []LeagueAccuracy         `json:"by_league"`
	ByType             []PredictionTypeAccuracy `json:"by_type"`
	Trends             []AccuracyTrend          `json:"trends"`
	PlatformComparison *PlatformStats           `json:"platform_comparison"`
	TimeRange          string                   `json:"time_range"`
}

// TimeRangeToDate converts time range string to start date
func TimeRangeToDate(timeRange string) time.Time {
	now := time.Now().UTC()
	switch timeRange {
	case "7d":
		return now.AddDate(0, 0, -7)
	case "30d":
		return now.AddDate(0, 0, -30)
	case "90d":
		return now.AddDate(0, 0, -90)
	default:
		return time.Time{} // Zero time for "all"
	}
}

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
	"gorm.io/gorm"
)

// AnalyticsRepositoryInterface defines analytics query methods
type AnalyticsRepositoryInterface interface {
	GetUserOverallStats(ctx context.Context, userID uint, since time.Time) (*models.UserAnalytics, error)
	GetAccuracyBySport(ctx context.Context, userID uint, since time.Time) ([]models.SportAccuracy, error)
	GetAccuracyByLeague(ctx context.Context, userID uint, since time.Time) ([]models.LeagueAccuracy, error)
	GetAccuracyByType(ctx context.Context, userID uint, since time.Time) ([]models.PredictionTypeAccuracy, error)
	GetAccuracyTrends(ctx context.Context, userID uint, since time.Time, groupBy string) ([]models.AccuracyTrend, error)
	GetPlatformStats(ctx context.Context, since time.Time) (*models.PlatformStats, error)
}

// AnalyticsRepository implements AnalyticsRepositoryInterface
type AnalyticsRepository struct {
	db *gorm.DB
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepositoryInterface {
	return &AnalyticsRepository{db: db}
}

// GetUserOverallStats retrieves overall user statistics
func (r *AnalyticsRepository) GetUserOverallStats(ctx context.Context, userID uint, since time.Time) (*models.UserAnalytics, error) {
	var result struct {
		TotalPredictions   int
		CorrectPredictions int
		TotalPoints        float64
	}

	query := r.db.WithContext(ctx).Table("scores").
		Select("COUNT(*) as total_predictions, SUM(CASE WHEN points > 0 THEN 1 ELSE 0 END) as correct_predictions, COALESCE(SUM(points), 0) as total_points").
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if !since.IsZero() {
		query = query.Where("scored_at >= ?", since)
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	accuracy := 0.0
	if result.TotalPredictions > 0 {
		accuracy = float64(result.CorrectPredictions) / float64(result.TotalPredictions) * 100
	}

	return &models.UserAnalytics{
		UserID:             userID,
		TotalPredictions:   result.TotalPredictions,
		CorrectPredictions: result.CorrectPredictions,
		OverallAccuracy:    accuracy,
		TotalPoints:        result.TotalPoints,
	}, nil
}

// GetAccuracyBySport retrieves accuracy grouped by sport type
func (r *AnalyticsRepository) GetAccuracyBySport(ctx context.Context, userID uint, since time.Time) ([]models.SportAccuracy, error) {
	var results []struct {
		SportType          string
		TotalPredictions   int
		CorrectPredictions int
		TotalPoints        float64
	}

	query := r.db.WithContext(ctx).Table("scores s").
		Select("COALESCE(e.sport_type, 'unknown') as sport_type, COUNT(*) as total_predictions, SUM(CASE WHEN s.points > 0 THEN 1 ELSE 0 END) as correct_predictions, COALESCE(SUM(s.points), 0) as total_points").
		Joins("LEFT JOIN predictions p ON s.prediction_id = p.id").
		Joins("LEFT JOIN events e ON p.event_id = e.id").
		Where("s.user_id = ? AND s.deleted_at IS NULL", userID)

	if !since.IsZero() {
		query = query.Where("s.scored_at >= ?", since)
	}

	query = query.Group("e.sport_type")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	accuracies := make([]models.SportAccuracy, len(results))
	for i, r := range results {
		accuracy := 0.0
		if r.TotalPredictions > 0 {
			accuracy = float64(r.CorrectPredictions) / float64(r.TotalPredictions) * 100
		}
		accuracies[i] = models.SportAccuracy{
			SportType:          r.SportType,
			TotalPredictions:   r.TotalPredictions,
			CorrectPredictions: r.CorrectPredictions,
			AccuracyPercentage: accuracy,
			TotalPoints:        r.TotalPoints,
		}
	}

	return accuracies, nil
}

// GetAccuracyByLeague retrieves accuracy grouped by league
// Note: This requires events to be linked to matches/leagues. Currently returns empty
// if the event->match relationship doesn't exist in the schema.
func (r *AnalyticsRepository) GetAccuracyByLeague(ctx context.Context, userID uint, since time.Time) ([]models.LeagueAccuracy, error) {
	var results []struct {
		LeagueID           uint
		LeagueName         string
		SportType          string
		TotalPredictions   int
		CorrectPredictions int
	}

	// Events don't have direct league relationship in current schema
	// This query attempts to join through matches table if events have match_id
	query := r.db.WithContext(ctx).Table("scores s").
		Select("l.id as league_id, l.name as league_name, sp.name as sport_type, COUNT(*) as total_predictions, SUM(CASE WHEN s.points > 0 THEN 1 ELSE 0 END) as correct_predictions").
		Joins("LEFT JOIN predictions p ON s.prediction_id = p.id").
		Joins("LEFT JOIN events e ON p.event_id = e.id").
		Joins("LEFT JOIN matches m ON m.id = e.id").
		Joins("LEFT JOIN leagues l ON m.league_id = l.id").
		Joins("LEFT JOIN sports sp ON l.sport_id = sp.id").
		Where("s.user_id = ? AND l.id IS NOT NULL AND s.deleted_at IS NULL", userID)

	if !since.IsZero() {
		query = query.Where("s.scored_at >= ?", since)
	}

	query = query.Group("l.id, l.name, sp.name")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	accuracies := make([]models.LeagueAccuracy, len(results))
	for i, r := range results {
		accuracy := 0.0
		if r.TotalPredictions > 0 {
			accuracy = float64(r.CorrectPredictions) / float64(r.TotalPredictions) * 100
		}
		accuracies[i] = models.LeagueAccuracy{
			LeagueID:           r.LeagueID,
			LeagueName:         r.LeagueName,
			SportType:          r.SportType,
			TotalPredictions:   r.TotalPredictions,
			CorrectPredictions: r.CorrectPredictions,
			AccuracyPercentage: accuracy,
		}
	}

	return accuracies, nil
}

// GetAccuracyByType retrieves accuracy grouped by prediction type
func (r *AnalyticsRepository) GetAccuracyByType(ctx context.Context, userID uint, since time.Time) ([]models.PredictionTypeAccuracy, error) {
	var results []struct {
		PredictionType     string
		TotalPredictions   int
		CorrectPredictions int
		TotalPoints        float64
	}

	query := r.db.WithContext(ctx).Table("scores s").
		Select("COALESCE(p.prediction_data::json->>'type', 'unknown') as prediction_type, COUNT(*) as total_predictions, SUM(CASE WHEN s.points > 0 THEN 1 ELSE 0 END) as correct_predictions, COALESCE(SUM(s.points), 0) as total_points").
		Joins("LEFT JOIN predictions p ON s.prediction_id = p.id").
		Where("s.user_id = ? AND s.deleted_at IS NULL", userID)

	if !since.IsZero() {
		query = query.Where("s.scored_at >= ?", since)
	}

	query = query.Group("p.prediction_data::json->>'type'")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	accuracies := make([]models.PredictionTypeAccuracy, len(results))
	for i, r := range results {
		accuracy := 0.0
		avgPoints := 0.0
		if r.TotalPredictions > 0 {
			accuracy = float64(r.CorrectPredictions) / float64(r.TotalPredictions) * 100
			avgPoints = r.TotalPoints / float64(r.TotalPredictions)
		}
		accuracies[i] = models.PredictionTypeAccuracy{
			PredictionType:     r.PredictionType,
			TotalPredictions:   r.TotalPredictions,
			CorrectPredictions: r.CorrectPredictions,
			AccuracyPercentage: accuracy,
			AveragePoints:      avgPoints,
		}
	}

	return accuracies, nil
}

// GetAccuracyTrends retrieves accuracy over time periods
func (r *AnalyticsRepository) GetAccuracyTrends(ctx context.Context, userID uint, since time.Time, groupBy string) ([]models.AccuracyTrend, error) {
	var results []struct {
		Period             string
		TotalPredictions   int
		CorrectPredictions int
		TotalPoints        float64
	}

	dateFormat := "YYYY-MM-DD"
	if groupBy == "week" {
		dateFormat = "IYYY-\"W\"IW"
	} else if groupBy == "month" {
		dateFormat = "YYYY-MM"
	}

	query := r.db.WithContext(ctx).Table("scores").
		Select(fmt.Sprintf("TO_CHAR(scored_at, '%s') as period, COUNT(*) as total_predictions, SUM(CASE WHEN points > 0 THEN 1 ELSE 0 END) as correct_predictions, COALESCE(SUM(points), 0) as total_points", dateFormat)).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group(fmt.Sprintf("TO_CHAR(scored_at, '%s')", dateFormat)).
		Order("period ASC")

	if !since.IsZero() {
		query = query.Where("scored_at >= ?", since)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	trends := make([]models.AccuracyTrend, len(results))
	for i, r := range results {
		accuracy := 0.0
		if r.TotalPredictions > 0 {
			accuracy = float64(r.CorrectPredictions) / float64(r.TotalPredictions) * 100
		}
		trends[i] = models.AccuracyTrend{
			Period:             r.Period,
			TotalPredictions:   r.TotalPredictions,
			CorrectPredictions: r.CorrectPredictions,
			AccuracyPercentage: accuracy,
			TotalPoints:        r.TotalPoints,
		}
	}

	return trends, nil
}

// GetPlatformStats retrieves platform-wide statistics
func (r *AnalyticsRepository) GetPlatformStats(ctx context.Context, since time.Time) (*models.PlatformStats, error) {
	var result struct {
		TotalUsers       int
		TotalPredictions int
		TotalCorrect     int
		TotalPoints      float64
	}

	query := r.db.WithContext(ctx).Table("scores").
		Select("COUNT(DISTINCT user_id) as total_users, COUNT(*) as total_predictions, SUM(CASE WHEN points > 0 THEN 1 ELSE 0 END) as total_correct, COALESCE(SUM(points), 0) as total_points").
		Where("deleted_at IS NULL")

	if !since.IsZero() {
		query = query.Where("scored_at >= ?", since)
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	avgAccuracy := 0.0
	avgPoints := 0.0
	if result.TotalPredictions > 0 {
		avgAccuracy = float64(result.TotalCorrect) / float64(result.TotalPredictions) * 100
		avgPoints = result.TotalPoints / float64(result.TotalPredictions)
	}

	return &models.PlatformStats{
		AverageAccuracy:            avgAccuracy,
		AveragePointsPerPrediction: avgPoints,
		TotalUsers:                 result.TotalUsers,
		TotalPredictions:           result.TotalPredictions,
	}, nil
}

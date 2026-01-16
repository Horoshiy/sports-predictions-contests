# Feature: User Analytics Dashboard

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

A comprehensive analytics dashboard showing users their prediction performance statistics including accuracy by sport/league, trends over time, comparison with platform averages, and data export functionality. This feature leverages existing prediction and scoring data to provide actionable insights.

## User Story

As a **sports prediction enthusiast**
I want to **view detailed statistics about my prediction accuracy and performance trends**
So that **I can identify my strengths, improve my predictions, and track my progress over time**

## Problem Statement

Users currently have no visibility into their prediction performance beyond basic leaderboard rankings. They cannot see:
- Which sports/leagues they perform best in
- How their accuracy changes over time
- How they compare to platform averages
- Detailed breakdown of their prediction types

## Solution Statement

Create a new Analytics page with backend aggregation endpoints that:
1. Aggregate user prediction data by sport, league, and prediction type
2. Calculate accuracy percentages and trends over configurable time periods
3. Compare user performance against platform averages
4. Provide CSV export functionality for detailed analysis

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium (4-6 hours)
**Primary Systems Affected**: scoring-service (backend), frontend (new page)
**Dependencies**: Existing predictions, scores, sports tables

---

## CONTEXT REFERENCES

### Relevant Codebase Files - READ BEFORE IMPLEMENTING

**Backend Patterns:**
- `backend/scoring-service/internal/service/scoring_service.go` - gRPC service pattern, response formatting
- `backend/scoring-service/internal/repository/score_repository.go` - Repository interface pattern, GORM queries
- `backend/scoring-service/internal/models/score.go` - Model structure with GORM hooks
- `backend/proto/scoring.proto` - Proto message and service definitions
- `backend/prediction-service/internal/models/prediction.go` - Prediction model with Event relationship

**Frontend Patterns:**
- `frontend/src/services/scoring-service.ts` - API service class pattern
- `frontend/src/types/scoring.types.ts` - TypeScript interface definitions
- `frontend/src/hooks/use-predictions.ts` - React Query hooks pattern
- `frontend/src/pages/ContestsPage.tsx` - Page component with tabs pattern
- `frontend/src/components/leaderboard/LeaderboardTable.tsx` - MaterialReactTable usage

**Database Schema:**
- `scripts/init-db.sql` - Table definitions for scores, predictions, sports, leagues

### New Files to Create

**Backend:**
- `backend/scoring-service/internal/models/analytics.go` - Analytics data models
- `backend/scoring-service/internal/repository/analytics_repository.go` - Analytics queries

**Frontend:**
- `frontend/src/types/analytics.types.ts` - Analytics TypeScript interfaces
- `frontend/src/services/analytics-service.ts` - Analytics API client
- `frontend/src/hooks/use-analytics.ts` - React Query hooks for analytics
- `frontend/src/pages/AnalyticsPage.tsx` - Main analytics page
- `frontend/src/components/analytics/AccuracyChart.tsx` - Accuracy over time chart
- `frontend/src/components/analytics/SportBreakdown.tsx` - Performance by sport
- `frontend/src/components/analytics/PlatformComparison.tsx` - User vs platform stats
- `frontend/src/components/analytics/ExportButton.tsx` - CSV export component

### Patterns to Follow

**Go Naming:**
- Structs: `PascalCase` (e.g., `UserAnalytics`, `SportAccuracy`)
- Functions: `PascalCase` for public, `camelCase` for private
- Files: `snake_case.go`

**TypeScript Naming:**
- Interfaces: `PascalCase` (e.g., `UserAnalytics`, `SportAccuracy`)
- Files: `kebab-case.ts` for services/utils, `PascalCase.tsx` for components

**gRPC Response Pattern:**
```go
return &pb.GetUserAnalyticsResponse{
    Response: &common.Response{
        Success:   true,
        Message:   "Analytics retrieved successfully",
        Code:      int32(common.ErrorCode_SUCCESS),
        Timestamp: timestamppb.Now(),
    },
    Analytics: analyticsProto,
}, nil
```

**React Query Hook Pattern:**
```typescript
export const useUserAnalytics = (userId: number, timeRange: string) => {
  return useQuery({
    queryKey: analyticsKeys.user(userId, timeRange),
    queryFn: () => analyticsService.getUserAnalytics(userId, timeRange),
    enabled: !!userId,
    staleTime: 5 * 60 * 1000,
  })
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Backend - Proto & Models

Define gRPC API contract and data models for analytics aggregation.

**Tasks:**
- Add analytics messages to scoring.proto
- Create analytics model structs
- Define repository interface

### Phase 2: Backend - Repository & Service

Implement database queries and gRPC service methods.

**Tasks:**
- Implement analytics repository with aggregate queries
- Add service methods for analytics endpoints
- Register new endpoints in main.go

### Phase 3: Frontend - Types & Service

Create TypeScript types and API client.

**Tasks:**
- Define analytics TypeScript interfaces
- Create analytics service class
- Create React Query hooks

### Phase 4: Frontend - Components & Page

Build UI components and integrate into app.

**Tasks:**
- Create chart components using Recharts
- Build analytics page with tabs
- Add route and navigation
- Implement CSV export

---

## STEP-BY-STEP TASKS

### Task 1: UPDATE `backend/proto/scoring.proto`

- **IMPLEMENT**: Add analytics messages and RPC methods
- **PATTERN**: Follow existing message structure in scoring.proto
- **ADD** after line 140 (after CalculateScoreResponse):

```protobuf
// Analytics messages
message SportAccuracy {
  string sport_type = 1;
  uint32 total_predictions = 2;
  uint32 correct_predictions = 3;
  double accuracy_percentage = 4;
  double total_points = 5;
}

message LeagueAccuracy {
  uint32 league_id = 1;
  string league_name = 2;
  string sport_type = 3;
  uint32 total_predictions = 4;
  uint32 correct_predictions = 5;
  double accuracy_percentage = 6;
}

message PredictionTypeAccuracy {
  string prediction_type = 1;  // "exact_score", "winner", "over_under"
  uint32 total_predictions = 2;
  uint32 correct_predictions = 3;
  double accuracy_percentage = 4;
  double average_points = 5;
}

message AccuracyTrend {
  string period = 1;  // "2026-01-15", "2026-W03", "2026-01"
  uint32 total_predictions = 2;
  uint32 correct_predictions = 3;
  double accuracy_percentage = 4;
  double total_points = 5;
}

message PlatformStats {
  double average_accuracy = 1;
  double average_points_per_prediction = 2;
  uint32 total_users = 3;
  uint32 total_predictions = 4;
}

message UserAnalytics {
  uint32 user_id = 1;
  uint32 total_predictions = 2;
  uint32 correct_predictions = 3;
  double overall_accuracy = 4;
  double total_points = 5;
  repeated SportAccuracy by_sport = 6;
  repeated LeagueAccuracy by_league = 7;
  repeated PredictionTypeAccuracy by_type = 8;
  repeated AccuracyTrend trends = 9;
  PlatformStats platform_comparison = 10;
  string time_range = 11;  // "7d", "30d", "90d", "all"
}

message GetUserAnalyticsRequest {
  uint32 user_id = 1;
  string time_range = 2;  // "7d", "30d", "90d", "all"
}

message GetUserAnalyticsResponse {
  common.Response response = 1;
  UserAnalytics analytics = 2;
}

message ExportAnalyticsRequest {
  uint32 user_id = 1;
  string time_range = 2;
  string format = 3;  // "csv"
}

message ExportAnalyticsResponse {
  common.Response response = 1;
  string data = 2;  // CSV content
  string filename = 3;
}
```

- **ADD** to ScoringService (after CalculateScore RPC):

```protobuf
  // Analytics
  rpc GetUserAnalytics(GetUserAnalyticsRequest) returns (GetUserAnalyticsResponse) {
    option (google.api.http) = {
      get: "/v1/users/{user_id}/analytics"
    };
  }
  rpc ExportAnalytics(ExportAnalyticsRequest) returns (ExportAnalyticsResponse) {
    option (google.api.http) = {
      get: "/v1/users/{user_id}/analytics/export"
    };
  }
```

- **VALIDATE**: `cat backend/proto/scoring.proto | grep -A5 "UserAnalytics"`

### Task 2: CREATE `backend/scoring-service/internal/models/analytics.go`

- **IMPLEMENT**: Analytics data structures (not GORM models, just structs for aggregation)
- **PATTERN**: Mirror proto messages as Go structs

```go
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
```

- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 3: CREATE `backend/scoring-service/internal/repository/analytics_repository.go`

- **IMPLEMENT**: Repository with aggregate SQL queries
- **PATTERN**: Follow ScoreRepository interface pattern
- **IMPORTS**: gorm.io/gorm, context, models package

```go
package repository

import (
	"context"
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
		Where("user_id = ?", userID)

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

	// Join scores with predictions and events to get sport_type
	query := r.db.WithContext(ctx).Table("scores s").
		Select("COALESCE(e.sport_type, 'unknown') as sport_type, COUNT(*) as total_predictions, SUM(CASE WHEN s.points > 0 THEN 1 ELSE 0 END) as correct_predictions, COALESCE(SUM(s.points), 0) as total_points").
		Joins("LEFT JOIN predictions p ON s.prediction_id = p.id").
		Joins("LEFT JOIN events e ON p.event_id = e.id").
		Where("s.user_id = ?", userID).
		Group("e.sport_type")

	if !since.IsZero() {
		query = query.Where("s.scored_at >= ?", since)
	}

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
func (r *AnalyticsRepository) GetAccuracyByLeague(ctx context.Context, userID uint, since time.Time) ([]models.LeagueAccuracy, error) {
	var results []struct {
		LeagueID           uint
		LeagueName         string
		SportType          string
		TotalPredictions   int
		CorrectPredictions int
	}

	query := r.db.WithContext(ctx).Table("scores s").
		Select("l.id as league_id, l.name as league_name, sp.name as sport_type, COUNT(*) as total_predictions, SUM(CASE WHEN s.points > 0 THEN 1 ELSE 0 END) as correct_predictions").
		Joins("LEFT JOIN predictions p ON s.prediction_id = p.id").
		Joins("LEFT JOIN events e ON p.event_id = e.id").
		Joins("LEFT JOIN matches m ON e.id = m.id").
		Joins("LEFT JOIN leagues l ON m.league_id = l.id").
		Joins("LEFT JOIN sports sp ON l.sport_id = sp.id").
		Where("s.user_id = ? AND l.id IS NOT NULL", userID).
		Group("l.id, l.name, sp.name")

	if !since.IsZero() {
		query = query.Where("s.scored_at >= ?", since)
	}

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

	// Extract prediction type from JSON prediction_data
	query := r.db.WithContext(ctx).Table("scores s").
		Select("COALESCE(p.prediction_data::json->>'type', 'unknown') as prediction_type, COUNT(*) as total_predictions, SUM(CASE WHEN s.points > 0 THEN 1 ELSE 0 END) as correct_predictions, COALESCE(SUM(s.points), 0) as total_points").
		Joins("LEFT JOIN predictions p ON s.prediction_id = p.id").
		Where("s.user_id = ?", userID).
		Group("p.prediction_data::json->>'type'")

	if !since.IsZero() {
		query = query.Where("s.scored_at >= ?", since)
	}

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

	// Group by day for short ranges, week for longer
	dateFormat := "YYYY-MM-DD"
	if groupBy == "week" {
		dateFormat = "IYYY-\"W\"IW"
	} else if groupBy == "month" {
		dateFormat = "YYYY-MM"
	}

	query := r.db.WithContext(ctx).Table("scores").
		Select("TO_CHAR(scored_at, ?) as period, COUNT(*) as total_predictions, SUM(CASE WHEN points > 0 THEN 1 ELSE 0 END) as correct_predictions, COALESCE(SUM(points), 0) as total_points", dateFormat).
		Where("user_id = ?", userID).
		Group("TO_CHAR(scored_at, ?)").
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
		Select("COUNT(DISTINCT user_id) as total_users, COUNT(*) as total_predictions, SUM(CASE WHEN points > 0 THEN 1 ELSE 0 END) as total_correct, COALESCE(SUM(points), 0) as total_points")

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
```

- **VALIDATE**: `cd backend/scoring-service && go build ./...`


### Task 4: UPDATE `backend/scoring-service/internal/service/scoring_service.go`

- **IMPLEMENT**: Add analytics service methods
- **PATTERN**: Follow existing GetUserScores method pattern
- **ADD** to ScoringService struct field:

```go
analyticsRepo repository.AnalyticsRepositoryInterface
```

- **UPDATE** NewScoringService to accept analyticsRepo parameter
- **ADD** these methods after GetUserScores:

```go
// GetUserAnalytics retrieves comprehensive analytics for a user
func (s *ScoringService) GetUserAnalytics(ctx context.Context, req *pb.GetUserAnalyticsRequest) (*pb.GetUserAnalyticsResponse, error) {
	// Validate user ID
	if req.UserId == 0 {
		return &pb.GetUserAnalyticsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "User ID is required",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Parse time range
	timeRange := req.TimeRange
	if timeRange == "" {
		timeRange = "30d"
	}
	since := models.TimeRangeToDate(timeRange)

	// Get overall stats
	analytics, err := s.analyticsRepo.GetUserOverallStats(ctx, uint(req.UserId), since)
	if err != nil {
		log.Printf("[ERROR] Failed to get user overall stats: %v", err)
		return &pb.GetUserAnalyticsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve analytics",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	analytics.TimeRange = timeRange

	// Get accuracy by sport
	bySport, err := s.analyticsRepo.GetAccuracyBySport(ctx, uint(req.UserId), since)
	if err != nil {
		log.Printf("[WARN] Failed to get accuracy by sport: %v", err)
	} else {
		analytics.BySport = bySport
	}

	// Get accuracy by league
	byLeague, err := s.analyticsRepo.GetAccuracyByLeague(ctx, uint(req.UserId), since)
	if err != nil {
		log.Printf("[WARN] Failed to get accuracy by league: %v", err)
	} else {
		analytics.ByLeague = byLeague
	}

	// Get accuracy by type
	byType, err := s.analyticsRepo.GetAccuracyByType(ctx, uint(req.UserId), since)
	if err != nil {
		log.Printf("[WARN] Failed to get accuracy by type: %v", err)
	} else {
		analytics.ByType = byType
	}

	// Get trends - use appropriate grouping based on time range
	groupBy := "day"
	if timeRange == "90d" || timeRange == "all" {
		groupBy = "week"
	}
	trends, err := s.analyticsRepo.GetAccuracyTrends(ctx, uint(req.UserId), since, groupBy)
	if err != nil {
		log.Printf("[WARN] Failed to get accuracy trends: %v", err)
	} else {
		analytics.Trends = trends
	}

	// Get platform stats for comparison
	platformStats, err := s.analyticsRepo.GetPlatformStats(ctx, since)
	if err != nil {
		log.Printf("[WARN] Failed to get platform stats: %v", err)
	} else {
		analytics.PlatformComparison = platformStats
	}

	return &pb.GetUserAnalyticsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Analytics retrieved successfully",
			Code:      int32(common.ErrorCode_SUCCESS),
			Timestamp: timestamppb.Now(),
		},
		Analytics: s.analyticsToProto(analytics),
	}, nil
}

// ExportAnalytics exports user analytics as CSV
func (s *ScoringService) ExportAnalytics(ctx context.Context, req *pb.ExportAnalyticsRequest) (*pb.ExportAnalyticsResponse, error) {
	// Get analytics first
	analyticsResp, err := s.GetUserAnalytics(ctx, &pb.GetUserAnalyticsRequest{
		UserId:    req.UserId,
		TimeRange: req.TimeRange,
	})
	if err != nil || !analyticsResp.Response.Success {
		return &pb.ExportAnalyticsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve analytics for export",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Generate CSV
	csv := s.generateCSV(analyticsResp.Analytics)
	filename := fmt.Sprintf("analytics_%d_%s.csv", req.UserId, req.TimeRange)

	return &pb.ExportAnalyticsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Export generated successfully",
			Code:      int32(common.ErrorCode_SUCCESS),
			Timestamp: timestamppb.Now(),
		},
		Data:     csv,
		Filename: filename,
	}, nil
}

// analyticsToProto converts analytics model to proto
func (s *ScoringService) analyticsToProto(a *models.UserAnalytics) *pb.UserAnalytics {
	proto := &pb.UserAnalytics{
		UserId:             uint32(a.UserID),
		TotalPredictions:   uint32(a.TotalPredictions),
		CorrectPredictions: uint32(a.CorrectPredictions),
		OverallAccuracy:    a.OverallAccuracy,
		TotalPoints:        a.TotalPoints,
		TimeRange:          a.TimeRange,
	}

	// Convert sport accuracies
	for _, s := range a.BySport {
		proto.BySport = append(proto.BySport, &pb.SportAccuracy{
			SportType:          s.SportType,
			TotalPredictions:   uint32(s.TotalPredictions),
			CorrectPredictions: uint32(s.CorrectPredictions),
			AccuracyPercentage: s.AccuracyPercentage,
			TotalPoints:        s.TotalPoints,
		})
	}

	// Convert league accuracies
	for _, l := range a.ByLeague {
		proto.ByLeague = append(proto.ByLeague, &pb.LeagueAccuracy{
			LeagueId:           uint32(l.LeagueID),
			LeagueName:         l.LeagueName,
			SportType:          l.SportType,
			TotalPredictions:   uint32(l.TotalPredictions),
			CorrectPredictions: uint32(l.CorrectPredictions),
			AccuracyPercentage: l.AccuracyPercentage,
		})
	}

	// Convert type accuracies
	for _, t := range a.ByType {
		proto.ByType = append(proto.ByType, &pb.PredictionTypeAccuracy{
			PredictionType:     t.PredictionType,
			TotalPredictions:   uint32(t.TotalPredictions),
			CorrectPredictions: uint32(t.CorrectPredictions),
			AccuracyPercentage: t.AccuracyPercentage,
			AveragePoints:      t.AveragePoints,
		})
	}

	// Convert trends
	for _, t := range a.Trends {
		proto.Trends = append(proto.Trends, &pb.AccuracyTrend{
			Period:             t.Period,
			TotalPredictions:   uint32(t.TotalPredictions),
			CorrectPredictions: uint32(t.CorrectPredictions),
			AccuracyPercentage: t.AccuracyPercentage,
			TotalPoints:        t.TotalPoints,
		})
	}

	// Convert platform stats
	if a.PlatformComparison != nil {
		proto.PlatformComparison = &pb.PlatformStats{
			AverageAccuracy:            a.PlatformComparison.AverageAccuracy,
			AveragePointsPerPrediction: a.PlatformComparison.AveragePointsPerPrediction,
			TotalUsers:                 uint32(a.PlatformComparison.TotalUsers),
			TotalPredictions:           uint32(a.PlatformComparison.TotalPredictions),
		}
	}

	return proto
}

// generateCSV creates CSV content from analytics
func (s *ScoringService) generateCSV(a *pb.UserAnalytics) string {
	var b strings.Builder

	// Header
	b.WriteString("User Analytics Report\n")
	b.WriteString(fmt.Sprintf("User ID,%d\n", a.UserId))
	b.WriteString(fmt.Sprintf("Time Range,%s\n", a.TimeRange))
	b.WriteString("\n")

	// Overall stats
	b.WriteString("Overall Statistics\n")
	b.WriteString("Metric,Value\n")
	b.WriteString(fmt.Sprintf("Total Predictions,%d\n", a.TotalPredictions))
	b.WriteString(fmt.Sprintf("Correct Predictions,%d\n", a.CorrectPredictions))
	b.WriteString(fmt.Sprintf("Overall Accuracy,%.2f%%\n", a.OverallAccuracy))
	b.WriteString(fmt.Sprintf("Total Points,%.2f\n", a.TotalPoints))
	b.WriteString("\n")

	// By Sport
	if len(a.BySport) > 0 {
		b.WriteString("Performance by Sport\n")
		b.WriteString("Sport,Total,Correct,Accuracy,Points\n")
		for _, s := range a.BySport {
			b.WriteString(fmt.Sprintf("%s,%d,%d,%.2f%%,%.2f\n",
				s.SportType, s.TotalPredictions, s.CorrectPredictions, s.AccuracyPercentage, s.TotalPoints))
		}
		b.WriteString("\n")
	}

	// By Type
	if len(a.ByType) > 0 {
		b.WriteString("Performance by Prediction Type\n")
		b.WriteString("Type,Total,Correct,Accuracy,Avg Points\n")
		for _, t := range a.ByType {
			b.WriteString(fmt.Sprintf("%s,%d,%d,%.2f%%,%.2f\n",
				t.PredictionType, t.TotalPredictions, t.CorrectPredictions, t.AccuracyPercentage, t.AveragePoints))
		}
		b.WriteString("\n")
	}

	// Trends
	if len(a.Trends) > 0 {
		b.WriteString("Accuracy Trends\n")
		b.WriteString("Period,Total,Correct,Accuracy,Points\n")
		for _, t := range a.Trends {
			b.WriteString(fmt.Sprintf("%s,%d,%d,%.2f%%,%.2f\n",
				t.Period, t.TotalPredictions, t.CorrectPredictions, t.AccuracyPercentage, t.TotalPoints))
		}
	}

	return b.String()
}
```

- **ADD** import for "strings" at top of file
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 5: UPDATE `backend/scoring-service/cmd/main.go`

- **IMPLEMENT**: Initialize analytics repository and pass to service
- **PATTERN**: Follow existing repository initialization pattern
- **ADD** after streakRepo initialization:

```go
analyticsRepo := repository.NewAnalyticsRepository(db)
```

- **UPDATE** NewScoringService call to include analyticsRepo:

```go
scoringService := service.NewScoringService(scoreRepo, leaderboardRepo, streakRepo, analyticsRepo)
```

- **VALIDATE**: `cd backend/scoring-service && go build ./cmd/main.go`

### Task 6: CREATE `frontend/src/types/analytics.types.ts`

- **IMPLEMENT**: TypeScript interfaces matching proto definitions
- **PATTERN**: Follow scoring.types.ts structure

```typescript
// Analytics types matching backend proto definitions

export interface SportAccuracy {
  sportType: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
  totalPoints: number
}

export interface LeagueAccuracy {
  leagueId: number
  leagueName: string
  sportType: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
}

export interface PredictionTypeAccuracy {
  predictionType: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
  averagePoints: number
}

export interface AccuracyTrend {
  period: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
  totalPoints: number
}

export interface PlatformStats {
  averageAccuracy: number
  averagePointsPerPrediction: number
  totalUsers: number
  totalPredictions: number
}

export interface UserAnalytics {
  userId: number
  totalPredictions: number
  correctPredictions: number
  overallAccuracy: number
  totalPoints: number
  bySport: SportAccuracy[]
  byLeague: LeagueAccuracy[]
  byType: PredictionTypeAccuracy[]
  trends: AccuracyTrend[]
  platformComparison: PlatformStats | null
  timeRange: string
}

// Request types
export interface GetUserAnalyticsRequest {
  userId: number
  timeRange?: TimeRange
}

export interface ExportAnalyticsRequest {
  userId: number
  timeRange?: TimeRange
  format?: 'csv'
}

// Response types
export interface GetUserAnalyticsResponse {
  response: ApiResponse
  analytics: UserAnalytics
}

export interface ExportAnalyticsResponse {
  response: ApiResponse
  data: string
  filename: string
}

// Common types
export type TimeRange = '7d' | '30d' | '90d' | 'all'

export interface ApiResponse {
  success: boolean
  message: string
  code: number
  timestamp: string
}
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/types/analytics.types.ts`

### Task 7: CREATE `frontend/src/services/analytics-service.ts`

- **IMPLEMENT**: API client for analytics endpoints
- **PATTERN**: Follow scoring-service.ts class pattern

```typescript
import grpcClient from './grpc-client'
import type {
  UserAnalytics,
  GetUserAnalyticsRequest,
  GetUserAnalyticsResponse,
  ExportAnalyticsRequest,
  ExportAnalyticsResponse,
  TimeRange,
} from '../types/analytics.types'

class AnalyticsService {
  private basePath = '/v1/users'

  // Get user analytics
  async getUserAnalytics(
    userId: number,
    timeRange: TimeRange = '30d'
  ): Promise<UserAnalytics> {
    const params = new URLSearchParams()
    params.append('time_range', timeRange)

    const response = await grpcClient.get<GetUserAnalyticsResponse>(
      `${this.basePath}/${userId}/analytics?${params.toString()}`
    )
    return response.analytics
  }

  // Export analytics as CSV
  async exportAnalytics(
    userId: number,
    timeRange: TimeRange = '30d'
  ): Promise<{ data: string; filename: string }> {
    const params = new URLSearchParams()
    params.append('time_range', timeRange)
    params.append('format', 'csv')

    const response = await grpcClient.get<ExportAnalyticsResponse>(
      `${this.basePath}/${userId}/analytics/export?${params.toString()}`
    )

    return {
      data: response.data,
      filename: response.filename,
    }
  }

  // Helper to download CSV
  downloadCSV(data: string, filename: string): void {
    const blob = new Blob([data], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', filename)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }
}

// Singleton instance
export const analyticsService = new AnalyticsService()
export default analyticsService
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/services/analytics-service.ts`

### Task 8: CREATE `frontend/src/hooks/use-analytics.ts`

- **IMPLEMENT**: React Query hooks for analytics
- **PATTERN**: Follow use-predictions.ts hook pattern

```typescript
import { useQuery } from '@tanstack/react-query'
import analyticsService from '../services/analytics-service'
import { useToast } from '../contexts/ToastContext'
import type { TimeRange } from '../types/analytics.types'

// Query keys
export const analyticsKeys = {
  all: ['analytics'] as const,
  user: (userId: number, timeRange: TimeRange) =>
    [...analyticsKeys.all, 'user', userId, timeRange] as const,
}

// Fetch user analytics
export const useUserAnalytics = (userId: number, timeRange: TimeRange = '30d') => {
  return useQuery({
    queryKey: analyticsKeys.user(userId, timeRange),
    queryFn: () => analyticsService.getUserAnalytics(userId, timeRange),
    enabled: !!userId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  })
}

// Export analytics hook
export const useExportAnalytics = () => {
  const { showToast } = useToast()

  const exportAnalytics = async (userId: number, timeRange: TimeRange = '30d') => {
    try {
      const { data, filename } = await analyticsService.exportAnalytics(userId, timeRange)
      analyticsService.downloadCSV(data, filename)
      showToast('Analytics exported successfully!', 'success')
    } catch (error) {
      showToast('Failed to export analytics', 'error')
      throw error
    }
  }

  return { exportAnalytics }
}
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/hooks/use-analytics.ts`


### Task 9: CREATE `frontend/src/components/analytics/AccuracyChart.tsx`

- **IMPLEMENT**: Line chart showing accuracy trends over time
- **PATTERN**: Use Recharts (already in package.json)

```tsx
import React from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { Card, CardContent, Typography, Box } from '@mui/material'
import type { AccuracyTrend } from '../../types/analytics.types'

interface AccuracyChartProps {
  trends: AccuracyTrend[]
  title?: string
}

export const AccuracyChart: React.FC<AccuracyChartProps> = ({
  trends,
  title = 'Accuracy Over Time',
}) => {
  if (!trends || trends.length === 0) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>{title}</Typography>
          <Box textAlign="center" py={4}>
            <Typography color="text.secondary">
              No trend data available yet
            </Typography>
          </Box>
        </CardContent>
      </Card>
    )
  }

  const data = trends.map((t) => ({
    period: t.period,
    accuracy: Number(t.accuracyPercentage.toFixed(1)),
    predictions: t.totalPredictions,
    points: Number(t.totalPoints.toFixed(1)),
  }))

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>{title}</Typography>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="period" tick={{ fontSize: 12 }} />
            <YAxis
              yAxisId="left"
              domain={[0, 100]}
              tick={{ fontSize: 12 }}
              label={{ value: 'Accuracy %', angle: -90, position: 'insideLeft' }}
            />
            <YAxis
              yAxisId="right"
              orientation="right"
              tick={{ fontSize: 12 }}
              label={{ value: 'Points', angle: 90, position: 'insideRight' }}
            />
            <Tooltip
              formatter={(value: number, name: string) => {
                if (name === 'accuracy') return [`${value}%`, 'Accuracy']
                if (name === 'points') return [value, 'Points']
                return [value, name]
              }}
            />
            <Legend />
            <Line
              yAxisId="left"
              type="monotone"
              dataKey="accuracy"
              stroke="#1976d2"
              strokeWidth={2}
              dot={{ r: 4 }}
              name="Accuracy %"
            />
            <Line
              yAxisId="right"
              type="monotone"
              dataKey="points"
              stroke="#2e7d32"
              strokeWidth={2}
              dot={{ r: 4 }}
              name="Points"
            />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}

export default AccuracyChart
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/analytics/AccuracyChart.tsx`

### Task 10: CREATE `frontend/src/components/analytics/SportBreakdown.tsx`

- **IMPLEMENT**: Bar chart showing performance by sport
- **PATTERN**: Use Recharts BarChart

```tsx
import React from 'react'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Cell,
} from 'recharts'
import { Card, CardContent, Typography, Box } from '@mui/material'
import type { SportAccuracy } from '../../types/analytics.types'

interface SportBreakdownProps {
  bySport: SportAccuracy[]
  title?: string
}

const COLORS = ['#1976d2', '#2e7d32', '#ed6c02', '#9c27b0', '#d32f2f', '#0288d1']

export const SportBreakdown: React.FC<SportBreakdownProps> = ({
  bySport,
  title = 'Performance by Sport',
}) => {
  if (!bySport || bySport.length === 0) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>{title}</Typography>
          <Box textAlign="center" py={4}>
            <Typography color="text.secondary">
              No sport data available yet
            </Typography>
          </Box>
        </CardContent>
      </Card>
    )
  }

  const data = bySport.map((s) => ({
    sport: s.sportType,
    accuracy: Number(s.accuracyPercentage.toFixed(1)),
    predictions: s.totalPredictions,
    points: Number(s.totalPoints.toFixed(1)),
  }))

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>{title}</Typography>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="sport" tick={{ fontSize: 12 }} />
            <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
            <Tooltip
              formatter={(value: number, name: string) => {
                if (name === 'accuracy') return [`${value}%`, 'Accuracy']
                return [value, name]
              }}
            />
            <Legend />
            <Bar dataKey="accuracy" name="Accuracy %" radius={[4, 4, 0, 0]}>
              {data.map((_, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
        
        {/* Stats table */}
        <Box mt={2}>
          <Typography variant="subtitle2" gutterBottom>Details</Typography>
          <Box component="table" sx={{ width: '100%', fontSize: 14 }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left' }}>Sport</th>
                <th style={{ textAlign: 'right' }}>Predictions</th>
                <th style={{ textAlign: 'right' }}>Correct</th>
                <th style={{ textAlign: 'right' }}>Points</th>
              </tr>
            </thead>
            <tbody>
              {bySport.map((s) => (
                <tr key={s.sportType}>
                  <td>{s.sportType}</td>
                  <td style={{ textAlign: 'right' }}>{s.totalPredictions}</td>
                  <td style={{ textAlign: 'right' }}>{s.correctPredictions}</td>
                  <td style={{ textAlign: 'right' }}>{s.totalPoints.toFixed(1)}</td>
                </tr>
              ))}
            </tbody>
          </Box>
        </Box>
      </CardContent>
    </Card>
  )
}

export default SportBreakdown
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/analytics/SportBreakdown.tsx`

### Task 11: CREATE `frontend/src/components/analytics/PlatformComparison.tsx`

- **IMPLEMENT**: Comparison cards showing user vs platform stats
- **PATTERN**: Use MUI Card components

```tsx
import React from 'react'
import {
  Card,
  CardContent,
  Typography,
  Box,
  Grid,
  Chip,
  LinearProgress,
} from '@mui/material'
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  Remove as NeutralIcon,
} from '@mui/icons-material'
import type { UserAnalytics, PlatformStats } from '../../types/analytics.types'

interface PlatformComparisonProps {
  userStats: {
    accuracy: number
    avgPoints: number
    totalPredictions: number
  }
  platformStats: PlatformStats | null
}

const ComparisonCard: React.FC<{
  label: string
  userValue: number
  platformValue: number
  format?: 'percent' | 'number'
  higherIsBetter?: boolean
}> = ({ label, userValue, platformValue, format = 'number', higherIsBetter = true }) => {
  const diff = userValue - platformValue
  const isAbove = diff > 0
  const isBetter = higherIsBetter ? isAbove : !isAbove

  const formatValue = (v: number) => {
    if (format === 'percent') return `${v.toFixed(1)}%`
    return v.toFixed(1)
  }

  return (
    <Card variant="outlined">
      <CardContent>
        <Typography variant="subtitle2" color="text.secondary" gutterBottom>
          {label}
        </Typography>
        <Box display="flex" alignItems="baseline" gap={1}>
          <Typography variant="h4" component="span">
            {formatValue(userValue)}
          </Typography>
          <Chip
            size="small"
            icon={
              Math.abs(diff) < 0.1 ? (
                <NeutralIcon />
              ) : isBetter ? (
                <TrendingUpIcon />
              ) : (
                <TrendingDownIcon />
              )
            }
            label={`${diff > 0 ? '+' : ''}${formatValue(diff)}`}
            color={Math.abs(diff) < 0.1 ? 'default' : isBetter ? 'success' : 'error'}
            variant="outlined"
          />
        </Box>
        <Box mt={2}>
          <Typography variant="caption" color="text.secondary">
            Platform Average: {formatValue(platformValue)}
          </Typography>
          <LinearProgress
            variant="determinate"
            value={Math.min((userValue / Math.max(platformValue, 1)) * 50, 100)}
            sx={{ mt: 1, height: 8, borderRadius: 4 }}
            color={isBetter ? 'success' : 'error'}
          />
        </Box>
      </CardContent>
    </Card>
  )
}

export const PlatformComparison: React.FC<PlatformComparisonProps> = ({
  userStats,
  platformStats,
}) => {
  if (!platformStats) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>Platform Comparison</Typography>
          <Box textAlign="center" py={4}>
            <Typography color="text.secondary">
              Platform statistics not available
            </Typography>
          </Box>
        </CardContent>
      </Card>
    )
  }

  const userAvgPoints = userStats.totalPredictions > 0
    ? userStats.avgPoints / userStats.totalPredictions
    : 0

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          How You Compare
        </Typography>
        <Typography variant="body2" color="text.secondary" gutterBottom>
          Your performance vs {platformStats.totalUsers.toLocaleString()} users
        </Typography>
        
        <Grid container spacing={2} mt={1}>
          <Grid item xs={12} md={6}>
            <ComparisonCard
              label="Accuracy"
              userValue={userStats.accuracy}
              platformValue={platformStats.averageAccuracy}
              format="percent"
              higherIsBetter={true}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <ComparisonCard
              label="Avg Points per Prediction"
              userValue={userAvgPoints}
              platformValue={platformStats.averagePointsPerPrediction}
              format="number"
              higherIsBetter={true}
            />
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  )
}

export default PlatformComparison
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/analytics/PlatformComparison.tsx`

### Task 12: CREATE `frontend/src/components/analytics/ExportButton.tsx`

- **IMPLEMENT**: Button component for CSV export
- **PATTERN**: Use MUI Button with loading state

```tsx
import React, { useState } from 'react'
import { Button, CircularProgress } from '@mui/material'
import { Download as DownloadIcon } from '@mui/icons-material'
import { useExportAnalytics } from '../../hooks/use-analytics'
import type { TimeRange } from '../../types/analytics.types'

interface ExportButtonProps {
  userId: number
  timeRange: TimeRange
}

export const ExportButton: React.FC<ExportButtonProps> = ({ userId, timeRange }) => {
  const [isExporting, setIsExporting] = useState(false)
  const { exportAnalytics } = useExportAnalytics()

  const handleExport = async () => {
    setIsExporting(true)
    try {
      await exportAnalytics(userId, timeRange)
    } finally {
      setIsExporting(false)
    }
  }

  return (
    <Button
      variant="outlined"
      startIcon={isExporting ? <CircularProgress size={20} /> : <DownloadIcon />}
      onClick={handleExport}
      disabled={isExporting}
    >
      {isExporting ? 'Exporting...' : 'Export CSV'}
    </Button>
  )
}

export default ExportButton
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/analytics/ExportButton.tsx`

### Task 13: CREATE `frontend/src/pages/AnalyticsPage.tsx`

- **IMPLEMENT**: Main analytics page with all components
- **PATTERN**: Follow ContestsPage.tsx structure with tabs

```tsx
import React, { useState } from 'react'
import {
  Box,
  Typography,
  Paper,
  Grid,
  Card,
  CardContent,
  ToggleButton,
  ToggleButtonGroup,
  CircularProgress,
  Alert,
} from '@mui/material'
import {
  Analytics as AnalyticsIcon,
  TrendingUp as TrendingUpIcon,
  EmojiEvents as TrophyIcon,
} from '@mui/icons-material'
import { useAuth } from '../contexts/AuthContext'
import { useUserAnalytics } from '../hooks/use-analytics'
import { AccuracyChart } from '../components/analytics/AccuracyChart'
import { SportBreakdown } from '../components/analytics/SportBreakdown'
import { PlatformComparison } from '../components/analytics/PlatformComparison'
import { ExportButton } from '../components/analytics/ExportButton'
import type { TimeRange } from '../types/analytics.types'

const StatCard: React.FC<{
  title: string
  value: string | number
  subtitle?: string
  icon: React.ReactNode
  color?: string
}> = ({ title, value, subtitle, icon, color = 'primary.main' }) => (
  <Card>
    <CardContent>
      <Box display="flex" alignItems="center" gap={2}>
        <Box
          sx={{
            p: 1.5,
            borderRadius: 2,
            bgcolor: `${color}15`,
            color: color,
          }}
        >
          {icon}
        </Box>
        <Box>
          <Typography variant="h4" component="div">
            {value}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {title}
          </Typography>
          {subtitle && (
            <Typography variant="caption" color="text.secondary">
              {subtitle}
            </Typography>
          )}
        </Box>
      </Box>
    </CardContent>
  </Card>
)

export const AnalyticsPage: React.FC = () => {
  const { user } = useAuth()
  const [timeRange, setTimeRange] = useState<TimeRange>('30d')

  const {
    data: analytics,
    isLoading,
    error,
  } = useUserAnalytics(user?.id || 0, timeRange)

  const handleTimeRangeChange = (
    _: React.MouseEvent<HTMLElement>,
    newRange: TimeRange | null
  ) => {
    if (newRange) {
      setTimeRange(newRange)
    }
  }

  if (!user) {
    return (
      <Alert severity="warning">
        Please log in to view your analytics.
      </Alert>
    )
  }

  if (error) {
    return (
      <Alert severity="error">
        Failed to load analytics. Please try again later.
      </Alert>
    )
  }

  return (
    <Box>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4" component="h1" gutterBottom>
            Your Analytics
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Track your prediction performance and identify areas for improvement
          </Typography>
        </Box>
        <Box display="flex" gap={2} alignItems="center">
          <ToggleButtonGroup
            value={timeRange}
            exclusive
            onChange={handleTimeRangeChange}
            size="small"
          >
            <ToggleButton value="7d">7 Days</ToggleButton>
            <ToggleButton value="30d">30 Days</ToggleButton>
            <ToggleButton value="90d">90 Days</ToggleButton>
            <ToggleButton value="all">All Time</ToggleButton>
          </ToggleButtonGroup>
          <ExportButton userId={user.id} timeRange={timeRange} />
        </Box>
      </Box>

      {isLoading ? (
        <Box display="flex" justifyContent="center" py={8}>
          <CircularProgress />
        </Box>
      ) : analytics ? (
        <>
          {/* Overview Stats */}
          <Grid container spacing={3} mb={3}>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Total Predictions"
                value={analytics.totalPredictions}
                icon={<AnalyticsIcon />}
                color="primary.main"
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Correct Predictions"
                value={analytics.correctPredictions}
                subtitle={`${analytics.overallAccuracy.toFixed(1)}% accuracy`}
                icon={<TrendingUpIcon />}
                color="success.main"
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Total Points"
                value={analytics.totalPoints.toFixed(1)}
                icon={<TrophyIcon />}
                color="warning.main"
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Avg Points/Prediction"
                value={
                  analytics.totalPredictions > 0
                    ? (analytics.totalPoints / analytics.totalPredictions).toFixed(2)
                    : '0'
                }
                icon={<AnalyticsIcon />}
                color="info.main"
              />
            </Grid>
          </Grid>

          {/* Charts Row */}
          <Grid container spacing={3} mb={3}>
            <Grid item xs={12} lg={8}>
              <AccuracyChart trends={analytics.trends} />
            </Grid>
            <Grid item xs={12} lg={4}>
              <PlatformComparison
                userStats={{
                  accuracy: analytics.overallAccuracy,
                  avgPoints: analytics.totalPoints,
                  totalPredictions: analytics.totalPredictions,
                }}
                platformStats={analytics.platformComparison}
              />
            </Grid>
          </Grid>

          {/* Sport Breakdown */}
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <SportBreakdown bySport={analytics.bySport} />
            </Grid>
          </Grid>

          {/* Prediction Type Breakdown */}
          {analytics.byType && analytics.byType.length > 0 && (
            <Paper sx={{ mt: 3, p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Performance by Prediction Type
              </Typography>
              <Grid container spacing={2}>
                {analytics.byType.map((t) => (
                  <Grid item xs={12} sm={6} md={4} key={t.predictionType}>
                    <Card variant="outlined">
                      <CardContent>
                        <Typography variant="subtitle1" gutterBottom>
                          {t.predictionType.replace('_', ' ').toUpperCase()}
                        </Typography>
                        <Typography variant="h5" color="primary">
                          {t.accuracyPercentage.toFixed(1)}%
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {t.correctPredictions} / {t.totalPredictions} correct
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Avg: {t.averagePoints.toFixed(2)} pts
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            </Paper>
          )}
        </>
      ) : (
        <Paper sx={{ p: 4, textAlign: 'center' }}>
          <AnalyticsIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No Analytics Data Yet
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Start making predictions to see your performance analytics here.
          </Typography>
        </Paper>
      )}
    </Box>
  )
}

export default AnalyticsPage
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit src/pages/AnalyticsPage.tsx`

### Task 14: UPDATE `frontend/src/App.tsx`

- **IMPLEMENT**: Add analytics route and navigation
- **PATTERN**: Follow existing route pattern
- **ADD** import at top:

```typescript
import AnalyticsPage from './pages/AnalyticsPage'
```

- **ADD** navigation button after Predictions in AppBarContent (around line 65):

```tsx
<Button color="inherit" component={Link} to="/analytics">Analytics</Button>
```

- **ADD** route after /predictions route (around line 115):

```tsx
<Route 
  path="/analytics" 
  element={
    <ProtectedRoute>
      <AnalyticsPage />
    </ProtectedRoute>
  } 
/>
```

- **VALIDATE**: `cd frontend && npx tsc --noEmit`

---

## TESTING STRATEGY

### Unit Tests

**Backend (Go):**
- Test `TimeRangeToDate` function with all time ranges
- Test analytics model struct initialization
- Test CSV generation with sample data

**Frontend (TypeScript):**
- Test analytics service methods
- Test React Query hooks with mock data
- Test chart components render correctly

### Integration Tests

- Test full analytics API flow: request → service → repository → response
- Test CSV export generates valid CSV format
- Test frontend displays data correctly from API

### Edge Cases

- User with no predictions (empty analytics)
- User with predictions but no scores yet
- Time range with no data
- Platform with single user (comparison edge case)
- Very large datasets (pagination/performance)

---

## VALIDATION COMMANDS

### Level 1: Syntax & Style

```bash
# Backend
cd backend/scoring-service && go fmt ./... && go vet ./...

# Frontend
cd frontend && npm run lint
```

### Level 2: Type Checking

```bash
# Backend
cd backend/scoring-service && go build ./...

# Frontend
cd frontend && npx tsc --noEmit
```

### Level 3: Unit Tests

```bash
# Backend
cd backend && go test ./scoring-service/...

# Frontend
cd frontend && npm test
```

### Level 4: Manual Validation

```bash
# Start services
docker-compose up -d postgres redis
cd backend/scoring-service && go run cmd/main.go &
cd frontend && npm run dev

# Test API endpoint
curl http://localhost:8080/v1/users/1/analytics?time_range=30d

# Test export
curl http://localhost:8080/v1/users/1/analytics/export?time_range=30d
```

---

## ACCEPTANCE CRITERIA

- [ ] Backend analytics endpoints return correct aggregated data
- [ ] Frontend displays accuracy trends chart
- [ ] Frontend displays sport breakdown chart
- [ ] Frontend displays platform comparison
- [ ] Time range selector filters data correctly
- [ ] CSV export downloads valid file
- [ ] Empty state displays when no data
- [ ] Loading states show during data fetch
- [ ] Error states handle API failures gracefully
- [ ] Navigation to /analytics works from header
- [ ] All TypeScript types match backend proto

---

## COMPLETION CHECKLIST

- [ ] All 14 tasks completed in order
- [ ] Each task validation passed
- [ ] Backend compiles without errors
- [ ] Frontend compiles without errors
- [ ] Manual testing confirms feature works
- [ ] Analytics page accessible from navigation
- [ ] Charts render with sample data
- [ ] Export downloads CSV file

---

## NOTES

### Design Decisions

1. **Aggregate queries in repository**: Keeps complex SQL in one place, service layer stays clean
2. **Time range as query param**: Allows caching by time range, easy to extend
3. **CSV export on backend**: Ensures consistent format, reduces frontend complexity
4. **Recharts for visualization**: Already in package.json, well-documented, responsive

### Trade-offs

1. **No real-time updates**: Analytics are point-in-time, 5-minute cache is acceptable
2. **Limited league data**: Depends on match→league relationship which may not exist for all predictions
3. **Platform stats include all users**: Could be filtered by contest in future

### Future Enhancements

- Add league-specific breakdown page
- Add prediction type deep-dive
- Add comparison with specific users
- Add goal-setting and progress tracking
- Add email reports (weekly/monthly)

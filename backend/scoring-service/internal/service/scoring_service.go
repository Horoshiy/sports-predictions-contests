package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
	"github.com/sports-prediction-contests/scoring-service/internal/repository"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/scoring"
	pb "github.com/sports-prediction-contests/shared/proto/scoring"
	"github.com/sports-prediction-contests/shared/proto/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScoringService implements the gRPC ScoringService
type ScoringService struct {
	pb.UnimplementedScoringServiceServer
	scoreRepo       repository.ScoreRepositoryInterface
	leaderboardRepo repository.LeaderboardRepositoryInterface
	streakRepo      repository.StreakRepositoryInterface
	analyticsRepo   repository.AnalyticsRepositoryInterface
}

// NewScoringService creates a new ScoringService instance
func NewScoringService(scoreRepo repository.ScoreRepositoryInterface, leaderboardRepo repository.LeaderboardRepositoryInterface, streakRepo repository.StreakRepositoryInterface, analyticsRepo repository.AnalyticsRepositoryInterface) *ScoringService {
	return &ScoringService{
		scoreRepo:       scoreRepo,
		leaderboardRepo: leaderboardRepo,
		streakRepo:      streakRepo,
		analyticsRepo:   analyticsRepo,
	}
}

// PredictionData represents the structure of prediction data
type PredictionData struct {
	Type       string      `json:"type"`        // "exact_score", "winner", "over_under"
	HomeScore  *int        `json:"home_score"`  // For exact score predictions
	AwayScore  *int        `json:"away_score"`  // For exact score predictions
	Winner     *string     `json:"winner"`      // "home", "away", "draw"
	OverUnder  *string     `json:"over_under"`  // "over", "under"
	Threshold  *float64    `json:"threshold"`   // For over/under predictions
	Value      interface{} `json:"value"`       // Generic value for other prediction types
	Props      []PropPrediction `json:"props,omitempty"` // Props predictions
}

// PropPrediction represents a single prop prediction
type PropPrediction struct {
	PropTypeID  uint    `json:"prop_type_id"`
	PropSlug    string  `json:"prop_slug"`
	Line        float64 `json:"line"`
	Selection   string  `json:"selection"`
	PlayerID    string  `json:"player_id,omitempty"`
	PointsValue float64 `json:"points_value"`
}

// ResultData represents the structure of event result data
type ResultData struct {
	HomeScore   int                    `json:"home_score"`
	AwayScore   int                    `json:"away_score"`
	Winner      string                 `json:"winner"`
	TotalGoals  int                    `json:"total_goals"`
	Stats       map[string]interface{} `json:"stats,omitempty"`
	PlayerStats map[string]interface{} `json:"player_stats,omitempty"`
}

// CreateScore creates a new score record
func (s *ScoringService) CreateScore(ctx context.Context, req *pb.CreateScoreRequest) (*pb.CreateScoreResponse, error) {
	// Extract user ID from JWT token for authorization
	_, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.CreateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get or create user streak
	streak, err := s.streakRepo.GetOrCreate(ctx, uint(req.ContestId), uint(req.UserId))
	if err != nil {
		log.Printf("[ERROR] Failed to get/create streak: %v", err)
		return &pb.CreateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to process streak",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Determine if prediction was correct and update streak
	basePoints := req.Points
	// Update streak first, then calculate multiplier based on new streak value
	isCorrect := basePoints > 0
	if isCorrect {
		streak.IncrementStreak(uint(req.PredictionId))
	} else {
		streak.ResetStreak(uint(req.PredictionId))
	}

	// Multiplier is based on the updated streak value
	multiplier := streak.GetMultiplier()

	// Calculate time coefficient based on submission time vs event date
	timeCoefficient := 1.0
	if req.SubmittedAt != nil && req.EventDate != nil {
		timeCoefficient = models.CalculateTimeCoefficient(
			req.SubmittedAt.AsTime(),
			req.EventDate.AsTime(),
		)
	}

	// Apply both multipliers
	finalPoints := basePoints * multiplier * timeCoefficient

	// Update streak in database - fail if this fails to maintain consistency
	if err := s.streakRepo.Update(ctx, streak); err != nil {
		log.Printf("[ERROR] Failed to update streak: %v", err)
		return &pb.CreateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to update streak",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Create score model with multiplied points
	score := &models.Score{
		UserID:          uint(req.UserId),
		ContestID:       uint(req.ContestId),
		PredictionID:    uint(req.PredictionId),
		Points:          finalPoints,
		TimeCoefficient: timeCoefficient,
	}

	// Save to database
	if err := s.scoreRepo.Create(ctx, score); err != nil {
		log.Printf("[ERROR] Failed to create score: %v", err)
		return &pb.CreateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to create score",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update leaderboard
	totalPoints, err := s.scoreRepo.GetTotalPointsByContestAndUser(ctx, uint(req.ContestId), uint(req.UserId))
	if err != nil {
		log.Printf("[ERROR] Failed to get total points: %v", err)
		return &pb.CreateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to update leaderboard after score creation",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	if err := s.leaderboardRepo.UpsertUserScore(ctx, uint(req.ContestId), uint(req.UserId), totalPoints); err != nil {
		log.Printf("[ERROR] Failed to update leaderboard: %v", err)
		return &pb.CreateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Score created but leaderboard update failed",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Score created: user=%d, contest=%d, base=%.2f, streak=%.2fx, time=%.2fx, final=%.2f",
		uint(req.UserId), uint(req.ContestId), basePoints, multiplier, timeCoefficient, finalPoints)

	return &pb.CreateScoreResponse{
		Response: &common.Response{
			Success:   true,
			Message:   fmt.Sprintf("Score created successfully (%.2fx streak, %.2fx time)", multiplier, timeCoefficient),
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Score: s.modelToProto(score),
	}, nil
}

// GetScore retrieves a score by ID
func (s *ScoringService) GetScore(ctx context.Context, req *pb.GetScoreRequest) (*pb.GetScoreResponse, error) {
	score, err := s.scoreRepo.GetByID(ctx, uint(req.Id))
	if err != nil {
		log.Printf("[ERROR] Failed to get score: %v", err)
		return &pb.GetScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Score not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetScoreResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Score retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Score: s.modelToProto(score),
	}, nil
}

// CalculateScore calculates points based on prediction accuracy
func (s *ScoringService) CalculateScore(ctx context.Context, req *pb.CalculateScoreRequest) (*pb.CalculateScoreResponse, error) {
	// Parse prediction data
	var predictionData PredictionData
	if err := json.Unmarshal([]byte(req.PredictionData), &predictionData); err != nil {
		return &pb.CalculateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Invalid prediction data format",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Parse result data
	var resultData ResultData
	if err := json.Unmarshal([]byte(req.ResultData), &resultData); err != nil {
		return &pb.CalculateScoreResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Invalid result data format",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Calculate points based on prediction type
	points, details := s.calculatePoints(predictionData, resultData)

	detailsJSON, _ := json.Marshal(details)

	return &pb.CalculateScoreResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Score calculated successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Points:             points,
		CalculationDetails: string(detailsJSON),
	}, nil
}

// calculatePoints implements the scoring algorithm
func (s *ScoringService) calculatePoints(prediction PredictionData, result ResultData) (float64, map[string]interface{}) {
	details := map[string]interface{}{
		"prediction_type": prediction.Type,
		"result":          result,
	}

	switch prediction.Type {
	case "exact_score":
		return s.calculateExactScorePoints(prediction, result, details)
	case "winner":
		return s.calculateWinnerPoints(prediction, result, details)
	case "over_under":
		return s.calculateOverUnderPoints(prediction, result, details)
	case "props":
		return s.calculatePropsPoints(prediction, result, details)
	default:
		details["error"] = "Unknown prediction type"
		return 0, details
	}
}

// calculateExactScorePoints calculates points for exact score predictions (uses default rules)
func (s *ScoringService) calculateExactScorePoints(prediction PredictionData, result ResultData, details map[string]interface{}) (float64, map[string]interface{}) {
	// Use default rules for backward compatibility
	defaultRules := scoring.DefaultStandardRules()
	return s.calculateExactScorePointsWithRules(prediction, result, details, &defaultRules)
}

// calculateExactScorePointsWithRules calculates points using custom scoring rules
func (s *ScoringService) calculateExactScorePointsWithRules(prediction PredictionData, result ResultData, details map[string]interface{}, rules *scoring.StandardScoringRules) (float64, map[string]interface{}) {
	if prediction.HomeScore == nil || prediction.AwayScore == nil {
		details["error"] = "Missing score prediction"
		return 0, details
	}

	predictedHome := *prediction.HomeScore
	predictedAway := *prediction.AwayScore

	details["predicted_score"] = fmt.Sprintf("%d-%d", predictedHome, predictedAway)
	details["actual_score"] = fmt.Sprintf("%d-%d", result.HomeScore, result.AwayScore)

	// Handle "any_other" prediction type
	if prediction.Type == "any_other" {
		isOther := result.HomeScore > 4 || result.AwayScore > 4
		details["result_is_other"] = isOther
		if isOther {
			details["match_type"] = "any_other_correct"
			return rules.AnyOther, details
		}
		details["match_type"] = "any_other_incorrect"
		return 0, details
	}

	// Exact match
	if predictedHome == result.HomeScore && predictedAway == result.AwayScore {
		details["match_type"] = "exact"
		return rules.ExactScore, details
	}

	// Correct goal difference
	predictedDiff := predictedHome - predictedAway
	actualDiff := result.HomeScore - result.AwayScore
	if predictedDiff == actualDiff {
		details["match_type"] = "goal_difference"
		return rules.GoalDifference, details
	}

	// Correct winner with possible team goals bonus
	predictedWinner := s.determineWinner(predictedHome, predictedAway)
	actualWinner := s.determineWinner(result.HomeScore, result.AwayScore)
	
	if predictedWinner == actualWinner {
		// Check if one team's goals match
		homeGoalsMatch := predictedHome == result.HomeScore
		awayGoalsMatch := predictedAway == result.AwayScore
		
		if homeGoalsMatch || awayGoalsMatch {
			details["match_type"] = "outcome_plus_team_goals"
			details["home_goals_match"] = homeGoalsMatch
			details["away_goals_match"] = awayGoalsMatch
			return rules.CorrectOutcome + rules.OutcomePlusTeamGoals, details
		}
		
		details["match_type"] = "correct_outcome"
		return rules.CorrectOutcome, details
	}

	details["match_type"] = "none"
	return 0, details
}

// CalculateWithContestRules calculates score using contest-specific rules
func (s *ScoringService) CalculateWithContestRules(predictionData, resultData, rulesJSON string) (float64, map[string]interface{}) {
	var prediction PredictionData
	if err := json.Unmarshal([]byte(predictionData), &prediction); err != nil {
		return 0, map[string]interface{}{"error": "Invalid prediction data"}
	}

	var result ResultData
	if err := json.Unmarshal([]byte(resultData), &result); err != nil {
		return 0, map[string]interface{}{"error": "Invalid result data"}
	}

	rules, err := scoring.ParseRules(rulesJSON)
	if err != nil {
		return 0, map[string]interface{}{"error": "Invalid rules: " + err.Error()}
	}

	details := map[string]interface{}{
		"prediction_type": prediction.Type,
		"contest_type":    rules.Type,
	}

	switch rules.Type {
	case scoring.ContestTypeStandard:
		return s.calculateExactScorePointsWithRules(prediction, result, details, rules.Standard)
	case scoring.ContestTypeRisky:
		return s.calculateRiskyPoints(prediction, result, details, rules.Risky)
	default:
		details["error"] = "Unknown contest type"
		return 0, details
	}
}

// calculateRiskyPoints calculates points for risky predictions
func (s *ScoringService) calculateRiskyPoints(prediction PredictionData, result ResultData, details map[string]interface{}, rules *scoring.RiskyScoringRules) (float64, map[string]interface{}) {
	// Risky predictions store selected events in prediction.Props or a custom field
	// For now, we'll use a simple calculation based on the prediction data
	
	// Parse risky selections from prediction data (expected format: {"risky_selections": ["penalty", "red_card"]})
	type RiskyPredictionData struct {
		Selections []string `json:"risky_selections"`
	}
	
	// Try to extract risky selections from the Value field
	selectionsData, ok := prediction.Value.(map[string]interface{})
	if !ok {
		details["error"] = "Invalid risky prediction format"
		return 0, details
	}
	
	selectionsRaw, ok := selectionsData["risky_selections"].([]interface{})
	if !ok {
		details["error"] = "Missing risky_selections"
		return 0, details
	}
	
	selections := make([]string, len(selectionsRaw))
	for i, v := range selectionsRaw {
		selections[i], _ = v.(string)
	}
	
	// Get outcomes from result stats
	outcomes := make(map[string]bool)
	if result.Stats != nil {
		for key, val := range result.Stats {
			if boolVal, ok := val.(bool); ok {
				outcomes[key] = boolVal
			}
		}
	}
	
	calc := scoring.NewCalculator(&scoring.ContestRules{Type: scoring.ContestTypeRisky, Risky: rules})
	calcResult := calc.CalculateRisky(selections, outcomes)
	
	for k, v := range calcResult.Details {
		details[k] = v
	}
	
	return calcResult.Points, details
}

// calculateWinnerPoints calculates points for winner predictions
func (s *ScoringService) calculateWinnerPoints(prediction PredictionData, result ResultData, details map[string]interface{}) (float64, map[string]interface{}) {
	if prediction.Winner == nil {
		details["error"] = "Missing winner prediction"
		return 0, details
	}

	predictedWinner := *prediction.Winner
	details["predicted_winner"] = predictedWinner
	details["actual_winner"] = result.Winner

	if predictedWinner == result.Winner {
		details["match"] = true
		return 3, details
	}

	details["match"] = false
	return 0, details
}

// calculateOverUnderPoints calculates points for over/under predictions
func (s *ScoringService) calculateOverUnderPoints(prediction PredictionData, result ResultData, details map[string]interface{}) (float64, map[string]interface{}) {
	if prediction.OverUnder == nil || prediction.Threshold == nil {
		details["error"] = "Missing over/under prediction or threshold"
		return 0, details
	}

	predictedOverUnder := *prediction.OverUnder
	threshold := *prediction.Threshold
	totalGoals := float64(result.TotalGoals)

	details["predicted"] = predictedOverUnder
	details["threshold"] = threshold
	details["total_goals"] = totalGoals

	var correct bool
	if predictedOverUnder == "over" {
		correct = totalGoals > threshold
	} else if predictedOverUnder == "under" {
		correct = totalGoals < threshold
	} else {
		details["error"] = "Invalid over/under value"
		return 0, details
	}

	details["correct"] = correct
	if correct {
		return 2, details
	}

	return 0, details
}

// determineWinner determines the winner based on scores
func (s *ScoringService) determineWinner(homeScore, awayScore int) string {
	if homeScore > awayScore {
		return "home"
	} else if awayScore > homeScore {
		return "away"
	}
	return "draw"
}

// GetUserScores retrieves all scores for a user in a contest
func (s *ScoringService) GetUserScores(ctx context.Context, req *pb.GetUserScoresRequest) (*pb.GetUserScoresResponse, error) {
	scores, err := s.scoreRepo.GetByContestAndUser(ctx, uint(req.ContestId), uint(req.UserId))
	if err != nil {
		log.Printf("[ERROR] Failed to get user scores: %v", err)
		return &pb.GetUserScoresResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve user scores",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Calculate total points
	totalPoints, err := s.scoreRepo.GetTotalPointsByContestAndUser(ctx, uint(req.ContestId), uint(req.UserId))
	if err != nil {
		log.Printf("[ERROR] Failed to get total points: %v", err)
		totalPoints = 0
	}

	// Convert to proto
	protoScores := make([]*pb.Score, len(scores))
	for i, score := range scores {
		protoScores[i] = s.modelToProto(score)
	}

	return &pb.GetUserScoresResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "User scores retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Scores:      protoScores,
		TotalPoints: totalPoints,
	}, nil
}

// modelToProto converts a Score model to protobuf message
func (s *ScoringService) modelToProto(score *models.Score) *pb.Score {
	return &pb.Score{
		Id:              uint32(score.ID),
		UserId:          uint32(score.UserID),
		ContestId:       uint32(score.ContestID),
		PredictionId:    uint32(score.PredictionID),
		Points:          score.Points,
		TimeCoefficient: score.TimeCoefficient,
		ScoredAt:        timestamppb.New(score.ScoredAt),
		CreatedAt:       timestamppb.New(score.CreatedAt),
		UpdatedAt:       timestamppb.New(score.UpdatedAt),
	}
}

// GetUserAnalytics retrieves comprehensive analytics for a user
func (s *ScoringService) GetUserAnalytics(ctx context.Context, req *pb.GetUserAnalyticsRequest) (*pb.GetUserAnalyticsResponse, error) {
	if uint(req.UserId) == 0 {
		return &pb.GetUserAnalyticsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "User ID is required",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	timeRange := req.TimeRange
	if timeRange == "" {
		timeRange = "30d"
	}
	since := models.TimeRangeToDate(timeRange)

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

	if bySport, err := s.analyticsRepo.GetAccuracyBySport(ctx, uint(req.UserId), since); err != nil {
		log.Printf("[WARN] Failed to get accuracy by sport: %v", err)
	} else {
		analytics.BySport = bySport
	}

	if byLeague, err := s.analyticsRepo.GetAccuracyByLeague(ctx, uint(req.UserId), since); err != nil {
		log.Printf("[WARN] Failed to get accuracy by league: %v", err)
	} else {
		analytics.ByLeague = byLeague
	}

	if byType, err := s.analyticsRepo.GetAccuracyByType(ctx, uint(req.UserId), since); err != nil {
		log.Printf("[WARN] Failed to get accuracy by type: %v", err)
	} else {
		analytics.ByType = byType
	}

	groupBy := "day"
	if timeRange == "90d" || timeRange == "all" {
		groupBy = "week"
	}
	if trends, err := s.analyticsRepo.GetAccuracyTrends(ctx, uint(req.UserId), since, groupBy); err != nil {
		log.Printf("[WARN] Failed to get accuracy trends: %v", err)
	} else {
		analytics.Trends = trends
	}

	if platformStats, err := s.analyticsRepo.GetPlatformStats(ctx, since); err != nil {
		log.Printf("[WARN] Failed to get platform stats: %v", err)
	} else {
		analytics.PlatformComparison = platformStats
	}

	return &pb.GetUserAnalyticsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Analytics retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Analytics: s.analyticsToProto(analytics),
	}, nil
}

// ExportAnalytics exports user analytics as CSV
func (s *ScoringService) ExportAnalytics(ctx context.Context, req *pb.ExportAnalyticsRequest) (*pb.ExportAnalyticsResponse, error) {
	timeRange := req.TimeRange
	if timeRange == "" {
		timeRange = "30d"
	}

	analyticsResp, err := s.GetUserAnalytics(ctx, &pb.GetUserAnalyticsRequest{
		UserId:    req.UserId,
		TimeRange: timeRange,
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

	csv := s.generateCSV(analyticsResp.Analytics)
	filename := fmt.Sprintf("analytics_%d_%s.csv", uint(req.UserId), timeRange)

	return &pb.ExportAnalyticsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Export generated successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Data:     csv,
		Filename: filename,
	}, nil
}

func (s *ScoringService) analyticsToProto(a *models.UserAnalytics) *pb.UserAnalytics {
	proto := &pb.UserAnalytics{
		UserId:             uint32(a.UserID),
		TotalPredictions:   uint32(a.TotalPredictions),
		CorrectPredictions: uint32(a.CorrectPredictions),
		OverallAccuracy:    a.OverallAccuracy,
		TotalPoints:        a.TotalPoints,
		TimeRange:          a.TimeRange,
	}

	for _, sp := range a.BySport {
		proto.BySport = append(proto.BySport, &pb.SportAccuracy{
			SportType:          sp.SportType,
			TotalPredictions:   uint32(sp.TotalPredictions),
			CorrectPredictions: uint32(sp.CorrectPredictions),
			AccuracyPercentage: sp.AccuracyPercentage,
			TotalPoints:        sp.TotalPoints,
		})
	}

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

	for _, t := range a.ByType {
		proto.ByType = append(proto.ByType, &pb.PredictionTypeAccuracy{
			PredictionType:     t.PredictionType,
			TotalPredictions:   uint32(t.TotalPredictions),
			CorrectPredictions: uint32(t.CorrectPredictions),
			AccuracyPercentage: t.AccuracyPercentage,
			AveragePoints:      t.AveragePoints,
		})
	}

	for _, tr := range a.Trends {
		proto.Trends = append(proto.Trends, &pb.AccuracyTrend{
			Period:             tr.Period,
			TotalPredictions:   uint32(tr.TotalPredictions),
			CorrectPredictions: uint32(tr.CorrectPredictions),
			AccuracyPercentage: tr.AccuracyPercentage,
			TotalPoints:        tr.TotalPoints,
		})
	}

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

func (s *ScoringService) generateCSV(a *pb.UserAnalytics) string {
	var b strings.Builder

	b.WriteString("User Analytics Report\n")
	b.WriteString(fmt.Sprintf("User ID,%d\n", a.UserId))
	b.WriteString(fmt.Sprintf("Time Range,%s\n\n", a.TimeRange))

	b.WriteString("Overall Statistics\n")
	b.WriteString("Metric,Value\n")
	b.WriteString(fmt.Sprintf("Total Predictions,%d\n", a.TotalPredictions))
	b.WriteString(fmt.Sprintf("Correct Predictions,%d\n", a.CorrectPredictions))
	b.WriteString(fmt.Sprintf("Overall Accuracy,%.2f%%\n", a.OverallAccuracy))
	b.WriteString(fmt.Sprintf("Total Points,%.2f\n\n", a.TotalPoints))

	if len(a.BySport) > 0 {
		b.WriteString("Performance by Sport\n")
		b.WriteString("Sport,Total,Correct,Accuracy,Points\n")
		for _, sp := range a.BySport {
			b.WriteString(fmt.Sprintf("%s,%d,%d,%.2f%%,%.2f\n",
				sp.SportType, sp.TotalPredictions, sp.CorrectPredictions, sp.AccuracyPercentage, sp.TotalPoints))
		}
		b.WriteString("\n")
	}

	if len(a.ByType) > 0 {
		b.WriteString("Performance by Prediction Type\n")
		b.WriteString("Type,Total,Correct,Accuracy,Avg Points\n")
		for _, t := range a.ByType {
			b.WriteString(fmt.Sprintf("%s,%d,%d,%.2f%%,%.2f\n",
				t.PredictionType, t.TotalPredictions, t.CorrectPredictions, t.AccuracyPercentage, t.AveragePoints))
		}
		b.WriteString("\n")
	}

	if len(a.Trends) > 0 {
		b.WriteString("Accuracy Trends\n")
		b.WriteString("Period,Total,Correct,Accuracy,Points\n")
		for _, tr := range a.Trends {
			b.WriteString(fmt.Sprintf("%s,%d,%d,%.2f%%,%.2f\n",
				tr.Period, tr.TotalPredictions, tr.CorrectPredictions, tr.AccuracyPercentage, tr.TotalPoints))
		}
	}

	return b.String()
}

// calculatePropsPoints calculates points for props predictions
func (s *ScoringService) calculatePropsPoints(prediction PredictionData, result ResultData, details map[string]interface{}) (float64, map[string]interface{}) {
	if len(prediction.Props) == 0 {
		details["error"] = "No props predictions found"
		return 0, details
	}

	var totalPoints float64
	propResults := make([]map[string]interface{}, 0)

	for _, prop := range prediction.Props {
		propResult := map[string]interface{}{
			"prop_slug": prop.PropSlug,
			"selection": prop.Selection,
			"line":      prop.Line,
		}

		correct := s.evaluateProp(prop, result)
		propResult["correct"] = correct

		if correct {
			points := prop.PointsValue
			if points == 0 {
				points = 2
			}
			totalPoints += points
			propResult["points"] = points
		} else {
			propResult["points"] = float64(0)
		}

		propResults = append(propResults, propResult)
	}

	details["props_results"] = propResults
	details["total_props"] = len(prediction.Props)
	details["correct_props"] = s.countCorrectProps(propResults)

	return totalPoints, details
}

func (s *ScoringService) evaluateProp(prop PropPrediction, result ResultData) bool {
	switch prop.PropSlug {
	case "total-goals-ou":
		totalGoals := float64(result.TotalGoals)
		if prop.Selection == "over" {
			return totalGoals > prop.Line
		}
		return totalGoals < prop.Line

	case "total-corners-ou":
		var corners float64
		switch v := result.Stats["corners"].(type) {
		case float64:
			corners = v
		case int:
			corners = float64(v)
		default:
			return false
		}
		if prop.Selection == "over" {
			return corners > prop.Line
		}
		return corners < prop.Line

	case "btts":
		btts := result.HomeScore > 0 && result.AwayScore > 0
		if prop.Selection == "yes" {
			return btts
		}
		return !btts

	case "first-to-score":
		if firstScorer, ok := result.Stats["first_to_score"].(string); ok {
			return firstScorer == prop.Selection
		}
		return false

	case "total-cards-ou":
		var cards float64
		switch v := result.Stats["cards"].(type) {
		case float64:
			cards = v
		case int:
			cards = float64(v)
		default:
			return false
		}
		if prop.Selection == "over" {
			return cards > prop.Line
		}
		return cards < prop.Line

	default:
		log.Printf("[WARN] Unknown prop slug: %s", prop.PropSlug)
		return false
	}
}

func (s *ScoringService) countCorrectProps(results []map[string]interface{}) int {
	count := 0
	for _, r := range results {
		if correct, ok := r["correct"].(bool); ok && correct {
			count++
		}
	}
	return count
}

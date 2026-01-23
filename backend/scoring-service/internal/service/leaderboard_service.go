package service

import (
	"context"
	"log"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
	"github.com/sports-prediction-contests/scoring-service/internal/repository"
	"github.com/sports-prediction-contests/shared/auth"
	pb "github.com/sports-prediction-contests/shared/proto/scoring"
	"github.com/sports-prediction-contests/shared/proto/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// LeaderboardService implements leaderboard-related gRPC methods
type LeaderboardService struct {
	pb.UnimplementedScoringServiceServer
	leaderboardRepo repository.LeaderboardRepositoryInterface
	scoreRepo       repository.ScoreRepositoryInterface
	streakRepo      repository.StreakRepositoryInterface
}

// NewLeaderboardService creates a new LeaderboardService instance
func NewLeaderboardService(leaderboardRepo repository.LeaderboardRepositoryInterface, scoreRepo repository.ScoreRepositoryInterface, streakRepo repository.StreakRepositoryInterface) *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: leaderboardRepo,
		scoreRepo:       scoreRepo,
		streakRepo:      streakRepo,
	}
}

// GetLeaderboard retrieves the leaderboard for a contest
func (s *LeaderboardService) GetLeaderboard(ctx context.Context, req *pb.GetLeaderboardRequest) (*pb.GetLeaderboardResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	// Get leaderboard entries
	leaderboards, err := s.leaderboardRepo.GetContestLeaderboard(ctx, uint(uint(req.ContestId)), limit)
	if err != nil {
		log.Printf("[ERROR] Failed to get leaderboard: %v", err)
		return &pb.GetLeaderboardResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve leaderboard",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Convert to proto with streak data
	entries := make([]*pb.LeaderboardEntry, len(leaderboards))
	
	// Batch fetch all streaks in one query to avoid N+1
	userIDs := make([]uint, len(leaderboards))
	for i, lb := range leaderboards {
		userIDs[i] = lb.UserID
	}
	streaks, _ := s.streakRepo.GetByContestAndUsers(ctx, uint(uint(req.ContestId)), userIDs)
	streakMap := make(map[uint]*models.UserStreak)
	for _, streak := range streaks {
		streakMap[streak.UserID] = streak
	}
	
	for i, lb := range leaderboards {
		entry := &pb.LeaderboardEntry{
			UserId:      uint32(lb.UserID),
			UserName:    "", // TODO: Fetch user name from user service
			TotalPoints: lb.TotalPoints,
			Rank:        uint32(lb.Rank),
			UpdatedAt:   timestamppb.New(lb.UpdatedAt),
		}
		// Use pre-fetched streak data
		if streak, ok := streakMap[lb.UserID]; ok {
			entry.CurrentStreak = uint32(streak.CurrentStreak)
			entry.MaxStreak = uint32(streak.MaxStreak)
			entry.Multiplier = streak.GetMultiplier()
		}
		entries[i] = entry
	}

	leaderboard := &pb.Leaderboard{
		ContestId: req.ContestId,
		Entries:   entries,
		UpdatedAt: timestamppb.Now(),
	}

	return &pb.GetLeaderboardResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Leaderboard retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Leaderboard: leaderboard,
	}, nil
}

// GetUserRank retrieves a user's rank in a contest
func (s *LeaderboardService) GetUserRank(ctx context.Context, req *pb.GetUserRankRequest) (*pb.GetUserRankResponse, error) {
	// Get user's leaderboard entry
	leaderboard, err := s.leaderboardRepo.GetByContestAndUser(ctx, uint(uint(req.ContestId)), uint(req.UserId))
	if err != nil {
		log.Printf("[ERROR] Failed to get user rank: %v", err)
		return &pb.GetUserRankResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "User not found in leaderboard",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetUserRankResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "User rank retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Rank:        uint32(leaderboard.Rank),
		TotalPoints: leaderboard.TotalPoints,
	}, nil
}

// UpdateLeaderboard recalculates and updates the leaderboard for a contest
func (s *LeaderboardService) UpdateLeaderboard(ctx context.Context, req *pb.UpdateLeaderboardRequest) (*pb.UpdateLeaderboardResponse, error) {
	// Extract user ID from JWT token for authorization
	_, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.UpdateLeaderboardResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Recalculate all user scores for the contest
	if err := s.recalculateContestScores(ctx, uint(uint(req.ContestId))); err != nil {
		log.Printf("[ERROR] Failed to recalculate contest scores: %v", err)
		return &pb.UpdateLeaderboardResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to recalculate scores",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update rankings
	if err := s.leaderboardRepo.UpdateRankings(ctx, uint(req.ContestId)); err != nil {
		log.Printf("[ERROR] Failed to update rankings: %v", err)
		return &pb.UpdateLeaderboardResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to update rankings",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get updated leaderboard
	leaderboards, err := s.leaderboardRepo.GetContestLeaderboard(ctx, uint(req.ContestId), 50)
	if err != nil {
		log.Printf("[ERROR] Failed to get updated leaderboard: %v", err)
		return &pb.UpdateLeaderboardResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve updated leaderboard",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Convert to proto
	entries := make([]*pb.LeaderboardEntry, len(leaderboards))
	for i, lb := range leaderboards {
		entries[i] = &pb.LeaderboardEntry{
			UserId:      uint32(lb.UserID),
			UserName:    "", // TODO: Fetch user name from user service
			TotalPoints: lb.TotalPoints,
			Rank:        uint32(lb.Rank),
			UpdatedAt:   timestamppb.New(lb.UpdatedAt),
		}
	}

	leaderboard := &pb.Leaderboard{
		ContestId: req.ContestId,
		Entries:   entries,
		UpdatedAt: timestamppb.Now(),
	}

	return &pb.UpdateLeaderboardResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Leaderboard updated successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Leaderboard: leaderboard,
	}, nil
}

// recalculateContestScores recalculates total scores for all users in a contest
func (s *LeaderboardService) recalculateContestScores(ctx context.Context, contestID uint) error {
	// Get all scores for the contest
	scores, err := s.scoreRepo.ListByContest(ctx, contestID)
	if err != nil {
		return err
	}

	// Group scores by user
	userScores := make(map[uint]float64)
	for _, score := range scores {
		userScores[score.UserID] += score.Points
	}

	// Update leaderboard entries
	for userID, totalPoints := range userScores {
		if err := s.leaderboardRepo.UpsertUserScore(ctx, contestID, userID, totalPoints); err != nil {
			log.Printf("[ERROR] Failed to update user score for user %d: %v", userID, err)
			return err
		}
	}

	return nil
}

// BatchUpdateUserScores updates multiple user scores efficiently
func (s *LeaderboardService) BatchUpdateUserScores(ctx context.Context, contestID uint, userScores map[uint]float64) error {
	for userID, totalPoints := range userScores {
		if err := s.leaderboardRepo.UpsertUserScore(ctx, contestID, userID, totalPoints); err != nil {
			log.Printf("[ERROR] Failed to batch update user score for user %d: %v", userID, err)
			return err
		}
	}

	// Recalculate rankings after batch update
	return s.leaderboardRepo.UpdateRankings(ctx, contestID)
}

// GetLeaderboardSize returns the number of participants in a contest leaderboard
func (s *LeaderboardService) GetLeaderboardSize(ctx context.Context, contestID uint) (int64, error) {
	_, _, err := s.leaderboardRepo.ListByContest(ctx, contestID, 1, 0)
	if err != nil {
		return 0, err
	}

	// This is a simple count - in production, you might want a dedicated count method
	_, total, err := s.leaderboardRepo.ListByContest(ctx, contestID, 1, 0)
	return total, err
}

// RefreshLeaderboardCache refreshes the cache for a contest leaderboard
func (s *LeaderboardService) RefreshLeaderboardCache(ctx context.Context, contestID uint) error {
	// Get fresh data from database
	leaderboards, err := s.leaderboardRepo.GetContestLeaderboard(ctx, contestID, 100)
	if err != nil {
		return err
	}

	// Cache will be updated automatically by the repository layer
	log.Printf("[INFO] Refreshed leaderboard cache for contest %d with %d entries", contestID, len(leaderboards))
	return nil
}

// GetUserStreak retrieves a user's streak information for a contest
func (s *LeaderboardService) GetUserStreak(ctx context.Context, req *pb.GetUserStreakRequest) (*pb.GetUserStreakResponse, error) {
	streak, err := s.streakRepo.GetByContestAndUser(ctx, uint(uint(req.ContestId)), uint(req.UserId))
	if err != nil {
		log.Printf("[ERROR] Failed to get user streak: %v", err)
		return &pb.GetUserStreakResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Streak not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetUserStreakResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Streak retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		CurrentStreak: uint32(streak.CurrentStreak),
		MaxStreak:     uint32(streak.MaxStreak),
		Multiplier:    streak.GetMultiplier(),
	}, nil
}

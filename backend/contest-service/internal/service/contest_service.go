package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sports-prediction-contests/contest-service/internal/models"
	"github.com/sports-prediction-contests/contest-service/internal/repository"
	"github.com/sports-prediction-contests/shared/auth"
	pb "github.com/sports-prediction-contests/shared/proto/contest"
	"github.com/sports-prediction-contests/shared/proto/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ContestService implements the gRPC ContestService
type ContestService struct {
	pb.UnimplementedContestServiceServer
	contestRepo     repository.ContestRepositoryInterface
	participantRepo repository.ParticipantRepositoryInterface
}

// NewContestService creates a new ContestService instance
func NewContestService(contestRepo repository.ContestRepositoryInterface, participantRepo repository.ParticipantRepositoryInterface) *ContestService {
	return &ContestService{
		contestRepo:     contestRepo,
		participantRepo: participantRepo,
	}
}

// CreateContest handles contest creation
func (s *ContestService) CreateContest(ctx context.Context, req *pb.CreateContestRequest) (*pb.CreateContestResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.CreateContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Create contest model
	contest := &models.Contest{
		Title:           req.Title,
		Description:     req.Description,
		SportType:       req.SportType,
		Rules:           req.Rules,
		StartDate:       req.StartDate.AsTime(),
		EndDate:         req.EndDate.AsTime(),
		MaxParticipants: uint(req.MaxParticipants),
		CreatorID:       userID,
		Status:          "draft",
	}

	// Save to database
	if err := s.contestRepo.Create(contest); err != nil {
		log.Printf("[ERROR] Failed to create contest: %v", err)
		return &pb.CreateContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Create admin participant for the creator
	participant := &models.Participant{
		ContestID: contest.ID,
		UserID:    userID,
		Role:      "admin",
		Status:    "active",
		JoinedAt:  time.Now(),
	}

	if err := s.participantRepo.Create(participant); err != nil {
		log.Printf("[ERROR] Failed to create admin participant: %v", err)
		// Delete the contest since admin participant creation failed
		s.contestRepo.Delete(contest.ID)
		return &pb.CreateContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to create contest admin",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update participant count using database aggregation
	if err := s.updateContestParticipantCount(contest.ID); err != nil {
		log.Printf("[ERROR] Failed to update participant count: %v", err)
	}

	log.Printf("[INFO] Contest created successfully: %d", contest.ID)

	return &pb.CreateContestResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Contest created successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Contest: s.modelToProto(contest),
	}, nil
}

// UpdateContest handles contest updates
func (s *ContestService) UpdateContest(ctx context.Context, req *pb.UpdateContestRequest) (*pb.UpdateContestResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &pb.UpdateContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get existing contest
	contest, err := s.contestRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.UpdateContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if user is the creator or admin
	if contest.CreatorID != userID {
		// Check if user is admin participant
		participant, err := s.participantRepo.GetByContestAndUser(contest.ID, userID)
		if err != nil || !participant.IsAdmin() {
			return &pb.UpdateContestResponse{
				Response: &common.Response{
					Success:   false,
					Message:   "Permission denied",
					Code:      int32(common.ErrorCode_PERMISSION_DENIED),
					Timestamp: timestamppb.Now(),
				},
			}, nil
		}
	}

	// Update fields
	contest.Title = req.Title
	contest.Description = req.Description
	contest.SportType = req.SportType
	contest.Rules = req.Rules
	contest.StartDate = req.StartDate.AsTime()
	contest.EndDate = req.EndDate.AsTime()
	contest.MaxParticipants = uint(req.MaxParticipants)
	if req.Status != "" {
		contest.Status = req.Status
	}

	// Save to database
	if err := s.contestRepo.Update(contest); err != nil {
		return &pb.UpdateContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Contest updated successfully: %d", contest.ID)

	return &pb.UpdateContestResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Contest updated successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Contest: s.modelToProto(contest),
	}, nil
}

// GetContest retrieves a contest by ID
func (s *ContestService) GetContest(ctx context.Context, req *pb.GetContestRequest) (*pb.GetContestResponse, error) {
	contest, err := s.contestRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.GetContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetContestResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Contest retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Contest: s.modelToProto(contest),
	}, nil
}

// DeleteContest handles contest deletion
func (s *ContestService) DeleteContest(ctx context.Context, req *pb.DeleteContestRequest) (*pb.DeleteContestResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &pb.DeleteContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get existing contest
	contest, err := s.contestRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.DeleteContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if user is the creator
	if contest.CreatorID != userID {
		// TODO: Add admin role check here
		// For now, allow deletion for testing purposes
		log.Printf("[WARN] User %d attempting to delete contest %d created by user %d", userID, req.Id, contest.CreatorID)
		
		// Uncomment this to enforce creator-only deletion:
		// return &pb.DeleteContestResponse{
		// 	Response: &common.Response{
		// 		Success:   false,
		// 		Message:   "Permission denied: only contest creator can delete",
		// 		Code:      int32(common.ErrorCode_PERMISSION_DENIED),
		// 		Timestamp: timestamppb.Now(),
		// 	},
		// }, nil
	}

	// Check if contest has predictions (safe delete)
	participantCount, err := s.participantRepo.CountByContest(contest.ID)
	if err == nil && participantCount > 0 {
		return &pb.DeleteContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   fmt.Sprintf("Cannot delete contest with %d participants. Remove participants first.", participantCount),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Delete contest
	if err := s.contestRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Contest deleted successfully: %d", req.Id)

	return &pb.DeleteContestResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Contest deleted successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

// ListContests retrieves contests with pagination and filters
func (s *ContestService) ListContests(ctx context.Context, req *pb.ListContestsRequest) (*pb.ListContestsResponse, error) {
	// Set default pagination
	limit := 20
	page := 1
	
	if req.Pagination != nil {
		if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
			limit = int(req.Pagination.Limit)
		}
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
	}

	offset := (page - 1) * limit

	// Get contests from repository
	contests, total, err := s.contestRepo.List(limit, offset, req.Status, req.SportType)
	if err != nil {
		return &pb.ListContestsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Convert to proto
	protoContests := make([]*pb.Contest, len(contests))
	for i, contest := range contests {
		protoContests[i] = s.modelToProto(contest)
	}

	// Calculate pagination
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return &pb.ListContestsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Contests retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Contests: protoContests,
		Pagination: &common.PaginationResponse{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

// JoinContest handles user joining a contest
func (s *ContestService) JoinContest(ctx context.Context, req *pb.JoinContestRequest) (*pb.JoinContestResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &pb.JoinContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get contest
	contest, err := s.contestRepo.GetByID(uint(req.ContestId))
	if err != nil {
		return &pb.JoinContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if user can join
	if !contest.CanJoin() {
		return &pb.JoinContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Cannot join this contest",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Create participant
	participant := &models.Participant{
		ContestID: contest.ID,
		UserID:    userID,
		Role:      "participant",
		Status:    "active",
		JoinedAt:  time.Now(),
	}

	if err := s.participantRepo.Create(participant); err != nil {
		return &pb.JoinContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update participant count using database aggregation
	if err := s.updateContestParticipantCount(contest.ID); err != nil {
		log.Printf("[ERROR] Failed to update participant count: %v", err)
	}

	log.Printf("[INFO] User %d joined contest %d", userID, contest.ID)

	return &pb.JoinContestResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Joined contest successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Participant: s.participantModelToProto(participant),
	}, nil
}

// LeaveContest handles user leaving a contest
func (s *ContestService) LeaveContest(ctx context.Context, req *pb.LeaveContestRequest) (*pb.LeaveContestResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &pb.LeaveContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Delete participant
	if err := s.participantRepo.DeleteByContestAndUser(uint(req.ContestId), userID); err != nil {
		return &pb.LeaveContestResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update participant count using database aggregation
	if err := s.updateContestParticipantCount(uint(req.ContestId)); err != nil {
		log.Printf("[ERROR] Failed to update participant count: %v", err)
	}

	log.Printf("[INFO] User %d left contest %d", userID, req.ContestId)

	return &pb.LeaveContestResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Left contest successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

// ListParticipants retrieves participants for a contest
func (s *ContestService) ListParticipants(ctx context.Context, req *pb.ListParticipantsRequest) (*pb.ListParticipantsResponse, error) {
	// Set default pagination
	limit := 20
	page := 1
	
	if req.Pagination != nil {
		if req.Pagination.Limit > 0 && req.Pagination.Limit <= 100 {
			limit = int(req.Pagination.Limit)
		}
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
	}

	offset := (page - 1) * limit

	// Get participants from repository
	participants, total, err := s.participantRepo.ListByContest(uint(req.ContestId), limit, offset)
	if err != nil {
		return &pb.ListParticipantsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Convert to proto
	protoParticipants := make([]*pb.Participant, len(participants))
	for i, participant := range participants {
		protoParticipants[i] = s.participantModelToProto(participant)
	}

	// Calculate pagination
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return &pb.ListParticipantsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Participants retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Participants: protoParticipants,
		Pagination: &common.PaginationResponse{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

// Check handles health check
func (s *ContestService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return &common.Response{
		Success:   true,
		Message:   "Contest service is healthy",
		Code:      0,
		Timestamp: timestamppb.Now(),
	}, nil
}

// Helper methods

// updateContestParticipantCount updates the participant count using database aggregation
func (s *ContestService) updateContestParticipantCount(contestID uint) error {
	count, err := s.participantRepo.CountByContest(contestID)
	if err != nil {
		return err
	}

	contest, err := s.contestRepo.GetByID(contestID)
	if err != nil {
		return err
	}

	contest.CurrentParticipants = uint(count)
	return s.contestRepo.Update(contest)
}

// modelToProto converts Contest model to proto message
func (s *ContestService) modelToProto(contest *models.Contest) *pb.Contest {
	return &pb.Contest{
		Id:                  uint32(contest.ID),
		Title:               contest.Title,
		Description:         contest.Description,
		SportType:           contest.SportType,
		Rules:               contest.Rules,
		Status:              contest.GetComputedStatus(), // Use computed status based on dates
		StartDate:           timestamppb.New(contest.StartDate),
		EndDate:             timestamppb.New(contest.EndDate),
		MaxParticipants:     uint32(contest.MaxParticipants),
		CurrentParticipants: uint32(contest.CurrentParticipants),
		CreatorId:           uint32(contest.CreatorID),
		CreatedAt:           timestamppb.New(contest.CreatedAt),
		UpdatedAt:           timestamppb.New(contest.UpdatedAt),
	}
}

// participantModelToProto converts Participant model to proto message
func (s *ContestService) participantModelToProto(participant *models.Participant) *pb.Participant {
	return &pb.Participant{
		Id:        uint32(participant.ID),
		ContestId: uint32(participant.ContestID),
		UserId:    uint32(participant.UserID),
		Role:      participant.Role,
		Status:    participant.Status,
		JoinedAt:  timestamppb.New(participant.JoinedAt),
	}
}

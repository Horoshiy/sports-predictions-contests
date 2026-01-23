package service

import (
	"context"
	"log"
	"time"

	"github.com/sports-prediction-contests/challenge-service/internal/models"
	"github.com/sports-prediction-contests/challenge-service/internal/repository"
	"github.com/sports-prediction-contests/shared/auth"
	pb "github.com/sports-prediction-contests/shared/proto/challenge"
	"github.com/sports-prediction-contests/shared/proto/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ChallengeService implements the gRPC ChallengeService
type ChallengeService struct {
	pb.UnimplementedChallengeServiceServer
	challengeRepo     repository.ChallengeRepositoryInterface
	participantRepo   repository.ChallengeParticipantRepositoryInterface
}

// NewChallengeService creates a new ChallengeService instance
func NewChallengeService(challengeRepo repository.ChallengeRepositoryInterface, participantRepo repository.ChallengeParticipantRepositoryInterface) *ChallengeService {
	return &ChallengeService{
		challengeRepo:   challengeRepo,
		participantRepo: participantRepo,
	}
}

// CreateChallenge handles challenge creation
func (s *ChallengeService) CreateChallenge(ctx context.Context, req *pb.CreateChallengeRequest) (*pb.CreateChallengeResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.CreateChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate that user is not challenging themselves
	if userID == uint(req.OpponentId) {
		log.Printf("[ERROR] User %d trying to challenge themselves", userID)
		return &pb.CreateChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Cannot challenge yourself",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate opponent ID is within valid range
	if req.OpponentId == 0 {
		log.Printf("[ERROR] Invalid opponent ID: %d", req.OpponentId)
		return &pb.CreateChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Invalid opponent ID",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate event ID is within valid range
	if req.EventId == 0 {
		log.Printf("[ERROR] Invalid event ID: %d", uint(req.EventId))
		return &pb.CreateChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Invalid event ID",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Create challenge model
	challenge := &models.Challenge{
		ChallengerID: userID,
		OpponentID:   uint(req.OpponentId),
		EventID:      uint(req.EventId),
		Message:      req.Message,
		Status:       "pending",
		ExpiresAt:    time.Now().UTC().Add(24 * time.Hour),
	}

	// Create challenge participants
	challengerParticipant := &models.ChallengeParticipant{
		UserID:   userID,
		Role:     "challenger",
		Status:   "active",
		JoinedAt: time.Now().UTC(),
	}

	opponentParticipant := &models.ChallengeParticipant{
		UserID:   uint(req.OpponentId),
		Role:     "opponent",
		Status:   "active",
		JoinedAt: time.Now().UTC(),
	}

	participants := []*models.ChallengeParticipant{challengerParticipant, opponentParticipant}

	// Save challenge and participants atomically
	if err := s.challengeRepo.CreateWithParticipants(challenge, participants); err != nil {
		log.Printf("[ERROR] Failed to create challenge with participants: %v", err)
		return &pb.CreateChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to create challenge",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Challenge created successfully: ID=%d, Challenger=%d, Opponent=%d", challenge.ID, userID, req.OpponentId)

	return &pb.CreateChallengeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Challenge created successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Challenge: s.challengeModelToProto(challenge),
	}, nil
}

// AcceptChallenge handles challenge acceptance
func (s *ChallengeService) AcceptChallenge(ctx context.Context, req *pb.AcceptChallengeRequest) (*pb.AcceptChallengeResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.AcceptChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get challenge from database
	challenge, err := s.challengeRepo.GetByID(uint(req.Id))
	if !ok {
		log.Printf("[ERROR] Failed to get challenge %d: %v", req.Id, err)
		return &pb.AcceptChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Challenge not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate that user is the opponent
	if challenge.OpponentID != userID {
		log.Printf("[ERROR] User %d is not the opponent for challenge %d", userID, uint(req.Id))
		return &pb.AcceptChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "You are not authorized to accept this challenge",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if challenge can be accepted
	if !challenge.CanAccept() {
		log.Printf("[ERROR] Challenge %d cannot be accepted (status: %s, expired: %v)", req.Id, challenge.Status, challenge.IsExpired())
		return &pb.AcceptChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Challenge cannot be accepted",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Accept the challenge
	challenge.Accept()

	// Update challenge in database
	if err := s.challengeRepo.Update(challenge); err != nil {
		log.Printf("[ERROR] Failed to update challenge %d: %v", req.Id, err)
		return &pb.AcceptChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to accept challenge",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Challenge accepted successfully: ID=%d, User=%d", req.Id, userID)

	return &pb.AcceptChallengeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Challenge accepted successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Challenge: s.challengeModelToProto(challenge),
	}, nil
}

// DeclineChallenge handles challenge decline
func (s *ChallengeService) DeclineChallenge(ctx context.Context, req *pb.DeclineChallengeRequest) (*pb.DeclineChallengeResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.DeclineChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get challenge from database
	challenge, err := s.challengeRepo.GetByID(uint(req.Id))
	if !ok {
		log.Printf("[ERROR] Failed to get challenge %d: %v", req.Id, err)
		return &pb.DeclineChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Challenge not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate that user is the opponent
	if challenge.OpponentID != userID {
		log.Printf("[ERROR] User %d is not the opponent for challenge %d", userID, uint(req.Id))
		return &pb.DeclineChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "You are not authorized to decline this challenge",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if challenge is still pending
	if challenge.Status != "pending" {
		log.Printf("[ERROR] Challenge %d is not pending (status: %s)", req.Id, challenge.Status)
		return &pb.DeclineChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Challenge is not pending",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Decline the challenge
	challenge.Status = "declined"

	// Update challenge in database
	if err := s.challengeRepo.Update(challenge); err != nil {
		log.Printf("[ERROR] Failed to update challenge %d: %v", req.Id, err)
		return &pb.DeclineChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to decline challenge",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Challenge declined successfully: ID=%d, User=%d", req.Id, userID)

	return &pb.DeclineChallengeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Challenge declined successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

// WithdrawChallenge handles challenge withdrawal
func (s *ChallengeService) WithdrawChallenge(ctx context.Context, req *pb.WithdrawChallengeRequest) (*pb.WithdrawChallengeResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.WithdrawChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get challenge from database
	challenge, err := s.challengeRepo.GetByID(uint(req.Id))
	if !ok {
		log.Printf("[ERROR] Failed to get challenge %d: %v", req.Id, err)
		return &pb.WithdrawChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Challenge not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate that user is the challenger
	if challenge.ChallengerID != userID {
		log.Printf("[ERROR] User %d is not the challenger for challenge %d", userID, uint(req.Id))
		return &pb.WithdrawChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "You are not authorized to withdraw this challenge",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if challenge can be withdrawn (only pending challenges)
	if challenge.Status != "pending" {
		log.Printf("[ERROR] Challenge %d cannot be withdrawn (status: %s)", req.Id, challenge.Status)
		return &pb.WithdrawChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Challenge cannot be withdrawn",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Delete the challenge
	if err := s.challengeRepo.Delete(uint(req.Id)); err != nil {
		log.Printf("[ERROR] Failed to delete challenge %d: %v", req.Id, err)
		return &pb.WithdrawChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to withdraw challenge",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Challenge withdrawn successfully: ID=%d, User=%d", req.Id, userID)

	return &pb.WithdrawChallengeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Challenge withdrawn successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

// GetChallenge retrieves a challenge by ID
func (s *ChallengeService) GetChallenge(ctx context.Context, req *pb.GetChallengeRequest) (*pb.GetChallengeResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.GetChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Get challenge from database
	challenge, err := s.challengeRepo.GetByID(uint(req.Id))
	if !ok {
		log.Printf("[ERROR] Failed to get challenge %d: %v", req.Id, err)
		return &pb.GetChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Challenge not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate that user is involved in the challenge
	if challenge.ChallengerID != userID && challenge.OpponentID != userID {
		log.Printf("[ERROR] User %d is not involved in challenge %d", userID, uint(req.Id))
		return &pb.GetChallengeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "You are not authorized to view this challenge",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetChallengeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Challenge retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Challenge: s.challengeModelToProto(challenge),
	}, nil
}

// ListUserChallenges retrieves challenges for a user
func (s *ChallengeService) ListUserChallenges(ctx context.Context, req *pb.ListUserChallengesRequest) (*pb.ListUserChallengesResponse, error) {
	// Extract user ID from JWT token
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		log.Printf("[ERROR] Failed to get user ID from context")
		return &pb.ListUserChallengesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Authentication required",
				Code:      int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Validate that user can only list their own challenges
	if uint32(userID) != req.UserId {
		log.Printf("[ERROR] User %d trying to list challenges for user %d", userID, req.UserId)
		return &pb.ListUserChallengesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "You can only list your own challenges",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Set default pagination
	limit := int(req.Pagination.Limit)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	// Validate page number and calculate offset safely
	page := req.Pagination.Page
	if page < 1 {
		page = 1
	}
	if page > 1000000 { // Reasonable upper limit to prevent overflow
		page = 1000000
	}
	
	offset := int(page-1) * limit
	if offset < 0 {
		offset = 0
	}

	// Get challenges from database
	challenges, total, err := s.challengeRepo.ListByUser(userID, req.Status, limit, offset)
	if !ok {
		log.Printf("[ERROR] Failed to list challenges for user %d: %v", userID, err)
		return &pb.ListUserChallengesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve challenges",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Convert to proto
	protoChallenges := make([]*pb.Challenge, len(challenges))
	for i, challenge := range challenges {
		protoChallenges[i] = s.challengeModelToProto(challenge)
	}

	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return &pb.ListUserChallengesResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Challenges retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Challenges: protoChallenges,
		Pagination: &common.PaginationResponse{
			Page:       req.Pagination.Page,
			Limit:      int32(limit),
			Total: int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

// ListOpenChallenges retrieves open challenges for an event
func (s *ChallengeService) ListOpenChallenges(ctx context.Context, req *pb.ListOpenChallengesRequest) (*pb.ListOpenChallengesResponse, error) {
	// Set default pagination
	limit := int(req.Pagination.Limit)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	// Validate page number and calculate offset safely
	page := req.Pagination.Page
	if page < 1 {
		page = 1
	}
	if page > 1000000 { // Reasonable upper limit to prevent overflow
		page = 1000000
	}
	
	offset := int(page-1) * limit
	if offset < 0 {
		offset = 0
	}

	// Get challenges from database
	challenges, total, err := s.challengeRepo.ListByEvent(uint(req.EventId), limit, offset)
	if err != nil {
		log.Printf("[ERROR] Failed to list challenges for event %d: %v", req.EventId, err)
		return &pb.ListOpenChallengesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve challenges",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Convert to proto
	protoChallenges := make([]*pb.Challenge, len(challenges))
	for i, challenge := range challenges {
		protoChallenges[i] = s.challengeModelToProto(challenge)
	}

	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return &pb.ListOpenChallengesResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Challenges retrieved successfully",
			Code:      int32(0),
			Timestamp: timestamppb.Now(),
		},
		Challenges: protoChallenges,
		Pagination: &common.PaginationResponse{
			Page:       req.Pagination.Page,
			Limit:      int32(limit),
			Total: int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

// Check handles health check
func (s *ChallengeService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return &common.Response{
		Success:   true,
		Message:   "Challenge service is healthy",
		Code:      int32(0),
		Timestamp: timestamppb.Now(),
	}, nil
}

// challengeModelToProto converts a challenge model to proto
func (s *ChallengeService) challengeModelToProto(challenge *models.Challenge) *pb.Challenge {
	if challenge == nil {
		return nil
	}
	
	protoChallenge := &pb.Challenge{
		Id:             uint32(challenge.ID),
		ChallengerId:   uint32(challenge.ChallengerID),
		OpponentId:     uint32(challenge.OpponentID),
		EventId:        uint32(challenge.EventID),
		Message:        challenge.Message,
		Status:         challenge.Status,
		ExpiresAt:      timestamppb.New(challenge.ExpiresAt),
		ChallengerScore: challenge.ChallengerScore,
		OpponentScore:  challenge.OpponentScore,
		CreatedAt:      timestamppb.New(challenge.CreatedAt),
		UpdatedAt:      timestamppb.New(challenge.UpdatedAt),
	}

	if challenge.AcceptedAt != nil {
		protoChallenge.AcceptedAt = timestamppb.New(*challenge.AcceptedAt)
	}

	if challenge.CompletedAt != nil {
		protoChallenge.CompletedAt = timestamppb.New(*challenge.CompletedAt)
	}

	if challenge.WinnerID != nil {
		protoChallenge.WinnerId = uint32(*challenge.WinnerID)
	}

	return protoChallenge
}

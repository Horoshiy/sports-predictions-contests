package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/sports-prediction-contests/prediction-service/internal/clients"
	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"github.com/sports-prediction-contests/prediction-service/internal/repository"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/coefficient"
	"github.com/sports-prediction-contests/shared/proto/common"
	pb "github.com/sports-prediction-contests/shared/proto/prediction"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PredictionService implements the gRPC PredictionService
type PredictionService struct {
	pb.UnimplementedPredictionServiceServer
	predictionRepo   repository.PredictionRepositoryInterface
	eventRepo        repository.EventRepositoryInterface
	propTypeRepo     *repository.PropTypeRepository
	riskyEventRepo   *repository.RiskyEventRepository
	relayRepo        repository.RelayRepositoryInterface
	contestClient    *clients.ContestClient
	teamClient       *clients.TeamClient
}

// NewPredictionService creates a new PredictionService instance
func NewPredictionService(
	predictionRepo repository.PredictionRepositoryInterface,
	eventRepo repository.EventRepositoryInterface,
	propTypeRepo *repository.PropTypeRepository,
	riskyEventRepo *repository.RiskyEventRepository,
	relayRepo repository.RelayRepositoryInterface,
	contestClient *clients.ContestClient,
	teamClient *clients.TeamClient,
) *PredictionService {
	return &PredictionService{
		predictionRepo:   predictionRepo,
		eventRepo:        eventRepo,
		propTypeRepo:     propTypeRepo,
		riskyEventRepo:   riskyEventRepo,
		relayRepo:        relayRepo,
		contestClient:    contestClient,
		teamClient:       teamClient,
	}
}

// SubmitPrediction handles prediction submission
func (s *PredictionService) SubmitPrediction(ctx context.Context, req *pb.SubmitPredictionRequest) (*pb.SubmitPredictionResponse, error) {
	// Extract user ID from JWT context
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Validate contest participation (non-blocking - auto-join on first prediction)
	// If validation fails, we still allow the prediction to be saved
	_ = s.contestClient.ValidateContestParticipation(ctx, req.ContestId, uint32(userID))

	// Check if contest is relay type - user can only predict assigned events
	contest, err := s.contestClient.GetContest(ctx, req.ContestId)
	if err == nil && contest != nil {
		contestType := parseContestType(contest.Rules)
		if contestType == "relay" {
			// For relay contests, validate user is assigned to this event
			canPredict, err := s.relayRepo.ValidateUserCanPredict(
				uint(req.ContestId), 0, uint(userID), uint(req.EventId),
			)
			if err != nil {
				return &pb.SubmitPredictionResponse{
					Response: &common.Response{
						Success:   false,
						Message:   "Failed to validate relay assignment",
						Code:      int32(common.ErrorCode_INTERNAL_ERROR),
						Timestamp: timestamppb.Now(),
					},
				}, nil
			}
			if !canPredict {
				return &pb.SubmitPredictionResponse{
					Response: &common.Response{
						Success:   false,
						Message:   "This event is not assigned to you in this relay contest",
						Code:      int32(common.ErrorCode_PERMISSION_DENIED),
						Timestamp: timestamppb.Now(),
					},
				}, nil
			}
		}
	}

	// Validate event exists and can accept predictions
	event, err := s.eventRepo.GetByID(uint(req.EventId))
	if err != nil {
		return &pb.SubmitPredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Event not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	if !event.CanAcceptPredictions() {
		return &pb.SubmitPredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Event cannot accept predictions",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check for existing prediction
	existingPrediction, err := s.predictionRepo.GetByUserContestAndEvent(userID, uint(req.ContestId), uint(req.EventId))
	if err != nil {
		return &pb.SubmitPredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// If prediction exists, update it (allow changing prediction before match starts)
	if existingPrediction != nil {
		existingPrediction.PredictionData = req.PredictionData
		existingPrediction.SubmittedAt = time.Now().UTC()
		
		if err := s.predictionRepo.Update(existingPrediction); err != nil {
			return &pb.SubmitPredictionResponse{
				Response: &common.Response{
					Success:   false,
					Message:   "Failed to update prediction",
					Code:      int32(common.ErrorCode_INTERNAL_ERROR),
					Timestamp: timestamppb.Now(),
				},
			}, nil
		}
		
		return &pb.SubmitPredictionResponse{
			Response: &common.Response{
				Success:   true,
				Message:   "Prediction updated successfully",
				Code:      0,
				Timestamp: timestamppb.Now(),
			},
			Prediction: s.modelToPB(existingPrediction),
		}, nil
	}

	// Create new prediction
	prediction := &models.Prediction{
		ContestID:      uint(req.ContestId),
		UserID:         userID,
		EventID:        uint(req.EventId),
		PredictionData: req.PredictionData,
		Status:         "pending",
		SubmittedAt:    time.Now().UTC(),
	}

	if err := s.predictionRepo.Create(prediction); err != nil {
		// Check if it's a unique constraint violation (duplicate)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return &pb.SubmitPredictionResponse{
				Response: &common.Response{
					Success:   false,
					Message:   "Prediction already exists for this event",
					Code:      int32(common.ErrorCode_ALREADY_EXISTS),
					Timestamp: timestamppb.Now(),
				},
			}, nil
		}
		
		return &pb.SubmitPredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to create prediction",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.SubmitPredictionResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Prediction submitted successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Prediction: s.modelToPB(prediction),
	}, nil
}

// GetPrediction retrieves a prediction by ID
func (s *PredictionService) GetPrediction(ctx context.Context, req *pb.GetPredictionRequest) (*pb.GetPredictionResponse, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	prediction, err := s.predictionRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.GetPredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Ensure user can only access their own predictions
	if prediction.UserID != userID {
		return &pb.GetPredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Access denied",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetPredictionResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Prediction retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Prediction: s.modelToPB(prediction),
	}, nil
}

// GetUserPredictions retrieves user predictions for a contest
func (s *PredictionService) GetUserPredictions(ctx context.Context, req *pb.GetUserPredictionsRequest) (*pb.GetUserPredictionsResponse, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Set default pagination if not provided
	page := int32(1)
	pageSize := int32(10)
	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = req.Pagination.Page
		}
		if req.Pagination.Limit > 0 {
			pageSize = req.Pagination.Limit
		}
	}

	// Calculate offset
	offset := int((page - 1) * pageSize)
	limit := int(pageSize)

	// Get predictions with pagination
	predictions, total, err := s.predictionRepo.List(limit, offset, uint(req.ContestId), userID)
	if err != nil {
		return &pb.GetUserPredictionsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve predictions",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	pbPredictions := make([]*pb.Prediction, len(predictions))
	for i, prediction := range predictions {
		pbPredictions[i] = s.modelToPB(prediction)
	}

	// Calculate total pages
	totalPages := int32((total + int64(pageSize) - 1) / int64(pageSize))

	return &pb.GetUserPredictionsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Predictions retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Predictions: pbPredictions,
		Pagination: &common.PaginationResponse{
			Page:       page,
			Limit:      pageSize,
			Total:      int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

// UpdatePrediction updates an existing prediction
func (s *PredictionService) UpdatePrediction(ctx context.Context, req *pb.UpdatePredictionRequest) (*pb.UpdatePredictionResponse, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	prediction, err := s.predictionRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.UpdatePredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Ensure user can only update their own predictions
	if prediction.UserID != userID {
		return &pb.UpdatePredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Access denied",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if prediction can be updated
	if !prediction.CanUpdate() {
		return &pb.UpdatePredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Prediction cannot be updated",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update prediction data
	prediction.PredictionData = req.PredictionData

	if err := s.predictionRepo.Update(prediction); err != nil {
		return &pb.UpdatePredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.UpdatePredictionResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Prediction updated successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Prediction: s.modelToPB(prediction),
	}, nil
}

// DeletePrediction deletes a prediction
func (s *PredictionService) DeletePrediction(ctx context.Context, req *pb.DeletePredictionRequest) (*pb.DeletePredictionResponse, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	prediction, err := s.predictionRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.DeletePredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Ensure user can only delete their own predictions
	if prediction.UserID != userID {
		return &pb.DeletePredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Access denied",
				Code:      int32(common.ErrorCode_PERMISSION_DENIED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	if err := s.predictionRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeletePredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.DeletePredictionResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Prediction deleted successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

// CreateEvent creates a new sports event
func (s *PredictionService) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	event := &models.Event{
		Title:     req.Title,
		SportType: req.SportType,
		HomeTeam:  req.HomeTeam,
		AwayTeam:  req.AwayTeam,
		EventDate: req.EventDate.AsTime(),
		Status:    "scheduled",
	}

	if err := s.eventRepo.Create(event); err != nil {
		return &pb.CreateEventResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.CreateEventResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Event created successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Event: s.eventModelToPB(event),
	}, nil
}

// GetEvent retrieves an event by ID
func (s *PredictionService) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	event, err := s.eventRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.GetEventResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetEventResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Event retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Event: s.eventModelToPB(event),
	}, nil
}

// ListEvents lists events with optional filters
func (s *PredictionService) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	limit := 10
	offset := 0
	
	if req.Pagination != nil {
		if req.Pagination.Limit > 0 {
			limit = int(req.Pagination.Limit)
		}
		if req.Pagination.Page > 1 {
			offset = int(req.Pagination.Page-1) * limit
		}
	}

	var events []*models.Event
	var total int64
	var err error

	// Filter by contest if contest_id is provided
	if req.ContestId > 0 {
		events, total, err = s.eventRepo.ListByContest(uint(req.ContestId), req.SportType, req.Status)
	} else {
		events, total, err = s.eventRepo.List(limit, offset, req.SportType, req.Status)
	}
	if err != nil {
		return &pb.ListEventsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = s.eventModelToPB(event)
	}

	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	page := int32(1)
	if req.Pagination != nil && req.Pagination.Page > 0 {
		page = req.Pagination.Page
	}

	return &pb.ListEventsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Events retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Events: pbEvents,
		Pagination: &common.PaginationResponse{
			Page:       page,
			Limit:      int32(limit),
			Total:      int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateEvent updates an existing event
func (s *PredictionService) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	event, err := s.eventRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.UpdateEventResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update fields if provided
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.HomeTeam != "" {
		event.HomeTeam = req.HomeTeam
	}
	if req.AwayTeam != "" {
		event.AwayTeam = req.AwayTeam
	}
	if req.EventDate != nil {
		event.EventDate = req.EventDate.AsTime()
	}
	if req.Status != "" {
		event.Status = req.Status
	}
	if req.ResultData != "" {
		event.ResultData = req.ResultData
	}

	if err := s.eventRepo.Update(event); err != nil {
		return &pb.UpdateEventResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.UpdateEventResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Event updated successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Event: s.eventModelToPB(event),
	}, nil
}

// Check implements health check
func (s *PredictionService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return &common.Response{
		Success:   true,
		Message:   "Prediction service is healthy",
		Code:      0,
		Timestamp: timestamppb.Now(),
	}, nil
}

// Helper methods for model conversion
func (s *PredictionService) modelToPB(prediction *models.Prediction) *pb.Prediction {
	if prediction == nil {
		return nil
	}
	
	pbPrediction := &pb.Prediction{
		Id:             uint32(prediction.ID),
		ContestId:      uint32(prediction.ContestID),
		UserId:         uint32(prediction.UserID),
		EventId:        uint32(prediction.EventID),
		PredictionData: prediction.PredictionData,
		Status:         prediction.Status,
		SubmittedAt:    timestamppb.New(prediction.SubmittedAt),
		CreatedAt:      timestamppb.New(prediction.CreatedAt),
		UpdatedAt:      timestamppb.New(prediction.UpdatedAt),
	}
	return pbPrediction
}

func (s *PredictionService) eventModelToPB(event *models.Event) *pb.Event {
	if event == nil {
		return nil
	}
	
	return &pb.Event{
		Id:         uint32(event.ID),
		Title:      event.Title,
		SportType:  event.SportType,
		HomeTeam:   event.HomeTeam,
		AwayTeam:   event.AwayTeam,
		EventDate:  timestamppb.New(event.EventDate),
		Status:     event.Status,
		ResultData: event.ResultData,
		CreatedAt:  timestamppb.New(event.CreatedAt),
		UpdatedAt:  timestamppb.New(event.UpdatedAt),
	}
}

// GetPropTypes returns prop types for a sport
func (s *PredictionService) GetPropTypes(ctx context.Context, req *pb.GetPropTypesRequest) (*pb.GetPropTypesResponse, error) {
	if req.SportType == "" {
		return &pb.GetPropTypesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Sport type is required",
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	propTypes, err := s.propTypeRepo.GetBySportType(ctx, req.SportType)
	if err != nil {
		return &pb.GetPropTypesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to retrieve prop types",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	protoPropTypes := make([]*pb.PropType, len(propTypes))
	for i, pt := range propTypes {
		protoPropTypes[i] = s.propTypeToProto(pt)
	}

	return &pb.GetPropTypesResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Prop types retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		PropTypes: protoPropTypes,
	}, nil
}

// ListPropTypes returns paginated prop types with filtering
func (s *PredictionService) ListPropTypes(ctx context.Context, req *pb.ListPropTypesRequest) (*pb.ListPropTypesResponse, error) {
	page, limit := 1, 20
	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
		if req.Pagination.Limit > 0 {
			limit = int(req.Pagination.Limit)
		}
	}

	propTypes, total, err := s.propTypeRepo.List(ctx, req.SportType, req.Category, req.ActiveOnly, page, limit)
	if err != nil {
		return &pb.ListPropTypesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to list prop types",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	protoPropTypes := make([]*pb.PropType, len(propTypes))
	for i, pt := range propTypes {
		protoPropTypes[i] = s.propTypeToProto(pt)
	}

	totalPages := (int(total) + limit - 1) / limit

	return &pb.ListPropTypesResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Prop types listed successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		PropTypes: protoPropTypes,
		Pagination: &common.PaginationResponse{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      int32(total),
			TotalPages: int32(totalPages),
		},
	}, nil
}

func (s *PredictionService) propTypeToProto(pt *models.PropType) *pb.PropType {
	proto := &pb.PropType{
		Id:            uint32(pt.ID),
		SportType:     pt.SportType,
		Name:          pt.Name,
		Slug:          pt.Slug,
		Description:   pt.Description,
		Category:      pt.Category,
		ValueType:     pt.ValueType,
		PointsCorrect: pt.PointsCorrect,
		IsActive:      pt.IsActive,
	}
	if pt.DefaultLine != nil {
		proto.DefaultLine = *pt.DefaultLine
	}
	if pt.MinValue != nil {
		proto.MinValue = *pt.MinValue
	}
	if pt.MaxValue != nil {
		proto.MaxValue = *pt.MaxValue
	}
	return proto
}

// GetPotentialCoefficient calculates the current time coefficient for an event
func (s *PredictionService) GetPotentialCoefficient(ctx context.Context, req *pb.GetPotentialCoefficientRequest) (*pb.GetPotentialCoefficientResponse, error) {
	event, err := s.eventRepo.GetByID(uint(req.EventId))
	if err != nil {
		return &pb.GetPotentialCoefficientResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Event not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	now := time.Now().UTC()
	hoursUntilEvent := event.EventDate.Sub(now).Hours()
	result := coefficient.Calculate(now, event.EventDate)

	return &pb.GetPotentialCoefficientResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Coefficient calculated",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Coefficient:     result.Coefficient,
		Tier:            result.Tier,
		HoursUntilEvent: hoursUntilEvent,
	}, nil
}

// ============= Risky Event Types =============

// ListRiskyEventTypes returns all risky event types
func (s *PredictionService) ListRiskyEventTypes(ctx context.Context, req *pb.ListRiskyEventTypesRequest) (*pb.ListRiskyEventTypesResponse, error) {
	var eventTypes []models.RiskyEventType
	var err error

	if req.IncludeInactive {
		eventTypes, err = s.riskyEventRepo.ListAllEventTypes(req.SportType)
	} else {
		eventTypes, err = s.riskyEventRepo.ListActiveEventTypes(req.SportType)
	}

	if err != nil {
		return &pb.ListRiskyEventTypesResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to list risky event types",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	protoTypes := make([]*pb.RiskyEventType, len(eventTypes))
	for i, et := range eventTypes {
		protoTypes[i] = s.riskyEventTypeToProto(&et)
	}

	return &pb.ListRiskyEventTypesResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Risky event types listed",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		EventTypes: protoTypes,
	}, nil
}

// CreateRiskyEventType creates a new risky event type (admin only)
func (s *PredictionService) CreateRiskyEventType(ctx context.Context, req *pb.CreateRiskyEventTypeRequest) (*pb.CreateRiskyEventTypeResponse, error) {
	// TODO: Add admin role check

	et := &models.RiskyEventType{
		Slug:          req.Slug,
		Name:          req.Name,
		NameEn:        req.NameEn,
		Description:   req.Description,
		DefaultPoints: req.DefaultPoints,
		SportType:     req.SportType,
		Category:      req.Category,
		Icon:          req.Icon,
		SortOrder:     int(req.SortOrder),
		IsActive:      true,
	}

	if et.SportType == "" {
		et.SportType = "football"
	}
	if et.Category == "" {
		et.Category = "general"
	}

	if err := s.riskyEventRepo.CreateEventType(et); err != nil {
		return &pb.CreateRiskyEventTypeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to create risky event type: " + err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.CreateRiskyEventTypeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Risky event type created",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		EventType: s.riskyEventTypeToProto(et),
	}, nil
}

// UpdateRiskyEventType updates an existing risky event type (admin only)
func (s *PredictionService) UpdateRiskyEventType(ctx context.Context, req *pb.UpdateRiskyEventTypeRequest) (*pb.UpdateRiskyEventTypeResponse, error) {
	// TODO: Add admin role check

	et, err := s.riskyEventRepo.GetEventType(uint(req.Id))
	if err != nil {
		return &pb.UpdateRiskyEventTypeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Risky event type not found",
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update fields
	if req.Name != "" {
		et.Name = req.Name
	}
	if req.NameEn != "" {
		et.NameEn = req.NameEn
	}
	if req.Description != "" {
		et.Description = req.Description
	}
	if req.DefaultPoints > 0 {
		et.DefaultPoints = req.DefaultPoints
	}
	if req.Category != "" {
		et.Category = req.Category
	}
	if req.Icon != "" {
		et.Icon = req.Icon
	}
	et.SortOrder = int(req.SortOrder)
	et.IsActive = req.IsActive

	if err := s.riskyEventRepo.UpdateEventType(et); err != nil {
		return &pb.UpdateRiskyEventTypeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to update risky event type",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.UpdateRiskyEventTypeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Risky event type updated",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		EventType: s.riskyEventTypeToProto(et),
	}, nil
}

// DeleteRiskyEventType soft-deletes a risky event type (admin only)
func (s *PredictionService) DeleteRiskyEventType(ctx context.Context, req *pb.DeleteRiskyEventTypeRequest) (*common.Response, error) {
	// TODO: Add admin role check

	if err := s.riskyEventRepo.DeleteEventType(uint(req.Id)); err != nil {
		return &common.Response{
			Success:   false,
			Message:   "Failed to delete risky event type",
			Code:      int32(common.ErrorCode_INTERNAL_ERROR),
			Timestamp: timestamppb.Now(),
		}, nil
	}

	return &common.Response{
		Success:   true,
		Message:   "Risky event type deleted",
		Code:      0,
		Timestamp: timestamppb.Now(),
	}, nil
}

// ============= Match Risky Events =============

// GetMatchRiskyEvents returns risky events for a specific match with overrides applied
func (s *PredictionService) GetMatchRiskyEvents(ctx context.Context, req *pb.GetMatchRiskyEventsRequest) (*pb.GetMatchRiskyEventsResponse, error) {
	// Get contest rules if contest_id provided
	var contestRulesJSON string
	maxSelections := 5 // default

	if req.ContestId > 0 {
		contest, err := s.contestClient.GetContest(ctx, req.ContestId)
		if err == nil && contest != nil {
			contestRulesJSON = contest.Rules
			// Parse max_selections from rules
			maxSelections = parseMaxSelections(contestRulesJSON)
		}
	}

	events, err := s.riskyEventRepo.GetMatchRiskyEvents(uint(req.EventId), contestRulesJSON)
	if err != nil {
		return &pb.GetMatchRiskyEventsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to get match risky events",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	protoEvents := make([]*pb.MatchRiskyEvent, len(events))
	for i, e := range events {
		protoEvents[i] = &pb.MatchRiskyEvent{
			RiskyEventTypeId: uint32(e.RiskyEventTypeID),
			Slug:             e.Slug,
			Name:             e.Name,
			NameEn:           e.NameEn,
			Icon:             e.Icon,
			Category:         e.Category,
			Points:           e.Points,
			IsEnabled:        e.IsEnabled,
			IsOverridden:     e.IsOverridden,
		}
		if e.Outcome != nil {
			protoEvents[i].Outcome = e.Outcome
		}
	}

	return &pb.GetMatchRiskyEventsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Match risky events retrieved",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Events:        protoEvents,
		MaxSelections: int32(maxSelections),
	}, nil
}

// SetMatchRiskyEventOverride sets a point override for a match (admin only)
func (s *PredictionService) SetMatchRiskyEventOverride(ctx context.Context, req *pb.SetMatchRiskyEventOverrideRequest) (*pb.SetMatchRiskyEventOverrideResponse, error) {
	// TODO: Add admin role check

	err := s.riskyEventRepo.SetMatchEventOverride(
		uint(req.EventId),
		uint(req.RiskyEventTypeId),
		req.Points,
		req.IsEnabled,
	)
	if err != nil {
		return &pb.SetMatchRiskyEventOverrideResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to set override: " + err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.SetMatchRiskyEventOverrideResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Override set",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

// SetMatchRiskyEventOutcome records the outcome of a risky event after match
func (s *PredictionService) SetMatchRiskyEventOutcome(ctx context.Context, req *pb.SetMatchRiskyEventOutcomeRequest) (*pb.SetMatchRiskyEventOutcomeResponse, error) {
	// TODO: Add admin role check

	err := s.riskyEventRepo.SetMatchEventOutcome(
		uint(req.EventId),
		uint(req.RiskyEventTypeId),
		req.Happened,
	)
	if err != nil {
		return &pb.SetMatchRiskyEventOutcomeResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Failed to set outcome: " + err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.SetMatchRiskyEventOutcomeResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Outcome recorded",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

// Helper: convert RiskyEventType to proto
func (s *PredictionService) riskyEventTypeToProto(et *models.RiskyEventType) *pb.RiskyEventType {
	return &pb.RiskyEventType{
		Id:            uint32(et.ID),
		Slug:          et.Slug,
		Name:          et.Name,
		NameEn:        et.NameEn,
		Description:   et.Description,
		DefaultPoints: et.DefaultPoints,
		SportType:     et.SportType,
		Category:      et.Category,
		Icon:          et.Icon,
		SortOrder:     int32(et.SortOrder),
		IsActive:      et.IsActive,
	}
}

// Helper: parse max_selections from contest rules JSON
func parseMaxSelections(rulesJSON string) int {
	if rulesJSON == "" {
		return 5
	}
	var rules struct {
		Risky *struct {
			MaxSelections int `json:"max_selections"`
		} `json:"risky"`
	}
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return 5
	}
	if rules.Risky != nil && rules.Risky.MaxSelections > 0 {
		return rules.Risky.MaxSelections
	}
	return 5
}

// Helper: parse contest type from rules JSON
func parseContestType(rulesJSON string) string {
	if rulesJSON == "" {
		return "standard"
	}
	var rules struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return "standard"
	}
	if rules.Type == "" {
		return "standard"
	}
	return rules.Type
}

// SetContestEvents sets the events for a contest (replaces existing)
func (s *PredictionService) SetContestEvents(ctx context.Context, req *pb.SetContestEventsRequest) (*pb.SetContestEventsResponse, error) {
	if req.ContestId == 0 {
		return &pb.SetContestEventsResponse{
			Response: &common.Response{
				Success: false,
				Message: "contest_id is required",
			},
		}, nil
	}

	// Convert uint64 to uint
	eventIDs := make([]uint, len(req.EventIds))
	for i, id := range req.EventIds {
		eventIDs[i] = uint(id)
	}

	// Set events for the contest
	err := s.eventRepo.SetContestEvents(uint(req.ContestId), eventIDs)
	if err != nil {
		return &pb.SetContestEventsResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.SetContestEventsResponse{
		Response: &common.Response{
			Success: true,
			Message: "Events set successfully",
		},
		EventCount: int32(len(eventIDs)),
	}, nil
}

// GetContestEventCount returns the number of events in a contest
func (s *PredictionService) GetContestEventCount(ctx context.Context, req *pb.GetContestEventCountRequest) (*pb.GetContestEventCountResponse, error) {
	if req.ContestId == 0 {
		return &pb.GetContestEventCountResponse{
			Response: &common.Response{
				Success: false,
				Message: "contest_id is required",
			},
		}, nil
	}

	count, err := s.eventRepo.GetContestEventCount(uint(req.ContestId))
	if err != nil {
		return &pb.GetContestEventCountResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.GetContestEventCountResponse{
		Response: &common.Response{
			Success: true,
		},
		EventCount: int32(count),
	}, nil
}

// SetRelayAssignments sets event assignments for a team (captain action)
func (s *PredictionService) SetRelayAssignments(ctx context.Context, req *pb.SetRelayAssignmentsRequest) (*pb.SetRelayAssignmentsResponse, error) {
	if req.ContestId == 0 || req.TeamId == 0 {
		return &pb.SetRelayAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: "contest_id and team_id are required",
			},
		}, nil
	}

	// Get user ID from context (the user making the request)
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &pb.SetRelayAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: "unauthorized: user ID not found in context",
			},
		}, nil
	}

	// Verify that user is actually the captain of this team
	isCaptain, err := s.teamClient.IsTeamCaptain(ctx, uint32(req.TeamId), uint64(userID))
	if err != nil {
		return &pb.SetRelayAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: "failed to verify captain status: " + err.Error(),
			},
		}, nil
	}
	if !isCaptain {
		return &pb.SetRelayAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: "only team captain can assign events",
			},
		}, nil
	}

	// Convert proto assignments to repository input
	assignments := make([]repository.RelayAssignmentInput, len(req.Assignments))
	for i, a := range req.Assignments {
		assignments[i] = repository.RelayAssignmentInput{
			UserID:  uint(a.UserId),
			EventID: uint(a.EventId),
		}
	}

	// Set assignments
	err = s.relayRepo.SetTeamAssignments(uint(req.ContestId), uint(req.TeamId), uint(userID), assignments)
	if err != nil {
		return &pb.SetRelayAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.SetRelayAssignmentsResponse{
		Response: &common.Response{
			Success: true,
			Message: "Assignments saved successfully",
		},
		AssignedCount: int32(len(assignments)),
	}, nil
}

// GetTeamAssignments retrieves all event assignments for a team
func (s *PredictionService) GetTeamAssignments(ctx context.Context, req *pb.GetTeamAssignmentsRequest) (*pb.GetTeamAssignmentsResponse, error) {
	if req.ContestId == 0 || req.TeamId == 0 {
		return &pb.GetTeamAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: "contest_id and team_id are required",
			},
		}, nil
	}

	// Get assignments
	assignments, err := s.relayRepo.GetTeamAssignments(uint(req.ContestId), uint(req.TeamId))
	if err != nil {
		return &pb.GetTeamAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	// Get stats
	stats, err := s.relayRepo.GetAssignmentStats(uint(req.ContestId), uint(req.TeamId))
	if err != nil {
		return &pb.GetTeamAssignmentsResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	// Convert to proto
	protoAssignments := make([]*pb.RelayAssignment, len(assignments))
	for i, a := range assignments {
		protoAssignments[i] = &pb.RelayAssignment{
			UserId:  uint64(a.UserID),
			EventId: uint64(a.EventID),
			Event:   s.eventModelToPB(&a.Event),
		}
	}

	return &pb.GetTeamAssignmentsResponse{
		Response: &common.Response{
			Success: true,
		},
		Assignments:    protoAssignments,
		TotalEvents:    int32(stats.TotalEvents),
		AssignedEvents: int32(stats.AssignedEvents),
	}, nil
}

// GetUserRelayEvents retrieves events assigned to the current user
func (s *PredictionService) GetUserRelayEvents(ctx context.Context, req *pb.GetUserRelayEventsRequest) (*pb.GetUserRelayEventsResponse, error) {
	if req.ContestId == 0 {
		return &pb.GetUserRelayEventsResponse{
			Response: &common.Response{
				Success: false,
				Message: "contest_id is required",
			},
		}, nil
	}

	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &pb.GetUserRelayEventsResponse{
			Response: &common.Response{
				Success: false,
				Message: "unauthorized: user ID not found in context",
			},
		}, nil
	}

	var assignments []*models.RelayEventAssignment
	var err error

	if req.TeamId > 0 {
		// Get assignments for specific team
		assignments, err = s.relayRepo.GetUserAssignmentsForTeam(uint(req.ContestId), uint(req.TeamId), uint(userID))
	} else {
		// Get all assignments for user in this contest
		assignments, err = s.relayRepo.GetUserAssignments(uint(req.ContestId), uint(userID))
	}

	if err != nil {
		return &pb.GetUserRelayEventsResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
			},
		}, nil
	}

	// Extract events
	events := make([]*pb.Event, len(assignments))
	var teamID uint64
	for i, a := range assignments {
		events[i] = s.eventModelToPB(&a.Event)
		if teamID == 0 {
			teamID = uint64(a.TeamID)
		}
	}

	return &pb.GetUserRelayEventsResponse{
		Response: &common.Response{
			Success: true,
		},
		Events: events,
		TeamId: teamID,
	}, nil
}

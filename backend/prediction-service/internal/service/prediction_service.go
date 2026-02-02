package service

import (
	"context"
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
	predictionRepo repository.PredictionRepositoryInterface
	eventRepo      repository.EventRepositoryInterface
	propTypeRepo   *repository.PropTypeRepository
	contestClient  *clients.ContestClient
}

// NewPredictionService creates a new PredictionService instance
func NewPredictionService(
	predictionRepo repository.PredictionRepositoryInterface,
	eventRepo repository.EventRepositoryInterface,
	propTypeRepo *repository.PropTypeRepository,
	contestClient *clients.ContestClient,
) *PredictionService {
	return &PredictionService{
		predictionRepo: predictionRepo,
		eventRepo:      eventRepo,
		propTypeRepo:   propTypeRepo,
		contestClient:  contestClient,
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

	if existingPrediction != nil {
		return &pb.SubmitPredictionResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Prediction already exists for this event",
				Code:      int32(common.ErrorCode_ALREADY_EXISTS),
				Timestamp: timestamppb.Now(),
			},
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

	events, total, err := s.eventRepo.List(limit, offset, req.SportType, req.Status)
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

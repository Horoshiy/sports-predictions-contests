package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	pb "github.com/sports-prediction-contests/shared/proto/prediction"
	"github.com/sports-prediction-contests/shared/proto/common"
)

// Mock repositories and clients
type MockPredictionRepository struct {
	mock.Mock
}

func (m *MockPredictionRepository) Create(prediction *models.Prediction) error {
	args := m.Called(prediction)
	return args.Error(0)
}

func (m *MockPredictionRepository) GetByID(id uint) (*models.Prediction, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Prediction), args.Error(1)
}

func (m *MockPredictionRepository) GetByUserAndContest(userID, contestID uint) ([]*models.Prediction, error) {
	args := m.Called(userID, contestID)
	return args.Get(0).([]*models.Prediction), args.Error(1)
}

func (m *MockPredictionRepository) GetByUserContestAndEvent(userID, contestID, eventID uint) (*models.Prediction, error) {
	args := m.Called(userID, contestID, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Prediction), args.Error(1)
}

func (m *MockPredictionRepository) Update(prediction *models.Prediction) error {
	args := m.Called(prediction)
	return args.Error(0)
}

func (m *MockPredictionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPredictionRepository) List(limit, offset int, contestID uint, userID uint) ([]*models.Prediction, int64, error) {
	args := m.Called(limit, offset, contestID, userID)
	return args.Get(0).([]*models.Prediction), args.Get(1).(int64), args.Error(2)
}

func (m *MockPredictionRepository) CountByContest(contestID uint) (int64, error) {
	args := m.Called(contestID)
	return args.Get(0).(int64), args.Error(1)
}

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(id uint) (*models.Event, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRepository) Update(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventRepository) List(limit, offset int, sportType, status string) ([]*models.Event, int64, error) {
	args := m.Called(limit, offset, sportType, status)
	return args.Get(0).([]*models.Event), args.Get(1).(int64), args.Error(2)
}

type MockContestClient struct {
	mock.Mock
}

func (m *MockContestClient) ValidateContestParticipation(ctx context.Context, contestID uint32, userID uint32) error {
	args := m.Called(ctx, contestID, userID)
	return args.Error(0)
}

func TestModelToPBNilHandling(t *testing.T) {
	service := &PredictionService{}

	// Test nil prediction
	result := service.modelToPB(nil)
	assert.Nil(t, result)

	// Test nil event
	eventResult := service.eventModelToPB(nil)
	assert.Nil(t, eventResult)

	// Test valid prediction
	prediction := &models.Prediction{
		ID:             1,
		ContestID:      1,
		UserID:         1,
		EventID:        1,
		PredictionData: `{"score": "2-1"}`,
		Status:         "pending",
		SubmittedAt:    time.Now(),
	}
	
	pbPrediction := service.modelToPB(prediction)
	assert.NotNil(t, pbPrediction)
	assert.Equal(t, uint32(1), pbPrediction.Id)
	assert.Equal(t, "pending", pbPrediction.Status)
}

func TestSubmitPredictionDuplicateHandling(t *testing.T) {
	mockPredRepo := &MockPredictionRepository{}
	mockEventRepo := &MockEventRepository{}
	mockContestClient := &MockContestClient{}

	service := NewPredictionService(mockPredRepo, mockEventRepo, mockContestClient)

	ctx := context.WithValue(context.Background(), "userID", uint(1))
	
	req := &pb.SubmitPredictionRequest{
		ContestId:      1,
		EventId:        1,
		PredictionData: `{"score": "2-1"}`,
	}

	// Mock contest validation success
	mockContestClient.On("ValidateContestParticipation", ctx, uint32(1), uint32(1)).Return(nil)

	// Mock event exists and can accept predictions
	event := &models.Event{
		ID:        1,
		Status:    "scheduled",
		EventDate: time.Now().Add(1 * time.Hour),
	}
	mockEventRepo.On("GetByID", uint(1)).Return(event, nil)

	// Mock no existing prediction
	mockPredRepo.On("GetByUserContestAndEvent", uint(1), uint(1), uint(1)).Return(nil, nil)

	// Mock database constraint violation
	mockPredRepo.On("Create", mock.AnythingOfType("*models.Prediction")).Return(
		errors.New("UNIQUE constraint failed: predictions.user_id, predictions.contest_id, predictions.event_id"))

	resp, err := service.SubmitPrediction(ctx, req)

	assert.NoError(t, err)
	assert.False(t, resp.Response.Success)
	assert.Equal(t, "Prediction already exists for this event", resp.Response.Message)
	assert.Equal(t, int32(common.ErrorCode_ALREADY_EXISTS), resp.Response.Code)
}

func TestGetUserPredictionsPagination(t *testing.T) {
	mockPredRepo := &MockPredictionRepository{}
	mockEventRepo := &MockEventRepository{}
	mockContestClient := &MockContestClient{}

	service := NewPredictionService(mockPredRepo, mockEventRepo, mockContestClient)

	ctx := context.WithValue(context.Background(), "userID", uint(1))
	
	req := &pb.GetUserPredictionsRequest{
		ContestId: 1,
		Pagination: &common.PaginationRequest{
			Page:     2,
			PageSize: 5,
		},
	}

	predictions := []*models.Prediction{
		{ID: 1, UserID: 1, ContestID: 1},
		{ID: 2, UserID: 1, ContestID: 1},
	}

	// Mock repository call with correct pagination
	mockPredRepo.On("List", 5, 5, uint(1), uint(1)).Return(predictions, int64(12), nil)

	resp, err := service.GetUserPredictions(ctx, req)

	assert.NoError(t, err)
	assert.True(t, resp.Response.Success)
	assert.Equal(t, int32(2), resp.Pagination.Page)
	assert.Equal(t, int32(5), resp.Pagination.Limit)
	assert.Equal(t, int32(12), resp.Pagination.Total)
	assert.Equal(t, int32(3), resp.Pagination.TotalPages) // ceil(12/5) = 3
	assert.Len(t, resp.Predictions, 2)
}

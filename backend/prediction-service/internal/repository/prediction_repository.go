package repository

import (
	"errors"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"gorm.io/gorm"
)

// PredictionRepositoryInterface defines the contract for prediction repository
type PredictionRepositoryInterface interface {
	Create(prediction *models.Prediction) error
	GetByID(id uint) (*models.Prediction, error)
	GetByUserAndContest(userID, contestID uint) ([]*models.Prediction, error)
	GetByUserContestAndEvent(userID, contestID, eventID uint) (*models.Prediction, error)
	Update(prediction *models.Prediction) error
	Delete(id uint) error
	List(limit, offset int, contestID uint, userID uint) ([]*models.Prediction, int64, error)
	CountByContest(contestID uint) (int64, error)
}

// PredictionRepository implements PredictionRepositoryInterface
type PredictionRepository struct {
	db *gorm.DB
}

// NewPredictionRepository creates a new prediction repository instance
func NewPredictionRepository(db *gorm.DB) PredictionRepositoryInterface {
	return &PredictionRepository{db: db}
}

// Create creates a new prediction
func (r *PredictionRepository) Create(prediction *models.Prediction) error {
	if prediction == nil {
		return errors.New("prediction cannot be nil")
	}

	return r.db.Create(prediction).Error
}

// GetByID retrieves a prediction by ID
func (r *PredictionRepository) GetByID(id uint) (*models.Prediction, error) {
	var prediction models.Prediction
	err := r.db.Preload("Event").First(&prediction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prediction not found")
		}
		return nil, err
	}
	return &prediction, nil
}

// GetByUserAndContest retrieves predictions by user ID and contest ID
func (r *PredictionRepository) GetByUserAndContest(userID, contestID uint) ([]*models.Prediction, error) {
	var predictions []*models.Prediction
	err := r.db.Preload("Event").Where("user_id = ? AND contest_id = ?", userID, contestID).Find(&predictions).Error
	return predictions, err
}

// GetByUserContestAndEvent retrieves a specific prediction by user, contest, and event
func (r *PredictionRepository) GetByUserContestAndEvent(userID, contestID, eventID uint) (*models.Prediction, error) {
	var prediction models.Prediction
	err := r.db.Preload("Event").Where("user_id = ? AND contest_id = ? AND event_id = ?", userID, contestID, eventID).First(&prediction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error for not found case
		}
		return nil, err
	}
	return &prediction, nil
}

// Update updates an existing prediction
func (r *PredictionRepository) Update(prediction *models.Prediction) error {
	if prediction == nil {
		return errors.New("prediction cannot be nil")
	}

	// Use raw SQL update to avoid GORM trying to save related Event
	return r.db.Exec(
		"UPDATE predictions SET prediction_data = ?, submitted_at = ?, updated_at = NOW() WHERE id = ?",
		prediction.PredictionData, prediction.SubmittedAt, prediction.ID,
	).Error
}

// Delete deletes a prediction by ID
func (r *PredictionRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Prediction{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("prediction not found")
	}

	return nil
}

// List retrieves predictions with pagination and optional filters
func (r *PredictionRepository) List(limit, offset int, contestID uint, userID uint) ([]*models.Prediction, int64, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Cap at 100 to prevent performance issues
	}
	if offset < 0 {
		offset = 0
	}

	var predictions []*models.Prediction
	var total int64

	query := r.db.Model(&models.Prediction{}).Preload("Event")

	// Apply filters
	if contestID > 0 {
		query = query.Where("contest_id = ?", contestID)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch records
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&predictions).Error
	return predictions, total, err
}

// CountByContest counts predictions for a specific contest
func (r *PredictionRepository) CountByContest(contestID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Prediction{}).Where("contest_id = ?", contestID).Count(&count).Error
	return count, err
}

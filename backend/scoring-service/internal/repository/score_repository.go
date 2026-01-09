package repository

import (
	"context"
	"errors"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
	"gorm.io/gorm"
)

// ScoreRepositoryInterface defines the contract for score repository
type ScoreRepositoryInterface interface {
	Create(ctx context.Context, score *models.Score) error
	GetByID(ctx context.Context, id uint) (*models.Score, error)
	Update(ctx context.Context, score *models.Score) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int, contestID, userID uint) ([]*models.Score, int64, error)
	GetByContestAndUser(ctx context.Context, contestID, userID uint) ([]*models.Score, error)
	GetByPrediction(ctx context.Context, predictionID uint) (*models.Score, error)
	BatchCreate(ctx context.Context, scores []*models.Score) error
	GetTotalPointsByContestAndUser(ctx context.Context, contestID, userID uint) (float64, error)
	ListByContest(ctx context.Context, contestID uint) ([]*models.Score, error)
}

// ScoreRepository implements ScoreRepositoryInterface
type ScoreRepository struct {
	db *gorm.DB
}

// NewScoreRepository creates a new score repository instance
func NewScoreRepository(db *gorm.DB) ScoreRepositoryInterface {
	return &ScoreRepository{db: db}
}

// Create creates a new score record
func (r *ScoreRepository) Create(ctx context.Context, score *models.Score) error {
	if err := r.db.WithContext(ctx).Create(score).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a score by its ID
func (r *ScoreRepository) GetByID(ctx context.Context, id uint) (*models.Score, error) {
	var score models.Score
	if err := r.db.WithContext(ctx).First(&score, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("score not found")
		}
		return nil, err
	}
	return &score, nil
}

// Update updates an existing score record
func (r *ScoreRepository) Update(ctx context.Context, score *models.Score) error {
	if err := r.db.WithContext(ctx).Save(score).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a score record by ID
func (r *ScoreRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Score{}, id).Error; err != nil {
		return err
	}
	return nil
}

// List retrieves scores with pagination and optional filtering
func (r *ScoreRepository) List(ctx context.Context, limit, offset int, contestID, userID uint) ([]*models.Score, int64, error) {
	var scores []*models.Score
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Score{})

	// Apply filters
	if contestID > 0 {
		query = query.Where("contest_id = ?", contestID)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Order("scored_at DESC").Find(&scores).Error; err != nil {
		return nil, 0, err
	}

	return scores, total, nil
}

// GetByContestAndUser retrieves all scores for a specific user in a contest
func (r *ScoreRepository) GetByContestAndUser(ctx context.Context, contestID, userID uint) ([]*models.Score, error) {
	var scores []*models.Score
	if err := r.db.WithContext(ctx).Where("contest_id = ? AND user_id = ?", contestID, userID).
		Order("scored_at DESC").Find(&scores).Error; err != nil {
		return nil, err
	}
	return scores, nil
}

// GetByPrediction retrieves a score by prediction ID
func (r *ScoreRepository) GetByPrediction(ctx context.Context, predictionID uint) (*models.Score, error) {
	var score models.Score
	if err := r.db.WithContext(ctx).Where("prediction_id = ?", predictionID).First(&score).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("score not found")
		}
		return nil, err
	}
	return &score, nil
}

// BatchCreate creates multiple score records in a single transaction
func (r *ScoreRepository) BatchCreate(ctx context.Context, scores []*models.Score) error {
	if len(scores) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, score := range scores {
			if err := tx.Create(score).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetTotalPointsByContestAndUser calculates total points for a user in a contest
func (r *ScoreRepository) GetTotalPointsByContestAndUser(ctx context.Context, contestID, userID uint) (float64, error) {
	var totalPoints float64
	if err := r.db.WithContext(ctx).Model(&models.Score{}).
		Where("contest_id = ? AND user_id = ?", contestID, userID).
		Select("COALESCE(SUM(points), 0)").Scan(&totalPoints).Error; err != nil {
		return 0, err
	}
	return totalPoints, nil
}

// ListByContest retrieves all scores for a specific contest
func (r *ScoreRepository) ListByContest(ctx context.Context, contestID uint) ([]*models.Score, error) {
	var scores []*models.Score
	if err := r.db.WithContext(ctx).Where("contest_id = ?", contestID).
		Order("scored_at DESC").Find(&scores).Error; err != nil {
		return nil, err
	}
	return scores, nil
}

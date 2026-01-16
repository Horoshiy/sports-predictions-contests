package repository

import (
	"context"
	"errors"

	"github.com/sports-prediction-contests/scoring-service/internal/models"
	"gorm.io/gorm"
)

// StreakRepositoryInterface defines the contract for streak repository
type StreakRepositoryInterface interface {
	GetOrCreate(ctx context.Context, contestID, userID uint) (*models.UserStreak, error)
	Update(ctx context.Context, streak *models.UserStreak) error
	GetByContestAndUser(ctx context.Context, contestID, userID uint) (*models.UserStreak, error)
	GetByContestAndUsers(ctx context.Context, contestID uint, userIDs []uint) ([]*models.UserStreak, error)
	GetTopStreaks(ctx context.Context, contestID uint, limit int) ([]*models.UserStreak, error)
}

// StreakRepository implements StreakRepositoryInterface
type StreakRepository struct {
	db *gorm.DB
}

// NewStreakRepository creates a new streak repository instance
func NewStreakRepository(db *gorm.DB) StreakRepositoryInterface {
	return &StreakRepository{db: db}
}

// GetOrCreate retrieves an existing streak or creates a new one (atomic)
func (r *StreakRepository) GetOrCreate(ctx context.Context, contestID, userID uint) (*models.UserStreak, error) {
	streak := models.UserStreak{
		UserID:        userID,
		ContestID:     contestID,
		CurrentStreak: 0,
		MaxStreak:     0,
	}
	// FirstOrCreate is atomic - handles race conditions
	err := r.db.WithContext(ctx).Where("contest_id = ? AND user_id = ?", contestID, userID).
		FirstOrCreate(&streak).Error
	if err != nil {
		return nil, err
	}
	return &streak, nil
}

// Update updates an existing streak record
func (r *StreakRepository) Update(ctx context.Context, streak *models.UserStreak) error {
	return r.db.WithContext(ctx).Save(streak).Error
}

// GetByContestAndUser retrieves a streak for a specific user in a contest
func (r *StreakRepository) GetByContestAndUser(ctx context.Context, contestID, userID uint) (*models.UserStreak, error) {
	var streak models.UserStreak
	if err := r.db.WithContext(ctx).Where("contest_id = ? AND user_id = ?", contestID, userID).First(&streak).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("streak not found")
		}
		return nil, err
	}
	return &streak, nil
}

// GetByContestAndUsers retrieves streaks for multiple users in a contest (batch)
func (r *StreakRepository) GetByContestAndUsers(ctx context.Context, contestID uint, userIDs []uint) ([]*models.UserStreak, error) {
	if len(userIDs) == 0 {
		return []*models.UserStreak{}, nil
	}
	var streaks []*models.UserStreak
	if err := r.db.WithContext(ctx).Where("contest_id = ? AND user_id IN ?", contestID, userIDs).Find(&streaks).Error; err != nil {
		return nil, err
	}
	return streaks, nil
}

// GetTopStreaks retrieves the top streaks for a contest
func (r *StreakRepository) GetTopStreaks(ctx context.Context, contestID uint, limit int) ([]*models.UserStreak, error) {
	if limit <= 0 {
		limit = 10
	}
	var streaks []*models.UserStreak
	if err := r.db.WithContext(ctx).Where("contest_id = ?", contestID).
		Order("current_streak DESC").Limit(limit).Find(&streaks).Error; err != nil {
		return nil, err
	}
	return streaks, nil
}

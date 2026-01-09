package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sports-prediction-contests/scoring-service/internal/cache"
	"github.com/sports-prediction-contests/scoring-service/internal/models"
	"gorm.io/gorm"
)

// LeaderboardRepositoryInterface defines the contract for leaderboard repository
type LeaderboardRepositoryInterface interface {
	Create(ctx context.Context, leaderboard *models.Leaderboard) error
	GetByID(ctx context.Context, id uint) (*models.Leaderboard, error)
	Update(ctx context.Context, leaderboard *models.Leaderboard) error
	Delete(ctx context.Context, id uint) error
	GetByContestAndUser(ctx context.Context, contestID, userID uint) (*models.Leaderboard, error)
	ListByContest(ctx context.Context, contestID uint, limit, offset int) ([]*models.Leaderboard, int64, error)
	UpdateRankings(ctx context.Context, contestID uint) error
	UpsertUserScore(ctx context.Context, contestID, userID uint, totalPoints float64) error
	GetContestLeaderboard(ctx context.Context, contestID uint, limit int) ([]*models.Leaderboard, error)
	RecalculateRanks(ctx context.Context, contestID uint) error
}

// LeaderboardRepository implements LeaderboardRepositoryInterface
type LeaderboardRepository struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

// NewLeaderboardRepository creates a new leaderboard repository instance
func NewLeaderboardRepository(db *gorm.DB, cache *cache.RedisCache) LeaderboardRepositoryInterface {
	return &LeaderboardRepository{
		db:    db,
		cache: cache,
	}
}

// Create creates a new leaderboard entry
func (r *LeaderboardRepository) Create(ctx context.Context, leaderboard *models.Leaderboard) error {
	if err := r.db.WithContext(ctx).Create(leaderboard).Error; err != nil {
		return err
	}

	// Update cache
	if r.cache != nil {
		_ = r.cache.SetLeaderboardScore(ctx, leaderboard.ContestID, leaderboard.UserID, leaderboard.TotalPoints)
	}

	return nil
}

// GetByID retrieves a leaderboard entry by its ID
func (r *LeaderboardRepository) GetByID(ctx context.Context, id uint) (*models.Leaderboard, error) {
	var leaderboard models.Leaderboard
	if err := r.db.WithContext(ctx).First(&leaderboard, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leaderboard entry not found")
		}
		return nil, err
	}
	return &leaderboard, nil
}

// Update updates an existing leaderboard entry
func (r *LeaderboardRepository) Update(ctx context.Context, leaderboard *models.Leaderboard) error {
	if err := r.db.WithContext(ctx).Save(leaderboard).Error; err != nil {
		return err
	}

	// Update cache
	if r.cache != nil {
		_ = r.cache.SetLeaderboardScore(ctx, leaderboard.ContestID, leaderboard.UserID, leaderboard.TotalPoints)
	}

	return nil
}

// Delete deletes a leaderboard entry by ID
func (r *LeaderboardRepository) Delete(ctx context.Context, id uint) error {
	// Get the entry first to update cache
	leaderboard, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.Leaderboard{}, id).Error; err != nil {
		return err
	}

	// Remove from cache
	if r.cache != nil {
		_ = r.cache.RemoveUserFromLeaderboard(ctx, leaderboard.ContestID, leaderboard.UserID)
	}

	return nil
}

// GetByContestAndUser retrieves a leaderboard entry for a specific user in a contest
func (r *LeaderboardRepository) GetByContestAndUser(ctx context.Context, contestID, userID uint) (*models.Leaderboard, error) {
	var leaderboard models.Leaderboard
	if err := r.db.WithContext(ctx).Where("contest_id = ? AND user_id = ?", contestID, userID).
		First(&leaderboard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leaderboard entry not found")
		}
		return nil, err
	}
	return &leaderboard, nil
}

// ListByContest retrieves leaderboard entries for a contest with pagination
func (r *LeaderboardRepository) ListByContest(ctx context.Context, contestID uint, limit, offset int) ([]*models.Leaderboard, int64, error) {
	var leaderboards []*models.Leaderboard
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Leaderboard{}).Where("contest_id = ?", contestID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results ordered by rank
	if err := query.Order("rank ASC").Limit(limit).Offset(offset).Find(&leaderboards).Error; err != nil {
		return nil, 0, err
	}

	return leaderboards, total, nil
}

// UpdateRankings recalculates and updates rankings for all users in a contest
func (r *LeaderboardRepository) UpdateRankings(ctx context.Context, contestID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get all leaderboard entries for the contest ordered by total points (descending)
		var leaderboards []*models.Leaderboard
		if err := tx.Where("contest_id = ?", contestID).
			Order("total_points DESC, updated_at ASC").
			Find(&leaderboards).Error; err != nil {
			return err
		}

		// Update ranks with tie-breaking logic
		currentRank := uint(1)
		var previousPoints float64 = -1
		var actualRank uint = 1

		for i, leaderboard := range leaderboards {
			if leaderboard.TotalPoints != previousPoints {
				currentRank = uint(i + 1)
			}

			leaderboard.Rank = currentRank
			previousPoints = leaderboard.TotalPoints
			actualRank++

			if err := tx.Save(leaderboard).Error; err != nil {
				return err
			}

			// Update cache
			if r.cache != nil {
				_ = r.cache.SetLeaderboardScore(ctx, contestID, leaderboard.UserID, leaderboard.TotalPoints)
			}
		}

		return nil
	})
}

// UpsertUserScore creates or updates a user's leaderboard entry
func (r *LeaderboardRepository) UpsertUserScore(ctx context.Context, contestID, userID uint, totalPoints float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var leaderboard models.Leaderboard

		// Try to find existing entry
		err := tx.Where("contest_id = ? AND user_id = ?", contestID, userID).First(&leaderboard).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new entry
			leaderboard = models.Leaderboard{
				ContestID:   contestID,
				UserID:      userID,
				TotalPoints: totalPoints,
				Rank:        0, // Will be calculated later
			}
			if err := tx.Create(&leaderboard).Error; err != nil {
				return err
			}
		} else {
			// Update existing entry
			leaderboard.TotalPoints = totalPoints
			if err := tx.Save(&leaderboard).Error; err != nil {
				return err
			}
		}

		// Update cache
		if r.cache != nil {
			_ = r.cache.SetLeaderboardScore(ctx, contestID, userID, totalPoints)
		}

		return nil
	})
}

// GetContestLeaderboard retrieves the top N entries from a contest leaderboard
func (r *LeaderboardRepository) GetContestLeaderboard(ctx context.Context, contestID uint, limit int) ([]*models.Leaderboard, error) {
	// Try cache first
	if r.cache != nil {
		cacheEntries, err := r.cache.GetLeaderboard(ctx, contestID, int64(limit))
		if err == nil && len(cacheEntries) > 0 {
			// Convert cache entries to leaderboard models
			leaderboards := make([]*models.Leaderboard, len(cacheEntries))
			for i, entry := range cacheEntries {
				leaderboards[i] = &models.Leaderboard{
					ContestID:   contestID,
					UserID:      entry.UserID,
					TotalPoints: entry.TotalPoints,
					Rank:        entry.Rank,
				}
			}
			return leaderboards, nil
		}
	}

	// Fallback to database
	var leaderboards []*models.Leaderboard
	if err := r.db.WithContext(ctx).Where("contest_id = ?", contestID).
		Order("rank ASC").Limit(limit).Find(&leaderboards).Error; err != nil {
		return nil, err
	}

	// Warm up cache
	if r.cache != nil {
		for _, lb := range leaderboards {
			_ = r.cache.SetLeaderboardScore(ctx, contestID, lb.UserID, lb.TotalPoints)
		}
	}

	return leaderboards, nil
}

// RecalculateRanks recalculates ranks for a contest based on current scores
func (r *LeaderboardRepository) RecalculateRanks(ctx context.Context, contestID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get all leaderboard entries for the contest ordered by total points (descending)
		var leaderboards []*models.Leaderboard
		if err := tx.Where("contest_id = ?", contestID).
			Order("total_points DESC, updated_at ASC").
			Find(&leaderboards).Error; err != nil {
			return err
		}

		// Update ranks with tie-breaking logic
		currentRank := uint(1)
		var previousPoints float64 = -1

		for i, leaderboard := range leaderboards {
			if leaderboard.TotalPoints != previousPoints {
				currentRank = uint(i + 1)
			}

			leaderboard.Rank = currentRank
			previousPoints = leaderboard.TotalPoints

			if err := tx.Save(leaderboard).Error; err != nil {
				return err
			}

			// Update cache
			if r.cache != nil {
				_ = r.cache.SetLeaderboardScore(ctx, contestID, leaderboard.UserID, leaderboard.TotalPoints)
			}
		}

		return nil
	})
}

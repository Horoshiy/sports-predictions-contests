package repository

import (
	"context"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"gorm.io/gorm"
)

type PropTypeRepository struct {
	db *gorm.DB
}

func NewPropTypeRepository(db *gorm.DB) *PropTypeRepository {
	return &PropTypeRepository{db: db}
}

func (r *PropTypeRepository) GetBySportType(ctx context.Context, sportType string) ([]*models.PropType, error) {
	var propTypes []*models.PropType
	err := r.db.WithContext(ctx).
		Where("sport_type = ? AND is_active = ? AND deleted_at IS NULL", sportType, true).
		Order("category, name").
		Find(&propTypes).Error
	return propTypes, err
}

func (r *PropTypeRepository) GetByID(ctx context.Context, id uint) (*models.PropType, error) {
	var propType models.PropType
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&propType).Error
	if err != nil {
		return nil, err
	}
	return &propType, nil
}

func (r *PropTypeRepository) List(ctx context.Context, sportType, category string, activeOnly bool, page, limit int) ([]*models.PropType, int64, error) {
	var propTypes []*models.PropType
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PropType{}).Where("deleted_at IS NULL")

	if sportType != "" {
		query = query.Where("sport_type = ?", sportType)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("sport_type, category, name").
		Offset(offset).Limit(limit).
		Find(&propTypes).Error

	return propTypes, total, err
}

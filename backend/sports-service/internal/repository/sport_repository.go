package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type SportRepositoryInterface interface {
	Create(sport *models.Sport) error
	GetByID(id uint) (*models.Sport, error)
	GetBySlug(slug string) (*models.Sport, error)
	GetByExternalID(externalID string) (*models.Sport, error)
	Update(sport *models.Sport) error
	Upsert(sport *models.Sport) error
	Delete(id uint) error
	List(limit, offset int, activeOnly bool) ([]*models.Sport, int64, error)
}

type SportRepository struct {
	db *gorm.DB
}

func NewSportRepository(db *gorm.DB) SportRepositoryInterface {
	return &SportRepository{db: db}
}

func (r *SportRepository) Create(sport *models.Sport) error {
	if sport == nil {
		return errors.New("sport cannot be nil")
	}
	return r.db.Create(sport).Error
}

func (r *SportRepository) GetByID(id uint) (*models.Sport, error) {
	if id == 0 {
		return nil, errors.New("invalid sport ID")
	}
	var sport models.Sport
	result := r.db.First(&sport, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport not found")
		}
		return nil, result.Error
	}
	return &sport, nil
}

func (r *SportRepository) GetBySlug(slug string) (*models.Sport, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}
	var sport models.Sport
	result := r.db.Where("slug = ?", slug).First(&sport)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport not found")
		}
		return nil, result.Error
	}
	return &sport, nil
}

func (r *SportRepository) Update(sport *models.Sport) error {
	if sport == nil {
		return errors.New("sport cannot be nil")
	}
	if sport.ID == 0 {
		return errors.New("sport ID cannot be zero")
	}
	result := r.db.Save(sport)
	if result.RowsAffected == 0 {
		return errors.New("sport not found")
	}
	return result.Error
}

func (r *SportRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid sport ID")
	}
	result := r.db.Delete(&models.Sport{}, id)
	if result.RowsAffected == 0 {
		return errors.New("sport not found")
	}
	return result.Error
}

func (r *SportRepository) List(limit, offset int, activeOnly bool) ([]*models.Sport, int64, error) {
	var sports []*models.Sport
	var total int64

	query := r.db.Model(&models.Sport{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&sports).Error; err != nil {
		return nil, 0, err
	}

	return sports, total, nil
}

func (r *SportRepository) GetByExternalID(externalID string) (*models.Sport, error) {
	if externalID == "" {
		return nil, errors.New("external ID cannot be empty")
	}
	var sport models.Sport
	result := r.db.Where("external_id = ?", externalID).First(&sport)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport not found")
		}
		return nil, result.Error
	}
	return &sport, nil
}

func (r *SportRepository) Upsert(sport *models.Sport) error {
	if sport == nil {
		return errors.New("sport cannot be nil")
	}
	if sport.ExternalID == "" {
		return r.Create(sport)
	}
	existing, err := r.GetByExternalID(sport.ExternalID)
	if err == nil {
		sport.ID = existing.ID
		return r.Update(sport)
	}
	return r.Create(sport)
}

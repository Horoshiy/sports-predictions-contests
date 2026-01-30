package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type LeagueRepositoryInterface interface {
	Create(league *models.League) error
	GetByID(id uint) (*models.League, error)
	GetByExternalID(externalID string) (*models.League, error)
	Update(league *models.League) error
	Upsert(league *models.League) error
	Delete(id uint) error
	List(limit, offset int, sportID uint, activeOnly bool) ([]*models.League, int64, error)
}

type LeagueRepository struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) LeagueRepositoryInterface {
	return &LeagueRepository{db: db}
}

func (r *LeagueRepository) Create(league *models.League) error {
	if league == nil {
		return errors.New("league cannot be nil")
	}
	return r.db.Create(league).Error
}

func (r *LeagueRepository) GetByID(id uint) (*models.League, error) {
	if id == 0 {
		return nil, errors.New("invalid league ID")
	}
	var league models.League
	result := r.db.Preload("Sport").First(&league, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("league not found")
		}
		return nil, result.Error
	}
	return &league, nil
}

func (r *LeagueRepository) Update(league *models.League) error {
	if league == nil {
		return errors.New("league cannot be nil")
	}
	if league.ID == 0 {
		return errors.New("league ID cannot be zero")
	}
	result := r.db.Omit("Sport").Model(league).Select("*").Updates(league)
	if result.RowsAffected == 0 {
		return errors.New("league not found")
	}
	return result.Error
}

func (r *LeagueRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid league ID")
	}
	result := r.db.Delete(&models.League{}, id)
	if result.RowsAffected == 0 {
		return errors.New("league not found")
	}
	return result.Error
}

func (r *LeagueRepository) List(limit, offset int, sportID uint, activeOnly bool) ([]*models.League, int64, error) {
	var leagues []*models.League
	var total int64

	query := r.db.Model(&models.League{})
	if sportID > 0 {
		query = query.Where("sport_id = ?", sportID)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Sport").Order("name ASC").Limit(limit).Offset(offset).Find(&leagues).Error; err != nil {
		return nil, 0, err
	}

	return leagues, total, nil
}

func (r *LeagueRepository) GetByExternalID(externalID string) (*models.League, error) {
	if externalID == "" {
		return nil, errors.New("external ID cannot be empty")
	}
	var league models.League
	result := r.db.Where("external_id = ?", externalID).First(&league)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("league not found")
		}
		return nil, result.Error
	}
	return &league, nil
}

func (r *LeagueRepository) Upsert(league *models.League) error {
	if league == nil {
		return errors.New("league cannot be nil")
	}
	if league.ExternalID == "" {
		return r.Create(league)
	}
	existing, err := r.GetByExternalID(league.ExternalID)
	if err == nil {
		league.ID = existing.ID
		return r.Update(league)
	}
	return r.Create(league)
}

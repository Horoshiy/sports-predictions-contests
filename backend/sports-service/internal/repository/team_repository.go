package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type TeamRepositoryInterface interface {
	Create(team *models.Team) error
	GetByID(id uint) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uint) error
	List(limit, offset int, sportID uint, activeOnly bool) ([]*models.Team, int64, error)
}

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepositoryInterface {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) error {
	if team == nil {
		return errors.New("team cannot be nil")
	}
	return r.db.Create(team).Error
}

func (r *TeamRepository) GetByID(id uint) (*models.Team, error) {
	if id == 0 {
		return nil, errors.New("invalid team ID")
	}
	var team models.Team
	result := r.db.Preload("Sport").First(&team, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, result.Error
	}
	return &team, nil
}

func (r *TeamRepository) Update(team *models.Team) error {
	if team == nil {
		return errors.New("team cannot be nil")
	}
	if team.ID == 0 {
		return errors.New("team ID cannot be zero")
	}
	result := r.db.Save(team)
	if result.RowsAffected == 0 {
		return errors.New("team not found")
	}
	return result.Error
}

func (r *TeamRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid team ID")
	}
	result := r.db.Delete(&models.Team{}, id)
	if result.RowsAffected == 0 {
		return errors.New("team not found")
	}
	return result.Error
}

func (r *TeamRepository) List(limit, offset int, sportID uint, activeOnly bool) ([]*models.Team, int64, error) {
	var teams []*models.Team
	var total int64

	query := r.db.Model(&models.Team{})
	if sportID > 0 {
		query = query.Where("sport_id = ?", sportID)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Sport").Order("name ASC").Limit(limit).Offset(offset).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

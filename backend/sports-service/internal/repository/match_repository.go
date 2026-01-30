package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type MatchRepositoryInterface interface {
	Create(match *models.Match) error
	GetByID(id uint) (*models.Match, error)
	GetByExternalID(externalID string) (*models.Match, error)
	Update(match *models.Match) error
	Upsert(match *models.Match) error
	Delete(id uint) error
	List(limit, offset int, leagueID, teamID uint, status string) ([]*models.Match, int64, error)
}

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepositoryInterface {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) Create(match *models.Match) error {
	if match == nil {
		return errors.New("match cannot be nil")
	}
	return r.db.Create(match).Error
}

func (r *MatchRepository) GetByID(id uint) (*models.Match, error) {
	if id == 0 {
		return nil, errors.New("invalid match ID")
	}
	var match models.Match
	result := r.db.Preload("League").Preload("HomeTeam").Preload("AwayTeam").First(&match, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("match not found")
		}
		return nil, result.Error
	}
	return &match, nil
}

func (r *MatchRepository) Update(match *models.Match) error {
	if match == nil {
		return errors.New("match cannot be nil")
	}
	if match.ID == 0 {
		return errors.New("match ID cannot be zero")
	}
	result := r.db.Omit("League", "HomeTeam", "AwayTeam").Model(match).Select("*").Updates(match)
	if result.RowsAffected == 0 {
		return errors.New("match not found")
	}
	return result.Error
}

func (r *MatchRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid match ID")
	}
	result := r.db.Delete(&models.Match{}, id)
	if result.RowsAffected == 0 {
		return errors.New("match not found")
	}
	return result.Error
}

func (r *MatchRepository) List(limit, offset int, leagueID, teamID uint, status string) ([]*models.Match, int64, error) {
	var matches []*models.Match
	var total int64

	query := r.db.Model(&models.Match{})
	if leagueID > 0 {
		query = query.Where("league_id = ?", leagueID)
	}
	if teamID > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("League").Preload("HomeTeam").Preload("AwayTeam").
		Order("scheduled_at DESC").Limit(limit).Offset(offset).Find(&matches).Error; err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}

func (r *MatchRepository) GetByExternalID(externalID string) (*models.Match, error) {
	if externalID == "" {
		return nil, errors.New("external ID cannot be empty")
	}
	var match models.Match
	result := r.db.Where("external_id = ?", externalID).First(&match)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("match not found")
		}
		return nil, result.Error
	}
	return &match, nil
}

func (r *MatchRepository) Upsert(match *models.Match) error {
	if match == nil {
		return errors.New("match cannot be nil")
	}
	if match.ExternalID == "" {
		return r.Create(match)
	}
	existing, err := r.GetByExternalID(match.ExternalID)
	if err == nil {
		match.ID = existing.ID
		return r.Update(match)
	}
	return r.Create(match)
}

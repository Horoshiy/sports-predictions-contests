package repository

import (
	"errors"

	"github.com/sports-prediction-contests/contest-service/internal/models"
	"gorm.io/gorm"
)

// ContestRepositoryInterface defines the contract for contest repository
type ContestRepositoryInterface interface {
	Create(contest *models.Contest) error
	GetByID(id uint) (*models.Contest, error)
	Update(contest *models.Contest) error
	Delete(id uint) error
	List(limit, offset int, status, sportType string) ([]*models.Contest, int64, error)
	GetByCreatorID(creatorID uint) ([]*models.Contest, error)
}

// ParticipantRepositoryInterface defines the contract for participant repository
type ParticipantRepositoryInterface interface {
	Create(participant *models.Participant) error
	GetByID(id uint) (*models.Participant, error)
	GetByContestAndUser(contestID, userID uint) (*models.Participant, error)
	Update(participant *models.Participant) error
	Delete(id uint) error
	ListByContest(contestID uint, limit, offset int) ([]*models.Participant, int64, error)
	CountByContest(contestID uint) (int64, error)
	DeleteByContestAndUser(contestID, userID uint) error
}

// ContestRepository implements ContestRepositoryInterface
type ContestRepository struct {
	db *gorm.DB
}

// ParticipantRepository implements ParticipantRepositoryInterface
type ParticipantRepository struct {
	db *gorm.DB
}

// NewContestRepository creates a new contest repository instance
func NewContestRepository(db *gorm.DB) ContestRepositoryInterface {
	return &ContestRepository{db: db}
}

// NewParticipantRepository creates a new participant repository instance
func NewParticipantRepository(db *gorm.DB) ParticipantRepositoryInterface {
	return &ParticipantRepository{db: db}
}

// Contest Repository Methods

// Create creates a new contest in the database
func (r *ContestRepository) Create(contest *models.Contest) error {
	if contest == nil {
		return errors.New("contest cannot be nil")
	}

	result := r.db.Create(contest)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetByID retrieves a contest by ID
func (r *ContestRepository) GetByID(id uint) (*models.Contest, error) {
	if id == 0 {
		return nil, errors.New("invalid contest ID")
	}

	var contest models.Contest
	result := r.db.First(&contest, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("contest not found")
		}
		return nil, result.Error
	}

	return &contest, nil
}

// Update updates an existing contest
func (r *ContestRepository) Update(contest *models.Contest) error {
	if contest == nil {
		return errors.New("contest cannot be nil")
	}

	if contest.ID == 0 {
		return errors.New("contest ID cannot be zero")
	}

	result := r.db.Save(contest)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("contest not found")
	}

	return nil
}

// Delete deletes a contest by ID
func (r *ContestRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid contest ID")
	}

	// Start transaction to delete contest and its participants
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete all participants first
	if err := tx.Where("contest_id = ?", id).Delete(&models.Participant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the contest
	result := tx.Delete(&models.Contest{}, id)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("contest not found")
	}

	return tx.Commit().Error
}

// List retrieves contests with pagination and filters
func (r *ContestRepository) List(limit, offset int, status, sportType string) ([]*models.Contest, int64, error) {
	var contests []*models.Contest
	var total int64

	query := r.db.Model(&models.Contest{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if sportType != "" {
		query = query.Where("sport_type = ?", sportType)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&contests).Error; err != nil {
		return nil, 0, err
	}

	return contests, total, nil
}

// GetByCreatorID retrieves contests created by a specific user
func (r *ContestRepository) GetByCreatorID(creatorID uint) ([]*models.Contest, error) {
	if creatorID == 0 {
		return nil, errors.New("invalid creator ID")
	}

	var contests []*models.Contest
	result := r.db.Where("creator_id = ?", creatorID).Order("created_at DESC").Find(&contests)
	if result.Error != nil {
		return nil, result.Error
	}

	return contests, nil
}

// Participant Repository Methods

// Create creates a new participant in the database
func (r *ParticipantRepository) Create(participant *models.Participant) error {
	if participant == nil {
		return errors.New("participant cannot be nil")
	}

	result := r.db.Create(participant)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetByID retrieves a participant by ID
func (r *ParticipantRepository) GetByID(id uint) (*models.Participant, error) {
	if id == 0 {
		return nil, errors.New("invalid participant ID")
	}

	var participant models.Participant
	result := r.db.Preload("Contest").First(&participant, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("participant not found")
		}
		return nil, result.Error
	}

	return &participant, nil
}

// GetByContestAndUser retrieves a participant by contest and user ID
func (r *ParticipantRepository) GetByContestAndUser(contestID, userID uint) (*models.Participant, error) {
	if contestID == 0 || userID == 0 {
		return nil, errors.New("invalid contest or user ID")
	}

	var participant models.Participant
	result := r.db.Where("contest_id = ? AND user_id = ?", contestID, userID).First(&participant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("participant not found")
		}
		return nil, result.Error
	}

	return &participant, nil
}

// Update updates an existing participant
func (r *ParticipantRepository) Update(participant *models.Participant) error {
	if participant == nil {
		return errors.New("participant cannot be nil")
	}

	if participant.ID == 0 {
		return errors.New("participant ID cannot be zero")
	}

	result := r.db.Save(participant)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("participant not found")
	}

	return nil
}

// Delete deletes a participant by ID
func (r *ParticipantRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid participant ID")
	}

	result := r.db.Delete(&models.Participant{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("participant not found")
	}

	return nil
}

// ListByContest retrieves participants for a specific contest with pagination
func (r *ParticipantRepository) ListByContest(contestID uint, limit, offset int) ([]*models.Participant, int64, error) {
	if contestID == 0 {
		return nil, 0, errors.New("invalid contest ID")
	}

	var participants []*models.Participant
	var total int64

	query := r.db.Where("contest_id = ?", contestID)

	// Count total records
	if err := query.Model(&models.Participant{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := query.Order("joined_at ASC").Limit(limit).Offset(offset).Find(&participants).Error; err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

// CountByContest counts participants for a specific contest
func (r *ParticipantRepository) CountByContest(contestID uint) (int64, error) {
	if contestID == 0 {
		return 0, errors.New("invalid contest ID")
	}

	var count int64
	result := r.db.Model(&models.Participant{}).Where("contest_id = ? AND status = ?", contestID, "active").Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// DeleteByContestAndUser deletes a participant by contest and user ID
func (r *ParticipantRepository) DeleteByContestAndUser(contestID, userID uint) error {
	if contestID == 0 || userID == 0 {
		return errors.New("invalid contest or user ID")
	}

	result := r.db.Where("contest_id = ? AND user_id = ?", contestID, userID).Delete(&models.Participant{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("participant not found")
	}

	return nil
}

package repository

import (
	"errors"

	"github.com/sports-prediction-contests/challenge-service/internal/models"
	"gorm.io/gorm"
)

// ChallengeRepositoryInterface defines the contract for challenge repository
type ChallengeRepositoryInterface interface {
	Create(challenge *models.Challenge) error
	CreateWithParticipants(challenge *models.Challenge, participants []*models.ChallengeParticipant) error
	GetByID(id uint) (*models.Challenge, error)
	Update(challenge *models.Challenge) error
	Delete(id uint) error
	ListByUser(userID uint, status string, limit, offset int) ([]*models.Challenge, int64, error)
	ListByEvent(eventID uint, limit, offset int) ([]*models.Challenge, int64, error)
	GetExpiredChallenges() ([]*models.Challenge, error)
	UpdateStatus(id uint, status string) error
}

// ChallengeParticipantRepositoryInterface defines the contract for challenge participant repository
type ChallengeParticipantRepositoryInterface interface {
	Create(participant *models.ChallengeParticipant) error
	GetByID(id uint) (*models.ChallengeParticipant, error)
	GetByChallengeAndUser(challengeID, userID uint) (*models.ChallengeParticipant, error)
	Update(participant *models.ChallengeParticipant) error
	Delete(id uint) error
	ListByChallenge(challengeID uint) ([]*models.ChallengeParticipant, error)
}

// ChallengeRepository implements ChallengeRepositoryInterface
type ChallengeRepository struct {
	db *gorm.DB
}

// ChallengeParticipantRepository implements ChallengeParticipantRepositoryInterface
type ChallengeParticipantRepository struct {
	db *gorm.DB
}

// NewChallengeRepository creates a new challenge repository instance
func NewChallengeRepository(db *gorm.DB) ChallengeRepositoryInterface {
	return &ChallengeRepository{db: db}
}

// NewChallengeParticipantRepository creates a new challenge participant repository instance
func NewChallengeParticipantRepository(db *gorm.DB) ChallengeParticipantRepositoryInterface {
	return &ChallengeParticipantRepository{db: db}
}

// Challenge Repository Methods

// Create creates a new challenge in the database
func (r *ChallengeRepository) Create(challenge *models.Challenge) error {
	if challenge == nil {
		return errors.New("challenge cannot be nil")
	}

	result := r.db.Create(challenge)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CreateWithParticipants creates a challenge and its participants in a single transaction
func (r *ChallengeRepository) CreateWithParticipants(challenge *models.Challenge, participants []*models.ChallengeParticipant) error {
	if challenge == nil {
		return errors.New("challenge cannot be nil")
	}
	if len(participants) == 0 {
		return errors.New("participants cannot be empty")
	}

	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			if tx.Error == nil {
				tx.Rollback()
			}
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// Create challenge
	if err := tx.Create(challenge).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update participant challenge IDs and create them
	for _, participant := range participants {
		participant.ChallengeID = challenge.ID
		if err := tx.Create(participant).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

// GetByID retrieves a challenge by ID
func (r *ChallengeRepository) GetByID(id uint) (*models.Challenge, error) {
	if id == 0 {
		return nil, errors.New("invalid challenge ID")
	}

	var challenge models.Challenge
	result := r.db.First(&challenge, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("challenge not found")
		}
		return nil, result.Error
	}

	return &challenge, nil
}

// Update updates an existing challenge
func (r *ChallengeRepository) Update(challenge *models.Challenge) error {
	if challenge == nil {
		return errors.New("challenge cannot be nil")
	}

	if challenge.ID == 0 {
		return errors.New("challenge ID cannot be zero")
	}

	result := r.db.Save(challenge)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("challenge not found")
	}

	return nil
}

// Delete deletes a challenge by ID
func (r *ChallengeRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid challenge ID")
	}

	// Start transaction to delete challenge and its participants
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			if tx.Error == nil {
				tx.Rollback()
			}
		}
	}()

	// Delete all participants first
	if err := tx.Where("challenge_id = ?", id).Delete(&models.ChallengeParticipant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the challenge
	result := tx.Delete(&models.Challenge{}, id)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("challenge not found")
	}

	return tx.Commit().Error
}

// ListByUser retrieves challenges for a specific user with pagination and filters
func (r *ChallengeRepository) ListByUser(userID uint, status string, limit, offset int) ([]*models.Challenge, int64, error) {
	if userID == 0 {
		return nil, 0, errors.New("invalid user ID")
	}

	var challenges []*models.Challenge
	var total int64

	query := r.db.Model(&models.Challenge{}).Where("challenger_id = ? OR opponent_id = ?", userID, userID)

	// Apply status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&challenges).Error; err != nil {
		return nil, 0, err
	}

	return challenges, total, nil
}

// ListByEvent retrieves challenges for a specific event with pagination
func (r *ChallengeRepository) ListByEvent(eventID uint, limit, offset int) ([]*models.Challenge, int64, error) {
	if eventID == 0 {
		return nil, 0, errors.New("invalid event ID")
	}

	var challenges []*models.Challenge
	var total int64

	query := r.db.Model(&models.Challenge{}).Where("event_id = ? AND status = ?", eventID, "pending")

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&challenges).Error; err != nil {
		return nil, 0, err
	}

	return challenges, total, nil
}

// GetExpiredChallenges retrieves challenges that have expired
func (r *ChallengeRepository) GetExpiredChallenges() ([]*models.Challenge, error) {
	var challenges []*models.Challenge
	result := r.db.Where("status = ? AND expires_at < NOW()", "pending").Find(&challenges)
	if result.Error != nil {
		return nil, result.Error
	}

	return challenges, nil
}

// UpdateStatus updates the status of a challenge
func (r *ChallengeRepository) UpdateStatus(id uint, status string) error {
	if id == 0 {
		return errors.New("invalid challenge ID")
	}

	result := r.db.Model(&models.Challenge{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("challenge not found")
	}

	return nil
}

// Challenge Participant Repository Methods

// Create creates a new challenge participant in the database
func (r *ChallengeParticipantRepository) Create(participant *models.ChallengeParticipant) error {
	if participant == nil {
		return errors.New("challenge participant cannot be nil")
	}

	result := r.db.Create(participant)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetByID retrieves a challenge participant by ID
func (r *ChallengeParticipantRepository) GetByID(id uint) (*models.ChallengeParticipant, error) {
	if id == 0 {
		return nil, errors.New("invalid challenge participant ID")
	}

	var participant models.ChallengeParticipant
	result := r.db.Preload("Challenge").First(&participant, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("challenge participant not found")
		}
		return nil, result.Error
	}

	return &participant, nil
}

// GetByChallengeAndUser retrieves a challenge participant by challenge and user ID
func (r *ChallengeParticipantRepository) GetByChallengeAndUser(challengeID, userID uint) (*models.ChallengeParticipant, error) {
	if challengeID == 0 || userID == 0 {
		return nil, errors.New("invalid challenge or user ID")
	}

	var participant models.ChallengeParticipant
	result := r.db.Where("challenge_id = ? AND user_id = ?", challengeID, userID).First(&participant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("challenge participant not found")
		}
		return nil, result.Error
	}

	return &participant, nil
}

// Update updates an existing challenge participant
func (r *ChallengeParticipantRepository) Update(participant *models.ChallengeParticipant) error {
	if participant == nil {
		return errors.New("challenge participant cannot be nil")
	}

	if participant.ID == 0 {
		return errors.New("challenge participant ID cannot be zero")
	}

	result := r.db.Save(participant)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("challenge participant not found")
	}

	return nil
}

// Delete deletes a challenge participant by ID
func (r *ChallengeParticipantRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid challenge participant ID")
	}

	result := r.db.Delete(&models.ChallengeParticipant{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("challenge participant not found")
	}

	return nil
}

// ListByChallenge retrieves participants for a specific challenge
func (r *ChallengeParticipantRepository) ListByChallenge(challengeID uint) ([]*models.ChallengeParticipant, error) {
	if challengeID == 0 {
		return nil, errors.New("invalid challenge ID")
	}

	var participants []*models.ChallengeParticipant
	result := r.db.Where("challenge_id = ?", challengeID).Order("joined_at ASC").Find(&participants)
	if result.Error != nil {
		return nil, result.Error
	}

	return participants, nil
}

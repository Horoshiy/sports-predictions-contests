package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Score represents a user's score for a specific prediction in a contest
type Score struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;uniqueIndex:idx_user_contest_prediction" json:"user_id"`
	ContestID       uint      `gorm:"not null;uniqueIndex:idx_user_contest_prediction;index:idx_contest_scores" json:"contest_id"`
	PredictionID    uint      `gorm:"not null;uniqueIndex:idx_user_contest_prediction" json:"prediction_id"`
	Points          float64   `gorm:"not null;default:0" json:"points"`
	TimeCoefficient float64   `gorm:"not null;default:1.0" json:"time_coefficient"`
	ScoredAt        time.Time `gorm:"not null" json:"scored_at"`
	gorm.Model
}

// ValidateUserID checks if the user ID is valid
func (s *Score) ValidateUserID() error {
	if s.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

// ValidateContestID checks if the contest ID is valid
func (s *Score) ValidateContestID() error {
	if s.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}
	return nil
}

// ValidatePredictionID checks if the prediction ID is valid
func (s *Score) ValidatePredictionID() error {
	if s.PredictionID == 0 {
		return errors.New("prediction ID cannot be empty")
	}
	return nil
}

// ValidatePoints checks if the points value is valid
func (s *Score) ValidatePoints() error {
	if s.Points < 0 {
		return errors.New("points cannot be negative")
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a score
func (s *Score) BeforeCreate(tx *gorm.DB) error {
	// Set scored time if not provided
	if s.ScoredAt.IsZero() {
		s.ScoredAt = time.Now().UTC()
	}

	// Validate fields
	if err := s.ValidateUserID(); err != nil {
		return err
	}

	if err := s.ValidateContestID(); err != nil {
		return err
	}

	if err := s.ValidatePredictionID(); err != nil {
		return err
	}

	if err := s.ValidatePoints(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a score
func (s *Score) BeforeUpdate(tx *gorm.DB) error {
	// Don't allow changing ScoredAt after creation
	if !s.ScoredAt.IsZero() {
		var existing Score
		if err := tx.Select("scored_at").First(&existing, s.ID).Error; err == nil {
			s.ScoredAt = existing.ScoredAt // Preserve original ScoredAt
		}
	}

	// Validate points on update
	if err := s.ValidatePoints(); err != nil {
		return err
	}

	return nil
}

// IsPositive checks if the score has positive points
func (s *Score) IsPositive() bool {
	return s.Points > 0
}

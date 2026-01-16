package models

import (
	"errors"

	"gorm.io/gorm"
)

// UserStreak tracks consecutive successful predictions for a user in a contest
type UserStreak struct {
	gorm.Model
	UserID                uint  `gorm:"not null;uniqueIndex:idx_user_contest_streak" json:"user_id"`
	ContestID             uint  `gorm:"not null;uniqueIndex:idx_user_contest_streak" json:"contest_id"`
	CurrentStreak         uint  `gorm:"not null;default:0" json:"current_streak"`
	MaxStreak             uint  `gorm:"not null;default:0" json:"max_streak"`
	LastPredictionID      *uint `json:"last_prediction_id"`
	LastPredictionCorrect *bool `json:"last_prediction_correct"`
}

// ValidateUserID checks if the user ID is valid
func (s *UserStreak) ValidateUserID() error {
	if s.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

// ValidateContestID checks if the contest ID is valid
func (s *UserStreak) ValidateContestID() error {
	if s.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a streak
func (s *UserStreak) BeforeCreate(tx *gorm.DB) error {
	if err := s.ValidateUserID(); err != nil {
		return err
	}
	if err := s.ValidateContestID(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a streak
func (s *UserStreak) BeforeUpdate(tx *gorm.DB) error {
	// GORM automatically updates UpdatedAt via gorm.Model
	return nil
}

// GetMultiplier returns the point multiplier based on current streak
func (s *UserStreak) GetMultiplier() float64 {
	switch {
	case s.CurrentStreak >= 10:
		return 2.0
	case s.CurrentStreak >= 7:
		return 1.75
	case s.CurrentStreak >= 5:
		return 1.5
	case s.CurrentStreak >= 3:
		return 1.25
	default:
		return 1.0
	}
}

// IncrementStreak increases the current streak and updates max if needed
func (s *UserStreak) IncrementStreak(predictionID uint) {
	s.CurrentStreak++
	if s.CurrentStreak > s.MaxStreak {
		s.MaxStreak = s.CurrentStreak
	}
	s.LastPredictionID = &predictionID
	correct := true
	s.LastPredictionCorrect = &correct
}

// ResetStreak sets the current streak to 0
func (s *UserStreak) ResetStreak(predictionID uint) {
	s.CurrentStreak = 0
	s.LastPredictionID = &predictionID
	correct := false
	s.LastPredictionCorrect = &correct
}

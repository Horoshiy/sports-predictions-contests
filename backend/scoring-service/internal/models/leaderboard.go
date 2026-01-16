package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Leaderboard represents a user's position and total points in a contest
type Leaderboard struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ContestID     uint      `gorm:"not null;uniqueIndex:idx_contest_user;index:idx_contest_rank" json:"contest_id"`
	UserID        uint      `gorm:"not null;uniqueIndex:idx_contest_user" json:"user_id"`
	TotalPoints   float64   `gorm:"not null;default:0" json:"total_points"`
	Rank          uint      `gorm:"not null;default:0;index:idx_contest_rank" json:"rank"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
	CurrentStreak uint      `gorm:"-" json:"current_streak"`
	MaxStreak     uint      `gorm:"-" json:"max_streak"`
	Multiplier    float64   `gorm:"-" json:"multiplier"`
	gorm.Model
}

// ValidateContestID checks if the contest ID is valid
func (l *Leaderboard) ValidateContestID() error {
	if l.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}
	return nil
}

// ValidateUserID checks if the user ID is valid
func (l *Leaderboard) ValidateUserID() error {
	if l.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

// ValidateTotalPoints checks if the total points value is valid
func (l *Leaderboard) ValidateTotalPoints() error {
	if l.TotalPoints < 0 {
		return errors.New("total points cannot be negative")
	}
	return nil
}

// ValidateRank checks if the rank value is valid
func (l *Leaderboard) ValidateRank() error {
	if l.Rank == 0 {
		return errors.New("rank cannot be zero")
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a leaderboard entry
func (l *Leaderboard) BeforeCreate(tx *gorm.DB) error {
	// Set updated time if not provided
	if l.UpdatedAt.IsZero() {
		l.UpdatedAt = time.Now().UTC()
	}

	// Validate fields
	if err := l.ValidateContestID(); err != nil {
		return err
	}

	if err := l.ValidateUserID(); err != nil {
		return err
	}

	if err := l.ValidateTotalPoints(); err != nil {
		return err
	}

	// Don't validate rank on create as it might be set later
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a leaderboard entry
func (l *Leaderboard) BeforeUpdate(tx *gorm.DB) error {
	// Update the timestamp
	l.UpdatedAt = time.Now().UTC()

	// Validate fields
	if err := l.ValidateTotalPoints(); err != nil {
		return err
	}

	return nil
}

// HasPoints checks if the leaderboard entry has any points
func (l *Leaderboard) HasPoints() bool {
	return l.TotalPoints > 0
}

// IsRanked checks if the leaderboard entry has a valid rank
func (l *Leaderboard) IsRanked() bool {
	return l.Rank > 0
}

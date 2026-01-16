package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TeamContestEntry struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TeamID      uint      `gorm:"not null;index" json:"team_id"`
	ContestID   uint      `gorm:"not null;index" json:"contest_id"`
	TotalPoints float64   `gorm:"default:0" json:"total_points"`
	Rank        uint      `gorm:"default:0" json:"rank"`
	JoinedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	gorm.Model

	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

func (TeamContestEntry) TableName() string {
	return "user_team_contest_entries"
}

func (e *TeamContestEntry) BeforeCreate(tx *gorm.DB) error {
	if e.TeamID == 0 {
		return errors.New("team ID cannot be empty")
	}
	if e.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}
	if e.JoinedAt.IsZero() {
		e.JoinedAt = time.Now()
	}

	var existing TeamContestEntry
	if tx.Where("team_id = ? AND contest_id = ?", e.TeamID, e.ContestID).First(&existing).Error == nil {
		return errors.New("team is already participating in this contest")
	}
	return nil
}

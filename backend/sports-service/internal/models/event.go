package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents a sports event in user-friendly format
type Event struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"not null" json:"title"`
	SportType  string         `gorm:"not null;index" json:"sport_type"`
	HomeTeam   string         `gorm:"not null" json:"home_team"`
	AwayTeam   string         `gorm:"not null" json:"away_team"`
	EventDate  time.Time      `gorm:"not null;index" json:"event_date"`
	Status     string         `gorm:"not null;default:'scheduled';index" json:"status"`
	ResultData string         `gorm:"type:jsonb" json:"result_data"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Event
func (Event) TableName() string {
	return "events"
}

// IsCompleted checks if event is completed
func (e *Event) IsCompleted() bool {
	return e.Status == "completed"
}

// IsScheduled checks if event is scheduled
func (e *Event) IsScheduled() bool {
	return e.Status == "scheduled"
}

// IsLive checks if event is live
func (e *Event) IsLive() bool {
	return e.Status == "live"
}

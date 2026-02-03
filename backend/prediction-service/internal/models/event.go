package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Event represents a sports event that can be predicted
type Event struct {
	gorm.Model
	Title      string    `gorm:"not null" json:"title"`
	SportType  string    `gorm:"not null;index" json:"sport_type"`
	HomeTeam   string    `gorm:"not null" json:"home_team"`
	AwayTeam   string    `gorm:"not null" json:"away_team"`
	EventDate  time.Time `gorm:"not null;index" json:"event_date"`
	Status     string    `gorm:"not null;default:'scheduled';index" json:"status"` // "scheduled", "live", "completed", "cancelled"
	ResultData string    `gorm:"type:jsonb" json:"result_data"` // JSON string for event results
}

// ValidateTitle checks if the title is valid
func (e *Event) ValidateTitle() error {
	if len(strings.TrimSpace(e.Title)) == 0 {
		return errors.New("title cannot be empty")
	}

	if len(e.Title) > 200 {
		return errors.New("title cannot exceed 200 characters")
	}

	return nil
}

// ValidateSportType checks if the sport type is valid
func (e *Event) ValidateSportType() error {
	if len(strings.TrimSpace(e.SportType)) == 0 {
		return errors.New("sport type cannot be empty")
	}

	return nil
}

// ValidateTeams checks if the team names are valid
func (e *Event) ValidateTeams() error {
	if len(strings.TrimSpace(e.HomeTeam)) == 0 {
		return errors.New("home team cannot be empty")
	}

	if len(strings.TrimSpace(e.AwayTeam)) == 0 {
		return errors.New("away team cannot be empty")
	}

	if len(e.HomeTeam) > 100 {
		return errors.New("home team name cannot exceed 100 characters")
	}

	if len(e.AwayTeam) > 100 {
		return errors.New("away team name cannot exceed 100 characters")
	}

	return nil
}

// ValidateEventDate checks if the event date is valid
func (e *Event) ValidateEventDate() error {
	if e.EventDate.IsZero() {
		return errors.New("event date cannot be empty")
	}

	// Use UTC for consistent timezone handling
	now := time.Now().UTC()
	// Allow events to be created up to 1 hour in the past for flexibility
	if e.EventDate.UTC().Before(now.Add(-1 * time.Hour)) {
		return errors.New("event date cannot be more than 1 hour in the past")
	}

	return nil
}

// ValidateStatus checks if the status is valid
func (e *Event) ValidateStatus() error {
	validStatuses := []string{"scheduled", "live", "completed", "cancelled"}
	for _, validStatus := range validStatuses {
		if e.Status == validStatus {
			return nil
		}
	}
	return errors.New("invalid status")
}

// ValidateResultData checks if the result data is valid
func (e *Event) ValidateResultData() error {
	if len(e.ResultData) > 5000 {
		return errors.New("result data cannot exceed 5000 characters")
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating an event
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not provided
	if e.Status == "" {
		e.Status = "scheduled"
	}

	// Validate fields
	if err := e.ValidateTitle(); err != nil {
		return err
	}

	if err := e.ValidateSportType(); err != nil {
		return err
	}

	if err := e.ValidateTeams(); err != nil {
		return err
	}

	if err := e.ValidateEventDate(); err != nil {
		return err
	}

	if err := e.ValidateStatus(); err != nil {
		return err
	}

	if err := e.ValidateResultData(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating an event
func (e *Event) BeforeUpdate(tx *gorm.DB) error {
	return e.BeforeCreate(tx)
}

// IsScheduled checks if the event is scheduled
func (e *Event) IsScheduled() bool {
	return e.Status == "scheduled"
}

// IsLive checks if the event is currently live
func (e *Event) IsLive() bool {
	return e.Status == "live"
}

// IsCompleted checks if the event is completed
func (e *Event) IsCompleted() bool {
	return e.Status == "completed"
}

// CanAcceptPredictions checks if the event can accept predictions
func (e *Event) CanAcceptPredictions() bool {
	return e.Status == "scheduled" && time.Now().UTC().Before(e.EventDate.UTC())
}

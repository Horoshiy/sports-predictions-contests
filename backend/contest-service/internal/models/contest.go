package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Contest represents a sports prediction contest
type Contest struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Title              string    `gorm:"not null" json:"title"`
	Description        string    `json:"description"`
	SportType          string    `gorm:"not null" json:"sport_type"`
	Rules              string    `gorm:"type:text" json:"rules"` // JSON string for flexible rule configuration
	PredictionSchema   []byte    `gorm:"type:jsonb" json:"prediction_schema"`
	Status             string    `gorm:"not null;default:'draft'" json:"status"` // "draft", "active", "completed", "cancelled"
	StartDate          time.Time `gorm:"not null" json:"start_date"`
	EndDate            time.Time `gorm:"not null" json:"end_date"`
	MaxParticipants    uint      `gorm:"default:0" json:"max_participants"` // 0 means unlimited
	CurrentParticipants uint     `gorm:"default:0" json:"current_participants"`
	CreatorID          uint      `gorm:"not null" json:"creator_id"`
	gorm.Model
}

// ValidateTitle checks if the title is valid
func (c *Contest) ValidateTitle() error {
	if len(strings.TrimSpace(c.Title)) == 0 {
		return errors.New("title cannot be empty")
	}

	if len(c.Title) > 200 {
		return errors.New("title cannot exceed 200 characters")
	}

	return nil
}

// ValidateDescription checks if the description is valid
func (c *Contest) ValidateDescription() error {
	if len(c.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	return nil
}

// ValidateSportType checks if the sport type is valid
func (c *Contest) ValidateSportType() error {
	if len(strings.TrimSpace(c.SportType)) == 0 {
		return errors.New("sport type cannot be empty")
	}

	// Allow any non-empty sport type for extensibility
	// Business logic validation can be added at the service layer if needed
	return nil
}

// ValidateStatus checks if the status is valid
func (c *Contest) ValidateStatus() error {
	validStatuses := []string{"draft", "active", "completed", "cancelled"}
	for _, validStatus := range validStatuses {
		if c.Status == validStatus {
			return nil
		}
	}

	return errors.New("invalid status")
}

// ValidateDates checks if the dates are valid
func (c *Contest) ValidateDates() error {
	if c.StartDate.IsZero() {
		return errors.New("start date cannot be empty")
	}

	if c.EndDate.IsZero() {
		return errors.New("end date cannot be empty")
	}

	if c.EndDate.Before(c.StartDate) {
		return errors.New("end date must be after start date")
	}

	// Use UTC for consistent timezone handling
	now := time.Now().UTC()
	if c.StartDate.UTC().Before(now.Add(-24 * time.Hour)) {
		return errors.New("start date cannot be more than 24 hours in the past")
	}

	return nil
}

// ValidateMaxParticipants checks if the max participants is valid
func (c *Contest) ValidateMaxParticipants() error {
	if c.MaxParticipants > 10000 {
		return errors.New("max participants cannot exceed 10000")
	}

	return nil
}

// BeforeCreate is a GORM hook that runs before creating a contest
func (c *Contest) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not provided
	if c.Status == "" {
		c.Status = "draft"
	}

	// Validate fields
	if err := c.ValidateTitle(); err != nil {
		return err
	}

	if err := c.ValidateDescription(); err != nil {
		return err
	}

	if err := c.ValidateSportType(); err != nil {
		return err
	}

	if err := c.ValidateStatus(); err != nil {
		return err
	}

	if err := c.ValidateDates(); err != nil {
		return err
	}

	if err := c.ValidateMaxParticipants(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a contest
func (c *Contest) BeforeUpdate(tx *gorm.DB) error {
	return c.BeforeCreate(tx)
}

// CanJoin checks if a user can join the contest
func (c *Contest) CanJoin() bool {
	if c.Status != "active" {
		return false
	}

	if c.MaxParticipants > 0 && c.CurrentParticipants >= c.MaxParticipants {
		return false
	}

	if time.Now().After(c.EndDate) {
		return false
	}

	return true
}

// IsActive checks if the contest is currently active
func (c *Contest) IsActive() bool {
	return c.Status == "active" && time.Now().Before(c.EndDate) && time.Now().After(c.StartDate)
}

// GetComputedStatus returns the actual status based on current time and dates
// This should be used for display purposes instead of the Status field
func (c *Contest) GetComputedStatus() string {
	now := time.Now().UTC()
	startUTC := c.StartDate.UTC()
	endUTC := c.EndDate.UTC()

	// If manually cancelled or draft, respect that status
	if c.Status == "cancelled" || c.Status == "draft" {
		return c.Status
	}

	// Compute status based on dates
	if now.Before(startUTC) {
		return "upcoming"
	}
	if now.After(endUTC) {
		return "completed"
	}
	return "active"
}

package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Prediction represents a user's prediction for a sports event
type Prediction struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ContestID      uint      `gorm:"not null;index:idx_contest_user;uniqueIndex:idx_user_contest_event" json:"contest_id"`
	UserID         uint      `gorm:"not null;index:idx_contest_user;uniqueIndex:idx_user_contest_event" json:"user_id"`
	EventID        uint      `gorm:"not null;uniqueIndex:idx_user_contest_event" json:"event_id"`
	PredictionData string    `gorm:"type:jsonb" json:"prediction_data"` // JSON string for flexible prediction data
	Status         string    `gorm:"not null;default:'pending'" json:"status"` // "pending", "scored", "cancelled"
	SubmittedAt    time.Time `gorm:"not null" json:"submitted_at"`
	gorm.Model

	// Relationships
	Event Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
}

// ValidateContestID checks if the contest ID is valid
func (p *Prediction) ValidateContestID() error {
	if p.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}
	return nil
}

// ValidateUserID checks if the user ID is valid
func (p *Prediction) ValidateUserID() error {
	if p.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

// ValidateEventID checks if the event ID is valid
func (p *Prediction) ValidateEventID() error {
	if p.EventID == 0 {
		return errors.New("event ID cannot be empty")
	}
	return nil
}

// ValidatePredictionData checks if the prediction data is valid
func (p *Prediction) ValidatePredictionData() error {
	if len(strings.TrimSpace(p.PredictionData)) == 0 {
		return errors.New("prediction data cannot be empty")
	}

	if len(p.PredictionData) > 5000 {
		return errors.New("prediction data cannot exceed 5000 characters")
	}

	return nil
}

// ValidateStatus checks if the status is valid
func (p *Prediction) ValidateStatus() error {
	validStatuses := []string{"pending", "scored", "cancelled"}
	for _, validStatus := range validStatuses {
		if p.Status == validStatus {
			return nil
		}
	}
	return errors.New("invalid status")
}

// BeforeCreate is a GORM hook that runs before creating a prediction
func (p *Prediction) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not provided
	if p.Status == "" {
		p.Status = "pending"
	}

	// Set submitted time if not provided
	if p.SubmittedAt.IsZero() {
		p.SubmittedAt = time.Now().UTC()
	}

	// Validate fields
	if err := p.ValidateContestID(); err != nil {
		return err
	}

	if err := p.ValidateUserID(); err != nil {
		return err
	}

	if err := p.ValidateEventID(); err != nil {
		return err
	}

	if err := p.ValidatePredictionData(); err != nil {
		return err
	}

	if err := p.ValidateStatus(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a prediction
func (p *Prediction) BeforeUpdate(tx *gorm.DB) error {
	// Don't reset SubmittedAt on update
	if err := p.ValidatePredictionData(); err != nil {
		return err
	}

	if err := p.ValidateStatus(); err != nil {
		return err
	}

	return nil
}

// CanUpdate checks if a prediction can be updated
func (p *Prediction) CanUpdate() bool {
	return p.Status == "pending"
}

// IsPending checks if the prediction is pending
func (p *Prediction) IsPending() bool {
	return p.Status == "pending"
}

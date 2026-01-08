package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Participant represents a contest participant
type Participant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ContestID uint      `gorm:"not null;index" json:"contest_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Role      string    `gorm:"not null;default:'participant'" json:"role"` // "admin", "participant"
	Status    string    `gorm:"not null;default:'active'" json:"status"`    // "active", "inactive", "banned"
	JoinedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	gorm.Model

	// Relationships
	Contest Contest `gorm:"foreignKey:ContestID" json:"contest,omitempty"`
}

// ValidateRole checks if the role is valid
func (p *Participant) ValidateRole() error {
	validRoles := []string{"admin", "participant"}
	for _, validRole := range validRoles {
		if p.Role == validRole {
			return nil
		}
	}

	return errors.New("invalid role")
}

// ValidateStatus checks if the status is valid
func (p *Participant) ValidateStatus() error {
	validStatuses := []string{"active", "inactive", "banned"}
	for _, validStatus := range validStatuses {
		if p.Status == validStatus {
			return nil
		}
	}

	return errors.New("invalid status")
}

// BeforeCreate is a GORM hook that runs before creating a participant
func (p *Participant) BeforeCreate(tx *gorm.DB) error {
	// Set default values if not provided
	if p.Role == "" {
		p.Role = "participant"
	}

	if p.Status == "" {
		p.Status = "active"
	}

	if p.JoinedAt.IsZero() {
		p.JoinedAt = time.Now()
	}

	// Validate fields
	if err := p.ValidateRole(); err != nil {
		return err
	}

	if err := p.ValidateStatus(); err != nil {
		return err
	}

	if p.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}

	if p.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}

	// Check for duplicate participation
	var existingParticipant Participant
	result := tx.Where("contest_id = ? AND user_id = ?", p.ContestID, p.UserID).First(&existingParticipant)
	if result.Error == nil {
		return errors.New("user is already a participant in this contest")
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a participant
func (p *Participant) BeforeUpdate(tx *gorm.DB) error {
	// Validate fields (but skip duplicate check for updates)
	if err := p.ValidateRole(); err != nil {
		return err
	}

	if err := p.ValidateStatus(); err != nil {
		return err
	}

	return nil
}

// IsActive checks if the participant is active
func (p *Participant) IsActive() bool {
	return p.Status == "active"
}

// IsAdmin checks if the participant is an admin
func (p *Participant) IsAdmin() bool {
	return p.Role == "admin"
}

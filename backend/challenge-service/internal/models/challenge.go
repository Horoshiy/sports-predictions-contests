package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// validStatuses is a map for O(1) status validation lookup
var validStatuses = map[string]bool{
	"pending":   true,
	"accepted":  true,
	"declined":  true,
	"expired":   true,
	"active":    true,
	"completed": true,
}

// Challenge represents a head-to-head challenge between users
type Challenge struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ChallengerID   uint      `gorm:"not null;index" json:"challenger_id"`
	OpponentID     uint      `gorm:"not null;index" json:"opponent_id"`
	EventID        uint      `gorm:"not null;index" json:"event_id"`
	Message        string    `gorm:"type:text" json:"message"`
	Status         string    `gorm:"not null;default:'pending'" json:"status"` // "pending", "accepted", "declined", "expired", "active", "completed"
	ExpiresAt      time.Time `gorm:"not null" json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	WinnerID       *uint     `json:"winner_id"`
	ChallengerScore float64  `gorm:"default:0" json:"challenger_score"`
	OpponentScore  float64   `gorm:"default:0" json:"opponent_score"`
	gorm.Model
}

// ChallengeParticipant represents a participant in a challenge
type ChallengeParticipant struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ChallengeID uint      `gorm:"not null;index" json:"challenge_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Role        string    `gorm:"not null" json:"role"` // "challenger", "opponent"
	Status      string    `gorm:"not null;default:'active'" json:"status"` // "active", "inactive"
	JoinedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	gorm.Model

	// Relationships
	Challenge Challenge `gorm:"foreignKey:ChallengeID" json:"challenge,omitempty"`
}

// ValidateUserIDs checks if the user IDs are valid
func (c *Challenge) ValidateUserIDs() error {
	if c.ChallengerID == 0 {
		return errors.New("challenger ID cannot be empty")
	}
	if c.OpponentID == 0 {
		return errors.New("opponent ID cannot be empty")
	}
	if c.ChallengerID == c.OpponentID {
		return errors.New("challenger and opponent cannot be the same user")
	}
	return nil
}

// ValidateEventID checks if the event ID is valid
func (c *Challenge) ValidateEventID() error {
	if c.EventID == 0 {
		return errors.New("event ID cannot be empty")
	}
	return nil
}

// ValidateStatus checks if the status is valid
func (c *Challenge) ValidateStatus() error {
	if !validStatuses[c.Status] {
		return errors.New("invalid status")
	}
	return nil
}

// ValidateMessage checks if the message is valid
func (c *Challenge) ValidateMessage() error {
	if len(c.Message) > 500 {
		return errors.New("message cannot exceed 500 characters")
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a challenge
func (c *Challenge) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not provided
	if c.Status == "" {
		c.Status = "pending"
	}

	// Set expiration time (24 hours from now)
	if c.ExpiresAt.IsZero() {
		c.ExpiresAt = time.Now().UTC().Add(24 * time.Hour)
	}

	// Validate fields
	if err := c.ValidateUserIDs(); err != nil {
		return err
	}
	if err := c.ValidateEventID(); err != nil {
		return err
	}
	if err := c.ValidateStatus(); err != nil {
		return err
	}
	if err := c.ValidateMessage(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a challenge
func (c *Challenge) BeforeUpdate(tx *gorm.DB) error {
	if err := c.ValidateStatus(); err != nil {
		return err
	}
	return nil
}

// CanAccept checks if the challenge can be accepted
func (c *Challenge) CanAccept() bool {
	return c.Status == "pending" && time.Now().UTC().Before(c.ExpiresAt)
}

// IsExpired checks if the challenge has expired
func (c *Challenge) IsExpired() bool {
	return c.Status == "pending" && time.Now().UTC().After(c.ExpiresAt)
}

// IsActive checks if the challenge is currently active
func (c *Challenge) IsActive() bool {
	return c.Status == "active" || c.Status == "accepted"
}

// IsCompleted checks if the challenge is completed
func (c *Challenge) IsCompleted() bool {
	return c.Status == "completed"
}

// Accept marks the challenge as accepted
func (c *Challenge) Accept() {
	c.Status = "accepted"
	now := time.Now().UTC()
	c.AcceptedAt = &now
}

// Complete marks the challenge as completed with scores
func (c *Challenge) Complete(challengerScore, opponentScore float64) {
	c.Status = "completed"
	c.ChallengerScore = challengerScore
	c.OpponentScore = opponentScore
	now := time.Now().UTC()
	c.CompletedAt = &now

	// Determine winner
	if challengerScore > opponentScore {
		c.WinnerID = &c.ChallengerID
	} else if opponentScore > challengerScore {
		c.WinnerID = &c.OpponentID
	}
	// If scores are equal, WinnerID remains nil (tie)
}

// ValidateUserID checks if the user ID is valid for ChallengeParticipant
func (cp *ChallengeParticipant) ValidateUserID() error {
	if cp.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

// ValidateChallengeID checks if the challenge ID is valid for ChallengeParticipant
func (cp *ChallengeParticipant) ValidateChallengeID() error {
	if cp.ChallengeID == 0 {
		return errors.New("challenge ID cannot be empty")
	}
	return nil
}

// ValidateRole checks if the role is valid for ChallengeParticipant
func (cp *ChallengeParticipant) ValidateRole() error {
	validRoles := []string{"challenger", "opponent"}
	for _, validRole := range validRoles {
		if cp.Role == validRole {
			return nil
		}
	}
	return errors.New("invalid role")
}

// BeforeCreate is a GORM hook that runs before creating a challenge participant
func (cp *ChallengeParticipant) BeforeCreate(tx *gorm.DB) error {
	if err := cp.ValidateUserID(); err != nil {
		return err
	}
	if err := cp.ValidateChallengeID(); err != nil {
		return err
	}
	if err := cp.ValidateRole(); err != nil {
		return err
	}
	return nil
}

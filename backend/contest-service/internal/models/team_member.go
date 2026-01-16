package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TeamMember struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	TeamID   uint      `gorm:"not null;index" json:"team_id"`
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	Role     string    `gorm:"not null;default:'member'" json:"role"`
	Status   string    `gorm:"not null;default:'active'" json:"status"`
	JoinedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	gorm.Model

	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

func (TeamMember) TableName() string {
	return "user_team_members"
}

func (m *TeamMember) ValidateRole() error {
	if m.Role != "captain" && m.Role != "member" {
		return errors.New("invalid role: must be 'captain' or 'member'")
	}
	return nil
}

func (m *TeamMember) ValidateStatus() error {
	if m.Status != "active" && m.Status != "inactive" {
		return errors.New("invalid status: must be 'active' or 'inactive'")
	}
	return nil
}

func (m *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if m.Role == "" {
		m.Role = "member"
	}
	if m.Status == "" {
		m.Status = "active"
	}
	if m.JoinedAt.IsZero() {
		m.JoinedAt = time.Now()
	}
	if err := m.ValidateRole(); err != nil {
		return err
	}
	if err := m.ValidateStatus(); err != nil {
		return err
	}
	if m.TeamID == 0 {
		return errors.New("team ID cannot be empty")
	}
	if m.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	// Note: Duplicate membership check is handled at service layer for better control
	return nil
}

func (m *TeamMember) BeforeUpdate(tx *gorm.DB) error {
	if err := m.ValidateRole(); err != nil {
		return err
	}
	return m.ValidateStatus()
}

func (m *TeamMember) IsActive() bool {
	return m.Status == "active"
}

func (m *TeamMember) IsCaptain() bool {
	return m.Role == "captain"
}

package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// RelayEventAssignment represents assignment of an event to a team member in relay contest
// Captain assigns which team member predicts which match
type RelayEventAssignment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ContestID  uint      `gorm:"not null;index:idx_relay_contest;uniqueIndex:idx_relay_unique" json:"contest_id"`
	TeamID     uint      `gorm:"not null;index:idx_relay_team;uniqueIndex:idx_relay_unique" json:"team_id"`
	UserID     uint      `gorm:"not null;index:idx_relay_user" json:"user_id"`
	EventID    uint      `gorm:"not null;uniqueIndex:idx_relay_unique" json:"event_id"`
	AssignedBy uint      `gorm:"not null" json:"assigned_by"` // captain who made the assignment
	AssignedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"assigned_at"`
	gorm.Model

	// Relationships (optional, for eager loading via Preload)
	Event Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
}

func (RelayEventAssignment) TableName() string {
	return "relay_event_assignments"
}

// ValidateContestID checks if contest ID is valid
func (r *RelayEventAssignment) ValidateContestID() error {
	if r.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}
	return nil
}

// ValidateTeamID checks if team ID is valid
func (r *RelayEventAssignment) ValidateTeamID() error {
	if r.TeamID == 0 {
		return errors.New("team ID cannot be empty")
	}
	return nil
}

// ValidateUserID checks if user ID is valid
func (r *RelayEventAssignment) ValidateUserID() error {
	if r.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

// ValidateEventID checks if event ID is valid
func (r *RelayEventAssignment) ValidateEventID() error {
	if r.EventID == 0 {
		return errors.New("event ID cannot be empty")
	}
	return nil
}

// ValidateAssignedBy checks if assigner (captain) ID is valid
func (r *RelayEventAssignment) ValidateAssignedBy() error {
	if r.AssignedBy == 0 {
		return errors.New("assigned_by (captain) ID cannot be empty")
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating an assignment
func (r *RelayEventAssignment) BeforeCreate(tx *gorm.DB) error {
	if err := r.ValidateContestID(); err != nil {
		return err
	}
	if err := r.ValidateTeamID(); err != nil {
		return err
	}
	if err := r.ValidateUserID(); err != nil {
		return err
	}
	if err := r.ValidateEventID(); err != nil {
		return err
	}
	if err := r.ValidateAssignedBy(); err != nil {
		return err
	}

	if r.AssignedAt.IsZero() {
		r.AssignedAt = time.Now().UTC()
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating
func (r *RelayEventAssignment) BeforeUpdate(tx *gorm.DB) error {
	return r.BeforeCreate(tx)
}

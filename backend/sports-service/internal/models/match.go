package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LeagueID    uint      `gorm:"not null;index" json:"league_id"`
	HomeTeamID  uint      `gorm:"not null;index" json:"home_team_id"`
	AwayTeamID  uint      `gorm:"not null;index" json:"away_team_id"`
	ScheduledAt time.Time `gorm:"not null;index" json:"scheduled_at"`
	Status      string    `gorm:"not null;default:'scheduled'" json:"status"`
	HomeScore   int       `gorm:"default:0" json:"home_score"`
	AwayScore   int       `gorm:"default:0" json:"away_score"`
	ResultData  string    `gorm:"type:text" json:"result_data"`
	ExternalID  string    `gorm:"uniqueIndex;size:50" json:"external_id,omitempty"`
	League      League    `gorm:"foreignKey:LeagueID" json:"league,omitempty"`
	HomeTeam    Team      `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeam    Team      `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	gorm.Model
}

func (m *Match) ValidateLeagueID() error {
	if m.LeagueID == 0 {
		return errors.New("league_id is required")
	}
	return nil
}

func (m *Match) ValidateTeams() error {
	if m.HomeTeamID == 0 {
		return errors.New("home_team_id is required")
	}
	if m.AwayTeamID == 0 {
		return errors.New("away_team_id is required")
	}
	if m.HomeTeamID == m.AwayTeamID {
		return errors.New("home and away teams must be different")
	}
	return nil
}

func (m *Match) ValidateScheduledAt() error {
	if m.ScheduledAt.IsZero() {
		return errors.New("scheduled_at is required")
	}
	return nil
}

func (m *Match) ValidateStatus() error {
	validStatuses := []string{"scheduled", "live", "completed", "cancelled", "postponed"}
	for _, s := range validStatuses {
		if m.Status == s {
			return nil
		}
	}
	return errors.New("invalid status")
}

func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if err := m.ValidateLeagueID(); err != nil {
		return err
	}
	if err := m.ValidateTeams(); err != nil {
		return err
	}
	if err := m.ValidateScheduledAt(); err != nil {
		return err
	}
	if m.Status == "" {
		m.Status = "scheduled"
	}
	return m.ValidateStatus()
}

func (m *Match) BeforeUpdate(tx *gorm.DB) error {
	return m.BeforeCreate(tx)
}

func (m *Match) IsCompleted() bool {
	return m.Status == "completed"
}

func (m *Match) IsLive() bool {
	return m.Status == "live"
}

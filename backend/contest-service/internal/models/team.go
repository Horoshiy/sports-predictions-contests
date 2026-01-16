package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Team struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"not null" json:"name"`
	Description    string `json:"description"`
	InviteCode     string `gorm:"uniqueIndex;not null" json:"invite_code"`
	CaptainID      uint   `gorm:"not null" json:"captain_id"`
	MaxMembers     uint   `gorm:"default:10" json:"max_members"`
	CurrentMembers uint   `gorm:"default:0" json:"current_members"`
	IsActive       bool   `gorm:"default:true" json:"is_active"`
	gorm.Model
}

func (Team) TableName() string {
	return "user_teams"
}

func GenerateInviteCode() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func (t *Team) ValidateName() error {
	if len(strings.TrimSpace(t.Name)) == 0 {
		return errors.New("team name cannot be empty")
	}
	if len(t.Name) > 100 {
		return errors.New("team name cannot exceed 100 characters")
	}
	return nil
}

func (t *Team) ValidateDescription() error {
	if len(t.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}
	return nil
}

func (t *Team) ValidateMaxMembers() error {
	if t.MaxMembers < 2 {
		return errors.New("team must allow at least 2 members")
	}
	if t.MaxMembers > 50 {
		return errors.New("team cannot exceed 50 members")
	}
	return nil
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if err := t.ValidateName(); err != nil {
		return err
	}
	if err := t.ValidateDescription(); err != nil {
		return err
	}
	if t.MaxMembers == 0 {
		t.MaxMembers = 10
	}
	if err := t.ValidateMaxMembers(); err != nil {
		return err
	}
	if t.CaptainID == 0 {
		return errors.New("captain ID cannot be empty")
	}
	if t.InviteCode == "" {
		code, err := GenerateInviteCode()
		if err != nil {
			return err
		}
		t.InviteCode = code
	}
	return nil
}

func (t *Team) BeforeUpdate(tx *gorm.DB) error {
	if err := t.ValidateName(); err != nil {
		return err
	}
	if err := t.ValidateDescription(); err != nil {
		return err
	}
	if err := t.ValidateMaxMembers(); err != nil {
		return err
	}
	return nil
}

func (t *Team) CanJoin() bool {
	return t.IsActive && t.CurrentMembers < t.MaxMembers
}

func (t *Team) IsCaptain(userID uint) bool {
	return t.CaptainID == userID
}

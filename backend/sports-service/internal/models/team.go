package models

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Team struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	SportID    uint   `gorm:"not null;index" json:"sport_id"`
	Name       string `gorm:"not null" json:"name"`
	Slug       string `gorm:"not null;uniqueIndex" json:"slug"`
	ShortName  string `json:"short_name"`
	LogoURL    string `json:"logo_url"`
	Country    string `json:"country"`
	ExternalID string `gorm:"uniqueIndex;size:50" json:"external_id,omitempty"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
	Sport      Sport  `gorm:"foreignKey:SportID" json:"sport,omitempty"`
	gorm.Model
}

func (t *Team) ValidateName() error {
	if len(strings.TrimSpace(t.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(t.Name) > 200 {
		return errors.New("name cannot exceed 200 characters")
	}
	return nil
}

func (t *Team) ValidateSlug() error {
	if len(strings.TrimSpace(t.Slug)) == 0 {
		return errors.New("slug cannot be empty")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, t.Slug)
	if !matched {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}

func (t *Team) ValidateSportID() error {
	if t.SportID == 0 {
		return errors.New("sport_id is required")
	}
	return nil
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if err := t.ValidateName(); err != nil {
		return err
	}
	if err := t.ValidateSportID(); err != nil {
		return err
	}
	if t.Slug == "" {
		slug := strings.ToLower(t.Name)
		slug = strings.ReplaceAll(slug, " ", "-")
		slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "")
		slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
		t.Slug = strings.Trim(slug, "-")
	}
	return t.ValidateSlug()
}

func (t *Team) BeforeUpdate(tx *gorm.DB) error {
	return t.BeforeCreate(tx)
}

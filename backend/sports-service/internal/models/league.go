package models

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type League struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	SportID  uint   `gorm:"not null;index" json:"sport_id"`
	Name     string `gorm:"not null" json:"name"`
	Slug     string `gorm:"not null;uniqueIndex" json:"slug"`
	Country  string `json:"country"`
	Season   string `json:"season"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	Sport    Sport  `gorm:"foreignKey:SportID" json:"sport,omitempty"`
	gorm.Model
}

func (l *League) ValidateName() error {
	if len(strings.TrimSpace(l.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(l.Name) > 200 {
		return errors.New("name cannot exceed 200 characters")
	}
	return nil
}

func (l *League) ValidateSlug() error {
	if len(strings.TrimSpace(l.Slug)) == 0 {
		return errors.New("slug cannot be empty")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, l.Slug)
	if !matched {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}

func (l *League) ValidateSportID() error {
	if l.SportID == 0 {
		return errors.New("sport_id is required")
	}
	return nil
}

func (l *League) BeforeCreate(tx *gorm.DB) error {
	if err := l.ValidateName(); err != nil {
		return err
	}
	if err := l.ValidateSportID(); err != nil {
		return err
	}
	if l.Slug == "" {
		slug := strings.ToLower(l.Name)
		slug = strings.ReplaceAll(slug, " ", "-")
		slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "")
		slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
		l.Slug = strings.Trim(slug, "-")
	}
	return l.ValidateSlug()
}

func (l *League) BeforeUpdate(tx *gorm.DB) error {
	return l.BeforeCreate(tx)
}

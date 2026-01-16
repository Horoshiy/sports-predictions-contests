package models

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Sport struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Slug        string `gorm:"not null;uniqueIndex" json:"slug"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	gorm.Model
}

func (s *Sport) ValidateName() error {
	if len(strings.TrimSpace(s.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(s.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	return nil
}

func (s *Sport) ValidateSlug() error {
	if len(strings.TrimSpace(s.Slug)) == 0 {
		return errors.New("slug cannot be empty")
	}
	if len(s.Slug) > 100 {
		return errors.New("slug cannot exceed 100 characters")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, s.Slug)
	if !matched {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}

func (s *Sport) BeforeCreate(tx *gorm.DB) error {
	if err := s.ValidateName(); err != nil {
		return err
	}
	if s.Slug == "" {
		// Sanitize: lowercase, replace spaces with hyphens, remove invalid chars
		slug := strings.ToLower(s.Name)
		slug = strings.ReplaceAll(slug, " ", "-")
		slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "")
		slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
		s.Slug = strings.Trim(slug, "-")
	}
	return s.ValidateSlug()
}

func (s *Sport) BeforeUpdate(tx *gorm.DB) error {
	return s.BeforeCreate(tx)
}

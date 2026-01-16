package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// PropType represents a type of prop prediction available for a sport
type PropType struct {
	gorm.Model
	SportType     string   `gorm:"not null;uniqueIndex:idx_prop_sport_slug" json:"sport_type"`
	Name          string   `gorm:"not null" json:"name"`
	Slug          string   `gorm:"not null;uniqueIndex:idx_prop_sport_slug" json:"slug"`
	Description   string   `json:"description"`
	Category      string   `gorm:"not null" json:"category"`
	ValueType     string   `gorm:"not null" json:"value_type"`
	DefaultLine   *float64 `json:"default_line"`
	MinValue      *float64 `json:"min_value"`
	MaxValue      *float64 `json:"max_value"`
	PointsCorrect float64  `gorm:"not null;default:2" json:"points_correct"`
	IsActive      bool     `gorm:"default:true" json:"is_active"`
}

func (PropType) TableName() string {
	return "prop_types"
}

func (p *PropType) ValidateSportType() error {
	if strings.TrimSpace(p.SportType) == "" {
		return errors.New("sport type cannot be empty")
	}
	return nil
}

func (p *PropType) ValidateName() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if len(p.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	return nil
}

func (p *PropType) ValidateCategory() error {
	validCategories := []string{"match", "player", "team"}
	for _, c := range validCategories {
		if p.Category == c {
			return nil
		}
	}
	return errors.New("invalid category: must be match, player, or team")
}

func (p *PropType) ValidateValueType() error {
	validTypes := []string{"over_under", "yes_no", "team_select", "player_select", "exact_value"}
	for _, t := range validTypes {
		if p.ValueType == t {
			return nil
		}
	}
	return errors.New("invalid value type")
}

func (p *PropType) ValidateSlug() error {
	if strings.TrimSpace(p.Slug) == "" {
		return errors.New("slug cannot be empty")
	}
	return nil
}

func (p *PropType) BeforeCreate(tx *gorm.DB) error {
	if err := p.ValidateSportType(); err != nil {
		return err
	}
	if err := p.ValidateName(); err != nil {
		return err
	}
	if err := p.ValidateSlug(); err != nil {
		return err
	}
	if err := p.ValidateCategory(); err != nil {
		return err
	}
	return p.ValidateValueType()
}

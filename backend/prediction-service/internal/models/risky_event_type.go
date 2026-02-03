package models

import "time"

// RiskyEventType represents a global risky event type that can be used in contests
type RiskyEventType struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Slug          string    `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	NameEn        string    `gorm:"size:100" json:"name_en"`
	Description   string    `json:"description"`
	DefaultPoints float64   `gorm:"type:decimal(5,2);not null;default:3" json:"default_points"`
	SportType     string    `gorm:"size:50;default:'football'" json:"sport_type"`
	Category      string    `gorm:"size:50;default:'general'" json:"category"`
	Icon          string    `gorm:"size:10" json:"icon"`
	SortOrder     int       `gorm:"default:0" json:"sort_order"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (RiskyEventType) TableName() string {
	return "risky_event_types"
}

// Categories for risky events
const (
	CategoryGoals   = "goals"
	CategoryCards   = "cards"
	CategoryDefense = "defense"
	CategoryTotals  = "totals"
	CategoryHalves  = "halves"
	CategoryTiming  = "timing"
	CategorySpecial = "special"
	CategoryGeneral = "general"
)

// AllCategories returns all available categories
func AllCategories() []string {
	return []string{
		CategoryGoals,
		CategoryCards,
		CategoryDefense,
		CategoryTotals,
		CategoryHalves,
		CategoryTiming,
		CategorySpecial,
		CategoryGeneral,
	}
}

package models

import "time"

// MatchRiskyEvent represents a risky event override for a specific match/event
// It allows customizing points and tracking outcomes for individual matches
type MatchRiskyEvent struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	EventID          uint    `gorm:"not null" json:"event_id"`
	RiskyEventTypeID uint    `gorm:"not null" json:"risky_event_type_id"`
	Points           float64 `gorm:"type:decimal(5,2);not null" json:"points"`
	IsEnabled        bool    `gorm:"default:true" json:"is_enabled"`
	Outcome          *bool   `json:"outcome"` // nil=pending, true=happened, false=didn't happen

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	RiskyEventType RiskyEventType `gorm:"foreignKey:RiskyEventTypeID" json:"risky_event_type,omitempty"`
	Event          Event          `gorm:"foreignKey:EventID" json:"event,omitempty"`
}

// TableName specifies the table name for GORM
func (MatchRiskyEvent) TableName() string {
	return "match_risky_events"
}

// MatchRiskyEventView is a view combining event type info with match-specific overrides
type MatchRiskyEventView struct {
	RiskyEventTypeID uint    `json:"risky_event_type_id"`
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	NameEn           string  `json:"name_en"`
	Icon             string  `json:"icon"`
	Category         string  `json:"category"`
	Points           float64 `json:"points"`    // Final points (override or default)
	IsEnabled        bool    `json:"is_enabled"`
	Outcome          *bool   `json:"outcome"`
	IsOverridden     bool    `json:"is_overridden"` // True if points differ from default
}

// RiskyPrediction represents a user's risky prediction for a match
type RiskyPrediction struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	PredictionID     uint    `gorm:"not null" json:"prediction_id"`
	RiskyEventTypeID uint    `gorm:"not null" json:"risky_event_type_id"`
	PointsIfCorrect  float64 `gorm:"type:decimal(5,2);not null" json:"points_if_correct"`

	CreatedAt time.Time `json:"created_at"`

	// Relations
	RiskyEventType RiskyEventType `gorm:"foreignKey:RiskyEventTypeID" json:"risky_event_type,omitempty"`
	Prediction     Prediction     `gorm:"foreignKey:PredictionID" json:"prediction,omitempty"`
}

// TableName specifies the table name for GORM
func (RiskyPrediction) TableName() string {
	return "risky_predictions"
}

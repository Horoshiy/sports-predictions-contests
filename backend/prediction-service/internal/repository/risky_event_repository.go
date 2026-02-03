package repository

import (
	"encoding/json"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"gorm.io/gorm"
)

type RiskyEventRepository struct {
	db *gorm.DB
}

func NewRiskyEventRepository(db *gorm.DB) *RiskyEventRepository {
	return &RiskyEventRepository{db: db}
}

// ListActiveEventTypes returns all active risky event types for a sport
func (r *RiskyEventRepository) ListActiveEventTypes(sportType string) ([]models.RiskyEventType, error) {
	var eventTypes []models.RiskyEventType
	query := r.db.Where("is_active = ?", true)
	if sportType != "" {
		query = query.Where("sport_type = ?", sportType)
	}
	err := query.Order("sort_order ASC, id ASC").Find(&eventTypes).Error
	return eventTypes, err
}

// ListAllEventTypes returns all risky event types (for admin)
func (r *RiskyEventRepository) ListAllEventTypes(sportType string) ([]models.RiskyEventType, error) {
	var eventTypes []models.RiskyEventType
	query := r.db
	if sportType != "" {
		query = query.Where("sport_type = ?", sportType)
	}
	err := query.Order("sort_order ASC, id ASC").Find(&eventTypes).Error
	return eventTypes, err
}

// GetEventType returns event type by ID
func (r *RiskyEventRepository) GetEventType(id uint) (*models.RiskyEventType, error) {
	var et models.RiskyEventType
	err := r.db.First(&et, id).Error
	if err != nil {
		return nil, err
	}
	return &et, nil
}

// GetEventTypeBySlug returns event type by slug
func (r *RiskyEventRepository) GetEventTypeBySlug(slug string) (*models.RiskyEventType, error) {
	var et models.RiskyEventType
	err := r.db.Where("slug = ?", slug).First(&et).Error
	if err != nil {
		return nil, err
	}
	return &et, nil
}

// CreateEventType creates a new risky event type
func (r *RiskyEventRepository) CreateEventType(et *models.RiskyEventType) error {
	return r.db.Create(et).Error
}

// UpdateEventType updates an existing risky event type
func (r *RiskyEventRepository) UpdateEventType(et *models.RiskyEventType) error {
	return r.db.Save(et).Error
}

// DeleteEventType deletes a risky event type (soft delete via is_active)
func (r *RiskyEventRepository) DeleteEventType(id uint) error {
	return r.db.Model(&models.RiskyEventType{}).Where("id = ?", id).Update("is_active", false).Error
}

// GetMatchRiskyEvents returns risky events for a match with overrides applied
// Priority: match override > contest rules > global default
func (r *RiskyEventRepository) GetMatchRiskyEvents(eventID uint, contestRulesJSON string) ([]models.MatchRiskyEventView, error) {
	// Parse contest rules to get event IDs and point overrides
	contestEvents := parseContestRiskyEvents(contestRulesJSON)

	// Get all active event types
	var eventTypes []models.RiskyEventType
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&eventTypes).Error
	if err != nil {
		return nil, err
	}

	// Get match-level overrides
	var matchOverrides []models.MatchRiskyEvent
	r.db.Where("event_id = ?", eventID).Find(&matchOverrides)

	// Build override map
	overrideMap := make(map[uint]models.MatchRiskyEvent)
	for _, mo := range matchOverrides {
		overrideMap[mo.RiskyEventTypeID] = mo
	}

	// Build result
	var result []models.MatchRiskyEventView

	// If contest has specific events, use only those
	if len(contestEvents) > 0 {
		for _, ce := range contestEvents {
			// Find the event type
			var et *models.RiskyEventType
			for i := range eventTypes {
				if eventTypes[i].ID == ce.RiskyEventTypeID || eventTypes[i].Slug == ce.Slug {
					et = &eventTypes[i]
					break
				}
			}
			if et == nil {
				continue
			}

			view := models.MatchRiskyEventView{
				RiskyEventTypeID: et.ID,
				Slug:             et.Slug,
				Name:             et.Name,
				NameEn:           et.NameEn,
				Icon:             et.Icon,
				Category:         et.Category,
				Points:           ce.Points, // Contest override
				IsEnabled:        true,
				IsOverridden:     ce.Points != et.DefaultPoints,
			}

			// Apply match-level override if exists
			if mo, ok := overrideMap[et.ID]; ok {
				view.Points = mo.Points
				view.IsEnabled = mo.IsEnabled
				view.Outcome = mo.Outcome
				view.IsOverridden = true
			}

			result = append(result, view)
		}
	} else {
		// No contest-specific events, use all active with defaults
		for _, et := range eventTypes {
			view := models.MatchRiskyEventView{
				RiskyEventTypeID: et.ID,
				Slug:             et.Slug,
				Name:             et.Name,
				NameEn:           et.NameEn,
				Icon:             et.Icon,
				Category:         et.Category,
				Points:           et.DefaultPoints,
				IsEnabled:        true,
				IsOverridden:     false,
			}

			// Apply match-level override if exists
			if mo, ok := overrideMap[et.ID]; ok {
				view.Points = mo.Points
				view.IsEnabled = mo.IsEnabled
				view.Outcome = mo.Outcome
				view.IsOverridden = true
			}

			result = append(result, view)
		}
	}

	return result, nil
}

// SetMatchEventOverride creates or updates a match-level point override
func (r *RiskyEventRepository) SetMatchEventOverride(eventID uint, riskyEventTypeID uint, points float64, isEnabled bool) error {
	var existing models.MatchRiskyEvent
	err := r.db.Where("event_id = ? AND risky_event_type_id = ?", eventID, riskyEventTypeID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new
		return r.db.Create(&models.MatchRiskyEvent{
			EventID:          eventID,
			RiskyEventTypeID: riskyEventTypeID,
			Points:           points,
			IsEnabled:        isEnabled,
		}).Error
	} else if err != nil {
		return err
	}

	// Update existing
	existing.Points = points
	existing.IsEnabled = isEnabled
	return r.db.Save(&existing).Error
}

// SetMatchEventOutcome records the outcome of a risky event after match completion
func (r *RiskyEventRepository) SetMatchEventOutcome(eventID uint, riskyEventTypeID uint, happened bool) error {
	return r.db.Model(&models.MatchRiskyEvent{}).
		Where("event_id = ? AND risky_event_type_id = ?", eventID, riskyEventTypeID).
		Update("outcome", happened).Error
}

// DeleteMatchEventOverride removes a match-level override
func (r *RiskyEventRepository) DeleteMatchEventOverride(eventID uint, riskyEventTypeID uint) error {
	return r.db.Where("event_id = ? AND risky_event_type_id = ?", eventID, riskyEventTypeID).
		Delete(&models.MatchRiskyEvent{}).Error
}

// contestRiskyEvent represents a risky event in contest rules JSON
type contestRiskyEvent struct {
	RiskyEventTypeID uint    `json:"risky_event_type_id"`
	Slug             string  `json:"slug"`
	Points           float64 `json:"points"`
}

// parseContestRiskyEvents extracts risky events from contest rules JSON
func parseContestRiskyEvents(rulesJSON string) []contestRiskyEvent {
	if rulesJSON == "" {
		return nil
	}

	var rules struct {
		Type  string `json:"type"`
		Risky *struct {
			Events []contestRiskyEvent `json:"events"`
		} `json:"risky"`
	}

	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return nil
	}

	if rules.Risky == nil || len(rules.Risky.Events) == 0 {
		return nil
	}

	return rules.Risky.Events
}

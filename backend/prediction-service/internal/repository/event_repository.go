package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"gorm.io/gorm"
)

// EventRepositoryInterface defines the contract for event repository
type EventRepositoryInterface interface {
	Create(event *models.Event) error
	GetByID(id uint) (*models.Event, error)
	Update(event *models.Event) error
	Delete(id uint) error
	List(limit, offset int, sportType, status string) ([]*models.Event, int64, error)
	ListByContest(contestID uint, sportType, status string) ([]*models.Event, int64, error)
	GetBySportType(sportType string) ([]*models.Event, error)
	GetByDateRange(startDate, endDate time.Time) ([]*models.Event, error)
	GetUpcoming(limit int) ([]*models.Event, error)
	// Contest-Event management
	AddEventsToContest(contestID uint, eventIDs []uint) error
	RemoveEventsFromContest(contestID uint, eventIDs []uint) error
	SetContestEvents(contestID uint, eventIDs []uint) error
	GetContestEventCount(contestID uint) (int64, error)
}

// EventRepository implements EventRepositoryInterface
type EventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new event repository instance
func NewEventRepository(db *gorm.DB) EventRepositoryInterface {
	return &EventRepository{db: db}
}

// Create creates a new event
func (r *EventRepository) Create(event *models.Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	return r.db.Create(event).Error
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.First(&event, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *models.Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	return r.db.Save(event).Error
}

// Delete deletes an event by ID
func (r *EventRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Event{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

// List retrieves events with pagination and optional filters
func (r *EventRepository) List(limit, offset int, sportType, status string) ([]*models.Event, int64, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Cap at 100 to prevent performance issues
	}
	if offset < 0 {
		offset = 0
	}

	var events []*models.Event
	var total int64

	query := r.db.Model(&models.Event{})

	// Apply filters
	if sportType != "" {
		query = query.Where("sport_type = ?", sportType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch records
	err := query.Limit(limit).Offset(offset).Order("event_date ASC").Find(&events).Error
	return events, total, err
}

// ListByContest retrieves events for a specific contest
func (r *EventRepository) ListByContest(contestID uint, sportType, status string) ([]*models.Event, int64, error) {
	var events []*models.Event
	var total int64

	// Join with contest_events table to filter by contest
	query := r.db.Model(&models.Event{}).
		Joins("INNER JOIN contest_events ce ON ce.event_id = events.id").
		Where("ce.contest_id = ?", contestID)

	// Apply additional filters
	if sportType != "" {
		query = query.Where("events.sport_type = ?", sportType)
	}
	if status != "" {
		query = query.Where("events.status = ?", status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch records
	err := query.Order("events.event_date ASC").Find(&events).Error
	return events, total, err
}

// GetBySportType retrieves events by sport type
func (r *EventRepository) GetBySportType(sportType string) ([]*models.Event, error) {
	var events []*models.Event
	err := r.db.Where("sport_type = ?", sportType).Order("event_date ASC").Find(&events).Error
	return events, err
}

// GetByDateRange retrieves events within a date range
func (r *EventRepository) GetByDateRange(startDate, endDate time.Time) ([]*models.Event, error) {
	var events []*models.Event
	err := r.db.Where("event_date BETWEEN ? AND ?", startDate, endDate).Order("event_date ASC").Find(&events).Error
	return events, err
}

// GetUpcoming retrieves upcoming events
func (r *EventRepository) GetUpcoming(limit int) ([]*models.Event, error) {
	var events []*models.Event
	now := time.Now().UTC()
	
	query := r.db.Where("event_date > ? AND status = ?", now, "scheduled").
		Order("event_date ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&events).Error
	return events, err
}

// AddEventsToContest adds events to a contest (many-to-many relationship)
func (r *EventRepository) AddEventsToContest(contestID uint, eventIDs []uint) error {
	if len(eventIDs) == 0 {
		return nil
	}

	// Build bulk insert query for better performance
	valueStrings := make([]string, len(eventIDs))
	valueArgs := make([]interface{}, 0, len(eventIDs)*2)
	
	for i, eventID := range eventIDs {
		valueStrings[i] = "(?, ?)"
		valueArgs = append(valueArgs, contestID, eventID)
	}

	query := "INSERT INTO contest_events (contest_id, event_id) VALUES " +
		strings.Join(valueStrings, ", ") +
		" ON CONFLICT DO NOTHING"

	return r.db.Exec(query, valueArgs...).Error
}

// RemoveEventsFromContest removes events from a contest
func (r *EventRepository) RemoveEventsFromContest(contestID uint, eventIDs []uint) error {
	if len(eventIDs) == 0 {
		return nil
	}

	return r.db.Exec(
		"DELETE FROM contest_events WHERE contest_id = ? AND event_id IN ?",
		contestID, eventIDs,
	).Error
}

// SetContestEvents replaces all events for a contest (removes old, adds new)
func (r *EventRepository) SetContestEvents(contestID uint, eventIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Remove all existing events for this contest
		if err := tx.Exec("DELETE FROM contest_events WHERE contest_id = ?", contestID).Error; err != nil {
			return err
		}

		// Add new events
		for _, eventID := range eventIDs {
			if err := tx.Exec(
				"INSERT INTO contest_events (contest_id, event_id) VALUES (?, ?)",
				contestID, eventID,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetContestEventCount returns the number of events in a contest
func (r *EventRepository) GetContestEventCount(contestID uint) (int64, error) {
	var count int64
	err := r.db.Raw("SELECT COUNT(*) FROM contest_events WHERE contest_id = ?", contestID).Scan(&count).Error
	return count, err
}

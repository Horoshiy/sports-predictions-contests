package repository

import (
	"errors"
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
	GetBySportType(sportType string) ([]*models.Event, error)
	GetByDateRange(startDate, endDate time.Time) ([]*models.Event, error)
	GetUpcoming(limit int) ([]*models.Event, error)
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

package repository

import (
	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

// EventRepositoryInterface defines the interface for event operations
type EventRepositoryInterface interface {
	Create(event *models.Event) error
	GetByID(id uint) (*models.Event, error)
	Update(event *models.Event) error
	Delete(id uint) error
	List(limit, offset int, sportType, status string) ([]*models.Event, int64, error)
}

// EventRepository implements EventRepositoryInterface
type EventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create creates a new event
func (r *EventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

// Delete soft deletes an event
func (r *EventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}

// List retrieves events with optional filtering
func (r *EventRepository) List(limit, offset int, sportType, status string) ([]*models.Event, int64, error) {
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

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results ordered by event date
	if err := query.Order("event_date ASC").Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

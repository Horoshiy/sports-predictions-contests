package repository

import (
	"context"
	"errors"
	"time"

	"github.com/sports-prediction-contests/notification-service/internal/models"
	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByID(ctx context.Context, id uint) (*models.Notification, error)
	Update(ctx context.Context, notification *models.Notification) error
	Delete(ctx context.Context, id uint) error
	GetByUser(ctx context.Context, userID uint, unreadOnly bool, limit, offset int) ([]*models.Notification, int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context, userID uint) (int64, error)
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	GetPreferences(ctx context.Context, userID uint) ([]*models.NotificationPreference, error)
	GetPreference(ctx context.Context, userID uint, channel string) (*models.NotificationPreference, error)
	UpdatePreference(ctx context.Context, pref *models.NotificationPreference) error
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepositoryInterface {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.WithContext(ctx).First(&notification, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

func (r *NotificationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Notification{}, id).Error
}

func (r *NotificationRepository) GetByUser(ctx context.Context, userID uint, unreadOnly bool, limit, offset int) ([]*models.Notification, int64, error) {
	var notifications []*models.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	}).Error
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uint) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{"is_read": true, "read_at": now})
	return result.RowsAffected, result.Error
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *NotificationRepository) GetPreferences(ctx context.Context, userID uint) ([]*models.NotificationPreference, error) {
	var prefs []*models.NotificationPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&prefs).Error
	return prefs, err
}

func (r *NotificationRepository) GetPreference(ctx context.Context, userID uint, channel string) (*models.NotificationPreference, error) {
	var pref models.NotificationPreference
	err := r.db.WithContext(ctx).Where("user_id = ? AND channel = ?", userID, channel).First(&pref).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pref, err
}

func (r *NotificationRepository) UpdatePreference(ctx context.Context, pref *models.NotificationPreference) error {
	return r.db.WithContext(ctx).Save(pref).Error
}

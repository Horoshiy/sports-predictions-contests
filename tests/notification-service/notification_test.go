package notification_test

import (
	"testing"

	"github.com/sports-prediction-contests/notification-service/internal/models"
)

func TestNotificationValidation(t *testing.T) {
	notification := &models.Notification{
		UserID:  1,
		Type:    "PREDICTION_RESULT",
		Title:   "Test Notification",
		Message: "This is a test message",
		Channel: "in_app",
	}

	if err := notification.ValidateUserID(); err != nil {
		t.Errorf("Expected valid user ID, got error: %v", err)
	}

	if err := notification.ValidateTitle(); err != nil {
		t.Errorf("Expected valid title, got error: %v", err)
	}

	if err := notification.ValidateMessage(); err != nil {
		t.Errorf("Expected valid message, got error: %v", err)
	}

	// Test invalid user ID
	notification.UserID = 0
	if err := notification.ValidateUserID(); err == nil {
		t.Error("Expected error for invalid user ID")
	}

	// Test empty title
	notification.UserID = 1
	notification.Title = ""
	if err := notification.ValidateTitle(); err == nil {
		t.Error("Expected error for empty title")
	}

	// Test empty message
	notification.Title = "Test"
	notification.Message = ""
	if err := notification.ValidateMessage(); err == nil {
		t.Error("Expected error for empty message")
	}
}

func TestNotificationPreferenceValidation(t *testing.T) {
	pref := &models.NotificationPreference{
		UserID:  1,
		Channel: "telegram",
		Enabled: true,
	}

	if err := pref.ValidateUserID(); err != nil {
		t.Errorf("Expected valid user ID, got error: %v", err)
	}

	if err := pref.ValidateChannel(); err != nil {
		t.Errorf("Expected valid channel, got error: %v", err)
	}

	// Test invalid user ID
	pref.UserID = 0
	if err := pref.ValidateUserID(); err == nil {
		t.Error("Expected error for invalid user ID")
	}

	// Test empty channel
	pref.UserID = 1
	pref.Channel = ""
	if err := pref.ValidateChannel(); err == nil {
		t.Error("Expected error for empty channel")
	}
}

func TestNotificationTitleLength(t *testing.T) {
	notification := &models.Notification{
		UserID:  1,
		Message: "Test",
		Channel: "in_app",
	}

	// Test title exceeding 255 characters
	longTitle := make([]byte, 256)
	for i := range longTitle {
		longTitle[i] = 'a'
	}
	notification.Title = string(longTitle)

	if err := notification.ValidateTitle(); err == nil {
		t.Error("Expected error for title exceeding 255 characters")
	}

	// Test valid title at boundary
	notification.Title = string(longTitle[:255])
	if err := notification.ValidateTitle(); err != nil {
		t.Errorf("Expected valid title at 255 chars, got error: %v", err)
	}
}

package notification_test

import (
	"testing"

	"github.com/sports-prediction-contests/notification-service/internal/channels"
)

func TestTelegramChannelDisabled(t *testing.T) {
	telegram, err := channels.NewTelegramChannel("", false)
	if err != nil {
		t.Errorf("Expected no error for disabled channel, got: %v", err)
	}

	if telegram.IsEnabled() {
		t.Error("Expected telegram channel to be disabled")
	}

	// Send should return nil when disabled
	err = telegram.Send(12345, "Test", "Message")
	if err != nil {
		t.Errorf("Expected nil error when sending to disabled channel, got: %v", err)
	}
}

func TestTelegramChannelNoToken(t *testing.T) {
	telegram, err := channels.NewTelegramChannel("", true)
	if err != nil {
		t.Errorf("Expected no error for empty token, got: %v", err)
	}

	if telegram.IsEnabled() {
		t.Error("Expected telegram channel to be disabled with empty token")
	}
}

func TestEmailChannelDisabled(t *testing.T) {
	email := channels.NewEmailChannel("", "", "", "", "", false)

	if email.IsEnabled() {
		t.Error("Expected email channel to be disabled")
	}

	// Send should return nil when disabled
	err := email.Send("test@example.com", "Test", "Message")
	if err != nil {
		t.Errorf("Expected nil error when sending to disabled channel, got: %v", err)
	}
}

func TestEmailChannelNoHost(t *testing.T) {
	email := channels.NewEmailChannel("", "587", "user", "pass", "from@test.com", true)

	if email.IsEnabled() {
		t.Error("Expected email channel to be disabled with empty host")
	}
}

func TestEmailChannelEnabled(t *testing.T) {
	email := channels.NewEmailChannel("smtp.example.com", "587", "user", "pass", "from@test.com", true)

	if !email.IsEnabled() {
		t.Error("Expected email channel to be enabled")
	}
}

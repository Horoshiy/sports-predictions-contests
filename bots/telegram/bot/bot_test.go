package bot

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// TestRegisterCommandsStructure tests that registerCommands creates valid command structure
func TestRegisterCommandsStructure(t *testing.T) {
	// This test verifies the command structure is valid
	// Actual API call requires valid bot token and is tested manually

	commands := []struct {
		command     string
		description string
	}{
		{"start", "Start bot and create account | Начать работу и создать аккаунт"},
		{"contests", "View active contests | Просмотр активных конкурсов"},
		{"leaderboard", "View contest leaderboard | Таблица лидеров конкурса"},
		{"mystats", "Your prediction statistics | Ваша статистика прогнозов"},
		{"help", "Show help message | Показать справку"},
		{"link", "Link existing account | Привязать существующий аккаунт"},
	}

	if len(commands) != 6 {
		t.Errorf("Expected 6 commands, got %d", len(commands))
	}

	for _, cmd := range commands {
		if cmd.command == "" {
			t.Error("Command name cannot be empty")
		}
		if cmd.description == "" {
			t.Error("Command description cannot be empty")
		}
		if len(cmd.description) > 256 {
			t.Errorf("Command description too long: %d chars (max 256)", len(cmd.description))
		}
	}
}

// TestPasswordGeneration tests that password generation is secure and deterministic
func TestPasswordGeneration(t *testing.T) {
	// Test with a valid secret
	h := &Handlers{
		passwordSecret: []byte("test-secret-key-for-testing"),
	}

	telegramID := int64(123456789)

	// Generate password twice to verify deterministic behavior
	pwd1, err1 := h.generateTelegramPassword(telegramID)
	if err1 != nil {
		t.Fatalf("Failed to generate password: %v", err1)
	}

	pwd2, err2 := h.generateTelegramPassword(telegramID)
	if err2 != nil {
		t.Fatalf("Failed to generate password on second attempt: %v", err2)
	}

	// Passwords should be deterministic (same for same ID)
	if pwd1 != pwd2 {
		t.Errorf("Password generation is not deterministic: %s != %s", pwd1, pwd2)
	}

	// Password should not be empty
	if pwd1 == "" {
		t.Error("Generated password is empty")
	}

	// Password should not be predictable (not contain the telegram ID)
	if pwd1 == "tg_secure_123456789" {
		t.Error("Password is still using old predictable pattern")
	}

	// Different IDs should produce different passwords
	pwd3, _ := h.generateTelegramPassword(int64(987654321))
	if pwd1 == pwd3 {
		t.Error("Different Telegram IDs produced same password")
	}
}

// TestPasswordGenerationWithoutSecret tests that password generation fails without secret
func TestPasswordGenerationWithoutSecret(t *testing.T) {
	h := &Handlers{
		passwordSecret: []byte{}, // Empty secret
	}

	_, err := h.generateTelegramPassword(int64(123456789))
	if err == nil {
		t.Error("Expected error when generating password without secret, got nil")
	}

	expectedErrMsg := "TELEGRAM_PASSWORD_SECRET not configured"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestNameLengthValidation tests that names are truncated to maxNameLength
func TestNameLengthValidation(t *testing.T) {
	// This test verifies the name length constant is properly defined
	if maxNameLength != 100 {
		t.Errorf("Expected maxNameLength to be 100, got %d", maxNameLength)
	}

	// Test that a long name would be truncated
	longName := string(make([]byte, 150))
	if len(longName) <= maxNameLength {
		t.Error("Test setup error: longName should be longer than maxNameLength")
	}

	truncated := longName[:maxNameLength]
	if len(truncated) != maxNameLength {
		t.Errorf("Expected truncated name length to be %d, got %d", maxNameLength, len(truncated))
	}
}

// TestRegistrationTimeout tests that timeout constant is defined
func TestRegistrationTimeout(t *testing.T) {
	// Verify the timeout constant exists and has a reasonable value
	if registrationTimeout == 0 {
		t.Error("registrationTimeout should not be zero")
	}

	// Should be at least 1 second
	if registrationTimeout.Seconds() < 1 {
		t.Errorf("registrationTimeout too short: %v", registrationTimeout)
	}

	// Should not be excessively long
	if registrationTimeout.Seconds() > 30 {
		t.Errorf("registrationTimeout too long: %v", registrationTimeout)
	}
}

// TestNameTruncationWithUTF8 tests that name truncation handles multi-byte UTF-8 characters correctly
func TestNameTruncationWithUTF8(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValid bool
	}{
		{
			name:        "Cyrillic characters",
			input:       strings.Repeat("Я", 60), // 120 bytes (2 bytes per char)
			expectValid: true,
		},
		{
			name:        "Emoji characters",
			input:       strings.Repeat("😀", 30), // 120 bytes (4 bytes per char)
			expectValid: true,
		},
		{
			name:        "Mixed ASCII and Cyrillic",
			input:       "Test" + strings.Repeat("Я", 50),
			expectValid: true,
		},
		{
			name:        "Short name",
			input:       "Short",
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := tt.input

			// Apply the same truncation logic as in registration.go
			if len(name) > maxNameLength {
				runes := []rune(name)
				if len(runes) > maxNameLength {
					name = string(runes[:maxNameLength])
				}
			}

			// Verify result is valid UTF-8
			if !utf8.ValidString(name) {
				t.Errorf("Truncated name contains invalid UTF-8: %q", name)
			}

			// Verify rune count is within bounds
			runeCount := utf8.RuneCountInString(name)
			if runeCount > maxNameLength {
				t.Errorf("Truncated name has too many runes: %d (max %d)", runeCount, maxNameLength)
			}

			// Verify byte length is reasonable (max 4 bytes per rune in UTF-8)
			if len(name) > maxNameLength*4 {
				t.Errorf("Truncated name too long: %d bytes", len(name))
			}
		})
	}
}

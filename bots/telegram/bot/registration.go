package bot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userpb "github.com/sports-prediction-contests/shared/proto/user"
)

const (
	registrationTimeout = 5 * time.Second
	maxNameLength       = 100
)

// generateTelegramPassword creates a secure deterministic password using HMAC
func (h *Handlers) generateTelegramPassword(telegramID int64) (string, error) {
	if len(h.passwordSecret) == 0 {
		return "", fmt.Errorf("TELEGRAM_PASSWORD_SECRET not configured")
	}

	mac := hmac.New(sha256.New, h.passwordSecret)
	mac.Write([]byte(fmt.Sprintf("%d", telegramID)))
	hash := mac.Sum(nil)

	return base64.URLEncoding.EncodeToString(hash), nil
}

// registerViaTelegram registers a new user using Telegram credentials
func (h *Handlers) registerViaTelegram(msg *tgbotapi.Message) {
	// Check if message has a sender
	if msg.From == nil {
		log.Printf("[ERROR] Cannot register: message has no sender")
		h.sendMessage(msg.Chat.ID, "❌ Registration failed: invalid message", nil)
		return
	}

	// Generate email from Telegram ID
	email := fmt.Sprintf("tg_%d@telegram.bot", msg.From.ID)

	// Use firstName + lastName or username as name
	name := strings.TrimSpace(msg.From.FirstName + " " + msg.From.LastName)
	if name == "" {
		name = msg.From.UserName
	}
	if name == "" {
		name = fmt.Sprintf("User%d", msg.From.ID)
	}

	// Limit name length to prevent database issues
	if len(name) > maxNameLength {
		// Truncate at rune boundary to handle multi-byte UTF-8 characters correctly
		runes := []rune(name)
		if len(runes) > maxNameLength {
			name = string(runes[:maxNameLength])
		}
	}

	// Generate secure deterministic password
	password, err := h.generateTelegramPassword(msg.From.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to generate password for Telegram user %d: %v", msg.From.ID, err)
		h.sendMessage(msg.Chat.ID, MsgServiceError, nil)
		return
	}

	// Try to register new user with separate context
	ctx1, cancel1 := context.WithTimeout(context.Background(), registrationTimeout)
	defer cancel1()

	resp, err := h.clients.User.Register(ctx1, &userpb.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})

	// If registration succeeds, create session
	if err == nil && resp != nil && resp.Response != nil && resp.Response.Success {
		now := time.Now()
		h.setSession(msg.Chat.ID, &UserSession{
			UserID:       resp.User.Id,
			Email:        email,
			LinkedAt:     now,
			LastActivity: now,
		})
		log.Printf("[INFO] New user %d registered via Telegram (chat=%d)", resp.User.Id, msg.Chat.ID)
		welcomeMsg := fmt.Sprintf("✅ Welcome, %s!\n\n%s", name, MsgWelcome)
		h.sendMessage(msg.Chat.ID, welcomeMsg, MainMenuKeyboard())
		return
	}

	// If registration failed, try login (user might already exist) with separate context
	ctx2, cancel2 := context.WithTimeout(context.Background(), registrationTimeout)
	defer cancel2()

	loginResp, loginErr := h.clients.User.Login(ctx2, &userpb.LoginRequest{
		Email:    email,
		Password: password,
	})

	if loginErr == nil && loginResp != nil && loginResp.Response != nil && loginResp.Response.Success {
		now := time.Now()
		h.setSession(msg.Chat.ID, &UserSession{
			UserID:       loginResp.User.Id,
			Email:        email,
			LinkedAt:     now,
			LastActivity: now,
		})
		log.Printf("[INFO] Existing user %d logged in via Telegram (chat=%d)", loginResp.User.Id, msg.Chat.ID)
		h.sendMessage(msg.Chat.ID, MsgWelcome, MainMenuKeyboard())
		return
	}

	// Both registration and login failed - log without exposing sensitive details
	log.Printf("[ERROR] Failed to register/login Telegram user %d: registration_failed=%t, login_failed=%t",
		msg.From.ID,
		err != nil,
		loginErr != nil)
	h.sendMessage(msg.Chat.ID, MsgRegistrationFailed, nil)
}

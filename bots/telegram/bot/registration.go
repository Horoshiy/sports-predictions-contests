package bot

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userpb "github.com/sports-prediction-contests/shared/proto/user"
)

// generateSecurePassword creates a cryptographically secure random password
func generateSecurePassword() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// registerViaTelegram registers a new user using Telegram credentials
func (h *Handlers) registerViaTelegram(msg *tgbotapi.Message) {
	// Check if message has a sender
	if msg.From == nil {
		log.Printf("[ERROR] Cannot register: message has no sender")
		h.sendMessage(msg.Chat.ID, "❌ Registration failed: invalid message", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	// Try to login first (user might already exist)
	// For existing users, we use a deterministic password for login attempts
	loginResp, err := h.clients.User.Login(ctx, &userpb.LoginRequest{
		Email:    email,
		Password: fmt.Sprintf("tg_%d", msg.From.ID),
	})

	if err == nil && loginResp != nil && loginResp.Response != nil && loginResp.Response.Success {
		// User already exists, create session
		now := time.Now()
		h.setSession(msg.Chat.ID, &UserSession{
			UserID:       loginResp.User.Id,
			Email:        email,
			LinkedAt:     now,
			LastActivity: now,
		})
		h.sendMessage(msg.Chat.ID, MsgWelcome, MainMenuKeyboard())
		return
	}

	// Generate secure random password for new user
	password, err := generateSecurePassword()
	if err != nil {
		log.Printf("[ERROR] Failed to generate secure password: %v", err)
		h.sendMessage(msg.Chat.ID, "❌ Registration failed: unable to generate credentials", nil)
		return
	}

	// Register new user
	resp, err := h.clients.User.Register(ctx, &userpb.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})

	if err != nil || resp == nil || resp.Response == nil || !resp.Response.Success {
		errMsg := "registration failed"
		if resp != nil && resp.Response != nil {
			errMsg = resp.Response.Message
		}
		log.Printf("[ERROR] Failed to register Telegram user: %v", err)
		h.sendMessage(msg.Chat.ID, fmt.Sprintf("❌ Registration failed: %s", errMsg), nil)
		return
	}

	// Create session
	now := time.Now()
	h.setSession(msg.Chat.ID, &UserSession{
		UserID:       resp.User.Id,
		Email:        email,
		LinkedAt:     now,
		LastActivity: now,
	})

	log.Printf("[INFO] User %d registered via Telegram (chat=%d)", resp.User.Id, msg.Chat.ID)
	welcomeMsg := fmt.Sprintf("✅ Welcome, %s!\n\n%s", name, MsgWelcome)
	h.sendMessage(msg.Chat.ID, welcomeMsg, MainMenuKeyboard())
}

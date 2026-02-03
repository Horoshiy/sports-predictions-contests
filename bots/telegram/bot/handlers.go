package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	contestpb "github.com/sports-prediction-contests/shared/proto/contest"
	notificationpb "github.com/sports-prediction-contests/shared/proto/notification"
	scoringpb "github.com/sports-prediction-contests/shared/proto/scoring"
	userpb "github.com/sports-prediction-contests/shared/proto/user"
	"github.com/sports-prediction-contests/telegram-bot/clients"
	"google.golang.org/grpc/metadata"
)

type Handlers struct {
	api               *tgbotapi.BotAPI
	clients           *clients.Clients
	passwordSecret    []byte
	sessions          map[int64]*UserSession
	mu                sync.RWMutex // Protects sessions map
	registrationLocks sync.Map     // Per-chat registration locks
	sessionTTL        time.Duration
	shutdownCh        chan struct{}
}

type UserSession struct {
	UserID       uint32
	Email        string
	LinkedAt     time.Time
	LastActivity time.Time
	// Navigation state
	CurrentContest uint32
	CurrentPage    int
}

func NewHandlers(api *tgbotapi.BotAPI, clients *clients.Clients, passwordSecret string) *Handlers {
	h := &Handlers{
		api:            api,
		clients:        clients,
		passwordSecret: []byte(passwordSecret),
		sessions:       make(map[int64]*UserSession),
		sessionTTL:     24 * time.Hour,
		shutdownCh:     make(chan struct{}),
	}

	// Start session cleanup goroutine
	go h.cleanupSessions()

	return h
}

// getSession safely retrieves a session and updates LastActivity
// If session doesn't exist, tries to restore it by logging in
func (h *Handlers) getSession(chatID int64) *UserSession {
	h.mu.Lock()
	session := h.sessions[chatID]
	if session != nil {
		session.LastActivity = time.Now()
		h.mu.Unlock()
		return session
	}
	h.mu.Unlock()

	// Try to restore session by logging in with Telegram credentials
	return h.tryRestoreSession(chatID)
}

// tryRestoreSession attempts to restore a session by logging in
func (h *Handlers) tryRestoreSession(chatID int64) *UserSession {
	email := fmt.Sprintf("tg_%d@telegram.bot", chatID)
	
	password, err := h.generateTelegramPassword(chatID)
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	loginResp, err := h.clients.User.Login(ctx, &userpb.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil || loginResp == nil || loginResp.Response == nil || !loginResp.Response.Success {
		return nil
	}

	now := time.Now()
	session := &UserSession{
		UserID:       loginResp.User.Id,
		Email:        email,
		LinkedAt:     now,
		LastActivity: now,
	}
	h.setSession(chatID, session)
	log.Printf("[INFO] Session restored for chat %d (user %d)", chatID, session.UserID)
	return session
}

// setSession safely stores a session
func (h *Handlers) setSession(chatID int64, session *UserSession) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sessions[chatID] = session
}

func (h *Handlers) HandleCommand(msg *tgbotapi.Message) {
	cmd := msg.Command()
	args := msg.CommandArguments()

	switch cmd {
	case "start":
		h.handleStart(msg)
	case "help":
		h.handleHelp(msg)
	case "contests":
		h.handleContests(msg)
	case "leaderboard":
		h.handleLeaderboard(msg, args)
	case "mystats":
		h.handleMyStats(msg)
	case "link":
		h.handleLink(msg, args)
	default:
		h.sendMessage(msg.Chat.ID, MsgUnknownCommand, nil)
	}
}

func (h *Handlers) HandleCallback(cb *tgbotapi.CallbackQuery) {
	data := cb.Data
	chatID := cb.Message.Chat.ID
	msgID := cb.Message.MessageID

	// Acknowledge callback
	if _, err := h.api.Request(tgbotapi.NewCallback(cb.ID, "")); err != nil {
		log.Printf("[WARN] Failed to acknowledge callback %s: %v", cb.ID, err)
	}

	switch {
	case data == "back_main":
		h.editMessage(chatID, msgID, MsgWelcome, MainMenuKeyboard())
	case data == "contests":
		h.handleContestsCallback(chatID, msgID)
	case data == "mystats":
		h.handleMyStatsCallback(chatID, msgID)
	case data == "help":
		h.editMessage(chatID, msgID, MsgHelp, BackToMainKeyboard())
	case strings.HasPrefix(data, "contest_"):
		id, err := strconv.ParseUint(strings.TrimPrefix(data, "contest_"), 10, 32)
		if err != nil {
			h.editMessage(chatID, msgID, "Invalid contest ID.", BackToMainKeyboard())
			return
		}
		h.handleContestDetail(chatID, msgID, uint32(id))
	case strings.HasPrefix(data, "leaderboard_"):
		id, err := strconv.ParseUint(strings.TrimPrefix(data, "leaderboard_"), 10, 32)
		if err != nil {
			h.editMessage(chatID, msgID, "Invalid contest ID.", BackToMainKeyboard())
			return
		}
		h.handleLeaderboardCallback(chatID, msgID, uint32(id))
	case strings.HasPrefix(data, "matches_"):
		// Format: matches_contestID_page
		parts := strings.Split(strings.TrimPrefix(data, "matches_"), "_")
		if len(parts) < 1 {
			log.Printf("[WARN] Invalid matches callback data: %s", data)
			return
		}
		contestID, _ := strconv.ParseUint(parts[0], 10, 32)
		page := 1
		if len(parts) > 1 {
			page, _ = strconv.Atoi(parts[1])
		}
		h.handleMatchList(chatID, msgID, uint32(contestID), page)
	case strings.HasPrefix(data, "match_"):
		// Format: match_matchID
		id, _ := strconv.ParseUint(strings.TrimPrefix(data, "match_"), 10, 32)
		h.handleMatchDetail(chatID, msgID, uint32(id))
	case strings.HasPrefix(data, "p_"):
		// Format: p_matchID_homeScore_awayScore
		parts := strings.Split(strings.TrimPrefix(data, "p_"), "_")
		if len(parts) < 3 {
			log.Printf("[WARN] Invalid prediction callback data: %s", data)
			return
		}
		matchID, _ := strconv.ParseUint(parts[0], 10, 32)
		homeScore, _ := strconv.Atoi(parts[1])
		awayScore, _ := strconv.Atoi(parts[2])
		h.handlePredictionSubmit(chatID, msgID, uint32(matchID), homeScore, awayScore)
	case strings.HasPrefix(data, "pany_"):
		// Format: pany_matchID
		id, _ := strconv.ParseUint(strings.TrimPrefix(data, "pany_"), 10, 32)
		h.handleAnyOtherScore(chatID, msgID, uint32(id))
	}
}

func (h *Handlers) handleStart(msg *tgbotapi.Message) {
	session := h.getSession(msg.Chat.ID)

	// If already registered, show menu
	if session != nil && session.UserID > 0 {
		h.sendMessage(msg.Chat.ID, MsgWelcome, MainMenuKeyboard())
		return
	}

	// Acquire lock for this chat ID to prevent race conditions
	lockInterface, _ := h.registrationLocks.LoadOrStore(msg.Chat.ID, &sync.Mutex{})
	lock, ok := lockInterface.(*sync.Mutex)
	if !ok {
		log.Printf("[ERROR] Invalid lock type for chat %d", msg.Chat.ID)
		h.sendMessage(msg.Chat.ID, MsgServiceError, nil)
		return
	}
	lock.Lock()
	defer func() {
		lock.Unlock()
		// Clean up lock after registration completes (one-time operation)
		h.registrationLocks.Delete(msg.Chat.ID)
	}()

	// Check again after acquiring lock
	session = h.getSession(msg.Chat.ID)
	if session != nil && session.UserID > 0 {
		h.sendMessage(msg.Chat.ID, MsgWelcome, MainMenuKeyboard())
		return
	}

	// Auto-register via Telegram
	h.registerViaTelegram(msg)
}

func (h *Handlers) handleHelp(msg *tgbotapi.Message) {
	h.sendMessage(msg.Chat.ID, MsgHelp, nil)
}

func (h *Handlers) handleContests(msg *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.clients.Contest.ListContests(ctx, &contestpb.ListContestsRequest{
		Status: "active",
	})
	if err != nil || resp == nil {
		log.Printf("[ERROR] Failed to list contests: %v", err)
		h.sendMessage(msg.Chat.ID, MsgServiceError, nil)
		return
	}

	if len(resp.Contests) == 0 {
		h.sendMessage(msg.Chat.ID, MsgNoContests, BackToMainKeyboard())
		return
	}

	text := MsgContestList
	var contests []ContestInfo
	for _, c := range resp.Contests {
		text += FormatContest(c.Id, c.Title, c.SportType, c.Status) + "\n"
		contests = append(contests, ContestInfo{
			ID:        c.Id,
			Title:     c.Title,
			SportType: c.SportType,
			Status:    c.Status,
		})
	}
	h.sendMessage(msg.Chat.ID, text, ContestListKeyboard(contests))
}

func (h *Handlers) handleContestsCallback(chatID int64, msgID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.clients.Contest.ListContests(ctx, &contestpb.ListContestsRequest{
		Status: "active",
	})
	if err != nil || resp == nil {
		log.Printf("[ERROR] Failed to list contests: %v", err)
		h.editMessage(chatID, msgID, MsgServiceError, BackToMainKeyboard())
		return
	}

	if len(resp.Contests) == 0 {
		h.editMessage(chatID, msgID, MsgNoContests, BackToMainKeyboard())
		return
	}

	text := MsgContestList
	var contests []ContestInfo
	for _, c := range resp.Contests {
		text += FormatContest(c.Id, c.Title, c.SportType, c.Status) + "\n"
		contests = append(contests, ContestInfo{
			ID:        c.Id,
			Title:     c.Title,
			SportType: c.SportType,
			Status:    c.Status,
		})
	}
	h.editMessage(chatID, msgID, text, ContestListKeyboard(contests))
}

func (h *Handlers) handleContestDetail(chatID int64, msgID int, contestID uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.clients.Contest.GetContest(ctx, &contestpb.GetContestRequest{Id: contestID})
	if err != nil || resp == nil || resp.Response == nil || !resp.Response.Success {
		h.editMessage(chatID, msgID, "Contest not found.", BackToMainKeyboard())
		return
	}

	// Update session with selected contest
	session := h.getSession(chatID)
	if session != nil {
		session.CurrentContest = contestID
		h.setSession(chatID, session)
	}

	c := resp.Contest
	text := fmt.Sprintf("🏆 <b>%s</b>\n\n%s\n\nSport: %s\nStatus: %s\nParticipants: %d",
		c.Title, c.Description, c.SportType, c.Status, c.CurrentParticipants)
	h.editMessage(chatID, msgID, text, ContestDetailKeyboard(contestID))
}

func (h *Handlers) handleLeaderboard(msg *tgbotapi.Message, args string) {
	contestID := uint32(0)
	if args != "" {
		id, err := strconv.ParseUint(args, 10, 32)
		if err == nil {
			contestID = uint32(id)
		}
	}

	if contestID == 0 {
		h.sendMessage(msg.Chat.ID, "Usage: /leaderboard <contest_id>", nil)
		return
	}

	h.showLeaderboard(msg.Chat.ID, 0, contestID)
}

func (h *Handlers) handleLeaderboardCallback(chatID int64, msgID int, contestID uint32) {
	if contestID == 0 {
		h.editMessage(chatID, msgID, "Select a contest first.", BackToMainKeyboard())
		return
	}
	h.showLeaderboard(chatID, msgID, contestID)
}

func (h *Handlers) showLeaderboard(chatID int64, msgID int, contestID uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.clients.Scoring.GetLeaderboard(ctx, &scoringpb.GetLeaderboardRequest{
		ContestId: contestID,
		Limit:     10,
	})
	if err != nil || resp == nil {
		log.Printf("[ERROR] Failed to get leaderboard: %v", err)
		if msgID > 0 {
			h.editMessage(chatID, msgID, MsgServiceError, BackToMainKeyboard())
		} else {
			h.sendMessage(chatID, MsgServiceError, nil)
		}
		return
	}

	if resp.Leaderboard == nil || len(resp.Leaderboard.Entries) == 0 {
		text := MsgLeaderboard + MsgEmptyLeaderboard
		if msgID > 0 {
			h.editMessage(chatID, msgID, text, BackToMainKeyboard())
		} else {
			h.sendMessage(chatID, text, BackToMainKeyboard())
		}
		return
	}

	text := MsgLeaderboard
	for i, e := range resp.Leaderboard.Entries {
		text += FormatLeaderboardEntry(i+1, e.UserName, e.TotalPoints, e.CurrentStreak)
	}

	if msgID > 0 {
		h.editMessage(chatID, msgID, text, BackToMainKeyboard())
	} else {
		h.sendMessage(chatID, text, BackToMainKeyboard())
	}
}

func (h *Handlers) handleMyStats(msg *tgbotapi.Message) {
	session := h.getSession(msg.Chat.ID)
	if session == nil {
		h.sendMessage(msg.Chat.ID, MsgNotLinked, nil)
		return
	}
	h.showUserStats(msg.Chat.ID, 0, session.UserID)
}

func (h *Handlers) handleMyStatsCallback(chatID int64, msgID int) {
	session := h.getSession(chatID)
	if session == nil {
		h.editMessage(chatID, msgID, MsgNotLinked, BackToMainKeyboard())
		return
	}
	h.showUserStats(chatID, msgID, session.UserID)
}

func (h *Handlers) showUserStats(chatID int64, msgID int, userID uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Add user_id to gRPC metadata for bot authentication
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", strconv.FormatUint(uint64(userID), 10))

	resp, err := h.clients.Scoring.GetUserAnalytics(ctx, &scoringpb.GetUserAnalyticsRequest{
		UserId:    userID,
		TimeRange: "all",
	})
	if err != nil || resp == nil {
		log.Printf("[ERROR] Failed to get user analytics: %v", err)
		if msgID > 0 {
			h.editMessage(chatID, msgID, MsgServiceError, BackToMainKeyboard())
		} else {
			h.sendMessage(chatID, MsgServiceError, nil)
		}
		return
	}

	totalPoints := 0.0
	if resp.Analytics != nil {
		totalPoints = resp.Analytics.TotalPoints
	}

	text := fmt.Sprintf(MsgStats, totalPoints, 0, 0)
	if msgID > 0 {
		h.editMessage(chatID, msgID, text, BackToMainKeyboard())
	} else {
		h.sendMessage(chatID, text, BackToMainKeyboard())
	}
}

func (h *Handlers) handleLink(msg *tgbotapi.Message, args string) {
	parts := strings.Fields(args)
	if len(parts) != 2 {
		h.sendMessage(msg.Chat.ID, MsgLinkUsage, nil)
		return
	}

	email, password := parts[0], parts[1]

	// Delete the message with password for security
	if _, err := h.api.Request(tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)); err != nil {
		log.Printf("[WARN] Failed to delete message with password: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Login to verify credentials
	loginResp, err := h.clients.User.Login(ctx, &userpb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil || loginResp == nil || loginResp.Response == nil || !loginResp.Response.Success {
		errMsg := "invalid credentials"
		if loginResp != nil && loginResp.Response != nil {
			errMsg = loginResp.Response.Message
		}
		h.sendMessage(msg.Chat.ID, fmt.Sprintf(MsgLinkFailed, errMsg), nil)
		return
	}

	// Update notification preferences with chat ID
	_, err = h.clients.Notification.UpdatePreference(ctx, &notificationpb.UpdatePreferenceRequest{
		UserId:         loginResp.User.Id,
		Channel:        notificationpb.NotificationChannel_TELEGRAM,
		Enabled:        true,
		TelegramChatId: msg.Chat.ID,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to update notification preference: %v", err)
	}

	// Store session (thread-safe)
	now := time.Now()
	h.setSession(msg.Chat.ID, &UserSession{
		UserID:       loginResp.User.Id,
		Email:        email,
		LinkedAt:     now,
		LastActivity: now,
	})

	log.Printf("[INFO] Session created (chat=%d, user=%d)", msg.Chat.ID, loginResp.User.Id)
	h.sendMessage(msg.Chat.ID, MsgLinkSuccess, MainMenuKeyboard())
}

func (h *Handlers) sendMessage(chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
		msg.ReplyMarkup = kb
	}
	if _, err := h.api.Send(msg); err != nil {
		log.Printf("[ERROR] Failed to send message: %v", err)
	}
}

func (h *Handlers) editMessage(chatID int64, msgID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	edit := tgbotapi.NewEditMessageText(chatID, msgID, text)
	edit.ParseMode = "HTML"
	edit.ReplyMarkup = &keyboard
	if _, err := h.api.Send(edit); err != nil {
		log.Printf("[ERROR] Failed to edit message: %v", err)
	}
}

// cleanupSessions periodically removes expired sessions to prevent memory leak
func (h *Handlers) cleanupSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.mu.Lock()
			now := time.Now()
			cleaned := 0
			for chatID, session := range h.sessions {
				if now.Sub(session.LastActivity) > h.sessionTTL {
					delete(h.sessions, chatID)
					cleaned++
				}
			}
			h.mu.Unlock()
			if cleaned > 0 {
				log.Printf("[INFO] Cleaned up %d expired sessions", cleaned)
			}
		case <-h.shutdownCh:
			log.Printf("[INFO] Session cleanup goroutine stopped")
			return
		}
	}
}

// Shutdown gracefully stops the handlers and cleanup goroutines
func (h *Handlers) Shutdown() {
	close(h.shutdownCh)
}

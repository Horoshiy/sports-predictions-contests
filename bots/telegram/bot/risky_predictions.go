package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	predictionpb "github.com/sports-prediction-contests/shared/proto/prediction"
)

// NOTE: Types below are duplicated from backend/shared/scoring/rules.go
// This is intentional because the bot is a separate Go module and importing
// the shared package would require complex module dependencies.
// Keep these types in sync with the backend definitions.

// RiskyEvent represents a risky event for prediction
type RiskyEvent struct {
	Slug   string  `json:"slug"`
	Name   string  `json:"name"`
	NameEn string  `json:"name_en,omitempty"`
	Points float64 `json:"points"`
}

// RiskyScoringRules represents risky contest rules
type RiskyScoringRules struct {
	MaxSelections int          `json:"max_selections"`
	Events        []RiskyEvent `json:"events,omitempty"`
}

// ContestRules represents contest scoring rules
type ContestRules struct {
	Type   string             `json:"type"`
	Risky  *RiskyScoringRules `json:"risky,omitempty"`
}

// Fallback risky events (used if API fails)
var fallbackRiskyEvents = []RiskyEvent{
	{Slug: "penalty", Name: "⚽ Пенальти", NameEn: "Penalty", Points: 3},
	{Slug: "red_card", Name: "🟥 Удаление", NameEn: "Red card", Points: 4},
	{Slug: "own_goal", Name: "🔙 Автогол", NameEn: "Own goal", Points: 5},
	{Slug: "hat_trick", Name: "🎩 Хет-трик", NameEn: "Hat-trick", Points: 6},
	{Slug: "clean_sheet_home", Name: "🏠 Хозяева на ноль", NameEn: "Home clean sheet", Points: 2},
}

// RiskyEventsCache caches API responses for risky events
type RiskyEventsCache struct {
	mu             sync.RWMutex
	globalEvents   []RiskyEvent
	globalExpiry   time.Time
	matchEvents    map[string]matchEventsCacheEntry // key: "eventId:contestId"
	cacheDuration  time.Duration
}

type matchEventsCacheEntry struct {
	events        []RiskyEvent
	maxSelections int
	expiry        time.Time
}

var riskyEventsCache = &RiskyEventsCache{
	matchEvents:   make(map[string]matchEventsCacheEntry),
	cacheDuration: 5 * time.Minute,
}

func init() {
	// Start cache cleanup goroutine
	go riskyEventsCache.startCleanup()
}

// startCleanup periodically removes expired cache entries to prevent memory leak
func (c *RiskyEventsCache) startCleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.matchEvents {
			if now.After(entry.expiry) {
				delete(c.matchEvents, key)
			}
		}
		c.mu.Unlock()
	}
}

// fetchGlobalRiskyEvents fetches all risky event types from API
func (h *Handlers) fetchGlobalRiskyEvents() ([]RiskyEvent, error) {
	riskyEventsCache.mu.RLock()
	if time.Now().Before(riskyEventsCache.globalExpiry) && len(riskyEventsCache.globalEvents) > 0 {
		events := riskyEventsCache.globalEvents
		riskyEventsCache.mu.RUnlock()
		return events, nil
	}
	riskyEventsCache.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.clients.Prediction.ListRiskyEventTypes(ctx, &predictionpb.ListRiskyEventTypesRequest{
		IncludeInactive: false, // Only fetch active events
	})
	if err != nil {
		log.Printf("[WARN] Failed to fetch risky events from API: %v", err)
		return fallbackRiskyEvents, err
	}

	events := make([]RiskyEvent, 0, len(resp.EventTypes))
	for _, e := range resp.EventTypes {
		events = append(events, RiskyEvent{
			Slug:   e.Slug,
			Name:   fmt.Sprintf("%s %s", e.Icon, e.Name),
			NameEn: e.NameEn,
			Points: e.DefaultPoints,
		})
	}

	// Update cache
	riskyEventsCache.mu.Lock()
	riskyEventsCache.globalEvents = events
	riskyEventsCache.globalExpiry = time.Now().Add(riskyEventsCache.cacheDuration)
	riskyEventsCache.mu.Unlock()

	return events, nil
}

// fetchMatchRiskyEvents fetches risky events for a specific match (with contest overrides)
func (h *Handlers) fetchMatchRiskyEvents(eventID, contestID uint32) ([]RiskyEvent, int, error) {
	cacheKey := fmt.Sprintf("%d:%d", eventID, contestID)

	riskyEventsCache.mu.RLock()
	if entry, ok := riskyEventsCache.matchEvents[cacheKey]; ok && time.Now().Before(entry.expiry) {
		riskyEventsCache.mu.RUnlock()
		return entry.events, entry.maxSelections, nil
	}
	riskyEventsCache.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.clients.Prediction.GetMatchRiskyEvents(ctx, &predictionpb.GetMatchRiskyEventsRequest{
		EventId:   eventID,
		ContestId: contestID,
	})
	if err != nil {
		log.Printf("[WARN] Failed to fetch match risky events from API: %v", err)
		// Fall back to global events
		globalEvents, _ := h.fetchGlobalRiskyEvents()
		return globalEvents, 5, err
	}

	events := make([]RiskyEvent, 0, len(resp.Events))
	for _, e := range resp.Events {
		if !e.IsEnabled {
			continue // Skip disabled events
		}
		events = append(events, RiskyEvent{
			Slug:   e.Slug,
			Name:   fmt.Sprintf("%s %s", e.Icon, e.Name),
			NameEn: e.NameEn,
			Points: e.Points,
		})
	}

	maxSelections := int(resp.MaxSelections)
	if maxSelections == 0 {
		maxSelections = 5
	}

	// Update cache
	riskyEventsCache.mu.Lock()
	riskyEventsCache.matchEvents[cacheKey] = matchEventsCacheEntry{
		events:        events,
		maxSelections: maxSelections,
		expiry:        time.Now().Add(riskyEventsCache.cacheDuration),
	}
	riskyEventsCache.mu.Unlock()

	return events, maxSelections, nil
}

// parseContestRules parses contest rules JSON
func parseContestRules(rulesJSON string) *ContestRules {
	if rulesJSON == "" {
		return &ContestRules{Type: "standard"}
	}
	var rules ContestRules
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return &ContestRules{Type: "standard"}
	}
	if rules.Type == "" {
		rules.Type = "standard"
	}
	return &rules
}

// isRiskyContest checks if contest has risky type
func isRiskyContest(rulesJSON string) bool {
	rules := parseContestRules(rulesJSON)
	return rules.Type == "risky"
}

// getRiskyEvents returns risky events for a contest (fallback when API not available)
func getRiskyEvents(rulesJSON string) []RiskyEvent {
	rules := parseContestRules(rulesJSON)
	if rules.Risky != nil && len(rules.Risky.Events) > 0 {
		return rules.Risky.Events
	}
	return fallbackRiskyEvents
}

// getMaxSelections returns max selections for risky contest
func getMaxSelections(rulesJSON string) int {
	rules := parseContestRules(rulesJSON)
	if rules.Risky != nil && rules.Risky.MaxSelections > 0 {
		return rules.Risky.MaxSelections
	}
	return 5
}

// RiskyEventsKeyboard creates keyboard for selecting risky events
func RiskyEventsKeyboard(matchID uint32, selectedSlugs []string, events []RiskyEvent, maxSelections int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, event := range events {
		isSelected := contains(selectedSlugs, event.Slug)
		
		// Format button text
		var btnText string
		if isSelected {
			btnText = fmt.Sprintf("✅ %s (+%.0f/−%.0f)", event.Name, event.Points, event.Points)
		} else {
			btnText = fmt.Sprintf("⬜ %s (+%.0f/−%.0f)", event.Name, event.Points, event.Points)
		}
		
		callbackData := fmt.Sprintf("risky_%d_%s", matchID, event.Slug)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnText, callbackData),
		))
	}

	// Info row
	selectedCount := len(selectedSlugs)
	infoText := fmt.Sprintf("📊 Выбрано: %d/%d", selectedCount, maxSelections)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(infoText, "noop"),
	))

	// Submit button (only if selections made)
	if selectedCount > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить прогноз", fmt.Sprintf("risky_submit_%d", matchID)),
		))
	}

	// Back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Назад", "back_to_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// formatRiskyPrediction formats risky prediction for display
func formatRiskyPrediction(selections []string, events []RiskyEvent) string {
	if len(selections) == 0 {
		return "нет"
	}
	
	var names []string
	for _, slug := range selections {
		for _, event := range events {
			if event.Slug == slug {
				names = append(names, event.Name)
				break
			}
		}
	}
	return strings.Join(names, ", ")
}

// contains checks if slice contains string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// toggleSelection adds or removes selection
func toggleSelection(selections []string, slug string, maxSelections int) []string {
	if contains(selections, slug) {
		// Remove
		result := make([]string, 0, len(selections)-1)
		for _, s := range selections {
			if s != slug {
				result = append(result, s)
			}
		}
		return result
	}
	
	// Add if under limit
	if len(selections) < maxSelections {
		return append(selections, slug)
	}
	return selections
}

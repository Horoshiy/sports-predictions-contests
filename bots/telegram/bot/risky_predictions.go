package bot

import (
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

// Default risky events
var defaultRiskyEvents = []RiskyEvent{
	{Slug: "penalty", Name: "⚽ Пенальти", NameEn: "Penalty", Points: 3},
	{Slug: "red_card", Name: "🟥 Удаление", NameEn: "Red card", Points: 4},
	{Slug: "own_goal", Name: "🔙 Автогол", NameEn: "Own goal", Points: 5},
	{Slug: "hat_trick", Name: "🎩 Хет-трик", NameEn: "Hat-trick", Points: 6},
	{Slug: "clean_sheet_home", Name: "🏠 Хозяева на ноль", NameEn: "Home clean sheet", Points: 2},
	{Slug: "clean_sheet_away", Name: "✈️ Гости на ноль", NameEn: "Away clean sheet", Points: 3},
	{Slug: "both_teams_score", Name: "⚽⚽ Обе забьют", NameEn: "Both teams score", Points: 2},
	{Slug: "over_3_goals", Name: "📈 Больше 3 голов", NameEn: "Over 3.5 goals", Points: 2},
	{Slug: "first_half_draw", Name: "🤝 Ничья в 1-м тайме", NameEn: "First half draw", Points: 2},
	{Slug: "comeback", Name: "🔄 Камбэк", NameEn: "Comeback", Points: 7},
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

// getRiskyEvents returns risky events for a contest
func getRiskyEvents(rulesJSON string) []RiskyEvent {
	rules := parseContestRules(rulesJSON)
	if rules.Risky != nil && len(rules.Risky.Events) > 0 {
		return rules.Risky.Events
	}
	return defaultRiskyEvents
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

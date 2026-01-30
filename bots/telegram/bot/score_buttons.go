package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ScorePredictionKeyboard creates score prediction buttons in 3-column layout
// Note: Callback data format "p_{matchID}_{home}_{away}" must stay under 64 bytes (Telegram limit)
// Current format supports match IDs up to ~10^15 safely (uint32 max is ~4.3 billion)
// Row 1: 0-0, 1-1, 2-2
// Row 2: 1-0, 2-0, 2-1
// Row 3: 3-0, 3-1, 3-2
// Row 4: 0-1, 0-2, 1-2
// Row 5: 0-3, 1-3, 2-3
// Row 6: Any Other (full width)
// Row 7: Back button
func ScorePredictionKeyboard(matchID uint32) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		// Row 1: Draws
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0-0", fmt.Sprintf("p_%d_0_0", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("1-1", fmt.Sprintf("p_%d_1_1", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("2-2", fmt.Sprintf("p_%d_2_2", matchID)),
		),
		// Row 2: Home wins (low)
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1-0", fmt.Sprintf("p_%d_1_0", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("2-0", fmt.Sprintf("p_%d_2_0", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("2-1", fmt.Sprintf("p_%d_2_1", matchID)),
		),
		// Row 3: Home wins (high)
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3-0", fmt.Sprintf("p_%d_3_0", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("3-1", fmt.Sprintf("p_%d_3_1", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("3-2", fmt.Sprintf("p_%d_3_2", matchID)),
		),
		// Row 4: Away wins (low)
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0-1", fmt.Sprintf("p_%d_0_1", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("0-2", fmt.Sprintf("p_%d_0_2", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("1-2", fmt.Sprintf("p_%d_1_2", matchID)),
		),
		// Row 5: Away wins (high)
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0-3", fmt.Sprintf("p_%d_0_3", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("1-3", fmt.Sprintf("p_%d_1_3", matchID)),
			tgbotapi.NewInlineKeyboardButtonData("2-3", fmt.Sprintf("p_%d_2_3", matchID)),
		),
		// Row 6: Any other score (full width)
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Any Other Score", fmt.Sprintf("pany_%d", matchID)),
		),
		// Row 7: Back button
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Back", fmt.Sprintf("match_%d", matchID)),
		),
	)
}

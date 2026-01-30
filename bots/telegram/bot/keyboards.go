package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 Contests", "contests"),
			tgbotapi.NewInlineKeyboardButtonData("🏅 Leaderboard", "leaderboard"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 My Stats", "mystats"),
			tgbotapi.NewInlineKeyboardButtonData("❓ Help", "help"),
		),
	)
}

func ContestListKeyboard(contests []ContestInfo) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range contests {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("🏆 %s", c.Title),
				fmt.Sprintf("contest_%d", c.ID),
			),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Back", "back_main"),
	))
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func ContestDetailKeyboard(contestID uint32) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚽ Matches", fmt.Sprintf("matches_%d_1", contestID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏅 Leaderboard", fmt.Sprintf("leaderboard_%d", contestID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Back", "contests"),
		),
	)
}

func BackToMainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Back to Menu", "back_main"),
		),
	)
}

type ContestInfo struct {
	ID        uint32
	Title     string
	SportType string
	Status    string
}

package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 Конкурсы", "contests"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Моя статистика", "mystats"),
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
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
		tgbotapi.NewInlineKeyboardButtonData("« Назад", "back_main"),
	))
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func ContestDetailKeyboard(contestID uint32) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚽ Матчи", fmt.Sprintf("matches_%d_1", contestID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏅 Таблица лидеров", fmt.Sprintf("leaderboard_%d", contestID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", "contests"),
		),
	)
}

func BackToMainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« В меню", "back_main"),
		),
	)
}

type ContestInfo struct {
	ID        uint32
	Title     string
	SportType string
	Status    string
}

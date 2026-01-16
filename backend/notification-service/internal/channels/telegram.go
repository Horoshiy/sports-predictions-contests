package channels

import (
	"fmt"
	"html"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramChannel struct {
	bot     *tgbotapi.BotAPI
	enabled bool
}

func NewTelegramChannel(token string, enabled bool) (*TelegramChannel, error) {
	if !enabled || token == "" {
		return &TelegramChannel{enabled: false}, nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	log.Printf("Telegram bot authorized: %s", bot.Self.UserName)
	return &TelegramChannel{bot: bot, enabled: true}, nil
}

func (t *TelegramChannel) Send(chatID int64, title, message string) error {
	if !t.enabled || t.bot == nil {
		return nil
	}

	text := fmt.Sprintf("<b>%s</b>\n\n%s", html.EscapeString(title), html.EscapeString(message))
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	return nil
}

func (t *TelegramChannel) IsEnabled() bool {
	return t.enabled
}

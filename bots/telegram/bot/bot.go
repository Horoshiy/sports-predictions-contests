package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sports-prediction-contests/telegram-bot/clients"
	"github.com/sports-prediction-contests/telegram-bot/config"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	handlers *Handlers
	stop     chan struct{}
}

func New(cfg *config.Config, clients *clients.Clients) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	// Log password secret configuration status (without exposing the secret)
	if cfg.TelegramPasswordSecret == "" {
		log.Printf("[WARN] TELEGRAM_PASSWORD_SECRET not configured - user registration will fail")
	} else {
		log.Printf("[INFO] TELEGRAM_PASSWORD_SECRET configured (length: %d bytes)", len(cfg.TelegramPasswordSecret))
	}

	bot := &Bot{
		api:      api,
		handlers: NewHandlers(api, clients, cfg.TelegramPasswordSecret),
		stop:     make(chan struct{}),
	}

	// Register commands with Telegram
	if err := bot.registerCommands(); err != nil {
		log.Printf("[WARN] Failed to register commands: %v", err)
		// Don't fail bot startup if command registration fails
	}

	return bot, nil
}

// registerCommands registers bot commands with Telegram API
func (b *Bot) registerCommands() error {
	commands := []tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "Start bot and create account | Начать работу и создать аккаунт",
		},
		{
			Command:     "contests",
			Description: "View active contests | Просмотр активных конкурсов",
		},
		{
			Command:     "leaderboard",
			Description: "View contest leaderboard | Таблица лидеров конкурса",
		},
		{
			Command:     "mystats",
			Description: "Your prediction statistics | Ваша статистика прогнозов",
		},
		{
			Command:     "help",
			Description: "Show help message | Показать справку",
		},
		{
			Command:     "link",
			Description: "Link existing account | Привязать существующий аккаунт",
		},
	}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	_, err := b.api.Request(cfg)
	if err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	log.Printf("[INFO] Successfully registered %d bot commands", len(commands))
	return nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Bot started, listening for updates...")

	for {
		select {
		case <-b.stop:
			log.Println("Bot stopped")
			return
		case update := <-updates:
			b.handleUpdate(update)
		}
	}
}

func (b *Bot) Stop() {
	close(b.stop)
	b.api.StopReceivingUpdates()
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil && update.Message.IsCommand() {
		b.handlers.HandleCommand(update.Message)
		return
	}

	if update.CallbackQuery != nil {
		b.handlers.HandleCallback(update.CallbackQuery)
		return
	}
}

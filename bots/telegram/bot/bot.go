package bot

import (
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

	return &Bot{
		api:      api,
		handlers: NewHandlers(api, clients),
		stop:     make(chan struct{}),
	}, nil
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

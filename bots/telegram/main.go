package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/telegram-bot/bot"
	"github.com/sports-prediction-contests/telegram-bot/clients"
	"github.com/sports-prediction-contests/telegram-bot/config"
)

func main() {
	cfg := config.Load()

	if cfg.TelegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	grpcClients, err := clients.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}
	defer grpcClients.Close()

	b, err := bot.New(cfg, grpcClients)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		b.Stop()
	}()

	b.Start()
}

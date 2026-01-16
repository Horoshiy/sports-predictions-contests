package config

import "os"

type Config struct {
	TelegramBotToken            string
	UserServiceEndpoint         string
	ContestServiceEndpoint      string
	PredictionServiceEndpoint   string
	ScoringServiceEndpoint      string
	NotificationServiceEndpoint string
	LogLevel                    string
}

func Load() *Config {
	return &Config{
		TelegramBotToken:            getEnvOrDefault("TELEGRAM_BOT_TOKEN", ""),
		UserServiceEndpoint:         getEnvOrDefault("USER_SERVICE_ENDPOINT", "localhost:8084"),
		ContestServiceEndpoint:      getEnvOrDefault("CONTEST_SERVICE_ENDPOINT", "localhost:8085"),
		PredictionServiceEndpoint:   getEnvOrDefault("PREDICTION_SERVICE_ENDPOINT", "localhost:8086"),
		ScoringServiceEndpoint:      getEnvOrDefault("SCORING_SERVICE_ENDPOINT", "localhost:8087"),
		NotificationServiceEndpoint: getEnvOrDefault("NOTIFICATION_SERVICE_ENDPOINT", "localhost:8089"),
		LogLevel:                    getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

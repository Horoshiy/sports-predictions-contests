package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	JWTSecret string
	LogLevel  string

	TelegramBotToken string
	TelegramEnabled  bool

	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	EmailEnabled bool

	WorkerPoolSize int
}

func Load() *Config {
	return &Config{
		Port:             getEnvOrDefault("NOTIFICATION_SERVICE_PORT", "8089"),
		JWTSecret:        getEnvOrDefault("JWT_SECRET", "your_jwt_secret_key_here"),
		LogLevel:         getEnvOrDefault("LOG_LEVEL", "info"),
		TelegramBotToken: getEnvOrDefault("TELEGRAM_BOT_TOKEN", ""),
		TelegramEnabled:  getEnvOrDefault("TELEGRAM_ENABLED", "false") == "true",
		SMTPHost:         getEnvOrDefault("SMTP_HOST", ""),
		SMTPPort:         getEnvOrDefault("SMTP_PORT", "587"),
		SMTPUser:         getEnvOrDefault("SMTP_USER", ""),
		SMTPPassword:     getEnvOrDefault("SMTP_PASSWORD", ""),
		SMTPFrom:         getEnvOrDefault("SMTP_FROM", ""),
		EmailEnabled:     getEnvOrDefault("EMAIL_ENABLED", "false") == "true",
		WorkerPoolSize:   parseIntOrDefault(getEnvOrDefault("WORKER_POOL_SIZE", "5"), 5),
	}
}

func (c *Config) Validate() error {
	if c.Port == "" {
		c.Port = "8089"
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseIntOrDefault(value string, defaultValue int) int {
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

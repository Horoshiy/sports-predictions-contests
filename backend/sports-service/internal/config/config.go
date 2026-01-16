package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
	LogLevel    string
}

func Load() *Config {
	return &Config{
		Port:        getEnvOrDefault("SPORTS_SERVICE_PORT", "8088"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "your_jwt_secret_key_here"),
		DatabaseURL: getEnvOrDefault("DATABASE_URL", ""),
		LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

func (c *Config) Validate() error {
	if os.Getenv("ENV") == "production" && (c.JWTSecret == "" || c.JWTSecret == "your_jwt_secret_key_here") {
		return errors.New("JWT_SECRET must be set in production")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

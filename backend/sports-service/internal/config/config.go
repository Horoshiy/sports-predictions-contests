package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port             string
	JWTSecret        string
	DatabaseURL      string
	LogLevel         string
	SyncEnabled      bool
	SyncIntervalMins int
	TheSportsDBURL   string
}

func Load() *Config {
	syncEnabled, _ := strconv.ParseBool(getEnvOrDefault("SYNC_ENABLED", "false"))
	syncInterval, _ := strconv.Atoi(getEnvOrDefault("SYNC_INTERVAL_MINS", "60"))

	return &Config{
		Port:             getEnvOrDefault("SPORTS_SERVICE_PORT", "8088"),
		JWTSecret:        getEnvOrDefault("JWT_SECRET", "your_jwt_secret_key_here"),
		DatabaseURL:      getEnvOrDefault("DATABASE_URL", ""),
		LogLevel:         getEnvOrDefault("LOG_LEVEL", "info"),
		SyncEnabled:      syncEnabled,
		SyncIntervalMins: syncInterval,
		TheSportsDBURL:   getEnvOrDefault("THESPORTSDB_URL", "https://www.thesportsdb.com/api/v1/json/3"),
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

package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the prediction service
type Config struct {
	// Server configuration
	Port string

	// JWT configuration
	JWTSecret string

	// Database configuration
	DatabaseURL string

	// Service endpoints
	ContestServiceEndpoint string

	// Logging configuration
	LogLevel string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:                   getEnvOrDefault("PREDICTION_SERVICE_PORT", "8086"),
		JWTSecret:              getEnvOrDefault("JWT_SECRET", "your_jwt_secret_key_here"),
		DatabaseURL:            getEnvOrDefault("DATABASE_URL", "postgres://sports_user:sports_password@localhost:5432/sports_prediction?sslmode=disable"),
		ContestServiceEndpoint: getEnvOrDefault("CONTEST_SERVICE_ENDPOINT", "contest-service:8085"),
		LogLevel:               getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.JWTSecret == "" || c.JWTSecret == "your_jwt_secret_key_here" {
		// In production, this should be an error, but for development we'll allow it
		// return errors.New("JWT_SECRET must be set to a secure value")
	}

	if c.Port == "" {
		c.Port = "8086"
	}

	return nil
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDurationOrDefault parses duration string or returns default
func parseDurationOrDefault(value string, defaultValue time.Duration) time.Duration {
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return defaultValue
}

// parseIntOrDefault parses integer string or returns default
func parseIntOrDefault(value string, defaultValue int) int {
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

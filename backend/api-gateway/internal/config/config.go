package config

import (
	"errors"
	"os"
	"time"
)

// Config holds all configuration for the API Gateway service
type Config struct {
	// Server configuration
	Port string

	// Service endpoints
	UserService       string
	ContestService    string
	PredictionService string

	// JWT configuration
	JWTSecret string

	// CORS configuration
	AllowedOrigins string

	// Logging configuration
	LogLevel string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:              getEnvOrDefault("API_GATEWAY_PORT", "8080"),
		UserService:       getEnvOrDefault("USER_SERVICE_ENDPOINT", "user-service:8084"),
		ContestService:    getEnvOrDefault("CONTEST_SERVICE_ENDPOINT", "contest-service:8085"),
		PredictionService: getEnvOrDefault("PREDICTION_SERVICE_ENDPOINT", "prediction-service:8086"),
		JWTSecret:         getEnvOrDefault("JWT_SECRET", "your_jwt_secret_key_here"),
		AllowedOrigins:    getEnvOrDefault("CORS_ALLOWED_ORIGINS", "*"),
		LogLevel:          getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Check JWT secret in production
	if os.Getenv("ENV") == "production" && (c.JWTSecret == "" || c.JWTSecret == "your_jwt_secret_key_here") {
		return errors.New("JWT_SECRET must be set to a secure value in production")
	}

	if c.Port == "" {
		c.Port = "8080"
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

package config_test

import (
	"os"
	"testing"

	"github.com/sports-prediction-contests/api-gateway/internal/config"
)

func TestConfig_JWTValidation(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		jwtSecret   string
		expectError bool
	}{
		{"Development with default secret", "development", "your_jwt_secret_key_here", false},
		{"Production with default secret", "production", "your_jwt_secret_key_here", true},
		{"Production with secure secret", "production", "secure_random_secret_key_123", false},
		{"No env set with default secret", "", "your_jwt_secret_key_here", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment
			if tt.env != "" {
				os.Setenv("ENV", tt.env)
				defer os.Unsetenv("ENV")
			}

			cfg := &config.Config{
				Port:      "8080",
				JWTSecret: tt.jwtSecret,
			}

			err := cfg.Validate()
			
			if tt.expectError && err == nil {
				t.Error("Expected validation error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

func TestConfig_CORSConfiguration(t *testing.T) {
	os.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com")
	defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

	cfg := config.Load()
	
	if cfg.AllowedOrigins != "https://example.com" {
		t.Errorf("Expected AllowedOrigins to be 'https://example.com', got '%s'", cfg.AllowedOrigins)
	}
}

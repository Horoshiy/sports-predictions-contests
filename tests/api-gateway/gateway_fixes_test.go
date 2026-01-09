package gateway_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sports-prediction-contests/api-gateway/internal/config"
	"github.com/sports-prediction-contests/api-gateway/internal/gateway"
)

func TestErrorResponse_NoDuplication(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		UserService:    "localhost:8084",
		ContestService: "localhost:8085",
		JWTSecret:      "test_secret",
		AllowedOrigins: "*",
	}

	server, err := gateway.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test that error responses don't have duplicate information
	req := httptest.NewRequest("GET", "/v1/nonexistent", nil)
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	var errorResp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&errorResp); err == nil {
		// Check that error and message fields are different
		if errorField, ok := errorResp["error"].(string); ok {
			if messageField, ok := errorResp["message"].(string); ok {
				if errorField == messageField && errorField != "" {
					t.Error("Error and Message fields should not be identical")
				}
			}
		}
	}
}

func TestHealthEndpoint_WithCORS(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		UserService:    "localhost:8084",
		ContestService: "localhost:8085",
		JWTSecret:      "test_secret",
		AllowedOrigins: "https://example.com",
	}

	server, err := gateway.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check CORS header is set correctly
	corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "https://example.com" {
		t.Errorf("Expected CORS origin 'https://example.com', got '%s'", corsOrigin)
	}
}

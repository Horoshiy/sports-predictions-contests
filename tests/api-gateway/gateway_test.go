package gateway_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sports-prediction-contests/api-gateway/internal/config"
	"github.com/sports-prediction-contests/api-gateway/internal/gateway"
)

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		UserService:    "localhost:8084",
		ContestService: "localhost:8085",
		JWTSecret:      "test_secret",
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

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

func TestCORSHeaders(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		UserService:    "localhost:8084",
		ContestService: "localhost:8085",
		JWTSecret:      "test_secret",
	}

	server, err := gateway.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req := httptest.NewRequest("OPTIONS", "/v1/auth/login", nil)
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}

	corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "*" {
		t.Errorf("Expected CORS origin '*', got '%s'", corsOrigin)
	}

	corsMethods := w.Header().Get("Access-Control-Allow-Methods")
	if corsMethods == "" {
		t.Error("Expected CORS methods header to be set")
	}
}

func TestAuthenticationRequired(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		UserService:    "localhost:8084",
		ContestService: "localhost:8085",
		JWTSecret:      "test_secret",
	}

	server, err := gateway.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test protected endpoint without auth
	req := httptest.NewRequest("GET", "/v1/users/profile", nil)
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for protected endpoint, got %d", w.Code)
	}
}

func TestPublicEndpointsNoAuth(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		UserService:    "localhost:8084",
		ContestService: "localhost:8085",
		JWTSecret:      "test_secret",
	}

	server, err := gateway.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test public login endpoint (should not require auth but will fail due to no backend)
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	// Should not be 401 (unauthorized), but may be 503 (service unavailable) due to no backend
	if w.Code == http.StatusUnauthorized {
		t.Errorf("Login endpoint should not require authentication, got status %d", w.Code)
	}
}

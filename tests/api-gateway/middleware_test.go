package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sports-prediction-contests/api-gateway/internal/middleware"
)

func TestJWTMiddleware_AuthBypass(t *testing.T) {
	secret := []byte("test_secret")
	middleware := middleware.JWTMiddleware(secret)
	
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{"Health endpoint should bypass auth", "/health", http.StatusOK},
		{"Auth register should bypass auth", "/v1/auth/register", http.StatusOK},
		{"Auth login should bypass auth", "/v1/auth/login", http.StatusOK},
		{"Protected endpoint should require auth", "/v1/users/profile", http.StatusUnauthorized},
		{"Malicious path should not bypass", "/something/auth/malicious", http.StatusUnauthorized},
		{"Unhealthy path should not bypass", "/unhealthy", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d for path %s, got %d", tt.expectedStatus, tt.path, w.Code)
			}
		})
	}
}

func TestCORSMiddleware_ConfigurableOrigins(t *testing.T) {
	tests := []struct {
		name            string
		allowedOrigins  string
		expectedOrigins string
	}{
		{"Default wildcard", "", "*"},
		{"Specific domain", "https://example.com", "https://example.com"},
		{"Multiple domains", "https://app.com,https://admin.com", "https://app.com,https://admin.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := middleware.CORSMiddleware(tt.allowedOrigins)
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if corsOrigin != tt.expectedOrigins {
				t.Errorf("Expected CORS origin '%s', got '%s'", tt.expectedOrigins, corsOrigin)
			}
		})
	}
}

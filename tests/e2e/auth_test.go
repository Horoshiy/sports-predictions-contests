//go:build e2e

package e2e

import (
	"net/http"
	"testing"
)

// TestAuthFlow tests the authentication workflow.
// Note: Subtests share state and MUST run sequentially (no t.Parallel()).
func TestAuthFlow(t *testing.T) {
	var registeredEmail, registeredPassword, authToken string
	var userID uint

	t.Run("RegisterUser", func(t *testing.T) {
		registeredEmail = generateTestEmail()
		registeredPassword = "testpassword123"

		body := map[string]string{
			"email":    registeredEmail,
			"password": registeredPassword,
			"name":     "Test User",
		}

		resp, err := makeRequest("POST", "/v1/auth/register", body, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[AuthResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if result.Token == "" {
			t.Error("Expected token in response")
		}
		if result.User.Email != registeredEmail {
			t.Errorf("Expected email %s, got %s", registeredEmail, result.User.Email)
		}

		authToken = result.Token
		userID = result.User.ID
		if userID == 0 {
			t.Error("Expected non-zero user ID")
		}
		t.Logf("Registered user ID: %d", userID)
	})

	t.Run("RegisterDuplicateEmail", func(t *testing.T) {
		body := map[string]string{
			"email":    registeredEmail,
			"password": "anotherpassword",
			"name":     "Another User",
		}

		resp, err := makeRequest("POST", "/v1/auth/register", body, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 409 or 400 for duplicate email, got %d", resp.StatusCode)
		}
	})

	t.Run("LoginUser", func(t *testing.T) {
		body := map[string]string{
			"email":    registeredEmail,
			"password": registeredPassword,
		}

		resp, err := makeRequest("POST", "/v1/auth/login", body, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[AuthResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if result.Token == "" {
			t.Error("Expected token in response")
		}
		authToken = result.Token
	})

	t.Run("LoginInvalidCredentials", func(t *testing.T) {
		body := map[string]string{
			"email":    registeredEmail,
			"password": "wrongpassword",
		}

		resp, err := makeRequest("POST", "/v1/auth/login", body, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 401 or 400 for invalid credentials, got %d", resp.StatusCode)
		}
	})

	t.Run("GetProfile", func(t *testing.T) {
		resp, err := makeRequest("GET", "/v1/users/profile", nil, authToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetProfileUnauthorized", func(t *testing.T) {
		resp, err := makeRequest("GET", "/v1/users/profile", nil, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

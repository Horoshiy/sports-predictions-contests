//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestContestFlow(t *testing.T) {
	// Setup: Register and get auth token
	email := generateTestEmail()
	body := map[string]string{
		"email":    email,
		"password": "testpassword123",
		"name":     "Contest Test User",
	}
	resp, err := makeRequest("POST", "/v1/auth/register", body, "")
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	authResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse auth response: %v", err)
	}
	token := authResp.Token

	var contestID uint

	t.Run("CreateContest", func(t *testing.T) {
		contestBody := map[string]interface{}{
			"title":            generateTestName("Test Contest"),
			"description":      "A test prediction contest",
			"sport_type":       "football",
			"rules":            `{"scoring": {"exact_score": 3, "correct_outcome": 1}}`,
			"start_date":       time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			"end_date":         time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"max_participants": 100,
		}

		resp, err := makeRequest("POST", "/v1/contests", contestBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[ContestResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		contestID = result.Contest.ID
		t.Logf("Created contest ID: %d", contestID)
	})

	t.Run("GetContest", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/contests/%d", contestID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[ContestResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if result.Contest.ID != contestID {
			t.Errorf("Expected contest ID %d, got %d", contestID, result.Contest.ID)
		}
	})

	t.Run("ListContests", func(t *testing.T) {
		resp, err := makeRequest("GET", "/v1/contests", nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("UpdateContest", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"id":          contestID,
			"title":       generateTestName("Updated Contest"),
			"description": "Updated description",
			"sport_type":  "football",
			"rules":       `{"scoring": {"exact_score": 5}}`,
			"start_date":  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			"end_date":    time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		}

		resp, err := makeRequest("PUT", fmt.Sprintf("/v1/contests/%d", contestID), updateBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("JoinContest", func(t *testing.T) {
		// Register another user to join
		email2 := generateTestEmail()
		body := map[string]string{
			"email":    email2,
			"password": "testpassword123",
			"name":     "Joiner User",
		}
		resp, err := makeRequest("POST", "/v1/auth/register", body, "")
		if err != nil {
			t.Fatalf("Failed to register second user: %v", err)
		}
		authResp2, err := parseResponse[AuthResponse](resp)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to parse auth response: %v", err)
		}
		token2 := authResp2.Token

		resp, err = makeRequest("POST", fmt.Sprintf("/v1/contests/%d/join", contestID), nil, token2)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("ListParticipants", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/contests/%d/participants", contestID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetContestNotFound", func(t *testing.T) {
		resp, err := makeRequest("GET", "/v1/contests/999999", nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 404 or 500 for non-existent contest, got %d", resp.StatusCode)
		}
	})
}

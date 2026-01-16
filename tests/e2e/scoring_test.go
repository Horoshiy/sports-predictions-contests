//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestScoringFlow(t *testing.T) {
	// Setup: Register and get auth token
	email := generateTestEmail()
	body := map[string]string{
		"email":    email,
		"password": "testpassword123",
		"name":     "Scoring Test User",
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
	userID := authResp.User.ID

	var contestID, predictionID uint

	// Setup contest and prediction
	t.Run("SetupContestAndPrediction", func(t *testing.T) {
		// Create contest
		contestBody := map[string]interface{}{
			"title":            generateTestName("Scoring Test Contest"),
			"description":      "Contest for scoring tests",
			"sport_type":       "football",
			"rules":            `{"scoring": {"exact_score": 3}}`,
			"start_date":       time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"end_date":         time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"max_participants": 100,
		}
		resp, err := makeRequest("POST", "/v1/contests", contestBody, token)
		if err != nil {
			t.Fatalf("Failed to create contest: %v", err)
		}
		contestResp, err := parseResponse[ContestResponse](resp)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to parse contest response: %v", err)
		}
		contestID = contestResp.Contest.ID

		// Create event
		eventBody := map[string]interface{}{
			"title":      generateTestName("Scoring Test Match"),
			"sport_type": "football",
			"home_team":  "Team A",
			"away_team":  "Team B",
			"event_date": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		}
		resp, err = makeRequest("POST", "/v1/events", eventBody, token)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
		eventResp, err := parseResponse[EventResponse](resp)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to parse event response: %v", err)
		}
		eventID := eventResp.Event.ID

		// Create prediction
		predBody := map[string]interface{}{
			"contest_id":      contestID,
			"event_id":        eventID,
			"prediction_data": `{"outcome": "home_win"}`,
		}
		resp, err = makeRequest("POST", "/v1/predictions", predBody, token)
		if err != nil {
			t.Fatalf("Failed to create prediction: %v", err)
		}
		predResp, err := parseResponse[PredictionResponse](resp)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to parse prediction response: %v", err)
		}
		predictionID = predResp.Prediction.ID

		t.Logf("Setup complete - Contest: %d, Prediction: %d", contestID, predictionID)
	})

	t.Run("CreateScore", func(t *testing.T) {
		scoreBody := map[string]interface{}{
			"user_id":       userID,
			"contest_id":    contestID,
			"prediction_id": predictionID,
			"points":        3.0,
		}

		resp, err := makeRequest("POST", "/v1/scores", scoreBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetLeaderboard", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/contests/%d/leaderboard", contestID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetUserRank", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/contests/%d/users/%d/rank", contestID, userID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetUserStreak", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/contests/%d/users/%d/streak", contestID, userID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetUserAnalytics", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/users/%d/analytics", userID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestPredictionFlow(t *testing.T) {
	// Setup: Register and get auth token
	email := generateTestEmail()
	body := map[string]string{
		"email":    email,
		"password": "testpassword123",
		"name":     "Prediction Test User",
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

	var eventID, predictionID, contestID uint

	// Create a contest first
	t.Run("SetupContest", func(t *testing.T) {
		contestBody := map[string]interface{}{
			"title":            generateTestName("Prediction Test Contest"),
			"description":      "Contest for prediction tests",
			"sport_type":       "football",
			"rules":            `{"scoring": {"exact_score": 3}}`,
			"start_date":       time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			"end_date":         time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"max_participants": 100,
		}

		resp, err := makeRequest("POST", "/v1/contests", contestBody, token)
		if err != nil {
			t.Fatalf("Failed to create contest: %v", err)
		}
		defer resp.Body.Close()

		result, err := parseResponse[ContestResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse contest response: %v", err)
		}
		contestID = result.Contest.ID
		t.Logf("Created contest ID: %d", contestID)
	})

	t.Run("CreateEvent", func(t *testing.T) {
		eventBody := map[string]interface{}{
			"title":      generateTestName("Man Utd vs Liverpool"),
			"sport_type": "football",
			"home_team":  "Manchester United",
			"away_team":  "Liverpool",
			"event_date": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		}

		resp, err := makeRequest("POST", "/v1/events", eventBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[EventResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		eventID = result.Event.ID
		t.Logf("Created event ID: %d", eventID)
	})

	t.Run("ListEvents", func(t *testing.T) {
		resp, err := makeRequest("GET", "/v1/events", nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("SubmitPrediction", func(t *testing.T) {
		predictionBody := map[string]interface{}{
			"contest_id":      contestID,
			"event_id":        eventID,
			"prediction_data": `{"outcome": "home_win", "home_score": 2, "away_score": 1}`,
		}

		resp, err := makeRequest("POST", "/v1/predictions", predictionBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[PredictionResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		predictionID = result.Prediction.ID
		t.Logf("Created prediction ID: %d", predictionID)
	})

	t.Run("GetPrediction", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/predictions/%d", predictionID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetUserPredictions", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/predictions/contest/%d", contestID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("UpdatePrediction", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"id":              predictionID,
			"prediction_data": `{"outcome": "away_win", "home_score": 1, "away_score": 3}`,
		}

		resp, err := makeRequest("PUT", fmt.Sprintf("/v1/predictions/%d", predictionID), updateBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetPotentialCoefficient", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/events/%d/coefficient", eventID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[CoefficientResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if result.Coefficient < 1.0 {
			t.Errorf("Expected coefficient >= 1.0, got %f", result.Coefficient)
		}
		t.Logf("Coefficient: %f, Tier: %s", result.Coefficient, result.Tier)
	})

	t.Run("DeletePrediction", func(t *testing.T) {
		resp, err := makeRequest("DELETE", fmt.Sprintf("/v1/predictions/%d", predictionID), nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

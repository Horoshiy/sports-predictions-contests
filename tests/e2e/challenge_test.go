//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// getTestEventID retrieves a valid event ID for testing
func getTestEventID(t *testing.T, token string) uint {
	// Try to get events from the sports service
	resp, err := makeRequest("GET", "/v1/events?limit=1", nil, token)
	if err != nil {
		t.Logf("Failed to get events, using fallback: %v", err)
		return 1 // Fallback to ID 1 if events endpoint is not available
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Logf("Events endpoint returned status %d, using fallback", resp.StatusCode)
		return 1 // Fallback to ID 1
	}

	// For now, return 1 as fallback since we don't have the events response structure
	// In a real implementation, we would parse the response and get the first event ID
	return 1
}

func TestChallengeFlow(t *testing.T) {
	// Setup: Register two users and get auth tokens
	challengerEmail := generateTestEmail()
	challengerBody := map[string]string{
		"email":    challengerEmail,
		"password": "testpassword123",
		"name":     "Challenge Test Challenger",
	}
	resp, err := makeRequest("POST", "/v1/auth/register", challengerBody, "")
	if err != nil {
		t.Fatalf("Failed to register challenger: %v", err)
	}
	challengerAuthResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse challenger auth response: %v", err)
	}
	challengerToken := challengerAuthResp.Token
	challengerID := challengerAuthResp.User.ID

	opponentEmail := generateTestEmail()
	opponentBody := map[string]string{
		"email":    opponentEmail,
		"password": "testpassword123",
		"name":     "Challenge Test Opponent",
	}
	resp, err = makeRequest("POST", "/v1/auth/register", opponentBody, "")
	if err != nil {
		t.Fatalf("Failed to register opponent: %v", err)
	}
	opponentAuthResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse opponent auth response: %v", err)
	}
	opponentToken := opponentAuthResp.Token
	opponentID := opponentAuthResp.User.ID

	// Get a valid event ID from the system
	eventID := getTestEventID(t, challengerToken)

	var challengeID uint

	t.Run("CreateChallenge", func(t *testing.T) {
		challengeBody := map[string]interface{}{
			"opponent_id": opponentID,
			"event_id":    eventID,
			"message":     "Let's see who's the better predictor!",
		}

		resp, err := makeRequest("POST", "/v1/challenges", challengeBody, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		challengeResp, err := parseResponse[ChallengeResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !challengeResp.Response.Success {
			t.Fatalf("Expected success, got: %s", challengeResp.Response.Message)
		}

		challengeID = challengeResp.Challenge.ID
		if challengeResp.Challenge.ChallengerId != challengerID {
			t.Errorf("Expected challenger ID %d, got %d", challengerID, challengeResp.Challenge.ChallengerId)
		}
		if challengeResp.Challenge.OpponentId != opponentID {
			t.Errorf("Expected opponent ID %d, got %d", opponentID, challengeResp.Challenge.OpponentId)
		}
		if challengeResp.Challenge.Status != "pending" {
			t.Errorf("Expected status 'pending', got %s", challengeResp.Challenge.Status)
		}
	})

	t.Run("GetChallenge", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/challenges/%d", challengeID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		challengeResp, err := parseResponse[ChallengeResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !challengeResp.Response.Success {
			t.Fatalf("Expected success, got: %s", challengeResp.Response.Message)
		}

		if challengeResp.Challenge.ID != challengeID {
			t.Errorf("Expected challenge ID %d, got %d", challengeID, challengeResp.Challenge.ID)
		}
	})

	t.Run("ListUserChallenges", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/users/%d/challenges", challengerID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		challengesResp, err := parseResponse[ChallengesResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !challengesResp.Response.Success {
			t.Fatalf("Expected success, got: %s", challengesResp.Response.Message)
		}

		if len(challengesResp.Challenges) == 0 {
			t.Error("Expected at least one challenge")
		}

		found := false
		for _, challenge := range challengesResp.Challenges {
			if challenge.ID == challengeID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created challenge not found in user's challenges")
		}
	})

	t.Run("AcceptChallenge", func(t *testing.T) {
		resp, err := makeRequest("PUT", fmt.Sprintf("/v1/challenges/%d/accept", challengeID), nil, opponentToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		challengeResp, err := parseResponse[ChallengeResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !challengeResp.Response.Success {
			t.Fatalf("Expected success, got: %s", challengeResp.Response.Message)
		}

		if challengeResp.Challenge.Status != "accepted" {
			t.Errorf("Expected status 'accepted', got %s", challengeResp.Challenge.Status)
		}
	})

	t.Run("VerifyAcceptedChallenge", func(t *testing.T) {
		resp, err := makeRequest("GET", fmt.Sprintf("/v1/challenges/%d", challengeID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		challengeResp, err := parseResponse[ChallengeResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if challengeResp.Challenge.Status != "accepted" {
			t.Errorf("Expected status 'accepted', got %s", challengeResp.Challenge.Status)
		}

		if challengeResp.Challenge.AcceptedAt == "" {
			t.Error("Expected AcceptedAt to be set")
		}
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		// Try to accept challenge as challenger (should fail)
		resp, err := makeRequest("PUT", fmt.Sprintf("/v1/challenges/%d/accept", challengeID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Error("Expected error when challenger tries to accept own challenge")
		}
	})

	t.Run("InvalidChallengeOperations", func(t *testing.T) {
		// Try to accept already accepted challenge
		resp, err := makeRequest("PUT", fmt.Sprintf("/v1/challenges/%d/accept", challengeID), nil, opponentToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			challengeResp, err := parseResponse[ChallengeResponse](resp)
			if err == nil && challengeResp.Response.Success {
				t.Error("Expected error when trying to accept already accepted challenge")
			}
		}

		// Try to withdraw accepted challenge
		resp, err = makeRequest("PUT", fmt.Sprintf("/v1/challenges/%d/withdraw", challengeID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			challengeResp, err := parseResponse[ChallengeResponse](resp)
			if err == nil && challengeResp.Response.Success {
				t.Error("Expected error when trying to withdraw accepted challenge")
			}
		}
	})
}

func TestChallengeDeclineFlow(t *testing.T) {
	// Setup: Register two users
	challengerEmail := generateTestEmail()
	challengerBody := map[string]string{
		"email":    challengerEmail,
		"password": "testpassword123",
		"name":     "Decline Test Challenger",
	}
	resp, err := makeRequest("POST", "/v1/auth/register", challengerBody, "")
	if err != nil {
		t.Fatalf("Failed to register challenger: %v", err)
	}
	challengerAuthResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse challenger auth response: %v", err)
	}
	challengerToken := challengerAuthResp.Token

	opponentEmail := generateTestEmail()
	opponentBody := map[string]string{
		"email":    opponentEmail,
		"password": "testpassword123",
		"name":     "Decline Test Opponent",
	}
	resp, err = makeRequest("POST", "/v1/auth/register", opponentBody, "")
	if err != nil {
		t.Fatalf("Failed to register opponent: %v", err)
	}
	opponentAuthResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse opponent auth response: %v", err)
	}
	opponentToken := opponentAuthResp.Token
	opponentID := opponentAuthResp.User.ID

	var challengeID uint

	// Get a valid event ID from the system
	eventID := getTestEventID(t, challengerToken)

	t.Run("CreateChallengeForDecline", func(t *testing.T) {
		challengeBody := map[string]interface{}{
			"opponent_id": opponentID,
			"event_id":    eventID,
			"message":     "Challenge to be declined",
		}

		resp, err := makeRequest("POST", "/v1/challenges", challengeBody, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		challengeResp, err := parseResponse[ChallengeResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		challengeID = challengeResp.Challenge.ID
	})

	t.Run("DeclineChallenge", func(t *testing.T) {
		resp, err := makeRequest("PUT", fmt.Sprintf("/v1/challenges/%d/decline", challengeID), nil, opponentToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// Verify challenge status
		resp, err = makeRequest("GET", fmt.Sprintf("/v1/challenges/%d", challengeID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to get challenge: %v", err)
		}
		defer resp.Body.Close()

		challengeResp, err := parseResponse[ChallengeResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if challengeResp.Challenge.Status != "declined" {
			t.Errorf("Expected status 'declined', got %s", challengeResp.Challenge.Status)
		}
	})
}

func TestChallengeWithdrawFlow(t *testing.T) {
	// Setup: Register two users
	challengerEmail := generateTestEmail()
	challengerBody := map[string]string{
		"email":    challengerEmail,
		"password": "testpassword123",
		"name":     "Withdraw Test Challenger",
	}
	resp, err := makeRequest("POST", "/v1/auth/register", challengerBody, "")
	if err != nil {
		t.Fatalf("Failed to register challenger: %v", err)
	}
	challengerAuthResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse challenger auth response: %v", err)
	}
	challengerToken := challengerAuthResp.Token

	opponentEmail := generateTestEmail()
	opponentBody := map[string]string{
		"email":    opponentEmail,
		"password": "testpassword123",
		"name":     "Withdraw Test Opponent",
	}
	resp, err = makeRequest("POST", "/v1/auth/register", opponentBody, "")
	if err != nil {
		t.Fatalf("Failed to register opponent: %v", err)
	}
	opponentAuthResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse opponent auth response: %v", err)
	}
	opponentID := opponentAuthResp.User.ID

	var challengeID uint

	// Get a valid event ID from the system
	eventID := getTestEventID(t, challengerToken)

	t.Run("CreateChallengeForWithdraw", func(t *testing.T) {
		challengeBody := map[string]interface{}{
			"opponent_id": opponentID,
			"event_id":    eventID,
			"message":     "Challenge to be withdrawn",
		}

		resp, err := makeRequest("POST", "/v1/challenges", challengeBody, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		challengeResp, err := parseResponse[ChallengeResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		challengeID = challengeResp.Challenge.ID
	})

	t.Run("WithdrawChallenge", func(t *testing.T) {
		resp, err := makeRequest("PUT", fmt.Sprintf("/v1/challenges/%d/withdraw", challengeID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// Verify challenge is deleted
		resp, err = makeRequest("GET", fmt.Sprintf("/v1/challenges/%d", challengeID), nil, challengerToken)
		if err != nil {
			t.Fatalf("Failed to get challenge: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected challenge to be deleted, but got status %d", resp.StatusCode)
		}
	})
}

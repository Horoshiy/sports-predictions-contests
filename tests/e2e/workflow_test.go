//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestCompleteUserJourney validates the entire platform workflow in a single test
func TestCompleteUserJourney(t *testing.T) {
	// Step 1: Register new user
	t.Log("Step 1: Registering new user...")
	email := generateTestEmail()
	registerBody := map[string]string{
		"email":    email,
		"password": "journey123",
		"name":     "Journey Test User",
	}

	resp, err := makeRequest("POST", "/v1/auth/register", registerBody, "")
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Registration failed with status %d", resp.StatusCode)
	}

	authResp, err := parseResponse[AuthResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse auth response: %v", err)
	}

	token := authResp.Token
	userID := authResp.User.ID
	t.Logf("Registered user ID: %d", userID)

	// Step 2: Create sport
	t.Log("Step 2: Creating sport...")
	sportBody := map[string]interface{}{
		"name":        generateTestName("Soccer"),
		"slug":        generateTestName("soccer"),
		"description": "Association football",
	}

	resp, err = makeRequest("POST", "/v1/sports", sportBody, token)
	if err != nil {
		t.Fatalf("Failed to create sport: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create sport failed with status %d", resp.StatusCode)
	}

	sportResp, err := parseResponse[SportResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse sport response: %v", err)
	}
	sportID := sportResp.Sport.ID
	t.Logf("Created sport ID: %d", sportID)

	// Step 3: Create league
	t.Log("Step 3: Creating league...")
	leagueBody := map[string]interface{}{
		"sport_id": sportID,
		"name":     generateTestName("Champions League"),
		"slug":     generateTestName("ucl"),
		"country":  "Europe",
		"season":   "2025-2026",
	}

	resp, err = makeRequest("POST", "/v1/leagues", leagueBody, token)
	if err != nil {
		t.Fatalf("Failed to create league: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create league failed with status %d", resp.StatusCode)
	}

	leagueResp, err := parseResponse[LeagueResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse league response: %v", err)
	}
	leagueID := leagueResp.League.ID
	t.Logf("Created league ID: %d", leagueID)

	// Step 4: Create teams
	t.Log("Step 4: Creating teams...")
	homeTeamBody := map[string]interface{}{
		"sport_id":   sportID,
		"name":       generateTestName("Real Madrid"),
		"slug":       generateTestName("real-madrid"),
		"short_name": "RMA",
		"country":    "Spain",
	}

	resp, err = makeRequest("POST", "/v1/teams", homeTeamBody, token)
	if err != nil {
		t.Fatalf("Failed to create home team: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create home team failed with status %d", resp.StatusCode)
	}
	homeTeamResp, err := parseResponse[TeamResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse home team response: %v", err)
	}
	homeTeamID := homeTeamResp.Team.ID

	awayTeamBody := map[string]interface{}{
		"sport_id":   sportID,
		"name":       generateTestName("Barcelona"),
		"slug":       generateTestName("barcelona"),
		"short_name": "BAR",
		"country":    "Spain",
	}

	resp, err = makeRequest("POST", "/v1/teams", awayTeamBody, token)
	if err != nil {
		t.Fatalf("Failed to create away team: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create away team failed with status %d", resp.StatusCode)
	}
	awayTeamResp, err := parseResponse[TeamResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse away team response: %v", err)
	}
	awayTeamID := awayTeamResp.Team.ID
	t.Logf("Created teams - Home: %d, Away: %d", homeTeamID, awayTeamID)

	// Step 5: Create match
	t.Log("Step 5: Creating match...")
	matchBody := map[string]interface{}{
		"league_id":    leagueID,
		"home_team_id": homeTeamID,
		"away_team_id": awayTeamID,
		"scheduled_at": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		"status":       "scheduled",
	}

	resp, err = makeRequest("POST", "/v1/matches", matchBody, token)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create match failed with status %d", resp.StatusCode)
	}
	matchResp, err := parseResponse[MatchResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse match response: %v", err)
	}
	t.Logf("Created match ID: %d", matchResp.Match.ID)

	// Step 6: Create contest
	t.Log("Step 6: Creating contest...")
	contestBody := map[string]interface{}{
		"title":            generateTestName("El Clasico Predictions"),
		"description":      "Predict the outcome of El Clasico",
		"sport_type":       "football",
		"rules":            `{"scoring": {"exact_score": 5, "correct_outcome": 2}}`,
		"start_date":       time.Now().Format(time.RFC3339),
		"end_date":         time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"max_participants": 1000,
	}

	resp, err = makeRequest("POST", "/v1/contests", contestBody, token)
	if err != nil {
		t.Fatalf("Failed to create contest: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create contest failed with status %d", resp.StatusCode)
	}

	contestResp, err := parseResponse[ContestResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse contest response: %v", err)
	}
	contestID := contestResp.Contest.ID
	t.Logf("Created contest ID: %d", contestID)

	// Step 7: Create event for prediction
	t.Log("Step 7: Creating event...")
	eventBody := map[string]interface{}{
		"title":      generateTestName("Real Madrid vs Barcelona"),
		"sport_type": "football",
		"home_team":  "Real Madrid",
		"away_team":  "Barcelona",
		"event_date": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}

	resp, err = makeRequest("POST", "/v1/events", eventBody, token)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create event failed with status %d", resp.StatusCode)
	}
	eventResp, err := parseResponse[EventResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse event response: %v", err)
	}
	eventID := eventResp.Event.ID
	t.Logf("Created event ID: %d", eventID)

	// Step 8: Submit prediction
	t.Log("Step 8: Submitting prediction...")
	predictionBody := map[string]interface{}{
		"contest_id":      contestID,
		"event_id":        eventID,
		"prediction_data": `{"outcome": "home_win", "home_score": 3, "away_score": 1}`,
	}

	resp, err = makeRequest("POST", "/v1/predictions", predictionBody, token)
	if err != nil {
		t.Fatalf("Failed to submit prediction: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Submit prediction failed with status %d", resp.StatusCode)
	}

	predResp, err := parseResponse[PredictionResponse](resp)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to parse prediction response: %v", err)
	}
	predictionID := predResp.Prediction.ID
	t.Logf("Created prediction ID: %d", predictionID)

	// Step 9: Verify prediction in user's predictions
	t.Log("Step 9: Verifying user predictions...")
	resp, err = makeRequest("GET", fmt.Sprintf("/v1/predictions/contest/%d", contestID), nil, token)
	if err != nil {
		t.Fatalf("Failed to get user predictions: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get predictions failed with status %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Step 10: Check leaderboard
	t.Log("Step 10: Checking leaderboard...")
	resp, err = makeRequest("GET", fmt.Sprintf("/v1/contests/%d/leaderboard", contestID), nil, token)
	if err != nil {
		t.Fatalf("Failed to get leaderboard: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get leaderboard failed with status %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Step 11: Get user analytics
	t.Log("Step 11: Getting user analytics...")
	resp, err = makeRequest("GET", fmt.Sprintf("/v1/users/%d/analytics", userID), nil, token)
	if err != nil {
		t.Fatalf("Failed to get analytics: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get analytics failed with status %d", resp.StatusCode)
	}
	resp.Body.Close()

	t.Log("✅ Complete user journey test passed!")
}

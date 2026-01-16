//go:build e2e

package e2e

import (
	"net/http"
	"testing"
	"time"
)

func TestSportsFlow(t *testing.T) {
	// Setup: Register and get auth token
	email := generateTestEmail()
	body := map[string]string{
		"email":    email,
		"password": "testpassword123",
		"name":     "Sports Test User",
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

	var sportID, leagueID, homeTeamID, awayTeamID uint

	t.Run("CreateSport", func(t *testing.T) {
		sportBody := map[string]interface{}{
			"name":        generateTestName("Football"),
			"slug":        generateTestName("football"),
			"description": "Test football sport",
		}

		resp, err := makeRequest("POST", "/v1/sports", sportBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[SportResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		sportID = result.Sport.ID
		t.Logf("Created sport ID: %d", sportID)
	})

	t.Run("ListSports", func(t *testing.T) {
		resp, err := makeRequest("GET", "/v1/sports", nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateLeague", func(t *testing.T) {
		leagueBody := map[string]interface{}{
			"sport_id": sportID,
			"name":     generateTestName("Premier League"),
			"slug":     generateTestName("premier-league"),
			"country":  "England",
			"season":   "2025-2026",
		}

		resp, err := makeRequest("POST", "/v1/leagues", leagueBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[LeagueResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		leagueID = result.League.ID
		t.Logf("Created league ID: %d", leagueID)
	})

	t.Run("CreateHomeTeam", func(t *testing.T) {
		teamBody := map[string]interface{}{
			"sport_id":   sportID,
			"name":       generateTestName("Manchester United"),
			"slug":       generateTestName("man-utd"),
			"short_name": "MUN",
			"country":    "England",
		}

		resp, err := makeRequest("POST", "/v1/teams", teamBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[TeamResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		homeTeamID = result.Team.ID
		t.Logf("Created home team ID: %d", homeTeamID)
	})

	t.Run("CreateAwayTeam", func(t *testing.T) {
		teamBody := map[string]interface{}{
			"sport_id":   sportID,
			"name":       generateTestName("Liverpool"),
			"slug":       generateTestName("liverpool"),
			"short_name": "LIV",
			"country":    "England",
		}

		resp, err := makeRequest("POST", "/v1/teams", teamBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[TeamResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		awayTeamID = result.Team.ID
		t.Logf("Created away team ID: %d", awayTeamID)
	})

	t.Run("CreateMatch", func(t *testing.T) {
		matchBody := map[string]interface{}{
			"league_id":    leagueID,
			"home_team_id": homeTeamID,
			"away_team_id": awayTeamID,
			"scheduled_at": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
			"status":       "scheduled",
		}

		resp, err := makeRequest("POST", "/v1/matches", matchBody, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		result, err := parseResponse[MatchResponse](resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		t.Logf("Created match ID: %d", result.Match.ID)
	})

	t.Run("ListMatches", func(t *testing.T) {
		resp, err := makeRequest("GET", "/v1/matches", nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

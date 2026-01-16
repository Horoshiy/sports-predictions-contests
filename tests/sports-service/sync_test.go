package sports_service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sports-prediction-contests/sports-service/internal/external"
)

func TestTheSportsDBClient_GetAllSports(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/all_sports.php" {
			t.Errorf("Expected path /all_sports.php, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sports":[{"idSport":"1","strSport":"Soccer","strFormat":"TeamvsTeam","strSportThumb":"https://example.com/soccer.png"}]}`))
	}))
	defer server.Close()

	client := external.NewClient(server.URL)
	sports, err := client.GetAllSports()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sports) != 1 {
		t.Fatalf("Expected 1 sport, got %d", len(sports))
	}

	if sports[0].IDSport != "1" {
		t.Errorf("Expected IDSport '1', got '%s'", sports[0].IDSport)
	}

	if sports[0].StrSport != "Soccer" {
		t.Errorf("Expected StrSport 'Soccer', got '%s'", sports[0].StrSport)
	}
}

func TestTheSportsDBClient_GetAllLeagues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/all_leagues.php" {
			t.Errorf("Expected path /all_leagues.php, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"leagues":[{"idLeague":"4328","strLeague":"English Premier League","strSport":"Soccer","strCountry":"England"}]}`))
	}))
	defer server.Close()

	client := external.NewClient(server.URL)
	leagues, err := client.GetAllLeagues()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(leagues) != 1 {
		t.Fatalf("Expected 1 league, got %d", len(leagues))
	}

	if leagues[0].IDLeague != "4328" {
		t.Errorf("Expected IDLeague '4328', got '%s'", leagues[0].IDLeague)
	}
}

func TestTheSportsDBClient_GetTeamsByLeague(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lookup_all_teams.php" {
			t.Errorf("Expected path /lookup_all_teams.php, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("id") != "4328" {
			t.Errorf("Expected id=4328, got %s", r.URL.Query().Get("id"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"teams":[{"idTeam":"133604","strTeam":"Arsenal","strTeamShort":"ARS","strCountry":"England"}]}`))
	}))
	defer server.Close()

	client := external.NewClient(server.URL)
	teams, err := client.GetTeamsByLeague("4328")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(teams) != 1 {
		t.Fatalf("Expected 1 team, got %d", len(teams))
	}

	if teams[0].StrTeam != "Arsenal" {
		t.Errorf("Expected StrTeam 'Arsenal', got '%s'", teams[0].StrTeam)
	}
}

func TestTheSportsDBClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := external.NewClient(server.URL)
	_, err := client.GetAllSports()

	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestTheSportsDBClient_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client := external.NewClient(server.URL)
	_, err := client.GetAllSports()

	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

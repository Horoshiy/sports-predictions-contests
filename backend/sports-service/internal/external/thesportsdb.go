package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is the TheSportsDB API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new TheSportsDB API client
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}

// API Response types

type APISport struct {
	IDSport       string `json:"idSport"`
	StrSport      string `json:"strSport"`
	StrFormat     string `json:"strFormat"`
	StrSportThumb string `json:"strSportThumb"`
}

type SportsResponse struct {
	Sports []APISport `json:"sports"`
}

type APILeague struct {
	IDLeague        string `json:"idLeague"`
	StrLeague       string `json:"strLeague"`
	StrSport        string `json:"strSport"`
	StrLeagueAlternate string `json:"strLeagueAlternate"`
	IDSoccerXML     string `json:"idSoccerXML"`
	StrCountry      string `json:"strCountry"`
	StrBadge        string `json:"strBadge"`
}

type LeaguesResponse struct {
	Leagues []APILeague `json:"leagues"`
}

type APITeam struct {
	IDTeam        string `json:"idTeam"`
	StrTeam       string `json:"strTeam"`
	StrTeamShort  string `json:"strTeamShort"`
	StrAlternate  string `json:"strAlternate"`
	IDLeague      string `json:"idLeague"`
	StrLeague     string `json:"strLeague"`
	IDSoccerXML   string `json:"idSoccerXML"`
	StrCountry    string `json:"strCountry"`
	StrTeamBadge  string `json:"strTeamBadge"`
	StrSport      string `json:"strSport"`
}

type TeamsResponse struct {
	Teams []APITeam `json:"teams"`
}

type APIEvent struct {
	IDEvent        string `json:"idEvent"`
	StrEvent       string `json:"strEvent"`
	StrSport       string `json:"strSport"`
	IDLeague       string `json:"idLeague"`
	StrLeague      string `json:"strLeague"`
	StrSeason      string `json:"strSeason"`
	IDHomeTeam     string `json:"idHomeTeam"`
	IDAwayTeam     string `json:"idAwayTeam"`
	StrHomeTeam    string `json:"strHomeTeam"`
	StrAwayTeam    string `json:"strAwayTeam"`
	IntHomeScore   string `json:"intHomeScore"`
	IntAwayScore   string `json:"intAwayScore"`
	StrStatus      string `json:"strStatus"`
	DateEvent      string `json:"dateEvent"`
	StrTime        string `json:"strTime"`
	StrTimestamp   string `json:"strTimestamp"`
}

type EventsResponse struct {
	Events []APIEvent `json:"events"`
}

// GetAllSports fetches all sports from the API
func (c *Client) GetAllSports() ([]APISport, error) {
	url := fmt.Sprintf("%s/all_sports.php", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sports: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result SportsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode sports response: %w", err)
	}
	return result.Sports, nil
}

// GetAllLeagues fetches all leagues from the API
func (c *Client) GetAllLeagues() ([]APILeague, error) {
	url := fmt.Sprintf("%s/all_leagues.php", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leagues: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result LeaguesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode leagues response: %w", err)
	}
	return result.Leagues, nil
}

// GetTeamsByLeague fetches teams for a specific league
func (c *Client) GetTeamsByLeague(leagueID string) ([]APITeam, error) {
	reqURL := fmt.Sprintf("%s/lookup_all_teams.php?id=%s", c.baseURL, url.QueryEscape(leagueID))
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result TeamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode teams response: %w", err)
	}
	return result.Teams, nil
}

// GetUpcomingEventsByTeam fetches upcoming events for a team
func (c *Client) GetUpcomingEventsByTeam(teamID string) ([]APIEvent, error) {
	reqURL := fmt.Sprintf("%s/eventsnext.php?id=%s", c.baseURL, url.QueryEscape(teamID))
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result EventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode events response: %w", err)
	}
	return result.Events, nil
}

// GetEventByID fetches a specific event by ID
func (c *Client) GetEventByID(eventID string) (*APIEvent, error) {
	reqURL := fmt.Sprintf("%s/lookupevent.php?id=%s", c.baseURL, url.QueryEscape(eventID))
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result EventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode event response: %w", err)
	}
	if len(result.Events) == 0 {
		return nil, fmt.Errorf("event not found")
	}
	return &result.Events[0], nil
}

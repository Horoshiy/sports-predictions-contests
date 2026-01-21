package e2e

import "time"

// Response represents a common API response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// User represents a user in the system
type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Response Response `json:"response"`
	User     User     `json:"user"`
	Token    string   `json:"token"`
}

// Contest represents a prediction contest
type Contest struct {
	ID                  uint      `json:"id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	SportType           string    `json:"sport_type"`
	Rules               string    `json:"rules"`
	Status              string    `json:"status"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	MaxParticipants     uint      `json:"max_participants"`
	CurrentParticipants uint      `json:"current_participants"`
	CreatorID           uint      `json:"creator_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Challenge represents a head-to-head challenge
type Challenge struct {
	ID              uint      `json:"id"`
	ChallengerId    uint      `json:"challenger_id"`
	OpponentId      uint      `json:"opponent_id"`
	EventId         uint      `json:"event_id"`
	Message         string    `json:"message"`
	Status          string    `json:"status"`
	ExpiresAt       time.Time `json:"expires_at"`
	AcceptedAt      string    `json:"accepted_at,omitempty"`
	CompletedAt     string    `json:"completed_at,omitempty"`
	WinnerId        uint      `json:"winner_id,omitempty"`
	ChallengerScore float64   `json:"challenger_score"`
	OpponentScore   float64   `json:"opponent_score"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ChallengeResponse represents a single challenge response
type ChallengeResponse struct {
	Response  Response  `json:"response"`
	Challenge Challenge `json:"challenge"`
}

// ChallengesResponse represents multiple challenges response
type ChallengesResponse struct {
	Response   Response    `json:"response"`
	Challenges []Challenge `json:"challenges"`
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

// ContestResponse represents a single contest response
type ContestResponse struct {
	Response Response `json:"response"`
	Contest  Contest  `json:"contest"`
}

// ContestsResponse represents a list of contests response
type ContestsResponse struct {
	Response   Response   `json:"response"`
	Contests   []Contest  `json:"contests"`
	Pagination Pagination `json:"pagination"`
}

// Participant represents a contest participant
type Participant struct {
	ID        uint      `json:"id"`
	ContestID uint      `json:"contest_id"`
	UserID    uint      `json:"user_id"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	JoinedAt  time.Time `json:"joined_at"`
}

// ParticipantResponse represents a participant response
type ParticipantResponse struct {
	Response    Response    `json:"response"`
	Participant Participant `json:"participant"`
}

// ParticipantsResponse represents a list of participants
type ParticipantsResponse struct {
	Response     Response      `json:"response"`
	Participants []Participant `json:"participants"`
	Pagination   Pagination    `json:"pagination"`
}

// Pagination represents pagination info
type Pagination struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// Event represents a sports event
type Event struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	SportType  string    `json:"sport_type"`
	HomeTeam   string    `json:"home_team"`
	AwayTeam   string    `json:"away_team"`
	EventDate  time.Time `json:"event_date"`
	Status     string    `json:"status"`
	ResultData string    `json:"result_data"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EventResponse represents a single event response
type EventResponse struct {
	Response Response `json:"response"`
	Event    Event    `json:"event"`
}

// EventsResponse represents a list of events
type EventsResponse struct {
	Response   Response   `json:"response"`
	Events     []Event    `json:"events"`
	Pagination Pagination `json:"pagination"`
}

// Prediction represents a user prediction
type Prediction struct {
	ID             uint      `json:"id"`
	ContestID      uint      `json:"contest_id"`
	UserID         uint      `json:"user_id"`
	EventID        uint      `json:"event_id"`
	PredictionData string    `json:"prediction_data"`
	Status         string    `json:"status"`
	SubmittedAt    time.Time `json:"submitted_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PredictionResponse represents a single prediction response
type PredictionResponse struct {
	Response   Response   `json:"response"`
	Prediction Prediction `json:"prediction"`
}

// PredictionsResponse represents a list of predictions
type PredictionsResponse struct {
	Response    Response     `json:"response"`
	Predictions []Prediction `json:"predictions"`
	Pagination  Pagination   `json:"pagination"`
}

// CoefficientResponse represents time coefficient response
type CoefficientResponse struct {
	Response        Response `json:"response"`
	Coefficient     float64  `json:"coefficient"`
	Tier            string   `json:"tier"`
	HoursUntilEvent float64  `json:"hours_until_event"`
}

// Score represents a user score
type Score struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	ContestID       uint      `json:"contest_id"`
	PredictionID    uint      `json:"prediction_id"`
	Points          float64   `json:"points"`
	TimeCoefficient float64   `json:"time_coefficient"`
	ScoredAt        time.Time `json:"scored_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ScoreResponse represents a single score response
type ScoreResponse struct {
	Response Response `json:"response"`
	Score    Score    `json:"score"`
}

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	UserID        uint    `json:"user_id"`
	UserName      string  `json:"user_name"`
	TotalPoints   float64 `json:"total_points"`
	Rank          uint    `json:"rank"`
	CurrentStreak uint    `json:"current_streak"`
	MaxStreak     uint    `json:"max_streak"`
	Multiplier    float64 `json:"multiplier"`
}

// Leaderboard represents a contest leaderboard
type Leaderboard struct {
	ContestID uint               `json:"contest_id"`
	Entries   []LeaderboardEntry `json:"entries"`
}

// LeaderboardResponse represents a leaderboard response
type LeaderboardResponse struct {
	Response    Response    `json:"response"`
	Leaderboard Leaderboard `json:"leaderboard"`
}

// UserRankResponse represents user rank response
type UserRankResponse struct {
	Response    Response `json:"response"`
	Rank        uint     `json:"rank"`
	TotalPoints float64  `json:"total_points"`
}

// UserStreakResponse represents user streak response
type UserStreakResponse struct {
	Response      Response `json:"response"`
	CurrentStreak uint     `json:"current_streak"`
	MaxStreak     uint     `json:"max_streak"`
	Multiplier    float64  `json:"multiplier"`
}

// UserAnalytics represents user analytics data
type UserAnalytics struct {
	UserID            uint    `json:"user_id"`
	TotalPredictions  uint    `json:"total_predictions"`
	CorrectPredictions uint   `json:"correct_predictions"`
	OverallAccuracy   float64 `json:"overall_accuracy"`
	TotalPoints       float64 `json:"total_points"`
}

// UserAnalyticsResponse represents analytics response
type UserAnalyticsResponse struct {
	Response  Response      `json:"response"`
	Analytics UserAnalytics `json:"analytics"`
}

// Sport represents a sport type
type Sport struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	IconURL     string    `json:"icon_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SportResponse represents a single sport response
type SportResponse struct {
	Response Response `json:"response"`
	Sport    Sport    `json:"sport"`
}

// SportsResponse represents a list of sports
type SportsResponse struct {
	Response   Response   `json:"response"`
	Sports     []Sport    `json:"sports"`
	Pagination Pagination `json:"pagination"`
}

// League represents a sports league
type League struct {
	ID        uint      `json:"id"`
	SportID   uint      `json:"sport_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Country   string    `json:"country"`
	Season    string    `json:"season"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LeagueResponse represents a single league response
type LeagueResponse struct {
	Response Response `json:"response"`
	League   League   `json:"league"`
}

// LeaguesResponse represents a list of leagues
type LeaguesResponse struct {
	Response   Response   `json:"response"`
	Leagues    []League   `json:"leagues"`
	Pagination Pagination `json:"pagination"`
}

// Team represents a sports team
type Team struct {
	ID        uint      `json:"id"`
	SportID   uint      `json:"sport_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ShortName string    `json:"short_name"`
	LogoURL   string    `json:"logo_url"`
	Country   string    `json:"country"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TeamResponse represents a single team response
type TeamResponse struct {
	Response Response `json:"response"`
	Team     Team     `json:"team"`
}

// TeamsResponse represents a list of teams
type TeamsResponse struct {
	Response   Response   `json:"response"`
	Teams      []Team     `json:"teams"`
	Pagination Pagination `json:"pagination"`
}

// Match represents a sports match
type Match struct {
	ID         uint      `json:"id"`
	LeagueID   uint      `json:"league_id"`
	HomeTeamID uint      `json:"home_team_id"`
	AwayTeamID uint      `json:"away_team_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status     string    `json:"status"`
	HomeScore  int       `json:"home_score"`
	AwayScore  int       `json:"away_score"`
	ResultData string    `json:"result_data"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MatchResponse represents a single match response
type MatchResponse struct {
	Response Response `json:"response"`
	Match    Match    `json:"match"`
}

// MatchesResponse represents a list of matches
type MatchesResponse struct {
	Response   Response   `json:"response"`
	Matches    []Match    `json:"matches"`
	Pagination Pagination `json:"pagination"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

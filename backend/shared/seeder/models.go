package seeder

import (
	"time"

	"gorm.io/gorm"
)

// SeederModels contains all the model definitions used by the seeder
// These mirror the actual models from each service but are defined here
// to avoid circular dependencies and provide a clean seeding interface

// User represents a user in the system (from user-service)
type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `gorm:"not null" json:"name"`
}

// Profile represents a user profile (from user-service)
type Profile struct {
	gorm.Model
	UserID            uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Bio               string `gorm:"type:text" json:"bio"`
	AvatarURL         string `gorm:"size:500" json:"avatar_url"`
	Location          string `gorm:"size:100" json:"location"`
	Website           string `gorm:"size:200" json:"website"`
	TwitterURL        string `gorm:"size:200" json:"twitter_url"`
	LinkedInURL       string `gorm:"size:200" json:"linkedin_url"`
	GitHubURL         string `gorm:"size:200" json:"github_url"`
	ProfileVisibility string `gorm:"size:20;default:'public'" json:"profile_visibility"`
}

// UserPreferences represents user preferences (from user-service)
type UserPreferences struct {
	gorm.Model
	UserID                uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Language              string `gorm:"size:10;default:'en'" json:"language"`
	Timezone              string `gorm:"size:50;default:'UTC'" json:"timezone"`
	EmailNotifications    bool   `gorm:"default:true" json:"email_notifications"`
	PushNotifications     bool   `gorm:"default:true" json:"push_notifications"`
	TelegramNotifications bool   `gorm:"default:false" json:"telegram_notifications"`
	Theme                 string `gorm:"size:20;default:'light'" json:"theme"`
}

// Contest represents a sports prediction contest (from contest-service)
type Contest struct {
	gorm.Model
	Title               string    `gorm:"not null" json:"title"`
	Description         string    `json:"description"`
	SportType           string    `gorm:"not null" json:"sport_type"`
	Rules               string    `gorm:"type:text" json:"rules"`
	PredictionSchema    []byte    `gorm:"type:jsonb" json:"prediction_schema"`
	Status              string    `gorm:"not null;default:'draft'" json:"status"`
	StartDate           time.Time `gorm:"not null" json:"start_date"`
	EndDate             time.Time `gorm:"not null" json:"end_date"`
	MaxParticipants     uint      `gorm:"default:0" json:"max_participants"`
	CurrentParticipants uint      `gorm:"default:0" json:"current_participants"`
	CreatorID           uint      `gorm:"not null" json:"creator_id"`
}

// Challenge represents a head-to-head challenge (from challenge-service)
type Challenge struct {
	gorm.Model
	ChallengerID    uint      `gorm:"not null;index" json:"challenger_id"`
	OpponentID      uint      `gorm:"not null;index" json:"opponent_id"`
	EventID         uint      `gorm:"not null;index" json:"event_id"`
	Message         string    `gorm:"type:text" json:"message"`
	Status          string    `gorm:"not null;default:'pending'" json:"status"`
	ExpiresAt       time.Time `gorm:"not null" json:"expires_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	WinnerID        *uint     `json:"winner_id"`
	ChallengerScore float64   `gorm:"default:0" json:"challenger_score"`
	OpponentScore   float64   `gorm:"default:0" json:"opponent_score"`
}

// ChallengeParticipant represents a participant in a challenge (from challenge-service)
type ChallengeParticipant struct {
	gorm.Model
	ChallengeID uint      `gorm:"not null;index" json:"challenge_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Role        string    `gorm:"not null" json:"role"`
	Status      string    `gorm:"not null;default:'active'" json:"status"`
	JoinedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
}

// Sport represents a sport type (from sports-service)
type Sport struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Slug        string `gorm:"not null;uniqueIndex" json:"slug"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	ExternalID  string `gorm:"uniqueIndex;size:50" json:"external_id,omitempty"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// League represents a sports league (from sports-service)
type League struct {
	gorm.Model
	SportID    uint   `gorm:"not null" json:"sport_id"`
	Name       string `gorm:"not null" json:"name"`
	Slug       string `gorm:"not null;uniqueIndex" json:"slug"`
	Country    string `gorm:"size:100" json:"country"`
	Season     string `gorm:"size:50" json:"season"`
	ExternalID string `gorm:"uniqueIndex;size:50" json:"external_id,omitempty"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
}

// Team represents a sports team (from sports-service)
type Team struct {
	gorm.Model
	SportID    uint   `gorm:"not null" json:"sport_id"`
	Name       string `gorm:"not null" json:"name"`
	Slug       string `gorm:"not null;uniqueIndex" json:"slug"`
	ShortName  string `gorm:"size:50" json:"short_name"`
	LogoURL    string `gorm:"size:500" json:"logo_url"`
	Country    string `gorm:"size:100" json:"country"`
	ExternalID string `gorm:"uniqueIndex;size:50" json:"external_id,omitempty"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
}

// Match represents a sports match (from sports-service)
type Match struct {
	gorm.Model
	LeagueID    uint      `gorm:"not null" json:"league_id"`
	HomeTeamID  uint      `gorm:"not null" json:"home_team_id"`
	AwayTeamID  uint      `gorm:"not null" json:"away_team_id"`
	ScheduledAt time.Time `gorm:"not null" json:"scheduled_at"`
	Status      string    `gorm:"not null;default:'scheduled'" json:"status"`
	HomeScore   int       `gorm:"default:0" json:"home_score"`
	AwayScore   int       `gorm:"default:0" json:"away_score"`
	ResultData  string    `gorm:"type:text" json:"result_data"`
	ExternalID  string    `gorm:"uniqueIndex;size:50" json:"external_id,omitempty"`
}

// Event represents a sports event for predictions (from prediction-service)
type Event struct {
	gorm.Model
	Title      string    `gorm:"not null" json:"title"`
	SportType  string    `gorm:"not null;index" json:"sport_type"`
	HomeTeam   string    `gorm:"not null" json:"home_team"`
	AwayTeam   string    `gorm:"not null" json:"away_team"`
	EventDate  time.Time `gorm:"not null;index" json:"event_date"`
	Status     string    `gorm:"not null;default:'scheduled';index" json:"status"`
	ResultData string    `gorm:"type:jsonb" json:"result_data"`
}

// Prediction represents a user prediction (from prediction-service)
type Prediction struct {
	gorm.Model
	UserID         uint      `gorm:"not null" json:"user_id"`
	ContestID      uint      `gorm:"not null" json:"contest_id"`
	EventID        uint      `gorm:"not null" json:"event_id"`
	MatchID        uint      `json:"match_id,omitempty"`
	PredictionType string    `gorm:"not null" json:"prediction_type"`
	PredictionData string    `gorm:"type:jsonb" json:"prediction_data"`
	SubmittedAt    time.Time `gorm:"not null" json:"submitted_at"`
	IsCorrect      *bool     `json:"is_correct,omitempty"`
	Points         float64   `gorm:"default:0" json:"points"`
}

// Score represents a user's score (from scoring-service)
type Score struct {
	gorm.Model
	UserID          uint      `gorm:"not null" json:"user_id"`
	ContestID       uint      `gorm:"not null" json:"contest_id"`
	PredictionID    uint      `gorm:"not null" json:"prediction_id"`
	Points          float64   `gorm:"not null;default:0" json:"points"`
	TimeCoefficient float64   `gorm:"not null;default:1.0" json:"time_coefficient"`
	ScoredAt        time.Time `gorm:"not null" json:"scored_at"`
}

// Leaderboard represents leaderboard entries (from scoring-service)
type Leaderboard struct {
	gorm.Model
	ContestID   uint    `gorm:"not null" json:"contest_id"`
	UserID      uint    `gorm:"not null" json:"user_id"`
	TotalPoints float64 `gorm:"not null;default:0" json:"total_points"`
	Rank        int     `gorm:"not null;default:0" json:"rank"`
	UpdatedAt   time.Time
}

// UserStreak represents user prediction streaks (from scoring-service)
type UserStreak struct {
	gorm.Model
	UserID                uint  `gorm:"not null" json:"user_id"`
	ContestID             uint  `gorm:"not null" json:"contest_id"`
	CurrentStreak         int   `gorm:"not null;default:0" json:"current_streak"`
	MaxStreak             int   `gorm:"not null;default:0" json:"max_streak"`
	LastPredictionID      uint  `json:"last_prediction_id,omitempty"`
	LastPredictionCorrect *bool `json:"last_prediction_correct,omitempty"`
}

// UserTeam represents team tournaments (from contest-service)
type UserTeam struct {
	gorm.Model
	Name           string `gorm:"not null" json:"name"`
	Description    string `gorm:"type:text" json:"description"`
	InviteCode     string `gorm:"uniqueIndex;not null" json:"invite_code"`
	CaptainID      uint   `gorm:"not null" json:"captain_id"`
	MaxMembers     int    `gorm:"default:10" json:"max_members"`
	CurrentMembers int    `gorm:"default:0" json:"current_members"`
	IsActive       bool   `gorm:"default:true" json:"is_active"`
}

// UserTeamMember represents team membership (from contest-service)
type UserTeamMember struct {
	gorm.Model
	TeamID   uint      `gorm:"not null" json:"team_id"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	Role     string    `gorm:"not null;default:'member'" json:"role"`
	Status   string    `gorm:"not null;default:'active'" json:"status"`
	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
}

// Notification represents user notifications (from notification-service)
type Notification struct {
	gorm.Model
	UserID  uint       `gorm:"not null" json:"user_id"`
	Type    string     `gorm:"not null" json:"type"`
	Title   string     `gorm:"not null" json:"title"`
	Message string     `gorm:"not null" json:"message"`
	Data    string     `gorm:"type:text" json:"data"`
	Channel string     `gorm:"not null;default:'in_app'" json:"channel"`
	IsRead  bool       `gorm:"default:false" json:"is_read"`
	SentAt  *time.Time `json:"sent_at,omitempty"`
	ReadAt  *time.Time `json:"read_at,omitempty"`
}

// NotificationPreference represents user notification preferences (from notification-service)
type NotificationPreference struct {
	gorm.Model
	UserID         uint   `gorm:"not null" json:"user_id"`
	Channel        string `gorm:"not null" json:"channel"`
	Enabled        bool   `gorm:"default:true" json:"enabled"`
	TelegramChatID *int64 `json:"telegram_chat_id,omitempty"`
	Email          string `json:"email,omitempty"`
}

// PropType represents prediction prop types (from sports-service)
type PropType struct {
	gorm.Model
	SportType     string   `gorm:"not null" json:"sport_type"`
	Name          string   `gorm:"not null" json:"name"`
	Slug          string   `gorm:"not null" json:"slug"`
	Description   string   `gorm:"type:text" json:"description"`
	Category      string   `gorm:"not null" json:"category"`
	ValueType     string   `gorm:"not null" json:"value_type"`
	DefaultLine   *float64 `json:"default_line,omitempty"`
	MinValue      *float64 `json:"min_value,omitempty"`
	MaxValue      *float64 `json:"max_value,omitempty"`
	PointsCorrect float64  `gorm:"not null;default:2" json:"points_correct"`
	IsActive      bool     `gorm:"default:true" json:"is_active"`
}

// TableName methods to ensure correct table names

func (User) TableName() string                   { return "users" }
func (Profile) TableName() string                { return "profiles" }
func (UserPreferences) TableName() string        { return "user_preferences" }
func (Contest) TableName() string                { return "contests" }
func (Sport) TableName() string                  { return "sports" }
func (League) TableName() string                 { return "leagues" }
func (Team) TableName() string                   { return "teams" }
func (Match) TableName() string                  { return "matches" }
func (Prediction) TableName() string             { return "predictions" }
func (Score) TableName() string                  { return "scores" }
func (Leaderboard) TableName() string            { return "leaderboards" }
func (UserStreak) TableName() string             { return "user_streaks" }
func (UserTeam) TableName() string               { return "user_teams" }
func (UserTeamMember) TableName() string         { return "user_team_members" }
func (Notification) TableName() string           { return "notifications" }
func (NotificationPreference) TableName() string { return "notification_preferences" }
func (PropType) TableName() string               { return "prop_types" }

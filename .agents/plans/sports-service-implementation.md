# Feature: Sports Service Implementation

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

Implement a comprehensive Sports Service microservice that manages sports types, teams, leagues, and matches. This service provides the foundational sports data layer that the Prediction Service and Contest Service depend on for event management.

## User Story

As a platform administrator
I want to manage sports, teams, leagues, and matches
So that users can make predictions on real sporting events

As a contest organizer
I want to associate contests with specific sports and events
So that participants can predict outcomes for relevant matches

## Problem Statement

The platform currently lacks a dedicated service for managing sports data. The Prediction Service has basic Event handling embedded, but there's no centralized management for:
- Sports types (football, basketball, tennis, etc.)
- Teams and their metadata
- Leagues and competitions
- Match scheduling and results

## Solution Statement

Create a dedicated Sports Service microservice following existing patterns that provides:
- CRUD operations for Sports, Teams, Leagues, and Matches
- Hierarchical data model (Sport → League → Team, Match)
- Integration with existing services via gRPC
- REST API exposure through API Gateway

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: Backend microservices, API Gateway, Database
**Dependencies**: PostgreSQL, gRPC, shared auth library

---

## CONTEXT REFERENCES

### Relevant Codebase Files - MUST READ BEFORE IMPLEMENTING

- `backend/contest-service/cmd/main.go` - Service initialization pattern with gRPC server setup
- `backend/contest-service/internal/config/config.go` - Config loading pattern
- `backend/contest-service/internal/models/contest.go` - GORM model with validation hooks
- `backend/contest-service/internal/repository/contest_repository.go` - Repository interface and implementation
- `backend/contest-service/internal/service/contest_service.go` (lines 1-150) - gRPC service implementation
- `backend/contest-service/Dockerfile` - Docker build pattern
- `backend/proto/common.proto` - Common message definitions
- `backend/proto/prediction.proto` - Event message pattern (lines 26-37)
- `backend/api-gateway/internal/config/config.go` - Gateway config pattern
- `backend/api-gateway/internal/gateway/gateway.go` - Service registration pattern
- `docker-compose.yml` - Service configuration pattern

### New Files to Create

```
backend/sports-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── sport.go
│   │   ├── team.go
│   │   ├── league.go
│   │   └── match.go
│   ├── repository/
│   │   ├── sport_repository.go
│   │   ├── team_repository.go
│   │   ├── league_repository.go
│   │   └── match_repository.go
│   └── service/
│       └── sports_service.go
├── go.mod
└── Dockerfile

backend/proto/sports.proto
backend/shared/proto/sports/sports.pb.gw.go (stub)
tests/sports-service/
├── sport_test.go
├── team_test.go
├── league_test.go
└── match_test.go
scripts/init-db.sql (update)
```

### Patterns to Follow

**Naming Conventions:**
- Go files: `snake_case.go`
- Packages: `lowercase`
- Structs/Interfaces: `PascalCase`
- Private functions: `camelCase`
- Proto messages: `PascalCase`
- Proto fields: `snake_case`

**Error Handling Pattern:**
```go
return &pb.Response{
    Response: &common.Response{
        Success:   false,
        Message:   err.Error(),
        Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
        Timestamp: timestamppb.Now(),
    },
}, nil  // Always return nil error, embed in response
```

**Logging Pattern:**
```go
log.Printf("[INFO] Operation successful: %d", id)
log.Printf("[ERROR] Failed operation: %v", err)
```

**Repository Interface Pattern:**
```go
type SportRepositoryInterface interface {
    Create(sport *models.Sport) error
    GetByID(id uint) (*models.Sport, error)
    Update(sport *models.Sport) error
    Delete(id uint) error
    List(limit, offset int, filters ...string) ([]*models.Sport, int64, error)
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation (Proto & Models)

Create gRPC proto definitions and GORM models following existing patterns.

### Phase 2: Data Layer (Repository)

Implement repository interfaces and database operations.

### Phase 3: Business Logic (Service)

Implement gRPC service with authentication and validation.

### Phase 4: Infrastructure (Docker & Gateway)

Configure Docker, update API Gateway, and database schema.

### Phase 5: Testing

Create unit and integration tests.

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `backend/proto/sports.proto`

- **IMPLEMENT**: gRPC service definition with Sport, Team, League, Match entities
- **PATTERN**: Mirror `backend/proto/contest.proto` structure
- **IMPORTS**: common.proto, google/protobuf/timestamp.proto, google/api/annotations.proto
- **VALIDATE**: `cat backend/proto/sports.proto | head -50`

```protobuf
syntax = "proto3";

package sports;

option go_package = "github.com/sports-prediction-contests/shared/proto/sports";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "google/api/annotations.proto";
import "common.proto";

// Sport represents a type of sport
message Sport {
  uint32 id = 1;
  string name = 2;
  string slug = 3;
  string description = 4;
  string icon_url = 5;
  bool is_active = 6;
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp updated_at = 8;
}

// League represents a sports league/competition
message League {
  uint32 id = 1;
  uint32 sport_id = 2;
  string name = 3;
  string slug = 4;
  string country = 5;
  string season = 6;
  bool is_active = 7;
  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;
}

// Team represents a sports team
message Team {
  uint32 id = 1;
  uint32 sport_id = 2;
  string name = 3;
  string slug = 4;
  string short_name = 5;
  string logo_url = 6;
  string country = 7;
  bool is_active = 8;
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Timestamp updated_at = 10;
}

// Match represents a sports match/event
message Match {
  uint32 id = 1;
  uint32 league_id = 2;
  uint32 home_team_id = 3;
  uint32 away_team_id = 4;
  google.protobuf.Timestamp scheduled_at = 5;
  string status = 6;
  int32 home_score = 7;
  int32 away_score = 8;
  string result_data = 9;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
}
```

### Task 2: CREATE `backend/proto/sports.proto` (continued - requests/responses)

- **IMPLEMENT**: Request and response messages for all CRUD operations
- **PATTERN**: Mirror `backend/proto/contest.proto` naming conventions
- **VALIDATE**: `grep -c "Request\|Response" backend/proto/sports.proto`

Add to sports.proto after messages:

```protobuf
// Sport requests
message CreateSportRequest {
  string name = 1;
  string slug = 2;
  string description = 3;
  string icon_url = 4;
}

message GetSportRequest {
  uint32 id = 1;
}

message UpdateSportRequest {
  uint32 id = 1;
  string name = 2;
  string slug = 3;
  string description = 4;
  string icon_url = 5;
  bool is_active = 6;
}

message DeleteSportRequest {
  uint32 id = 1;
}

message ListSportsRequest {
  common.PaginationRequest pagination = 1;
  bool active_only = 2;
}

// League requests
message CreateLeagueRequest {
  uint32 sport_id = 1;
  string name = 2;
  string slug = 3;
  string country = 4;
  string season = 5;
}

message GetLeagueRequest {
  uint32 id = 1;
}

message UpdateLeagueRequest {
  uint32 id = 1;
  uint32 sport_id = 2;
  string name = 3;
  string slug = 4;
  string country = 5;
  string season = 6;
  bool is_active = 7;
}

message DeleteLeagueRequest {
  uint32 id = 1;
}

message ListLeaguesRequest {
  common.PaginationRequest pagination = 1;
  uint32 sport_id = 2;
  bool active_only = 3;
}

// Team requests
message CreateTeamRequest {
  uint32 sport_id = 1;
  string name = 2;
  string slug = 3;
  string short_name = 4;
  string logo_url = 5;
  string country = 6;
}

message GetTeamRequest {
  uint32 id = 1;
}

message UpdateTeamRequest {
  uint32 id = 1;
  uint32 sport_id = 2;
  string name = 3;
  string slug = 4;
  string short_name = 5;
  string logo_url = 6;
  string country = 7;
  bool is_active = 8;
}

message DeleteTeamRequest {
  uint32 id = 1;
}

message ListTeamsRequest {
  common.PaginationRequest pagination = 1;
  uint32 sport_id = 2;
  bool active_only = 3;
}

// Match requests
message CreateMatchRequest {
  uint32 league_id = 1;
  uint32 home_team_id = 2;
  uint32 away_team_id = 3;
  google.protobuf.Timestamp scheduled_at = 4;
}

message GetMatchRequest {
  uint32 id = 1;
}

message UpdateMatchRequest {
  uint32 id = 1;
  uint32 league_id = 2;
  uint32 home_team_id = 3;
  uint32 away_team_id = 4;
  google.protobuf.Timestamp scheduled_at = 5;
  string status = 6;
  int32 home_score = 7;
  int32 away_score = 8;
  string result_data = 9;
}

message DeleteMatchRequest {
  uint32 id = 1;
}

message ListMatchesRequest {
  common.PaginationRequest pagination = 1;
  uint32 league_id = 2;
  uint32 team_id = 3;
  string status = 4;
}

// Responses
message SportResponse {
  common.Response response = 1;
  Sport sport = 2;
}

message ListSportsResponse {
  common.Response response = 1;
  repeated Sport sports = 2;
  common.PaginationResponse pagination = 3;
}

message LeagueResponse {
  common.Response response = 1;
  League league = 2;
}

message ListLeaguesResponse {
  common.Response response = 1;
  repeated League leagues = 2;
  common.PaginationResponse pagination = 3;
}

message TeamResponse {
  common.Response response = 1;
  Team team = 2;
}

message ListTeamsResponse {
  common.Response response = 1;
  repeated Team teams = 2;
  common.PaginationResponse pagination = 3;
}

message MatchResponse {
  common.Response response = 1;
  Match match = 2;
}

message ListMatchesResponse {
  common.Response response = 1;
  repeated Match matches = 2;
  common.PaginationResponse pagination = 3;
}

message DeleteResponse {
  common.Response response = 1;
}
```

### Task 3: CREATE `backend/proto/sports.proto` (continued - service definition)

- **IMPLEMENT**: SportsService with all RPC methods and HTTP annotations
- **PATTERN**: Mirror `backend/proto/prediction.proto` HTTP annotations
- **VALIDATE**: `grep -c "rpc " backend/proto/sports.proto`

Add service definition:

```protobuf
// Sports Service
service SportsService {
  // Sport management
  rpc CreateSport(CreateSportRequest) returns (SportResponse) {
    option (google.api.http) = {
      post: "/v1/sports"
      body: "*"
    };
  }
  rpc GetSport(GetSportRequest) returns (SportResponse) {
    option (google.api.http) = {
      get: "/v1/sports/{id}"
    };
  }
  rpc UpdateSport(UpdateSportRequest) returns (SportResponse) {
    option (google.api.http) = {
      put: "/v1/sports/{id}"
      body: "*"
    };
  }
  rpc DeleteSport(DeleteSportRequest) returns (DeleteResponse) {
    option (google.api.http) = {
      delete: "/v1/sports/{id}"
    };
  }
  rpc ListSports(ListSportsRequest) returns (ListSportsResponse) {
    option (google.api.http) = {
      get: "/v1/sports"
    };
  }

  // League management
  rpc CreateLeague(CreateLeagueRequest) returns (LeagueResponse) {
    option (google.api.http) = {
      post: "/v1/leagues"
      body: "*"
    };
  }
  rpc GetLeague(GetLeagueRequest) returns (LeagueResponse) {
    option (google.api.http) = {
      get: "/v1/leagues/{id}"
    };
  }
  rpc UpdateLeague(UpdateLeagueRequest) returns (LeagueResponse) {
    option (google.api.http) = {
      put: "/v1/leagues/{id}"
      body: "*"
    };
  }
  rpc DeleteLeague(DeleteLeagueRequest) returns (DeleteResponse) {
    option (google.api.http) = {
      delete: "/v1/leagues/{id}"
    };
  }
  rpc ListLeagues(ListLeaguesRequest) returns (ListLeaguesResponse) {
    option (google.api.http) = {
      get: "/v1/leagues"
    };
  }

  // Team management
  rpc CreateTeam(CreateTeamRequest) returns (TeamResponse) {
    option (google.api.http) = {
      post: "/v1/teams"
      body: "*"
    };
  }
  rpc GetTeam(GetTeamRequest) returns (TeamResponse) {
    option (google.api.http) = {
      get: "/v1/teams/{id}"
    };
  }
  rpc UpdateTeam(UpdateTeamRequest) returns (TeamResponse) {
    option (google.api.http) = {
      put: "/v1/teams/{id}"
      body: "*"
    };
  }
  rpc DeleteTeam(DeleteTeamRequest) returns (DeleteResponse) {
    option (google.api.http) = {
      delete: "/v1/teams/{id}"
    };
  }
  rpc ListTeams(ListTeamsRequest) returns (ListTeamsResponse) {
    option (google.api.http) = {
      get: "/v1/teams"
    };
  }

  // Match management
  rpc CreateMatch(CreateMatchRequest) returns (MatchResponse) {
    option (google.api.http) = {
      post: "/v1/matches"
      body: "*"
    };
  }
  rpc GetMatch(GetMatchRequest) returns (MatchResponse) {
    option (google.api.http) = {
      get: "/v1/matches/{id}"
    };
  }
  rpc UpdateMatch(UpdateMatchRequest) returns (MatchResponse) {
    option (google.api.http) = {
      put: "/v1/matches/{id}"
      body: "*"
    };
  }
  rpc DeleteMatch(DeleteMatchRequest) returns (DeleteResponse) {
    option (google.api.http) = {
      delete: "/v1/matches/{id}"
    };
  }
  rpc ListMatches(ListMatchesRequest) returns (ListMatchesResponse) {
    option (google.api.http) = {
      get: "/v1/matches"
    };
  }

  // Health check
  rpc Check(google.protobuf.Empty) returns (common.Response) {
    option (google.api.http) = {
      get: "/v1/sports/health"
    };
  }
}
```

### Task 4: CREATE `backend/sports-service/go.mod`

- **IMPLEMENT**: Go module definition with dependencies
- **PATTERN**: Mirror `backend/contest-service/go.mod`
- **VALIDATE**: `cat backend/sports-service/go.mod`

```go
module github.com/sports-prediction-contests/sports-service

go 1.21

require (
	github.com/sports-prediction-contests/shared v0.0.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

replace github.com/sports-prediction-contests/shared => ../shared
```

### Task 5: CREATE `backend/sports-service/internal/config/config.go`

- **IMPLEMENT**: Configuration loading from environment
- **PATTERN**: Mirror `backend/contest-service/internal/config/config.go`
- **VALIDATE**: `go build ./backend/sports-service/internal/config/`

```go
package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
	LogLevel    string
}

func Load() *Config {
	return &Config{
		Port:        getEnvOrDefault("SPORTS_SERVICE_PORT", "8088"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "your_jwt_secret_key_here"),
		DatabaseURL: getEnvOrDefault("DATABASE_URL", ""),
		LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

func (c *Config) Validate() error {
	if os.Getenv("ENV") == "production" && (c.JWTSecret == "" || c.JWTSecret == "your_jwt_secret_key_here") {
		return errors.New("JWT_SECRET must be set in production")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

### Task 6: CREATE `backend/sports-service/internal/models/sport.go`

- **IMPLEMENT**: Sport GORM model with validation hooks
- **PATTERN**: Mirror `backend/contest-service/internal/models/contest.go`
- **VALIDATE**: `go build ./backend/sports-service/internal/models/`

```go
package models

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Sport struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Slug        string `gorm:"not null;uniqueIndex" json:"slug"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	gorm.Model
}

func (s *Sport) ValidateName() error {
	if len(strings.TrimSpace(s.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(s.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	return nil
}

func (s *Sport) ValidateSlug() error {
	if len(strings.TrimSpace(s.Slug)) == 0 {
		return errors.New("slug cannot be empty")
	}
	if len(s.Slug) > 100 {
		return errors.New("slug cannot exceed 100 characters")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, s.Slug)
	if !matched {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}

func (s *Sport) BeforeCreate(tx *gorm.DB) error {
	if err := s.ValidateName(); err != nil {
		return err
	}
	if s.Slug == "" {
		s.Slug = strings.ToLower(strings.ReplaceAll(s.Name, " ", "-"))
	}
	return s.ValidateSlug()
}

func (s *Sport) BeforeUpdate(tx *gorm.DB) error {
	return s.BeforeCreate(tx)
}
```

### Task 7: CREATE `backend/sports-service/internal/models/league.go`

- **IMPLEMENT**: League GORM model with sport relationship
- **PATTERN**: Mirror Sport model structure
- **VALIDATE**: `go build ./backend/sports-service/internal/models/`

```go
package models

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type League struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	SportID  uint   `gorm:"not null;index" json:"sport_id"`
	Name     string `gorm:"not null" json:"name"`
	Slug     string `gorm:"not null;uniqueIndex" json:"slug"`
	Country  string `json:"country"`
	Season   string `json:"season"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	Sport    Sport  `gorm:"foreignKey:SportID" json:"sport,omitempty"`
	gorm.Model
}

func (l *League) ValidateName() error {
	if len(strings.TrimSpace(l.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(l.Name) > 200 {
		return errors.New("name cannot exceed 200 characters")
	}
	return nil
}

func (l *League) ValidateSlug() error {
	if len(strings.TrimSpace(l.Slug)) == 0 {
		return errors.New("slug cannot be empty")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, l.Slug)
	if !matched {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}

func (l *League) ValidateSportID() error {
	if l.SportID == 0 {
		return errors.New("sport_id is required")
	}
	return nil
}

func (l *League) BeforeCreate(tx *gorm.DB) error {
	if err := l.ValidateName(); err != nil {
		return err
	}
	if err := l.ValidateSportID(); err != nil {
		return err
	}
	if l.Slug == "" {
		l.Slug = strings.ToLower(strings.ReplaceAll(l.Name, " ", "-"))
	}
	return l.ValidateSlug()
}

func (l *League) BeforeUpdate(tx *gorm.DB) error {
	return l.BeforeCreate(tx)
}
```

### Task 8: CREATE `backend/sports-service/internal/models/team.go`

- **IMPLEMENT**: Team GORM model with sport relationship
- **PATTERN**: Mirror League model structure
- **VALIDATE**: `go build ./backend/sports-service/internal/models/`

```go
package models

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Team struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SportID   uint   `gorm:"not null;index" json:"sport_id"`
	Name      string `gorm:"not null" json:"name"`
	Slug      string `gorm:"not null;uniqueIndex" json:"slug"`
	ShortName string `json:"short_name"`
	LogoURL   string `json:"logo_url"`
	Country   string `json:"country"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	Sport     Sport  `gorm:"foreignKey:SportID" json:"sport,omitempty"`
	gorm.Model
}

func (t *Team) ValidateName() error {
	if len(strings.TrimSpace(t.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(t.Name) > 200 {
		return errors.New("name cannot exceed 200 characters")
	}
	return nil
}

func (t *Team) ValidateSlug() error {
	if len(strings.TrimSpace(t.Slug)) == 0 {
		return errors.New("slug cannot be empty")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, t.Slug)
	if !matched {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}

func (t *Team) ValidateSportID() error {
	if t.SportID == 0 {
		return errors.New("sport_id is required")
	}
	return nil
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if err := t.ValidateName(); err != nil {
		return err
	}
	if err := t.ValidateSportID(); err != nil {
		return err
	}
	if t.Slug == "" {
		t.Slug = strings.ToLower(strings.ReplaceAll(t.Name, " ", "-"))
	}
	return t.ValidateSlug()
}

func (t *Team) BeforeUpdate(tx *gorm.DB) error {
	return t.BeforeCreate(tx)
}
```

### Task 9: CREATE `backend/sports-service/internal/models/match.go`

- **IMPLEMENT**: Match GORM model with league and team relationships
- **PATTERN**: Mirror existing model patterns
- **VALIDATE**: `go build ./backend/sports-service/internal/models/`

```go
package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LeagueID    uint      `gorm:"not null;index" json:"league_id"`
	HomeTeamID  uint      `gorm:"not null;index" json:"home_team_id"`
	AwayTeamID  uint      `gorm:"not null;index" json:"away_team_id"`
	ScheduledAt time.Time `gorm:"not null;index" json:"scheduled_at"`
	Status      string    `gorm:"not null;default:'scheduled'" json:"status"`
	HomeScore   int       `gorm:"default:0" json:"home_score"`
	AwayScore   int       `gorm:"default:0" json:"away_score"`
	ResultData  string    `gorm:"type:text" json:"result_data"`
	League      League    `gorm:"foreignKey:LeagueID" json:"league,omitempty"`
	HomeTeam    Team      `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeam    Team      `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	gorm.Model
}

func (m *Match) ValidateLeagueID() error {
	if m.LeagueID == 0 {
		return errors.New("league_id is required")
	}
	return nil
}

func (m *Match) ValidateTeams() error {
	if m.HomeTeamID == 0 {
		return errors.New("home_team_id is required")
	}
	if m.AwayTeamID == 0 {
		return errors.New("away_team_id is required")
	}
	if m.HomeTeamID == m.AwayTeamID {
		return errors.New("home and away teams must be different")
	}
	return nil
}

func (m *Match) ValidateScheduledAt() error {
	if m.ScheduledAt.IsZero() {
		return errors.New("scheduled_at is required")
	}
	return nil
}

func (m *Match) ValidateStatus() error {
	validStatuses := []string{"scheduled", "live", "completed", "cancelled", "postponed"}
	for _, s := range validStatuses {
		if m.Status == s {
			return nil
		}
	}
	return errors.New("invalid status")
}

func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if err := m.ValidateLeagueID(); err != nil {
		return err
	}
	if err := m.ValidateTeams(); err != nil {
		return err
	}
	if err := m.ValidateScheduledAt(); err != nil {
		return err
	}
	if m.Status == "" {
		m.Status = "scheduled"
	}
	return m.ValidateStatus()
}

func (m *Match) BeforeUpdate(tx *gorm.DB) error {
	return m.BeforeCreate(tx)
}

func (m *Match) IsCompleted() bool {
	return m.Status == "completed"
}

func (m *Match) IsLive() bool {
	return m.Status == "live"
}
```

Continue in next tasks...


### Task 10: CREATE `backend/sports-service/internal/repository/sport_repository.go`

- **IMPLEMENT**: Sport repository with CRUD operations
- **PATTERN**: Mirror `backend/contest-service/internal/repository/contest_repository.go`
- **VALIDATE**: `go build ./backend/sports-service/internal/repository/`

```go
package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type SportRepositoryInterface interface {
	Create(sport *models.Sport) error
	GetByID(id uint) (*models.Sport, error)
	GetBySlug(slug string) (*models.Sport, error)
	Update(sport *models.Sport) error
	Delete(id uint) error
	List(limit, offset int, activeOnly bool) ([]*models.Sport, int64, error)
}

type SportRepository struct {
	db *gorm.DB
}

func NewSportRepository(db *gorm.DB) SportRepositoryInterface {
	return &SportRepository{db: db}
}

func (r *SportRepository) Create(sport *models.Sport) error {
	if sport == nil {
		return errors.New("sport cannot be nil")
	}
	return r.db.Create(sport).Error
}

func (r *SportRepository) GetByID(id uint) (*models.Sport, error) {
	if id == 0 {
		return nil, errors.New("invalid sport ID")
	}
	var sport models.Sport
	result := r.db.First(&sport, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport not found")
		}
		return nil, result.Error
	}
	return &sport, nil
}

func (r *SportRepository) GetBySlug(slug string) (*models.Sport, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}
	var sport models.Sport
	result := r.db.Where("slug = ?", slug).First(&sport)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport not found")
		}
		return nil, result.Error
	}
	return &sport, nil
}

func (r *SportRepository) Update(sport *models.Sport) error {
	if sport == nil {
		return errors.New("sport cannot be nil")
	}
	if sport.ID == 0 {
		return errors.New("sport ID cannot be zero")
	}
	result := r.db.Save(sport)
	if result.RowsAffected == 0 {
		return errors.New("sport not found")
	}
	return result.Error
}

func (r *SportRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid sport ID")
	}
	result := r.db.Delete(&models.Sport{}, id)
	if result.RowsAffected == 0 {
		return errors.New("sport not found")
	}
	return result.Error
}

func (r *SportRepository) List(limit, offset int, activeOnly bool) ([]*models.Sport, int64, error) {
	var sports []*models.Sport
	var total int64

	query := r.db.Model(&models.Sport{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&sports).Error; err != nil {
		return nil, 0, err
	}

	return sports, total, nil
}
```

### Task 11: CREATE `backend/sports-service/internal/repository/league_repository.go`

- **IMPLEMENT**: League repository with CRUD and sport filtering
- **PATTERN**: Mirror sport repository structure
- **VALIDATE**: `go build ./backend/sports-service/internal/repository/`

```go
package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type LeagueRepositoryInterface interface {
	Create(league *models.League) error
	GetByID(id uint) (*models.League, error)
	Update(league *models.League) error
	Delete(id uint) error
	List(limit, offset int, sportID uint, activeOnly bool) ([]*models.League, int64, error)
}

type LeagueRepository struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) LeagueRepositoryInterface {
	return &LeagueRepository{db: db}
}

func (r *LeagueRepository) Create(league *models.League) error {
	if league == nil {
		return errors.New("league cannot be nil")
	}
	return r.db.Create(league).Error
}

func (r *LeagueRepository) GetByID(id uint) (*models.League, error) {
	if id == 0 {
		return nil, errors.New("invalid league ID")
	}
	var league models.League
	result := r.db.Preload("Sport").First(&league, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("league not found")
		}
		return nil, result.Error
	}
	return &league, nil
}

func (r *LeagueRepository) Update(league *models.League) error {
	if league == nil {
		return errors.New("league cannot be nil")
	}
	if league.ID == 0 {
		return errors.New("league ID cannot be zero")
	}
	result := r.db.Save(league)
	if result.RowsAffected == 0 {
		return errors.New("league not found")
	}
	return result.Error
}

func (r *LeagueRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid league ID")
	}
	result := r.db.Delete(&models.League{}, id)
	if result.RowsAffected == 0 {
		return errors.New("league not found")
	}
	return result.Error
}

func (r *LeagueRepository) List(limit, offset int, sportID uint, activeOnly bool) ([]*models.League, int64, error) {
	var leagues []*models.League
	var total int64

	query := r.db.Model(&models.League{})
	if sportID > 0 {
		query = query.Where("sport_id = ?", sportID)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Sport").Order("name ASC").Limit(limit).Offset(offset).Find(&leagues).Error; err != nil {
		return nil, 0, err
	}

	return leagues, total, nil
}
```

### Task 12: CREATE `backend/sports-service/internal/repository/team_repository.go`

- **IMPLEMENT**: Team repository with CRUD and sport filtering
- **PATTERN**: Mirror league repository structure
- **VALIDATE**: `go build ./backend/sports-service/internal/repository/`

```go
package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type TeamRepositoryInterface interface {
	Create(team *models.Team) error
	GetByID(id uint) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uint) error
	List(limit, offset int, sportID uint, activeOnly bool) ([]*models.Team, int64, error)
}

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepositoryInterface {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) error {
	if team == nil {
		return errors.New("team cannot be nil")
	}
	return r.db.Create(team).Error
}

func (r *TeamRepository) GetByID(id uint) (*models.Team, error) {
	if id == 0 {
		return nil, errors.New("invalid team ID")
	}
	var team models.Team
	result := r.db.Preload("Sport").First(&team, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, result.Error
	}
	return &team, nil
}

func (r *TeamRepository) Update(team *models.Team) error {
	if team == nil {
		return errors.New("team cannot be nil")
	}
	if team.ID == 0 {
		return errors.New("team ID cannot be zero")
	}
	result := r.db.Save(team)
	if result.RowsAffected == 0 {
		return errors.New("team not found")
	}
	return result.Error
}

func (r *TeamRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid team ID")
	}
	result := r.db.Delete(&models.Team{}, id)
	if result.RowsAffected == 0 {
		return errors.New("team not found")
	}
	return result.Error
}

func (r *TeamRepository) List(limit, offset int, sportID uint, activeOnly bool) ([]*models.Team, int64, error) {
	var teams []*models.Team
	var total int64

	query := r.db.Model(&models.Team{})
	if sportID > 0 {
		query = query.Where("sport_id = ?", sportID)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Sport").Order("name ASC").Limit(limit).Offset(offset).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}
```

### Task 13: CREATE `backend/sports-service/internal/repository/match_repository.go`

- **IMPLEMENT**: Match repository with CRUD and filtering by league/team/status
- **PATTERN**: Mirror team repository structure with additional filters
- **VALIDATE**: `go build ./backend/sports-service/internal/repository/`

```go
package repository

import (
	"errors"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"gorm.io/gorm"
)

type MatchRepositoryInterface interface {
	Create(match *models.Match) error
	GetByID(id uint) (*models.Match, error)
	Update(match *models.Match) error
	Delete(id uint) error
	List(limit, offset int, leagueID, teamID uint, status string) ([]*models.Match, int64, error)
}

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepositoryInterface {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) Create(match *models.Match) error {
	if match == nil {
		return errors.New("match cannot be nil")
	}
	return r.db.Create(match).Error
}

func (r *MatchRepository) GetByID(id uint) (*models.Match, error) {
	if id == 0 {
		return nil, errors.New("invalid match ID")
	}
	var match models.Match
	result := r.db.Preload("League").Preload("HomeTeam").Preload("AwayTeam").First(&match, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("match not found")
		}
		return nil, result.Error
	}
	return &match, nil
}

func (r *MatchRepository) Update(match *models.Match) error {
	if match == nil {
		return errors.New("match cannot be nil")
	}
	if match.ID == 0 {
		return errors.New("match ID cannot be zero")
	}
	result := r.db.Save(match)
	if result.RowsAffected == 0 {
		return errors.New("match not found")
	}
	return result.Error
}

func (r *MatchRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid match ID")
	}
	result := r.db.Delete(&models.Match{}, id)
	if result.RowsAffected == 0 {
		return errors.New("match not found")
	}
	return result.Error
}

func (r *MatchRepository) List(limit, offset int, leagueID, teamID uint, status string) ([]*models.Match, int64, error) {
	var matches []*models.Match
	var total int64

	query := r.db.Model(&models.Match{})
	if leagueID > 0 {
		query = query.Where("league_id = ?", leagueID)
	}
	if teamID > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("League").Preload("HomeTeam").Preload("AwayTeam").
		Order("scheduled_at DESC").Limit(limit).Offset(offset).Find(&matches).Error; err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}
```

### Task 14: CREATE `backend/sports-service/internal/service/sports_service.go`

- **IMPLEMENT**: gRPC service implementation with all CRUD methods
- **PATTERN**: Mirror `backend/contest-service/internal/service/contest_service.go`
- **IMPORTS**: Repository interfaces, proto definitions, auth package
- **VALIDATE**: `go build ./backend/sports-service/internal/service/`

```go
package service

import (
	"context"
	"log"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
	"github.com/sports-prediction-contests/shared/proto/common"
	pb "github.com/sports-prediction-contests/shared/proto/sports"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SportsService struct {
	pb.UnimplementedSportsServiceServer
	sportRepo  repository.SportRepositoryInterface
	leagueRepo repository.LeagueRepositoryInterface
	teamRepo   repository.TeamRepositoryInterface
	matchRepo  repository.MatchRepositoryInterface
}

func NewSportsService(
	sportRepo repository.SportRepositoryInterface,
	leagueRepo repository.LeagueRepositoryInterface,
	teamRepo repository.TeamRepositoryInterface,
	matchRepo repository.MatchRepositoryInterface,
) *SportsService {
	return &SportsService{
		sportRepo:  sportRepo,
		leagueRepo: leagueRepo,
		teamRepo:   teamRepo,
		matchRepo:  matchRepo,
	}
}

// Health check
func (s *SportsService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return &common.Response{
		Success:   true,
		Message:   "Sports Service is healthy",
		Code:      0,
		Timestamp: timestamppb.Now(),
	}, nil
}

// Sport methods
func (s *SportsService) CreateSport(ctx context.Context, req *pb.CreateSportRequest) (*pb.SportResponse, error) {
	sport := &models.Sport{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IconURL:     req.IconUrl,
		IsActive:    true,
	}

	if err := s.sportRepo.Create(sport); err != nil {
		log.Printf("[ERROR] Failed to create sport: %v", err)
		return &pb.SportResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] Sport created: %d", sport.ID)
	return &pb.SportResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Sport created successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Sport: s.sportToProto(sport),
	}, nil
}

func (s *SportsService) GetSport(ctx context.Context, req *pb.GetSportRequest) (*pb.SportResponse, error) {
	sport, err := s.sportRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.SportResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.SportResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Sport retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Sport: s.sportToProto(sport),
	}, nil
}

func (s *SportsService) UpdateSport(ctx context.Context, req *pb.UpdateSportRequest) (*pb.SportResponse, error) {
	sport, err := s.sportRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.SportResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	sport.Name = req.Name
	sport.Slug = req.Slug
	sport.Description = req.Description
	sport.IconURL = req.IconUrl
	sport.IsActive = req.IsActive

	if err := s.sportRepo.Update(sport); err != nil {
		return &pb.SportResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.SportResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Sport updated successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Sport: s.sportToProto(sport),
	}, nil
}

func (s *SportsService) DeleteSport(ctx context.Context, req *pb.DeleteSportRequest) (*pb.DeleteResponse, error) {
	if err := s.sportRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.DeleteResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Sport deleted successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

func (s *SportsService) ListSports(ctx context.Context, req *pb.ListSportsRequest) (*pb.ListSportsResponse, error) {
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := (int(req.Pagination.Page) - 1) * limit
	if offset < 0 {
		offset = 0
	}

	sports, total, err := s.sportRepo.List(limit, offset, req.ActiveOnly)
	if err != nil {
		return &pb.ListSportsResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	pbSports := make([]*pb.Sport, len(sports))
	for i, sport := range sports {
		pbSports[i] = s.sportToProto(sport)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListSportsResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Sports retrieved successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		Sports: pbSports,
		Pagination: &common.PaginationResponse{
			Page:       req.Pagination.Page,
			Limit:      req.Pagination.Limit,
			Total:      int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

// Helper methods
func (s *SportsService) sportToProto(sport *models.Sport) *pb.Sport {
	return &pb.Sport{
		Id:          uint32(sport.ID),
		Name:        sport.Name,
		Slug:        sport.Slug,
		Description: sport.Description,
		IconUrl:     sport.IconURL,
		IsActive:    sport.IsActive,
		CreatedAt:   timestamppb.New(sport.CreatedAt),
		UpdatedAt:   timestamppb.New(sport.UpdatedAt),
	}
}
```


### Task 15: ADD League/Team/Match methods to `backend/sports-service/internal/service/sports_service.go`

- **IMPLEMENT**: Remaining CRUD methods for League, Team, Match entities
- **PATTERN**: Mirror Sport methods structure
- **VALIDATE**: `go build ./backend/sports-service/internal/service/`

Add these methods to the SportsService:

```go
// League methods
func (s *SportsService) CreateLeague(ctx context.Context, req *pb.CreateLeagueRequest) (*pb.LeagueResponse, error) {
	league := &models.League{
		SportID:  uint(req.SportId),
		Name:     req.Name,
		Slug:     req.Slug,
		Country:  req.Country,
		Season:   req.Season,
		IsActive: true,
	}

	if err := s.leagueRepo.Create(league); err != nil {
		log.Printf("[ERROR] Failed to create league: %v", err)
		return &pb.LeagueResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	log.Printf("[INFO] League created: %d", league.ID)
	return &pb.LeagueResponse{
		Response: &common.Response{Success: true, Message: "League created successfully", Code: 0, Timestamp: timestamppb.Now()},
		League:   s.leagueToProto(league),
	}, nil
}

func (s *SportsService) GetLeague(ctx context.Context, req *pb.GetLeagueRequest) (*pb.LeagueResponse, error) {
	league, err := s.leagueRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.LeagueResponse{
		Response: &common.Response{Success: true, Message: "League retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		League:   s.leagueToProto(league),
	}, nil
}

func (s *SportsService) UpdateLeague(ctx context.Context, req *pb.UpdateLeagueRequest) (*pb.LeagueResponse, error) {
	league, err := s.leagueRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	league.SportID = uint(req.SportId)
	league.Name = req.Name
	league.Slug = req.Slug
	league.Country = req.Country
	league.Season = req.Season
	league.IsActive = req.IsActive

	if err := s.leagueRepo.Update(league); err != nil {
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.LeagueResponse{
		Response: &common.Response{Success: true, Message: "League updated successfully", Code: 0, Timestamp: timestamppb.Now()},
		League:   s.leagueToProto(league),
	}, nil
}

func (s *SportsService) DeleteLeague(ctx context.Context, req *pb.DeleteLeagueRequest) (*pb.DeleteResponse, error) {
	if err := s.leagueRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.DeleteResponse{
		Response: &common.Response{Success: true, Message: "League deleted successfully", Code: 0, Timestamp: timestamppb.Now()},
	}, nil
}

func (s *SportsService) ListLeagues(ctx context.Context, req *pb.ListLeaguesRequest) (*pb.ListLeaguesResponse, error) {
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := (int(req.Pagination.Page) - 1) * limit
	if offset < 0 {
		offset = 0
	}

	leagues, total, err := s.leagueRepo.List(limit, offset, uint(req.SportId), req.ActiveOnly)
	if err != nil {
		return &pb.ListLeaguesResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INTERNAL_ERROR), Timestamp: timestamppb.Now()},
		}, nil
	}

	pbLeagues := make([]*pb.League, len(leagues))
	for i, league := range leagues {
		pbLeagues[i] = s.leagueToProto(league)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListLeaguesResponse{
		Response:   &common.Response{Success: true, Message: "Leagues retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Leagues:    pbLeagues,
		Pagination: &common.PaginationResponse{Page: req.Pagination.Page, Limit: req.Pagination.Limit, Total: int32(total), TotalPages: totalPages},
	}, nil
}

// Team methods
func (s *SportsService) CreateTeam(ctx context.Context, req *pb.CreateTeamRequest) (*pb.TeamResponse, error) {
	team := &models.Team{
		SportID:   uint(req.SportId),
		Name:      req.Name,
		Slug:      req.Slug,
		ShortName: req.ShortName,
		LogoURL:   req.LogoUrl,
		Country:   req.Country,
		IsActive:  true,
	}

	if err := s.teamRepo.Create(team); err != nil {
		log.Printf("[ERROR] Failed to create team: %v", err)
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	log.Printf("[INFO] Team created: %d", team.ID)
	return &pb.TeamResponse{
		Response: &common.Response{Success: true, Message: "Team created successfully", Code: 0, Timestamp: timestamppb.Now()},
		Team:     s.teamToProto(team),
	}, nil
}

func (s *SportsService) GetTeam(ctx context.Context, req *pb.GetTeamRequest) (*pb.TeamResponse, error) {
	team, err := s.teamRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.TeamResponse{
		Response: &common.Response{Success: true, Message: "Team retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Team:     s.teamToProto(team),
	}, nil
}

func (s *SportsService) UpdateTeam(ctx context.Context, req *pb.UpdateTeamRequest) (*pb.TeamResponse, error) {
	team, err := s.teamRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	team.SportID = uint(req.SportId)
	team.Name = req.Name
	team.Slug = req.Slug
	team.ShortName = req.ShortName
	team.LogoURL = req.LogoUrl
	team.Country = req.Country
	team.IsActive = req.IsActive

	if err := s.teamRepo.Update(team); err != nil {
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.TeamResponse{
		Response: &common.Response{Success: true, Message: "Team updated successfully", Code: 0, Timestamp: timestamppb.Now()},
		Team:     s.teamToProto(team),
	}, nil
}

func (s *SportsService) DeleteTeam(ctx context.Context, req *pb.DeleteTeamRequest) (*pb.DeleteResponse, error) {
	if err := s.teamRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.DeleteResponse{
		Response: &common.Response{Success: true, Message: "Team deleted successfully", Code: 0, Timestamp: timestamppb.Now()},
	}, nil
}

func (s *SportsService) ListTeams(ctx context.Context, req *pb.ListTeamsRequest) (*pb.ListTeamsResponse, error) {
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := (int(req.Pagination.Page) - 1) * limit
	if offset < 0 {
		offset = 0
	}

	teams, total, err := s.teamRepo.List(limit, offset, uint(req.SportId), req.ActiveOnly)
	if err != nil {
		return &pb.ListTeamsResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INTERNAL_ERROR), Timestamp: timestamppb.Now()},
		}, nil
	}

	pbTeams := make([]*pb.Team, len(teams))
	for i, team := range teams {
		pbTeams[i] = s.teamToProto(team)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListTeamsResponse{
		Response:   &common.Response{Success: true, Message: "Teams retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Teams:      pbTeams,
		Pagination: &common.PaginationResponse{Page: req.Pagination.Page, Limit: req.Pagination.Limit, Total: int32(total), TotalPages: totalPages},
	}, nil
}

// Match methods
func (s *SportsService) CreateMatch(ctx context.Context, req *pb.CreateMatchRequest) (*pb.MatchResponse, error) {
	match := &models.Match{
		LeagueID:    uint(req.LeagueId),
		HomeTeamID:  uint(req.HomeTeamId),
		AwayTeamID:  uint(req.AwayTeamId),
		ScheduledAt: req.ScheduledAt.AsTime(),
		Status:      "scheduled",
	}

	if err := s.matchRepo.Create(match); err != nil {
		log.Printf("[ERROR] Failed to create match: %v", err)
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	log.Printf("[INFO] Match created: %d", match.ID)
	return &pb.MatchResponse{
		Response: &common.Response{Success: true, Message: "Match created successfully", Code: 0, Timestamp: timestamppb.Now()},
		Match:    s.matchToProto(match),
	}, nil
}

func (s *SportsService) GetMatch(ctx context.Context, req *pb.GetMatchRequest) (*pb.MatchResponse, error) {
	match, err := s.matchRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.MatchResponse{
		Response: &common.Response{Success: true, Message: "Match retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Match:    s.matchToProto(match),
	}, nil
}

func (s *SportsService) UpdateMatch(ctx context.Context, req *pb.UpdateMatchRequest) (*pb.MatchResponse, error) {
	match, err := s.matchRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	match.LeagueID = uint(req.LeagueId)
	match.HomeTeamID = uint(req.HomeTeamId)
	match.AwayTeamID = uint(req.AwayTeamId)
	match.ScheduledAt = req.ScheduledAt.AsTime()
	match.Status = req.Status
	match.HomeScore = int(req.HomeScore)
	match.AwayScore = int(req.AwayScore)
	match.ResultData = req.ResultData

	if err := s.matchRepo.Update(match); err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.MatchResponse{
		Response: &common.Response{Success: true, Message: "Match updated successfully", Code: 0, Timestamp: timestamppb.Now()},
		Match:    s.matchToProto(match),
	}, nil
}

func (s *SportsService) DeleteMatch(ctx context.Context, req *pb.DeleteMatchRequest) (*pb.DeleteResponse, error) {
	if err := s.matchRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.DeleteResponse{
		Response: &common.Response{Success: true, Message: "Match deleted successfully", Code: 0, Timestamp: timestamppb.Now()},
	}, nil
}

func (s *SportsService) ListMatches(ctx context.Context, req *pb.ListMatchesRequest) (*pb.ListMatchesResponse, error) {
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := (int(req.Pagination.Page) - 1) * limit
	if offset < 0 {
		offset = 0
	}

	matches, total, err := s.matchRepo.List(limit, offset, uint(req.LeagueId), uint(req.TeamId), req.Status)
	if err != nil {
		return &pb.ListMatchesResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INTERNAL_ERROR), Timestamp: timestamppb.Now()},
		}, nil
	}

	pbMatches := make([]*pb.Match, len(matches))
	for i, match := range matches {
		pbMatches[i] = s.matchToProto(match)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListMatchesResponse{
		Response:   &common.Response{Success: true, Message: "Matches retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Matches:    pbMatches,
		Pagination: &common.PaginationResponse{Page: req.Pagination.Page, Limit: req.Pagination.Limit, Total: int32(total), TotalPages: totalPages},
	}, nil
}

// Helper conversion methods
func (s *SportsService) leagueToProto(league *models.League) *pb.League {
	return &pb.League{
		Id:        uint32(league.ID),
		SportId:   uint32(league.SportID),
		Name:      league.Name,
		Slug:      league.Slug,
		Country:   league.Country,
		Season:    league.Season,
		IsActive:  league.IsActive,
		CreatedAt: timestamppb.New(league.CreatedAt),
		UpdatedAt: timestamppb.New(league.UpdatedAt),
	}
}

func (s *SportsService) teamToProto(team *models.Team) *pb.Team {
	return &pb.Team{
		Id:        uint32(team.ID),
		SportId:   uint32(team.SportID),
		Name:      team.Name,
		Slug:      team.Slug,
		ShortName: team.ShortName,
		LogoUrl:   team.LogoURL,
		Country:   team.Country,
		IsActive:  team.IsActive,
		CreatedAt: timestamppb.New(team.CreatedAt),
		UpdatedAt: timestamppb.New(team.UpdatedAt),
	}
}

func (s *SportsService) matchToProto(match *models.Match) *pb.Match {
	return &pb.Match{
		Id:          uint32(match.ID),
		LeagueId:    uint32(match.LeagueID),
		HomeTeamId:  uint32(match.HomeTeamID),
		AwayTeamId:  uint32(match.AwayTeamID),
		ScheduledAt: timestamppb.New(match.ScheduledAt),
		Status:      match.Status,
		HomeScore:   int32(match.HomeScore),
		AwayScore:   int32(match.AwayScore),
		ResultData:  match.ResultData,
		CreatedAt:   timestamppb.New(match.CreatedAt),
		UpdatedAt:   timestamppb.New(match.UpdatedAt),
	}
}
```


### Task 16: CREATE `backend/sports-service/cmd/main.go`

- **IMPLEMENT**: Service entry point with gRPC server setup
- **PATTERN**: Mirror `backend/contest-service/cmd/main.go`
- **VALIDATE**: `go build ./backend/sports-service/cmd/`

```go
package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/sports-service/internal/config"
	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
	"github.com/sports-prediction-contests/sports-service/internal/service"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	pb "github.com/sports-prediction-contests/shared/proto/sports"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	db, err := database.NewConnectionFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.Sport{}, &models.League{}, &models.Team{}, &models.Match{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	sportRepo := repository.NewSportRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	sportsService := service.NewSportsService(sportRepo, leagueRepo, teamRepo, matchRepo)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
	)

	pb.RegisterSportsServiceServer(server, sportsService)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("[INFO] Sports Service starting on port %s", cfg.Port)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("[INFO] Shutting down Sports Service...")
	server.GracefulStop()

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("[INFO] Sports Service stopped")
}
```

### Task 17: CREATE `backend/sports-service/Dockerfile`

- **IMPLEMENT**: Multi-stage Docker build
- **PATTERN**: Mirror `backend/contest-service/Dockerfile`
- **VALIDATE**: `docker build -t sports-service ./backend/sports-service`

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY ../shared ../shared

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o sports-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/sports-service .

EXPOSE 8088

CMD ["./sports-service"]
```

### Task 18: UPDATE `docker-compose.yml`

- **IMPLEMENT**: Add sports-service configuration
- **PATTERN**: Mirror existing service definitions
- **GOTCHA**: Use port 8088 (next available)
- **VALIDATE**: `docker-compose config`

Add after scoring-service:

```yaml
  # Sports service
  sports-service:
    build:
      context: ./backend/sports-service
      dockerfile: Dockerfile
    container_name: sports-sports-service
    ports:
      - "8088:8088"
    environment:
      - DATABASE_URL=postgres://sports_user:sports_password@postgres:5432/sports_prediction?sslmode=disable
      - JWT_SECRET=your_jwt_secret_key_here
      - SPORTS_SERVICE_PORT=8088
      - LOG_LEVEL=info
    depends_on:
      - postgres
    networks:
      - sports-network
    profiles:
      - services
```

### Task 19: UPDATE `backend/api-gateway/internal/config/config.go`

- **IMPLEMENT**: Add SportsService endpoint configuration
- **PATTERN**: Mirror existing service endpoint fields
- **VALIDATE**: `go build ./backend/api-gateway/internal/config/`

Add to Config struct:
```go
SportsService string
```

Add to Load() function:
```go
SportsService: getEnvOrDefault("SPORTS_SERVICE_ENDPOINT", "sports-service:8088"),
```

### Task 20: UPDATE `backend/api-gateway/internal/gateway/gateway.go`

- **IMPLEMENT**: Register sports service with gateway
- **PATTERN**: Mirror existing service registrations
- **IMPORTS**: Add `sportspb "github.com/sports-prediction-contests/shared/proto/sports"`
- **VALIDATE**: `go build ./backend/api-gateway/...`

Add import:
```go
sportspb "github.com/sports-prediction-contests/shared/proto/sports"
```

Add registration after scoring service:
```go
// Register sports service
err = sportspb.RegisterSportsServiceHandlerFromEndpoint(ctx, mux, cfg.SportsService, opts)
if err != nil {
    return nil, err
}
```

### Task 21: CREATE `backend/shared/proto/sports/sports.pb.gw.go` (stub)

- **IMPLEMENT**: gRPC-gateway stub for compilation
- **PATTERN**: Mirror `backend/shared/proto/contest/contest.pb.gw.go`
- **VALIDATE**: `go build ./backend/shared/proto/sports/`

```go
package sports

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterSportsServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return nil
}
```

### Task 22: UPDATE `scripts/init-db.sql`

- **IMPLEMENT**: Add sports, leagues, teams, matches tables
- **PATTERN**: Mirror existing table definitions
- **VALIDATE**: `psql -f scripts/init-db.sql` (or Docker exec)

Add after leaderboards table:

```sql
-- Create sports table
CREATE TABLE IF NOT EXISTS sports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_sports_slug ON sports(slug);
CREATE INDEX IF NOT EXISTS idx_sports_is_active ON sports(is_active);
CREATE INDEX IF NOT EXISTS idx_sports_deleted_at ON sports(deleted_at);

-- Create leagues table
CREATE TABLE IF NOT EXISTS leagues (
    id SERIAL PRIMARY KEY,
    sport_id INTEGER NOT NULL REFERENCES sports(id),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    country VARCHAR(100),
    season VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_leagues_sport_id ON leagues(sport_id);
CREATE INDEX IF NOT EXISTS idx_leagues_slug ON leagues(slug);
CREATE INDEX IF NOT EXISTS idx_leagues_is_active ON leagues(is_active);
CREATE INDEX IF NOT EXISTS idx_leagues_deleted_at ON leagues(deleted_at);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    sport_id INTEGER NOT NULL REFERENCES sports(id),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    short_name VARCHAR(50),
    logo_url VARCHAR(500),
    country VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_teams_sport_id ON teams(sport_id);
CREATE INDEX IF NOT EXISTS idx_teams_slug ON teams(slug);
CREATE INDEX IF NOT EXISTS idx_teams_is_active ON teams(is_active);
CREATE INDEX IF NOT EXISTS idx_teams_deleted_at ON teams(deleted_at);

-- Create matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    league_id INTEGER NOT NULL REFERENCES leagues(id),
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    scheduled_at TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    home_score INTEGER DEFAULT 0,
    away_score INTEGER DEFAULT 0,
    result_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_matches_league_id ON matches(league_id);
CREATE INDEX IF NOT EXISTS idx_matches_home_team_id ON matches(home_team_id);
CREATE INDEX IF NOT EXISTS idx_matches_away_team_id ON matches(away_team_id);
CREATE INDEX IF NOT EXISTS idx_matches_scheduled_at ON matches(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status);
CREATE INDEX IF NOT EXISTS idx_matches_deleted_at ON matches(deleted_at);
```

### Task 23: UPDATE `docker-compose.yml` API Gateway environment

- **IMPLEMENT**: Add SPORTS_SERVICE_ENDPOINT to api-gateway
- **VALIDATE**: `docker-compose config | grep SPORTS`

Add to api-gateway environment:
```yaml
- SPORTS_SERVICE_ENDPOINT=sports-service:8088
```

### Task 24: CREATE `tests/sports-service/sport_test.go`

- **IMPLEMENT**: Unit tests for Sport model and repository
- **PATTERN**: Mirror `tests/contest-service/contest_test.go`
- **VALIDATE**: `go test ./tests/sports-service/...`

```go
package sports_service_test

import (
	"testing"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Sport{}, &models.League{}, &models.Team{}, &models.Match{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	return db
}

func TestSportValidation(t *testing.T) {
	sport := &models.Sport{Name: "Football", Slug: "football"}
	if err := sport.ValidateName(); err != nil {
		t.Errorf("Expected valid name, got error: %v", err)
	}

	sport.Name = ""
	if err := sport.ValidateName(); err == nil {
		t.Error("Expected error for empty name")
	}
}

func TestSportSlugValidation(t *testing.T) {
	sport := &models.Sport{Name: "Football", Slug: "football"}
	if err := sport.ValidateSlug(); err != nil {
		t.Errorf("Expected valid slug, got error: %v", err)
	}

	sport.Slug = "Invalid Slug!"
	if err := sport.ValidateSlug(); err == nil {
		t.Error("Expected error for invalid slug")
	}
}

func TestSportRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSportRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football", Description: "Association football"}
	if err := repo.Create(sport); err != nil {
		t.Fatalf("Failed to create sport: %v", err)
	}

	if sport.ID == 0 {
		t.Error("Expected sport ID to be set")
	}

	retrieved, err := repo.GetByID(sport.ID)
	if err != nil {
		t.Fatalf("Failed to get sport: %v", err)
	}

	if retrieved.Name != sport.Name {
		t.Errorf("Expected name %s, got %s", sport.Name, retrieved.Name)
	}
}

func TestSportList(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSportRepository(db)

	sports := []*models.Sport{
		{Name: "Football", Slug: "football", IsActive: true},
		{Name: "Basketball", Slug: "basketball", IsActive: true},
		{Name: "Tennis", Slug: "tennis", IsActive: false},
	}

	for _, s := range sports {
		if err := repo.Create(s); err != nil {
			t.Fatalf("Failed to create sport: %v", err)
		}
	}

	list, total, err := repo.List(10, 0, false)
	if err != nil {
		t.Fatalf("Failed to list sports: %v", err)
	}

	if total != 3 {
		t.Errorf("Expected 3 sports, got %d", total)
	}

	activeList, activeTotal, err := repo.List(10, 0, true)
	if err != nil {
		t.Fatalf("Failed to list active sports: %v", err)
	}

	if activeTotal != 2 {
		t.Errorf("Expected 2 active sports, got %d", activeTotal)
	}

	if len(activeList) != 2 {
		t.Errorf("Expected 2 active sports in list, got %d", len(activeList))
	}
}
```

### Task 25: CREATE `tests/sports-service/match_test.go`

- **IMPLEMENT**: Unit tests for Match model and repository
- **PATTERN**: Mirror sport_test.go structure
- **VALIDATE**: `go test ./tests/sports-service/...`

```go
package sports_service_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
)

func TestMatchValidation(t *testing.T) {
	match := &models.Match{
		LeagueID:    1,
		HomeTeamID:  1,
		AwayTeamID:  2,
		ScheduledAt: time.Now().Add(24 * time.Hour),
		Status:      "scheduled",
	}

	if err := match.ValidateTeams(); err != nil {
		t.Errorf("Expected valid teams, got error: %v", err)
	}

	match.HomeTeamID = match.AwayTeamID
	if err := match.ValidateTeams(); err == nil {
		t.Error("Expected error for same home and away team")
	}
}

func TestMatchStatusValidation(t *testing.T) {
	match := &models.Match{Status: "scheduled"}
	if err := match.ValidateStatus(); err != nil {
		t.Errorf("Expected valid status, got error: %v", err)
	}

	match.Status = "invalid"
	if err := match.ValidateStatus(); err == nil {
		t.Error("Expected error for invalid status")
	}
}

func TestMatchRepository(t *testing.T) {
	db := setupTestDB(t)
	sportRepo := repository.NewSportRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football"}
	sportRepo.Create(sport)

	league := &models.League{SportID: sport.ID, Name: "Premier League", Slug: "premier-league"}
	leagueRepo.Create(league)

	homeTeam := &models.Team{SportID: sport.ID, Name: "Team A", Slug: "team-a"}
	awayTeam := &models.Team{SportID: sport.ID, Name: "Team B", Slug: "team-b"}
	teamRepo.Create(homeTeam)
	teamRepo.Create(awayTeam)

	match := &models.Match{
		LeagueID:    league.ID,
		HomeTeamID:  homeTeam.ID,
		AwayTeamID:  awayTeam.ID,
		ScheduledAt: time.Now().Add(24 * time.Hour),
		Status:      "scheduled",
	}

	if err := matchRepo.Create(match); err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}

	retrieved, err := matchRepo.GetByID(match.ID)
	if err != nil {
		t.Fatalf("Failed to get match: %v", err)
	}

	if retrieved.LeagueID != league.ID {
		t.Errorf("Expected league ID %d, got %d", league.ID, retrieved.LeagueID)
	}
}
```

---

## TESTING STRATEGY

### Unit Tests

- Model validation tests for all entities (Sport, League, Team, Match)
- Repository CRUD operation tests with in-memory SQLite
- Service method tests with mocked repositories

### Integration Tests

- Full service workflow tests with real database
- API Gateway routing tests
- gRPC endpoint tests

### Edge Cases

- Empty/null field validation
- Duplicate slug handling
- Foreign key constraint violations
- Invalid status transitions
- Same team as home and away

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
# Build all sports service packages
go build ./backend/sports-service/...

# Build API gateway with new registration
go build ./backend/api-gateway/...

# Verify proto file syntax
cat backend/proto/sports.proto | head -20
```

### Level 2: Unit Tests

```bash
# Run sports service tests
go test ./tests/sports-service/... -v

# Run with coverage
go test ./tests/sports-service/... -cover
```

### Level 3: Docker Build

```bash
# Build sports service image
docker build -t sports-service ./backend/sports-service

# Validate docker-compose config
docker-compose config
```

### Level 4: Integration Tests

```bash
# Start infrastructure
docker-compose up -d postgres redis

# Start sports service
docker-compose --profile services up -d sports-service

# Test health endpoint
curl http://localhost:8088/v1/sports/health
```

### Level 5: Manual API Testing

```bash
# Create sport
curl -X POST http://localhost:8080/v1/sports \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Football", "slug": "football", "description": "Association football"}'

# List sports
curl http://localhost:8080/v1/sports

# Create league
curl -X POST http://localhost:8080/v1/leagues \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"sport_id": 1, "name": "Premier League", "slug": "premier-league", "country": "England", "season": "2025-26"}'
```

---

## ACCEPTANCE CRITERIA

- [ ] Sports Service implements all CRUD operations for Sport, League, Team, Match
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ for models and repositories
- [ ] Service integrates with API Gateway successfully
- [ ] Docker container builds and runs correctly
- [ ] Database schema includes all required tables and indexes
- [ ] gRPC service follows existing authentication patterns
- [ ] Code follows project conventions (naming, error handling, logging)

---

## COMPLETION CHECKLIST

- [ ] All 25 tasks completed in order
- [ ] Each task validation passed
- [ ] All validation commands executed successfully
- [ ] Full test suite passes
- [ ] No linting or type checking errors
- [ ] Manual API testing confirms feature works
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality

---

## NOTES

### Design Decisions

1. **Separate entities vs embedded**: Chose separate Sport, League, Team, Match entities for flexibility and proper normalization
2. **Slug fields**: Added for SEO-friendly URLs and human-readable identifiers
3. **Soft deletes**: Using GORM's DeletedAt for data recovery capability
4. **Status field on Match**: Supports full match lifecycle (scheduled → live → completed)

### Trade-offs

1. **No Redis caching initially**: Can be added later for frequently accessed data
2. **Simple authorization**: All authenticated users can manage sports data; role-based access can be added later
3. **No external sports API integration**: Manual data entry first; external APIs can be integrated later

### Future Enhancements

1. Add Redis caching for sports/leagues/teams lists
2. Implement role-based access control for admin operations
3. Add external sports data API integration
4. Add real-time match status updates via WebSocket

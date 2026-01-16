# Feature: Team Tournaments

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

Team Tournaments enables users to create and join teams that compete together in prediction contests. Teams have shared rankings, member roles (captain, member), and an invitation system with unique codes/links. This strengthens the social aspect of the platform and enables natural viral growth through team invitations.

## User Story

As a sports prediction enthusiast
I want to create or join a team and compete together with friends
So that I can enjoy a social prediction experience and compare our collective performance against other teams

## Problem Statement

Currently, all contest participation is individual. Users cannot collaborate, share strategies, or compete as groups. This limits social engagement and viral growth potential.

## Solution Statement

Implement a team system that allows:
- Creating teams with customizable names and descriptions
- Inviting members via unique invite codes
- Team-based leaderboards aggregating member scores
- Role-based permissions (captain can manage team, members can participate)
- Team activity tracking within contests

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium (4-6 hours)
**Primary Systems Affected**: contest-service, scoring-service, frontend
**Dependencies**: Existing contest, participant, and leaderboard infrastructure

---

## CONTEXT REFERENCES

### Relevant Codebase Files - READ BEFORE IMPLEMENTING

**Backend Proto Patterns:**
- `backend/proto/contest.proto` - Contest/Participant message patterns to mirror
- `backend/proto/scoring.proto` - Leaderboard message patterns for team leaderboards
- `backend/proto/common.proto` - Common response/pagination patterns

**Backend Model Patterns:**
- `backend/contest-service/internal/models/contest.go` - GORM model with validation hooks
- `backend/contest-service/internal/models/participant.go` - Participant model with roles/status

**Backend Repository Patterns:**
- `backend/contest-service/internal/repository/contest_repository.go` - Repository interface pattern

**Backend Service Patterns:**
- `backend/contest-service/internal/service/contest_service.go` - gRPC service implementation

**Frontend Type Patterns:**
- `frontend/src/types/contest.types.ts` - TypeScript interface patterns
- `frontend/src/types/scoring.types.ts` - Leaderboard types

**Frontend Hook Patterns:**
- `frontend/src/hooks/use-contests.ts` - React Query hooks with query keys

**Frontend Component Patterns:**
- `frontend/src/components/contests/ContestList.tsx` - MaterialReactTable pattern
- `frontend/src/components/contests/ContestForm.tsx` - Form with Zod validation

### New Files to Create

**Backend:**
- `backend/proto/team.proto` - Team gRPC definitions
- `backend/contest-service/internal/models/team.go` - Team model
- `backend/contest-service/internal/models/team_member.go` - TeamMember model
- `backend/contest-service/internal/repository/team_repository.go` - Team repository
- `backend/contest-service/internal/service/team_service.go` - Team gRPC service
- `backend/shared/proto/team/team.pb.go` - Generated proto (stub)
- `backend/shared/proto/team/team.pb.gw.go` - Generated gateway (stub)

**Frontend:**
- `frontend/src/types/team.types.ts` - Team TypeScript types
- `frontend/src/utils/team-validation.ts` - Zod schemas
- `frontend/src/services/team-service.ts` - API service
- `frontend/src/hooks/use-teams.ts` - React Query hooks
- `frontend/src/components/teams/TeamList.tsx` - Team list component
- `frontend/src/components/teams/TeamForm.tsx` - Create/edit team form
- `frontend/src/components/teams/TeamMembers.tsx` - Member management
- `frontend/src/components/teams/TeamInvite.tsx` - Invite code/link component
- `frontend/src/components/teams/TeamLeaderboard.tsx` - Team rankings
- `frontend/src/pages/TeamsPage.tsx` - Main teams page

**Database:**
- Update `scripts/init-db.sql` - Add teams and team_members tables

**Tests:**
- `tests/contest-service/team_test.go` - Team model tests

### Patterns to Follow

**Naming Conventions:**
- Go: PascalCase for exported, camelCase for private
- TypeScript: camelCase for variables, PascalCase for types/interfaces
- Proto: snake_case for fields, PascalCase for messages

**Error Handling (Go):**
```go
if err != nil {
    log.Printf("[ERROR] Failed to X: %v", err)
    return &pb.XResponse{
        Response: &common.Response{
            Success:   false,
            Message:   err.Error(),
            Code:      int32(common.ErrorCode_INTERNAL_ERROR),
            Timestamp: timestamppb.Now(),
        },
    }, nil
}
```

**React Query Pattern:**
```typescript
export const teamKeys = {
  all: ['teams'] as const,
  lists: () => [...teamKeys.all, 'list'] as const,
  list: (filters: ListTeamsRequest) => [...teamKeys.lists(), filters] as const,
  details: () => [...teamKeys.all, 'detail'] as const,
  detail: (id: number) => [...teamKeys.details(), id] as const,
}
```

**GORM Model Pattern:**
```go
type Team struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    // fields...
    gorm.Model
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
    // validation
    return nil
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Database & Proto Foundation
- Add database tables for teams and team_members
- Create team.proto with all messages and service definition
- Create proto stub files for compilation

### Phase 2: Backend Models & Repository
- Create Team and TeamMember GORM models with validation
- Create TeamRepository with CRUD operations
- Generate invite codes with crypto/rand

### Phase 3: Backend Service
- Implement TeamService gRPC handlers
- Add team leaderboard aggregation to scoring-service
- Register team service in API gateway

### Phase 4: Frontend Types & Services
- Create TypeScript types matching proto
- Create Zod validation schemas
- Create team-service API client
- Create React Query hooks

### Phase 5: Frontend Components
- Create TeamList with MaterialReactTable
- Create TeamForm for create/edit
- Create TeamMembers for member management
- Create TeamInvite for invite codes
- Create TeamLeaderboard for rankings
- Create TeamsPage with tabs

### Phase 6: Integration & Testing
- Add /teams route to App.tsx
- Create unit tests for models
- Manual validation

---

## STEP-BY-STEP TASKS

### Task 1: UPDATE scripts/init-db.sql

- **IMPLEMENT**: Add teams and team_members tables after notifications section
- **PATTERN**: Follow existing table patterns (id, created_at, updated_at, deleted_at, indexes)
- **GOTCHA**: Use ON DELETE RESTRICT for foreign keys to prevent orphaned records

```sql
-- Add after notification_preferences table:

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    invite_code VARCHAR(20) UNIQUE NOT NULL,
    captain_id INTEGER NOT NULL,
    max_members INTEGER DEFAULT 10,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_teams_invite_code ON teams(invite_code);
CREATE INDEX IF NOT EXISTS idx_teams_captain_id ON teams(captain_id);
CREATE INDEX IF NOT EXISTS idx_teams_is_active ON teams(is_active);
CREATE INDEX IF NOT EXISTS idx_teams_deleted_at ON teams(deleted_at);

-- Create team_members table
CREATE TABLE IF NOT EXISTS team_members (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    user_id INTEGER NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(team_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_team_members_team_id ON team_members(team_id);
CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members(user_id);
CREATE INDEX IF NOT EXISTS idx_team_members_deleted_at ON team_members(deleted_at);

-- Create team_contest_entries table (links teams to contests)
CREATE TABLE IF NOT EXISTS team_contest_entries (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    contest_id INTEGER NOT NULL,
    total_points DECIMAL(10,2) NOT NULL DEFAULT 0,
    rank INTEGER NOT NULL DEFAULT 0,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(team_id, contest_id)
);

CREATE INDEX IF NOT EXISTS idx_team_contest_entries_team_id ON team_contest_entries(team_id);
CREATE INDEX IF NOT EXISTS idx_team_contest_entries_contest_id ON team_contest_entries(contest_id);
CREATE INDEX IF NOT EXISTS idx_team_contest_entries_rank ON team_contest_entries(contest_id, rank);
CREATE INDEX IF NOT EXISTS idx_team_contest_entries_deleted_at ON team_contest_entries(deleted_at);
```

- **VALIDATE**: `cat scripts/init-db.sql | grep -A 5 "teams table"`

---

### Task 2: CREATE backend/proto/team.proto

- **IMPLEMENT**: Define Team, TeamMember messages and TeamService
- **PATTERN**: Mirror contest.proto structure exactly
- **IMPORTS**: common.proto, google/protobuf/timestamp.proto, google/api/annotations.proto

```protobuf
syntax = "proto3";

package team;

option go_package = "github.com/sports-prediction-contests/shared/proto/team";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "google/api/annotations.proto";
import "common.proto";

message Team {
  uint32 id = 1;
  string name = 2;
  string description = 3;
  string invite_code = 4;
  uint32 captain_id = 5;
  uint32 max_members = 6;
  uint32 current_members = 7;
  bool is_active = 8;
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Timestamp updated_at = 10;
}

message TeamMember {
  uint32 id = 1;
  uint32 team_id = 2;
  uint32 user_id = 3;
  string user_name = 4;
  string role = 5;
  string status = 6;
  google.protobuf.Timestamp joined_at = 7;
}

message TeamLeaderboardEntry {
  uint32 team_id = 1;
  string team_name = 2;
  double total_points = 3;
  uint32 rank = 4;
  uint32 member_count = 5;
}

// Request messages
message CreateTeamRequest {
  string name = 1;
  string description = 2;
  uint32 max_members = 3;
}

message UpdateTeamRequest {
  uint32 id = 1;
  string name = 2;
  string description = 3;
  uint32 max_members = 4;
}

message GetTeamRequest {
  uint32 id = 1;
}

message DeleteTeamRequest {
  uint32 id = 1;
}

message ListTeamsRequest {
  common.PaginationRequest pagination = 1;
  bool my_teams_only = 2;
}

message JoinTeamRequest {
  string invite_code = 1;
}

message LeaveTeamRequest {
  uint32 team_id = 1;
}

message RemoveMemberRequest {
  uint32 team_id = 1;
  uint32 user_id = 2;
}

message ListMembersRequest {
  uint32 team_id = 1;
  common.PaginationRequest pagination = 2;
}

message RegenerateInviteCodeRequest {
  uint32 team_id = 1;
}

message JoinContestAsTeamRequest {
  uint32 team_id = 1;
  uint32 contest_id = 2;
}

message LeaveContestAsTeamRequest {
  uint32 team_id = 1;
  uint32 contest_id = 2;
}

message GetTeamLeaderboardRequest {
  uint32 contest_id = 1;
  uint32 limit = 2;
}

// Response messages
message CreateTeamResponse {
  common.Response response = 1;
  Team team = 2;
}

message UpdateTeamResponse {
  common.Response response = 1;
  Team team = 2;
}

message GetTeamResponse {
  common.Response response = 1;
  Team team = 2;
}

message DeleteTeamResponse {
  common.Response response = 1;
}

message ListTeamsResponse {
  common.Response response = 1;
  repeated Team teams = 2;
  common.PaginationResponse pagination = 3;
}

message JoinTeamResponse {
  common.Response response = 1;
  TeamMember member = 2;
}

message LeaveTeamResponse {
  common.Response response = 1;
}

message RemoveMemberResponse {
  common.Response response = 1;
}

message ListMembersResponse {
  common.Response response = 1;
  repeated TeamMember members = 2;
  common.PaginationResponse pagination = 3;
}

message RegenerateInviteCodeResponse {
  common.Response response = 1;
  string invite_code = 2;
}

message JoinContestAsTeamResponse {
  common.Response response = 1;
}

message LeaveContestAsTeamResponse {
  common.Response response = 1;
}

message GetTeamLeaderboardResponse {
  common.Response response = 1;
  repeated TeamLeaderboardEntry entries = 2;
}

// Team Service
service TeamService {
  rpc CreateTeam(CreateTeamRequest) returns (CreateTeamResponse) {
    option (google.api.http) = {
      post: "/v1/teams"
      body: "*"
    };
  }
  rpc UpdateTeam(UpdateTeamRequest) returns (UpdateTeamResponse) {
    option (google.api.http) = {
      put: "/v1/teams/{id}"
      body: "*"
    };
  }
  rpc GetTeam(GetTeamRequest) returns (GetTeamResponse) {
    option (google.api.http) = {
      get: "/v1/teams/{id}"
    };
  }
  rpc DeleteTeam(DeleteTeamRequest) returns (DeleteTeamResponse) {
    option (google.api.http) = {
      delete: "/v1/teams/{id}"
    };
  }
  rpc ListTeams(ListTeamsRequest) returns (ListTeamsResponse) {
    option (google.api.http) = {
      get: "/v1/teams"
    };
  }
  
  // Member management
  rpc JoinTeam(JoinTeamRequest) returns (JoinTeamResponse) {
    option (google.api.http) = {
      post: "/v1/teams/join"
      body: "*"
    };
  }
  rpc LeaveTeam(LeaveTeamRequest) returns (LeaveTeamResponse) {
    option (google.api.http) = {
      post: "/v1/teams/{team_id}/leave"
      body: "*"
    };
  }
  rpc RemoveMember(RemoveMemberRequest) returns (RemoveMemberResponse) {
    option (google.api.http) = {
      delete: "/v1/teams/{team_id}/members/{user_id}"
    };
  }
  rpc ListMembers(ListMembersRequest) returns (ListMembersResponse) {
    option (google.api.http) = {
      get: "/v1/teams/{team_id}/members"
    };
  }
  rpc RegenerateInviteCode(RegenerateInviteCodeRequest) returns (RegenerateInviteCodeResponse) {
    option (google.api.http) = {
      post: "/v1/teams/{team_id}/regenerate-invite"
      body: "*"
    };
  }
  
  // Contest participation
  rpc JoinContestAsTeam(JoinContestAsTeamRequest) returns (JoinContestAsTeamResponse) {
    option (google.api.http) = {
      post: "/v1/teams/{team_id}/contests/{contest_id}/join"
      body: "*"
    };
  }
  rpc LeaveContestAsTeam(LeaveContestAsTeamRequest) returns (LeaveContestAsTeamResponse) {
    option (google.api.http) = {
      post: "/v1/teams/{team_id}/contests/{contest_id}/leave"
      body: "*"
    };
  }
  rpc GetTeamLeaderboard(GetTeamLeaderboardRequest) returns (GetTeamLeaderboardResponse) {
    option (google.api.http) = {
      get: "/v1/contests/{contest_id}/team-leaderboard"
    };
  }
  
  // Health check
  rpc Check(google.protobuf.Empty) returns (common.Response) {
    option (google.api.http) = {
      get: "/v1/teams/health"
    };
  }
}
```

- **VALIDATE**: `cat backend/proto/team.proto | head -50`

---

### Task 3: CREATE backend/shared/proto/team/team.pb.go (stub)

- **IMPLEMENT**: Create minimal stub for compilation
- **PATTERN**: Mirror backend/shared/proto/contest/contest.pb.go structure

```go
package team

// Stub file - will be replaced by protoc generation
// This allows the project to compile before proto generation

type Team struct{}
type TeamMember struct{}
type TeamLeaderboardEntry struct{}
```

- **VALIDATE**: `cat backend/shared/proto/team/team.pb.go`

---

### Task 4: CREATE backend/shared/proto/team/team.pb.gw.go (stub)

- **IMPLEMENT**: Create gateway stub for API gateway registration
- **PATTERN**: Mirror other gateway stubs

```go
package team

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// RegisterTeamServiceHandlerFromEndpoint stub
func RegisterTeamServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return nil
}

// RegisterTeamServiceHandler stub
func RegisterTeamServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return nil
}
```

- **VALIDATE**: `cat backend/shared/proto/team/team.pb.gw.go`

---

### Task 5: CREATE backend/contest-service/internal/models/team.go

- **IMPLEMENT**: Team GORM model with validation hooks
- **PATTERN**: Mirror contest.go model structure exactly
- **IMPORTS**: errors, strings, time, crypto/rand, encoding/hex, gorm.io/gorm

```go
package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Team represents a group of users competing together
type Team struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	Description    string    `json:"description"`
	InviteCode     string    `gorm:"uniqueIndex;not null" json:"invite_code"`
	CaptainID      uint      `gorm:"not null" json:"captain_id"`
	MaxMembers     uint      `gorm:"default:10" json:"max_members"`
	CurrentMembers uint      `gorm:"default:0" json:"current_members"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	gorm.Model
}

// GenerateInviteCode creates a random 8-character invite code
func GenerateInviteCode() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

// ValidateName checks if the team name is valid
func (t *Team) ValidateName() error {
	name := strings.TrimSpace(t.Name)
	if len(name) == 0 {
		return errors.New("team name cannot be empty")
	}
	if len(name) > 100 {
		return errors.New("team name cannot exceed 100 characters")
	}
	return nil
}

// ValidateDescription checks if the description is valid
func (t *Team) ValidateDescription() error {
	if len(t.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}
	return nil
}

// ValidateMaxMembers checks if max members is valid
func (t *Team) ValidateMaxMembers() error {
	if t.MaxMembers < 2 {
		return errors.New("team must allow at least 2 members")
	}
	if t.MaxMembers > 50 {
		return errors.New("team cannot exceed 50 members")
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a team
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if err := t.ValidateName(); err != nil {
		return err
	}
	if err := t.ValidateDescription(); err != nil {
		return err
	}
	if err := t.ValidateMaxMembers(); err != nil {
		return err
	}

	if t.CaptainID == 0 {
		return errors.New("captain ID cannot be empty")
	}

	// Generate invite code if not set
	if t.InviteCode == "" {
		code, err := GenerateInviteCode()
		if err != nil {
			return err
		}
		t.InviteCode = code
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a team
func (t *Team) BeforeUpdate(tx *gorm.DB) error {
	if err := t.ValidateName(); err != nil {
		return err
	}
	if err := t.ValidateDescription(); err != nil {
		return err
	}
	if err := t.ValidateMaxMembers(); err != nil {
		return err
	}
	return nil
}

// CanJoin checks if a new member can join the team
func (t *Team) CanJoin() bool {
	return t.IsActive && t.CurrentMembers < t.MaxMembers
}

// IsCaptain checks if a user is the team captain
func (t *Team) IsCaptain(userID uint) bool {
	return t.CaptainID == userID
}
```

- **VALIDATE**: `cat backend/contest-service/internal/models/team.go | head -30`

---

### Task 6: CREATE backend/contest-service/internal/models/team_member.go

- **IMPLEMENT**: TeamMember GORM model
- **PATTERN**: Mirror participant.go model structure

```go
package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// TeamMember represents a user's membership in a team
type TeamMember struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	TeamID   uint      `gorm:"not null;index" json:"team_id"`
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	Role     string    `gorm:"not null;default:'member'" json:"role"` // "captain", "member"
	Status   string    `gorm:"not null;default:'active'" json:"status"` // "active", "inactive"
	JoinedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	gorm.Model

	// Relationships
	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

// ValidateRole checks if the role is valid
func (m *TeamMember) ValidateRole() error {
	validRoles := []string{"captain", "member"}
	for _, r := range validRoles {
		if m.Role == r {
			return nil
		}
	}
	return errors.New("invalid role: must be 'captain' or 'member'")
}

// ValidateStatus checks if the status is valid
func (m *TeamMember) ValidateStatus() error {
	validStatuses := []string{"active", "inactive"}
	for _, s := range validStatuses {
		if m.Status == s {
			return nil
		}
	}
	return errors.New("invalid status: must be 'active' or 'inactive'")
}

// BeforeCreate is a GORM hook that runs before creating a team member
func (m *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if m.Role == "" {
		m.Role = "member"
	}
	if m.Status == "" {
		m.Status = "active"
	}
	if m.JoinedAt.IsZero() {
		m.JoinedAt = time.Now()
	}

	if err := m.ValidateRole(); err != nil {
		return err
	}
	if err := m.ValidateStatus(); err != nil {
		return err
	}

	if m.TeamID == 0 {
		return errors.New("team ID cannot be empty")
	}
	if m.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}

	// Check for duplicate membership
	var existing TeamMember
	result := tx.Where("team_id = ? AND user_id = ?", m.TeamID, m.UserID).First(&existing)
	if result.Error == nil {
		return errors.New("user is already a member of this team")
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating
func (m *TeamMember) BeforeUpdate(tx *gorm.DB) error {
	if err := m.ValidateRole(); err != nil {
		return err
	}
	if err := m.ValidateStatus(); err != nil {
		return err
	}
	return nil
}

// IsActive checks if the member is active
func (m *TeamMember) IsActive() bool {
	return m.Status == "active"
}

// IsCaptain checks if the member is a captain
func (m *TeamMember) IsCaptain() bool {
	return m.Role == "captain"
}
```

- **VALIDATE**: `cat backend/contest-service/internal/models/team_member.go | head -30`

---

### Task 7: CREATE backend/contest-service/internal/models/team_contest_entry.go

- **IMPLEMENT**: TeamContestEntry model for team participation in contests
- **PATTERN**: Mirror leaderboard.go structure

```go
package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// TeamContestEntry represents a team's participation in a contest
type TeamContestEntry struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TeamID      uint      `gorm:"not null;index" json:"team_id"`
	ContestID   uint      `gorm:"not null;index" json:"contest_id"`
	TotalPoints float64   `gorm:"default:0" json:"total_points"`
	Rank        uint      `gorm:"default:0" json:"rank"`
	JoinedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	gorm.Model

	// Relationships
	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

// BeforeCreate validates before creating
func (e *TeamContestEntry) BeforeCreate(tx *gorm.DB) error {
	if e.TeamID == 0 {
		return errors.New("team ID cannot be empty")
	}
	if e.ContestID == 0 {
		return errors.New("contest ID cannot be empty")
	}
	if e.JoinedAt.IsZero() {
		e.JoinedAt = time.Now()
	}

	// Check for duplicate entry
	var existing TeamContestEntry
	result := tx.Where("team_id = ? AND contest_id = ?", e.TeamID, e.ContestID).First(&existing)
	if result.Error == nil {
		return errors.New("team is already participating in this contest")
	}

	return nil
}
```

- **VALIDATE**: `cat backend/contest-service/internal/models/team_contest_entry.go`


---

### Task 8: CREATE backend/contest-service/internal/repository/team_repository.go

- **IMPLEMENT**: Team and TeamMember repository interfaces and implementations
- **PATTERN**: Mirror contest_repository.go exactly
- **IMPORTS**: errors, gorm.io/gorm, models package

```go
package repository

import (
	"errors"

	"github.com/sports-prediction-contests/contest-service/internal/models"
	"gorm.io/gorm"
)

// TeamRepositoryInterface defines the contract for team repository
type TeamRepositoryInterface interface {
	Create(team *models.Team) error
	GetByID(id uint) (*models.Team, error)
	GetByInviteCode(code string) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uint) error
	List(limit, offset int, userID uint, myTeamsOnly bool) ([]*models.Team, int64, error)
}

// TeamMemberRepositoryInterface defines the contract for team member repository
type TeamMemberRepositoryInterface interface {
	Create(member *models.TeamMember) error
	GetByID(id uint) (*models.TeamMember, error)
	GetByTeamAndUser(teamID, userID uint) (*models.TeamMember, error)
	Update(member *models.TeamMember) error
	Delete(id uint) error
	DeleteByTeamAndUser(teamID, userID uint) error
	ListByTeam(teamID uint, limit, offset int) ([]*models.TeamMember, int64, error)
	CountByTeam(teamID uint) (int64, error)
	GetUserTeams(userID uint) ([]*models.TeamMember, error)
}

// TeamContestEntryRepositoryInterface defines the contract
type TeamContestEntryRepositoryInterface interface {
	Create(entry *models.TeamContestEntry) error
	GetByTeamAndContest(teamID, contestID uint) (*models.TeamContestEntry, error)
	Delete(teamID, contestID uint) error
	ListByContest(contestID uint, limit int) ([]*models.TeamContestEntry, error)
	UpdatePoints(teamID, contestID uint, points float64) error
	UpdateRanks(contestID uint) error
}

// TeamRepository implements TeamRepositoryInterface
type TeamRepository struct {
	db *gorm.DB
}

// TeamMemberRepository implements TeamMemberRepositoryInterface
type TeamMemberRepository struct {
	db *gorm.DB
}

// TeamContestEntryRepository implements TeamContestEntryRepositoryInterface
type TeamContestEntryRepository struct {
	db *gorm.DB
}

// NewTeamRepository creates a new team repository
func NewTeamRepository(db *gorm.DB) TeamRepositoryInterface {
	return &TeamRepository{db: db}
}

// NewTeamMemberRepository creates a new team member repository
func NewTeamMemberRepository(db *gorm.DB) TeamMemberRepositoryInterface {
	return &TeamMemberRepository{db: db}
}

// NewTeamContestEntryRepository creates a new team contest entry repository
func NewTeamContestEntryRepository(db *gorm.DB) TeamContestEntryRepositoryInterface {
	return &TeamContestEntryRepository{db: db}
}

// Team Repository Methods

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
	result := r.db.First(&team, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, result.Error
	}
	return &team, nil
}

func (r *TeamRepository) GetByInviteCode(code string) (*models.Team, error) {
	if code == "" {
		return nil, errors.New("invite code cannot be empty")
	}
	var team models.Team
	result := r.db.Where("invite_code = ? AND is_active = ?", code, true).First(&team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid invite code")
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
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete team members first
	if err := tx.Where("team_id = ?", id).Delete(&models.TeamMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete team contest entries
	if err := tx.Where("team_id = ?", id).Delete(&models.TeamContestEntry{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete team
	result := tx.Delete(&models.Team{}, id)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("team not found")
	}

	return tx.Commit().Error
}

func (r *TeamRepository) List(limit, offset int, userID uint, myTeamsOnly bool) ([]*models.Team, int64, error) {
	var teams []*models.Team
	var total int64

	query := r.db.Model(&models.Team{}).Where("is_active = ?", true)

	if myTeamsOnly && userID > 0 {
		// Get teams where user is a member
		query = query.Where("id IN (SELECT team_id FROM team_members WHERE user_id = ? AND status = ?)", userID, "active")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

// TeamMember Repository Methods

func (r *TeamMemberRepository) Create(member *models.TeamMember) error {
	if member == nil {
		return errors.New("member cannot be nil")
	}
	return r.db.Create(member).Error
}

func (r *TeamMemberRepository) GetByID(id uint) (*models.TeamMember, error) {
	if id == 0 {
		return nil, errors.New("invalid member ID")
	}
	var member models.TeamMember
	result := r.db.Preload("Team").First(&member, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, result.Error
	}
	return &member, nil
}

func (r *TeamMemberRepository) GetByTeamAndUser(teamID, userID uint) (*models.TeamMember, error) {
	if teamID == 0 || userID == 0 {
		return nil, errors.New("invalid team or user ID")
	}
	var member models.TeamMember
	result := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, result.Error
	}
	return &member, nil
}

func (r *TeamMemberRepository) Update(member *models.TeamMember) error {
	if member == nil {
		return errors.New("member cannot be nil")
	}
	if member.ID == 0 {
		return errors.New("member ID cannot be zero")
	}
	result := r.db.Save(member)
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return result.Error
}

func (r *TeamMemberRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid member ID")
	}
	result := r.db.Delete(&models.TeamMember{}, id)
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return result.Error
}

func (r *TeamMemberRepository) DeleteByTeamAndUser(teamID, userID uint) error {
	if teamID == 0 || userID == 0 {
		return errors.New("invalid team or user ID")
	}
	result := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamMember{})
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return result.Error
}

func (r *TeamMemberRepository) ListByTeam(teamID uint, limit, offset int) ([]*models.TeamMember, int64, error) {
	if teamID == 0 {
		return nil, 0, errors.New("invalid team ID")
	}
	var members []*models.TeamMember
	var total int64

	query := r.db.Where("team_id = ?", teamID)

	if err := query.Model(&models.TeamMember{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("role DESC, joined_at ASC").Limit(limit).Offset(offset).Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *TeamMemberRepository) CountByTeam(teamID uint) (int64, error) {
	if teamID == 0 {
		return 0, errors.New("invalid team ID")
	}
	var count int64
	result := r.db.Model(&models.TeamMember{}).Where("team_id = ? AND status = ?", teamID, "active").Count(&count)
	return count, result.Error
}

func (r *TeamMemberRepository) GetUserTeams(userID uint) ([]*models.TeamMember, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	var members []*models.TeamMember
	result := r.db.Preload("Team").Where("user_id = ? AND status = ?", userID, "active").Find(&members)
	return members, result.Error
}

// TeamContestEntry Repository Methods

func (r *TeamContestEntryRepository) Create(entry *models.TeamContestEntry) error {
	if entry == nil {
		return errors.New("entry cannot be nil")
	}
	return r.db.Create(entry).Error
}

func (r *TeamContestEntryRepository) GetByTeamAndContest(teamID, contestID uint) (*models.TeamContestEntry, error) {
	if teamID == 0 || contestID == 0 {
		return nil, errors.New("invalid team or contest ID")
	}
	var entry models.TeamContestEntry
	result := r.db.Where("team_id = ? AND contest_id = ?", teamID, contestID).First(&entry)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("entry not found")
		}
		return nil, result.Error
	}
	return &entry, nil
}

func (r *TeamContestEntryRepository) Delete(teamID, contestID uint) error {
	if teamID == 0 || contestID == 0 {
		return errors.New("invalid team or contest ID")
	}
	result := r.db.Where("team_id = ? AND contest_id = ?", teamID, contestID).Delete(&models.TeamContestEntry{})
	if result.RowsAffected == 0 {
		return errors.New("entry not found")
	}
	return result.Error
}

func (r *TeamContestEntryRepository) ListByContest(contestID uint, limit int) ([]*models.TeamContestEntry, error) {
	if contestID == 0 {
		return nil, errors.New("invalid contest ID")
	}
	if limit <= 0 {
		limit = 10
	}
	var entries []*models.TeamContestEntry
	result := r.db.Preload("Team").Where("contest_id = ?", contestID).Order("total_points DESC").Limit(limit).Find(&entries)
	return entries, result.Error
}

func (r *TeamContestEntryRepository) UpdatePoints(teamID, contestID uint, points float64) error {
	if teamID == 0 || contestID == 0 {
		return errors.New("invalid team or contest ID")
	}
	result := r.db.Model(&models.TeamContestEntry{}).
		Where("team_id = ? AND contest_id = ?", teamID, contestID).
		Update("total_points", points)
	return result.Error
}

func (r *TeamContestEntryRepository) UpdateRanks(contestID uint) error {
	if contestID == 0 {
		return errors.New("invalid contest ID")
	}
	// Update ranks based on total_points
	sql := `
		UPDATE team_contest_entries 
		SET rank = subquery.new_rank
		FROM (
			SELECT id, ROW_NUMBER() OVER (ORDER BY total_points DESC) as new_rank
			FROM team_contest_entries
			WHERE contest_id = ? AND deleted_at IS NULL
		) AS subquery
		WHERE team_contest_entries.id = subquery.id
	`
	return r.db.Exec(sql, contestID).Error
}
```

- **VALIDATE**: `cat backend/contest-service/internal/repository/team_repository.go | head -50`


---

### Task 9: CREATE backend/contest-service/internal/service/team_service.go

- **IMPLEMENT**: TeamService gRPC implementation
- **PATTERN**: Mirror contest_service.go exactly
- **IMPORTS**: context, log, time, models, repository, auth, proto packages

The service should implement all RPC methods from team.proto:
- CreateTeam: Create team, add captain as member, generate invite code
- UpdateTeam: Only captain can update
- GetTeam: Get team by ID
- DeleteTeam: Only captain can delete
- ListTeams: List all or user's teams
- JoinTeam: Join via invite code
- LeaveTeam: Leave team (captain cannot leave)
- RemoveMember: Captain removes member
- ListMembers: List team members
- RegenerateInviteCode: Captain regenerates code
- JoinContestAsTeam: Register team for contest
- LeaveContestAsTeam: Unregister team
- GetTeamLeaderboard: Get team rankings for contest
- Check: Health check

Key implementation notes:
- Extract userID from JWT context using `auth.GetUserIDFromContext(ctx)`
- Use `timestamppb.Now()` for timestamps
- Return proper error codes from `common.ErrorCode`
- Update CurrentMembers count after join/leave operations
- Only captain can update/delete team or remove members

- **VALIDATE**: `cat backend/contest-service/internal/service/team_service.go | head -100`

---

### Task 10: UPDATE backend/api-gateway/internal/config/config.go

- **IMPLEMENT**: Add TeamService endpoint configuration
- **PATTERN**: Mirror existing service endpoint patterns

Add to Config struct:
```go
TeamService string `envconfig:"TEAM_SERVICE_ENDPOINT" default:"localhost:8085"`
```

Note: Team service runs on same port as contest-service (8085) since it's part of contest-service.

- **VALIDATE**: `grep -n "TeamService" backend/api-gateway/internal/config/config.go`

---

### Task 11: UPDATE backend/api-gateway/internal/gateway/gateway.go

- **IMPLEMENT**: Register TeamService handler
- **PATTERN**: Mirror existing service registrations
- **IMPORTS**: Add team proto import

Add import:
```go
teampb "github.com/sports-prediction-contests/shared/proto/team"
```

Add registration in NewServer or registerServices:
```go
if err := teampb.RegisterTeamServiceHandlerFromEndpoint(ctx, mux, cfg.TeamService, opts); err != nil {
    return nil, fmt.Errorf("failed to register team service: %w", err)
}
```

- **VALIDATE**: `grep -n "team" backend/api-gateway/internal/gateway/gateway.go`

---

### Task 12: UPDATE backend/contest-service/cmd/main.go

- **IMPLEMENT**: Initialize and register TeamService
- **PATTERN**: Mirror ContestService initialization

Add after ContestService initialization:
```go
// Initialize team repositories
teamRepo := repository.NewTeamRepository(db)
teamMemberRepo := repository.NewTeamMemberRepository(db)
teamContestEntryRepo := repository.NewTeamContestEntryRepository(db)

// Initialize team service
teamService := service.NewTeamService(teamRepo, teamMemberRepo, teamContestEntryRepo)

// Register team service (note: requires adding to proto registration)
// teampb.RegisterTeamServiceServer(grpcServer, teamService)
```

- **VALIDATE**: `grep -n "team" backend/contest-service/cmd/main.go`

---

### Task 13: CREATE frontend/src/types/team.types.ts

- **IMPLEMENT**: TypeScript interfaces matching team.proto
- **PATTERN**: Mirror contest.types.ts exactly

```typescript
// Team types matching backend proto definitions

export interface Team {
  id: number
  name: string
  description: string
  inviteCode: string
  captainId: number
  maxMembers: number
  currentMembers: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface TeamMember {
  id: number
  teamId: number
  userId: number
  userName: string
  role: 'captain' | 'member'
  status: 'active' | 'inactive'
  joinedAt: string
}

export interface TeamLeaderboardEntry {
  teamId: number
  teamName: string
  totalPoints: number
  rank: number
  memberCount: number
}

// Request types
export interface CreateTeamRequest {
  name: string
  description: string
  maxMembers: number
}

export interface UpdateTeamRequest {
  id: number
  name: string
  description: string
  maxMembers: number
}

export interface ListTeamsRequest {
  pagination?: PaginationRequest
  myTeamsOnly?: boolean
}

export interface JoinTeamRequest {
  inviteCode: string
}

export interface LeaveTeamRequest {
  teamId: number
}

export interface RemoveMemberRequest {
  teamId: number
  userId: number
}

export interface ListMembersRequest {
  teamId: number
  pagination?: PaginationRequest
}

export interface JoinContestAsTeamRequest {
  teamId: number
  contestId: number
}

export interface GetTeamLeaderboardRequest {
  contestId: number
  limit?: number
}

// Response types
export interface ApiResponse {
  success: boolean
  message: string
  code: number
  timestamp: string
}

export interface CreateTeamResponse {
  response: ApiResponse
  team: Team
}

export interface ListTeamsResponse {
  response: ApiResponse
  teams: Team[]
  pagination: PaginationResponse
}

export interface ListMembersResponse {
  response: ApiResponse
  members: TeamMember[]
  pagination: PaginationResponse
}

export interface GetTeamLeaderboardResponse {
  response: ApiResponse
  entries: TeamLeaderboardEntry[]
}

// Common types (import from common.types.ts if available)
export interface PaginationRequest {
  page: number
  limit: number
}

export interface PaginationResponse {
  page: number
  limit: number
  total: number
  totalPages: number
}

// Form types
export interface TeamFormData {
  name: string
  description: string
  maxMembers: number
}
```

- **VALIDATE**: `cat frontend/src/types/team.types.ts | head -30`

---

### Task 14: CREATE frontend/src/utils/team-validation.ts

- **IMPLEMENT**: Zod validation schemas for team forms
- **PATTERN**: Mirror sports-validation.ts

```typescript
import { z } from 'zod'
import type { TeamFormData } from '../types/team.types'

export const teamSchema = z.object({
  name: z
    .string()
    .min(1, 'Team name is required')
    .max(100, 'Team name cannot exceed 100 characters'),
  description: z
    .string()
    .max(500, 'Description cannot exceed 500 characters')
    .optional()
    .default(''),
  maxMembers: z
    .number()
    .min(2, 'Team must allow at least 2 members')
    .max(50, 'Team cannot exceed 50 members')
    .default(10),
})

export const joinTeamSchema = z.object({
  inviteCode: z
    .string()
    .min(1, 'Invite code is required')
    .max(20, 'Invalid invite code'),
})

export type TeamSchemaType = z.infer<typeof teamSchema>
export type JoinTeamSchemaType = z.infer<typeof joinTeamSchema>

export const teamFormDataToRequest = (data: TeamSchemaType): TeamFormData => ({
  name: data.name,
  description: data.description || '',
  maxMembers: data.maxMembers,
})
```

- **VALIDATE**: `cat frontend/src/utils/team-validation.ts`

---

### Task 15: CREATE frontend/src/services/team-service.ts

- **IMPLEMENT**: Team API service class
- **PATTERN**: Mirror contest-service.ts exactly

```typescript
import grpcClient from './grpc-client'
import type {
  Team,
  TeamMember,
  TeamLeaderboardEntry,
  CreateTeamRequest,
  UpdateTeamRequest,
  ListTeamsRequest,
  JoinTeamRequest,
  LeaveTeamRequest,
  RemoveMemberRequest,
  ListMembersRequest,
  JoinContestAsTeamRequest,
  GetTeamLeaderboardRequest,
  CreateTeamResponse,
  ListTeamsResponse,
  ListMembersResponse,
  GetTeamLeaderboardResponse,
  PaginationResponse,
} from '../types/team.types'

class TeamService {
  private basePath = '/v1/teams'

  async createTeam(request: CreateTeamRequest): Promise<Team> {
    const response = await grpcClient.post<CreateTeamResponse>(this.basePath, request)
    return response.team
  }

  async updateTeam(request: UpdateTeamRequest): Promise<Team> {
    const response = await grpcClient.put<CreateTeamResponse>(
      `${this.basePath}/${request.id}`,
      request
    )
    return response.team
  }

  async getTeam(id: number): Promise<Team> {
    const response = await grpcClient.get<{ response: any; team: Team }>(
      `${this.basePath}/${id}`
    )
    return response.team
  }

  async deleteTeam(id: number): Promise<void> {
    await grpcClient.delete(`${this.basePath}/${id}`)
  }

  async listTeams(request: ListTeamsRequest = {}): Promise<{
    teams: Team[]
    pagination: PaginationResponse
  }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.myTeamsOnly) {
      params.append('my_teams_only', 'true')
    }
    const queryString = params.toString()
    const url = queryString ? `${this.basePath}?${queryString}` : this.basePath
    const response = await grpcClient.get<ListTeamsResponse>(url)
    return {
      teams: response.teams || [],
      pagination: response.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 },
    }
  }

  async joinTeam(request: JoinTeamRequest): Promise<TeamMember> {
    const response = await grpcClient.post<{ response: any; member: TeamMember }>(
      `${this.basePath}/join`,
      request
    )
    return response.member
  }

  async leaveTeam(teamId: number): Promise<void> {
    await grpcClient.post(`${this.basePath}/${teamId}/leave`, {})
  }

  async removeMember(teamId: number, userId: number): Promise<void> {
    await grpcClient.delete(`${this.basePath}/${teamId}/members/${userId}`)
  }

  async listMembers(request: ListMembersRequest): Promise<{
    members: TeamMember[]
    pagination: PaginationResponse
  }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    const queryString = params.toString()
    const url = queryString
      ? `${this.basePath}/${request.teamId}/members?${queryString}`
      : `${this.basePath}/${request.teamId}/members`
    const response = await grpcClient.get<ListMembersResponse>(url)
    return {
      members: response.members || [],
      pagination: response.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 },
    }
  }

  async regenerateInviteCode(teamId: number): Promise<string> {
    const response = await grpcClient.post<{ response: any; inviteCode: string }>(
      `${this.basePath}/${teamId}/regenerate-invite`,
      {}
    )
    return response.inviteCode
  }

  async joinContestAsTeam(teamId: number, contestId: number): Promise<void> {
    await grpcClient.post(`${this.basePath}/${teamId}/contests/${contestId}/join`, {})
  }

  async leaveContestAsTeam(teamId: number, contestId: number): Promise<void> {
    await grpcClient.post(`${this.basePath}/${teamId}/contests/${contestId}/leave`, {})
  }

  async getTeamLeaderboard(contestId: number, limit: number = 10): Promise<TeamLeaderboardEntry[]> {
    const response = await grpcClient.get<GetTeamLeaderboardResponse>(
      `/v1/contests/${contestId}/team-leaderboard?limit=${limit}`
    )
    return response.entries || []
  }
}

export const teamService = new TeamService()
export default teamService
```

- **VALIDATE**: `cat frontend/src/services/team-service.ts | head -50`


---

### Task 16: CREATE frontend/src/hooks/use-teams.ts

- **IMPLEMENT**: React Query hooks for teams
- **PATTERN**: Mirror use-contests.ts exactly

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import teamService from '../services/team-service'
import { useToast } from '../contexts/ToastContext'
import type {
  CreateTeamRequest,
  UpdateTeamRequest,
  ListTeamsRequest,
  JoinTeamRequest,
  ListMembersRequest,
} from '../types/team.types'

export const teamKeys = {
  all: ['teams'] as const,
  lists: () => [...teamKeys.all, 'list'] as const,
  list: (filters: ListTeamsRequest) => [...teamKeys.lists(), filters] as const,
  details: () => [...teamKeys.all, 'detail'] as const,
  detail: (id: number) => [...teamKeys.details(), id] as const,
  members: (teamId: number) => [...teamKeys.all, 'members', teamId] as const,
  leaderboard: (contestId: number) => ['team-leaderboard', contestId] as const,
}

export const useTeams = (request: ListTeamsRequest = {}) => {
  return useQuery({
    queryKey: teamKeys.list(request),
    queryFn: () => teamService.listTeams(request),
    staleTime: 5 * 60 * 1000,
  })
}

export const useTeam = (id: number) => {
  return useQuery({
    queryKey: teamKeys.detail(id),
    queryFn: () => teamService.getTeam(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}

export const useTeamMembers = (request: ListMembersRequest) => {
  return useQuery({
    queryKey: teamKeys.members(request.teamId),
    queryFn: () => teamService.listMembers(request),
    enabled: !!request.teamId,
    staleTime: 2 * 60 * 1000,
  })
}

export const useTeamLeaderboard = (contestId: number, limit: number = 10) => {
  return useQuery({
    queryKey: teamKeys.leaderboard(contestId),
    queryFn: () => teamService.getTeamLeaderboard(contestId, limit),
    enabled: !!contestId,
    staleTime: 30 * 1000,
  })
}

export const useCreateTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: CreateTeamRequest) => teamService.createTeam(request),
    onSuccess: (newTeam) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      queryClient.setQueryData(teamKeys.detail(newTeam.id), newTeam)
      showToast('Team created successfully!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to create team: ${error.message}`, 'error')
    },
  })
}

export const useUpdateTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: UpdateTeamRequest) => teamService.updateTeam(request),
    onSuccess: (updatedTeam) => {
      queryClient.setQueryData(teamKeys.detail(updatedTeam.id), updatedTeam)
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Team updated successfully!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to update team: ${error.message}`, 'error')
    },
  })
}

export const useDeleteTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (id: number) => teamService.deleteTeam(id),
    onSuccess: (_, deletedId) => {
      queryClient.removeQueries({ queryKey: teamKeys.detail(deletedId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Team deleted successfully!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to delete team: ${error.message}`, 'error')
    },
  })
}

export const useJoinTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: JoinTeamRequest) => teamService.joinTeam(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Successfully joined team!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to join team: ${error.message}`, 'error')
    },
  })
}

export const useLeaveTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (teamId: number) => teamService.leaveTeam(teamId),
    onSuccess: (_, teamId) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.detail(teamId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.members(teamId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Successfully left team!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to leave team: ${error.message}`, 'error')
    },
  })
}

export const useRemoveMember = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: ({ teamId, userId }: { teamId: number; userId: number }) =>
      teamService.removeMember(teamId, userId),
    onSuccess: (_, { teamId }) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.detail(teamId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.members(teamId) })
      showToast('Member removed successfully!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to remove member: ${error.message}`, 'error')
    },
  })
}

export const useRegenerateInviteCode = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (teamId: number) => teamService.regenerateInviteCode(teamId),
    onSuccess: (_, teamId) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.detail(teamId) })
      showToast('Invite code regenerated!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to regenerate code: ${error.message}`, 'error')
    },
  })
}

export const useJoinContestAsTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: ({ teamId, contestId }: { teamId: number; contestId: number }) =>
      teamService.joinContestAsTeam(teamId, contestId),
    onSuccess: (_, { contestId }) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.leaderboard(contestId) })
      showToast('Team joined contest!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to join contest: ${error.message}`, 'error')
    },
  })
}
```

- **VALIDATE**: `cat frontend/src/hooks/use-teams.ts | head -50`

---

### Task 17: CREATE frontend/src/components/teams/TeamList.tsx

- **IMPLEMENT**: Team list with MaterialReactTable
- **PATTERN**: Mirror ContestList.tsx exactly

Key features:
- Display team name, members count, captain indicator
- Actions: View Members, Edit (captain only), Delete (captain only)
- Create Team button in toolbar
- Pagination with URL sync

- **VALIDATE**: `cat frontend/src/components/teams/TeamList.tsx | head -50`

---

### Task 18: CREATE frontend/src/components/teams/TeamForm.tsx

- **IMPLEMENT**: Create/Edit team dialog form
- **PATTERN**: Mirror ContestForm.tsx
- **IMPORTS**: react-hook-form, @hookform/resolvers/zod, team-validation

Key features:
- Name field (required, max 100 chars)
- Description field (optional, max 500 chars)
- Max Members field (2-50, default 10)
- Dialog with Cancel/Submit buttons

- **VALIDATE**: `cat frontend/src/components/teams/TeamForm.tsx | head -50`

---

### Task 19: CREATE frontend/src/components/teams/TeamMembers.tsx

- **IMPLEMENT**: Team members list with management
- **PATTERN**: Mirror ParticipantList.tsx

Key features:
- List members with role badges (Captain/Member)
- Remove button for captain (not self)
- Show joined date
- Pagination

- **VALIDATE**: `cat frontend/src/components/teams/TeamMembers.tsx | head -50`

---

### Task 20: CREATE frontend/src/components/teams/TeamInvite.tsx

- **IMPLEMENT**: Invite code display and copy component

Key features:
- Display invite code prominently
- Copy to clipboard button
- Regenerate button (captain only)
- Share link generation

```typescript
import React from 'react'
import { Box, Typography, Button, TextField, IconButton, Tooltip } from '@mui/material'
import { ContentCopy, Refresh } from '@mui/icons-material'
import { useToast } from '../../contexts/ToastContext'
import { useRegenerateInviteCode } from '../../hooks/use-teams'

interface TeamInviteProps {
  teamId: number
  inviteCode: string
  isCaptain: boolean
}

export const TeamInvite: React.FC<TeamInviteProps> = ({ teamId, inviteCode, isCaptain }) => {
  const { showToast } = useToast()
  const regenerateMutation = useRegenerateInviteCode()

  const handleCopy = () => {
    navigator.clipboard.writeText(inviteCode)
    showToast('Invite code copied!', 'success')
  }

  const handleRegenerate = () => {
    if (window.confirm('Regenerate invite code? The old code will stop working.')) {
      regenerateMutation.mutate(teamId)
    }
  }

  return (
    <Box sx={{ p: 2, bgcolor: 'background.paper', borderRadius: 1 }}>
      <Typography variant="subtitle2" color="text.secondary" gutterBottom>
        Invite Code
      </Typography>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <TextField
          value={inviteCode}
          InputProps={{ readOnly: true }}
          size="small"
          sx={{ fontFamily: 'monospace', flex: 1 }}
        />
        <Tooltip title="Copy code">
          <IconButton onClick={handleCopy} color="primary">
            <ContentCopy />
          </IconButton>
        </Tooltip>
        {isCaptain && (
          <Tooltip title="Regenerate code">
            <IconButton
              onClick={handleRegenerate}
              disabled={regenerateMutation.isPending}
              color="warning"
            >
              <Refresh />
            </IconButton>
          </Tooltip>
        )}
      </Box>
    </Box>
  )
}

export default TeamInvite
```

- **VALIDATE**: `cat frontend/src/components/teams/TeamInvite.tsx`

---

### Task 21: CREATE frontend/src/components/teams/TeamLeaderboard.tsx

- **IMPLEMENT**: Team rankings table for contests
- **PATTERN**: Mirror LeaderboardTable.tsx

Key features:
- Rank, Team Name, Total Points, Member Count columns
- Highlight user's team
- Trophy icons for top 3

- **VALIDATE**: `cat frontend/src/components/teams/TeamLeaderboard.tsx | head -50`

---

### Task 22: CREATE frontend/src/pages/TeamsPage.tsx

- **IMPLEMENT**: Main teams page with tabs
- **PATTERN**: Mirror ContestsPage.tsx structure

Tabs:
1. My Teams - Teams user is member of
2. All Teams - Browse all teams
3. Join Team - Enter invite code

Key features:
- Tab navigation
- Create Team button
- Team detail dialog with members and invite code
- Join team form

- **VALIDATE**: `cat frontend/src/pages/TeamsPage.tsx | head -50`

---

### Task 23: UPDATE frontend/src/App.tsx

- **IMPLEMENT**: Add /teams route and navigation
- **PATTERN**: Mirror existing route additions

Add import:
```typescript
import TeamsPage from './pages/TeamsPage'
```

Add route:
```typescript
<Route path="/teams" element={<ProtectedRoute><TeamsPage /></ProtectedRoute>} />
```

Add navigation link in header/menu.

- **VALIDATE**: `grep -n "teams" frontend/src/App.tsx`

---

### Task 24: CREATE tests/contest-service/team_test.go

- **IMPLEMENT**: Unit tests for team models
- **PATTERN**: Mirror contest_test.go

Test cases:
- Team name validation (empty, too long)
- Description validation (too long)
- MaxMembers validation (< 2, > 50)
- InviteCode generation
- TeamMember role validation
- TeamMember duplicate check

- **VALIDATE**: `cat tests/contest-service/team_test.go | head -50`

---

## TESTING STRATEGY

### Unit Tests
- Team model validation (name, description, maxMembers)
- TeamMember model validation (role, status, duplicate check)
- InviteCode generation uniqueness
- Repository CRUD operations

### Integration Tests
- Create team → captain becomes member
- Join team via invite code
- Leave team updates member count
- Captain cannot leave team
- Only captain can remove members
- Team leaderboard aggregation

### Edge Cases
- Join with invalid invite code
- Join full team (maxMembers reached)
- Remove last member (should fail)
- Delete team with active contest entries
- Concurrent join requests

---

## VALIDATION COMMANDS

### Level 1: Syntax & Style
```bash
cd backend && go build ./...
cd frontend && npm run lint
```

### Level 2: Unit Tests
```bash
cd tests/contest-service && go test -v ./...
cd frontend && npm test
```

### Level 3: Integration Tests
```bash
# Start services
docker-compose up -d postgres redis
cd backend/contest-service && go run cmd/main.go &

# Test API endpoints
curl -X POST http://localhost:8080/v1/teams -H "Authorization: Bearer $TOKEN" -d '{"name":"Test Team","maxMembers":10}'
curl http://localhost:8080/v1/teams
```

### Level 4: Manual Validation
1. Create a team and verify invite code generated
2. Copy invite code and join from another account
3. Verify member count updates
4. Test captain-only actions (edit, delete, remove member)
5. Join contest as team and verify team leaderboard

---

## ACCEPTANCE CRITERIA

- [ ] Teams can be created with name, description, max members
- [ ] Unique invite codes are generated automatically
- [ ] Users can join teams via invite code
- [ ] Captain can manage team (edit, delete, remove members)
- [ ] Member count updates correctly on join/leave
- [ ] Teams can join contests
- [ ] Team leaderboard shows aggregated scores
- [ ] All validation commands pass
- [ ] Unit tests cover model validation
- [ ] Frontend displays teams with CRUD operations

---

## COMPLETION CHECKLIST

- [ ] Database tables created (teams, team_members, team_contest_entries)
- [ ] Proto definitions complete
- [ ] Backend models with validation
- [ ] Repository with CRUD operations
- [ ] gRPC service implementation
- [ ] API gateway registration
- [ ] Frontend types and validation
- [ ] Frontend service and hooks
- [ ] Frontend components (List, Form, Members, Invite, Leaderboard)
- [ ] Teams page with routing
- [ ] Unit tests passing
- [ ] Manual testing complete

---

## NOTES

### Design Decisions
- Teams are managed within contest-service (not separate service) to simplify architecture
- Invite codes are 8-character hex strings (4 bytes = 4 billion combinations)
- Captain role is immutable (cannot transfer ownership in MVP)
- Team leaderboard aggregates individual member scores

### Future Enhancements
- Transfer captain role
- Team chat/activity feed
- Team achievements/badges
- Private teams (invite-only visibility)
- Team statistics dashboard

### Risks
- Race condition on concurrent joins - mitigated by database unique constraint
- Orphaned team entries if captain deletes account - handle in user deletion flow

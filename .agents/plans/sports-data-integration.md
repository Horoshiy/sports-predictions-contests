# Feature: Sports Data Integration

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

Integrate with TheSportsDB free API to automatically sync sports data (leagues, teams, matches, and results) into the platform. This enables the platform to have real sports data without manual entry, and automatically updates match results for prediction scoring.

**Key Capabilities:**
- Sync sports, leagues, and teams from external API
- Import upcoming fixtures/matches
- Fetch and update match results when completed
- Background worker for periodic sync
- Manual sync trigger via admin API

## User Story

As a platform administrator
I want to automatically sync sports data from external APIs
So that users have real matches to predict on without manual data entry

## Problem Statement

Currently, all sports data (leagues, teams, matches) must be manually entered through the admin UI. This is time-consuming and error-prone. Match results also need manual updates, which delays prediction scoring.

## Solution Statement

Create a sports data sync service that:
1. Connects to TheSportsDB free API (no auth required for basic data)
2. Syncs leagues, teams, and upcoming matches on schedule
3. Updates match results when games complete
4. Stores external IDs for deduplication
5. Provides manual sync endpoints for admins

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium (4-6 hours)
**Primary Systems Affected**: sports-service
**Dependencies**: TheSportsDB API (free tier, no API key required for basic endpoints)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/sports-service/internal/service/sports_service.go` - Existing sports service with CRUD operations
- `backend/sports-service/internal/models/match.go` - Match model with status field
- `backend/sports-service/internal/models/sport.go` - Sport model structure
- `backend/sports-service/internal/models/league.go` - League model structure
- `backend/sports-service/internal/models/team.go` - Team model structure
- `backend/sports-service/internal/repository/match_repository.go` - Match repository interface
- `backend/sports-service/internal/config/config.go` - Service configuration pattern
- `backend/notification-service/internal/worker/worker.go` - Worker pool pattern to follow
- `backend/proto/sports.proto` - gRPC proto definitions

### New Files to Create

- `backend/sports-service/internal/external/thesportsdb.go` - TheSportsDB API client
- `backend/sports-service/internal/sync/sync_service.go` - Sync orchestration service
- `backend/sports-service/internal/sync/worker.go` - Background sync worker
- `tests/sports-service/sync_test.go` - Unit tests for sync logic

### Files to Modify

- `backend/sports-service/internal/models/sport.go` - Add ExternalID field
- `backend/sports-service/internal/models/league.go` - Add ExternalID field
- `backend/sports-service/internal/models/team.go` - Add ExternalID field
- `backend/sports-service/internal/models/match.go` - Add ExternalID field
- `backend/sports-service/internal/config/config.go` - Add sync configuration
- `backend/sports-service/cmd/main.go` - Initialize sync worker
- `backend/proto/sports.proto` - Add sync RPC methods
- `scripts/init-db.sql` - Add external_id columns and indexes

### Relevant Documentation

- [TheSportsDB API](https://www.thesportsdb.com/api.php)
  - Free tier endpoints (no API key needed for basic data)
  - Key endpoints: all_sports, all_leagues, lookup_all_teams, eventsnext, lookupevent
- [Go HTTP Client Best Practices](https://pkg.go.dev/net/http)
  - Timeout configuration, connection pooling

### Patterns to Follow

**Config Pattern (from config.go):**
```go
type Config struct {
    Port        string
    JWTSecret   string
    DatabaseURL string
    LogLevel    string
}

func Load() *Config {
    return &Config{
        Port: getEnvOrDefault("SPORTS_SERVICE_PORT", "8088"),
        // ...
    }
}
```

**Worker Pool Pattern (from notification-service/worker.go):**
```go
type WorkerPool struct {
    jobs chan Job
    quit chan bool
    wg   sync.WaitGroup
}

func (w *WorkerPool) Start() {
    w.wg.Add(1)
    go w.worker()
}

func (w *WorkerPool) Stop() {
    close(w.quit)
    w.wg.Wait()
}
```

**Repository Pattern (from match_repository.go):**
```go
type MatchRepositoryInterface interface {
    Create(match *models.Match) error
    GetByID(id uint) (*models.Match, error)
    Update(match *models.Match) error
}
```

**Naming Conventions:**
- Go files: `snake_case.go`
- Go structs: `PascalCase`
- Go functions: `PascalCase` (public), `camelCase` (private)
- External IDs: `external_id` in DB, `ExternalID` in Go

---

## IMPLEMENTATION PLAN

### Phase 1: Database Schema Updates

Add external_id fields to existing models for deduplication when syncing from external API.

### Phase 2: External API Client

Create HTTP client for TheSportsDB API with proper timeout handling and response parsing.

### Phase 3: Sync Service

Build sync orchestration that maps external data to internal models and handles upserts.

### Phase 4: Background Worker

Implement periodic sync worker that runs on configurable schedule.

### Phase 5: API Endpoints

Add gRPC endpoints for manual sync triggers and sync status.

---

## STEP-BY-STEP TASKS

### Task 1: UPDATE `scripts/init-db.sql` - Add external_id columns

- **IMPLEMENT**: Add `external_id` column to sports, leagues, teams, matches tables
- **PATTERN**: Follow existing column patterns in init-db.sql
- **GOTCHA**: Use VARCHAR for external_id (TheSportsDB uses string IDs)
- **VALIDATE**: `grep -n "external_id" scripts/init-db.sql`

```sql
-- Add to sports table after icon_url
ALTER TABLE sports ADD COLUMN IF NOT EXISTS external_id VARCHAR(50) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_sports_external_id ON sports(external_id);

-- Add to leagues table after season
ALTER TABLE leagues ADD COLUMN IF NOT EXISTS external_id VARCHAR(50) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_leagues_external_id ON leagues(external_id);

-- Add to teams table after country
ALTER TABLE teams ADD COLUMN IF NOT EXISTS external_id VARCHAR(50) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_teams_external_id ON teams(external_id);

-- Add to matches table after result_data
ALTER TABLE matches ADD COLUMN IF NOT EXISTS external_id VARCHAR(50) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_matches_external_id ON matches(external_id);
```

### Task 2: UPDATE `backend/sports-service/internal/models/sport.go` - Add ExternalID

- **IMPLEMENT**: Add ExternalID field to Sport struct
- **PATTERN**: Follow existing field patterns with gorm tags
- **VALIDATE**: `grep -n "ExternalID" backend/sports-service/internal/models/sport.go`

Add field after IconURL:
```go
ExternalID string `gorm:"uniqueIndex;size:50" json:"external_id,omitempty"`
```

### Task 3: UPDATE `backend/sports-service/internal/models/league.go` - Add ExternalID

- **IMPLEMENT**: Add ExternalID field to League struct
- **PATTERN**: Same as sport.go
- **VALIDATE**: `grep -n "ExternalID" backend/sports-service/internal/models/league.go`

### Task 4: UPDATE `backend/sports-service/internal/models/team.go` - Add ExternalID

- **IMPLEMENT**: Add ExternalID field to Team struct
- **PATTERN**: Same as sport.go
- **VALIDATE**: `grep -n "ExternalID" backend/sports-service/internal/models/team.go`

### Task 5: UPDATE `backend/sports-service/internal/models/match.go` - Add ExternalID

- **IMPLEMENT**: Add ExternalID field to Match struct
- **PATTERN**: Same as sport.go
- **VALIDATE**: `grep -n "ExternalID" backend/sports-service/internal/models/match.go`

### Task 6: UPDATE `backend/sports-service/internal/config/config.go` - Add sync config

- **IMPLEMENT**: Add sync-related configuration fields
- **PATTERN**: Follow existing getEnvOrDefault pattern
- **VALIDATE**: `grep -n "Sync" backend/sports-service/internal/config/config.go`

Add fields:
```go
SyncEnabled      bool
SyncIntervalMins int
TheSportsDBURL   string
```

### Task 7: CREATE `backend/sports-service/internal/external/thesportsdb.go` - API client

- **IMPLEMENT**: HTTP client for TheSportsDB API
- **IMPORTS**: `net/http`, `encoding/json`, `time`, `context`
- **GOTCHA**: Use 10-second timeout, handle rate limits gracefully
- **VALIDATE**: `go build ./backend/sports-service/...`

Key structures:
```go
type Client struct {
    httpClient *http.Client
    baseURL    string
}

type SportsResponse struct {
    Sports []APISport `json:"sports"`
}

type APISport struct {
    IDSport     string `json:"idSport"`
    StrSport    string `json:"strSport"`
    StrFormat   string `json:"strFormat"`
    StrSportThumb string `json:"strSportThumb"`
}

// Similar for leagues, teams, events
```

Key methods:
- `GetAllSports() ([]APISport, error)`
- `GetLeaguesBySport(sportID string) ([]APILeague, error)`
- `GetTeamsByLeague(leagueID string) ([]APITeam, error)`
- `GetUpcomingEvents(leagueID string) ([]APIEvent, error)`
- `GetEventByID(eventID string) (*APIEvent, error)`

### Task 8: CREATE `backend/sports-service/internal/sync/sync_service.go` - Sync orchestration

- **IMPLEMENT**: Service that orchestrates syncing from external API to database
- **PATTERN**: Follow service pattern from sports_service.go
- **IMPORTS**: Repository interfaces, external client, models
- **GOTCHA**: Use upsert logic - update if external_id exists, create if not
- **VALIDATE**: `go build ./backend/sports-service/...`

Key methods:
```go
type SyncService struct {
    client     *external.Client
    sportRepo  repository.SportRepositoryInterface
    leagueRepo repository.LeagueRepositoryInterface
    teamRepo   repository.TeamRepositoryInterface
    matchRepo  repository.MatchRepositoryInterface
}

func (s *SyncService) SyncSports(ctx context.Context) (int, error)
func (s *SyncService) SyncLeagues(ctx context.Context, sportID uint) (int, error)
func (s *SyncService) SyncTeams(ctx context.Context, leagueID uint) (int, error)
func (s *SyncService) SyncUpcomingMatches(ctx context.Context, leagueID uint) (int, error)
func (s *SyncService) SyncMatchResults(ctx context.Context) (int, error)
```

### Task 9: UPDATE repositories - Add GetByExternalID and Upsert methods

- **IMPLEMENT**: Add methods to find by external_id and upsert
- **FILES**: 
  - `backend/sports-service/internal/repository/sport_repository.go`
  - `backend/sports-service/internal/repository/league_repository.go`
  - `backend/sports-service/internal/repository/team_repository.go`
  - `backend/sports-service/internal/repository/match_repository.go`
- **PATTERN**: Follow existing repository method patterns
- **VALIDATE**: `go build ./backend/sports-service/...`

Add to each repository interface:
```go
GetByExternalID(externalID string) (*models.Sport, error)
Upsert(sport *models.Sport) error
```

### Task 10: CREATE `backend/sports-service/internal/sync/worker.go` - Background worker

- **IMPLEMENT**: Periodic sync worker using ticker
- **PATTERN**: Follow notification-service worker pattern
- **IMPORTS**: `time`, `sync`, `context`
- **GOTCHA**: Graceful shutdown, don't overlap syncs
- **VALIDATE**: `go build ./backend/sports-service/...`

```go
type SyncWorker struct {
    syncService *SyncService
    interval    time.Duration
    quit        chan bool
    wg          sync.WaitGroup
    running     bool
    mu          sync.Mutex
}

func (w *SyncWorker) Start()
func (w *SyncWorker) Stop()
func (w *SyncWorker) runSync()
```

### Task 11: UPDATE `backend/proto/sports.proto` - Add sync RPC methods

- **IMPLEMENT**: Add manual sync trigger endpoints
- **PATTERN**: Follow existing RPC patterns in sports.proto
- **VALIDATE**: `grep -n "Sync" backend/proto/sports.proto`

Add messages:
```protobuf
message SyncRequest {
  string entity_type = 1; // "sports", "leagues", "teams", "matches", "results"
  uint32 parent_id = 2;   // Optional: sport_id for leagues, league_id for teams/matches
}

message SyncResponse {
  common.Response response = 1;
  int32 synced_count = 2;
  string entity_type = 3;
}

message SyncStatusRequest {}

message SyncStatusResponse {
  common.Response response = 1;
  bool sync_enabled = 2;
  string last_sync_at = 3;
  int32 sync_interval_mins = 4;
}
```

Add RPCs:
```protobuf
rpc TriggerSync(SyncRequest) returns (SyncResponse) {
  option (google.api.http) = {
    post: "/v1/sports/sync"
    body: "*"
  };
}

rpc GetSyncStatus(SyncStatusRequest) returns (SyncStatusResponse) {
  option (google.api.http) = {
    get: "/v1/sports/sync/status"
  };
}
```

### Task 12: UPDATE `backend/sports-service/internal/service/sports_service.go` - Add sync methods

- **IMPLEMENT**: Implement TriggerSync and GetSyncStatus gRPC methods
- **PATTERN**: Follow existing service method patterns
- **GOTCHA**: Only allow admins to trigger sync (check user role if available)
- **VALIDATE**: `go build ./backend/sports-service/...`

### Task 13: UPDATE `backend/sports-service/cmd/main.go` - Initialize sync worker

- **IMPLEMENT**: Create and start sync worker on service startup
- **PATTERN**: Follow notification-service main.go pattern
- **GOTCHA**: Graceful shutdown - stop worker before server
- **VALIDATE**: `go build ./backend/sports-service/cmd/...`

### Task 14: CREATE `tests/sports-service/sync_test.go` - Unit tests

- **IMPLEMENT**: Tests for sync service and external client
- **PATTERN**: Follow existing test patterns in tests/sports-service/
- **VALIDATE**: `go test ./tests/sports-service/...`

Test cases:
- TestTheSportsDBClient_GetAllSports
- TestSyncService_SyncSports
- TestSyncService_MapExternalToInternal
- TestSyncWorker_StartStop

---

## TESTING STRATEGY

### Unit Tests

Based on existing patterns in `tests/sports-service/`:

```go
func TestTheSportsDBClient_GetAllSports(t *testing.T) {
    // Mock HTTP server
    // Verify response parsing
}

func TestSyncService_MapSportFromExternal(t *testing.T) {
    // Test mapping from APISport to models.Sport
}

func TestSyncWorker_GracefulShutdown(t *testing.T) {
    // Verify worker stops cleanly
}
```

### Integration Tests

- Start sports-service with sync enabled
- Verify sync worker runs on schedule
- Verify manual sync endpoint works

### Edge Cases

- External API timeout handling
- Duplicate external_id handling (upsert)
- Empty response from API
- Invalid data from API (missing required fields)
- Worker already running when sync triggered

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
# Build sports service
cd backend/sports-service && go build ./...

# Check for syntax errors
go vet ./backend/sports-service/...
```

### Level 2: Unit Tests

```bash
# Run sports service tests
go test ./tests/sports-service/... -v
```

### Level 3: Integration Tests

```bash
# Start services
docker-compose up -d postgres redis sports-service

# Test sync endpoint (requires running service)
curl -X POST http://localhost:8080/v1/sports/sync \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"entity_type": "sports"}'

# Check sync status
curl http://localhost:8080/v1/sports/sync/status \
  -H "Authorization: Bearer <token>"
```

### Level 4: Manual Validation

1. Start services with `docker-compose up`
2. Trigger sync via API or wait for scheduled sync
3. Verify sports/leagues/teams appear in database
4. Check logs for sync activity

---

## ACCEPTANCE CRITERIA

- [ ] External API client successfully fetches data from TheSportsDB
- [ ] Sports, leagues, teams sync correctly with deduplication
- [ ] Upcoming matches sync with proper team/league associations
- [ ] Match results update when games complete
- [ ] Background worker runs on configured schedule
- [ ] Manual sync endpoint works for admins
- [ ] Graceful shutdown stops worker cleanly
- [ ] All unit tests pass
- [ ] No build errors or warnings

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed
- [ ] All validation commands executed successfully
- [ ] Unit tests pass
- [ ] Manual testing confirms sync works
- [ ] Code follows project conventions
- [ ] No regressions in existing functionality

---

## NOTES

### TheSportsDB Free Tier Limitations

- No API key required for basic endpoints
- Rate limited (be respectful with request frequency)
- Some endpoints are Patreon-only (searchteams, etc.)
- Free endpoints available:
  - `all_sports.php` - List all sports
  - `all_leagues.php` - List all leagues
  - `lookup_all_teams.php?id={league_id}` - Teams by league
  - `eventsnext.php?id={team_id}` - Upcoming events for team
  - `lookupevent.php?id={event_id}` - Event details

### Design Decisions

1. **Upsert over Delete+Insert**: Preserves internal IDs and relationships
2. **External ID as String**: TheSportsDB uses string IDs
3. **Configurable Sync Interval**: Default 60 mins, adjustable via env
4. **Single Worker**: Prevents overlapping syncs, simpler than pool

### Future Enhancements

- Webhook support for real-time result updates
- Multiple API provider support (API-Football, etc.)
- Selective sync (specific leagues only)
- Sync history/audit log

### Confidence Score: 8/10

High confidence due to:
- Clear API documentation
- Existing patterns to follow (worker, repository)
- Well-defined scope

Risks:
- TheSportsDB API availability/rate limits
- Data mapping edge cases

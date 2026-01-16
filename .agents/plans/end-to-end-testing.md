# Feature: End-to-End Testing with Running Services

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to existing test patterns, API endpoints from proto files, and Docker service configuration.

## Feature Description

Comprehensive end-to-end (E2E) testing suite that validates the complete user journey through the Sports Prediction Contests platform with all microservices running. Tests will exercise the HTTP API Gateway endpoints, verify service-to-service communication, and validate the full prediction workflow from user registration to leaderboard updates.

## User Story

As a **developer/QA engineer**
I want to **run automated E2E tests against running services**
So that **I can validate the complete platform functionality before deployment and catch integration issues early**

## Problem Statement

The platform has 8 microservices with comprehensive unit tests, but lacks automated E2E tests that validate:
- Service-to-service communication via gRPC
- HTTP API Gateway routing and authentication
- Complete user workflows (register → create contest → predict → score → leaderboard)
- Database state consistency across services
- Error handling and edge cases in production-like environment

## Solution Statement

Create a Go-based E2E test suite that:
1. Uses Docker Compose to spin up all services
2. Tests HTTP endpoints via the API Gateway (port 8080)
3. Validates complete user journeys with real database operations
4. Includes setup/teardown for test isolation
5. Provides clear test output and failure diagnostics

## Feature Metadata

**Feature Type**: New Capability (Testing Infrastructure)
**Estimated Complexity**: Medium
**Primary Systems Affected**: All services via API Gateway
**Dependencies**: Docker, Docker Compose, Go testing framework, net/http

---

## CONTEXT REFERENCES

### Relevant Codebase Files - MUST READ BEFORE IMPLEMENTING

- `docker-compose.yml` - Service configuration, ports, environment variables
- `backend/proto/user.proto` - User API endpoints: `/v1/auth/register`, `/v1/auth/login`, `/v1/users/profile`
- `backend/proto/contest.proto` - Contest API: `/v1/contests`, `/v1/contests/{id}/join`
- `backend/proto/prediction.proto` - Prediction API: `/v1/predictions`, `/v1/events`
- `backend/proto/scoring.proto` - Scoring API: `/v1/scores`, `/v1/contests/{contest_id}/leaderboard`
- `backend/proto/sports.proto` - Sports API: `/v1/sports`, `/v1/leagues`, `/v1/teams`, `/v1/matches`
- `backend/api-gateway/internal/gateway/gateway.go` - Gateway routing, error handling, health check
- `tests/prediction-service/integration_test.go` - Existing integration test pattern with `//go:build integration`
- `tests/contest-service/integration_test.go` - SQLite in-memory test pattern
- `Makefile` - Existing test commands and Docker targets

### New Files to Create

- `tests/e2e/go.mod` - E2E test module
- `tests/e2e/main_test.go` - Test setup, Docker health checks, cleanup
- `tests/e2e/auth_test.go` - User registration and login tests
- `tests/e2e/contest_test.go` - Contest CRUD and participant tests
- `tests/e2e/prediction_test.go` - Prediction workflow tests
- `tests/e2e/scoring_test.go` - Scoring and leaderboard tests
- `tests/e2e/sports_test.go` - Sports data management tests
- `tests/e2e/helpers.go` - HTTP client helpers, test utilities
- `scripts/e2e-test.sh` - Script to run E2E tests with Docker

### Relevant Documentation

- [Go Testing Package](https://pkg.go.dev/testing)
  - TestMain for setup/teardown
  - Why: Standard Go testing patterns
- [Docker Compose CLI](https://docs.docker.com/compose/reference/)
  - `--profile` flag for service groups
  - Why: Starting services for tests

### Patterns to Follow

**Test File Naming**: `*_test.go` with `//go:build e2e` build tag

**HTTP Client Pattern** (from existing codebase):
```go
func makeRequest(method, url string, body interface{}, token string) (*http.Response, error) {
    var reqBody io.Reader
    if body != nil {
        jsonBody, _ := json.Marshal(body)
        reqBody = bytes.NewBuffer(jsonBody)
    }
    req, _ := http.NewRequest(method, url, reqBody)
    req.Header.Set("Content-Type", "application/json")
    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }
    return http.DefaultClient.Do(req)
}
```

**Test Isolation Pattern**:
```go
func TestMain(m *testing.M) {
    // Setup: wait for services
    // Run tests
    code := m.Run()
    // Teardown: cleanup test data
    os.Exit(code)
}
```

**Error Response Pattern** (from gateway.go):
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

**Tasks:**
- Create E2E test module with dependencies
- Implement test helpers for HTTP requests
- Create Docker startup/health check utilities
- Set up TestMain with service readiness checks

### Phase 2: Core Test Implementation

**Tasks:**
- Implement authentication flow tests (register, login, profile)
- Implement contest management tests (CRUD, join/leave)
- Implement sports data tests (sports, leagues, teams, matches)
- Implement prediction workflow tests (submit, update, delete)
- Implement scoring and leaderboard tests

### Phase 3: Integration

**Tasks:**
- Create shell script for running E2E tests
- Add Makefile targets for E2E testing
- Document test execution process

### Phase 4: Validation

**Tasks:**
- Run full E2E suite against Docker services
- Verify all tests pass
- Document any discovered issues

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `tests/e2e/go.mod`

- **IMPLEMENT**: Go module for E2E tests
- **PATTERN**: Mirror `tests/notification-service/go.mod` structure
- **CONTENT**:
```go
module github.com/sports-prediction-contests/e2e

go 1.21
```
- **VALIDATE**: `cd tests/e2e && go mod tidy`

### Task 2: CREATE `tests/e2e/helpers.go`

- **IMPLEMENT**: HTTP client helpers and test utilities
- **IMPORTS**: `net/http`, `encoding/json`, `bytes`, `io`, `fmt`, `time`
- **FUNCTIONS**:
  - `BaseURL() string` - returns `http://localhost:8080`
  - `makeRequest(method, path string, body interface{}, token string) (*http.Response, error)`
  - `parseResponse[T any](resp *http.Response) (T, error)`
  - `waitForService(url string, timeout time.Duration) error`
  - `generateTestEmail() string` - unique email for test isolation
- **GOTCHA**: Use `localhost:8080` not container names (tests run on host)
- **VALIDATE**: `cd tests/e2e && go build ./...`

### Task 3: CREATE `tests/e2e/types.go`

- **IMPLEMENT**: Response types matching proto definitions
- **TYPES**:
  - `Response` with `Success bool`, `Message string`
  - `User` with `ID`, `Email`, `Name`, `CreatedAt`
  - `AuthResponse` with `Response`, `User`, `Token`
  - `Contest` with all fields from proto
  - `ContestResponse`, `ContestsResponse`
  - `Prediction`, `PredictionResponse`
  - `Event`, `EventResponse`, `EventsResponse`
  - `LeaderboardEntry`, `Leaderboard`, `LeaderboardResponse`
  - `Sport`, `League`, `Team`, `Match` and their responses
- **PATTERN**: Match JSON field names from proto (snake_case)
- **VALIDATE**: `cd tests/e2e && go build ./...`

### Task 4: CREATE `tests/e2e/main_test.go`

- **IMPLEMENT**: TestMain with service health checks
- **BUILD TAG**: `//go:build e2e`
- **IMPORTS**: `testing`, `os`, `time`, `log`
- **FUNCTIONS**:
  - `TestMain(m *testing.M)` - wait for services, run tests, exit
  - `waitForServices()` - check `/health` endpoint with retries
- **TIMEOUT**: 60 seconds for services to be ready
- **HEALTH ENDPOINT**: `GET http://localhost:8080/health`
- **GOTCHA**: Services need time to start - use exponential backoff
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -c`

### Task 5: CREATE `tests/e2e/auth_test.go`

- **IMPLEMENT**: Authentication flow tests
- **BUILD TAG**: `//go:build e2e`
- **TEST CASES**:
  - `TestRegisterUser` - POST `/v1/auth/register` with email, password, name
  - `TestRegisterDuplicateEmail` - expect 409 Conflict
  - `TestLoginUser` - POST `/v1/auth/login` with credentials
  - `TestLoginInvalidCredentials` - expect 401 Unauthorized
  - `TestGetProfile` - GET `/v1/users/profile` with JWT token
  - `TestGetProfileUnauthorized` - expect 401 without token
- **PATTERN**: Each test uses unique email via `generateTestEmail()`
- **ASSERTIONS**: Check status codes, response body, token presence
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -run TestRegister -v` (requires services)

### Task 6: CREATE `tests/e2e/sports_test.go`

- **IMPLEMENT**: Sports data management tests
- **BUILD TAG**: `//go:build e2e`
- **TEST CASES**:
  - `TestCreateSport` - POST `/v1/sports` (requires auth)
  - `TestListSports` - GET `/v1/sports`
  - `TestCreateLeague` - POST `/v1/leagues` with sport_id
  - `TestCreateTeam` - POST `/v1/teams` with sport_id
  - `TestCreateMatch` - POST `/v1/matches` with league_id, team_ids
  - `TestListMatches` - GET `/v1/matches` with filters
- **DEPENDENCIES**: Tests must run in order (sport → league → team → match)
- **PATTERN**: Use `t.Run()` for subtests to ensure order
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -run TestSports -v`

### Task 7: CREATE `tests/e2e/contest_test.go`

- **IMPLEMENT**: Contest management tests
- **BUILD TAG**: `//go:build e2e`
- **TEST CASES**:
  - `TestCreateContest` - POST `/v1/contests` with all fields
  - `TestGetContest` - GET `/v1/contests/{id}`
  - `TestListContests` - GET `/v1/contests` with pagination
  - `TestUpdateContest` - PUT `/v1/contests/{id}`
  - `TestJoinContest` - POST `/v1/contests/{id}/join`
  - `TestLeaveContest` - POST `/v1/contests/{id}/leave`
  - `TestListParticipants` - GET `/v1/contests/{id}/participants`
  - `TestDeleteContest` - DELETE `/v1/contests/{id}`
- **SETUP**: Register user and get token before tests
- **DATES**: Use future dates for start_date/end_date (RFC3339 format)
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -run TestContest -v`

### Task 8: CREATE `tests/e2e/prediction_test.go`

- **IMPLEMENT**: Prediction workflow tests
- **BUILD TAG**: `//go:build e2e`
- **TEST CASES**:
  - `TestCreateEvent` - POST `/v1/events`
  - `TestListEvents` - GET `/v1/events`
  - `TestSubmitPrediction` - POST `/v1/predictions`
  - `TestGetPrediction` - GET `/v1/predictions/{id}`
  - `TestGetUserPredictions` - GET `/v1/predictions/contest/{contest_id}`
  - `TestUpdatePrediction` - PUT `/v1/predictions/{id}`
  - `TestDeletePrediction` - DELETE `/v1/predictions/{id}`
  - `TestGetPotentialCoefficient` - GET `/v1/events/{id}/coefficient`
- **SETUP**: Create contest and event before prediction tests
- **PREDICTION_DATA**: JSON string like `{"outcome": "home_win", "home_score": 2, "away_score": 1}`
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -run TestPrediction -v`

### Task 9: CREATE `tests/e2e/scoring_test.go`

- **IMPLEMENT**: Scoring and leaderboard tests
- **BUILD TAG**: `//go:build e2e`
- **TEST CASES**:
  - `TestCreateScore` - POST `/v1/scores`
  - `TestGetLeaderboard` - GET `/v1/contests/{id}/leaderboard`
  - `TestGetUserRank` - GET `/v1/contests/{id}/users/{user_id}/rank`
  - `TestGetUserStreak` - GET `/v1/contests/{id}/users/{user_id}/streak`
  - `TestGetUserAnalytics` - GET `/v1/users/{user_id}/analytics`
- **SETUP**: Create contest, user, prediction, then score
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -run TestScoring -v`

### Task 10: CREATE `tests/e2e/workflow_test.go`

- **IMPLEMENT**: Complete user journey test
- **BUILD TAG**: `//go:build e2e`
- **TEST CASE**: `TestCompleteUserJourney`
  1. Register new user → get token
  2. Create sport, league, teams
  3. Create match (scheduled for future)
  4. Create contest
  5. Join contest
  6. Submit prediction for match
  7. Verify prediction in user's predictions
  8. Check leaderboard shows user
  9. Get user analytics
- **PURPOSE**: Validates entire platform workflow in single test
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -run TestCompleteUserJourney -v`

### Task 11: CREATE `scripts/e2e-test.sh`

- **IMPLEMENT**: Shell script to run E2E tests
- **CONTENT**:
```bash
#!/bin/bash
set -e

echo "Starting services..."
docker-compose --profile services up -d

echo "Waiting for services to be healthy..."
sleep 10

# Wait for API Gateway
until curl -s http://localhost:8080/health > /dev/null; do
    echo "Waiting for API Gateway..."
    sleep 2
done

echo "Running E2E tests..."
cd tests/e2e
go test -tags=e2e -v -timeout 5m ./...
TEST_EXIT=$?

echo "Stopping services..."
docker-compose --profile services down

exit $TEST_EXIT
```
- **PERMISSIONS**: Make executable
- **VALIDATE**: `chmod +x scripts/e2e-test.sh && ./scripts/e2e-test.sh`

### Task 12: UPDATE `Makefile`

- **ADD**: E2E test targets
- **CONTENT** to append:
```makefile
e2e-test: ## Run end-to-end tests with Docker services
	@echo "Running E2E tests..."
	@./scripts/e2e-test.sh

e2e-test-only: ## Run E2E tests (assumes services are running)
	@echo "Running E2E tests against running services..."
	@cd tests/e2e && go test -tags=e2e -v -timeout 5m ./...
```
- **VALIDATE**: `make help | grep e2e`

---

## TESTING STRATEGY

### Unit Tests
Not applicable - this feature IS the test infrastructure.

### Integration Tests
The E2E tests themselves are integration tests that validate:
- HTTP API Gateway routing
- gRPC service-to-service communication
- Database operations across services
- JWT authentication flow

### Edge Cases

1. **Service Unavailable**: Test graceful handling when a service is down
2. **Invalid JWT**: Test 401 responses for expired/invalid tokens
3. **Duplicate Resources**: Test 409 Conflict for duplicate emails, contest joins
4. **Not Found**: Test 404 for non-existent resources
5. **Validation Errors**: Test 400 for invalid request bodies
6. **Pagination**: Test limit/offset parameters
7. **Concurrent Access**: Test multiple users joining same contest

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
# Build E2E test module
cd tests/e2e && go build ./...

# Compile tests without running
cd tests/e2e && go test -tags=e2e -c
```

### Level 2: Service Health

```bash
# Start infrastructure only
make docker-up

# Start all services
docker-compose --profile services up -d

# Check API Gateway health
curl http://localhost:8080/health
```

### Level 3: E2E Tests

```bash
# Run all E2E tests (with service startup)
make e2e-test

# Run specific test file
cd tests/e2e && go test -tags=e2e -v -run TestAuth ./...

# Run with verbose output
cd tests/e2e && go test -tags=e2e -v -timeout 5m ./...
```

### Level 4: Manual Validation

```bash
# Test registration manually
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Test login manually
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Test authenticated endpoint
curl http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer <token>"
```

---

## ACCEPTANCE CRITERIA

- [ ] E2E test module compiles without errors
- [ ] All E2E tests pass against running Docker services
- [ ] Authentication flow tests cover register, login, profile
- [ ] Contest tests cover full CRUD and participant management
- [ ] Prediction tests cover submit, update, delete workflow
- [ ] Scoring tests verify leaderboard and analytics
- [ ] Complete user journey test validates end-to-end workflow
- [ ] `make e2e-test` runs full suite with Docker orchestration
- [ ] Tests are isolated (each test uses unique data)
- [ ] Tests clean up after themselves (no leftover test data)
- [ ] Test output is clear and diagnostic on failure

---

## COMPLETION CHECKLIST

- [ ] All 12 tasks completed in order
- [ ] `tests/e2e/go.mod` created and dependencies resolved
- [ ] Helper functions implemented and tested
- [ ] Type definitions match proto responses
- [ ] TestMain waits for service health
- [ ] Auth tests pass
- [ ] Sports tests pass
- [ ] Contest tests pass
- [ ] Prediction tests pass
- [ ] Scoring tests pass
- [ ] Workflow test passes
- [ ] Shell script works end-to-end
- [ ] Makefile targets added
- [ ] All validation commands executed successfully

---

## NOTES

### Design Decisions

1. **Go-based tests**: Chosen over shell scripts for better assertions, error handling, and IDE support
2. **Build tags**: Using `//go:build e2e` to separate from unit tests
3. **Host networking**: Tests run on host machine, connect to `localhost:8080`
4. **Test isolation**: Each test generates unique emails/names to avoid conflicts
5. **Sequential execution**: Some tests depend on others (sport → league → team → match)

### Known Limitations

1. **Database cleanup**: Tests don't clean up database - relies on fresh Docker volumes
2. **Parallel execution**: Tests should run sequentially due to shared state
3. **Service startup time**: May need to increase timeout on slower machines

### Future Improvements

1. Add database cleanup between test runs
2. Add performance/load testing
3. Add negative test cases for all error conditions
4. Add frontend E2E tests with Playwright
5. Add CI/CD integration for automated E2E testing

### Estimated Time

- **Implementation**: 2-3 hours
- **Validation**: 30 minutes
- **Total**: ~3 hours

### Confidence Score: 8/10

High confidence due to:
- Clear API endpoints from proto files
- Existing test patterns to follow
- Well-documented Docker setup
- Standard Go testing approach

Risks:
- Service startup timing may need tuning
- Some endpoints may have undocumented behavior
- Database state between tests could cause flakiness

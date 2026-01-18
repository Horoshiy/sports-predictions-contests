# End-to-End Testing Guide

Comprehensive guide for running E2E tests on the Sports Prediction Contests platform.

## Overview

The E2E test suite validates complete user workflows across all microservices, ensuring the platform works correctly from end to end.

## Test Architecture

### Test Structure
```
tests/e2e/
├── main_test.go          # Test runner and setup
├── helpers.go            # Test utilities and helpers
├── types.go              # Test data structures
├── auth_test.go          # Authentication workflows
├── contest_test.go       # Contest management tests
├── prediction_test.go    # Prediction submission tests
├── scoring_test.go       # Scoring and leaderboard tests
├── sports_test.go        # Sports data management tests
└── workflow_test.go      # Complete user workflows
```

### Test Environment
- **Isolated Docker environment**: Tests run against containerized services
- **Fresh database**: Each test run starts with a clean database
- **Real service communication**: Tests use actual gRPC/HTTP communication
- **Automated setup/teardown**: Infrastructure managed automatically

## Prerequisites

### Required Software
- **Docker** and **Docker Compose**
- **Go** 1.21+ for running tests
- **curl** for manual API testing
- **jq** for JSON processing (optional)

### Environment Setup
```bash
# Ensure you're in the project root
cd /path/to/sports-prediction-contests

# Verify Docker is running
docker --version
docker-compose --version

# Check Go installation
go version
```

## Running E2E Tests

### Automated Test Execution

#### Full E2E Test Suite
```bash
# Run complete E2E test suite with infrastructure setup
make e2e-test

# This command will:
# 1. Start PostgreSQL and Redis
# 2. Wait for database readiness
# 3. Start all microservices
# 4. Wait for service health checks
# 5. Run E2E tests
# 6. Clean up all services
```

#### E2E Tests Only (Services Running)
```bash
# If services are already running
make e2e-test-only

# Or run directly
cd tests/e2e
go test -tags=e2e -v -timeout 5m ./...
```

### Manual Test Execution

#### Step 1: Start Infrastructure
```bash
# Start PostgreSQL and Redis
docker-compose up -d postgres redis

# Wait for database readiness
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction
```

#### Step 2: Start All Services
```bash
# Start all microservices
docker-compose --profile services up -d

# Wait for services to be ready (15-30 seconds)
sleep 15

# Verify API Gateway health
curl http://localhost:8080/health
```

#### Step 3: Run Tests
```bash
# Run specific test files
cd tests/e2e

# Authentication tests
go test -tags=e2e -v -run TestAuth

# Contest management tests
go test -tags=e2e -v -run TestContest

# Complete workflow tests
go test -tags=e2e -v -run TestWorkflow

# All tests with verbose output
go test -tags=e2e -v -timeout 5m ./...
```

#### Step 4: Cleanup
```bash
# Stop all services
docker-compose --profile services down

# Stop infrastructure
docker-compose down
```

## Test Scenarios

### Authentication Workflow Tests

#### User Registration and Login
```go
func TestUserRegistrationAndLogin(t *testing.T) {
    // Test user registration
    user := registerTestUser(t, "testuser", "test@example.com", "password123")
    
    // Test user login
    token := loginUser(t, "test@example.com", "password123")
    
    // Verify token validity
    profile := getUserProfile(t, token)
    assert.Equal(t, "testuser", profile.Username)
}
```

**Manual Test:**
```bash
# Register user
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "password123",
    "full_name": "Test User"
  }'

# Login user
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Save token for subsequent requests
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Contest Management Tests

#### Contest Creation and Participation
```go
func TestContestCreationAndParticipation(t *testing.T) {
    // Create user and get token
    token := setupTestUser(t)
    
    // Create contest
    contest := createTestContest(t, token, "Test Contest")
    
    // Join contest
    joinContest(t, token, contest.ID)
    
    // Verify participation
    participants := getContestParticipants(t, token, contest.ID)
    assert.Len(t, participants, 1)
}
```

**Manual Test:**
```bash
# Create contest
curl -X POST http://localhost:8080/v1/contests \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Contest",
    "description": "E2E test contest",
    "sport_type": "football",
    "rules": "{\"scoring_system\": \"standard\"}",
    "start_date": "2026-01-20T00:00:00Z",
    "end_date": "2026-05-20T23:59:59Z",
    "max_participants": 100
  }'

# Join contest (save contest_id from previous response)
curl -X POST http://localhost:8080/v1/contests/1/join \
  -H "Authorization: Bearer $JWT_TOKEN"

# List participants
curl -H "Authorization: Bearer $JWT_TOKEN" \
     http://localhost:8080/v1/contests/1/participants
```

### Prediction Workflow Tests

#### Prediction Submission and Scoring
```go
func TestPredictionWorkflow(t *testing.T) {
    // Setup: user, contest, event
    token := setupTestUser(t)
    contest := createTestContest(t, token, "Prediction Test")
    event := createTestEvent(t, token, "Test Match")
    
    // Submit prediction
    prediction := submitPrediction(t, token, contest.ID, event.ID, "home_win")
    
    // Update event result
    updateEventResult(t, token, event.ID, "home_win")
    
    // Calculate score
    score := calculateScore(t, token, prediction.ID)
    assert.Greater(t, score.Points, 0)
}
```

**Manual Test:**
```bash
# Create event
curl -X POST http://localhost:8080/v1/events \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Match",
    "description": "E2E test match",
    "sport_type": "football",
    "start_time": "2026-01-25T15:00:00Z",
    "end_time": "2026-01-25T17:00:00Z"
  }'

# Submit prediction
curl -X POST http://localhost:8080/v1/predictions \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "contest_id": 1,
    "event_id": 1,
    "prediction_type": "match_outcome",
    "prediction_value": "home_win"
  }'
```

### Leaderboard and Scoring Tests

#### Leaderboard Updates
```go
func TestLeaderboardUpdates(t *testing.T) {
    // Setup multiple users and predictions
    users := setupMultipleUsers(t, 3)
    contest := createTestContest(t, users[0].Token, "Leaderboard Test")
    
    // Submit predictions with different outcomes
    submitMultiplePredictions(t, users, contest.ID)
    
    // Get leaderboard
    leaderboard := getLeaderboard(t, users[0].Token, contest.ID)
    assert.Len(t, leaderboard, 3)
    
    // Verify ranking order
    assert.Greater(t, leaderboard[0].TotalScore, leaderboard[1].TotalScore)
}
```

### Team Tournament Tests

#### Team Creation and Management
```go
func TestTeamTournaments(t *testing.T) {
    // Create team leader
    leaderToken := setupTestUser(t)
    team := createTestTeam(t, leaderToken, "Test Team")
    
    // Create team member
    memberToken := setupTestUser(t)
    joinTeamByInviteCode(t, memberToken, team.InviteCode)
    
    // Create team contest
    contest := createTestContest(t, leaderToken, "Team Contest")
    joinContestAsTeam(t, leaderToken, team.ID, contest.ID)
    
    // Verify team participation
    teamLeaderboard := getTeamLeaderboard(t, leaderToken, contest.ID)
    assert.Len(t, teamLeaderboard, 1)
}
```

## Test Data Management

### Test User Creation
```go
type TestUser struct {
    ID       uint32
    Username string
    Email    string
    Token    string
}

func setupTestUser(t *testing.T) *TestUser {
    username := fmt.Sprintf("testuser_%d", time.Now().Unix())
    email := fmt.Sprintf("%s@example.com", username)
    
    // Register user
    user := registerTestUser(t, username, email, "password123")
    
    // Login and get token
    token := loginUser(t, email, "password123")
    
    return &TestUser{
        ID:       user.ID,
        Username: username,
        Email:    email,
        Token:    token,
    }
}
```

### Test Data Cleanup
```go
func cleanupTestData(t *testing.T) {
    // Tests run against isolated Docker environment
    // Cleanup happens automatically when containers are stopped
    // No manual cleanup required
}
```

## Test Configuration

### Environment Variables
```bash
# Test configuration in tests/e2e/.env
BASE_URL=http://localhost:8080
DATABASE_URL=postgres://sports_user:sports_password@localhost:5432/sports_prediction?sslmode=disable
REDIS_URL=redis://localhost:6379
TEST_TIMEOUT=5m
LOG_LEVEL=debug
```

### Test Flags
```bash
# Run with verbose output
go test -tags=e2e -v

# Run specific test pattern
go test -tags=e2e -run TestAuth

# Run with timeout
go test -tags=e2e -timeout 10m

# Run with race detection
go test -tags=e2e -race

# Generate coverage report
go test -tags=e2e -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Debugging Failed Tests

### Common Issues and Solutions

#### Services Not Ready
```bash
# Check service status
docker-compose ps

# Check service logs
docker-compose logs api-gateway
docker-compose logs user-service

# Wait longer for services to start
sleep 30
curl http://localhost:8080/health
```

#### Database Connection Issues
```bash
# Check PostgreSQL status
docker-compose logs postgres

# Test database connection
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

#### Port Conflicts
```bash
# Check what's using ports
lsof -i :8080
lsof -i :5432

# Kill conflicting processes
kill -9 <PID>
```

### Test Debugging
```go
func TestWithDebugging(t *testing.T) {
    // Enable debug logging
    log.SetLevel(log.DebugLevel)
    
    // Add debug prints
    t.Logf("Starting test: %s", t.Name())
    
    // Use test helpers with error details
    user := registerTestUserWithDetails(t, "debug_user", "debug@example.com")
    t.Logf("Created user: %+v", user)
}
```

## Performance Testing

### Load Testing with E2E Tests
```bash
# Run tests multiple times to check for race conditions
for i in {1..10}; do
    echo "Run $i"
    go test -tags=e2e -run TestWorkflow
done

# Run tests in parallel
go test -tags=e2e -parallel 4
```

### Memory and Resource Monitoring
```bash
# Monitor Docker resource usage during tests
docker stats

# Check memory usage
docker-compose exec api-gateway ps aux
```

## Continuous Integration

### GitHub Actions E2E Tests
```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Run E2E Tests
        run: make e2e-test
      
      - name: Upload Test Results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: e2e-test-results
          path: tests/e2e/test-results.xml
```

## Test Reporting

### Generate Test Report
```bash
# Run tests with JSON output
go test -tags=e2e -json ./... > test-results.json

# Generate HTML report (requires gotestsum)
gotestsum --junitfile test-results.xml --format testname

# View coverage
go test -tags=e2e -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Troubleshooting Checklist

- [ ] Docker and Docker Compose installed and running
- [ ] All required ports available (8080-8089, 5432, 6379)
- [ ] Services started in correct order (infrastructure first)
- [ ] Database initialized and accessible
- [ ] API Gateway health check passing
- [ ] JWT tokens valid and not expired
- [ ] Test data cleanup between runs
- [ ] Sufficient system resources (memory, disk)
- [ ] Network connectivity between services
- [ ] Environment variables correctly set

---

**E2E tests provide confidence that the entire platform works correctly. Run them regularly during development and always before deployment.**

# Feature: Complete Team Service Implementation

## Feature Description

Implement a fully-featured, standalone Team Service microservice with complete gRPC backend, frontend integration, comprehensive testing, and bilingual documentation. The Team Service enables users to create teams, manage members, join contests as teams, and compete in team-based leaderboards.

Currently, team functionality exists within contest-service but lacks:
- Standalone microservice architecture
- gRPC server registration in contest-service main.go
- Complete E2E testing coverage
- Dedicated API documentation
- README for the service
- Integration tests for team workflows

## User Story

**As a** sports prediction platform user  
**I want to** create and manage teams with other users  
**So that** I can compete in team-based tournaments and collaborate on predictions

## Problem Statement

The platform has team models, repositories, and service logic implemented in contest-service, but the team functionality is not fully integrated:

1. **Backend**: Team service exists but is not registered in the gRPC server
2. **Testing**: Only basic unit tests exist; missing integration and E2E tests
3. **Documentation**: No dedicated README or API documentation for team endpoints
4. **Deployment**: No clear separation or documentation of team service responsibilities

This creates confusion about team service ownership and makes it difficult for developers to understand and extend team functionality.

## Solution Statement

Complete the Team Service implementation by:

1. **Backend Integration**: Register TeamService in contest-service gRPC server with proper migration
2. **Testing**: Add comprehensive integration tests and expand E2E test coverage
3. **Documentation**: Create dedicated README and expand API documentation
4. **Validation**: Ensure all team workflows work end-to-end with proper error handling

The solution maintains team service within contest-service (as designed) but makes it a first-class citizen with complete documentation and testing.

## Feature Metadata

**Feature Type**: Enhancement  
**Estimated Complexity**: Medium  
**Primary Systems Affected**: 
- backend/contest-service (team service registration)
- tests/contest-service (integration tests)
- frontend/tests/e2e (E2E tests)
- docs/en and docs/ru (documentation)

**Dependencies**: 
- gRPC and Protocol Buffers (already configured)
- GORM (already configured)
- Playwright (already configured)
- Existing team models, repositories, and service logic

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Backend - Team Service Implementation:**
- `backend/contest-service/internal/service/team_service.go` (lines 1-460) - Complete team service implementation with all methods
- `backend/contest-service/internal/models/team.go` (lines 1-108) - Team model with validation
- `backend/contest-service/internal/models/team_member.go` (lines 1-79) - TeamMember model
- `backend/contest-service/internal/models/team_contest_entry.go` (lines 1-42) - TeamContestEntry model
- `backend/contest-service/internal/repository/team_repository.go` (lines 1-277) - All team repositories

**Backend - Service Registration Pattern:**
- `backend/user-service/cmd/main.go` (lines 1-75) - Pattern for service registration and gRPC server setup
- `backend/challenge-service/cmd/main.go` (lines 1-86) - Another example of service registration
- `backend/contest-service/cmd/main.go` (lines 1-86) - Current contest service main.go (needs team service registration)

**Backend - Proto and API Gateway:**
- `backend/proto/team.proto` (lines 1-200) - Complete team service proto definition
- `backend/api-gateway/internal/gateway/gateway.go` (lines 80-85) - Team service already registered in gateway
- `backend/api-gateway/internal/config/config.go` (line 40) - Team service endpoint config (points to contest-service:8085)

**Frontend - Team Service:**
- `frontend/src/services/team-service.ts` (lines 1-82) - Complete frontend team service client
- `frontend/src/pages/TeamsPage.tsx` (lines 1-153) - Teams page implementation
- `frontend/src/components/teams/TeamList.tsx` - Team list component
- `frontend/src/components/teams/TeamForm.tsx` - Team creation/edit form
- `frontend/src/components/teams/TeamMembers.tsx` - Team members management
- `frontend/src/hooks/use-teams.ts` - React hooks for team operations

**Testing - Existing Tests:**
- `tests/contest-service/team_test.go` (lines 1-155) - Unit tests for team models
- `frontend/tests/e2e/teams.spec.ts` (lines 1-17) - Basic E2E tests (needs expansion)
- `frontend/tests/helpers/selectors.ts` - Test selectors (check for team selectors)

**Documentation - Existing:**
- `backend/contest-service/README.md` - Contest service README (no team documentation)
- `docs/en/api/services-overview.md` (lines 142-200) - Partial team API documentation
- `.kiro/steering/innovations.md` (lines 112-140) - Team Tournaments innovation description

### New Files to Create

**Documentation:**
- `docs/en/api/team-service.md` - Complete Team Service API documentation
- `docs/ru/api/team-service.md` - Russian translation of Team Service API documentation

**Testing:**
- `tests/contest-service/team_integration_test.go` - Integration tests for team service
- `tests/e2e/team_workflow_test.go` - Backend E2E tests for complete team workflows

### Files to Update

**Backend:**
- `backend/contest-service/cmd/main.go` - Register TeamService in gRPC server
- `backend/contest-service/README.md` - Add team service documentation section

**Frontend Testing:**
- `frontend/tests/e2e/teams.spec.ts` - Expand E2E test coverage

**Documentation:**
- `docs/en/api/services-overview.md` - Expand team service section
- `docs/ru/api/services-overview.md` - Add Russian team service documentation
- `docs/en/README.md` - Update with team service information
- `docs/ru/README.md` - Update with team service information

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

**gRPC and Protocol Buffers:**
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
  - Section: Server implementation
  - Why: Shows how to register multiple services on one gRPC server
- [Protocol Buffers Go Tutorial](https://protobuf.dev/getting-started/gotutorial/)
  - Section: Compiling protocol buffers
  - Why: Understanding proto compilation for team.proto

**Testing:**
- [Go Testing Package](https://pkg.go.dev/testing)
  - Section: Integration testing patterns
  - Why: Writing integration tests for team service
- [Playwright Test](https://playwright.dev/docs/writing-tests)
  - Section: Test fixtures and page objects
  - Why: Expanding E2E tests for team workflows

**GORM:**
- [GORM Associations](https://gorm.io/docs/associations.html)
  - Section: Has Many, Belongs To
  - Why: Understanding team-member relationships

### Patterns to Follow

**Service Registration Pattern (from user-service and challenge-service):**

```go
// In cmd/main.go
func main() {
    // ... config and db setup ...
    
    // Initialize repositories
    teamRepo := repository.NewTeamRepository(db)
    memberRepo := repository.NewTeamMemberRepository(db)
    contestEntryRepo := repository.NewTeamContestEntryRepository(db)
    
    // Initialize services
    contestService := service.NewContestService(contestRepo, participantRepo)
    teamService := service.NewTeamService(teamRepo, memberRepo, contestEntryRepo)
    
    // Create gRPC server
    server := grpc.NewServer(
        grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
    )
    
    // Register BOTH services
    contestpb.RegisterContestServiceServer(server, contestService)
    teampb.RegisterTeamServiceServer(server, teamService)
    
    // ... start server ...
}
```

**Database Migration Pattern:**

```go
// Auto-migrate all models including team models
if err := db.AutoMigrate(
    &models.Contest{}, 
    &models.Participant{},
    &models.Team{},
    &models.TeamMember{},
    &models.TeamContestEntry{},
); err != nil {
    log.Fatalf("Failed to migrate database: %v", err)
}
```

**Integration Test Pattern (from contest_service_test):**

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Migrate all models
    db.AutoMigrate(&models.Team{}, &models.TeamMember{})
    
    return db
}

func TestTeamServiceIntegration(t *testing.T) {
    db := setupTestDB(t)
    // ... test implementation ...
}
```

**E2E Test Pattern (from auth.spec.ts):**

```typescript
test.describe('Team Management', () => {
  test('should create team successfully', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.click(SELECTORS.teams.createButton)
    
    await authenticatedPage.fill(SELECTORS.teams.nameInput, 'Test Team')
    await authenticatedPage.fill(SELECTORS.teams.descriptionInput, 'Test Description')
    await authenticatedPage.click(SELECTORS.teams.submitButton)
    
    await expect(authenticatedPage.locator(SELECTORS.teams.successMessage)).toBeVisible()
  })
})
```

**Error Handling Pattern (from team_service.go):**

```go
func (s *TeamService) CreateTeam(ctx context.Context, name, description string, maxMembers uint) (*TeamResponse, error) {
    userID, ok := auth.GetUserIDFromContext(ctx)
    if !ok {
        return &TeamResponse{
            Response: errorResponse("Authentication required", common.ErrorCode_UNAUTHENTICATED)
        }, nil
    }
    
    // ... business logic ...
    
    if err := s.teamRepo.CreateWithMember(team, member); err != nil {
        log.Printf("[ERROR] Failed to create team: %v", err)
        return &TeamResponse{
            Response: errorResponse(err.Error(), common.ErrorCode_INVALID_ARGUMENT)
        }, nil
    }
    
    return &TeamResponse{
        Response: successResponse("Team created successfully"),
        Team:     s.modelToProto(team),
    }, nil
}
```

**Logging Pattern:**

```go
log.Printf("[INFO] Team created: %d by user %d", team.ID, userID)
log.Printf("[ERROR] Failed to create team: %v", err)
```

---

## IMPLEMENTATION PLAN

### Phase 1: Backend Integration (30 minutes)

Register TeamService in contest-service gRPC server and ensure proper database migration.

**Tasks:**
1. Update contest-service main.go to register TeamService
2. Add team models to database migration
3. Verify gRPC server starts successfully
4. Test health check endpoint

### Phase 2: Integration Testing (45 minutes)

Create comprehensive integration tests for team service covering all workflows.

**Tasks:**
1. Create team_integration_test.go with test database setup
2. Test team CRUD operations
3. Test member management workflows
4. Test contest entry workflows
5. Test error cases and edge conditions

### Phase 3: E2E Testing Expansion (30 minutes)

Expand frontend E2E tests to cover complete team workflows.

**Tasks:**
1. Add team creation test
2. Add team member management test
3. Add team join via invite code test
4. Add team leaderboard test
5. Verify all tests pass in all browsers

### Phase 4: Documentation (45 minutes)

Create comprehensive bilingual documentation for Team Service.

**Tasks:**
1. Create dedicated Team Service API documentation (EN)
2. Translate to Russian
3. Update contest-service README
4. Update main documentation index
5. Add usage examples and code snippets

---

## STEP-BY-STEP TASKS

### Task 1: UPDATE backend/contest-service/cmd/main.go

**IMPLEMENT**: Register TeamService in gRPC server alongside ContestService

**PATTERN**: Mirror user-service/cmd/main.go (lines 35-55) and challenge-service/cmd/main.go (lines 35-50)

**IMPORTS**: Add `teampb "github.com/sports-prediction-contests/shared/proto/team"`

**CHANGES**:
```go
// After line 33 (after db.AutoMigrate for Contest and Participant)
if err := db.AutoMigrate(
    &models.Contest{}, 
    &models.Participant{},
    &models.Team{},
    &models.TeamMember{},
    &models.TeamContestEntry{},
); err != nil {
    log.Fatalf("Failed to migrate database: %v", err)
}

// After line 37 (after participantRepo initialization)
teamRepo := repository.NewTeamRepository(db)
memberRepo := repository.NewTeamMemberRepository(db)
contestEntryRepo := repository.NewTeamContestEntryRepository(db)

// After line 40 (after contestService initialization)
teamService := service.NewTeamService(teamRepo, memberRepo, contestEntryRepo)

// After line 49 (after RegisterContestServiceServer)
teampb.RegisterTeamServiceServer(server, teamService)
```

**GOTCHA**: Ensure team models are migrated BEFORE creating repositories

**VALIDATE**: `cd backend/contest-service && go build ./cmd/main.go`

### Task 2: UPDATE backend/contest-service/README.md

**IMPLEMENT**: Add Team Service section to README

**PATTERN**: Mirror contest-service README structure (lines 1-50)

**ADD**: After "Participant Operations" section (around line 30), add:

```markdown
### Team Operations

- `CreateTeam` - Create a new team
- `UpdateTeam` - Update team details (captain only)
- `GetTeam` - Retrieve team by ID
- `DeleteTeam` - Delete team (captain only)
- `ListTeams` - List teams with pagination and filters
- `JoinTeam` - Join team using invite code
- `LeaveTeam` - Leave a team
- `RemoveMember` - Remove team member (captain only)
- `ListMembers` - List team members with pagination
- `RegenerateInviteCode` - Generate new invite code (captain only)
- `JoinContestAsTeam` - Join contest as a team (captain only)
- `LeaveContestAsTeam` - Leave contest as a team (captain only)
- `GetTeamLeaderboard` - Get team rankings for a contest

### Team Management

Teams allow users to collaborate in contests. Each team has:
- **Captain**: User who created the team, has full management rights
- **Members**: Users who joined via invite code
- **Invite Code**: Unique 8-character code for joining
- **Max Members**: Configurable limit (2-50, default 10)

#### Team Workflow

1. User creates team (becomes captain)
2. Captain shares invite code with others
3. Users join team using invite code
4. Captain joins contests on behalf of team
5. Team members' predictions contribute to team score
6. Team appears in contest leaderboard

## Usage Examples

### Create Team

\`\`\`bash
grpcurl -plaintext -d '{
  "name": "Dream Team",
  "description": "Best predictors united",
  "max_members": 10
}' localhost:8085 team.TeamService/CreateTeam
\`\`\`

### Join Team

\`\`\`bash
grpcurl -plaintext -d '{
  "invite_code": "A1B2C3D4"
}' localhost:8085 team.TeamService/JoinTeam
\`\`\`

### List Team Members

\`\`\`bash
grpcurl -plaintext -d '{
  "team_id": 1,
  "pagination": {"page": 1, "limit": 10}
}' localhost:8085 team.TeamService/ListMembers
\`\`\`
```

**VALIDATE**: `cat backend/contest-service/README.md | grep -A 5 "Team Operations"`


### Task 3: CREATE tests/contest-service/team_integration_test.go

**CREATE**: Comprehensive integration tests for team service

**PATTERN**: Mirror tests/contest-service/integration_test.go structure

**IMPORTS**:
```go
package contest_service_test

import (
    "context"
    "testing"
    "time"

    "github.com/sports-prediction-contests/contest-service/internal/models"
    "github.com/sports-prediction-contests/contest-service/internal/repository"
    "github.com/sports-prediction-contests/contest-service/internal/service"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)
```

**IMPLEMENT**: Test cases for:
1. `TestTeamServiceCreateTeam` - Create team with captain
2. `TestTeamServiceJoinTeam` - Join team via invite code
3. `TestTeamServiceLeaveTeam` - Leave team (non-captain)
4. `TestTeamServiceRemoveMember` - Captain removes member
5. `TestTeamServiceRegenerateInviteCode` - Captain regenerates code
6. `TestTeamServiceJoinContest` - Team joins contest
7. `TestTeamServiceLeaderboard` - Get team leaderboard
8. `TestTeamServicePermissions` - Verify captain-only operations
9. `TestTeamServiceValidation` - Test validation errors
10. `TestTeamServiceConcurrency` - Test concurrent operations

**GOTCHA**: Use in-memory SQLite for tests, not production database

**VALIDATE**: `cd tests/contest-service && go test -v -run TestTeamService`

### Task 4: CREATE tests/e2e/team_workflow_test.go

**CREATE**: Backend E2E tests for complete team workflows

**PATTERN**: Mirror tests/e2e/contest_test.go structure (lines 1-171)

**IMPORTS**:
```go
package e2e_test

import (
    "encoding/json"
    "testing"
)
```

**IMPLEMENT**: Test complete workflows:
1. `TestTeamCreationWorkflow` - Create team, verify in database
2. `TestTeamJoinWorkflow` - Create team, join with second user
3. `TestTeamContestWorkflow` - Create team, join contest, verify leaderboard
4. `TestTeamMemberManagement` - Add/remove members, verify counts
5. `TestTeamPermissions` - Verify non-captain cannot perform captain actions

**GOTCHA**: Ensure services are running before tests (check with waitForService)

**VALIDATE**: `cd tests/e2e && go test -tags=e2e -v -run TestTeam`

### Task 5: UPDATE frontend/tests/e2e/teams.spec.ts

**UPDATE**: Expand E2E test coverage for team workflows

**PATTERN**: Mirror frontend/tests/e2e/auth.spec.ts structure (lines 1-100)

**ADD**: After existing tests (line 17), add:

```typescript
test('should create team successfully', async ({ authenticatedPage }) => {
  await authenticatedPage.goto('/teams')
  await authenticatedPage.click(SELECTORS.teams.createButton)
  
  await fillAntdInput(authenticatedPage, SELECTORS.teams.nameInput, 'E2E Test Team')
  await fillAntdInput(authenticatedPage, SELECTORS.teams.descriptionInput, 'Created by E2E test')
  await fillAntdInput(authenticatedPage, SELECTORS.teams.maxMembersInput, '10')
  
  await clickAntdButton(authenticatedPage, SELECTORS.teams.submitButton)
  await waitForAntdNotification(authenticatedPage, 'success')
  
  // Verify team appears in list
  const teamCard = authenticatedPage.locator(SELECTORS.teams.teamCard).filter({ hasText: 'E2E Test Team' })
  await expect(teamCard).toBeVisible()
})

test('should join team with invite code', async ({ authenticatedPage }) => {
  await authenticatedPage.goto('/teams')
  await authenticatedPage.click(SELECTORS.teams.joinButton)
  
  await fillAntdInput(authenticatedPage, SELECTORS.teams.inviteCodeInput, 'TESTCODE')
  await clickAntdButton(authenticatedPage, SELECTORS.teams.joinSubmitButton)
  
  await waitForAntdNotification(authenticatedPage, 'success')
})

test('should display team members', async ({ authenticatedPage }) => {
  await authenticatedPage.goto('/teams')
  
  // Click first team card
  const firstTeam = authenticatedPage.locator(SELECTORS.teams.teamCard).first()
  await firstTeam.click()
  
  // Verify members list is visible
  await expect(authenticatedPage.locator(SELECTORS.teams.membersList)).toBeVisible()
  
  const memberCount = await authenticatedPage.locator(SELECTORS.teams.memberItem).count()
  expect(memberCount).toBeGreaterThan(0)
})

test('should regenerate invite code', async ({ authenticatedPage }) => {
  await authenticatedPage.goto('/teams')
  
  // Find team where user is captain
  const captainTeam = authenticatedPage.locator(SELECTORS.teams.teamCard).filter({ hasText: 'Captain' }).first()
  await captainTeam.click()
  
  // Click regenerate button
  await authenticatedPage.click(SELECTORS.teams.regenerateButton)
  
  // Confirm in modal
  await clickAntdButton(authenticatedPage, SELECTORS.common.confirmButton)
  await waitForAntdNotification(authenticatedPage, 'success')
})
```

**IMPORTS**: Add helper functions from `../helpers/test-utils`

**GOTCHA**: Use Ant Design-specific helpers (fillAntdInput, clickAntdButton) not standard Playwright methods

**VALIDATE**: `cd frontend && npm run test:e2e -- teams.spec.ts`

### Task 6: CREATE docs/en/api/team-service.md

**CREATE**: Complete English API documentation for Team Service

**PATTERN**: Mirror docs/en/api/services-overview.md structure

**IMPLEMENT**: Full API reference with:
- Service overview
- All endpoints with request/response examples
- Authentication requirements
- Error codes
- Usage examples
- Best practices

**CONTENT**:
```markdown
# Team Service API

## Overview

The Team Service enables users to create and manage teams for collaborative prediction contests. Teams have a captain (creator) who manages membership and contest participation.

**Base URL**: `/v1/teams`  
**Port**: 8085 (via Contest Service)  
**Authentication**: Required for all endpoints except health check

## Endpoints

### Create Team

Create a new team. The creator becomes the team captain.

**Endpoint**: `POST /v1/teams`  
**Authentication**: Required

**Request Body**:
```json
{
  "name": "Dream Team",
  "description": "Best predictors united",
  "max_members": 10
}
```

**Response**:
```json
{
  "response": {
    "success": true,
    "message": "Team created successfully",
    "code": 0,
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "team": {
    "id": 1,
    "name": "Dream Team",
    "description": "Best predictors united",
    "invite_code": "A1B2C3D4",
    "captain_id": 123,
    "max_members": 10,
    "current_members": 1,
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Validation**:
- `name`: Required, 1-100 characters
- `description`: Optional, max 500 characters
- `max_members`: Optional, 2-50 (default: 10)

**Error Codes**:
- `400`: Invalid input (name too long, invalid max_members)
- `401`: Authentication required
- `500`: Internal server error

---

### Update Team

Update team details. Only the captain can update the team.

**Endpoint**: `PUT /v1/teams/{id}`  
**Authentication**: Required (Captain only)

**Request Body**:
```json
{
  "id": 1,
  "name": "Updated Team Name",
  "description": "Updated description",
  "max_members": 15
}
```

**Response**: Same as Create Team

**Error Codes**:
- `400`: Invalid input
- `401`: Authentication required
- `403`: Permission denied (not captain)
- `404`: Team not found

---

### Get Team

Retrieve team details by ID.

**Endpoint**: `GET /v1/teams/{id}`  
**Authentication**: Optional

**Response**: Same as Create Team

---

### Delete Team

Delete a team. Only the captain can delete the team.

**Endpoint**: `DELETE /v1/teams/{id}`  
**Authentication**: Required (Captain only)

**Response**:
```json
{
  "success": true,
  "message": "Team deleted successfully",
  "code": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error Codes**:
- `401`: Authentication required
- `403`: Permission denied (not captain)
- `404`: Team not found

---

### List Teams

List teams with pagination and filters.

**Endpoint**: `GET /v1/teams`  
**Authentication**: Optional

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)
- `my_teams_only`: Filter to user's teams (default: false)

**Response**:
```json
{
  "response": {
    "success": true,
    "message": "Teams retrieved",
    "code": 0,
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "teams": [
    {
      "id": 1,
      "name": "Dream Team",
      "description": "Best predictors",
      "invite_code": "A1B2C3D4",
      "captain_id": 123,
      "max_members": 10,
      "current_members": 5,
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 50,
    "total_pages": 3
  }
}
```

---

### Join Team

Join a team using an invite code.

**Endpoint**: `POST /v1/teams/join`  
**Authentication**: Required

**Request Body**:
```json
{
  "invite_code": "A1B2C3D4"
}
```

**Response**:
```json
{
  "response": {
    "success": true,
    "message": "Joined team successfully",
    "code": 0,
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "member": {
    "id": 10,
    "team_id": 1,
    "user_id": 456,
    "user_name": "John Doe",
    "role": "member",
    "status": "active",
    "joined_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Codes**:
- `400`: Invalid invite code, team full, already a member
- `401`: Authentication required
- `404`: Team not found

---

### Leave Team

Leave a team. Captains cannot leave their own team.

**Endpoint**: `POST /v1/teams/{team_id}/leave`  
**Authentication**: Required

**Response**:
```json
{
  "success": true,
  "message": "Left team successfully",
  "code": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error Codes**:
- `400`: Captain cannot leave team
- `401`: Authentication required
- `404`: Team or membership not found

---

### Remove Member

Remove a member from the team. Only the captain can remove members.

**Endpoint**: `DELETE /v1/teams/{team_id}/members/{user_id}`  
**Authentication**: Required (Captain only)

**Response**:
```json
{
  "success": true,
  "message": "Member removed successfully",
  "code": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error Codes**:
- `400`: Cannot remove captain
- `401`: Authentication required
- `403`: Permission denied (not captain)
- `404`: Team or member not found

---

### List Members

List team members with pagination.

**Endpoint**: `GET /v1/teams/{team_id}/members`  
**Authentication**: Optional

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

**Response**:
```json
{
  "response": {
    "success": true,
    "message": "Members retrieved",
    "code": 0,
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "members": [
    {
      "id": 1,
      "team_id": 1,
      "user_id": 123,
      "user_name": "Captain User",
      "role": "captain",
      "status": "active",
      "joined_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "team_id": 1,
      "user_id": 456,
      "user_name": "Member User",
      "role": "member",
      "status": "active",
      "joined_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

---

### Regenerate Invite Code

Generate a new invite code for the team. Only the captain can regenerate.

**Endpoint**: `POST /v1/teams/{team_id}/regenerate-invite`  
**Authentication**: Required (Captain only)

**Response**:
```json
{
  "response": {
    "success": true,
    "message": "Invite code regenerated",
    "code": 0,
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "invite_code": "X9Y8Z7W6"
}
```

**Error Codes**:
- `401`: Authentication required
- `403`: Permission denied (not captain)
- `404`: Team not found

---

### Join Contest as Team

Join a contest as a team. Only the captain can join contests.

**Endpoint**: `POST /v1/teams/{team_id}/contests/{contest_id}/join`  
**Authentication**: Required (Captain only)

**Response**:
```json
{
  "success": true,
  "message": "Team joined contest successfully",
  "code": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error Codes**:
- `400`: Team already in contest, contest full
- `401`: Authentication required
- `403`: Permission denied (not captain)
- `404`: Team or contest not found

---

### Leave Contest as Team

Leave a contest as a team. Only the captain can leave contests.

**Endpoint**: `POST /v1/teams/{team_id}/contests/{contest_id}/leave`  
**Authentication**: Required (Captain only)

**Response**:
```json
{
  "success": true,
  "message": "Team left contest successfully",
  "code": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### Get Team Leaderboard

Get team rankings for a contest.

**Endpoint**: `GET /v1/contests/{contest_id}/team-leaderboard`  
**Authentication**: Optional

**Query Parameters**:
- `limit`: Number of teams to return (default: 10)

**Response**:
```json
{
  "response": {
    "success": true,
    "message": "Leaderboard retrieved",
    "code": 0,
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "entries": [
    {
      "team_id": 1,
      "team_name": "Dream Team",
      "total_points": 150.5,
      "rank": 1,
      "member_count": 5
    },
    {
      "team_id": 2,
      "team_name": "Prediction Masters",
      "total_points": 145.0,
      "rank": 2,
      "member_count": 8
    }
  ]
}
```

---

### Health Check

Check if the Team Service is healthy.

**Endpoint**: `GET /v1/teams/health`  
**Authentication**: Not required

**Response**:
```json
{
  "success": true,
  "message": "Team service is healthy",
  "code": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Data Models

### Team

| Field | Type | Description |
|-------|------|-------------|
| `id` | uint32 | Unique team identifier |
| `name` | string | Team name (1-100 chars) |
| `description` | string | Team description (max 500 chars) |
| `invite_code` | string | 8-character invite code |
| `captain_id` | uint32 | User ID of team captain |
| `max_members` | uint32 | Maximum team size (2-50) |
| `current_members` | uint32 | Current number of members |
| `is_active` | bool | Whether team is active |
| `created_at` | timestamp | Creation timestamp |
| `updated_at` | timestamp | Last update timestamp |

### TeamMember

| Field | Type | Description |
|-------|------|-------------|
| `id` | uint32 | Unique member identifier |
| `team_id` | uint32 | Team ID |
| `user_id` | uint32 | User ID |
| `user_name` | string | User display name |
| `role` | string | "captain" or "member" |
| `status` | string | "active" or "inactive" |
| `joined_at` | timestamp | Join timestamp |

### TeamLeaderboardEntry

| Field | Type | Description |
|-------|------|-------------|
| `team_id` | uint32 | Team ID |
| `team_name` | string | Team name |
| `total_points` | double | Total team points |
| `rank` | uint32 | Team rank in contest |
| `member_count` | uint32 | Number of team members |

---

## Best Practices

### Team Management

1. **Invite Codes**: Share invite codes securely; regenerate if compromised
2. **Team Size**: Set appropriate max_members based on contest type
3. **Captain Responsibilities**: Captain should coordinate with members before joining/leaving contests
4. **Member Activity**: Inactive members can be removed by captain

### Error Handling

Always check the `response.success` field before processing data:

```typescript
const response = await teamService.createTeam(request)
if (!response.response.success) {
  console.error('Error:', response.response.message)
  return
}
// Process response.team
```

### Pagination

Use pagination for large lists to improve performance:

```typescript
const { teams, pagination } = await teamService.listTeams({
  pagination: { page: 1, limit: 20 },
  myTeamsOnly: true
})
```

---

## Common Error Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 400 | Invalid request (validation error) |
| 401 | Authentication required |
| 403 | Permission denied |
| 404 | Resource not found |
| 409 | Conflict (duplicate, already exists) |
| 500 | Internal server error |

---

## Examples

### Complete Team Workflow

```bash
# 1. Create team
curl -X POST http://localhost:8080/v1/teams \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dream Team",
    "description": "Best predictors",
    "max_members": 10
  }'

# Response includes invite_code: "A1B2C3D4"

# 2. Share invite code with friends

# 3. Friend joins team
curl -X POST http://localhost:8080/v1/teams/join \
  -H "Authorization: Bearer $FRIEND_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"invite_code": "A1B2C3D4"}'

# 4. Captain joins contest
curl -X POST http://localhost:8080/v1/teams/1/contests/5/join \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

# 5. Check team leaderboard
curl http://localhost:8080/v1/contests/5/team-leaderboard?limit=10
```

---

## Related Documentation

- [Contest Service API](contest-service.md)
- [User Service API](user-service.md)
- [Services Overview](services-overview.md)
```

**VALIDATE**: `cat docs/en/api/team-service.md | grep -c "###"`


### Task 7: CREATE docs/ru/api/team-service.md

**CREATE**: Russian translation of Team Service API documentation

**PATTERN**: Mirror docs/en/api/team-service.md structure with Russian translation

**TRANSLATE**: All content from English version to Russian

**KEY TRANSLATIONS**:
- Team Service → Сервис Команд
- Captain → Капитан
- Member → Участник
- Invite Code → Код приглашения
- Leaderboard → Таблица лидеров

**VALIDATE**: `cat docs/ru/api/team-service.md | grep -c "###"`

### Task 8: UPDATE docs/en/api/services-overview.md

**UPDATE**: Expand Team Management section with complete endpoint list

**PATTERN**: Mirror other service sections in the same file

**REPLACE**: Lines 142-200 (current team section) with comprehensive documentation

**ADD**: Complete endpoint list, examples, and cross-references to team-service.md

**VALIDATE**: `grep -A 50 "Team Management" docs/en/api/services-overview.md`

### Task 9: UPDATE docs/ru/api/services-overview.md

**UPDATE**: Add Russian Team Service documentation to services overview

**PATTERN**: Mirror English version structure

**TRANSLATE**: Team service section from English to Russian

**VALIDATE**: `grep -A 50 "Управление командами" docs/ru/api/services-overview.md`

### Task 10: UPDATE docs/en/README.md

**UPDATE**: Add Team Service to documentation index

**PATTERN**: Mirror existing service entries

**ADD**: After Contest Service entry (around line 30):

```markdown
- [Team Service](api/team-service.md) - Team management and tournaments
```

**VALIDATE**: `grep "Team Service" docs/en/README.md`

### Task 11: UPDATE docs/ru/README.md

**UPDATE**: Add Team Service to Russian documentation index

**PATTERN**: Mirror English version

**ADD**: Russian translation of Team Service entry

**VALIDATE**: `grep "Сервис Команд" docs/ru/README.md`

---

## TESTING STRATEGY

### Unit Tests (Existing)

**Scope**: Model validation and business logic

**Location**: `tests/contest-service/team_test.go`

**Coverage**: Already implemented
- Team name validation
- Description validation
- Max members validation
- Invite code generation
- Team member role validation
- Team member status validation

**Status**: ✅ Complete (155 lines, 8 test functions)

### Integration Tests (New)

**Scope**: Service layer with database interactions

**Location**: `tests/contest-service/team_integration_test.go`

**Coverage**: Must implement
- Team CRUD operations with database
- Member management workflows
- Contest entry workflows
- Permission checks (captain vs member)
- Concurrent operations
- Error handling and edge cases

**Test Database**: In-memory SQLite for isolation

**Fixtures**: 
- Test users (captain, members)
- Test teams with various configurations
- Test contests for team participation

### Backend E2E Tests (New)

**Scope**: Complete workflows across services

**Location**: `tests/e2e/team_workflow_test.go`

**Coverage**: Must implement
- Full team creation to contest participation workflow
- Multi-user team join workflow
- Team leaderboard calculation
- Permission enforcement across API calls

**Prerequisites**: All services must be running

### Frontend E2E Tests (Expand)

**Scope**: User interface workflows

**Location**: `frontend/tests/e2e/teams.spec.ts`

**Coverage**: Must expand
- Team creation via UI
- Team join via invite code
- Member management UI
- Invite code regeneration
- Team leaderboard display

**Browsers**: Chromium, Firefox, WebKit

**Fixtures**: Authenticated user session

### Edge Cases to Test

1. **Concurrency**: Multiple users joining team simultaneously
2. **Validation**: Invalid invite codes, full teams, duplicate joins
3. **Permissions**: Non-captain attempting captain-only operations
4. **Cascading Deletes**: Team deletion removes members and contest entries
5. **Member Limits**: Joining when team is at max capacity
6. **Captain Constraints**: Captain cannot leave or be removed
7. **Inactive Teams**: Cannot join inactive teams

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Build

**Backend Build**:
```bash
cd backend/contest-service && go build ./cmd/main.go
```

**Expected**: No compilation errors, binary created

**Proto Generation** (if proto changed):
```bash
cd backend && make proto
```

**Expected**: No errors, proto files regenerated

**Frontend Build**:
```bash
cd frontend && npm run build
```

**Expected**: No TypeScript errors, build succeeds

### Level 2: Unit Tests

**Backend Unit Tests**:
```bash
cd tests/contest-service && go test -v -run TestTeam
```

**Expected**: All 8 existing team tests pass

**Backend Integration Tests**:
```bash
cd tests/contest-service && go test -v -run TestTeamService
```

**Expected**: All 10 new integration tests pass

### Level 3: Service Integration

**Start Services**:
```bash
make docker-services
```

**Expected**: All services start successfully

**Check Service Health**:
```bash
curl http://localhost:8080/v1/teams/health
```

**Expected**: `{"success":true,"message":"Team service is healthy"}`

**Test Team Creation**:
```bash
# Get auth token first
TOKEN=$(curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user1@example.com","password":"password123"}' \
  | jq -r '.token')

# Create team
curl -X POST http://localhost:8080/v1/teams \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Team","description":"Test","max_members":10}'
```

**Expected**: Team created successfully with invite code

### Level 4: E2E Tests

**Backend E2E Tests**:
```bash
cd tests/e2e && go test -tags=e2e -v -run TestTeam
```

**Expected**: All 5 team workflow tests pass

**Frontend E2E Tests**:
```bash
cd frontend && npm run test:e2e -- teams.spec.ts
```

**Expected**: All 6 team UI tests pass in all browsers

**Full E2E Suite**:
```bash
make playwright-test
```

**Expected**: All tests pass including new team tests

### Level 5: Manual Validation

**Team Creation Workflow**:
1. Navigate to http://localhost:3000/teams
2. Click "Create Team" button
3. Fill form: Name="Manual Test Team", Description="Test", Max Members=10
4. Submit form
5. Verify team appears in list with invite code

**Team Join Workflow**:
1. Copy invite code from created team
2. Open incognito window, login as different user
3. Navigate to /teams
4. Click "Join Team"
5. Enter invite code
6. Verify joined successfully and appears in "My Teams"

**Team Leaderboard**:
1. Create or join a team
2. Captain joins a contest
3. Make predictions in the contest
4. Navigate to contest page
5. Click "Team Leaderboard" tab
6. Verify team appears with correct points

### Level 6: Documentation Validation

**Check Documentation Links**:
```bash
# Check all markdown files for broken links
cd docs && find . -name "*.md" -exec grep -l "team-service.md" {} \;
```

**Expected**: team-service.md referenced in README files

**Verify API Examples**:
```bash
# Test each curl example from documentation
# (Copy examples from docs/en/api/team-service.md)
```

**Expected**: All examples work as documented

---

## ACCEPTANCE CRITERIA

- [x] TeamService registered in contest-service gRPC server
- [x] Team models migrated in database
- [x] Service starts without errors
- [x] Health check endpoint responds correctly
- [ ] 10 integration tests implemented and passing
- [ ] 5 backend E2E tests implemented and passing
- [ ] 6 frontend E2E tests implemented and passing
- [ ] All existing tests still pass (no regressions)
- [ ] English API documentation complete
- [ ] Russian API documentation complete
- [ ] README updated with team service section
- [ ] Documentation index updated
- [ ] All validation commands pass
- [ ] Manual testing confirms all workflows work
- [ ] No TypeScript errors
- [ ] No Go compilation errors
- [ ] No linting errors

---

## COMPLETION CHECKLIST

### Backend Integration
- [ ] contest-service/cmd/main.go updated with TeamService registration
- [ ] Team models added to database migration
- [ ] Service builds successfully
- [ ] Service starts without errors
- [ ] Health check endpoint works

### Testing
- [ ] team_integration_test.go created with 10 test cases
- [ ] All integration tests pass
- [ ] team_workflow_test.go created with 5 E2E tests
- [ ] All backend E2E tests pass
- [ ] teams.spec.ts expanded with 4 new tests
- [ ] All frontend E2E tests pass in all browsers
- [ ] No test regressions

### Documentation
- [ ] docs/en/api/team-service.md created
- [ ] docs/ru/api/team-service.md created
- [ ] backend/contest-service/README.md updated
- [ ] docs/en/api/services-overview.md updated
- [ ] docs/ru/api/services-overview.md updated
- [ ] docs/en/README.md updated
- [ ] docs/ru/README.md updated
- [ ] All documentation links work

### Validation
- [ ] All Level 1 validation commands pass
- [ ] All Level 2 validation commands pass
- [ ] All Level 3 validation commands pass
- [ ] All Level 4 validation commands pass
- [ ] All Level 5 manual tests pass
- [ ] All Level 6 documentation checks pass

### Quality
- [ ] Code follows project conventions
- [ ] No hardcoded values
- [ ] Proper error handling
- [ ] Logging follows patterns
- [ ] No security issues
- [ ] Performance is acceptable

---

## NOTES

### Design Decisions

**Why Team Service in Contest Service?**
- Teams are tightly coupled with contests (team contest entries)
- Sharing database connection and transaction management
- Reduces inter-service communication overhead
- Follows existing architecture pattern (both services on port 8085)

**Why Not Standalone Team Service?**
- Would require additional service coordination for contest participation
- Would need separate database or shared database with contest-service
- API Gateway already configured to route teams to contest-service:8085
- Current architecture is simpler and more maintainable

### Implementation Trade-offs

**Pros of Current Approach**:
- Single database transaction for team + contest operations
- No network latency between team and contest operations
- Simpler deployment (one service instead of two)
- Easier testing (no service mocking needed)

**Cons of Current Approach**:
- Contest service has more responsibilities
- Slightly larger service codebase
- Team and contest concerns not fully separated

**Decision**: Keep team service in contest-service as designed, but make it a first-class citizen with complete documentation and testing.

### Future Enhancements

**Potential Improvements**:
1. **Team Chat**: Add messaging between team members
2. **Team Statistics**: Aggregate team performance across all contests
3. **Team Achievements**: Badges and rewards for team milestones
4. **Team Invitations**: Email/SMS invitations instead of just codes
5. **Team Roles**: Add more roles beyond captain/member (co-captain, analyst)
6. **Team Templates**: Pre-configured team settings for different contest types

**Migration Path to Standalone Service** (if needed in future):
1. Create new team-service directory
2. Move team models, repositories, and service
3. Update API Gateway to route to new service
4. Add gRPC client in contest-service for team operations
5. Update docker-compose with new service
6. Migrate team data to new service database

### Known Limitations

1. **Captain Transfer**: No mechanism to transfer captain role to another member
2. **Team Merging**: Cannot merge two teams
3. **Team History**: No audit log of team changes
4. **Team Privacy**: All teams are visible in list (no private teams)
5. **Team Search**: No search functionality for teams

These limitations are acceptable for MVP and can be addressed in future iterations.

---

## CONFIDENCE SCORE

**Estimated Confidence for One-Pass Success**: 9/10

**Reasoning**:
- ✅ All code already exists and is tested (models, repositories, service)
- ✅ Proto definitions are complete and generated
- ✅ Frontend integration is complete and working
- ✅ API Gateway routing is configured
- ✅ Clear patterns to follow from other services
- ⚠️ Only missing: service registration, tests, and documentation
- ⚠️ Minor risk: Integration test database setup

**Risk Mitigation**:
- Follow exact patterns from user-service and challenge-service
- Use in-memory SQLite for integration tests (proven pattern)
- Test incrementally after each task
- Validate with curl commands before E2E tests

**Expected Implementation Time**: 2.5-3 hours
- Backend Integration: 30 minutes
- Integration Testing: 45 minutes
- E2E Testing: 30 minutes
- Documentation: 45 minutes
- Validation: 30 minutes

---

## RELATED DOCUMENTATION

### Internal References
- `.kiro/steering/innovations.md` (lines 112-140) - Team Tournaments innovation
- `.agents/plans/team-tournaments-implementation.md` - Original team implementation plan
- `.agents/code-reviews/team-tournaments-final-review.md` - Team implementation review

### External References
- [gRPC Go Basics](https://grpc.io/docs/languages/go/basics/)
- [GORM Associations](https://gorm.io/docs/associations.html)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)

---

**END OF IMPLEMENTATION PLAN**

This plan provides complete context for implementing full Team Service integration with testing and documentation. Execute tasks in order, validate after each task, and ensure all acceptance criteria are met before considering the feature complete.

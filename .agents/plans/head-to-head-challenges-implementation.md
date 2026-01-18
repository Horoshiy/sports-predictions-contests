# Feature: Head-to-Head Challenges

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Head-to-Head Challenges enable direct duels between two users on specific matches or series of matches. Users can challenge friends or accept open challenges, creating personalized competition that increases emotional engagement and platform stickiness. The feature integrates seamlessly with existing contest infrastructure and notification systems.

## User Story

As a sports prediction enthusiast
I want to challenge specific users to head-to-head prediction duels
So that I can compete directly with friends and prove my prediction skills in personalized matchups

## Problem Statement

Current contest system only supports large group competitions, missing the intimate competitive dynamic of 1v1 challenges. Users want to:
- Challenge specific friends to prediction duels
- Accept challenges from other users
- Track head-to-head records and statistics
- Receive notifications about challenge updates
- Compete on individual matches rather than full contests

## Solution Statement

Implement a challenge system that allows users to create direct challenges or open challenges on specific sports events. The system will handle invitation flows, dedicated scoring for 1v1 matchups, real-time notifications, and integration with existing Telegram bot functionality.

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Low-Medium
**Primary Systems Affected**: contest-service, notification-service, frontend, telegram-bot
**Dependencies**: Existing contest and notification infrastructure

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/contest-service/internal/models/contest.go` (lines 12-25) - Why: Contest model structure to mirror for challenges
- `backend/contest-service/internal/service/contest_service.go` (lines 18-22) - Why: Service pattern for challenge service
- `backend/contest-service/internal/repository/contest_repository.go` (lines 11-18) - Why: Repository interface pattern
- `backend/notification-service/internal/models/notification.go` (lines 10-22) - Why: Notification model structure for challenge notifications
- `backend/notification-service/internal/service/notification_service.go` (lines 26-64) - Why: Notification sending patterns
- `backend/proto/contest.proto` - Why: gRPC service definition patterns to follow
- `backend/scoring-service/internal/models/score.go` (lines 11-20) - Why: Scoring model for challenge scoring
- `bots/telegram/bot/handlers.go` (lines 1-50) - Why: Telegram bot integration patterns
- `frontend/src/types/contest.types.ts` (lines 3-17) - Why: TypeScript interface patterns
- `frontend/src/services/contest-service.ts` (lines 22-149) - Why: Frontend service patterns

### New Files to Create

- `backend/proto/challenge.proto` - gRPC service definition for challenges
- `backend/contest-service/internal/models/challenge.go` - Challenge data model
- `backend/contest-service/internal/repository/challenge_repository.go` - Challenge repository
- `backend/contest-service/internal/service/challenge_service.go` - Challenge business logic
- `frontend/src/types/challenge.types.ts` - TypeScript interfaces for challenges
- `frontend/src/services/challenge-service.ts` - Frontend challenge service
- `frontend/src/components/challenges/ChallengeCard.tsx` - Challenge display component
- `frontend/src/components/challenges/ChallengeForm.tsx` - Challenge creation form
- `frontend/src/components/challenges/ChallengeList.tsx` - List of challenges
- `frontend/src/pages/ChallengesPage.tsx` - Main challenges page
- `tests/contest-service/challenge_test.go` - Unit tests for challenge service
- `tests/e2e/challenge_test.go` - E2E tests for challenge workflows

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Head-to-Head Challenges Research](research/head-to-head-challenges-research.md)
  - Specific sections: Challenge flows, scoring mechanisms, notification patterns
  - Why: Comprehensive best practices and implementation patterns
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
  - Specific section: Service definition and implementation
  - Why: Required for implementing challenge gRPC service
- [GORM Documentation](https://gorm.io/docs/models.html)
  - Specific section: Model definition and relationships
  - Why: Database model patterns for challenges

### Patterns to Follow

**Naming Conventions:**
```go
// Go backend - follow existing patterns
type Challenge struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `gorm:"not null" json:"title"`
    Status      string    `gorm:"not null;default:'pending'" json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    gorm.Model
}
```

**Error Handling:**
```go
// Follow existing contest service pattern
func (s *ChallengeService) CreateChallenge(ctx context.Context, req *pb.CreateChallengeRequest) (*pb.CreateChallengeResponse, error) {
    if err := validateCreateChallengeRequest(req); err != nil {
        return &pb.CreateChallengeResponse{
            Response: &commonpb.Response{
                Success: false,
                Message: err.Error(),
            },
        }, nil
    }
    // Implementation...
}
```

**Logging Pattern:**
```go
// Use structured logging like other services
log.Printf("Creating challenge: challenger_id=%d, opponent_id=%d, event_id=%d", 
    req.ChallengerId, req.OpponentId, req.EventId)
```

**Frontend Service Pattern:**
```typescript
// Follow existing contest service structure
class ChallengeService {
  private client: GrpcClient;
  
  async createChallenge(request: CreateChallengeRequest): Promise<CreateChallengeResponse> {
    // Implementation following contest service pattern
  }
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up core challenge infrastructure with database models, gRPC service definitions, and basic repository patterns.

**Tasks:**
- Define challenge proto messages and service
- Create challenge database model with proper relationships
- Implement challenge repository with CRUD operations
- Set up challenge service with basic validation

### Phase 2: Core Implementation

Implement challenge business logic including creation, acceptance, and scoring mechanisms.

**Tasks:**
- Implement challenge creation and invitation flow
- Add challenge acceptance/decline logic
- Create challenge-specific scoring system
- Integrate with existing notification service

### Phase 3: Integration

Connect challenge system to frontend UI and Telegram bot integration.

**Tasks:**
- Create frontend challenge components and pages
- Add challenge endpoints to API gateway
- Integrate challenge notifications
- Add Telegram bot challenge commands

### Phase 4: Testing & Validation

Comprehensive testing of challenge workflows and edge cases.

**Tasks:**
- Implement unit tests for challenge service
- Create integration tests for challenge flows
- Add E2E tests for complete challenge workflows
- Test notification delivery and edge cases

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/proto/challenge.proto

- **IMPLEMENT**: gRPC service definition for challenges with CRUD operations
- **PATTERN**: Mirror contest.proto structure - file:backend/proto/contest.proto:1-130
- **IMPORTS**: google/protobuf/timestamp.proto, google/protobuf/empty.proto, google/api/annotations.proto, common.proto
- **GOTCHA**: Use consistent field naming with existing proto files (snake_case)
- **VALIDATE**: `cd backend && buf generate`

### CREATE backend/contest-service/internal/models/challenge.go

- **IMPLEMENT**: Challenge model with GORM tags and JSON serialization
- **PATTERN**: Mirror Contest model structure - file:backend/contest-service/internal/models/contest.go:12-25
- **IMPORTS**: gorm.io/gorm, time
- **GOTCHA**: Include proper database indexes for performance (challenger_id, opponent_id, status)
- **VALIDATE**: `cd backend/contest-service && go build ./...`

### CREATE backend/contest-service/internal/repository/challenge_repository.go

- **IMPLEMENT**: Challenge repository interface and implementation
- **PATTERN**: Mirror ContestRepository pattern - file:backend/contest-service/internal/repository/contest_repository.go:11-35
- **IMPORTS**: gorm.io/gorm, context, models package
- **GOTCHA**: Include methods for finding challenges by user ID and status
- **VALIDATE**: `cd backend/contest-service && go test ./internal/repository/...`

### CREATE backend/contest-service/internal/service/challenge_service.go

- **IMPLEMENT**: Challenge service with gRPC method implementations
- **PATTERN**: Mirror ContestService structure - file:backend/contest-service/internal/service/contest_service.go:18-22
- **IMPORTS**: context, proto packages, repository, models
- **GOTCHA**: Implement proper validation for challenge constraints (user can't challenge themselves)
- **VALIDATE**: `cd backend/contest-service && go test ./internal/service/...`

### UPDATE backend/contest-service/cmd/main.go

- **IMPLEMENT**: Register challenge service with gRPC server
- **PATTERN**: Follow existing service registration pattern in main.go
- **IMPORTS**: Add challenge service import
- **GOTCHA**: Ensure challenge service is initialized with proper dependencies
- **VALIDATE**: `cd backend/contest-service && go run cmd/main.go --help`

### UPDATE backend/proto/contest.proto

- **IMPLEMENT**: Add challenge-related imports and service registration
- **PATTERN**: Follow existing service definition patterns
- **IMPORTS**: Import challenge.proto
- **GOTCHA**: Maintain backward compatibility with existing contest service
- **VALIDATE**: `cd backend && buf generate`

### CREATE frontend/src/types/challenge.types.ts

- **IMPLEMENT**: TypeScript interfaces for challenge entities and requests
- **PATTERN**: Mirror contest.types.ts structure - file:frontend/src/types/contest.types.ts:3-17
- **IMPORTS**: Common types from common.types.ts
- **GOTCHA**: Ensure interface names match proto message names
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/services/challenge-service.ts

- **IMPLEMENT**: Frontend service for challenge API calls
- **PATTERN**: Mirror ContestService class - file:frontend/src/services/contest-service.ts:22-149
- **IMPORTS**: GrpcClient, challenge types
- **GOTCHA**: Handle gRPC-Web conversion properly for challenge messages
- **VALIDATE**: `cd frontend && npm run lint`

### CREATE frontend/src/components/challenges/ChallengeCard.tsx

- **IMPLEMENT**: Challenge display component with action buttons
- **PATTERN**: Mirror ContestCard component - file:frontend/src/components/contests/ContestCard.tsx:22-31
- **IMPORTS**: React, Material-UI components, challenge types
- **GOTCHA**: Handle different challenge states (pending, accepted, completed)
- **VALIDATE**: `cd frontend && npm test -- ChallengeCard`

### CREATE frontend/src/components/challenges/ChallengeForm.tsx

- **IMPLEMENT**: Challenge creation form with opponent selection
- **PATTERN**: Mirror ContestForm component - file:frontend/src/components/contests/ContestForm.tsx:24-30
- **IMPORTS**: React Hook Form, Zod validation, Material-UI
- **GOTCHA**: Include event selection and opponent search functionality
- **VALIDATE**: `cd frontend && npm test -- ChallengeForm`

### CREATE frontend/src/components/challenges/ChallengeList.tsx

- **IMPLEMENT**: List component for displaying user challenges
- **PATTERN**: Mirror ContestList component - file:frontend/src/components/contests/ContestList.tsx:27-31
- **IMPORTS**: React, challenge components, hooks
- **GOTCHA**: Support filtering by challenge status and type
- **VALIDATE**: `cd frontend && npm test -- ChallengeList`

### CREATE frontend/src/pages/ChallengesPage.tsx

- **IMPLEMENT**: Main challenges page with tabs for different challenge views
- **PATTERN**: Mirror ContestsPage structure - file:frontend/src/pages/ContestsPage.tsx
- **IMPORTS**: React, challenge components, Material-UI tabs
- **GOTCHA**: Include tabs for "My Challenges", "Pending Invites", "Create Challenge"
- **VALIDATE**: `cd frontend && npm run build`

### UPDATE frontend/src/App.tsx

- **IMPLEMENT**: Add challenges route to application routing
- **PATTERN**: Follow existing route definition pattern - file:frontend/src/App.tsx
- **IMPORTS**: ChallengesPage component
- **GOTCHA**: Ensure route is protected and requires authentication
- **VALIDATE**: `cd frontend && npm start`

### UPDATE backend/api-gateway/internal/gateway/gateway.go

- **IMPLEMENT**: Add challenge service proxy configuration
- **PATTERN**: Follow existing service proxy pattern in gateway.go
- **IMPORTS**: Challenge proto package
- **GOTCHA**: Ensure proper endpoint routing for challenge APIs
- **VALIDATE**: `cd backend/api-gateway && go test ./...`

### UPDATE backend/notification-service/internal/service/notification_service.go

- **IMPLEMENT**: Add challenge-specific notification templates
- **PATTERN**: Follow existing notification sending pattern - file:backend/notification-service/internal/service/notification_service.go:26-64
- **IMPORTS**: Challenge notification types
- **GOTCHA**: Include notification types for challenge_received, challenge_accepted, challenge_completed
- **VALIDATE**: `cd backend/notification-service && go test ./...`

### UPDATE bots/telegram/bot/handlers.go

- **IMPLEMENT**: Add challenge-related Telegram bot commands
- **PATTERN**: Follow existing handler pattern - file:bots/telegram/bot/handlers.go:1-50
- **IMPORTS**: Challenge service client
- **GOTCHA**: Include commands for /challenge, /accept_challenge, /my_challenges
- **VALIDATE**: `cd bots/telegram && go test ./...`

### CREATE tests/contest-service/challenge_test.go

- **IMPLEMENT**: Unit tests for challenge service methods
- **PATTERN**: Follow existing test patterns in contest service tests
- **IMPORTS**: Testing framework, challenge service, test helpers
- **GOTCHA**: Test edge cases like self-challenge prevention, duplicate challenges
- **VALIDATE**: `cd tests/contest-service && go test -v ./...`

### CREATE tests/e2e/challenge_test.go

- **IMPLEMENT**: End-to-end tests for complete challenge workflows
- **PATTERN**: Follow existing E2E test structure - file:tests/e2e/contest_test.go
- **IMPORTS**: E2E test helpers, HTTP client
- **GOTCHA**: Test full flow from challenge creation to completion
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -v ./...`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Ensure challenge service environment variables are configured
- **PATTERN**: Follow existing service configuration pattern
- **IMPORTS**: None
- **GOTCHA**: Verify all services can communicate with challenge endpoints
- **VALIDATE**: `make docker-services && make status`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests with fixtures and assertions following existing Go testing patterns in the project. Focus on:

- Challenge creation validation
- Challenge acceptance/decline logic
- Challenge scoring calculations
- Repository CRUD operations
- Service method error handling

### Integration Tests

Test service interactions and database operations:

- Challenge service with notification service integration
- Challenge repository with database operations
- gRPC service method implementations
- Frontend service API calls

### Edge Cases

Specific edge cases that must be tested for this feature:

- User attempting to challenge themselves
- Accepting already accepted/expired challenges
- Challenge timeout handling
- Concurrent challenge acceptance
- Invalid opponent ID handling
- Challenge scoring with missing predictions

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Backend linting and formatting
cd backend && go fmt ./...
cd backend && go vet ./...
cd backend/contest-service && golangci-lint run

# Frontend linting and formatting
cd frontend && npm run lint
cd frontend && npm run lint:fix
```

### Level 2: Unit Tests

```bash
# Backend unit tests
cd backend/contest-service && go test ./internal/...
cd backend/notification-service && go test ./internal/...
cd tests/contest-service && go test -v ./...

# Frontend unit tests
cd frontend && npm test -- --coverage
```

### Level 3: Integration Tests

```bash
# Service integration tests
cd tests/notification-service && go test -v ./...
cd tests/telegram-bot && go test -v ./...

# Build verification
cd backend && go build ./...
cd frontend && npm run build
```

### Level 4: Manual Validation

```bash
# Start services
make docker-up
make docker-services

# Test API endpoints
curl -X POST http://localhost:8080/v1/challenges \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"challenger_id": 1, "opponent_id": 2, "event_id": 1}'

# Test frontend
cd frontend && npm start
# Navigate to http://localhost:3000/challenges
```

### Level 5: End-to-End Tests

```bash
# Full E2E test suite
make e2e-test

# Specific challenge workflow tests
cd tests/e2e && go test -tags=e2e -v -run TestChallengeWorkflow ./...
```

---

## ACCEPTANCE CRITERIA

- [ ] Users can create direct challenges to specific opponents
- [ ] Users can accept or decline received challenges
- [ ] Challenge notifications are sent via all configured channels
- [ ] Challenge scoring works independently from contest scoring
- [ ] Telegram bot supports challenge commands (/challenge, /accept, /my_challenges)
- [ ] Frontend displays challenges with proper status indicators
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ requirement
- [ ] Integration tests verify service interactions
- [ ] E2E tests confirm complete challenge workflows
- [ ] No regressions in existing contest functionality
- [ ] Challenge timeouts are handled properly
- [ ] Database performance is maintained with proper indexing

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration + E2E)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms feature works end-to-end
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability
- [ ] Documentation updated for new API endpoints
- [ ] Performance impact assessed and acceptable

---

## NOTES

**Design Decisions:**
- Challenges are implemented as a separate entity from contests to maintain clear separation of concerns
- Challenge scoring is independent but follows similar patterns to contest scoring
- Notification integration reuses existing infrastructure for consistency
- Telegram bot integration follows established command patterns

**Trade-offs:**
- Simple direct challenges only (no open challenge board in MVP)
- Basic point-based scoring (no ELO rating system initially)
- Challenge expiration handled via application logic (not database triggers)

**Future Enhancements:**
- Open challenge boards for public matchmaking
- ELO rating system for skill-based matching
- Tournament bracket integration
- Advanced scoring mechanisms with confidence levels

# Feature: Head-to-Head Challenges

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Head-to-Head Challenges enable direct 1v1 competitions between users on specific sports events or series of events. Users can challenge each other directly, accept/decline challenges, and compete for points with dedicated scoring and leaderboards. This feature adds personalized competition mechanics that increase user engagement through direct social interaction.

## User Story

As a sports prediction enthusiast
I want to challenge another user to a direct 1v1 prediction competition
So that I can test my skills against specific opponents and enjoy personalized competition beyond general contests

## Problem Statement

The current platform only supports general contests where users compete against all participants. There's no mechanism for direct user-to-user challenges, which limits personalized engagement and social interaction. Users cannot create focused competitions with friends or rivals, missing opportunities for targeted engagement and retention.

## Solution Statement

Implement a comprehensive Head-to-Head Challenge system that allows users to:
- Send direct challenges to specific users on upcoming sports events
- Accept/decline challenge invitations with timeout handling
- Compete in isolated 1v1 scoring environments
- Receive real-time notifications for challenge events
- View challenge history and statistics
- Integrate seamlessly with existing contest, prediction, and scoring systems

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium (4-6 hours)
**Primary Systems Affected**: New challenge-service, contest-service, prediction-service, scoring-service, notification-service, frontend
**Dependencies**: Existing gRPC infrastructure, PostgreSQL, Redis, notification system

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/contest-service/internal/models/contest.go` (lines 12-25) - Why: Contest model structure to mirror for challenges
- `backend/contest-service/internal/models/participant.go` (lines 11-22) - Why: Participant model pattern for challenge participants
- `backend/contest-service/internal/service/contest_service.go` (lines 1-50) - Why: Service layer patterns and gRPC implementation
- `backend/contest-service/internal/repository/contest_repository.go` - Why: Repository interface patterns
- `backend/notification-service/internal/models/notification.go` - Why: Notification model structure and validation patterns
- `backend/proto/notification.proto` - Why: Notification types enum to extend for challenges
- `backend/proto/contest.proto` - Why: gRPC service patterns and message structures
- `backend/scoring-service/internal/models/streak.go` - Why: User streak model for challenge scoring integration
- `backend/api-gateway/internal/gateway/gateway.go` (lines 1-80) - Why: Service registration patterns for new challenge service
- `scripts/init-db.sql` - Why: Database schema patterns and indexing strategies
- `frontend/src/types/contest.types.ts` - Why: TypeScript interface patterns for frontend integration

### New Files to Create

- `backend/challenge-service/` - Complete microservice for challenge management
- `backend/proto/challenge.proto` - gRPC service definitions for challenges
- `backend/challenge-service/internal/models/challenge.go` - Challenge and ChallengeParticipant models
- `backend/challenge-service/internal/service/challenge_service.go` - gRPC service implementation
- `backend/challenge-service/internal/repository/challenge_repository.go` - Database operations
- `frontend/src/types/challenge.types.ts` - TypeScript interfaces for challenges
- `frontend/src/services/challenge-service.ts` - Frontend API client
- `frontend/src/components/challenges/` - Challenge UI components
- `tests/challenge-service/` - Unit tests for challenge service

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
  - Specific section: Service implementation patterns
  - Why: Required for implementing gRPC challenge service
- [GORM Documentation](https://gorm.io/docs/models.html)
  - Specific section: Model definition and relationships
  - Why: Shows proper model structure and validation patterns
- [Material-UI Components](https://mui.com/material-ui/react-dialog/)
  - Specific section: Dialog and form components
  - Why: Required for challenge invitation UI components

### Patterns to Follow

**Naming Conventions:**
- Go structs: `PascalCase` (Challenge, ChallengeParticipant)
- Go functions: `PascalCase` for public, `camelCase` for private
- Database tables: `snake_case` (challenges, challenge_participants)
- Proto messages: `PascalCase` with descriptive suffixes (CreateChallengeRequest)

**Error Handling:**
```go
return &pb.CreateChallengeResponse{
    Response: &common.Response{
        Success:   false,
        Message:   "Validation error message",
        Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
        Timestamp: timestamppb.Now(),
    },
}, nil
```

**Logging Pattern:**
```go
log.Printf("[ERROR] Failed to create challenge: %v", err)
log.Printf("[INFO] Challenge created successfully: ID=%d", challenge.ID)
```

**Model Validation Pattern:**
```go
func (c *Challenge) BeforeCreate(tx *gorm.DB) error {
    if err := c.ValidateUserIDs(); err != nil {
        return err
    }
    return c.ValidateEventID()
}
```

**Repository Interface Pattern:**
```go
type ChallengeRepositoryInterface interface {
    Create(challenge *Challenge) error
    GetByID(id uint) (*Challenge, error)
    // ... other methods
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up the challenge microservice infrastructure and core data models.

**Tasks:**
- Create challenge-service directory structure following existing service patterns
- Define Protocol Buffers schema for challenge operations
- Implement core Challenge and ChallengeParticipant models with validation
- Set up repository interfaces and database operations

### Phase 2: Core Implementation

Implement the main challenge business logic and gRPC service.

**Tasks:**
- Implement challenge service with CRUD operations
- Add challenge state management (pending, accepted, active, completed)
- Integrate with existing notification system for challenge events
- Implement timeout handling for challenge expiration

### Phase 3: Integration

Connect the challenge service to existing platform services and API gateway.

**Tasks:**
- Register challenge service in API Gateway
- Integrate with contest and prediction services for event data
- Connect to scoring service for challenge-specific scoring
- Add database migrations and indexes

### Phase 4: Frontend Implementation

Build user interface components for challenge management.

**Tasks:**
- Create TypeScript interfaces and API client
- Implement challenge creation and invitation UI
- Build challenge management dashboard
- Add challenge notifications to existing notification system

### Phase 5: Testing & Validation

Ensure comprehensive testing and system integration.

**Tasks:**
- Implement unit tests for all challenge service components
- Create integration tests for challenge workflows
- Add E2E tests for complete challenge lifecycle
- Validate performance and error handling

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/challenge-service directory structure

- **IMPLEMENT**: Complete microservice directory structure following existing patterns
- **PATTERN**: Mirror `backend/contest-service/` structure exactly
- **IMPORTS**: Copy Dockerfile, go.mod template from contest-service
- **GOTCHA**: Ensure go.mod has correct module path for challenge-service
- **VALIDATE**: `ls -la backend/challenge-service/` shows cmd/, internal/, Dockerfile, go.mod

### CREATE backend/proto/challenge.proto

- **IMPLEMENT**: gRPC service definition with all challenge operations
- **PATTERN**: Follow `backend/proto/contest.proto` message and service patterns
- **IMPORTS**: Import common.proto, google/protobuf/timestamp.proto, google/api/annotations.proto
- **GOTCHA**: Use consistent field numbering and proper HTTP annotations
- **VALIDATE**: `protoc --proto_path=backend/proto --go_out=backend/shared --go-grpc_out=backend/shared backend/proto/challenge.proto`

### UPDATE backend/proto/notification.proto

- **IMPLEMENT**: Add challenge notification types to NotificationType enum
- **PATTERN**: Follow existing enum pattern with sequential numbering
- **IMPORTS**: No new imports needed
- **GOTCHA**: Start numbering from 7 (after NEW_CONTEST = 6)
- **VALIDATE**: `grep -n "CHALLENGE_" backend/proto/notification.proto` shows new types

### CREATE backend/challenge-service/internal/models/challenge.go

- **IMPLEMENT**: Challenge and ChallengeParticipant models with full validation
- **PATTERN**: Mirror `backend/contest-service/internal/models/contest.go` structure
- **IMPORTS**: gorm.io/gorm, time, errors packages
- **GOTCHA**: Use proper GORM tags and validation methods like existing models
- **VALIDATE**: `go build backend/challenge-service/internal/models/`

### CREATE backend/challenge-service/internal/repository/challenge_repository.go

- **IMPLEMENT**: Repository interface and implementation for challenge operations
- **PATTERN**: Follow `backend/contest-service/internal/repository/contest_repository.go` exactly
- **IMPORTS**: gorm.io/gorm, models package, context
- **GOTCHA**: Include proper error handling and transaction support
- **VALIDATE**: `go build backend/challenge-service/internal/repository/`

### CREATE backend/challenge-service/internal/service/challenge_service.go

- **IMPLEMENT**: gRPC service implementation with all challenge operations
- **PATTERN**: Mirror `backend/contest-service/internal/service/contest_service.go` structure
- **IMPORTS**: context, log, time, proto packages, auth package
- **GOTCHA**: Use auth.GetUserIDFromContext for authentication like existing services
- **VALIDATE**: `go build backend/challenge-service/internal/service/`

### CREATE backend/challenge-service/internal/config/config.go

- **IMPLEMENT**: Configuration structure for challenge service
- **PATTERN**: Copy from `backend/contest-service/internal/config/` exactly
- **IMPORTS**: os package for environment variables
- **GOTCHA**: Update service name and port environment variables
- **VALIDATE**: `go build backend/challenge-service/internal/config/`

### CREATE backend/challenge-service/cmd/main.go

- **IMPLEMENT**: Main entry point for challenge service
- **PATTERN**: Copy from `backend/contest-service/cmd/main.go` and adapt
- **IMPORTS**: All internal packages, gRPC packages
- **GOTCHA**: Update service registration and port configuration
- **VALIDATE**: `go build backend/challenge-service/cmd/`

### UPDATE backend/go.work

- **IMPLEMENT**: Add challenge-service to Go workspace
- **PATTERN**: Follow existing service entries in go.work
- **IMPORTS**: No imports needed
- **GOTCHA**: Maintain alphabetical order of services
- **VALIDATE**: `cd backend && go work sync`

### UPDATE scripts/init-db.sql

- **IMPLEMENT**: Add challenges and challenge_participants tables
- **PATTERN**: Follow existing table creation patterns with proper indexes
- **IMPORTS**: No imports needed
- **GOTCHA**: Use proper foreign key constraints and unique indexes
- **VALIDATE**: `psql -h localhost -U sports_user -d sports_prediction -f scripts/init-db.sql`

### UPDATE backend/api-gateway/internal/gateway/gateway.go

- **IMPLEMENT**: Register challenge service in API Gateway
- **PATTERN**: Follow existing service registration pattern (lines 40-60)
- **IMPORTS**: Add challengepb import
- **GOTCHA**: Add to both import and registration sections
- **VALIDATE**: `go build backend/api-gateway/`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Add challenge-service container configuration
- **PATTERN**: Follow existing service container patterns
- **IMPORTS**: No imports needed
- **GOTCHA**: Use correct port (8090) and environment variables
- **VALIDATE**: `docker-compose config` validates YAML syntax

### CREATE frontend/src/types/challenge.types.ts

- **IMPLEMENT**: TypeScript interfaces for all challenge-related types
- **PATTERN**: Mirror `frontend/src/types/contest.types.ts` structure exactly
- **IMPORTS**: No imports needed for type definitions
- **GOTCHA**: Match proto field names with camelCase conversion
- **VALIDATE**: `cd frontend && npm run build` compiles without errors

### CREATE frontend/src/services/challenge-service.ts

- **IMPLEMENT**: Frontend API client for challenge operations
- **PATTERN**: Follow `frontend/src/services/contest-service.ts` class structure
- **IMPORTS**: axios, challenge types, API response types
- **GOTCHA**: Use proper error handling and response transformation
- **VALIDATE**: `cd frontend && npm run lint` passes without errors

### CREATE frontend/src/components/challenges/ChallengeDialog.tsx

- **IMPLEMENT**: Modal dialog for creating new challenges
- **PATTERN**: Follow `frontend/src/components/contests/ContestForm.tsx` patterns
- **IMPORTS**: React, Material-UI Dialog, react-hook-form, zod
- **GOTCHA**: Include proper form validation and error handling
- **VALIDATE**: `cd frontend && npm run build` compiles component successfully

### CREATE frontend/src/components/challenges/ChallengeList.tsx

- **IMPLEMENT**: List component for displaying user challenges
- **PATTERN**: Follow `frontend/src/components/contests/ContestList.tsx` structure
- **IMPORTS**: React, Material-UI components, challenge types
- **GOTCHA**: Include proper loading states and error handling
- **VALIDATE**: Component renders without console errors in development

### CREATE frontend/src/components/challenges/ChallengeCard.tsx

- **IMPLEMENT**: Card component for individual challenge display
- **PATTERN**: Follow `frontend/src/components/contests/ContestCard.tsx` layout
- **IMPORTS**: React, Material-UI Card components, challenge types
- **GOTCHA**: Include challenge status indicators and action buttons
- **VALIDATE**: Component displays challenge data correctly

### UPDATE frontend/src/components/contests/ContestDetail.tsx

- **IMPLEMENT**: Add "Challenge User" button to contest participant list
- **PATTERN**: Follow existing button patterns in the component
- **IMPORTS**: Add ChallengeDialog import
- **GOTCHA**: Only show button for other participants, not current user
- **VALIDATE**: Button appears and opens challenge dialog correctly

### CREATE tests/challenge-service/challenge_test.go

- **IMPLEMENT**: Unit tests for challenge service operations
- **PATTERN**: Follow `tests/contest-service/contest_test.go` test structure
- **IMPORTS**: testing, testify, challenge service packages
- **GOTCHA**: Include tests for all CRUD operations and validation
- **VALIDATE**: `cd tests/challenge-service && go test -v`

### CREATE tests/e2e/challenge_test.go

- **IMPLEMENT**: End-to-end tests for complete challenge workflow
- **PATTERN**: Follow existing E2E test patterns in tests/e2e/
- **IMPORTS**: testing, HTTP client, E2E test utilities
- **GOTCHA**: Test complete challenge lifecycle from creation to completion
- **VALIDATE**: `cd tests/e2e && go test -tags=e2e -v ./challenge_test.go`

### UPDATE backend/shared/seeder/coordinator.go

- **IMPLEMENT**: Add challenge seeding to fake data generation
- **PATTERN**: Follow existing seeding patterns for contests and participants
- **IMPORTS**: Add challenge models import
- **GOTCHA**: Create realistic challenge data with proper relationships
- **VALIDATE**: `make seed-small` includes challenge data in output

---

## TESTING STRATEGY

### Unit Tests

Design unit tests with fixtures and assertions following existing Go testing patterns:

- **Challenge Model Tests**: Validation, state transitions, business logic
- **Repository Tests**: CRUD operations, error handling, transaction safety
- **Service Tests**: gRPC method implementations, authentication, authorization
- **Frontend Component Tests**: User interactions, form validation, API integration

### Integration Tests

- **Database Integration**: Challenge creation, participant management, state updates
- **Service Integration**: Challenge service communication with contest, prediction, scoring services
- **Notification Integration**: Challenge event notifications through existing notification system
- **API Gateway Integration**: HTTP routing and gRPC proxy functionality

### Edge Cases

- **Concurrent Challenge Operations**: Multiple users accepting same challenge
- **Challenge Timeout Handling**: Automatic expiration and cleanup
- **Invalid Challenge States**: Preventing invalid state transitions
- **User Permission Validation**: Ensuring users can only manage their own challenges
- **Event Validation**: Challenges on non-existent or past events

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Go code formatting and linting
cd backend/challenge-service && go fmt ./...
cd backend/challenge-service && go vet ./...

# Frontend linting
cd frontend && npm run lint
cd frontend && npm run type-check
```

### Level 2: Unit Tests

```bash
# Backend unit tests
cd backend/challenge-service && go test ./...
cd tests/challenge-service && go test -v

# Frontend unit tests
cd frontend && npm test -- --coverage --watchAll=false
```

### Level 3: Integration Tests

```bash
# Database integration
make docker-up
cd tests/challenge-service && go test -tags=integration -v

# E2E tests
make e2e-test
```

### Level 4: Manual Validation

```bash
# Start all services
make docker-services

# Test challenge creation API
curl -X POST http://localhost:8080/v1/challenges \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"opponent_id": 2, "event_id": 1, "message": "Test challenge"}'

# Test challenge acceptance
curl -X PUT http://localhost:8080/v1/challenges/1/accept \
  -H "Authorization: Bearer <token>"

# Verify frontend challenge components
# Navigate to http://localhost:3000 and test challenge creation flow
```

### Level 5: Additional Validation

```bash
# Protocol buffer compilation
make proto

# Docker build validation
docker-compose build challenge-service

# Database migration validation
make seed-small
```

---

## ACCEPTANCE CRITERIA

- [ ] Users can create challenges targeting specific opponents on upcoming events
- [ ] Challenge invitations are sent via notification system (in-app, Telegram, email)
- [ ] Recipients can accept or decline challenges with proper state management
- [ ] Challenges automatically expire after 24 hours if not accepted
- [ ] Accepted challenges create isolated scoring environments for participants
- [ ] Challenge results are calculated and displayed with winner determination
- [ ] Challenge history and statistics are accessible to participants
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ requirement
- [ ] Integration tests verify end-to-end challenge workflows
- [ ] Code follows project conventions and patterns
- [ ] No regressions in existing functionality
- [ ] Frontend components integrate seamlessly with existing UI
- [ ] Performance meets requirements (sub-200ms API responses)
- [ ] Security considerations addressed (user authorization, input validation)

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
- [ ] Database migrations applied successfully
- [ ] API Gateway properly routes challenge requests
- [ ] Frontend components render and function correctly
- [ ] Notification system delivers challenge events
- [ ] Challenge timeout handling works correctly

---

## NOTES

### Design Decisions

- **Microservice Architecture**: Challenge service follows existing microservice patterns for consistency and scalability
- **State Management**: Challenge states (pending, accepted, active, completed, declined, expired) provide clear lifecycle management
- **Notification Integration**: Leverages existing multi-channel notification system for challenge events
- **Database Design**: Separate challenges and challenge_participants tables mirror contest service patterns
- **Frontend Integration**: Challenge components integrate with existing contest UI for familiar user experience

### Trade-offs

- **Complexity vs Features**: Starting with basic 1v1 challenges, can extend to group challenges later
- **Real-time vs Polling**: Using existing notification system rather than WebSocket for simplicity
- **Storage vs Performance**: Storing challenge history for analytics vs database size considerations

### Future Enhancements

- **ELO Rating System**: Add skill-based matchmaking for challenges
- **Challenge Templates**: Pre-configured challenge types for common scenarios
- **Tournament Brackets**: Multi-round elimination challenges
- **Wagering System**: Optional point stakes for challenges
- **Challenge Leagues**: Seasonal challenge competitions with rankings

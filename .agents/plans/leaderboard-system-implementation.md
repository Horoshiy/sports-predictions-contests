# Feature: Leaderboard System

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement a comprehensive real-time leaderboard system for the sports prediction contests platform. The system will track user scores, calculate rankings, and provide real-time updates across multiple contests. It includes scoring logic for predictions, leaderboard calculations with proper tie-breaking, caching for performance, and both API endpoints and frontend UI components.

## User Story

As a contest participant
I want to see real-time leaderboards showing my ranking and points
So that I can track my performance against other participants and stay engaged in the competition

## Problem Statement

The current platform lacks any scoring or ranking system. Users can make predictions but have no way to see how they're performing compared to others. There's no scoring service, no leaderboard calculations, and no way to track contest standings, which significantly reduces user engagement and competitive motivation.

## Solution Statement

Build a complete scoring and leaderboard system with:
1. **Scoring Service**: Calculate points based on prediction accuracy when events complete
2. **Leaderboard Engine**: Real-time ranking calculations with Redis caching
3. **API Layer**: gRPC endpoints for score management and leaderboard queries
4. **Frontend Integration**: Leaderboard components with live updates
5. **Performance Optimization**: Multi-layer caching and efficient ranking algorithms

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: New scoring-service, contest-service, prediction-service, frontend
**Dependencies**: Redis (existing), PostgreSQL (existing), gRPC (existing)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/contest-service/internal/models/contest.go` (lines 1-150) - Why: Contest model structure and validation patterns to follow
- `backend/contest-service/internal/models/participant.go` (lines 1-100) - Why: Participant model for leaderboard user relationships
- `backend/prediction-service/internal/models/prediction.go` (lines 1-120) - Why: Prediction model with scoring status field
- `backend/prediction-service/internal/models/event.go` (lines 1-100) - Why: Event model with result data structure
- `backend/shared/database/connection.go` (lines 1-50) - Why: Database connection patterns and GORM setup
- `backend/shared/auth/jwt.go` (lines 1-80) - Why: JWT claims structure for user authentication
- `backend/proto/contest.proto` (lines 1-150) - Why: Existing gRPC service patterns and message structures
- `backend/proto/common.proto` (lines 1-30) - Why: Common response and pagination patterns
- `backend/api-gateway/internal/gateway/gateway.go` (lines 1-100) - Why: Service registration and routing patterns
- `frontend/src/services/contest-service.ts` (lines 1-150) - Why: Frontend service patterns and API client structure
- `tests/contest-service/contest_test.go` (lines 1-50) - Why: Testing patterns and validation approaches

### New Files to Create

- `backend/scoring-service/cmd/main.go` - Main entry point for scoring service
- `backend/scoring-service/internal/config/config.go` - Service configuration
- `backend/scoring-service/internal/models/score.go` - Score data model
- `backend/scoring-service/internal/models/leaderboard.go` - Leaderboard data model
- `backend/scoring-service/internal/repository/score_repository.go` - Score database operations
- `backend/scoring-service/internal/repository/leaderboard_repository.go` - Leaderboard database operations
- `backend/scoring-service/internal/service/scoring_service.go` - Core scoring business logic
- `backend/scoring-service/internal/service/leaderboard_service.go` - Leaderboard business logic
- `backend/scoring-service/internal/cache/redis_cache.go` - Redis caching layer
- `backend/scoring-service/go.mod` - Go module definition
- `backend/scoring-service/Dockerfile` - Container configuration
- `backend/proto/scoring.proto` - gRPC service definitions for scoring
- `frontend/src/services/scoring-service.ts` - Frontend scoring API client
- `frontend/src/components/leaderboard/LeaderboardTable.tsx` - Leaderboard display component
- `frontend/src/components/leaderboard/UserScore.tsx` - Individual user score component
- `frontend/src/types/scoring.types.ts` - TypeScript type definitions
- `tests/scoring-service/scoring_test.go` - Unit tests for scoring logic
- `tests/scoring-service/leaderboard_test.go` - Unit tests for leaderboard logic

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [GORM Documentation - Models](https://gorm.io/docs/models.html)
  - Specific section: Model definition and relationships
  - Why: Required for implementing Score and Leaderboard models with proper GORM tags
- [Redis Go Client](https://redis.uptrace.dev/guide/server.html#connecting-to-redis-server)
  - Specific section: Sorted Sets operations
  - Why: Needed for efficient leaderboard caching with O(log N) operations
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
  - Specific section: Service implementation patterns
  - Why: Shows proper gRPC service structure for scoring endpoints
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers/docs/proto3)
  - Specific section: Message definitions and field types
  - Why: Required for defining scoring service proto messages
- [React Query Documentation](https://tanstack.com/query/latest/docs/react/overview)
  - Specific section: Real-time updates and polling
  - Why: Needed for implementing live leaderboard updates in frontend

### Patterns to Follow

**GORM Model Pattern** (from contest.go):
```go
type Model struct {
    ID uint `gorm:"primaryKey" json:"id"`
    // fields with validation tags
    gorm.Model
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
    return m.validateFields()
}
```

**gRPC Service Pattern** (from contest.proto):
```proto
service ServiceName {
  rpc MethodName(RequestMessage) returns (ResponseMessage) {
    option (google.api.http) = {
      post: "/v1/path"
      body: "*"
    };
  }
}
```

**Repository Pattern** (from contest-service):
```go
type Repository interface {
    Create(ctx context.Context, entity *Model) error
    GetByID(ctx context.Context, id uint) (*Model, error)
    Update(ctx context.Context, entity *Model) error
    Delete(ctx context.Context, id uint) error
}
```

**Frontend Service Pattern** (from contest-service.ts):
```typescript
class Service {
  private basePath = '/v1/resource'
  
  async method(request: RequestType): Promise<ResponseType> {
    const response = await grpcClient.post<ApiResponse>(this.basePath, request)
    return response.data
  }
}
```

**Error Handling Pattern** (from existing models):
```go
func (m *Model) ValidateField() error {
    if condition {
        return errors.New("validation message")
    }
    return nil
}
```

**Testing Pattern** (from contest_test.go):
```go
func TestFunction(t *testing.T) {
    // Setup
    entity := &Model{...}
    
    // Test
    err := entity.Method()
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up the scoring service infrastructure, data models, and database schema. This includes creating the service structure, defining scoring and leaderboard models with proper GORM relationships, and establishing the Redis caching layer.

**Tasks:**
- Create scoring service directory structure and Go module
- Implement Score and Leaderboard data models with GORM
- Set up Redis caching infrastructure for leaderboard operations
- Create database migration scripts for new tables

### Phase 2: Core Scoring Logic

Implement the core business logic for calculating scores based on prediction accuracy and managing leaderboard rankings. This includes scoring algorithms, ranking calculations with tie-breaking, and batch processing capabilities.

**Tasks:**
- Implement scoring algorithms for different prediction types
- Create leaderboard calculation engine with ranking logic
- Add score repository with database operations
- Implement leaderboard repository with caching integration

### Phase 3: gRPC API Layer

Create the gRPC service definitions and implement the API endpoints for score management and leaderboard queries. This includes proto definitions, service implementation, and integration with existing API gateway.

**Tasks:**
- Define scoring service proto with all required messages and endpoints
- Implement gRPC service handlers for scoring operations
- Add leaderboard query endpoints with pagination support
- Register scoring service in API gateway routing

### Phase 4: Frontend Integration

Build React components for displaying leaderboards and user scores, integrate with the backend API, and implement real-time updates for live leaderboard functionality.

**Tasks:**
- Create TypeScript types for scoring API responses
- Implement frontend scoring service client
- Build leaderboard table component with sorting and pagination
- Add real-time updates using React Query polling

### Phase 5: Testing & Validation

Implement comprehensive testing for all components, including unit tests, integration tests, and performance validation for the caching layer.

**Tasks:**
- Create unit tests for scoring algorithms and leaderboard logic
- Add integration tests for gRPC endpoints
- Implement performance tests for Redis caching operations
- Add frontend component tests for leaderboard UI

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/scoring-service/go.mod

- **IMPLEMENT**: Go module with dependencies for gRPC, GORM, Redis, and JWT
- **PATTERN**: Mirror backend/contest-service/go.mod structure
- **IMPORTS**: gorm.io/gorm, gorm.io/driver/postgres, github.com/go-redis/redis/v8, google.golang.org/grpc
- **GOTCHA**: Use same Go version (1.21) and compatible dependency versions as other services
- **VALIDATE**: `cd backend/scoring-service && go mod tidy && go mod verify`

### CREATE backend/scoring-service/internal/config/config.go

- **IMPLEMENT**: Service configuration struct with database, Redis, and gRPC settings
- **PATTERN**: Mirror backend/contest-service/internal/config/config.go:1-50
- **IMPORTS**: os package for environment variables, validation methods
- **GOTCHA**: Include JWT secret and service port configuration
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/config`

### CREATE backend/scoring-service/internal/models/score.go

- **IMPLEMENT**: Score model with UserID, ContestID, PredictionID, Points, ScoredAt fields
- **PATTERN**: Mirror backend/contest-service/internal/models/contest.go:1-30 for GORM structure
- **IMPORTS**: gorm.io/gorm, time, errors packages
- **GOTCHA**: Add unique constraint on (UserID, ContestID, PredictionID), include validation methods
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/models`

### CREATE backend/scoring-service/internal/models/leaderboard.go

- **IMPLEMENT**: Leaderboard model with ContestID, UserID, TotalPoints, Rank, UpdatedAt fields
- **PATTERN**: Mirror backend/contest-service/internal/models/participant.go:1-40 for relationships
- **IMPORTS**: gorm.io/gorm, time package
- **GOTCHA**: Add unique constraint on (ContestID, UserID), index on (ContestID, Rank)
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/models`

### CREATE backend/scoring-service/internal/cache/redis_cache.go

- **IMPLEMENT**: Redis client wrapper with leaderboard-specific operations (ZADD, ZRANGE, ZRANK)
- **PATTERN**: Mirror backend/shared/database/connection.go:1-40 for connection setup
- **IMPORTS**: github.com/go-redis/redis/v8, context, fmt, strconv
- **GOTCHA**: Use Redis sorted sets for O(log N) leaderboard operations, implement connection pooling
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/cache`

### CREATE backend/proto/scoring.proto

- **IMPLEMENT**: gRPC service definition with Score, Leaderboard messages and CRUD operations
- **PATTERN**: Mirror backend/proto/contest.proto:1-150 for service structure and HTTP annotations
- **IMPORTS**: google/protobuf/timestamp.proto, google/api/annotations.proto, common.proto
- **GOTCHA**: Include pagination for leaderboard queries, add real-time update endpoints
- **VALIDATE**: `cd backend && protoc --proto_path=proto --go_out=shared --go-grpc_out=shared proto/scoring.proto`

### UPDATE backend/go.work

- **IMPLEMENT**: Add ./scoring-service to the workspace use directive
- **PATTERN**: Mirror existing service entries in backend/go.work:3-10
- **IMPORTS**: None required
- **GOTCHA**: Maintain alphabetical order of services
- **VALIDATE**: `cd backend && go work sync`

### CREATE backend/scoring-service/internal/repository/score_repository.go

- **IMPLEMENT**: Score repository interface and implementation with CRUD operations
- **PATTERN**: Mirror backend/contest-service/internal/repository/contest_repository.go:1-100
- **IMPORTS**: context, gorm.io/gorm, models package
- **GOTCHA**: Include batch operations for bulk score updates, add GetByContestAndUser method
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/repository`

### CREATE backend/scoring-service/internal/repository/leaderboard_repository.go

- **IMPLEMENT**: Leaderboard repository with ranking calculations and cache integration
- **PATTERN**: Mirror backend/contest-service/internal/repository/contest_repository.go:1-100
- **IMPORTS**: context, gorm.io/gorm, models and cache packages
- **GOTCHA**: Implement UpdateRankings method with transaction support, cache invalidation
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/repository`

### CREATE backend/scoring-service/internal/service/scoring_service.go

- **IMPLEMENT**: Core scoring business logic with prediction evaluation algorithms
- **PATTERN**: Mirror backend/contest-service/internal/service/contest_service.go:1-150
- **IMPORTS**: context, models, repository packages, prediction service client
- **GOTCHA**: Implement different scoring algorithms (exact match, close match, bonus points)
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/service`

### CREATE backend/scoring-service/internal/service/leaderboard_service.go

- **IMPLEMENT**: Leaderboard business logic with ranking calculations and real-time updates
- **PATTERN**: Mirror backend/contest-service/internal/service/contest_service.go:1-150
- **IMPORTS**: context, models, repository, cache packages
- **GOTCHA**: Implement tie-breaking logic, batch ranking updates, cache warming
- **VALIDATE**: `cd backend/scoring-service && go build ./internal/service`

### CREATE backend/scoring-service/cmd/main.go

- **IMPLEMENT**: Main service entry point with gRPC server setup and graceful shutdown
- **PATTERN**: Mirror backend/contest-service/cmd/main.go:1-80
- **IMPORTS**: config, service packages, gRPC server setup
- **GOTCHA**: Include health check endpoint, proper signal handling for shutdown
- **VALIDATE**: `cd backend/scoring-service && go build ./cmd`

### CREATE backend/scoring-service/Dockerfile

- **IMPLEMENT**: Multi-stage Docker build for scoring service
- **PATTERN**: Mirror backend/contest-service/Dockerfile:1-30
- **IMPORTS**: None required
- **GOTCHA**: Use same base image and Go version as other services
- **VALIDATE**: `cd backend/scoring-service && docker build -t scoring-service .`

### UPDATE backend/api-gateway/internal/gateway/gateway.go

- **IMPLEMENT**: Add scoring service client and route registration
- **PATTERN**: Mirror existing service registration in gateway.go:50-80
- **IMPORTS**: scoring service proto package
- **GOTCHA**: Add scoring service endpoint configuration, include in health checks
- **VALIDATE**: `cd backend/api-gateway && go build ./internal/gateway`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Add scoring-service container configuration
- **PATTERN**: Mirror existing service definitions in docker-compose.yml:50-80
- **IMPORTS**: None required
- **GOTCHA**: Include proper environment variables and network configuration
- **VALIDATE**: `docker-compose config`

### CREATE frontend/src/types/scoring.types.ts

- **IMPLEMENT**: TypeScript interfaces for Score, Leaderboard, and API request/response types
- **PATTERN**: Mirror frontend/src/types/contest.types.ts:1-100
- **IMPORTS**: None required
- **GOTCHA**: Include pagination types, match backend proto message structure
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/services/scoring-service.ts

- **IMPLEMENT**: Frontend API client for scoring service with all CRUD operations
- **PATTERN**: Mirror frontend/src/services/contest-service.ts:1-150
- **IMPORTS**: grpc-client, scoring types
- **GOTCHA**: Include real-time polling methods, proper error handling
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/leaderboard/LeaderboardTable.tsx

- **IMPLEMENT**: React component for displaying contest leaderboards with sorting and pagination
- **PATTERN**: Mirror frontend/src/components/contests/ContestList.tsx:1-100
- **IMPORTS**: React, Material-UI components, scoring service and types
- **GOTCHA**: Include real-time updates, user highlighting, responsive design
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/leaderboard/UserScore.tsx

- **IMPLEMENT**: Component for displaying individual user score and rank
- **PATTERN**: Mirror frontend/src/components/contests/ContestCard.tsx:1-80
- **IMPORTS**: React, Material-UI components, scoring types
- **GOTCHA**: Include score breakdown, rank change indicators, loading states
- **VALIDATE**: `cd frontend && npm run build`

### UPDATE frontend/src/pages/ContestsPage.tsx

- **IMPLEMENT**: Add leaderboard tab/section to contest page
- **PATTERN**: Mirror existing tab structure in ContestsPage.tsx:30-60
- **IMPORTS**: LeaderboardTable component
- **GOTCHA**: Include conditional rendering based on contest status
- **VALIDATE**: `cd frontend && npm run build`

### CREATE tests/scoring-service/scoring_test.go

- **IMPLEMENT**: Unit tests for scoring algorithms and business logic
- **PATTERN**: Mirror tests/contest-service/contest_test.go:1-50
- **IMPORTS**: testing, scoring service packages
- **GOTCHA**: Test different prediction types, edge cases, error conditions
- **VALIDATE**: `cd tests/scoring-service && go test -v`

### CREATE tests/scoring-service/leaderboard_test.go

- **IMPLEMENT**: Unit tests for leaderboard calculations and ranking logic
- **PATTERN**: Mirror tests/contest-service/contest_test.go:1-50
- **IMPORTS**: testing, leaderboard service packages
- **GOTCHA**: Test tie-breaking, ranking updates, cache consistency
- **VALIDATE**: `cd tests/scoring-service && go test -v`

### CREATE tests/scoring-service/integration_test.go

- **IMPLEMENT**: Integration tests for gRPC endpoints and database operations
- **PATTERN**: Mirror tests/contest-service/integration_test.go:1-100
- **IMPORTS**: testing, gRPC client, database setup
- **GOTCHA**: Include Redis cache testing, concurrent access scenarios
- **VALIDATE**: `cd tests/scoring-service && go test -v -tags=integration`

### UPDATE scripts/init-db.sql

- **IMPLEMENT**: Add database schema for scores and leaderboards tables
- **PATTERN**: Mirror existing table creation in init-db.sql:10-20
- **IMPORTS**: None required
- **GOTCHA**: Include proper indexes for performance, foreign key constraints
- **VALIDATE**: `docker-compose exec postgres psql -U sports_user -d sports_prediction -f /docker-entrypoint-initdb.d/init-db.sql`

---

## TESTING STRATEGY

### Unit Tests

**Scope**: Test individual components in isolation - scoring algorithms, leaderboard calculations, data models, and business logic methods.

Design unit tests with fixtures and assertions following existing testing approaches:
- Score calculation accuracy for different prediction types
- Leaderboard ranking logic with tie-breaking scenarios
- Model validation methods and GORM hooks
- Repository methods with mock database interactions
- Cache operations with Redis mock

### Integration Tests

**Scope**: Test complete workflows across service boundaries - gRPC endpoints, database operations, cache consistency, and inter-service communication.

Integration test scenarios:
- End-to-end scoring workflow from prediction to leaderboard update
- gRPC API endpoints with real database and Redis instances
- Concurrent score updates and ranking calculations
- Cache invalidation and consistency across multiple operations
- Performance testing for leaderboard queries with large datasets

### Edge Cases

**Specific edge cases that must be tested for this feature:**
- Simultaneous score updates for the same contest
- Leaderboard calculations with identical scores (tie-breaking)
- Cache failures and fallback to database queries
- Invalid prediction data and scoring error handling
- Contest completion and final leaderboard freezing
- User leaving contest after scoring has begun
- Redis connection failures and recovery scenarios
- Large contest leaderboards (10,000+ participants)

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Go formatting and linting
cd backend/scoring-service && go fmt ./...
cd backend/scoring-service && go vet ./...
cd backend && go work sync

# Frontend linting
cd frontend && npm run lint
cd frontend && npm run lint:fix
```

### Level 2: Unit Tests

```bash
# Backend unit tests
cd backend/scoring-service && go test -v ./...
cd tests/scoring-service && go test -v

# Frontend unit tests
cd frontend && npm test -- --coverage
```

### Level 3: Integration Tests

```bash
# Start test environment
make docker-up

# Run integration tests
cd tests/scoring-service && go test -v -tags=integration

# Test gRPC endpoints
cd backend && go test -v ./api-gateway/...
```

### Level 4: Manual Validation

```bash
# Start all services
make docker-services

# Test scoring API endpoints
curl -X POST http://localhost:8080/v1/scores \
  -H "Content-Type: application/json" \
  -d '{"contest_id": 1, "user_id": 1, "prediction_id": 1, "points": 10}'

# Test leaderboard endpoint
curl http://localhost:8080/v1/contests/1/leaderboard

# Test frontend build
cd frontend && npm run build && npm run preview
```

### Level 5: Additional Validation (Optional)

```bash
# Performance testing
cd backend/scoring-service && go test -bench=. -benchmem

# Redis operations testing
redis-cli -h localhost -p 6379 ping
redis-cli -h localhost -p 6379 zrange contest:1:leaderboard 0 -1 withscores

# Database schema validation
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "\dt"
```

---

## ACCEPTANCE CRITERIA

- [ ] Scoring service calculates points accurately for different prediction types
- [ ] Leaderboards display real-time rankings with proper tie-breaking
- [ ] Redis caching provides sub-100ms leaderboard query performance
- [ ] gRPC API endpoints handle all CRUD operations for scores and leaderboards
- [ ] Frontend components display leaderboards with live updates
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ requirement for all new code
- [ ] Integration tests verify end-to-end scoring workflows
- [ ] Code follows existing project conventions and patterns
- [ ] No regressions in existing contest and prediction functionality
- [ ] Performance requirements met: <100ms API responses, 10,000+ concurrent users
- [ ] Security considerations addressed: proper authentication and input validation
- [ ] Database migrations execute successfully without data loss

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in dependency order
- [ ] Each task validation passed immediately after implementation
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms leaderboard functionality works end-to-end
- [ ] Performance benchmarks meet requirements
- [ ] Acceptance criteria all verified and met
- [ ] Code reviewed for quality, security, and maintainability
- [ ] Documentation updated with API endpoints and usage examples

---

## NOTES

**Design Decisions:**
- Used Redis sorted sets for O(log N) leaderboard operations instead of database-only approach
- Implemented multi-layer caching (Redis + in-memory) for optimal performance
- Chose standard competition ranking (1224) with tie-breaking by submission time
- Separated scoring and leaderboard services for better scalability and maintainability

**Performance Considerations:**
- Batch processing for score updates to reduce database load
- Precomputed leaderboard snapshots for high-traffic contests
- Connection pooling for Redis and database connections
- Pagination for large leaderboard queries

**Security Considerations:**
- JWT authentication required for all scoring operations
- Input validation for all score calculations and leaderboard queries
- Rate limiting on leaderboard API endpoints to prevent abuse
- Audit logging for all score modifications

**Future Enhancements:**
- WebSocket integration for real-time leaderboard updates
- Advanced scoring algorithms (ELO rating, weighted scoring)
- Leaderboard analytics and historical tracking
- Mobile push notifications for rank changes

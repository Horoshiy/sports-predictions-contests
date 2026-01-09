# Feature: Prediction Service Implementation

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement a comprehensive Prediction Service microservice that handles user predictions for sports contests. The service manages prediction submission, validation, storage, and retrieval with support for multiple prediction types (match outcomes, exact scores, player statistics) and flexible scoring rules. It integrates with the existing Contest Service to validate contest participation and enforce prediction deadlines.

## User Story

As a contest participant
I want to submit predictions for sports events in active contests
So that I can compete with other users and earn points based on prediction accuracy

## Problem Statement

The platform currently has Contest and User services but lacks the core functionality for users to submit predictions on sports events. Without a Prediction Service, users cannot participate in contests by making predictions, which is the fundamental feature of the sports prediction platform.

## Solution Statement

Implement a dedicated Prediction Service microservice that provides gRPC APIs for prediction management, validates predictions against contest rules and deadlines, stores predictions with flexible data structures, and integrates with existing services for authentication and contest validation.

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: Prediction Service (new), Contest Service (integration), API Gateway (routing)
**Dependencies**: GORM, gRPC, JWT authentication, Contest Service client

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/contest-service/internal/models/contest.go` (lines 1-150) - Why: Contest model structure and validation patterns to follow
- `backend/user-service/internal/service/user_service.go` (lines 1-200) - Why: gRPC service implementation patterns and authentication handling
- `backend/contest-service/internal/service/contest_service.go` (lines 1-300) - Why: Service layer patterns and repository integration
- `backend/shared/auth/interceptors.go` (lines 1-50) - Why: JWT authentication context extraction patterns
- `backend/contest-service/internal/repository/contest_repository.go` (lines 1-200) - Why: Repository interface and GORM patterns
- `backend/proto/contest.proto` (lines 1-150) - Why: Protocol buffer definition patterns and HTTP annotations
- `backend/proto/common.proto` (lines 1-40) - Why: Common response structures and error codes

### New Files to Create

- `backend/prediction-service/cmd/main.go` - Main entry point for Prediction Service
- `backend/prediction-service/internal/config/config.go` - Configuration management
- `backend/prediction-service/internal/models/prediction.go` - Prediction data model with GORM
- `backend/prediction-service/internal/models/event.go` - Sports event data model
- `backend/prediction-service/internal/repository/prediction_repository.go` - Prediction repository interface and implementation
- `backend/prediction-service/internal/repository/event_repository.go` - Event repository interface and implementation
- `backend/prediction-service/internal/service/prediction_service.go` - gRPC service implementation
- `backend/prediction-service/internal/clients/contest_client.go` - Contest service gRPC client
- `backend/prediction-service/go.mod` - Go module definition
- `backend/prediction-service/Dockerfile` - Container configuration
- `backend/proto/prediction.proto` - Protocol buffer definitions with HTTP annotations
- `tests/prediction-service/prediction_test.go` - Unit tests
- `tests/prediction-service/integration_test.go` - Integration tests

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [GORM Documentation](https://gorm.io/docs/)
  - Specific section: Models and Associations
  - Why: Required for implementing prediction and event models with proper relationships
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
  - Specific section: Service Implementation
  - Why: Shows proper gRPC service patterns and error handling
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers/docs/proto3)
  - Specific section: Message Types and Services
  - Why: Required for defining prediction service API

### Patterns to Follow

**Model Validation Pattern:**
```go
func (p *Prediction) BeforeCreate(tx *gorm.DB) error {
    if err := p.ValidateContestID(); err != nil {
        return err
    }
    // Additional validations...
    return nil
}
```

**gRPC Service Pattern:**
```go
func (s *PredictionService) SubmitPrediction(ctx context.Context, req *pb.SubmitPredictionRequest) (*pb.SubmitPredictionResponse, error) {
    userID, ok := auth.GetUserIDFromContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "user not authenticated")
    }
    // Implementation...
}
```

**Repository Interface Pattern:**
```go
type PredictionRepositoryInterface interface {
    Create(prediction *models.Prediction) error
    GetByID(id uint) (*models.Prediction, error)
    GetByUserAndContest(userID, contestID uint) ([]*models.Prediction, error)
    Update(prediction *models.Prediction) error
}
```

**Error Handling Pattern:**
```go
return &pb.SubmitPredictionResponse{
    Response: &common.Response{
        Success: false,
        Message: err.Error(),
        Code:    int32(common.ErrorCode_INVALID_ARGUMENT),
        Timestamp: timestamppb.Now(),
    },
}, nil
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up the basic Prediction Service structure with configuration, models, and repository patterns.

**Tasks:**
- Create Go module and directory structure following existing service patterns
- Implement configuration management with environment variables
- Define Protocol Buffer schemas for prediction operations
- Create database models for predictions and events

### Phase 2: Core Implementation

Implement the main prediction service logic with gRPC endpoints and business rules.

**Tasks:**
- Implement repository layer with GORM integration
- Create gRPC service with prediction CRUD operations
- Add contest service client for validation
- Implement prediction validation and deadline enforcement

### Phase 3: Integration

Integrate with existing services and add HTTP endpoints through API Gateway.

**Tasks:**
- Update API Gateway to route prediction endpoints
- Add JWT authentication integration
- Implement cross-service communication
- Update Docker Compose configuration

### Phase 4: Testing & Validation

Create comprehensive tests and validate the complete prediction workflow.

**Tasks:**
- Implement unit tests for all components
- Create integration tests for prediction workflows
- Add edge case tests for validation rules
- Validate against acceptance criteria

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/prediction-service/go.mod

- **IMPLEMENT**: Go module with gRPC and GORM dependencies
- **PATTERN**: Mirror backend/contest-service/go.mod structure
- **IMPORTS**: gRPC, GORM, shared module, protobuf
- **GOTCHA**: Use replace directive for shared module
- **VALIDATE**: `cd backend/prediction-service && go mod tidy`

### CREATE backend/prediction-service/internal/config/config.go

- **IMPLEMENT**: Configuration struct with database, gRPC, and service endpoints
- **PATTERN**: Mirror backend/contest-service/internal/config/config.go structure
- **IMPORTS**: os, time packages
- **GOTCHA**: Use same environment variable naming convention
- **VALIDATE**: `go build ./internal/config`

### CREATE backend/proto/prediction.proto

- **IMPLEMENT**: Protocol buffer definitions for prediction operations
- **PATTERN**: Follow backend/proto/contest.proto structure with HTTP annotations
- **IMPORTS**: google/protobuf/timestamp.proto, google/api/annotations.proto, common.proto
- **GOTCHA**: Include proper HTTP method mappings for REST endpoints
- **VALIDATE**: `protoc --proto_path=backend/proto --go_out=backend/shared backend/proto/prediction.proto`

### CREATE backend/prediction-service/internal/models/prediction.go

- **IMPLEMENT**: Prediction model with GORM tags and validation
- **PATTERN**: Mirror backend/contest-service/internal/models/contest.go validation patterns
- **IMPORTS**: gorm.io/gorm, time, errors packages
- **GOTCHA**: Use JSON field for flexible prediction data, add proper indexes
- **VALIDATE**: `go build ./internal/models`

### CREATE backend/prediction-service/internal/models/event.go

- **IMPLEMENT**: Sports event model for prediction targets
- **PATTERN**: Follow Contest model structure with validation hooks
- **IMPORTS**: gorm.io/gorm, time packages
- **GOTCHA**: Include sport type, teams, and event metadata
- **VALIDATE**: `go build ./internal/models`

### CREATE backend/prediction-service/internal/repository/prediction_repository.go

- **IMPLEMENT**: Repository interface and GORM implementation
- **PATTERN**: Mirror backend/contest-service/internal/repository/contest_repository.go structure
- **IMPORTS**: gorm.io/gorm, models package
- **GOTCHA**: Include methods for user-contest queries and deadline validation
- **VALIDATE**: `go build ./internal/repository`

### CREATE backend/prediction-service/internal/repository/event_repository.go

- **IMPLEMENT**: Event repository with CRUD operations
- **PATTERN**: Follow prediction repository interface pattern
- **IMPORTS**: gorm.io/gorm, models package
- **GOTCHA**: Add methods for event lookup by sport type and date range
- **VALIDATE**: `go build ./internal/repository`

### CREATE backend/prediction-service/internal/clients/contest_client.go

- **IMPLEMENT**: gRPC client for contest service validation
- **PATTERN**: Create client connection with proper error handling
- **IMPORTS**: google.golang.org/grpc, contest proto package
- **GOTCHA**: Use service discovery endpoint from config
- **VALIDATE**: `go build ./internal/clients`

### CREATE backend/prediction-service/internal/service/prediction_service.go

- **IMPLEMENT**: Main gRPC service with prediction operations
- **PATTERN**: Mirror backend/user-service/internal/service/user_service.go structure
- **IMPORTS**: context, gRPC packages, auth package, proto packages
- **GOTCHA**: Extract user ID from JWT context, validate contest participation
- **VALIDATE**: `go build ./internal/service`

### CREATE backend/prediction-service/cmd/main.go

- **IMPLEMENT**: Main entry point with gRPC server and graceful shutdown
- **PATTERN**: Mirror backend/contest-service/cmd/main.go structure
- **IMPORTS**: net, gRPC, signal packages, internal packages
- **GOTCHA**: Include JWT interceptor and database migration
- **VALIDATE**: `go build ./cmd`

### CREATE backend/prediction-service/Dockerfile

- **IMPLEMENT**: Multi-stage Docker build
- **PATTERN**: Mirror backend/contest-service/Dockerfile
- **IMPORTS**: golang:1.21-alpine base image
- **GOTCHA**: Copy shared module, expose correct port
- **VALIDATE**: `docker build -t prediction-service backend/prediction-service`

### UPDATE backend/go.work

- **IMPLEMENT**: Add prediction-service to workspace
- **PATTERN**: Add to existing use directive
- **IMPORTS**: None
- **GOTCHA**: Maintain alphabetical order
- **VALIDATE**: `cd backend && go work sync`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Add prediction service configuration
- **PATTERN**: Follow contest-service configuration structure
- **IMPORTS**: None
- **GOTCHA**: Use correct port and environment variables
- **VALIDATE**: `docker-compose config`

### UPDATE backend/proto/prediction.proto

- **IMPLEMENT**: Add HTTP annotations for API Gateway integration
- **PATTERN**: Follow contest.proto HTTP annotation patterns
- **IMPORTS**: Already included
- **GOTCHA**: Use proper REST paths like /v1/predictions
- **VALIDATE**: `buf generate`

### UPDATE backend/api-gateway/internal/gateway/gateway.go

- **IMPLEMENT**: Register prediction service handler
- **PATTERN**: Follow existing service registration pattern
- **IMPORTS**: prediction proto package
- **GOTCHA**: Add to service registration list with proper endpoint
- **VALIDATE**: `go build ./internal/gateway`

### CREATE tests/prediction-service/prediction_test.go

- **IMPLEMENT**: Unit tests for prediction service operations
- **PATTERN**: Follow tests/contest-service/contest_test.go structure
- **IMPORTS**: testing, testify packages
- **GOTCHA**: Mock contest service client for isolated testing
- **VALIDATE**: `go test ./tests/prediction-service`

### CREATE tests/prediction-service/integration_test.go

- **IMPLEMENT**: Integration tests for complete prediction workflow
- **PATTERN**: Follow existing integration test patterns
- **IMPORTS**: testing, database packages
- **GOTCHA**: Use test database and proper cleanup
- **VALIDATE**: `go test ./tests/prediction-service -tags=integration`

---

## TESTING STRATEGY

### Unit Tests

Test prediction service components in isolation:
- Prediction model validation with various input scenarios
- Repository operations with mocked database
- Service layer business logic with mocked dependencies
- Contest client integration with mocked responses

### Integration Tests

Test complete prediction workflow:
- End-to-end prediction submission and retrieval
- Contest validation and deadline enforcement
- Cross-service communication with contest service
- Database operations with real PostgreSQL instance

### Edge Cases

- Prediction submission after contest deadline
- Duplicate predictions for same event
- Invalid contest ID or user authentication
- Malformed prediction data and validation errors
- Contest service unavailability scenarios

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
cd backend/prediction-service && go fmt ./...
cd backend/prediction-service && go vet ./...
cd backend && go work sync
```

### Level 2: Unit Tests

```bash
cd backend/prediction-service && go test ./internal/...
cd tests && go test ./prediction-service/...
```

### Level 3: Integration Tests

```bash
make docker-up
cd tests && go test ./prediction-service/... -tags=integration
```

### Level 4: Manual Validation

```bash
# Test prediction submission
curl -X POST http://localhost:8080/v1/predictions \
  -H "Authorization: Bearer <token>" \
  -d '{"contest_id":1,"event_id":1,"prediction_data":"{\"outcome\":\"home_win\"}"}'

# Test prediction retrieval
curl -X GET http://localhost:8080/v1/predictions/contest/1 \
  -H "Authorization: Bearer <token>"

# Test prediction validation
curl -X POST http://localhost:8080/v1/predictions \
  -H "Authorization: Bearer <token>" \
  -d '{"contest_id":999,"event_id":1,"prediction_data":"invalid"}'
```

### Level 5: Docker Validation

```bash
make docker-services
docker-compose ps
docker-compose logs prediction-service
curl -X GET http://localhost:8080/v1/predictions/health
```

---

## ACCEPTANCE CRITERIA

- [ ] Prediction Service serves gRPC endpoints on configured port
- [ ] All prediction operations accessible via REST API through gateway
- [ ] JWT authentication works for all protected endpoints
- [ ] Contest validation prevents predictions for invalid/expired contests
- [ ] Prediction data stored with flexible JSON structure
- [ ] Deadline enforcement prevents late predictions
- [ ] User can retrieve their predictions for specific contests
- [ ] Cross-service communication with Contest Service works
- [ ] Database models include proper indexes and relationships
- [ ] Error responses properly formatted with gRPC status codes
- [ ] Docker container builds and runs successfully
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ requirement
- [ ] Integration tests verify end-to-end workflows

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms feature works
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability
- [ ] Docker Compose integration verified
- [ ] API Gateway routing tested and working

---

## NOTES

### Design Decisions

**Flexible Prediction Data**: Using JSON field for prediction_data allows support for different sports and prediction types without schema changes.

**Contest Integration**: gRPC client to Contest Service ensures real-time validation of contest status and user participation.

**Deadline Enforcement**: Server-side validation prevents late predictions and maintains contest integrity.

**Repository Pattern**: Consistent with existing services for maintainability and testability.

### Performance Considerations

- Database indexes on (contest_id, user_id) for fast user prediction queries
- JSON field indexing for common prediction data queries
- Connection pooling for Contest Service client
- Proper error handling to prevent service cascading failures

### Security Considerations

- JWT authentication for all prediction operations
- Contest participation validation to prevent unauthorized predictions
- Input validation for prediction data to prevent injection attacks
- Rate limiting consideration for prediction submission endpoints

### Scalability Design

- Stateless service design for horizontal scaling
- Database partitioning consideration for large prediction volumes
- Caching strategy for frequently accessed contest data
- Async processing consideration for prediction scoring updates

# Feature: Contest Service Implementation

Система создания и управления конкурсами спортивных прогнозов с поддержкой настраиваемых правил, видов спорта и систем подсчёта очков.

## Feature Description

Contest Service - это микросервис для создания, настройки и управления конкурсами спортивных прогнозов. Сервис позволяет администраторам создавать конкурсы с гибкими правилами, поддерживает различные виды спорта и системы подсчёта очков, обеспечивает управление участниками и отслеживание статуса конкурсов.

## User Story

As a **contest administrator**
I want to **create and manage sports prediction contests with customizable rules**
So that **I can engage users with flexible prediction competitions across different sports**

## Problem Statement

Платформе нужна система управления конкурсами, которая позволит:
- Быстро создавать конкурсы без программирования
- Настраивать правила подсчёта очков для разных видов спорта
- Управлять участниками и их статусами
- Отслеживать жизненный цикл конкурсов
- Интегрироваться с другими микросервисами через gRPC

## Solution Statement

Реализация Contest Service как gRPC микросервиса с:
- CRUD операциями для конкурсов
- Гибкой системой правил и настроек
- Управлением участниками
- Интеграцией с User Service для аутентификации
- PostgreSQL для хранения данных
- Следованием паттернам существующего User Service

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: Contest Service, Database, gRPC API
**Dependencies**: User Service (authentication), PostgreSQL, gRPC, GORM

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/user-service/internal/models/user.go` (lines 1-100) - Why: GORM model patterns, validation methods, hooks
- `backend/user-service/internal/service/user_service.go` (lines 1-50) - Why: gRPC service implementation pattern
- `backend/user-service/internal/config/config.go` (lines 1-70) - Why: Configuration loading and validation patterns
- `backend/user-service/cmd/main.go` (lines 1-50) - Why: Service startup and gRPC server setup
- `backend/proto/common.proto` (lines 1-50) - Why: Common proto definitions and error codes
- `backend/user-service/internal/repository/user_repository.go` - Why: Repository pattern implementation
- `backend/shared/database/database.go` - Why: Database connection patterns
- `backend/shared/auth/auth.go` - Why: Authentication middleware patterns

### New Files to Create

- `backend/proto/contest.proto` - gRPC service definitions for contest operations
- `backend/contest-service/go.mod` - Go module configuration
- `backend/contest-service/cmd/main.go` - Service entry point
- `backend/contest-service/internal/models/contest.go` - Contest data models
- `backend/contest-service/internal/models/participant.go` - Participant data models
- `backend/contest-service/internal/config/config.go` - Service configuration
- `backend/contest-service/internal/repository/contest_repository.go` - Data access layer
- `backend/contest-service/internal/service/contest_service.go` - Business logic layer
- `backend/contest-service/Dockerfile` - Container configuration
- `backend/contest-service/README.md` - Service documentation
- `tests/contest-service/contest_test.go` - Unit tests

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [gRPC Go Documentation](https://grpc.io/docs/languages/go/quickstart/)
  - Specific section: Service implementation patterns
  - Why: Required for implementing gRPC service methods
- [GORM Documentation](https://gorm.io/docs/models.html)
  - Specific section: Model definition and associations
  - Why: Shows proper model relationships and validation
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers/docs/proto3)
  - Specific section: Message definitions and services
  - Why: Required for defining contest API schema

### Patterns to Follow

**Naming Conventions:**
- Files: `snake_case.go`
- Structs: `PascalCase`
- Functions: `PascalCase` (public), `camelCase` (private)
- Proto messages: `PascalCase`

**Error Handling:**
```go
if err != nil {
    return &pb.CreateContestResponse{
        Response: &common.Response{
            Success: false,
            Message: err.Error(),
            Code: int32(common.ErrorCode_INVALID_ARGUMENT),
            Timestamp: timestamppb.Now(),
        },
    }, nil
}
```

**Logging Pattern:**
```go
log.Printf("[INFO] Contest created successfully: %s", contest.ID)
log.Printf("[ERROR] Failed to create contest: %v", err)
```

**GORM Model Pattern:**
```go
type Contest struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Title       string `gorm:"not null" json:"title"`
    Description string `json:"description"`
    gorm.Model
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Создание базовой структуры сервиса и proto определений

**Tasks:**
- Создать proto схему для Contest Service
- Настроить Go модуль и зависимости
- Создать базовую структуру директорий
- Настроить конфигурацию сервиса

### Phase 2: Core Implementation

Реализация основной бизнес-логики и моделей данных

**Tasks:**
- Создать модели данных (Contest, Participant)
- Реализовать repository слой
- Создать gRPC service implementation
- Добавить валидацию и обработку ошибок

### Phase 3: Integration

Интеграция с существующими сервисами и базой данных

**Tasks:**
- Настроить подключение к PostgreSQL
- Добавить аутентификацию через User Service
- Создать миграции базы данных
- Настроить Docker контейнер

### Phase 4: Testing & Validation

Создание тестов и валидация функциональности

**Tasks:**
- Написать unit тесты для всех компонентов
- Создать integration тесты
- Добавить тесты для edge cases
- Валидировать через gRPC клиент

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/proto/contest.proto

- **IMPLEMENT**: gRPC service definitions for contest management
- **PATTERN**: Mirror `backend/proto/common.proto` structure and imports
- **IMPORTS**: `google/protobuf/timestamp.proto`, `common.proto`
- **GOTCHA**: Use consistent naming with existing proto files
- **VALIDATE**: `protoc --go_out=. --go-grpc_out=. backend/proto/contest.proto`

### CREATE backend/contest-service/go.mod

- **IMPLEMENT**: Go module with required dependencies
- **PATTERN**: Mirror `backend/user-service/go.mod` structure
- **IMPORTS**: gRPC, GORM, PostgreSQL driver, shared modules
- **GOTCHA**: Use same versions as user-service for consistency
- **VALIDATE**: `cd backend/contest-service && go mod tidy`

### CREATE backend/contest-service/internal/models/contest.go

- **IMPLEMENT**: Contest model with GORM annotations and validation
- **PATTERN**: Mirror `backend/user-service/internal/models/user.go:1-100`
- **IMPORTS**: GORM, validation helpers, time package
- **GOTCHA**: Include proper JSON tags and GORM relationships
- **VALIDATE**: `cd backend/contest-service && go build ./internal/models`

### CREATE backend/contest-service/internal/models/participant.go

- **IMPLEMENT**: Participant model for contest membership
- **PATTERN**: Follow Contest model structure with foreign keys
- **IMPORTS**: GORM, time package
- **GOTCHA**: Proper relationship setup with Contest and User models
- **VALIDATE**: `cd backend/contest-service && go build ./internal/models`

### CREATE backend/contest-service/internal/config/config.go

- **IMPLEMENT**: Configuration loading from environment variables
- **PATTERN**: Mirror `backend/user-service/internal/config/config.go:1-70`
- **IMPORTS**: os, time, strconv packages
- **GOTCHA**: Use CONTEST_SERVICE_PORT instead of USER_SERVICE_PORT
- **VALIDATE**: `cd backend/contest-service && go build ./internal/config`

### CREATE backend/contest-service/internal/repository/contest_repository.go

- **IMPLEMENT**: Data access layer with CRUD operations
- **PATTERN**: Mirror user-service repository structure
- **IMPORTS**: GORM, models, context
- **GOTCHA**: Include proper error handling and transaction support
- **VALIDATE**: `cd backend/contest-service && go build ./internal/repository`

### CREATE backend/contest-service/internal/service/contest_service.go

- **IMPLEMENT**: gRPC service implementation with business logic
- **PATTERN**: Mirror `backend/user-service/internal/service/user_service.go:1-50`
- **IMPORTS**: gRPC, proto definitions, repository, common
- **GOTCHA**: Proper error response formatting with common.Response
- **VALIDATE**: `cd backend/contest-service && go build ./internal/service`

### CREATE backend/contest-service/cmd/main.go

- **IMPLEMENT**: Service entry point with gRPC server setup
- **PATTERN**: Mirror `backend/user-service/cmd/main.go:1-50`
- **IMPORTS**: gRPC, service, config, database
- **GOTCHA**: Use correct port and service registration
- **VALIDATE**: `cd backend/contest-service && go build ./cmd`

### UPDATE backend/go.work

- **IMPLEMENT**: Add contest-service to Go workspace
- **PATTERN**: Add `./contest-service` to use directive
- **IMPORTS**: None
- **GOTCHA**: Maintain alphabetical order of services
- **VALIDATE**: `cd backend && go work sync`

### CREATE backend/contest-service/Dockerfile

- **IMPLEMENT**: Multi-stage Docker build for contest service
- **PATTERN**: Mirror `backend/user-service/Dockerfile:1-30`
- **IMPORTS**: None
- **GOTCHA**: Update binary name and port exposure
- **VALIDATE**: `docker build -t contest-service backend/contest-service`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Add contest-service to Docker Compose
- **PATTERN**: Mirror user-service configuration
- **IMPORTS**: None
- **GOTCHA**: Use unique port (8085) and proper environment variables
- **VALIDATE**: `docker-compose config`

### CREATE backend/contest-service/README.md

- **IMPLEMENT**: Service documentation with API examples
- **PATTERN**: Mirror `backend/user-service/README.md` structure
- **IMPORTS**: None
- **GOTCHA**: Update service-specific information and port numbers
- **VALIDATE**: Manual review for completeness

### UPDATE .env.example

- **IMPLEMENT**: Add contest service environment variables
- **PATTERN**: Add CONTEST_SERVICE_* variables
- **IMPORTS**: None
- **GOTCHA**: Use port 8085 for contest service
- **VALIDATE**: Check all required variables are documented

---

## TESTING STRATEGY

### Unit Tests

Тестирование каждого компонента изолированно с использованием Go testing framework и testify для assertions.

**Scope:**
- Model validation methods
- Repository CRUD operations (with test database)
- Service business logic (with mocked dependencies)
- Configuration loading and validation

### Integration Tests

Тестирование взаимодействия между компонентами с реальной базой данных.

**Scope:**
- gRPC service endpoints with authentication
- Database operations with transactions
- Error handling and edge cases

### Edge Cases

**Specific edge cases to test:**
- Contest creation with invalid data
- Participant management with non-existent users
- Concurrent access to contest operations
- Database connection failures
- Authentication token validation

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Go formatting and linting
cd backend/contest-service && go fmt ./...
cd backend/contest-service && go vet ./...
cd backend && go work sync
```

### Level 2: Unit Tests

```bash
# Run unit tests with coverage
cd backend/contest-service && go test ./... -v -cover
cd backend/contest-service && go test ./... -race
```

### Level 3: Integration Tests

```bash
# Build and test service
cd backend/contest-service && go build ./cmd
docker-compose up -d postgres redis
cd backend/contest-service && go test ./... -tags=integration
```

### Level 4: Manual Validation

```bash
# Start service and test gRPC endpoints
docker-compose up contest-service
grpcurl -plaintext -d '{"title":"Test Contest","description":"Test Description"}' localhost:8085 contest.ContestService/CreateContest
grpcurl -plaintext localhost:8085 contest.ContestService/ListContests
```

### Level 5: Additional Validation

```bash
# Protocol buffer validation
protoc --go_out=. --go-grpc_out=. backend/proto/contest.proto
# Docker build validation
docker build -t contest-service backend/contest-service
```

---

## ACCEPTANCE CRITERIA

- [ ] Contest Service implements full CRUD operations for contests
- [ ] gRPC API follows existing service patterns and conventions
- [ ] Database models support flexible contest rules and participant management
- [ ] Authentication integration with User Service works correctly
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage exceeds 80% for all components
- [ ] Integration tests verify end-to-end contest workflows
- [ ] Docker containerization works with existing compose setup
- [ ] Service follows Go and gRPC best practices
- [ ] Error handling provides meaningful responses
- [ ] Configuration supports environment-based deployment
- [ ] Documentation includes API usage examples

---

## COMPLETION CHECKLIST

- [ ] All 13 tasks completed in dependency order
- [ ] Each task validation passed immediately after implementation
- [ ] All 5 levels of validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No Go linting or formatting errors
- [ ] gRPC service responds correctly to test calls
- [ ] Docker Compose integration works without issues
- [ ] All acceptance criteria verified and met
- [ ] Code follows established patterns from User Service
- [ ] Service documentation is complete and accurate

---

## NOTES

**Design Decisions:**
- Contest Service uses port 8085 to avoid conflicts with User Service (8084)
- Models support flexible rule configuration through JSON fields
- Participant management allows for different user roles (admin, participant)
- Authentication leverages existing User Service JWT validation

**Performance Considerations:**
- Database indexes on frequently queried fields (contest_id, user_id)
- Pagination support for large contest lists
- Efficient participant lookup with proper foreign key relationships

**Security Considerations:**
- All endpoints require valid JWT authentication
- Contest ownership validation for modification operations
- Input validation prevents SQL injection and data corruption
- Proper error messages without sensitive information exposure

**Future Extensibility:**
- Contest rules stored as JSON for flexible configuration
- Plugin architecture support for custom scoring systems
- Multi-language support through internationalization fields
- Integration points for external sports data APIs

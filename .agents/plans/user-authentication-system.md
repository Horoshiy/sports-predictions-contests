# Feature: User Authentication System

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement a comprehensive JWT-based authentication system for the Sports Prediction Contests platform. This system will provide secure user registration, login, profile management, and token-based authentication across all microservices. The authentication service will serve as the foundation for all user-related functionality in the platform.

## User Story

As a sports prediction platform user
I want to securely register, login, and manage my account
So that I can participate in prediction contests with personalized experience and data protection

## Problem Statement

The platform currently lacks any user authentication mechanism, making it impossible to:
- Identify users across different services
- Secure API endpoints and protect user data
- Implement personalized features like contest participation
- Track user statistics and achievements
- Ensure data privacy and security compliance

## Solution Statement

Implement a JWT-based authentication system using Go gRPC microservices with:
- Secure user registration and login with bcrypt password hashing
- JWT token generation and validation with configurable expiration
- gRPC interceptors for automatic authentication across services
- PostgreSQL user storage with GORM ORM
- Redis-based token blacklisting for logout functionality
- Profile management and user information retrieval

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: user-service, api-gateway, shared libraries
**Dependencies**: golang-jwt/jwt, GORM, bcrypt, Redis client

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/proto/common.proto` (lines 45-55) - Why: Contains error codes including UNAUTHENTICATED that we'll use
- `backend/shared/go.mod` - Why: Shows existing gRPC dependencies and module structure
- `backend/go.work` - Why: Shows workspace structure for user-service integration
- `docker-compose.yml` (lines 5-25) - Why: PostgreSQL and Redis configuration for authentication
- `.env.example` (lines 20-22) - Why: JWT configuration variables and user service port

### New Files to Create

- `backend/proto/user.proto` - gRPC service definitions for authentication
- `backend/user-service/go.mod` - User service module configuration
- `backend/user-service/cmd/main.go` - Service entry point
- `backend/user-service/internal/models/user.go` - User data model with GORM
- `backend/user-service/internal/repository/user_repository.go` - Database operations
- `backend/user-service/internal/service/auth_service.go` - Authentication business logic
- `backend/user-service/internal/service/user_service.go` - gRPC service implementation
- `backend/user-service/internal/config/config.go` - Configuration management
- `backend/shared/auth/jwt.go` - JWT utilities for all services
- `backend/shared/auth/interceptors.go` - gRPC authentication interceptors
- `backend/shared/database/connection.go` - Database connection utilities
- `tests/backend/user-service/auth_test.go` - Authentication unit tests

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [golang-jwt/jwt v5 Documentation](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)
  - Specific section: Token creation and validation
  - Why: Required for implementing secure JWT tokens
- [gRPC Go Authentication Guide](https://grpc.io/docs/guides/auth/)
  - Specific section: Interceptors and metadata
  - Why: Shows proper gRPC authentication patterns
- [GORM Documentation](https://gorm.io/docs/)
  - Specific section: Models and associations
  - Why: Required for PostgreSQL user model implementation
- [bcrypt Package](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
  - Specific section: Password hashing
  - Why: Secure password storage implementation

### Patterns to Follow

**Module Structure Pattern:**
```go
// From backend/shared/go.mod
module github.com/sports-prediction-contests/shared
go 1.21
```

**gRPC Service Pattern:**
```go
// From backend/proto/common.proto
service HealthService {
  rpc Check(google.protobuf.Empty) returns (Response);
}
```

**Error Handling Pattern:**
```go
// From backend/proto/common.proto
enum ErrorCode {
  UNAUTHENTICATED = 5;
  PERMISSION_DENIED = 4;
}
```

**Environment Configuration Pattern:**
```bash
# From .env.example
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION=24h
USER_SERVICE_PORT=8084
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up Protocol Buffers definitions, Go modules, and basic project structure for the user authentication service.

**Tasks:**
- Define gRPC service contracts for authentication operations
- Create user-service Go module with proper dependencies
- Set up shared authentication utilities for cross-service use
- Configure database connection utilities

### Phase 2: Core Implementation

Implement the core authentication logic including user models, JWT handling, and password security.

**Tasks:**
- Create User model with GORM annotations and password hashing
- Implement JWT token generation and validation utilities
- Build user repository for database operations
- Create authentication service with login/register logic

### Phase 3: gRPC Service Integration

Build the gRPC service implementation and authentication interceptors for cross-service security.

**Tasks:**
- Implement UserService gRPC server with all endpoints
- Create JWT authentication interceptors for gRPC
- Set up service configuration and environment handling
- Build main service entry point with proper initialization

### Phase 4: Testing & Validation

Comprehensive testing of authentication flows and security measures.

**Tasks:**
- Implement unit tests for authentication logic
- Create integration tests for gRPC endpoints
- Test JWT token lifecycle and validation
- Validate password security and error handling

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/proto/user.proto

- **IMPLEMENT**: Complete gRPC service definition for user authentication
- **PATTERN**: Mirror service structure from `backend/proto/common.proto`
- **IMPORTS**: Use common.proto Response and error patterns
- **GOTCHA**: Include proper go_package option for code generation
- **VALIDATE**: `protoc --proto_path=backend/proto --go_out=backend/shared --go-grpc_out=backend/shared backend/proto/user.proto`

### CREATE backend/user-service/go.mod

- **IMPLEMENT**: Go module with authentication dependencies
- **PATTERN**: Follow module structure from `backend/shared/go.mod`
- **IMPORTS**: golang-jwt/jwt/v5, gorm.io/gorm, gorm.io/driver/postgres, golang.org/x/crypto/bcrypt
- **GOTCHA**: Use go 1.21 to match workspace requirement
- **VALIDATE**: `cd backend/user-service && go mod tidy`

### CREATE backend/shared/auth/jwt.go

- **IMPLEMENT**: JWT token generation and validation utilities
- **PATTERN**: Use research findings for golang-jwt/jwt v5 implementation
- **IMPORTS**: github.com/golang-jwt/jwt/v5, time, errors
- **GOTCHA**: Use RegisteredClaims for standard JWT fields
- **VALIDATE**: `cd backend/shared && go build ./auth`

### CREATE backend/shared/auth/interceptors.go

- **IMPLEMENT**: gRPC unary interceptor for JWT authentication
- **PATTERN**: Follow gRPC interceptor pattern from research
- **IMPORTS**: google.golang.org/grpc, context, metadata, status
- **GOTCHA**: Skip authentication for Login/Register methods
- **VALIDATE**: `cd backend/shared && go build ./auth`

### CREATE backend/shared/database/connection.go

- **IMPLEMENT**: PostgreSQL connection utility with GORM
- **PATTERN**: Use environment variables from `.env.example`
- **IMPORTS**: gorm.io/gorm, gorm.io/driver/postgres, os
- **GOTCHA**: Handle connection errors gracefully
- **VALIDATE**: `cd backend/shared && go build ./database`

### CREATE backend/user-service/internal/models/user.go

- **IMPLEMENT**: User model with GORM annotations and password methods
- **PATTERN**: Follow GORM model conventions with gorm.Model embedding
- **IMPORTS**: gorm.io/gorm, golang.org/x/crypto/bcrypt
- **GOTCHA**: Never expose password field in JSON, use json:"-" tag
- **VALIDATE**: `cd backend/user-service && go build ./internal/models`

### CREATE backend/user-service/internal/repository/user_repository.go

- **IMPLEMENT**: User repository with CRUD operations
- **PATTERN**: Repository pattern with interface and struct implementation
- **IMPORTS**: gorm.io/gorm, models package
- **GOTCHA**: Handle GORM errors properly, especially record not found
- **VALIDATE**: `cd backend/user-service && go build ./internal/repository`

### CREATE backend/user-service/internal/service/auth_service.go

- **IMPLEMENT**: Authentication business logic service
- **PATTERN**: Service layer pattern with dependency injection
- **IMPORTS**: repository, models, shared/auth packages
- **GOTCHA**: Validate email format and password strength
- **VALIDATE**: `cd backend/user-service && go build ./internal/service`

### CREATE backend/user-service/internal/service/user_service.go

- **IMPLEMENT**: gRPC UserService server implementation
- **PATTERN**: Implement all methods from user.proto service definition
- **IMPORTS**: context, generated proto package, auth_service
- **GOTCHA**: Convert between proto messages and internal models
- **VALIDATE**: `cd backend/user-service && go build ./internal/service`

### CREATE backend/user-service/internal/config/config.go

- **IMPLEMENT**: Configuration struct with environment variable loading
- **PATTERN**: Use struct tags for environment variable mapping
- **IMPORTS**: os, strconv for type conversions
- **GOTCHA**: Provide sensible defaults for development
- **VALIDATE**: `cd backend/user-service && go build ./internal/config`

### CREATE backend/user-service/cmd/main.go

- **IMPLEMENT**: Service main entry point with gRPC server setup
- **PATTERN**: Standard Go service main with graceful shutdown
- **IMPORTS**: net, grpc, config, service, database packages
- **GOTCHA**: Register authentication interceptor and handle signals
- **VALIDATE**: `cd backend/user-service && go build ./cmd`

### UPDATE backend/go.work

- **IMPLEMENT**: Add user-service to workspace
- **PATTERN**: Add "./user-service" to existing use block
- **IMPORTS**: No imports needed
- **GOTCHA**: Maintain existing workspace structure
- **VALIDATE**: `cd backend && go work sync`

### CREATE tests/backend/user-service/auth_test.go

- **IMPLEMENT**: Unit tests for authentication flows
- **PATTERN**: Table-driven tests with setup/teardown
- **IMPORTS**: testing, testify/assert, service packages
- **GOTCHA**: Use test database or mocks for isolation
- **VALIDATE**: `cd backend/user-service && go test ./...`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Add user-service container configuration
- **PATTERN**: Mirror existing service configuration structure
- **IMPORTS**: No imports needed
- **GOTCHA**: Use correct port 8084 and environment variables
- **VALIDATE**: `docker-compose config`

---

## TESTING STRATEGY

### Unit Tests

**Scope**: Test authentication logic, JWT utilities, and user repository operations in isolation

Design unit tests with fixtures and assertions following Go testing conventions:
- Test JWT token generation and validation with various scenarios
- Test password hashing and verification with bcrypt
- Test user repository CRUD operations with mock database
- Test authentication service business logic with dependency injection

### Integration Tests

**Scope**: Test complete gRPC service endpoints with real database connections

- Test user registration flow end-to-end
- Test login authentication with valid/invalid credentials
- Test profile retrieval with JWT authentication
- Test gRPC interceptor authentication across service calls

### Edge Cases

- Invalid JWT tokens (expired, malformed, wrong signature)
- Duplicate user registration attempts
- Password validation edge cases (empty, too short, special characters)
- Database connection failures and recovery
- Concurrent user operations and race conditions

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Protocol Buffers compilation
cd backend && make proto

# Go module validation
cd backend && go work sync
cd backend/user-service && go mod tidy
cd backend/shared && go mod tidy

# Go build validation
cd backend/user-service && go build ./...
cd backend/shared && go build ./...
```

### Level 2: Unit Tests

```bash
# Run all user service tests
cd backend/user-service && go test ./... -v

# Run shared package tests
cd backend/shared && go test ./... -v

# Test coverage validation
cd backend/user-service && go test ./... -cover
```

### Level 3: Integration Tests

```bash
# Start test environment
make docker-up

# Run integration tests
cd backend/user-service && go test ./... -tags=integration

# Test gRPC service manually
grpcurl -plaintext -d '{"email":"test@example.com","password":"password123","name":"Test User"}' localhost:8084 user.UserService/Register
```

### Level 4: Manual Validation

```bash
# Start user service
cd backend/user-service && go run cmd/main.go

# Test registration endpoint
grpcurl -plaintext -d '{"email":"user@test.com","password":"securepass","name":"Test User"}' localhost:8084 user.UserService/Register

# Test login endpoint
grpcurl -plaintext -d '{"email":"user@test.com","password":"securepass"}' localhost:8084 user.UserService/Login

# Test authenticated profile endpoint (use token from login)
grpcurl -plaintext -H "authorization: Bearer <JWT_TOKEN>" localhost:8084 user.UserService/GetProfile
```

### Level 5: Additional Validation (Optional)

```bash
# Database validation
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "SELECT * FROM users;"

# Redis validation (if token blacklisting implemented)
docker exec -it sports-redis redis-cli KEYS "blacklist:*"

# Security validation
cd backend/user-service && go run cmd/main.go &
sleep 2
curl -X POST http://localhost:8084 -d '{}' # Should fail without proper gRPC
```

---

## ACCEPTANCE CRITERIA

- [ ] User can register with email, password, and name
- [ ] User can login with email and password to receive JWT token
- [ ] JWT tokens are properly signed and include user information
- [ ] Passwords are securely hashed using bcrypt
- [ ] gRPC interceptor validates JWT tokens on protected endpoints
- [ ] User profile can be retrieved using valid JWT token
- [ ] Invalid authentication attempts return proper error codes
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ requirement
- [ ] Integration tests verify end-to-end authentication flows
- [ ] Code follows Go conventions and project patterns
- [ ] No security vulnerabilities in authentication implementation
- [ ] Service integrates properly with existing infrastructure
- [ ] Database migrations create user table correctly
- [ ] Environment configuration works across development/production

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in dependency order
- [ ] Each task validation passed immediately after implementation
- [ ] Protocol Buffers generate Go code successfully
- [ ] All Go modules build without errors
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual gRPC testing confirms all endpoints work
- [ ] JWT token lifecycle works correctly
- [ ] Password security implemented properly
- [ ] Database integration functions correctly
- [ ] Service starts and responds to health checks
- [ ] Acceptance criteria all verified
- [ ] Code reviewed for security and maintainability

---

## NOTES

### Security Considerations
- JWT secret must be loaded from environment variables, never hardcoded
- Use bcrypt cost of 12 for password hashing (balance security/performance)
- Implement token expiration and consider refresh token mechanism
- Use TLS in production for all gRPC communications
- Validate all input data to prevent injection attacks

### Performance Considerations  
- Database connection pooling for concurrent requests
- JWT validation should be fast (avoid database lookups when possible)
- Consider Redis caching for frequently accessed user data
- Use prepared statements (GORM handles this automatically)

### Future Enhancements
- Token blacklisting with Redis for logout functionality
- Refresh token mechanism for extended sessions
- Role-based access control (RBAC) for different user types
- OAuth integration for social login
- Multi-factor authentication (MFA) support
- Account verification via email

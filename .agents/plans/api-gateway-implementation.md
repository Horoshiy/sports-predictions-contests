# Feature: API Gateway Implementation

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement a comprehensive HTTP-to-gRPC API Gateway that serves as the central entry point for the Sports Prediction Contests platform. The gateway will provide RESTful HTTP endpoints that translate to gRPC calls for user-service and contest-service, handle authentication, error formatting, and service discovery.

## User Story

As a frontend developer or external API consumer
I want to interact with the platform through standard HTTP/REST endpoints
So that I can build web applications, mobile apps, and integrations without needing gRPC clients

## Problem Statement

Currently, the platform only exposes gRPC services (user-service on port 8084, contest-service on port 8085) which are not directly accessible from web browsers or standard HTTP clients. Frontend applications need HTTP/REST endpoints to communicate with the backend services.

## Solution Statement

Implement a grpc-gateway based API Gateway that:
- Translates HTTP/REST requests to gRPC calls
- Provides unified authentication and error handling
- Serves as single entry point on port 8080
- Maintains existing JWT authentication flow
- Supports service discovery and load balancing

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: API Gateway, Authentication, Service Communication
**Dependencies**: grpc-gateway, existing user-service, contest-service

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/shared/auth/interceptors.go` (lines 1-50) - Why: JWT authentication patterns and context handling
- `backend/shared/auth/jwt.go` (lines 1-100) - Why: JWT validation and token parsing logic
- `backend/user-service/internal/config/config.go` (lines 1-70) - Why: Configuration patterns for environment variables
- `backend/proto/user.proto` (lines 1-50) - Why: User service gRPC definitions for HTTP mapping
- `backend/proto/contest.proto` (lines 1-100) - Why: Contest service gRPC definitions for HTTP mapping
- `backend/proto/common.proto` (lines 1-40) - Why: Common response structures and error codes
- `docker-compose.yml` (lines 20-35) - Why: API Gateway service configuration and networking

### New Files to Create

- `backend/api-gateway/cmd/main.go` - Main entry point for API Gateway service
- `backend/api-gateway/internal/config/config.go` - Configuration management
- `backend/api-gateway/internal/gateway/gateway.go` - Gateway server implementation
- `backend/api-gateway/internal/middleware/auth.go` - HTTP authentication middleware
- `backend/api-gateway/internal/middleware/logging.go` - Request logging middleware
- `backend/api-gateway/internal/middleware/cors.go` - CORS handling middleware
- `backend/api-gateway/go.mod` - Go module definition
- `backend/api-gateway/Dockerfile` - Container configuration
- `backend/proto/user.proto` - Updated with HTTP annotations
- `backend/proto/contest.proto` - Updated with HTTP annotations
- `buf.gen.yaml` - Protocol buffer generation configuration
- `tests/api-gateway/gateway_test.go` - Integration tests

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [gRPC-Gateway Guide](https://blog.logrocket.com/guide-to-grpc-gateway/)
  - Specific section: HTTP annotations and service registration
  - Why: Required for implementing HTTP-to-gRPC translation
- [gRPC-Gateway GitHub](https://github.com/grpc-ecosystem/grpc-gateway)
  - Specific section: Installation and basic setup
  - Why: Shows proper dependency management and code generation
- [Google API Annotations](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto)
  - Specific section: HTTP rule definitions
  - Why: Required for defining REST endpoints in proto files

### Patterns to Follow

**Configuration Pattern:**
```go
type Config struct {
    Port        string
    UserService string
    ContestService string
    JWTSecret   string
}

func Load() *Config {
    return &Config{
        Port: getEnvOrDefault("API_GATEWAY_PORT", "8080"),
        // ... other fields
    }
}
```

**Error Handling Pattern:**
```go
func ErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
    w.Header().Set("Content-Type", "application/json")
    // Convert gRPC status to HTTP status
}
```

**JWT Middleware Pattern:**
```go
func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract and validate JWT token
            // Add to request context
        })
    }
}
```

**Service Registration Pattern:**
```go
err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userServiceAddr, opts)
err := pb.RegisterContestServiceHandlerFromEndpoint(ctx, mux, contestServiceAddr, opts)
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up the basic API Gateway structure with configuration management and dependency setup.

**Tasks:**
- Create Go module and directory structure
- Set up configuration management following existing patterns
- Install and configure grpc-gateway dependencies
- Create basic Dockerfile following existing service patterns

### Phase 2: Protocol Buffer Enhancement

Add HTTP annotations to existing proto files to define REST endpoints.

**Tasks:**
- Update proto files with google.api.http annotations
- Configure buf.gen.yaml for grpc-gateway code generation
- Generate gateway stubs and validate compilation
- Update shared module with generated code

### Phase 3: Core Gateway Implementation

Implement the main gateway server with service registration and middleware.

**Tasks:**
- Create gateway server with gRPC service registration
- Implement HTTP middleware for authentication, logging, and CORS
- Add error handling and response formatting
- Configure service discovery for user-service and contest-service

### Phase 4: Testing & Integration

Create comprehensive tests and integrate with existing infrastructure.

**Tasks:**
- Implement unit tests for middleware components
- Create integration tests for HTTP-to-gRPC translation
- Update Docker Compose configuration
- Add health check endpoints

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/api-gateway/go.mod

- **IMPLEMENT**: Go module with grpc-gateway dependencies
- **PATTERN**: Mirror user-service/go.mod structure
- **IMPORTS**: grpc-gateway/v2, grpc, shared module
- **GOTCHA**: Use replace directive for shared module
- **VALIDATE**: `cd backend/api-gateway && go mod tidy`

### CREATE backend/api-gateway/internal/config/config.go

- **IMPLEMENT**: Configuration struct with service endpoints and JWT settings
- **PATTERN**: Mirror backend/user-service/internal/config/config.go:1-70
- **IMPORTS**: os, time packages
- **GOTCHA**: Use same environment variable naming convention
- **VALIDATE**: `go build ./internal/config`

### UPDATE backend/proto/user.proto

- **IMPLEMENT**: Add google.api.http annotations for REST endpoints
- **PATTERN**: POST /v1/auth/register, POST /v1/auth/login, GET /v1/users/profile, PUT /v1/users/profile
- **IMPORTS**: google/api/annotations.proto
- **GOTCHA**: Use body: "*" for POST requests, path parameters for GET
- **VALIDATE**: `protoc --proto_path=backend/proto --go_out=backend/shared backend/proto/user.proto`

### UPDATE backend/proto/contest.proto

- **IMPLEMENT**: Add google.api.http annotations for contest endpoints
- **PATTERN**: POST /v1/contests, GET /v1/contests/{id}, PUT /v1/contests/{id}, DELETE /v1/contests/{id}
- **IMPORTS**: google/api/annotations.proto
- **GOTCHA**: Use path parameters for resource IDs, query parameters for filters
- **VALIDATE**: `protoc --proto_path=backend/proto --go_out=backend/shared backend/proto/contest.proto`

### CREATE buf.gen.yaml

- **IMPLEMENT**: Protocol buffer generation configuration for grpc-gateway
- **PATTERN**: Generate both gRPC and gateway stubs
- **IMPORTS**: grpc, grpc-gateway, openapiv2 plugins
- **GOTCHA**: Output to backend/shared directory
- **VALIDATE**: `buf generate`

### CREATE backend/api-gateway/internal/middleware/auth.go

- **IMPLEMENT**: HTTP JWT authentication middleware
- **PATTERN**: Extract Bearer token, validate using shared/auth/jwt.go
- **IMPORTS**: shared/auth, http, context packages
- **GOTCHA**: Skip auth for login/register endpoints, add user context
- **VALIDATE**: `go build ./internal/middleware`

### CREATE backend/api-gateway/internal/middleware/logging.go

- **IMPLEMENT**: HTTP request logging middleware
- **PATTERN**: Log method, path, status code, duration
- **IMPORTS**: log, time, http packages
- **GOTCHA**: Use structured logging format
- **VALIDATE**: `go build ./internal/middleware`

### CREATE backend/api-gateway/internal/middleware/cors.go

- **IMPLEMENT**: CORS handling middleware for web clients
- **PATTERN**: Allow common headers and methods
- **IMPORTS**: http package
- **GOTCHA**: Handle preflight OPTIONS requests
- **VALIDATE**: `go build ./internal/middleware`

### CREATE backend/api-gateway/internal/gateway/gateway.go

- **IMPLEMENT**: Main gateway server with service registration
- **PATTERN**: Create ServeMux, register services, add middleware chain
- **IMPORTS**: grpc-gateway/runtime, generated proto packages
- **GOTCHA**: Use insecure connection for development, add error handler
- **VALIDATE**: `go build ./internal/gateway`

### CREATE backend/api-gateway/cmd/main.go

- **IMPLEMENT**: Main entry point with graceful shutdown
- **PATTERN**: Mirror backend/user-service/cmd/main.go structure
- **IMPORTS**: gateway, config, signal packages
- **GOTCHA**: Listen on configured port, handle SIGINT/SIGTERM
- **VALIDATE**: `go build ./cmd`

### CREATE backend/api-gateway/Dockerfile

- **IMPLEMENT**: Multi-stage Docker build
- **PATTERN**: Mirror backend/user-service/Dockerfile
- **IMPORTS**: golang:1.21-alpine base image
- **GOTCHA**: Copy shared module, expose port 8080
- **VALIDATE**: `docker build -t api-gateway backend/api-gateway`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Remove placeholder API Gateway, add real service
- **PATTERN**: Use build context, add environment variables
- **IMPORTS**: None
- **GOTCHA**: Ensure network connectivity to other services
- **VALIDATE**: `docker-compose config`

### CREATE tests/api-gateway/gateway_test.go

- **IMPLEMENT**: Integration tests for HTTP endpoints
- **PATTERN**: Test auth flow, contest CRUD operations
- **IMPORTS**: testing, http, json packages
- **GOTCHA**: Use test containers for dependencies
- **VALIDATE**: `go test ./tests/api-gateway`

### UPDATE backend/api-gateway/go.mod

- **IMPLEMENT**: Add missing dependencies discovered during implementation
- **PATTERN**: Add testify, httptest for testing
- **IMPORTS**: Required test dependencies
- **GOTCHA**: Keep versions compatible with other services
- **VALIDATE**: `cd backend/api-gateway && go mod tidy && go test ./...`

---

## TESTING STRATEGY

### Unit Tests

Test middleware components in isolation:
- JWT authentication middleware with valid/invalid tokens
- CORS middleware with different origins
- Logging middleware output verification
- Error handler with various gRPC status codes

### Integration Tests

Test complete HTTP-to-gRPC flow:
- User registration and login via HTTP endpoints
- Contest CRUD operations via REST API
- Authentication flow with JWT tokens
- Error responses and status code mapping

### Edge Cases

- Invalid JWT tokens and missing authorization headers
- Malformed JSON requests and validation errors
- Service unavailability and timeout handling
- Concurrent requests and rate limiting

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
cd backend/api-gateway && go fmt ./...
cd backend/api-gateway && go vet ./...
cd backend && go work sync
```

### Level 2: Unit Tests

```bash
cd backend/api-gateway && go test ./internal/...
cd tests && go test ./api-gateway/...
```

### Level 3: Integration Tests

```bash
make docker-up
cd backend/api-gateway && go run cmd/main.go &
curl -X POST http://localhost:8080/v1/auth/register -d '{"email":"test@example.com","password":"password","name":"Test User"}'
curl -X POST http://localhost:8080/v1/auth/login -d '{"email":"test@example.com","password":"password"}'
```

### Level 4: Manual Validation

```bash
# Test user endpoints
curl -X GET http://localhost:8080/v1/users/profile -H "Authorization: Bearer <token>"

# Test contest endpoints
curl -X POST http://localhost:8080/v1/contests -H "Authorization: Bearer <token>" -d '{"title":"Test Contest","sport_type":"football"}'
curl -X GET http://localhost:8080/v1/contests/1

# Test error handling
curl -X GET http://localhost:8080/v1/users/profile
```

### Level 5: Docker Validation

```bash
make docker-services
docker-compose ps
docker-compose logs api-gateway
curl -X GET http://localhost:8080/health
```

---

## ACCEPTANCE CRITERIA

- [ ] API Gateway serves HTTP endpoints on port 8080
- [ ] All user service endpoints accessible via REST API
- [ ] All contest service endpoints accessible via REST API
- [ ] JWT authentication works for protected endpoints
- [ ] Public endpoints (login, register) work without authentication
- [ ] Error responses properly formatted as JSON with HTTP status codes
- [ ] CORS headers allow frontend access
- [ ] Request logging captures method, path, status, duration
- [ ] Health check endpoint returns service status
- [ ] Docker container builds and runs successfully
- [ ] Integration with existing services works without modification
- [ ] All validation commands pass with zero errors
- [ ] Performance meets <200ms response time requirement

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
- [ ] Service discovery and communication tested

---

## NOTES

### Design Decisions

**grpc-gateway vs Custom Proxy**: Chose grpc-gateway for:
- Automatic code generation from proto files
- Built-in HTTP status code mapping
- OpenAPI documentation generation
- Battle-tested in production environments

**Authentication Strategy**: HTTP middleware approach:
- Extracts JWT from Authorization header
- Validates using existing shared/auth package
- Adds user context for downstream gRPC calls
- Maintains compatibility with existing services

**Service Discovery**: Static configuration for development:
- Environment variables for service endpoints
- Docker Compose networking for service resolution
- Extensible to Consul/Kubernetes for production

### Performance Considerations

- Connection pooling to gRPC services
- Middleware ordering for optimal performance
- Structured logging to avoid I/O blocking
- Graceful shutdown to prevent connection leaks

### Security Considerations

- JWT validation using existing secret
- CORS configuration for web client access
- Request size limits to prevent DoS
- Rate limiting consideration for future enhancement

### Extensibility

- Middleware chain allows easy addition of features
- Protocol buffer annotations support new endpoints
- Service registration pattern scales to additional services
- Error handling supports custom status codes and messages

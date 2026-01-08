# Contest Service

Contest Service is a gRPC microservice for managing sports prediction contests. It provides comprehensive contest management functionality including creation, updates, participant management, and flexible rule configuration.

## Features

- **Contest Management**: Create, update, delete, and list contests
- **Participant Management**: Join/leave contests, list participants
- **Flexible Rules**: JSON-based rule configuration for different sports
- **Authentication**: JWT-based authentication integration
- **Pagination**: Efficient pagination for large datasets
- **Status Management**: Contest lifecycle management (draft, active, completed, cancelled)

## API Endpoints

### Contest Operations

- `CreateContest` - Create a new contest
- `UpdateContest` - Update existing contest (creator/admin only)
- `GetContest` - Retrieve contest by ID
- `DeleteContest` - Delete contest (creator only)
- `ListContests` - List contests with pagination and filters

### Participant Operations

- `JoinContest` - Join an active contest
- `LeaveContest` - Leave a contest
- `ListParticipants` - List contest participants with pagination

### Health Check

- `Check` - Service health check

## Configuration

The service is configured via environment variables:

```bash
CONTEST_SERVICE_PORT=8085          # Service port (default: 8085)
DATABASE_URL=postgres://...        # PostgreSQL connection string
JWT_SECRET=your_secret_key         # JWT signing secret
LOG_LEVEL=info                     # Logging level
```

## Database Schema

### Contests Table

- `id` - Primary key
- `title` - Contest title (required, max 200 chars)
- `description` - Contest description (max 1000 chars)
- `sport_type` - Sport type (football, basketball, tennis, etc.)
- `rules` - JSON string for flexible rule configuration
- `status` - Contest status (draft, active, completed, cancelled)
- `start_date` - Contest start date
- `end_date` - Contest end date
- `max_participants` - Maximum participants (0 = unlimited)
- `current_participants` - Current participant count
- `creator_id` - User ID of contest creator
- `created_at`, `updated_at` - Timestamps

### Participants Table

- `id` - Primary key
- `contest_id` - Foreign key to contests table
- `user_id` - Foreign key to users table
- `role` - Participant role (admin, participant)
- `status` - Participant status (active, inactive, banned)
- `joined_at` - Join timestamp
- `created_at`, `updated_at` - Timestamps

## Usage Examples

### Create Contest

```bash
grpcurl -plaintext -d '{
  "title": "Premier League Predictions",
  "description": "Predict Premier League match outcomes",
  "sport_type": "football",
  "rules": "{\"scoring\": {\"exact_score\": 3, \"correct_outcome\": 1}}",
  "start_date": "2024-01-15T00:00:00Z",
  "end_date": "2024-05-15T23:59:59Z",
  "max_participants": 100
}' localhost:8085 contest.ContestService/CreateContest
```

### List Contests

```bash
grpcurl -plaintext -d '{
  "pagination": {"page": 1, "limit": 10},
  "status": "active",
  "sport_type": "football"
}' localhost:8085 contest.ContestService/ListContests
```

### Join Contest

```bash
grpcurl -plaintext -d '{
  "contest_id": 1
}' localhost:8085 contest.ContestService/JoinContest
```

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Protocol Buffers compiler

### Setup

```bash
# Install dependencies
go mod tidy

# Generate proto files (if needed)
protoc --go_out=. --go-grpc_out=. ../proto/contest.proto

# Run service
go run cmd/main.go
```

### Testing

```bash
# Run unit tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run integration tests
go test ./... -tags=integration
```

### Docker

```bash
# Build image
docker build -t contest-service .

# Run with Docker Compose
docker-compose up contest-service
```

## Authentication

All endpoints except health check require JWT authentication. Include the JWT token in the gRPC metadata:

```
authorization: Bearer <jwt_token>
```

## Error Handling

The service returns structured error responses with:

- `success` - Boolean indicating success/failure
- `message` - Human-readable error message
- `code` - Error code from common.ErrorCode enum
- `timestamp` - Response timestamp

## Logging

The service uses structured logging with levels:

- `INFO` - General information
- `ERROR` - Error conditions
- `DEBUG` - Debug information (development only)

Log format: `[LEVEL] Message: details`

## Performance

- Database indexes on frequently queried fields
- Pagination limits (max 100 items per page)
- Connection pooling for database connections
- Efficient participant counting with database aggregation

## Security

- JWT token validation on all authenticated endpoints
- Input validation and sanitization
- SQL injection prevention through GORM
- Contest ownership validation for modifications
- Participant duplicate prevention

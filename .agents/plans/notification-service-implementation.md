# Feature: Notification Service

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

A multi-channel notification microservice for the Sports Prediction Contests platform that delivers real-time alerts via Telegram bots, email (SMTP), and in-app notifications. The service handles event-driven notifications for prediction results, leaderboard changes, contest updates, and match reminders.

## User Story

As a sports prediction contest participant
I want to receive notifications about my predictions, contest updates, and match reminders
So that I stay engaged and never miss important events

## Problem Statement

Users currently have no way to receive alerts about:
- Prediction results and score updates
- Leaderboard position changes
- Contest start/end times
- Upcoming match reminders
- New contest invitations

## Solution Statement

Implement a notification microservice that:
1. Stores in-app notifications in PostgreSQL with read/unread status
2. Sends Telegram messages via Bot API for instant alerts
3. Sends email notifications via SMTP for important updates
4. Processes notifications asynchronously via worker pool
5. Exposes gRPC API for other services to trigger notifications
6. Integrates with API Gateway for frontend access

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: notification-service (new), api-gateway, frontend
**Dependencies**: telegram-bot-api/v5, jordan-wright/email, Redis (queue), PostgreSQL

---

## CONTEXT REFERENCES

### Relevant Codebase Files - READ BEFORE IMPLEMENTING

- `backend/scoring-service/cmd/main.go` - Service initialization pattern with gRPC, health checks, graceful shutdown
- `backend/scoring-service/internal/config/config.go` - Config loading pattern with `getEnvOrDefault()`
- `backend/scoring-service/internal/models/score.go` - GORM model with validation hooks
- `backend/scoring-service/internal/repository/score_repository.go` - Repository interface pattern
- `backend/proto/scoring.proto` - Proto definition pattern with HTTP annotations
- `backend/proto/common.proto` - Common Response wrapper and pagination
- `backend/api-gateway/internal/gateway/gateway.go` - Service registration pattern
- `backend/api-gateway/internal/config/config.go` - Gateway config with service endpoints
- `backend/shared/auth/interceptors.go` - JWT interceptor pattern
- `docker-compose.yml` - Service configuration pattern
- `scripts/init-db.sql` - Database schema pattern

### New Files to Create

```
backend/notification-service/
├── cmd/main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/config.go
    ├── models/
    │   └── notification.go
    ├── repository/
    │   └── notification_repository.go
    ├── service/
    │   └── notification_service.go
    ├── channels/
    │   ├── telegram.go
    │   └── email.go
    └── worker/
        └── worker.go

backend/proto/notification.proto
backend/shared/proto/notification/  (generated)
tests/notification-service/
├── notification_test.go
└── channels_test.go
```

### Relevant Documentation

- [Telegram Bot API](https://core.telegram.org/bots/api) - Bot creation and message sending
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Go library for Telegram
- [jordan-wright/email](https://github.com/jordan-wright/email) - Simple SMTP email library
- [gRPC Go](https://grpc.io/docs/languages/go/) - gRPC service patterns

### Patterns to Follow

**Config Loading:**
```go
func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

**GORM Model with Hooks:**
```go
type Model struct {
    ID uint `gorm:"primaryKey" json:"id"`
    gorm.Model
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
    // validation logic
    return nil
}
```

**Repository Interface:**
```go
type RepositoryInterface interface {
    Create(ctx context.Context, model *Model) error
    GetByID(ctx context.Context, id uint) (*Model, error)
    // ...
}
```

**gRPC Service Registration:**
```go
server := grpc.NewServer(grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))))
pb.RegisterServiceServer(server, serviceImpl)
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation
- Create notification-service directory structure
- Define proto schema for notification API
- Set up Go module and dependencies
- Create config loader

### Phase 2: Core Implementation
- Implement Notification model with GORM
- Create notification repository
- Build Telegram channel sender
- Build Email channel sender
- Implement worker pool for async processing
- Create main notification service

### Phase 3: Integration
- Generate proto code
- Register with API Gateway
- Add to Docker Compose
- Update database schema

### Phase 4: Testing
- Unit tests for models
- Unit tests for channels
- Integration tests

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `backend/proto/notification.proto`

Define gRPC service and messages for notifications.

**IMPLEMENT**: Proto schema with notification types, channels, CRUD operations
**PATTERN**: Mirror `backend/proto/scoring.proto` structure
**IMPORTS**: common.proto, google/protobuf/timestamp.proto, google/api/annotations.proto

```protobuf
syntax = "proto3";
package notification;
option go_package = "github.com/sports-prediction-contests/shared/proto/notification";

// NotificationType enum
enum NotificationType {
  NOTIFICATION_TYPE_UNSPECIFIED = 0;
  PREDICTION_RESULT = 1;
  LEADERBOARD_UPDATE = 2;
  CONTEST_START = 3;
  CONTEST_END = 4;
  MATCH_REMINDER = 5;
  NEW_CONTEST = 6;
}

// NotificationChannel enum
enum NotificationChannel {
  CHANNEL_UNSPECIFIED = 0;
  IN_APP = 1;
  TELEGRAM = 2;
  EMAIL = 3;
}

// Notification message, NotificationPreference message
// SendNotificationRequest/Response, GetNotificationsRequest/Response
// MarkAsReadRequest/Response, UpdatePreferencesRequest/Response
// NotificationService with HTTP annotations
```

**VALIDATE**: `cat backend/proto/notification.proto | head -50`

---

### Task 2: CREATE `backend/notification-service/go.mod`

**IMPLEMENT**: Go module with required dependencies

```go
module github.com/sports-prediction-contests/notification-service

go 1.21

require (
    github.com/sports-prediction-contests/shared v0.0.0
    google.golang.org/grpc v1.59.0
    google.golang.org/protobuf v1.31.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
)

replace github.com/sports-prediction-contests/shared => ../shared
```

**VALIDATE**: `cat backend/notification-service/go.mod`

---

### Task 3: CREATE `backend/notification-service/internal/config/config.go`

**IMPLEMENT**: Configuration struct with env loading
**PATTERN**: Mirror `backend/scoring-service/internal/config/config.go`

Config fields needed:
- Port, JWTSecret, DatabaseURL, RedisURL, LogLevel
- TelegramBotToken, TelegramEnabled
- SMTPHost, SMTPPort, SMTPUser, SMTPPassword, SMTPFrom, EmailEnabled
- WorkerPoolSize

**VALIDATE**: `cat backend/notification-service/internal/config/config.go`

---

### Task 4: CREATE `backend/notification-service/internal/models/notification.go`

**IMPLEMENT**: GORM models for Notification and NotificationPreference
**PATTERN**: Mirror `backend/scoring-service/internal/models/score.go`

Notification fields:
- ID, UserID, Type, Title, Message, Data (JSON), Channel, IsRead, SentAt, ReadAt
- GORM hooks for validation

NotificationPreference fields:
- ID, UserID, Channel, Enabled, TelegramChatID, Email

**VALIDATE**: `cat backend/notification-service/internal/models/notification.go`

---

### Task 5: CREATE `backend/notification-service/internal/repository/notification_repository.go`

**IMPLEMENT**: Repository with CRUD operations
**PATTERN**: Mirror `backend/scoring-service/internal/repository/score_repository.go`

Methods:
- Create, GetByID, Update, Delete
- GetByUser (with pagination, filter by read status)
- MarkAsRead, MarkAllAsRead
- GetUnreadCount
- GetPreferences, UpdatePreferences

**VALIDATE**: `cat backend/notification-service/internal/repository/notification_repository.go`

---

### Task 6: CREATE `backend/notification-service/internal/channels/telegram.go`

**IMPLEMENT**: Telegram Bot API integration
**IMPORTS**: github.com/go-telegram-bot-api/telegram-bot-api/v5

```go
type TelegramChannel struct {
    bot     *tgbotapi.BotAPI
    enabled bool
}

func NewTelegramChannel(token string, enabled bool) (*TelegramChannel, error)
func (t *TelegramChannel) Send(chatID int64, title, message string) error
func (t *TelegramChannel) IsEnabled() bool
```

**GOTCHA**: Handle rate limiting (Telegram allows ~30 messages/second)
**VALIDATE**: `cat backend/notification-service/internal/channels/telegram.go`

---

### Task 7: CREATE `backend/notification-service/internal/channels/email.go`

**IMPLEMENT**: SMTP email sender
**IMPORTS**: github.com/jordan-wright/email, net/smtp

```go
type EmailChannel struct {
    host, port, user, password, from string
    enabled bool
}

func NewEmailChannel(host, port, user, password, from string, enabled bool) *EmailChannel
func (e *EmailChannel) Send(to, subject, body string) error
func (e *EmailChannel) IsEnabled() bool
```

**VALIDATE**: `cat backend/notification-service/internal/channels/email.go`

---

### Task 8: CREATE `backend/notification-service/internal/worker/worker.go`

**IMPLEMENT**: Worker pool for async notification processing

```go
type NotificationJob struct {
    Notification *models.Notification
    Preferences  *models.NotificationPreference
}

type WorkerPool struct {
    jobs     chan NotificationJob
    quit     chan bool
    telegram *channels.TelegramChannel
    email    *channels.EmailChannel
    repo     repository.NotificationRepositoryInterface
}

func NewWorkerPool(size int, telegram, email, repo) *WorkerPool
func (w *WorkerPool) Start()
func (w *WorkerPool) Stop()
func (w *WorkerPool) Submit(job NotificationJob)
func (w *WorkerPool) processJob(job NotificationJob)
```

**VALIDATE**: `cat backend/notification-service/internal/worker/worker.go`

---

### Task 9: CREATE `backend/notification-service/internal/service/notification_service.go`

**IMPLEMENT**: gRPC service implementation
**PATTERN**: Mirror `backend/scoring-service/internal/service/scoring_service.go`

Implement all proto service methods:
- SendNotification - creates notification, submits to worker pool
- GetNotifications - paginated list for user
- GetNotification - single notification by ID
- MarkAsRead, MarkAllAsRead
- GetUnreadCount
- GetPreferences, UpdatePreferences
- Check (health)

**VALIDATE**: `cat backend/notification-service/internal/service/notification_service.go`

---

### Task 10: CREATE `backend/notification-service/cmd/main.go`

**IMPLEMENT**: Service entry point with gRPC server
**PATTERN**: Mirror `backend/scoring-service/cmd/main.go` exactly

Flow:
1. Load config, validate
2. Connect to database, auto-migrate models
3. Initialize channels (Telegram, Email)
4. Initialize repository
5. Initialize worker pool, start workers
6. Create notification service
7. Create gRPC server with JWT interceptor
8. Register service and health server
9. Handle graceful shutdown (stop workers, graceful stop server)
10. Start serving

**VALIDATE**: `cat backend/notification-service/cmd/main.go`

---

### Task 11: CREATE `backend/notification-service/Dockerfile`

**IMPLEMENT**: Multi-stage Docker build
**PATTERN**: Mirror `backend/scoring-service/Dockerfile`

**VALIDATE**: `cat backend/notification-service/Dockerfile`

---

### Task 12: UPDATE `backend/go.work`

**IMPLEMENT**: Add notification-service to Go workspace

Add line: `./notification-service`

**VALIDATE**: `cat backend/go.work`

---

### Task 13: UPDATE `scripts/init-db.sql`

**IMPLEMENT**: Add notifications and notification_preferences tables

```sql
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data TEXT,
    channel VARCHAR(20) NOT NULL DEFAULT 'in_app',
    is_read BOOLEAN DEFAULT false,
    sent_at TIMESTAMP,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications(deleted_at);

CREATE TABLE IF NOT EXISTS notification_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    channel VARCHAR(20) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    telegram_chat_id BIGINT,
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(user_id, channel)
);

CREATE INDEX IF NOT EXISTS idx_notification_preferences_user ON notification_preferences(user_id);
```

**VALIDATE**: `grep -A 5 "notifications" scripts/init-db.sql`

---

### Task 14: UPDATE `backend/api-gateway/internal/config/config.go`

**IMPLEMENT**: Add NotificationService endpoint to Config struct and Load()

Add field: `NotificationService string`
Add in Load(): `NotificationService: getEnvOrDefault("NOTIFICATION_SERVICE_ENDPOINT", "notification-service:8089"),`

**VALIDATE**: `grep -i notification backend/api-gateway/internal/config/config.go`

---

### Task 15: UPDATE `backend/api-gateway/internal/gateway/gateway.go`

**IMPLEMENT**: Register notification service with gateway

Add import: `notificationpb "github.com/sports-prediction-contests/shared/proto/notification"`

Add registration after sports service:
```go
err = notificationpb.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, cfg.NotificationService, opts)
if err != nil {
    return nil, err
}
```

**VALIDATE**: `grep -i notification backend/api-gateway/internal/gateway/gateway.go`

---

### Task 16: UPDATE `docker-compose.yml`

**IMPLEMENT**: Add notification-service container

Add after sports-service:
```yaml
notification-service:
  build:
    context: ./backend/notification-service
    dockerfile: Dockerfile
  container_name: sports-notification-service
  ports:
    - "8089:8089"
  environment:
    - DATABASE_URL=postgres://sports_user:sports_password@postgres:5432/sports_prediction?sslmode=disable
    - REDIS_URL=redis://redis:6379
    - JWT_SECRET=your_jwt_secret_key_here
    - NOTIFICATION_SERVICE_PORT=8089
    - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN:-}
    - TELEGRAM_ENABLED=${TELEGRAM_ENABLED:-false}
    - SMTP_HOST=${SMTP_HOST:-}
    - SMTP_PORT=${SMTP_PORT:-587}
    - SMTP_USER=${SMTP_USER:-}
    - SMTP_PASSWORD=${SMTP_PASSWORD:-}
    - SMTP_FROM=${SMTP_FROM:-}
    - EMAIL_ENABLED=${EMAIL_ENABLED:-false}
    - WORKER_POOL_SIZE=5
    - LOG_LEVEL=info
  depends_on:
    - postgres
    - redis
  networks:
    - sports-network
  profiles:
    - services
```

Also add to api-gateway environment:
`- NOTIFICATION_SERVICE_ENDPOINT=notification-service:8089`

**VALIDATE**: `grep -A 5 notification docker-compose.yml`

---

### Task 17: GENERATE Proto Code

**IMPLEMENT**: Generate Go code from notification.proto

Run: `cd backend && protoc --proto_path=proto --go_out=shared --go-grpc_out=shared --grpc-gateway_out=shared proto/notification.proto`

Or use buf if configured: `buf generate`

**VALIDATE**: `ls -la backend/shared/proto/notification/`

---

### Task 18: CREATE `tests/notification-service/notification_test.go`

**IMPLEMENT**: Unit tests for notification model validation
**PATTERN**: Mirror `tests/scoring-service/scoring_test.go`

Test cases:
- TestNotificationValidation (valid notification, empty title, empty message)
- TestNotificationPreferenceValidation
- TestNotificationBeforeCreate hook

**VALIDATE**: `cat tests/notification-service/notification_test.go`

---

### Task 19: CREATE `tests/notification-service/channels_test.go`

**IMPLEMENT**: Unit tests for channel implementations

Test cases:
- TestTelegramChannelDisabled (returns nil when disabled)
- TestEmailChannelDisabled (returns nil when disabled)
- TestTelegramChannelEnabled (mock test)
- TestEmailChannelEnabled (mock test)

**VALIDATE**: `cat tests/notification-service/channels_test.go`

---

## TESTING STRATEGY

### Unit Tests
- Model validation (notification.go)
- Repository methods (mock DB)
- Channel send methods (mock external APIs)
- Worker pool job processing

### Integration Tests
- Full notification flow: create → worker → channel
- Database operations
- gRPC endpoint responses

### Edge Cases
- Empty notification title/message
- Invalid user ID
- Disabled channels (should skip gracefully)
- Worker pool shutdown during processing
- Database connection failure
- Telegram rate limiting

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build
```bash
cd backend/notification-service && go mod tidy
cd backend/notification-service && go build ./...
cd backend && go work sync
```

### Level 2: Unit Tests
```bash
cd tests/notification-service && go test -v ./...
```

### Level 3: Proto Generation
```bash
cd backend && protoc --proto_path=proto \
  --go_out=shared --go_opt=paths=source_relative \
  --go-grpc_out=shared --go-grpc_opt=paths=source_relative \
  proto/notification.proto
```

### Level 4: Docker Build
```bash
docker-compose build notification-service
```

### Level 5: Integration Test
```bash
docker-compose up -d postgres redis
docker-compose --profile services up -d notification-service
# Test health endpoint
curl http://localhost:8089/health || grpcurl -plaintext localhost:8089 list
```

---

## ACCEPTANCE CRITERIA

- [ ] Notification service starts and registers with API Gateway
- [ ] In-app notifications stored in PostgreSQL with read/unread status
- [ ] Telegram notifications sent when enabled and chat ID configured
- [ ] Email notifications sent when enabled and SMTP configured
- [ ] Worker pool processes notifications asynchronously
- [ ] gRPC API exposes all CRUD operations
- [ ] User preferences stored and respected per channel
- [ ] Health check endpoint responds correctly
- [ ] Graceful shutdown stops workers cleanly
- [ ] All unit tests pass
- [ ] Docker container builds successfully

---

## COMPLETION CHECKLIST

- [ ] All 19 tasks completed in order
- [ ] Proto code generated successfully
- [ ] Go modules sync without errors
- [ ] Docker build succeeds
- [ ] Unit tests pass
- [ ] Service starts and health check responds
- [ ] API Gateway registers notification service

---

## NOTES

**Design Decisions:**
1. Worker pool pattern chosen over message queue (Redis Streams) for simplicity - can upgrade later
2. Telegram and Email channels are optional (disabled by default) - service works with just in-app
3. Notification preferences per channel allow granular user control
4. Data field stores JSON for flexible notification payloads (prediction details, scores, etc.)

**Future Enhancements:**
- Redis Streams for distributed queue
- Push notifications (FCM/APNs)
- Notification templates
- Batch sending optimization
- Rate limiting per user

**Security Considerations:**
- Telegram bot token stored in environment variable
- SMTP credentials stored in environment variables
- JWT authentication required for all endpoints except health

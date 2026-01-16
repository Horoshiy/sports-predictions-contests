# Code Review: Notification Service Implementation

**Date**: 2026-01-16
**Reviewer**: Kiro CLI
**Scope**: Notification service implementation and related changes

## Stats

- Files Modified: 4
- Files Added: 13
- Files Deleted: 0
- New lines: ~1,200
- Deleted lines: 7

---

## Issues Found

### CRITICAL

```
severity: critical
file: backend/notification-service/Dockerfile
line: 8
issue: Invalid COPY path - Docker build will fail
detail: `COPY ../shared ../shared` is invalid in Docker. Docker COPY cannot access files outside the build context. The scoring-service Dockerfile has the same issue but uses a different binary path. This will cause the Docker build to fail.
suggestion: Use a multi-stage build with the build context set to the backend directory, or copy shared module differently. Example fix:
  - Change build context in docker-compose.yml to `./backend`
  - Update Dockerfile paths accordingly
  - Or use go mod vendor approach
```

```
severity: critical
file: backend/notification-service/internal/worker/worker.go
line: 80-81
issue: SentAt update is never persisted to database
detail: The worker sets `n.SentAt = &now` but never calls the repository to save this update. The notification's SentAt field will remain NULL in the database even after successful delivery.
suggestion: Add repository call to update the notification after successful send:
  ctx := context.Background()
  s.repo.Update(ctx, n)
```

### HIGH

```
severity: high
file: backend/notification-service/internal/service/notification_service.go
line: 80-81
issue: Negative offset calculation when page is 0 or not provided
detail: When `req.Pagination.Page` is 0 (default protobuf value), the offset calculation `(0 - 1) * limit` results in a negative offset (-20), which will cause database errors or unexpected behavior.
suggestion: Add validation for page number:
  page := int(req.Pagination.Page)
  if page < 1 {
      page = 1
  }
  offset = (page - 1) * limit
```

```
severity: high
file: backend/notification-service/internal/worker/worker.go
line: 55-58
issue: Jobs channel not drained on shutdown - potential goroutine leak
detail: When Stop() is called, the quit channel is closed but pending jobs in the jobs channel are never processed or drained. Workers exit immediately, leaving jobs unprocessed.
suggestion: Drain the jobs channel before exiting or process remaining jobs:
  func (w *WorkerPool) Stop() {
      close(w.quit)
      close(w.jobs) // Signal no more jobs
      w.wg.Wait()
  }
  // And in worker, handle closed channel
```

```
severity: high
file: backend/notification-service/internal/service/notification_service.go
line: 51
issue: Silently ignoring GetPreference error
detail: `pref, _ := s.repo.GetPreference(...)` ignores the error. If there's a database error (not just "not found"), it will be silently ignored and the notification may fail to send without any indication.
suggestion: Handle the error properly:
  pref, err := s.repo.GetPreference(ctx, uint(req.UserId), channelToString(ch))
  if err != nil {
      log.Printf("Warning: failed to get preference: %v", err)
  }
```

### MEDIUM

```
severity: medium
file: backend/notification-service/internal/channels/telegram.go
line: 32
issue: Markdown injection vulnerability in Telegram messages
detail: User-provided title and message are directly interpolated into Markdown format without escaping. Special characters like `*`, `_`, `[`, etc. could break formatting or be used for injection.
suggestion: Escape Markdown special characters or use HTML parse mode with proper escaping:
  text := fmt.Sprintf("<b>%s</b>\n\n%s", html.EscapeString(title), html.EscapeString(message))
  msg.ParseMode = "HTML"
```

```
severity: medium
file: backend/notification-service/internal/service/notification_service.go
line: 163
issue: Silently ignoring GetPreference error in UpdatePreference
detail: Same issue as line 51 - error from GetPreference is ignored with `pref, _ := ...`
suggestion: Log or handle the error appropriately
```

```
severity: medium
file: backend/notification-service/go.mod
line: 1-15
issue: Missing go.sum file - dependencies not locked
detail: The go.mod file exists but there's no go.sum file to lock dependency versions. This could lead to inconsistent builds.
suggestion: Run `go mod tidy` to generate go.sum file
```

```
severity: medium
file: backend/notification-service/internal/config/config.go
line: 32
issue: DatabaseURL and RedisURL loaded but never used
detail: Config loads DatabaseURL and RedisURL but main.go uses `database.NewConnectionFromEnv()` which reads directly from environment. These config fields are dead code.
suggestion: Either use the config values or remove them to avoid confusion
```

### LOW

```
severity: low
file: backend/notification-service/cmd/main.go
line: 68-69
issue: Unused context variable
detail: `ctx, cancel := context.WithTimeout(...)` is created but `ctx` is only used in `_ = ctx` (explicitly ignored). The timeout context serves no purpose.
suggestion: Remove the unused context or use it for graceful shutdown operations:
  // Remove these lines if not needed:
  ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  defer cancel()
  _ = ctx
```

```
severity: low
file: backend/notification-service/Dockerfile
line: 13
issue: Binary path inconsistent with scoring-service pattern
detail: Notification service builds to `/notification-service` while scoring-service builds to `/app/scoring-service`. This inconsistency could cause confusion.
suggestion: Use consistent path: `RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service ./cmd/main.go`
```

```
severity: low
file: backend/notification-service/internal/worker/worker.go
line: 27
issue: Hardcoded job queue buffer size
detail: Job queue buffer is hardcoded to 100. This should be configurable or at least documented.
suggestion: Make buffer size configurable via WorkerPool constructor or config
```

```
severity: low
file: tests/notification-service/go.mod
line: 1-9
issue: Missing go.sum file for tests
detail: Test module also missing go.sum file
suggestion: Run `go mod tidy` in tests directory
```

---

## Summary

The notification service implementation follows the established patterns from other services (scoring-service, user-service) well. However, there are 2 critical issues that must be fixed before deployment:

1. **Dockerfile COPY path issue** - Will cause build failure
2. **SentAt never persisted** - Notifications will appear as never sent in database

Additionally, the pagination offset bug (HIGH) will cause issues when clients don't explicitly set page number.

### Recommended Priority

1. Fix Dockerfile build context issue
2. Fix SentAt persistence in worker
3. Fix pagination offset calculation
4. Handle GetPreference errors properly
5. Address remaining medium/low issues

### Positive Observations

- Good separation of concerns (channels, worker, service, repository)
- Proper use of interfaces for repository
- Graceful shutdown pattern implemented
- Consistent error response format
- Good test coverage for model validation and channel disabled states

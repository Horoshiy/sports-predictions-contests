# Code Review: Telegram Bot Implementation

**Date**: 2026-01-16
**Reviewer**: Kiro
**Feature**: Telegram Bot for Sports Prediction Platform

## Stats

- Files Modified: 3
- Files Added: 12
- Files Deleted: 0
- New lines: ~450
- Deleted lines: 0

---

## Issues Found

### CRITICAL

```
severity: critical
file: bots/telegram/bot/handlers.go
line: 83
issue: Potential nil pointer dereference on gRPC response
detail: When ListContests returns an error, the code logs and returns. But if err is nil and resp itself is nil (edge case), accessing resp.Contests will panic.
suggestion: Add nil check: if err != nil || resp == nil || resp.Response == nil || !resp.Response.Success
```

```
severity: critical
file: bots/telegram/bot/handlers.go
line: 259-260
issue: Race condition on sessions map
detail: The sessions map is accessed concurrently from multiple goroutines (each update is processed) without synchronization. This can cause data races and crashes.
suggestion: Use sync.RWMutex to protect sessions map access, or use sync.Map
```

### HIGH

```
severity: high
file: bots/telegram/bot/handlers.go
line: 76-77
issue: Ignored error from strconv.ParseUint
detail: When parsing contest_id or leaderboard_id from callback data, errors are silently ignored with `id, _ := strconv.ParseUint(...)`. Malformed callback data will result in id=0 being used.
suggestion: Handle parse error and show user-friendly message instead of proceeding with id=0
```

```
severity: high
file: bots/telegram/clients/clients.go
line: 23-55
issue: No connection cleanup on partial failure
detail: If connecting to the 3rd service fails, the first 2 connections are leaked because Close() is never called on partial initialization failure.
suggestion: Add cleanup in error paths: close already-opened connections before returning error
```

```
severity: high
file: bots/telegram/bot/handlers.go
line: 247
issue: Password visible in Telegram chat history
detail: While the message is deleted after processing, there's a race window where the password is visible. Also, message deletion may fail silently (error not checked).
suggestion: Check delete error and warn user if deletion failed. Consider using inline keyboard flow instead of plain text password.
```

### MEDIUM

```
severity: medium
file: bots/telegram/bot/handlers.go
line: 220-232
issue: Unused streak variables in showUserStats
detail: currentStreak and maxStreak are declared and set to 0, but never populated from analytics response. The comment says "Get streak from first contest if available" but no code implements this.
suggestion: Either implement streak retrieval from analytics or remove the misleading comment and variables
```

```
severity: medium
file: bots/telegram/Dockerfile
line: 11
issue: go mod download before copying source may fail
detail: go.mod has a replace directive pointing to ../../backend/shared which doesn't exist at download time. The build may fail or behave unexpectedly.
suggestion: Copy go.mod and go.sum first, then run go mod download, or restructure to avoid replace directive issues
```

```
severity: medium
file: bots/telegram/bot/bot.go
line: 37-47
issue: Updates channel not drained on stop
detail: When Stop() is called, the updates channel from GetUpdatesChan may still have pending updates. The select loop exits but doesn't drain the channel, potentially causing goroutine leaks in the telegram library.
suggestion: Drain updates channel after StopReceivingUpdates() or use context-based cancellation
```

```
severity: medium
file: tests/telegram-bot/bot_test.go
line: 87-91
issue: Incorrect rank formatting for ranks > 9
detail: formatLeaderboardEntry uses `string(rune('0'+rank))` which only works for single digits. Rank 10 would produce ':' character instead of "10."
suggestion: Use fmt.Sprintf("%d.", rank) for all non-medal ranks (already correct in actual bot code, but test helper is wrong)
```

### LOW

```
severity: low
file: bots/telegram/bot/handlers.go
line: 31
issue: In-memory sessions lost on restart
detail: User sessions are stored in memory map and will be lost when bot restarts. Users will need to re-link their accounts.
suggestion: Document this limitation or consider persisting sessions to Redis/database for production use
```

```
severity: low
file: bots/telegram/go.mod
line: 6-9
issue: Missing indirect dependencies
detail: go.mod doesn't include indirect dependencies that will be added by go mod tidy. This is fine but running go mod tidy is required before build.
suggestion: Run go mod tidy to populate go.sum and indirect deps
```

```
severity: low
file: docker-compose.yml
line: 196-217
issue: No restart policy for telegram-bot service
detail: Other services don't have restart policies either, but for a long-running bot, restart: unless-stopped would be beneficial.
suggestion: Add restart: unless-stopped for production resilience
```

---

## Summary

The Telegram bot implementation is functional and follows project patterns well. The main concerns are:

1. **Race condition on sessions map** - Critical for production stability
2. **Partial connection cleanup** - Memory/resource leak on startup failures  
3. **Password handling** - Security concern with visible passwords

### Recommended Priority Fixes

1. Add mutex protection for sessions map (critical)
2. Fix connection cleanup in clients.go (high)
3. Handle strconv parse errors properly (high)
4. Fix Dockerfile build order (medium)

### Positive Observations

- Clean separation of concerns (config, clients, bot, handlers)
- Proper context timeouts on gRPC calls
- Good error logging with [ERROR] prefix
- Graceful shutdown handling in main.go
- HTML message formatting consistent with notification-service

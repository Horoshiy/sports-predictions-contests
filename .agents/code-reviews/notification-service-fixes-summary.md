# Code Review Fixes Summary - Notification Service

**Date**: 2026-01-16
**Fixes Applied**: 9 issues resolved

## Critical Fixes

### 1. Dockerfile Invalid COPY Path
**File**: `backend/notification-service/Dockerfile`
**Problem**: `COPY ../shared ../shared` is invalid - Docker cannot access files outside build context
**Fix**: 
- Changed build context in docker-compose.yml to `./backend`
- Updated Dockerfile to use relative paths from backend directory
- Binary now builds to `/app/notification-service` for consistency

### 2. SentAt Never Persisted
**File**: `backend/notification-service/internal/worker/worker.go`
**Problem**: Worker set `n.SentAt` but never saved to database
**Fix**: Added `w.repo.Update(ctx, n)` call after successful send, with proper error logging

## High Priority Fixes

### 3. Negative Pagination Offset
**File**: `backend/notification-service/internal/service/notification_service.go`
**Problem**: When page=0 (protobuf default), offset became -20
**Fix**: Added validation to ensure page >= 1 and limit > 0 before calculating offset

### 4. Jobs Not Drained on Shutdown
**File**: `backend/notification-service/internal/worker/worker.go`
**Problem**: Pending jobs lost when Stop() called
**Fix**: Added drain loop in Stop() to process remaining jobs before exiting

### 5. Silently Ignored GetPreference Errors
**File**: `backend/notification-service/internal/service/notification_service.go`
**Problem**: Errors from GetPreference ignored with `pref, _ := ...`
**Fix**: Added proper error logging for both SendNotification and UpdatePreference methods

## Medium Priority Fixes

### 6. Markdown Injection Vulnerability
**File**: `backend/notification-service/internal/channels/telegram.go`
**Problem**: User input not escaped, could break Markdown formatting
**Fix**: Switched to HTML parse mode with `html.EscapeString()` for title and message

### 7. Dead Config Code Removed
**File**: `backend/notification-service/internal/config/config.go`
**Problem**: DatabaseURL and RedisURL loaded but never used
**Fix**: Removed unused fields from Config struct and Load() function

## Low Priority Fixes

### 8. Unused Context Variable
**File**: `backend/notification-service/cmd/main.go`
**Problem**: Context created but never used in shutdown handler
**Fix**: Removed unused context creation and related imports (context, time)

### 9. Hardcoded Job Queue Buffer
**File**: `backend/notification-service/internal/worker/worker.go`
**Problem**: Buffer size hardcoded to 100
**Fix**: Made buffer size proportional to worker count (20 jobs per worker, minimum 100)

## Files Modified

- `backend/notification-service/Dockerfile`
- `backend/notification-service/cmd/main.go`
- `backend/notification-service/internal/config/config.go`
- `backend/notification-service/internal/channels/telegram.go`
- `backend/notification-service/internal/service/notification_service.go`
- `backend/notification-service/internal/worker/worker.go`
- `docker-compose.yml`

## Files Created

- `backend/notification-service/go.sum` (placeholder)
- `tests/notification-service/go.sum` (placeholder)

## Validation

Run these commands to validate fixes:
```bash
cd backend/notification-service && go mod tidy && go build ./...
cd backend && go work sync
docker-compose build notification-service
```

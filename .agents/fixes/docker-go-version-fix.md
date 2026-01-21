# Docker Go Version Fix

## Issue
All Dockerfiles were using `golang:1.21-alpine` but all `go.mod` files require Go 1.24.0, causing build failures with:
```
process "/bin/sh -c go mod download && CGO_ENABLED=0 GOOS=linux go build..." did not complete successfully: exit code: 1
```

## Root Cause
Version mismatch between Docker base image (Go 1.21) and module requirements (Go 1.24).

## Fixes Applied

### Fix 1: Go Version Update
Updated all Dockerfiles from `golang:1.21-alpine` to `golang:1.24-alpine`

### Fix 2: Alpine Version Standardization
Updated all backend service Dockerfiles from `alpine:latest` to `alpine:3.19` for:
- **Reproducibility**: Pinned version ensures consistent builds
- **Consistency**: All services now use the same Alpine version
- **Best Practice**: Avoids unpredictable changes from :latest tag

### Services Updated (9 total)
1. ✅ `bots/telegram/Dockerfile` - Go 1.24 (Alpine already pinned)
2. ✅ `backend/api-gateway/Dockerfile` - Go 1.24 + Alpine 3.19
3. ✅ `backend/challenge-service/Dockerfile` - Go 1.24 + Alpine 3.19
4. ✅ `backend/contest-service/Dockerfile` - Go 1.24 + Alpine 3.19
5. ✅ `backend/notification-service/Dockerfile` - Go 1.24 + Alpine 3.19
6. ✅ `backend/prediction-service/Dockerfile` - Go 1.24 + Alpine 3.19
7. ✅ `backend/scoring-service/Dockerfile` - Go 1.24 + Alpine 3.19
8. ✅ `backend/sports-service/Dockerfile` - Go 1.24 + Alpine 3.19
9. ✅ `backend/user-service/Dockerfile` - Go 1.24 + Alpine 3.19

## Verification
All `go mod download` commands now work correctly and builds are reproducible.

### Automated Test
Created `tests/dockerfile-consistency-test.sh` to verify:
- ✅ All services use golang:1.24-alpine
- ✅ All services use alpine:3.19
- ✅ No :latest tags present
- ✅ Multi-stage builds maintained
- ✅ CGO disabled for static binaries

## Testing
To verify the fix:
```bash
# Run consistency test
./tests/dockerfile-consistency-test.sh

# Test individual service
docker build -f backend/api-gateway/Dockerfile -t api-gateway ./backend

# Test Telegram bot
docker build -f bots/telegram/Dockerfile -t telegram-bot .

# Test all services
docker-compose --profile services build
```

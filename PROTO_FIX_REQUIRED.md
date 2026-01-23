# Protocol Buffer Generation Required

## Problem
The Docker build is failing because Protocol Buffer files are not fully generated. The error you're seeing:

```
target api-gateway: failed to solve: process "/bin/sh -c CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go" did not complete successfully: exit code: 1
```

This happens because Go services are trying to import proto packages that don't exist yet.

## Root Cause
Most proto files in `backend/shared/proto/` are missing their generated `.pb.go` and `_grpc.pb.go` files:

### Missing Files:
- `backend/shared/proto/scoring/` - **COMPLETELY MISSING** (needed by api-gateway, scoring-service)
- `backend/shared/proto/profile/` - **COMPLETELY MISSING** (needed by user-service)
- `backend/shared/proto/contest/` - Missing `contest.pb.go` and `contest_grpc.pb.go`
- `backend/shared/proto/prediction/` - Missing `prediction.pb.go` and `prediction_grpc.pb.go`
- `backend/shared/proto/sports/` - Missing `sports.pb.go` and `sports_grpc.pb.go`
- `backend/shared/proto/user/` - Missing `user.pb.go` and `user_grpc.pb.go`

### Complete Files (OK):
- `backend/shared/proto/challenge/` ✓
- `backend/shared/proto/common/` ✓
- `backend/shared/proto/notification/` ✓
- `backend/shared/proto/team/` ✓

## Solution

### Option 1: Install protoc and generate (RECOMMENDED)

1. **Install Protocol Buffers compiler:**
   ```bash
   # macOS
   brew install protobuf
   
   # Linux
   sudo apt-get install -y protobuf-compiler
   ```

2. **Install Go plugins:**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. **Generate proto files:**
   ```bash
   ./scripts/generate-protos.sh
   ```
   
   Or manually:
   ```bash
   cd backend
   protoc --proto_path=proto \
       --go_out=shared \
       --go-grpc_out=shared \
       proto/*.proto
   ```

4. **Verify generation:**
   ```bash
   find backend/shared/proto -name "*.pb.go" | sort
   ```
   
   You should see at least 18 files (2 per service: .pb.go and _grpc.pb.go for 9 services).

### Option 2: Use Docker to generate

If you can't install protoc locally, use Docker:

```bash
# Start Docker daemon first
docker run --rm \
    -v $(pwd)/backend/proto:/protos \
    -v $(pwd)/backend/shared:/output \
    namely/protoc-all:latest \
    -f /protos/*.proto \
    -l go \
    -o /output
```

### Option 3: Use Makefile target

```bash
make proto
```

This runs the same protoc command but through the Makefile.

## Verification

After generating proto files, verify the build works:

```bash
# Test local build
cd backend/api-gateway
go build -o /tmp/test-api-gateway ./cmd/main.go

# Test Docker build
docker-compose build api-gateway
```

## Why This Happened

The proto source files (`.proto`) exist in `backend/proto/`, but the generated Go code wasn't committed to the repository. This is common in projects where generated files are excluded from version control.

However, for this project to work out-of-the-box, either:
1. Generated files should be committed, OR
2. Proto generation should be part of the Docker build process

## Quick Fix for Docker Build

If you want Docker to generate protos automatically, modify each service's Dockerfile to include proto generation:

```dockerfile
FROM golang:1.24-alpine AS builder

# Install protoc
RUN apk add --no-cache protobuf protobuf-dev

# Install Go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

# Copy proto files and generate
COPY proto ./proto
COPY shared/go.mod shared/go.sum ./shared/
RUN cd proto && protoc --proto_path=. \
    --go_out=../shared \
    --go-grpc_out=../shared \
    *.proto

# ... rest of Dockerfile
```

## Next Steps

1. Generate proto files using one of the options above
2. Verify all services build successfully
3. Test Docker Compose startup: `make docker-services`
4. Consider committing generated files or updating Dockerfiles

FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY challenge-service/go.mod challenge-service/go.sum ./challenge-service/
COPY shared ./shared

# Set working directory to service
WORKDIR /app/challenge-service

# Download dependencies
RUN go mod download

# Copy source code
COPY challenge-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o challenge-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/challenge-service/challenge-service .

# Expose port
EXPOSE 8090

# Command to run
CMD ["./challenge-service"]

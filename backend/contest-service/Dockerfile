FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY contest-service/go.mod contest-service/go.sum ./contest-service/
COPY shared ./shared

# Set working directory to service
WORKDIR /app/contest-service

# Download dependencies
RUN go mod download

# Copy source code
COPY contest-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o contest-service ./cmd/main.go

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/contest-service/contest-service .

# Expose port
EXPOSE 8085

# Command to run
CMD ["./contest-service"]

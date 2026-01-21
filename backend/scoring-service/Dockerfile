FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY scoring-service/go.mod scoring-service/go.sum ./scoring-service/
COPY shared ./shared

# Set working directory to service
WORKDIR /app/scoring-service

# Download dependencies
RUN go mod download

# Copy source code
COPY scoring-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o scoring-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/scoring-service/scoring-service .

# Expose port
EXPOSE 8087

# Command to run
CMD ["./scoring-service"]

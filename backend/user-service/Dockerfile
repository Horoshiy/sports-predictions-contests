FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY user-service/go.mod user-service/go.sum ./user-service/
COPY shared ./shared

# Set working directory to service
WORKDIR /app/user-service

# Download dependencies
RUN go mod download

# Copy source code
COPY user-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/user-service/user-service .

# Expose port
EXPOSE 8084

# Run the binary
CMD ["./user-service"]

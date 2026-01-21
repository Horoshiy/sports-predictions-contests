FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY shared ./shared

# Set working directory to service
WORKDIR /app/api-gateway

# Download dependencies
RUN go mod download

# Copy source code
COPY api-gateway/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway/api-gateway .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY scoring-service/go.mod scoring-service/go.sum ./scoring-service/
COPY shared/go.mod shared/go.sum ./shared/

# Set working directory to service
WORKDIR /app/scoring-service

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Copy only necessary shared files
WORKDIR /app
COPY shared ./shared

# Copy source code
WORKDIR /app/scoring-service
COPY scoring-service/ .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o scoring-service ./cmd/main.go

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/scoring-service/scoring-service .

# Expose port
EXPOSE 8087

# Command to run
CMD ["./scoring-service"]

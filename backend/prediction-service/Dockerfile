FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY prediction-service/go.mod prediction-service/go.sum ./prediction-service/
COPY shared ./shared

# Set working directory to service
WORKDIR /app/prediction-service

# Download dependencies
RUN go mod download

# Copy source code
COPY prediction-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o prediction-service ./cmd/main.go

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/prediction-service/prediction-service .

# Expose port
EXPOSE 8086

# Run the binary
CMD ["./prediction-service"]

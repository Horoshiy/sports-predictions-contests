FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy go mod files
COPY notification-service/go.mod notification-service/go.sum* ./notification-service/
COPY shared ./shared

WORKDIR /app/notification-service

# Download dependencies
RUN go mod download

# Copy source code (excluding go.mod files already copied)
COPY notification-service/cmd ./cmd
COPY notification-service/internal ./internal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/notification-service/notification-service .

EXPOSE 8089

CMD ["./notification-service"]

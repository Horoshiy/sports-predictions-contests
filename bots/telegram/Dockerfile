FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy shared module first (required for replace directive)
COPY backend/shared ./backend/shared

# Copy bot source
COPY bots/telegram ./bots/telegram

WORKDIR /app/bots/telegram

# Download deps and build (shared already available for replace)
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /telegram-bot .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

COPY --from=builder /telegram-bot /telegram-bot

ENTRYPOINT ["/telegram-bot"]

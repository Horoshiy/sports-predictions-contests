FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY shared ./shared

WORKDIR /app/notification-service

# Download dependencies
RUN go mod download

# Copy source code
COPY notification-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o notification-service ./cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/notification-service/notification-service .

EXPOSE 8089

CMD ["./notification-service"]

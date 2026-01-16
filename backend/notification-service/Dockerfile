FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy the entire backend directory (build context should be ./backend)
COPY notification-service/go.mod notification-service/go.sum* ./notification-service/
COPY shared ./shared

WORKDIR /app/notification-service
RUN go mod download

COPY notification-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/notification-service .

EXPOSE 8089

CMD ["./notification-service"]

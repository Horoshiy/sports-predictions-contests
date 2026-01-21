FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY sports-service/go.mod sports-service/go.sum ./sports-service/
COPY shared ./shared

WORKDIR /app/sports-service

RUN go mod download

COPY sports-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o sports-service ./cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/sports-service/sports-service .

EXPOSE 8088

CMD ["./sports-service"]

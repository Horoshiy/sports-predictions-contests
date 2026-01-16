FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY ../shared ../shared

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o sports-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/sports-service .

EXPOSE 8088

CMD ["./sports-service"]

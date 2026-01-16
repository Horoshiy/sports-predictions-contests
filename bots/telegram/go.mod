module github.com/sports-prediction-contests/telegram-bot

go 1.21

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/sports-prediction-contests/shared v0.0.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
)

replace github.com/sports-prediction-contests/shared => ../../backend/shared

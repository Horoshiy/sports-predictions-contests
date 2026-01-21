module github.com/sports-prediction-contests/telegram-bot

go 1.24.0

toolchain go1.24.12

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/sports-prediction-contests/shared v0.0.0
	google.golang.org/grpc v1.60.1
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace github.com/sports-prediction-contests/shared => ../../backend/shared

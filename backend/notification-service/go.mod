module github.com/sports-prediction-contests/notification-service

go 1.21

require (
	github.com/sports-prediction-contests/shared v0.0.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
)

replace github.com/sports-prediction-contests/shared => ../shared

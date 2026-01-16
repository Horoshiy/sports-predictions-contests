module github.com/sports-prediction-contests/sports-service

go 1.21

require (
	github.com/sports-prediction-contests/shared v0.0.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

replace github.com/sports-prediction-contests/shared => ../shared

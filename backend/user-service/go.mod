module github.com/sports-prediction-contests/user-service

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/sports-prediction-contests/shared v0.0.0
	golang.org/x/crypto v0.18.0
	google.golang.org/grpc v1.60.1
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

replace github.com/sports-prediction-contests/shared => ../shared

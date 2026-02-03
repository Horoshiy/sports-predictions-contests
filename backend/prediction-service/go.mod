module github.com/sports-prediction-contests/prediction-service

go 1.24.0

toolchain go1.24.12

require (
	github.com/sports-prediction-contests/shared v0.0.0
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
	gorm.io/gorm v1.30.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	gorm.io/driver/postgres v1.5.4 // indirect
)

replace github.com/sports-prediction-contests/shared => ../shared

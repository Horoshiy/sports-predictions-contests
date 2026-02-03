module github.com/sports-prediction-contests/api-gateway

go 1.24.0

toolchain go1.24.12

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.5
	github.com/sports-prediction-contests/shared v0.0.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/sports-prediction-contests/shared => ../shared

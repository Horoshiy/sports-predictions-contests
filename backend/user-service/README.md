# User Service

JWT-based authentication service for the Sports Prediction Contests platform.

## Features

- User registration with email validation
- Secure login with bcrypt password hashing
- JWT token generation and validation
- Profile management (get/update)
- gRPC API with authentication interceptors

## API Endpoints

### Public Endpoints (No Authentication Required)
- `Register(RegisterRequest)` - Create new user account
- `Login(LoginRequest)` - Authenticate user and get JWT token

### Protected Endpoints (JWT Token Required)
- `GetProfile(GetProfileRequest)` - Get current user profile
- `UpdateProfile(UpdateProfileRequest)` - Update user profile

## Environment Variables

```bash
USER_SERVICE_PORT=8084
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION=24h
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
LOG_LEVEL=info
```

## Running the Service

### Development
```bash
cd backend/user-service
go run cmd/main.go
```

### Docker
```bash
docker-compose up user-service
```

## Testing

```bash
# Unit tests
go test ./...

# With coverage
go test ./... -cover
```

## gRPC Testing

```bash
# Register user
grpcurl -plaintext -d '{"email":"test@example.com","password":"password123","name":"Test User"}' localhost:8084 user.UserService/Register

# Login
grpcurl -plaintext -d '{"email":"test@example.com","password":"password123"}' localhost:8084 user.UserService/Login

# Get profile (use token from login)
grpcurl -plaintext -H "authorization: Bearer <JWT_TOKEN>" localhost:8084 user.UserService/GetProfile
```

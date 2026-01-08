package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	pb "github.com/sports-prediction-contests/shared/proto/user"
	"github.com/sports-prediction-contests/user-service/internal/config"
	"github.com/sports-prediction-contests/user-service/internal/models"
	"github.com/sports-prediction-contests/user-service/internal/repository"
	"github.com/sports-prediction-contests/user-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewConnectionFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, []byte(cfg.JWTSecret), cfg.JWTExpiration)
	userService := service.NewUserService(authService, userRepo)

	// Create gRPC server with JWT interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
	)

	// Register services
	pb.RegisterUserServiceServer(server, userService)

	// Start listening
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("User service starting on port %s", cfg.Port)

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down user service...")
		server.GracefulStop()
	}()

	// Start server
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

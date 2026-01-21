package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/challenge-service/internal/config"
	"github.com/sports-prediction-contests/challenge-service/internal/models"
	"github.com/sports-prediction-contests/challenge-service/internal/repository"
	"github.com/sports-prediction-contests/challenge-service/internal/service"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	pb "github.com/sports-prediction-contests/shared/proto/challenge"
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
	if err := db.AutoMigrate(&models.Challenge{}, &models.ChallengeParticipant{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	challengeRepo := repository.NewChallengeRepository(db)
	participantRepo := repository.NewChallengeParticipantRepository(db)

	// Initialize services
	challengeService := service.NewChallengeService(challengeRepo, participantRepo)

	// Create gRPC server with JWT interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
	)

	// Register services
	pb.RegisterChallengeServiceServer(server, challengeService)

	// Start listening
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("[INFO] Challenge Service starting on port %s", cfg.Port)

	// Start server in a goroutine
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-c
	log.Println("[INFO] Shutting down Challenge Service...")

	// Gracefully stop the server
	server.GracefulStop()
	
	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
	
	log.Println("[INFO] Challenge Service stopped")
}

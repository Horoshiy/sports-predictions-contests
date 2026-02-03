package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/prediction-service/internal/clients"
	"github.com/sports-prediction-contests/prediction-service/internal/config"
	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"github.com/sports-prediction-contests/prediction-service/internal/repository"
	"github.com/sports-prediction-contests/prediction-service/internal/service"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	pb "github.com/sports-prediction-contests/shared/proto/prediction"
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
	if err := db.AutoMigrate(
		&models.Prediction{},
		&models.Event{},
		&models.PropType{},
		&models.RiskyEventType{},
		&models.MatchRiskyEvent{},
		&models.RiskyPrediction{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize contest client
	contestClient, err := clients.NewContestClient(cfg.ContestServiceEndpoint)
	if err != nil {
		log.Fatalf("Failed to create contest client: %v", err)
	}
	defer contestClient.Close()

	// Initialize repositories
	predictionRepo := repository.NewPredictionRepository(db)
	eventRepo := repository.NewEventRepository(db)
	propTypeRepo := repository.NewPropTypeRepository(db)
	riskyEventRepo := repository.NewRiskyEventRepository(db)

	// Initialize services
	predictionService := service.NewPredictionService(predictionRepo, eventRepo, propTypeRepo, riskyEventRepo, contestClient)

	// Create gRPC server with JWT interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
	)

	// Register services
	pb.RegisterPredictionServiceServer(server, predictionService)

	// Start listening
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("[INFO] Prediction Service starting on port %s", cfg.Port)

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
	log.Println("[INFO] Shutting down Prediction Service...")

	// Gracefully stop the server
	server.GracefulStop()

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("[INFO] Prediction Service stopped")
}

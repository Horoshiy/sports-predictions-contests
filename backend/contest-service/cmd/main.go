package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/contest-service/internal/config"
	"github.com/sports-prediction-contests/contest-service/internal/models"
	"github.com/sports-prediction-contests/contest-service/internal/repository"
	"github.com/sports-prediction-contests/contest-service/internal/service"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	contestpb "github.com/sports-prediction-contests/shared/proto/contest"
	teampb "github.com/sports-prediction-contests/shared/proto/team"
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
		&models.Contest{},
		&models.Participant{},
		&models.Team{},
		&models.TeamMember{},
		&models.TeamContestEntry{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Printf("[INFO] Database migration completed successfully")

	// Initialize repositories
	contestRepo := repository.NewContestRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	memberRepo := repository.NewTeamMemberRepository(db)
	contestEntryRepo := repository.NewTeamContestEntryRepository(db)

	// Initialize services
	contestService := service.NewContestService(contestRepo, participantRepo)
	teamService := service.NewTeamService(teamRepo, memberRepo, contestEntryRepo)
	teamServiceGRPC := service.NewTeamServiceGRPC(teamService)

	// Create gRPC server with JWT interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
	)

	// Register services
	contestpb.RegisterContestServiceServer(server, contestService)
	teampb.RegisterTeamServiceServer(server, teamServiceGRPC)

	// Start listening
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("[INFO] Contest & Team Service starting on port %s", cfg.Port)

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
	log.Println("[INFO] Shutting down Contest & Team Service...")

	// Gracefully stop the server
	server.GracefulStop()
	
	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
	
	log.Println("[INFO] Contest & Team Service stopped")
}

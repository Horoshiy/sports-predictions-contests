package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/sports-service/internal/config"
	"github.com/sports-prediction-contests/sports-service/internal/external"
	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
	"github.com/sports-prediction-contests/sports-service/internal/service"
	"github.com/sports-prediction-contests/sports-service/internal/sync"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	pb "github.com/sports-prediction-contests/shared/proto/sports"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	db, err := database.NewConnectionFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.Sport{}, &models.League{}, &models.Team{}, &models.Match{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	sportRepo := repository.NewSportRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	sportsService := service.NewSportsService(sportRepo, leagueRepo, teamRepo, matchRepo)

	// Initialize sync worker if enabled
	var syncWorker *sync.SyncWorker
	if cfg.SyncEnabled {
		apiClient := external.NewClient(cfg.TheSportsDBURL)
		syncService := sync.NewSyncService(apiClient, sportRepo, leagueRepo, teamRepo, matchRepo)
		syncWorker = sync.NewSyncWorker(syncService, cfg.SyncIntervalMins)
		syncWorker.Start()
		log.Printf("[INFO] Sync worker started with %d minute interval", cfg.SyncIntervalMins)
	}
	sportsService.SetSyncWorker(syncWorker, cfg.SyncEnabled, cfg.SyncIntervalMins)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
	)

	pb.RegisterSportsServiceServer(server, sportsService)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("[INFO] Sports Service starting on port %s", cfg.Port)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("[INFO] Shutting down Sports Service...")

	// Stop sync worker first
	if syncWorker != nil {
		syncWorker.Stop()
	}

	server.GracefulStop()

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("[INFO] Sports Service stopped")
}

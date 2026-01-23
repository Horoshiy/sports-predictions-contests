package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sports-prediction-contests/scoring-service/internal/cache"
	"github.com/sports-prediction-contests/scoring-service/internal/config"
	"github.com/sports-prediction-contests/scoring-service/internal/models"
	"github.com/sports-prediction-contests/scoring-service/internal/repository"
	"github.com/sports-prediction-contests/scoring-service/internal/service"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	"github.com/sports-prediction-contests/shared/proto/common"
	pb "github.com/sports-prediction-contests/shared/proto/scoring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	if err := db.AutoMigrate(&models.Score{}, &models.Leaderboard{}, &models.UserStreak{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to Redis
	redisCache, err := cache.NewRedisCacheFromURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()

	// Initialize repositories
	scoreRepo := repository.NewScoreRepository(db)
	leaderboardRepo := repository.NewLeaderboardRepository(db, redisCache)
	streakRepo := repository.NewStreakRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	// Initialize services
	scoringService := service.NewScoringService(scoreRepo, leaderboardRepo, streakRepo, analyticsRepo)
	leaderboardService := service.NewLeaderboardService(leaderboardRepo, scoreRepo, streakRepo)

	// Create combined service that implements all methods
	combinedService := &CombinedScoringService{
		ScoringService:     scoringService,
		LeaderboardService: leaderboardService,
	}

	// Create gRPC server with JWT interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))),
	)

	// Register services
	pb.RegisterScoringServiceServer(server, combinedService)

	// Register health service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Start server
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("Scoring Service starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.DatabaseURL)
	log.Printf("Redis URL: %s", cfg.RedisURL)

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down Scoring Service...")
		
		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Set health status to not serving
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

		// Graceful stop with context (for future use with cleanup operations)
		done := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			log.Println("Scoring Service stopped gracefully")
		case <-ctx.Done():
			log.Println("Shutdown timeout exceeded, forcing stop")
			server.Stop()
		}
		
		os.Exit(0)
	}()

	// Start serving
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// CombinedScoringService combines both scoring and leaderboard services
type CombinedScoringService struct {
	pb.UnimplementedScoringServiceServer
	ScoringService     *service.ScoringService
	LeaderboardService *service.LeaderboardService
}

// Implement all scoring service methods by delegating to appropriate services

func (s *CombinedScoringService) CreateScore(ctx context.Context, req *pb.CreateScoreRequest) (*pb.CreateScoreResponse, error) {
	return s.ScoringService.CreateScore(ctx, req)
}

func (s *CombinedScoringService) GetScore(ctx context.Context, req *pb.GetScoreRequest) (*pb.GetScoreResponse, error) {
	return s.ScoringService.GetScore(ctx, req)
}

func (s *CombinedScoringService) GetUserScores(ctx context.Context, req *pb.GetUserScoresRequest) (*pb.GetUserScoresResponse, error) {
	return s.ScoringService.GetUserScores(ctx, req)
}

func (s *CombinedScoringService) CalculateScore(ctx context.Context, req *pb.CalculateScoreRequest) (*pb.CalculateScoreResponse, error) {
	return s.ScoringService.CalculateScore(ctx, req)
}

func (s *CombinedScoringService) GetLeaderboard(ctx context.Context, req *pb.GetLeaderboardRequest) (*pb.GetLeaderboardResponse, error) {
	return s.LeaderboardService.GetLeaderboard(ctx, req)
}

func (s *CombinedScoringService) GetUserRank(ctx context.Context, req *pb.GetUserRankRequest) (*pb.GetUserRankResponse, error) {
	return s.LeaderboardService.GetUserRank(ctx, req)
}

func (s *CombinedScoringService) GetUserStreak(ctx context.Context, req *pb.GetUserStreakRequest) (*pb.GetUserStreakResponse, error) {
	return s.LeaderboardService.GetUserStreak(ctx, req)
}

func (s *CombinedScoringService) UpdateLeaderboard(ctx context.Context, req *pb.UpdateLeaderboardRequest) (*pb.UpdateLeaderboardResponse, error) {
	return s.LeaderboardService.UpdateLeaderboard(ctx, req)
}

func (s *CombinedScoringService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return &common.Response{
		Success:   true,
		Message:   "Scoring service is healthy",
		Code:      0,
		Timestamp: timestamppb.Now(),
	}, nil
}

func (s *CombinedScoringService) GetUserAnalytics(ctx context.Context, req *pb.GetUserAnalyticsRequest) (*pb.GetUserAnalyticsResponse, error) {
	return s.ScoringService.GetUserAnalytics(ctx, req)
}

func (s *CombinedScoringService) ExportAnalytics(ctx context.Context, req *pb.ExportAnalyticsRequest) (*pb.ExportAnalyticsResponse, error) {
	return s.ScoringService.ExportAnalytics(ctx, req)
}

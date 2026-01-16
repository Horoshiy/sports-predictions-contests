package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sports-prediction-contests/notification-service/internal/channels"
	"github.com/sports-prediction-contests/notification-service/internal/config"
	"github.com/sports-prediction-contests/notification-service/internal/models"
	"github.com/sports-prediction-contests/notification-service/internal/repository"
	"github.com/sports-prediction-contests/notification-service/internal/service"
	"github.com/sports-prediction-contests/notification-service/internal/worker"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/database"
	pb "github.com/sports-prediction-contests/shared/proto/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	if err := db.AutoMigrate(&models.Notification{}, &models.NotificationPreference{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	telegram, err := channels.NewTelegramChannel(cfg.TelegramBotToken, cfg.TelegramEnabled)
	if err != nil {
		log.Printf("Warning: Telegram channel disabled: %v", err)
	}

	email := channels.NewEmailChannel(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom, cfg.EmailEnabled)

	repo := repository.NewNotificationRepository(db)
	workerPool := worker.NewWorkerPool(cfg.WorkerPoolSize, telegram, email, repo)
	workerPool.Start()

	notificationService := service.NewNotificationService(repo, workerPool)

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.JWTUnaryInterceptor([]byte(cfg.JWTSecret))))
	pb.RegisterNotificationServiceServer(server, notificationService)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("Notification Service starting on port %s", cfg.Port)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down Notification Service...")

		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		workerPool.Stop()
		server.GracefulStop()
		log.Println("Notification Service stopped")
		os.Exit(0)
	}()

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

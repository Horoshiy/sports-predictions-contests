package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sports-prediction-contests/api-gateway/internal/config"
	"github.com/sports-prediction-contests/api-gateway/internal/gateway"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create gateway server
	server, err := gateway.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create gateway server: %v", err)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: server.Handler(),
	}

	log.Printf("API Gateway starting on port %s", cfg.Port)
	log.Printf("User Service: %s", cfg.UserService)
	log.Printf("Contest Service: %s", cfg.ContestService)

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down API Gateway...")
		
		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start server
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to serve: %v", err)
	}
}

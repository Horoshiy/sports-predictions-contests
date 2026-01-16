//go:build e2e

package e2e

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.Println("E2E Tests: Waiting for services to be ready...")

	if err := waitForServices(); err != nil {
		log.Fatalf("Services not ready: %v", err)
	}

	log.Println("E2E Tests: Services are ready, running tests...")
	code := m.Run()

	log.Println("E2E Tests: Completed")
	os.Exit(code)
}

func waitForServices() error {
	healthURL := BaseURL() + "/health"
	timeout := 60 * time.Second

	log.Printf("Checking API Gateway health at %s", healthURL)
	return waitForService(healthURL, timeout)
}

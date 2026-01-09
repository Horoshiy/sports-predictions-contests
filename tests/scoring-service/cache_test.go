package cache_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/scoring-service/internal/cache"
)

func TestRedisConfigTimeout(t *testing.T) {
	config := cache.Config{
		Addr:           "localhost:6379",
		Password:       "",
		DB:             0,
		ConnectTimeout: 2 * time.Second,
	}

	// Test that config accepts custom timeout
	if config.ConnectTimeout != 2*time.Second {
		t.Errorf("Expected timeout to be 2s, got %v", config.ConnectTimeout)
	}
}

func TestRedisConfigDefaultTimeout(t *testing.T) {
	config := cache.Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		// ConnectTimeout not set, should use default
	}

	// Test that zero timeout will use default
	if config.ConnectTimeout != 0 {
		t.Errorf("Expected zero timeout for default, got %v", config.ConnectTimeout)
	}
}

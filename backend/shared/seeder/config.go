package seeder

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// DataSize represents the size of the dataset to generate
type DataSize string

const (
	DataSizeSmall  DataSize = "small"
	DataSizeMedium DataSize = "medium"
	DataSizeLarge  DataSize = "large"
)

// Config holds the configuration for data seeding
type Config struct {
	Size        DataSize `json:"size"`
	Seed        int64    `json:"seed"`
	BatchSize   int      `json:"batch_size"`
	DatabaseURL string   `json:"database_url"`
}

// DataCounts holds the number of entities to generate for each size
type DataCounts struct {
	Users       int
	Sports      int
	Leagues     int
	Teams       int
	Matches     int
	Contests    int
	Challenges  int
	Predictions int
	UserTeams   int
}

// LoadConfig loads seeding configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Size:        DataSize(getEnv("SEED_SIZE", "small")),
		Seed:        getEnvInt64("SEED_VALUE", 42),
		BatchSize:   getEnvInt("BATCH_SIZE", 100),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	return config
}

// GetDataCounts returns the number of entities to generate based on the data size
func (c *Config) GetDataCounts() DataCounts {
	switch c.Size {
	case DataSizeSmall:
		return DataCounts{
			Users:       20,
			Sports:      3,
			Leagues:     6,
			Teams:       24,
			Matches:     50,
			Contests:    8,
			Challenges:  15,
			Predictions: 200,
			UserTeams:   4,
		}
	case DataSizeMedium:
		return DataCounts{
			Users:       100,
			Sports:      5,
			Leagues:     15,
			Teams:       60,
			Matches:     200,
			Contests:    25,
			Challenges:  60,
			Predictions: 1000,
			UserTeams:   15,
		}
	case DataSizeLarge:
		return DataCounts{
			Users:       500,
			Sports:      8,
			Leagues:     30,
			Teams:       120,
			Matches:     800,
			Contests:    50,
			Challenges:  200,
			Predictions: 5000,
			UserTeams:   30,
		}
	default:
		return DataCounts{
			Users:       20,
			Sports:      3,
			Leagues:     6,
			Teams:       24,
			Matches:     50,
			Contests:    8,
			Challenges:  15,
			Predictions: 200,
			UserTeams:   4,
		}
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Size != DataSizeSmall && c.Size != DataSizeMedium && c.Size != DataSizeLarge {
		return fmt.Errorf("invalid data size: %s (must be small, medium, or large)", c.Size)
	}

	if c.BatchSize <= 0 {
		return fmt.Errorf("batch size must be positive, got: %d", c.BatchSize)
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required for seeding operations")
	}

	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		} else {
			log.Printf("seeder: failed to parse environment variable %s as integer: %v, using default %d", key, err, defaultValue)
		}
	}
	return defaultValue
}

// getEnvInt64 gets an int64 environment variable with a default value
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		} else {
			log.Printf("seeder: failed to parse environment variable %s as int64: %v, using default %d", key, err, defaultValue)
		}
	}
	return defaultValue
}

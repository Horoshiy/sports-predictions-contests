package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/sports-prediction-contests/shared/seeder"
)

func main() {
	// Command line flags
	var (
		size     = flag.String("size", "small", "Data size: small, medium, large")
		seed     = flag.Int64("seed", 42, "Random seed for reproducible data")
		dbURL    = flag.String("db", "", "Database URL (overrides environment)")
		testMode = flag.Bool("test", false, "Test mode - validate seeding without committing")
		help     = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Load configuration
	config := seeder.LoadConfig()
	
	// Override with command line flags
	if *size != "small" {
		config.Size = seeder.DataSize(*size)
	}
	if *seed != 42 {
		config.Seed = *seed
	}
	if *dbURL != "" {
		config.DatabaseURL = *dbURL
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Printf("Starting data seeding with configuration:")
	log.Printf("  Size: %s", config.Size)
	log.Printf("  Seed: %d", config.Seed)
	log.Printf("  Batch Size: %d", config.BatchSize)
	log.Printf("  Test Mode: %t", *testMode)

	// Create coordinator
	coordinator, err := seeder.NewCoordinator(config)
	if err != nil {
		log.Fatalf("Failed to create coordinator: %v", err)
	}
	defer coordinator.Close()

	// Start seeding
	startTime := time.Now()
	
	if *testMode {
		log.Println("Running in test mode - seeding will be rolled back")
		if err := coordinator.TestSeed(); err != nil {
			log.Fatalf("Test seeding failed: %v", err)
		}
		log.Println("Test mode completed successfully - all seeding logic validated")
		return
	}

	if err := coordinator.SeedAll(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Seeding completed successfully in %v", duration)

	// Print summary
	counts := config.GetDataCounts()
	fmt.Println("\n=== SEEDING SUMMARY ===")
	fmt.Printf("Data Size: %s\n", config.Size)
	fmt.Printf("Random Seed: %d\n", config.Seed)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Println("\nGenerated Entities:")
	fmt.Printf("  Users: %d (with profiles and preferences)\n", counts.Users)
	fmt.Printf("  Sports: %d\n", counts.Sports)
	fmt.Printf("  Leagues: %d\n", counts.Leagues)
	fmt.Printf("  Teams: %d\n", counts.Teams)
	fmt.Printf("  Matches: %d\n", counts.Matches)
	fmt.Printf("  Contests: %d\n", counts.Contests)
	fmt.Printf("  Predictions: %d\n", counts.Predictions)
	fmt.Printf("  User Teams: %d\n", counts.UserTeams)
	fmt.Println("\nAdditional Data:")
	fmt.Println("  - User profiles and preferences")
	fmt.Println("  - Scoring data (scores, leaderboards, streaks)")
	fmt.Println("  - Notifications and preferences")
	fmt.Println("  - Team memberships")
	fmt.Println("\n=== READY FOR DEVELOPMENT ===")
	fmt.Println("Your Sports Prediction Contests platform is now populated with realistic data!")
	fmt.Println("You can start the frontend and explore all features immediately.")
}

func showHelp() {
	fmt.Println("Sports Prediction Contests - Data Seeder")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("This tool populates the database with realistic fake data for development and testing.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  go run scripts/seed-data.go [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -size string")
	fmt.Println("        Data size preset: small, medium, large (default: small)")
	fmt.Println("        small:  20 users, 8 contests, 200 predictions")
	fmt.Println("        medium: 100 users, 25 contests, 1000 predictions")
	fmt.Println("        large:  500 users, 50 contests, 5000 predictions")
	fmt.Println()
	fmt.Println("  -seed int")
	fmt.Println("        Random seed for reproducible data generation (default: 42)")
	fmt.Println()
	fmt.Println("  -db string")
	fmt.Println("        Database URL (overrides DATABASE_URL environment variable)")
	fmt.Println()
	fmt.Println("  -test")
	fmt.Println("        Test mode - validate configuration without seeding data")
	fmt.Println()
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("ENVIRONMENT VARIABLES:")
	fmt.Println("  DATABASE_URL     Database connection string")
	fmt.Println("  SEED_SIZE        Data size preset (small/medium/large)")
	fmt.Println("  SEED_VALUE       Random seed value")
	fmt.Println("  BATCH_SIZE       Batch size for database operations (default: 100)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Seed with small dataset")
	fmt.Println("  go run scripts/seed-data.go")
	fmt.Println()
	fmt.Println("  # Seed with medium dataset and custom seed")
	fmt.Println("  go run scripts/seed-data.go -size medium -seed 123")
	fmt.Println()
	fmt.Println("  # Test configuration without seeding")
	fmt.Println("  go run scripts/seed-data.go -test")
	fmt.Println()
	fmt.Println("  # Seed with custom database URL")
	fmt.Println("  go run scripts/seed-data.go -db \"postgres://user:pass@localhost/db\"")
	fmt.Println()
	fmt.Println("NOTES:")
	fmt.Println("  - Seeding will clear existing data in some tables")
	fmt.Println("  - Use the same seed value to generate identical datasets")
	fmt.Println("  - Large datasets may take several minutes to generate")
	fmt.Println("  - Ensure database is running and accessible before seeding")
}

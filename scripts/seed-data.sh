#!/bin/bash

# Sports Prediction Contests - Data Seeding Script
# This script provides a convenient wrapper for seeding the database with fake data

set -e

# Default values
SIZE="small"
SEED="42"
TEST_MODE=false
HELP=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show help
show_help() {
    cat << EOF
Sports Prediction Contests - Data Seeding Script
===============================================

This script populates the database with realistic fake data for development and testing.

USAGE:
    ./scripts/seed-data.sh [OPTIONS]

OPTIONS:
    --size SIZE         Data size preset: small, medium, large (default: small)
                        small:  20 users, 8 contests, 200 predictions
                        medium: 100 users, 25 contests, 1000 predictions  
                        large:  500 users, 50 contests, 5000 predictions

    --seed SEED         Random seed for reproducible data (default: 42)

    --test              Test mode - validate setup without seeding data

    --help              Show this help message

ENVIRONMENT VARIABLES:
    DATABASE_URL        Database connection string
    SEED_SIZE          Data size preset (small/medium/large)
    SEED_VALUE         Random seed value
    BATCH_SIZE         Batch size for database operations (default: 100)

EXAMPLES:
    # Seed with small dataset
    ./scripts/seed-data.sh

    # Seed with medium dataset and custom seed
    ./scripts/seed-data.sh --size medium --seed 123

    # Test configuration without seeding
    ./scripts/seed-data.sh --test

PREREQUISITES:
    - PostgreSQL database running
    - Go 1.21+ installed
    - All backend dependencies available

NOTES:
    - Seeding will add data to existing tables
    - Use the same seed value to generate identical datasets
    - Large datasets may take several minutes to generate
    - Ensure database is accessible before running
EOF
}

# Function to check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi

    # Check Go version
    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
    MAJOR_VERSION=$(echo $GO_VERSION | cut -d. -f1)
    MINOR_VERSION=$(echo $GO_VERSION | cut -d. -f2)
    
    if [ "$MAJOR_VERSION" -lt 1 ] || ([ "$MAJOR_VERSION" -eq 1 ] && [ "$MINOR_VERSION" -lt 21 ]); then
        print_error "Go 1.21+ is required, found version $GO_VERSION"
        exit 1
    fi

    print_success "Go $GO_VERSION found"

    # Check if we're in the right directory
    if [ ! -f "backend/shared/go.mod" ]; then
        print_error "Must be run from project root directory"
        print_error "Expected to find backend/shared/go.mod"
        exit 1
    fi

    # Check database connection
    if [ -z "$DATABASE_URL" ]; then
        DATABASE_URL="postgres://sports_user:sports_password@localhost:5432/sports_prediction?sslmode=disable"
        print_warning "DATABASE_URL not set, using development default: $DATABASE_URL"
        print_warning "⚠️  NEVER use this default in production environments!"
    fi

    print_success "Prerequisites check passed"
}

# Function to wait for database
wait_for_database() {
    print_info "Waiting for database to be ready..."
    
    # Extract host and port from DATABASE_URL
    # This is a simplified parser - works for standard postgres URLs
    HOST=$(echo $DATABASE_URL | sed -n 's/.*@\([^:]*\):.*/\1/p')
    PORT=$(echo $DATABASE_URL | sed -n 's/.*:\([0-9]*\)\/.*/\1/p')
    
    if [ -z "$HOST" ] || [ -z "$PORT" ]; then
        HOST="localhost"
        PORT="5432"
        print_warning "Could not parse DATABASE_URL, using default host:port ($HOST:$PORT)"
    fi

    # Wait for database to be ready (max 30 seconds)
    for i in {1..30}; do
        if nc -z "$HOST" "$PORT" 2>/dev/null; then
            print_success "Database is ready"
            return 0
        fi
        print_info "Waiting for database... ($i/30)"
        sleep 1
    done

    print_error "Database is not ready after 30 seconds"
    print_error "Please ensure PostgreSQL is running on $HOST:$PORT"
    exit 1
}

# Function to run the seeder
run_seeder() {
    print_info "Starting data seeding..."
    print_info "Size: $SIZE, Seed: $SEED, Test Mode: $TEST_MODE"

    # Set environment variables
    export SEED_SIZE="$SIZE"
    export SEED_VALUE="$SEED"
    export DATABASE_URL="$DATABASE_URL"

    # Build and run the seeder
    cd backend/shared
    
    if [ "$TEST_MODE" = true ]; then
        print_info "Running in test mode..."
        go run ./cmd/seeder/main.go -test -size "$SIZE" -seed "$SEED"
    else
        print_info "Running data seeding..."
        go run ./cmd/seeder/main.go -size "$SIZE" -seed "$SEED"
    fi

    cd ../..
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --size)
            SIZE="$2"
            shift 2
            ;;
        --seed)
            SEED="$2"
            shift 2
            ;;
        --test)
            TEST_MODE=true
            shift
            ;;
        --help)
            HELP=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Show help if requested
if [ "$HELP" = true ]; then
    show_help
    exit 0
fi

# Validate size parameter
if [[ ! "$SIZE" =~ ^(small|medium|large)$ ]]; then
    print_error "Invalid size: $SIZE"
    print_error "Valid sizes: small, medium, large"
    exit 1
fi

# Validate seed parameter
if ! [[ "$SEED" =~ ^[0-9]+$ ]]; then
    print_error "Invalid seed: $SEED"
    print_error "Seed must be a positive integer"
    exit 1
fi

# Main execution
main() {
    echo "Sports Prediction Contests - Data Seeder"
    echo "========================================"
    echo

    check_prerequisites
    
    if [ "$TEST_MODE" = false ]; then
        wait_for_database
    fi
    
    run_seeder
    
    echo
    print_success "Seeding completed successfully!"
    
    if [ "$TEST_MODE" = false ]; then
        echo
        print_info "Next steps:"
        print_info "1. Start the frontend: make dev"
        print_info "2. Open http://localhost:3000"
        print_info "3. Explore the populated platform!"
    fi
}

# Run main function
main "$@"

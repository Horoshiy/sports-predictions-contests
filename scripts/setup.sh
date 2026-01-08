#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running in dry-run mode
DRY_RUN=false
if [[ "$1" == "--dry-run" ]]; then
    DRY_RUN=true
    log_info "Running in dry-run mode - no changes will be made"
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check dependencies
check_dependencies() {
    log_info "Checking dependencies..."
    
    local missing_deps=()
    
    # Check Go
    if ! command_exists go; then
        missing_deps+=("go (1.21+)")
    else
        GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
        log_success "Go $GO_VERSION found"
    fi
    
    # Check Node.js
    if ! command_exists node; then
        missing_deps+=("node (18+)")
    else
        NODE_VERSION=$(node --version | sed 's/v//')
        log_success "Node.js $NODE_VERSION found"
    fi
    
    # Check Docker
    if ! command_exists docker; then
        missing_deps+=("docker")
    else
        log_success "Docker found"
    fi
    
    # Check Docker Compose
    if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
        missing_deps+=("docker-compose")
    else
        log_success "Docker Compose found"
    fi
    
    # Check Protocol Buffers compiler
    if ! command_exists protoc; then
        log_warning "protoc not found - Protocol Buffers generation will not work"
    else
        log_success "Protocol Buffers compiler found"
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo ""
        echo "Please install the missing dependencies and run this script again."
        echo "Installation guides:"
        echo "  - Go: https://golang.org/doc/install"
        echo "  - Node.js: https://nodejs.org/en/download/"
        echo "  - Docker: https://docs.docker.com/get-docker/"
        echo "  - Protocol Buffers: https://grpc.io/docs/protoc-installation/"
        exit 1
    fi
}

# Function to setup environment file
setup_env() {
    log_info "Setting up environment file..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY-RUN] Would copy .env.example to .env"
        return
    fi
    
    if [ ! -f .env ]; then
        cp .env.example .env
        log_success "Created .env file from template"
        log_warning "Please review and update .env file with your specific configuration"
    else
        log_info ".env file already exists, skipping"
    fi
}

# Function to initialize Go modules
setup_go() {
    log_info "Setting up Go modules..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY-RUN] Would initialize Go workspace and modules"
        return
    fi
    
    if command_exists go; then
        cd backend
        if [ -f go.work ]; then
            go work sync
            log_success "Go workspace synchronized"
        fi
        
        if [ -f shared/go.mod ]; then
            cd shared
            go mod tidy
            log_success "Shared module dependencies updated"
            cd ..
        fi
        cd ..
    else
        log_warning "Go not found, skipping Go setup"
    fi
}

# Function to setup frontend
setup_frontend() {
    log_info "Setting up frontend..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY-RUN] Would install frontend dependencies"
        return
    fi
    
    if command_exists npm; then
        cd frontend
        npm install
        log_success "Frontend dependencies installed"
        cd ..
    else
        log_warning "npm not found, skipping frontend setup"
    fi
}

# Function to create initial database script
create_db_script() {
    log_info "Creating database initialization script..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY-RUN] Would create database initialization script"
        return
    fi
    
    cat > scripts/init-db.sql << 'EOF'
-- Sports Prediction Contests Database Initialization

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS contests;
CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS sports;
CREATE SCHEMA IF NOT EXISTS predictions;

-- Set default search path
ALTER DATABASE sports_prediction SET search_path TO public, contests, users, sports, predictions;

-- Create basic tables (will be expanded by migrations)
CREATE TABLE IF NOT EXISTS users.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contests.contests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users.users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users.users(email);
CREATE INDEX IF NOT EXISTS idx_contests_created_by ON contests.contests(created_by);
EOF
    
    log_success "Database initialization script created"
}

# Main setup function
main() {
    log_info "Starting Sports Prediction Contests development environment setup..."
    echo ""
    
    check_dependencies
    echo ""
    
    setup_env
    echo ""
    
    setup_go
    echo ""
    
    setup_frontend
    echo ""
    
    create_db_script
    echo ""
    
    if [[ "$DRY_RUN" == "false" ]]; then
        log_success "Setup completed successfully!"
        echo ""
        log_info "Next steps:"
        echo "  1. Review and update .env file with your configuration"
        echo "  2. Start development environment: make dev"
        echo "  3. Check service status: make status"
        echo "  4. View logs: make logs"
        echo ""
        log_info "For help with available commands: make help"
    else
        log_info "Dry-run completed - no changes were made"
    fi
}

# Run main function
main "$@"

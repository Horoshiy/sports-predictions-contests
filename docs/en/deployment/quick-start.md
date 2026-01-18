# Quick Start Guide

Get the Sports Prediction Contests platform up and running in minutes with this comprehensive quick start guide.

## Prerequisites

Before starting, ensure you have the following installed on your system:

### Required Software
- **Docker** (20.10+) and **Docker Compose** (2.0+)
- **Go** (1.21+) for backend development
- **Node.js** (18+) and **npm** (8+) for frontend development
- **Git** for version control

### System Requirements
- **RAM**: Minimum 4GB, recommended 8GB
- **Storage**: At least 2GB free space
- **OS**: Linux, macOS, or Windows with WSL2

### Verification Commands
```bash
# Check Docker installation
docker --version
docker-compose --version

# Check Go installation
go version

# Check Node.js installation
node --version
npm --version

# Check Git installation
git --version
```

## Installation Steps

### Step 1: Clone the Repository

```bash
# Clone the repository
git clone https://github.com/coleam00/dynamous-kiro-hackathon
cd dynamous-kiro-hackathon

# Verify you're in the correct directory
ls -la
```

### Step 2: Environment Setup

```bash
# Run the automated setup script
make setup

# This command will:
# - Copy .env.example to .env
# - Install Go dependencies
# - Install Node.js dependencies
# - Set up development environment
```

### Step 3: Start Infrastructure Services

```bash
# Start PostgreSQL and Redis
make docker-up

# Verify services are running
docker-compose ps
```

Expected output:
```
NAME                COMMAND                  SERVICE             STATUS              PORTS
sports-postgres     "docker-entrypoint.s…"   postgres            running             0.0.0.0:5432->5432/tcp
sports-redis        "docker-entrypoint.s…"   redis               running             0.0.0.0:6379->6379/tcp
```

### Step 4: Database Initialization

The database will be automatically initialized with the required schema when PostgreSQL starts. You can verify the connection:

```bash
# Test database connection
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction
```

### Step 5: Start All Services (Development Mode)

```bash
# Start all microservices and frontend
make dev

# This will start:
# - All 7 microservices
# - React frontend
# - Telegram bot (if configured)
```

### Step 6: Verify Installation

#### Check Service Health
```bash
# Check API Gateway health
curl http://localhost:8080/health

# Expected response: {"status": "ok", "timestamp": "..."}
```

#### Check Individual Services
```bash
# User Service
curl http://localhost:8080/v1/auth/health

# Contest Service  
curl http://localhost:8080/v1/contests/health

# Prediction Service
curl http://localhost:8080/v1/predictions/health

# Scoring Service
curl http://localhost:8080/v1/scores/health

# Sports Service
curl http://localhost:8080/v1/sports/health

# Notification Service
curl http://localhost:8080/v1/notifications/health
```

#### Access the Frontend
Open your browser and navigate to:
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080

## Basic Usage

### Create Your First User

```bash
# Register a new user
# ⚠️ SECURITY NOTE: Use strong passwords in production (min 12 chars, mixed case, numbers, symbols)
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "SecureP@ssw0rd2026!",
    "full_name": "Test User"
  }'
```

### Login and Get Token

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecureP@ssw0rd2026!"
  }'

# Save the returned token for subsequent requests
export JWT_TOKEN="your_jwt_token_here"
```

### Create Your First Contest

```bash
# Create a contest
curl -X POST http://localhost:8080/v1/contests \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "title": "Premier League Predictions",
    "description": "Predict Premier League match outcomes",
    "sport_type": "football",
    "rules": "{\"scoring_system\": \"standard\", \"max_predictions_per_user\": 10}",
    "start_date": "2026-01-20T00:00:00Z",
    "end_date": "2026-05-20T23:59:59Z",
    "max_participants": 100
  }'
```

## Development Commands

### Essential Make Commands
```bash
# Show all available commands
make help

# Setup development environment
make setup

# Start development environment (infrastructure only)
make dev

# Start all services including microservices
make docker-services

# Build all services
make build

# Run tests
make test

# Run E2E tests
make e2e-test

# Clean up
make clean

# Show service status
make status

# View logs
make logs
```

### Docker Commands
```bash
# Start infrastructure only
docker-compose up -d postgres redis

# Start all services
docker-compose --profile services up -d

# Stop all services
docker-compose --profile services down

# View logs
docker-compose logs -f [service-name]

# Restart a service
docker-compose restart [service-name]
```

## Configuration

### Environment Variables

Key environment variables in `.env`:

```bash
# Database
DATABASE_URL=postgres://sports_user:sports_password@localhost:5432/sports_prediction?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# API Gateway
API_GATEWAY_PORT=8080

# JWT
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION=24h

# Telegram Bot (optional)
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
TELEGRAM_ENABLED=false
```

For a complete list, see [Environment Variables Reference](environment-variables.md).

## Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Check what's using the port
lsof -i :8080

# Kill the process if needed
kill -9 <PID>
```

#### Database Connection Issues
```bash
# Check PostgreSQL status
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

#### Service Not Starting
```bash
# Check service logs
docker-compose logs [service-name]

# Rebuild service
docker-compose build [service-name]
docker-compose up -d [service-name]
```

#### Frontend Not Loading
```bash
# Check frontend logs
docker-compose logs frontend

# Rebuild frontend
cd frontend
npm install
npm run build
```

### Getting Help

If you encounter issues:

1. Check the [Common Issues Guide](../troubleshooting/common-issues.md)
2. Review service logs: `docker-compose logs [service-name]`
3. Verify environment variables in `.env`
4. Ensure all prerequisites are installed
5. Try restarting services: `make clean && make dev`

## Next Steps

Once you have the platform running:

1. **Explore the API**: Check out the [API Documentation](../api/services-overview.md)
2. **Run Tests**: Execute the test suite with `make test`
3. **Configure Telegram Bot**: Set up the Telegram integration
4. **Customize Settings**: Modify environment variables for your needs
5. **Deploy to Production**: Follow the [Production Deployment Guide](production.md)

## Validation Checklist

- [ ] All prerequisites installed and verified
- [ ] Repository cloned successfully
- [ ] Environment setup completed (`make setup`)
- [ ] Infrastructure services running (PostgreSQL, Redis)
- [ ] All microservices healthy
- [ ] Frontend accessible at http://localhost:3000
- [ ] API Gateway responding at http://localhost:8080
- [ ] User registration and login working
- [ ] Contest creation successful

---

**Congratulations!** You now have a fully functional Sports Prediction Contests platform running locally. Explore the [API documentation](../api/services-overview.md) to learn more about available features.

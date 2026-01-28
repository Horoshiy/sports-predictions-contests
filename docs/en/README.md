# Sports Prediction Contests - English Documentation

Welcome to the comprehensive English documentation for the Sports Prediction Contests platform.

## Table of Contents

### 🚀 Getting Started
- [Quick Start Guide](deployment/quick-start.md) - Get up and running in minutes
- [Production Deployment](deployment/production.md) - Production-ready setup
- [Environment Variables](deployment/environment-variables.md) - Complete configuration reference

### 📖 API Documentation
- [Services Overview](api/services-overview.md) - All microservices and endpoints
- [Authentication](api/authentication.md) - JWT authentication and authorization
- [User Service](api/user-service.md) - User management and authentication
- [Contest Service](api/contest-service.md) - Contest and team management
- [Prediction Service](api/prediction-service.md) - Predictions and events
- [Scoring Service](api/scoring-service.md) - Scoring and leaderboards
- [Sports Service](api/sports-service.md) - Sports data and synchronization
- [Notification Service](api/notification-service.md) - Multi-channel notifications

### 🧪 Testing
- [E2E Testing](testing/e2e-testing.md) - Backend end-to-end testing procedures
- [Playwright Testing](testing/playwright-testing.md) - Frontend E2E testing with Playwright
- [Unit Testing](testing/unit-testing.md) - Unit testing for all services
- [Performance Testing](testing/performance-testing.md) - Load testing and benchmarks

### 🛠️ Development
- [Fake Data Seeding](development/fake-data-seeding.md) - Generate realistic test data for development

### 🔧 Troubleshooting
- [Common Issues](troubleshooting/common-issues.md) - Frequently encountered problems
- [Diagnostic Tools](troubleshooting/diagnostic-tools.md) - Debugging and diagnostics

## Platform Overview

The Sports Prediction Contests platform is a comprehensive microservices-based system designed for creating and managing sports prediction competitions. It features:

### Core Features
- **Contest Constructor**: Create customizable prediction contests with flexible rules
- **Multi-Sport Support**: Football, basketball, tennis, and more
- **Team Tournaments**: Collaborative team-based competitions
- **Real-time Scoring**: Live leaderboards and streak tracking
- **Props Predictions**: Detailed statistics-based predictions
- **Multi-channel Notifications**: In-app, email, and Telegram notifications

### Architecture
- **7 Microservices**: Modular, scalable architecture
- **gRPC Communication**: High-performance inter-service communication
- **PostgreSQL Database**: Reliable data persistence
- **Redis Caching**: Fast session and leaderboard caching
- **Docker Deployment**: Containerized for easy deployment

### Supported Platforms
- **Web Application**: React-based responsive web interface
- **Telegram Bot**: Interactive bot for contest participation
- **REST API**: Open API for third-party integrations

## Quick Start

1. **Prerequisites**: Ensure you have Docker, Docker Compose, Go 1.21+, and Node.js 18+ installed
2. **Clone Repository**: Get the latest code from the repository
3. **Setup Environment**: Run `make setup` to configure the development environment
4. **Start Services**: Use `make dev` to start the platform
5. **Verify Installation**: Check that all services are running correctly

For detailed instructions, see the [Quick Start Guide](deployment/quick-start.md).

## API Overview

The platform exposes a comprehensive REST API through the API Gateway on port 8080:

| Service | Base URL | Purpose |
|---------|----------|---------|
| API Gateway | `http://localhost:8080` | Main entry point |
| User Management | `/v1/auth/*`, `/v1/users/*` | Authentication and user operations |
| Contests | `/v1/contests/*` | Contest and team management |
| Predictions | `/v1/predictions/*`, `/v1/events/*` | Prediction submission and events |
| Scoring | `/v1/scores/*`, `/v1/leaderboard/*` | Scoring and rankings |
| Sports Data | `/v1/sports/*`, `/v1/leagues/*` | Sports information |
| Notifications | `/v1/notifications/*` | Notification management |

## Development Workflow

1. **Load Context**: Use `@prime` to understand the codebase
2. **Plan Features**: Use `@plan-feature` for new functionality
3. **Implement**: Use `@execute` for systematic development
4. **Review**: Use `@code-review` for quality assurance
5. **Test**: Run comprehensive test suites
6. **Deploy**: Use Docker Compose for deployment

## Support and Resources

- **Troubleshooting**: Check the [common issues guide](troubleshooting/common-issues.md)
- **API Reference**: Complete endpoint documentation in [services overview](api/services-overview.md)
- **Testing**: Comprehensive testing guides in the [testing section](testing/)
- **Deployment**: Production deployment guides in [deployment section](deployment/)

## Language Options

- **English**: You are here
- **Русский**: [Русская документация](../ru/README.md)

---

*For the latest updates and community support, visit the project repository.*

# Fake Data Seeding System

The Sports Prediction Contests platform includes a comprehensive fake data seeding system that populates your development database with realistic test data. This system generates users, contests, predictions, sports data, and all related entities with proper relationships and realistic values.

## Overview

The seeding system provides:

- **Realistic Data Generation**: Uses the gofakeit library to create authentic-looking sports data
- **Configurable Data Volumes**: Three preset sizes (small/medium/large) for different development needs
- **Referential Integrity**: Ensures all foreign key relationships are properly maintained
- **Reproducible Results**: Use the same seed value to generate identical datasets
- **Cross-Service Coordination**: Populates data across all microservices consistently

## Quick Start

### Prerequisites

- PostgreSQL database running (localhost:5432 by default)
- Go 1.21+ installed
- Project dependencies installed (`make setup`)

### Basic Usage

```bash
# Seed with small dataset (recommended for development)
make seed-small

# Seed with medium dataset (good for testing)
make seed-medium

# Seed with large dataset (performance testing)
make seed-large

# Test configuration without seeding
make seed-test
```

### Custom Seeding

```bash
# Use shell script directly with custom parameters
./scripts/seed-data.sh --size medium --seed 123

# Set environment variables for custom seeding
SEED_SIZE=medium SEED_VALUE=456 make seed-custom
```

## Data Size Presets

### Small Dataset (Default)
- **Users**: 20 (with profiles and preferences)
- **Sports**: 3 (Football, Soccer, Basketball)
- **Leagues**: 6 (2 per sport)
- **Teams**: 24 (8 per sport)
- **Matches**: 50 (various schedules and statuses)
- **Contests**: 8 (different sports and configurations)
- **Predictions**: 200 (realistic user predictions)
- **User Teams**: 4 (team tournaments)

**Use Case**: Daily development, quick testing, feature development

### Medium Dataset
- **Users**: 100
- **Sports**: 5
- **Leagues**: 15
- **Teams**: 60
- **Matches**: 200
- **Contests**: 25
- **Predictions**: 1000
- **User Teams**: 15

**Use Case**: Integration testing, UI stress testing, demo preparation

### Large Dataset
- **Users**: 500
- **Sports**: 8
- **Leagues**: 30
- **Teams**: 120
- **Matches**: 800
- **Contests**: 50
- **Predictions**: 5000
- **User Teams**: 30

**Use Case**: Performance testing, load testing, production simulation

## Generated Data Types

### Core Entities

#### Users
- Realistic names and email addresses
- Secure password hashes (bcrypt)
- Complete user profiles with bio, avatar, location, social links
- User preferences (language, timezone, notifications, theme)

#### Sports Data
- Authentic sport names and descriptions
- Realistic team names (city + mascot combinations)
- Professional league structures
- Match schedules with realistic timing
- Match results with detailed statistics

#### Contests
- Varied contest titles and descriptions
- Different sports types and rules
- Realistic date ranges and participant limits
- Multiple status types (draft, active, completed)

#### Predictions
- Multiple prediction types (match result, exact score, over/under, etc.)
- Realistic prediction data in JSON format
- Proper user-contest-match relationships
- Scoring results with correctness indicators

### Advanced Features

#### Scoring System
- Individual prediction scores with time coefficients
- Leaderboard entries with rankings
- User streak tracking (current and maximum streaks)
- Realistic point distributions

#### Team Tournaments
- User-created teams with captains and members
- Team invite codes and membership management
- Team-based contest participation

#### Notifications
- Various notification types (contest updates, results, achievements)
- Multi-channel preferences (in-app, email, Telegram)
- Read/unread status tracking
- Realistic notification history

## Configuration

### Environment Variables

```bash
# Data generation settings
SEED_SIZE=small          # Data size preset (small/medium/large)
SEED_VALUE=42           # Random seed for reproducible data
BATCH_SIZE=100          # Database batch insert size

# Database connection
DATABASE_URL=postgres://sports_user:sports_password@localhost:5432/sports_prediction?sslmode=disable
```

### Command Line Options

```bash
./scripts/seed-data.sh [OPTIONS]

Options:
  --size SIZE    Data size preset: small, medium, large
  --seed SEED    Random seed for reproducible data
  --test         Test mode - validate without seeding
  --help         Show help message
```

## Implementation Details

### Architecture

The seeding system consists of several components:

- **Config**: Manages seeding configuration and data size presets
- **Factory**: Generates fake data using gofakeit library
- **SportsDataGenerator**: Specialized generator for sports-specific data
- **Coordinator**: Orchestrates seeding across all services with proper dependencies
- **Models**: Seeding-specific data structures mirroring service models

### Dependency Order

Data is seeded in the correct dependency order to maintain referential integrity:

1. **Users** (referenced by everything)
2. **Sports** (referenced by leagues, teams, contests)
3. **Leagues** (depends on sports)
4. **Teams** (depends on sports)
5. **Matches** (depends on leagues and teams)
6. **Contests** (depends on users and sports)
7. **Predictions** (depends on users, contests, matches)
8. **Scoring Data** (depends on predictions)
9. **User Teams** (depends on users)
10. **Notifications** (depends on users)

### Transaction Safety

All seeding operations are wrapped in database transactions to ensure:
- **Atomicity**: Either all data is seeded or none
- **Consistency**: Database remains in valid state if seeding fails
- **Rollback**: Automatic cleanup on errors

## Usage Examples

### Development Workflow

```bash
# Start fresh development environment
make docker-up
make seed-small

# Start frontend and explore
make dev
# Open http://localhost:3000
```

### Testing Different Scenarios

```bash
# Generate consistent test data
./scripts/seed-data.sh --seed 123

# Later, regenerate identical data
./scripts/seed-data.sh --seed 123
```

### Performance Testing

```bash
# Generate large dataset for load testing
make seed-large

# Monitor database performance
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "
  SELECT schemaname,tablename,n_tup_ins,n_tup_upd,n_tup_del 
  FROM pg_stat_user_tables 
  ORDER BY n_tup_ins DESC;
"
```

## Troubleshooting

### Common Issues

#### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Start database if needed
make docker-up

# Test connection manually
psql postgres://sports_user:sports_password@localhost:5432/sports_prediction
```

#### Go Module Issues
```bash
# Ensure you're in project root
pwd  # Should show project root directory

# Update dependencies
cd backend/shared && go mod tidy
```

#### Permission Denied on Script
```bash
# Make script executable
chmod +x scripts/seed-data.sh
```

### Validation Commands

```bash
# Test seeding configuration
make seed-test

# Check generated data counts
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "
  SELECT 
    'users' as table_name, COUNT(*) as count FROM users
  UNION ALL
  SELECT 'contests', COUNT(*) FROM contests
  UNION ALL  
  SELECT 'predictions', COUNT(*) FROM predictions;
"

# Verify data relationships
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "
  SELECT c.title, COUNT(p.id) as prediction_count 
  FROM contests c 
  LEFT JOIN predictions p ON c.id = p.contest_id 
  GROUP BY c.id, c.title 
  ORDER BY prediction_count DESC;
"
```

### Performance Monitoring

```bash
# Monitor seeding progress
tail -f /tmp/seeding.log

# Check database size after seeding
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "
  SELECT pg_size_pretty(pg_database_size('sports_prediction')) as database_size;
"
```

## Best Practices

### Development

1. **Use Small Dataset**: Start with `make seed-small` for daily development
2. **Consistent Seeds**: Use the same seed value across team members for consistent test data
3. **Regular Refresh**: Re-seed periodically to test with fresh data
4. **Test Mode First**: Use `make seed-test` to validate configuration before seeding

### Testing

1. **Isolated Environments**: Use different databases for different test scenarios
2. **Reproducible Tests**: Use fixed seed values for consistent test results
3. **Performance Baselines**: Use large datasets to establish performance benchmarks
4. **Data Validation**: Verify data integrity after seeding

### Production Considerations

⚠️ **Warning**: Never run the seeding system against production databases!

- Seeding is designed for development and testing only
- Generated data includes fake email addresses and passwords
- Use environment variables to prevent accidental production seeding
- Consider data privacy regulations when using generated data

## Integration with Development Workflow

### Docker Compose Integration

The seeding system integrates seamlessly with the existing Docker development workflow:

```bash
# Start infrastructure
make docker-up

# Seed data
make seed-small

# Start all services
make dev
```

### CI/CD Integration

For automated testing environments:

```bash
# In CI pipeline
make docker-up
make seed-test  # Validate configuration
make seed-small # Generate test data
make e2e-test   # Run tests against seeded data
```

### Team Collaboration

Share consistent development environments:

```bash
# Team lead generates data
SEED_VALUE=12345 make seed-medium

# Team members use same seed
SEED_VALUE=12345 make seed-medium
```

## Advanced Usage

### Custom Data Generation

For specialized testing scenarios, you can extend the seeding system:

1. **Custom Sports**: Add new sports to the predefined list in `sports_data.go`
2. **Custom Rules**: Modify contest rule generation in `factory.go`
3. **Custom Relationships**: Adjust prediction-contest relationships for specific test cases

### Performance Optimization

For large datasets:

```bash
# Increase batch size for faster insertion
BATCH_SIZE=500 make seed-large

# Monitor memory usage during seeding
top -p $(pgrep seed-data)
```

### Data Export/Import

```bash
# Export seeded data for sharing
pg_dump -U sports_user -h localhost sports_prediction > seeded_data.sql

# Import on another machine
psql -U sports_user -h localhost sports_prediction < seeded_data.sql
```

## Support

For issues with the seeding system:

1. Check the [Troubleshooting](#troubleshooting) section
2. Verify prerequisites and configuration
3. Review logs for specific error messages
4. Consult the main project documentation

The fake data seeding system is designed to make development and testing easier by providing realistic, consistent data across all platform features. Use it regularly to ensure your development environment stays populated with meaningful test data.

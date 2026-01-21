# Feature: Comprehensive Fake Data Seeding System

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement a comprehensive fake data seeding system that populates all microservices with realistic test data on first application launch. The system will generate users, contests, predictions, sports data, teams, matches, and all related entities with proper relationships and realistic values. Include bilingual documentation (English/Russian) explaining the seeding system, data structure, and usage.

## User Story

As a developer or demo user
I want the application to be populated with realistic fake data on first launch
So that I can immediately explore all features, test functionality, and demonstrate the platform without manually creating data

## Problem Statement

Currently, the Sports Prediction Contests platform starts with minimal data (only basic prop types for Soccer), making it difficult to:
- Demonstrate platform capabilities to stakeholders
- Test features that require existing data relationships
- Develop frontend components with realistic data
- Onboard new developers who need populated environments
- Showcase the full user experience without manual data entry

## Solution Statement

Create a comprehensive fake data seeding system using Go's gofakeit library that generates realistic sports prediction data across all microservices. The system will:
- Generate interconnected data with proper foreign key relationships
- Support configurable data volumes (small/medium/large datasets)
- Provide realistic sports-specific data (teams, leagues, matches, predictions)
- Include bilingual documentation for international users
- Integrate seamlessly with existing Docker development workflow

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: All backend services, documentation, development workflow
**Dependencies**: gofakeit library, existing database schema, Docker Compose

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `scripts/init-db.sql` (lines 1-300) - Why: Current database schema and existing seeding patterns for prop types
- `backend/user-service/internal/models/user.go` (lines 13-21) - Why: User model structure with Profile and Preferences relationships
- `backend/contest-service/internal/models/contest.go` (lines 12-25) - Why: Contest model with rules, dates, and participant tracking
- `backend/prediction-service/internal/models/prediction.go` - Why: Prediction data structure and JSONB fields
- `backend/sports-service/internal/models/sport.go` - Why: Sports, leagues, teams, matches model relationships
- `backend/scoring-service/internal/models/streak.go` - Why: UserStreak model for gamification features
- `tests/e2e/helpers.go` (lines 75-85) - Why: Existing test data generation patterns to follow
- `docker-compose.yml` (lines 1-50) - Why: Database initialization and service dependencies
- `Makefile` (lines 1-50) - Why: Existing development commands and setup patterns
- `docs/en/deployment/quick-start.md` - Why: Current setup documentation structure to extend

### New Files to Create

- `backend/shared/seeder/factory.go` - Central data factory with gofakeit integration
- `backend/shared/seeder/coordinator.go` - Cross-service data coordination and relationships
- `backend/shared/seeder/config.go` - Seeding configuration and data size management
- `scripts/seed-data.go` - Standalone seeding executable
- `scripts/seed-data.sh` - Shell script for Docker integration
- `docs/en/development/fake-data-seeding.md` - English documentation
- `docs/ru/development/fake-data-seeding.md` - Russian documentation
- `backend/shared/seeder/sports_data.go` - Sports-specific realistic data generation
- `backend/shared/seeder/models.go` - Seeding-specific data structures

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [gofakeit Documentation](https://github.com/brianvoe/gofakeit#readme)
  - Specific section: Data generation methods and sports data
  - Why: Primary library for generating realistic fake data
- [GORM Documentation](https://gorm.io/docs/create.html#batch-insert)
  - Specific section: Batch insert and transaction handling
  - Why: Efficient database seeding with proper relationships
- [Docker Compose Documentation](https://docs.docker.com/compose/startup-order/)
  - Specific section: Service dependencies and initialization
  - Why: Proper integration with existing development workflow

### Patterns to Follow

**Naming Conventions:**
- Go files: `snake_case.go`
- Structs: `PascalCase`
- Functions: `PascalCase` (public), `camelCase` (private)
- Database tables: Follow existing GORM model patterns

**Error Handling:**
```go
if err != nil {
    return fmt.Errorf("failed to seed users: %w", err)
}
```

**Logging Pattern:**
```go
log.Printf("Seeding %d users...", count)
log.Printf("Successfully seeded %d users", len(users))
```

**GORM Transaction Pattern:**
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()
if err := tx.Error; err != nil {
    return err
}
// ... operations
return tx.Commit().Error
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation Setup

Set up the core seeding infrastructure with gofakeit integration and configuration management.

**Tasks:**
- Install gofakeit dependency in shared module
- Create seeding configuration structure
- Set up data factory with realistic generators
- Create coordinator for cross-service data relationships

### Phase 2: Core Data Generation

Implement realistic data generators for all major entities with proper relationships.

**Tasks:**
- Generate users with profiles and preferences
- Create sports, leagues, teams with realistic data
- Generate contests with various configurations
- Create matches with realistic schedules and results
- Generate predictions with proper user-contest relationships

### Phase 3: Advanced Features Integration

Add support for advanced platform features like streaks, analytics, and team tournaments.

**Tasks:**
- Generate user streaks and scoring data
- Create team tournament data
- Add notification preferences and history
- Generate props predictions data
- Create realistic analytics data

### Phase 4: Documentation and Integration

Create comprehensive bilingual documentation and integrate with development workflow.

**Tasks:**
- Write English documentation with examples
- Create Russian documentation translation
- Integrate seeding with Docker Compose
- Add Makefile commands for data management
- Create validation and testing procedures

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/shared/seeder/config.go

- **IMPLEMENT**: Seeding configuration with data size presets and environment variables
- **PATTERN**: Follow existing config patterns in backend services
- **IMPORTS**: Required imports for environment handling and validation
- **GOTCHA**: Use environment variables with sensible defaults for different deployment scenarios
- **VALIDATE**: `go build backend/shared/seeder/config.go`

### UPDATE backend/shared/go.mod

- **IMPLEMENT**: Add gofakeit dependency to shared module
- **PATTERN**: Follow existing dependency management in go.mod files
- **IMPORTS**: `github.com/brianvoe/gofakeit/v6`
- **GOTCHA**: Ensure version compatibility with existing Go version (1.21+)
- **VALIDATE**: `cd backend/shared && go mod tidy && go mod verify`

### CREATE backend/shared/seeder/factory.go

- **IMPLEMENT**: Central data factory with gofakeit integration for all entity types
- **PATTERN**: Mirror existing model structures from each service
- **IMPORTS**: gofakeit, time, existing model packages, GORM
- **GOTCHA**: Generate data with proper foreign key relationships and realistic constraints
- **VALIDATE**: `go build backend/shared/seeder/factory.go`

### CREATE backend/shared/seeder/sports_data.go

- **IMPLEMENT**: Sports-specific realistic data generation (teams, leagues, matches)
- **PATTERN**: Use existing sports service model structures
- **IMPORTS**: gofakeit, sports service models, time for match scheduling
- **GOTCHA**: Generate realistic team names, league structures, and match schedules
- **VALIDATE**: `go build backend/shared/seeder/sports_data.go`

### CREATE backend/shared/seeder/models.go

- **IMPLEMENT**: Seeding-specific data structures and helper types
- **PATTERN**: Follow existing model patterns with GORM tags
- **IMPORTS**: GORM, time, existing service models
- **GOTCHA**: Create structures that facilitate cross-service data coordination
- **VALIDATE**: `go build backend/shared/seeder/models.go`

### CREATE backend/shared/seeder/coordinator.go

- **IMPLEMENT**: Cross-service data coordination ensuring referential integrity
- **PATTERN**: Use GORM transaction patterns from existing services
- **IMPORTS**: GORM, all service models, factory, config
- **GOTCHA**: Seed data in dependency order (users → sports → contests → predictions)
- **VALIDATE**: `go build backend/shared/seeder/coordinator.go`

### CREATE scripts/seed-data.go

- **IMPLEMENT**: Standalone executable for running data seeding
- **PATTERN**: Follow existing script patterns in scripts/ directory
- **IMPORTS**: Database connection, seeder coordinator, configuration
- **GOTCHA**: Handle database connection and graceful error handling
- **VALIDATE**: `go build -o scripts/seed-data scripts/seed-data.go`

### CREATE scripts/seed-data.sh

- **IMPLEMENT**: Shell script wrapper for Docker integration
- **PATTERN**: Follow existing shell script patterns in scripts/
- **IMPORTS**: Environment variable handling, Docker service dependencies
- **GOTCHA**: Wait for database availability before seeding
- **VALIDATE**: `chmod +x scripts/seed-data.sh && ./scripts/seed-data.sh --help`

### UPDATE docker-compose.yml

- **IMPLEMENT**: Add data seeding service and integration with existing services
- **PATTERN**: Follow existing service definitions and dependency patterns
- **IMPORTS**: Existing database and service configurations
- **GOTCHA**: Ensure seeding runs after database initialization but before other services
- **VALIDATE**: `docker-compose config`

### UPDATE Makefile

- **IMPLEMENT**: Add seeding commands for development workflow
- **PATTERN**: Follow existing Makefile command patterns and help formatting
- **IMPORTS**: Existing make targets and Docker commands
- **GOTCHA**: Provide commands for different data sizes and reset functionality
- **VALIDATE**: `make help | grep seed`

### CREATE docs/en/development/fake-data-seeding.md

- **IMPLEMENT**: Comprehensive English documentation for seeding system
- **PATTERN**: Follow existing documentation structure in docs/en/
- **IMPORTS**: Code examples, configuration options, troubleshooting
- **GOTCHA**: Include practical examples and common use cases
- **VALIDATE**: Check markdown syntax and internal links

### CREATE docs/ru/development/fake-data-seeding.md

- **IMPLEMENT**: Russian translation of seeding documentation
- **PATTERN**: Mirror English documentation structure with cultural adaptations
- **IMPORTS**: Translated code examples and Russian technical terminology
- **GOTCHA**: Maintain technical accuracy while using appropriate Russian terms
- **VALIDATE**: Check markdown syntax and consistency with English version

### UPDATE docs/en/README.md

- **IMPLEMENT**: Add seeding system to main English documentation index
- **PATTERN**: Follow existing documentation linking and structure
- **IMPORTS**: Link to new seeding documentation
- **GOTCHA**: Integrate seamlessly with existing documentation flow
- **VALIDATE**: Check all links work and structure is consistent

### UPDATE docs/ru/README.md

- **IMPLEMENT**: Add seeding system to main Russian documentation index
- **PATTERN**: Follow existing Russian documentation structure
- **IMPORTS**: Link to new Russian seeding documentation
- **GOTCHA**: Maintain consistency with English version while using proper Russian
- **VALIDATE**: Check all links work and structure is consistent

### UPDATE README.md

- **IMPLEMENT**: Add seeding information to main project README
- **PATTERN**: Follow existing README structure and command examples
- **IMPORTS**: New make commands and quick start information
- **GOTCHA**: Keep main README concise while highlighting key seeding features
- **VALIDATE**: Check markdown rendering and command accuracy

---

## TESTING STRATEGY

### Unit Tests

Design unit tests for data generation functions following existing test patterns in tests/ directory:

- Test data factory functions generate valid data structures
- Test configuration loading and validation
- Test cross-service data relationship integrity
- Test different data size configurations

### Integration Tests

Create integration tests that verify seeding works with actual database:

- Test full seeding process with small dataset
- Verify referential integrity across all services
- Test seeding performance with different data sizes
- Validate data cleanup and reset functionality

### Edge Cases

Test specific edge cases for the seeding system:

- Database connection failures during seeding
- Partial seeding failures and rollback behavior
- Concurrent seeding attempts
- Invalid configuration handling
- Large dataset memory usage

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
cd backend/shared && go fmt ./seeder/...
cd backend/shared && go vet ./seeder/...
go build scripts/seed-data.go
```

### Level 2: Unit Tests

```bash
cd backend/shared && go test ./seeder/... -v
go test scripts/seed-data.go -v
```

### Level 3: Integration Tests

```bash
make docker-up
./scripts/seed-data.sh --size small --test
make e2e-test-only
```

### Level 4: Manual Validation

```bash
# Test seeding with different sizes
make seed-small
make seed-medium
make seed-large

# Verify data in database
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "SELECT COUNT(*) FROM users;"
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "SELECT COUNT(*) FROM contests;"
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "SELECT COUNT(*) FROM predictions;"

# Test frontend with seeded data
make dev
# Open http://localhost:3000 and verify data appears in UI
```

### Level 5: Documentation Validation

```bash
# Check documentation links and formatting
markdownlint docs/en/development/fake-data-seeding.md
markdownlint docs/ru/development/fake-data-seeding.md

# Verify bilingual consistency
diff -u <(grep "^#" docs/en/development/fake-data-seeding.md) <(grep "^#" docs/ru/development/fake-data-seeding.md | wc -l)
```

---

## ACCEPTANCE CRITERIA

- [ ] Seeding system generates realistic data for all microservices
- [ ] Data maintains referential integrity across service boundaries
- [ ] Three data size presets (small/medium/large) work correctly
- [ ] Docker Compose integration allows one-command seeding
- [ ] Makefile provides convenient seeding commands
- [ ] English documentation is comprehensive and accurate
- [ ] Russian documentation is complete and properly translated
- [ ] Seeding performance is acceptable for development use
- [ ] Data can be reset and re-seeded without issues
- [ ] Frontend displays seeded data correctly
- [ ] All validation commands pass with zero errors
- [ ] No regressions in existing functionality

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms seeding works
- [ ] Documentation is complete in both languages
- [ ] Docker integration works seamlessly
- [ ] Performance is acceptable for all data sizes
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

### Design Decisions

1. **gofakeit Library**: Chosen for comprehensive fake data generation with sports-specific capabilities
2. **Shared Module**: Seeding logic in shared module allows reuse across all services
3. **Configuration-Driven**: Environment variables and presets support different deployment scenarios
4. **Transaction Safety**: All seeding operations use database transactions for consistency
5. **Bilingual Documentation**: Supports international development team and users

### Performance Considerations

- Use batch inserts for large datasets to improve performance
- Generate data in dependency order to minimize foreign key constraint issues
- Implement progress logging for long-running seeding operations
- Consider memory usage for large dataset generation

### Security Considerations

- Seeding should only run in development/staging environments
- Generated passwords use secure hashing (bcrypt)
- No real email addresses or sensitive data in generated content
- Clear documentation about not using seeded data in production

### Future Enhancements

- Add support for custom seeding scenarios
- Implement incremental seeding (add data without full reset)
- Add data export/import functionality for sharing datasets
- Create web UI for seeding management

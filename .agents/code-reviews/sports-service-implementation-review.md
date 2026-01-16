# Code Review: Sports Service Implementation

**Date**: 2026-01-15
**Reviewer**: Kiro CLI
**Scope**: Sports Service microservice implementation

## Stats

- Files Modified: 4
- Files Added: 18
- Files Deleted: 0
- New lines: ~1,800
- Deleted lines: 0

---

## Issues Found

### CRITICAL

```
severity: critical
file: backend/sports-service/internal/service/sports_service.go
line: 119-122, 224-227, 329-332, 434-437
issue: Potential nil pointer dereference on pagination request
detail: ListSports, ListLeagues, ListTeams, and ListMatches access req.Pagination.Limit without checking if Pagination is nil. If a client sends a request without pagination, this will cause a panic.
suggestion: Add nil check before accessing pagination fields:
  if req.Pagination == nil {
    req.Pagination = &common.PaginationRequest{Page: 1, Limit: 10}
  }
```

```
severity: critical
file: backend/sports-service/internal/service/sports_service.go
line: 350
issue: Potential nil pointer dereference on CreateMatch
detail: req.ScheduledAt.AsTime() will panic if ScheduledAt is nil. No validation before accessing.
suggestion: Add nil check:
  if req.ScheduledAt == nil {
    return &pb.MatchResponse{Response: &common.Response{Success: false, Message: "scheduled_at is required"...}}, nil
  }
```

### HIGH

```
severity: high
file: backend/sports-service/internal/service/sports_service.go
line: 380
issue: Potential nil pointer dereference on UpdateMatch
detail: req.ScheduledAt.AsTime() called without nil check. Will panic if client doesn't provide scheduled_at.
suggestion: Add nil check before accessing ScheduledAt
```

```
severity: high
file: backend/shared/proto/sports/sports.pb.gw.go
line: 10-17
issue: Stub implementation does nothing - gateway registration will silently fail
detail: RegisterSportsServiceHandlerFromEndpoint returns nil without actually registering the service. The API Gateway will appear to work but sports endpoints won't be accessible via HTTP.
suggestion: This is a placeholder stub. For production, generate proper gRPC-gateway code using protoc with grpc-gateway plugin, or implement actual gRPC client connection and handler registration.
```

```
severity: high
file: backend/sports-service/internal/models/sport.go
line: 45-48
issue: Auto-generated slug may produce invalid slugs
detail: The auto-slug generation `strings.ToLower(strings.ReplaceAll(s.Name, " ", "-"))` doesn't handle special characters. A name like "Fútbol (Soccer)" would produce "fútbol-(soccer)" which fails the slug validation regex `^[a-z0-9-]+$`.
suggestion: Use a proper slugify function that removes/replaces special characters:
  slug := regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(strings.ToLower(s.Name), "-")
  slug = strings.Trim(slug, "-")
```

### MEDIUM

```
severity: medium
file: backend/sports-service/internal/service/sports_service.go
line: 47-65, 155-173, 263-281, 345-363
issue: Create operations don't validate foreign key existence
detail: CreateLeague, CreateTeam, CreateMatch don't verify that referenced SportID/LeagueID/TeamIDs exist before attempting to create. This will result in database foreign key constraint errors with unclear error messages.
suggestion: Add existence checks before create:
  if _, err := s.sportRepo.GetByID(uint(req.SportId)); err != nil {
    return &pb.LeagueResponse{Response: &common.Response{Success: false, Message: "sport not found"...}}, nil
  }
```

```
severity: medium
file: backend/sports-service/internal/models/league.go
line: 51-54
issue: Same auto-slug issue as Sport model
detail: Auto-generated slug doesn't sanitize special characters properly.
suggestion: Apply same fix as Sport model
```

```
severity: medium
file: backend/sports-service/internal/models/team.go
line: 51-54
issue: Same auto-slug issue as Sport and League models
detail: Auto-generated slug doesn't sanitize special characters properly.
suggestion: Apply same fix as Sport model
```

```
severity: medium
file: backend/sports-service/internal/service/sports_service.go
line: 131-134, 236-239, 341-344, 446-449
issue: Division by zero possible in pagination calculation
detail: totalPages calculation divides by req.Pagination.Limit. If Limit is 0 (after the default assignment), this causes integer division by zero panic.
suggestion: The default limit assignment on line 119 handles this, but add explicit check:
  if req.Pagination.Limit <= 0 {
    req.Pagination.Limit = 10
  }
```

```
severity: medium
file: scripts/init-db.sql
line: 57-136
issue: Missing ON DELETE behavior for foreign keys
detail: Foreign key constraints don't specify ON DELETE behavior. Deleting a sport will fail if leagues/teams reference it, with unclear error messages.
suggestion: Add ON DELETE CASCADE or ON DELETE RESTRICT explicitly:
  sport_id INTEGER NOT NULL REFERENCES sports(id) ON DELETE RESTRICT
```

### LOW

```
severity: low
file: backend/sports-service/internal/service/sports_service.go
line: all Create methods
issue: Inconsistent with contest-service pattern - missing auth check
detail: Contest service extracts userID from JWT context for audit/ownership. Sports service Create methods don't extract user context, which may be intentional for admin-only operations but differs from established pattern.
suggestion: Consider adding auth.GetUserIDFromContext() for audit logging, or document that sports management is admin-only.
```

```
severity: low
file: docker-compose.yml
line: 152
issue: Container name typo
detail: Container name is "sports-sports-service" (double "sports")
suggestion: Change to "sports-service" for consistency with other services
```

```
severity: low
file: backend/sports-service/Dockerfile
line: 5
issue: COPY path may fail in Docker build context
detail: `COPY ../shared ../shared` copies from parent directory which may not be in build context depending on how docker build is invoked.
suggestion: Verify build context includes shared directory, or use multi-stage build with proper context. Other services have same pattern so this may work, but verify.
```

---

## Security Analysis

- ✅ No SQL injection vulnerabilities (using GORM parameterized queries)
- ✅ No hardcoded secrets (uses environment variables)
- ✅ JWT authentication interceptor applied
- ⚠️ No authorization checks - any authenticated user can create/modify sports data

## Performance Analysis

- ✅ Proper database indexing in init-db.sql
- ✅ Pagination implemented for list operations
- ✅ Preloading relationships to avoid N+1 queries
- ⚠️ No caching layer (acceptable for initial implementation)

## Code Quality

- ✅ Follows existing codebase patterns
- ✅ Consistent error handling
- ✅ Proper logging format
- ✅ Repository interface pattern maintained
- ⚠️ Some code duplication in pagination handling (could extract helper)

---

## Summary

**Total Issues**: 12
- Critical: 2
- High: 3
- Medium: 5
- Low: 3

**Recommendation**: Fix critical and high issues before merging. The nil pointer dereferences on pagination and timestamps will cause runtime panics. The gateway stub needs proper implementation for HTTP API access.

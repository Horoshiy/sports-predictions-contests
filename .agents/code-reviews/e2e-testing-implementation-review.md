# Code Review: E2E Testing Implementation

**Date**: 2026-01-16
**Reviewer**: Kiro CLI

## Stats

- Files Modified: 1
- Files Added: 11
- Files Deleted: 0
- New lines: ~550
- Deleted lines: 0

---

## Issues Found

### HIGH

```
severity: high
file: tests/e2e/helpers.go
line: 44-47
issue: Response body read after defer close - potential data loss
detail: In parseResponse(), io.ReadAll is called, then defer resp.Body.Close() is set. 
        The defer is placed AFTER reading, which is correct for execution order, but 
        the body is already consumed. However, if the caller also has a defer close, 
        it will try to close an already-closed body (harmless but wasteful).
suggestion: Remove the defer in parseResponse since caller is responsible for closing,
            or document that parseResponse consumes and closes the body.
```

```
severity: high
file: tests/e2e/sports_test.go
line: 14-18
issue: Ignoring error from parseResponse in setup
detail: The error from parseResponse[AuthResponse](resp) is ignored with `_`. If 
        registration fails or returns invalid JSON, token will be empty string and 
        all subsequent tests will fail with confusing errors.
suggestion: Check the error: `authResp, err := parseResponse[AuthResponse](resp); if err != nil { t.Fatalf(...) }`
```

```
severity: high
file: tests/e2e/contest_test.go
line: 14-18
issue: Same ignored error pattern in setup
detail: Same issue as sports_test.go - parseResponse error ignored in setup.
suggestion: Add error handling for parseResponse in setup block.
```

```
severity: high
file: tests/e2e/prediction_test.go
line: 14-18
issue: Same ignored error pattern in setup
detail: Same issue - parseResponse error ignored.
suggestion: Add error handling for parseResponse in setup block.
```

```
severity: high
file: tests/e2e/scoring_test.go
line: 14-18
issue: Same ignored error pattern in setup
detail: Same issue - parseResponse error ignored.
suggestion: Add error handling for parseResponse in setup block.
```

### MEDIUM

```
severity: medium
file: tests/e2e/helpers.go
line: 77
issue: Weak random seed - predictable test emails
detail: math/rand without explicit seeding uses a deterministic seed. In Go 1.20+
        this is auto-seeded, but for Go 1.21 it's good practice to be explicit.
        Could cause test collisions if tests run in quick succession.
suggestion: Add rand.Seed(time.Now().UnixNano()) in init() or use crypto/rand for uniqueness.
```

```
severity: medium
file: tests/e2e/auth_test.go
line: 12-13
issue: Closure variables shared across subtests
detail: Variables registeredEmail, authToken, userID are shared across t.Run() subtests.
        If subtests run in parallel (t.Parallel()), this would cause race conditions.
        Currently safe because no t.Parallel() is called, but fragile.
suggestion: Document that these tests must run sequentially, or restructure to pass 
            values explicitly between subtests.
```

```
severity: medium
file: tests/e2e/workflow_test.go
line: 89-95
issue: Missing status code check after creating teams
detail: After creating home team and away team, resp.StatusCode is not checked.
        If team creation fails, the test continues with zero IDs.
suggestion: Add status code validation after each team creation.
```

```
severity: medium
file: scripts/e2e-test.sh
line: 17
issue: Fixed sleep duration instead of health check for database
detail: `sleep 5` assumes database is ready in 5 seconds. On slow machines or 
        cold starts, PostgreSQL may take longer to initialize.
suggestion: Add a database health check loop similar to the API Gateway check:
            `until pg_isready -h localhost -p 5432; do sleep 1; done`
```

### LOW

```
severity: low
file: tests/e2e/types.go
line: 1
issue: Unused types defined
detail: Several response types are defined but never used in tests:
        ContestsResponse, ParticipantsResponse, EventsResponse, PredictionsResponse,
        SportsResponse, LeaguesResponse, TeamsResponse, MatchesResponse, HealthResponse
suggestion: Remove unused types or add tests that use list endpoints with pagination.
```

```
severity: low
file: tests/e2e/main_test.go
line: 12
issue: Unused import - testing package imported but only m.Run() used
detail: The `testing` import is required for TestMain signature, but the `time` 
        import in main_test.go is unused (time is used in waitForServices via helpers).
suggestion: Remove unused `time` import from main_test.go if not needed.
```

```
severity: low
file: tests/e2e/go.mod
line: 1-3
issue: Missing go.sum file
detail: go.mod exists but go.sum is missing. Running `go mod tidy` will create it.
suggestion: Run `go mod tidy` to generate go.sum before committing.
```

---

## Summary

The E2E test implementation is well-structured and covers the main user workflows. The test isolation via unique emails is good practice. However, there are several instances of ignored errors in setup blocks that could lead to confusing test failures.

**Critical Issues**: 0
**High Issues**: 5 (all related to ignored parseResponse errors)
**Medium Issues**: 4
**Low Issues**: 3

### Recommended Fixes (Priority Order)

1. **Fix ignored errors in all test setup blocks** - Add error handling for parseResponse calls
2. **Add database health check to e2e-test.sh** - Replace fixed sleep with proper check
3. **Add status checks in workflow_test.go** - Validate team creation responses
4. **Run `go mod tidy`** - Generate go.sum file

### Positive Observations

- Good use of build tags (`//go:build e2e`) for test separation
- Proper test isolation with unique emails/names
- Comprehensive coverage of API endpoints
- Clear test naming and logging
- Proper cleanup with defer resp.Body.Close()
- Good error messages in t.Fatalf() calls

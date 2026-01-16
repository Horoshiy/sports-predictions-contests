# Code Review: E2E Testing - Post-Fix Review

**Date**: 2026-01-16
**Reviewer**: Kiro CLI

## Stats

- Files Modified: 1 (Makefile)
- Files Added: 10 (tests/e2e/*, scripts/e2e-test.sh)
- Files Deleted: 0
- New lines: ~600
- Deleted lines: 0

---

## Issues Found

### MEDIUM

```
severity: medium
file: tests/e2e/workflow_test.go
line: 55
issue: Ignored error from parseResponse for sportResp
detail: After status code check, the error from parseResponse[SportResponse] is ignored 
        with `_`. If JSON parsing fails, sportID will be 0 and subsequent operations 
        will fail with confusing errors.
suggestion: Add error check: `sportResp, err := parseResponse[SportResponse](resp); if err != nil { t.Fatalf(...) }`
```

```
severity: medium
file: tests/e2e/workflow_test.go
line: 75
issue: Ignored error from parseResponse for leagueResp
detail: Same pattern - error ignored after status check.
suggestion: Add error handling for parseResponse.
```

```
severity: medium
file: tests/e2e/workflow_test.go
line: 97, 116
issue: Ignored error from parseResponse for team responses
detail: Both homeTeamResp and awayTeamResp ignore parseResponse errors.
suggestion: Add error handling for both parseResponse calls.
```

```
severity: medium
file: tests/e2e/workflow_test.go
line: 131
issue: Ignored error from parseResponse for matchResp
detail: Error ignored, could cause confusing failures.
suggestion: Add error handling.
```

```
severity: medium
file: tests/e2e/workflow_test.go
line: 152
issue: Ignored error from parseResponse for contestResp
detail: Error ignored after status check.
suggestion: Add error handling.
```

```
severity: medium
file: tests/e2e/workflow_test.go
line: 168
issue: Ignored error from parseResponse for eventResp
detail: Error ignored, eventID could be 0.
suggestion: Add error handling.
```

```
severity: medium
file: tests/e2e/workflow_test.go
line: 184
issue: Ignored error from parseResponse for predResp
detail: Error ignored after status check.
suggestion: Add error handling.
```

### LOW

```
severity: low
file: tests/e2e/helpers.go
line: 14
issue: Deprecated rand.Seed in Go 1.20+
detail: In Go 1.20+, math/rand is automatically seeded. rand.Seed() is deprecated 
        and will be removed in future versions. For Go 1.21, this works but generates 
        a deprecation warning.
suggestion: Remove rand.Seed() call or use crypto/rand for truly unique values.
```

```
severity: low
file: tests/e2e/auth_test.go
line: 14
issue: Variable userID declared but only used for logging
detail: userID is assigned but only used in t.Logf(). If logging is removed, 
        this becomes an unused variable error.
suggestion: Consider using `_ = result.User.ID` if only for verification, or 
            add assertions that use userID.
```

---

## Summary

The previous HIGH priority issues have been fixed. The code now properly handles errors in the setup blocks of test files. However, `workflow_test.go` still has multiple instances of ignored `parseResponse` errors after status code checks.

**Critical Issues**: 0
**High Issues**: 0 (all fixed from previous review)
**Medium Issues**: 7 (all in workflow_test.go - ignored parseResponse errors)
**Low Issues**: 2

### Positive Changes Since Last Review

✅ All test setup blocks now check parseResponse errors
✅ rand.Seed added for unique test emails
✅ Database health check added to e2e-test.sh
✅ Team creation status checks added in workflow_test.go
✅ Sequential test execution documented in auth_test.go

### Remaining Work

The `workflow_test.go` file has a pattern where status code is checked but parseResponse error is ignored. While the status check catches HTTP errors, JSON parsing errors would still cause silent failures with zero values.

### Recommendation

Either:
1. Fix all 7 ignored parseResponse errors in workflow_test.go, OR
2. Accept as-is since status code checks provide reasonable error detection

The code is functional and the remaining issues are medium priority - they won't cause incorrect behavior, just potentially confusing error messages if JSON parsing fails (which is unlikely if status code is 200).

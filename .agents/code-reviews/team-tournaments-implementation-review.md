# Code Review: Team Tournaments Feature

**Date**: 2026-01-16
**Reviewer**: Kiro CLI
**Feature**: Team Tournaments Implementation

## Stats

- Files Modified: 4
- Files Added: 19
- Files Deleted: 0
- New lines: ~1,800
- Deleted lines: 0

---

## Issues Found

### CRITICAL

```
severity: critical
file: backend/contest-service/internal/repository/team_repository.go
line: 107
issue: Wrong table name in subquery causes SQL error
detail: The List() method uses "team_members" in the subquery but the actual table is "user_team_members" (as defined by TableName() method). This will cause a SQL error when myTeamsOnly=true.
suggestion: Change `SELECT team_id FROM team_members` to `SELECT team_id FROM user_team_members`
```

---

### HIGH

```
severity: high
file: backend/contest-service/internal/service/team_service.go
line: 327
issue: updateMemberCount silently ignores Update error
detail: The updateMemberCount helper calls s.teamRepo.Update(team) but doesn't check or log the error. If the update fails, the CurrentMembers count will be stale.
suggestion: Add error handling: `if err := s.teamRepo.Update(team); err != nil { log.Printf("[ERROR] Failed to update member count: %v", err) }`
```

```
severity: high
file: backend/contest-service/internal/models/team.go
line: 60
issue: ValidateMaxMembers allows 0 which bypasses validation
detail: When MaxMembers is 0 (default uint value), the validation passes because 0 < 2 returns true. However, in BeforeCreate, if maxMembers is 0, it's not set to a default, so teams could be created with MaxMembers=0 which breaks CanJoin() logic.
suggestion: Add check in BeforeCreate: `if t.MaxMembers == 0 { t.MaxMembers = 10 }` before validation, or change validation to `if t.MaxMembers < 2 || t.MaxMembers == 0`
```

```
severity: high
file: frontend/src/components/teams/TeamList.tsx
line: 27
issue: useEffect creates infinite loop with searchParams dependency
detail: The useEffect updates searchParams, which triggers re-render, which triggers useEffect again. This is the same pattern that was fixed in other components.
suggestion: Remove `searchParams` from the dependency array, keep only `pagination`. Or use a ref to track if the update came from user action vs URL sync.
```

---

### MEDIUM

```
severity: medium
file: backend/contest-service/internal/repository/team_repository.go
line: 91
issue: Delete transaction doesn't use defer for rollback recovery
detail: The Delete method starts a transaction but the defer recovery only handles panics, not regular errors. If tx.Commit() fails, the transaction may be left in an inconsistent state.
suggestion: Add `defer tx.Rollback()` at the start (GORM ignores rollback on committed transactions)
```

```
severity: medium
file: backend/contest-service/internal/service/team_service.go
line: 193
issue: JoinTeam returns raw error to client exposing internal details
detail: When memberRepo.Create fails (e.g., duplicate membership), the raw database error is returned which may expose internal schema details.
suggestion: Wrap the error: `return nil, fmt.Errorf("failed to join team: %w", err)` or return a user-friendly message
```

```
severity: medium
file: frontend/src/hooks/use-teams.ts
line: 28
issue: useTeamMembers query key doesn't include pagination
detail: The query key only includes teamId, so changing pagination won't trigger a refetch. This causes stale data when paginating through members.
suggestion: Change to `queryKey: [teamKeys.members(request.teamId), request.pagination]`
```

```
severity: medium
file: frontend/src/components/teams/TeamMembers.tsx
line: 14
issue: Hardcoded pagination limit of 50 may cause performance issues
detail: Fetching 50 members at once could be slow for large teams. Also, there's no pagination UI to load more members.
suggestion: Reduce default limit to 20 and add pagination controls, or use infinite scroll
```

---

### LOW

```
severity: low
file: backend/contest-service/internal/service/team_service.go
line: 12
issue: Unused import warning potential
detail: The `emptypb` import is only used in the Check method. If proto generation changes, this could become unused.
suggestion: No action needed, just noting for awareness
```

```
severity: low
file: frontend/src/types/team.types.ts
line: 51-60
issue: Duplicate PaginationRequest/PaginationResponse types
detail: These types are already defined in common.types.ts. Having duplicates can lead to inconsistencies.
suggestion: Import from common.types.ts instead: `import type { PaginationRequest, PaginationResponse } from './common.types'`
```

```
severity: low
file: frontend/src/components/teams/TeamInvite.tsx
line: 20
issue: navigator.clipboard.writeText may fail without HTTPS
detail: The Clipboard API requires a secure context (HTTPS) in production. In development HTTP, this will fail silently.
suggestion: Add try-catch and fallback: `try { await navigator.clipboard.writeText(inviteCode) } catch { /* fallback or show error */ }`
```

```
severity: low
file: tests/contest-service/team_test.go
line: 1
issue: Test package naming inconsistency
detail: Package is `contest_service_test` but file is in `tests/contest-service/`. This may cause import issues.
suggestion: Ensure go.mod exists in tests/contest-service/ with correct module path
```

---

## Summary

The Team Tournaments feature implementation is generally well-structured and follows existing codebase patterns. However, there are several issues that should be addressed before merging:

**Must Fix (Critical/High):**
1. Wrong table name in repository subquery - will cause runtime SQL errors
2. Silent error in updateMemberCount - data consistency risk
3. MaxMembers=0 validation gap - logic error
4. useEffect infinite loop in TeamList - performance/UX issue

**Should Fix (Medium):**
1. Transaction handling in Delete
2. Error message exposure in JoinTeam
3. Query key missing pagination
4. Hardcoded member limit

**Nice to Have (Low):**
1. Duplicate type definitions
2. Clipboard API error handling
3. Test package naming

Total: 1 Critical, 3 High, 4 Medium, 4 Low issues identified.

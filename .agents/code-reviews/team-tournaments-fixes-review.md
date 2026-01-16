# Code Review: Team Tournaments Fixes

**Date**: 2026-01-16
**Reviewer**: Kiro CLI
**Scope**: Post-fix review of Team Tournaments implementation

## Stats

- Files Modified: 4
- Files Added: 17 (excluding review/plan files)
- Files Deleted: 0
- New lines: ~77 (in modified files)
- Total new code: ~1,800 lines

---

## Issues Found

### HIGH

```
severity: high
file: backend/contest-service/internal/service/team_service.go
line: 87-92
issue: Race condition in CreateTeam - captain member creation failure leaves orphan team
detail: If memberRepo.Create fails after teamRepo.Create succeeds, the code calls s.teamRepo.Delete(team.ID). However, Delete() starts a new transaction and deletes members first, which may fail or leave inconsistent state. The original team creation and member addition should be in a single transaction.
suggestion: Wrap team creation and captain member addition in a single database transaction, or use a saga pattern with proper compensation.
```

```
severity: high
file: backend/contest-service/internal/repository/team_repository.go
line: 107
issue: List query doesn't filter by deleted_at for soft-deleted teams
detail: The query filters by `is_active = true` but doesn't check `deleted_at IS NULL`. GORM's soft delete requires explicit handling when using raw WHERE clauses. Soft-deleted teams may appear in listings.
suggestion: Add `.Where("deleted_at IS NULL")` to the query, or use `r.db.Model(&models.Team{}).Unscoped().Where(...)` pattern consistently.
```

---

### MEDIUM

```
severity: medium
file: backend/contest-service/internal/service/team_service.go
line: 195-197
issue: JoinTeam doesn't verify user isn't already a member before attempting create
detail: The duplicate check happens in BeforeCreate hook, but this means a database query is made even when we could check earlier. Also, the error message from BeforeCreate ("user is already a member") is wrapped, making it less clear to the user.
suggestion: Add explicit check: `if _, err := s.memberRepo.GetByTeamAndUser(team.ID, userID); err == nil { return nil, &joinError{"Already a member"} }`
```

```
severity: medium
file: backend/contest-service/internal/models/team_member.go
line: 55-58
issue: BeforeCreate performs database query which may cause issues in batch operations
detail: The duplicate check query in BeforeCreate hook runs for every insert, which could cause N+1 queries in batch scenarios and may behave unexpectedly with transactions.
suggestion: Move duplicate checking to service layer where it can be controlled, or document this behavior clearly.
```

```
severity: medium
file: frontend/src/components/teams/TeamList.tsx
line: 23-24
issue: parseInt without radix and no NaN handling
detail: `parseInt(searchParams.get('page') || '0')` doesn't specify radix and doesn't handle NaN case. If URL contains `?page=abc`, this will result in NaN.
suggestion: Use `parseInt(searchParams.get('page') || '0', 10) || 0` to handle both radix and NaN.
```

```
severity: medium
file: frontend/src/pages/TeamsPage.tsx
line: 37-40
issue: Unused leaveTeamMutation variable in component scope
detail: leaveTeamMutation is declared but only used inside handleLeaveTeam. The mutation state (isPending) is used in the button, but if the dialog closes, the pending state is lost.
suggestion: Consider keeping dialog open until mutation completes, or move mutation closer to usage.
```

---

### LOW

```
severity: low
file: backend/contest-service/internal/service/team_service.go
line: 14
issue: Unused import - emptypb only used in Check method
detail: The emptypb import is only used for the health check method. If proto generation changes or Check is removed, this becomes unused.
suggestion: No action needed, just noting for awareness during future refactoring.
```

```
severity: low
file: frontend/src/types/team.types.ts
line: 60
issue: Re-exporting types that are already exported from common.types
detail: `export type { PaginationRequest, PaginationResponse }` re-exports types that consumers could import directly from common.types.ts. This creates two import paths for the same types.
suggestion: Remove re-export or document that team.types.ts is the canonical import location for team-related code.
```

```
severity: low
file: scripts/init-db.sql
line: 218
issue: user_team_members foreign key uses ON DELETE RESTRICT
detail: ON DELETE RESTRICT prevents deleting teams that have members. While this is intentional for data integrity, it means the application must delete all members before deleting a team. The repository Delete() method handles this, but direct SQL operations would fail.
suggestion: Document this constraint, or consider ON DELETE CASCADE if orphan prevention is handled at application level.
```

```
severity: low
file: backend/contest-service/internal/models/team.go
line: 53-56
issue: ValidateMaxMembers has redundant condition
detail: `if t.MaxMembers == 0 || t.MaxMembers < 2` - the `< 2` check already covers `== 0` since 0 < 2 is true.
suggestion: Simplify to `if t.MaxMembers < 2` for clarity, or keep explicit 0 check with comment explaining it's for documentation purposes.
```

---

## Previously Fixed Issues (Verified)

The following issues from the previous review have been correctly fixed:

1. ✅ **CRITICAL**: Wrong table name `team_members` → `user_team_members` in repository subquery
2. ✅ **HIGH**: Silent error in updateMemberCount - now logs errors
3. ✅ **HIGH**: MaxMembers=0 validation gap - now sets default before validation
4. ✅ **HIGH**: useEffect infinite loop - removed searchParams from dependency array
5. ✅ **MEDIUM**: Transaction rollback - now uses defer pattern
6. ✅ **MEDIUM**: Raw error exposure in JoinTeam - now wrapped with fmt.Errorf
7. ✅ **MEDIUM**: Query key missing pagination - now includes request.pagination
8. ✅ **MEDIUM**: Hardcoded limit 50 → 20
9. ✅ **LOW**: Duplicate types - now imports from common.types.ts
10. ✅ **LOW**: Clipboard API error handling - now has try-catch

---

## Summary

The Team Tournaments feature implementation is solid after the previous round of fixes. The remaining issues are primarily:

**Must Address (High):**
1. Race condition in team creation - needs transactional handling
2. Soft delete filtering in List query

**Should Address (Medium):**
1. Duplicate member check optimization
2. BeforeCreate hook database query concern
3. parseInt NaN handling in TeamList
4. Mutation state handling in TeamsPage

**Nice to Have (Low):**
1. Unused import awareness
2. Type re-export cleanup
3. Foreign key constraint documentation
4. Redundant validation condition

Total: 2 High, 4 Medium, 4 Low issues identified.

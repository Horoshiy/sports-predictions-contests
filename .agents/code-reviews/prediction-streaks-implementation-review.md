# Code Review: Prediction Streaks with Multipliers

**Date**: 2026-01-16
**Feature**: Prediction Streaks with Multipliers
**Reviewer**: Kiro CLI

## Stats

- Files Modified: 10
- Files Added: 6
- Files Deleted: 0
- New lines: ~350
- Deleted lines: ~16

---

## Issues Found

### CRITICAL

```
severity: critical
file: backend/scoring-service/internal/service/leaderboard_service.go
line: 60-65
issue: N+1 query problem in GetLeaderboard
detail: For each leaderboard entry, a separate database query is made to fetch streak data. With 50 entries (default limit), this results in 51 database queries instead of 2. This will cause significant performance degradation under load.
suggestion: Batch fetch all streaks for the contest in a single query before the loop, then look up from a map:
```go
// Fetch all streaks for users in one query
userIDs := make([]uint, len(leaderboards))
for i, lb := range leaderboards {
    userIDs[i] = lb.UserID
}
streaks, _ := s.streakRepo.GetByContestAndUsers(ctx, req.ContestId, userIDs)
streakMap := make(map[uint]*models.UserStreak)
for _, streak := range streaks {
    streakMap[streak.UserID] = streak
}
// Then use streakMap[lb.UserID] in the loop
```
```

---

### HIGH

```
severity: high
file: backend/scoring-service/internal/repository/streak_repository.go
line: 30-47
issue: Race condition in GetOrCreate
detail: Between checking if streak exists and creating a new one, another concurrent request could create the same streak, causing a unique constraint violation. This is especially likely during high-traffic periods when multiple predictions are submitted simultaneously.
suggestion: Use GORM's FirstOrCreate with proper locking, or handle the unique constraint error gracefully:
```go
func (r *StreakRepository) GetOrCreate(ctx context.Context, contestID, userID uint) (*models.UserStreak, error) {
    streak := models.UserStreak{
        UserID:    userID,
        ContestID: contestID,
    }
    err := r.db.WithContext(ctx).Where("contest_id = ? AND user_id = ?", contestID, userID).
        FirstOrCreate(&streak).Error
    if err != nil {
        return nil, err
    }
    return &streak, nil
}
```
```

```
severity: high
file: backend/scoring-service/internal/service/scoring_service.go
line: 97-99
issue: Silent failure on streak update doesn't rollback score creation
detail: If streak update fails, the error is only logged but the score is still created with the multiplied points. This creates data inconsistency - user gets bonus points but streak isn't recorded. Next prediction will start from streak 0 again.
suggestion: Either wrap in a transaction or return error to prevent score creation:
```go
if err := s.streakRepo.Update(ctx, streak); err != nil {
    log.Printf("[ERROR] Failed to update streak: %v", err)
    return &pb.CreateScoreResponse{
        Response: &common.Response{
            Success:   false,
            Message:   "Failed to update streak",
            Code:      int32(common.ErrorCode_INTERNAL_ERROR),
            Timestamp: timestamppb.Now(),
        },
    }, nil
}
```
```

```
severity: high
file: backend/scoring-service/internal/models/streak.go
line: 51-54
issue: BeforeUpdate hook sets UpdatedAt but GORM already handles this
detail: The UserStreak struct embeds gorm.Model which already has UpdatedAt field with auto-update. Manually setting UpdatedAt in BeforeUpdate may conflict with GORM's automatic handling and could cause issues with soft deletes.
suggestion: Remove the manual UpdatedAt assignment - GORM handles this automatically:
```go
func (s *UserStreak) BeforeUpdate(tx *gorm.DB) error {
    // GORM automatically updates UpdatedAt via gorm.Model
    return nil
}
```
Or remove the hook entirely if no other validation is needed.
```

---

### MEDIUM

```
severity: medium
file: backend/scoring-service/internal/service/scoring_service.go
line: 82-88
issue: Multiplier applied before streak is incremented for first correct prediction
detail: When a user makes their first correct prediction (streak goes 0→1), the multiplier is calculated AFTER IncrementStreak(), so they get 1.0x. But the log message shows streak=1 with multiplier=1.0x which is correct. However, the code flow is confusing - multiplier should be calculated after streak update for clarity.
suggestion: The current logic is actually correct (multiplier based on NEW streak value), but add a comment to clarify:
```go
// Update streak first, then calculate multiplier based on new streak value
if isCorrect {
    streak.IncrementStreak(uint(req.PredictionId))
} else {
    streak.ResetStreak(uint(req.PredictionId))
}
// Multiplier is based on the updated streak value
multiplier := streak.GetMultiplier()
```
```

```
severity: medium
file: frontend/src/components/leaderboard/LeaderboardTable.tsx
line: 104-112
issue: userStreak query doesn't handle error state
detail: The userStreak query destructures only `data`, ignoring potential errors. If the API call fails, the UI will silently show no streak data without any indication to the user.
suggestion: Add error handling:
```typescript
const {
  data: userStreak,
  error: userStreakError,
} = useQuery({
  // ...
})
// Then optionally show error state or log it
```
```

```
severity: medium
file: backend/scoring-service/internal/models/streak.go
line: 11-19
issue: Duplicate ID field - gorm.Model already includes ID
detail: UserStreak defines its own ID field with `gorm:"primaryKey"` but also embeds gorm.Model which already has an ID field. This creates ambiguity and potential issues.
suggestion: Remove the explicit ID field since gorm.Model provides it:
```go
type UserStreak struct {
    gorm.Model
    UserID                uint  `gorm:"not null;uniqueIndex:idx_user_contest_streak" json:"user_id"`
    ContestID             uint  `gorm:"not null;uniqueIndex:idx_user_contest_streak" json:"contest_id"`
    // ...
}
```
Or keep explicit ID and remove gorm.Model, using only CreatedAt/UpdatedAt/DeletedAt explicitly.
```

```
severity: medium
file: scripts/init-db.sql
line: 57-71
issue: user_streaks table missing foreign key constraints
detail: The user_streaks table references user_id and contest_id but doesn't have foreign key constraints. This allows orphaned records if users or contests are deleted.
suggestion: Add foreign key constraints (if users/contests tables exist):
```sql
CREATE TABLE IF NOT EXISTS user_streaks (
    -- ... existing columns ...
    CONSTRAINT fk_user_streaks_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_streaks_contest FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE
);
```
Note: This may require contests table to exist first.
```

---

### LOW

```
severity: low
file: frontend/src/components/leaderboard/LeaderboardTable.tsx
line: 22-23
issue: Unused imports TrendingUpIcon and TrendingDownIcon
detail: These icons are imported but never used in the component.
suggestion: Remove unused imports:
```typescript
import {
  Refresh as RefreshIcon,
  EmojiEvents as TrophyIcon,
} from '@mui/icons-material'
```
```

```
severity: low
file: frontend/src/components/leaderboard/LeaderboardTable.tsx
line: 28
issue: Unused import 'Leaderboard' type
detail: The Leaderboard type is imported but the component uses the response directly.
suggestion: Remove if not needed, or keep if used elsewhere in the file.
```

```
severity: low
file: backend/scoring-service/internal/repository/streak_repository.go
line: 68-74
issue: GetTopStreaks doesn't validate limit parameter
detail: If limit is 0 or negative, the query may return unexpected results.
suggestion: Add validation:
```go
func (r *StreakRepository) GetTopStreaks(ctx context.Context, contestID uint, limit int) ([]*models.UserStreak, error) {
    if limit <= 0 {
        limit = 10 // default
    }
    // ...
}
```
```

---

## Security Analysis

✅ No SQL injection vulnerabilities - using GORM parameterized queries
✅ No XSS vulnerabilities - React handles escaping
✅ No exposed secrets or API keys
✅ Proper authentication check in CreateScore

---

## Performance Analysis

⚠️ **N+1 Query Issue** (Critical): GetLeaderboard makes N+1 queries for streak data
✅ Database indexes properly defined for user_streaks table
✅ Redis caching preserved for leaderboard data

---

## Code Quality Summary

**Strengths:**
- Clean separation of concerns (model, repository, service)
- Consistent error handling patterns
- Good test coverage for streak model
- Proper use of GORM hooks for validation

**Areas for Improvement:**
- N+1 query in GetLeaderboard needs batch optimization
- Race condition in GetOrCreate needs transaction handling
- Silent failure on streak update should be addressed

---

## Recommendation

**Do not merge** until Critical and High issues are addressed:

1. Fix N+1 query in GetLeaderboard (Critical)
2. Fix race condition in GetOrCreate (High)
3. Handle streak update failure properly (High)
4. Remove duplicate ID field in UserStreak model (Medium)

Medium and Low issues can be addressed in follow-up commits.

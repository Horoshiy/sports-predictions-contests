# Feature: Prediction Streaks with Multipliers

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

Implement a gamification system where consecutive successful predictions increase a point multiplier. A series of correct predictions builds a "streak" that multiplies earned points, but resets to 1x on any failed prediction. This creates excitement, encourages regular platform visits, and rewards consistent accuracy.

**Multiplier Formula**:
- Streak 0-2: 1.0x (base)
- Streak 3-4: 1.25x
- Streak 5-6: 1.5x
- Streak 7-9: 1.75x
- Streak 10+: 2.0x (max)

## User Story

As a prediction contest participant
I want to earn bonus points for consecutive correct predictions
So that I'm rewarded for consistent accuracy and motivated to return regularly

## Problem Statement

Current scoring system treats each prediction independently without rewarding consistency. Users who make many correct predictions in a row receive the same per-prediction points as users who alternate between correct and incorrect. This misses an opportunity for gamification and engagement.

## Solution Statement

Add streak tracking to the scoring system that:
1. Tracks current and maximum streak per user per contest
2. Applies a multiplier to points based on current streak length
3. Resets streak on failed predictions
4. Displays streak information in leaderboard and user profile
5. Stores streak history for analytics

## Feature Metadata

**Feature Type**: Enhancement
**Estimated Complexity**: Low
**Primary Systems Affected**: scoring-service, frontend (leaderboard)
**Dependencies**: None (uses existing infrastructure)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/scoring-service/internal/models/score.go` - Score model structure to extend
- `backend/scoring-service/internal/models/leaderboard.go` - Leaderboard model with validation hooks
- `backend/scoring-service/internal/service/scoring_service.go` - Scoring logic with `calculatePoints()` method
- `backend/scoring-service/internal/repository/score_repository.go` - Score repository interface
- `backend/scoring-service/internal/repository/leaderboard_repository.go` - Leaderboard repository with cache
- `backend/proto/scoring.proto` - gRPC proto definitions
- `frontend/src/types/scoring.types.ts` - Frontend TypeScript types
- `frontend/src/components/leaderboard/LeaderboardTable.tsx` - Leaderboard UI component
- `scripts/init-db.sql` (lines 21-55) - Database schema for scores and leaderboards

### New Files to Create

- `backend/scoring-service/internal/models/streak.go` - UserStreak model
- `backend/scoring-service/internal/repository/streak_repository.go` - Streak repository
- `tests/scoring-service/streak_test.go` - Unit tests for streak logic

### Files to Modify

- `backend/scoring-service/internal/service/scoring_service.go` - Add streak logic to CreateScore
- `backend/proto/scoring.proto` - Add streak fields to LeaderboardEntry
- `backend/scoring-service/internal/models/leaderboard.go` - Add streak fields
- `scripts/init-db.sql` - Add user_streaks table
- `frontend/src/types/scoring.types.ts` - Add streak types
- `frontend/src/components/leaderboard/LeaderboardTable.tsx` - Display streak info

### Patterns to Follow

**Model Validation (from score.go)**:
```go
func (s *Score) ValidateUserID() error {
    if s.UserID == 0 {
        return errors.New("user ID cannot be empty")
    }
    return nil
}

func (s *Score) BeforeCreate(tx *gorm.DB) error {
    // Validate fields before insert
}
```

**Repository Interface Pattern (from score_repository.go)**:
```go
type ScoreRepositoryInterface interface {
    Create(ctx context.Context, score *models.Score) error
    GetByID(ctx context.Context, id uint) (*models.Score, error)
    // ...
}
```

**Service Response Pattern (from scoring_service.go)**:
```go
return &pb.CreateScoreResponse{
    Response: &common.Response{
        Success:   true,
        Message:   "Score created successfully",
        Code:      int32(common.ErrorCode_SUCCESS),
        Timestamp: timestamppb.Now(),
    },
    Score: s.modelToProto(score),
}, nil
```

**Naming Conventions**:
- Go files: `snake_case.go`
- Go structs: `PascalCase`
- Go functions: `PascalCase` (public), `camelCase` (private)
- Proto messages: `PascalCase`
- TypeScript interfaces: `PascalCase`

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation (Database & Models)

Create the streak tracking infrastructure:
- Database table for user streaks
- Go model with validation
- Repository interface and implementation

### Phase 2: Core Implementation (Scoring Logic)

Integrate streak tracking into scoring flow:
- Update streak on score creation
- Calculate multiplier based on streak
- Apply multiplier to points

### Phase 3: API & Proto Updates

Expose streak data through API:
- Add streak fields to proto messages
- Update leaderboard response

### Phase 4: Frontend Integration

Display streak information:
- Update TypeScript types
- Add streak column to leaderboard
- Show streak badge/indicator

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `scripts/init-db.sql` - Add user_streaks table

- **IMPLEMENT**: Add new table after leaderboards table definition
- **PATTERN**: Follow existing table structure (scores, leaderboards)
- **SQL**:
```sql
-- Create user_streaks table for streak tracking
CREATE TABLE IF NOT EXISTS user_streaks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    contest_id INTEGER NOT NULL,
    current_streak INTEGER NOT NULL DEFAULT 0,
    max_streak INTEGER NOT NULL DEFAULT 0,
    last_prediction_id INTEGER,
    last_prediction_correct BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(user_id, contest_id)
);

CREATE INDEX IF NOT EXISTS idx_user_streaks_contest_user ON user_streaks(contest_id, user_id);
CREATE INDEX IF NOT EXISTS idx_user_streaks_deleted_at ON user_streaks(deleted_at);
```
- **VALIDATE**: `grep -A 20 "user_streaks" scripts/init-db.sql`

### Task 2: CREATE `backend/scoring-service/internal/models/streak.go`

- **IMPLEMENT**: UserStreak model with validation hooks
- **PATTERN**: Mirror `score.go` and `leaderboard.go` structure
- **IMPORTS**: `errors`, `time`, `gorm.io/gorm`
- **FIELDS**: ID, UserID, ContestID, CurrentStreak, MaxStreak, LastPredictionID, LastPredictionCorrect, gorm.Model
- **METHODS**: 
  - `ValidateUserID()`, `ValidateContestID()` 
  - `BeforeCreate()`, `BeforeUpdate()` hooks
  - `GetMultiplier() float64` - returns multiplier based on CurrentStreak
  - `IncrementStreak()` - increases streak and updates max
  - `ResetStreak()` - sets current streak to 0
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 3: CREATE `backend/scoring-service/internal/repository/streak_repository.go`

- **IMPLEMENT**: Repository interface and implementation
- **PATTERN**: Mirror `score_repository.go` structure
- **IMPORTS**: `context`, `errors`, `gorm.io/gorm`, models package
- **INTERFACE**: `StreakRepositoryInterface`
  - `GetOrCreate(ctx, contestID, userID uint) (*models.UserStreak, error)`
  - `Update(ctx, streak *models.UserStreak) error`
  - `GetByContestAndUser(ctx, contestID, userID uint) (*models.UserStreak, error)`
  - `GetTopStreaks(ctx, contestID uint, limit int) ([]*models.UserStreak, error)`
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 4: UPDATE `backend/scoring-service/internal/service/scoring_service.go`

- **IMPLEMENT**: Add streak repository and integrate into CreateScore
- **CHANGES**:
  1. Add `streakRepo` field to `ScoringService` struct
  2. Update `NewScoringService` constructor to accept streak repository
  3. In `CreateScore`:
     - Get or create user streak
     - Determine if prediction was correct (points > 0)
     - Update streak (increment or reset)
     - Apply multiplier to points before saving
     - Include multiplier in response
- **PATTERN**: Follow existing repository injection pattern
- **GOTCHA**: Apply multiplier BEFORE saving score, not after
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 5: UPDATE `backend/scoring-service/cmd/main.go`

- **IMPLEMENT**: Initialize streak repository and pass to service
- **PATTERN**: Follow existing repository initialization
- **CHANGES**:
  1. Auto-migrate UserStreak model
  2. Create StreakRepository instance
  3. Pass to NewScoringService
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 6: UPDATE `backend/proto/scoring.proto`

- **IMPLEMENT**: Add streak fields to LeaderboardEntry and new messages
- **CHANGES**:
  1. Add to `LeaderboardEntry`:
     ```protobuf
     uint32 current_streak = 6;
     uint32 max_streak = 7;
     double multiplier = 8;
     ```
  2. Add new message:
     ```protobuf
     message GetUserStreakRequest {
       uint32 contest_id = 1;
       uint32 user_id = 2;
     }
     
     message GetUserStreakResponse {
       common.Response response = 1;
       uint32 current_streak = 2;
       uint32 max_streak = 3;
       double multiplier = 4;
     }
     ```
  3. Add RPC:
     ```protobuf
     rpc GetUserStreak(GetUserStreakRequest) returns (GetUserStreakResponse) {
       option (google.api.http) = {
         get: "/v1/contests/{contest_id}/users/{user_id}/streak"
       };
     }
     ```
- **VALIDATE**: `cd backend && buf lint proto/`

### Task 7: UPDATE `backend/scoring-service/internal/models/leaderboard.go`

- **IMPLEMENT**: Add streak fields to Leaderboard model
- **CHANGES**: Add fields (not persisted, populated from join):
  ```go
  CurrentStreak uint    `gorm:"-" json:"current_streak"`
  MaxStreak     uint    `gorm:"-" json:"max_streak"`
  Multiplier    float64 `gorm:"-" json:"multiplier"`
  ```
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 8: UPDATE `backend/scoring-service/internal/service/leaderboard_service.go`

- **IMPLEMENT**: Populate streak data in leaderboard responses
- **CHANGES**:
  1. Add streakRepo to LeaderboardService
  2. In GetLeaderboard, fetch streak data for each user
  3. Implement GetUserStreak RPC method
- **PATTERN**: Follow existing service method patterns
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

### Task 9: UPDATE `frontend/src/types/scoring.types.ts`

- **IMPLEMENT**: Add streak types
- **CHANGES**:
  1. Add to `LeaderboardEntry`:
     ```typescript
     currentStreak: number
     maxStreak: number
     multiplier: number
     ```
  2. Add new interfaces:
     ```typescript
     export interface GetUserStreakRequest {
       contestId: number
       userId: number
     }
     
     export interface GetUserStreakResponse {
       response: ApiResponse
       currentStreak: number
       maxStreak: number
       multiplier: number
     }
     ```
- **VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 10: UPDATE `frontend/src/services/scoring-service.ts`

- **IMPLEMENT**: Add getUserStreak method
- **PATTERN**: Follow existing service method patterns
- **VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 11: UPDATE `frontend/src/components/leaderboard/LeaderboardTable.tsx`

- **IMPLEMENT**: Display streak information in leaderboard
- **CHANGES**:
  1. Add streak column after points:
     ```typescript
     {
       accessorKey: 'currentStreak',
       header: 'Streak',
       size: 100,
       Cell: ({ cell, row }) => {
         const streak = cell.getValue<number>()
         const multiplier = row.original.multiplier
         return (
           <Box display="flex" alignItems="center" gap={1}>
             <Typography variant="body2">
               🔥 {streak}
             </Typography>
             {multiplier > 1 && (
               <Chip 
                 label={`${multiplier}x`} 
                 size="small" 
                 color="warning"
               />
             )}
           </Box>
         )
       },
     }
     ```
  2. Update user rank card to show streak
- **VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 12: CREATE `tests/scoring-service/streak_test.go`

- **IMPLEMENT**: Unit tests for streak model and multiplier calculation
- **PATTERN**: Follow `scoring_test.go` structure
- **TEST CASES**:
  1. `TestStreakValidation` - validate UserID, ContestID
  2. `TestGetMultiplier` - verify multiplier tiers (0-2: 1.0, 3-4: 1.25, etc.)
  3. `TestIncrementStreak` - verify streak increment and max update
  4. `TestResetStreak` - verify streak reset to 0
- **VALIDATE**: `cd tests/scoring-service && go test -v ./...`

---

## TESTING STRATEGY

### Unit Tests

**Streak Model Tests** (`tests/scoring-service/streak_test.go`):
- Validation methods (UserID, ContestID)
- GetMultiplier returns correct values for each tier
- IncrementStreak updates both current and max
- ResetStreak sets current to 0, preserves max

**Scoring Service Tests**:
- CreateScore applies multiplier correctly
- Streak increments on successful prediction (points > 0)
- Streak resets on failed prediction (points = 0)

### Integration Tests

- Create multiple scores for same user/contest
- Verify streak increments correctly
- Verify leaderboard includes streak data

### Edge Cases

1. First prediction (no existing streak) - should create streak with 0
2. Exact 0 points - should reset streak
3. Negative points (if possible) - should reset streak
4. Concurrent score creation - streak should be consistent
5. Max streak preservation - max should never decrease

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
# Backend
cd backend/scoring-service && go build ./...
cd backend && buf lint proto/

# Frontend
cd frontend && npx tsc --noEmit
```

### Level 2: Unit Tests

```bash
cd tests/scoring-service && go test -v ./...
```

### Level 3: Integration Tests

```bash
# Start services
docker-compose up -d postgres redis

# Run integration tests
cd tests/scoring-service && go test -v -tags=integration ./...
```

### Level 4: Manual Validation

1. Create a score with points > 0 → streak should be 1
2. Create another score with points > 0 → streak should be 2
3. Create score with points = 0 → streak should reset to 0
4. Verify leaderboard shows streak column
5. Verify multiplier applies to points (streak 3+ should show 1.25x)

---

## ACCEPTANCE CRITERIA

- [ ] UserStreak model created with validation
- [ ] Streak repository with GetOrCreate, Update methods
- [ ] CreateScore updates streak and applies multiplier
- [ ] Proto updated with streak fields in LeaderboardEntry
- [ ] GetUserStreak RPC endpoint implemented
- [ ] Frontend displays streak in leaderboard
- [ ] Multiplier badge shown when > 1.0x
- [ ] Unit tests pass for streak logic
- [ ] Multiplier tiers work correctly (1.0, 1.25, 1.5, 1.75, 2.0)
- [ ] Max streak preserved when current resets

---

## COMPLETION CHECKLIST

- [ ] All 12 tasks completed in order
- [ ] Backend builds without errors
- [ ] Frontend compiles without TypeScript errors
- [ ] Unit tests pass
- [ ] Manual testing confirms streak tracking works
- [ ] Leaderboard displays streak information
- [ ] Multiplier applies correctly to points

---

## NOTES

### Multiplier Tiers Rationale

| Streak | Multiplier | Rationale |
|--------|------------|-----------|
| 0-2    | 1.0x       | Base - no bonus for short streaks |
| 3-4    | 1.25x      | Small reward for consistency |
| 5-6    | 1.5x       | Moderate reward |
| 7-9    | 1.75x      | Significant reward |
| 10+    | 2.0x       | Maximum - prevents runaway scores |

### Design Decisions

1. **Streak per contest**: Each contest has independent streaks to prevent cross-contest gaming
2. **Points > 0 = success**: Simple rule - any positive points counts as correct
3. **Multiplier applied before save**: Ensures stored points reflect actual earned value
4. **Max streak preserved**: Allows users to see their best performance even after reset
5. **No streak decay**: Streaks only reset on failure, not on inactivity (simpler UX)

### Future Enhancements (Out of Scope)

- Streak milestones with achievements/badges
- Streak insurance (one "free" miss)
- Streak leaderboard (separate from points)
- Streak notifications via notification-service

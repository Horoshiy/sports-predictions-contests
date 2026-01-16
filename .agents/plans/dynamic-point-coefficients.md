# Feature: Dynamic Point Coefficients

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

Points for predictions change based on submission time — earlier predictions earn more points. This creates a time-decay coefficient that rewards users who submit predictions early (when less information is available) with higher point multipliers, while later predictions closer to event start receive lower multipliers.

## User Story

As a **sports prediction participant**
I want to **earn more points for submitting predictions early**
So that **I'm rewarded for taking risks with less information and encouraged to engage with the platform regularly**

## Problem Statement

Currently, all predictions are scored equally regardless of when they're submitted. This doesn't incentivize early engagement or reward users who make predictions when outcomes are more uncertain. Users can wait until the last moment to gather maximum information before predicting.

## Solution Statement

Implement a time-based coefficient system that:
1. Calculates a multiplier based on how early a prediction is submitted relative to event start
2. Applies this coefficient during score calculation (alongside existing streak multiplier)
3. Displays potential points in the prediction form so users understand the incentive
4. Shows the coefficient applied in score breakdowns

## Feature Metadata

**Feature Type**: Enhancement
**Estimated Complexity**: Low (2-4 hours)
**Primary Systems Affected**: scoring-service, prediction-service, frontend
**Dependencies**: None (builds on existing scoring infrastructure)

---

## CONTEXT REFERENCES

### Relevant Codebase Files - MUST READ BEFORE IMPLEMENTING

- `backend/scoring-service/internal/service/scoring_service.go` (lines 63-130) - Why: CreateScore method where coefficient will be applied alongside streak multiplier
- `backend/scoring-service/internal/models/streak.go` (lines 55-68) - Why: GetMultiplier pattern to mirror for time coefficient
- `backend/prediction-service/internal/models/prediction.go` (lines 11-24) - Why: Prediction model with SubmittedAt field already exists
- `backend/prediction-service/internal/models/event.go` (lines 11-22) - Why: Event model with EventDate field for time calculation
- `backend/prediction-service/internal/service/prediction_service.go` (lines 43-115) - Why: SubmitPrediction method where coefficient can be calculated
- `frontend/src/components/predictions/PredictionForm.tsx` - Why: Form where potential points will be displayed
- `frontend/src/components/predictions/EventCard.tsx` - Why: Card where coefficient indicator can be shown
- `backend/proto/scoring.proto` (lines 11-21) - Why: Score message structure for adding coefficient field
- `backend/proto/prediction.proto` (lines 11-21) - Why: Prediction message for potential coefficient response

### New Files to Create

- `backend/scoring-service/internal/models/coefficient.go` - Time coefficient calculation logic
- `frontend/src/components/predictions/CoefficientIndicator.tsx` - Visual coefficient display component

### Files to Modify

- `backend/scoring-service/internal/service/scoring_service.go` - Apply time coefficient in CreateScore
- `backend/scoring-service/internal/models/score.go` - Add TimeCoefficient field
- `backend/proto/scoring.proto` - Add time_coefficient to Score message
- `backend/proto/prediction.proto` - Add GetPotentialCoefficient RPC
- `backend/prediction-service/internal/service/prediction_service.go` - Add coefficient calculation endpoint
- `frontend/src/types/scoring.types.ts` - Add timeCoefficient to Score type
- `frontend/src/types/prediction.types.ts` - Add coefficient types
- `frontend/src/components/predictions/PredictionForm.tsx` - Display potential points
- `frontend/src/components/predictions/EventCard.tsx` - Show coefficient indicator
- `scripts/init-db.sql` - Add time_coefficient column to scores table

### Patterns to Follow

**Multiplier Calculation Pattern** (from streak.go):
```go
// GetMultiplier returns the point multiplier based on current streak
func (s *UserStreak) GetMultiplier() float64 {
    switch {
    case s.CurrentStreak >= 10:
        return 2.0
    case s.CurrentStreak >= 7:
        return 1.75
    // ...
    }
}
```

**Score Creation Pattern** (from scoring_service.go):
```go
// Multiplier is based on the updated streak value
multiplier := streak.GetMultiplier()
finalPoints := basePoints * multiplier
```

**Error Handling Pattern**:
```go
if err != nil {
    log.Printf("[ERROR] Failed to X: %v", err)
    return &pb.Response{
        Success: false,
        Message: "Error message",
        Code:    int32(common.ErrorCode_INTERNAL_ERROR),
    }, nil
}
```

**Frontend Component Pattern** (from EventCard.tsx):
```tsx
export const ComponentName: React.FC<Props> = ({ prop1, prop2 }) => {
  return (
    <Box sx={{ ... }}>
      {/* content */}
    </Box>
  )
}
export default ComponentName
```

---

## IMPLEMENTATION PLAN

### Phase 1: Backend Coefficient Model

Create the time coefficient calculation logic as a standalone model following the streak multiplier pattern.

**Tasks:**
- Create coefficient model with calculation function
- Add time_coefficient field to Score model
- Update database schema

### Phase 2: Scoring Service Integration

Integrate coefficient calculation into the scoring flow alongside existing streak multiplier.

**Tasks:**
- Modify CreateScore to calculate and apply time coefficient
- Update proto definitions for coefficient field
- Add logging for coefficient application

### Phase 3: Prediction Service Enhancement

Add endpoint to calculate potential coefficient for UI display.

**Tasks:**
- Add GetPotentialCoefficient RPC to prediction service
- Calculate coefficient based on current time vs event date

### Phase 4: Frontend Integration

Display coefficient information in prediction UI.

**Tasks:**
- Create CoefficientIndicator component
- Update PredictionForm to show potential points
- Update EventCard to show coefficient status
- Add coefficient to scoring types

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `backend/scoring-service/internal/models/coefficient.go`

- **IMPLEMENT**: Time coefficient calculation model
- **PATTERN**: Mirror `streak.go` GetMultiplier pattern
- **FORMULA**: 
  - 7+ days before event: 2.0x
  - 3-7 days before: 1.5x
  - 1-3 days before: 1.25x
  - 12-24 hours before: 1.1x
  - <12 hours before: 1.0x
- **IMPORTS**: `time`
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

```go
package models

import "time"

// CalculateTimeCoefficient returns point multiplier based on prediction timing
// Earlier predictions relative to event start earn higher multipliers
func CalculateTimeCoefficient(submittedAt, eventDate time.Time) float64 {
    hoursUntilEvent := eventDate.Sub(submittedAt).Hours()
    
    switch {
    case hoursUntilEvent >= 168: // 7+ days
        return 2.0
    case hoursUntilEvent >= 72: // 3-7 days
        return 1.5
    case hoursUntilEvent >= 24: // 1-3 days
        return 1.25
    case hoursUntilEvent >= 12: // 12-24 hours
        return 1.1
    default:
        return 1.0
    }
}

// GetCoefficientTier returns human-readable tier name
func GetCoefficientTier(coefficient float64) string {
    switch coefficient {
    case 2.0:
        return "Early Bird"
    case 1.5:
        return "Ahead of Time"
    case 1.25:
        return "Timely"
    case 1.1:
        return "Last Minute"
    default:
        return "Standard"
    }
}
```

### Task 2: UPDATE `backend/scoring-service/internal/models/score.go`

- **IMPLEMENT**: Add TimeCoefficient field to Score struct
- **PATTERN**: Follow existing field patterns with GORM tags
- **IMPORTS**: None additional needed
- **GOTCHA**: Field must have default value 1.0 for backward compatibility
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

Add after line 17 (after Points field):
```go
TimeCoefficient float64 `gorm:"not null;default:1.0" json:"time_coefficient"`
```

### Task 3: UPDATE `backend/proto/scoring.proto`

- **IMPLEMENT**: Add time_coefficient field to Score message
- **PATTERN**: Follow existing field numbering
- **GOTCHA**: Use field number 9 (after updated_at which is 8)
- **VALIDATE**: `cd backend && protoc --proto_path=proto --go_out=shared --go-grpc_out=shared proto/scoring.proto`

Add to Score message after line 18:
```protobuf
double time_coefficient = 9;
```

### Task 4: UPDATE `backend/scoring-service/internal/service/scoring_service.go`

- **IMPLEMENT**: Calculate and apply time coefficient in CreateScore
- **PATTERN**: Mirror streak multiplier application pattern
- **IMPORTS**: Add prediction-service client or pass event_date in request
- **GOTCHA**: Need event date - add to CreateScoreRequest or fetch from prediction
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

Modify CreateScore method (around line 90-100):

1. First, update CreateScoreRequest in proto to include event_date:
```protobuf
message CreateScoreRequest {
  uint32 user_id = 1;
  uint32 contest_id = 2;
  uint32 prediction_id = 3;
  double points = 4;
  google.protobuf.Timestamp submitted_at = 5;
  google.protobuf.Timestamp event_date = 6;
}
```

2. Then in scoring_service.go, after streak multiplier calculation (around line 100):
```go
// Calculate time coefficient
timeCoefficient := 1.0
if req.SubmittedAt != nil && req.EventDate != nil {
    timeCoefficient = models.CalculateTimeCoefficient(
        req.SubmittedAt.AsTime(),
        req.EventDate.AsTime(),
    )
}

// Apply both multipliers
finalPoints := basePoints * multiplier * timeCoefficient

// Create score model with both multipliers
score := &models.Score{
    UserID:          req.UserId,
    ContestID:       req.ContestId,
    PredictionID:    req.PredictionId,
    Points:          finalPoints,
    TimeCoefficient: timeCoefficient,
}
```

3. Update log message:
```go
log.Printf("[INFO] Score created: user=%d, contest=%d, base=%.2f, streak=%.2fx, time=%.2fx, final=%.2f",
    req.UserId, req.ContestId, basePoints, multiplier, timeCoefficient, finalPoints)
```

### Task 5: UPDATE `backend/scoring-service/internal/service/scoring_service.go` modelToProto

- **IMPLEMENT**: Add TimeCoefficient to proto conversion
- **PATTERN**: Follow existing field mapping
- **VALIDATE**: `cd backend/scoring-service && go build ./...`

Update modelToProto function (around line 180):
```go
func (s *ScoringService) modelToProto(score *models.Score) *pb.Score {
    return &pb.Score{
        Id:              uint32(score.ID),
        UserId:          uint32(score.UserID),
        ContestId:       uint32(score.ContestID),
        PredictionId:    uint32(score.PredictionID),
        Points:          score.Points,
        TimeCoefficient: score.TimeCoefficient,
        ScoredAt:        timestamppb.New(score.ScoredAt),
        CreatedAt:       timestamppb.New(score.CreatedAt),
        UpdatedAt:       timestamppb.New(score.UpdatedAt),
    }
}
```

### Task 6: UPDATE `backend/proto/prediction.proto`

- **IMPLEMENT**: Add GetPotentialCoefficient RPC for frontend
- **PATTERN**: Follow existing RPC patterns
- **VALIDATE**: `cd backend && protoc --proto_path=proto --go_out=shared --go-grpc_out=shared proto/prediction.proto`

Add messages after ListPropTypesResponse:
```protobuf
message GetPotentialCoefficientRequest {
  uint32 event_id = 1;
}

message GetPotentialCoefficientResponse {
  common.Response response = 1;
  double coefficient = 2;
  string tier = 3;
  double hours_until_event = 4;
}
```

Add RPC to PredictionService:
```protobuf
rpc GetPotentialCoefficient(GetPotentialCoefficientRequest) returns (GetPotentialCoefficientResponse) {
  option (google.api.http) = {
    get: "/v1/events/{event_id}/coefficient"
  };
}
```

### Task 7: UPDATE `backend/prediction-service/internal/service/prediction_service.go`

- **IMPLEMENT**: Add GetPotentialCoefficient method
- **PATTERN**: Follow existing service method patterns
- **IMPORTS**: Add coefficient calculation (copy function or import)
- **VALIDATE**: `cd backend/prediction-service && go build ./...`

Add method:
```go
// GetPotentialCoefficient calculates the current time coefficient for an event
func (s *PredictionService) GetPotentialCoefficient(ctx context.Context, req *pb.GetPotentialCoefficientRequest) (*pb.GetPotentialCoefficientResponse, error) {
    event, err := s.eventRepo.GetByID(uint(req.EventId))
    if err != nil {
        return &pb.GetPotentialCoefficientResponse{
            Response: &common.Response{
                Success:   false,
                Message:   "Event not found",
                Code:      int32(common.ErrorCode_NOT_FOUND),
                Timestamp: timestamppb.Now(),
            },
        }, nil
    }

    now := time.Now().UTC()
    hoursUntilEvent := event.EventDate.Sub(now).Hours()
    coefficient := calculateTimeCoefficient(now, event.EventDate)
    tier := getCoefficientTier(coefficient)

    return &pb.GetPotentialCoefficientResponse{
        Response: &common.Response{
            Success:   true,
            Message:   "Coefficient calculated",
            Code:      0,
            Timestamp: timestamppb.Now(),
        },
        Coefficient:     coefficient,
        Tier:            tier,
        HoursUntilEvent: hoursUntilEvent,
    }, nil
}

// Helper functions (add at bottom of file)
func calculateTimeCoefficient(submittedAt, eventDate time.Time) float64 {
    hoursUntilEvent := eventDate.Sub(submittedAt).Hours()
    switch {
    case hoursUntilEvent >= 168:
        return 2.0
    case hoursUntilEvent >= 72:
        return 1.5
    case hoursUntilEvent >= 24:
        return 1.25
    case hoursUntilEvent >= 12:
        return 1.1
    default:
        return 1.0
    }
}

func getCoefficientTier(coefficient float64) string {
    switch coefficient {
    case 2.0:
        return "Early Bird"
    case 1.5:
        return "Ahead of Time"
    case 1.25:
        return "Timely"
    case 1.1:
        return "Last Minute"
    default:
        return "Standard"
    }
}
```

### Task 8: UPDATE `scripts/init-db.sql`

- **IMPLEMENT**: Add time_coefficient column to scores table
- **PATTERN**: Follow existing column definitions
- **GOTCHA**: Add default value for existing rows
- **VALIDATE**: Check SQL syntax

Add to scores table definition or add migration:
```sql
-- Add time_coefficient column to scores table
ALTER TABLE scores ADD COLUMN IF NOT EXISTS time_coefficient DOUBLE PRECISION NOT NULL DEFAULT 1.0;
```

### Task 9: UPDATE `frontend/src/types/scoring.types.ts`

- **IMPLEMENT**: Add timeCoefficient to Score interface
- **PATTERN**: Follow existing type patterns
- **VALIDATE**: `cd frontend && npm run lint`

Add to Score interface:
```typescript
export interface Score {
  id: number
  userId: number
  contestId: number
  predictionId: number
  points: number
  timeCoefficient: number  // Add this field
  scoredAt: string
  createdAt: string
  updatedAt: string
}
```

### Task 10: UPDATE `frontend/src/types/prediction.types.ts`

- **IMPLEMENT**: Add coefficient types
- **PATTERN**: Follow existing type patterns
- **VALIDATE**: `cd frontend && npm run lint`

Add at end of file:
```typescript
// Time coefficient types
export interface PotentialCoefficientResponse {
  response: ApiResponse
  coefficient: number
  tier: string
  hoursUntilEvent: number
}
```

### Task 11: UPDATE `frontend/src/services/prediction-service.ts`

- **IMPLEMENT**: Add getPotentialCoefficient method
- **PATTERN**: Follow existing service method patterns
- **VALIDATE**: `cd frontend && npm run lint`

Add method to PredictionService class:
```typescript
// Get potential coefficient for an event
async getPotentialCoefficient(eventId: number): Promise<{
  coefficient: number
  tier: string
  hoursUntilEvent: number
}> {
  const response = await grpcClient.get<PotentialCoefficientResponse>(
    `${this.eventsPath}/${eventId}/coefficient`
  )
  return {
    coefficient: response.coefficient,
    tier: response.tier,
    hoursUntilEvent: response.hoursUntilEvent,
  }
}
```

Add import at top:
```typescript
import type { PotentialCoefficientResponse } from '../types/prediction.types'
```

### Task 12: CREATE `frontend/src/components/predictions/CoefficientIndicator.tsx`

- **IMPLEMENT**: Visual component showing current coefficient
- **PATTERN**: Follow EventCard component pattern
- **IMPORTS**: MUI components
- **VALIDATE**: `cd frontend && npm run lint`

```tsx
import React from 'react'
import { Box, Chip, Tooltip, Typography } from '@mui/material'
import { TrendingUp, AccessTime } from '@mui/icons-material'

interface CoefficientIndicatorProps {
  coefficient: number
  tier: string
  hoursUntilEvent: number
  compact?: boolean
}

const getColor = (coefficient: number): 'success' | 'info' | 'warning' | 'default' => {
  if (coefficient >= 2.0) return 'success'
  if (coefficient >= 1.5) return 'info'
  if (coefficient >= 1.25) return 'warning'
  return 'default'
}

export const CoefficientIndicator: React.FC<CoefficientIndicatorProps> = ({
  coefficient,
  tier,
  hoursUntilEvent,
  compact = false,
}) => {
  const formatTimeRemaining = (hours: number): string => {
    if (hours >= 168) return `${Math.floor(hours / 24)} days`
    if (hours >= 24) return `${Math.floor(hours / 24)} days ${Math.floor(hours % 24)}h`
    return `${Math.floor(hours)}h ${Math.floor((hours % 1) * 60)}m`
  }

  if (compact) {
    return (
      <Tooltip title={`${tier} - ${coefficient}x points (${formatTimeRemaining(hoursUntilEvent)} left)`}>
        <Chip
          icon={<TrendingUp />}
          label={`${coefficient}x`}
          size="small"
          color={getColor(coefficient)}
          variant="outlined"
        />
      </Tooltip>
    )
  }

  return (
    <Box sx={{ p: 1.5, bgcolor: 'action.hover', borderRadius: 1 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
        <TrendingUp color={getColor(coefficient)} fontSize="small" />
        <Typography variant="subtitle2" fontWeight="bold">
          {coefficient}x Points
        </Typography>
        <Chip label={tier} size="small" color={getColor(coefficient)} />
      </Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
        <AccessTime fontSize="small" color="action" />
        <Typography variant="caption" color="text.secondary">
          {formatTimeRemaining(hoursUntilEvent)} until coefficient drops
        </Typography>
      </Box>
    </Box>
  )
}

export default CoefficientIndicator
```

### Task 13: UPDATE `frontend/src/hooks/use-predictions.ts` (or create if needed)

- **IMPLEMENT**: Add usePotentialCoefficient hook
- **PATTERN**: Follow existing React Query hook patterns
- **VALIDATE**: `cd frontend && npm run lint`

Check if file exists, then add:
```typescript
import { useQuery } from '@tanstack/react-query'
import { predictionService } from '../services/prediction-service'

export const usePotentialCoefficient = (eventId: number | undefined) => {
  return useQuery({
    queryKey: ['coefficient', eventId],
    queryFn: () => predictionService.getPotentialCoefficient(eventId!),
    enabled: !!eventId,
    refetchInterval: 60000, // Refresh every minute
  })
}
```

### Task 14: UPDATE `frontend/src/components/predictions/PredictionForm.tsx`

- **IMPLEMENT**: Display potential points with coefficient
- **PATTERN**: Follow existing form patterns
- **IMPORTS**: Add CoefficientIndicator, usePotentialCoefficient
- **VALIDATE**: `cd frontend && npm run lint`

Add imports:
```typescript
import { CoefficientIndicator } from './CoefficientIndicator'
import { usePotentialCoefficient } from '../../hooks/use-predictions'
```

Add hook call inside component (after existing hooks):
```typescript
const { data: coefficientData } = usePotentialCoefficient(event?.id)
```

Add coefficient display after event info box (around line 75, after the Divider):
```tsx
{coefficientData && coefficientData.coefficient > 1 && (
  <Box sx={{ mb: 2 }}>
    <CoefficientIndicator
      coefficient={coefficientData.coefficient}
      tier={coefficientData.tier}
      hoursUntilEvent={coefficientData.hoursUntilEvent}
    />
  </Box>
)}
```

### Task 15: UPDATE `frontend/src/components/predictions/EventCard.tsx`

- **IMPLEMENT**: Show coefficient indicator on event cards
- **PATTERN**: Follow existing card layout
- **IMPORTS**: Add CoefficientIndicator, usePotentialCoefficient
- **VALIDATE**: `cd frontend && npm run lint`

Add imports:
```typescript
import { CoefficientIndicator } from './CoefficientIndicator'
import { usePotentialCoefficient } from '../../hooks/use-predictions'
```

Add hook inside component:
```typescript
const { data: coefficientData } = usePotentialCoefficient(event.id)
```

Add coefficient display after event date (around line 65, before hasPrediction check):
```tsx
{coefficientData && coefficientData.coefficient > 1 && (
  <Box sx={{ mt: 1 }}>
    <CoefficientIndicator
      coefficient={coefficientData.coefficient}
      tier={coefficientData.tier}
      hoursUntilEvent={coefficientData.hoursUntilEvent}
      compact
    />
  </Box>
)}
```

### Task 16: Regenerate Proto Files

- **IMPLEMENT**: Regenerate all proto-generated Go files
- **VALIDATE**: `cd backend && make proto` or manual protoc commands

```bash
cd backend
protoc --proto_path=proto \
  --go_out=shared --go_opt=paths=source_relative \
  --go-grpc_out=shared --go-grpc_opt=paths=source_relative \
  proto/scoring.proto proto/prediction.proto
```

---

## TESTING STRATEGY

### Unit Tests

**Backend Tests** (`tests/scoring-service/coefficient_test.go`):
```go
func TestCalculateTimeCoefficient(t *testing.T) {
    eventDate := time.Now().Add(8 * 24 * time.Hour) // 8 days from now
    
    tests := []struct {
        name        string
        submittedAt time.Time
        expected    float64
    }{
        {"7+ days early", time.Now(), 2.0},
        {"5 days early", time.Now().Add(3 * 24 * time.Hour), 1.5},
        {"2 days early", time.Now().Add(6 * 24 * time.Hour), 1.25},
        {"18 hours early", time.Now().Add(7*24*time.Hour + 6*time.Hour), 1.1},
        {"6 hours early", time.Now().Add(7*24*time.Hour + 18*time.Hour), 1.0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := models.CalculateTimeCoefficient(tt.submittedAt, eventDate)
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

**Frontend Tests** (`tests/frontend/src/components/CoefficientIndicator.test.tsx`):
```tsx
import { render, screen } from '@testing-library/react'
import { CoefficientIndicator } from '../../../frontend/src/components/predictions/CoefficientIndicator'

describe('CoefficientIndicator', () => {
  it('displays coefficient value', () => {
    render(<CoefficientIndicator coefficient={2.0} tier="Early Bird" hoursUntilEvent={200} />)
    expect(screen.getByText('2x Points')).toBeInTheDocument()
  })

  it('shows compact version', () => {
    render(<CoefficientIndicator coefficient={1.5} tier="Ahead of Time" hoursUntilEvent={100} compact />)
    expect(screen.getByText('1.5x')).toBeInTheDocument()
  })
})
```

### Integration Tests

- Test full flow: submit prediction → calculate score with coefficient
- Verify coefficient is stored in database
- Test API endpoint returns correct coefficient for different time ranges

### Edge Cases

- Event date in the past (should return 1.0)
- Event date exactly at boundary times
- Null/missing event date handling
- Coefficient calculation with different timezones

---

## VALIDATION COMMANDS

### Level 1: Syntax & Style

```bash
# Backend
cd backend/scoring-service && go build ./...
cd backend/prediction-service && go build ./...
cd backend && go vet ./...

# Frontend
cd frontend && npm run lint
cd frontend && npx tsc --noEmit
```

### Level 2: Unit Tests

```bash
# Backend
cd backend && go test ./scoring-service/... -v
cd backend && go test ./prediction-service/... -v

# Frontend
cd frontend && npm test
```

### Level 3: Integration Tests

```bash
# Start services
make docker-up

# Test coefficient endpoint
curl http://localhost:8080/v1/events/1/coefficient

# Test score creation with coefficient
curl -X POST http://localhost:8080/v1/scores \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"user_id": 1, "contest_id": 1, "prediction_id": 1, "points": 10}'
```

### Level 4: Manual Validation

1. Create an event 8 days in the future
2. Check coefficient endpoint returns 2.0x
3. Submit prediction, verify coefficient displayed in form
4. Wait/modify event date, verify coefficient changes
5. Score prediction, verify final points include coefficient

---

## ACCEPTANCE CRITERIA

- [ ] Time coefficient calculated based on submission time vs event date
- [ ] Coefficient applied alongside existing streak multiplier
- [ ] Score model stores time_coefficient value
- [ ] API endpoint returns current potential coefficient for events
- [ ] Frontend displays coefficient in prediction form
- [ ] Frontend shows coefficient indicator on event cards
- [ ] Coefficient tiers: 2.0x (7+ days), 1.5x (3-7 days), 1.25x (1-3 days), 1.1x (12-24h), 1.0x (<12h)
- [ ] All validation commands pass with zero errors
- [ ] Unit tests cover coefficient calculation edge cases
- [ ] No regressions in existing scoring functionality

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms feature works
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

### Design Decisions

1. **Coefficient tiers chosen for balance**: 2.0x max rewards very early predictions without being overpowering. The tiers create meaningful decision points for users.

2. **Multiplicative with streak**: Coefficients multiply together (base × streak × time) rather than adding, creating compound rewards for skilled early predictors.

3. **Stored on score record**: TimeCoefficient is stored for audit trail and analytics, not recalculated.

4. **Real-time coefficient display**: Frontend polls coefficient every minute to show users the current multiplier and urgency.

### Trade-offs

- **Complexity vs Value**: Adding another multiplier increases scoring complexity but significantly improves engagement incentives.
- **API call overhead**: Each event card fetches coefficient, but caching and batching can optimize this later.

### Future Enhancements

- Contest-specific coefficient configurations
- Different coefficient curves per sport type
- Coefficient history visualization in analytics
- Push notifications when coefficient is about to drop

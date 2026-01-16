# Feature: Predictions UI

## Feature Description

Frontend interface for users to submit, view, update, and manage their sports predictions within contests. Users can browse available events, submit predictions with flexible data formats (winner, exact score, etc.), track their prediction history, and see prediction status (pending/scored/cancelled).

## User Story

As a **contest participant**
I want to **submit and manage predictions for sports events**
So that **I can compete in prediction contests and track my performance**

## Problem Statement

The platform has a complete backend prediction service but no frontend UI. Users cannot:
- View available events for prediction
- Submit predictions for events within contests
- View/edit their existing predictions
- Track prediction status and results

## Solution Statement

Build a comprehensive Predictions UI with:
- Event browsing with filtering by sport/status
- Prediction submission form with flexible data input
- User predictions list with status indicators
- Edit/delete capabilities for pending predictions
- Integration with existing contest context

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: Frontend (React), API Gateway integration
**Dependencies**: prediction-service (backend), contest-service (for context)

---

## CONTEXT REFERENCES

### Relevant Codebase Files - MUST READ BEFORE IMPLEMENTING

- `backend/proto/prediction.proto` - gRPC API definitions for predictions and events
- `frontend/src/types/contest.types.ts` - Type patterns to mirror
- `frontend/src/services/contest-service.ts` - Service class pattern to follow
- `frontend/src/hooks/use-contests.ts` - React Query hooks pattern
- `frontend/src/utils/validation.ts` - Zod validation schema patterns
- `frontend/src/components/contests/ContestList.tsx` - MaterialReactTable pattern
- `frontend/src/components/contests/ContestForm.tsx` - Dialog form pattern with react-hook-form
- `frontend/src/pages/ContestsPage.tsx` - Page structure with tabs pattern

### New Files to Create

- `frontend/src/types/prediction.types.ts` - TypeScript interfaces
- `frontend/src/utils/prediction-validation.ts` - Zod schemas
- `frontend/src/services/prediction-service.ts` - API service class
- `frontend/src/hooks/use-predictions.ts` - React Query hooks
- `frontend/src/components/predictions/EventList.tsx` - Events table
- `frontend/src/components/predictions/EventCard.tsx` - Event display card
- `frontend/src/components/predictions/PredictionForm.tsx` - Submit/edit form
- `frontend/src/components/predictions/PredictionList.tsx` - User predictions table
- `frontend/src/components/predictions/PredictionCard.tsx` - Prediction display
- `frontend/src/pages/PredictionsPage.tsx` - Main predictions page

### Patterns to Follow

**Naming Conventions:**
- Types: `PascalCase` (e.g., `Prediction`, `Event`, `SubmitPredictionRequest`)
- Hooks: `use[Entity]s`, `use[Entity]`, `useCreate[Entity]`, `useUpdate[Entity]`, `useDelete[Entity]`
- Services: `[entity]Service` singleton class
- Components: `PascalCase.tsx`

**Service Pattern (from contest-service.ts):**
```typescript
class PredictionService {
  private basePath = '/v1/predictions'
  
  async submitPrediction(request: SubmitPredictionRequest): Promise<Prediction> {
    const response = await grpcClient.post<SubmitPredictionResponse>(this.basePath, request)
    return response.prediction
  }
}
export const predictionService = new PredictionService()
```

**Hook Pattern (from use-contests.ts):**
```typescript
export const predictionKeys = {
  all: ['predictions'] as const,
  lists: () => [...predictionKeys.all, 'list'] as const,
  list: (filters: ListPredictionsRequest) => [...predictionKeys.lists(), filters] as const,
}

export const usePredictions = (request: GetUserPredictionsRequest) => {
  return useQuery({
    queryKey: predictionKeys.list(request),
    queryFn: () => predictionService.getUserPredictions(request),
  })
}
```

**Validation Pattern (from validation.ts):**
```typescript
export const predictionSchema = z.object({
  eventId: z.number().min(1, 'Event is required'),
  predictionData: z.string().min(1, 'Prediction is required').max(5000),
})
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation (Types, Validation, Service)

Create TypeScript types matching proto definitions, Zod validation schemas, and API service class.

### Phase 2: React Query Hooks

Build data fetching and mutation hooks following existing patterns.

### Phase 3: UI Components

Create event browsing, prediction form, and prediction list components.

### Phase 4: Page Integration

Build PredictionsPage and integrate into App routing.

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `frontend/src/types/prediction.types.ts`

**IMPLEMENT**: TypeScript interfaces matching `backend/proto/prediction.proto`

```typescript
// Prediction entity
export interface Prediction {
  id: number
  contestId: number
  userId: number
  eventId: number
  predictionData: string // JSON string
  status: 'pending' | 'scored' | 'cancelled'
  submittedAt: string
  createdAt: string
  updatedAt: string
}

// Event entity
export interface Event {
  id: number
  title: string
  sportType: string
  homeTeam: string
  awayTeam: string
  eventDate: string
  status: 'scheduled' | 'live' | 'completed' | 'cancelled'
  resultData: string
  createdAt: string
  updatedAt: string
}

// Request/Response types matching proto
export interface SubmitPredictionRequest {
  contestId: number
  eventId: number
  predictionData: string
}

export interface GetUserPredictionsRequest {
  contestId: number
  pagination?: PaginationRequest
}

export interface UpdatePredictionRequest {
  id: number
  predictionData: string
}

export interface ListEventsRequest {
  sportType?: string
  status?: string
  pagination?: PaginationRequest
}

// Import PaginationRequest/Response from contest.types.ts or common.types.ts
```

**PATTERN**: Mirror `frontend/src/types/contest.types.ts` structure
**IMPORTS**: `import type { PaginationRequest, PaginationResponse, ApiResponse } from './common.types'`
**VALIDATE**: `npx tsc --noEmit`

---

### Task 2: CREATE `frontend/src/utils/prediction-validation.ts`

**IMPLEMENT**: Zod validation schemas for prediction forms

```typescript
import { z } from 'zod'

export const predictionSchema = z.object({
  eventId: z.number().min(1, 'Event selection is required'),
  predictionData: z
    .string()
    .min(1, 'Prediction is required')
    .max(5000, 'Prediction cannot exceed 5000 characters'),
})

export type PredictionFormData = z.infer<typeof predictionSchema>

// For score-based predictions
export const scorePredictionSchema = z.object({
  homeScore: z.number().min(0, 'Score must be 0 or greater'),
  awayScore: z.number().min(0, 'Score must be 0 or greater'),
})

// For winner predictions
export const winnerPredictionSchema = z.object({
  winner: z.enum(['home', 'away', 'draw'], {
    errorMap: () => ({ message: 'Select a winner' }),
  }),
})
```

**PATTERN**: Mirror `frontend/src/utils/validation.ts`
**VALIDATE**: `npx tsc --noEmit`

---

### Task 3: CREATE `frontend/src/services/prediction-service.ts`

**IMPLEMENT**: API service class with all prediction endpoints

```typescript
import grpcClient from './grpc-client'
import type { /* all types */ } from '../types/prediction.types'

class PredictionService {
  private predictionsPath = '/v1/predictions'
  private eventsPath = '/v1/events'

  // Predictions
  async submitPrediction(request: SubmitPredictionRequest): Promise<Prediction>
  async getPrediction(id: number): Promise<Prediction>
  async getUserPredictions(request: GetUserPredictionsRequest): Promise<{predictions, pagination}>
  async updatePrediction(request: UpdatePredictionRequest): Promise<Prediction>
  async deletePrediction(id: number): Promise<void>

  // Events
  async getEvent(id: number): Promise<Event>
  async listEvents(request: ListEventsRequest): Promise<{events, pagination}>
}

export const predictionService = new PredictionService()
```

**PATTERN**: Mirror `frontend/src/services/contest-service.ts` exactly
**IMPORTS**: `import grpcClient from './grpc-client'`
**VALIDATE**: `npx tsc --noEmit`

---

### Task 4: CREATE `frontend/src/hooks/use-predictions.ts`

**IMPLEMENT**: React Query hooks for predictions and events

```typescript
// Query keys
export const predictionKeys = {
  all: ['predictions'] as const,
  lists: () => [...predictionKeys.all, 'list'] as const,
  list: (contestId: number) => [...predictionKeys.lists(), contestId] as const,
  details: () => [...predictionKeys.all, 'detail'] as const,
  detail: (id: number) => [...predictionKeys.details(), id] as const,
}

export const eventKeys = {
  all: ['events'] as const,
  lists: () => [...eventKeys.all, 'list'] as const,
  list: (filters: ListEventsRequest) => [...eventKeys.lists(), filters] as const,
  details: () => [...eventKeys.all, 'detail'] as const,
  detail: (id: number) => [...eventKeys.details(), id] as const,
}

// Hooks: useUserPredictions, usePrediction, useEvents, useEvent
// Mutations: useSubmitPrediction, useUpdatePrediction, useDeletePrediction
```

**PATTERN**: Mirror `frontend/src/hooks/use-contests.ts` exactly
**IMPORTS**: `useQuery, useMutation, useQueryClient` from `@tanstack/react-query`
**VALIDATE**: `npx tsc --noEmit`

---

### Task 5: CREATE `frontend/src/components/predictions/EventCard.tsx`

**IMPLEMENT**: Card component displaying event details with prediction action

```typescript
interface EventCardProps {
  event: Event
  onPredict: (event: Event) => void
  existingPrediction?: Prediction
  disabled?: boolean
}
```

Display: title, teams (home vs away), event date, status chip, predict button
**PATTERN**: Use MUI Card, Chip, Button components
**VALIDATE**: `npx tsc --noEmit`

---

### Task 6: CREATE `frontend/src/components/predictions/EventList.tsx`

**IMPLEMENT**: Grid/list of EventCards with filtering

```typescript
interface EventListProps {
  contestId: number
  onPredict: (event: Event) => void
  userPredictions: Prediction[]
}
```

Features: sport type filter, status filter, pagination, loading states
**PATTERN**: Mirror filtering from `ContestList.tsx`
**VALIDATE**: `npx tsc --noEmit`

---

### Task 7: CREATE `frontend/src/components/predictions/PredictionForm.tsx`

**IMPLEMENT**: Dialog form for submitting/editing predictions

```typescript
interface PredictionFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: PredictionFormData) => void
  event: Event | null
  prediction?: Prediction | null
  loading?: boolean
}
```

Features:
- Display event info (teams, date)
- Flexible prediction input (JSON or structured)
- Winner selection radio buttons
- Score input fields
- Validation with Zod

**PATTERN**: Mirror `frontend/src/components/contests/ContestForm.tsx`
**IMPORTS**: Dialog, TextField, RadioGroup from MUI; useForm, Controller from react-hook-form
**VALIDATE**: `npx tsc --noEmit`

---

### Task 8: CREATE `frontend/src/components/predictions/PredictionList.tsx`

**IMPLEMENT**: Table showing user's predictions with actions

```typescript
interface PredictionListProps {
  contestId: number
  onEdit: (prediction: Prediction) => void
}
```

Columns: Event, Prediction, Status, Submitted At, Actions (edit/delete for pending)
**PATTERN**: Mirror `frontend/src/components/contests/ContestList.tsx` with MaterialReactTable
**VALIDATE**: `npx tsc --noEmit`

---

### Task 9: CREATE `frontend/src/pages/PredictionsPage.tsx`

**IMPLEMENT**: Main predictions page with tabs

```typescript
// Two tabs:
// 1. "Available Events" - EventList for making predictions
// 2. "My Predictions" - PredictionList showing user's predictions

// State: selectedContest (from URL or dropdown), form dialogs
// Integration: Contest selector dropdown at top
```

**PATTERN**: Mirror `frontend/src/pages/ContestsPage.tsx` tab structure
**VALIDATE**: `npx tsc --noEmit`

---

### Task 10: UPDATE `frontend/src/App.tsx`

**IMPLEMENT**: Add predictions route and navigation

```typescript
// Add import
import PredictionsPage from './pages/PredictionsPage'

// Add navigation button in AppBarContent (after Sports)
<Button color="inherit" component={Link} to="/predictions">Predictions</Button>

// Add route
<Route 
  path="/predictions" 
  element={
    <ProtectedRoute>
      <PredictionsPage />
    </ProtectedRoute>
  } 
/>
```

**PATTERN**: Mirror existing route structure
**VALIDATE**: `npm run build`

---

## TESTING STRATEGY

### Unit Tests

Test validation schemas with valid/invalid inputs:
- `predictionSchema` - required fields, max length
- `scorePredictionSchema` - non-negative scores
- `winnerPredictionSchema` - valid enum values

### Integration Tests

Test hooks with mock API responses:
- `useUserPredictions` - fetches and caches correctly
- `useSubmitPrediction` - invalidates queries on success
- `useEvents` - pagination works correctly

### Manual Testing

1. Navigate to /predictions
2. Select a contest from dropdown
3. View available events
4. Click "Predict" on an event
5. Submit prediction with score/winner
6. Verify prediction appears in "My Predictions" tab
7. Edit a pending prediction
8. Delete a pending prediction
9. Verify cannot edit/delete scored predictions

---

## VALIDATION COMMANDS

### Level 1: TypeScript Compilation
```bash
cd frontend && npx tsc --noEmit
```

### Level 2: Linting
```bash
cd frontend && npm run lint
```

### Level 3: Build
```bash
cd frontend && npm run build
```

### Level 4: Unit Tests
```bash
cd frontend && npm run test
```

### Level 5: Manual Validation
1. Start backend: `docker-compose up -d`
2. Start frontend: `cd frontend && npm run dev`
3. Login and navigate to /predictions
4. Complete full prediction workflow

---

## ACCEPTANCE CRITERIA

- [ ] Users can view available events filtered by sport/status
- [ ] Users can submit predictions for events in active contests
- [ ] Users can view their prediction history per contest
- [ ] Users can edit pending predictions
- [ ] Users can delete pending predictions
- [ ] Scored/cancelled predictions are read-only
- [ ] Form validation prevents invalid submissions
- [ ] Loading and error states display correctly
- [ ] Navigation integrates with existing app structure
- [ ] All TypeScript compiles without errors
- [ ] All existing tests continue to pass

---

## COMPLETION CHECKLIST

- [ ] All 10 tasks completed in order
- [ ] TypeScript compilation passes
- [ ] Linting passes
- [ ] Build succeeds
- [ ] Manual testing confirms feature works
- [ ] No regressions in existing functionality

---

## NOTES

**Prediction Data Format**: The `predictionData` field is a JSON string allowing flexible prediction types:
- Winner: `{"winner": "home"}` or `{"winner": "away"}` or `{"winner": "draw"}`
- Score: `{"homeScore": 2, "awayScore": 1}`
- Combined: `{"winner": "home", "homeScore": 2, "awayScore": 1}`

**Event Status Logic**:
- Only `scheduled` events with future `eventDate` accept predictions
- `live`, `completed`, `cancelled` events are view-only

**Contest Context**: Predictions are always within a contest context. The page should:
1. Allow contest selection via dropdown
2. Remember selected contest in URL params
3. Filter events/predictions by selected contest

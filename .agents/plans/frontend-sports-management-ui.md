# Feature: Frontend Sports Management UI

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Complete React TypeScript frontend for managing sports data entities (Sports, Leagues, Teams, Matches). This UI enables administrators to create, view, edit, and delete sports-related data that powers the prediction contests platform. The implementation follows existing patterns from the contest management UI.

## User Story

As a platform administrator
I want to manage sports, leagues, teams, and matches through a web interface
So that I can configure the sports data needed for prediction contests

## Problem Statement

The backend Sports Service (port 8088) is fully implemented with gRPC endpoints for CRUD operations on Sports, Leagues, Teams, and Matches. However, there is no frontend UI to interact with these endpoints, requiring administrators to use direct API calls.

## Solution Statement

Build a comprehensive Sports Management UI following the established patterns from ContestsPage, including:
- TypeScript types matching the sports.proto definitions
- Service layer using the existing grpcClient pattern
- React Query hooks for data fetching and mutations
- Material-UI components with MaterialReactTable for listings
- Form dialogs with Zod validation for create/edit operations
- Tab-based navigation for Sports → Leagues → Teams → Matches

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: Frontend (React), API Gateway integration
**Dependencies**: @tanstack/react-query, @mui/material, material-react-table, zod, react-hook-form

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/proto/sports.proto` (full file) - Why: Defines all Sports Service messages and RPCs - source of truth for types
- `frontend/src/types/contest.types.ts` (full file) - Why: Pattern for TypeScript type definitions matching proto
- `frontend/src/services/grpc-client.ts` (full file) - Why: Base HTTP client to use for API calls
- `frontend/src/services/contest-service.ts` (full file) - Why: Pattern for service class implementation
- `frontend/src/hooks/use-contests.ts` (full file) - Why: Pattern for React Query hooks with mutations
- `frontend/src/components/contests/ContestList.tsx` (full file) - Why: Pattern for MaterialReactTable list component
- `frontend/src/components/contests/ContestForm.tsx` (full file) - Why: Pattern for form dialog with Zod validation
- `frontend/src/pages/ContestsPage.tsx` (full file) - Why: Pattern for page structure with tabs
- `frontend/src/utils/validation.ts` (full file) - Why: Pattern for Zod schemas
- `frontend/src/utils/date-utils.ts` (full file) - Why: Date formatting utilities to reuse
- `frontend/src/contexts/ToastContext.tsx` (full file) - Why: Toast notifications pattern
- `frontend/src/types/common.types.ts` (full file) - Why: Shared types (ApiResponse, PaginationRequest/Response)

### New Files to Create

- `frontend/src/types/sports.types.ts` - TypeScript types for Sport, League, Team, Match entities
- `frontend/src/services/sports-service.ts` - Service class for Sports API calls
- `frontend/src/hooks/use-sports.ts` - React Query hooks for sports data
- `frontend/src/components/sports/SportList.tsx` - MaterialReactTable for sports
- `frontend/src/components/sports/SportForm.tsx` - Create/Edit sport dialog
- `frontend/src/components/sports/LeagueList.tsx` - MaterialReactTable for leagues
- `frontend/src/components/sports/LeagueForm.tsx` - Create/Edit league dialog
- `frontend/src/components/sports/TeamList.tsx` - MaterialReactTable for teams
- `frontend/src/components/sports/TeamForm.tsx` - Create/Edit team dialog
- `frontend/src/components/sports/MatchList.tsx` - MaterialReactTable for matches
- `frontend/src/components/sports/MatchForm.tsx` - Create/Edit match dialog
- `frontend/src/pages/SportsPage.tsx` - Main page with tab navigation
- `frontend/src/utils/sports-validation.ts` - Zod schemas for sports entities

### Files to Modify

- `frontend/src/App.tsx` - Add /sports route and navigation link

### Relevant Documentation

- [Material React Table Docs](https://www.material-react-table.com/) - Table component patterns
- [React Hook Form + Zod](https://react-hook-form.com/get-started#SchemaValidation) - Form validation integration
- [TanStack Query Mutations](https://tanstack.com/query/latest/docs/react/guides/mutations) - Mutation patterns

### Patterns to Follow

**Naming Conventions:**
- Types: PascalCase (Sport, League, Team, Match)
- Interfaces: PascalCase with descriptive suffix (CreateSportRequest, ListSportsResponse)
- Hooks: camelCase with `use` prefix (useSports, useCreateSport)
- Services: camelCase class with singleton export (sportsService)
- Components: PascalCase (SportList, SportForm)

**Service Pattern (from contest-service.ts):**
```typescript
class SportsService {
  private basePath = '/v1/sports'
  
  async listSports(request: ListSportsRequest = {}): Promise<{...}> {
    const params = new URLSearchParams()
    // Build query params
    const response = await grpcClient.get<ListSportsResponse>(url)
    return { sports: response.sports, pagination: response.pagination }
  }
}
export const sportsService = new SportsService()
```

**Hook Pattern (from use-contests.ts):**
```typescript
export const sportsKeys = {
  all: ['sports'] as const,
  lists: () => [...sportsKeys.all, 'list'] as const,
  list: (filters: ListSportsRequest) => [...sportsKeys.lists(), filters] as const,
}

export const useSports = (request: ListSportsRequest = {}) => {
  return useQuery({
    queryKey: sportsKeys.list(request),
    queryFn: () => sportsService.listSports(request),
    staleTime: 5 * 60 * 1000,
  })
}
```

**Form Pattern (from ContestForm.tsx):**
```typescript
const { control, handleSubmit, reset, formState: { errors, isValid } } = useForm<SportFormData>({
  resolver: zodResolver(sportSchema),
  defaultValues,
  mode: 'onChange',
})
```

**List Pattern (from ContestList.tsx):**
```typescript
const columns = useMemo<MRT_ColumnDef<Sport>[]>(() => [...], [])
const table = useMaterialReactTable({
  columns,
  data: data?.sports ?? [],
  manualPagination: true,
  // ... config
})
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation (Types & Service)

Create TypeScript types matching sports.proto and service layer for API communication.

**Tasks:**
- Define Sport, League, Team, Match interfaces
- Define request/response types for all CRUD operations
- Implement SportsService class with all API methods
- Create Zod validation schemas

### Phase 2: Data Layer (Hooks)

Implement React Query hooks for data fetching and mutations.

**Tasks:**
- Create query keys structure for all entities
- Implement list/detail queries for each entity
- Implement create/update/delete mutations with cache invalidation
- Add toast notifications for success/error states

### Phase 3: UI Components (Lists & Forms)

Build Material-UI components for displaying and editing data.

**Tasks:**
- Create list components with MaterialReactTable
- Create form dialogs with validation
- Implement entity-specific features (status chips, foreign key selects)

### Phase 4: Integration (Page & Routing)

Integrate components into page and add routing.

**Tasks:**
- Create SportsPage with tab navigation
- Add route to App.tsx
- Add navigation link in AppBar

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### Task 1: CREATE `frontend/src/types/sports.types.ts`

- **IMPLEMENT**: TypeScript interfaces matching sports.proto messages
- **PATTERN**: Mirror `frontend/src/types/contest.types.ts` structure
- **IMPORTS**: Import from `./common.types` for ApiResponse, PaginationRequest, PaginationResponse
- **GOTCHA**: Use `number` for uint32 proto fields, `string` for timestamps (ISO format)
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/types/sports.types.ts`

```typescript
// Required interfaces:
// - Sport, League, Team, Match (entities)
// - Create/Update/List/Get/Delete request types for each
// - Response types matching proto (SportResponse, ListSportsResponse, etc.)
// - Form data types for each entity
// Match status enum: 'scheduled' | 'live' | 'finished' | 'cancelled' | 'postponed'
```

### Task 2: CREATE `frontend/src/utils/sports-validation.ts`

- **IMPLEMENT**: Zod schemas for Sport, League, Team, Match forms
- **PATTERN**: Mirror `frontend/src/utils/validation.ts` contestSchema pattern
- **IMPORTS**: `import { z } from 'zod'`
- **GOTCHA**: Slug validation should allow alphanumeric and hyphens only
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/utils/sports-validation.ts`

```typescript
// Required schemas:
// - sportSchema: name (required, max 100), slug (auto-generated or manual), description, iconUrl
// - leagueSchema: sportId (required), name (required), slug, country, season
// - teamSchema: sportId (required), name (required), slug, shortName, logoUrl, country
// - matchSchema: leagueId (required), homeTeamId (required), awayTeamId (required), scheduledAt (required, future date for create)
// Export form data types: SportFormData, LeagueFormData, TeamFormData, MatchFormData
```

### Task 3: CREATE `frontend/src/services/sports-service.ts`

- **IMPLEMENT**: SportsService class with CRUD methods for all 4 entities
- **PATTERN**: Mirror `frontend/src/services/contest-service.ts` exactly
- **IMPORTS**: `import grpcClient from './grpc-client'`, types from `../types/sports.types`
- **GOTCHA**: Different base paths: /v1/sports, /v1/leagues, /v1/teams, /v1/matches
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/services/sports-service.ts`

```typescript
// Required methods per entity:
// Sports: createSport, getSport, updateSport, deleteSport, listSports
// Leagues: createLeague, getLeague, updateLeague, deleteLeague, listLeagues
// Teams: createTeam, getTeam, updateTeam, deleteTeam, listTeams
// Matches: createMatch, getMatch, updateMatch, deleteMatch, listMatches
// Export singleton: export const sportsService = new SportsService()
```

### Task 4: CREATE `frontend/src/hooks/use-sports.ts`

- **IMPLEMENT**: React Query hooks for all sports entities
- **PATTERN**: Mirror `frontend/src/hooks/use-contests.ts` exactly
- **IMPORTS**: `@tanstack/react-query`, `../services/sports-service`, `../contexts/ToastContext`, types
- **GOTCHA**: Separate query keys for each entity type to avoid cache conflicts
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/hooks/use-sports.ts`

```typescript
// Required exports:
// Query keys: sportsKeys, leaguesKeys, teamsKeys, matchesKeys
// List hooks: useSports, useLeagues, useTeams, useMatches
// Detail hooks: useSport, useLeague, useTeam, useMatch
// Mutation hooks: useCreateSport, useUpdateSport, useDeleteSport (same for League, Team, Match)
// Total: 4 key objects + 4 list hooks + 4 detail hooks + 12 mutation hooks = 24 exports
```

### Task 5: CREATE `frontend/src/components/sports/SportList.tsx`

- **IMPLEMENT**: MaterialReactTable component for sports listing
- **PATTERN**: Mirror `frontend/src/components/contests/ContestList.tsx`
- **IMPORTS**: material-react-table, @mui/material, @mui/icons-material, hooks, types
- **GOTCHA**: Include is_active status chip (green for active, gray for inactive)
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/SportList.tsx`

```typescript
// Columns: id, name, slug, description (truncated), isActive (chip), createdAt
// Actions: Edit, Delete, Toggle Active
// Top toolbar: "Add Sport" button
// Props: onCreateSport, onEditSport callbacks
```

### Task 6: CREATE `frontend/src/components/sports/SportForm.tsx`

- **IMPLEMENT**: Dialog form for creating/editing sports
- **PATTERN**: Mirror `frontend/src/components/contests/ContestForm.tsx`
- **IMPORTS**: @mui/material, react-hook-form, @hookform/resolvers/zod, validation schema
- **GOTCHA**: Auto-generate slug from name if not provided (lowercase, replace spaces with hyphens)
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/SportForm.tsx`

```typescript
// Fields: name (required), slug (optional, auto-generated), description (textarea), iconUrl
// Props: open, onClose, onSubmit, sport (for edit mode), loading
```

### Task 7: CREATE `frontend/src/components/sports/LeagueList.tsx`

- **IMPLEMENT**: MaterialReactTable component for leagues listing
- **PATTERN**: Mirror SportList.tsx with league-specific columns
- **IMPORTS**: Same as SportList plus useSports for sport name lookup
- **GOTCHA**: Show sport name instead of sportId (fetch sports for lookup)
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/LeagueList.tsx`

```typescript
// Columns: id, name, sportName (from lookup), country, season, isActive, createdAt
// Filter: sportId dropdown to filter by sport
// Actions: Edit, Delete
// Props: onCreateLeague, onEditLeague, sportFilter (optional)
```

### Task 8: CREATE `frontend/src/components/sports/LeagueForm.tsx`

- **IMPLEMENT**: Dialog form for creating/editing leagues
- **PATTERN**: Mirror SportForm.tsx with league fields
- **IMPORTS**: Same as SportForm plus useSports for sport dropdown
- **GOTCHA**: Sport dropdown must load sports list for selection
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/LeagueForm.tsx`

```typescript
// Fields: sportId (dropdown, required), name (required), slug, country, season
// Sport dropdown: Load active sports using useSports hook
// Props: open, onClose, onSubmit, league (for edit), loading
```

### Task 9: CREATE `frontend/src/components/sports/TeamList.tsx`

- **IMPLEMENT**: MaterialReactTable component for teams listing
- **PATTERN**: Mirror LeagueList.tsx with team-specific columns
- **IMPORTS**: Same as LeagueList
- **GOTCHA**: Show logo as small image if logoUrl exists
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/TeamList.tsx`

```typescript
// Columns: id, logo (Avatar), name, shortName, sportName, country, isActive, createdAt
// Filter: sportId dropdown
// Actions: Edit, Delete
// Props: onCreateTeam, onEditTeam, sportFilter
```

### Task 10: CREATE `frontend/src/components/sports/TeamForm.tsx`

- **IMPLEMENT**: Dialog form for creating/editing teams
- **PATTERN**: Mirror LeagueForm.tsx with team fields
- **IMPORTS**: Same as LeagueForm
- **GOTCHA**: Logo preview if URL is provided
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/TeamForm.tsx`

```typescript
// Fields: sportId (dropdown), name, slug, shortName, logoUrl (with preview), country
// Props: open, onClose, onSubmit, team (for edit), loading
```

### Task 11: CREATE `frontend/src/components/sports/MatchList.tsx`

- **IMPLEMENT**: MaterialReactTable component for matches listing
- **PATTERN**: Mirror TeamList.tsx with match-specific columns
- **IMPORTS**: Same as TeamList plus useLeagues, useTeams for lookups
- **GOTCHA**: Show team names, format scheduledAt with formatDateTime, status chip colors
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/MatchList.tsx`

```typescript
// Columns: id, homeTeam vs awayTeam, league, scheduledAt, status (chip), score (if finished)
// Status colors: scheduled=default, live=warning, finished=success, cancelled=error, postponed=info
// Filters: leagueId, status dropdowns
// Actions: Edit, Delete, Update Score (for finished matches)
// Props: onCreateMatch, onEditMatch, leagueFilter, statusFilter
```

### Task 12: CREATE `frontend/src/components/sports/MatchForm.tsx`

- **IMPLEMENT**: Dialog form for creating/editing matches
- **PATTERN**: Mirror TeamForm.tsx with match fields
- **IMPORTS**: Same as TeamForm plus DateTimePicker, useLeagues, useTeams
- **GOTCHA**: Team dropdowns should filter by sport of selected league
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/components/sports/MatchForm.tsx`

```typescript
// Fields: leagueId (dropdown), homeTeamId (dropdown), awayTeamId (dropdown), scheduledAt (DateTimePicker)
// Edit mode additional: status (dropdown), homeScore, awayScore, resultData (JSON textarea)
// Validation: homeTeamId !== awayTeamId
// Props: open, onClose, onSubmit, match (for edit), loading
```

### Task 13: CREATE `frontend/src/pages/SportsPage.tsx`

- **IMPLEMENT**: Main page with tabs for Sports, Leagues, Teams, Matches
- **PATTERN**: Mirror `frontend/src/pages/ContestsPage.tsx` tab structure
- **IMPORTS**: All list and form components, hooks, @mui/material
- **GOTCHA**: Maintain selected entity state for edit operations across tab switches
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/pages/SportsPage.tsx`

```typescript
// Tabs: Sports (0), Leagues (1), Teams (2), Matches (3)
// State: activeTab, formOpen, selectedEntity (Sport|League|Team|Match|null), entityType
// Each tab renders corresponding List component
// Form dialogs rendered based on entityType
// Handle create/edit/delete for all entity types
```

### Task 14: UPDATE `frontend/src/App.tsx`

- **IMPLEMENT**: Add /sports route and navigation link in AppBar
- **PATTERN**: Mirror existing /contests route setup
- **IMPORTS**: Add SportsPage import
- **GOTCHA**: Add Sports link only for authenticated users (next to Contests)
- **VALIDATE**: `cd frontend && npx tsc --noEmit src/App.tsx`

```typescript
// Add import: import SportsPage from './pages/SportsPage'
// Add route: <Route path="/sports" element={<ProtectedRoute><SportsPage /></ProtectedRoute>} />
// Add AppBar button: <Button color="inherit" component={Link} to="/sports">Sports</Button>
// Position: After "Sports Prediction Contests" title, before user menu
```

---

## TESTING STRATEGY

### Unit Tests

Based on project patterns in `frontend/src/tests/`, create tests for:
- Validation schemas (sports-validation.ts)
- Service methods (mock grpcClient)
- Hook behavior (mock service responses)

### Integration Tests

- Form submission flows
- List pagination and filtering
- CRUD operations with cache invalidation

### Edge Cases

- Empty lists (no sports/leagues/teams/matches)
- Form validation errors
- API error handling
- Loading states
- Pagination boundary conditions
- Foreign key validation (league must have valid sport, etc.)

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
cd frontend && npx tsc --noEmit
cd frontend && npm run lint
```

### Level 2: Unit Tests

```bash
cd frontend && npm run test
```

### Level 3: Build Verification

```bash
cd frontend && npm run build
```

### Level 4: Manual Validation

1. Start backend services: `docker-compose up -d postgres redis`
2. Start API gateway: `cd backend/api-gateway && go run cmd/main.go`
3. Start sports service: `cd backend/sports-service && go run cmd/main.go`
4. Start frontend: `cd frontend && npm run dev`
5. Navigate to http://localhost:3000/sports
6. Test CRUD operations for each entity type:
   - Create a sport → verify in list
   - Edit the sport → verify changes
   - Create a league for the sport → verify sport dropdown works
   - Create teams for the sport → verify sport dropdown works
   - Create a match with teams → verify league/team dropdowns work
   - Delete entities → verify cascade behavior

---

## ACCEPTANCE CRITERIA

- [ ] All 4 entity types (Sport, League, Team, Match) have full CRUD UI
- [ ] TypeScript compiles without errors
- [ ] ESLint passes without warnings
- [ ] Forms validate input before submission
- [ ] Lists support pagination and filtering
- [ ] Toast notifications show for success/error states
- [ ] Foreign key relationships work (league→sport, team→sport, match→league/teams)
- [ ] Navigation between Sports page and Contests page works
- [ ] Protected route requires authentication
- [ ] UI follows existing Material-UI design patterns

---

## COMPLETION CHECKLIST

- [ ] All 14 tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] TypeScript compilation passes
- [ ] ESLint passes
- [ ] Manual testing confirms all CRUD operations work
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

### Design Decisions

1. **Single SportsService class**: All entity methods in one service class for simplicity, matching the backend service structure.

2. **Separate query keys per entity**: Prevents cache conflicts and allows granular invalidation.

3. **Tab-based navigation**: Matches ContestsPage pattern and provides clear entity separation.

4. **Foreign key dropdowns**: Load related entities (sports for leagues, teams for matches) to provide user-friendly selection.

5. **Slug auto-generation**: Reduces user input burden while allowing manual override.

### Potential Improvements (Out of Scope)

- Bulk operations (delete multiple)
- Advanced filtering (date ranges, search)
- Drag-and-drop reordering
- Import/export functionality
- Real-time updates via WebSocket

### Dependencies on Backend

- Sports Service must be running on port 8088
- API Gateway must proxy /v1/sports, /v1/leagues, /v1/teams, /v1/matches
- Database must have sports, leagues, teams, matches tables (already in init-db.sql)

# Feature: Complete Team Service Frontend Integration

The following plan completes the Team Service frontend integration. The backend gRPC service is fully implemented and registered in the API Gateway. The frontend has partial implementation with types, services, hooks, and basic components, but lacks full integration and comprehensive testing.

**IMPORTANT**: Validate documentation and codebase patterns before implementing. Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

Complete the frontend integration for the Team Tournaments feature, enabling users to create teams, invite members, manage team rosters, join contests as teams, and view team leaderboards. This feature strengthens the social aspect of the platform and enables viral growth through team invitations.

## User Story

**As a** sports prediction enthusiast  
**I want to** create and join teams with other users  
**So that** I can compete in contests as a group and share the experience with friends

## Problem Statement

The Team Service backend is fully implemented with gRPC endpoints, database models, and business logic. However, the frontend integration is incomplete:
- TeamsPage exists but has minimal functionality
- Team components are partially implemented but not fully integrated
- Team leaderboard integration in contests is missing
- E2E tests are minimal (only 2 basic tests)
- No integration with contest participation flow

This prevents users from experiencing the full Team Tournaments feature, which is a critical innovation for social engagement and user retention.

## Solution Statement

Complete the frontend integration by:
1. Enhancing TeamsPage with full CRUD operations and better UX
2. Integrating team leaderboards into ContestsPage
3. Adding team-based contest participation flow
4. Implementing comprehensive E2E tests covering all team workflows
5. Ensuring proper error handling and loading states throughout

## Feature Metadata

**Feature Type**: Enhancement (completing partial implementation)  
**Estimated Complexity**: Medium (4-6 hours)  
**Primary Systems Affected**: Frontend (pages, components, tests)  
**Dependencies**: 
- Backend Team Service (✅ Complete)
- API Gateway registration (✅ Complete)
- Frontend team-service.ts (✅ Complete)
- Frontend use-teams.ts hooks (✅ Complete)

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Backend (Reference Only - Already Complete)**:
- `backend/proto/team.proto` - gRPC service definition with all endpoints
- `backend/contest-service/internal/models/team.go` - Team model with validation
- `backend/contest-service/internal/service/team_service_grpc.go` (lines 1-100) - gRPC implementation
- `backend/api-gateway/internal/gateway/gateway.go` (lines 84-87) - Team service registration

**Frontend Types & Services (Already Complete)**:
- `frontend/src/types/team.types.ts` - All TypeScript interfaces for teams
- `frontend/src/services/team-service.ts` - Complete gRPC-Web client
- `frontend/src/hooks/use-teams.ts` - React Query hooks for all operations
- `frontend/src/utils/team-validation.ts` - Zod schemas for form validation

**Frontend Components (Partially Complete)**:
- `frontend/src/pages/TeamsPage.tsx` - Main page, needs enhancement
- `frontend/src/components/teams/TeamList.tsx` - List component, functional
- `frontend/src/components/teams/TeamForm.tsx` - Create/edit form, functional
- `frontend/src/components/teams/TeamMembers.tsx` - Members list, functional
- `frontend/src/components/teams/TeamInvite.tsx` - Invite code display, functional
- `frontend/src/components/teams/TeamLeaderboard.tsx` - Leaderboard component, needs integration

**Pages to Integrate With**:
- `frontend/src/pages/ContestsPage.tsx` - Add team leaderboard tab
- `frontend/src/App.tsx` - Verify teams route is registered

**Testing**:
- `frontend/tests/e2e/teams.spec.ts` - Minimal tests, needs expansion
- `frontend/tests/helpers/selectors.ts` - Test selectors defined
- `frontend/playwright.config.ts` - E2E test configuration

### New Files to Create

- None (all infrastructure exists)

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

**Ant Design Components** (already in use):
- [Table Component](https://ant.design/components/table) - Used in TeamList
- [Modal Component](https://ant.design/components/modal) - Used in TeamForm
- [Tabs Component](https://ant.design/components/tabs) - Used in TeamsPage and ContestsPage
- [Form Component](https://ant.design/components/form) - Used in TeamForm
- [List Component](https://ant.design/components/list) - Used in TeamMembers

**React Query** (TanStack Query v5):
- [Mutations](https://tanstack.com/query/v5/docs/react/guides/mutations) - Already used in hooks
- [Query Invalidation](https://tanstack.com/query/v5/docs/react/guides/query-invalidation) - For cache updates

**Playwright Testing**:
- [Test Fixtures](https://playwright.dev/docs/test-fixtures) - Already using auth.fixture
- [Assertions](https://playwright.dev/docs/test-assertions) - For E2E validation

### Patterns to Follow

**Naming Conventions**:
```typescript
// Component files: PascalCase.tsx
TeamList.tsx, TeamForm.tsx

// Hook files: use-[name].ts
use-teams.ts, use-contests.ts

// Type files: [name].types.ts
team.types.ts, contest.types.ts

// Utility files: kebab-case.ts
team-validation.ts, date-utils.ts
```

**Error Handling Pattern** (from existing code):
```typescript
// In mutations (use-teams.ts)
onError: (error: Error) => showToast(`Failed to create team: ${error.message}`, 'error')

// In components (TeamsPage.tsx)
try {
  await joinTeamMutation.mutateAsync({ inviteCode: data.inviteCode })
  showSuccess('Successfully joined team!')
} catch (error: any) {
  showError(error?.message || 'Failed to join team')
}
```

**Loading State Pattern** (from TeamList.tsx):
```typescript
const { data, isLoading, isError, error } = useTeams({ pagination })

// In render
if (isError) {
  return <Alert message="Error" description={error?.message} type="error" showIcon />
}

// In Table
<Table loading={isLoading || deleteTeamMutation.isPending} />
```

**Modal Pattern** (from ContestsPage.tsx):
```typescript
const [formOpen, setFormOpen] = useState(false)
const [editItem, setEditItem] = useState<Team | null>(null)

const handleCreate = () => {
  setEditItem(null)
  setFormOpen(true)
}

const handleEdit = (item: Team) => {
  setEditItem(item)
  setFormOpen(true)
}

<TeamForm
  open={formOpen}
  onClose={() => {
    setFormOpen(false)
    setEditItem(null)
  }}
  team={editItem}
/>
```

**Tabs Pattern** (from ContestsPage.tsx):
```typescript
const [tabValue, setTabValue] = useState('1')

<Tabs
  activeKey={tabValue}
  onChange={setTabValue}
  items={[
    { key: '1', label: 'Tab 1', children: <Component1 /> },
    { key: '2', label: 'Tab 2', children: <Component2 /> },
  ]}
/>
```

**Table Actions Pattern** (from TeamList.tsx):
```typescript
{
  title: 'Actions',
  key: 'actions',
  width: 150,
  render: (_, record) => (
    <Space>
      <Tooltip title="View">
        <Button icon={<EyeOutlined />} size="small" onClick={() => handleView(record)} />
      </Tooltip>
      <Tooltip title="Edit">
        <Button type="primary" icon={<EditOutlined />} size="small" onClick={() => handleEdit(record)} />
      </Tooltip>
      <Tooltip title="Delete">
        <Button danger icon={<DeleteOutlined />} size="small" onClick={() => handleDelete(record)} />
      </Tooltip>
    </Space>
  ),
}
```

**Date Formatting** (from existing utils):
```typescript
import { formatRelativeTime } from '../../utils/date-utils'

// Usage
formatRelativeTime(team.createdAt) // "2 hours ago"
```

---

## IMPLEMENTATION PLAN

### Phase 1: Enhance TeamsPage UX

Improve the TeamsPage with better layout, proper modal integration, and enhanced user experience.

**Tasks:**
- Refactor TeamsPage to use consistent modal patterns
- Add proper team details modal with members and invite code
- Improve tab organization and content
- Add empty states for better UX

### Phase 2: Integrate Team Leaderboards in Contests

Add team leaderboard functionality to the ContestsPage so users can view team rankings alongside individual rankings.

**Tasks:**
- Add "Team Leaderboard" tab to ContestsPage
- Integrate TeamLeaderboard component
- Add team participation indicator
- Handle team vs individual contest modes

### Phase 3: Enhance Team Components

Improve existing team components with better error handling, loading states, and user feedback.

**Tasks:**
- Add confirmation modals for destructive actions
- Improve empty states in TeamList and TeamMembers
- Add loading skeletons for better perceived performance
- Enhance TeamInvite component integration

### Phase 4: Comprehensive E2E Testing

Expand E2E test coverage to validate all team workflows and edge cases.

**Tasks:**
- Test team creation and editing
- Test joining teams with invite codes
- Test member management (add/remove)
- Test team leaderboard display
- Test error scenarios and validation

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### Task 1: UPDATE TeamsPage.tsx - Enhance UX and Modal Integration

- **IMPLEMENT**: Refactor to use proper modal for team details instead of inline display
- **PATTERN**: Mirror ContestsPage modal pattern (ContestsPage.tsx:20-50)
- **CHANGES**:
  - Add `viewTeamModalOpen` state separate from `viewTeam`
  - Create proper team details modal with tabs for members and invite code
  - Improve empty states for each tab
  - Add better error handling for join team flow
- **IMPORTS**: No new imports needed
- **GOTCHA**: TeamInvite requires `isCaptain` prop - calculate from `team.captainId === user?.id`
- **VALIDATE**: `npm run dev` and navigate to `/teams`, verify modals open/close properly

### Task 2: UPDATE TeamsPage.tsx - Add Team Details Modal

- **IMPLEMENT**: Create comprehensive team details modal with multiple sections
- **PATTERN**: Use Ant Design Modal with Tabs for sections (similar to ContestsPage participant modal)
- **CHANGES**:
  - Add modal showing team info, members, and invite code
  - Include TeamMembers component
  - Include TeamInvite component with proper props
  - Add "Leave Team" button for non-captains
  - Add "Delete Team" button for captains
- **IMPORTS**: Import `Modal`, `Tabs`, `Button`, `Descriptions` from antd
- **GOTCHA**: Check user role before showing captain-only actions
- **VALIDATE**: Click "View Members" button, verify modal shows all sections

### Task 3: UPDATE ContestsPage.tsx - Add Team Leaderboard Tab

- **IMPLEMENT**: Add team leaderboard tab to contest details view
- **PATTERN**: Mirror existing leaderboard tab structure (ContestsPage.tsx:80-120)
- **CHANGES**:
  - Add new tab item for "Team Leaderboard"
  - Import and use TeamLeaderboard component
  - Pass contestId and userTeamId props
  - Show only when contest supports team participation
- **IMPORTS**: `import TeamLeaderboard from '../components/teams/TeamLeaderboard'`
- **GOTCHA**: Need to fetch user's team for the contest to highlight in leaderboard
- **VALIDATE**: Navigate to contest details, verify team leaderboard tab appears

### Task 4: CREATE useUserTeams hook - Get User's Teams for Contest

- **IMPLEMENT**: Add hook to fetch teams user belongs to for a specific contest
- **PATTERN**: Follow existing hook patterns in use-teams.ts
- **FILE**: `frontend/src/hooks/use-teams.ts`
- **CHANGES**:
  - Add `useUserTeamsForContest(contestId: number)` hook
  - Use existing `useTeams` with `myTeamsOnly: true` filter
  - Return first team (users typically in one team per contest)
- **IMPORTS**: Already have all needed imports
- **GOTCHA**: Handle case where user is not in any team
- **VALIDATE**: `npm run build` to check TypeScript compilation

### Task 5: UPDATE TeamList.tsx - Add Confirmation Modals

- **IMPLEMENT**: Replace window.confirm with Ant Design Modal.confirm
- **PATTERN**: Use Modal.confirm for destructive actions (common Ant Design pattern)
- **CHANGES**:
  - Import `Modal` from antd
  - Replace `window.confirm` with `Modal.confirm`
  - Add proper title, content, and button labels
  - Show team details in confirmation message
- **IMPORTS**: `import { Modal } from 'antd'`
- **GOTCHA**: Modal.confirm is async, handle promise properly
- **VALIDATE**: Click delete button, verify modal appears with proper styling

### Task 6: UPDATE TeamList.tsx - Add Empty State

- **IMPLEMENT**: Add proper empty state when no teams exist
- **PATTERN**: Use Ant Design Empty component (similar to other lists)
- **CHANGES**:
  - Import `Empty` from antd
  - Add custom empty state with "Create Team" call-to-action
  - Show different message for "My Teams" vs "All Teams"
- **IMPORTS**: `import { Empty } from 'antd'`
- **GOTCHA**: Check if data is loaded before showing empty state
- **VALIDATE**: View page with no teams, verify empty state displays

### Task 7: UPDATE TeamMembers.tsx - Add Loading Skeleton

- **IMPLEMENT**: Add skeleton loading state for better UX
- **PATTERN**: Use Ant Design Skeleton component
- **CHANGES**:
  - Import `Skeleton` from antd
  - Replace simple Spin with List.Item.Meta skeleton
  - Show 3 skeleton items while loading
- **IMPORTS**: `import { Skeleton } from 'antd'`
- **GOTCHA**: Skeleton should match the actual content structure
- **VALIDATE**: Reload page, verify skeleton appears briefly

### Task 8: UPDATE TeamInvite.tsx - Fix Notification Imports

- **IMPLEMENT**: Ensure consistent notification usage
- **PATTERN**: Use showSuccess/showError from utils/notification.ts
- **CHANGES**:
  - Verify imports are from correct file
  - Ensure error handling is consistent
- **IMPORTS**: Check current imports match project pattern
- **GOTCHA**: Project uses custom notification utils, not Ant Design message
- **VALIDATE**: Copy invite code, verify notification appears

### Task 9: CREATE teams.spec.ts - Comprehensive E2E Tests

- **IMPLEMENT**: Expand E2E test coverage for all team workflows
- **PATTERN**: Follow existing test patterns from contests.spec.ts and auth.spec.ts
- **FILE**: `frontend/tests/e2e/teams.spec.ts`
- **CHANGES**:
  - Add test: "should create a new team"
  - Add test: "should edit team details"
  - Add test: "should delete team"
  - Add test: "should join team with invite code"
  - Add test: "should view team members"
  - Add test: "should leave team"
  - Add test: "should display team leaderboard in contest"
- **IMPORTS**: Use existing test fixtures and selectors
- **GOTCHA**: Tests need proper cleanup to avoid state pollution
- **VALIDATE**: `npm run test:e2e -- teams.spec.ts`

### Task 10: UPDATE selectors.ts - Add Missing Team Selectors

- **IMPLEMENT**: Add any missing selectors for new team functionality
- **PATTERN**: Follow existing selector structure
- **FILE**: `frontend/tests/helpers/selectors.ts`
- **CHANGES**:
  - Add selector for team details modal
  - Add selector for team leaderboard tab
  - Add selector for leave team button
  - Add selector for invite code input
- **IMPORTS**: None needed
- **GOTCHA**: Use data-testid attributes if needed for unique identification
- **VALIDATE**: Run E2E tests to verify selectors work

### Task 11: UPDATE TeamsPage.tsx - Add Data Test IDs

- **IMPLEMENT**: Add data-testid attributes for E2E testing
- **PATTERN**: Use kebab-case for test IDs
- **CHANGES**:
  - Add `data-testid="teams-page"` to main container
  - Add `data-testid="create-team-button"` to create button
  - Add `data-testid="team-details-modal"` to details modal
  - Add `data-testid="join-team-form"` to join form
- **IMPORTS**: None needed
- **GOTCHA**: Don't add test IDs to every element, only key interaction points
- **VALIDATE**: Inspect elements in browser dev tools

### Task 12: UPDATE ContestsPage.tsx - Add Team Leaderboard Data Test ID

- **IMPLEMENT**: Add test ID for team leaderboard tab
- **PATTERN**: Consistent with other test IDs
- **CHANGES**:
  - Add `data-testid="team-leaderboard-tab"` to tab item
  - Add `data-testid="team-leaderboard"` to TeamLeaderboard component wrapper
- **IMPORTS**: None needed
- **GOTCHA**: Ensure tab is visible before testing
- **VALIDATE**: Inspect contest details page

### Task 13: CREATE E2E Test - Team Creation Workflow

- **IMPLEMENT**: Test complete team creation flow
- **PATTERN**: Use authenticatedPage fixture from auth.fixture.ts
- **FILE**: `frontend/tests/e2e/teams.spec.ts`
- **TEST STEPS**:
  1. Navigate to /teams
  2. Click "Create Team" button
  3. Fill in team name, description, max members
  4. Submit form
  5. Verify team appears in list
  6. Verify success notification
- **IMPORTS**: `import { test, expect } from '../fixtures/auth.fixture'`
- **GOTCHA**: Use unique team names to avoid conflicts
- **VALIDATE**: `npm run test:e2e -- teams.spec.ts -g "should create"`

### Task 14: CREATE E2E Test - Join Team with Invite Code

- **IMPLEMENT**: Test joining team via invite code
- **PATTERN**: Multi-step workflow test
- **FILE**: `frontend/tests/e2e/teams.spec.ts`
- **TEST STEPS**:
  1. Create a team (setup)
  2. Copy invite code
  3. Switch to "Join Team" tab
  4. Enter invite code
  5. Submit
  6. Verify team appears in "My Teams"
- **IMPORTS**: Already imported
- **GOTCHA**: Need to handle async clipboard operations
- **VALIDATE**: `npm run test:e2e -- teams.spec.ts -g "should join"`

### Task 15: CREATE E2E Test - Team Member Management

- **IMPLEMENT**: Test viewing and managing team members
- **PATTERN**: Test with captain permissions
- **FILE**: `frontend/tests/e2e/teams.spec.ts`
- **TEST STEPS**:
  1. Create team (user becomes captain)
  2. Open team details modal
  3. Verify members list shows captain
  4. Verify invite code is displayed
  5. Test regenerate invite code (captain only)
- **IMPORTS**: Already imported
- **GOTCHA**: Captain-only actions should not be visible to regular members
- **VALIDATE**: `npm run test:e2e -- teams.spec.ts -g "member management"`

### Task 16: CREATE E2E Test - Team Leaderboard Display

- **IMPLEMENT**: Test team leaderboard in contest context
- **PATTERN**: Integration test across pages
- **FILE**: `frontend/tests/e2e/teams.spec.ts`
- **TEST STEPS**:
  1. Navigate to contests page
  2. Select a contest
  3. Click "Team Leaderboard" tab
  4. Verify leaderboard displays
  5. Verify columns (rank, team name, members, points)
- **IMPORTS**: Already imported
- **GOTCHA**: Contest must support team participation
- **VALIDATE**: `npm run test:e2e -- teams.spec.ts -g "leaderboard"`

### Task 17: UPDATE teams.spec.ts - Add Error Scenario Tests

- **IMPLEMENT**: Test error handling and validation
- **PATTERN**: Negative test cases
- **FILE**: `frontend/tests/e2e/teams.spec.ts`
- **TEST CASES**:
  - Invalid invite code shows error
  - Empty team name shows validation error
  - Max members below minimum shows error
  - Cannot delete team with active members (if applicable)
- **IMPORTS**: Already imported
- **GOTCHA**: Error messages should be user-friendly
- **VALIDATE**: `npm run test:e2e -- teams.spec.ts -g "error"`

### Task 18: UPDATE package.json - Add Team Test Script

- **IMPLEMENT**: Add convenience script for running team tests
- **PATTERN**: Follow existing test script patterns
- **FILE**: `frontend/package.json`
- **CHANGES**:
  - Add `"test:e2e:teams": "playwright test teams.spec.ts"` to scripts
- **IMPORTS**: None needed
- **GOTCHA**: Ensure script name doesn't conflict with existing scripts
- **VALIDATE**: `npm run test:e2e:teams`

### Task 19: VALIDATE - Run Full E2E Test Suite

- **IMPLEMENT**: Execute all E2E tests to ensure no regressions
- **PATTERN**: Full test suite validation
- **COMMAND**: `npm run test:e2e`
- **EXPECTED**: All tests pass, including new team tests
- **GOTCHA**: May need to increase timeout for slower tests
- **VALIDATE**: Check test report for any failures

### Task 20: VALIDATE - Manual Testing Checklist

- **IMPLEMENT**: Perform manual testing of all team features
- **PATTERN**: User acceptance testing
- **CHECKLIST**:
  - [ ] Create team with valid data
  - [ ] Edit team details
  - [ ] Delete team (with confirmation)
  - [ ] Join team with invite code
  - [ ] View team members
  - [ ] Leave team (non-captain)
  - [ ] Regenerate invite code (captain)
  - [ ] View team leaderboard in contest
  - [ ] All loading states display correctly
  - [ ] All error messages are clear
  - [ ] All success notifications appear
- **GOTCHA**: Test with different user roles (captain vs member)
- **VALIDATE**: Complete checklist with no issues

---

## TESTING STRATEGY

### Unit Tests

**Scope**: Not required for this enhancement (components already have implicit coverage through E2E tests)

**Rationale**: The project focuses on E2E testing with Playwright. Unit tests would be redundant given comprehensive E2E coverage.

### Integration Tests

**Scope**: E2E tests serve as integration tests, validating full user workflows

**Coverage Requirements**:
- Team CRUD operations (create, read, update, delete)
- Team membership management (join, leave, remove members)
- Invite code functionality (display, copy, regenerate)
- Team leaderboard display in contests
- Error handling and validation
- Cross-page navigation and state management

### Edge Cases

**Specific edge cases to test**:
1. **Empty States**: No teams exist, no members in team
2. **Permission Checks**: Captain vs member actions
3. **Invalid Input**: Malformed invite codes, invalid team names
4. **Concurrent Actions**: Multiple users joining same team
5. **Network Errors**: Failed API calls, timeout scenarios
6. **State Consistency**: Team list updates after create/delete
7. **Navigation**: Deep linking to teams page with query params

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# TypeScript compilation
cd frontend && npm run build

# ESLint checks
cd frontend && npm run lint

# Fix auto-fixable issues
cd frontend && npm run lint:fix
```

**Expected**: Zero TypeScript errors, zero ESLint errors

### Level 2: Development Server

```bash
# Start frontend dev server
cd frontend && npm run dev

# Verify teams page loads
# Open http://localhost:3000/teams in browser
```

**Expected**: Page loads without console errors, all components render

### Level 3: E2E Tests - Teams Only

```bash
# Run team-specific E2E tests
cd frontend && npm run test:e2e -- teams.spec.ts

# Run in headed mode for debugging
cd frontend && npm run test:e2e:headed -- teams.spec.ts

# Run in UI mode for interactive debugging
cd frontend && npm run test:e2e:ui -- teams.spec.ts
```

**Expected**: All team tests pass on all browsers (Chromium, Firefox, WebKit)

### Level 4: Full E2E Test Suite

```bash
# Run complete E2E test suite
cd frontend && npm run test:e2e

# Generate and view report
cd frontend && npm run test:e2e:report
```

**Expected**: All tests pass, no regressions in other features

### Level 5: Manual Validation

**Test with running backend services**:

```bash
# Start all services
make docker-services

# Verify API Gateway is running
curl http://localhost:8080/health

# Verify team endpoints are accessible
curl http://localhost:8080/v1/teams/health
```

**Manual Testing Steps**:
1. Register new user or login
2. Navigate to Teams page
3. Create a new team
4. Copy invite code
5. Open incognito window, login as different user
6. Join team using invite code
7. Return to first user, view team members
8. Navigate to Contests page
9. View team leaderboard tab
10. Verify all data displays correctly

**Expected**: All operations complete successfully, no errors in browser console

---

## ACCEPTANCE CRITERIA

- [x] Backend Team Service fully implemented and operational
- [x] API Gateway routes team requests to contest-service
- [x] Frontend types and service client complete
- [x] React Query hooks for all team operations complete
- [ ] TeamsPage provides full CRUD functionality with good UX
- [ ] Team details modal shows members and invite code
- [ ] Team leaderboard integrated into ContestsPage
- [ ] All team components have proper loading states
- [ ] All team components have proper error handling
- [ ] Confirmation modals for destructive actions
- [ ] Empty states for better UX
- [ ] E2E tests cover all team workflows (minimum 8 tests)
- [ ] E2E tests include error scenarios
- [ ] All E2E tests pass on all browsers
- [ ] No regressions in existing features
- [ ] Manual testing checklist completed
- [ ] Code follows project conventions (Ant Design, React Query patterns)
- [ ] TypeScript compilation succeeds with no errors
- [ ] ESLint passes with no warnings

---

## COMPLETION CHECKLIST

- [ ] Task 1: TeamsPage UX enhancement completed
- [ ] Task 2: Team details modal implemented
- [ ] Task 3: Team leaderboard tab added to ContestsPage
- [ ] Task 4: useUserTeams hook created
- [ ] Task 5: Confirmation modals added to TeamList
- [ ] Task 6: Empty states added to TeamList
- [ ] Task 7: Loading skeletons added to TeamMembers
- [ ] Task 8: Notification imports verified
- [ ] Task 9: E2E test file structure created
- [ ] Task 10: Test selectors updated
- [ ] Task 11: Data test IDs added to TeamsPage
- [ ] Task 12: Data test IDs added to ContestsPage
- [ ] Task 13: Team creation E2E test implemented
- [ ] Task 14: Join team E2E test implemented
- [ ] Task 15: Member management E2E test implemented
- [ ] Task 16: Team leaderboard E2E test implemented
- [ ] Task 17: Error scenario tests implemented
- [ ] Task 18: Test script added to package.json
- [ ] Task 19: Full E2E suite passes
- [ ] Task 20: Manual testing completed
- [ ] All validation commands executed successfully
- [ ] All acceptance criteria met

---

## NOTES

### Design Decisions

**Why enhance TeamsPage instead of rewrite?**
- Existing structure is sound, just needs better UX
- Components are already functional
- Incremental improvement reduces risk

**Why integrate team leaderboard into ContestsPage?**
- Teams compete in contests, so leaderboard belongs in contest context
- Consistent with individual leaderboard placement
- Better user experience (one place to view all rankings)

**Why focus on E2E tests over unit tests?**
- Project already uses Playwright for E2E testing
- E2E tests provide more value for UI components
- Backend has unit tests, frontend focuses on integration

### Trade-offs

**Comprehensive vs Minimal**:
- Chose comprehensive E2E coverage over minimal tests
- Rationale: Teams feature is critical for social engagement
- Risk: Longer test execution time (acceptable for quality)

**Modal vs Inline Display**:
- Chose modal for team details over inline expansion
- Rationale: Consistent with ContestsPage pattern, cleaner UX
- Trade-off: One extra click, but better information hierarchy

### Implementation Risks

**Low Risk**:
- Backend is complete and tested
- API Gateway integration verified
- Service and hooks already functional

**Medium Risk**:
- E2E test flakiness (mitigated by proper waits and selectors)
- Cross-browser compatibility (mitigated by testing on all browsers)

**Mitigation Strategies**:
- Use data-testid for stable selectors
- Add proper loading state waits in tests
- Test on all three browsers (Chromium, Firefox, WebKit)
- Manual testing before considering complete

### Future Enhancements (Out of Scope)

- Team chat or activity feed
- Team achievements and badges
- Team vs team challenges
- Team statistics and analytics
- Team profile customization (avatar, banner)
- Team-based notifications

These can be added in future iterations after core functionality is validated.

---

## CONFIDENCE SCORE: 9/10

**Rationale for High Confidence**:
- Backend fully implemented and operational ✅
- Frontend infrastructure (types, services, hooks) complete ✅
- Components exist and are functional ✅
- Clear patterns established in codebase ✅
- Comprehensive test strategy defined ✅
- All dependencies satisfied ✅

**Risk Factors** (-1 point):
- E2E tests may require iteration to handle async operations properly
- Team leaderboard integration needs careful state management
- Cross-browser testing may reveal minor UI inconsistencies

**Success Probability**: Very High - This is primarily integration work with clear patterns to follow. The heavy lifting (backend, types, services) is already done.

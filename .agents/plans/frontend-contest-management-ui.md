# Feature: Frontend User Interface for Contest Management

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Create a comprehensive React TypeScript frontend interface for contest management that allows users to view, create, edit, and manage sports prediction contests. The interface will provide full CRUD operations with real-time data synchronization, form validation, and responsive Material-UI design patterns.

## User Story

As a contest organizer
I want to manage sports prediction contests through an intuitive web interface
So that I can create, configure, and monitor contests without technical knowledge

## Problem Statement

The sports prediction platform currently has a complete backend microservices architecture with gRPC APIs, but lacks a user-friendly frontend interface. Users cannot interact with the contest system, create new contests, or manage existing ones without direct API calls.

## Solution Statement

Build a modern React TypeScript frontend with Material-UI components that provides:
- Contest listing with search, filtering, and pagination
- Contest creation and editing forms with validation
- Participant management interface
- Real-time data synchronization using React Query
- Responsive design for desktop and mobile
- Integration with existing gRPC backend services

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: Frontend application, API Gateway integration
**Dependencies**: React 18, Material-UI v5, React Hook Form, React Query, gRPC-Web

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/proto/contest.proto` (lines 1-150) - Why: Contest data structure and API definitions
- `backend/proto/common.proto` - Why: Common types like PaginationRequest/Response
- `backend/api-gateway/internal/gateway/gateway.go` (lines 1-100) - Why: HTTP endpoint patterns and error handling
- `backend/contest-service/internal/models/contest.go` (lines 1-80) - Why: Contest validation rules and field constraints
- `frontend/package.json` - Why: Existing dependencies and project configuration
- `docker-compose.yml` (lines 40-70) - Why: API Gateway endpoint configuration

### New Files to Create

- `frontend/src/types/contest.types.ts` - TypeScript interfaces for contest data
- `frontend/src/services/grpc-client.ts` - gRPC-Web client configuration
- `frontend/src/services/contest-service.ts` - Contest service API wrapper
- `frontend/src/hooks/use-contests.ts` - React Query hooks for contest operations
- `frontend/src/components/contests/ContestList.tsx` - Contest listing component
- `frontend/src/components/contests/ContestForm.tsx` - Contest creation/editing form
- `frontend/src/components/contests/ContestCard.tsx` - Individual contest display
- `frontend/src/components/contests/ParticipantList.tsx` - Participant management
- `frontend/src/pages/ContestsPage.tsx` - Main contests page
- `frontend/src/utils/validation.ts` - Form validation schemas
- `frontend/src/utils/date-utils.ts` - Date formatting utilities
- `frontend/vite.config.ts` - Vite configuration for development
- `frontend/src/main.tsx` - Application entry point
- `frontend/src/App.tsx` - Main application component
- `frontend/index.html` - HTML template

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Material-UI Data Grid CRUD](https://www.material-react-table.com/docs/examples/editing-crud-inline-table)
  - Specific section: Inline table editing with TypeScript
  - Why: Pattern for contest management table with CRUD operations
- [React Hook Form with Material-UI](https://react-hook-form.com/get-started#IntegratingwithUIlibraries)
  - Specific section: Controller integration patterns
  - Why: Form validation and Material-UI component integration
- [React Query Best Practices](https://tanstack.com/query/latest/docs/react/guides/best-practices)
  - Specific section: Optimistic updates and error handling
  - Why: Data synchronization with backend services
- [gRPC-Web Documentation](https://github.com/grpc/grpc-web)
  - Specific section: TypeScript client generation
  - Why: Integration with existing gRPC backend services

### Patterns to Follow

**Component Structure Pattern:**
```typescript
// From Material-UI examples
interface ComponentProps {
  data: DataType[];
  onAction: (item: DataType) => void;
}

export const Component: React.FC<ComponentProps> = ({ data, onAction }) => {
  // Component implementation
};
```

**React Query Hook Pattern:**
```typescript
// From React Query documentation
export const useContests = () => {
  return useQuery({
    queryKey: ['contests'],
    queryFn: () => contestService.getContests(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
```

**Form Validation Pattern:**
```typescript
// From React Hook Form + Zod
const schema = z.object({
  title: z.string().min(1, 'Title is required').max(200),
  description: z.string().max(1000).optional(),
});

type FormData = z.infer<typeof schema>;
```

**Error Handling Pattern:**
```typescript
// From API Gateway error format
interface ErrorResponse {
  error: string;
  code: number;
  message: string;
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation Setup

Set up the core infrastructure for React application with TypeScript, routing, and build configuration.

**Tasks:**
- Configure Vite build system with TypeScript and React
- Set up React Router for navigation
- Configure Material-UI theme and global styles
- Create application entry points and basic structure

### Phase 2: gRPC Integration Layer

Establish communication layer between React frontend and gRPC backend services.

**Tasks:**
- Configure gRPC-Web client with TypeScript generation
- Create service wrappers for Contest and User APIs
- Implement error handling and response transformation
- Set up React Query client with proper configuration

### Phase 3: Core Contest Components

Build the main contest management interface components with Material-UI.

**Tasks:**
- Create contest data types and interfaces
- Implement contest listing with search and filtering
- Build contest creation and editing forms
- Add participant management interface

### Phase 4: Integration and Polish

Connect all components with data layer and add advanced features.

**Tasks:**
- Integrate React Query hooks with components
- Add form validation and error handling
- Implement responsive design patterns
- Add loading states and user feedback

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE frontend/vite.config.ts

- **IMPLEMENT**: Vite configuration with React plugin and development server
- **PATTERN**: Standard Vite React TypeScript configuration
- **IMPORTS**: `@vitejs/plugin-react`, `vite`
- **GOTCHA**: Configure proxy for API Gateway at localhost:8080
- **VALIDATE**: `cd frontend && npm run dev`

### CREATE frontend/index.html

- **IMPLEMENT**: HTML template with proper meta tags and root div
- **PATTERN**: Standard React HTML template
- **IMPORTS**: None (HTML file)
- **GOTCHA**: Include Material-UI font imports and viewport meta tag
- **VALIDATE**: File exists and has proper structure

### CREATE frontend/src/main.tsx

- **IMPLEMENT**: React application entry point with providers
- **PATTERN**: React 18 createRoot pattern
- **IMPORTS**: `react`, `react-dom/client`, `@tanstack/react-query`
- **GOTCHA**: Wrap with QueryClient provider and StrictMode
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/App.tsx

- **IMPLEMENT**: Main application component with routing and theme
- **PATTERN**: Material-UI ThemeProvider with React Router
- **IMPORTS**: `@mui/material`, `react-router-dom`
- **GOTCHA**: Include CssBaseline for consistent styling
- **VALIDATE**: Component renders without errors

### CREATE frontend/src/types/contest.types.ts

- **IMPLEMENT**: TypeScript interfaces matching backend proto definitions
- **PATTERN**: Mirror backend/proto/contest.proto structure
- **IMPORTS**: None (type definitions only)
- **GOTCHA**: Use camelCase for JavaScript, match proto field types exactly
- **VALIDATE**: TypeScript compilation passes

### CREATE frontend/src/services/grpc-client.ts

- **IMPLEMENT**: gRPC-Web client configuration and base setup
- **PATTERN**: Singleton client with error handling
- **IMPORTS**: `grpc-web`, `google-protobuf`
- **GOTCHA**: Configure transport for development (localhost:8080)
- **VALIDATE**: Client instantiates without errors

### CREATE frontend/src/services/contest-service.ts

- **IMPLEMENT**: Contest service API wrapper with TypeScript
- **PATTERN**: Promise-based service methods with error transformation
- **IMPORTS**: gRPC client, contest types
- **GOTCHA**: Transform gRPC responses to JavaScript objects
- **VALIDATE**: Service methods return proper types

### CREATE frontend/src/hooks/use-contests.ts

- **IMPLEMENT**: React Query hooks for contest CRUD operations
- **PATTERN**: Separate hooks for queries and mutations
- **IMPORTS**: `@tanstack/react-query`, contest service
- **GOTCHA**: Implement optimistic updates for better UX
- **VALIDATE**: Hooks return proper loading/error states

### CREATE frontend/src/utils/validation.ts

- **IMPLEMENT**: Zod schemas for form validation
- **PATTERN**: Mirror backend validation rules from contest.go
- **IMPORTS**: `zod`
- **GOTCHA**: Match backend constraints (title max 200, description max 1000)
- **VALIDATE**: Schemas validate test data correctly

### CREATE frontend/src/utils/date-utils.ts

- **IMPLEMENT**: Date formatting and validation utilities
- **PATTERN**: Use date-fns for consistent date handling
- **IMPORTS**: `date-fns`
- **GOTCHA**: Handle timezone conversion for UTC backend dates
- **VALIDATE**: Functions format dates correctly

### CREATE frontend/src/components/contests/ContestCard.tsx

- **IMPLEMENT**: Material-UI card component for individual contest display
- **PATTERN**: Card with header, content, and action buttons
- **IMPORTS**: `@mui/material`, contest types
- **GOTCHA**: Show contest status with appropriate colors
- **VALIDATE**: Component renders with mock data

### CREATE frontend/src/components/contests/ContestForm.tsx

- **IMPLEMENT**: React Hook Form with Material-UI for contest creation/editing
- **PATTERN**: Controller-based form with Zod validation
- **IMPORTS**: `react-hook-form`, `@hookform/resolvers/zod`, `@mui/material`
- **GOTCHA**: Handle date picker for start/end dates properly
- **VALIDATE**: Form validates and submits correctly

### CREATE frontend/src/components/contests/ContestList.tsx

- **IMPLEMENT**: Material React Table with CRUD operations
- **PATTERN**: Data table with inline editing and actions
- **IMPORTS**: `material-react-table`, contest hooks
- **GOTCHA**: Implement search, filtering, and pagination
- **VALIDATE**: Table displays data and handles interactions

### CREATE frontend/src/components/contests/ParticipantList.tsx

- **IMPLEMENT**: Participant management interface with Material-UI
- **PATTERN**: List component with add/remove actions
- **IMPORTS**: `@mui/material`, participant types
- **GOTCHA**: Handle participant status changes
- **VALIDATE**: Component manages participant list correctly

### CREATE frontend/src/pages/ContestsPage.tsx

- **IMPLEMENT**: Main page component combining all contest components
- **PATTERN**: Page layout with tabs or sections
- **IMPORTS**: All contest components, Material-UI layout
- **GOTCHA**: Handle loading and error states properly
- **VALIDATE**: Page renders and navigates correctly

### UPDATE frontend/package.json

- **IMPLEMENT**: Add missing dependencies for complete functionality
- **PATTERN**: Add to existing dependencies object
- **IMPORTS**: `zod`, `@hookform/resolvers`, `material-react-table`
- **GOTCHA**: Ensure version compatibility with existing packages
- **VALIDATE**: `cd frontend && npm install`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests using Vitest and React Testing Library following existing project patterns:

- Component rendering tests for all contest components
- Form validation tests for ContestForm
- Hook behavior tests for use-contests
- Service method tests for contest-service
- Utility function tests for validation and date utilities

### Integration Tests

- Full contest creation workflow from form to backend
- Contest listing with search and filtering
- Participant management operations
- Error handling scenarios

### Edge Cases

- Network connectivity issues and retry logic
- Invalid form data submission
- Concurrent contest modifications
- Large dataset pagination performance
- Mobile responsive behavior

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
cd frontend && npm run lint
cd frontend && npm run lint:fix
cd frontend && npx tsc --noEmit
```

### Level 2: Unit Tests

```bash
cd frontend && npm test
cd frontend && npm run test:ui
```

### Level 3: Integration Tests

```bash
cd frontend && npm run build
cd frontend && npm run preview
```

### Level 4: Manual Validation

```bash
# Start backend services
make docker-up
make docker-services

# Start frontend development server
cd frontend && npm run dev

# Test contest creation workflow
# 1. Navigate to http://localhost:3000
# 2. Create new contest with valid data
# 3. Verify contest appears in list
# 4. Edit contest and save changes
# 5. Test participant management
```

### Level 5: Additional Validation (Optional)

```bash
# Performance testing
cd frontend && npm run build && npm run preview
# Check bundle size and loading performance

# Accessibility testing
# Use browser dev tools to check ARIA compliance
```

---

## ACCEPTANCE CRITERIA

- [ ] Contest listing displays all contests with search and filtering
- [ ] Contest creation form validates input and submits successfully
- [ ] Contest editing allows modification of all editable fields
- [ ] Participant management allows adding/removing participants
- [ ] Real-time data synchronization with backend services
- [ ] Responsive design works on desktop and mobile devices
- [ ] Form validation matches backend validation rules
- [ ] Error handling provides user-friendly messages
- [ ] Loading states provide appropriate user feedback
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets requirements (80%+)
- [ ] Integration tests verify end-to-end workflows
- [ ] Code follows project conventions and patterns
- [ ] No regressions in existing functionality
- [ ] Performance meets requirements (< 3s initial load)

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
- [ ] Documentation updated (if applicable)
- [ ] Performance requirements met

---

## NOTES

### Design Decisions

- **Material React Table**: Chosen for comprehensive CRUD operations with minimal custom code
- **React Hook Form + Zod**: Provides type-safe validation matching backend constraints
- **React Query**: Enables optimistic updates and automatic cache management
- **gRPC-Web**: Direct integration with existing backend services without REST translation

### Performance Considerations

- Implement virtual scrolling for large contest lists
- Use React Query's stale-while-revalidate for better perceived performance
- Lazy load participant lists for contests with many participants
- Optimize bundle size with code splitting

### Security Considerations

- Validate all user input on frontend and backend
- Implement proper authentication state management
- Sanitize display data to prevent XSS attacks
- Use HTTPS in production environment

### Future Enhancements

- Real-time updates using WebSocket connections
- Offline support with service workers
- Advanced filtering and sorting options
- Bulk operations for contest management
- Export functionality for contest data

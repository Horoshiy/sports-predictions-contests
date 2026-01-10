# Feature: Frontend Authentication UI

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Complete frontend authentication system with login/register pages, authentication context, protected routes, and JWT token management. This feature bridges the gap between the existing backend authentication service and the frontend UI, enabling users to securely access the Sports Prediction Contests platform.

## User Story

As a user of the Sports Prediction Contests platform
I want to register an account and login securely
So that I can create contests, make predictions, and track my performance with personalized access

## Problem Statement

The platform currently has a fully functional backend authentication system with JWT tokens, user registration, and login endpoints, but lacks the frontend UI components and authentication flow. Users cannot currently:
- Register new accounts through the web interface
- Login to access protected features
- Maintain authenticated sessions
- Access user-specific functionality like creating contests
- Have their identity verified for predictions and scoring

## Solution Statement

Implement a comprehensive frontend authentication system using React Context for state management, Material-UI for consistent UI components, and secure JWT token handling. The solution includes login/register forms, protected route components, authentication context with automatic token refresh, and integration with the existing gRPC-Gateway backend.

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: Frontend React App, User Service Integration
**Dependencies**: Material-UI, React Hook Form, Zod, React Router, existing gRPC client

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `frontend/src/contexts/ToastContext.tsx` - Why: Pattern for React Context implementation with TypeScript
- `frontend/src/components/contests/ContestForm.tsx` (lines 1-50) - Why: Material-UI form patterns with React Hook Form and Zod validation
- `frontend/src/utils/validation.ts` - Why: Zod validation patterns and helper functions to mirror
- `frontend/src/services/grpc-client.ts` - Why: Existing API client patterns and token handling (localStorage)
- `frontend/src/App.tsx` - Why: Current routing structure and Material-UI theme setup
- `backend/proto/user.proto` - Why: API contract for authentication endpoints
- `backend/user-service/internal/service/auth_service.go` - Why: Backend authentication logic and response structure
- `frontend/src/types/contest.types.ts` - Why: TypeScript type definition patterns

### New Files to Create

- `frontend/src/contexts/AuthContext.tsx` - Authentication context provider and hooks
- `frontend/src/components/auth/LoginForm.tsx` - Login form component
- `frontend/src/components/auth/RegisterForm.tsx` - Registration form component
- `frontend/src/components/auth/ProtectedRoute.tsx` - Route protection wrapper
- `frontend/src/pages/LoginPage.tsx` - Login page layout
- `frontend/src/pages/RegisterPage.tsx` - Registration page layout
- `frontend/src/services/auth-service.ts` - Authentication API service
- `frontend/src/types/auth.types.ts` - Authentication TypeScript types
- `frontend/src/utils/auth-validation.ts` - Authentication-specific validation schemas
- `frontend/src/hooks/use-auth.ts` - Custom authentication hooks
- `tests/frontend/src/components/auth/LoginForm.test.tsx` - Login form tests
- `tests/frontend/src/contexts/AuthContext.test.tsx` - Auth context tests

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [React Hook Form Documentation](https://react-hook-form.com/get-started)
  - Specific section: Controller component for Material-UI integration
  - Why: Required for form handling patterns matching existing codebase
- [Material-UI TextField Documentation](https://mui.com/material-ui/react-text-field/)
  - Specific section: Validation and error states
  - Why: Consistent form styling and error handling
- [React Router v6 Protected Routes](https://reactrouter.com/en/main/examples/auth)
  - Specific section: Authentication example with Navigate
  - Why: Modern protected route implementation
- [Zod Validation Documentation](https://zod.dev/?id=basic-usage)
  - Specific section: Schema composition and error handling
  - Why: Validation patterns matching existing validation.ts

### Patterns to Follow

**Naming Conventions:**
- Components: PascalCase (LoginForm, RegisterForm)
- Files: kebab-case (auth-service.ts, use-auth.ts)
- Types: PascalCase interfaces (AuthUser, LoginFormData)
- Context: PascalCase with Context suffix (AuthContext)

**Error Handling:**
```typescript
// Pattern from grpc-client.ts
try {
  const response = await fetch(url, config)
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({
      error: 'Request failed',
      message: response.statusText,
    }))
    throw new Error(errorData.message || `HTTP ${response.status}`)
  }
  return await response.json()
} catch (error) {
  console.error('Request failed:', error)
  throw error
}
```

**Context Pattern (from ToastContext.tsx):**
```typescript
const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
```

**Form Validation Pattern (from ContestForm.tsx):**
```typescript
const { control, handleSubmit, formState: { errors } } = useForm<FormData>({
  resolver: zodResolver(validationSchema),
  defaultValues: { /* defaults */ }
})
```

**Material-UI Styling Pattern:**
- Use `sx` prop for component styling
- Follow theme spacing units
- Use theme palette colors
- Consistent component composition

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up authentication types, validation schemas, and API service layer to establish the data contracts and communication patterns.

**Tasks:**
- Create TypeScript interfaces for authentication data
- Implement Zod validation schemas for login/register forms
- Build authentication API service using existing gRPC client patterns
- Set up authentication context structure

### Phase 2: Core Implementation

Build the authentication context, forms, and core authentication logic with proper state management and error handling.

**Tasks:**
- Implement AuthContext with login/logout/register methods
- Create login and registration form components
- Add authentication state management with token handling
- Implement automatic token refresh logic

### Phase 3: Integration

Connect authentication to routing system, protect existing routes, and integrate with current application structure.

**Tasks:**
- Create ProtectedRoute component for route protection
- Add authentication pages to routing system
- Update existing components to use authentication context
- Integrate with existing toast notification system

### Phase 4: Testing & Validation

Comprehensive testing of authentication flows, edge cases, and integration with existing functionality.

**Tasks:**
- Unit tests for authentication components and context
- Integration tests for authentication flows
- Manual testing of complete user journeys
- Validation of security best practices

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE frontend/src/types/auth.types.ts

- **IMPLEMENT**: TypeScript interfaces for authentication data matching backend proto definitions
- **PATTERN**: Mirror type structure from `frontend/src/types/contest.types.ts`
- **IMPORTS**: No external imports needed
- **GOTCHA**: Ensure User interface excludes password field for security
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/utils/auth-validation.ts

- **IMPLEMENT**: Zod validation schemas for login and registration forms
- **PATTERN**: Follow validation patterns from `frontend/src/utils/validation.ts`
- **IMPORTS**: `import { z } from 'zod'`
- **GOTCHA**: Email validation must match backend requirements, password minimum 8 characters
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/services/auth-service.ts

- **IMPLEMENT**: Authentication API service with login, register, and profile methods
- **PATTERN**: Mirror service structure from `frontend/src/services/contest-service.ts`
- **IMPORTS**: Import grpcClient and auth types
- **GOTCHA**: Use exact API paths from backend proto: `/v1/auth/login`, `/v1/auth/register`
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/contexts/AuthContext.tsx

- **IMPLEMENT**: Authentication context with login, logout, register, and user state
- **PATTERN**: Follow context pattern from `frontend/src/contexts/ToastContext.tsx`
- **IMPORTS**: React context hooks, auth types, auth service, toast context
- **GOTCHA**: Handle token storage in localStorage, implement automatic logout on token expiry
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/hooks/use-auth.ts

- **IMPLEMENT**: Custom hook for accessing authentication context
- **PATTERN**: Simple context consumer hook like useToast pattern
- **IMPORTS**: AuthContext from contexts
- **GOTCHA**: Throw error if used outside AuthProvider
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/auth/LoginForm.tsx

- **IMPLEMENT**: Login form with email/password fields and validation
- **PATTERN**: Mirror form structure from `frontend/src/components/contests/ContestForm.tsx`
- **IMPORTS**: Material-UI components, React Hook Form, Zod resolver, auth validation
- **GOTCHA**: Use Controller for Material-UI integration, handle loading states
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/auth/RegisterForm.tsx

- **IMPLEMENT**: Registration form with email, password, confirm password, and name fields
- **PATTERN**: Mirror form structure from LoginForm.tsx with additional fields
- **IMPORTS**: Same as LoginForm plus additional validation for password confirmation
- **GOTCHA**: Password confirmation validation, proper error display for each field
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/pages/LoginPage.tsx

- **IMPLEMENT**: Login page layout with form, navigation links, and branding
- **PATTERN**: Follow page structure from `frontend/src/pages/ContestsPage.tsx`
- **IMPORTS**: Material-UI layout components, LoginForm, Link from react-router-dom
- **GOTCHA**: Center form on page, add link to registration, handle authentication redirects
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/pages/RegisterPage.tsx

- **IMPLEMENT**: Registration page layout with form and navigation
- **PATTERN**: Mirror LoginPage.tsx structure with RegisterForm
- **IMPORTS**: Material-UI layout components, RegisterForm, Link from react-router-dom
- **GOTCHA**: Add link to login page, consistent styling with LoginPage
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/auth/ProtectedRoute.tsx

- **IMPLEMENT**: Route wrapper that redirects unauthenticated users to login
- **PATTERN**: Use React Router v6 Navigate component for redirects
- **IMPORTS**: React, Navigate from react-router-dom, useAuth hook
- **GOTCHA**: Preserve location state for redirect after login, handle loading states
- **VALIDATE**: `cd frontend && npm run build`

### UPDATE frontend/src/App.tsx

- **IMPLEMENT**: Add AuthProvider wrapper and authentication routes
- **PATTERN**: Wrap existing providers, add routes for /login and /register
- **IMPORTS**: AuthProvider, LoginPage, RegisterPage, ProtectedRoute
- **GOTCHA**: AuthProvider must wrap Router, protect existing routes with ProtectedRoute
- **VALIDATE**: `cd frontend && npm run build`

### UPDATE frontend/src/App.tsx

- **IMPLEMENT**: Add logout button to AppBar and user display
- **PATTERN**: Add user menu to existing AppBar in header
- **IMPORTS**: Material-UI Menu components, useAuth hook
- **GOTCHA**: Show user name when logged in, handle logout action
- **VALIDATE**: `cd frontend && npm run build`

### UPDATE frontend/src/pages/ContestsPage.tsx

- **IMPLEMENT**: Replace hardcoded currentUserId with authenticated user ID
- **PATTERN**: Use useAuth hook to get current user
- **IMPORTS**: useAuth hook
- **GOTCHA**: Remove TODO comment, handle case when user is not loaded yet
- **VALIDATE**: `cd frontend && npm run build`

### CREATE tests/frontend/src/components/auth/LoginForm.test.tsx

- **IMPLEMENT**: Unit tests for LoginForm component
- **PATTERN**: Follow testing patterns from existing test files
- **IMPORTS**: React Testing Library, Jest, LoginForm component
- **GOTCHA**: Mock auth context and form submission
- **VALIDATE**: `cd frontend && npm test`

### CREATE tests/frontend/src/contexts/AuthContext.test.tsx

- **IMPLEMENT**: Unit tests for AuthContext functionality
- **PATTERN**: Test context provider and consumer patterns
- **IMPORTS**: React Testing Library, AuthContext, mock auth service
- **GOTCHA**: Test login, logout, and registration flows
- **VALIDATE**: `cd frontend && npm test`

---

## TESTING STRATEGY

### Unit Tests

**Scope**: Individual components and context functionality
- LoginForm component with form validation and submission
- RegisterForm component with password confirmation
- AuthContext with login/logout/register methods
- ProtectedRoute component with redirect logic
- Auth service API calls with mocked responses

**Requirements**: 
- Test form validation errors
- Test successful and failed authentication flows
- Test token storage and retrieval
- Test automatic logout on token expiry
- Mock external dependencies (API calls, localStorage)

### Integration Tests

**Scope**: Complete authentication flows and route protection
- Full login flow from form submission to authenticated state
- Registration flow with automatic login
- Protected route access and redirects
- Token refresh and session management
- Integration with existing toast notifications

### Edge Cases

**Specific edge cases that must be tested:**
- Invalid credentials handling
- Network errors during authentication
- Expired token scenarios
- Malformed API responses
- Concurrent login attempts
- Browser refresh with stored tokens
- Logout from multiple tabs

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
cd frontend && npm run lint
cd frontend && npm run build
cd frontend && npx tsc --noEmit
```

### Level 2: Unit Tests

```bash
cd frontend && npm test -- --coverage --watchAll=false
cd frontend && npm test -- auth --watchAll=false
```

### Level 3: Integration Tests

```bash
cd frontend && npm run test:integration
cd backend && go test ./user-service/... -v
```

### Level 4: Manual Validation

**Authentication Flow Testing:**
```bash
# Start development environment
make dev
cd frontend && npm run dev

# Test in browser:
# 1. Navigate to http://localhost:3000/login
# 2. Test registration with valid data
# 3. Test login with created account
# 4. Verify protected route access
# 5. Test logout functionality
# 6. Verify redirect after logout
```

**API Integration Testing:**
```bash
# Test backend authentication endpoints
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Level 5: Additional Validation (Optional)

```bash
# Security validation
cd frontend && npm audit
cd frontend && npm run build && npx serve -s build

# Performance validation
cd frontend && npm run build
# Test bundle size and loading performance
```

---

## ACCEPTANCE CRITERIA

- [ ] Users can register new accounts with email, password, and name
- [ ] Users can login with valid credentials and receive JWT tokens
- [ ] Authentication state persists across browser sessions
- [ ] Protected routes redirect unauthenticated users to login
- [ ] Users can logout and clear authentication state
- [ ] Form validation provides clear error messages
- [ ] Authentication errors are displayed with user-friendly messages
- [ ] UI follows existing Material-UI design patterns
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ requirement
- [ ] Integration tests verify end-to-end authentication flows
- [ ] No regressions in existing contest functionality
- [ ] Security best practices implemented (no password storage, secure token handling)
- [ ] Performance meets requirements (fast form interactions, minimal bundle impact)

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms authentication works end-to-end
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability
- [ ] Security review completed (no sensitive data exposure)
- [ ] Performance impact assessed and acceptable

---

## NOTES

**Security Considerations:**
- JWT tokens stored in localStorage (existing pattern) - consider HttpOnly cookies for production
- Password validation enforces minimum 8 characters
- No password storage in frontend state or localStorage
- Automatic logout on token expiry prevents stale sessions

**Design Decisions:**
- Following existing Material-UI patterns for consistency
- Using React Hook Form + Zod for validation (matches existing forms)
- Context-based state management for simplicity and existing patterns
- Protected routes using React Router v6 Navigate component

**Future Enhancements:**
- Remember me functionality with refresh tokens
- Social login integration (Google, Facebook)
- Two-factor authentication
- Password reset functionality
- Account verification via email

**Performance Considerations:**
- Lazy loading of authentication pages
- Minimal bundle size impact with tree shaking
- Efficient re-renders with proper context optimization
- Token refresh without full page reload

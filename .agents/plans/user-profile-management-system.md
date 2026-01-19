# Feature: User Profile Management System

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

A comprehensive user profile management system that extends the existing basic user authentication to include detailed profile information, avatar uploads, privacy settings, profile completion tracking, and user preferences. This system will provide users with full control over their profile data while maintaining security and privacy standards.

## User Story

As a registered user
I want to manage my detailed profile information including avatar, bio, preferences, and privacy settings
So that I can personalize my experience, control my privacy, and track my engagement with the platform

## Problem Statement

The current user system only supports basic authentication with minimal profile data (name, email). Users need:
- Detailed profile information (bio, location, social links)
- Avatar/profile picture management
- Privacy controls for profile visibility
- Profile completion tracking for gamification
- User preferences for notifications and UI customization
- Secure and intuitive profile management interface

## Solution Statement

Extend the existing user service and frontend to include comprehensive profile management with:
- Enhanced user model with additional profile fields
- Secure avatar upload system with image processing
- Granular privacy settings with field-level controls
- Profile completion tracking with progress indicators
- User preferences management with hierarchical settings
- Modern React UI with Material-UI components following existing patterns

## Feature Metadata

**Feature Type**: Enhancement
**Estimated Complexity**: Medium
**Primary Systems Affected**: user-service, frontend, api-gateway
**Dependencies**: Image processing library, file storage system

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/user-service/internal/models/user.go` (lines 13-19) - Why: Current User model structure to extend
- `backend/user-service/internal/models/user.go` (lines 41-97) - Why: Validation patterns to follow for new fields
- `backend/proto/user.proto` (lines 1-75) - Why: gRPC service definition to extend
- `frontend/src/types/auth.types.ts` - Why: TypeScript interfaces to extend
- `frontend/src/contexts/AuthContext.tsx` (lines 1-40) - Why: Auth context pattern to extend
- `frontend/src/components/auth/RegisterForm.tsx` - Why: Form handling patterns with React Hook Form + Zod
- `frontend/src/App.tsx` (lines 50-130) - Why: Navigation structure and user menu integration
- `backend/user-service/internal/service/user_service.go` - Why: Service layer patterns to follow
- `backend/user-service/internal/repository/user_repository.go` - Why: Repository patterns for data access

### New Files to Create

- `backend/user-service/internal/models/profile.go` - Extended profile model with validation
- `backend/user-service/internal/models/user_preferences.go` - User preferences model
- `backend/user-service/internal/service/profile_service.go` - Profile management service
- `backend/user-service/internal/handlers/upload_handler.go` - Avatar upload handler
- `backend/proto/profile.proto` - Profile management gRPC definitions
- `frontend/src/pages/ProfilePage.tsx` - Main profile management page
- `frontend/src/components/profile/ProfileForm.tsx` - Profile editing form
- `frontend/src/components/profile/AvatarUpload.tsx` - Avatar upload component
- `frontend/src/components/profile/PrivacySettings.tsx` - Privacy controls component
- `frontend/src/components/profile/ProfileCompletion.tsx` - Progress tracking component
- `frontend/src/types/profile.types.ts` - Profile-related TypeScript types
- `frontend/src/services/profile-service.ts` - Profile API service client

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Go GORM Documentation](https://gorm.io/docs/models.html#Embedded-Struct)
  - Specific section: Model definition and validation
  - Why: Required for extending User model with profile fields
- [Material-UI File Upload](https://mui.com/material-ui/react-button/#file-upload)
  - Specific section: File input components
  - Why: Shows proper file upload UI patterns
- [React Hook Form Documentation](https://react-hook-form.com/docs/useform)
  - Specific section: Form validation and submission
  - Why: Required for profile form implementation
- [gRPC HTTP Gateway](https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/)
  - Specific section: HTTP annotations for file uploads
  - Why: Shows how to handle multipart uploads in gRPC

### Patterns to Follow

**Naming Conventions:**
```go
// Model fields: PascalCase with GORM tags
type Profile struct {
    Bio        string `gorm:"type:text" json:"bio"`
    AvatarURL  string `gorm:"size:500" json:"avatar_url"`
    Location   string `gorm:"size:100" json:"location"`
}
```

**Error Handling:**
```go
// Follow existing validation pattern
func (p *Profile) ValidateBio() error {
    if len(p.Bio) > 500 {
        return errors.New("bio must be less than 500 characters")
    }
    return nil
}
```

**Logging Pattern:**
```go
// Use structured logging like existing services
log.WithFields(log.Fields{
    "user_id": userID,
    "action": "update_profile",
}).Info("Profile updated successfully")
```

**Frontend Form Pattern:**
```typescript
// Follow RegisterForm.tsx pattern with React Hook Form + Zod
const schema = z.object({
  bio: z.string().max(500, "Bio must be less than 500 characters"),
  location: z.string().max(100, "Location must be less than 100 characters"),
})
```

---

## IMPLEMENTATION PLAN

### Phase 1: Backend Foundation

Extend the user service with profile management capabilities, following existing patterns for model definition, validation, and service layer implementation.

**Tasks:**
- Extend User model with profile fields (bio, avatar_url, location, etc.)
- Create UserPreferences model for settings management
- Add validation methods following existing patterns
- Update database migrations

### Phase 2: gRPC Service Extension

Add profile management endpoints to the existing user service, maintaining consistency with current API patterns.

**Tasks:**
- Extend user.proto with profile management messages
- Add profile service methods to UserService
- Implement avatar upload endpoint with proper validation
- Add privacy settings and preferences endpoints

### Phase 3: Service Layer Implementation

Implement business logic for profile management, following the existing service and repository patterns.

**Tasks:**
- Create ProfileService with CRUD operations
- Implement avatar upload handling with security validation
- Add profile completion calculation logic
- Create preferences management service

### Phase 4: Frontend Integration

Build React components for profile management using Material-UI and existing form patterns.

**Tasks:**
- Create ProfilePage with tabbed interface
- Implement profile editing forms with validation
- Add avatar upload component with preview
- Create privacy settings interface
- Add profile completion progress indicator

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE backend/user-service/internal/models/profile.go

- **IMPLEMENT**: Extended profile fields with GORM tags and JSON serialization
- **PATTERN**: Mirror User struct pattern from `backend/user-service/internal/models/user.go:13-19`
- **IMPORTS**: `gorm.io/gorm`, `errors`, `strings`, `regexp`
- **GOTCHA**: Use proper GORM field types (text for bio, varchar for URLs)
- **VALIDATE**: `cd backend/user-service && go build ./...`

### UPDATE backend/user-service/internal/models/user.go

- **IMPLEMENT**: Add profile relationship and privacy settings fields
- **PATTERN**: Follow existing GORM relationship patterns
- **IMPORTS**: Add profile model import
- **GOTCHA**: Use proper GORM association tags for has-one relationship
- **VALIDATE**: `cd backend/user-service && go build ./...`

### CREATE backend/user-service/internal/models/user_preferences.go

- **IMPLEMENT**: User preferences model with JSONB storage
- **PATTERN**: Mirror User validation methods from `backend/user-service/internal/models/user.go:41-97`
- **IMPORTS**: `gorm.io/gorm`, `encoding/json`, `errors`
- **GOTCHA**: Use JSONB type for PostgreSQL, proper JSON tags
- **VALIDATE**: `cd backend/user-service && go build ./...`

### CREATE backend/proto/profile.proto

- **IMPLEMENT**: Profile management gRPC service definitions
- **PATTERN**: Mirror structure from `backend/proto/user.proto`
- **IMPORTS**: `google/protobuf/timestamp.proto`, `google/api/annotations.proto`, `common.proto`
- **GOTCHA**: Use proper HTTP annotations for file uploads, consistent naming
- **VALIDATE**: `make proto`

### UPDATE backend/proto/user.proto

- **IMPLEMENT**: Add profile fields to User message and new profile endpoints
- **PATTERN**: Follow existing message structure and HTTP annotations
- **IMPORTS**: Import profile.proto if needed
- **GOTCHA**: Maintain backward compatibility, use optional fields
- **VALIDATE**: `make proto`

### CREATE backend/user-service/internal/service/profile_service.go

- **IMPLEMENT**: Profile CRUD operations and business logic
- **PATTERN**: Mirror UserService structure from `backend/user-service/internal/service/user_service.go`
- **IMPORTS**: Models, repository interfaces, context, gRPC status codes
- **GOTCHA**: Use context for user identification, proper error handling
- **VALIDATE**: `cd backend/user-service && go build ./...`

### CREATE backend/user-service/internal/handlers/upload_handler.go

- **IMPLEMENT**: Secure avatar upload with validation and processing
- **PATTERN**: Follow existing handler patterns in user service
- **IMPORTS**: `mime/multipart`, `image`, `image/jpeg`, `image/png`, `path/filepath`
- **GOTCHA**: Validate file types, size limits, generate secure filenames
- **VALIDATE**: `cd backend/user-service && go build ./...`

### UPDATE backend/user-service/internal/repository/user_repository.go

- **IMPLEMENT**: Add profile-related database operations
- **PATTERN**: Follow existing repository methods and error handling
- **IMPORTS**: Add profile models
- **GOTCHA**: Use proper GORM preloading for profile associations
- **VALIDATE**: `cd backend/user-service && go build ./...`

### UPDATE backend/user-service/cmd/main.go

- **IMPLEMENT**: Register new profile service handlers
- **PATTERN**: Follow existing service registration pattern
- **IMPORTS**: Add profile service imports
- **GOTCHA**: Maintain existing initialization order
- **VALIDATE**: `cd backend/user-service && go run cmd/main.go --help`

### CREATE frontend/src/types/profile.types.ts

- **IMPLEMENT**: TypeScript interfaces for profile management
- **PATTERN**: Mirror structure from `frontend/src/types/auth.types.ts`
- **IMPORTS**: Common types, auth types
- **GOTCHA**: Match backend proto definitions exactly
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/services/profile-service.ts

- **IMPLEMENT**: Profile API client with gRPC-Web integration
- **PATTERN**: Mirror AuthService class from `frontend/src/services/auth-service.ts`
- **IMPORTS**: gRPC client, profile types, error handling
- **GOTCHA**: Handle file uploads properly, use FormData for avatar
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/profile/ProfileForm.tsx

- **IMPLEMENT**: Profile editing form with React Hook Form + Zod validation
- **PATTERN**: Mirror RegisterForm structure from `frontend/src/components/auth/RegisterForm.tsx`
- **IMPORTS**: React Hook Form, Zod, Material-UI components
- **GOTCHA**: Use proper validation rules, handle async updates
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/profile/AvatarUpload.tsx

- **IMPLEMENT**: Avatar upload component with preview and validation
- **PATTERN**: Follow Material-UI file upload patterns
- **IMPORTS**: React, Material-UI, file validation utilities
- **GOTCHA**: Validate file types client-side, show upload progress
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/profile/PrivacySettings.tsx

- **IMPLEMENT**: Privacy controls with toggle switches and explanations
- **PATTERN**: Use Material-UI Switch and FormControl components
- **IMPORTS**: Material-UI components, profile types
- **GOTCHA**: Clear privacy level explanations, immediate save on change
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/components/profile/ProfileCompletion.tsx

- **IMPLEMENT**: Progress indicator with completion percentage and suggestions
- **PATTERN**: Use Material-UI LinearProgress and List components
- **IMPORTS**: Material-UI components, profile completion logic
- **GOTCHA**: Calculate completion percentage accurately, show actionable suggestions
- **VALIDATE**: `cd frontend && npm run build`

### CREATE frontend/src/pages/ProfilePage.tsx

- **IMPLEMENT**: Main profile page with tabbed interface
- **PATTERN**: Mirror page structure from `frontend/src/pages/AnalyticsPage.tsx`
- **IMPORTS**: React, Material-UI Tabs, profile components
- **GOTCHA**: Use proper tab navigation, handle loading states
- **VALIDATE**: `cd frontend && npm run build`

### UPDATE frontend/src/App.tsx

- **IMPLEMENT**: Add profile route and update navigation menu
- **PATTERN**: Follow existing route structure and menu item pattern from lines 50-130
- **IMPORTS**: Add ProfilePage import
- **GOTCHA**: Make profile menu item functional, add proper route protection
- **VALIDATE**: `cd frontend && npm run dev` (check navigation works)

### UPDATE frontend/src/contexts/AuthContext.tsx

- **IMPLEMENT**: Add profile management methods to auth context
- **PATTERN**: Follow existing context method patterns
- **IMPORTS**: Add profile service and types
- **GOTCHA**: Update user state when profile changes, handle errors properly
- **VALIDATE**: `cd frontend && npm run build`

### CREATE tests/backend/user-service/profile_test.go

- **IMPLEMENT**: Unit tests for profile service operations
- **PATTERN**: Follow existing test patterns in user service tests
- **IMPORTS**: Testing framework, testify, profile models
- **GOTCHA**: Test validation rules, error cases, file upload scenarios
- **VALIDATE**: `cd tests/backend/user-service && go test -v`

### CREATE tests/frontend/components/ProfileForm.test.tsx

- **IMPLEMENT**: Component tests for profile form
- **PATTERN**: Follow existing frontend test patterns
- **IMPORTS**: React Testing Library, Jest, profile components
- **GOTCHA**: Test form validation, submission, error handling
- **VALIDATE**: `cd frontend && npm test`

---

## TESTING STRATEGY

### Unit Tests

**Backend Testing:**
- Profile model validation methods
- Profile service CRUD operations
- Avatar upload handler security validation
- Privacy settings logic
- Profile completion calculation

**Frontend Testing:**
- Profile form validation and submission
- Avatar upload component functionality
- Privacy settings toggle behavior
- Profile completion progress calculation

### Integration Tests

**API Integration:**
- Profile CRUD endpoints with authentication
- Avatar upload with file validation
- Privacy settings persistence
- Profile completion tracking

**Frontend Integration:**
- Profile page navigation and tab switching
- Form submission with API integration
- Avatar upload with progress indication
- Real-time profile completion updates

### Edge Cases

- Invalid file types and oversized uploads
- Concurrent profile updates
- Privacy setting conflicts
- Profile completion with partial data
- Network failures during avatar upload
- XSS prevention in bio and location fields

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Backend validation
cd backend/user-service && go fmt ./...
cd backend/user-service && go vet ./...
cd backend && go work sync

# Frontend validation  
cd frontend && npm run lint
cd frontend && npm run build
```

### Level 2: Unit Tests

```bash
# Backend unit tests
cd backend/user-service && go test ./... -v
cd tests/backend/user-service && go test -v

# Frontend unit tests
cd frontend && npm test -- --coverage
```

### Level 3: Integration Tests

```bash
# Protocol buffer generation
make proto

# Service integration tests
cd tests/backend && go test -v ./...

# Frontend integration tests
cd frontend && npm run test:integration
```

### Level 4: Manual Validation

```bash
# Start development environment
make dev

# Test profile endpoints
curl -X GET http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer <token>"

curl -X PUT http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"bio":"Test bio","location":"Test location"}'

# Test avatar upload
curl -X POST http://localhost:8080/v1/users/avatar \
  -H "Authorization: Bearer <token>" \
  -F "avatar=@test-image.jpg"
```

### Level 5: Additional Validation (Optional)

```bash
# E2E tests if available
make e2e-test

# Performance testing
cd tests && go test -bench=. -benchmem
```

---

## ACCEPTANCE CRITERIA

- [ ] Users can view and edit comprehensive profile information (bio, location, social links)
- [ ] Avatar upload system works with proper validation and security
- [ ] Privacy settings provide granular control over profile visibility
- [ ] Profile completion tracking shows accurate progress and suggestions
- [ ] User preferences are saved and applied across the application
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets 80%+ requirement
- [ ] Integration tests verify end-to-end profile workflows
- [ ] Code follows existing project conventions and patterns
- [ ] No regressions in existing authentication functionality
- [ ] Profile page integrates seamlessly with existing navigation
- [ ] Avatar uploads are secure and properly validated
- [ ] Privacy settings are enforced across all profile views
- [ ] Profile completion gamification encourages user engagement

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in dependency order
- [ ] Each task validation passed immediately after implementation
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms all profile features work
- [ ] Profile navigation integrated into main app
- [ ] Avatar upload security validated
- [ ] Privacy settings properly enforced
- [ ] Profile completion tracking accurate
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

**Security Considerations:**
- Avatar uploads must validate file types and sizes server-side
- Profile data should be sanitized to prevent XSS attacks
- Privacy settings must be enforced at the API level
- File storage should use secure, isolated domains

**Performance Considerations:**
- Avatar images should be resized and optimized automatically
- Profile completion should be calculated efficiently
- Privacy settings should be cached for performance
- Large profile updates should be handled asynchronously

**User Experience:**
- Profile completion should provide clear, actionable suggestions
- Avatar upload should show progress and preview
- Privacy settings should have clear explanations
- Form validation should provide immediate feedback

**Future Extensibility:**
- Profile system designed to support additional fields easily
- Privacy settings architecture supports new privacy levels
- Preferences system can accommodate new user settings
- Avatar system can support multiple image sizes and formats

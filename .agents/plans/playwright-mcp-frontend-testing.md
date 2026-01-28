# Feature: Playwright MCP Frontend Testing System

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement a comprehensive end-to-end testing system for the React frontend using Playwright with Model Context Protocol (MCP) integration. This system will enable AI-driven browser automation testing through the Playwright MCP server, providing visual testing capabilities, automated test generation, and comprehensive coverage of all user workflows.

The testing framework will cover authentication flows, contest management, predictions, team tournaments, analytics dashboards, and all interactive UI components built with Ant Design.

## User Story

As a **developer**
I want to **have automated E2E tests for the frontend using Playwright MCP**
So that **I can ensure UI functionality works correctly, catch regressions early, and leverage AI-assisted test generation and debugging**

## Problem Statement

Currently, the project has:
- Backend E2E tests in Go (tests/e2e/)
- Minimal frontend unit tests (1 file: frontend/src/tests/fixes.test.ts)
- No comprehensive frontend E2E testing
- No visual regression testing
- No automated browser interaction testing
- MCP server configured but not utilized for testing

This creates risks:
- UI regressions go undetected until manual testing
- Complex user workflows (auth → contest → prediction → leaderboard) not validated
- Ant Design component interactions not tested
- No automated validation of responsive design
- AI-assisted debugging capabilities unused

## Solution Statement

Implement a complete Playwright testing infrastructure with MCP integration that:
1. **Configures Playwright** for React + Vite + TypeScript + Ant Design
2. **Creates test fixtures** for authentication, API mocking, and data seeding
3. **Implements E2E test suites** for all 8 pages and critical user workflows
4. **Integrates Playwright MCP** for AI-driven test generation and debugging
5. **Adds visual regression testing** with screenshot comparison
6. **Provides CI/CD integration** with automated test execution
7. **Generates test reports** with HTML and trace viewing capabilities

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium-High (6-8 hours)
**Primary Systems Affected**: 
- Frontend testing infrastructure
- CI/CD pipeline
- Development workflow
- MCP integration

**Dependencies**: 
- @playwright/test (latest)
- @playwright/mcp (already configured)
- Playwright browsers (chromium, firefox, webkit)
- Docker services (for integration testing)

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Frontend Structure:**
- `frontend/package.json` - Current dependencies and scripts (Vitest configured)
- `frontend/vite.config.ts` - Vite configuration with proxy setup
- `frontend/tsconfig.json` - TypeScript configuration (strict mode disabled)
- `frontend/src/App.tsx` (lines 1-100) - Main app structure, routing, auth integration
- `frontend/src/pages/LoginPage.tsx` - Authentication flow entry point
- `frontend/src/pages/ContestsPage.tsx` - Contest management UI
- `frontend/src/pages/PredictionsPage.tsx` - Prediction submission UI
- `frontend/src/pages/TeamsPage.tsx` - Team tournament UI
- `frontend/src/pages/AnalyticsPage.tsx` - Analytics dashboard
- `frontend/src/pages/ProfilePage.tsx` - User profile management
- `frontend/src/components/auth/LoginForm.tsx` - Login form component
- `frontend/src/components/auth/RegisterForm.tsx` - Registration form
- `frontend/src/components/auth/ProtectedRoute.tsx` - Route protection logic
- `frontend/src/contexts/AuthContext.tsx` - Authentication state management

**Testing Infrastructure:**
- `tests/e2e/` - Existing Go-based backend E2E tests (pattern reference)
- `scripts/e2e-test.sh` - E2E test execution script (Docker orchestration)
- `Makefile` (lines 60-90) - Test commands and infrastructure setup
- `.kiro/settings/mcp.json` - MCP server configuration (Playwright already configured)

**Backend Services:**
- `docker-compose.yml` - Service orchestration (ports, dependencies)
- `backend/api-gateway/` - API Gateway on port 8080
- `scripts/seed-data.sh` - Test data generation

### New Files to Create

**Playwright Configuration:**
- `frontend/playwright.config.ts` - Main Playwright configuration
- `frontend/.env.test` - Test environment variables
- `frontend/playwright/index.html` - Component test mount point (if using component testing)

**Test Utilities:**
- `frontend/tests/fixtures/auth.fixture.ts` - Authentication test fixtures
- `frontend/tests/fixtures/api.fixture.ts` - API mocking fixtures
- `frontend/tests/fixtures/data.fixture.ts` - Test data generators
- `frontend/tests/helpers/test-utils.ts` - Common test utilities
- `frontend/tests/helpers/selectors.ts` - Reusable selectors
- `frontend/tests/helpers/assertions.ts` - Custom assertions

**E2E Test Suites:**
- `frontend/tests/e2e/auth.spec.ts` - Authentication flows (login, register, logout)
- `frontend/tests/e2e/contests.spec.ts` - Contest management (create, view, join)
- `frontend/tests/e2e/predictions.spec.ts` - Prediction submission and viewing
- `frontend/tests/e2e/teams.spec.ts` - Team tournament workflows
- `frontend/tests/e2e/analytics.spec.ts` - Analytics dashboard interactions
- `frontend/tests/e2e/profile.spec.ts` - User profile management
- `frontend/tests/e2e/navigation.spec.ts` - Navigation and routing
- `frontend/tests/e2e/workflows.spec.ts` - Complete user workflows

**Visual Testing:**
- `frontend/tests/visual/snapshots.spec.ts` - Visual regression tests
- `frontend/tests/visual/.gitignore` - Ignore generated screenshots

**Scripts:**
- `scripts/playwright-test.sh` - Playwright test execution script
- `scripts/playwright-install.sh` - Browser installation script

**Documentation:**
- `docs/en/testing/playwright-testing.md` - English documentation
- `docs/ru/testing/playwright-testing.md` - Russian documentation

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

**Playwright Official Documentation:**
- [Playwright Test Configuration](https://playwright.dev/docs/test-configuration)
  - Section: Configuration file, test options, projects
  - Why: Required for proper Playwright setup with TypeScript and Vite
- [Playwright Fixtures](https://playwright.dev/docs/test-fixtures)
  - Section: Built-in fixtures, custom fixtures, fixture composition
  - Why: Essential for creating reusable test setup (auth, API mocking)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
  - Section: Test isolation, locators, assertions, debugging
  - Why: Ensures tests are reliable and maintainable
- [Playwright with React](https://playwright.dev/docs/test-components)
  - Section: Component testing, mounting components
  - Why: Optional component testing setup for isolated component tests

**Playwright MCP Integration:**
- [Playwright MCP Server](https://github.com/microsoft/playwright-mcp)
  - Section: Installation, available tools, AI integration
  - Why: Enables AI-driven test generation and debugging through MCP
- [MCP Protocol Documentation](https://modelcontextprotocol.io/introduction)
  - Section: Server configuration, tool invocation
  - Why: Understanding how MCP integrates with Kiro CLI

**React + Vite + Playwright:**
- [Playwright with Vite](https://playwright.dev/docs/test-webserver)
  - Section: webServer configuration, dev server integration
  - Why: Proper integration with Vite dev server for testing
- [Testing React Applications](https://playwright.dev/docs/test-react)
  - Section: Locators for React, testing patterns
  - Why: React-specific testing strategies

**Ant Design Testing:**
- [Ant Design Testing](https://ant.design/docs/react/testing)
  - Section: Component testing, accessibility
  - Why: Proper selectors and interaction patterns for Ant Design components

### Patterns to Follow

**Test File Naming Convention:**
```typescript
// Pattern from existing codebase
// frontend/src/tests/fixes.test.ts - Vitest unit tests
// tests/e2e/*_test.go - Go E2E tests

// New Playwright pattern:
// frontend/tests/e2e/*.spec.ts - E2E tests
// frontend/tests/visual/*.spec.ts - Visual tests
```

**TypeScript Configuration:**
```typescript
// From frontend/tsconfig.json
// Note: strict mode is DISABLED, noImplicitAny is false
// Tests should follow same relaxed TypeScript settings
{
  "strict": false,
  "noImplicitAny": false
}
```

**API Proxy Pattern:**
```typescript
// From frontend/vite.config.ts
server: {
  proxy: {
    '/v1': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      secure: false,
    }
  }
}
// Tests should use same proxy or mock API calls
```

**Authentication Pattern:**
```typescript
// From frontend/src/contexts/AuthContext.tsx
// Login flow: POST /v1/auth/login → JWT token → localStorage
// Protected routes check isAuthenticated state
// Pattern: login() → navigate() → protected page
```

**Error Handling Pattern:**
```typescript
// From frontend/src/pages/LoginPage.tsx
try {
  await login(data.email, data.password)
  navigate(from, { replace: true })
} catch (error) {
  // Error handled by auth context (toast notification)
}
```

**Component Testing Pattern:**
```typescript
// From frontend/src/tests/fixes.test.ts
import { describe, it, expect } from 'vitest'

describe('Feature Name', () => {
  it('should do something specific', () => {
    expect(result).toBe(expected)
  })
})
```

**Docker Service Pattern:**
```bash
# From scripts/e2e-test.sh
# 1. Start infrastructure (postgres, redis)
# 2. Wait for database readiness
# 3. Start all microservices
# 4. Wait for API Gateway health check
# 5. Run tests
# 6. Cleanup services
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation & Configuration

Set up Playwright infrastructure, install dependencies, and configure test environment.

**Tasks:**
1. Install Playwright and dependencies
2. Create Playwright configuration file
3. Set up test directory structure
4. Configure TypeScript for tests
5. Create environment configuration
6. Install browser binaries

### Phase 2: Test Fixtures & Utilities

Create reusable test fixtures for authentication, API mocking, and common test utilities.

**Tasks:**
1. Create authentication fixtures (login, register, logout)
2. Build API mocking fixtures for backend services
3. Implement test data generators
4. Create common test utilities (wait helpers, selectors)
5. Build custom assertions for Ant Design components
6. Set up page object models (optional but recommended)

### Phase 3: Core E2E Test Suites

Implement comprehensive E2E tests for all pages and user workflows.

**Tasks:**
1. Authentication tests (login, register, logout, protected routes)
2. Contest management tests (create, view, join, leave)
3. Prediction tests (submit, view, edit)
4. Team tournament tests (create team, invite, manage)
5. Analytics dashboard tests (view stats, charts)
6. Profile management tests (update profile, preferences)
7. Navigation tests (routing, menu, breadcrumbs)
8. Complete workflow tests (end-to-end user journeys)

### Phase 4: Visual Testing & MCP Integration

Add visual regression testing and integrate with Playwright MCP server.

**Tasks:**
1. Configure visual testing with screenshot comparison
2. Create baseline screenshots for all pages
3. Implement visual regression tests
4. Document MCP integration for AI-assisted testing
5. Create examples of MCP-driven test generation

### Phase 5: CI/CD Integration & Documentation

Integrate tests into CI/CD pipeline and create comprehensive documentation.

**Tasks:**
1. Create test execution scripts
2. Add Makefile commands for Playwright tests
3. Configure GitHub Actions (if applicable)
4. Generate HTML test reports
5. Create bilingual documentation (EN/RU)
6. Add troubleshooting guide

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.


### Task 1: CREATE frontend/playwright.config.ts

- **IMPLEMENT**: Playwright configuration for React + Vite + TypeScript
- **PATTERN**: Similar to tests/e2e/main_test.go (service orchestration)
- **IMPORTS**: @playwright/test, path, dotenv
- **CONFIG**:
  - testDir: './tests/e2e'
  - timeout: 30000ms
  - retries: 2 (CI), 0 (local)
  - workers: 4 (parallel execution)
  - reporter: ['html', 'list', 'json']
  - use: baseURL 'http://localhost:3000', screenshot 'only-on-failure', trace 'retain-on-failure'
  - projects: chromium, firefox, webkit
  - webServer: command 'npm run dev', port 3000, reuseExistingServer true
- **GOTCHA**: Vite dev server must be running or auto-started by webServer config
- **VALIDATE**: `npx playwright test --config=frontend/playwright.config.ts --list`

### Task 2: UPDATE frontend/package.json

- **IMPLEMENT**: Add Playwright dependencies and test scripts
- **PATTERN**: Existing scripts structure (dev, build, test)
- **DEPENDENCIES**:
  ```json
  "devDependencies": {
    "@playwright/test": "^1.48.0",
    "dotenv": "^16.4.5"
  }
  ```
- **SCRIPTS**:
  ```json
  "test:e2e": "playwright test",
  "test:e2e:ui": "playwright test --ui",
  "test:e2e:debug": "playwright test --debug",
  "test:e2e:headed": "playwright test --headed",
  "test:e2e:report": "playwright show-report",
  "playwright:install": "playwright install --with-deps"
  ```
- **VALIDATE**: `cd frontend && npm install && npm run playwright:install`

### Task 3: CREATE frontend/.env.test

- **IMPLEMENT**: Test environment variables
- **PATTERN**: Similar to .env.example pattern
- **CONTENT**:
  ```env
  VITE_API_URL=http://localhost:8080
  VITE_DEFAULT_REDIRECT=/contests
  # Test user credentials
  TEST_USER_EMAIL=test@example.com
  TEST_USER_PASSWORD=TestPass123!
  TEST_ADMIN_EMAIL=admin@example.com
  TEST_ADMIN_PASSWORD=AdminPass123!
  ```
- **GOTCHA**: Must not commit real credentials, use test-only accounts
- **VALIDATE**: `cat frontend/.env.test | grep -E "^(VITE_|TEST_)"`

### Task 4: CREATE frontend/tests/helpers/test-utils.ts

- **IMPLEMENT**: Common test utilities and helper functions
- **PATTERN**: Reusable utilities pattern
- **FUNCTIONS**:
  - `waitForPageLoad(page)` - Wait for network idle
  - `fillAntdInput(page, label, value)` - Fill Ant Design input by label
  - `clickAntdButton(page, text)` - Click Ant Design button by text
  - `selectAntdOption(page, label, option)` - Select from Ant Design dropdown
  - `waitForAntdNotification(page, type)` - Wait for success/error notification
  - `getTableRowCount(page, tableSelector)` - Count Ant Design table rows
  - `clearLocalStorage(page)` - Clear auth tokens
- **IMPORTS**: @playwright/test, Page, Locator
- **VALIDATE**: `npx tsc --noEmit frontend/tests/helpers/test-utils.ts`

### Task 5: CREATE frontend/tests/helpers/selectors.ts

- **IMPLEMENT**: Centralized selectors for Ant Design components
- **PATTERN**: Selector constants for maintainability
- **SELECTORS**:
  ```typescript
  export const SELECTORS = {
    auth: {
      emailInput: 'input[type="email"]',
      passwordInput: 'input[type="password"]',
      loginButton: 'button:has-text("Login")',
      registerButton: 'button:has-text("Register")',
      logoutButton: 'button:has-text("Logout")',
    },
    navigation: {
      contestsLink: 'a[href="/contests"]',
      predictionsLink: 'a[href="/predictions"]',
      teamsLink: 'a[href="/teams"]',
      analyticsLink: 'a[href="/analytics"]',
      profileLink: 'a[href="/profile"]',
    },
    contests: {
      createButton: 'button:has-text("Create Contest")',
      contestCard: '.ant-card',
      joinButton: 'button:has-text("Join")',
      leaveButton: 'button:has-text("Leave")',
    },
    // ... more selectors
  }
  ```
- **VALIDATE**: `npx tsc --noEmit frontend/tests/helpers/selectors.ts`

### Task 6: CREATE frontend/tests/fixtures/auth.fixture.ts

- **IMPLEMENT**: Authentication test fixtures
- **PATTERN**: Playwright fixtures pattern from docs
- **FIXTURES**:
  - `authenticatedPage` - Page with logged-in user
  - `adminPage` - Page with admin user
  - `testUser` - Test user credentials
  - `testAdmin` - Admin credentials
- **IMPLEMENTATION**:
  ```typescript
  import { test as base } from '@playwright/test'
  
  export const test = base.extend({
    authenticatedPage: async ({ page }, use) => {
      // Navigate to login
      await page.goto('/login')
      // Fill credentials
      await page.fill('input[type="email"]', process.env.TEST_USER_EMAIL!)
      await page.fill('input[type="password"]', process.env.TEST_USER_PASSWORD!)
      // Submit
      await page.click('button:has-text("Login")')
      // Wait for redirect
      await page.waitForURL('/contests')
      await use(page)
    }
  })
  ```
- **GOTCHA**: Must handle auth token in localStorage
- **VALIDATE**: `npx tsc --noEmit frontend/tests/fixtures/auth.fixture.ts`

### Task 7: CREATE frontend/tests/fixtures/api.fixture.ts

- **IMPLEMENT**: API mocking fixtures for isolated testing
- **PATTERN**: Playwright route mocking
- **FIXTURES**:
  - `mockedAPI` - Page with mocked API responses
  - `mockContests` - Mock contest data
  - `mockPredictions` - Mock prediction data
- **IMPLEMENTATION**:
  ```typescript
  export const test = base.extend({
    mockedAPI: async ({ page }, use) => {
      // Mock successful login
      await page.route('**/v1/auth/login', route => {
        route.fulfill({
          status: 200,
          body: JSON.stringify({ token: 'mock-jwt-token', user: { id: 1, email: 'test@example.com' } })
        })
      })
      // Mock contests list
      await page.route('**/v1/contests', route => {
        route.fulfill({
          status: 200,
          body: JSON.stringify({ contests: [] })
        })
      })
      await use(page)
    }
  })
  ```
- **VALIDATE**: `npx tsc --noEmit frontend/tests/fixtures/api.fixture.ts`

### Task 8: CREATE frontend/tests/fixtures/data.fixture.ts

- **IMPLEMENT**: Test data generators using gofakeit pattern
- **PATTERN**: Similar to backend/shared/seeder/factory.go
- **FUNCTIONS**:
  - `generateUser()` - Random user data
  - `generateContest()` - Random contest data
  - `generatePrediction()` - Random prediction data
  - `generateTeam()` - Random team data
- **IMPLEMENTATION**:
  ```typescript
  export function generateUser() {
    return {
      email: `test-${Date.now()}@example.com`,
      password: 'TestPass123!',
      username: `user_${Date.now()}`,
      displayName: 'Test User'
    }
  }
  
  export function generateContest() {
    const now = new Date()
    const startDate = new Date(now.getTime() + 24 * 60 * 60 * 1000)
    const endDate = new Date(startDate.getTime() + 7 * 24 * 60 * 60 * 1000)
    
    return {
      title: `Test Contest ${Date.now()}`,
      description: 'Automated test contest',
      sportType: 'Football',
      startDate: startDate.toISOString(),
      endDate: endDate.toISOString(),
      maxParticipants: 100
    }
  }
  ```
- **VALIDATE**: `npx tsc --noEmit frontend/tests/fixtures/data.fixture.ts`

### Task 9: CREATE frontend/tests/e2e/auth.spec.ts

- **IMPLEMENT**: Authentication flow E2E tests
- **PATTERN**: Playwright test structure with describe/test blocks
- **TEST CASES**:
  1. User can view login page
  2. User can login with valid credentials
  3. User cannot login with invalid credentials
  4. User can register new account
  5. User cannot register with existing email
  6. User can logout
  7. Protected routes redirect to login
  8. Authenticated user can access protected routes
- **IMPLEMENTATION**:
  ```typescript
  import { test, expect } from '@playwright/test'
  import { SELECTORS } from '../helpers/selectors'
  
  test.describe('Authentication', () => {
    test('should login with valid credentials', async ({ page }) => {
      await page.goto('/login')
      await page.fill(SELECTORS.auth.emailInput, process.env.TEST_USER_EMAIL!)
      await page.fill(SELECTORS.auth.passwordInput, process.env.TEST_USER_PASSWORD!)
      await page.click(SELECTORS.auth.loginButton)
      await expect(page).toHaveURL('/contests')
    })
    
    test('should show error with invalid credentials', async ({ page }) => {
      await page.goto('/login')
      await page.fill(SELECTORS.auth.emailInput, 'wrong@example.com')
      await page.fill(SELECTORS.auth.passwordInput, 'wrongpass')
      await page.click(SELECTORS.auth.loginButton)
      await expect(page.locator('.ant-notification-notice-error')).toBeVisible()
    })
  })
  ```
- **VALIDATE**: `cd frontend && npm run test:e2e -- auth.spec.ts`

### Task 10: CREATE frontend/tests/e2e/contests.spec.ts

- **IMPLEMENT**: Contest management E2E tests
- **PATTERN**: Use authenticatedPage fixture
- **TEST CASES**:
  1. User can view contests list
  2. User can create new contest
  3. User can view contest details
  4. User can join contest
  5. User can leave contest
  6. User can filter contests by sport type
  7. User can search contests by title
  8. Contest form validation works correctly
- **IMPORTS**: test from auth.fixture, expect, generateContest
- **VALIDATE**: `cd frontend && npm run test:e2e -- contests.spec.ts`

### Task 11: CREATE frontend/tests/e2e/predictions.spec.ts

- **IMPLEMENT**: Prediction submission and viewing tests
- **PATTERN**: Use authenticatedPage fixture
- **TEST CASES**:
  1. User can view predictions page
  2. User can submit prediction for event
  3. User can view their predictions
  4. User can edit prediction before deadline
  5. User cannot edit prediction after deadline
  6. Prediction form validation works
  7. Props predictions work correctly
- **VALIDATE**: `cd frontend && npm run test:e2e -- predictions.spec.ts`

### Task 12: CREATE frontend/tests/e2e/teams.spec.ts

- **IMPLEMENT**: Team tournament workflow tests
- **PATTERN**: Use authenticatedPage fixture
- **TEST CASES**:
  1. User can view teams page
  2. User can create new team
  3. User can invite members to team
  4. User can accept team invitation
  5. User can leave team
  6. User can view team leaderboard
  7. Team admin can manage members
- **VALIDATE**: `cd frontend && npm run test:e2e -- teams.spec.ts`

### Task 13: CREATE frontend/tests/e2e/analytics.spec.ts

- **IMPLEMENT**: Analytics dashboard interaction tests
- **PATTERN**: Use authenticatedPage fixture
- **TEST CASES**:
  1. User can view analytics dashboard
  2. Charts render correctly
  3. User can filter by date range
  4. User can filter by sport type
  5. Statistics display correctly
  6. Accuracy metrics calculate properly
- **VALIDATE**: `cd frontend && npm run test:e2e -- analytics.spec.ts`

### Task 14: CREATE frontend/tests/e2e/profile.spec.ts

- **IMPLEMENT**: User profile management tests
- **PATTERN**: Use authenticatedPage fixture
- **TEST CASES**:
  1. User can view profile page
  2. User can update display name
  3. User can update preferences
  4. User can view prediction history
  5. User can view achievements
  6. Profile form validation works
- **VALIDATE**: `cd frontend && npm run test:e2e -- profile.spec.ts`

### Task 15: CREATE frontend/tests/e2e/navigation.spec.ts

- **IMPLEMENT**: Navigation and routing tests
- **PATTERN**: Use authenticatedPage fixture
- **TEST CASES**:
  1. All navigation links work correctly
  2. Breadcrumbs display correctly
  3. Back button navigation works
  4. Deep linking works
  5. 404 page displays for invalid routes
  6. Menu highlights active page
- **VALIDATE**: `cd frontend && npm run test:e2e -- navigation.spec.ts`

### Task 16: CREATE frontend/tests/e2e/workflows.spec.ts

- **IMPLEMENT**: Complete end-to-end user workflow tests
- **PATTERN**: Multi-step workflows combining multiple features
- **TEST CASES**:
  1. Complete user journey: Register → Login → Create Contest → Join → Predict → View Leaderboard
  2. Team workflow: Create Team → Invite Members → Join Contest → Compete
  3. Analytics workflow: Make Predictions → View Stats → Track Progress
  4. Challenge workflow: Create Challenge → Accept → Compete → View Results
- **IMPLEMENTATION**:
  ```typescript
  test('complete user journey', async ({ page }) => {
    // Register
    await page.goto('/register')
    const user = generateUser()
    await page.fill(SELECTORS.auth.emailInput, user.email)
    await page.fill(SELECTORS.auth.passwordInput, user.password)
    await page.click(SELECTORS.auth.registerButton)
    
    // Create contest
    await page.goto('/contests')
    await page.click(SELECTORS.contests.createButton)
    const contest = generateContest()
    await fillContestForm(page, contest)
    await page.click('button:has-text("Create")')
    
    // Join contest
    await page.click(SELECTORS.contests.joinButton)
    
    // Make prediction
    await page.goto('/predictions')
    await submitPrediction(page)
    
    // View leaderboard
    await page.goto('/contests')
    await page.click('.ant-tabs-tab:has-text("Leaderboard")')
    await expect(page.locator('.leaderboard-table')).toBeVisible()
  })
  ```
- **VALIDATE**: `cd frontend && npm run test:e2e -- workflows.spec.ts`

### Task 17: CREATE frontend/tests/visual/snapshots.spec.ts

- **IMPLEMENT**: Visual regression tests with screenshot comparison
- **PATTERN**: Playwright screenshot testing
- **TEST CASES**:
  1. Login page visual snapshot
  2. Contests page visual snapshot
  3. Predictions page visual snapshot
  4. Teams page visual snapshot
  5. Analytics page visual snapshot
  6. Profile page visual snapshot
  7. Mobile responsive snapshots
- **IMPLEMENTATION**:
  ```typescript
  import { test, expect } from '@playwright/test'
  
  test.describe('Visual Regression', () => {
    test('login page matches snapshot', async ({ page }) => {
      await page.goto('/login')
      await expect(page).toHaveScreenshot('login-page.png')
    })
    
    test('contests page matches snapshot', async ({ page }) => {
      // Login first
      await loginUser(page)
      await page.goto('/contests')
      await page.waitForLoadState('networkidle')
      await expect(page).toHaveScreenshot('contests-page.png')
    })
  })
  
  test.describe('Mobile Visual Regression', () => {
    test.use({ viewport: { width: 375, height: 667 } })
    
    test('mobile login page matches snapshot', async ({ page }) => {
      await page.goto('/login')
      await expect(page).toHaveScreenshot('mobile-login-page.png')
    })
  })
  ```
- **GOTCHA**: First run generates baseline, subsequent runs compare
- **VALIDATE**: `cd frontend && npm run test:e2e -- visual/snapshots.spec.ts --update-snapshots`

### Task 18: CREATE frontend/tests/visual/.gitignore

- **IMPLEMENT**: Ignore generated screenshots except baselines
- **CONTENT**:
  ```
  # Ignore test results
  *-actual.png
  *-diff.png
  
  # Keep baseline snapshots
  !*-expected.png
  ```
- **VALIDATE**: `cat frontend/tests/visual/.gitignore`

### Task 19: CREATE scripts/playwright-install.sh

- **IMPLEMENT**: Browser installation script
- **PATTERN**: Similar to scripts/setup.sh
- **IMPLEMENTATION**:
  ```bash
  #!/bin/bash
  set -e
  
  echo "Installing Playwright browsers..."
  cd frontend
  
  # Install dependencies if needed
  if [ ! -d "node_modules" ]; then
    echo "Installing npm dependencies..."
    npm install
  fi
  
  # Install Playwright browsers
  echo "Installing Playwright browsers (chromium, firefox, webkit)..."
  npx playwright install --with-deps
  
  echo "✅ Playwright browsers installed successfully!"
  ```
- **VALIDATE**: `chmod +x scripts/playwright-install.sh && ./scripts/playwright-install.sh`

### Task 20: CREATE scripts/playwright-test.sh

- **IMPLEMENT**: Playwright test execution script with service orchestration
- **PATTERN**: Mirror scripts/e2e-test.sh structure
- **IMPLEMENTATION**:
  ```bash
  #!/bin/bash
  set -e
  
  echo "=========================================="
  echo "Frontend E2E Tests - Playwright"
  echo "=========================================="
  
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  YELLOW='\033[1;33m'
  NC='\033[0m'
  
  cd "$(dirname "$0")/.."
  
  # Parse arguments
  HEADED=false
  UI_MODE=false
  DEBUG=false
  TEST_FILE=""
  
  while [[ $# -gt 0 ]]; do
    case $1 in
      --headed) HEADED=true; shift ;;
      --ui) UI_MODE=true; shift ;;
      --debug) DEBUG=true; shift ;;
      *) TEST_FILE="$1"; shift ;;
    esac
  done
  
  echo -e "${YELLOW}Starting infrastructure services...${NC}"
  docker-compose up -d postgres redis
  
  echo -e "${YELLOW}Waiting for database...${NC}"
  MAX_RETRIES=30
  RETRY_COUNT=0
  until docker-compose exec -T postgres pg_isready -U sports_user -d sports_prediction > /dev/null 2>&1; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
      echo -e "${RED}Database not ready${NC}"
      docker-compose down
      exit 1
    fi
    sleep 1
  done
  
  echo -e "${YELLOW}Starting microservices...${NC}"
  docker-compose --profile services up -d
  sleep 15
  
  echo -e "${YELLOW}Checking API Gateway health...${NC}"
  RETRY_COUNT=0
  until curl -s http://localhost:8080/health > /dev/null 2>&1; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
      echo -e "${RED}API Gateway not ready${NC}"
      docker-compose --profile services down
      exit 1
    fi
    sleep 2
  done
  
  echo -e "${GREEN}Services ready!${NC}"
  
  # Run Playwright tests
  cd frontend
  
  if [ "$UI_MODE" = true ]; then
    echo -e "${YELLOW}Running Playwright in UI mode...${NC}"
    npm run test:e2e:ui
  elif [ "$DEBUG" = true ]; then
    echo -e "${YELLOW}Running Playwright in debug mode...${NC}"
    npm run test:e2e:debug $TEST_FILE
  elif [ "$HEADED" = true ]; then
    echo -e "${YELLOW}Running Playwright in headed mode...${NC}"
    npm run test:e2e:headed $TEST_FILE
  else
    echo -e "${YELLOW}Running Playwright tests...${NC}"
    if npm run test:e2e $TEST_FILE; then
      TEST_EXIT=0
      echo -e "${GREEN}All tests passed!${NC}"
    else
      TEST_EXIT=1
      echo -e "${RED}Some tests failed!${NC}"
    fi
  fi
  
  cd ..
  
  echo -e "${YELLOW}Stopping services...${NC}"
  docker-compose --profile services down
  docker-compose down
  
  exit $TEST_EXIT
  ```
- **VALIDATE**: `chmod +x scripts/playwright-test.sh && ./scripts/playwright-test.sh --help`

### Task 21: UPDATE Makefile

- **IMPLEMENT**: Add Playwright test commands
- **PATTERN**: Existing test commands (e2e-test, e2e-test-only)
- **ADD TARGETS**:
  ```makefile
  playwright-install: ## Install Playwright browsers
  	@echo "Installing Playwright browsers..."
  	@./scripts/playwright-install.sh
  
  playwright-test: ## Run Playwright E2E tests with services
  	@echo "Running Playwright E2E tests..."
  	@./scripts/playwright-test.sh
  
  playwright-test-ui: ## Run Playwright tests in UI mode
  	@echo "Running Playwright in UI mode..."
  	@./scripts/playwright-test.sh --ui
  
  playwright-test-headed: ## Run Playwright tests in headed mode
  	@echo "Running Playwright in headed mode..."
  	@./scripts/playwright-test.sh --headed
  
  playwright-test-only: ## Run Playwright tests (assumes services running)
  	@echo "Running Playwright tests..."
  	@cd frontend && npm run test:e2e
  
  playwright-report: ## Show Playwright test report
  	@cd frontend && npm run test:e2e:report
  ```
- **VALIDATE**: `make help | grep playwright`

### Task 22: CREATE docs/en/testing/playwright-testing.md

- **IMPLEMENT**: English documentation for Playwright testing
- **PATTERN**: Similar to docs/en/testing/e2e-testing.md structure
- **SECTIONS**:
  1. Overview and benefits
  2. Installation and setup
  3. Running tests (commands and options)
  4. Writing new tests (patterns and examples)
  5. Test fixtures and utilities
  6. Visual regression testing
  7. MCP integration for AI-assisted testing
  8. Debugging tests
  9. CI/CD integration
  10. Troubleshooting
- **CONTENT**:
  ```markdown
  # Playwright Frontend Testing
  
  ## Overview
  
  The Sports Prediction Contests platform uses Playwright for comprehensive end-to-end testing of the React frontend. Playwright provides reliable, fast, and cross-browser testing capabilities with AI-assisted test generation through MCP integration.
  
  ## Installation
  
  ### Prerequisites
  - Node.js 18+
  - Docker and Docker Compose (for service integration)
  
  ### Install Playwright
  
  \`\`\`bash
  # Install Playwright and browsers
  make playwright-install
  
  # Or manually
  cd frontend
  npm install
  npx playwright install --with-deps
  \`\`\`
  
  ## Running Tests
  
  ### Full Test Suite with Services
  
  \`\`\`bash
  # Run all tests with automatic service orchestration
  make playwright-test
  
  # Run specific test file
  ./scripts/playwright-test.sh auth.spec.ts
  \`\`\`
  
  ### Quick Testing (Services Already Running)
  
  \`\`\`bash
  # Start services first
  make docker-services
  
  # Run tests only
  make playwright-test-only
  
  # Or specific tests
  cd frontend
  npm run test:e2e -- auth.spec.ts
  \`\`\`
  
  ### Interactive Modes
  
  \`\`\`bash
  # UI Mode - Interactive test runner
  make playwright-test-ui
  
  # Headed Mode - See browser during tests
  make playwright-test-headed
  
  # Debug Mode - Step through tests
  cd frontend
  npm run test:e2e:debug auth.spec.ts
  \`\`\`
  
  ## Test Structure
  
  ### Test Organization
  
  \`\`\`
  frontend/tests/
  ├── e2e/              # End-to-end test suites
  │   ├── auth.spec.ts
  │   ├── contests.spec.ts
  │   ├── predictions.spec.ts
  │   ├── teams.spec.ts
  │   ├── analytics.spec.ts
  │   ├── profile.spec.ts
  │   ├── navigation.spec.ts
  │   └── workflows.spec.ts
  ├── visual/           # Visual regression tests
  │   └── snapshots.spec.ts
  ├── fixtures/         # Test fixtures
  │   ├── auth.fixture.ts
  │   ├── api.fixture.ts
  │   └── data.fixture.ts
  └── helpers/          # Test utilities
      ├── test-utils.ts
      ├── selectors.ts
      └── assertions.ts
  \`\`\`
  
  ## Writing Tests
  
  ### Basic Test Example
  
  \`\`\`typescript
  import { test, expect } from '@playwright/test'
  import { SELECTORS } from '../helpers/selectors'
  
  test.describe('Feature Name', () => {
    test('should do something', async ({ page }) => {
      await page.goto('/page')
      await page.click(SELECTORS.someButton)
      await expect(page.locator('.result')).toBeVisible()
    })
  })
  \`\`\`
  
  ### Using Authentication Fixture
  
  \`\`\`typescript
  import { test } from '../fixtures/auth.fixture'
  import { expect } from '@playwright/test'
  
  test('authenticated user can access protected page', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/contests')
    await expect(authenticatedPage).toHaveURL('/contests')
  })
  \`\`\`
  
  ## MCP Integration
  
  The Playwright MCP server is configured in \`.kiro/settings/mcp.json\` and enables AI-assisted testing through Kiro CLI.
  
  ### Using MCP for Test Generation
  
  In Kiro CLI, you can leverage the Playwright MCP server to:
  - Generate test code from natural language descriptions
  - Debug failing tests with AI assistance
  - Capture screenshots and analyze UI issues
  - Automate test maintenance
  
  Example workflow:
  \`\`\`
  # In Kiro CLI
  > Use Playwright MCP to create a test for user login flow
  > Debug the failing contest creation test
  > Capture screenshot of the analytics dashboard
  \`\`\`
  
  ## Visual Regression Testing
  
  ### Creating Baseline Screenshots
  
  \`\`\`bash
  cd frontend
  npm run test:e2e -- visual/snapshots.spec.ts --update-snapshots
  \`\`\`
  
  ### Running Visual Tests
  
  \`\`\`bash
  npm run test:e2e -- visual/snapshots.spec.ts
  \`\`\`
  
  ## Troubleshooting
  
  ### Tests Fail with "Connection Refused"
  - Ensure services are running: \`make docker-services\`
  - Check API Gateway health: \`curl http://localhost:8080/health\`
  
  ### Browser Installation Issues
  - Run: \`npx playwright install --with-deps\`
  - On Linux, may need system dependencies
  
  ### Flaky Tests
  - Increase timeout in playwright.config.ts
  - Add explicit waits: \`await page.waitForLoadState('networkidle')\`
  - Use \`test.retry(2)\` for specific tests
  
  ## CI/CD Integration
  
  Playwright tests can be integrated into GitHub Actions or other CI/CD pipelines:
  
  \`\`\`yaml
  - name: Run Playwright Tests
    run: |
      make playwright-install
      make playwright-test
  
  - name: Upload Test Report
    if: always()
    uses: actions/upload-artifact@v3
    with:
      name: playwright-report
      path: frontend/playwright-report/
  \`\`\`
  ```
- **VALIDATE**: `cat docs/en/testing/playwright-testing.md | wc -l`

### Task 23: CREATE docs/ru/testing/playwright-testing.md

- **IMPLEMENT**: Russian translation of Playwright testing documentation
- **PATTERN**: Mirror English documentation structure
- **CONTENT**: Full Russian translation of Task 22 content
- **VALIDATE**: `cat docs/ru/testing/playwright-testing.md | wc -l`

### Task 24: UPDATE docs/en/README.md

- **IMPLEMENT**: Add Playwright testing to documentation index
- **PATTERN**: Existing documentation structure
- **ADD SECTION**:
  ```markdown
  ### 🧪 Testing
  - [E2E Testing](testing/e2e-testing.md) - Backend end-to-end testing
  - [Playwright Testing](testing/playwright-testing.md) - Frontend E2E testing with Playwright
  - [Unit Testing](testing/unit-testing.md) - Unit testing for all services
  - [Performance Testing](testing/performance-testing.md) - Load testing and benchmarks
  ```
- **VALIDATE**: `grep -A 5 "Testing" docs/en/README.md`

### Task 25: UPDATE docs/ru/README.md

- **IMPLEMENT**: Add Playwright testing to Russian documentation index
- **PATTERN**: Mirror English documentation update
- **VALIDATE**: `grep -A 5 "Тестирование" docs/ru/README.md`

### Task 26: CREATE .github/workflows/playwright.yml (Optional)

- **IMPLEMENT**: GitHub Actions workflow for Playwright tests
- **PATTERN**: CI/CD automation
- **IMPLEMENTATION**:
  ```yaml
  name: Playwright Tests
  
  on:
    push:
      branches: [ master, develop ]
    pull_request:
      branches: [ master, develop ]
  
  jobs:
    test:
      timeout-minutes: 60
      runs-on: ubuntu-latest
      steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 18
      
      - name: Install dependencies
        run: |
          cd frontend
          npm ci
      
      - name: Install Playwright Browsers
        run: npx playwright install --with-deps
      
      - name: Start services
        run: |
          docker-compose up -d postgres redis
          docker-compose --profile services up -d
      
      - name: Wait for services
        run: |
          timeout 60 bash -c 'until curl -s http://localhost:8080/health; do sleep 2; done'
      
      - name: Run Playwright tests
        run: |
          cd frontend
          npm run test:e2e
      
      - name: Upload test results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: playwright-report
          path: frontend/playwright-report/
          retention-days: 30
      
      - name: Cleanup
        if: always()
        run: |
          docker-compose --profile services down
          docker-compose down
  ```
- **VALIDATE**: `cat .github/workflows/playwright.yml | grep -E "^name:"`

---

## TESTING STRATEGY

### Unit Tests

**Scope**: Individual utility functions and helpers
- Test utilities in `frontend/tests/helpers/test-utils.ts`
- Selector constants validation
- Data generators in `frontend/tests/fixtures/data.fixture.ts`

**Framework**: Vitest (already configured)

**Pattern**:
```typescript
import { describe, it, expect } from 'vitest'
import { waitForPageLoad, fillAntdInput } from './test-utils'

describe('Test Utilities', () => {
  it('should generate valid user data', () => {
    const user = generateUser()
    expect(user.email).toMatch(/^test-\d+@example\.com$/)
    expect(user.password).toBe('TestPass123!')
  })
})
```

### Integration Tests

**Scope**: E2E test suites covering complete user workflows
- Authentication flows (login, register, logout)
- Contest management (CRUD operations)
- Prediction submission and viewing
- Team tournament workflows
- Analytics dashboard interactions
- Profile management
- Navigation and routing

**Framework**: Playwright Test

**Pattern**: Use fixtures for authentication and API mocking

### Visual Regression Tests

**Scope**: Screenshot comparison for all major pages
- Desktop viewport (1920x1080)
- Tablet viewport (768x1024)
- Mobile viewport (375x667)

**Baseline Management**:
- First run: `npm run test:e2e -- visual/snapshots.spec.ts --update-snapshots`
- Subsequent runs: Automatic comparison
- Update baselines after intentional UI changes

### Edge Cases

**Authentication Edge Cases**:
- Invalid credentials
- Expired tokens
- Network errors during login
- Concurrent login sessions
- Logout during active request

**Contest Edge Cases**:
- Creating contest with past dates
- Joining full contest
- Leaving contest with active predictions
- Contest deadline edge cases

**Prediction Edge Cases**:
- Submitting after deadline
- Editing locked predictions
- Invalid prediction values
- Network errors during submission

**Form Validation Edge Cases**:
- Empty required fields
- Invalid email formats
- Weak passwords
- Whitespace-only inputs
- Maximum length violations

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Installation & Configuration

```bash
# Install Playwright and browsers
cd frontend && npm install
npx playwright install --with-deps

# Verify Playwright configuration
npx playwright test --config=playwright.config.ts --list

# Check TypeScript compilation
npx tsc --noEmit
```

### Level 2: Unit Tests (Helpers & Utilities)

```bash
# Run Vitest unit tests
cd frontend && npm test

# Run with coverage
npm test -- --coverage
```

### Level 3: E2E Tests (Individual Suites)

```bash
# Test authentication flows
npm run test:e2e -- auth.spec.ts

# Test contest management
npm run test:e2e -- contests.spec.ts

# Test predictions
npm run test:e2e -- predictions.spec.ts

# Test teams
npm run test:e2e -- teams.spec.ts

# Test analytics
npm run test:e2e -- analytics.spec.ts

# Test profile
npm run test:e2e -- profile.spec.ts

# Test navigation
npm run test:e2e -- navigation.spec.ts

# Test complete workflows
npm run test:e2e -- workflows.spec.ts
```

### Level 4: Visual Regression Tests

```bash
# Run visual tests
npm run test:e2e -- visual/snapshots.spec.ts

# Update baselines (after intentional UI changes)
npm run test:e2e -- visual/snapshots.spec.ts --update-snapshots
```

### Level 5: Full Test Suite

```bash
# Run all tests with service orchestration
cd .. && make playwright-test

# Run in UI mode for debugging
make playwright-test-ui

# Run in headed mode to see browser
make playwright-test-headed
```

### Level 6: Cross-Browser Testing

```bash
# Run on all browsers (chromium, firefox, webkit)
cd frontend && npm run test:e2e

# Run on specific browser
npm run test:e2e -- --project=firefox
npm run test:e2e -- --project=webkit
```

### Level 7: Test Reports

```bash
# Generate and view HTML report
npm run test:e2e:report

# View trace for failed tests
npx playwright show-trace playwright-report/trace.zip
```

---

## ACCEPTANCE CRITERIA

- [ ] Playwright installed and configured for React + Vite + TypeScript
- [ ] All 8 E2E test suites implemented and passing
- [ ] Authentication fixtures working correctly
- [ ] API mocking fixtures functional
- [ ] Test data generators producing valid data
- [ ] Visual regression tests capturing baselines
- [ ] All validation commands pass with zero errors
- [ ] Test coverage includes all major user workflows
- [ ] Cross-browser testing works (chromium, firefox, webkit)
- [ ] Test execution scripts functional
- [ ] Makefile commands integrated
- [ ] Documentation complete (English and Russian)
- [ ] MCP integration documented
- [ ] CI/CD workflow configured (optional)
- [ ] Test reports generated successfully
- [ ] No flaky tests (all tests pass consistently)

---

## COMPLETION CHECKLIST

- [ ] All 26 tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + E2E + visual)
- [ ] No TypeScript compilation errors
- [ ] Cross-browser tests pass
- [ ] Test reports generated and viewable
- [ ] Documentation reviewed and complete
- [ ] MCP integration tested
- [ ] Scripts executable and functional
- [ ] Makefile commands working
- [ ] Services orchestration reliable
- [ ] No regressions in existing functionality

---

## NOTES

### Design Decisions

**Why Playwright over Cypress?**
- Better TypeScript support
- Native multi-browser testing (chromium, firefox, webkit)
- Faster execution with parallel testing
- MCP integration for AI-assisted testing
- Better debugging tools (trace viewer, inspector)

**Why Separate E2E and Visual Tests?**
- Visual tests are slower (screenshot comparison)
- E2E tests focus on functionality
- Allows running E2E tests more frequently
- Visual tests run on UI changes only

**Why Test Fixtures?**
- Reduces code duplication
- Ensures consistent test setup
- Makes tests more maintainable
- Follows Playwright best practices

### Trade-offs

**Service Orchestration**:
- **Pro**: Tests run against real backend services
- **Con**: Slower test execution
- **Mitigation**: API mocking fixtures for unit-level E2E tests

**Cross-Browser Testing**:
- **Pro**: Catches browser-specific issues
- **Con**: 3x test execution time
- **Mitigation**: Run chromium-only locally, all browsers in CI

**Visual Regression**:
- **Pro**: Catches unintended UI changes
- **Con**: Baseline management overhead
- **Mitigation**: Update baselines only on intentional UI changes

### MCP Integration Benefits

The Playwright MCP server enables:
1. **AI-Assisted Test Generation**: Describe test in natural language, get Playwright code
2. **Visual Debugging**: AI can analyze screenshots and suggest fixes
3. **Test Maintenance**: AI helps update tests when UI changes
4. **Flaky Test Analysis**: AI identifies patterns in flaky tests
5. **Coverage Gaps**: AI suggests missing test scenarios

### Future Enhancements

**Phase 2 (Post-Implementation)**:
- Component testing with Playwright CT
- Performance testing with Lighthouse
- Accessibility testing with axe-core
- API contract testing
- Load testing with k6

**Phase 3 (Advanced)**:
- Visual AI testing (Percy, Applitools)
- Chaos engineering tests
- Security testing (OWASP)
- Mobile app testing (React Native)

---

## CONFIDENCE SCORE

**Estimated One-Pass Success**: 8.5/10

**Reasoning**:
- ✅ Clear task breakdown with validation commands
- ✅ Existing patterns from backend E2E tests
- ✅ MCP server already configured
- ✅ Comprehensive documentation and examples
- ⚠️ Ant Design component selectors may need adjustment
- ⚠️ Service orchestration timing may need tuning
- ⚠️ Visual regression baselines need initial generation

**Risk Mitigation**:
- Start with authentication tests (simplest)
- Validate service orchestration early
- Generate visual baselines incrementally
- Use Playwright UI mode for debugging
- Leverage MCP for AI-assisted troubleshooting

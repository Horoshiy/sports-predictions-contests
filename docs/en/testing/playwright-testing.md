# Playwright Frontend Testing

## Overview

The Sports Prediction Contests platform uses Playwright for comprehensive end-to-end testing of the React frontend. Playwright provides reliable, fast, and cross-browser testing capabilities with AI-assisted test generation through MCP integration.

## Installation

### Prerequisites
- Node.js 18+
- Docker and Docker Compose (for service integration)

### Install Playwright

```bash
# Install Playwright and browsers
make playwright-install

# Or manually
cd frontend
npm install
npx playwright install --with-deps
```

## Running Tests

### Full Test Suite with Services

```bash
# Run all tests with automatic service orchestration
make playwright-test

# Run specific test file
./scripts/playwright-test.sh auth.spec.ts
```

### Quick Testing (Services Already Running)

```bash
# Start services first
make docker-services

# Run tests only
make playwright-test-only

# Or specific tests
cd frontend
npm run test:e2e -- auth.spec.ts
```

### Interactive Modes

```bash
# UI Mode - Interactive test runner
make playwright-test-ui

# Headed Mode - See browser during tests
make playwright-test-headed

# Debug Mode - Step through tests
cd frontend
npm run test:e2e:debug auth.spec.ts
```

## Test Structure

### Test Organization

```
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
```

## Writing Tests

### Basic Test Example

```typescript
import { test, expect } from '@playwright/test'
import { SELECTORS } from '../helpers/selectors'

test.describe('Feature Name', () => {
  test('should do something', async ({ page }) => {
    await page.goto('/page')
    await page.click(SELECTORS.someButton)
    await expect(page.locator('.result')).toBeVisible()
  })
})
```

### Using Authentication Fixture

```typescript
import { test } from '../fixtures/auth.fixture'
import { expect } from '@playwright/test'

test('authenticated user can access protected page', async ({ authenticatedPage }) => {
  await authenticatedPage.goto('/contests')
  await expect(authenticatedPage).toHaveURL('/contests')
})
```

## MCP Integration

The Playwright MCP server is configured in `.kiro/settings/mcp.json` and enables AI-assisted testing through Kiro CLI.

### Using MCP for Test Generation

In Kiro CLI, you can leverage the Playwright MCP server to:
- Generate test code from natural language descriptions
- Debug failing tests with AI assistance
- Capture screenshots and analyze UI issues
- Automate test maintenance

Example workflow:
```
# In Kiro CLI
> Use Playwright MCP to create a test for user login flow
> Debug the failing contest creation test
> Capture screenshot of the analytics dashboard
```

## Visual Regression Testing

### Creating Baseline Screenshots

```bash
cd frontend
npm run test:e2e -- visual/snapshots.spec.ts --update-snapshots
```

### Running Visual Tests

```bash
npm run test:e2e -- visual/snapshots.spec.ts
```

## Troubleshooting

### Tests Fail with "Connection Refused"
- Ensure services are running: `make docker-services`
- Check API Gateway health: `curl http://localhost:8080/health`

### Browser Installation Issues
- Run: `npx playwright install --with-deps`
- On Linux, may need system dependencies

### Flaky Tests
- Increase timeout in playwright.config.ts
- Add explicit waits: `await page.waitForLoadState('networkidle')`
- Use `test.retry(2)` for specific tests

## CI/CD Integration

Playwright tests can be integrated into GitHub Actions or other CI/CD pipelines:

```yaml
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
```

## Test Reports

```bash
# View HTML report
npm run test:e2e:report

# View trace for failed tests
npx playwright show-trace playwright-report/trace.zip
```

## Available Test Suites

- **auth.spec.ts** - Authentication flows (login, register, logout)
- **contests.spec.ts** - Contest management (create, view, join)
- **predictions.spec.ts** - Prediction submission and viewing
- **teams.spec.ts** - Team tournament workflows
- **analytics.spec.ts** - Analytics dashboard interactions
- **profile.spec.ts** - User profile management
- **navigation.spec.ts** - Navigation and routing
- **workflows.spec.ts** - Complete user workflows
- **visual/snapshots.spec.ts** - Visual regression tests

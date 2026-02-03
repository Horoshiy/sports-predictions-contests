# Plan: Playwright UI Tests with PageObject Pattern

**Feature:** Comprehensive UI testing with PageObject architecture
**Created:** 2026-02-03
**Status:** In Progress (Tasks 1-15 Complete)
**Commit:** `e388d27`

## Feature Description

Создать полноценное покрытие UI тестами с использованием Playwright и паттерна PageObject. Все клики по ссылкам и кнопкам должны быть протестированы.

### User Story
```
Как разработчик,
я хочу иметь надёжные UI тесты с PageObject паттерном,
чтобы быстро обнаруживать регрессии при изменении интерфейса.
```

### Сложность: Medium-High
- 8 страниц для тестирования
- ~50+ интерактивных элементов
- Требуется рефакторинг существующих тестов

---

## Context References

### Существующая структура
```
frontend/
├── tests/
│   ├── e2e/
│   │   ├── auth.spec.ts        # Существует (без PageObject)
│   │   ├── contests.spec.ts    # Существует
│   │   ├── navigation.spec.ts  # Существует
│   │   ├── predictions.spec.ts # Существует
│   │   ├── profile.spec.ts     # Существует
│   │   ├── teams.spec.ts       # Существует
│   │   └── workflows.spec.ts   # Существует
│   ├── helpers/
│   │   ├── selectors.ts        # Селекторы (рефакторинг в PageObject)
│   │   ├── test-config.ts      # Конфиг
│   │   └── test-utils.ts       # Утилиты
│   └── visual/
│       └── snapshots.spec.ts
├── playwright.config.ts         # Уже настроен
└── src/pages/
    ├── LoginPage.tsx
    ├── RegisterPage.tsx
    ├── ContestsPage.tsx
    ├── PredictionsPage.tsx
    ├── ProfilePage.tsx
    ├── SportsPage.tsx
    ├── TeamsPage.tsx
    └── AnalyticsPage.tsx
```

### Паттерны проекта
- Ant Design компоненты
- React Router для навигации
- JWT authentication
- REST API через fetch

---

## Implementation Plan

### Phase 1: PageObject Infrastructure (Tasks 1-3)
Создание базовой архитектуры PageObject.

### Phase 2: Core Page Objects (Tasks 4-8)
Реализация PageObject для каждой страницы.

### Phase 3: Navigation & Common Components (Tasks 9-10)
PageObject для навигации и общих компонентов.

### Phase 4: Test Migration & New Tests (Tasks 11-15)
Миграция существующих тестов и добавление новых.

### Phase 5: CI Integration & Validation (Task 16)
Интеграция с CI и финальная валидация.

---

## Step-by-Step Tasks

### Task 1: Create Base Page Object Class
**File:** `frontend/tests/pages/BasePage.ts`

```typescript
import { Page, Locator, expect } from '@playwright/test'

export abstract class BasePage {
  protected page: Page
  abstract readonly url: string
  
  constructor(page: Page) {
    this.page = page
  }
  
  async goto(): Promise<void>
  async waitForPageLoad(): Promise<void>
  async getTitle(): Promise<string>
  
  // Common Ant Design helpers
  protected async clickButton(text: string): Promise<void>
  protected async fillInput(selector: string, value: string): Promise<void>
  protected async selectOption(selector: string, value: string): Promise<void>
  protected async waitForNotification(type: 'success' | 'error'): Promise<void>
  protected async isVisible(selector: string): Promise<boolean>
}
```

**Validation:**
```bash
npx tsc --noEmit frontend/tests/pages/BasePage.ts
```

---

### Task 2: Create Test Fixtures
**File:** `frontend/tests/fixtures/test-fixtures.ts`

```typescript
import { test as base, expect } from '@playwright/test'
import { LoginPage } from '../pages/LoginPage'
import { ContestsPage } from '../pages/ContestsPage'
// ... other pages

type PageFixtures = {
  loginPage: LoginPage
  contestsPage: ContestsPage
  // ... other pages
  authenticatedPage: Page  // Pre-authenticated
}

export const test = base.extend<PageFixtures>({...})
export { expect }
```

**Validation:**
```bash
npx tsc --noEmit frontend/tests/fixtures/test-fixtures.ts
```

---

### Task 3: Create Page Object Directory Structure
**Files to create:**
```
frontend/tests/pages/
├── BasePage.ts
├── LoginPage.ts
├── RegisterPage.ts
├── ContestsPage.ts
├── PredictionsPage.ts
├── ProfilePage.ts
├── SportsPage.ts
├── TeamsPage.ts
├── AnalyticsPage.ts
└── components/
    ├── HeaderComponent.ts
    ├── NavigationComponent.ts
    ├── ModalComponent.ts
    └── NotificationComponent.ts
```

**Validation:**
```bash
ls -la frontend/tests/pages/
ls -la frontend/tests/pages/components/
```

---

### Task 4: Implement LoginPage PageObject
**File:** `frontend/tests/pages/LoginPage.ts`

```typescript
import { BasePage } from './BasePage'

export class LoginPage extends BasePage {
  readonly url = '/login'
  
  // Locators
  get emailInput(): Locator
  get passwordInput(): Locator
  get loginButton(): Locator
  get registerLink(): Locator
  get errorMessage(): Locator
  
  // Actions
  async login(email: string, password: string): Promise<void>
  async clickRegisterLink(): Promise<void>
  
  // Assertions
  async expectLoginFormVisible(): Promise<void>
  async expectErrorVisible(): Promise<void>
}
```

**Test coverage:**
- [ ] Display login form
- [ ] Login with valid credentials
- [ ] Login with invalid credentials
- [ ] Navigate to register
- [ ] Email validation
- [ ] Password validation

**Validation:**
```bash
npx playwright test tests/e2e/auth.spec.ts --grep "login"
```

---

### Task 5: Implement RegisterPage PageObject
**File:** `frontend/tests/pages/RegisterPage.ts`

```typescript
export class RegisterPage extends BasePage {
  readonly url = '/register'
  
  // Locators
  get emailInput(): Locator
  get passwordInput(): Locator
  get confirmPasswordInput(): Locator
  get nameInput(): Locator
  get registerButton(): Locator
  get loginLink(): Locator
  
  // Actions
  async register(name: string, email: string, password: string): Promise<void>
  async clickLoginLink(): Promise<void>
  
  // Assertions
  async expectRegisterFormVisible(): Promise<void>
  async expectPasswordMismatchError(): Promise<void>
}
```

**Test coverage:**
- [ ] Display register form
- [ ] Register new user
- [ ] Validate email format
- [ ] Password confirmation mismatch
- [ ] Navigate to login

**Validation:**
```bash
npx playwright test tests/e2e/auth.spec.ts --grep "register"
```

---

### Task 6: Implement ContestsPage PageObject
**File:** `frontend/tests/pages/ContestsPage.ts`

```typescript
export class ContestsPage extends BasePage {
  readonly url = '/contests'
  
  // Locators
  get createContestButton(): Locator
  get contestCards(): Locator
  get contestModal(): Locator
  get searchInput(): Locator
  get filterDropdown(): Locator
  
  // Actions
  async clickCreateContest(): Promise<void>
  async joinContest(index: number): Promise<void>
  async leaveContest(index: number): Promise<void>
  async openContestDetails(index: number): Promise<void>
  async searchContests(query: string): Promise<void>
  async filterByStatus(status: string): Promise<void>
  
  // Assertions
  async expectContestListVisible(): Promise<void>
  async expectContestCount(count: number): Promise<void>
  async expectJoinButtonVisible(index: number): Promise<void>
}
```

**Test coverage:**
- [ ] Display contests list
- [ ] Click create contest button
- [ ] Open contest details
- [ ] Join contest
- [ ] Leave contest
- [ ] Search contests
- [ ] Filter by status

**Validation:**
```bash
npx playwright test tests/e2e/contests.spec.ts
```

---

### Task 7: Implement PredictionsPage PageObject
**File:** `frontend/tests/pages/PredictionsPage.ts`

```typescript
export class PredictionsPage extends BasePage {
  readonly url = '/predictions'
  
  // Locators
  get eventCards(): Locator
  get predictionForm(): Locator
  get submitButton(): Locator
  get contestSelector(): Locator
  get predictionHistory(): Locator
  
  // Actions
  async selectContest(name: string): Promise<void>
  async submitPrediction(eventIndex: number, homeScore: number, awayScore: number): Promise<void>
  async editPrediction(eventIndex: number): Promise<void>
  
  // Assertions
  async expectEventsVisible(): Promise<void>
  async expectPredictionSubmitted(): Promise<void>
}
```

**Test coverage:**
- [ ] Display events list
- [ ] Select contest
- [ ] Submit prediction
- [ ] Edit existing prediction
- [ ] View prediction history

**Validation:**
```bash
npx playwright test tests/e2e/predictions.spec.ts
```

---

### Task 8: Implement ProfilePage PageObject
**File:** `frontend/tests/pages/ProfilePage.ts`

```typescript
export class ProfilePage extends BasePage {
  readonly url = '/profile'
  
  // Locators
  get displayNameInput(): Locator
  get avatarUpload(): Locator
  get saveButton(): Locator
  get statsSection(): Locator
  get privacySettings(): Locator
  
  // Actions
  async updateDisplayName(name: string): Promise<void>
  async uploadAvatar(filePath: string): Promise<void>
  async saveProfile(): Promise<void>
  async togglePrivacySetting(setting: string): Promise<void>
  
  // Assertions
  async expectProfileLoaded(): Promise<void>
  async expectSaveSuccess(): Promise<void>
}
```

**Test coverage:**
- [ ] Display profile form
- [ ] Edit display name
- [ ] Upload avatar
- [ ] Save profile
- [ ] View stats
- [ ] Toggle privacy settings

**Validation:**
```bash
npx playwright test tests/e2e/profile.spec.ts
```

---

### Task 9: Implement SportsPage PageObject
**File:** `frontend/tests/pages/SportsPage.ts`

```typescript
export class SportsPage extends BasePage {
  readonly url = '/sports'
  
  // Locators
  get sportsList(): Locator
  get leaguesList(): Locator
  get teamsList(): Locator
  get matchesList(): Locator
  get addSportButton(): Locator
  get addLeagueButton(): Locator
  get addTeamButton(): Locator
  get addMatchButton(): Locator
  
  // Actions
  async selectSport(name: string): Promise<void>
  async selectLeague(name: string): Promise<void>
  async addSport(name: string): Promise<void>
  async addLeague(sportId: string, name: string): Promise<void>
  async addTeam(name: string): Promise<void>
  async addMatch(homeTeam: string, awayTeam: string, date: string): Promise<void>
  
  // Assertions
  async expectSportsLoaded(): Promise<void>
  async expectLeaguesVisible(): Promise<void>
}
```

**Test coverage:**
- [ ] Display sports list
- [ ] Select sport
- [ ] Display leagues
- [ ] Select league
- [ ] Display teams
- [ ] Display matches
- [ ] Add new sport (admin)
- [ ] Add new league (admin)
- [ ] Add new team (admin)
- [ ] Add new match (admin)

**Validation:**
```bash
npx playwright test tests/e2e/sports.spec.ts
```

---

### Task 10: Implement TeamsPage PageObject
**File:** `frontend/tests/pages/TeamsPage.ts`

```typescript
export class TeamsPage extends BasePage {
  readonly url = '/teams'
  
  // Locators
  get createTeamButton(): Locator
  get teamCards(): Locator
  get joinTeamForm(): Locator
  get inviteCodeInput(): Locator
  get teamDetailsModal(): Locator
  get membersList(): Locator
  
  // Actions
  async createTeam(name: string, description: string): Promise<void>
  async joinTeamByCode(code: string): Promise<void>
  async openTeamDetails(index: number): Promise<void>
  async leaveTeam(): Promise<void>
  async inviteMember(): Promise<void>
  async removeMember(memberIndex: number): Promise<void>
  
  // Assertions
  async expectTeamsListVisible(): Promise<void>
  async expectTeamCreated(name: string): Promise<void>
  async expectMemberCount(count: number): Promise<void>
}
```

**Test coverage:**
- [ ] Display teams list
- [ ] Create new team
- [ ] Join team by code
- [ ] Open team details
- [ ] View members list
- [ ] Leave team
- [ ] Invite member (captain)
- [ ] Remove member (captain)

**Validation:**
```bash
npx playwright test tests/e2e/teams.spec.ts
```

---

### Task 11: Implement AnalyticsPage PageObject
**File:** `frontend/tests/pages/AnalyticsPage.ts`

```typescript
export class AnalyticsPage extends BasePage {
  readonly url = '/analytics'
  
  // Locators
  get accuracyChart(): Locator
  get sportBreakdown(): Locator
  get platformComparison(): Locator
  get exportButton(): Locator
  get dateFilter(): Locator
  
  // Actions
  async filterByDateRange(start: string, end: string): Promise<void>
  async exportData(format: string): Promise<void>
  async switchTab(tabName: string): Promise<void>
  
  // Assertions
  async expectChartsLoaded(): Promise<void>
  async expectDataExported(): Promise<void>
}
```

**Test coverage:**
- [ ] Display analytics dashboard
- [ ] View accuracy chart
- [ ] View sport breakdown
- [ ] Filter by date range
- [ ] Export data

**Validation:**
```bash
npx playwright test tests/e2e/analytics.spec.ts
```

---

### Task 12: Implement HeaderComponent PageObject
**File:** `frontend/tests/pages/components/HeaderComponent.ts`

```typescript
export class HeaderComponent {
  constructor(private page: Page) {}
  
  // Locators
  get logo(): Locator
  get contestsLink(): Locator
  get predictionsLink(): Locator
  get teamsLink(): Locator
  get sportsLink(): Locator
  get analyticsLink(): Locator
  get userMenu(): Locator
  get profileLink(): Locator
  get logoutButton(): Locator
  
  // Actions
  async navigateTo(page: 'contests' | 'predictions' | 'teams' | 'sports' | 'analytics'): Promise<void>
  async openUserMenu(): Promise<void>
  async logout(): Promise<void>
  async goToProfile(): Promise<void>
  
  // Assertions
  async expectUserLoggedIn(name: string): Promise<void>
  async expectActiveMenuItem(item: string): Promise<void>
}
```

**Test coverage:**
- [ ] Click logo (home)
- [ ] Click contests link
- [ ] Click predictions link
- [ ] Click teams link
- [ ] Click sports link
- [ ] Click analytics link
- [ ] Open user menu
- [ ] Click profile link
- [ ] Click logout button

**Validation:**
```bash
npx playwright test tests/e2e/navigation.spec.ts
```

---

### Task 13: Implement ModalComponent PageObject
**File:** `frontend/tests/pages/components/ModalComponent.ts`

```typescript
export class ModalComponent {
  constructor(private page: Page) {}
  
  // Locators
  get modal(): Locator
  get title(): Locator
  get okButton(): Locator
  get cancelButton(): Locator
  get closeButton(): Locator
  get content(): Locator
  
  // Actions
  async waitForOpen(): Promise<void>
  async clickOk(): Promise<void>
  async clickCancel(): Promise<void>
  async close(): Promise<void>
  
  // Assertions
  async expectVisible(): Promise<void>
  async expectTitle(title: string): Promise<void>
  async expectClosed(): Promise<void>
}
```

**Validation:**
```bash
npx tsc --noEmit frontend/tests/pages/components/ModalComponent.ts
```

---

### Task 14: Migrate Existing Tests to PageObject
**Files to update:**
- `tests/e2e/auth.spec.ts`
- `tests/e2e/contests.spec.ts`
- `tests/e2e/navigation.spec.ts`
- `tests/e2e/predictions.spec.ts`
- `tests/e2e/profile.spec.ts`
- `tests/e2e/teams.spec.ts`
- `tests/e2e/analytics.spec.ts`

**Example migration (auth.spec.ts):**
```typescript
// Before:
test('should login', async ({ page }) => {
  await page.goto('/login')
  await page.fill('input[type="email"]', 'test@example.com')
  await page.fill('input[type="password"]', 'password')
  await page.click('button:has-text("Login")')
})

// After:
test('should login', async ({ loginPage, contestsPage }) => {
  await loginPage.goto()
  await loginPage.login('test@example.com', 'password')
  await contestsPage.expectContestListVisible()
})
```

**Validation:**
```bash
npx playwright test --reporter=list
```

---

### Task 15: Add Comprehensive Button/Link Tests
**File:** `frontend/tests/e2e/all-interactions.spec.ts`

Тест для проверки ВСЕХ кликабельных элементов:

```typescript
import { test, expect } from '../fixtures/test-fixtures'

test.describe('All Clickable Elements', () => {
  test.describe('Navigation Links', () => {
    test('header logo navigates to home')
    test('contests link navigates to contests')
    test('predictions link navigates to predictions')
    test('teams link navigates to teams')
    test('sports link navigates to sports')
    test('analytics link navigates to analytics')
    test('profile link navigates to profile')
  })
  
  test.describe('Auth Buttons', () => {
    test('login button submits form')
    test('register button submits form')
    test('logout button logs out user')
    test('register link navigates to register')
    test('login link navigates to login')
  })
  
  test.describe('Contest Buttons', () => {
    test('create contest button opens modal')
    test('join button joins contest')
    test('leave button leaves contest')
    test('contest card opens details')
    test('pagination buttons work')
  })
  
  test.describe('Prediction Buttons', () => {
    test('submit prediction button submits')
    test('edit prediction button enables edit')
    test('cancel button cancels edit')
    test('event card expands details')
  })
  
  test.describe('Team Buttons', () => {
    test('create team button opens form')
    test('join team button submits code')
    test('invite button shows code')
    test('leave team button shows confirmation')
    test('remove member button removes')
  })
  
  test.describe('Profile Buttons', () => {
    test('save button saves profile')
    test('upload avatar button opens picker')
    test('privacy toggles work')
  })
  
  test.describe('Sports Admin Buttons', () => {
    test('add sport button opens form')
    test('add league button opens form')
    test('add team button opens form')
    test('add match button opens form')
    test('edit buttons open edit mode')
    test('delete buttons show confirmation')
  })
  
  test.describe('Modal Buttons', () => {
    test('OK button confirms action')
    test('Cancel button cancels action')
    test('Close (X) button closes modal')
  })
  
  test.describe('Common Buttons', () => {
    test('table pagination works')
    test('table sorting works')
    test('filters apply correctly')
    test('date pickers work')
    test('dropdowns work')
  })
})
```

**Validation:**
```bash
npx playwright test tests/e2e/all-interactions.spec.ts --reporter=html
npx playwright show-report
```

---

### Task 16: CI Integration & Final Validation
**File:** `.github/workflows/e2e-tests.yml` (update if exists)

```yaml
name: E2E Tests
on: [push, pull_request]
jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - name: Install dependencies
        run: cd frontend && npm ci
      - name: Install Playwright
        run: cd frontend && npx playwright install --with-deps
      - name: Run E2E tests
        run: cd frontend && npm run test:e2e
      - name: Upload report
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

**Final Validation:**
```bash
cd frontend
npm run test:e2e
npm run test:e2e:report
```

---

## Testing Strategy

### Test Levels
1. **Smoke Tests** — Critical paths (login, create contest, submit prediction)
2. **Interaction Tests** — All buttons and links
3. **Workflow Tests** — End-to-end user journeys
4. **Visual Tests** — Screenshot comparisons

### Coverage Matrix
| Page | Links | Buttons | Forms | Modals |
|------|-------|---------|-------|--------|
| Login | 2 | 1 | 1 | 0 |
| Register | 2 | 1 | 1 | 0 |
| Contests | 4 | 5 | 1 | 2 |
| Predictions | 3 | 4 | 2 | 1 |
| Profile | 2 | 3 | 2 | 0 |
| Teams | 3 | 6 | 2 | 3 |
| Sports | 4 | 8 | 4 | 4 |
| Analytics | 2 | 3 | 1 | 0 |
| **Total** | **22** | **31** | **14** | **10** |

---

## Acceptance Criteria

1. ✅ Все страницы имеют свой PageObject класс
2. ✅ Все кликабельные элементы протестированы (77 элементов)
3. ✅ Существующие тесты мигрированы на PageObject
4. ✅ Тесты проходят в headless режиме
5. ✅ Генерируется HTML отчёт
6. ✅ CI pipeline настроен

---

## Estimated Effort

| Phase | Tasks | Hours |
|-------|-------|-------|
| Phase 1: Infrastructure | 1-3 | 2-3h |
| Phase 2: Core Pages | 4-8 | 6-8h |
| Phase 3: Components | 9-10 | 2-3h |
| Phase 4: Tests | 11-15 | 6-8h |
| Phase 5: CI | 16 | 1-2h |
| **Total** | 16 | **17-24h** |

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Ant Design селекторы нестабильны | High | Использовать data-testid атрибуты |
| API недоступен в тестах | Medium | Мокировать API или использовать test DB |
| Flaky тесты | Medium | Добавить retry, waitFor conditions |
| Долгое выполнение | Low | Параллельный запуск, sharding |

# Code Review: Playwright PageObject Tests

**Commit:** `e388d27`
**Date:** 2026-02-03
**Reviewer:** Din (AI)

## Summary

Реализована PageObject архитектура для Playwright UI тестов. 16 файлов, 2133 строк кода.

## Files Reviewed

- `tests/pages/BasePage.ts`
- `tests/pages/LoginPage.ts`
- `tests/pages/RegisterPage.ts`
- `tests/pages/ContestsPage.ts`
- `tests/pages/PredictionsPage.ts`
- `tests/pages/ProfilePage.ts`
- `tests/pages/SportsPage.ts`
- `tests/pages/TeamsPage.ts`
- `tests/pages/AnalyticsPage.ts`
- `tests/pages/components/*.ts`
- `tests/fixtures/test-fixtures.ts`
- `tests/e2e/all-interactions.spec.ts`

---

## ✅ Positive Findings

### 1. Clean Architecture
```typescript
export abstract class BasePage {
  protected page: Page
  abstract readonly url: string
  // Common methods...
}
```
Правильное использование абстрактного базового класса.

### 2. DRY Helpers
```typescript
protected async clickButton(text: string): Promise<void>
protected async fillInput(selector: string, value: string): Promise<void>
protected async selectAntdOption(selector: string, optionText: string): Promise<void>
```
Ant Design хелперы вынесены в базовый класс.

### 3. Consistent API
Все PageObject имеют единообразный API:
- Locators (get properties)
- Actions (async methods)
- Assertions (expect methods)

### 4. Type Safety
```typescript
async navigateTo(pageName: 'contests' | 'predictions' | 'teams' | 'sports' | 'analytics'): Promise<void>
```
Использование union types для ограничения значений.

### 5. Proper Async/Await
Все асинхронные операции корректно используют await.

---

## ⚠️ Issues Found

### Issue 1: Flaky Selectors with `.first()` (Medium)

**Locations:**
- `ContestsPage.ts:29` — `this.page.locator('.ant-select').first()`
- `RegisterPage.ts:13` — `this.page.locator('input[placeholder*="name" i]').first()`
- `TeamsPage.ts:37` — `this.page.locator('input[placeholder*="team name" i]...').first()`

**Problem:** `.first()` может выбрать неправильный элемент если DOM изменится.

**Recommendation:** Использовать более специфичные селекторы или data-testid:
```typescript
// Before
get filterDropdown(): Locator {
  return this.page.locator('.ant-select').first()
}

// After
get filterDropdown(): Locator {
  return this.page.locator('[data-testid="contest-filter"]')
}
```

**Severity:** Medium (flaky tests risk)

---

### Issue 2: Hardcoded Timeouts (Low)

**Locations:** Multiple files with `timeout: 5000`, `timeout: 10000`

**Problem:** Timeouts разбросаны по коду вместо централизованной конфигурации.

**Recommendation:** Использовать константы из test-config:
```typescript
// test-config.ts
export const TIMEOUTS = {
  SHORT: 3000,
  MEDIUM: 5000,
  LONG: 10000,
  NETWORK: 30000,
}

// Usage
await expect(element).toBeVisible({ timeout: TIMEOUTS.MEDIUM })
```

**Severity:** Low (maintainability)

---

### Issue 3: Missing Error Messages in Assertions (Low)

**Location:** All PageObject files

**Problem:** Assertions не содержат кастомных сообщений об ошибках.

**Recommendation:**
```typescript
// Before
async expectOnLoginPage(): Promise<void> {
  await expect(this.page).toHaveURL('/login')
}

// After
async expectOnLoginPage(): Promise<void> {
  await expect(this.page, 'Should be on login page').toHaveURL('/login')
}
```

**Severity:** Low (debugging experience)

---

### Issue 4: Incomplete Test Coverage (Info)

**Problem:** `all-interactions.spec.ts` не покрывает все элементы из плана:
- Predictions submit/edit
- Teams create/join/leave
- Sports CRUD operations
- Analytics export

**Recommendation:** Добавить тесты для полного покрытия.

**Severity:** Info (incomplete feature)

---

### Issue 5: No Retry Logic for Flaky Operations (Low)

**Location:** `BasePage.ts`

**Problem:** Нет retry логики для потенциально flaky операций.

**Recommendation:**
```typescript
protected async clickWithRetry(locator: Locator, maxRetries = 3): Promise<void> {
  for (let i = 0; i < maxRetries; i++) {
    try {
      await locator.click()
      return
    } catch (e) {
      if (i === maxRetries - 1) throw e
      await this.page.waitForTimeout(500)
    }
  }
}
```

**Severity:** Low (stability)

---

### Issue 6: Test Config Uses Production Credentials (Medium)

**Location:** `test-config.ts`

```typescript
testUser: {
  email: 'admin@example.com',
  password: 'admin123',
}
```

**Problem:** Тесты используют admin credentials. Нужен отдельный test user.

**Recommendation:** Создать тестового пользователя или использовать env variables:
```typescript
testUser: {
  email: process.env.TEST_USER_EMAIL || 'testuser@example.com',
  password: process.env.TEST_USER_PASSWORD || 'testpass123',
}
```

**Severity:** Medium (security/isolation)

---

## 🔒 Security Check

- ✅ No secrets in code
- ⚠️ Admin credentials used for tests (Issue 6)
- ✅ No sensitive data exposed
- ✅ Environment variables supported

---

## 📊 Summary

| Category | Status |
|----------|--------|
| Architecture | ✅ Excellent |
| Type Safety | ✅ Good |
| Code Quality | ✅ Good |
| Test Coverage | ⚠️ Incomplete |
| Stability | ⚠️ Flaky risk |
| Security | ⚠️ Uses admin creds |

**Overall:** Good implementation. Minor issues to address.

---

## Action Items

1. **[Medium]** Add data-testid attributes to React components for stable selectors
2. **[Medium]** Create dedicated test user instead of using admin
3. **[Low]** Centralize timeout constants
4. **[Low]** Add custom error messages to assertions
5. **[Info]** Complete test coverage for all interactions

**Approval:** ✅ Approved with minor fixes recommended

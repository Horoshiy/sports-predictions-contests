import { test, expect } from '@playwright/test'
import { SELECTORS } from '../helpers/selectors'
import { TEST_CONFIG, TIMEOUTS } from '../helpers/test-config'

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    // Clear storage before each test
    await page.evaluate(() => localStorage.clear())
  })

  test('should display login page', async ({ page }) => {
    await page.goto('/login')
    await expect(page).toHaveURL('/login')
    await expect(page.locator('h1')).toContainText('Sports Prediction Contests')
    await expect(page.locator(SELECTORS.auth.emailInput)).toBeVisible()
    await expect(page.locator(SELECTORS.auth.passwordInput)).toBeVisible()
    await expect(page.locator(SELECTORS.auth.loginButton)).toBeVisible()
  })

  test('should login with valid credentials', async ({ page }) => {
    await page.goto('/login')
    
    await page.fill(SELECTORS.auth.emailInput, TEST_CONFIG.testUser.email)
    await page.fill(SELECTORS.auth.passwordInput, TEST_CONFIG.testUser.password)
    await page.click(SELECTORS.auth.loginButton)
    
    // Should redirect to contests page
    await expect(page).toHaveURL('/contests', { timeout: TIMEOUTS.MEDIUM })
  })

  test('should show error with invalid credentials', async ({ page }) => {
    await page.goto('/login')
    
    await page.fill(SELECTORS.auth.emailInput, 'wrong@example.com')
    await page.fill(SELECTORS.auth.passwordInput, 'wrongpassword')
    await page.click(SELECTORS.auth.loginButton)
    
    // Should show error notification
    await expect(page.locator(SELECTORS.common.notificationError)).toBeVisible({ timeout: TIMEOUTS.SHORT })
  })

  test('should display register page', async ({ page }) => {
    await page.goto('/register')
    await expect(page).toHaveURL('/register')
    await expect(page.locator(SELECTORS.auth.emailInput)).toBeVisible()
    await expect(page.locator(SELECTORS.auth.passwordInput)).toBeVisible()
    await expect(page.locator(SELECTORS.auth.registerButton)).toBeVisible()
  })

  test('should validate email format on login', async ({ page }) => {
    await page.goto('/login')
    
    await page.fill(SELECTORS.auth.emailInput, 'invalid-email')
    await page.fill(SELECTORS.auth.passwordInput, 'password123')
    await page.click(SELECTORS.auth.loginButton)
    
    // Should show validation error
    const emailInput = page.locator(SELECTORS.auth.emailInput)
    await expect(emailInput).toHaveAttribute('type', 'email')
  })

  test('should redirect to login when accessing protected route', async ({ page }) => {
    await page.goto('/contests')
    
    // Should redirect to login
    await expect(page).toHaveURL('/login', { timeout: TIMEOUTS.SHORT })
  })

  test('should logout successfully', async ({ page }) => {
    // First login
    await page.goto('/login')
    
    await page.fill(SELECTORS.auth.emailInput, TEST_CONFIG.testUser.email)
    await page.fill(SELECTORS.auth.passwordInput, TEST_CONFIG.testUser.password)
    await page.click(SELECTORS.auth.loginButton)
    await page.waitForURL('/contests')
    
    // Then logout
    await page.click('button:has-text("Logout")')
    
    // Should redirect to login
    await expect(page).toHaveURL('/login', { timeout: TIMEOUTS.SHORT })
  })

  test('should persist authentication after page reload', async ({ page }) => {
    // Login
    await page.goto('/login')
    
    await page.fill(SELECTORS.auth.emailInput, TEST_CONFIG.testUser.email)
    await page.fill(SELECTORS.auth.passwordInput, TEST_CONFIG.testUser.password)
    await page.click(SELECTORS.auth.loginButton)
    await page.waitForURL('/contests')
    
    // Reload page
    await page.reload()
    
    // Should still be on contests page
    await expect(page).toHaveURL('/contests')
  })
})

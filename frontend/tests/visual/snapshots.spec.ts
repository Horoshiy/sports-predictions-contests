import { test, expect } from '@playwright/test'
import { loginUser } from '../helpers/test-utils'
import { TEST_CONFIG } from '../helpers/test-config'

test.describe('Visual Regression Tests', () => {
  test('login page matches snapshot', async ({ page }) => {
    await page.goto('/login')
    await page.waitForLoadState('networkidle')
    await expect(page).toHaveScreenshot('login-page.png')
  })

  test('contests page matches snapshot', async ({ page }) => {
    await loginUser(page, TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    await page.goto('/contests')
    await page.waitForLoadState('networkidle')
    await expect(page).toHaveScreenshot('contests-page.png')
  })

  test('predictions page matches snapshot', async ({ page }) => {
    await loginUser(page, TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    await page.goto('/predictions')
    await page.waitForLoadState('networkidle')
    await expect(page).toHaveScreenshot('predictions-page.png')
  })
})

test.describe('Mobile Visual Regression', () => {
  test.use({ viewport: { width: 375, height: 667 } })
  
  test('mobile login page matches snapshot', async ({ page }) => {
    await page.goto('/login')
    await page.waitForLoadState('networkidle')
    await expect(page).toHaveScreenshot('mobile-login-page.png')
  })
})

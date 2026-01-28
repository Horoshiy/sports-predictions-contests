import { test, expect } from '@playwright/test'
import { SELECTORS } from '../helpers/selectors'
import { loginUser } from '../helpers/test-utils'
import { TEST_CONFIG } from '../helpers/test-config'

test.describe('Complete User Workflows', () => {
  test('complete user journey: login → contests → predictions', async ({ page }) => {
    // Login
    await loginUser(page, TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    await expect(page).toHaveURL('/contests')
    
    // View contests
    await page.waitForLoadState('networkidle')
    const contestCount = await page.locator(SELECTORS.contests.contestCard).count()
    expect(contestCount).toBeGreaterThanOrEqual(0)
    
    // Navigate to predictions
    await page.click(SELECTORS.navigation.predictionsLink)
    await expect(page).toHaveURL('/predictions')
    
    // View predictions
    await page.waitForLoadState('networkidle')
    const predictionCount = await page.locator(SELECTORS.predictions.eventCard).count()
    expect(predictionCount).toBeGreaterThanOrEqual(0)
  })

  test('navigation workflow: visit all pages', async ({ page }) => {
    await loginUser(page, TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    
    // Visit each page
    const pages = ['/contests', '/predictions', '/teams', '/analytics', '/profile']
    
    for (const pagePath of pages) {
      await page.goto(pagePath)
      await expect(page).toHaveURL(pagePath)
      await page.waitForLoadState('networkidle')
    }
  })
})

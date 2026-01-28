import { test, expect } from '../fixtures/auth.fixture'
import { SELECTORS } from '../helpers/selectors'

test.describe('Navigation', () => {
  test('should navigate to all main pages', async ({ authenticatedPage }) => {
    // Contests
    await authenticatedPage.click(SELECTORS.navigation.contestsLink)
    await expect(authenticatedPage).toHaveURL('/contests')
    
    // Predictions
    await authenticatedPage.click(SELECTORS.navigation.predictionsLink)
    await expect(authenticatedPage).toHaveURL('/predictions')
    
    // Teams
    await authenticatedPage.click(SELECTORS.navigation.teamsLink)
    await expect(authenticatedPage).toHaveURL('/teams')
    
    // Analytics
    await authenticatedPage.click(SELECTORS.navigation.analyticsLink)
    await expect(authenticatedPage).toHaveURL('/analytics')
    
    // Profile
    await authenticatedPage.click(SELECTORS.navigation.profileLink)
    await expect(authenticatedPage).toHaveURL('/profile')
  })

  test('should display navigation menu', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/contests')
    await expect(authenticatedPage.locator(SELECTORS.navigation.header)).toBeVisible()
    await expect(authenticatedPage.locator(SELECTORS.navigation.menu)).toBeVisible()
  })
})

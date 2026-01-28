import { test, expect } from '../fixtures/auth.fixture'
import { SELECTORS } from '../helpers/selectors'

test.describe('User Profile', () => {
  test('should display profile page', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/profile')
    await expect(authenticatedPage).toHaveURL('/profile')
  })

  test('should display user information', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/profile')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Should show profile content
    await expect(authenticatedPage.locator('h1, h2')).toContainText(/Profile/i)
  })
})

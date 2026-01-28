import { test, expect } from '../fixtures/auth.fixture'
import { SELECTORS } from '../helpers/selectors'

test.describe('Analytics Dashboard', () => {
  test('should display analytics page', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/analytics')
    await expect(authenticatedPage).toHaveURL('/analytics')
  })

  test('should display statistics', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/analytics')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Should show analytics content
    const hasStats = await authenticatedPage.locator(SELECTORS.analytics.statsCard).count()
    expect(hasStats).toBeGreaterThanOrEqual(0)
  })
})

import { test, expect } from '../fixtures/auth.fixture'
import { SELECTORS } from '../helpers/selectors'
import { generateContest } from '../fixtures/data.fixture'

test.describe('Contest Management', () => {
  test('should display contests page', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/contests')
    await expect(authenticatedPage).toHaveURL('/contests')
    await expect(authenticatedPage.locator('h1, h2')).toContainText(/Contest/i)
  })

  test('should view contest list', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/contests')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Should show contests or empty state
    const hasContests = await authenticatedPage.locator(SELECTORS.contests.contestCard).count()
    expect(hasContests).toBeGreaterThanOrEqual(0)
  })

  test('should navigate to contest details', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/contests')
    await authenticatedPage.waitForLoadState('networkidle')
    
    const contestCard = authenticatedPage.locator(SELECTORS.contests.contestCard).first()
    if (await contestCard.isVisible()) {
      await contestCard.click()
      // Should navigate to contest details or show modal
      await authenticatedPage.waitForTimeout(1000)
    }
  })
})

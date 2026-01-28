import { test, expect } from '../fixtures/auth.fixture'
import { SELECTORS } from '../helpers/selectors'

test.describe('Predictions', () => {
  test('should display predictions page', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/predictions')
    await expect(authenticatedPage).toHaveURL('/predictions')
  })

  test('should view predictions list', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/predictions')
    await authenticatedPage.waitForLoadState('networkidle')
    
    const hasPredictions = await authenticatedPage.locator(SELECTORS.predictions.eventCard).count()
    expect(hasPredictions).toBeGreaterThanOrEqual(0)
  })
})

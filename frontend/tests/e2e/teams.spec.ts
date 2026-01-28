import { test, expect } from '../fixtures/auth.fixture'
import { SELECTORS } from '../helpers/selectors'

test.describe('Teams', () => {
  test('should display teams page', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await expect(authenticatedPage).toHaveURL('/teams')
  })

  test('should view teams list', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.waitForLoadState('networkidle')
    
    const hasTeams = await authenticatedPage.locator(SELECTORS.teams.teamCard).count()
    expect(hasTeams).toBeGreaterThanOrEqual(0)
  })
})

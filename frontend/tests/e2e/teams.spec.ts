import { test, expect } from '../fixtures/auth.fixture'
import { SELECTORS } from '../helpers/selectors'

test.describe('Teams', () => {
  test('should display teams page', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await expect(authenticatedPage).toHaveURL('/teams')
    await expect(authenticatedPage.locator('[data-testid="teams-page"]')).toBeVisible()
  })

  test('should view teams list', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Check if table or empty state is visible
    const hasTable = await authenticatedPage.locator('.ant-table').isVisible().catch(() => false)
    const hasEmpty = await authenticatedPage.locator('.ant-empty').isVisible().catch(() => false)
    expect(hasTable || hasEmpty).toBeTruthy()
  })

  test('should create a new team', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Click create team button
    await authenticatedPage.locator('[data-testid="create-team-button"]').first().click()
    
    // Wait for modal
    await authenticatedPage.waitForSelector('.ant-modal', { state: 'visible' })
    
    // Fill form with unique team name
    const teamName = `Test Team ${Date.now()}`
    await authenticatedPage.fill('input[id="name"]', teamName)
    await authenticatedPage.fill('textarea[id="description"]', 'Test team description')
    
    // Submit form
    await authenticatedPage.locator('.ant-modal-footer button:has-text("Create")').click()
    
    // Wait for success notification
    await authenticatedPage.waitForSelector('.ant-notification-notice-success', { timeout: 5000 })
    
    // Wait for team to appear in list
    await authenticatedPage.waitForSelector(`text=${teamName}`, { timeout: 5000 }).catch(() => {})
    const teamExists = await authenticatedPage.locator(`text=${teamName}`).isVisible().catch(() => false)
    expect(teamExists).toBeTruthy()
  })

  test('should view team members', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Check if there are any teams
    const hasTeams = await authenticatedPage.locator('.ant-table-row').count() > 0
    
    if (hasTeams) {
      // Click view members button using data-testid
      await authenticatedPage.locator('[data-testid="view-members-button"]').first().click()
      
      // Wait for modal
      await authenticatedPage.waitForSelector('[data-testid="team-details-modal"]', { state: 'visible', timeout: 5000 })
      
      // Verify members tab is visible
      const membersTab = authenticatedPage.locator('div[role="tab"]:has-text("Members")')
      await expect(membersTab).toBeVisible()
    }
  })

  test('should display team leaderboard in contest', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/contests')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Check if there are contests
    const hasContests = await authenticatedPage.locator('.ant-table-row').count() > 0
    
    if (hasContests) {
      // Click on first contest to select it
      await authenticatedPage.locator('.ant-table-row').first().click()
      
      // Switch to Team Leaderboard tab
      await authenticatedPage.locator('div[role="tab"]:has-text("Team Leaderboard")').click()
      
      // Verify team leaderboard is displayed
      await expect(authenticatedPage.locator('[data-testid="team-leaderboard"]')).toBeVisible()
    }
  })

  test('should show empty state when no teams', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Switch to "My Teams" tab which might be empty
    await authenticatedPage.locator('div[role="tab"]:has-text("My Teams")').click()
    await authenticatedPage.waitForTimeout(500)
    
    // Check for either table or empty state
    const hasContent = await authenticatedPage.locator('.ant-table, .ant-empty').isVisible()
    expect(hasContent).toBeTruthy()
  })

  test('should navigate between tabs', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Click My Teams tab
    await authenticatedPage.locator('div[role="tab"]:has-text("My Teams")').click()
    await expect(authenticatedPage.locator('div[role="tab"]:has-text("My Teams")').locator('..')).toHaveClass(/ant-tabs-tab-active/)
    
    // Click All Teams tab
    await authenticatedPage.locator('div[role="tab"]:has-text("All Teams")').click()
    await expect(authenticatedPage.locator('div[role="tab"]:has-text("All Teams")').locator('..')).toHaveClass(/ant-tabs-tab-active/)
    
    // Click Join Team tab
    await authenticatedPage.locator('div[role="tab"]:has-text("Join Team")').click()
    await expect(authenticatedPage.locator('div[role="tab"]:has-text("Join Team")').locator('..')).toHaveClass(/ant-tabs-tab-active/)
    
    // Verify join form is visible
    await expect(authenticatedPage.locator('[data-testid="join-team-form"]')).toBeVisible()
  })

  test('should validate empty team name', async ({ authenticatedPage }) => {
    await authenticatedPage.goto('/teams')
    await authenticatedPage.waitForLoadState('networkidle')
    
    // Click create team button
    await authenticatedPage.locator('[data-testid="create-team-button"]').first().click()
    
    // Wait for modal
    await authenticatedPage.waitForSelector('.ant-modal', { state: 'visible' })
    
    // Try to submit without filling name
    await authenticatedPage.locator('.ant-modal-footer button:has-text("Create")').click()
    
    // Verify validation error appears
    const hasError = await authenticatedPage.locator('.ant-form-item-explain-error').isVisible()
    expect(hasError).toBeTruthy()
  })
})

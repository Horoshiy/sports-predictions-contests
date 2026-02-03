import { test, expect } from '../fixtures/test-fixtures'
import { TEST_CONFIG } from '../helpers/test-config'

/**
 * Comprehensive tests for ALL clickable elements
 * Tests every button, link, and interactive element in the application
 */

test.describe('All Clickable Elements', () => {
  
  // ==================== Navigation Links ====================
  
  test.describe('Navigation Links', () => {
    test.beforeEach(async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    })

    test('contests link navigates to contests', async ({ header, contestsPage }) => {
      await header.goToContests()
      await contestsPage.expectOnContestsPage()
    })

    test('predictions link navigates to predictions', async ({ header, predictionsPage }) => {
      await header.goToPredictions()
      await predictionsPage.expectOnPredictionsPage()
    })

    test('sports link navigates to sports', async ({ header, sportsPage }) => {
      await header.goToSports()
      await sportsPage.expectOnSportsPage()
    })

    test('analytics link navigates to analytics', async ({ header, analyticsPage }) => {
      await header.goToAnalytics()
      await analyticsPage.expectOnAnalyticsPage()
    })

    test('profile link navigates to profile', async ({ header, profilePage }) => {
      await header.goToProfile()
      await profilePage.expectOnProfilePage()
    })

    test('user menu opens on click', async ({ header }) => {
      await header.openUserMenu()
      await expect(header.logoutButton).toBeVisible()
    })
  })

  // ==================== Auth Buttons ====================

  test.describe('Auth Buttons', () => {
    test('login button submits form', async ({ loginPage, contestsPage }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
      await contestsPage.expectOnContestsPage()
    })

    test('register link navigates to register', async ({ loginPage, registerPage }) => {
      await loginPage.goto()
      await loginPage.clickRegisterLink()
      await registerPage.expectOnRegisterPage()
    })

    test('login link from register navigates to login', async ({ registerPage, loginPage }) => {
      await registerPage.goto()
      await registerPage.clickLoginLink()
      await loginPage.expectOnLoginPage()
    })

    test('logout button logs out user', async ({ loginPage, header }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
      await header.logout()
      await header.expectLoggedOut()
    })
  })

  // ==================== Contest Buttons ====================

  test.describe('Contest Buttons', () => {
    test.beforeEach(async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    })

    test('create contest button opens modal', async ({ contestsPage }) => {
      await contestsPage.goto()
      await contestsPage.clickCreateContest()
      await contestsPage.expectCreateModalVisible()
    })

    test('contest card is clickable', async ({ contestsPage, page }) => {
      await contestsPage.goto()
      await contestsPage.expectContestListVisible()
      const firstCard = contestsPage.getContestCard(0)
      await expect(firstCard).toBeVisible()
      // Card should be clickable
      await firstCard.click()
    })
  })

  // ==================== Profile Buttons ====================

  test.describe('Profile Buttons', () => {
    test.beforeEach(async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    })

    test('save button is clickable', async ({ profilePage }) => {
      await profilePage.goto()
      await profilePage.expectProfileLoaded()
      await expect(profilePage.saveButton).toBeVisible()
    })

    test('display name input accepts text', async ({ profilePage }) => {
      await profilePage.goto()
      await profilePage.expectProfileLoaded()
      await profilePage.updateDisplayName('Test Name')
      await profilePage.expectDisplayName('Test Name')
    })
  })

  // ==================== Sports Admin Buttons ====================

  test.describe('Sports Admin Buttons', () => {
    test.beforeEach(async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testAdmin.email, TEST_CONFIG.testAdmin.password)
    })

    test('add sport button is visible for admin', async ({ sportsPage }) => {
      await sportsPage.goto()
      await expect(sportsPage.addSportButton).toBeVisible()
    })

    test('add sport button opens modal', async ({ sportsPage }) => {
      await sportsPage.goto()
      await sportsPage.clickAddSport()
      await sportsPage.expectModalVisible()
    })
  })

  // ==================== Modal Buttons ====================

  test.describe('Modal Buttons', () => {
    test.beforeEach(async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    })

    test('cancel button closes modal', async ({ contestsPage, page }) => {
      await contestsPage.goto()
      await contestsPage.clickCreateContest()
      await contestsPage.expectCreateModalVisible()
      await page.click('.ant-modal-footer button:not(.ant-btn-primary)')
      await expect(page.locator('.ant-modal')).not.toBeVisible()
    })

    test('close (X) button closes modal', async ({ contestsPage, page }) => {
      await contestsPage.goto()
      await contestsPage.clickCreateContest()
      await contestsPage.expectCreateModalVisible()
      await page.click('.ant-modal-close')
      await expect(page.locator('.ant-modal')).not.toBeVisible()
    })
  })

  // ==================== Form Inputs ====================

  test.describe('Form Inputs', () => {
    test('login email input accepts email', async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.emailInput.fill('test@example.com')
      await expect(loginPage.emailInput).toHaveValue('test@example.com')
    })

    test('login password input accepts password', async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.passwordInput.fill('password123')
      await expect(loginPage.passwordInput).toHaveValue('password123')
    })
  })

  // ==================== Protected Routes ====================

  test.describe('Protected Routes', () => {
    test('contests page redirects to login when not authenticated', async ({ contestsPage, page }) => {
      // Navigate first, then clear to ensure no auth, then navigate again
      await page.goto('/login')
      await page.evaluate(() => localStorage.clear())
      await contestsPage.goto()
      await expect(page).toHaveURL('/login')
    })

    test('predictions page redirects to login when not authenticated', async ({ predictionsPage, page }) => {
      await page.goto('/login')
      await page.evaluate(() => localStorage.clear())
      await predictionsPage.goto()
      await expect(page).toHaveURL('/login')
    })

    test('profile page redirects to login when not authenticated', async ({ profilePage, page }) => {
      await page.goto('/login')
      await page.evaluate(() => localStorage.clear())
      await profilePage.goto()
      await expect(page).toHaveURL('/login')
    })
  })

  // ==================== Error States ====================

  test.describe('Error States', () => {
    test('login with invalid credentials shows error', async ({ loginPage }) => {
      await loginPage.goto()
      await loginPage.login('wrong@email.com', 'wrongpassword')
      await loginPage.expectErrorVisible()
    })
  })

  // ==================== Session Persistence ====================

  test.describe('Session Persistence', () => {
    test('user stays logged in after page reload', async ({ loginPage, contestsPage, page }) => {
      await loginPage.goto()
      await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
      await contestsPage.expectOnContestsPage()
      
      await page.reload()
      await contestsPage.expectOnContestsPage()
    })
  })
})

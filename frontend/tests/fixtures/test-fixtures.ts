import { test as base, expect, Page } from '@playwright/test'
import { LoginPage } from '../pages/LoginPage'
import { RegisterPage } from '../pages/RegisterPage'
import { ContestsPage } from '../pages/ContestsPage'
import { PredictionsPage } from '../pages/PredictionsPage'
import { ProfilePage } from '../pages/ProfilePage'
import { SportsPage } from '../pages/SportsPage'
import { TeamsPage } from '../pages/TeamsPage'
import { AnalyticsPage } from '../pages/AnalyticsPage'
import { HeaderComponent } from '../pages/components/HeaderComponent'
import { TEST_CONFIG } from '../helpers/test-config'

/**
 * Extended test fixtures with PageObject instances
 */
type PageFixtures = {
  loginPage: LoginPage
  registerPage: RegisterPage
  contestsPage: ContestsPage
  predictionsPage: PredictionsPage
  profilePage: ProfilePage
  sportsPage: SportsPage
  teamsPage: TeamsPage
  analyticsPage: AnalyticsPage
  header: HeaderComponent
  authenticatedPage: Page
}

export const test = base.extend<PageFixtures>({
  // Page Objects
  loginPage: async ({ page }, use) => {
    await use(new LoginPage(page))
  },

  registerPage: async ({ page }, use) => {
    await use(new RegisterPage(page))
  },

  contestsPage: async ({ page }, use) => {
    await use(new ContestsPage(page))
  },

  predictionsPage: async ({ page }, use) => {
    await use(new PredictionsPage(page))
  },

  profilePage: async ({ page }, use) => {
    await use(new ProfilePage(page))
  },

  sportsPage: async ({ page }, use) => {
    await use(new SportsPage(page))
  },

  teamsPage: async ({ page }, use) => {
    await use(new TeamsPage(page))
  },

  analyticsPage: async ({ page }, use) => {
    await use(new AnalyticsPage(page))
  },

  header: async ({ page }, use) => {
    await use(new HeaderComponent(page))
  },

  // Pre-authenticated page fixture
  authenticatedPage: async ({ page }, use) => {
    const loginPage = new LoginPage(page)
    await loginPage.goto()
    await loginPage.login(TEST_CONFIG.testUser.email, TEST_CONFIG.testUser.password)
    await page.waitForURL('/contests')
    await use(page)
  },
})

export { expect }

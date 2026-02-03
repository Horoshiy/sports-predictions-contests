import { Page, Locator, expect } from '@playwright/test'

/**
 * Header Component Page Object
 */
export class HeaderComponent {
  private page: Page

  constructor(page: Page) {
    this.page = page
  }

  // ==================== Locators ====================

  get logo(): Locator {
    return this.page.locator('header .logo, header a:first-child')
  }

  get contestsLink(): Locator {
    return this.page.locator('a[href="/contests"]')
  }

  get predictionsLink(): Locator {
    return this.page.locator('a[href="/predictions"]')
  }

  get teamsLink(): Locator {
    return this.page.locator('a[href="/teams"]')
  }

  get sportsLink(): Locator {
    return this.page.locator('a[href="/sports"]')
  }

  get analyticsLink(): Locator {
    return this.page.locator('a[href="/analytics"]')
  }

  get profileLink(): Locator {
    return this.page.locator('a[href="/profile"]')
  }

  get userMenu(): Locator {
    return this.page.locator('.ant-dropdown-trigger, .ant-avatar')
  }

  get logoutButton(): Locator {
    return this.page.locator('text=Logout')
  }

  get header(): Locator {
    return this.page.locator('header')
  }

  get menuItems(): Locator {
    return this.page.locator('.ant-menu-item')
  }

  // ==================== Actions ====================

  /**
   * Click logo to go home
   */
  async clickLogo(): Promise<void> {
    await this.logo.click()
  }

  /**
   * Navigate to a page
   */
  async navigateTo(pageName: 'contests' | 'predictions' | 'teams' | 'sports' | 'analytics'): Promise<void> {
    const links = {
      contests: this.contestsLink,
      predictions: this.predictionsLink,
      teams: this.teamsLink,
      sports: this.sportsLink,
      analytics: this.analyticsLink,
    }
    await links[pageName].click()
  }

  /**
   * Navigate to contests
   */
  async goToContests(): Promise<void> {
    await this.contestsLink.click()
  }

  /**
   * Navigate to predictions
   */
  async goToPredictions(): Promise<void> {
    await this.predictionsLink.click()
  }

  /**
   * Navigate to teams
   */
  async goToTeams(): Promise<void> {
    await this.teamsLink.click()
  }

  /**
   * Navigate to sports
   */
  async goToSports(): Promise<void> {
    await this.sportsLink.click()
  }

  /**
   * Navigate to analytics
   */
  async goToAnalytics(): Promise<void> {
    await this.analyticsLink.click()
  }

  /**
   * Open user menu dropdown
   */
  async openUserMenu(): Promise<void> {
    await this.userMenu.click()
  }

  /**
   * Go to profile (via user menu)
   */
  async goToProfile(): Promise<void> {
    await this.openUserMenu()
    await this.profileLink.click()
  }

  /**
   * Logout
   */
  async logout(): Promise<void> {
    await this.openUserMenu()
    await this.logoutButton.click()
  }

  // ==================== Assertions ====================

  /**
   * Expect user to be logged in
   */
  async expectUserLoggedIn(userName?: string): Promise<void> {
    await expect(this.userMenu).toBeVisible()
    if (userName) {
      await expect(this.page.locator(`text=Welcome, ${userName}`)).toBeVisible()
    }
  }

  /**
   * Expect active menu item
   */
  async expectActiveMenuItem(item: string): Promise<void> {
    await expect(this.page.locator(`.ant-menu-item-selected:has-text("${item}")`)).toBeVisible()
  }

  /**
   * Expect header visible
   */
  async expectHeaderVisible(): Promise<void> {
    await expect(this.header).toBeVisible()
  }

  /**
   * Expect navigation links visible
   */
  async expectNavigationVisible(): Promise<void> {
    await expect(this.contestsLink).toBeVisible()
    await expect(this.predictionsLink).toBeVisible()
  }

  /**
   * Expect logout success (redirected to login)
   */
  async expectLoggedOut(): Promise<void> {
    await expect(this.page).toHaveURL('/login')
  }
}

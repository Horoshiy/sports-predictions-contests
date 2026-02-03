import { Page, Locator, expect } from '@playwright/test'
import { TIMEOUTS } from '../../helpers/test-config'

/**
 * Header Component Object
 * Handles navigation menu and user dropdown
 */
export class HeaderComponent {
  private readonly page: Page

  // Header selectors
  private readonly header = 'header'
  private readonly logo = 'header div:has-text("Sports Prediction Contests")'
  private readonly navigationMenu = 'header .ant-menu'
  private readonly userDropdown = 'header .ant-dropdown-trigger, header .ant-avatar'
  private readonly welcomeText = 'header span:has-text("Welcome")'

  // Menu items
  private readonly menuItems = {
    contests: '.ant-menu-item:has-text("Contests")',
    teams: '.ant-menu-item:has-text("Teams")',
    predictions: '.ant-menu-item:has-text("Predictions")',
    sports: '.ant-menu-item:has-text("Sports")',
    analytics: '.ant-menu-item:has-text("Analytics")',
  }

  // Dropdown menu items
  private readonly dropdownMenu = '.ant-dropdown'
  private readonly profileMenuItem = '.ant-dropdown-menu-item:has-text("Profile")'
  private readonly logoutMenuItem = '.ant-dropdown-menu-item:has-text("Logout")'

  constructor(page: Page) {
    this.page = page
  }

  // ==================== Navigation ====================

  /**
   * Navigate to Contests page
   */
  async goToContests(): Promise<void> {
    await this.page.click(this.menuItems.contests)
    await this.page.waitForURL('**/contests')
  }

  /**
   * Navigate to Teams page
   */
  async goToTeams(): Promise<void> {
    await this.page.click(this.menuItems.teams)
    await this.page.waitForURL('**/teams')
  }

  /**
   * Navigate to Predictions page
   */
  async goToPredictions(): Promise<void> {
    await this.page.click(this.menuItems.predictions)
    await this.page.waitForURL('**/predictions')
  }

  /**
   * Navigate to Sports page
   */
  async goToSports(): Promise<void> {
    await this.page.click(this.menuItems.sports)
    await this.page.waitForURL('**/sports')
  }

  /**
   * Navigate to Analytics page
   */
  async goToAnalytics(): Promise<void> {
    await this.page.click(this.menuItems.analytics)
    await this.page.waitForURL('**/analytics')
  }

  /**
   * Click logo to go home
   */
  async clickLogo(): Promise<void> {
    await this.page.click(this.logo)
  }

  // ==================== User Menu ====================

  /**
   * Open user dropdown menu
   */
  async openUserMenu(): Promise<void> {
    await this.page.click(this.userDropdown)
    await this.page.waitForSelector(this.dropdownMenu, { state: 'visible' })
  }

  /**
   * Navigate to profile page via dropdown
   */
  async goToProfile(): Promise<void> {
    await this.openUserMenu()
    await this.page.click(this.profileMenuItem)
    await this.page.waitForURL('**/profile')
  }

  /**
   * Logout via dropdown
   */
  async logout(): Promise<void> {
    await this.openUserMenu()
    await this.page.click(this.logoutMenuItem)
    await this.page.waitForURL('**/login')
  }

  // ==================== State Checks ====================

  /**
   * Check if user is logged in (header shows user menu)
   */
  async isLoggedIn(): Promise<boolean> {
    try {
      await this.page.waitForSelector(this.welcomeText, { timeout: TIMEOUTS.SHORT })
      return true
    } catch {
      return false
    }
  }

  /**
   * Get welcome text (includes username)
   */
  async getWelcomeText(): Promise<string | null> {
    const element = this.page.locator(this.welcomeText)
    if (await element.isVisible()) {
      return await element.textContent()
    }
    return null
  }

  /**
   * Get username from welcome text
   */
  async getUsername(): Promise<string | null> {
    const welcomeText = await this.getWelcomeText()
    if (welcomeText) {
      // Extract name from "Welcome, <name>"
      const match = welcomeText.match(/Welcome,\s*(.+)/)
      return match ? match[1].trim() : null
    }
    return null
  }

  /**
   * Check if navigation menu is visible
   */
  async isNavigationVisible(): Promise<boolean> {
    return await this.page.locator(this.navigationMenu).isVisible()
  }

  /**
   * Get active menu item
   */
  async getActiveMenuItem(): Promise<string | null> {
    const activeItem = this.page.locator('.ant-menu-item-selected')
    if (await activeItem.isVisible()) {
      return await activeItem.textContent()
    }
    return null
  }

  /**
   * Check if menu item is visible
   */
  async isMenuItemVisible(item: 'contests' | 'teams' | 'predictions' | 'sports' | 'analytics'): Promise<boolean> {
    return await this.page.locator(this.menuItems[item]).isVisible()
  }

  // ==================== Assertions ====================

  /**
   * Assert header is visible
   */
  async expectHeaderVisible(): Promise<void> {
    await expect(
      this.page.locator(this.header),
      'Expected header to be visible'
    ).toBeVisible()
  }

  /**
   * Assert user is logged in
   */
  async expectLoggedIn(): Promise<void> {
    await expect(
      this.page.locator(this.welcomeText),
      'Expected user to be logged in'
    ).toBeVisible()
  }

  /**
   * Assert user is logged out
   */
  async expectLoggedOut(): Promise<void> {
    await expect(
      this.page.locator(this.welcomeText),
      'Expected user to be logged out'
    ).toBeHidden()
  }

  /**
   * Assert menu item is active
   */
  async expectMenuItemActive(item: 'Contests' | 'Teams' | 'Predictions' | 'Sports' | 'Analytics'): Promise<void> {
    await expect(
      this.page.locator(`.ant-menu-item-selected:has-text("${item}")`),
      `Expected "${item}" menu item to be active`
    ).toBeVisible()
  }

  /**
   * Assert username matches
   */
  async expectUsername(expectedName: string): Promise<void> {
    const username = await this.getUsername()
    expect(username).toBe(expectedName)
  }

  /**
   * Assert navigation contains expected items
   */
  async expectNavigationItems(): Promise<void> {
    await expect(
      this.page.locator(this.menuItems.contests),
      'Expected Contests menu item'
    ).toBeVisible()
    await expect(
      this.page.locator(this.menuItems.predictions),
      'Expected Predictions menu item'
    ).toBeVisible()
    await expect(
      this.page.locator(this.menuItems.sports),
      'Expected Sports menu item'
    ).toBeVisible()
    await expect(
      this.page.locator(this.menuItems.analytics),
      'Expected Analytics menu item'
    ).toBeVisible()
  }
}

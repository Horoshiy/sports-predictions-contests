import { Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'

/**
 * Login Page Object
 */
export class LoginPage extends BasePage {
  readonly url = '/login'

  // ==================== Locators ====================

  get emailInput(): Locator {
    return this.page.locator('input[type="email"]')
  }

  get passwordInput(): Locator {
    return this.page.locator('input[type="password"]')
  }

  get loginButton(): Locator {
    return this.page.locator('button:has-text("Login")')
  }

  get registerLink(): Locator {
    return this.page.locator('a:has-text("Register")')
  }

  get errorMessage(): Locator {
    return this.page.locator('.ant-notification-notice-error')
  }

  get pageTitle(): Locator {
    return this.page.locator('h1')
  }

  // ==================== Actions ====================

  /**
   * Login with email and password
   */
  async login(email: string, password: string): Promise<void> {
    await this.emailInput.fill(email)
    await this.passwordInput.fill(password)
    await this.loginButton.click()
  }

  /**
   * Click register link to navigate to registration
   */
  async clickRegisterLink(): Promise<void> {
    await this.registerLink.click()
  }

  /**
   * Clear login form
   */
  async clearForm(): Promise<void> {
    await this.emailInput.clear()
    await this.passwordInput.clear()
  }

  // ==================== Assertions ====================

  /**
   * Expect login form to be visible
   */
  async expectLoginFormVisible(): Promise<void> {
    await expect(this.emailInput).toBeVisible()
    await expect(this.passwordInput).toBeVisible()
    await expect(this.loginButton).toBeVisible()
  }

  /**
   * Expect error notification to be visible
   */
  async expectErrorVisible(): Promise<void> {
    await expect(this.errorMessage).toBeVisible({ timeout: 5000 })
  }

  /**
   * Expect to be on login page
   */
  async expectOnLoginPage(): Promise<void> {
    await expect(this.page).toHaveURL('/login')
  }

  /**
   * Expect page title
   */
  async expectTitle(title: string): Promise<void> {
    await expect(this.pageTitle).toContainText(title)
  }
}

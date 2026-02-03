import { Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'

/**
 * Register Page Object
 */
export class RegisterPage extends BasePage {
  readonly url = '/register'

  // ==================== Locators ====================

  get nameInput(): Locator {
    return this.page.locator('input[placeholder*="name" i]').first()
  }

  get emailInput(): Locator {
    return this.page.locator('input[type="email"]')
  }

  get passwordInput(): Locator {
    return this.page.locator('input[type="password"]').first()
  }

  get confirmPasswordInput(): Locator {
    return this.page.locator('input[type="password"]').nth(1)
  }

  get registerButton(): Locator {
    return this.page.locator('button:has-text("Register")')
  }

  get loginLink(): Locator {
    return this.page.locator('a:has-text("Login")')
  }

  get errorMessage(): Locator {
    return this.page.locator('.ant-notification-notice-error')
  }

  get successMessage(): Locator {
    return this.page.locator('.ant-notification-notice-success')
  }

  // ==================== Actions ====================

  /**
   * Register a new user
   */
  async register(name: string, email: string, password: string, confirmPassword?: string): Promise<void> {
    await this.nameInput.fill(name)
    await this.emailInput.fill(email)
    await this.passwordInput.fill(password)
    if (this.confirmPasswordInput) {
      await this.confirmPasswordInput.fill(confirmPassword || password)
    }
    await this.registerButton.click()
  }

  /**
   * Click login link
   */
  async clickLoginLink(): Promise<void> {
    await this.loginLink.click()
  }

  // ==================== Assertions ====================

  /**
   * Expect register form to be visible
   */
  async expectRegisterFormVisible(): Promise<void> {
    await expect(this.emailInput).toBeVisible()
    await expect(this.passwordInput).toBeVisible()
    await expect(this.registerButton).toBeVisible()
  }

  /**
   * Expect password mismatch error
   */
  async expectPasswordMismatchError(): Promise<void> {
    await expect(this.errorMessage).toBeVisible({ timeout: 5000 })
  }

  /**
   * Expect registration success
   */
  async expectRegistrationSuccess(): Promise<void> {
    await expect(this.successMessage).toBeVisible({ timeout: 5000 })
  }

  /**
   * Expect to be on register page
   */
  async expectOnRegisterPage(): Promise<void> {
    await expect(this.page).toHaveURL('/register')
  }
}

import { Page, Locator, expect } from '@playwright/test'
import { TIMEOUTS } from '../../helpers/test-config'

/**
 * Notification Component Page Object (Ant Design Notification)
 */
export class NotificationComponent {
  private page: Page

  constructor(page: Page) {
    this.page = page
  }

  // ==================== Locators ====================

  get notification(): Locator {
    return this.page.locator('.ant-notification-notice')
  }

  get successNotification(): Locator {
    return this.page.locator('.ant-notification-notice-success')
  }

  get errorNotification(): Locator {
    return this.page.locator('.ant-notification-notice-error')
  }

  get infoNotification(): Locator {
    return this.page.locator('.ant-notification-notice-info')
  }

  get warningNotification(): Locator {
    return this.page.locator('.ant-notification-notice-warning')
  }

  get message(): Locator {
    return this.page.locator('.ant-notification-notice-message')
  }

  get description(): Locator {
    return this.page.locator('.ant-notification-notice-description')
  }

  get closeButton(): Locator {
    return this.page.locator('.ant-notification-notice-close')
  }

  // ==================== Actions ====================

  /**
   * Wait for notification
   */
  async waitFor(type: 'success' | 'error' | 'info' | 'warning' = 'success'): Promise<void> {
    const locators = {
      success: this.successNotification,
      error: this.errorNotification,
      info: this.infoNotification,
      warning: this.warningNotification,
    }
    await locators[type].waitFor({ state: 'visible', timeout: TIMEOUTS.MEDIUM })
  }

  /**
   * Close notification
   */
  async close(): Promise<void> {
    await this.closeButton.click()
  }

  /**
   * Get notification message text
   */
  async getMessageText(): Promise<string> {
    return await this.message.textContent() || ''
  }

  /**
   * Get notification description text
   */
  async getDescriptionText(): Promise<string> {
    return await this.description.textContent() || ''
  }

  // ==================== Assertions ====================

  /**
   * Expect success notification
   */
  async expectSuccess(message?: string): Promise<void> {
    await expect(
      this.successNotification,
      'Success notification should be visible'
    ).toBeVisible({ timeout: TIMEOUTS.MEDIUM })
    if (message) {
      await expect(this.message, `Notification should contain "${message}"`).toContainText(message)
    }
  }

  /**
   * Expect error notification
   */
  async expectError(message?: string): Promise<void> {
    await expect(
      this.errorNotification,
      'Error notification should be visible'
    ).toBeVisible({ timeout: TIMEOUTS.MEDIUM })
    if (message) {
      await expect(this.message, `Notification should contain "${message}"`).toContainText(message)
    }
  }

  /**
   * Expect info notification
   */
  async expectInfo(message?: string): Promise<void> {
    await expect(
      this.infoNotification,
      'Info notification should be visible'
    ).toBeVisible({ timeout: TIMEOUTS.MEDIUM })
    if (message) {
      await expect(this.message, `Notification should contain "${message}"`).toContainText(message)
    }
  }

  /**
   * Expect warning notification
   */
  async expectWarning(message?: string): Promise<void> {
    await expect(
      this.warningNotification,
      'Warning notification should be visible'
    ).toBeVisible({ timeout: TIMEOUTS.MEDIUM })
    if (message) {
      await expect(this.message, `Notification should contain "${message}"`).toContainText(message)
    }
  }

  /**
   * Expect notification closed
   */
  async expectClosed(): Promise<void> {
    await expect(this.notification).not.toBeVisible()
  }
}

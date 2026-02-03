import { Page, Locator, expect } from '@playwright/test'
import { TIMEOUTS } from '../helpers/test-config'

/**
 * Base Page Object class - all page objects extend this
 */
export abstract class BasePage {
  protected page: Page
  abstract readonly url: string
  abstract readonly pageName: string  // For better error messages

  constructor(page: Page) {
    this.page = page
  }

  /**
   * Navigate to this page
   */
  async goto(): Promise<void> {
    await this.page.goto(this.url)
    await this.waitForPageLoad()
  }

  /**
   * Wait for page to be fully loaded
   */
  async waitForPageLoad(): Promise<void> {
    await this.page.waitForLoadState('networkidle')
  }

  /**
   * Get the page title
   */
  async getTitle(): Promise<string> {
    return await this.page.title()
  }

  /**
   * Get current URL
   */
  async getCurrentUrl(): Promise<string> {
    return this.page.url()
  }

  /**
   * Check if element is visible
   */
  protected async isVisible(selector: string): Promise<boolean> {
    try {
      await this.page.waitForSelector(selector, { timeout: TIMEOUTS.SHORT })
      return true
    } catch {
      return false
    }
  }

  // ==================== Ant Design Helpers ====================

  /**
   * Click a button by text
   */
  protected async clickButton(text: string): Promise<void> {
    await this.page.click(`button:has-text("${text}")`)
  }

  /**
   * Click a link by text
   */
  protected async clickLink(text: string): Promise<void> {
    await this.page.click(`a:has-text("${text}")`)
  }

  /**
   * Fill an input field
   */
  protected async fillInput(selector: string, value: string): Promise<void> {
    await this.page.fill(selector, value)
  }

  /**
   * Fill Ant Design input by label
   */
  protected async fillAntdInputByLabel(label: string, value: string): Promise<void> {
    const input = this.page.locator(`label:has-text("${label}")`).locator('..').locator('input')
    await input.fill(value)
  }

  /**
   * Select option in Ant Design Select
   */
  protected async selectAntdOption(selector: string, optionText: string): Promise<void> {
    await this.page.click(selector)
    await this.page.click(`.ant-select-dropdown .ant-select-item:has-text("${optionText}")`)
  }

  /**
   * Select option by label
   */
  protected async selectAntdOptionByLabel(label: string, optionText: string): Promise<void> {
    const selectWrapper = this.page.locator(`label:has-text("${label}")`).locator('..').locator('.ant-select')
    await selectWrapper.click()
    await this.page.click(`.ant-select-dropdown .ant-select-item:has-text("${optionText}")`)
  }

  /**
   * Wait for Ant Design notification
   */
  protected async waitForNotification(type: 'success' | 'error' | 'info' | 'warning'): Promise<void> {
    await this.page.waitForSelector(`.ant-notification-notice-${type}`, { timeout: TIMEOUTS.MEDIUM })
  }

  /**
   * Expect notification to appear
   */
  protected async expectNotification(type: 'success' | 'error'): Promise<void> {
    await expect(
      this.page.locator(`.ant-notification-notice-${type}`),
      `Expected ${type} notification on ${this.pageName}`
    ).toBeVisible({ timeout: TIMEOUTS.MEDIUM })
  }

  /**
   * Click Ant Design modal OK button
   */
  protected async clickModalOk(): Promise<void> {
    await this.page.click('.ant-modal-footer button.ant-btn-primary')
  }

  /**
   * Click Ant Design modal Cancel button
   */
  protected async clickModalCancel(): Promise<void> {
    await this.page.click('.ant-modal-footer button:not(.ant-btn-primary)')
  }

  /**
   * Wait for modal to be visible
   */
  protected async waitForModal(): Promise<void> {
    await this.page.waitForSelector('.ant-modal', { state: 'visible' })
  }

  /**
   * Wait for modal to be hidden
   */
  protected async waitForModalClosed(): Promise<void> {
    await this.page.waitForSelector('.ant-modal', { state: 'hidden' })
  }

  /**
   * Get table row count
   */
  protected async getTableRowCount(tableSelector: string = '.ant-table'): Promise<number> {
    return await this.page.locator(`${tableSelector} tbody tr`).count()
  }

  /**
   * Wait for loading spinner to disappear
   */
  protected async waitForLoadingComplete(): Promise<void> {
    await this.page.waitForSelector('.ant-spin', { state: 'hidden', timeout: TIMEOUTS.LONG }).catch(() => {})
  }

  /**
   * Get locator with data-testid
   */
  protected getByTestId(testId: string): Locator {
    return this.page.locator(`[data-testid="${testId}"]`)
  }

  /**
   * Click with retry for flaky operations
   */
  protected async clickWithRetry(locator: Locator, maxRetries = 3): Promise<void> {
    for (let i = 0; i < maxRetries; i++) {
      try {
        await locator.click({ timeout: TIMEOUTS.SHORT })
        return
      } catch (e) {
        if (i === maxRetries - 1) throw e
        await this.page.waitForTimeout(500)
      }
    }
  }
}

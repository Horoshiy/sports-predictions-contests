import { Page, Locator, expect } from '@playwright/test'

/**
 * Base Page Object class - all page objects extend this
 */
export abstract class BasePage {
  protected page: Page
  abstract readonly url: string

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
      await this.page.waitForSelector(selector, { timeout: 3000 })
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
    await this.page.waitForSelector(`.ant-notification-notice-${type}`, { timeout: 5000 })
  }

  /**
   * Expect notification to appear
   */
  protected async expectNotification(type: 'success' | 'error'): Promise<void> {
    await expect(this.page.locator(`.ant-notification-notice-${type}`)).toBeVisible({ timeout: 5000 })
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
    await this.page.waitForSelector('.ant-spin', { state: 'hidden', timeout: 10000 }).catch(() => {})
  }
}

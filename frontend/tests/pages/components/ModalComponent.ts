import { Page, Locator, expect } from '@playwright/test'

/**
 * Modal Component Page Object (Ant Design Modal)
 */
export class ModalComponent {
  private page: Page

  constructor(page: Page) {
    this.page = page
  }

  // ==================== Locators ====================

  get modal(): Locator {
    return this.page.locator('.ant-modal')
  }

  get title(): Locator {
    return this.page.locator('.ant-modal-title')
  }

  get content(): Locator {
    return this.page.locator('.ant-modal-body')
  }

  get okButton(): Locator {
    return this.page.locator('.ant-modal-footer button.ant-btn-primary')
  }

  get cancelButton(): Locator {
    return this.page.locator('.ant-modal-footer button:not(.ant-btn-primary)')
  }

  get closeButton(): Locator {
    return this.page.locator('.ant-modal-close')
  }

  get footer(): Locator {
    return this.page.locator('.ant-modal-footer')
  }

  // ==================== Actions ====================

  /**
   * Wait for modal to open
   */
  async waitForOpen(): Promise<void> {
    await this.modal.waitFor({ state: 'visible' })
  }

  /**
   * Click OK button
   */
  async clickOk(): Promise<void> {
    await this.okButton.click()
  }

  /**
   * Click Cancel button
   */
  async clickCancel(): Promise<void> {
    await this.cancelButton.click()
  }

  /**
   * Close modal via X button
   */
  async close(): Promise<void> {
    await this.closeButton.click()
  }

  /**
   * Get modal title text
   */
  async getTitleText(): Promise<string> {
    return await this.title.textContent() || ''
  }

  /**
   * Fill input in modal by placeholder
   */
  async fillInput(placeholder: string, value: string): Promise<void> {
    await this.content.locator(`input[placeholder*="${placeholder}" i]`).fill(value)
  }

  /**
   * Fill textarea in modal
   */
  async fillTextarea(value: string): Promise<void> {
    await this.content.locator('textarea').fill(value)
  }

  // ==================== Assertions ====================

  /**
   * Expect modal to be visible
   */
  async expectVisible(): Promise<void> {
    await expect(this.modal).toBeVisible()
  }

  /**
   * Expect modal to be closed
   */
  async expectClosed(): Promise<void> {
    await expect(this.modal).not.toBeVisible()
  }

  /**
   * Expect modal title
   */
  async expectTitle(title: string): Promise<void> {
    await expect(this.title).toHaveText(title)
  }

  /**
   * Expect OK button enabled
   */
  async expectOkEnabled(): Promise<void> {
    await expect(this.okButton).toBeEnabled()
  }

  /**
   * Expect OK button disabled
   */
  async expectOkDisabled(): Promise<void> {
    await expect(this.okButton).toBeDisabled()
  }

  /**
   * Expect content contains text
   */
  async expectContentContains(text: string): Promise<void> {
    await expect(this.content).toContainText(text)
  }
}

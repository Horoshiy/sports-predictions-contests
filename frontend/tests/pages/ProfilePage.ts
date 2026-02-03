import { Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'

/**
 * Profile Page Object
 */
export class ProfilePage extends BasePage {
  readonly url = '/profile'

  // ==================== Locators ====================

  get displayNameInput(): Locator {
    return this.page.locator('input[placeholder*="name" i]').first()
  }

  get emailDisplay(): Locator {
    return this.page.locator('.ant-descriptions-item:has-text("Email")')
  }

  get avatarUpload(): Locator {
    return this.page.locator('.ant-upload, input[type="file"]')
  }

  get saveButton(): Locator {
    return this.page.locator('button:has-text("Save")')
  }

  get statsSection(): Locator {
    return this.page.locator('.ant-statistic, .stats-section')
  }

  get predictionHistory(): Locator {
    return this.page.locator('.prediction-history')
  }

  get privacyToggle(): Locator {
    return this.page.locator('.ant-switch')
  }

  get successMessage(): Locator {
    return this.page.locator('.ant-notification-notice-success')
  }

  get profileForm(): Locator {
    return this.page.locator('form')
  }

  // ==================== Actions ====================

  /**
   * Update display name
   */
  async updateDisplayName(name: string): Promise<void> {
    await this.displayNameInput.clear()
    await this.displayNameInput.fill(name)
  }

  /**
   * Upload avatar
   */
  async uploadAvatar(filePath: string): Promise<void> {
    const fileInput = this.page.locator('input[type="file"]')
    await fileInput.setInputFiles(filePath)
  }

  /**
   * Save profile
   */
  async saveProfile(): Promise<void> {
    await this.saveButton.click()
  }

  /**
   * Toggle privacy setting
   */
  async togglePrivacySetting(index: number = 0): Promise<void> {
    await this.privacyToggle.nth(index).click()
  }

  /**
   * Update and save profile
   */
  async updateProfile(name: string): Promise<void> {
    await this.updateDisplayName(name)
    await this.saveProfile()
  }

  // ==================== Assertions ====================

  /**
   * Expect profile to be loaded
   */
  async expectProfileLoaded(): Promise<void> {
    await expect(this.displayNameInput).toBeVisible({ timeout: 10000 })
  }

  /**
   * Expect save success notification
   */
  async expectSaveSuccess(): Promise<void> {
    await expect(this.successMessage).toBeVisible({ timeout: 5000 })
  }

  /**
   * Expect stats section visible
   */
  async expectStatsVisible(): Promise<void> {
    await expect(this.statsSection.first()).toBeVisible()
  }

  /**
   * Expect display name value
   */
  async expectDisplayName(name: string): Promise<void> {
    await expect(this.displayNameInput).toHaveValue(name)
  }

  /**
   * Expect to be on profile page
   */
  async expectOnProfilePage(): Promise<void> {
    await expect(this.page).toHaveURL('/profile')
  }

  /**
   * Expect form visible
   */
  async expectFormVisible(): Promise<void> {
    await expect(this.profileForm).toBeVisible()
  }
}

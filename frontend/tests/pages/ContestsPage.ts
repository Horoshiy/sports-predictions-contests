import { Locator, expect } from '@playwright/test'
import { TIMEOUTS } from '../helpers/test-config'
import { BasePage } from './BasePage'

/**
 * Contests Page Object
 */
export class ContestsPage extends BasePage {
  readonly url = '/contests'
  readonly pageName = 'Contests Page'

  // ==================== Locators ====================

  get createContestButton(): Locator {
    return this.page.locator('button:has-text("Create Contest")')
  }

  get contestCards(): Locator {
    return this.page.locator('.ant-card')
  }

  get contestModal(): Locator {
    return this.page.locator('.ant-modal')
  }

  get searchInput(): Locator {
    return this.page.locator('input[placeholder*="search" i]')
  }

  get filterDropdown(): Locator {
    return this.page.locator('.ant-select').first()
  }

  get contestTitleInput(): Locator {
    return this.page.locator('input[placeholder*="title" i]')
  }

  get contestDescriptionInput(): Locator {
    return this.page.locator('textarea')
  }

  get submitContestButton(): Locator {
    return this.page.locator('.ant-modal-footer button.ant-btn-primary')
  }

  get pagination(): Locator {
    return this.page.locator('.ant-pagination')
  }

  // ==================== Actions ====================

  /**
   * Click create contest button
   */
  async clickCreateContest(): Promise<void> {
    await this.createContestButton.click()
    await this.waitForModal()
  }

  /**
   * Get contest card by index
   */
  getContestCard(index: number): Locator {
    return this.contestCards.nth(index)
  }

  /**
   * Join contest by index
   */
  async joinContest(index: number): Promise<void> {
    const card = this.getContestCard(index)
    await card.locator('button:has-text("Join")').click()
  }

  /**
   * Leave contest by index
   */
  async leaveContest(index: number): Promise<void> {
    const card = this.getContestCard(index)
    await card.locator('button:has-text("Leave")').click()
  }

  /**
   * Open contest details by index
   */
  async openContestDetails(index: number): Promise<void> {
    const card = this.getContestCard(index)
    await card.click()
  }

  /**
   * Search contests
   */
  async searchContests(query: string): Promise<void> {
    await this.searchInput.fill(query)
    await this.page.keyboard.press('Enter')
    await this.waitForLoadingComplete()
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: string): Promise<void> {
    await this.filterDropdown.click()
    await this.page.click(`.ant-select-dropdown .ant-select-item:has-text("${status}")`)
    await this.waitForLoadingComplete()
  }

  /**
   * Create a new contest
   */
  async createContest(title: string, description: string): Promise<void> {
    await this.clickCreateContest()
    await this.contestTitleInput.fill(title)
    await this.contestDescriptionInput.fill(description)
    await this.submitContestButton.click()
    await this.waitForModalClosed()
  }

  /**
   * Go to next page
   */
  async goToNextPage(): Promise<void> {
    await this.pagination.locator('.ant-pagination-next').click()
    await this.waitForLoadingComplete()
  }

  /**
   * Go to previous page
   */
  async goToPreviousPage(): Promise<void> {
    await this.pagination.locator('.ant-pagination-prev').click()
    await this.waitForLoadingComplete()
  }

  // ==================== Assertions ====================

  /**
   * Expect contests list to be visible
   */
  async expectContestListVisible(): Promise<void> {
    await expect(
      this.contestCards.first(),
      'Contest cards should be visible on Contests Page'
    ).toBeVisible({ timeout: TIMEOUTS.LONG })
  }

  /**
   * Expect specific contest count
   */
  async expectContestCount(count: number): Promise<void> {
    await expect(this.contestCards).toHaveCount(count)
  }

  /**
   * Expect contest count at least
   */
  async expectContestCountAtLeast(count: number): Promise<void> {
    const actualCount = await this.contestCards.count()
    expect(actualCount).toBeGreaterThanOrEqual(count)
  }

  /**
   * Expect join button visible on contest
   */
  async expectJoinButtonVisible(index: number): Promise<void> {
    const card = this.getContestCard(index)
    await expect(card.locator('button:has-text("Join")')).toBeVisible()
  }

  /**
   * Expect leave button visible on contest
   */
  async expectLeaveButtonVisible(index: number): Promise<void> {
    const card = this.getContestCard(index)
    await expect(card.locator('button:has-text("Leave")')).toBeVisible()
  }

  /**
   * Expect to be on contests page
   */
  async expectOnContestsPage(): Promise<void> {
    await expect(this.page).toHaveURL('/contests')
  }

  /**
   * Expect create contest modal visible
   */
  async expectCreateModalVisible(): Promise<void> {
    await expect(this.contestModal).toBeVisible()
  }
}

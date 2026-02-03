import { Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'

/**
 * Predictions Page Object
 */
export class PredictionsPage extends BasePage {
  readonly url = '/predictions'

  // ==================== Locators ====================

  get eventCards(): Locator {
    return this.page.locator('.ant-card')
  }

  get predictionForm(): Locator {
    return this.page.locator('form')
  }

  get submitButton(): Locator {
    return this.page.locator('button:has-text("Submit")')
  }

  get contestSelector(): Locator {
    return this.page.locator('.ant-select').first()
  }

  get predictionHistory(): Locator {
    return this.page.locator('.prediction-history, .ant-list')
  }

  get homeScoreInput(): Locator {
    return this.page.locator('input[type="number"]').first()
  }

  get awayScoreInput(): Locator {
    return this.page.locator('input[type="number"]').nth(1)
  }

  get editButton(): Locator {
    return this.page.locator('button:has-text("Edit")')
  }

  get cancelButton(): Locator {
    return this.page.locator('button:has-text("Cancel")')
  }

  get saveButton(): Locator {
    return this.page.locator('button:has-text("Save")')
  }

  // ==================== Actions ====================

  /**
   * Get event card by index
   */
  getEventCard(index: number): Locator {
    return this.eventCards.nth(index)
  }

  /**
   * Select contest
   */
  async selectContest(contestName: string): Promise<void> {
    await this.contestSelector.click()
    await this.page.click(`.ant-select-dropdown .ant-select-item:has-text("${contestName}")`)
    await this.waitForLoadingComplete()
  }

  /**
   * Submit prediction for an event
   */
  async submitPrediction(eventIndex: number, homeScore: number, awayScore: number): Promise<void> {
    const card = this.getEventCard(eventIndex)
    await card.click()
    
    await this.homeScoreInput.fill(homeScore.toString())
    await this.awayScoreInput.fill(awayScore.toString())
    await this.submitButton.click()
  }

  /**
   * Edit existing prediction
   */
  async editPrediction(eventIndex: number, newHomeScore: number, newAwayScore: number): Promise<void> {
    const card = this.getEventCard(eventIndex)
    await card.locator('button:has-text("Edit")').click()
    
    await this.homeScoreInput.clear()
    await this.homeScoreInput.fill(newHomeScore.toString())
    await this.awayScoreInput.clear()
    await this.awayScoreInput.fill(newAwayScore.toString())
    await this.saveButton.click()
  }

  /**
   * Cancel editing
   */
  async cancelEdit(): Promise<void> {
    await this.cancelButton.click()
  }

  /**
   * Click on event card to expand
   */
  async expandEvent(index: number): Promise<void> {
    await this.getEventCard(index).click()
  }

  // ==================== Assertions ====================

  /**
   * Expect events to be visible
   */
  async expectEventsVisible(): Promise<void> {
    await expect(this.eventCards.first()).toBeVisible({ timeout: 10000 })
  }

  /**
   * Expect prediction submitted notification
   */
  async expectPredictionSubmitted(): Promise<void> {
    await this.expectNotification('success')
  }

  /**
   * Expect event count
   */
  async expectEventCount(count: number): Promise<void> {
    await expect(this.eventCards).toHaveCount(count)
  }

  /**
   * Expect to be on predictions page
   */
  async expectOnPredictionsPage(): Promise<void> {
    await expect(this.page).toHaveURL('/predictions')
  }

  /**
   * Expect prediction form visible
   */
  async expectPredictionFormVisible(): Promise<void> {
    await expect(this.homeScoreInput).toBeVisible()
    await expect(this.awayScoreInput).toBeVisible()
  }
}

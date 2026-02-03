import { Page, Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'
import { TIMEOUTS } from '../helpers/test-config'

/**
 * Analytics Page Object
 * Handles user analytics dashboard with charts and statistics
 */
export class AnalyticsPage extends BasePage {
  readonly url = '/analytics'
  readonly pageName = 'Analytics Dashboard'

  // Time range selector
  private readonly timeRangeSelector = '.ant-segmented'

  // Statistic cards selectors
  private readonly statisticCards = '.ant-statistic'
  private readonly accuracyCard = '.ant-card:has(.ant-statistic-title:has-text("Accuracy"))'
  private readonly totalPointsCard = '.ant-card:has(.ant-statistic-title:has-text("Total Points"))'
  private readonly totalPredictionsCard = '.ant-card:has(.ant-statistic-title:has-text("Total Predictions"))'
  private readonly correctPredictionsCard = '.ant-card:has(.ant-statistic-title:has-text("Correct"))'

  // Chart containers
  private readonly accuracyChart = '.ant-card:has-text("Accuracy")'
  private readonly sportBreakdownChart = '.ant-card:has-text("Sport")'
  private readonly platformComparisonSection = '.ant-card:has-text("Platform")'

  // Export button
  private readonly exportButton = 'button:has-text("Export")'

  constructor(page: Page) {
    super(page)
  }

  // ==================== Time Range Selection ====================

  /**
   * Select time range
   */
  async selectTimeRange(range: '7d' | '30d' | '90d' | 'all'): Promise<void> {
    const labels: Record<string, string> = {
      '7d': '7 Days',
      '30d': '30 Days',
      '90d': '90 Days',
      'all': 'All Time',
    }
    await this.page.click(`${this.timeRangeSelector} .ant-segmented-item:has-text("${labels[range]}")`)
    await this.waitForLoadingComplete()
  }

  /**
   * Get current time range
   */
  async getCurrentTimeRange(): Promise<string> {
    const selected = this.page.locator(`${this.timeRangeSelector} .ant-segmented-item-selected`)
    return await selected.textContent() || ''
  }

  // ==================== Statistics Retrieval ====================

  /**
   * Get accuracy percentage
   */
  async getAccuracy(): Promise<string | null> {
    const value = this.page.locator(`${this.accuracyCard} .ant-statistic-content-value`)
    return await value.textContent()
  }

  /**
   * Get total points
   */
  async getTotalPoints(): Promise<string | null> {
    const value = this.page.locator(`${this.totalPointsCard} .ant-statistic-content-value`)
    return await value.textContent()
  }

  /**
   * Get total predictions count
   */
  async getTotalPredictions(): Promise<string | null> {
    const value = this.page.locator(`${this.totalPredictionsCard} .ant-statistic-content-value`)
    return await value.textContent()
  }

  /**
   * Get correct predictions count
   */
  async getCorrectPredictions(): Promise<string | null> {
    const value = this.page.locator(`${this.correctPredictionsCard} .ant-statistic-content-value`)
    return await value.textContent()
  }

  /**
   * Get all statistics as object
   */
  async getAllStatistics(): Promise<{
    accuracy: string | null
    totalPoints: string | null
    totalPredictions: string | null
    correctPredictions: string | null
  }> {
    return {
      accuracy: await this.getAccuracy(),
      totalPoints: await this.getTotalPoints(),
      totalPredictions: await this.getTotalPredictions(),
      correctPredictions: await this.getCorrectPredictions(),
    }
  }

  // ==================== Charts ====================

  /**
   * Check if accuracy chart is visible
   */
  async isAccuracyChartVisible(): Promise<boolean> {
    return await this.isVisible(this.accuracyChart)
  }

  /**
   * Check if sport breakdown chart is visible
   */
  async isSportBreakdownVisible(): Promise<boolean> {
    return await this.isVisible(this.sportBreakdownChart)
  }

  /**
   * Check if platform comparison is visible
   */
  async isPlatformComparisonVisible(): Promise<boolean> {
    return await this.isVisible(this.platformComparisonSection)
  }

  // ==================== Export ====================

  /**
   * Click export button
   */
  async clickExport(): Promise<void> {
    await this.page.click(this.exportButton)
  }

  /**
   * Export analytics data
   * Returns the download promise for verification
   */
  async exportData(): Promise<void> {
    const downloadPromise = this.page.waitForEvent('download')
    await this.clickExport()
    await downloadPromise
  }

  // ==================== Page State ====================

  /**
   * Check if page is loading
   */
  async isLoading(): Promise<boolean> {
    return await this.isVisible('.ant-spin')
  }

  /**
   * Check if error alert is shown
   */
  async hasError(): Promise<boolean> {
    return await this.isVisible('.ant-alert-error')
  }

  /**
   * Get error message
   */
  async getErrorMessage(): Promise<string | null> {
    const alert = this.page.locator('.ant-alert-error .ant-alert-description')
    if (await alert.isVisible()) {
      return await alert.textContent()
    }
    return null
  }

  /**
   * Check if no data alert is shown
   */
  async hasNoData(): Promise<boolean> {
    return await this.isVisible('.ant-alert-info:has-text("No data")')
  }

  // ==================== Assertions ====================

  /**
   * Assert page loaded successfully
   */
  async expectPageLoaded(): Promise<void> {
    await expect(
      this.page.locator('h2:has-text("Analytics Dashboard")'),
      'Expected Analytics Dashboard title'
    ).toBeVisible()
  }

  /**
   * Assert statistics cards are visible
   */
  async expectStatisticsVisible(): Promise<void> {
    await expect(
      this.page.locator(this.accuracyCard),
      'Expected Accuracy card'
    ).toBeVisible()
    await expect(
      this.page.locator(this.totalPointsCard),
      'Expected Total Points card'
    ).toBeVisible()
    await expect(
      this.page.locator(this.totalPredictionsCard),
      'Expected Total Predictions card'
    ).toBeVisible()
    await expect(
      this.page.locator(this.correctPredictionsCard),
      'Expected Correct Predictions card'
    ).toBeVisible()
  }

  /**
   * Assert time range is selected
   */
  async expectTimeRangeSelected(range: '7 Days' | '30 Days' | '90 Days' | 'All Time'): Promise<void> {
    await expect(
      this.page.locator(`${this.timeRangeSelector} .ant-segmented-item-selected:has-text("${range}")`),
      `Expected "${range}" to be selected`
    ).toBeVisible()
  }

  /**
   * Assert accuracy value is valid (contains %)
   */
  async expectValidAccuracy(): Promise<void> {
    const accuracy = await this.getAccuracy()
    expect(accuracy).not.toBeNull()
    expect(accuracy).toContain('%')
  }

  /**
   * Assert no error state
   */
  async expectNoError(): Promise<void> {
    await expect(
      this.page.locator('.ant-alert-error'),
      'Expected no error alert'
    ).toBeHidden()
  }

  /**
   * Assert charts are rendered
   */
  async expectChartsRendered(): Promise<void> {
    // Wait for charts to load (they contain canvas or svg elements)
    await expect(
      this.page.locator(`${this.accuracyChart} canvas, ${this.accuracyChart} svg`),
      'Expected accuracy chart to be rendered'
    ).toBeVisible({ timeout: TIMEOUTS.MEDIUM })
  }
}

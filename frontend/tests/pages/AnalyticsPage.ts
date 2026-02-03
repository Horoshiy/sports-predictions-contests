import { Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'

/**
 * Analytics Page Object
 */
export class AnalyticsPage extends BasePage {
  readonly url = '/analytics'

  // ==================== Locators ====================

  get accuracyChart(): Locator {
    return this.page.locator('.recharts-wrapper, .accuracy-chart').first()
  }

  get sportBreakdown(): Locator {
    return this.page.locator('.sport-breakdown, .recharts-pie')
  }

  get platformComparison(): Locator {
    return this.page.locator('.platform-comparison')
  }

  get exportButton(): Locator {
    return this.page.locator('button:has-text("Export")')
  }

  get dateFilter(): Locator {
    return this.page.locator('.ant-picker-range')
  }

  get statsCards(): Locator {
    return this.page.locator('.ant-statistic')
  }

  get tabs(): Locator {
    return this.page.locator('.ant-tabs-tab')
  }

  get charts(): Locator {
    return this.page.locator('.recharts-wrapper')
  }

  get filterDropdown(): Locator {
    return this.page.locator('.ant-select')
  }

  // ==================== Actions ====================

  /**
   * Filter by date range
   */
  async filterByDateRange(startDate: string, endDate: string): Promise<void> {
    await this.dateFilter.click()
    await this.page.fill('.ant-picker-input input', startDate)
    await this.page.keyboard.press('Tab')
    await this.page.fill('.ant-picker-input:nth-child(2) input', endDate)
    await this.page.keyboard.press('Enter')
    await this.waitForLoadingComplete()
  }

  /**
   * Export data
   */
  async exportData(format: string = 'CSV'): Promise<void> {
    await this.exportButton.click()
    if (format !== 'CSV') {
      await this.page.click(`.ant-dropdown-menu-item:has-text("${format}")`)
    }
  }

  /**
   * Switch to a tab
   */
  async switchTab(tabName: string): Promise<void> {
    await this.page.click(`.ant-tabs-tab:has-text("${tabName}")`)
    await this.waitForLoadingComplete()
  }

  /**
   * Apply filter
   */
  async applyFilter(filterName: string, value: string): Promise<void> {
    await this.filterDropdown.click()
    await this.page.click(`.ant-select-item:has-text("${value}")`)
    await this.waitForLoadingComplete()
  }

  // ==================== Assertions ====================

  /**
   * Expect charts to be loaded
   */
  async expectChartsLoaded(): Promise<void> {
    await expect(this.charts.first()).toBeVisible({ timeout: 10000 })
  }

  /**
   * Expect accuracy chart visible
   */
  async expectAccuracyChartVisible(): Promise<void> {
    await expect(this.accuracyChart).toBeVisible()
  }

  /**
   * Expect sport breakdown visible
   */
  async expectSportBreakdownVisible(): Promise<void> {
    await expect(this.sportBreakdown).toBeVisible()
  }

  /**
   * Expect stats cards visible
   */
  async expectStatsVisible(): Promise<void> {
    await expect(this.statsCards.first()).toBeVisible()
  }

  /**
   * Expect data exported notification
   */
  async expectDataExported(): Promise<void> {
    await this.expectNotification('success')
  }

  /**
   * Expect to be on analytics page
   */
  async expectOnAnalyticsPage(): Promise<void> {
    await expect(this.page).toHaveURL('/analytics')
  }

  /**
   * Expect chart count
   */
  async expectChartCount(count: number): Promise<void> {
    await expect(this.charts).toHaveCount(count)
  }
}

import { Locator, expect } from '@playwright/test'
import { TIMEOUTS } from '../helpers/test-config'
import { BasePage } from './BasePage'

/**
 * Sports Page Object (Admin sports management)
 */
export class SportsPage extends BasePage {
  readonly url = '/sports'
  readonly pageName = 'Sports Page'

  // ==================== Locators ====================

  get sportsList(): Locator {
    return this.page.locator('.sports-list, .ant-list').first()
  }

  get leaguesList(): Locator {
    return this.page.locator('.leagues-list, .ant-list').nth(1)
  }

  get teamsList(): Locator {
    return this.page.locator('.teams-list, .ant-list')
  }

  get matchesList(): Locator {
    return this.page.locator('.matches-list, .ant-table')
  }

  get addSportButton(): Locator {
    return this.page.locator('button:has-text("Add Sport")')
  }

  get addLeagueButton(): Locator {
    return this.page.locator('button:has-text("Add League")')
  }

  get addTeamButton(): Locator {
    return this.page.locator('button:has-text("Add Team")')
  }

  get addMatchButton(): Locator {
    return this.page.locator('button:has-text("Add Match")')
  }

  get sportCards(): Locator {
    return this.page.locator('.ant-card')
  }

  get modal(): Locator {
    return this.page.locator('.ant-modal')
  }

  get nameInput(): Locator {
    return this.page.locator('input[placeholder*="name" i]').first()
  }

  get submitButton(): Locator {
    return this.page.locator('.ant-modal-footer button.ant-btn-primary')
  }

  get tabs(): Locator {
    return this.page.locator('.ant-tabs-tab')
  }

  // ==================== Actions ====================

  /**
   * Select a sport
   */
  async selectSport(name: string): Promise<void> {
    await this.page.click(`.ant-card:has-text("${name}")`)
    await this.waitForLoadingComplete()
  }

  /**
   * Select a league
   */
  async selectLeague(name: string): Promise<void> {
    await this.page.click(`.ant-list-item:has-text("${name}")`)
    await this.waitForLoadingComplete()
  }

  /**
   * Click add sport button
   */
  async clickAddSport(): Promise<void> {
    await this.addSportButton.click()
    await this.waitForModal()
  }

  /**
   * Add a new sport
   */
  async addSport(name: string): Promise<void> {
    await this.clickAddSport()
    await this.nameInput.fill(name)
    await this.submitButton.click()
    await this.waitForModalClosed()
  }

  /**
   * Click add league button
   */
  async clickAddLeague(): Promise<void> {
    await this.addLeagueButton.click()
    await this.waitForModal()
  }

  /**
   * Add a new league
   */
  async addLeague(name: string): Promise<void> {
    await this.clickAddLeague()
    await this.nameInput.fill(name)
    await this.submitButton.click()
    await this.waitForModalClosed()
  }

  /**
   * Click add team button
   */
  async clickAddTeam(): Promise<void> {
    await this.addTeamButton.click()
    await this.waitForModal()
  }

  /**
   * Add a new team
   */
  async addTeam(name: string): Promise<void> {
    await this.clickAddTeam()
    await this.nameInput.fill(name)
    await this.submitButton.click()
    await this.waitForModalClosed()
  }

  /**
   * Click add match button
   */
  async clickAddMatch(): Promise<void> {
    await this.addMatchButton.click()
    await this.waitForModal()
  }

  /**
   * Switch to tab
   */
  async switchTab(tabName: string): Promise<void> {
    await this.page.click(`.ant-tabs-tab:has-text("${tabName}")`)
    await this.waitForLoadingComplete()
  }

  /**
   * Edit item by clicking edit button
   */
  async clickEdit(index: number): Promise<void> {
    await this.page.locator('button:has-text("Edit")').nth(index).click()
    await this.waitForModal()
  }

  /**
   * Delete item by clicking delete button
   */
  async clickDelete(index: number): Promise<void> {
    await this.page.locator('button:has-text("Delete")').nth(index).click()
  }

  // ==================== Assertions ====================

  /**
   * Expect sports list loaded
   */
  async expectSportsLoaded(): Promise<void> {
    await expect(this.sportCards.first()).toBeVisible({ timeout: 10000 })
  }

  /**
   * Expect leagues visible
   */
  async expectLeaguesVisible(): Promise<void> {
    await expect(this.leaguesList).toBeVisible()
  }

  /**
   * Expect modal visible
   */
  async expectModalVisible(): Promise<void> {
    await expect(this.modal).toBeVisible()
  }

  /**
   * Expect to be on sports page
   */
  async expectOnSportsPage(): Promise<void> {
    await expect(this.page).toHaveURL('/sports')
  }

  /**
   * Expect teams visible
   */
  async expectTeamsVisible(): Promise<void> {
    await expect(this.teamsList).toBeVisible()
  }

  /**
   * Expect matches visible
   */
  async expectMatchesVisible(): Promise<void> {
    await expect(this.matchesList).toBeVisible()
  }
}

import { Page, Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'
import { TIMEOUTS } from '../helpers/test-config'

/**
 * Sports Management Page Object
 * Handles Sports, Leagues, Teams, and Matches tabs
 */
export class SportsPage extends BasePage {
  readonly url = '/sports'
  readonly pageName = 'Sports Management'

  // Tab selectors
  private readonly tabs = {
    sports: '.ant-tabs-tab:has-text("Sports")',
    leagues: '.ant-tabs-tab:has-text("Leagues")',
    teams: '.ant-tabs-tab:has-text("Teams")',
    matches: '.ant-tabs-tab:has-text("Matches")',
  }

  // Common selectors
  private readonly addButton = 'button:has-text("Add"), button:has-text("Create")'
  private readonly table = '.ant-table'
  private readonly modal = '.ant-modal'
  private readonly modalTitle = '.ant-modal-title'
  private readonly submitButton = '.ant-modal-footer button.ant-btn-primary'

  constructor(page: Page) {
    super(page)
  }

  // ==================== Navigation ====================

  /**
   * Switch to Sports tab
   */
  async goToSportsTab(): Promise<void> {
    await this.page.click(this.tabs.sports)
    await this.waitForLoadingComplete()
  }

  /**
   * Switch to Leagues tab
   */
  async goToLeaguesTab(): Promise<void> {
    await this.page.click(this.tabs.leagues)
    await this.waitForLoadingComplete()
  }

  /**
   * Switch to Teams tab
   */
  async goToTeamsTab(): Promise<void> {
    await this.page.click(this.tabs.teams)
    await this.waitForLoadingComplete()
  }

  /**
   * Switch to Matches tab
   */
  async goToMatchesTab(): Promise<void> {
    await this.page.click(this.tabs.matches)
    await this.waitForLoadingComplete()
  }

  /**
   * Get active tab name
   */
  async getActiveTab(): Promise<string> {
    const activeTab = this.page.locator('.ant-tabs-tab-active')
    return await activeTab.textContent() || ''
  }

  // ==================== Sports CRUD ====================

  /**
   * Open create sport form
   */
  async openCreateSportForm(): Promise<void> {
    await this.goToSportsTab()
    await this.page.click('button:has-text("Add Sport"), button:has-text("Create Sport")')
    await this.waitForModal()
  }

  /**
   * Create a new sport
   */
  async createSport(name: string, description?: string, iconUrl?: string): Promise<void> {
    await this.openCreateSportForm()
    await this.fillSportForm(name, description, iconUrl)
    await this.submitForm()
  }

  /**
   * Edit an existing sport
   */
  async editSport(sportName: string): Promise<void> {
    await this.goToSportsTab()
    await this.clickTableRowAction(sportName, 'Edit')
    await this.waitForModal()
  }

  /**
   * Fill sport form fields
   */
  async fillSportForm(name: string, description?: string, iconUrl?: string): Promise<void> {
    await this.fillAntdInputByLabel('Name', name)
    if (description) {
      await this.fillAntdInputByLabel('Description', description)
    }
    if (iconUrl) {
      await this.fillAntdInputByLabel('Icon URL', iconUrl)
    }
  }

  /**
   * Get sports count in table
   */
  async getSportsCount(): Promise<number> {
    await this.goToSportsTab()
    return await this.getTableRowCount()
  }

  // ==================== Leagues CRUD ====================

  /**
   * Open create league form
   */
  async openCreateLeagueForm(): Promise<void> {
    await this.goToLeaguesTab()
    await this.page.click('button:has-text("Add League"), button:has-text("Create League")')
    await this.waitForModal()
  }

  /**
   * Create a new league
   */
  async createLeague(name: string, sportName: string, country?: string, season?: string): Promise<void> {
    await this.openCreateLeagueForm()
    await this.fillLeagueForm(name, sportName, country, season)
    await this.submitForm()
  }

  /**
   * Edit an existing league
   */
  async editLeague(leagueName: string): Promise<void> {
    await this.goToLeaguesTab()
    await this.clickTableRowAction(leagueName, 'Edit')
    await this.waitForModal()
  }

  /**
   * Fill league form fields
   */
  async fillLeagueForm(name: string, sportName: string, country?: string, season?: string): Promise<void> {
    await this.fillAntdInputByLabel('Name', name)
    await this.selectAntdOptionByLabel('Sport', sportName)
    if (country) {
      await this.fillAntdInputByLabel('Country', country)
    }
    if (season) {
      await this.fillAntdInputByLabel('Season', season)
    }
  }

  /**
   * Get leagues count in table
   */
  async getLeaguesCount(): Promise<number> {
    await this.goToLeaguesTab()
    return await this.getTableRowCount()
  }

  // ==================== Teams CRUD ====================

  /**
   * Open create team form
   */
  async openCreateTeamForm(): Promise<void> {
    await this.goToTeamsTab()
    await this.page.click('button:has-text("Add Team"), button:has-text("Create Team")')
    await this.waitForModal()
  }

  /**
   * Create a new team
   */
  async createTeam(name: string, sportName: string, country?: string, shortName?: string): Promise<void> {
    await this.openCreateTeamForm()
    await this.fillTeamForm(name, sportName, country, shortName)
    await this.submitForm()
  }

  /**
   * Edit an existing team
   */
  async editTeam(teamName: string): Promise<void> {
    await this.goToTeamsTab()
    await this.clickTableRowAction(teamName, 'Edit')
    await this.waitForModal()
  }

  /**
   * Fill team form fields
   */
  async fillTeamForm(name: string, sportName: string, country?: string, shortName?: string): Promise<void> {
    await this.fillAntdInputByLabel('Name', name)
    await this.selectAntdOptionByLabel('Sport', sportName)
    if (country) {
      await this.fillAntdInputByLabel('Country', country)
    }
    if (shortName) {
      await this.fillAntdInputByLabel('Short Name', shortName)
    }
  }

  /**
   * Get teams count in table
   */
  async getTeamsCount(): Promise<number> {
    await this.goToTeamsTab()
    return await this.getTableRowCount()
  }

  // ==================== Matches CRUD ====================

  /**
   * Open create match form
   */
  async openCreateMatchForm(): Promise<void> {
    await this.goToMatchesTab()
    await this.page.click('button:has-text("Add Match"), button:has-text("Create Match")')
    await this.waitForModal()
  }

  /**
   * Create a new match
   */
  async createMatch(leagueName: string, homeTeam: string, awayTeam: string, scheduledAt?: Date): Promise<void> {
    await this.openCreateMatchForm()
    await this.fillMatchForm(leagueName, homeTeam, awayTeam, scheduledAt)
    await this.submitForm()
  }

  /**
   * Edit an existing match
   */
  async editMatch(matchDescription: string): Promise<void> {
    await this.goToMatchesTab()
    await this.clickTableRowAction(matchDescription, 'Edit')
    await this.waitForModal()
  }

  /**
   * Fill match form fields
   */
  async fillMatchForm(leagueName: string, homeTeam: string, awayTeam: string, scheduledAt?: Date): Promise<void> {
    await this.selectAntdOptionByLabel('League', leagueName)
    await this.selectAntdOptionByLabel('Home Team', homeTeam)
    await this.selectAntdOptionByLabel('Away Team', awayTeam)
    if (scheduledAt) {
      await this.setDateTimePicker('Scheduled At', scheduledAt)
    }
  }

  /**
   * Get matches count in table
   */
  async getMatchesCount(): Promise<number> {
    await this.goToMatchesTab()
    return await this.getTableRowCount()
  }

  /**
   * Set match result (for editing)
   */
  async setMatchResult(homeScore: number, awayScore: number): Promise<void> {
    await this.fillAntdInputByLabel('Home Score', homeScore.toString())
    await this.fillAntdInputByLabel('Away Score', awayScore.toString())
  }

  // ==================== Common Actions ====================

  /**
   * Submit form in modal
   */
  async submitForm(): Promise<void> {
    await this.page.click(this.submitButton)
    await this.waitForModalClosed()
  }

  /**
   * Cancel form
   */
  async cancelForm(): Promise<void> {
    await this.clickModalCancel()
    await this.waitForModalClosed()
  }

  /**
   * Click action button in table row
   */
  private async clickTableRowAction(rowText: string, actionText: string): Promise<void> {
    const row = this.page.locator(`${this.table} tbody tr:has-text("${rowText}")`)
    await row.locator(`button:has-text("${actionText}"), a:has-text("${actionText}")`).click()
  }

  /**
   * Delete entity from table
   */
  async deleteEntity(entityName: string): Promise<void> {
    await this.clickTableRowAction(entityName, 'Delete')
    // Confirm deletion in popconfirm or modal
    await this.page.click('.ant-popconfirm-buttons button.ant-btn-primary, .ant-modal-footer button.ant-btn-primary')
    await this.waitForLoadingComplete()
  }

  /**
   * Set date time picker value
   */
  private async setDateTimePicker(label: string, date: Date): Promise<void> {
    const datePickerWrapper = this.page.locator(`label:has-text("${label}")`).locator('..').locator('.ant-picker')
    await datePickerWrapper.click()
    
    // Format date for input
    const formattedDate = date.toISOString().split('T')[0]
    await this.page.fill('.ant-picker-input input', formattedDate)
    await this.page.keyboard.press('Enter')
  }

  /**
   * Search in table
   */
  async searchInTable(query: string): Promise<void> {
    const searchInput = this.page.locator('input[placeholder*="Search"], .ant-input-search input')
    await searchInput.fill(query)
    await this.page.keyboard.press('Enter')
    await this.waitForLoadingComplete()
  }

  /**
   * Clear search
   */
  async clearSearch(): Promise<void> {
    const searchInput = this.page.locator('input[placeholder*="Search"], .ant-input-search input')
    await searchInput.clear()
    await this.page.keyboard.press('Enter')
    await this.waitForLoadingComplete()
  }

  // ==================== Assertions ====================

  /**
   * Assert entity exists in table
   */
  async expectEntityInTable(entityName: string): Promise<void> {
    await expect(
      this.page.locator(`${this.table} tbody tr:has-text("${entityName}")`),
      `Expected "${entityName}" to be in table on ${this.pageName}`
    ).toBeVisible()
  }

  /**
   * Assert entity does not exist in table
   */
  async expectEntityNotInTable(entityName: string): Promise<void> {
    await expect(
      this.page.locator(`${this.table} tbody tr:has-text("${entityName}")`),
      `Expected "${entityName}" to NOT be in table on ${this.pageName}`
    ).toBeHidden()
  }

  /**
   * Assert modal is visible with title
   */
  async expectModalWithTitle(title: string): Promise<void> {
    await expect(
      this.page.locator(`${this.modalTitle}:has-text("${title}")`),
      `Expected modal with title "${title}"`
    ).toBeVisible()
  }

  /**
   * Assert tab is active
   */
  async expectTabActive(tabName: 'Sports' | 'Leagues' | 'Teams' | 'Matches'): Promise<void> {
    await expect(
      this.page.locator(`.ant-tabs-tab-active:has-text("${tabName}")`),
      `Expected "${tabName}" tab to be active`
    ).toBeVisible()
  }
}

import { Page, Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'
import { TIMEOUTS } from '../helpers/test-config'

/**
 * Teams Page Object
 * Handles team management: My Teams, All Teams, Join Team
 */
export class TeamsPage extends BasePage {
  readonly url = '/teams'
  readonly pageName = 'Teams'

  // Tab selectors
  private readonly tabs = {
    myTeams: '.ant-tabs-tab:has-text("My Teams")',
    allTeams: '.ant-tabs-tab:has-text("All Teams")',
    joinTeam: '.ant-tabs-tab:has-text("Join Team")',
  }

  // Common selectors
  private readonly table = '.ant-table'
  private readonly createTeamButton = 'button:has-text("Create Team")'
  private readonly teamFormModal = '.ant-modal'
  private readonly teamDetailsModal = '[data-testid="team-details-modal"]'
  private readonly joinTeamForm = '[data-testid="join-team-form"]'

  constructor(page: Page) {
    super(page)
  }

  // ==================== Navigation ====================

  /**
   * Switch to My Teams tab
   */
  async goToMyTeamsTab(): Promise<void> {
    await this.page.click(this.tabs.myTeams)
    await this.waitForLoadingComplete()
  }

  /**
   * Switch to All Teams tab
   */
  async goToAllTeamsTab(): Promise<void> {
    await this.page.click(this.tabs.allTeams)
    await this.waitForLoadingComplete()
  }

  /**
   * Switch to Join Team tab
   */
  async goToJoinTeamTab(): Promise<void> {
    await this.page.click(this.tabs.joinTeam)
  }

  /**
   * Get active tab name
   */
  async getActiveTab(): Promise<string> {
    const activeTab = this.page.locator('.ant-tabs-tab-active')
    return await activeTab.textContent() || ''
  }

  // ==================== Team Creation ====================

  /**
   * Open create team form
   */
  async openCreateTeamForm(): Promise<void> {
    await this.page.click(this.createTeamButton)
    await this.waitForModal()
  }

  /**
   * Create a new team
   */
  async createTeam(name: string, description?: string, maxMembers?: number): Promise<void> {
    await this.openCreateTeamForm()
    await this.fillTeamForm(name, description, maxMembers)
    await this.submitForm()
  }

  /**
   * Fill team form fields
   */
  async fillTeamForm(name: string, description?: string, maxMembers?: number): Promise<void> {
    await this.fillAntdInputByLabel('Name', name)
    if (description) {
      await this.fillAntdInputByLabel('Description', description)
    }
    if (maxMembers) {
      await this.fillAntdInputByLabel('Max Members', maxMembers.toString())
    }
  }

  /**
   * Submit form in modal
   */
  async submitForm(): Promise<void> {
    await this.page.click('.ant-modal-footer button.ant-btn-primary')
    await this.waitForModalClosed()
  }

  /**
   * Cancel form
   */
  async cancelForm(): Promise<void> {
    await this.clickModalCancel()
    await this.waitForModalClosed()
  }

  // ==================== Team Actions ====================

  /**
   * View team members
   */
  async viewTeamMembers(teamName: string): Promise<void> {
    const row = this.page.locator(`${this.table} tbody tr:has-text("${teamName}")`)
    await row.locator('button:has-text("Members"), button:has-text("View")').click()
    await this.page.waitForSelector(this.teamDetailsModal, { state: 'visible' })
  }

  /**
   * Edit team
   */
  async editTeam(teamName: string): Promise<void> {
    const row = this.page.locator(`${this.table} tbody tr:has-text("${teamName}")`)
    await row.locator('button:has-text("Edit")').click()
    await this.waitForModal()
  }

  /**
   * Close team details modal
   */
  async closeTeamDetails(): Promise<void> {
    await this.page.click(`${this.teamDetailsModal} button:has-text("Close")`)
    await this.page.waitForSelector(this.teamDetailsModal, { state: 'hidden' })
  }

  // ==================== Join Team ====================

  /**
   * Join team with invite code
   */
  async joinTeamWithCode(inviteCode: string): Promise<void> {
    await this.goToJoinTeamTab()
    await this.page.fill(`${this.joinTeamForm} input[placeholder*="invite"]`, inviteCode)
    await this.page.click(`${this.joinTeamForm} button[type="submit"]`)
  }

  /**
   * Get join team error message
   */
  async getJoinTeamError(): Promise<string | null> {
    const errorSpan = this.page.locator(`${this.joinTeamForm} span[style*="color"]`)
    if (await errorSpan.isVisible()) {
      return await errorSpan.textContent()
    }
    return null
  }

  // ==================== Team Details Modal ====================

  /**
   * Go to Team Info tab in modal
   */
  async goToTeamInfoTab(): Promise<void> {
    await this.page.click(`${this.teamDetailsModal} .ant-tabs-tab:has-text("Team Info")`)
  }

  /**
   * Go to Members tab in modal
   */
  async goToMembersTab(): Promise<void> {
    await this.page.click(`${this.teamDetailsModal} .ant-tabs-tab:has-text("Members")`)
  }

  /**
   * Go to Invite Code tab in modal
   */
  async goToInviteCodeTab(): Promise<void> {
    await this.page.click(`${this.teamDetailsModal} .ant-tabs-tab:has-text("Invite Code")`)
  }

  /**
   * Get invite code from modal
   */
  async getInviteCode(): Promise<string | null> {
    await this.goToInviteCodeTab()
    const codeElement = this.page.locator(`${this.teamDetailsModal} [data-testid="invite-code"], ${this.teamDetailsModal} code, ${this.teamDetailsModal} .ant-typography-copy-success`)
    if (await codeElement.isVisible()) {
      return await codeElement.textContent()
    }
    return null
  }

  /**
   * Leave team (non-captain)
   */
  async leaveTeam(): Promise<void> {
    await this.page.click(`${this.teamDetailsModal} button:has-text("Leave Team")`)
    await this.page.click('.ant-popconfirm-buttons button.ant-btn-primary')
    await this.waitForLoadingComplete()
  }

  /**
   * Delete team (captain only)
   */
  async deleteTeam(): Promise<void> {
    await this.page.click(`${this.teamDetailsModal} button:has-text("Delete Team")`)
    await this.page.click('.ant-popconfirm-buttons button.ant-btn-primary')
    await this.waitForLoadingComplete()
  }

  // ==================== Table Operations ====================

  /**
   * Get teams count in current tab
   */
  async getTeamsCount(): Promise<number> {
    return await this.getTableRowCount()
  }

  /**
   * Check if team exists in table
   */
  async teamExistsInTable(teamName: string): Promise<boolean> {
    return await this.isVisible(`${this.table} tbody tr:has-text("${teamName}")`)
  }

  /**
   * Get team member count from table
   */
  async getTeamMemberCount(teamName: string): Promise<string | null> {
    const row = this.page.locator(`${this.table} tbody tr:has-text("${teamName}")`)
    const membersCell = row.locator('td').nth(2) // Assuming members is 3rd column
    return await membersCell.textContent()
  }

  // ==================== Assertions ====================

  /**
   * Assert team is in table
   */
  async expectTeamInTable(teamName: string): Promise<void> {
    await expect(
      this.page.locator(`${this.table} tbody tr:has-text("${teamName}")`),
      `Expected "${teamName}" to be in teams table`
    ).toBeVisible()
  }

  /**
   * Assert team is NOT in table
   */
  async expectTeamNotInTable(teamName: string): Promise<void> {
    await expect(
      this.page.locator(`${this.table} tbody tr:has-text("${teamName}")`),
      `Expected "${teamName}" to NOT be in teams table`
    ).toBeHidden()
  }

  /**
   * Assert team details modal is visible
   */
  async expectTeamDetailsVisible(): Promise<void> {
    await expect(
      this.page.locator(this.teamDetailsModal),
      'Expected team details modal to be visible'
    ).toBeVisible()
  }

  /**
   * Assert tab is active
   */
  async expectTabActive(tabName: 'My Teams' | 'All Teams' | 'Join Team'): Promise<void> {
    await expect(
      this.page.locator(`.ant-tabs-tab-active:has-text("${tabName}")`),
      `Expected "${tabName}" tab to be active`
    ).toBeVisible()
  }

  /**
   * Assert success notification after team operation
   */
  async expectSuccessNotification(): Promise<void> {
    await this.expectNotification('success')
  }
}

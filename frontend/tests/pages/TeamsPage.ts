import { Locator, expect } from '@playwright/test'
import { BasePage } from './BasePage'

/**
 * Teams Page Object
 */
export class TeamsPage extends BasePage {
  readonly url = '/teams'

  // ==================== Locators ====================

  get createTeamButton(): Locator {
    return this.page.locator('button:has-text("Create Team")')
  }

  get teamCards(): Locator {
    return this.page.locator('.ant-card')
  }

  get joinTeamForm(): Locator {
    return this.page.locator('[data-testid="join-team-form"], form:has(input[placeholder*="code" i])')
  }

  get inviteCodeInput(): Locator {
    return this.page.locator('input[placeholder*="invite" i], input[placeholder*="code" i]')
  }

  get teamDetailsModal(): Locator {
    return this.page.locator('.ant-modal, [data-testid="team-details-modal"]')
  }

  get membersList(): Locator {
    return this.page.locator('.ant-list, .members-list')
  }

  get teamNameInput(): Locator {
    return this.page.locator('input[placeholder*="team name" i], input[placeholder*="name" i]').first()
  }

  get teamDescriptionInput(): Locator {
    return this.page.locator('textarea')
  }

  get joinButton(): Locator {
    return this.page.locator('button:has-text("Join")')
  }

  get leaveTeamButton(): Locator {
    return this.page.locator('button:has-text("Leave")')
  }

  get inviteButton(): Locator {
    return this.page.locator('button:has-text("Invite")')
  }

  get deleteTeamButton(): Locator {
    return this.page.locator('button:has-text("Delete")')
  }

  get submitButton(): Locator {
    return this.page.locator('.ant-modal-footer button.ant-btn-primary')
  }

  get inviteCodeDisplay(): Locator {
    return this.page.locator('.invite-code, code, .ant-typography-copy')
  }

  // ==================== Actions ====================

  /**
   * Click create team button
   */
  async clickCreateTeam(): Promise<void> {
    await this.createTeamButton.click()
    await this.waitForModal()
  }

  /**
   * Create a new team
   */
  async createTeam(name: string, description: string = ''): Promise<void> {
    await this.clickCreateTeam()
    await this.teamNameInput.fill(name)
    if (description) {
      await this.teamDescriptionInput.fill(description)
    }
    await this.submitButton.click()
    await this.waitForModalClosed()
  }

  /**
   * Join team by invite code
   */
  async joinTeamByCode(code: string): Promise<void> {
    await this.inviteCodeInput.fill(code)
    await this.joinButton.click()
  }

  /**
   * Get team card by index
   */
  getTeamCard(index: number): Locator {
    return this.teamCards.nth(index)
  }

  /**
   * Open team details by index
   */
  async openTeamDetails(index: number): Promise<void> {
    await this.getTeamCard(index).click()
    await this.waitForModal()
  }

  /**
   * Leave current team
   */
  async leaveTeam(): Promise<void> {
    await this.leaveTeamButton.click()
    // Confirm in modal
    await this.clickModalOk()
  }

  /**
   * Click invite button to show invite code
   */
  async clickInvite(): Promise<void> {
    await this.inviteButton.click()
  }

  /**
   * Remove a team member (captain action)
   */
  async removeMember(memberIndex: number): Promise<void> {
    await this.membersList.locator('button:has-text("Remove")').nth(memberIndex).click()
    await this.clickModalOk()
  }

  /**
   * Delete team (captain action)
   */
  async deleteTeam(): Promise<void> {
    await this.deleteTeamButton.click()
    await this.clickModalOk()
  }

  // ==================== Assertions ====================

  /**
   * Expect teams list visible
   */
  async expectTeamsListVisible(): Promise<void> {
    await expect(this.teamCards.first()).toBeVisible({ timeout: 10000 })
  }

  /**
   * Expect team created notification
   */
  async expectTeamCreated(name: string): Promise<void> {
    await this.expectNotification('success')
    await expect(this.page.locator(`.ant-card:has-text("${name}")`)).toBeVisible()
  }

  /**
   * Expect member count
   */
  async expectMemberCount(count: number): Promise<void> {
    const members = this.membersList.locator('.ant-list-item')
    await expect(members).toHaveCount(count)
  }

  /**
   * Expect team details modal visible
   */
  async expectTeamDetailsVisible(): Promise<void> {
    await expect(this.teamDetailsModal).toBeVisible()
  }

  /**
   * Expect to be on teams page
   */
  async expectOnTeamsPage(): Promise<void> {
    await expect(this.page).toHaveURL('/teams')
  }

  /**
   * Expect invite code visible
   */
  async expectInviteCodeVisible(): Promise<void> {
    await expect(this.inviteCodeDisplay).toBeVisible()
  }

  /**
   * Expect no teams message
   */
  async expectNoTeams(): Promise<void> {
    await expect(this.page.locator('text=No teams')).toBeVisible()
  }
}

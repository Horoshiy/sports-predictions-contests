import { Page } from '@playwright/test'

export async function waitForPageLoad(page: Page) {
  await page.waitForLoadState('networkidle')
}

export async function fillAntdInput(page: Page, label: string, value: string) {
  const input = page.locator(`label:has-text("${label}")`).locator('..').locator('input')
  await input.fill(value)
}

export async function clickAntdButton(page: Page, text: string) {
  await page.click(`button:has-text("${text}")`)
}

export async function selectAntdOption(page: Page, label: string, option: string) {
  const selectWrapper = page.locator(`label:has-text("${label}")`).locator('..').locator('.ant-select')
  await selectWrapper.click()
  await page.click(`.ant-select-dropdown .ant-select-item:has-text("${option}")`)
}

export async function waitForAntdNotification(page: Page, type: 'success' | 'error' | 'info' | 'warning') {
  await page.waitForSelector(`.ant-notification-notice-${type}`, { timeout: 5000 })
}

export async function getTableRowCount(page: Page, tableSelector: string) {
  const rows = await page.locator(`${tableSelector} tbody tr`).count()
  return rows
}

export async function clearLocalStorage(page: Page) {
  await page.evaluate(() => localStorage.clear())
}

export async function loginUser(page: Page, email: string, password: string) {
  await page.goto('/login')
  await page.fill('input[type="email"]', email)
  await page.fill('input[type="password"]', password)
  await page.click('button:has-text("Login")')
  await page.waitForURL('/contests')
}

import { test as base } from '@playwright/test'
import type { Page } from '@playwright/test'
import { TEST_CONFIG, TIMEOUTS } from '../helpers/test-config'

type AuthFixtures = {
  authenticatedPage: Page
  adminPage: Page
  testUser: { email: string; password: string }
  testAdmin: { email: string; password: string }
}

export const test = base.extend<AuthFixtures>({
  testUser: async ({}, use) => {
    await use(TEST_CONFIG.testUser)
  },

  testAdmin: async ({}, use) => {
    await use(TEST_CONFIG.testAdmin)
  },

  authenticatedPage: async ({ page, testUser }, use) => {
    try {
      // Navigate to login
      await page.goto('/login')
      
      // Fill credentials
      await page.fill('input[type="email"]', testUser.email)
      await page.fill('input[type="password"]', testUser.password)
      
      // Submit login form
      await page.click('button:has-text("Login")')
      
      // Wait for redirect to contests page
      await page.waitForURL('/contests', { timeout: TIMEOUTS.MEDIUM })
    } catch (error) {
      throw new Error(`Failed to authenticate test user (${testUser.email}): ${error.message}`)
    }
    
    await use(page)
    
    // Cleanup: logout after test
    await page.evaluate(() => localStorage.clear())
  },

  adminPage: async ({ page, testAdmin }, use) => {
    try {
      // Navigate to login
      await page.goto('/login')
      
      // Fill admin credentials
      await page.fill('input[type="email"]', testAdmin.email)
      await page.fill('input[type="password"]', testAdmin.password)
      
      // Submit login form
      await page.click('button:has-text("Login")')
      
      // Wait for redirect
      await page.waitForURL('/contests', { timeout: TIMEOUTS.MEDIUM })
    } catch (error) {
      throw new Error(`Failed to authenticate admin user (${testAdmin.email}): ${error.message}`)
    }
    
    await use(page)
    
    // Cleanup
    await page.evaluate(() => localStorage.clear())
  },
})

export { expect } from '@playwright/test'

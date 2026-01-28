import { test as base } from '@playwright/test'
import type { Page } from '@playwright/test'

type APIFixtures = {
  mockedAPI: Page
  mockContests: any[]
  mockPredictions: any[]
}

export const test = base.extend<APIFixtures>({
  mockContests: async ({}, use) => {
    await use([
      {
        id: 1,
        title: 'Mock Contest 1',
        description: 'Test contest',
        sportType: 'Football',
        startDate: new Date(Date.now() + 86400000).toISOString(),
        endDate: new Date(Date.now() + 604800000).toISOString(),
        maxParticipants: 100,
        participantCount: 10,
      },
    ])
  },

  mockPredictions: async ({}, use) => {
    await use([
      {
        id: 1,
        contestId: 1,
        eventId: 1,
        prediction: 'Team A wins',
        points: 10,
        status: 'pending',
      },
    ])
  },

  mockedAPI: async ({ page, mockContests, mockPredictions }, use) => {
    // Mock successful login
    await page.route('**/v1/auth/login', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          token: 'mock-jwt-token',
          user: {
            id: 1,
            email: 'test@example.com',
            username: 'testuser',
            displayName: 'Test User',
          },
        }),
      })
    })

    // Mock contests list
    await page.route('**/v1/contests', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ contests: mockContests }),
      })
    })

    // Mock predictions list
    await page.route('**/v1/predictions', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ predictions: mockPredictions }),
      })
    })

    // Mock user profile
    await page.route('**/v1/users/me', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 1,
          email: 'test@example.com',
          username: 'testuser',
          displayName: 'Test User',
        }),
      })
    })

    await use(page)
  },
})

export { expect } from '@playwright/test'

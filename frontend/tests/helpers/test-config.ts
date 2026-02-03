export const TEST_CONFIG = {
  testUser: {
    email: process.env.TEST_USER_EMAIL || 'admin@example.com',
    password: process.env.TEST_USER_PASSWORD || 'admin123',
  },
  testAdmin: {
    email: process.env.TEST_ADMIN_EMAIL || 'admin@example.com',
    password: process.env.TEST_ADMIN_PASSWORD || 'admin123',
  },
}

export const TIMEOUTS = {
  SHORT: 5000,
  MEDIUM: 10000,
  LONG: 30000,
}

/**
 * Test configuration
 */
export const TEST_CONFIG = {
  // Regular test user (non-admin)
  testUser: {
    email: process.env.TEST_USER_EMAIL || 'testuser@example.com',
    password: process.env.TEST_USER_PASSWORD || 'TestUser123!',
    name: 'Test User',
  },
  // Admin user for admin-only features
  testAdmin: {
    email: process.env.TEST_ADMIN_EMAIL || 'admin@example.com',
    password: process.env.TEST_ADMIN_PASSWORD || 'admin123',
    name: 'Admin User',
  },
  // Base URL for tests
  baseUrl: process.env.BASE_URL || 'http://localhost:3000',
}

/**
 * Centralized timeout constants
 */
export const TIMEOUTS = {
  /** Short timeout for quick operations (3s) */
  SHORT: 3000,
  /** Medium timeout for standard operations (5s) */
  MEDIUM: 5000,
  /** Long timeout for slow operations (10s) */
  LONG: 10000,
  /** Network timeout for API calls (30s) */
  NETWORK: 30000,
  /** Page load timeout (60s) */
  PAGE_LOAD: 60000,
}

/**
 * Data test IDs for stable selectors
 * These should match data-testid attributes in React components
 */
export const TEST_IDS = {
  // Auth
  loginForm: 'login-form',
  registerForm: 'register-form',
  emailInput: 'email-input',
  passwordInput: 'password-input',
  loginButton: 'login-button',
  registerButton: 'register-button',
  
  // Navigation
  header: 'app-header',
  navContests: 'nav-contests',
  navPredictions: 'nav-predictions',
  navTeams: 'nav-teams',
  navSports: 'nav-sports',
  navAnalytics: 'nav-analytics',
  userMenu: 'user-menu',
  
  // Contests
  contestList: 'contest-list',
  contestCard: 'contest-card',
  contestFilter: 'contest-filter',
  createContestBtn: 'create-contest-btn',
  
  // Predictions
  eventList: 'event-list',
  eventCard: 'event-card',
  predictionForm: 'prediction-form',
  
  // Teams
  teamList: 'team-list',
  teamCard: 'team-card',
  createTeamBtn: 'create-team-btn',
  joinTeamForm: 'join-team-form',
  
  // Profile
  profileForm: 'profile-form',
  displayNameInput: 'display-name-input',
  
  // Common
  modal: 'modal',
  modalOk: 'modal-ok',
  modalCancel: 'modal-cancel',
}

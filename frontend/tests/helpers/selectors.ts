export const SELECTORS = {
  auth: {
    emailInput: 'input[type="email"]',
    passwordInput: 'input[type="password"]',
    loginButton: 'button:has-text("Sign In")',
    registerButton: 'button:has-text("Sign Up"), button:has-text("Register")',
    logoutButton: 'button:has-text("Logout")',
    usernameInput: 'input[placeholder*="username"]',
    displayNameInput: 'input[placeholder*="display name"]',
  },
  
  navigation: {
    contestsLink: 'a[href="/contests"]',
    predictionsLink: 'a[href="/predictions"]',
    teamsLink: 'a[href="/teams"]',
    analyticsLink: 'a[href="/analytics"]',
    profileLink: 'a[href="/profile"]',
    sportsLink: 'a[href="/sports"]',
    header: 'header',
    menu: '.ant-menu',
  },
  
  contests: {
    createButton: 'button:has-text("Create Contest")',
    contestCard: '.ant-card',
    joinButton: 'button:has-text("Join")',
    leaveButton: 'button:has-text("Leave")',
    titleInput: 'input[placeholder*="title"]',
    descriptionInput: 'textarea[placeholder*="description"]',
    sportTypeSelect: '.ant-select',
    startDatePicker: 'input[placeholder*="Start"]',
    endDatePicker: 'input[placeholder*="End"]',
    maxParticipantsInput: 'input[type="number"]',
  },
  
  predictions: {
    submitButton: 'button:has-text("Submit")',
    predictionForm: 'form',
    eventCard: '.ant-card',
    predictionInput: 'input[type="text"]',
    scoreInput: 'input[type="number"]',
  },
  
  teams: {
    createTeamButton: 'button:has-text("Create Team")',
    teamCard: '.ant-card',
    inviteButton: 'button:has-text("Invite")',
    leaveTeamButton: 'button:has-text("Leave Team")',
    teamNameInput: 'input[placeholder*="team name"]',
    membersList: '.ant-list',
    teamDetailsModal: '[data-testid="team-details-modal"]',
    teamLeaderboardTab: 'div[data-testid="team-leaderboard"]',
    joinTeamForm: '[data-testid="join-team-form"]',
    inviteCodeInput: 'input[placeholder*="invite code"]',
    deleteTeamButton: 'button:has-text("Delete Team")',
  },
  
  analytics: {
    dashboard: '.analytics-dashboard',
    chart: '.recharts-wrapper',
    filterButton: 'button:has-text("Filter")',
    dateRangePicker: '.ant-picker-range',
    statsCard: '.ant-statistic',
  },
  
  profile: {
    displayNameInput: 'input[placeholder*="display name"]',
    emailDisplay: '.ant-descriptions-item:has-text("Email")',
    saveButton: 'button:has-text("Save")',
    predictionHistory: '.prediction-history',
    achievements: '.achievements',
  },
  
  common: {
    notification: '.ant-notification',
    notificationSuccess: '.ant-notification-notice-success',
    notificationError: '.ant-notification-notice-error',
    modal: '.ant-modal',
    modalOk: '.ant-modal-footer button:has-text("OK")',
    modalCancel: '.ant-modal-footer button:has-text("Cancel")',
    table: '.ant-table',
    pagination: '.ant-pagination',
    loading: '.ant-spin',
  },
}

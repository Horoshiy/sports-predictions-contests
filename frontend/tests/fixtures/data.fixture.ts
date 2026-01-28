export function generateUser() {
  const timestamp = Date.now()
  return {
    email: `test-${timestamp}@example.com`,
    password: 'TestPass123!',
    username: `user_${timestamp}`,
    displayName: `Test User ${timestamp}`,
  }
}

export function generateContest() {
  const now = new Date()
  const startDate = new Date(now.getTime() + 24 * 60 * 60 * 1000) // Tomorrow
  const endDate = new Date(startDate.getTime() + 7 * 24 * 60 * 60 * 1000) // 7 days later
  
  return {
    title: `Test Contest ${Date.now()}`,
    description: 'Automated test contest for E2E testing',
    sportType: 'Football',
    startDate: startDate.toISOString(),
    endDate: endDate.toISOString(),
    maxParticipants: 100,
    rules: JSON.stringify({ scoring: 'standard' }),
  }
}

export function generatePrediction() {
  return {
    eventId: 1,
    prediction: 'Team A wins',
    score: '2-1',
    confidence: 80,
  }
}

export function generateTeam() {
  const timestamp = Date.now()
  return {
    name: `Test Team ${timestamp}`,
    description: 'Automated test team',
    maxMembers: 10,
  }
}

export function generateChallenge() {
  const now = new Date()
  const deadline = new Date(now.getTime() + 48 * 60 * 60 * 1000) // 2 days
  
  return {
    opponentId: 2,
    eventId: 1,
    wager: 100,
    deadline: deadline.toISOString(),
  }
}

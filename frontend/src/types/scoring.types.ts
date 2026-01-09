// Scoring types matching backend proto definitions

export interface Score {
  id: number
  userId: number
  contestId: number
  predictionId: number
  points: number
  scoredAt: string // ISO string
  createdAt: string // ISO string
  updatedAt: string // ISO string
}

export interface LeaderboardEntry {
  userId: number
  userName: string
  totalPoints: number
  rank: number
  updatedAt: string // ISO string
}

export interface Leaderboard {
  contestId: number
  entries: LeaderboardEntry[]
  updatedAt: string // ISO string
}

// Request types
export interface CreateScoreRequest {
  userId: number
  contestId: number
  predictionId: number
  points: number
}

export interface UpdateScoreRequest {
  id: number
  points: number
}

export interface GetScoreRequest {
  id: number
}

export interface DeleteScoreRequest {
  id: number
}

export interface ListScoresRequest {
  pagination?: PaginationRequest
  contestId?: number
  userId?: number
}

export interface GetUserScoresRequest {
  userId: number
  contestId: number
}

export interface GetLeaderboardRequest {
  contestId: number
  limit?: number // Number of top entries to return
}

export interface GetUserRankRequest {
  contestId: number
  userId: number
}

export interface UpdateLeaderboardRequest {
  contestId: number
}

export interface CalculateScoreRequest {
  predictionId: number
  predictionData: string // JSON string
  resultData: string // JSON string
}

// Response types
export interface CreateScoreResponse {
  response: ApiResponse
  score: Score
}

export interface UpdateScoreResponse {
  response: ApiResponse
  score: Score
}

export interface GetScoreResponse {
  response: ApiResponse
  score: Score
}

export interface DeleteScoreResponse {
  response: ApiResponse
}

export interface ListScoresResponse {
  response: ApiResponse
  scores: Score[]
  pagination: PaginationResponse
}

export interface GetUserScoresResponse {
  response: ApiResponse
  scores: Score[]
  totalPoints: number
}

export interface GetLeaderboardResponse {
  response: ApiResponse
  leaderboard: Leaderboard
}

export interface GetUserRankResponse {
  response: ApiResponse
  rank: number
  totalPoints: number
}

export interface UpdateLeaderboardResponse {
  response: ApiResponse
  leaderboard: Leaderboard
}

export interface CalculateScoreResponse {
  response: ApiResponse
  points: number
  calculationDetails: string // JSON string with calculation breakdown
}

// Common types
export interface PaginationRequest {
  page: number
  limit: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface PaginationResponse {
  page: number
  limit: number
  total: number
  totalPages: number
}

export interface ApiResponse {
  success: boolean
  message: string
  code: number
  timestamp: string // ISO string
}

// Prediction data structures for scoring
export interface PredictionData {
  type: 'exact_score' | 'winner' | 'over_under'
  homeScore?: number // For exact score predictions
  awayScore?: number // For exact score predictions
  winner?: 'home' | 'away' | 'draw' // For winner predictions
  overUnder?: 'over' | 'under' // For over/under predictions
  threshold?: number // For over/under predictions
  value?: any // Generic value for other prediction types
}

export interface ResultData {
  homeScore: number
  awayScore: number
  winner: 'home' | 'away' | 'draw'
  totalGoals: number
}

export interface CalculationDetails {
  predictionType: string
  result: ResultData
  predictedScore?: string
  actualScore?: string
  matchType?: 'exact' | 'goal_difference' | 'winner' | 'none'
  predictedWinner?: string
  actualWinner?: string
  predicted?: string
  threshold?: number
  totalGoals?: number
  correct?: boolean
  match?: boolean
  error?: string
}

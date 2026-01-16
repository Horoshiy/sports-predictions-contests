// Analytics types matching backend proto definitions

export interface SportAccuracy {
  sportType: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
  totalPoints: number
}

export interface LeagueAccuracy {
  leagueId: number
  leagueName: string
  sportType: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
}

export interface PredictionTypeAccuracy {
  predictionType: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
  averagePoints: number
}

export interface AccuracyTrend {
  period: string
  totalPredictions: number
  correctPredictions: number
  accuracyPercentage: number
  totalPoints: number
}

export interface PlatformStats {
  averageAccuracy: number
  averagePointsPerPrediction: number
  totalUsers: number
  totalPredictions: number
}

export interface UserAnalytics {
  userId: number
  totalPredictions: number
  correctPredictions: number
  overallAccuracy: number
  totalPoints: number
  bySport: SportAccuracy[]
  byLeague: LeagueAccuracy[]
  byType: PredictionTypeAccuracy[]
  trends: AccuracyTrend[]
  platformComparison: PlatformStats | null
  timeRange: string
}

export type TimeRange = '7d' | '30d' | '90d' | 'all'

export interface GetUserAnalyticsResponse {
  response: {
    success: boolean
    message: string
    code: number
    timestamp: string
  }
  analytics: UserAnalytics
}

export interface ExportAnalyticsResponse {
  response: {
    success: boolean
    message: string
    code: number
    timestamp: string
  }
  data: string
  filename: string
}

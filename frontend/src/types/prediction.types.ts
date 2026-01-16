// Prediction types matching backend proto definitions
import type { PaginationRequest, PaginationResponse, ApiResponse } from './common.types'

// Prediction entity
export interface Prediction {
  id: number
  contestId: number
  userId: number
  eventId: number
  predictionData: string // JSON string for flexible prediction data
  status: 'pending' | 'scored' | 'cancelled'
  submittedAt: string
  createdAt: string
  updatedAt: string
}

// Event entity
export interface Event {
  id: number
  title: string
  sportType: string
  homeTeam: string
  awayTeam: string
  eventDate: string
  status: 'scheduled' | 'live' | 'completed' | 'cancelled'
  resultData: string
  createdAt: string
  updatedAt: string
}

// Request types
export interface SubmitPredictionRequest {
  contestId: number
  eventId: number
  predictionData: string
}

export interface GetUserPredictionsRequest {
  contestId: number
  pagination?: PaginationRequest
}

export interface UpdatePredictionRequest {
  id: number
  predictionData: string
}

export interface ListEventsRequest {
  sportType?: string
  status?: string
  pagination?: PaginationRequest
}

// Response types
export interface SubmitPredictionResponse {
  response: ApiResponse
  prediction: Prediction
}

export interface GetPredictionResponse {
  response: ApiResponse
  prediction: Prediction
}

export interface GetUserPredictionsResponse {
  response: ApiResponse
  predictions: Prediction[]
  pagination: PaginationResponse
}

export interface UpdatePredictionResponse {
  response: ApiResponse
  prediction: Prediction
}

export interface DeletePredictionResponse {
  response: ApiResponse
}

export interface GetEventResponse {
  response: ApiResponse
  event: Event
}

export interface ListEventsResponse {
  response: ApiResponse
  events: Event[]
  pagination: PaginationResponse
}

// Parsed prediction data types
export interface WinnerPrediction {
  winner: 'home' | 'away' | 'draw'
}

export interface ScorePrediction {
  homeScore: number
  awayScore: number
}

export interface CombinedPrediction extends WinnerPrediction, ScorePrediction {}

export type ParsedPredictionData = WinnerPrediction | ScorePrediction | CombinedPrediction

// Time coefficient types
export interface PotentialCoefficientResponse {
  response: ApiResponse
  coefficient: number
  tier: string
  hoursUntilEvent: number
}

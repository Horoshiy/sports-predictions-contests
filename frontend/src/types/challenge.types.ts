// Challenge types matching backend proto definitions

export interface Challenge {
  id: number
  challengerId: number
  opponentId: number
  eventId: number
  message: string
  status: 'pending' | 'accepted' | 'declined' | 'expired' | 'active' | 'completed'
  expiresAt: string // ISO string
  acceptedAt?: string // ISO string
  completedAt?: string // ISO string
  winnerId?: number
  challengerScore: number
  opponentScore: number
  createdAt: string // ISO string
  updatedAt: string // ISO string
}

export interface ChallengeParticipant {
  id: number
  challengeId: number
  userId: number
  role: 'challenger' | 'opponent'
  status: 'active' | 'inactive'
  joinedAt: string // ISO string
}

// Request types
export interface CreateChallengeRequest {
  opponentId: number
  eventId: number
  message: string
}

export interface AcceptChallengeRequest {
  id: number
}

export interface DeclineChallengeRequest {
  id: number
}

export interface WithdrawChallengeRequest {
  id: number
}

export interface GetChallengeRequest {
  id: number
}

export interface ListUserChallengesRequest {
  userId: number
  status?: string
  pagination?: PaginationRequest
}

export interface ListOpenChallengesRequest {
  eventId: number
  pagination?: PaginationRequest
}

// Response types
export interface ApiResponse {
  success: boolean
  message: string
  code: number
  timestamp: string
}

export interface CreateChallengeResponse {
  response: ApiResponse
  challenge: Challenge
}

export interface AcceptChallengeResponse {
  response: ApiResponse
  challenge: Challenge
}

export interface DeclineChallengeResponse {
  response: ApiResponse
}

export interface WithdrawChallengeResponse {
  response: ApiResponse
}

export interface GetChallengeResponse {
  response: ApiResponse
  challenge: Challenge
}

export interface ListUserChallengesResponse {
  response: ApiResponse
  challenges: Challenge[]
  pagination: PaginationResponse
}

export interface ListOpenChallengesResponse {
  response: ApiResponse
  challenges: Challenge[]
  pagination: PaginationResponse
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

// Error types
export interface ErrorResponse {
  error: string
  code: number
  message: string
}

// Form types
export interface ChallengeFormData {
  opponentId: number
  eventId: number
  message: string
}

// UI helper types
export interface ChallengeWithUserInfo extends Challenge {
  challengerName?: string
  opponentName?: string
  eventTitle?: string
}

export interface ChallengeStatus {
  label: string
  color: 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning'
  description: string
}

export const CHALLENGE_STATUSES: Record<Challenge['status'], ChallengeStatus> = {
  pending: {
    label: 'Pending',
    color: 'warning',
    description: 'Waiting for opponent response'
  },
  accepted: {
    label: 'Accepted',
    color: 'info',
    description: 'Challenge accepted, waiting for event'
  },
  declined: {
    label: 'Declined',
    color: 'error',
    description: 'Challenge declined by opponent'
  },
  expired: {
    label: 'Expired',
    color: 'default',
    description: 'Challenge expired without response'
  },
  active: {
    label: 'Active',
    color: 'primary',
    description: 'Challenge is currently active'
  },
  completed: {
    label: 'Completed',
    color: 'success',
    description: 'Challenge completed with results'
  }
}

// Contest types matching backend proto definitions

export interface Contest {
  id: number
  title: string
  description: string
  sportType: string
  rules: string // JSON string for flexible rule configuration
  status: 'draft' | 'active' | 'completed' | 'cancelled'
  startDate: string // ISO string
  endDate: string // ISO string
  maxParticipants: number
  currentParticipants: number
  creatorId: number
  createdAt: string // ISO string
  updatedAt: string // ISO string
}

export interface Participant {
  id: number
  contestId: number
  userId: number
  role: 'admin' | 'participant'
  status: 'active' | 'inactive' | 'banned'
  joinedAt: string // ISO string
}

// Request types
export interface CreateContestRequest {
  title: string
  description: string
  sportType: string
  rules: string
  startDate: string // ISO string
  endDate: string // ISO string
  maxParticipants: number
}

export interface UpdateContestRequest {
  id: number
  title: string
  description: string
  sportType: string
  rules: string
  startDate: string // ISO string
  endDate: string // ISO string
  maxParticipants: number
  status: 'draft' | 'active' | 'completed' | 'cancelled'
}

export interface ListContestsRequest {
  pagination?: PaginationRequest
  status?: string
  sportType?: string
}

export interface JoinContestRequest {
  contestId: number
}

export interface LeaveContestRequest {
  contestId: number
}

export interface ListParticipantsRequest {
  contestId: number
  pagination?: PaginationRequest
}

// Response types
export interface ApiResponse {
  success: boolean
  message: string
  code: number
  timestamp: string
}

export interface CreateContestResponse {
  response: ApiResponse
  contest: Contest
}

export interface UpdateContestResponse {
  response: ApiResponse
  contest: Contest
}

export interface GetContestResponse {
  response: ApiResponse
  contest: Contest
}

export interface DeleteContestResponse {
  response: ApiResponse
}

export interface ListContestsResponse {
  response: ApiResponse
  contests: Contest[]
  pagination: PaginationResponse
}

export interface JoinContestResponse {
  response: ApiResponse
  participant: Participant
}

export interface LeaveContestResponse {
  response: ApiResponse
}

export interface ListParticipantsResponse {
  response: ApiResponse
  participants: Participant[]
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
export interface ContestFormData {
  title: string
  description?: string
  sportType: string
  rules?: string
  startDate: Date
  endDate: Date
  maxParticipants: number
}

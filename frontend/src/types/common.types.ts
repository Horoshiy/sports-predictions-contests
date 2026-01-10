// Common types shared across the application

export interface ApiResponse {
  success: boolean
  message: string
  code: number
  timestamp: string
}

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

export interface ErrorResponse {
  error: string
  code: number
  message: string
}

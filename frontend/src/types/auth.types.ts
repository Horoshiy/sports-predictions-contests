// Authentication types matching backend proto definitions

import type { ApiResponse } from './common.types'

export interface AuthUser {
  id: number
  email: string
  name: string
  createdAt: string // ISO string
  updatedAt: string // ISO string
}

// Request types
export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface GetProfileRequest {
  // User ID will be extracted from JWT token
}

export interface UpdateProfileRequest {
  name: string
  email: string
}

// Response types
export interface RegisterResponse {
  response: ApiResponse
  user: AuthUser
  token: string
}

export interface LoginResponse {
  response: ApiResponse
  user: AuthUser
  token: string
}

export interface GetProfileResponse {
  response: ApiResponse
  user: AuthUser
}

export interface UpdateProfileResponse {
  response: ApiResponse
  user: AuthUser
}

// Form types
export interface LoginFormData {
  email: string
  password: string
}

export interface RegisterFormData {
  email: string
  password: string
  confirmPassword: string
  name: string
}

export interface UpdateProfileFormData {
  name: string
  email: string
}

// Auth context types
export interface AuthContextType {
  user: AuthUser | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => void
  updateProfile: (data: UpdateProfileFormData) => Promise<void>
}

// Error types
export interface AuthError {
  message: string
  code?: number
  field?: string
}

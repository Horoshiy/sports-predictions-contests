// Profile types matching backend proto definitions

import type { ApiResponse } from './common.types'

export interface Profile {
  id: number
  userId: number
  bio: string
  avatarUrl: string
  location: string
  website: string
  twitterUrl: string
  linkedinUrl: string
  githubUrl: string
  profileVisibility: 'public' | 'friends' | 'private'
  createdAt: string // ISO string
  updatedAt: string // ISO string
}

export interface UserPreferences {
  id: number
  userId: number
  emailNotifications: boolean
  pushNotifications: boolean
  contestNotifications: boolean
  predictionReminders: boolean
  weeklyDigest: boolean
  theme: 'light' | 'dark' | 'auto'
  language: 'en' | 'ru' | 'es' | 'fr' | 'de'
  timezone: string
  customSettings: Record<string, any>
  createdAt: string // ISO string
  updatedAt: string // ISO string
}

// Request types
export interface GetProfileRequest {
  // User ID will be extracted from JWT token
}

export interface UpdateProfileRequest {
  bio: string
  location: string
  website: string
  twitterUrl: string
  linkedinUrl: string
  githubUrl: string
  profileVisibility: 'public' | 'friends' | 'private'
}

export interface UploadAvatarRequest {
  file: File
}

export interface GetPreferencesRequest {
  // User ID will be extracted from JWT token
}

export interface UpdatePreferencesRequest {
  emailNotifications: boolean
  pushNotifications: boolean
  contestNotifications: boolean
  predictionReminders: boolean
  weeklyDigest: boolean
  theme: 'light' | 'dark' | 'auto'
  language: 'en' | 'ru' | 'es' | 'fr' | 'de'
  timezone: string
  customSettings: Record<string, any>
}

export interface GetProfileCompletionRequest {
  // User ID will be extracted from JWT token
}

// Response types
export interface GetProfileResponse {
  response: ApiResponse
  profile: Profile
}

export interface UpdateProfileResponse {
  response: ApiResponse
  profile: Profile
}

export interface UploadAvatarResponse {
  response: ApiResponse
  avatarUrl: string
}

export interface GetPreferencesResponse {
  response: ApiResponse
  preferences: UserPreferences
}

export interface UpdatePreferencesResponse {
  response: ApiResponse
  preferences: UserPreferences
}

export interface GetProfileCompletionResponse {
  response: ApiResponse
  completionPercentage: number
  missingFields: string[]
  suggestions: string[]
}

// Form types
export interface ProfileFormData {
  bio: string
  location: string
  website: string
  twitterUrl: string
  linkedinUrl: string
  githubUrl: string
  profileVisibility: 'public' | 'friends' | 'private'
}

export interface PreferencesFormData {
  emailNotifications: boolean
  pushNotifications: boolean
  contestNotifications: boolean
  predictionReminders: boolean
  weeklyDigest: boolean
  theme: 'light' | 'dark' | 'auto'
  language: 'en' | 'ru' | 'es' | 'fr' | 'de'
  timezone: string
}

// Profile completion types
export interface ProfileCompletion {
  percentage: number
  missingFields: string[]
  suggestions: string[]
}

// Avatar upload types
export interface AvatarUploadState {
  isUploading: boolean
  progress: number
  error: string | null
}

// Privacy visibility options
export const PROFILE_VISIBILITY_OPTIONS = [
  { value: 'public', label: 'Public - Visible to everyone' },
  { value: 'friends', label: 'Friends Only - Visible to friends' },
  { value: 'private', label: 'Private - Only visible to you' }
] as const

// Theme options
export const THEME_OPTIONS = [
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' },
  { value: 'auto', label: 'Auto (System)' }
] as const

// Language options
export const LANGUAGE_OPTIONS = [
  { value: 'en', label: 'English' },
  { value: 'ru', label: 'Русский' },
  { value: 'es', label: 'Español' },
  { value: 'fr', label: 'Français' },
  { value: 'de', label: 'Deutsch' }
] as const

// Error types
export interface ProfileError {
  message: string
  code?: number
  field?: string
}

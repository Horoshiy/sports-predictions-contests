import grpcClient from './grpc-client'
import type {
  Profile,
  UserPreferences,
  ProfileFormData,
  PreferencesFormData,
  GetProfileResponse,
  UpdateProfileResponse,
  UploadAvatarResponse,
  GetPreferencesResponse,
  UpdatePreferencesResponse,
  GetProfileCompletionResponse,
  ProfileCompletion,
} from '../types/profile.types'

class ProfileService {
  private basePath = '/v1/profile'

  // Validate API response structure
  private validateResponse(response: any): boolean {
    return response && 
           typeof response === 'object' && 
           response.response && 
           typeof response.response.success === 'boolean'
  }

  // Get user profile
  async getProfile(): Promise<Profile> {
    const response = await grpcClient.get<GetProfileResponse>(this.basePath)

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to get profile')
    }

    return response.profile
  }

  // Update user profile
  async updateProfile(profileData: ProfileFormData): Promise<Profile> {
    const response = await grpcClient.put<UpdateProfileResponse>(
      this.basePath,
      profileData
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to update profile')
    }

    return response.profile
  }

  // Upload avatar
  async uploadAvatar(file: File): Promise<string> {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('fileType', file.type)
    formData.append('fileSize', file.size.toString())

    const response = await grpcClient.postFormData<UploadAvatarResponse>(
      `${this.basePath}/avatar`,
      formData
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to upload avatar')
    }

    return response.avatarUrl
  }

  // Get user preferences
  async getPreferences(): Promise<UserPreferences> {
    const response = await grpcClient.get<GetPreferencesResponse>(
      `${this.basePath}/preferences`
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to get preferences')
    }

    return response.preferences
  }

  // Update user preferences
  async updatePreferences(preferencesData: PreferencesFormData): Promise<UserPreferences> {
    const response = await grpcClient.put<UpdatePreferencesResponse>(
      `${this.basePath}/preferences`,
      preferencesData
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to update preferences')
    }

    return response.preferences
  }

  // Get profile completion status
  async getProfileCompletion(): Promise<ProfileCompletion> {
    const response = await grpcClient.get<GetProfileCompletionResponse>(
      `${this.basePath}/completion`
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to get profile completion')
    }

    return {
      percentage: response.completionPercentage,
      missingFields: response.missingFields,
      suggestions: response.suggestions,
    }
  }

  // Validate file for avatar upload
  validateAvatarFile(file: File): { isValid: boolean; error?: string } {
    const maxSize = 5 * 1024 * 1024 // 5MB
    const allowedTypes = ['image/jpeg', 'image/png', 'image/gif']

    if (file.size > maxSize) {
      return {
        isValid: false,
        error: 'File size must be less than 5MB'
      }
    }

    if (!allowedTypes.includes(file.type)) {
      return {
        isValid: false,
        error: 'File must be a JPEG, PNG, or GIF image'
      }
    }

    return { isValid: true }
  }

  // Generate avatar preview URL
  generateAvatarPreview(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = (e) => {
        if (e.target?.result) {
          resolve(e.target.result as string)
        } else {
          reject(new Error('Failed to read file'))
        }
      }
      reader.onerror = () => reject(new Error('Failed to read file'))
      reader.readAsDataURL(file)
    })
  }

  // Get default profile data
  getDefaultProfile(): Partial<Profile> {
    return {
      bio: '',
      location: '',
      website: '',
      twitterUrl: '',
      linkedinUrl: '',
      githubUrl: '',
      profileVisibility: 'public',
      avatarUrl: '',
    }
  }

  // Get default preferences
  getDefaultPreferences(): Partial<UserPreferences> {
    return {
      emailNotifications: true,
      pushNotifications: true,
      contestNotifications: true,
      predictionReminders: true,
      weeklyDigest: true,
      theme: 'light',
      language: 'en',
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
      customSettings: {},
    }
  }
}

export const profileService = new ProfileService()
export default profileService

import grpcClient from './grpc-client'
import type {
  AuthUser,
  RegisterRequest,
  LoginRequest,
  UpdateProfileRequest,
  RegisterResponse,
  LoginResponse,
  GetProfileResponse,
  UpdateProfileResponse,
} from '../types/auth.types'

class AuthService {
  private basePath = '/v1/auth'
  private userPath = '/v1/users'

  // Validate API response structure
  private validateResponse(response: any): boolean {
    return response && 
           typeof response === 'object' && 
           response.response && 
           typeof response.response.success === 'boolean'
  }

  // Register a new user
  async register(email: string, password: string, name: string): Promise<{ user: AuthUser; token: string }> {
    const request: RegisterRequest = {
      email,
      password,
      name,
    }

    const response = await grpcClient.post<RegisterResponse>(
      `${this.basePath}/register`,
      request
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Registration failed')
    }

    return {
      user: response.user,
      token: response.token,
    }
  }

  // Login user
  async login(email: string, password: string): Promise<{ user: AuthUser; token: string }> {
    const request: LoginRequest = {
      email,
      password,
    }

    const response = await grpcClient.post<LoginResponse>(
      `${this.basePath}/login`,
      request
    )

    console.log('Login response:', response)

    if (!this.validateResponse(response)) {
      console.error('Invalid response structure:', response)
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Login failed')
    }

    if (!response.token || !response.user) {
      console.error('Missing token or user in response:', response)
      throw new Error('Invalid login response: missing token or user data')
    }

    return {
      user: response.user,
      token: response.token,
    }
  }

  // Get user profile
  async getProfile(): Promise<AuthUser> {
    const response = await grpcClient.get<GetProfileResponse>(
      `${this.userPath}/profile`
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to get profile')
    }

    return response.user
  }

  // Update user profile
  async updateProfile(name: string, email: string): Promise<AuthUser> {
    const request: UpdateProfileRequest = {
      name,
      email,
    }

    const response = await grpcClient.put<UpdateProfileResponse>(
      `${this.userPath}/profile`,
      request
    )

    if (!this.validateResponse(response)) {
      throw new Error('Invalid API response format')
    }

    if (!response.response.success) {
      throw new Error(response.response.message || 'Failed to update profile')
    }

    return response.user
  }

  // Verify token validity
  async verifyToken(): Promise<AuthUser> {
    try {
      return await this.getProfile()
    } catch (error) {
      throw new Error('Token is invalid or expired')
    }
  }
}

// Singleton instance
export const authService = new AuthService()
export default authService

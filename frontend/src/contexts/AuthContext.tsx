import React, { createContext, useContext, useState, useEffect, useRef, ReactNode } from 'react'
import { authService } from '../services/auth-service'
import { useToast } from './ToastContext'
import type { AuthUser, AuthContextType, UpdateProfileFormData } from '../types/auth.types'

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const { showToast } = useToast()
  const isMountedRef = useRef(true)

  const isAuthenticated = !!user

  // Debug logging
  useEffect(() => {
    console.log('AuthProvider state changed:', { user, isAuthenticated, isLoading })
  }, [user, isAuthenticated, isLoading])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      isMountedRef.current = false
    }
  }, [])

  // Initialize auth state from stored token
  useEffect(() => {
    const initializeAuth = async () => {
      const token = localStorage.getItem('authToken')
      if (token) {
        try {
          // Add timeout to token verification
          const userData = await Promise.race([
            authService.verifyToken(),
            new Promise<never>((_, reject) => 
              setTimeout(() => reject(new Error('Token verification timeout')), 5000)
            )
          ])
          if (isMountedRef.current) {
            setUser(userData)
          }
        } catch (error) {
          // Token is invalid or verification failed, remove it
          localStorage.removeItem('authToken')
          console.error('Token verification failed:', error)
        }
      }
      // Always set loading to false, even if no token
      if (isMountedRef.current) {
        setIsLoading(false)
      }
    }

    initializeAuth()
  }, [])

  const login = async (email: string, password: string): Promise<void> => {
    try {
      setIsLoading(true)
      console.log('Attempting login for:', email)
      const { user: userData, token } = await authService.login(email, password)
      
      console.log('Login successful, user:', userData)
      console.log('Token received:', token ? 'Yes' : 'No')
      
      // Store token and user data
      localStorage.setItem('authToken', token)
      // Always set user, even if component is unmounting (navigation will happen)
      console.log('Setting user in state:', userData)
      setUser(userData)
      setIsLoading(false) // Set loading to false immediately after setting user
      console.log('User set, isAuthenticated should be:', !!userData)
      showToast('Login successful!', 'success')
    } catch (error) {
      console.error('Login error:', error)
      setIsLoading(false) // Also set loading to false on error
      if (isMountedRef.current) {
        const message = error instanceof Error ? error.message : 'Login failed'
        showToast(message, 'error')
      }
      throw error
    }
  }

  const register = async (email: string, password: string, name: string): Promise<void> => {
    try {
      setIsLoading(true)
      const { user: userData, token } = await authService.register(email, password, name)
      
      // Store token and user data
      localStorage.setItem('authToken', token)
      setUser(userData)
      setIsLoading(false)
      showToast('Registration successful! Welcome!', 'success')
    } catch (error) {
      setIsLoading(false)
      if (isMountedRef.current) {
        const message = error instanceof Error ? error.message : 'Registration failed'
        showToast(message, 'error')
      }
      throw error
    }
  }

  const logout = (): void => {
    localStorage.removeItem('authToken')
    if (isMountedRef.current) {
      setUser(null)
      showToast('Logged out successfully', 'info')
    }
  }

  const updateProfile = async (data: UpdateProfileFormData): Promise<void> => {
    try {
      if (isMountedRef.current) {
        setIsLoading(true)
      }
      const updatedUser = await authService.updateProfile(data.name, data.email)
      if (isMountedRef.current) {
        setUser(updatedUser)
        showToast('Profile updated successfully!', 'success')
      }
    } catch (error) {
      if (isMountedRef.current) {
        const message = error instanceof Error ? error.message : 'Failed to update profile'
        showToast(message, 'error')
      }
      throw error
    } finally {
      if (isMountedRef.current) {
        setIsLoading(false)
      }
    }
  }

  const value: AuthContextType = {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    updateProfile,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

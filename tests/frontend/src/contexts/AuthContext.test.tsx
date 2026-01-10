import React from 'react'
import { render, screen, waitFor, act } from '@testing-library/react'
import { AuthProvider, useAuth } from '../../../../../frontend/src/contexts/AuthContext'
import { ToastProvider } from '../../../../../frontend/src/contexts/ToastContext'
import { authService } from '../../../../../frontend/src/services/auth-service'
import type { AuthUser } from '../../../../../frontend/src/types/auth.types'

// Mock the auth service
jest.mock('../../../../../frontend/src/services/auth-service', () => ({
  authService: {
    login: jest.fn(),
    register: jest.fn(),
    verifyToken: jest.fn(),
    updateProfile: jest.fn(),
  },
}))

// Mock localStorage
const mockLocalStorage = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
}
Object.defineProperty(window, 'localStorage', {
  value: mockLocalStorage,
})

// Test component that uses the auth context
const TestComponent: React.FC = () => {
  const { user, isAuthenticated, isLoading, login, register, logout, updateProfile } = useAuth()

  return (
    <div>
      <div data-testid="user">{user ? user.name : 'No user'}</div>
      <div data-testid="authenticated">{isAuthenticated ? 'true' : 'false'}</div>
      <div data-testid="loading">{isLoading ? 'true' : 'false'}</div>
      <button onClick={() => login('test@example.com', 'password123')}>Login</button>
      <button onClick={() => register('test@example.com', 'password123', 'Test User')}>Register</button>
      <button onClick={logout}>Logout</button>
      <button onClick={() => updateProfile({ name: 'Updated Name', email: 'updated@example.com' })}>
        Update Profile
      </button>
    </div>
  )
}

const renderWithProviders = (component: React.ReactElement) => {
  return render(
    <ToastProvider>
      <AuthProvider>
        {component}
      </AuthProvider>
    </ToastProvider>
  )
}

describe('AuthContext', () => {
  const mockUser: AuthUser = {
    id: 1,
    email: 'test@example.com',
    name: 'Test User',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  }

  beforeEach(() => {
    jest.clearAllMocks()
    mockLocalStorage.getItem.mockReturnValue(null)
  })

  it('initializes with no user when no token in localStorage', async () => {
    mockLocalStorage.getItem.mockReturnValue(null)

    renderWithProviders(<TestComponent />)

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('No user')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
      expect(screen.getByTestId('loading')).toHaveTextContent('false')
    })
  })

  it('verifies token on initialization when token exists', async () => {
    mockLocalStorage.getItem.mockReturnValue('mock-token')
    ;(authService.verifyToken as jest.Mock).mockResolvedValue(mockUser)

    renderWithProviders(<TestComponent />)

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('Test User')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
      expect(screen.getByTestId('loading')).toHaveTextContent('false')
    })

    expect(authService.verifyToken).toHaveBeenCalled()
  })

  it('removes invalid token on initialization', async () => {
    mockLocalStorage.getItem.mockReturnValue('invalid-token')
    ;(authService.verifyToken as jest.Mock).mockRejectedValue(new Error('Invalid token'))

    renderWithProviders(<TestComponent />)

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('No user')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
      expect(screen.getByTestId('loading')).toHaveTextContent('false')
    })

    expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('authToken')
  })

  it('handles successful login', async () => {
    ;(authService.login as jest.Mock).mockResolvedValue({
      user: mockUser,
      token: 'new-token',
    })

    renderWithProviders(<TestComponent />)

    await act(async () => {
      screen.getByText('Login').click()
    })

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('Test User')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
    })

    expect(authService.login).toHaveBeenCalledWith('test@example.com', 'password123')
    expect(mockLocalStorage.setItem).toHaveBeenCalledWith('authToken', 'new-token')
  })

  it('handles failed login', async () => {
    ;(authService.login as jest.Mock).mockRejectedValue(new Error('Invalid credentials'))

    renderWithProviders(<TestComponent />)

    await act(async () => {
      screen.getByText('Login').click()
    })

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('No user')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
    })

    expect(mockLocalStorage.setItem).not.toHaveBeenCalled()
  })

  it('handles successful registration', async () => {
    ;(authService.register as jest.Mock).mockResolvedValue({
      user: mockUser,
      token: 'new-token',
    })

    renderWithProviders(<TestComponent />)

    await act(async () => {
      screen.getByText('Register').click()
    })

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('Test User')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
    })

    expect(authService.register).toHaveBeenCalledWith('test@example.com', 'password123', 'Test User')
    expect(mockLocalStorage.setItem).toHaveBeenCalledWith('authToken', 'new-token')
  })

  it('handles logout', async () => {
    // Start with authenticated user
    mockLocalStorage.getItem.mockReturnValue('mock-token')
    ;(authService.verifyToken as jest.Mock).mockResolvedValue(mockUser)

    renderWithProviders(<TestComponent />)

    await waitFor(() => {
      expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
    })

    await act(async () => {
      screen.getByText('Logout').click()
    })

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('No user')
      expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
    })

    expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('authToken')
  })

  it('handles profile update', async () => {
    const updatedUser = { ...mockUser, name: 'Updated Name', email: 'updated@example.com' }
    
    // Start with authenticated user
    mockLocalStorage.getItem.mockReturnValue('mock-token')
    ;(authService.verifyToken as jest.Mock).mockResolvedValue(mockUser)
    ;(authService.updateProfile as jest.Mock).mockResolvedValue(updatedUser)

    renderWithProviders(<TestComponent />)

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('Test User')
    })

    await act(async () => {
      screen.getByText('Update Profile').click()
    })

    await waitFor(() => {
      expect(screen.getByTestId('user')).toHaveTextContent('Updated Name')
    })

    expect(authService.updateProfile).toHaveBeenCalledWith('Updated Name', 'updated@example.com')
  })

  it('throws error when useAuth is used outside AuthProvider', () => {
    const TestComponentWithoutProvider = () => {
      useAuth()
      return <div>Test</div>
    }

    // Suppress console.error for this test
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation()

    expect(() => render(<TestComponentWithoutProvider />)).toThrow(
      'useAuth must be used within an AuthProvider'
    )

    consoleSpy.mockRestore()
  })
})

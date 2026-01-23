import React, { useEffect, useState } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { Spin, Result, Button, Space } from 'antd'
import { useAuth } from '../../contexts/AuthContext'

interface ProtectedRouteProps {
  children: React.ReactNode
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth()
  const location = useLocation()
  const [hasError, setHasError] = useState(false)

  // Add timeout for loading state
  useEffect(() => {
    if (isLoading) {
      const timeout = setTimeout(() => {
        setHasError(true)
      }, 10000) // 10 second timeout

      return () => clearTimeout(timeout)
    }
  }, [isLoading])

  // Show error state if timeout occurred
  if (hasError) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
        <Result
          status="error"
          title="Authentication check timed out"
          subTitle="Please try refreshing the page or login again."
          extra={
            <Space>
              <Button type="primary" onClick={() => window.location.reload()}>
                Refresh Page
              </Button>
              <Button onClick={() => window.location.href = '/login'}>
                Go to Login
              </Button>
            </Space>
          }
        />
      </div>
    )
  }

  // Show loading spinner while checking authentication
  if (isLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
        <Spin size="large" />
      </div>
    )
  }

  // Redirect to login if not authenticated, preserving the intended destination
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  // Render children if authenticated
  return <>{children}</>
}

export default ProtectedRoute

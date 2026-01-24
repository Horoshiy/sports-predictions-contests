import React from 'react'
import { Space, Typography } from 'antd'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { LoginForm } from '../components/auth/LoginForm'
import { useAuth } from '../contexts/AuthContext'
import type { LoginFormData } from '../utils/auth-validation'

const { Title, Text } = Typography

const LoginPage: React.FC = () => {
  const { login, isLoading } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  // Get the intended destination from location state, default to configurable path
  const defaultPath = import.meta.env.VITE_DEFAULT_REDIRECT || '/contests'
  const from = location.state?.from?.pathname || defaultPath

  const handleLogin = async (data: LoginFormData) => {
    try {
      console.log('handleLogin: Starting login...')
      await login(data.email, data.password)
      console.log('handleLogin: Login successful, navigating to:', from)
      // Redirect to intended destination after successful login
      navigate(from, { replace: true })
      console.log('handleLogin: Navigation complete')
    } catch (error) {
      console.error('handleLogin: Login failed:', error)
      // Error is handled by the auth context (toast notification)
    }
  }

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '24px' }}>
      <Space direction="vertical" size="large" align="center" style={{ width: '100%', maxWidth: '500px' }}>
        <div style={{ textAlign: 'center', marginBottom: '24px' }}>
          <Title level={1}>Sports Prediction Contests</Title>
          <Title level={4} type="secondary" style={{ fontWeight: 'normal' }}>
            Make predictions, compete with friends, and climb the leaderboards!
          </Title>
        </div>

        {/* Development credentials notice */}
        {import.meta.env.DEV && (
          <div style={{ 
            width: '100%', 
            padding: '12px', 
            background: '#e6f7ff', 
            border: '1px solid #91d5ff',
            borderRadius: '4px',
            marginBottom: '16px'
          }}>
            <Text strong style={{ display: 'block', marginBottom: '8px' }}>
              🔑 Development Admin Credentials:
            </Text>
            <Text style={{ display: 'block', fontSize: '13px' }}>
              Email: <Text code>admin@sportsprediction.com</Text>
            </Text>
            <Text style={{ display: 'block', fontSize: '13px' }}>
              Password: <Text code>admin123</Text>
            </Text>
          </div>
        )}

        <LoginForm onSubmit={handleLogin} loading={isLoading} />

        <Text type="secondary">
          Don't have an account?{' '}
          <Link to="/register" style={{ color: '#1976d2' }}>
            Sign up here
          </Link>
        </Text>
      </Space>
    </div>
  )
}

export default LoginPage

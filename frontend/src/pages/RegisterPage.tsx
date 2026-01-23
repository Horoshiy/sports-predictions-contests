import React from 'react'
import { Space, Typography } from 'antd'
import { Link, useNavigate } from 'react-router-dom'
import { RegisterForm } from '../components/auth/RegisterForm'
import { useAuth } from '../contexts/AuthContext'
import type { RegisterFormData } from '../utils/auth-validation'

const { Title, Text } = Typography

const RegisterPage: React.FC = () => {
  const { register, isLoading } = useAuth()
  const navigate = useNavigate()

  const handleRegister = async (data: RegisterFormData) => {
    try {
      await register(data.email, data.password, data.name)
      // Redirect to configurable default page after successful registration
      const defaultPath = import.meta.env.VITE_DEFAULT_REDIRECT || '/contests'
      navigate(defaultPath, { replace: true })
    } catch (error) {
      // Error is handled by the auth context (toast notification)
    }
  }

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '24px' }}>
      <Space direction="vertical" size="large" align="center" style={{ width: '100%', maxWidth: '500px' }}>
        <div style={{ textAlign: 'center', marginBottom: '24px' }}>
          <Title level={1}>Sports Prediction Contests</Title>
          <Title level={4} type="secondary" style={{ fontWeight: 'normal' }}>
            Join the community and start making predictions today!
          </Title>
        </div>

        <RegisterForm onSubmit={handleRegister} loading={isLoading} />

        <Text type="secondary">
          Already have an account?{' '}
          <Link to="/login" style={{ color: '#1976d2' }}>
            Sign in here
          </Link>
        </Text>
      </Space>
    </div>
  )
}

export default RegisterPage

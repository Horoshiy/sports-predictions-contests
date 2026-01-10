import React from 'react'
import { Box, Container, Typography, Link as MuiLink } from '@mui/material'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { LoginForm } from '../components/auth/LoginForm'
import { useAuth } from '../contexts/AuthContext'
import type { LoginFormData } from '../utils/auth-validation'

const LoginPage: React.FC = () => {
  const { login, isLoading } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  // Get the intended destination from location state, default to configurable path
  const defaultPath = import.meta.env.VITE_DEFAULT_REDIRECT || '/contests'
  const from = location.state?.from?.pathname || defaultPath

  const handleLogin = async (data: LoginFormData) => {
    try {
      await login(data.email, data.password)
      // Redirect to intended destination after successful login
      navigate(from, { replace: true })
    } catch (error) {
      // Error is handled by the auth context (toast notification)
      console.error('Login failed:', error)
    }
  }

  return (
    <Container component="main" maxWidth="sm">
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          py: 4,
        }}
      >
        <Box sx={{ mb: 4, textAlign: 'center' }}>
          <Typography variant="h3" component="h1" gutterBottom>
            Sports Prediction Contests
          </Typography>
          <Typography variant="h6" color="text.secondary">
            Make predictions, compete with friends, and climb the leaderboards!
          </Typography>
        </Box>

        <LoginForm onSubmit={handleLogin} loading={isLoading} />

        <Box sx={{ mt: 3, textAlign: 'center' }}>
          <Typography variant="body2" color="text.secondary">
            Don't have an account?{' '}
            <MuiLink component={Link} to="/register" underline="hover">
              Sign up here
            </MuiLink>
          </Typography>
        </Box>
      </Box>
    </Container>
  )
}

export default LoginPage

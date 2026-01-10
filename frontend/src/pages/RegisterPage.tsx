import React from 'react'
import { Box, Container, Typography, Link as MuiLink } from '@mui/material'
import { Link, useNavigate } from 'react-router-dom'
import { RegisterForm } from '../components/auth/RegisterForm'
import { useAuth } from '../contexts/AuthContext'
import type { RegisterFormData } from '../utils/auth-validation'

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
      console.error('Registration failed:', error)
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
            Join the community and start making predictions today!
          </Typography>
        </Box>

        <RegisterForm onSubmit={handleRegister} loading={isLoading} />

        <Box sx={{ mt: 3, textAlign: 'center' }}>
          <Typography variant="body2" color="text.secondary">
            Already have an account?{' '}
            <MuiLink component={Link} to="/login" underline="hover">
              Sign in here
            </MuiLink>
          </Typography>
        </Box>
      </Box>
    </Container>
  )
}

export default RegisterPage

import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate, Link } from 'react-router-dom'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import { 
  CssBaseline, 
  Container, 
  AppBar, 
  Toolbar, 
  Typography, 
  Button, 
  Box,
  Menu,
  MenuItem,
  IconButton,
  Avatar
} from '@mui/material'
import { AccountCircle, ExitToApp } from '@mui/icons-material'
import { ToastProvider } from './contexts/ToastContext'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { ProtectedRoute } from './components/auth/ProtectedRoute'
import ContestsPage from './pages/ContestsPage'
import SportsPage from './pages/SportsPage'
import PredictionsPage from './pages/PredictionsPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
  typography: {
    fontFamily: 'Roboto, Arial, sans-serif',
  },
})

const AppBarContent: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth()
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = () => {
    logout()
    handleClose()
  }

  return (
    <Toolbar>
      <Typography variant="h6" component="div" sx={{ mr: 2 }}>
        Sports Prediction Contests
      </Typography>
      
      {isAuthenticated && user ? (
        <Box sx={{ display: 'flex', alignItems: 'center', flexGrow: 1 }}>
          <Button color="inherit" component={Link} to="/contests">Contests</Button>
          <Button color="inherit" component={Link} to="/predictions">Predictions</Button>
          <Button color="inherit" component={Link} to="/sports">Sports</Button>
          <Box sx={{ flexGrow: 1 }} />
          <Typography variant="body2" sx={{ mr: 2 }}>
            Welcome, {user.name}
          </Typography>
          <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleMenu}
            color="inherit"
          >
            <Avatar sx={{ width: 32, height: 32 }}>
              {user.name.charAt(0).toUpperCase()}
            </Avatar>
          </IconButton>
          <Menu
            id="menu-appbar"
            anchorEl={anchorEl}
            anchorOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            keepMounted
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem onClick={handleClose}>
              <AccountCircle sx={{ mr: 1 }} />
              Profile
            </MenuItem>
            <MenuItem onClick={handleLogout}>
              <ExitToApp sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Box>
      ) : null}
    </Toolbar>
  )
}

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <ToastProvider>
        <AuthProvider>
          <Router>
            <AppBar position="static">
              <AppBarContent />
            </AppBar>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
              <Routes>
                <Route path="/login" element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />
                <Route 
                  path="/contests" 
                  element={
                    <ProtectedRoute>
                      <ContestsPage />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/sports" 
                  element={
                    <ProtectedRoute>
                      <SportsPage />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/predictions" 
                  element={
                    <ProtectedRoute>
                      <PredictionsPage />
                    </ProtectedRoute>
                  } 
                />
                <Route path="/" element={<Navigate to="/contests" replace />} />
              </Routes>
            </Container>
          </Router>
        </AuthProvider>
      </ToastProvider>
    </ThemeProvider>
  )
}

export default App

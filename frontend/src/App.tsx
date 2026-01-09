import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import { CssBaseline, Container, AppBar, Toolbar, Typography } from '@mui/material'
import { ToastProvider } from './contexts/ToastContext'
import ContestsPage from './pages/ContestsPage'

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

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <ToastProvider>
        <Router>
          <AppBar position="static">
            <Toolbar>
              <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                Sports Prediction Contests
              </Typography>
            </Toolbar>
          </AppBar>
          <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
            <Routes>
              <Route path="/contests" element={<ContestsPage />} />
              <Route path="/" element={<Navigate to="/contests" replace />} />
            </Routes>
          </Container>
        </Router>
      </ToastProvider>
    </ThemeProvider>
  )
}

export default App

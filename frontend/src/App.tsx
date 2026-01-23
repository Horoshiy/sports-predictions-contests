import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate, Link } from 'react-router-dom'
import { ConfigProvider, Layout, Menu, Dropdown, Avatar, Space, Typography } from 'antd'
import { UserOutlined, LogoutOutlined, TrophyOutlined, TeamOutlined, LineChartOutlined, FundOutlined, BarChartOutlined } from '@ant-design/icons'
import type { MenuProps } from 'antd'
import { antdTheme } from './theme/antd-theme'
import { ToastProvider } from './contexts/ToastContext'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { ProtectedRoute } from './components/auth/ProtectedRoute'
import ContestsPage from './pages/ContestsPage'
import SportsPage from './pages/SportsPage'
import PredictionsPage from './pages/PredictionsPage'
import AnalyticsPage from './pages/AnalyticsPage'
import TeamsPage from './pages/TeamsPage'
import ProfilePage from './pages/ProfilePage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'

const { Header, Content } = Layout
const { Text } = Typography

const AppHeader: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth()

  const menuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: <Link to="/profile">Profile</Link>,
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'Logout',
      onClick: logout,
    },
  ]

  return (
    <Header style={{ display: 'flex', alignItems: 'center', background: '#1976d2' }}>
      <div style={{ color: 'white', fontSize: '20px', fontWeight: 'bold', marginRight: '40px' }}>
        Sports Prediction Contests
      </div>
      
      {isAuthenticated && user ? (
        <>
          <Menu
            theme="dark"
            mode="horizontal"
            style={{ flex: 1, minWidth: 0, background: '#1976d2' }}
            items={[
              { key: 'contests', icon: <TrophyOutlined />, label: <Link to="/contests">Contests</Link> },
              { key: 'teams', icon: <TeamOutlined />, label: <Link to="/teams">Teams</Link> },
              { key: 'predictions', icon: <LineChartOutlined />, label: <Link to="/predictions">Predictions</Link> },
              { key: 'sports', icon: <FundOutlined />, label: <Link to="/sports">Sports</Link> },
              { key: 'analytics', icon: <BarChartOutlined />, label: <Link to="/analytics">Analytics</Link> },
            ]}
          />
          <Space style={{ marginLeft: 'auto' }}>
            <Text style={{ color: 'white' }}>Welcome, {user.name}</Text>
            <Dropdown menu={{ items: menuItems }} placement="bottomRight">
              <Avatar style={{ cursor: 'pointer', backgroundColor: '#1565c0' }}>
                {user.name.charAt(0).toUpperCase()}
              </Avatar>
            </Dropdown>
          </Space>
        </>
      ) : null}
    </Header>
  )
}

function App() {
  return (
    <ConfigProvider theme={antdTheme}>
      <ToastProvider>
        <AuthProvider>
          <Router>
            <Layout style={{ minHeight: '100vh' }}>
              <AppHeader />
              <Content style={{ padding: '24px 50px', maxWidth: '1600px', width: '100%', margin: '0 auto' }}>
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
                  <Route 
                    path="/analytics" 
                    element={
                      <ProtectedRoute>
                        <AnalyticsPage />
                      </ProtectedRoute>
                    } 
                  />
                  <Route 
                    path="/teams" 
                    element={
                      <ProtectedRoute>
                        <TeamsPage />
                      </ProtectedRoute>
                    } 
                  />
                  <Route 
                    path="/profile" 
                    element={
                      <ProtectedRoute>
                        <ProfilePage />
                      </ProtectedRoute>
                    } 
                  />
                  <Route path="/" element={<Navigate to="/contests" replace />} />
                </Routes>
              </Content>
            </Layout>
          </Router>
        </AuthProvider>
      </ToastProvider>
    </ConfigProvider>
  )
}

export default App

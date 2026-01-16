import React, { useState } from 'react'
import {
  Box,
  Typography,
  Paper,
  Grid,
  Card,
  CardContent,
  ToggleButton,
  ToggleButtonGroup,
  CircularProgress,
  Alert,
} from '@mui/material'
import {
  Analytics as AnalyticsIcon,
  TrendingUp as TrendingUpIcon,
  EmojiEvents as TrophyIcon,
} from '@mui/icons-material'
import { useAuth } from '../contexts/AuthContext'
import { useUserAnalytics } from '../hooks/use-analytics'
import { AccuracyChart } from '../components/analytics/AccuracyChart'
import { SportBreakdown } from '../components/analytics/SportBreakdown'
import { PlatformComparison } from '../components/analytics/PlatformComparison'
import { ExportButton } from '../components/analytics/ExportButton'
import type { TimeRange } from '../types/analytics.types'

const StatCard: React.FC<{
  title: string
  value: string | number
  subtitle?: string
  icon: React.ReactNode
  color?: string
}> = ({ title, value, subtitle, icon, color = 'primary.main' }) => (
  <Card>
    <CardContent>
      <Box display="flex" alignItems="center" gap={2}>
        <Box
          sx={{
            p: 1.5,
            borderRadius: 2,
            bgcolor: `${color}15`,
            color: color,
          }}
        >
          {icon}
        </Box>
        <Box>
          <Typography variant="h4" component="div">
            {value}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {title}
          </Typography>
          {subtitle && (
            <Typography variant="caption" color="text.secondary">
              {subtitle}
            </Typography>
          )}
        </Box>
      </Box>
    </CardContent>
  </Card>
)

export const AnalyticsPage: React.FC = () => {
  const { user } = useAuth()
  const [timeRange, setTimeRange] = useState<TimeRange>('30d')

  const {
    data: analytics,
    isLoading,
    error,
  } = useUserAnalytics(user?.id || 0, timeRange)

  const handleTimeRangeChange = (
    _: React.MouseEvent<HTMLElement>,
    newRange: TimeRange | null
  ) => {
    if (newRange) {
      setTimeRange(newRange)
    }
  }

  if (!user) {
    return (
      <Alert severity="warning">
        Please log in to view your analytics.
      </Alert>
    )
  }

  if (error) {
    return (
      <Alert severity="error">
        Failed to load analytics. Please try again later.
      </Alert>
    )
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4" component="h1" gutterBottom>
            Your Analytics
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Track your prediction performance and identify areas for improvement
          </Typography>
        </Box>
        <Box display="flex" gap={2} alignItems="center">
          <ToggleButtonGroup
            value={timeRange}
            exclusive
            onChange={handleTimeRangeChange}
            size="small"
          >
            <ToggleButton value="7d">7 Days</ToggleButton>
            <ToggleButton value="30d">30 Days</ToggleButton>
            <ToggleButton value="90d">90 Days</ToggleButton>
            <ToggleButton value="all">All Time</ToggleButton>
          </ToggleButtonGroup>
          <ExportButton userId={user.id} timeRange={timeRange} />
        </Box>
      </Box>

      {isLoading ? (
        <Box display="flex" justifyContent="center" py={8}>
          <CircularProgress />
        </Box>
      ) : analytics ? (
        <>
          <Grid container spacing={3} mb={3}>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Total Predictions"
                value={analytics.totalPredictions}
                icon={<AnalyticsIcon />}
                color="primary.main"
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Correct Predictions"
                value={analytics.correctPredictions}
                subtitle={`${analytics.overallAccuracy.toFixed(1)}% accuracy`}
                icon={<TrendingUpIcon />}
                color="success.main"
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Total Points"
                value={analytics.totalPoints.toFixed(1)}
                icon={<TrophyIcon />}
                color="warning.main"
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <StatCard
                title="Avg Points/Prediction"
                value={
                  analytics.totalPredictions > 0
                    ? (analytics.totalPoints / analytics.totalPredictions).toFixed(2)
                    : '0'
                }
                icon={<AnalyticsIcon />}
                color="info.main"
              />
            </Grid>
          </Grid>

          <Grid container spacing={3} mb={3}>
            <Grid item xs={12} lg={8}>
              <AccuracyChart trends={analytics.trends} />
            </Grid>
            <Grid item xs={12} lg={4}>
              <PlatformComparison
                userStats={{
                  accuracy: analytics.overallAccuracy,
                  avgPoints: analytics.totalPoints,
                  totalPredictions: analytics.totalPredictions,
                }}
                platformStats={analytics.platformComparison}
              />
            </Grid>
          </Grid>

          <Grid container spacing={3}>
            <Grid item xs={12}>
              <SportBreakdown bySport={analytics.bySport} />
            </Grid>
          </Grid>

          {analytics.byType && analytics.byType.length > 0 && (
            <Paper sx={{ mt: 3, p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Performance by Prediction Type
              </Typography>
              <Grid container spacing={2}>
                {analytics.byType.map((t) => (
                  <Grid item xs={12} sm={6} md={4} key={t.predictionType}>
                    <Card variant="outlined">
                      <CardContent>
                        <Typography variant="subtitle1" gutterBottom>
                          {t.predictionType.replace('_', ' ').toUpperCase()}
                        </Typography>
                        <Typography variant="h5" color="primary">
                          {t.accuracyPercentage.toFixed(1)}%
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {t.correctPredictions} / {t.totalPredictions} correct
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Avg: {t.averagePoints.toFixed(2)} pts
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            </Paper>
          )}
        </>
      ) : (
        <Paper sx={{ p: 4, textAlign: 'center' }}>
          <AnalyticsIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No Analytics Data Yet
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Start making predictions to see your performance analytics here.
          </Typography>
        </Paper>
      )}
    </Box>
  )
}

export default AnalyticsPage

import React from 'react'
import {
  Card,
  CardContent,
  Typography,
  Box,
  Grid,
  Chip,
  LinearProgress,
} from '@mui/material'
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  Remove as NeutralIcon,
} from '@mui/icons-material'
import type { PlatformStats } from '../../types/analytics.types'

interface PlatformComparisonProps {
  userStats: {
    accuracy: number
    avgPoints: number
    totalPredictions: number
  }
  platformStats: PlatformStats | null
}

const ComparisonCard: React.FC<{
  label: string
  userValue: number
  platformValue: number
  format?: 'percent' | 'number'
  higherIsBetter?: boolean
}> = ({ label, userValue, platformValue, format = 'number', higherIsBetter = true }) => {
  const diff = userValue - platformValue
  const isAbove = diff > 0
  const isBetter = higherIsBetter ? isAbove : !isAbove

  const formatValue = (v: number) => {
    if (format === 'percent') return `${v.toFixed(1)}%`
    return v.toFixed(1)
  }

  return (
    <Card variant="outlined">
      <CardContent>
        <Typography variant="subtitle2" color="text.secondary" gutterBottom>
          {label}
        </Typography>
        <Box display="flex" alignItems="baseline" gap={1}>
          <Typography variant="h4" component="span">
            {formatValue(userValue)}
          </Typography>
          <Chip
            size="small"
            icon={
              Math.abs(diff) < 0.1 ? (
                <NeutralIcon />
              ) : isBetter ? (
                <TrendingUpIcon />
              ) : (
                <TrendingDownIcon />
              )
            }
            label={`${diff > 0 ? '+' : ''}${formatValue(diff)}`}
            color={Math.abs(diff) < 0.1 ? 'default' : isBetter ? 'success' : 'error'}
            variant="outlined"
          />
        </Box>
        <Box mt={2}>
          <Typography variant="caption" color="text.secondary">
            Platform Average: {formatValue(platformValue)}
          </Typography>
          <LinearProgress
            variant="determinate"
            value={Math.min((userValue / Math.max(platformValue, 1)) * 50, 100)}
            sx={{ mt: 1, height: 8, borderRadius: 4 }}
            color={isBetter ? 'success' : 'error'}
          />
        </Box>
      </CardContent>
    </Card>
  )
}

export const PlatformComparison: React.FC<PlatformComparisonProps> = ({
  userStats,
  platformStats,
}) => {
  if (!platformStats) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>Platform Comparison</Typography>
          <Box textAlign="center" py={4}>
            <Typography color="text.secondary">
              Platform statistics not available
            </Typography>
          </Box>
        </CardContent>
      </Card>
    )
  }

  const userAvgPoints = userStats.totalPredictions > 0
    ? userStats.avgPoints / userStats.totalPredictions
    : 0

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          How You Compare
        </Typography>
        <Typography variant="body2" color="text.secondary" gutterBottom>
          Your performance vs {platformStats.totalUsers.toLocaleString()} users
        </Typography>
        
        <Grid container spacing={2} mt={1}>
          <Grid item xs={12}>
            <ComparisonCard
              label="Accuracy"
              userValue={userStats.accuracy}
              platformValue={platformStats.averageAccuracy}
              format="percent"
              higherIsBetter={true}
            />
          </Grid>
          <Grid item xs={12}>
            <ComparisonCard
              label="Avg Points per Prediction"
              userValue={userAvgPoints}
              platformValue={platformStats.averagePointsPerPrediction}
              format="number"
              higherIsBetter={true}
            />
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  )
}

export default PlatformComparison

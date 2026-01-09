import React from 'react'
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  LinearProgress,
  Avatar,
  Tooltip,
  IconButton,
  Divider,
} from '@mui/material'
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  Remove as StableIcon,
  EmojiEvents as TrophyIcon,
  Timeline as TimelineIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material'
import { useQuery } from '@tanstack/react-query'
import scoringService from '../../services/scoring-service'
import type { Score } from '../../types/scoring.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface UserScoreProps {
  userId: number
  contestId: number
  userName?: string
  showDetails?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
  onRefresh?: () => void
}

interface ScoreBreakdown {
  totalPredictions: number
  scoredPredictions: number
  averagePoints: number
  bestScore: number
  recentScores: Score[]
}

const getRankColor = (rank: number) => {
  if (rank === 1) return 'gold'
  if (rank === 2) return 'silver'
  if (rank === 3) return '#CD7F32' // Bronze
  if (rank <= 10) return 'primary'
  return 'text.secondary'
}

const getTrendIcon = (change: number) => {
  if (change > 0) return <TrendingUpIcon color="success" fontSize="small" />
  if (change < 0) return <TrendingDownIcon color="error" fontSize="small" />
  return <StableIcon color="disabled" fontSize="small" />
}

export const UserScore: React.FC<UserScoreProps> = ({
  userId,
  contestId,
  userName,
  showDetails = true,
  autoRefresh = true,
  refreshInterval = 30000,
  onRefresh,
}) => {
  // Query for user scores
  const {
    data: userScores,
    isLoading: isLoadingScores,
    error: scoresError,
    refetch: refetchScores,
  } = useQuery({
    queryKey: ['userScores', userId, contestId],
    queryFn: () => scoringService.getUserScores({ userId, contestId }),
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  // Query for user rank
  const {
    data: userRank,
    isLoading: isLoadingRank,
    error: rankError,
    refetch: refetchRank,
  } = useQuery({
    queryKey: ['userRank', userId, contestId],
    queryFn: () => scoringService.getUserRank({ contestId, userId }),
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  // Calculate score breakdown
  const scoreBreakdown: ScoreBreakdown | null = React.useMemo(() => {
    if (!userScores?.scores) return null

    const scores = userScores.scores
    const totalPredictions = scores.length
    const scoredPredictions = scores.filter(s => s.points > 0).length
    const averagePoints = totalPredictions > 0 
      ? userScores.totalPoints / totalPredictions 
      : 0
    const bestScore = Math.max(...scores.map(s => s.points), 0)
    const recentScores = scores
      .sort((a, b) => new Date(b.scoredAt).getTime() - new Date(a.scoredAt).getTime())
      .slice(0, 5)

    return {
      totalPredictions,
      scoredPredictions,
      averagePoints,
      bestScore,
      recentScores,
    }
  }, [userScores])

  const handleRefresh = () => {
    refetchScores()
    refetchRank()
    onRefresh?.()
  }

  const isLoading = isLoadingScores || isLoadingRank
  const hasError = scoresError || rankError

  if (hasError) {
    return (
      <Card variant="outlined" sx={{ bgcolor: 'error.50' }}>
        <CardContent>
          <Typography color="error" variant="body2">
            Failed to load user score data
          </Typography>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card variant="outlined">
      <CardContent>
        {/* Header */}
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Box display="flex" alignItems="center" gap={2}>
            <Avatar sx={{ width: 40, height: 40 }}>
              {userName ? userName.charAt(0).toUpperCase() : 'U'}
            </Avatar>
            <Box>
              <Typography variant="h6">
                {userName || `User ${userId}`}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Contest Performance
              </Typography>
            </Box>
          </Box>
          <Tooltip title="Refresh scores">
            <IconButton onClick={handleRefresh} size="small">
              <RefreshIcon />
            </IconButton>
          </Tooltip>
        </Box>

        {/* Loading State */}
        {isLoading && (
          <Box>
            <LinearProgress sx={{ mb: 2 }} />
            <Typography variant="body2" color="text.secondary" textAlign="center">
              Loading score data...
            </Typography>
          </Box>
        )}

        {/* Main Score Display */}
        {!isLoading && userRank && (
          <Box>
            {/* Rank and Points */}
            <Box display="flex" justifyContent="space-around" mb={2}>
              <Box textAlign="center">
                <Box display="flex" alignItems="center" justifyContent="center" gap={1}>
                  <TrophyIcon sx={{ color: getRankColor(userRank.rank) }} />
                  <Typography 
                    variant="h4" 
                    sx={{ color: getRankColor(userRank.rank), fontWeight: 'bold' }}
                  >
                    #{userRank.rank}
                  </Typography>
                </Box>
                <Typography variant="caption" color="text.secondary">
                  Current Rank
                </Typography>
              </Box>
              
              <Divider orientation="vertical" flexItem />
              
              <Box textAlign="center">
                <Typography variant="h4" color="primary" fontWeight="bold">
                  {userRank.totalPoints.toFixed(1)}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  Total Points
                </Typography>
              </Box>
            </Box>

            {/* Score Breakdown */}
            {showDetails && scoreBreakdown && (
              <>
                <Divider sx={{ my: 2 }} />
                
                <Box display="flex" justifyContent="space-between" mb={2}>
                  <Box textAlign="center" flex={1}>
                    <Typography variant="h6" color="text.primary">
                      {scoreBreakdown.totalPredictions}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Predictions
                    </Typography>
                  </Box>
                  
                  <Box textAlign="center" flex={1}>
                    <Typography variant="h6" color="success.main">
                      {scoreBreakdown.scoredPredictions}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Scored
                    </Typography>
                  </Box>
                  
                  <Box textAlign="center" flex={1}>
                    <Typography variant="h6" color="info.main">
                      {scoreBreakdown.averagePoints.toFixed(1)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Avg Points
                    </Typography>
                  </Box>
                  
                  <Box textAlign="center" flex={1}>
                    <Typography variant="h6" color="warning.main">
                      {scoreBreakdown.bestScore.toFixed(1)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Best Score
                    </Typography>
                  </Box>
                </Box>

                {/* Success Rate Progress */}
                <Box mb={2}>
                  <Box display="flex" justifyContent="space-between" mb={1}>
                    <Typography variant="body2">Success Rate</Typography>
                    <Typography variant="body2" color="text.secondary">
                      {scoreBreakdown.totalPredictions > 0 
                        ? Math.round((scoreBreakdown.scoredPredictions / scoreBreakdown.totalPredictions) * 100)
                        : 0}%
                    </Typography>
                  </Box>
                  <LinearProgress 
                    variant="determinate" 
                    value={scoreBreakdown.totalPredictions > 0 
                      ? (scoreBreakdown.scoredPredictions / scoreBreakdown.totalPredictions) * 100
                      : 0}
                    sx={{ height: 8, borderRadius: 4 }}
                  />
                </Box>

                {/* Recent Scores */}
                {scoreBreakdown.recentScores.length > 0 && (
                  <Box>
                    <Typography variant="subtitle2" gutterBottom>
                      Recent Scores
                    </Typography>
                    <Box display="flex" gap={1} flexWrap="wrap">
                      {scoreBreakdown.recentScores.map((score) => (
                        <Tooltip 
                          key={score.id}
                          title={`Prediction ${score.predictionId} - ${formatRelativeTime(score.scoredAt)}`}
                        >
                          <Chip
                            label={score.points.toFixed(1)}
                            size="small"
                            color={score.points > 0 ? 'success' : 'default'}
                            variant="outlined"
                          />
                        </Tooltip>
                      ))}
                    </Box>
                  </Box>
                )}
              </>
            )}
          </Box>
        )}

        {/* Empty State */}
        {!isLoading && !userRank && (
          <Box textAlign="center" py={2}>
            <TimelineIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 1 }} />
            <Typography variant="body2" color="text.secondary">
              No scores yet. Make some predictions to see your performance!
            </Typography>
          </Box>
        )}
      </CardContent>
    </Card>
  )
}

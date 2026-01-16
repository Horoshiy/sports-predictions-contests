import React, { useMemo, useState, useEffect } from 'react'
import {
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
} from 'material-react-table'
import {
  Box,
  Typography,
  Chip,
  Avatar,
  IconButton,
  Tooltip,
  Card,
  CardContent,
  CircularProgress,
  Alert,
  Button,
} from '@mui/material'
import {
  Refresh as RefreshIcon,
  EmojiEvents as TrophyIcon,
} from '@mui/icons-material'
import { useQuery } from '@tanstack/react-query'
import scoringService from '../../services/scoring-service'
import type { LeaderboardEntry } from '../../types/scoring.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface LeaderboardTableProps {
  contestId: number
  currentUserId?: number
  limit?: number
  autoRefresh?: boolean
  refreshInterval?: number
}

const getRankColor = (rank: number) => {
  switch (rank) {
    case 1:
      return '#FFD700' // Gold
    case 2:
      return '#C0C0C0' // Silver
    case 3:
      return '#CD7F32' // Bronze
    default:
      return 'inherit'
  }
}

const getRankIcon = (rank: number) => {
  if (rank <= 3) {
    return <TrophyIcon sx={{ color: getRankColor(rank), fontSize: 20 }} />
  }
  return null
}

export const LeaderboardTable: React.FC<LeaderboardTableProps> = ({
  contestId,
  currentUserId,
  limit = 50,
  autoRefresh = true,
  refreshInterval = 30000, // 30 seconds
}) => {
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date())

  // Query for leaderboard data
  const {
    data: leaderboard,
    isLoading,
    error,
    refetch,
    isRefetching,
  } = useQuery({
    queryKey: ['leaderboard', contestId, limit],
    queryFn: () => scoringService.getLeaderboard({ contestId, limit }),
    refetchInterval: autoRefresh ? refreshInterval : false,
    refetchIntervalInBackground: true,
  })

  // Update lastUpdated when data changes
  useEffect(() => {
    if (leaderboard) {
      setLastUpdated(new Date())
    }
  }, [leaderboard])

  // Query for current user's rank if provided
  const {
    data: userRank,
    isLoading: isLoadingUserRank,
  } = useQuery({
    queryKey: ['userRank', contestId, currentUserId],
    queryFn: () => currentUserId 
      ? scoringService.getUserRank({ contestId, userId: currentUserId })
      : null,
    enabled: !!currentUserId,
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  // Query for current user's streak if provided
  const {
    data: userStreak,
    error: userStreakError,
  } = useQuery({
    queryKey: ['userStreak', contestId, currentUserId],
    queryFn: () => currentUserId 
      ? scoringService.getUserStreak({ contestId, userId: currentUserId })
      : null,
    enabled: !!currentUserId,
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  // Log streak error if any (silent fail for UI)
  if (userStreakError) {
    console.error('Failed to fetch user streak:', userStreakError)
  }

  // Define table columns
  const columns = useMemo<MRT_ColumnDef<LeaderboardEntry>[]>(
    () => [
      {
        accessorKey: 'rank',
        header: 'Rank',
        size: 80,
        Cell: ({ cell, row }) => {
          const rank = cell.getValue<number>()
          const isCurrentUser = currentUserId === row.original.userId
          
          return (
            <Box 
              display="flex" 
              alignItems="center" 
              gap={1}
              sx={{
                fontWeight: isCurrentUser ? 'bold' : 'normal',
                color: isCurrentUser ? 'primary.main' : 'inherit',
              }}
            >
              {getRankIcon(rank)}
              <Typography 
                variant="body2" 
                sx={{ 
                  fontWeight: rank <= 3 ? 'bold' : 'normal',
                  color: getRankColor(rank),
                }}
              >
                #{rank}
              </Typography>
            </Box>
          )
        },
      },
      {
        accessorKey: 'userName',
        header: 'Player',
        size: 200,
        Cell: ({ cell, row }) => {
          const userName = cell.getValue<string>() || `User ${row.original.userId}`
          const isCurrentUser = currentUserId === row.original.userId
          
          return (
            <Box display="flex" alignItems="center" gap={2}>
              <Avatar sx={{ width: 32, height: 32, fontSize: 14 }}>
                {userName.charAt(0).toUpperCase()}
              </Avatar>
              <Box>
                <Typography 
                  variant="body2" 
                  sx={{ 
                    fontWeight: isCurrentUser ? 'bold' : 'normal',
                    color: isCurrentUser ? 'primary.main' : 'inherit',
                  }}
                >
                  {userName}
                  {isCurrentUser && (
                    <Chip 
                      label="You" 
                      size="small" 
                      color="primary" 
                      sx={{ ml: 1, height: 20 }}
                    />
                  )}
                </Typography>
              </Box>
            </Box>
          )
        },
      },
      {
        accessorKey: 'totalPoints',
        header: 'Points',
        size: 120,
        Cell: ({ cell }) => {
          const points = cell.getValue<number>()
          return (
            <Typography 
              variant="body2" 
              sx={{ 
                fontWeight: 'bold',
                color: points > 0 ? 'success.main' : 'text.secondary',
              }}
            >
              {points.toFixed(1)}
            </Typography>
          )
        },
      },
      {
        accessorKey: 'currentStreak',
        header: 'Streak',
        size: 100,
        Cell: ({ cell, row }) => {
          const streak = cell.getValue<number>() || 0
          const multiplier = row.original.multiplier || 1
          return (
            <Box display="flex" alignItems="center" gap={1}>
              <Typography variant="body2">
                🔥 {streak}
              </Typography>
              {multiplier > 1 && (
                <Chip 
                  label={`${multiplier}x`} 
                  size="small" 
                  color="warning"
                />
              )}
            </Box>
          )
        },
      },
      {
        accessorKey: 'updatedAt',
        header: 'Last Updated',
        size: 150,
        Cell: ({ cell }) => {
          const updatedAt = cell.getValue<string>()
          return (
            <Typography variant="body2" color="text.secondary">
              {formatRelativeTime(updatedAt)}
            </Typography>
          )
        },
      },
    ],
    [currentUserId]
  )

  const table = useMaterialReactTable({
    columns,
    data: leaderboard?.entries || [],
    enableColumnActions: false,
    enableColumnFilters: false,
    enablePagination: false,
    enableSorting: false,
    enableBottomToolbar: false,
    enableTopToolbar: false,
    muiTableBodyRowProps: ({ row }) => ({
      sx: {
        backgroundColor: currentUserId === row.original.userId 
          ? 'action.selected' 
          : 'inherit',
      },
    }),
    muiTableProps: {
      sx: {
        '& .MuiTableHead-root': {
          '& .MuiTableCell-root': {
            backgroundColor: 'grey.50',
            fontWeight: 'bold',
          },
        },
      },
    },
  })

  const handleRefresh = () => {
    refetch()
    setLastUpdated(new Date())
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        Failed to load leaderboard. Please try again.
        <Button onClick={handleRefresh} sx={{ ml: 2 }}>
          Retry
        </Button>
      </Alert>
    )
  }

  return (
    <Card>
      <CardContent>
        {/* Header */}
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h6" component="h2">
            Leaderboard
          </Typography>
          <Box display="flex" alignItems="center" gap={1}>
            <Typography variant="caption" color="text.secondary">
              Last updated: {formatRelativeTime(lastUpdated.toISOString())}
            </Typography>
            <Tooltip title="Refresh leaderboard">
              <IconButton 
                onClick={handleRefresh} 
                disabled={isRefetching}
                size="small"
              >
                {isRefetching ? (
                  <CircularProgress size={20} />
                ) : (
                  <RefreshIcon />
                )}
              </IconButton>
            </Tooltip>
          </Box>
        </Box>

        {/* Current User Rank Card */}
        {currentUserId && userRank && !isLoadingUserRank && (
          <Card variant="outlined" sx={{ mb: 2, bgcolor: 'primary.50' }}>
            <CardContent sx={{ py: 1.5 }}>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Typography variant="subtitle2" color="primary">
                  Your Position
                </Typography>
                <Box display="flex" alignItems="center" gap={2}>
                  <Box textAlign="center">
                    <Typography variant="h6" color="primary">
                      #{userRank.rank}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Rank
                    </Typography>
                  </Box>
                  <Box textAlign="center">
                    <Typography variant="h6" color="primary">
                      {userRank.totalPoints.toFixed(1)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Points
                    </Typography>
                  </Box>
                  {userStreak && (
                    <Box textAlign="center">
                      <Typography variant="h6" color="warning.main">
                        🔥 {userStreak.currentStreak}
                        {userStreak.multiplier > 1 && (
                          <Chip 
                            label={`${userStreak.multiplier}x`} 
                            size="small" 
                            color="warning"
                            sx={{ ml: 0.5, height: 20 }}
                          />
                        )}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        Streak
                      </Typography>
                    </Box>
                  )}
                </Box>
              </Box>
            </CardContent>
          </Card>
        )}

        {/* Loading State */}
        {isLoading && (
          <Box display="flex" justifyContent="center" py={4}>
            <CircularProgress />
          </Box>
        )}

        {/* Empty State */}
        {!isLoading && (!leaderboard?.entries || leaderboard.entries.length === 0) && (
          <Box textAlign="center" py={4}>
            <TrophyIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" gutterBottom>
              No rankings yet
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Scores will appear here once predictions are evaluated.
            </Typography>
          </Box>
        )}

        {/* Leaderboard Table */}
        {!isLoading && leaderboard?.entries && leaderboard.entries.length > 0 && (
          <MaterialReactTable table={table} />
        )}

        {/* Auto-refresh indicator */}
        {autoRefresh && (
          <Box display="flex" justifyContent="center" mt={2}>
            <Chip 
              label={`Auto-refreshing every ${refreshInterval / 1000}s`}
              size="small"
              variant="outlined"
              color="primary"
            />
          </Box>
        )}
      </CardContent>
    </Card>
  )
}

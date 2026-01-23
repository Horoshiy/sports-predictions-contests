import React, { useMemo, useState } from 'react'
import {
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
} from 'material-react-table'
import {
  Box,
  Button,
  IconButton,
  Tooltip,
  Chip,
  Typography,
  Alert,
  CircularProgress,
} from '@mui/material'
import {
  Check as AcceptIcon,
  Close as DeclineIcon,
  Delete as WithdrawIcon,
  Visibility as ViewIcon,
} from '@mui/icons-material'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import challengeService from '../../services/challenge-service'
import type { Challenge } from '../../types/challenge.types'
import { CHALLENGE_STATUSES } from '../../types/challenge.types'
import { formatDate, formatRelativeTime } from '../../utils/date-utils'

interface ChallengeListProps {
  userId: number
  statusFilter?: Challenge['status']
  onViewChallenge?: (challenge: Challenge) => void
  onCreateChallenge?: () => void
}

export const ChallengeList: React.FC<ChallengeListProps> = ({
  userId,
  statusFilter,
  onViewChallenge,
  onCreateChallenge,
}) => {
  const queryClient = useQueryClient()
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  })

  // Fetch challenges
  const {
    data: challengesData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['challenges', userId, statusFilter, pagination],
    queryFn: () =>
      challengeService.listUserChallenges({
        userId,
        status: statusFilter,
        pagination: {
          page: pagination.pageIndex + 1,
          limit: pagination.pageSize,
        },
      }),
  })

  // Accept challenge mutation
  const acceptChallengeMutation = useMutation({
    mutationFn: (challengeId: number) =>
      challengeService.acceptChallenge({ id: challengeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['challenges'] })
    },
  })

  // Decline challenge mutation
  const declineChallengeMutation = useMutation({
    mutationFn: (challengeId: number) =>
      challengeService.declineChallenge({ id: challengeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['challenges'] })
    },
  })

  // Withdraw challenge mutation
  const withdrawChallengeMutation = useMutation({
    mutationFn: (challengeId: number) =>
      challengeService.withdrawChallenge({ id: challengeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['challenges'] })
    },
  })

  const handleAcceptChallenge = (challenge: Challenge) => {
    acceptChallengeMutation.mutate(challenge.id)
  }

  const handleDeclineChallenge = (challenge: Challenge) => {
    declineChallengeMutation.mutate(challenge.id)
  }

  const handleWithdrawChallenge = (challenge: Challenge) => {
    withdrawChallengeMutation.mutate(challenge.id)
  }

  const columns = useMemo<MRT_ColumnDef<Challenge>[]>(
    () => [
      {
        accessorKey: 'id',
        header: 'ID',
        size: 80,
      },
      {
        accessorKey: 'challengerId',
        header: 'Challenger',
        Cell: ({ row }) => {
          const challenge = row.original
          const isCurrentUserChallenger = challenge.challengerId === userId
          return (
            <Typography variant="body2">
              {isCurrentUserChallenger ? 'You' : `User ${challenge.challengerId}`}
            </Typography>
          )
        },
      },
      {
        accessorKey: 'opponentId',
        header: 'Opponent',
        Cell: ({ row }) => {
          const challenge = row.original
          const isCurrentUserOpponent = challenge.opponentId === userId
          return (
            <Typography variant="body2">
              {isCurrentUserOpponent ? 'You' : `User ${challenge.opponentId}`}
            </Typography>
          )
        },
      },
      {
        accessorKey: 'eventId',
        header: 'Event',
        Cell: ({ row }) => (
          <Typography variant="body2">Event {row.original.eventId}</Typography>
        ),
      },
      {
        accessorKey: 'message',
        header: 'Message',
        Cell: ({ row }) => (
          <Typography
            variant="body2"
            sx={{
              maxWidth: 200,
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            }}
          >
            {row.original.message || 'No message'}
          </Typography>
        ),
      },
      {
        accessorKey: 'status',
        header: 'Status',
        Cell: ({ row }) => {
          const status = row.original.status
          const statusInfo = CHALLENGE_STATUSES[status]
          return (
            <Chip
              label={statusInfo.label}
              color={statusInfo.color}
              size="small"
              title={statusInfo.description}
            />
          )
        },
      },
      {
        accessorKey: 'expiresAt',
        header: 'Expires',
        Cell: ({ row }) => {
          const challenge = row.original
          if (challenge.status !== 'pending') return null
          
          const expiresAt = new Date(challenge.expiresAt)
          const now = new Date()
          
          if (now > expiresAt) {
            return (
              <Typography variant="caption" color="error">
                Expired
              </Typography>
            )
          }
          
          return (
            <Typography variant="caption" color="warning.main">
              {formatRelativeTime(expiresAt)}
            </Typography>
          )
        },
      },
      {
        accessorKey: 'createdAt',
        header: 'Created',
        Cell: ({ row }) => (
          <Typography variant="caption">
            {formatDate(row.original.createdAt)}
          </Typography>
        ),
      },
    ],
    [userId]
  )

  const table = useMaterialReactTable({
    columns,
    data: challengesData?.challenges || [],
    enableRowActions: true,
    positionActionsColumn: 'last',
    renderRowActions: ({ row }) => {
      const challenge = row.original
      const isCurrentUserChallenger = challenge.challengerId === userId
      const isCurrentUserOpponent = challenge.opponentId === userId
      const statusInfo = challengeService.getChallengeStatusInfo(challenge)

      return (
        <Box sx={{ display: 'flex', gap: 1 }}>
          {onViewChallenge && (
            <Tooltip title="View Details">
              <IconButton
                size="small"
                onClick={() => onViewChallenge(challenge)}
              >
                <ViewIcon />
              </IconButton>
            </Tooltip>
          )}

          {/* Accept button for opponents on pending challenges */}
          {isCurrentUserOpponent && statusInfo.canAccept && (
            <Tooltip title="Accept Challenge">
              <IconButton
                size="small"
                color="success"
                onClick={() => handleAcceptChallenge(challenge)}
                disabled={acceptChallengeMutation.isPending}
              >
                <AcceptIcon />
              </IconButton>
            </Tooltip>
          )}

          {/* Decline button for opponents on pending challenges */}
          {isCurrentUserOpponent && statusInfo.canDecline && (
            <Tooltip title="Decline Challenge">
              <IconButton
                size="small"
                color="error"
                onClick={() => handleDeclineChallenge(challenge)}
                disabled={declineChallengeMutation.isPending}
              >
                <DeclineIcon />
              </IconButton>
            </Tooltip>
          )}

          {/* Withdraw button for challengers on pending challenges */}
          {isCurrentUserChallenger && statusInfo.canWithdraw && (
            <Tooltip title="Withdraw Challenge">
              <IconButton
                size="small"
                color="warning"
                onClick={() => handleWithdrawChallenge(challenge)}
                disabled={withdrawChallengeMutation.isPending}
              >
                <WithdrawIcon />
              </IconButton>
            </Tooltip>
          )}
        </Box>
      )
    },
    manualPagination: true,
    rowCount: challengesData?.pagination.total || 0,
    onPaginationChange: setPagination,
    state: {
      isLoading,
      pagination,
    },
    renderTopToolbarCustomActions: () => (
      <Box sx={{ display: 'flex', gap: 1 }}>
        {onCreateChallenge && (
          <Button
            variant="contained"
            startIcon={<AcceptIcon />}
            onClick={onCreateChallenge}
          >
            Create Challenge
          </Button>
        )}
      </Box>
    ),
  })

  if (error) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        Failed to load challenges: {error.message}
      </Alert>
    )
  }

  return (
    <Box>
      <MaterialReactTable table={table} />
      
      {/* Loading overlay for mutations */}
      {(acceptChallengeMutation.isPending ||
        declineChallengeMutation.isPending ||
        withdrawChallengeMutation.isPending) && (
        <Box
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            bgcolor: 'rgba(255, 255, 255, 0.7)',
            zIndex: 1000,
          }}
        >
          <CircularProgress />
        </Box>
      )}
    </Box>
  )
}

export default ChallengeList

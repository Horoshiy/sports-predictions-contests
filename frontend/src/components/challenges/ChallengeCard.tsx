import React from 'react'
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Chip,
  Box,
  IconButton,
  Tooltip,
  Avatar,
} from '@mui/material'
import {
  Check as AcceptIcon,
  Close as DeclineIcon,
  Delete as WithdrawIcon,
  Visibility as ViewIcon,
  Person as PersonIcon,
  Event as EventIcon,
  Timer as TimerIcon,
} from '@mui/icons-material'
import type { Challenge, ChallengeWithUserInfo } from '../../types/challenge.types'
import { CHALLENGE_STATUSES } from '../../types/challenge.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface ChallengeCardProps {
  challenge: ChallengeWithUserInfo
  currentUserId: number
  onAccept?: (challenge: Challenge) => void
  onDecline?: (challenge: Challenge) => void
  onWithdraw?: (challenge: Challenge) => void
  onView?: (challenge: Challenge) => void
  loading?: boolean
}

const getStatusColor = (status: Challenge['status']) => {
  return CHALLENGE_STATUSES[status]?.color || 'default'
}

export const ChallengeCard: React.FC<ChallengeCardProps> = ({
  challenge,
  currentUserId,
  onAccept,
  onDecline,
  onWithdraw,
  onView,
  loading = false,
}) => {
  const isCurrentUserChallenger = challenge.challengerId === currentUserId
  const isCurrentUserOpponent = challenge.opponentId === currentUserId
  const statusInfo = CHALLENGE_STATUSES[challenge.status]

  // Check if challenge is expired
  const isExpired = challenge.status === 'pending' && new Date() > new Date(challenge.expiresAt)
  const canAccept = challenge.status === 'pending' && !isExpired && isCurrentUserOpponent
  const canDecline = challenge.status === 'pending' && isCurrentUserOpponent
  const canWithdraw = challenge.status === 'pending' && isCurrentUserChallenger

  const getOpponentName = () => {
    if (isCurrentUserChallenger) {
      return challenge.opponentName || `User ${challenge.opponentId}`
    } else {
      return challenge.challengerName || `User ${challenge.challengerId}`
    }
  }

  const getOpponentRole = () => {
    if (isCurrentUserChallenger) {
      return 'Opponent'
    } else {
      return 'Challenger'
    }
  }

  const getScoreDisplay = () => {
    if (challenge.status !== 'completed') return null

    const userScore = isCurrentUserChallenger ? challenge.challengerScore : challenge.opponentScore
    const opponentScore = isCurrentUserChallenger ? challenge.opponentScore : challenge.challengerScore
    
    let result = 'Tie'
    let resultColor = 'info'
    
    if (challenge.winnerId === currentUserId) {
      result = 'Won'
      resultColor = 'success'
    } else if (challenge.winnerId && challenge.winnerId !== currentUserId) {
      result = 'Lost'
      resultColor = 'error'
    }

    return (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 1 }}>
        <Typography variant="body2" fontWeight="bold">
          Score: {userScore} - {opponentScore}
        </Typography>
        <Chip
          label={result}
          color={resultColor as any}
          size="small"
        />
      </Box>
    )
  }

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <CardContent sx={{ flexGrow: 1 }}>
        {/* Header with status */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Typography variant="h6" component="h3">
            Challenge #{challenge.id}
          </Typography>
          <Chip
            label={isExpired ? 'Expired' : statusInfo.label}
            color={isExpired ? 'default' : getStatusColor(challenge.status)}
            size="small"
          />
        </Box>

        {/* Opponent info */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
          <Avatar sx={{ width: 32, height: 32 }}>
            <PersonIcon />
          </Avatar>
          <Box>
            <Typography variant="body2" fontWeight="medium">
              {getOpponentName()}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {getOpponentRole()}
            </Typography>
          </Box>
        </Box>

        {/* Event info */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
          <EventIcon color="action" />
          <Typography variant="body2">
            {challenge.eventTitle || `Event #${challenge.eventId}`}
          </Typography>
        </Box>

        {/* Message */}
        {challenge.message && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="body2" color="text.secondary" fontStyle="italic">
              "{challenge.message}"
            </Typography>
          </Box>
        )}

        {/* Timing info */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
          <TimerIcon color="action" />
          <Typography variant="caption" color="text.secondary">
            Created {formatRelativeTime(new Date(challenge.createdAt))}
          </Typography>
        </Box>

        {/* Expiration for pending challenges */}
        {challenge.status === 'pending' && (
          <Box sx={{ mb: 1 }}>
            <Typography variant="caption" color={isExpired ? 'error' : 'warning.main'}>
              {isExpired 
                ? 'Expired' 
                : `Expires ${formatRelativeTime(new Date(challenge.expiresAt))}`
              }
            </Typography>
          </Box>
        )}

        {/* Acceptance date */}
        {challenge.acceptedAt && (
          <Box sx={{ mb: 1 }}>
            <Typography variant="caption" color="text.secondary">
              Accepted {formatRelativeTime(new Date(challenge.acceptedAt))}
            </Typography>
          </Box>
        )}

        {/* Score display for completed challenges */}
        {getScoreDisplay()}
      </CardContent>

      <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
        <Box sx={{ display: 'flex', gap: 1 }}>
          {/* Accept button */}
          {canAccept && onAccept && (
            <Button
              size="small"
              variant="contained"
              color="success"
              startIcon={<AcceptIcon />}
              onClick={() => onAccept(challenge)}
              disabled={loading}
            >
              Accept
            </Button>
          )}

          {/* Decline button */}
          {canDecline && onDecline && (
            <Button
              size="small"
              variant="outlined"
              color="error"
              startIcon={<DeclineIcon />}
              onClick={() => onDecline(challenge)}
              disabled={loading}
            >
              Decline
            </Button>
          )}

          {/* Withdraw button */}
          {canWithdraw && onWithdraw && (
            <Button
              size="small"
              variant="outlined"
              color="warning"
              startIcon={<WithdrawIcon />}
              onClick={() => onWithdraw(challenge)}
              disabled={loading}
            >
              Withdraw
            </Button>
          )}
        </Box>

        {/* View details button */}
        {onView && (
          <Tooltip title="View Details">
            <IconButton
              size="small"
              onClick={() => onView(challenge)}
              disabled={loading}
            >
              <ViewIcon />
            </IconButton>
          </Tooltip>
        )}
      </CardActions>
    </Card>
  )
}

export default ChallengeCard

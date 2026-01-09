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
} from '@mui/material'
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  People as PeopleIcon,
  Sports as SportsIcon,
} from '@mui/icons-material'
import type { Contest } from '../../types/contest.types'
import { formatDate, formatRelativeTime, getContestStatusByDate } from '../../utils/date-utils'

interface ContestCardProps {
  contest: Contest
  onEdit?: (contest: Contest) => void
  onDelete?: (contest: Contest) => void
  onViewParticipants?: (contest: Contest) => void
  onJoin?: (contest: Contest) => void
  onLeave?: (contest: Contest) => void
  isParticipant?: boolean
  canEdit?: boolean
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'draft':
      return 'default'
    case 'active':
      return 'success'
    case 'completed':
      return 'info'
    case 'cancelled':
      return 'error'
    default:
      return 'default'
  }
}

const getStatusByDateColor = (status: 'upcoming' | 'active' | 'completed') => {
  switch (status) {
    case 'upcoming':
      return 'warning'
    case 'active':
      return 'success'
    case 'completed':
      return 'info'
    default:
      return 'default'
  }
}

export const ContestCard: React.FC<ContestCardProps> = ({
  contest,
  onEdit,
  onDelete,
  onViewParticipants,
  onJoin,
  onLeave,
  isParticipant = false,
  canEdit = false,
}) => {
  const dateStatus = getContestStatusByDate(contest.startDate, contest.endDate)
  const isActive = contest.status === 'active' && dateStatus === 'active'
  const canJoin = contest.status === 'active' && !isParticipant && 
    (contest.maxParticipants === 0 || contest.currentParticipants < contest.maxParticipants)

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <CardContent sx={{ flexGrow: 1 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Typography variant="h6" component="h2" sx={{ flexGrow: 1, mr: 1 }}>
            {contest.title}
          </Typography>
          <Box sx={{ display: 'flex', gap: 0.5 }}>
            <Chip
              label={contest.status}
              color={getStatusColor(contest.status)}
              size="small"
            />
            <Chip
              label={dateStatus}
              color={getStatusByDateColor(dateStatus)}
              size="small"
            />
          </Box>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
          <SportsIcon sx={{ mr: 1, color: 'text.secondary' }} />
          <Typography variant="body2" color="text.secondary">
            {contest.sportType}
          </Typography>
        </Box>

        {contest.description && (
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            {contest.description.length > 100
              ? `${contest.description.substring(0, 100)}...`
              : contest.description}
          </Typography>
        )}

        <Box sx={{ mb: 2 }}>
          <Typography variant="body2" color="text.secondary">
            <strong>Start:</strong> {formatDate(contest.startDate)}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            <strong>End:</strong> {formatDate(contest.endDate)}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
          <PeopleIcon sx={{ mr: 1, color: 'text.secondary' }} />
          <Typography variant="body2" color="text.secondary">
            {contest.currentParticipants}
            {contest.maxParticipants > 0 && ` / ${contest.maxParticipants}`} participants
          </Typography>
        </Box>

        <Typography variant="caption" color="text.secondary">
          Created {formatRelativeTime(contest.createdAt)}
        </Typography>
      </CardContent>

      <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
        <Box>
          {canJoin && onJoin && (
            <Button
              size="small"
              variant="contained"
              color="primary"
              onClick={() => onJoin(contest)}
            >
              Join Contest
            </Button>
          )}
          {isParticipant && onLeave && (
            <Button
              size="small"
              variant="outlined"
              color="secondary"
              onClick={() => onLeave(contest)}
            >
              Leave Contest
            </Button>
          )}
          {onViewParticipants && (
            <Button
              size="small"
              variant="text"
              startIcon={<PeopleIcon />}
              onClick={() => onViewParticipants(contest)}
            >
              View Participants
            </Button>
          )}
        </Box>

        {canEdit && (
          <Box>
            {onEdit && (
              <Tooltip title="Edit Contest">
                <IconButton
                  size="small"
                  onClick={() => onEdit(contest)}
                  color="primary"
                >
                  <EditIcon />
                </IconButton>
              </Tooltip>
            )}
            {onDelete && (
              <Tooltip title="Delete Contest">
                <IconButton
                  size="small"
                  onClick={() => onDelete(contest)}
                  color="error"
                >
                  <DeleteIcon />
                </IconButton>
              </Tooltip>
            )}
          </Box>
        )}
      </CardActions>
    </Card>
  )
}

export default ContestCard

import React, { useState } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
  Typography,
  Box,
  Chip,
  Divider,
  CircularProgress,
} from '@mui/material'
import {
  PersonRemove as RemoveIcon,
  AdminPanelSettings as AdminIcon,
  Person as PersonIcon,
} from '@mui/icons-material'
import { useContestParticipants, useLeaveContest } from '../../hooks/use-contests'
import type { Contest, Participant } from '../../types/contest.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface ParticipantListProps {
  open: boolean
  onClose: () => void
  contest: Contest | null
}

const getRoleIcon = (role: string) => {
  switch (role) {
    case 'admin':
      return <AdminIcon color="primary" />
    case 'participant':
      return <PersonIcon color="action" />
    default:
      return <PersonIcon color="action" />
  }
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'default'
    case 'banned':
      return 'error'
    default:
      return 'default'
  }
}

export const ParticipantList: React.FC<ParticipantListProps> = ({
  open,
  onClose,
  contest,
}) => {
  const [pagination] = useState({
    page: 1,
    limit: 50,
  })

  const { data, isLoading, isError, error } = useContestParticipants({
    contestId: contest?.id || 0,
    pagination,
  })

  const leaveContestMutation = useLeaveContest()

  const handleRemoveParticipant = (participant: Participant) => {
    if (window.confirm('Are you sure you want to remove this participant?')) {
      leaveContestMutation.mutate({
        contestId: participant.contestId,
      })
    }
  }

  if (!contest) {
    return null
  }

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">
            Participants - {contest.title}
          </Typography>
          <Chip
            label={`${contest.currentParticipants} participants`}
            color="primary"
            variant="outlined"
          />
        </Box>
      </DialogTitle>

      <DialogContent>
        {isLoading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}

        {isError && (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography color="error" variant="h6">
              Failed to load participants
            </Typography>
            <Typography color="text.secondary">
              {error?.message || 'An unknown error occurred'}
            </Typography>
          </Box>
        )}

        {data && data.participants.length === 0 && (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography color="text.secondary">
              No participants yet
            </Typography>
          </Box>
        )}

        {data && data.participants.length > 0 && (
          <List>
            {data.participants.map((participant, index) => (
              <React.Fragment key={participant.id}>
                <ListItem>
                  <Box sx={{ mr: 2 }}>
                    {getRoleIcon(participant.role)}
                  </Box>
                  <ListItemText
                    primary={
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Typography variant="body1">
                          User #{participant.userId}
                        </Typography>
                        <Chip
                          label={participant.role}
                          size="small"
                          color={participant.role === 'admin' ? 'primary' : 'default'}
                        />
                        <Chip
                          label={participant.status}
                          size="small"
                          color={getStatusColor(participant.status)}
                        />
                      </Box>
                    }
                    secondary={
                      <Typography variant="body2" color="text.secondary">
                        Joined {formatRelativeTime(participant.joinedAt)}
                      </Typography>
                    }
                  />
                  <ListItemSecondaryAction>
                    {participant.role !== 'admin' && (
                      <IconButton
                        edge="end"
                        color="error"
                        onClick={() => handleRemoveParticipant(participant)}
                        disabled={leaveContestMutation.isPending}
                        title="Remove participant"
                      >
                        <RemoveIcon />
                      </IconButton>
                    )}
                  </ListItemSecondaryAction>
                </ListItem>
                {index < data.participants.length - 1 && <Divider />}
              </React.Fragment>
            ))}
          </List>
        )}

        {data && data.pagination && data.pagination.totalPages > 1 && (
          <Box sx={{ mt: 2, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              Showing {data.participants.length} of {data.pagination.total} participants
            </Typography>
          </Box>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  )
}

export default ParticipantList

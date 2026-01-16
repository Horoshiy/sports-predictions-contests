import React from 'react'
import { Box, List, ListItem, ListItemText, ListItemSecondaryAction, IconButton, Chip, Typography, CircularProgress } from '@mui/material'
import { Delete as DeleteIcon } from '@mui/icons-material'
import { useTeamMembers, useRemoveMember } from '../../hooks/use-teams'
import { useAuth } from '../../contexts/AuthContext'
import { formatRelativeTime } from '../../utils/date-utils'
import type { Team } from '../../types/team.types'

interface TeamMembersProps {
  team: Team
}

export const TeamMembers: React.FC<TeamMembersProps> = ({ team }) => {
  const { user } = useAuth()
  const { data, isLoading, isError } = useTeamMembers({ teamId: team.id, pagination: { page: 1, limit: 20 } })
  const removeMemberMutation = useRemoveMember()

  const isCaptain = team.captainId === user?.id

  const handleRemove = (userId: number, userName: string) => {
    if (window.confirm(`Remove ${userName || `User #${userId}`} from the team?`)) {
      removeMemberMutation.mutate({ teamId: team.id, userId })
    }
  }

  if (isLoading) return <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}><CircularProgress /></Box>
  if (isError) return <Typography color="error">Failed to load members</Typography>

  return (
    <List>
      {data?.members.map((member) => (
        <ListItem key={member.id} divider>
          <ListItemText
            primary={
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Typography>{member.userName || `User #${member.userId}`}</Typography>
                <Chip label={member.role} size="small" color={member.role === 'captain' ? 'primary' : 'default'} />
              </Box>
            }
            secondary={`Joined ${formatRelativeTime(member.joinedAt)}`}
          />
          {isCaptain && member.role !== 'captain' && (
            <ListItemSecondaryAction>
              <IconButton edge="end" color="error" onClick={() => handleRemove(member.userId, member.userName)} disabled={removeMemberMutation.isPending}>
                <DeleteIcon />
              </IconButton>
            </ListItemSecondaryAction>
          )}
        </ListItem>
      ))}
      {(!data?.members || data.members.length === 0) && (
        <ListItem><ListItemText primary="No members found" /></ListItem>
      )}
    </List>
  )
}

export default TeamMembers

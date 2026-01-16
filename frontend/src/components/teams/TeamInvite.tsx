import React from 'react'
import { Box, Typography, TextField, IconButton, Tooltip } from '@mui/material'
import { ContentCopy, Refresh } from '@mui/icons-material'
import { useToast } from '../../contexts/ToastContext'
import { useRegenerateInviteCode } from '../../hooks/use-teams'

interface TeamInviteProps {
  teamId: number
  inviteCode: string
  isCaptain: boolean
}

export const TeamInvite: React.FC<TeamInviteProps> = ({ teamId, inviteCode, isCaptain }) => {
  const { showToast } = useToast()
  const regenerateMutation = useRegenerateInviteCode()

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(inviteCode)
      showToast('Invite code copied!', 'success')
    } catch {
      showToast('Failed to copy - please copy manually', 'warning')
    }
  }

  const handleRegenerate = () => {
    if (window.confirm('Regenerate invite code? The old code will stop working.')) {
      regenerateMutation.mutate(teamId)
    }
  }

  return (
    <Box sx={{ p: 2, bgcolor: 'grey.100', borderRadius: 1 }}>
      <Typography variant="subtitle2" color="text.secondary" gutterBottom>Invite Code</Typography>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <TextField value={inviteCode} InputProps={{ readOnly: true }} size="small" sx={{ fontFamily: 'monospace', flex: 1 }} />
        <Tooltip title="Copy code">
          <IconButton onClick={handleCopy} color="primary"><ContentCopy /></IconButton>
        </Tooltip>
        {isCaptain && (
          <Tooltip title="Regenerate code">
            <IconButton onClick={handleRegenerate} disabled={regenerateMutation.isPending} color="warning"><Refresh /></IconButton>
          </Tooltip>
        )}
      </Box>
    </Box>
  )
}

export default TeamInvite

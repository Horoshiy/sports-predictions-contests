import React from 'react'
import { Box, Chip, Tooltip, Typography } from '@mui/material'
import { TrendingUp, AccessTime } from '@mui/icons-material'

interface CoefficientIndicatorProps {
  coefficient: number
  tier: string
  hoursUntilEvent: number
  compact?: boolean
}

const getChipColor = (coefficient: number): 'success' | 'info' | 'warning' | 'default' => {
  if (coefficient >= 2.0) return 'success'
  if (coefficient >= 1.5) return 'info'
  if (coefficient >= 1.25) return 'warning'
  return 'default'
}

const getIconColor = (coefficient: number): 'success' | 'info' | 'warning' | 'inherit' => {
  if (coefficient >= 2.0) return 'success'
  if (coefficient >= 1.5) return 'info'
  if (coefficient >= 1.25) return 'warning'
  return 'inherit'
}

export const CoefficientIndicator: React.FC<CoefficientIndicatorProps> = ({
  coefficient,
  tier,
  hoursUntilEvent,
  compact = false,
}) => {
  const formatTimeRemaining = (hours: number): string => {
    if (hours >= 168) return `${Math.floor(hours / 24)} days`
    if (hours >= 24) return `${Math.floor(hours / 24)}d ${Math.floor(hours % 24)}h`
    return `${Math.floor(hours)}h ${Math.round((hours % 1) * 60)}m`
  }

  if (compact) {
    return (
      <Tooltip title={`${tier} - ${coefficient}x points (${formatTimeRemaining(hoursUntilEvent)} left)`}>
        <Chip
          icon={<TrendingUp />}
          label={`${coefficient}x`}
          size="small"
          color={getChipColor(coefficient)}
          variant="outlined"
        />
      </Tooltip>
    )
  }

  return (
    <Box sx={{ p: 1.5, bgcolor: 'action.hover', borderRadius: 1 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
        <TrendingUp color={getIconColor(coefficient)} fontSize="small" />
        <Typography variant="subtitle2" fontWeight="bold">
          {coefficient}x Points
        </Typography>
        <Chip label={tier} size="small" color={getChipColor(coefficient)} />
      </Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
        <AccessTime fontSize="small" color="action" />
        <Typography variant="caption" color="text.secondary">
          {formatTimeRemaining(hoursUntilEvent)} until coefficient drops
        </Typography>
      </Box>
    </Box>
  )
}

export default CoefficientIndicator

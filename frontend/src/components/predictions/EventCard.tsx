import React from 'react'
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Chip,
  Box,
} from '@mui/material'
import { SportsSoccer, Edit, CheckCircle } from '@mui/icons-material'
import type { Event, Prediction } from '../../types/prediction.types'
import { formatDate } from '../../utils/date-utils'
import { CoefficientIndicator } from './CoefficientIndicator'
import { usePotentialCoefficient } from '../../hooks/use-predictions'

interface EventCardProps {
  event: Event
  onPredict: (event: Event) => void
  existingPrediction?: Prediction
  disabled?: boolean
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'scheduled': return 'info'
    case 'live': return 'warning'
    case 'completed': return 'success'
    case 'cancelled': return 'error'
    default: return 'default'
  }
}

const canAcceptPredictions = (event: Event): boolean => {
  return event.status === 'scheduled' && new Date(event.eventDate) > new Date()
}

export const EventCard: React.FC<EventCardProps> = ({
  event,
  onPredict,
  existingPrediction,
  disabled = false,
}) => {
  const isPredictable = canAcceptPredictions(event) && !disabled
  const hasPrediction = !!existingPrediction
  const { data: coefficientData } = usePotentialCoefficient(isPredictable ? event.id : undefined)

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <CardContent sx={{ flexGrow: 1 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
          <Chip
            label={event.sportType}
            size="small"
            icon={<SportsSoccer />}
            variant="outlined"
          />
          <Chip
            label={event.status}
            size="small"
            color={getStatusColor(event.status)}
          />
        </Box>
        
        <Typography variant="h6" component="h3" gutterBottom>
          {event.title}
        </Typography>
        
        <Box sx={{ my: 2, textAlign: 'center' }}>
          <Typography variant="body1" fontWeight="medium">
            {event.homeTeam}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ my: 0.5 }}>
            vs
          </Typography>
          <Typography variant="body1" fontWeight="medium">
            {event.awayTeam}
          </Typography>
        </Box>
        
        <Typography variant="body2" color="text.secondary">
          {formatDate(event.eventDate)}
        </Typography>

        {coefficientData && coefficientData.coefficient > 1 && (
          <Box sx={{ mt: 1 }}>
            <CoefficientIndicator
              coefficient={coefficientData.coefficient}
              tier={coefficientData.tier}
              hoursUntilEvent={coefficientData.hoursUntilEvent}
              compact
            />
          </Box>
        )}

        {hasPrediction && (
          <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <CheckCircle color="success" fontSize="small" />
            <Typography variant="body2" color="success.main">
              Prediction submitted
            </Typography>
          </Box>
        )}
      </CardContent>
      
      <CardActions>
        {isPredictable && (
          <Button
            size="small"
            variant={hasPrediction ? 'outlined' : 'contained'}
            startIcon={hasPrediction ? <Edit /> : undefined}
            onClick={() => onPredict(event)}
            fullWidth
          >
            {hasPrediction ? 'Edit Prediction' : 'Make Prediction'}
          </Button>
        )}
        {!isPredictable && !hasPrediction && (
          <Typography variant="body2" color="text.secondary" sx={{ px: 1 }}>
            Predictions closed
          </Typography>
        )}
      </CardActions>
    </Card>
  )
}

export default EventCard

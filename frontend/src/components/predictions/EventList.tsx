import React, { useState } from 'react'
import {
  Box,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Typography,
  CircularProgress,
  Pagination,
  Alert,
} from '@mui/material'
import { useEvents } from '../../hooks/use-predictions'
import { useSports } from '../../hooks/use-sports'
import EventCard from './EventCard'
import type { Event, Prediction, ListEventsRequest } from '../../types/prediction.types'

interface EventListProps {
  contestId: number
  onPredict: (event: Event) => void
  userPredictions: Prediction[]
}

const statusOptions = ['', 'scheduled', 'live', 'completed', 'cancelled']

export const EventList: React.FC<EventListProps> = ({
  contestId,
  onPredict,
  userPredictions,
}) => {
  const [sportType, setSportType] = useState('')
  const [status, setStatus] = useState('scheduled')
  const [page, setPage] = useState(1)
  const pageSize = 12

  // Fetch sports dynamically from backend
  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 50 } })
  const sports = sportsData?.sports || []

  const request: ListEventsRequest = {
    sportType: sportType || undefined,
    status: status || undefined,
    pagination: { page, limit: pageSize },
  }

  const { data, isLoading, isError, error } = useEvents(request)

  const getPredictionForEvent = (eventId: number): Prediction | undefined => {
    return userPredictions.find(p => p.eventId === eventId)
  }

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
        <CircularProgress />
      </Box>
    )
  }

  if (isError) {
    return (
      <Alert severity="error">
        Failed to load events: {error?.message || 'Unknown error'}
      </Alert>
    )
  }

  const events = data?.events || []
  const totalPages = data?.pagination?.totalPages || 1

  return (
    <Box>
      <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
        <FormControl size="small" sx={{ minWidth: 150 }}>
          <InputLabel>Sport Type</InputLabel>
          <Select
            value={sportType}
            label="Sport Type"
            onChange={(e) => { setSportType(e.target.value); setPage(1) }}
          >
            <MenuItem value="">All Sports</MenuItem>
            {sports.map((sport) => (
              <MenuItem key={sport.id} value={sport.name}>{sport.name}</MenuItem>
            ))}
          </Select>
        </FormControl>

        <FormControl size="small" sx={{ minWidth: 150 }}>
          <InputLabel>Status</InputLabel>
          <Select
            value={status}
            label="Status"
            onChange={(e) => { setStatus(e.target.value); setPage(1) }}
          >
            <MenuItem value="">All Status</MenuItem>
            {statusOptions.filter(s => s).map((s) => (
              <MenuItem key={s} value={s}>{s.charAt(0).toUpperCase() + s.slice(1)}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>

      {events.length === 0 ? (
        <Typography color="text.secondary" sx={{ textAlign: 'center', py: 4 }}>
          No events found matching your filters
        </Typography>
      ) : (
        <>
          <Grid container spacing={2}>
            {events.map((event) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={event.id}>
                <EventCard
                  event={event}
                  onPredict={onPredict}
                  existingPrediction={getPredictionForEvent(event.id)}
                  disabled={!contestId}
                />
              </Grid>
            ))}
          </Grid>

          {totalPages > 1 && (
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
              <Pagination
                count={totalPages}
                page={page}
                onChange={(_, value) => setPage(value)}
                color="primary"
              />
            </Box>
          )}
        </>
      )}
    </Box>
  )
}

export default EventList

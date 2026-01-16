import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
} from '@mui/material'
import { useContests } from '../hooks/use-contests'
import { useUserPredictions, useSubmitPrediction, useUpdatePrediction, useEvent } from '../hooks/use-predictions'
import EventList from '../components/predictions/EventList'
import PredictionList from '../components/predictions/PredictionList'
import PredictionForm from '../components/predictions/PredictionForm'
import type { Event, Prediction } from '../types/prediction.types'

export const PredictionsPage: React.FC = () => {
  const [tabValue, setTabValue] = useState(0)
  const [selectedContestId, setSelectedContestId] = useState<number>(0)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null)
  const [selectedPrediction, setSelectedPrediction] = useState<Prediction | null>(null)
  const [editingEventId, setEditingEventId] = useState<number>(0)

  // Fetch active contests for dropdown
  const { data: contestsData } = useContests({ status: 'active' })
  const contests = contestsData?.contests || []

  // Fetch user predictions for selected contest (reduced limit for performance)
  const { data: predictionsData } = useUserPredictions({
    contestId: selectedContestId,
    pagination: { page: 1, limit: 20 },
  })
  const userPredictions = predictionsData?.predictions || []

  // Fetch event data when editing a prediction
  const { data: fetchedEvent } = useEvent(editingEventId)
  
  // Update selectedEvent when fetchedEvent loads
  useEffect(() => {
    if (fetchedEvent && editingEventId) {
      setSelectedEvent(fetchedEvent)
    }
  }, [fetchedEvent, editingEventId])

  const submitPredictionMutation = useSubmitPrediction()
  const updatePredictionMutation = useUpdatePrediction()

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue)
  }

  const handleContestChange = (contestId: number) => {
    setSelectedContestId(contestId)
  }

  const handlePredict = (event: Event) => {
    const existingPrediction = userPredictions.find(p => p.eventId === event.id)
    setSelectedEvent(event)
    setSelectedPrediction(existingPrediction || null)
    setIsFormOpen(true)
  }

  const handleEditPrediction = (prediction: Prediction) => {
    setSelectedPrediction(prediction)
    setEditingEventId(prediction.eventId)
    setIsFormOpen(true)
  }

  const handleFormSubmit = async (predictionData: string) => {
    if (!selectedEvent || !selectedContestId) return

    try {
      if (selectedPrediction) {
        await updatePredictionMutation.mutateAsync({
          id: selectedPrediction.id,
          predictionData,
        })
      } else {
        await submitPredictionMutation.mutateAsync({
          contestId: selectedContestId,
          eventId: selectedEvent.id,
          predictionData,
        })
      }
      setIsFormOpen(false)
      setSelectedEvent(null)
      setSelectedPrediction(null)
    } catch (error) {
      // Error handled by mutation
    }
  }

  const handleFormClose = () => {
    setIsFormOpen(false)
    setSelectedEvent(null)
    setSelectedPrediction(null)
    setEditingEventId(0)
  }

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Predictions
      </Typography>

      <Paper sx={{ p: 2, mb: 3 }}>
        <FormControl fullWidth size="small">
          <InputLabel>Select Contest</InputLabel>
          <Select
            value={selectedContestId || ''}
            label="Select Contest"
            onChange={(e) => handleContestChange(Number(e.target.value))}
          >
            <MenuItem value="">
              <em>Select a contest</em>
            </MenuItem>
            {contests.map((contest) => (
              <MenuItem key={contest.id} value={contest.id}>
                {contest.title} ({contest.sportType})
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Paper>

      {!selectedContestId ? (
        <Alert severity="info">
          Please select a contest to view events and make predictions.
        </Alert>
      ) : (
        <Paper sx={{ width: '100%' }}>
          <Tabs value={tabValue} onChange={handleTabChange}>
            <Tab label="Available Events" />
            <Tab label="My Predictions" />
          </Tabs>

          <Box sx={{ p: 3 }}>
            {tabValue === 0 && (
              <EventList
                contestId={selectedContestId}
                onPredict={handlePredict}
                userPredictions={userPredictions}
              />
            )}
            {tabValue === 1 && (
              <PredictionList
                contestId={selectedContestId}
                onEdit={handleEditPrediction}
              />
            )}
          </Box>
        </Paper>
      )}

      <PredictionForm
        open={isFormOpen}
        onClose={handleFormClose}
        onSubmit={handleFormSubmit}
        event={selectedEvent}
        prediction={selectedPrediction}
        loading={submitPredictionMutation.isPending || updatePredictionMutation.isPending}
      />
    </Box>
  )
}

export default PredictionsPage

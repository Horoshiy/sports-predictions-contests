import React, { useState, useEffect } from 'react'
import { Space, Typography, Tabs, Select, Alert } from 'antd'
import { useContests } from '../hooks/use-contests'
import { useUserPredictions, useSubmitPrediction, useUpdatePrediction, useEvent } from '../hooks/use-predictions'
import { showError, showSuccess } from '../utils/notification'
import EventList from '../components/predictions/EventList'
import PredictionList from '../components/predictions/PredictionList'
import PredictionForm from '../components/predictions/PredictionForm'
import type { Event, Prediction } from '../types/prediction.types'

const { Title } = Typography

export const PredictionsPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState('events')
  const [selectedContestId, setSelectedContestId] = useState<number>(0)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null)
  const [selectedPrediction, setSelectedPrediction] = useState<Prediction | null>(null)
  const [editingEventId, setEditingEventId] = useState<number>(0)

  const { data: contestsData } = useContests({ status: 'active' })
  const contests = contestsData?.contests || []

  const { data: predictionsData } = useUserPredictions({
    contestId: selectedContestId,
    pagination: { page: 1, limit: 20 },
  })
  const userPredictions = predictionsData?.predictions || []

  const { data: fetchedEvent } = useEvent(editingEventId)
  
  useEffect(() => {
    if (fetchedEvent && editingEventId) {
      setSelectedEvent(fetchedEvent)
    }
  }, [fetchedEvent, editingEventId])

  const submitPredictionMutation = useSubmitPrediction()
  const updatePredictionMutation = useUpdatePrediction()

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
        showSuccess('Prediction updated successfully')
      } else {
        await submitPredictionMutation.mutateAsync({
          contestId: selectedContestId,
          eventId: selectedEvent.id,
          predictionData,
        })
        showSuccess('Prediction submitted successfully')
      }
      setIsFormOpen(false)
      setSelectedEvent(null)
      setSelectedPrediction(null)
    } catch (error: any) {
      showError(error?.message || 'Failed to submit prediction')
    }
  }

  const handleFormClose = () => {
    setIsFormOpen(false)
    setSelectedEvent(null)
    setSelectedPrediction(null)
    setEditingEventId(0)
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%', padding: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={2}>Predictions</Title>
        <Select
          style={{ width: 250 }}
          placeholder="Select Contest"
          value={selectedContestId || undefined}
          onChange={setSelectedContestId}
        >
          {contests.map(c => (
            <Select.Option key={c.id} value={c.id}>{c.title}</Select.Option>
          ))}
        </Select>
      </div>

      {!selectedContestId && (
        <Alert message="Please select a contest to view events and make predictions" type="info" showIcon />
      )}

      {selectedContestId && (
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={[
            {
              key: 'events',
              label: 'Available Events',
              children: (
                <EventList
                  contestId={selectedContestId}
                  onPredict={handlePredict}
                  userPredictions={userPredictions}
                />
              ),
            },
            {
              key: 'predictions',
              label: 'My Predictions',
              children: (
                <PredictionList
                  contestId={selectedContestId}
                  onEdit={handleEditPrediction}
                />
              ),
            },
          ]}
        />
      )}

      <PredictionForm
        open={isFormOpen}
        onClose={handleFormClose}
        onSubmit={handleFormSubmit}
        event={selectedEvent}
        existingPrediction={selectedPrediction}
        loading={submitPredictionMutation.isPending || updatePredictionMutation.isPending}
      />
    </Space>
  )
}

export default PredictionsPage

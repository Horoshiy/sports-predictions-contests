import React from 'react'
import { Modal, Form, Input, Radio, Button, Alert, Typography, Divider, Space, InputNumber } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { 
  predictionSchema, 
  type PredictionFormData,
  formDataToPredictionData,
} from '../../utils/prediction-validation'
import type { Event, Prediction } from '../../types/prediction.types'
import { formatDate } from '../../utils/date-utils'

const { Text, Title } = Typography

interface PredictionFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (predictionData: string) => void
  event: Event | null
  existingPrediction?: Prediction | null
  loading?: boolean
}

export const PredictionForm: React.FC<PredictionFormProps> = ({
  open,
  onClose,
  onSubmit,
  event,
  existingPrediction,
  loading = false,
}) => {
  const isEdit = !!existingPrediction

  const defaultValues: PredictionFormData = React.useMemo(() => {
    if (existingPrediction) {
      try {
        const parsed = JSON.parse(existingPrediction.predictionData)
        return {
          eventId: existingPrediction.eventId,
          predictionType: parsed.predictionType || 'winner',
          winner: parsed.winner,
          homeScore: parsed.homeScore,
          awayScore: parsed.awayScore,
        }
      } catch {
        return {
          eventId: event?.id || 0,
          predictionType: 'winner',
        }
      }
    }
    return {
      eventId: event?.id || 0,
      predictionType: 'winner',
    }
  }, [existingPrediction, event])

  const { control, handleSubmit, reset, watch, formState: { errors } } = useForm<PredictionFormData>({
    resolver: zodResolver(predictionSchema),
    defaultValues,
    mode: 'onChange',
  })

  const predictionType = watch('predictionType')

  React.useEffect(() => {
    if (event) {
      reset({ ...defaultValues, eventId: event.id })
    }
  }, [event, defaultValues, reset])

  const handleFormSubmit = (data: PredictionFormData) => {
    const predictionData = formDataToPredictionData(data)
    onSubmit(predictionData)
  }

  const handleClose = () => {
    reset()
    onClose()
  }

  if (!event) return null

  return (
    <Modal
      open={open}
      title={isEdit ? 'Edit Prediction' : 'Make Prediction'}
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose} disabled={loading}>Cancel</Button>,
        <Button key="submit" type="primary" onClick={handleSubmit(handleFormSubmit)} loading={loading}>
          {isEdit ? 'Update' : 'Submit'} Prediction
        </Button>,
      ]}
      width={600}
    >
      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
        <div>
          <Title level={5} style={{ margin: 0 }}>{event.title}</Title>
          <Text type="secondary">{event.homeTeam} vs {event.awayTeam}</Text>
          <br />
          <Text type="secondary">{formatDate(event.eventDate)}</Text>
        </div>

        <Divider />

        <Form layout="vertical">
          <Controller name="predictionType" control={control} render={({ field }) => (
            <Form.Item label="Prediction Type">
              <Radio.Group {...field}>
                <Radio value="winner">Winner</Radio>
                <Radio value="score">Exact Score</Radio>
                <Radio value="combined">Winner + Score</Radio>
              </Radio.Group>
            </Form.Item>
          )} />

          {(predictionType === 'winner' || predictionType === 'combined') && (
            <Controller name="winner" control={control} render={({ field }) => (
              <Form.Item label="Winner" required validateStatus={errors.winner ? 'error' : ''} help={errors.winner?.message}>
                <Radio.Group {...field}>
                  <Radio value="home">{event.homeTeam}</Radio>
                  <Radio value="away">{event.awayTeam}</Radio>
                  <Radio value="draw">Draw</Radio>
                </Radio.Group>
              </Form.Item>
            )} />
          )}

          {(predictionType === 'score' || predictionType === 'combined') && (
            <Space>
              <Controller name="homeScore" control={control} render={({ field }) => (
                <Form.Item label={`${event.homeTeam} Score`} validateStatus={errors.homeScore ? 'error' : ''} help={errors.homeScore?.message}>
                  <InputNumber {...field} min={0} style={{ width: 100 }} />
                </Form.Item>
              )} />
              <Controller name="awayScore" control={control} render={({ field }) => (
                <Form.Item label={`${event.awayTeam} Score`} validateStatus={errors.awayScore ? 'error' : ''} help={errors.awayScore?.message}>
                  <InputNumber {...field} min={0} style={{ width: 100 }} />
                </Form.Item>
              )} />
            </Space>
          )}
        </Form>
      </Space>
    </Modal>
  )
}

export default PredictionForm

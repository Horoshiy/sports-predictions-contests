import React from 'react'
import { Card, Button, Tag, Space, Typography } from 'antd'
import { TrophyOutlined, EditOutlined, CheckCircleOutlined } from '@ant-design/icons'
import type { Event, Prediction } from '../../types/prediction.types'
import { formatDate } from '../../utils/date-utils'
import { CoefficientIndicator } from './CoefficientIndicator'
import { usePotentialCoefficient } from '../../hooks/use-predictions'

const { Text, Title } = Typography

interface EventCardProps {
  event: Event
  onPredict: (event: Event) => void
  existingPrediction?: Prediction
  disabled?: boolean
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    scheduled: 'processing',
    live: 'warning',
    completed: 'success',
    cancelled: 'error',
  }
  return colors[status] || 'default'
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
  const { data: coefficientData } = usePotentialCoefficient(event.id)
  const coefficient = typeof coefficientData === 'number' ? coefficientData : coefficientData?.coefficient
  const canPredict = canAcceptPredictions(event) && !disabled

  return (
    <Card
      actions={[
        canPredict ? (
          <Button
            key="predict"
            type={existingPrediction ? 'default' : 'primary'}
            icon={existingPrediction ? <EditOutlined /> : <TrophyOutlined />}
            onClick={() => onPredict(event)}
          >
            {existingPrediction ? 'Edit Prediction' : 'Make Prediction'}
          </Button>
        ) : existingPrediction ? (
          <Space key="predicted">
            <CheckCircleOutlined style={{ color: '#52c41a' }} />
            <Text type="success">Predicted</Text>
          </Space>
        ) : null,
      ]}
    >
      <Space direction="vertical" size="small" style={{ width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Tag color={getStatusColor(event.status)}>{event.status}</Tag>
          {coefficient && <CoefficientIndicator coefficient={coefficient} tier="standard" hoursUntilEvent={0} />}
        </div>
        <Title level={5} style={{ margin: 0 }}>{event.title}</Title>
        <Text type="secondary">{event.homeTeam} vs {event.awayTeam}</Text>
        <Text type="secondary" style={{ fontSize: 12 }}>
          {formatDate(event.eventDate)}
        </Text>
      </Space>
    </Card>
  )
}

export default EventCard

import React, { useState } from 'react'
import { Card, Checkbox, Button, Typography, Space, Tag, Alert, Divider } from 'antd'
import { ThunderboltOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons'

const { Title, Text } = Typography

export interface RiskyEvent {
  slug: string
  name: string
  name_en?: string
  points: number
}

interface RiskyPredictionFormProps {
  events: RiskyEvent[]
  maxSelections: number
  onSubmit: (selections: string[]) => void
  loading?: boolean
  existingSelections?: string[]
  matchTitle?: string
}

export const RiskyPredictionForm: React.FC<RiskyPredictionFormProps> = ({
  events,
  maxSelections,
  onSubmit,
  loading = false,
  existingSelections = [],
  matchTitle,
}) => {
  const [selections, setSelections] = useState<string[]>(existingSelections)

  const handleToggle = (slug: string) => {
    setSelections(prev => {
      if (prev.includes(slug)) {
        return prev.filter(s => s !== slug)
      }
      if (prev.length >= maxSelections) {
        return prev // Don't add if at max
      }
      return [...prev, slug]
    })
  }

  const calculatePotential = () => {
    let potential = 0
    let risk = 0
    for (const slug of selections) {
      const event = events.find(e => e.slug === slug)
      if (event) {
        potential += event.points
        risk += event.points
      }
    }
    return { potential, risk }
  }

  const { potential, risk } = calculatePotential()

  return (
    <Card 
      title={
        <Space>
          <ThunderboltOutlined style={{ color: '#faad14' }} />
          <span>Рисковый прогноз</span>
        </Space>
      }
    >
      {matchTitle && (
        <Title level={5} style={{ marginBottom: 16 }}>{matchTitle}</Title>
      )}

      <Alert
        type="warning"
        message={`Выбери до ${maxSelections} событий`}
        description="Угадал событие → получаешь очки. Не угадал → теряешь очки."
        showIcon
        style={{ marginBottom: 16 }}
      />

      <Space direction="vertical" style={{ width: '100%' }}>
        {events.map(event => {
          const isSelected = selections.includes(event.slug)
          const isDisabled = !isSelected && selections.length >= maxSelections

          return (
            <Card
              key={event.slug}
              size="small"
              style={{
                borderColor: isSelected ? '#1890ff' : undefined,
                backgroundColor: isSelected ? '#e6f7ff' : undefined,
                opacity: isDisabled ? 0.5 : 1,
              }}
              hoverable={!isDisabled}
              onClick={() => !isDisabled && handleToggle(event.slug)}
            >
              <Space style={{ width: '100%', justifyContent: 'space-between' }}>
                <Space>
                  <Checkbox 
                    checked={isSelected} 
                    disabled={isDisabled}
                    onChange={() => handleToggle(event.slug)}
                  />
                  <div>
                    <Text strong>{event.name}</Text>
                    {event.name_en && (
                      <Text type="secondary" style={{ display: 'block', fontSize: 12 }}>
                        {event.name_en}
                      </Text>
                    )}
                  </div>
                </Space>
                <Space>
                  <Tag color="green" icon={<CheckCircleOutlined />}>
                    +{event.points}
                  </Tag>
                  <Tag color="red" icon={<CloseCircleOutlined />}>
                    −{event.points}
                  </Tag>
                </Space>
              </Space>
            </Card>
          )
        })}
      </Space>

      <Divider />

      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <div>
          <Text>Выбрано: </Text>
          <Text strong>{selections.length} / {maxSelections}</Text>
        </div>
        <Space>
          <Tag color="green">Потенциал: +{potential}</Tag>
          <Tag color="red">Риск: −{risk}</Tag>
        </Space>
      </Space>

      <Button
        type="primary"
        icon={<ThunderboltOutlined />}
        onClick={() => onSubmit(selections)}
        loading={loading}
        disabled={selections.length === 0}
        style={{ marginTop: 16, width: '100%' }}
        size="large"
      >
        {existingSelections.length > 0 ? 'Изменить прогноз' : 'Сделать прогноз'}
      </Button>
    </Card>
  )
}

export default RiskyPredictionForm

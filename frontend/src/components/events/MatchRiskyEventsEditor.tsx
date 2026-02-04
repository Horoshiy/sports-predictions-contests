import React, { useState, useEffect } from 'react'
import {
  Card,
  List,
  InputNumber,
  Switch,
  Space,
  Typography,
  Tag,
  Button,
  Spin,
  Empty,
  Alert,
  Popconfirm,
  Divider,
} from 'antd'
import {
  SaveOutlined,
  UndoOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons'
import {
  useMatchRiskyEvents,
  useSetMatchRiskyEvents,
  useSetMatchEventOutcome,
} from '../../hooks/use-risky-events'
import type { MatchRiskyEvent } from '../../types/risky-events.types'

const { Text, Title } = Typography

interface MatchRiskyEventsEditorProps {
  eventId: number
  contestId: number
  isMatchFinished?: boolean
  readOnly?: boolean
}

interface EditableEvent {
  riskyEventTypeId: number
  slug: string
  name: string
  icon?: string
  points: number
  isEnabled: boolean
  outcome?: boolean | null
  originalPoints: number
  originalEnabled: boolean
}

const MatchRiskyEventsEditor: React.FC<MatchRiskyEventsEditorProps> = ({
  eventId,
  contestId,
  isMatchFinished = false,
  readOnly = false,
}) => {
  const [editableEvents, setEditableEvents] = useState<EditableEvent[]>([])
  const [hasChanges, setHasChanges] = useState(false)

  const {
    data: matchEventsData,
    isLoading,
    isError,
    error,
    refetch,
  } = useMatchRiskyEvents(eventId, contestId)

  const setMatchEventsMutation = useSetMatchRiskyEvents()
  const setOutcomeMutation = useSetMatchEventOutcome()

  // Initialize editable events when data loads
  useEffect(() => {
    if (matchEventsData?.events) {
      setEditableEvents(
        matchEventsData.events.map((e) => ({
          riskyEventTypeId: e.riskyEventTypeId,
          slug: e.slug,
          name: e.name,
          icon: e.icon,
          points: e.points,
          isEnabled: e.isEnabled,
          outcome: e.outcome,
          originalPoints: e.points,
          originalEnabled: e.isEnabled,
        }))
      )
      setHasChanges(false)
    }
  }, [matchEventsData])

  // Check for changes
  useEffect(() => {
    const changed = editableEvents.some(
      (e) => e.points !== e.originalPoints || e.isEnabled !== e.originalEnabled
    )
    setHasChanges(changed)
  }, [editableEvents])

  const handlePointsChange = (riskyEventTypeId: number, points: number | null) => {
    setEditableEvents((prev) =>
      prev.map((e) =>
        e.riskyEventTypeId === riskyEventTypeId
          ? { ...e, points: points ?? 0 }
          : e
      )
    )
  }

  const handleEnabledChange = (riskyEventTypeId: number, enabled: boolean) => {
    setEditableEvents((prev) =>
      prev.map((e) =>
        e.riskyEventTypeId === riskyEventTypeId
          ? { ...e, isEnabled: enabled }
          : e
      )
    )
  }

  const handleSave = () => {
    setMatchEventsMutation.mutate(
      {
        eventId,
        events: editableEvents.map((e) => ({
          riskyEventTypeId: e.riskyEventTypeId,
          points: e.points,
          isEnabled: e.isEnabled,
        })),
      },
      {
        onSuccess: () => {
          refetch()
        },
      }
    )
  }

  const handleReset = () => {
    setEditableEvents((prev) =>
      prev.map((e) => ({
        ...e,
        points: e.originalPoints,
        isEnabled: e.originalEnabled,
      }))
    )
  }

  const handleSetOutcome = (riskyEventTypeId: number, outcome: boolean) => {
    setOutcomeMutation.mutate(
      {
        eventId,
        riskyEventTypeId,
        outcome,
      },
      {
        onSuccess: () => {
          refetch()
        },
      }
    )
  }

  const renderOutcomeStatus = (event: EditableEvent) => {
    if (event.outcome === true) {
      return (
        <Tag icon={<CheckCircleOutlined />} color="success">
          Произошло
        </Tag>
      )
    }
    if (event.outcome === false) {
      return (
        <Tag icon={<CloseCircleOutlined />} color="error">
          Не произошло
        </Tag>
      )
    }
    return (
      <Tag icon={<QuestionCircleOutlined />} color="default">
        Ожидает
      </Tag>
    )
  }

  if (isLoading) {
    return (
      <Card size="small">
        <div style={{ textAlign: 'center', padding: 24 }}>
          <Spin />
          <div style={{ marginTop: 8 }}>
            <Text type="secondary">Загрузка событий...</Text>
          </div>
        </div>
      </Card>
    )
  }

  if (isError) {
    return (
      <Alert
        message="Ошибка загрузки"
        description={error?.message || 'Не удалось загрузить рисковые события'}
        type="error"
        showIcon
      />
    )
  }

  if (!matchEventsData?.events?.length) {
    return (
      <Card size="small">
        <Empty
          description="Нет рисковых событий для этого матча"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      </Card>
    )
  }

  return (
    <Card
      size="small"
      title={
        <Space>
          <span>⚡ Рисковые события</span>
          {matchEventsData.maxSelections > 0 && (
            <Tag color="blue">
              Макс. выбор: {matchEventsData.maxSelections}
            </Tag>
          )}
        </Space>
      }
      extra={
        !readOnly && (
          <Space>
            {hasChanges && (
              <Button
                icon={<UndoOutlined />}
                size="small"
                onClick={handleReset}
              >
                Сброс
              </Button>
            )}
            <Button
              type="primary"
              icon={<SaveOutlined />}
              size="small"
              onClick={handleSave}
              disabled={!hasChanges}
              loading={setMatchEventsMutation.isPending}
            >
              Сохранить
            </Button>
          </Space>
        )
      }
    >
      <List
        size="small"
        dataSource={editableEvents}
        renderItem={(event) => (
          <List.Item
            style={{
              opacity: event.isEnabled ? 1 : 0.5,
              background: event.isEnabled ? 'transparent' : '#f5f5f5',
              padding: '8px 12px',
              borderRadius: 4,
              marginBottom: 4,
            }}
          >
            <Space style={{ width: '100%', justifyContent: 'space-between' }}>
              {/* Event name with icon */}
              <Space>
                <Text style={{ fontSize: 16 }}>{event.icon || '⚡'}</Text>
                <Text
                  strong={event.isEnabled}
                  delete={!event.isEnabled}
                  style={{ minWidth: 150 }}
                >
                  {event.name}
                </Text>
              </Space>

              {/* Points and controls */}
              <Space size="middle">
                {/* Points editor */}
                {!readOnly ? (
                  <Space size="small">
                    <Text type="secondary">Очки:</Text>
                    <InputNumber
                      value={event.points}
                      onChange={(val) =>
                        handlePointsChange(event.riskyEventTypeId, val)
                      }
                      min={-100}
                      max={100}
                      size="small"
                      style={{ width: 70 }}
                      disabled={!event.isEnabled}
                    />
                    {event.points !== event.originalPoints && (
                      <Tag color="orange" style={{ marginLeft: 4 }}>
                        было: {event.originalPoints}
                      </Tag>
                    )}
                  </Space>
                ) : (
                  <Tag color={event.points > 0 ? 'green' : event.points < 0 ? 'red' : 'default'}>
                    {event.points > 0 ? '+' : ''}{event.points} очков
                  </Tag>
                )}

                {/* Enabled toggle */}
                {!readOnly && (
                  <Switch
                    checked={event.isEnabled}
                    onChange={(checked) =>
                      handleEnabledChange(event.riskyEventTypeId, checked)
                    }
                    checkedChildren="Вкл"
                    unCheckedChildren="Выкл"
                    size="small"
                  />
                )}

                {/* Outcome section (only after match finished) */}
                {isMatchFinished && event.isEnabled && (
                  <>
                    <Divider type="vertical" />
                    {!readOnly ? (
                      <Space size="small">
                        {renderOutcomeStatus(event)}
                        <Popconfirm
                          title="Отметить как произошедшее?"
                          onConfirm={() =>
                            handleSetOutcome(event.riskyEventTypeId, true)
                          }
                          okText="Да"
                          cancelText="Нет"
                        >
                          <Button
                            size="small"
                            type={event.outcome === true ? 'primary' : 'default'}
                            icon={<CheckCircleOutlined />}
                            loading={setOutcomeMutation.isPending}
                          />
                        </Popconfirm>
                        <Popconfirm
                          title="Отметить как НЕ произошедшее?"
                          onConfirm={() =>
                            handleSetOutcome(event.riskyEventTypeId, false)
                          }
                          okText="Да"
                          cancelText="Нет"
                        >
                          <Button
                            size="small"
                            type={event.outcome === false ? 'primary' : 'default'}
                            danger={event.outcome === false}
                            icon={<CloseCircleOutlined />}
                            loading={setOutcomeMutation.isPending}
                          />
                        </Popconfirm>
                      </Space>
                    ) : (
                      renderOutcomeStatus(event)
                    )}
                  </>
                )}
              </Space>
            </Space>
          </List.Item>
        )}
      />

      {/* Info about max selections */}
      {matchEventsData.maxSelections > 0 && (
        <div style={{ marginTop: 12, textAlign: 'center' }}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            Участники могут выбрать до {matchEventsData.maxSelections} событий
          </Text>
        </div>
      )}
    </Card>
  )
}

export default MatchRiskyEventsEditor

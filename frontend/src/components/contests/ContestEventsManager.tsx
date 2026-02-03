import React, { useState, useMemo } from 'react'
import {
  Card,
  Table,
  Tag,
  Space,
  Button,
  Typography,
  Empty,
  Alert,
  Tooltip,
  Modal,
} from 'antd'
import {
  ThunderboltOutlined,
  CalendarOutlined,
  TrophyOutlined,
  EyeOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useEvents } from '../../hooks/use-predictions'
import type { Event } from '../../types/prediction.types'
import type { Contest } from '../../types/contest.types'
import { formatDateTime } from '../../utils/date-utils'
import MatchRiskyEventsEditor from '../events/MatchRiskyEventsEditor'

const { Text, Title } = Typography

interface ContestEventsManagerProps {
  contest: Contest
}

const statusColors: Record<string, string> = {
  scheduled: 'default',
  live: 'warning',
  finished: 'success',
  cancelled: 'error',
  postponed: 'processing',
}

const ContestEventsManager: React.FC<ContestEventsManagerProps> = ({ contest }) => {
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null)
  const [isRiskyModalOpen, setIsRiskyModalOpen] = useState(false)

  // Parse rules JSON to check contest type
  const parsedRules = React.useMemo(() => {
    if (!contest.rules) return null
    try {
      return JSON.parse(contest.rules) as { contestType?: string }
    } catch {
      return null
    }
  }, [contest.rules])

  const isRiskyContest = parsedRules?.contestType === 'risky'

  const { data, isLoading, isError, error } = useEvents({
    contestId: contest.id,
    pagination: { page: 1, limit: 100 },
  })

  const handleOpenRiskyEditor = (event: Event) => {
    setSelectedEvent(event)
    setIsRiskyModalOpen(true)
  }

  const handleCloseRiskyEditor = () => {
    setSelectedEvent(null)
    setIsRiskyModalOpen(false)
  }

  const columns: ColumnsType<Event> = useMemo(() => {
    const baseColumns: ColumnsType<Event> = [
      {
        title: 'ID',
        dataIndex: 'id',
        key: 'id',
        width: 60,
      },
      {
        title: 'Матч',
        key: 'title',
        width: 250,
        render: (_, event) => (
          <Space direction="vertical" size={0}>
            <Text strong>{event.title || `${event.homeTeam} vs ${event.awayTeam}`}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {event.homeTeam} — {event.awayTeam}
            </Text>
          </Space>
        ),
      },
      {
        title: 'Дата',
        dataIndex: 'eventDate',
        key: 'eventDate',
        width: 150,
        render: (date: string) => (
          <Space>
            <CalendarOutlined />
            {formatDateTime(date)}
          </Space>
        ),
      },
      {
        title: 'Статус',
        dataIndex: 'status',
        key: 'status',
        width: 100,
        render: (status: string) => (
          <Tag color={statusColors[status] || 'default'}>{status}</Tag>
        ),
      },
      {
        title: 'Результат',
        key: 'result',
        width: 100,
        render: (_, event) => {
          if (event.status === 'completed' && event.resultData) {
            try {
              const result = JSON.parse(event.resultData)
              return (
                <Tag color="green" icon={<TrophyOutlined />}>
                  {result.homeScore ?? '-'} : {result.awayScore ?? '-'}
                </Tag>
              )
            } catch {
              return <Text type="secondary">—</Text>
            }
          }
          return <Text type="secondary">—</Text>
        },
      },
    ]

    // Add risky events column for risky contests
    if (isRiskyContest) {
      baseColumns.push({
        title: 'Risky Events',
        key: 'riskyEvents',
        width: 150,
        render: (_, event) => (
          <Button
            type="link"
            icon={<ThunderboltOutlined />}
            onClick={() => handleOpenRiskyEditor(event)}
          >
            Настроить
          </Button>
        ),
      })
    }

    return baseColumns
  }, [isRiskyContest])

  if (isError) {
    return (
      <Alert
        message="Ошибка загрузки"
        description={error?.message || 'Не удалось загрузить события конкурса'}
        type="error"
        showIcon
      />
    )
  }

  if (!isLoading && (!data?.events || data.events.length === 0)) {
    return (
      <Card>
        <Empty
          description={
            <Space direction="vertical" size={0}>
              <Text>Нет событий в этом конкурсе</Text>
              <Text type="secondary">
                Добавьте матчи через раздел "События" или API
              </Text>
            </Space>
          }
        />
      </Card>
    )
  }

  return (
    <>
      <Card
        title={
          <Space>
            <TrophyOutlined />
            <span>События конкурса: {contest.title}</span>
            {isRiskyContest && (
              <Tag color="orange" icon={<ThunderboltOutlined />}>
                Risky
              </Tag>
            )}
          </Space>
        }
        extra={
          <Text type="secondary">
            Всего: {data?.pagination?.total || data?.events?.length || 0} событий
          </Text>
        }
      >
        <Table
          columns={columns}
          dataSource={data?.events ?? []}
          rowKey="id"
          loading={isLoading}
          pagination={false}
          size="small"
        />
      </Card>

      {/* Risky Events Editor Modal */}
      <Modal
        title={
          <Space>
            <ThunderboltOutlined />
            <span>
              Risky события:{' '}
              {selectedEvent?.title ||
                `${selectedEvent?.homeTeam} vs ${selectedEvent?.awayTeam}`}
            </span>
          </Space>
        }
        open={isRiskyModalOpen}
        onCancel={handleCloseRiskyEditor}
        footer={null}
        width={700}
        destroyOnClose
      >
        {selectedEvent && (
          <MatchRiskyEventsEditor
            eventId={selectedEvent.id}
            contestId={contest.id}
            isMatchFinished={selectedEvent.status === 'completed'}
          />
        )}
      </Modal>
    </>
  )
}

export default ContestEventsManager

import React, { useEffect, useState } from 'react'
import { Table, Input, Tag, Space, Typography, Alert, Spin, DatePicker } from 'antd'
import { SearchOutlined, CalendarOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import dayjs from 'dayjs'

const { Text } = Typography
const { RangePicker } = DatePicker

interface Event {
  id: number
  title: string
  sport_type: string
  home_team: string
  away_team: string
  event_date: string
  status: string
}

interface EventSelectorProps {
  value?: number[]
  onChange?: (eventIds: number[]) => void
  maxEvents?: number
  minEvents?: number
}

export const EventSelector: React.FC<EventSelectorProps> = ({
  value = [],
  onChange,
  maxEvents = 30,
  minEvents = 5,
}) => {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchText, setSearchText] = useState('')
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs | null, dayjs.Dayjs | null] | null>(null)

  useEffect(() => {
    fetchEvents()
  }, [])

  const fetchEvents = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await fetch('/api/v1/events?limit=500&status=scheduled')
      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`)
      }
      const data = await response.json()
      if (data.events) {
        setEvents(data.events)
      } else {
        setEvents([])
      }
    } catch (err) {
      console.error('Failed to fetch events:', err)
      setError('Не удалось загрузить матчи. Попробуйте обновить страницу.')
    } finally {
      setLoading(false)
    }
  }

  const filteredEvents = events.filter(event => {
    // Text search
    const searchLower = searchText.toLowerCase()
    const matchesSearch = !searchText || 
      event.title.toLowerCase().includes(searchLower) ||
      event.home_team.toLowerCase().includes(searchLower) ||
      event.away_team.toLowerCase().includes(searchLower) ||
      event.sport_type.toLowerCase().includes(searchLower)
    
    // Date range filter
    let matchesDate = true
    if (dateRange && dateRange[0] && dateRange[1]) {
      const eventDate = dayjs(event.event_date)
      matchesDate = eventDate.isAfter(dateRange[0].startOf('day')) && 
                    eventDate.isBefore(dateRange[1].endOf('day'))
    }
    
    return matchesSearch && matchesDate
  })

  const columns: ColumnsType<Event> = [
    {
      title: 'Матч',
      key: 'match',
      render: (_, record) => (
        <Space direction="vertical" size={0}>
          <Text strong>{record.home_team} vs {record.away_team}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>{record.title}</Text>
        </Space>
      ),
    },
    {
      title: 'Лига',
      dataIndex: 'sport_type',
      key: 'sport_type',
      width: 120,
      render: (type: string) => <Tag color="blue">{type}</Tag>,
      filters: [...new Set(events.map(e => e.sport_type))].map(t => ({ text: t, value: t })),
      onFilter: (value, record) => record.sport_type === value,
    },
    {
      title: 'Дата',
      dataIndex: 'event_date',
      key: 'event_date',
      width: 150,
      render: (date: string) => (
        <Space>
          <CalendarOutlined />
          {dayjs(date).format('DD.MM.YYYY HH:mm')}
        </Space>
      ),
      sorter: (a, b) => dayjs(a.event_date).unix() - dayjs(b.event_date).unix(),
      defaultSortOrder: 'ascend',
    },
    {
      title: 'Статус',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => {
        const colors: Record<string, string> = {
          scheduled: 'green',
          live: 'orange',
          completed: 'default',
          cancelled: 'red',
        }
        return <Tag color={colors[status] || 'default'}>{status}</Tag>
      },
    },
  ]

  const rowSelection = {
    selectedRowKeys: value,
    onChange: (selectedRowKeys: React.Key[]) => {
      const ids = selectedRowKeys as number[]
      if (ids.length <= maxEvents) {
        onChange?.(ids)
      }
    },
    getCheckboxProps: (record: Event) => ({
      disabled: !value.includes(record.id) && value.length >= maxEvents,
    }),
  }

  const isValidSelection = value.length >= minEvents && value.length <= maxEvents

  return (
    <div>
      <Space direction="vertical" style={{ width: '100%', marginBottom: 16 }}>
        <Space wrap>
          <Input
            placeholder="Поиск по командам..."
            prefix={<SearchOutlined />}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 250 }}
            allowClear
          />
          <RangePicker
            value={dateRange}
            onChange={(dates) => setDateRange(dates)}
            format="DD.MM.YYYY"
            placeholder={['От', 'До']}
          />
        </Space>
        
        {error && (
          <Alert
            type="error"
            message={error}
            showIcon
            closable
            onClose={() => setError(null)}
          />
        )}
        
        <Alert
          type={isValidSelection ? 'success' : 'warning'}
          message={
            <Space>
              <Text>
                Выбрано матчей: <Text strong>{value.length}</Text>
              </Text>
              <Text type="secondary">
                (мин. {minEvents}, макс. {maxEvents})
              </Text>
            </Space>
          }
          showIcon
        />
      </Space>

      <Spin spinning={loading}>
        <Table
          rowKey="id"
          columns={columns}
          dataSource={filteredEvents}
          rowSelection={rowSelection}
          pagination={{
            pageSize: 20,
            showSizeChanger: true,
            showTotal: (total, range) => `${range[0]}-${range[1]} из ${total} матчей`,
          }}
          size="small"
          scroll={{ y: 400 }}
        />
      </Spin>
    </div>
  )
}

export default EventSelector

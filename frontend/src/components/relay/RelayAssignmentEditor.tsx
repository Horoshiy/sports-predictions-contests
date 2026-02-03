import React, { useEffect, useState, useMemo } from 'react'
import { Card, Table, Select, Button, Alert, Space, Typography, Tag, Spin, message, Divider } from 'antd'
import { UserOutlined, SwapOutlined, SaveOutlined, CalendarOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import dayjs from 'dayjs'

const { Text, Title } = Typography

interface Event {
  id: number
  title: string
  sport_type: string
  home_team: string
  away_team: string
  event_date: string
  status: string
}

interface TeamMember {
  id: number
  user_id: number
  username: string
  role: 'captain' | 'member'
}

interface Assignment {
  user_id: number
  event_id: number
}

interface RelayAssignmentEditorProps {
  contestId: number
  teamId: number
  events: Event[]
  members: TeamMember[]
  initialAssignments?: Assignment[]
  onSave?: (assignments: Assignment[]) => Promise<void>
  loading?: boolean
}

export const RelayAssignmentEditor: React.FC<RelayAssignmentEditorProps> = ({
  contestId,
  teamId,
  events,
  members,
  initialAssignments = [],
  onSave,
  loading = false,
}) => {
  // Map: eventId -> userId
  const [assignments, setAssignments] = useState<Map<number, number>>(new Map())
  const [saving, setSaving] = useState(false)
  const [hasChanges, setHasChanges] = useState(false)

  // Initialize assignments from props
  useEffect(() => {
    const map = new Map<number, number>()
    initialAssignments.forEach(a => {
      map.set(a.event_id, a.user_id)
    })
    setAssignments(map)
    setHasChanges(false)
  }, [initialAssignments])

  // Stats per member
  const memberStats = useMemo(() => {
    const stats = new Map<number, number>()
    members.forEach(m => stats.set(m.user_id, 0))
    assignments.forEach((userId) => {
      stats.set(userId, (stats.get(userId) || 0) + 1)
    })
    return stats
  }, [assignments, members])

  // Handle assignment change
  const handleAssign = (eventId: number, userId: number | null) => {
    const newAssignments = new Map(assignments)
    if (userId === null) {
      newAssignments.delete(eventId)
    } else {
      newAssignments.set(eventId, userId)
    }
    setAssignments(newAssignments)
    setHasChanges(true)
  }

  // Auto-distribute events evenly
  const handleAutoDistribute = () => {
    const newAssignments = new Map<number, number>()
    const memberIds = members.map(m => m.user_id)
    
    events.forEach((event, index) => {
      const memberIndex = index % memberIds.length
      newAssignments.set(event.id, memberIds[memberIndex])
    })
    
    setAssignments(newAssignments)
    setHasChanges(true)
    message.success('Матчи распределены равномерно')
  }

  // Save assignments
  const handleSave = async () => {
    if (!onSave) return
    
    setSaving(true)
    try {
      const assignmentList: Assignment[] = []
      assignments.forEach((userId, eventId) => {
        assignmentList.push({ user_id: userId, event_id: eventId })
      })
      await onSave(assignmentList)
      setHasChanges(false)
      message.success('Распределение сохранено!')
    } catch (error) {
      message.error('Ошибка сохранения: ' + (error as Error).message)
    } finally {
      setSaving(false)
    }
  }

  // Table columns
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
    },
    {
      title: 'Назначен',
      key: 'assigned',
      width: 200,
      render: (_, record) => {
        const assignedUserId = assignments.get(record.id)
        return (
          <Select
            style={{ width: '100%' }}
            placeholder="Выбери участника"
            value={assignedUserId}
            onChange={(value) => handleAssign(record.id, value)}
            allowClear
          >
            {members.map(member => (
              <Select.Option key={member.user_id} value={member.user_id}>
                <Space>
                  <UserOutlined />
                  {member.username}
                  {member.role === 'captain' && <Tag color="gold" style={{ marginLeft: 4 }}>К</Tag>}
                </Space>
              </Select.Option>
            ))}
          </Select>
        )
      },
    },
  ]

  const assignedCount = assignments.size
  const totalCount = events.length
  const isComplete = assignedCount === totalCount

  return (
    <Card 
      title={
        <Space>
          <SwapOutlined />
          <span>Распределение матчей</span>
        </Space>
      }
      extra={
        <Space>
          <Button 
            icon={<SwapOutlined />} 
            onClick={handleAutoDistribute}
            disabled={loading || saving}
          >
            Авто-распределить
          </Button>
          <Button 
            type="primary" 
            icon={<SaveOutlined />}
            onClick={handleSave}
            loading={saving}
            disabled={!hasChanges || !isComplete}
          >
            Сохранить
          </Button>
        </Space>
      }
    >
      {/* Stats */}
      <Alert
        type={isComplete ? 'success' : 'warning'}
        message={
          <Space split={<Divider type="vertical" />}>
            <Text>
              Назначено: <Text strong>{assignedCount}</Text> / {totalCount}
            </Text>
            {members.map(member => (
              <Text key={member.user_id}>
                {member.username}: <Text strong>{memberStats.get(member.user_id) || 0}</Text>
              </Text>
            ))}
          </Space>
        }
        showIcon
        style={{ marginBottom: 16 }}
      />

      {!isComplete && (
        <Alert
          type="info"
          message="Назначь участника для каждого матча. Только после этого можно сохранить распределение."
          style={{ marginBottom: 16 }}
        />
      )}

      <Spin spinning={loading}>
        <Table
          rowKey="id"
          columns={columns}
          dataSource={events}
          pagination={{
            pageSize: 15,
            showSizeChanger: true,
            showTotal: (total, range) => `${range[0]}-${range[1]} из ${total}`,
          }}
          size="small"
        />
      </Spin>
    </Card>
  )
}

export default RelayAssignmentEditor

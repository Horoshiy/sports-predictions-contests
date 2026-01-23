import React from 'react'
import { Card, Button, Tag, Space, Typography, Tooltip } from 'antd'
import { EditOutlined, DeleteOutlined, TeamOutlined, TrophyOutlined } from '@ant-design/icons'
import type { Contest } from '../../types/contest.types'
import { formatDate, formatRelativeTime, getContestStatusByDate } from '../../utils/date-utils'

const { Text, Title } = Typography

interface ContestCardProps {
  contest: Contest
  onEdit?: (contest: Contest) => void
  onDelete?: (contest: Contest) => void
  onViewParticipants?: (contest: Contest) => void
  onJoin?: (contest: Contest) => void
  onLeave?: (contest: Contest) => void
  isParticipant?: boolean
  canEdit?: boolean
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'draft':
      return 'default'
    case 'active':
      return 'success'
    case 'completed':
      return 'blue'
    case 'cancelled':
      return 'error'
    default:
      return 'default'
  }
}

const getStatusByDateColor = (status: 'upcoming' | 'active' | 'completed') => {
  switch (status) {
    case 'upcoming':
      return 'warning'
    case 'active':
      return 'success'
    case 'completed':
      return 'blue'
    default:
      return 'default'
  }
}

export const ContestCard: React.FC<ContestCardProps> = ({
  contest,
  onEdit,
  onDelete,
  onViewParticipants,
  onJoin,
  onLeave,
  isParticipant = false,
  canEdit = false,
}) => {
  const dateStatus = getContestStatusByDate(contest.startDate, contest.endDate)
  const canJoin = contest.status === 'active' && !isParticipant && 
    (contest.maxParticipants === 0 || contest.currentParticipants < contest.maxParticipants)

  const actions = []
  
  if (canJoin && onJoin) {
    actions.push(
      <Button type="primary" size="small" onClick={() => onJoin(contest)}>
        Join Contest
      </Button>
    )
  }
  
  if (isParticipant && onLeave) {
    actions.push(
      <Button danger size="small" onClick={() => onLeave(contest)}>
        Leave Contest
      </Button>
    )
  }
  
  if (onViewParticipants) {
    actions.push(
      <Button type="link" size="small" icon={<TeamOutlined />} onClick={() => onViewParticipants(contest)}>
        View Participants
      </Button>
    )
  }

  const extra = canEdit ? (
    <Space>
      {onEdit && (
        <Tooltip title="Edit Contest">
          <Button type="text" size="small" icon={<EditOutlined />} onClick={() => onEdit(contest)} />
        </Tooltip>
      )}
      {onDelete && (
        <Tooltip title="Delete Contest">
          <Button type="text" size="small" danger icon={<DeleteOutlined />} onClick={() => onDelete(contest)} />
        </Tooltip>
      )}
    </Space>
  ) : undefined

  return (
    <Card
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Title level={5} style={{ margin: 0 }}>{contest.title}</Title>
          <Space size={4}>
            <Tag color={getStatusColor(contest.status)}>{contest.status}</Tag>
            <Tag color={getStatusByDateColor(dateStatus)}>{dateStatus}</Tag>
          </Space>
        </div>
      }
      extra={extra}
      actions={actions.length > 0 ? actions : undefined}
      style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
      bodyStyle={{ flexGrow: 1 }}
    >
      <Space direction="vertical" size="small" style={{ width: '100%' }}>
        <Space>
          <TrophyOutlined style={{ color: '#8c8c8c' }} />
          <Text type="secondary">{contest.sportType}</Text>
        </Space>

        {contest.description && (
          <Text type="secondary">
            {contest.description.length > 100
              ? `${contest.description.substring(0, 100)}...`
              : contest.description}
          </Text>
        )}

        <div>
          <Text type="secondary"><strong>Start:</strong> {formatDate(contest.startDate)}</Text>
          <br />
          <Text type="secondary"><strong>End:</strong> {formatDate(contest.endDate)}</Text>
        </div>

        <Space>
          <TeamOutlined style={{ color: '#8c8c8c' }} />
          <Text type="secondary">
            {contest.currentParticipants}
            {contest.maxParticipants > 0 && ` / ${contest.maxParticipants}`} participants
          </Text>
        </Space>

        <Text type="secondary" style={{ fontSize: '12px' }}>
          Created {formatRelativeTime(contest.createdAt)}
        </Text>
      </Space>
    </Card>
  )
}

export default ContestCard

import React from 'react'
import { Card, Typography, Button, Tag, Space, Avatar, Tooltip } from 'antd'
import {
  CheckOutlined,
  CloseOutlined,
  DeleteOutlined,
  EyeOutlined,
  UserOutlined,
  CalendarOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons'
import type { Challenge, ChallengeWithUserInfo } from '../../types/challenge.types'
import { CHALLENGE_STATUSES } from '../../types/challenge.types'
import { formatRelativeTime } from '../../utils/date-utils'

const { Text, Paragraph } = Typography

interface ChallengeCardProps {
  challenge: ChallengeWithUserInfo
  currentUserId: number
  onAccept?: (challenge: Challenge) => void
  onDecline?: (challenge: Challenge) => void
  onWithdraw?: (challenge: Challenge) => void
  onView?: (challenge: Challenge) => void
  loading?: boolean
}

const getStatusColor = (status: Challenge['status']) => {
  return CHALLENGE_STATUSES[status]?.color || 'default'
}

export const ChallengeCard: React.FC<ChallengeCardProps> = ({
  challenge,
  currentUserId,
  onAccept,
  onDecline,
  onWithdraw,
  onView,
  loading = false,
}) => {
  const isCurrentUserChallenger = challenge.challengerId === currentUserId
  const isCurrentUserOpponent = challenge.opponentId === currentUserId
  const statusInfo = CHALLENGE_STATUSES[challenge.status]

  // Check if challenge is expired
  const isExpired = challenge.status === 'pending' && new Date() > new Date(challenge.expiresAt)
  const canAccept = challenge.status === 'pending' && !isExpired && isCurrentUserOpponent
  const canDecline = challenge.status === 'pending' && isCurrentUserOpponent
  const canWithdraw = challenge.status === 'pending' && isCurrentUserChallenger

  const getOpponentName = () => {
    if (isCurrentUserChallenger) {
      return challenge.opponentName || `User ${challenge.opponentId}`
    } else {
      return challenge.challengerName || `User ${challenge.challengerId}`
    }
  }

  const getOpponentRole = () => {
    if (isCurrentUserChallenger) {
      return 'Opponent'
    } else {
      return 'Challenger'
    }
  }

  const getScoreDisplay = () => {
    if (challenge.status !== 'completed') return null

    const userScore = isCurrentUserChallenger ? challenge.challengerScore : challenge.opponentScore
    const opponentScore = isCurrentUserChallenger ? challenge.opponentScore : challenge.challengerScore
    
    let result = 'Tie'
    let resultColor: 'default' | 'success' | 'error' = 'default'
    
    if (challenge.winnerId === currentUserId) {
      result = 'Won'
      resultColor = 'success'
    } else if (challenge.winnerId && challenge.winnerId !== currentUserId) {
      result = 'Lost'
      resultColor = 'error'
    }

    return (
      <Space style={{ marginTop: 8 }}>
        <Text strong>Score: {userScore} - {opponentScore}</Text>
        <Tag color={resultColor}>{result}</Tag>
      </Space>
    )
  }

  const actions = []
  
  // Accept button
  if (canAccept && onAccept) {
    actions.push(
      <Button
        key="accept"
        type="primary"
        icon={<CheckOutlined />}
        onClick={() => onAccept(challenge)}
        loading={loading}
        style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
      >
        Accept
      </Button>
    )
  }

  // Decline button
  if (canDecline && onDecline) {
    actions.push(
      <Button
        key="decline"
        danger
        icon={<CloseOutlined />}
        onClick={() => onDecline(challenge)}
        loading={loading}
      >
        Decline
      </Button>
    )
  }

  // Withdraw button
  if (canWithdraw && onWithdraw) {
    actions.push(
      <Button
        key="withdraw"
        icon={<DeleteOutlined />}
        onClick={() => onWithdraw(challenge)}
        loading={loading}
      >
        Withdraw
      </Button>
    )
  }

  // View details button
  if (onView) {
    actions.push(
      <Tooltip key="view" title="View Details">
        <Button
          icon={<EyeOutlined />}
          onClick={() => onView(challenge)}
          disabled={loading}
        />
      </Tooltip>
    )
  }

  return (
    <Card
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Text strong>Challenge #{challenge.id}</Text>
          <Tag color={isExpired ? 'default' : getStatusColor(challenge.status)}>
            {isExpired ? 'Expired' : statusInfo.label}
          </Tag>
        </div>
      }
      actions={actions}
      style={{ height: '100%' }}
    >
      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
        {/* Opponent info */}
        <Space>
          <Avatar icon={<UserOutlined />} size={32} />
          <div>
            <Text strong>{getOpponentName()}</Text>
            <br />
            <Text type="secondary" style={{ fontSize: 12 }}>{getOpponentRole()}</Text>
          </div>
        </Space>

        {/* Event info */}
        <Space>
          <CalendarOutlined />
          <Text>{challenge.eventTitle || `Event #${challenge.eventId}`}</Text>
        </Space>

        {/* Message */}
        {challenge.message && (
          <Paragraph italic type="secondary" style={{ marginBottom: 0 }}>
            "{challenge.message}"
          </Paragraph>
        )}

        {/* Timing info */}
        <Space>
          <ClockCircleOutlined />
          <Text type="secondary" style={{ fontSize: 12 }}>
            Created {formatRelativeTime(new Date(challenge.createdAt))}
          </Text>
        </Space>

        {/* Expiration for pending challenges */}
        {challenge.status === 'pending' && (
          <Text type={isExpired ? 'danger' : 'warning'} style={{ fontSize: 12 }}>
            {isExpired 
              ? 'Expired' 
              : `Expires ${formatRelativeTime(new Date(challenge.expiresAt))}`
            }
          </Text>
        )}

        {/* Acceptance date */}
        {challenge.acceptedAt && (
          <Text type="secondary" style={{ fontSize: 12 }}>
            Accepted {formatRelativeTime(new Date(challenge.acceptedAt))}
          </Text>
        )}

        {/* Score display for completed challenges */}
        {getScoreDisplay()}
      </Space>
    </Card>
  )
}

export default ChallengeCard

import React, { useState } from 'react'
import { Table, Button, Tag, Space, Alert, Spin, Tooltip } from 'antd'
import {
  CheckOutlined,
  CloseOutlined,
  DeleteOutlined,
  EyeOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import challengeService from '../../services/challenge-service'
import type { Challenge } from '../../types/challenge.types'
import { CHALLENGE_STATUSES } from '../../types/challenge.types'
import { formatDate, formatRelativeTime } from '../../utils/date-utils'

interface ChallengeListProps {
  userId: number
  statusFilter?: Challenge['status']
  onViewChallenge?: (challenge: Challenge) => void
  onCreateChallenge?: () => void
}

export const ChallengeList: React.FC<ChallengeListProps> = ({
  userId,
  statusFilter,
  onViewChallenge,
  onCreateChallenge,
}) => {
  const queryClient = useQueryClient()
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  })

  // Fetch challenges
  const {
    data: challengesData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['challenges', userId, statusFilter, pagination],
    queryFn: () =>
      challengeService.listUserChallenges({
        userId,
        status: statusFilter,
        pagination: {
          page: pagination.pageIndex + 1,
          limit: pagination.pageSize,
        },
      }),
  })

  // Accept challenge mutation
  const acceptChallengeMutation = useMutation({
    mutationFn: (challengeId: number) =>
      challengeService.acceptChallenge({ id: challengeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['challenges'] })
    },
  })

  // Decline challenge mutation
  const declineChallengeMutation = useMutation({
    mutationFn: (challengeId: number) =>
      challengeService.declineChallenge({ id: challengeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['challenges'] })
    },
  })

  // Withdraw challenge mutation
  const withdrawChallengeMutation = useMutation({
    mutationFn: (challengeId: number) =>
      challengeService.withdrawChallenge({ id: challengeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['challenges'] })
    },
  })

  const handleAcceptChallenge = (challenge: Challenge) => {
    acceptChallengeMutation.mutate(challenge.id)
  }

  const handleDeclineChallenge = (challenge: Challenge) => {
    declineChallengeMutation.mutate(challenge.id)
  }

  const handleWithdrawChallenge = (challenge: Challenge) => {
    withdrawChallengeMutation.mutate(challenge.id)
  }

  const columns: ColumnsType<Challenge> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: 'Challenger',
      dataIndex: 'challengerId',
      key: 'challengerId',
      render: (challengerId: number) => challengerId === userId ? 'You' : `User ${challengerId}`,
    },
    {
      title: 'Opponent',
      dataIndex: 'opponentId',
      key: 'opponentId',
      render: (opponentId: number) => opponentId === userId ? 'You' : `User ${opponentId}`,
    },
    {
      title: 'Event',
      dataIndex: 'eventId',
      key: 'eventId',
      render: (eventId: number) => `Event ${eventId}`,
    },
    {
      title: 'Message',
      dataIndex: 'message',
      key: 'message',
      ellipsis: true,
      render: (message: string) => message || 'No message',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: Challenge['status']) => {
        const statusInfo = CHALLENGE_STATUSES[status]
        return <Tag color={statusInfo.color}>{statusInfo.label}</Tag>
      },
    },
    {
      title: 'Expires',
      dataIndex: 'expiresAt',
      key: 'expiresAt',
      render: (expiresAt: string, record) => {
        if (record.status !== 'pending') return null
        const expires = new Date(expiresAt)
        const now = new Date()
        return now > expires ? (
          <Tag color="error">Expired</Tag>
        ) : (
          <Tag color="warning">{formatRelativeTime(expires)}</Tag>
        )
      },
    },
    {
      title: 'Created',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (createdAt: string) => formatDate(createdAt),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, challenge) => {
        const isCurrentUserChallenger = challenge.challengerId === userId
        const isCurrentUserOpponent = challenge.opponentId === userId
        const statusInfo = challengeService.getChallengeStatusInfo(challenge)

        return (
          <Space>
            {onViewChallenge && (
              <Tooltip title="View Details">
                <Button
                  size="small"
                  icon={<EyeOutlined />}
                  onClick={() => onViewChallenge(challenge)}
                />
              </Tooltip>
            )}
            {isCurrentUserOpponent && statusInfo.canAccept && (
              <Tooltip title="Accept Challenge">
                <Button
                  size="small"
                  type="primary"
                  icon={<CheckOutlined />}
                  onClick={() => handleAcceptChallenge(challenge)}
                  loading={acceptChallengeMutation.isPending}
                  style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
                />
              </Tooltip>
            )}
            {isCurrentUserOpponent && statusInfo.canDecline && (
              <Tooltip title="Decline Challenge">
                <Button
                  size="small"
                  danger
                  icon={<CloseOutlined />}
                  onClick={() => handleDeclineChallenge(challenge)}
                  loading={declineChallengeMutation.isPending}
                />
              </Tooltip>
            )}
            {isCurrentUserChallenger && statusInfo.canWithdraw && (
              <Tooltip title="Withdraw Challenge">
                <Button
                  size="small"
                  icon={<DeleteOutlined />}
                  onClick={() => handleWithdrawChallenge(challenge)}
                  loading={withdrawChallengeMutation.isPending}
                />
              </Tooltip>
            )}
          </Space>
        )
      },
    },
  ]

  if (error) {
    return (
      <Alert
        message="Failed to load challenges"
        description={error.message}
        type="error"
        showIcon
        style={{ marginTop: 16 }}
      />
    )
  }

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      {onCreateChallenge && (
        <Button
          type="primary"
          icon={<CheckOutlined />}
          onClick={onCreateChallenge}
        >
          Create Challenge
        </Button>
      )}
      <Table
        columns={columns}
        dataSource={challengesData?.challenges || []}
        rowKey="id"
        loading={isLoading || acceptChallengeMutation.isPending || declineChallengeMutation.isPending || withdrawChallengeMutation.isPending}
        pagination={{
          current: pagination.pageIndex + 1,
          pageSize: pagination.pageSize,
          total: challengesData?.pagination.total || 0,
          onChange: (page, pageSize) => setPagination({ pageIndex: page - 1, pageSize }),
        }}
      />
    </Space>
  )
}

export default ChallengeList

import React, { useState } from 'react'
import { Modal, Button, List, Tag, Space, Spin, Alert, Typography, Avatar } from 'antd'
import { UserDeleteOutlined, CrownOutlined, UserOutlined } from '@ant-design/icons'
import { useContestParticipants, useLeaveContest } from '../../hooks/use-contests'
import type { Contest, Participant } from '../../types/contest.types'
import { formatRelativeTime } from '../../utils/date-utils'

const { Text } = Typography

interface ParticipantListProps {
  open: boolean
  onClose: () => void
  contest: Contest | null
}

const getRoleIcon = (role: string) => {
  switch (role) {
    case 'admin':
      return <CrownOutlined style={{ color: '#1890ff' }} />
    case 'participant':
      return <UserOutlined />
    default:
      return <UserOutlined />
  }
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'default'
    case 'banned':
      return 'error'
    default:
      return 'default'
  }
}

export const ParticipantList: React.FC<ParticipantListProps> = ({
  open,
  onClose,
  contest,
}) => {
  const [pagination] = useState({
    page: 1,
    limit: 50,
  })

  const { data, isLoading, isError, error } = useContestParticipants({
    contestId: contest?.id || 0,
    pagination,
  })

  const leaveContestMutation = useLeaveContest()

  const handleRemoveParticipant = (participant: Participant) => {
    if (window.confirm('Are you sure you want to remove this participant?')) {
      leaveContestMutation.mutate({
        contestId: participant.contestId,
      })
    }
  }

  return (
    <Modal
      open={open}
      title={`Participants - ${contest?.title || 'Contest'}`}
      onCancel={onClose}
      footer={[
        <Button key="close" onClick={onClose}>Close</Button>,
      ]}
      width={600}
    >
      {isLoading && (
        <div style={{ textAlign: 'center', padding: '32px 0' }}>
          <Spin size="large" />
        </div>
      )}

      {isError && (
        <Alert
          message="Error loading participants"
          description={error?.message}
          type="error"
          showIcon
        />
      )}

      {!isLoading && !isError && (
        <>
          <div style={{ marginBottom: 16 }}>
            <Text type="secondary">
              {data?.participants?.length || 0} participant(s)
              {contest?.maxParticipants && ` / ${contest.maxParticipants} max`}
            </Text>
          </div>

          <List
            dataSource={data?.participants || []}
            renderItem={(participant) => (
              <List.Item
                actions={[
                  <Button
                    key="remove"
                    danger
                    size="small"
                    icon={<UserDeleteOutlined />}
                    onClick={() => handleRemoveParticipant(participant)}
                    loading={leaveContestMutation.isPending}
                  >
                    Remove
                  </Button>,
                ]}
              >
                <List.Item.Meta
                  avatar={<Avatar icon={getRoleIcon(participant.role)} />}
                  title={
                    <Space>
                      <Text>User {participant.userId}</Text>
                      <Tag color={getStatusColor(participant.status)}>{participant.status}</Tag>
                      {participant.role === 'admin' && <Tag color="blue">Admin</Tag>}
                    </Space>
                  }
                  description={
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      Joined {formatRelativeTime(participant.joinedAt)}
                    </Text>
                  }
                />
              </List.Item>
            )}
          />
        </>
      )}
    </Modal>
  )
}

export default ParticipantList

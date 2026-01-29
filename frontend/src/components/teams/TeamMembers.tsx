import React from 'react'
import { List, Tag, Button, Spin, Typography, Popconfirm, Skeleton } from 'antd'
import { DeleteOutlined, CrownOutlined } from '@ant-design/icons'
import { useTeamMembers, useRemoveMember } from '../../hooks/use-teams'
import { useAuth } from '../../contexts/AuthContext'
import { formatRelativeTime } from '../../utils/date-utils'
import type { Team } from '../../types/team.types'

const { Text } = Typography

interface TeamMembersProps {
  team: Team
}

export const TeamMembers: React.FC<TeamMembersProps> = ({ team }) => {
  const { user } = useAuth()
  const { data, isLoading, isError } = useTeamMembers({ teamId: team.id, pagination: { page: 1, limit: 20 } })
  const removeMemberMutation = useRemoveMember()

  const isCaptain = team.captainId === user?.id

  const handleRemove = (userId: number) => {
    removeMemberMutation.mutate({ teamId: team.id, userId })
  }

  if (isLoading) {
    return (
      <List
        dataSource={[1, 2, 3]}
        renderItem={() => (
          <List.Item>
            <Skeleton avatar active />
          </List.Item>
        )}
      />
    )
  }

  if (isError) return <Text type="danger">Failed to load members</Text>

  return (
    <List
      dataSource={data?.members || []}
      locale={{ emptyText: 'No members found' }}
      renderItem={(member) => (
        <List.Item
          actions={
            isCaptain && member.role !== 'captain'
              ? [
                  <Popconfirm
                    key="remove"
                    title="Remove member"
                    description={`Remove ${member.userName || `User #${member.userId}`} from the team?`}
                    onConfirm={() => handleRemove(member.userId)}
                    okText="Yes"
                    cancelText="No"
                  >
                    <Button type="text" danger icon={<DeleteOutlined />} loading={removeMemberMutation.isPending} />
                  </Popconfirm>,
                ]
              : []
          }
        >
          <List.Item.Meta
            title={
              <span>
                {member.userName || `User #${member.userId}`}
                {' '}
                <Tag color={member.role === 'captain' ? 'gold' : 'default'} icon={member.role === 'captain' ? <CrownOutlined /> : undefined}>
                  {member.role}
                </Tag>
              </span>
            }
            description={`Joined ${formatRelativeTime(member.joinedAt)}`}
          />
        </List.Item>
      )}
    />
  )
}

export default TeamMembers

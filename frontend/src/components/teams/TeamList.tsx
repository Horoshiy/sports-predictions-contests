import React, { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Table, Button, Tag, Tooltip, Space, Alert, Modal, Empty } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined, TeamOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useTeams, useDeleteTeam } from '../../hooks/use-teams'
import { useAuth } from '../../contexts/AuthContext'
import type { Team } from '../../types/team.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface TeamListProps {
  onCreateTeam: () => void
  onEditTeam: (team: Team) => void
  onViewMembers: (team: Team) => void
  myTeamsOnly?: boolean
}

export const TeamList: React.FC<TeamListProps> = ({ onCreateTeam, onEditTeam, onViewMembers, myTeamsOnly = false }) => {
  const { user } = useAuth()
  const [searchParams, setSearchParams] = useSearchParams()
  const [pagination, setPagination] = useState({
    pageIndex: parseInt(searchParams.get('page') || '0', 10) || 0,
    pageSize: parseInt(searchParams.get('limit') || '10', 10) || 10,
  })

  useEffect(() => {
    const params = new URLSearchParams(searchParams)
    params.set('page', pagination.pageIndex.toString())
    params.set('limit', pagination.pageSize.toString())
    setSearchParams(params, { replace: true })
  }, [pagination, searchParams, setSearchParams])

  const { data, isLoading, isError, error } = useTeams({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    myTeamsOnly,
  })

  const deleteTeamMutation = useDeleteTeam()

  const handleDelete = (team: Team) => {
    Modal.confirm({
      title: 'Delete Team',
      content: `Are you sure you want to delete "${team.name}"? This action cannot be undone and all members will be removed.`,
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: () => {
        deleteTeamMutation.mutate(team.id)
      },
    })
  }

  const columns: ColumnsType<Team> = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Name', dataIndex: 'name', key: 'name', width: 200 },
    {
      title: 'Members',
      key: 'members',
      width: 100,
      render: (_, team) => {
        const count = team.currentMembers ?? 0
        return `${count}/${team.maxMembers}`
      },
    },
    {
      title: 'Status',
      dataIndex: 'isActive',
      key: 'isActive',
      width: 100,
      render: (isActive: boolean) => <Tag color={isActive ? 'success' : 'default'}>{isActive ? 'Active' : 'Inactive'}</Tag>,
    },
    {
      title: 'Created',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 120,
      render: (date: string) => formatRelativeTime(date),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 150,
      render: (_, team) => (
        <Space>
          <Tooltip title="View Members">
            <Button 
              icon={<TeamOutlined />} 
              size="small" 
              onClick={() => onViewMembers(team)}
              data-testid="view-members-button"
            />
          </Tooltip>
          <Tooltip title="Edit">
            <Button type="primary" icon={<EditOutlined />} size="small" onClick={() => onEditTeam(team)} />
          </Tooltip>
          <Tooltip title="Delete">
            <Button danger icon={<DeleteOutlined />} size="small" onClick={() => handleDelete(team)} loading={deleteTeamMutation.isPending} />
          </Tooltip>
        </Space>
      ),
    },
  ]

  if (isError) {
    return <Alert message="Error" description={error?.message} type="error" showIcon />
  }

  const hasTeams = data?.teams && data.teams.length > 0

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Button type="primary" icon={<PlusOutlined />} onClick={onCreateTeam} data-testid="create-team-button">
        Create Team
      </Button>
      {!isLoading && !hasTeams ? (
        <Empty
          description={myTeamsOnly ? "You haven't joined any teams yet" : "No teams available"}
          style={{ padding: '48px 0' }}
        >
          <Button type="primary" onClick={onCreateTeam}>Create Your First Team</Button>
        </Empty>
      ) : (
        <Table
          columns={columns}
          dataSource={data?.teams ?? []}
          rowKey="id"
          loading={isLoading || deleteTeamMutation.isPending}
          pagination={{
            current: pagination.pageIndex + 1,
            pageSize: pagination.pageSize,
            total: data?.pagination?.total ?? 0,
            onChange: (page, pageSize) => setPagination({ pageIndex: page - 1, pageSize }),
          }}
        />
      )}
    </Space>
  )
}

export default TeamList

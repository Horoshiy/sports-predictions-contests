import React, { useMemo, useState } from 'react'
import { Table, Button, Tag, Tooltip, Space, Select, Avatar, Alert } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useTeams, useDeleteTeam, useSports } from '../../hooks/use-sports'
import type { Team } from '../../types/sports.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface TeamListProps {
  onCreateTeam: () => void
  onEditTeam: (team: Team) => void
}

export const TeamList: React.FC<TeamListProps> = ({ onCreateTeam, onEditTeam }) => {
  const [sportFilter, setSportFilter] = useState<number | ''>('')
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })

  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 } })
  const sportsMap = useMemo(() => new Map(sportsData?.sports?.map(s => [s.id, s.name]) || []), [sportsData])

  const { data, isLoading, isError, error } = useTeams({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    sportId: sportFilter || undefined,
  })

  const deleteMutation = useDeleteTeam()

  const handleDelete = (team: Team) => {
    if (window.confirm(`Delete "${team.name}"?`)) {
      deleteMutation.mutate(team.id)
    }
  }

  const columns: ColumnsType<Team> = useMemo(() => [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    {
      title: '',
      dataIndex: 'logoUrl',
      key: 'logo',
      width: 50,
      render: (logoUrl: string, team) => <Avatar src={logoUrl}>{team.name[0]}</Avatar>,
    },
    { title: 'Name', dataIndex: 'name', key: 'name', width: 150 },
    { title: 'Short', dataIndex: 'shortName', key: 'shortName', width: 80 },
    {
      title: 'Sport',
      dataIndex: 'sportId',
      key: 'sportId',
      width: 100,
      render: (sportId: number) => sportsMap.get(sportId) || '-',
    },
    { title: 'Country', dataIndex: 'country', key: 'country', width: 100 },
    {
      title: 'Status',
      dataIndex: 'isActive',
      key: 'isActive',
      width: 90,
      render: (isActive: boolean) => <Tag color={isActive ? 'success' : 'default'}>{isActive ? 'Active' : 'Inactive'}</Tag>,
    },
    {
      title: 'Created',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 110,
      render: (createdAt: string) => formatRelativeTime(createdAt),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 120,
      render: (_, team) => (
        <Space>
          <Tooltip title="Edit">
            <Button type="primary" icon={<EditOutlined />} size="small" onClick={() => onEditTeam(team)} />
          </Tooltip>
          <Tooltip title="Delete">
            <Button danger icon={<DeleteOutlined />} size="small" onClick={() => handleDelete(team)} loading={deleteMutation.isPending} />
          </Tooltip>
        </Space>
      ),
    },
  ], [sportsMap, deleteMutation.isPending])

  if (isError) {
    return <Alert message="Error" description={error?.message} type="error" showIcon />
  }

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Space>
        <Button type="primary" icon={<PlusOutlined />} onClick={onCreateTeam}>Add Team</Button>
        <Select
          style={{ minWidth: 150 }}
          placeholder="Filter by Sport"
          value={sportFilter}
          onChange={setSportFilter}
          allowClear
        >
          <Select.Option value="">All Sports</Select.Option>
          {sportsData?.sports?.map(s => <Select.Option key={s.id} value={s.id}>{s.name}</Select.Option>)}
        </Select>
      </Space>
      <Table
        columns={columns}
        dataSource={data?.teams ?? []}
        rowKey="id"
        loading={isLoading}
        pagination={{
          current: pagination.pageIndex + 1,
          pageSize: pagination.pageSize,
          total: data?.pagination?.total ?? 0,
          onChange: (page, pageSize) => setPagination({ pageIndex: page - 1, pageSize }),
        }}
      />
    </Space>
  )
}

export default TeamList

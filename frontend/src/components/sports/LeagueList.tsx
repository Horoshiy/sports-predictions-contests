import React, { useMemo, useState } from 'react'
import { Table, Button, Tag, Tooltip, Space, Select, Alert } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useLeagues, useDeleteLeague, useSports } from '../../hooks/use-sports'
import type { League } from '../../types/sports.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface LeagueListProps {
  onCreateLeague: () => void
  onEditLeague: (league: League) => void
}

export const LeagueList: React.FC<LeagueListProps> = ({ onCreateLeague, onEditLeague }) => {
  const [sportFilter, setSportFilter] = useState<number | ''>('')
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })

  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 } })
  const sportsMap = useMemo(() => new Map(sportsData?.sports?.map(s => [s.id, s.name]) || []), [sportsData])

  const { data, isLoading, isError, error } = useLeagues({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    sportId: sportFilter || undefined,
  })

  const deleteMutation = useDeleteLeague()

  const handleDelete = (league: League) => {
    if (window.confirm(`Delete "${league.name}"?`)) {
      deleteMutation.mutate(league.id)
    }
  }

  const columns: ColumnsType<League> = useMemo(() => [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Name', dataIndex: 'name', key: 'name', width: 180 },
    {
      title: 'Sport',
      dataIndex: 'sportId',
      key: 'sportId',
      width: 120,
      render: (sportId: number) => sportsMap.get(sportId) || '-',
    },
    { title: 'Country', dataIndex: 'country', key: 'country', width: 100 },
    { title: 'Season', dataIndex: 'season', key: 'season', width: 100 },
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
      render: (createdAt: string) => formatRelativeTime(createdAt),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 120,
      render: (_, league) => (
        <Space>
          <Tooltip title="Edit">
            <Button 
              type="primary" 
              icon={<EditOutlined />} 
              size="small" 
              onClick={() => onEditLeague(league)}
              aria-label={`Edit ${league.name}`}
            />
          </Tooltip>
          <Tooltip title="Delete">
            <Button 
              danger 
              icon={<DeleteOutlined />} 
              size="small" 
              onClick={() => handleDelete(league)} 
              loading={deleteMutation.isPending}
              aria-label={`Delete ${league.name}`}
            />
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
        <Button type="primary" icon={<PlusOutlined />} onClick={onCreateLeague}>Add League</Button>
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
        dataSource={data?.leagues ?? []}
        rowKey="id"
        loading={isLoading || deleteMutation.isPending}
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

export default LeagueList

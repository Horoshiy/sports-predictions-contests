import React, { useMemo, useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Table, Button, Tag, Tooltip, Space, Alert } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useSports, useDeleteSport } from '../../hooks/use-sports'
import type { Sport } from '../../types/sports.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface SportListProps {
  onCreateSport: () => void
  onEditSport: (sport: Sport) => void
}

export const SportList: React.FC<SportListProps> = ({ onCreateSport, onEditSport }) => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [pagination, setPagination] = useState({
    pageIndex: parseInt(searchParams.get('page') || '0'),
    pageSize: parseInt(searchParams.get('limit') || '10'),
  })

  useEffect(() => {
    const params = new URLSearchParams()
    params.set('page', pagination.pageIndex.toString())
    params.set('limit', pagination.pageSize.toString())
    setSearchParams(params, { replace: true })
  }, [pagination.pageIndex, pagination.pageSize, setSearchParams])

  const { data, isLoading, isError, error } = useSports({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
  })

  const deleteMutation = useDeleteSport()

  const handleDelete = (sport: Sport) => {
    if (window.confirm(`Delete "${sport.name}"?`)) {
      deleteMutation.mutate(sport.id)
    }
  }

  const columns: ColumnsType<Sport> = useMemo(() => [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Name', dataIndex: 'name', key: 'name', width: 150 },
    { title: 'Slug', dataIndex: 'slug', key: 'slug', width: 120 },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
      width: 200,
      ellipsis: true,
      render: (desc: string) => desc || '-',
    },
    {
      title: 'Status',
      dataIndex: 'isActive',
      key: 'isActive',
      width: 100,
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'success' : 'default'}>{isActive ? 'Active' : 'Inactive'}</Tag>
      ),
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
      render: (_, sport) => (
        <Space>
          <Tooltip title="Edit">
            <Button 
              type="primary" 
              icon={<EditOutlined />} 
              size="small" 
              onClick={() => onEditSport(sport)}
              aria-label={`Edit ${sport.name}`}
            />
          </Tooltip>
          <Tooltip title="Delete">
            <Button 
              danger 
              icon={<DeleteOutlined />} 
              size="small" 
              onClick={() => handleDelete(sport)} 
              loading={deleteMutation.isPending}
              aria-label={`Delete ${sport.name}`}
            />
          </Tooltip>
        </Space>
      ),
    },
  ], [deleteMutation.isPending])

  if (isError) {
    return <Alert message="Error" description={error?.message} type="error" showIcon />
  }

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Button type="primary" icon={<PlusOutlined />} onClick={onCreateSport}>
        Add Sport
      </Button>
      <Table
        columns={columns}
        dataSource={data?.sports ?? []}
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

export default SportList

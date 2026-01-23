import React, { useMemo, useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Table, Button, Tag, Tooltip, Space, Alert } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined, TeamOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useContests, useDeleteContest } from '../../hooks/use-contests'
import type { Contest } from '../../types/contest.types'
import { formatDate, formatRelativeTime } from '../../utils/date-utils'

interface ContestListProps {
  onCreateContest: () => void
  onEditContest: (contest: Contest) => void
  onViewParticipants: (contest: Contest) => void
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'draft':
      return 'default'
    case 'active':
      return 'success'
    case 'completed':
      return 'processing'
    case 'cancelled':
      return 'error'
    default:
      return 'default'
  }
}

export const ContestList: React.FC<ContestListProps> = ({
  onCreateContest,
  onEditContest,
  onViewParticipants,
}) => {
  const [searchParams, setSearchParams] = useSearchParams()
  
  const [pagination, setPagination] = useState({
    pageIndex: parseInt(searchParams.get('page') || '0'),
    pageSize: parseInt(searchParams.get('limit') || '10'),
  })

  useEffect(() => {
    const params = new URLSearchParams(searchParams)
    params.set('page', pagination.pageIndex.toString())
    params.set('limit', pagination.pageSize.toString())
    setSearchParams(params, { replace: true })
  }, [pagination, searchParams, setSearchParams])

  const { data, isLoading, isError, error } = useContests({
    pagination: {
      page: pagination.pageIndex + 1,
      limit: pagination.pageSize,
    },
  })

  const deleteContestMutation = useDeleteContest()

  const handleDeleteContest = (contest: Contest) => {
    if (window.confirm(`Are you sure you want to delete "${contest.title}"?`)) {
      deleteContestMutation.mutate(contest.id)
    }
  }

  const columns: ColumnsType<Contest> = useMemo(() => [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Title', dataIndex: 'title', key: 'title', width: 200 },
    { title: 'Sport', dataIndex: 'sportType', key: 'sportType', width: 100 },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => <Tag color={getStatusColor(status)}>{status}</Tag>,
    },
    {
      title: 'Start Date',
      dataIndex: 'startDate',
      key: 'startDate',
      width: 120,
      render: (date: string) => formatDate(date),
    },
    {
      title: 'End Date',
      dataIndex: 'endDate',
      key: 'endDate',
      width: 120,
      render: (date: string) => formatDate(date),
    },
    {
      title: 'Participants',
      dataIndex: 'participantCount',
      key: 'participantCount',
      width: 100,
      render: (count: number, contest) => `${count}${contest.maxParticipants ? `/${contest.maxParticipants}` : ''}`,
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
      render: (_, contest) => (
        <Space>
          <Tooltip title="View Participants">
            <Button icon={<TeamOutlined />} size="small" onClick={() => onViewParticipants(contest)} />
          </Tooltip>
          <Tooltip title="Edit">
            <Button type="primary" icon={<EditOutlined />} size="small" onClick={() => onEditContest(contest)} />
          </Tooltip>
          <Tooltip title="Delete">
            <Button danger icon={<DeleteOutlined />} size="small" onClick={() => handleDeleteContest(contest)} loading={deleteContestMutation.isPending} />
          </Tooltip>
        </Space>
      ),
    },
  ], [deleteContestMutation.isPending])

  if (isError) {
    return <Alert message="Error" description={error?.message} type="error" showIcon />
  }

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Button type="primary" icon={<PlusOutlined />} onClick={onCreateContest}>
        Create Contest
      </Button>
      <Table
        columns={columns}
        dataSource={data?.contests ?? []}
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

export default ContestList

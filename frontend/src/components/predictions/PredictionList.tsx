import React, { useMemo, useState } from 'react'
import { Table, Button, Tag, Tooltip, Space } from 'antd'
import { EditOutlined, DeleteOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useUserPredictions, useDeletePrediction } from '../../hooks/use-predictions'
import type { Prediction } from '../../types/prediction.types'
import { formatDate, formatRelativeTime } from '../../utils/date-utils'

interface PredictionListProps {
  contestId: number
  onEdit: (prediction: Prediction) => void
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'warning',
    scored: 'success',
    cancelled: 'error',
  }
  return colors[status] || 'default'
}

const formatPredictionData = (data: string): string => {
  try {
    const parsed = JSON.parse(data)
    const parts: string[] = []
    if (parsed.winner) parts.push(`Winner: ${parsed.winner}`)
    if (parsed.homeScore !== undefined && parsed.awayScore !== undefined) {
      parts.push(`Score: ${parsed.homeScore}-${parsed.awayScore}`)
    }
    if (parsed.overUnder) parts.push(`O/U: ${parsed.overUnder > 0 ? 'Over' : 'Under'} ${parsed.overUnderValue}`)
    return parts.join(', ') || 'N/A'
  } catch {
    return 'Invalid data'
  }
}

export const PredictionList: React.FC<PredictionListProps> = ({
  contestId,
  onEdit,
}) => {
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })

  const { data, isLoading } = useUserPredictions({
    contestId,
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
  })

  const deleteMutation = useDeletePrediction()

  const handleDelete = (prediction: Prediction) => {
    if (window.confirm('Delete this prediction?')) {
      deleteMutation.mutate(prediction.id)
    }
  }

  const columns: ColumnsType<Prediction> = useMemo(() => [
    { title: 'Event', dataIndex: 'eventTitle', key: 'eventTitle', width: 200 },
    {
      title: 'Prediction',
      dataIndex: 'predictionData',
      key: 'predictionData',
      width: 200,
      render: (data: string) => formatPredictionData(data),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => <Tag color={getStatusColor(status)}>{status}</Tag>,
    },
    {
      title: 'Points',
      dataIndex: 'pointsEarned',
      key: 'pointsEarned',
      width: 100,
      render: (points: number | null) => points !== null ? points.toFixed(1) : '-',
    },
    {
      title: 'Submitted',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 120,
      render: (date: string) => formatRelativeTime(date),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 120,
      render: (_, prediction) => (
        <Space>
          {prediction.status === 'pending' && (
            <>
              <Tooltip title="Edit">
                <Button type="primary" icon={<EditOutlined />} size="small" onClick={() => onEdit(prediction)} />
              </Tooltip>
              <Tooltip title="Delete">
                <Button danger icon={<DeleteOutlined />} size="small" onClick={() => handleDelete(prediction)} loading={deleteMutation.isPending} />
              </Tooltip>
            </>
          )}
        </Space>
      ),
    },
  ], [deleteMutation.isPending])

  return (
    <Table
      columns={columns}
      dataSource={data?.predictions ?? []}
      rowKey="id"
      loading={isLoading}
      pagination={{
        current: pagination.pageIndex + 1,
        pageSize: pagination.pageSize,
        total: data?.pagination?.total ?? 0,
        onChange: (page, pageSize) => setPagination({ pageIndex: page - 1, pageSize }),
      }}
    />
  )
}

export default PredictionList

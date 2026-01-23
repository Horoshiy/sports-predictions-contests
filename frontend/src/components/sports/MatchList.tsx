import React, { useMemo, useState } from 'react'
import { Table, Button, Tag, Tooltip, Space, Select, Alert, Typography } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useMatches, useDeleteMatch, useLeagues, useTeams } from '../../hooks/use-sports'
import type { Match, MatchStatus } from '../../types/sports.types'
import { formatDateTime } from '../../utils/date-utils'

const { Text } = Typography

interface MatchListProps {
  onCreateMatch: () => void
  onEditMatch: (match: Match) => void
}

const statusColors: Record<MatchStatus, string> = {
  scheduled: 'default',
  live: 'warning',
  finished: 'success',
  cancelled: 'error',
  postponed: 'processing',
}

export const MatchList: React.FC<MatchListProps> = ({ onCreateMatch, onEditMatch }) => {
  const [leagueFilter, setLeagueFilter] = useState<number | ''>('')
  const [statusFilter, setStatusFilter] = useState<MatchStatus | ''>('')
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })

  const { data: leaguesData } = useLeagues({ pagination: { page: 1, limit: 100 } })
  const { data: teamsData } = useTeams({ pagination: { page: 1, limit: 200 } })

  const leaguesMap = useMemo(() => new Map(leaguesData?.leagues?.map(l => [l.id, l.name]) || []), [leaguesData])
  const teamsMap = useMemo(() => new Map(teamsData?.teams?.map(t => [t.id, t.name]) || []), [teamsData])

  const { data, isLoading, isError, error } = useMatches({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    leagueId: leagueFilter || undefined,
    status: statusFilter || undefined,
  })

  const deleteMutation = useDeleteMatch()

  const handleDelete = (match: Match) => {
    if (window.confirm('Delete this match?')) {
      deleteMutation.mutate(match.id)
    }
  }

  const columns: ColumnsType<Match> = useMemo(() => [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    {
      title: 'Match',
      key: 'matchup',
      width: 250,
      render: (_, match) => (
        <Text>{teamsMap.get(match.homeTeamId) || 'TBD'} vs {teamsMap.get(match.awayTeamId) || 'TBD'}</Text>
      ),
    },
    {
      title: 'League',
      dataIndex: 'leagueId',
      key: 'leagueId',
      width: 150,
      render: (leagueId: number) => leaguesMap.get(leagueId) || '-',
    },
    {
      title: 'Scheduled',
      dataIndex: 'scheduledAt',
      key: 'scheduledAt',
      width: 150,
      render: (scheduledAt: string) => formatDateTime(scheduledAt),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: MatchStatus) => <Tag color={statusColors[status]}>{status}</Tag>,
    },
    {
      title: 'Score',
      key: 'score',
      width: 80,
      render: (_, match) => match.status === 'finished' || match.status === 'live'
        ? `${match.homeScore} - ${match.awayScore}`
        : '-',
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 120,
      render: (_, match) => (
        <Space>
          <Tooltip title="Edit">
            <Button type="primary" icon={<EditOutlined />} size="small" onClick={() => onEditMatch(match)} />
          </Tooltip>
          <Tooltip title="Delete">
            <Button danger icon={<DeleteOutlined />} size="small" onClick={() => handleDelete(match)} loading={deleteMutation.isPending} />
          </Tooltip>
        </Space>
      ),
    },
  ], [leaguesMap, teamsMap, deleteMutation.isPending])

  if (isError) {
    return <Alert message="Error" description={error?.message} type="error" showIcon />
  }

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Space wrap>
        <Button type="primary" icon={<PlusOutlined />} onClick={onCreateMatch}>Add Match</Button>
        <Select
          style={{ minWidth: 150 }}
          placeholder="League"
          value={leagueFilter}
          onChange={setLeagueFilter}
          allowClear
        >
          <Select.Option value="">All Leagues</Select.Option>
          {leaguesData?.leagues?.map(l => <Select.Option key={l.id} value={l.id}>{l.name}</Select.Option>)}
        </Select>
        <Select
          style={{ minWidth: 120 }}
          placeholder="Status"
          value={statusFilter}
          onChange={setStatusFilter}
          allowClear
        >
          <Select.Option value="">All</Select.Option>
          <Select.Option value="scheduled">Scheduled</Select.Option>
          <Select.Option value="live">Live</Select.Option>
          <Select.Option value="finished">Finished</Select.Option>
          <Select.Option value="cancelled">Cancelled</Select.Option>
          <Select.Option value="postponed">Postponed</Select.Option>
        </Select>
      </Space>
      <Table
        columns={columns}
        dataSource={data?.matches ?? []}
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

export default MatchList

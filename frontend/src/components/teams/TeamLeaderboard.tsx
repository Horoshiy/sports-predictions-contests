import React from 'react'
import { Table, Typography, Tag, Space } from 'antd'
import { TrophyOutlined, TeamOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useTeamLeaderboard } from '../../hooks/use-teams'

const { Text } = Typography

interface TeamLeaderboardProps {
  contestId: number
  userTeamId?: number
}

const getRankColor = (rank: number) => {
  if (rank === 1) return 'gold'
  if (rank === 2) return 'default'
  if (rank === 3) return 'orange'
  return undefined
}

interface LeaderboardEntry {
  rank: number
  teamId: number
  teamName: string
  totalPoints: number
  memberCount: number
}

export const TeamLeaderboard: React.FC<TeamLeaderboardProps> = ({ contestId, userTeamId }) => {
  const { data: entries, isLoading, isError } = useTeamLeaderboard(contestId, 20)

  const columns: ColumnsType<LeaderboardEntry> = [
    {
      title: 'Rank',
      dataIndex: 'rank',
      key: 'rank',
      width: 80,
      render: (rank: number) => (
        <Space>
          {rank <= 3 && <TrophyOutlined style={{ color: rank === 1 ? '#FFD700' : rank === 2 ? '#C0C0C0' : '#CD7F32' }} />}
          <Text strong={rank <= 3}>#{rank}</Text>
        </Space>
      ),
    },
    {
      title: 'Team',
      dataIndex: 'teamName',
      key: 'teamName',
      render: (name: string, record) => (
        <Space>
          <TeamOutlined />
          <Text strong={record.teamId === userTeamId}>{name}</Text>
          {record.teamId === userTeamId && <Tag color="blue">Your Team</Tag>}
        </Space>
      ),
    },
    {
      title: 'Members',
      dataIndex: 'memberCount',
      key: 'memberCount',
      width: 100,
    },
    {
      title: 'Points',
      dataIndex: 'totalPoints',
      key: 'totalPoints',
      width: 120,
      render: (points: number) => <Text strong>{points.toFixed(1)}</Text>,
    },
  ]

  if (isLoading) return <Text>Loading team leaderboard...</Text>
  if (isError) return <Text type="danger">Failed to load team leaderboard</Text>
  if (!entries || entries.length === 0) return <Text type="secondary">No teams in this contest yet</Text>

  return (
    <Table
      columns={columns}
      dataSource={entries}
      rowKey="teamId"
      size="small"
      pagination={false}
      rowClassName={(record) => record.teamId === userTeamId ? 'ant-table-row-selected' : ''}
    />
  )
}

export default TeamLeaderboard

import React, { useState, useEffect } from 'react'
import { Table, Card, Typography, Tag, Avatar, Button, Tooltip, Space, Spin, Alert, Statistic, Row, Col } from 'antd'
import { ReloadOutlined, TrophyOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useQuery } from '@tanstack/react-query'
import scoringService from '../../services/scoring-service'
import type { LeaderboardEntry } from '../../types/scoring.types'
import { formatRelativeTime } from '../../utils/date-utils'
import { DEFAULT_REFRESH_INTERVAL, MAX_PARTICIPANTS_DISPLAY } from '../../utils/constants'

const { Text, Title } = Typography

interface LeaderboardTableProps {
  contestId: number
  currentUserId?: number
  limit?: number
  autoRefresh?: boolean
  refreshInterval?: number
}

const getRankColor = (rank: number) => {
  switch (rank) {
    case 1:
      return 'gold'
    case 2:
      return 'default'
    case 3:
      return 'orange'
    default:
      return 'default'
  }
}

const getRankIcon = (rank: number) => {
  if (rank <= 3) {
    const colors = { 1: '#FFD700', 2: '#C0C0C0', 3: '#CD7F32' }
    return <TrophyOutlined style={{ color: colors[rank as 1 | 2 | 3], fontSize: 20 }} />
  }
  return null
}

export const LeaderboardTable: React.FC<LeaderboardTableProps> = ({
  contestId,
  currentUserId,
  limit = MAX_PARTICIPANTS_DISPLAY,
  autoRefresh = true,
  refreshInterval = DEFAULT_REFRESH_INTERVAL,
}) => {
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date())

  // Query for leaderboard data
  const {
    data: leaderboard,
    isLoading,
    error,
    refetch,
    isRefetching,
  } = useQuery({
    queryKey: ['leaderboard', contestId, limit],
    queryFn: () => scoringService.getLeaderboard({ contestId, limit }),
    refetchInterval: autoRefresh ? refreshInterval : false,
    refetchIntervalInBackground: true,
  })

  // Update lastUpdated when data changes
  useEffect(() => {
    if (leaderboard) {
      setLastUpdated(new Date())
    }
  }, [leaderboard])

  // Query for current user's rank if provided
  const {
    data: userRank,
    isLoading: isLoadingUserRank,
  } = useQuery({
    queryKey: ['userRank', contestId, currentUserId],
    queryFn: () => currentUserId 
      ? scoringService.getUserRank({ contestId, userId: currentUserId })
      : null,
    enabled: !!currentUserId,
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  // Query for current user's streak if provided
  const {
    data: userStreak,
    error: userStreakError,
  } = useQuery({
    queryKey: ['userStreak', contestId, currentUserId],
    queryFn: () => currentUserId 
      ? scoringService.getUserStreak({ contestId, userId: currentUserId })
      : null,
    enabled: !!currentUserId,
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  // Log streak error if any (silent fail for UI)
  if (userStreakError) {
    console.error('Failed to fetch user streak:', userStreakError)
  }

  // Define table columns
  const columns: ColumnsType<LeaderboardEntry> = [
    {
      title: 'Rank',
      dataIndex: 'rank',
      key: 'rank',
      width: 80,
      render: (rank: number, record) => {
        const isCurrentUser = currentUserId === record.userId
        return (
          <Space>
            {getRankIcon(rank)}
            <Text strong={rank <= 3 || isCurrentUser} style={{ color: isCurrentUser ? '#1890ff' : undefined }}>
              #{rank}
            </Text>
          </Space>
        )
      },
    },
    {
      title: 'Player',
      dataIndex: 'userName',
      key: 'userName',
      width: 200,
      render: (userName: string, record) => {
        const displayName = userName || `User ${record.userId}`
        const isCurrentUser = currentUserId === record.userId
        return (
          <Space>
            <Avatar size={32}>{displayName.charAt(0).toUpperCase()}</Avatar>
            <Space direction="vertical" size={0}>
              <Text strong={isCurrentUser} style={{ color: isCurrentUser ? '#1890ff' : undefined }}>
                {displayName}
                {isCurrentUser && <Tag color="blue" style={{ marginLeft: 8 }}>You</Tag>}
              </Text>
            </Space>
          </Space>
        )
      },
    },
    {
      title: 'Points',
      dataIndex: 'totalPoints',
      key: 'totalPoints',
      width: 120,
      render: (points: number) => (
        <Text strong style={{ color: points > 0 ? '#52c41a' : undefined }}>
          {points.toFixed(1)}
        </Text>
      ),
    },
    {
      title: 'Streak',
      dataIndex: 'currentStreak',
      key: 'currentStreak',
      width: 100,
      render: (streak: number, record) => {
        const multiplier = record.multiplier ?? 1
        return (
          <Space>
            <Text>🔥 {streak ?? 0}</Text>
            {multiplier > 1 && <Tag color="warning">{multiplier}x</Tag>}
          </Space>
        )
      },
    },
    {
      title: 'Last Updated',
      dataIndex: 'updatedAt',
      key: 'updatedAt',
      width: 150,
      render: (updatedAt: string) => (
        <Text type="secondary">{formatRelativeTime(updatedAt)}</Text>
      ),
    },
  ]

  const handleRefresh = () => {
    refetch()
    setLastUpdated(new Date())
  }

  if (error) {
    return (
      <Alert
        message="Failed to load leaderboard"
        description="Please try again."
        type="error"
        showIcon
        action={
          <Button size="small" onClick={handleRefresh}>
            Retry
          </Button>
        }
        style={{ marginBottom: 16 }}
      />
    )
  }

  return (
    <Card>
      {/* Header */}
      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Title level={4} style={{ margin: 0 }}>Leaderboard</Title>
          <Space>
            <Text type="secondary" style={{ fontSize: 12 }}>
              Last updated: {formatRelativeTime(lastUpdated.toISOString())}
            </Text>
            <Tooltip title="Refresh leaderboard">
              <Button
                icon={isRefetching ? <Spin size="small" /> : <ReloadOutlined />}
                onClick={handleRefresh}
                disabled={isRefetching}
                size="small"
              />
            </Tooltip>
          </Space>
        </div>

        {/* Current User Rank Card */}
        {currentUserId && userRank && !isLoadingUserRank && (
          <Card size="small" style={{ backgroundColor: '#e6f7ff', borderColor: '#1890ff' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Text strong style={{ color: '#1890ff' }}>Your Position</Text>
              <Row gutter={32}>
                <Col>
                  <Statistic
                    title="Rank"
                    value={userRank.rank}
                    prefix="#"
                    valueStyle={{ color: '#1890ff', fontSize: 20 }}
                  />
                </Col>
                <Col>
                  <Statistic
                    title="Points"
                    value={userRank.totalPoints}
                    precision={1}
                    valueStyle={{ color: '#1890ff', fontSize: 20 }}
                  />
                </Col>
                {userStreak && (
                  <Col>
                    <Statistic
                      title="Streak"
                      value={userStreak.currentStreak}
                      prefix="🔥"
                      suffix={userStreak.multiplier > 1 ? <Tag color="warning">{userStreak.multiplier}x</Tag> : undefined}
                      valueStyle={{ color: '#fa8c16', fontSize: 20 }}
                    />
                  </Col>
                )}
              </Row>
            </div>
          </Card>
        )}

        {/* Loading State */}
        {isLoading && (
          <div style={{ textAlign: 'center', padding: '32px 0' }}>
            <Spin size="large" />
          </div>
        )}

        {/* Empty State */}
        {!isLoading && (!leaderboard?.entries || leaderboard.entries.length === 0) && (
          <div style={{ textAlign: 'center', padding: '32px 0' }}>
            <TrophyOutlined style={{ fontSize: 48, color: '#bfbfbf', marginBottom: 16 }} />
            <Title level={5} type="secondary">No rankings yet</Title>
            <Text type="secondary">Scores will appear here once predictions are evaluated.</Text>
          </div>
        )}

        {/* Leaderboard Table */}
        {!isLoading && leaderboard?.entries && leaderboard.entries.length > 0 && (
          <Table
            columns={columns}
            dataSource={leaderboard.entries}
            rowKey="userId"
            pagination={false}
            rowClassName={(record) => currentUserId === record.userId ? 'ant-table-row-selected' : ''}
          />
        )}

        {/* Auto-refresh indicator */}
        {autoRefresh && (
          <div style={{ textAlign: 'center' }}>
            <Tag color="blue">Auto-refreshing every {refreshInterval / 1000}s</Tag>
          </div>
        )}
      </Space>
    </Card>
  )
}

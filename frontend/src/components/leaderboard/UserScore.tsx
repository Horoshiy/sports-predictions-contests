import React from 'react'
import { Card, Statistic, Row, Col, Avatar, Button, Tooltip, Progress, Tag, Space, Divider, Typography } from 'antd'
import { TrophyOutlined, RiseOutlined, FallOutlined, ReloadOutlined, LineChartOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import scoringService from '../../services/scoring-service'
import type { Score } from '../../types/scoring.types'
import { formatRelativeTime } from '../../utils/date-utils'

const { Text } = Typography

interface UserScoreProps {
  userId: number
  contestId: number
  userName?: string
  showDetails?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
  onRefresh?: () => void
}

interface ScoreBreakdown {
  totalPredictions: number
  scoredPredictions: number
  averagePoints: number
  bestScore: number
  recentScores: Score[]
}

const getRankColor = (rank: number) => {
  if (rank === 1) return 'gold'
  if (rank === 2) return 'silver'
  if (rank === 3) return 'orange'
  if (rank <= 10) return 'blue'
  return 'default'
}

export const UserScore: React.FC<UserScoreProps> = ({
  userId,
  contestId,
  userName,
  showDetails = true,
  autoRefresh = true,
  refreshInterval = 30000,
  onRefresh,
}) => {
  const {
    data: userScores,
    isLoading: isLoadingScores,
    error: scoresError,
    refetch: refetchScores,
  } = useQuery({
    queryKey: ['userScores', userId, contestId],
    queryFn: () => scoringService.getUserScores({ userId, contestId }),
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  const {
    data: userRank,
    isLoading: isLoadingRank,
    error: rankError,
    refetch: refetchRank,
  } = useQuery({
    queryKey: ['userRank', userId, contestId],
    queryFn: () => scoringService.getUserRank({ contestId, userId }),
    refetchInterval: autoRefresh ? refreshInterval : false,
  })

  const scoreBreakdown: ScoreBreakdown | null = React.useMemo(() => {
    if (!userScores?.scores) return null

    const scores = userScores.scores
    const totalPredictions = scores.length
    const scoredPredictions = scores.filter(s => s.points > 0).length
    const averagePoints = totalPredictions > 0 ? userScores.totalPoints / totalPredictions : 0
    const bestScore = Math.max(...scores.map(s => s.points), 0)
    const recentScores = scores
      .sort((a, b) => new Date(b.scoredAt).getTime() - new Date(a.scoredAt).getTime())
      .slice(0, 5)

    return { totalPredictions, scoredPredictions, averagePoints, bestScore, recentScores }
  }, [userScores])

  const handleRefresh = () => {
    refetchScores()
    refetchRank()
    onRefresh?.()
  }

  const isLoading = isLoadingScores || isLoadingRank
  const hasError = scoresError || rankError

  if (hasError) {
    return (
      <Card>
        <Text type="danger">Failed to load user score data</Text>
      </Card>
    )
  }

  return (
    <Card
      loading={isLoading}
      title={
        <Space>
          <Avatar>{userName ? userName.charAt(0).toUpperCase() : 'U'}</Avatar>
          <div>
            <div>{userName || `User ${userId}`}</div>
            <Text type="secondary" style={{ fontSize: '12px' }}>Contest Performance</Text>
          </div>
        </Space>
      }
      extra={
        <Tooltip title="Refresh scores">
          <Button type="text" icon={<ReloadOutlined />} onClick={handleRefresh} />
        </Tooltip>
      }
    >
      {!isLoading && userRank && (
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Row gutter={16}>
            <Col span={12}>
              <Statistic
                title="Current Rank"
                value={userRank.rank}
                prefix={<TrophyOutlined style={{ color: getRankColor(userRank.rank) === 'gold' ? '#FFD700' : undefined }} />}
                valueStyle={{ color: getRankColor(userRank.rank) === 'gold' ? '#FFD700' : undefined }}
              />
            </Col>
            <Col span={12}>
              <Statistic
                title="Total Points"
                value={userRank.totalPoints}
                precision={1}
                valueStyle={{ color: '#1976d2' }}
              />
            </Col>
          </Row>

          {showDetails && scoreBreakdown && (
            <>
              <Divider />
              <Row gutter={16}>
                <Col span={6}>
                  <Statistic title="Predictions" value={scoreBreakdown.totalPredictions} />
                </Col>
                <Col span={6}>
                  <Statistic title="Scored" value={scoreBreakdown.scoredPredictions} valueStyle={{ color: '#52c41a' }} />
                </Col>
                <Col span={6}>
                  <Statistic title="Avg Points" value={scoreBreakdown.averagePoints} precision={1} valueStyle={{ color: '#1890ff' }} />
                </Col>
                <Col span={6}>
                  <Statistic title="Best Score" value={scoreBreakdown.bestScore} precision={1} valueStyle={{ color: '#faad14' }} />
                </Col>
              </Row>

              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <Text>Success Rate</Text>
                  <Text type="secondary">
                    {scoreBreakdown.totalPredictions > 0 
                      ? Math.round((scoreBreakdown.scoredPredictions / scoreBreakdown.totalPredictions) * 100)
                      : 0}%
                  </Text>
                </div>
                <Progress 
                  percent={scoreBreakdown.totalPredictions > 0 
                    ? Math.round((scoreBreakdown.scoredPredictions / scoreBreakdown.totalPredictions) * 100)
                    : 0}
                  strokeColor="#52c41a"
                />
              </div>

              {scoreBreakdown.recentScores.length > 0 && (
                <div>
                  <Text strong>Recent Scores</Text>
                  <div style={{ marginTop: 8 }}>
                    <Space wrap>
                      {scoreBreakdown.recentScores.map((score) => (
                        <Tooltip key={score.id} title={`Prediction ${score.predictionId} - ${formatRelativeTime(score.scoredAt)}`}>
                          <Tag color={score.points > 0 ? 'success' : 'default'}>
                            {score.points.toFixed(1)}
                          </Tag>
                        </Tooltip>
                      ))}
                    </Space>
                  </div>
                </div>
              )}
            </>
          )}
        </Space>
      )}

      {!isLoading && !userRank && (
        <div style={{ textAlign: 'center', padding: '24px 0' }}>
          <LineChartOutlined style={{ fontSize: 48, color: '#8c8c8c', marginBottom: 8 }} />
          <div>
            <Text type="secondary">No scores yet. Make some predictions to see your performance!</Text>
          </div>
        </div>
      )}
    </Card>
  )
}

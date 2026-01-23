import React, { useState } from 'react'
import { Space, Typography, Card, Row, Col, Segmented, Spin, Alert, Statistic } from 'antd'
import { LineChartOutlined, TrophyOutlined, RiseOutlined } from '@ant-design/icons'
import { useAuth } from '../contexts/AuthContext'
import { useUserAnalytics } from '../hooks/use-analytics'
import { AccuracyChart } from '../components/analytics/AccuracyChart'
import { SportBreakdown } from '../components/analytics/SportBreakdown'
import { PlatformComparison } from '../components/analytics/PlatformComparison'
import { ExportButton } from '../components/analytics/ExportButton'
import type { TimeRange } from '../types/analytics.types'

const { Title } = Typography

export const AnalyticsPage: React.FC = () => {
  const { user } = useAuth()
  const [timeRange, setTimeRange] = useState<TimeRange>('30d')

  const {
    data: analytics,
    isLoading,
    error,
  } = useUserAnalytics(user?.id || 0, timeRange)

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: '64px 0' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (error) {
    return <Alert message="Error" description="Failed to load analytics data" type="error" showIcon />
  }

  if (!analytics) {
    return <Alert message="No data available" type="info" showIcon />
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%', padding: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={2}>Analytics Dashboard</Title>
        <Space>
          <Segmented
            options={[
              { label: '7 Days', value: '7d' },
              { label: '30 Days', value: '30d' },
              { label: '90 Days', value: '90d' },
              { label: 'All Time', value: 'all' },
            ]}
            value={timeRange}
            onChange={(value) => setTimeRange(value as TimeRange)}
          />
          <ExportButton userId={user?.id || 0} timeRange={timeRange} />
        </Space>
      </div>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Accuracy"
              value={analytics.overallAccuracy}
              precision={1}
              suffix="%"
              prefix={<LineChartOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Points"
              value={analytics.totalPoints}
              precision={1}
              prefix={<TrophyOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Predictions"
              value={analytics.totalPredictions}
              prefix={<RiseOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Correct"
              value={analytics.correctPredictions}
              suffix={`/ ${analytics.totalPredictions}`}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <AccuracyChart trends={analytics.trends} />
        </Col>
        <Col xs={24} lg={12}>
          <SportBreakdown bySport={analytics.bySport} />
        </Col>
      </Row>

      {analytics.platformComparison && (
        <Row gutter={[16, 16]}>
          <Col xs={24}>
            <PlatformComparison
              comparison={{
                userAccuracy: analytics.overallAccuracy,
                platformAverage: analytics.platformComparison.averageAccuracy,
              }}
            />
          </Col>
        </Row>
      )}
    </Space>
  )
}

export default AnalyticsPage

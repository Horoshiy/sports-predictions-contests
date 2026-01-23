import React from 'react'
import { Card, Typography, Space, Statistic, Row, Col } from 'antd'
import { ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons'

const { Title, Text } = Typography

interface PlatformComparisonType {
  userAccuracy: number
  platformAverage: number
}

interface PlatformComparisonProps {
  comparison: PlatformComparisonType
  title?: string
}

export const PlatformComparison: React.FC<PlatformComparisonProps> = ({
  comparison,
  title = 'Platform Comparison',
}) => {
  if (!comparison) {
    return (
      <Card>
        <Title level={5}>{title}</Title>
        <div style={{ textAlign: 'center', padding: '32px 0' }}>
          <Text type="secondary">No comparison data available</Text>
        </div>
      </Card>
    )
  }

  const diff = comparison.userAccuracy - comparison.platformAverage
  const isAboveAverage = diff > 0

  return (
    <Card>
      <Title level={5}>{title}</Title>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <Row gutter={16}>
          <Col span={12}>
            <Statistic
              title="Your Accuracy"
              value={comparison.userAccuracy}
              precision={1}
              suffix="%"
              valueStyle={{ color: '#1890ff' }}
            />
          </Col>
          <Col span={12}>
            <Statistic
              title="Platform Average"
              value={comparison.platformAverage}
              precision={1}
              suffix="%"
              valueStyle={{ color: '#8c8c8c' }}
            />
          </Col>
        </Row>
        <Statistic
          title="Difference"
          value={Math.abs(diff)}
          precision={1}
          suffix="%"
          prefix={isAboveAverage ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
          valueStyle={{ color: isAboveAverage ? '#3f8600' : '#cf1322' }}
        />
        <Text type="secondary">
          {isAboveAverage
            ? `You're performing ${diff.toFixed(1)}% better than the platform average!`
            : `You're ${Math.abs(diff).toFixed(1)}% below the platform average.`}
        </Text>
      </Space>
    </Card>
  )
}

export default PlatformComparison

import React from 'react'
import { Tag, Tooltip, Space, Typography } from 'antd'
import { RiseOutlined, ClockCircleOutlined } from '@ant-design/icons'

const { Text } = Typography

interface CoefficientIndicatorProps {
  coefficient: number
  tier: string
  hoursUntilEvent: number
  compact?: boolean
}

const getTagColor = (coefficient: number): string => {
  if (coefficient >= 2.0) return 'success'
  if (coefficient >= 1.5) return 'blue'
  if (coefficient >= 1.25) return 'warning'
  return 'default'
}

export const CoefficientIndicator: React.FC<CoefficientIndicatorProps> = ({
  coefficient,
  tier,
  hoursUntilEvent,
  compact = false,
}) => {
  const formatTimeRemaining = (hours: number): string => {
    if (hours >= 168) return `${Math.floor(hours / 24)} days`
    if (hours >= 24) return `${Math.floor(hours / 24)}d ${Math.floor(hours % 24)}h`
    return `${Math.floor(hours)}h ${Math.round((hours % 1) * 60)}m`
  }

  if (compact) {
    return (
      <Tooltip title={`${tier} - ${coefficient}x points (${formatTimeRemaining(hoursUntilEvent)} left)`}>
        <Tag icon={<RiseOutlined />} color={getTagColor(coefficient)}>
          {coefficient}x
        </Tag>
      </Tooltip>
    )
  }

  return (
    <div style={{ padding: 12, backgroundColor: '#f5f5f5', borderRadius: 4 }}>
      <Space direction="vertical" size={4} style={{ width: '100%' }}>
        <Space>
          <RiseOutlined style={{ color: getTagColor(coefficient) === 'success' ? '#52c41a' : undefined }} />
          <Text strong>{coefficient}x Points</Text>
          <Tag color={getTagColor(coefficient)}>{tier}</Tag>
        </Space>
        <Space size={4}>
          <ClockCircleOutlined style={{ fontSize: 12 }} />
          <Text type="secondary" style={{ fontSize: 12 }}>
            {formatTimeRemaining(hoursUntilEvent)} until coefficient drops
          </Text>
        </Space>
      </Space>
    </div>
  )
}

export default CoefficientIndicator

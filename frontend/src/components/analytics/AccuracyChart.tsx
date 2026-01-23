import React from 'react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { Card, Typography } from 'antd'
import type { AccuracyTrend } from '../../types/analytics.types'

const { Title, Text } = Typography

interface AccuracyChartProps {
  trends: AccuracyTrend[]
  title?: string
}

export const AccuracyChart: React.FC<AccuracyChartProps> = ({
  trends,
  title = 'Accuracy Over Time',
}) => {
  if (!trends || trends.length === 0) {
    return (
      <Card>
        <Title level={5}>{title}</Title>
        <div style={{ textAlign: 'center', padding: '32px 0' }}>
          <Text type="secondary">No data available yet</Text>
        </div>
      </Card>
    )
  }

  return (
    <Card>
      <Title level={5}>{title}</Title>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={trends}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="date" />
          <YAxis domain={[0, 100]} />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="accuracy" stroke="#1890ff" name="Accuracy %" />
        </LineChart>
      </ResponsiveContainer>
    </Card>
  )
}

export default AccuracyChart

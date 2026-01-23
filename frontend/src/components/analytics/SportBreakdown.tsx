import React from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, Cell } from 'recharts'
import { Card, Typography } from 'antd'
import type { SportAccuracy } from '../../types/analytics.types'

const { Title, Text } = Typography

interface SportBreakdownProps {
  bySport: SportAccuracy[]
  title?: string
}

const COLORS = ['#1976d2', '#2e7d32', '#ed6c02', '#9c27b0', '#d32f2f', '#0288d1']

export const SportBreakdown: React.FC<SportBreakdownProps> = ({
  bySport,
  title = 'Performance by Sport',
}) => {
  if (!bySport || bySport.length === 0) {
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
        <BarChart data={bySport}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="sport" />
          <YAxis domain={[0, 100]} />
          <Tooltip />
          <Legend />
          <Bar dataKey="accuracy" name="Accuracy %">
            {bySport.map((_, index) => (
              <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </Card>
  )
}

export default SportBreakdown

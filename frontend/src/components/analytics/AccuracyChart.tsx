import React from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { Card, CardContent, Typography, Box } from '@mui/material'
import type { AccuracyTrend } from '../../types/analytics.types'

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
        <CardContent>
          <Typography variant="h6" gutterBottom>{title}</Typography>
          <Box textAlign="center" py={4}>
            <Typography color="text.secondary">
              No trend data available yet
            </Typography>
          </Box>
        </CardContent>
      </Card>
    )
  }

  const data = trends.map((t) => ({
    period: t.period,
    accuracy: Number(t.accuracyPercentage.toFixed(1)),
    predictions: t.totalPredictions,
    points: Number(t.totalPoints.toFixed(1)),
  }))

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>{title}</Typography>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="period" tick={{ fontSize: 12 }} />
            <YAxis
              yAxisId="left"
              domain={[0, 100]}
              tick={{ fontSize: 12 }}
              label={{ value: 'Accuracy %', angle: -90, position: 'insideLeft' }}
            />
            <YAxis
              yAxisId="right"
              orientation="right"
              tick={{ fontSize: 12 }}
              label={{ value: 'Points', angle: 90, position: 'insideRight' }}
            />
            <Tooltip
              formatter={(value: number, name: string) => {
                if (name === 'Accuracy %') return [`${value}%`, 'Accuracy']
                if (name === 'Points') return [value, 'Points']
                return [value, name]
              }}
            />
            <Legend />
            <Line
              yAxisId="left"
              type="monotone"
              dataKey="accuracy"
              stroke="#1976d2"
              strokeWidth={2}
              dot={{ r: 4 }}
              name="Accuracy %"
            />
            <Line
              yAxisId="right"
              type="monotone"
              dataKey="points"
              stroke="#2e7d32"
              strokeWidth={2}
              dot={{ r: 4 }}
              name="Points"
            />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}

export default AccuracyChart

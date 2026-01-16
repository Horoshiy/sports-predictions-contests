import React from 'react'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Cell,
} from 'recharts'
import { Card, CardContent, Typography, Box } from '@mui/material'
import type { SportAccuracy } from '../../types/analytics.types'

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
        <CardContent>
          <Typography variant="h6" gutterBottom>{title}</Typography>
          <Box textAlign="center" py={4}>
            <Typography color="text.secondary">
              No sport data available yet
            </Typography>
          </Box>
        </CardContent>
      </Card>
    )
  }

  const data = bySport.map((s) => ({
    sport: s.sportType,
    accuracy: Number(s.accuracyPercentage.toFixed(1)),
    predictions: s.totalPredictions,
    points: Number(s.totalPoints.toFixed(1)),
  }))

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>{title}</Typography>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="sport" tick={{ fontSize: 12 }} />
            <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
            <Tooltip
              formatter={(value: number, name: string) => {
                if (name === 'Accuracy %') return [`${value}%`, 'Accuracy']
                return [value, name]
              }}
            />
            <Legend />
            <Bar dataKey="accuracy" name="Accuracy %" radius={[4, 4, 0, 0]}>
              {data.map((_, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
        
        <Box mt={2}>
          <Typography variant="subtitle2" gutterBottom>Details</Typography>
          <Box component="table" sx={{ width: '100%', fontSize: 14 }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left' }}>Sport</th>
                <th style={{ textAlign: 'right' }}>Predictions</th>
                <th style={{ textAlign: 'right' }}>Correct</th>
                <th style={{ textAlign: 'right' }}>Points</th>
              </tr>
            </thead>
            <tbody>
              {bySport.map((s) => (
                <tr key={s.sportType}>
                  <td>{s.sportType}</td>
                  <td style={{ textAlign: 'right' }}>{s.totalPredictions}</td>
                  <td style={{ textAlign: 'right' }}>{s.correctPredictions}</td>
                  <td style={{ textAlign: 'right' }}>{s.totalPoints.toFixed(1)}</td>
                </tr>
              ))}
            </tbody>
          </Box>
        </Box>
      </CardContent>
    </Card>
  )
}

export default SportBreakdown

import React, { useState } from 'react'
import { Button, CircularProgress } from '@mui/material'
import { Download as DownloadIcon } from '@mui/icons-material'
import { useExportAnalytics } from '../../hooks/use-analytics'
import type { TimeRange } from '../../types/analytics.types'

interface ExportButtonProps {
  userId: number
  timeRange: TimeRange
}

export const ExportButton: React.FC<ExportButtonProps> = ({ userId, timeRange }) => {
  const [isExporting, setIsExporting] = useState(false)
  const { exportAnalytics } = useExportAnalytics()

  const handleExport = async () => {
    setIsExporting(true)
    try {
      await exportAnalytics(userId, timeRange)
    } finally {
      setIsExporting(false)
    }
  }

  return (
    <Button
      variant="outlined"
      startIcon={isExporting ? <CircularProgress size={20} /> : <DownloadIcon />}
      onClick={handleExport}
      disabled={isExporting}
    >
      {isExporting ? 'Exporting...' : 'Export CSV'}
    </Button>
  )
}

export default ExportButton

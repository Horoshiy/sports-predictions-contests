import React, { useState } from 'react'
import { Button } from 'antd'
import { DownloadOutlined } from '@ant-design/icons'
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
      icon={<DownloadOutlined />}
      onClick={handleExport}
      loading={isExporting}
    >
      {isExporting ? 'Exporting...' : 'Export CSV'}
    </Button>
  )
}

export default ExportButton

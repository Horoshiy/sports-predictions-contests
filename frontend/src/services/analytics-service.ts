import grpcClient from './grpc-client'
import type {
  UserAnalytics,
  GetUserAnalyticsResponse,
  ExportAnalyticsResponse,
  TimeRange,
} from '../types/analytics.types'

class AnalyticsService {
  private basePath = '/v1/users'

  async getUserAnalytics(
    userId: number,
    timeRange: TimeRange = '30d'
  ): Promise<UserAnalytics> {
    const params = new URLSearchParams()
    params.append('time_range', timeRange)

    const response = await grpcClient.get<GetUserAnalyticsResponse>(
      `${this.basePath}/${userId}/analytics?${params.toString()}`
    )
    
    if (!response.response?.success) {
      throw new Error(response.response?.message || 'Failed to fetch analytics')
    }
    
    return response.analytics
  }

  async exportAnalytics(
    userId: number,
    timeRange: TimeRange = '30d'
  ): Promise<{ data: string; filename: string }> {
    const params = new URLSearchParams()
    params.append('time_range', timeRange)
    params.append('format', 'csv')

    const response = await grpcClient.get<ExportAnalyticsResponse>(
      `${this.basePath}/${userId}/analytics/export?${params.toString()}`
    )

    if (!response.response?.success) {
      throw new Error(response.response?.message || 'Failed to export analytics')
    }

    return {
      data: response.data,
      filename: response.filename,
    }
  }

  downloadCSV(data: string, filename: string): void {
    const blob = new Blob([data], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', filename)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }
}

export const analyticsService = new AnalyticsService()
export default analyticsService

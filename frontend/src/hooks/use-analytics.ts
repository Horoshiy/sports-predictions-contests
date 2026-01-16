import { useQuery } from '@tanstack/react-query'
import analyticsService from '../services/analytics-service'
import { useToast } from '../contexts/ToastContext'
import type { TimeRange } from '../types/analytics.types'

export const analyticsKeys = {
  all: ['analytics'] as const,
  user: (userId: number, timeRange: TimeRange) =>
    [...analyticsKeys.all, 'user', userId, timeRange] as const,
}

export const useUserAnalytics = (userId: number, timeRange: TimeRange = '30d') => {
  return useQuery({
    queryKey: analyticsKeys.user(userId, timeRange),
    queryFn: () => analyticsService.getUserAnalytics(userId, timeRange),
    enabled: !!userId,
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
  })
}

export const useExportAnalytics = () => {
  const { showToast } = useToast()

  const exportAnalytics = async (userId: number, timeRange: TimeRange = '30d') => {
    try {
      const { data, filename } = await analyticsService.exportAnalytics(userId, timeRange)
      analyticsService.downloadCSV(data, filename)
      showToast('Analytics exported successfully!', 'success')
    } catch (error) {
      showToast('Failed to export analytics', 'error')
      throw error
    }
  }

  return { exportAnalytics }
}

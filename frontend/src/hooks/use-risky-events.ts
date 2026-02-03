import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { riskyEventsService } from '../services/risky-events-service'
import { useToast } from '../contexts/ToastContext'
import type {
  RiskyEventType,
  ListRiskyEventTypesRequest,
  CreateRiskyEventTypeRequest,
  UpdateRiskyEventTypeRequest,
  SetMatchRiskyEventsRequest,
  SetMatchRiskyEventOutcomeRequest,
} from '../types/risky-events.types'

// ==================== Query Keys ====================

export const riskyEventTypesKeys = {
  all: ['riskyEventTypes'] as const,
  lists: () => [...riskyEventTypesKeys.all, 'list'] as const,
  list: (filters: ListRiskyEventTypesRequest) => [...riskyEventTypesKeys.lists(), filters] as const,
  details: () => [...riskyEventTypesKeys.all, 'detail'] as const,
  detail: (id: number) => [...riskyEventTypesKeys.details(), id] as const,
}

export const matchRiskyEventsKeys = {
  all: ['matchRiskyEvents'] as const,
  match: (eventId: number, contestId: number) => [...matchRiskyEventsKeys.all, eventId, contestId] as const,
}

// ==================== Event Types Hooks ====================

/**
 * Fetch list of risky event types
 */
export const useRiskyEventTypes = (request: ListRiskyEventTypesRequest = {}) => {
  return useQuery({
    queryKey: riskyEventTypesKeys.list(request),
    queryFn: () => riskyEventsService.listEventTypes(request),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

/**
 * Fetch single risky event type
 */
export const useRiskyEventType = (id: number) => {
  return useQuery({
    queryKey: riskyEventTypesKeys.detail(id),
    queryFn: () => riskyEventsService.getEventType(id),
    enabled: !!id,
  })
}

/**
 * Create new risky event type
 */
export const useCreateRiskyEventType = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: CreateRiskyEventTypeRequest) => 
      riskyEventsService.createEventType(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.lists() })
      showToast('Событие создано!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Ошибка создания: ${error.message}`, 'error')
    },
  })
}

/**
 * Update risky event type
 */
export const useUpdateRiskyEventType = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: UpdateRiskyEventTypeRequest) => 
      riskyEventsService.updateEventType(request),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.lists() })
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.detail(variables.id) })
      showToast('Событие обновлено!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Ошибка обновления: ${error.message}`, 'error')
    },
  })
}

/**
 * Delete risky event type
 */
export const useDeleteRiskyEventType = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (id: number) => riskyEventsService.deleteEventType(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.lists() })
      showToast('Событие удалено!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Ошибка удаления: ${error.message}`, 'error')
    },
  })
}

/**
 * Toggle risky event type active status
 */
export const useToggleRiskyEventType = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: ({ id, isActive }: { id: number; isActive: boolean }) =>
      riskyEventsService.updateEventType({ id, isActive }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.lists() })
      showToast(
        variables.isActive ? 'Событие активировано' : 'Событие деактивировано',
        'success'
      )
    },
    onError: (error: Error) => {
      showToast(`Ошибка: ${error.message}`, 'error')
    },
  })
}

// ==================== Match Events Hooks ====================

/**
 * Fetch risky events for a specific match
 */
export const useMatchRiskyEvents = (eventId: number, contestId: number, enabled = true) => {
  return useQuery({
    queryKey: matchRiskyEventsKeys.match(eventId, contestId),
    queryFn: () => riskyEventsService.getMatchEvents({ eventId, contestId }),
    enabled: enabled && !!eventId && !!contestId,
  })
}

/**
 * Set/update risky events for a match
 */
export const useSetMatchRiskyEvents = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: SetMatchRiskyEventsRequest) =>
      riskyEventsService.setMatchEvents(request),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ 
        queryKey: matchRiskyEventsKeys.all 
      })
      showToast('События матча обновлены!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Ошибка: ${error.message}`, 'error')
    },
  })
}

/**
 * Set outcome for a match risky event (after match finished)
 */
export const useSetMatchEventOutcome = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: SetMatchRiskyEventOutcomeRequest) =>
      riskyEventsService.setMatchEventOutcome(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: matchRiskyEventsKeys.all })
      showToast('Результат события сохранён!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Ошибка: ${error.message}`, 'error')
    },
  })
}

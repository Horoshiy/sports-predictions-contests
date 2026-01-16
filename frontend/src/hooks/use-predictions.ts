import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import predictionService from '../services/prediction-service'
import { useToast } from '../contexts/ToastContext'
import type {
  SubmitPredictionRequest,
  GetUserPredictionsRequest,
  UpdatePredictionRequest,
  ListEventsRequest,
} from '../types/prediction.types'

// Query keys
export const predictionKeys = {
  all: ['predictions'] as const,
  lists: () => [...predictionKeys.all, 'list'] as const,
  list: (contestId: number, pagination?: { page: number; limit: number }) => 
    [...predictionKeys.lists(), contestId, pagination] as const,
  details: () => [...predictionKeys.all, 'detail'] as const,
  detail: (id: number) => [...predictionKeys.details(), id] as const,
}

export const eventKeys = {
  all: ['events'] as const,
  lists: () => [...eventKeys.all, 'list'] as const,
  list: (filters: ListEventsRequest) => [...eventKeys.lists(), filters] as const,
  details: () => [...eventKeys.all, 'detail'] as const,
  detail: (id: number) => [...eventKeys.details(), id] as const,
}

export const propTypeKeys = {
  all: ['propTypes'] as const,
  lists: () => [...propTypeKeys.all, 'list'] as const,
  list: (sportType: string) => [...propTypeKeys.lists(), sportType] as const,
}

// Fetch prop types for a sport
export const usePropTypes = (sportType: string) => {
  return useQuery({
    queryKey: propTypeKeys.list(sportType),
    queryFn: () => predictionService.getPropTypes(sportType),
    enabled: !!sportType,
    staleTime: 10 * 60 * 1000,
  })
}

// Fetch user predictions for a contest
export const useUserPredictions = (request: GetUserPredictionsRequest) => {
  return useQuery({
    queryKey: predictionKeys.list(request.contestId, request.pagination),
    queryFn: () => predictionService.getUserPredictions(request),
    enabled: !!request.contestId,
    staleTime: 2 * 60 * 1000,
  })
}

// Fetch single prediction
export const usePrediction = (id: number) => {
  return useQuery({
    queryKey: predictionKeys.detail(id),
    queryFn: () => predictionService.getPrediction(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}

// Fetch events with filtering
export const useEvents = (request: ListEventsRequest = {}) => {
  return useQuery({
    queryKey: eventKeys.list(request),
    queryFn: () => predictionService.listEvents(request),
    staleTime: 2 * 60 * 1000,
  })
}

// Fetch single event
export const useEvent = (id: number) => {
  return useQuery({
    queryKey: eventKeys.detail(id),
    queryFn: () => predictionService.getEvent(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}

// Submit prediction mutation
export const useSubmitPrediction = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: SubmitPredictionRequest) => 
      predictionService.submitPrediction(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: predictionKeys.lists() })
      showToast('Prediction submitted successfully!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to submit prediction: ${error.message}`, 'error')
    },
  })
}

// Update prediction mutation
export const useUpdatePrediction = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: UpdatePredictionRequest) => 
      predictionService.updatePrediction(request),
    onSuccess: (updatedPrediction) => {
      queryClient.setQueryData(predictionKeys.detail(updatedPrediction.id), updatedPrediction)
      queryClient.invalidateQueries({ queryKey: predictionKeys.lists() })
      showToast('Prediction updated successfully!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to update prediction: ${error.message}`, 'error')
    },
  })
}

// Delete prediction mutation
export const useDeletePrediction = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (id: number) => predictionService.deletePrediction(id),
    onSuccess: (_, deletedId) => {
      queryClient.removeQueries({ queryKey: predictionKeys.detail(deletedId) })
      queryClient.invalidateQueries({ queryKey: predictionKeys.lists() })
      showToast('Prediction deleted successfully!', 'success')
    },
    onError: (error: Error) => {
      showToast(`Failed to delete prediction: ${error.message}`, 'error')
    },
  })
}

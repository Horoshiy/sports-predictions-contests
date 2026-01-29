import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import contestService from '../services/contest-service'
import { useToast } from '../contexts/ToastContext'
import type {
  Contest,
  CreateContestRequest,
  UpdateContestRequest,
  ListContestsRequest,
  JoinContestRequest,
  LeaveContestRequest,
  ListParticipantsRequest,
} from '../types/contest.types'

// Query keys
export const contestKeys = {
  all: ['contests'] as const,
  lists: () => [...contestKeys.all, 'list'] as const,
  list: (filters: ListContestsRequest) => [...contestKeys.lists(), filters] as const,
  details: () => [...contestKeys.all, 'detail'] as const,
  detail: (id: number) => [...contestKeys.details(), id] as const,
  participants: (contestId: number) => [...contestKeys.all, 'participants', contestId] as const,
}

// Fetch contests with optional filtering
export const useContests = (request: ListContestsRequest = {}) => {
  return useQuery({
    queryKey: contestKeys.list(request),
    queryFn: () => contestService.listContests(request),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Fetch single contest
export const useContest = (id: number) => {
  return useQuery({
    queryKey: contestKeys.detail(id),
    queryFn: () => contestService.getContest(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Fetch contest participants
export const useContestParticipants = (request: ListParticipantsRequest) => {
  return useQuery({
    queryKey: contestKeys.participants(request.contestId),
    queryFn: () => contestService.listParticipants(request),
    enabled: !!request.contestId,
    staleTime: 2 * 60 * 1000, // 2 minutes
  })
}

// Create contest mutation
export const useCreateContest = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: CreateContestRequest) => contestService.createContest(request),
    onSuccess: (newContest) => {
      console.log('Contest created:', newContest)
      
      // Invalidate and refetch contests list
      queryClient.invalidateQueries({ queryKey: contestKeys.lists() })
      
      // Add the new contest to the cache if it has an id
      if (newContest?.id) {
        queryClient.setQueryData(contestKeys.detail(newContest.id), newContest)
      }
      
      showToast('Contest created successfully!', 'success')
    },
    onError: (error) => {
      console.error('Failed to create contest:', error)
      showToast(`Failed to create contest: ${error.message}`, 'error')
    },
  })
}

// Update contest mutation
export const useUpdateContest = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: UpdateContestRequest) => contestService.updateContest(request),
    onSuccess: (updatedContest) => {
      // Update the contest in the cache
      queryClient.setQueryData(contestKeys.detail(updatedContest.id), updatedContest)
      
      // Invalidate lists to ensure consistency
      queryClient.invalidateQueries({ queryKey: contestKeys.lists() })
      
      showToast('Contest updated successfully!', 'success')
    },
    onError: (error) => {
      console.error('Failed to update contest:', error)
      showToast(`Failed to update contest: ${error.message}`, 'error')
    },
  })
}

// Delete contest mutation
export const useDeleteContest = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (id: number) => {
      console.log('Deleting contest:', id)
      return contestService.deleteContest(id)
    },
    onSuccess: (_, deletedId) => {
      console.log('Contest deleted successfully:', deletedId)
      
      // Get all queries with contests lists
      const queries = queryClient.getQueriesData({ queryKey: contestKeys.lists() })
      console.log('Found queries to update:', queries.length)
      
      // Update each query to remove the deleted contest
      queryClient.setQueriesData({ queryKey: contestKeys.lists() }, (old: any) => {
        console.log('Updating query data, old:', old)
        if (!old?.contests) {
          console.log('No contests in old data')
          return old
        }
        
        const newContests = old.contests.filter((c: Contest) => c.id !== deletedId)
        console.log('Filtered contests:', old.contests.length, '->', newContests.length)
        
        return {
          ...old,
          contests: newContests,
          pagination: {
            ...old.pagination,
            total: Math.max(0, (old.pagination?.total || 0) - 1),
          },
        }
      })
      
      // Remove the contest from the detail cache
      queryClient.removeQueries({ queryKey: contestKeys.detail(deletedId) })
      
      showToast('Contest deleted successfully!', 'success')
    },
    onError: (error) => {
      console.error('Failed to delete contest:', error)
      showToast(`Failed to delete contest: ${error.message}`, 'error')
    },
  })
}

// Join contest mutation
export const useJoinContest = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: JoinContestRequest) => contestService.joinContest(request),
    onSuccess: (_, { contestId }) => {
      // Invalidate contest details and participants
      queryClient.invalidateQueries({ queryKey: contestKeys.detail(contestId) })
      queryClient.invalidateQueries({ queryKey: contestKeys.participants(contestId) })
      
      // Invalidate lists to update participant counts
      queryClient.invalidateQueries({ queryKey: contestKeys.lists() })
      
      showToast('Successfully joined contest!', 'success')
    },
    onError: (error) => {
      console.error('Failed to join contest:', error)
      showToast(`Failed to join contest: ${error.message}`, 'error')
    },
  })
}

// Leave contest mutation
export const useLeaveContest = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

  return useMutation({
    mutationFn: (request: LeaveContestRequest) => contestService.leaveContest(request),
    onSuccess: (_, { contestId }) => {
      // Invalidate contest details and participants
      queryClient.invalidateQueries({ queryKey: contestKeys.detail(contestId) })
      queryClient.invalidateQueries({ queryKey: contestKeys.participants(contestId) })
      
      // Invalidate lists to update participant counts
      queryClient.invalidateQueries({ queryKey: contestKeys.lists() })
      
      showToast('Successfully left contest!', 'success')
    },
    onError: (error) => {
      console.error('Failed to leave contest:', error)
      showToast(`Failed to leave contest: ${error.message}`, 'error')
    },
  })
}

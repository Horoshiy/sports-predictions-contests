import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import teamService from '../services/team-service'
import { useToast } from '../contexts/ToastContext'
import type { CreateTeamRequest, UpdateTeamRequest, ListTeamsRequest, JoinTeamRequest, ListMembersRequest } from '../types/team.types'

export const teamKeys = {
  all: ['teams'] as const,
  lists: () => [...teamKeys.all, 'list'] as const,
  list: (filters: ListTeamsRequest) => [...teamKeys.lists(), filters] as const,
  details: () => [...teamKeys.all, 'detail'] as const,
  detail: (id: number) => [...teamKeys.details(), id] as const,
  members: (teamId: number) => [...teamKeys.all, 'members', teamId] as const,
  leaderboard: (contestId: number) => ['team-leaderboard', contestId] as const,
}

export const useTeams = (request: ListTeamsRequest = {}) => useQuery({
  queryKey: teamKeys.list(request),
  queryFn: () => teamService.listTeams(request),
  staleTime: 5 * 60 * 1000,
})

export const useTeam = (id: number) => useQuery({
  queryKey: teamKeys.detail(id),
  queryFn: () => teamService.getTeam(id),
  enabled: !!id,
  staleTime: 5 * 60 * 1000,
})

export const useTeamMembers = (request: ListMembersRequest) => useQuery({
  queryKey: [teamKeys.members(request.teamId), request.pagination],
  queryFn: () => teamService.listMembers(request),
  enabled: !!request.teamId,
  staleTime: 2 * 60 * 1000,
})

export const useTeamLeaderboard = (contestId: number, limit = 10) => useQuery({
  queryKey: teamKeys.leaderboard(contestId),
  queryFn: () => teamService.getTeamLeaderboard(contestId, limit),
  enabled: !!contestId,
  staleTime: 30 * 1000,
})

export const useCreateTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: CreateTeamRequest) => teamService.createTeam(request),
    onSuccess: (newTeam) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      queryClient.setQueryData(teamKeys.detail(newTeam.id), newTeam)
      showToast('Team created successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to create team: ${error.message}`, 'error'),
  })
}

export const useUpdateTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: UpdateTeamRequest) => teamService.updateTeam(request),
    onSuccess: (updatedTeam) => {
      queryClient.setQueryData(teamKeys.detail(updatedTeam.id), updatedTeam)
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Team updated successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to update team: ${error.message}`, 'error'),
  })
}

export const useDeleteTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (id: number) => teamService.deleteTeam(id),
    onSuccess: (_, deletedId) => {
      queryClient.removeQueries({ queryKey: teamKeys.detail(deletedId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Team deleted successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to delete team: ${error.message}`, 'error'),
  })
}

export const useJoinTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: JoinTeamRequest) => teamService.joinTeam(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Successfully joined team!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to join team: ${error.message}`, 'error'),
  })
}

export const useLeaveTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (teamId: number) => teamService.leaveTeam(teamId),
    onSuccess: (_, teamId) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.detail(teamId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.members(teamId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.lists() })
      showToast('Successfully left team!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to leave team: ${error.message}`, 'error'),
  })
}

export const useRemoveMember = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: ({ teamId, userId }: { teamId: number; userId: number }) => teamService.removeMember(teamId, userId),
    onSuccess: (_, { teamId }) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.detail(teamId) })
      queryClient.invalidateQueries({ queryKey: teamKeys.members(teamId) })
      showToast('Member removed successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to remove member: ${error.message}`, 'error'),
  })
}

export const useRegenerateInviteCode = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (teamId: number) => teamService.regenerateInviteCode(teamId),
    onSuccess: (_, teamId) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.detail(teamId) })
      showToast('Invite code regenerated!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to regenerate code: ${error.message}`, 'error'),
  })
}

export const useJoinContestAsTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: ({ teamId, contestId }: { teamId: number; contestId: number }) => teamService.joinContestAsTeam(teamId, contestId),
    onSuccess: (_, { contestId }) => {
      queryClient.invalidateQueries({ queryKey: teamKeys.leaderboard(contestId) })
      showToast('Team joined contest!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to join contest: ${error.message}`, 'error'),
  })
}

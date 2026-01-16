import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { sportsService } from '../services/sports-service'
import { useToast } from '../contexts/ToastContext'
import type {
  Sport, League, Team, Match,
  CreateSportRequest, UpdateSportRequest, ListSportsRequest,
  CreateLeagueRequest, UpdateLeagueRequest, ListLeaguesRequest,
  CreateTeamRequest, UpdateTeamRequest, ListTeamsRequest,
  CreateMatchRequest, UpdateMatchRequest, ListMatchesRequest,
} from '../types/sports.types'

// Query keys
export const sportsKeys = {
  all: ['sports'] as const,
  lists: () => [...sportsKeys.all, 'list'] as const,
  list: (filters: ListSportsRequest) => [...sportsKeys.lists(), filters] as const,
  details: () => [...sportsKeys.all, 'detail'] as const,
  detail: (id: number) => [...sportsKeys.details(), id] as const,
}

export const leaguesKeys = {
  all: ['leagues'] as const,
  lists: () => [...leaguesKeys.all, 'list'] as const,
  list: (filters: ListLeaguesRequest) => [...leaguesKeys.lists(), filters] as const,
  details: () => [...leaguesKeys.all, 'detail'] as const,
  detail: (id: number) => [...leaguesKeys.details(), id] as const,
}

export const teamsKeys = {
  all: ['teams'] as const,
  lists: () => [...teamsKeys.all, 'list'] as const,
  list: (filters: ListTeamsRequest) => [...teamsKeys.lists(), filters] as const,
  details: () => [...teamsKeys.all, 'detail'] as const,
  detail: (id: number) => [...teamsKeys.details(), id] as const,
}

export const matchesKeys = {
  all: ['matches'] as const,
  lists: () => [...matchesKeys.all, 'list'] as const,
  list: (filters: ListMatchesRequest) => [...matchesKeys.lists(), filters] as const,
  details: () => [...matchesKeys.all, 'detail'] as const,
  detail: (id: number) => [...matchesKeys.details(), id] as const,
}

// Sports hooks
export const useSports = (request: ListSportsRequest = {}) => {
  return useQuery({
    queryKey: sportsKeys.list(request),
    queryFn: () => sportsService.listSports(request),
    staleTime: 5 * 60 * 1000,
  })
}

export const useSport = (id: number) => {
  return useQuery({
    queryKey: sportsKeys.detail(id),
    queryFn: () => sportsService.getSport(id),
    enabled: !!id,
  })
}

export const useCreateSport = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: CreateSportRequest) => sportsService.createSport(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: sportsKeys.lists() })
      showToast('Sport created successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to create sport: ${error.message}`, 'error'),
  })
}

export const useUpdateSport = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: UpdateSportRequest) => sportsService.updateSport(request),
    onSuccess: (sport) => {
      queryClient.setQueryData(sportsKeys.detail(sport.id), sport)
      queryClient.invalidateQueries({ queryKey: sportsKeys.lists() })
      showToast('Sport updated successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to update sport: ${error.message}`, 'error'),
  })
}

export const useDeleteSport = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (id: number) => sportsService.deleteSport(id),
    onSuccess: (_, id) => {
      queryClient.removeQueries({ queryKey: sportsKeys.detail(id) })
      queryClient.invalidateQueries({ queryKey: sportsKeys.lists() })
      showToast('Sport deleted successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to delete sport: ${error.message}`, 'error'),
  })
}

// Leagues hooks
export const useLeagues = (request: ListLeaguesRequest = {}) => {
  return useQuery({
    queryKey: leaguesKeys.list(request),
    queryFn: () => sportsService.listLeagues(request),
    staleTime: 5 * 60 * 1000,
  })
}

export const useLeague = (id: number) => {
  return useQuery({
    queryKey: leaguesKeys.detail(id),
    queryFn: () => sportsService.getLeague(id),
    enabled: !!id,
  })
}

export const useCreateLeague = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: CreateLeagueRequest) => sportsService.createLeague(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: leaguesKeys.lists() })
      showToast('League created successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to create league: ${error.message}`, 'error'),
  })
}

export const useUpdateLeague = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: UpdateLeagueRequest) => sportsService.updateLeague(request),
    onSuccess: (league) => {
      queryClient.setQueryData(leaguesKeys.detail(league.id), league)
      queryClient.invalidateQueries({ queryKey: leaguesKeys.lists() })
      showToast('League updated successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to update league: ${error.message}`, 'error'),
  })
}

export const useDeleteLeague = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (id: number) => sportsService.deleteLeague(id),
    onSuccess: (_, id) => {
      queryClient.removeQueries({ queryKey: leaguesKeys.detail(id) })
      queryClient.invalidateQueries({ queryKey: leaguesKeys.lists() })
      showToast('League deleted successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to delete league: ${error.message}`, 'error'),
  })
}

// Teams hooks
export const useTeams = (request: ListTeamsRequest = {}) => {
  return useQuery({
    queryKey: teamsKeys.list(request),
    queryFn: () => sportsService.listTeams(request),
    staleTime: 5 * 60 * 1000,
  })
}

export const useTeam = (id: number) => {
  return useQuery({
    queryKey: teamsKeys.detail(id),
    queryFn: () => sportsService.getTeam(id),
    enabled: !!id,
  })
}

export const useCreateTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: CreateTeamRequest) => sportsService.createTeam(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: teamsKeys.lists() })
      showToast('Team created successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to create team: ${error.message}`, 'error'),
  })
}

export const useUpdateTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: UpdateTeamRequest) => sportsService.updateTeam(request),
    onSuccess: (team) => {
      queryClient.setQueryData(teamsKeys.detail(team.id), team)
      queryClient.invalidateQueries({ queryKey: teamsKeys.lists() })
      showToast('Team updated successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to update team: ${error.message}`, 'error'),
  })
}

export const useDeleteTeam = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (id: number) => sportsService.deleteTeam(id),
    onSuccess: (_, id) => {
      queryClient.removeQueries({ queryKey: teamsKeys.detail(id) })
      queryClient.invalidateQueries({ queryKey: teamsKeys.lists() })
      showToast('Team deleted successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to delete team: ${error.message}`, 'error'),
  })
}

// Matches hooks
export const useMatches = (request: ListMatchesRequest = {}) => {
  return useQuery({
    queryKey: matchesKeys.list(request),
    queryFn: () => sportsService.listMatches(request),
    staleTime: 5 * 60 * 1000,
  })
}

export const useMatch = (id: number) => {
  return useQuery({
    queryKey: matchesKeys.detail(id),
    queryFn: () => sportsService.getMatch(id),
    enabled: !!id,
  })
}

export const useCreateMatch = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: CreateMatchRequest) => sportsService.createMatch(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: matchesKeys.lists() })
      showToast('Match created successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to create match: ${error.message}`, 'error'),
  })
}

export const useUpdateMatch = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: UpdateMatchRequest) => sportsService.updateMatch(request),
    onSuccess: (match) => {
      queryClient.setQueryData(matchesKeys.detail(match.id), match)
      queryClient.invalidateQueries({ queryKey: matchesKeys.lists() })
      showToast('Match updated successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to update match: ${error.message}`, 'error'),
  })
}

export const useDeleteMatch = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (id: number) => sportsService.deleteMatch(id),
    onSuccess: (_, id) => {
      queryClient.removeQueries({ queryKey: matchesKeys.detail(id) })
      queryClient.invalidateQueries({ queryKey: matchesKeys.lists() })
      showToast('Match deleted successfully!', 'success')
    },
    onError: (error: Error) => showToast(`Failed to delete match: ${error.message}`, 'error'),
  })
}

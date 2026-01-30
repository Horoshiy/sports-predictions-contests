import grpcClient from './grpc-client'
import type {
  Sport, League, Team, Match,
  CreateSportRequest, UpdateSportRequest, ListSportsRequest,
  CreateLeagueRequest, UpdateLeagueRequest, ListLeaguesRequest,
  CreateTeamRequest, UpdateTeamRequest, ListTeamsRequest,
  CreateMatchRequest, UpdateMatchRequest, ListMatchesRequest,
  SportResponse, ListSportsResponse,
  LeagueResponse, ListLeaguesResponse,
  TeamResponse, ListTeamsResponse,
  MatchResponse, ListMatchesResponse,
  DeleteResponse,
} from '../types/sports.types'
import type { PaginationResponse } from '../types/common.types'

const defaultPagination = { page: 1, limit: 10, total: 0, totalPages: 0 }

class SportsService {
  // Sports
  async createSport(request: CreateSportRequest): Promise<Sport> {
    const response = await grpcClient.post<SportResponse>('/v1/sports', request)
    return response.sport
  }

  async getSport(id: number): Promise<Sport> {
    const response = await grpcClient.get<SportResponse>(`/v1/sports/${id}`)
    return response.sport
  }

  async updateSport(request: UpdateSportRequest): Promise<Sport> {
    const response = await grpcClient.put<SportResponse>(`/v1/sports/${request.id}`, request)
    return response.sport
  }

  async deleteSport(id: number): Promise<void> {
    await grpcClient.delete<DeleteResponse>(`/v1/sports/${id}`)
  }

  async listSports(request: ListSportsRequest = {}): Promise<{ sports: Sport[]; pagination: PaginationResponse }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.activeOnly) params.append('active_only', 'true')
    const url = params.toString() ? `/v1/sports?${params}` : '/v1/sports'
    const response = await grpcClient.get<ListSportsResponse>(url)
    return { sports: response.sports || [], pagination: response.pagination || defaultPagination }
  }

  // Leagues
  async createLeague(request: CreateLeagueRequest): Promise<League> {
    const response = await grpcClient.post<LeagueResponse>('/v1/leagues', request)
    return response.league
  }

  async getLeague(id: number): Promise<League> {
    const response = await grpcClient.get<LeagueResponse>(`/v1/leagues/${id}`)
    return response.league
  }

  async updateLeague(request: UpdateLeagueRequest): Promise<League> {
    const response = await grpcClient.put<LeagueResponse>(`/v1/leagues/${request.id}`, request)
    return response.league
  }

  async deleteLeague(id: number): Promise<void> {
    await grpcClient.delete<DeleteResponse>(`/v1/leagues/${id}`)
  }

  async listLeagues(request: ListLeaguesRequest = {}): Promise<{ leagues: League[]; pagination: PaginationResponse }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.sportId) params.append('sport_id', request.sportId.toString())
    if (request.activeOnly) params.append('active_only', 'true')
    const url = params.toString() ? `/v1/leagues?${params}` : '/v1/leagues'
    const response = await grpcClient.get<ListLeaguesResponse>(url)
    return { leagues: response.leagues || [], pagination: response.pagination || defaultPagination }
  }

  // Teams
  async createTeam(request: CreateTeamRequest): Promise<Team> {
    const response = await grpcClient.post<TeamResponse>('/v1/sports/teams', request)
    return response.team
  }

  async getTeam(id: number): Promise<Team> {
    const response = await grpcClient.get<TeamResponse>(`/v1/sports/teams/${id}`)
    return response.team
  }

  async updateTeam(request: UpdateTeamRequest): Promise<Team> {
    const response = await grpcClient.put<TeamResponse>(`/v1/sports/teams/${request.id}`, request)
    return response.team
  }

  async deleteTeam(id: number): Promise<void> {
    await grpcClient.delete<DeleteResponse>(`/v1/sports/teams/${id}`)
  }

  async listTeams(request: ListTeamsRequest = {}): Promise<{ teams: Team[]; pagination: PaginationResponse }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.sportId) params.append('sport_id', request.sportId.toString())
    if (request.activeOnly) params.append('active_only', 'true')
    const url = params.toString() ? `/v1/sports/teams?${params}` : '/v1/sports/teams'
    const response = await grpcClient.get<ListTeamsResponse>(url)
    return { teams: response.teams || [], pagination: response.pagination || defaultPagination }
  }

  // Matches
  async createMatch(request: CreateMatchRequest): Promise<Match> {
    const response = await grpcClient.post<MatchResponse>('/v1/matches', request)
    return response.match
  }

  async getMatch(id: number): Promise<Match> {
    const response = await grpcClient.get<MatchResponse>(`/v1/matches/${id}`)
    return response.match
  }

  async updateMatch(request: UpdateMatchRequest): Promise<Match> {
    const response = await grpcClient.put<MatchResponse>(`/v1/matches/${request.id}`, request)
    return response.match
  }

  async deleteMatch(id: number): Promise<void> {
    await grpcClient.delete<DeleteResponse>(`/v1/matches/${id}`)
  }

  async listMatches(request: ListMatchesRequest = {}): Promise<{ matches: Match[]; pagination: PaginationResponse }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.leagueId) params.append('league_id', request.leagueId.toString())
    if (request.teamId) params.append('team_id', request.teamId.toString())
    if (request.status) params.append('status', request.status)
    const url = params.toString() ? `/v1/matches?${params}` : '/v1/matches'
    const response = await grpcClient.get<ListMatchesResponse>(url)
    return { matches: response.matches || [], pagination: response.pagination || defaultPagination }
  }
}

export const sportsService = new SportsService()
export default sportsService

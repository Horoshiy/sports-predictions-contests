import grpcClient from './grpc-client'
import type {
  Team, TeamMember, TeamLeaderboardEntry,
  CreateTeamRequest, UpdateTeamRequest, ListTeamsRequest,
  JoinTeamRequest, ListMembersRequest, PaginationResponse,
} from '../types/team.types'

class TeamService {
  private basePath = '/v1/teams'

  async createTeam(request: CreateTeamRequest): Promise<Team> {
    const response = await grpcClient.post<{ team: Team }>(this.basePath, request)
    return response.team
  }

  async updateTeam(request: UpdateTeamRequest): Promise<Team> {
    const response = await grpcClient.put<{ team: Team }>(`${this.basePath}/${request.id}`, request)
    return response.team
  }

  async getTeam(id: number): Promise<Team> {
    const response = await grpcClient.get<{ team: Team }>(`${this.basePath}/${id}`)
    return response.team
  }

  async deleteTeam(id: number): Promise<void> {
    await grpcClient.delete(`${this.basePath}/${id}`)
  }

  async listTeams(request: ListTeamsRequest = {}): Promise<{ teams: Team[]; pagination: PaginationResponse }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.myTeamsOnly) params.append('my_teams_only', 'true')
    const url = params.toString() ? `${this.basePath}?${params}` : this.basePath
    const response = await grpcClient.get<{ teams: Team[]; pagination: PaginationResponse }>(url)
    return { teams: response.teams || [], pagination: response.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 } }
  }

  async joinTeam(request: JoinTeamRequest): Promise<TeamMember> {
    const response = await grpcClient.post<{ member: TeamMember }>(`${this.basePath}/join`, request)
    return response.member
  }

  async leaveTeam(teamId: number): Promise<void> {
    await grpcClient.post(`${this.basePath}/${teamId}/leave`, {})
  }

  async removeMember(teamId: number, userId: number): Promise<void> {
    await grpcClient.delete(`${this.basePath}/${teamId}/members/${userId}`)
  }

  async listMembers(request: ListMembersRequest): Promise<{ members: TeamMember[]; pagination: PaginationResponse }> {
    const params = new URLSearchParams()
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    const url = params.toString() ? `${this.basePath}/${request.teamId}/members?${params}` : `${this.basePath}/${request.teamId}/members`
    const response = await grpcClient.get<{ members: TeamMember[]; pagination: PaginationResponse }>(url)
    return { members: response.members || [], pagination: response.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 } }
  }

  async regenerateInviteCode(teamId: number): Promise<string> {
    const response = await grpcClient.post<{ inviteCode: string }>(`${this.basePath}/${teamId}/regenerate-invite`, {})
    return response.inviteCode
  }

  async joinContestAsTeam(teamId: number, contestId: number): Promise<void> {
    await grpcClient.post(`${this.basePath}/${teamId}/contests/${contestId}/join`, {})
  }

  async getTeamLeaderboard(contestId: number, limit = 10): Promise<TeamLeaderboardEntry[]> {
    const response = await grpcClient.get<{ entries: TeamLeaderboardEntry[] }>(`/v1/contests/${contestId}/team-leaderboard?limit=${limit}`)
    return response.entries || []
  }
}

export const teamService = new TeamService()
export default teamService

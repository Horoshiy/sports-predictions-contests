import grpcClient from './grpc-client'
import type {
  Challenge,
  CreateChallengeRequest,
  AcceptChallengeRequest,
  DeclineChallengeRequest,
  WithdrawChallengeRequest,
  GetChallengeRequest,
  ListUserChallengesRequest,
  ListOpenChallengesRequest,
  CreateChallengeResponse,
  AcceptChallengeResponse,
  DeclineChallengeResponse,
  WithdrawChallengeResponse,
  GetChallengeResponse,
  ListUserChallengesResponse,
  ListOpenChallengesResponse,
  PaginationRequest,
} from '../types/challenge.types'

class ChallengeService {
  private basePath = '/v1/challenges'

  // Create a new challenge
  async createChallenge(request: CreateChallengeRequest): Promise<Challenge> {
    const response = await grpcClient.post<CreateChallengeResponse>(
      this.basePath,
      request
    )
    return response.challenge
  }

  // Accept a challenge
  async acceptChallenge(request: AcceptChallengeRequest): Promise<Challenge> {
    const response = await grpcClient.put<AcceptChallengeResponse>(
      `${this.basePath}/${request.id}/accept`,
      request
    )
    return response.challenge
  }

  // Decline a challenge
  async declineChallenge(request: DeclineChallengeRequest): Promise<void> {
    await grpcClient.put<DeclineChallengeResponse>(
      `${this.basePath}/${request.id}/decline`,
      request
    )
  }

  // Withdraw a challenge
  async withdrawChallenge(request: WithdrawChallengeRequest): Promise<void> {
    await grpcClient.put<WithdrawChallengeResponse>(
      `${this.basePath}/${request.id}/withdraw`,
      request
    )
  }

  // Get a challenge by ID
  async getChallenge(id: number): Promise<Challenge> {
    const response = await grpcClient.get<GetChallengeResponse>(
      `${this.basePath}/${id}`
    )
    return response.challenge
  }

  // List challenges for a user
  async listUserChallenges(request: ListUserChallengesRequest): Promise<{
    challenges: Challenge[]
    pagination: {
      page: number
      limit: number
      total: number
      totalPages: number
    }
  }> {
    const params = new URLSearchParams()
    
    if (request.status) {
      params.append('status', request.status)
    }
    
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
      if (request.pagination.sortBy) {
        params.append('sort_by', request.pagination.sortBy)
      }
      if (request.pagination.sortOrder) {
        params.append('sort_order', request.pagination.sortOrder)
      }
    }

    const queryString = params.toString()
    const url = `/v1/users/${request.userId}/challenges${queryString ? `?${queryString}` : ''}`
    
    const response = await grpcClient.get<ListUserChallengesResponse>(url)
    
    return {
      challenges: response.challenges,
      pagination: response.pagination
    }
  }

  // List open challenges for an event
  async listOpenChallenges(request: ListOpenChallengesRequest): Promise<{
    challenges: Challenge[]
    pagination: {
      page: number
      limit: number
      total: number
      totalPages: number
    }
  }> {
    const params = new URLSearchParams()
    
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
      if (request.pagination.sortBy) {
        params.append('sort_by', request.pagination.sortBy)
      }
      if (request.pagination.sortOrder) {
        params.append('sort_order', request.pagination.sortOrder)
      }
    }

    const queryString = params.toString()
    const url = `/v1/events/${request.eventId}/challenges${queryString ? `?${queryString}` : ''}`
    
    const response = await grpcClient.get<ListOpenChallengesResponse>(url)
    
    return {
      challenges: response.challenges,
      pagination: response.pagination
    }
  }

  // Helper method to get challenges by status
  async getChallengesByStatus(userId: number, status: Challenge['status']): Promise<Challenge[]> {
    const result = await this.listUserChallenges({
      userId,
      status,
      pagination: { page: 1, limit: 100 }
    })
    return result.challenges
  }

  // Helper method to get pending challenges for user
  async getPendingChallenges(userId: number): Promise<Challenge[]> {
    return this.getChallengesByStatus(userId, 'pending')
  }

  // Helper method to get active challenges for user
  async getActiveChallenges(userId: number): Promise<Challenge[]> {
    return this.getChallengesByStatus(userId, 'active')
  }

  // Helper method to get completed challenges for user
  async getCompletedChallenges(userId: number): Promise<Challenge[]> {
    return this.getChallengesByStatus(userId, 'completed')
  }

  // Check if user can challenge another user
  canChallenge(challengerId: number, opponentId: number): boolean {
    return challengerId !== opponentId
  }

  // Get challenge status display info
  getChallengeStatusInfo(challenge: Challenge) {
    const now = new Date()
    const expiresAt = new Date(challenge.expiresAt)
    
    if (challenge.status === 'pending' && now > expiresAt) {
      return {
        status: 'expired' as const,
        canAccept: false,
        canDecline: false,
        canWithdraw: false
      }
    }

    return {
      status: challenge.status,
      canAccept: challenge.status === 'pending' && now < expiresAt,
      canDecline: challenge.status === 'pending',
      canWithdraw: challenge.status === 'pending'
    }
  }
}

// Export singleton instance
const challengeService = new ChallengeService()
export default challengeService

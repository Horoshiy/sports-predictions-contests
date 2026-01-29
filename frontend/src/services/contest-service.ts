import grpcClient from './grpc-client'
import type {
  Contest,
  Participant,
  CreateContestRequest,
  UpdateContestRequest,
  ListContestsRequest,
  JoinContestRequest,
  LeaveContestRequest,
  ListParticipantsRequest,
  CreateContestResponse,
  UpdateContestResponse,
  GetContestResponse,
  DeleteContestResponse,
  ListContestsResponse,
  JoinContestResponse,
  LeaveContestResponse,
  ListParticipantsResponse,
  PaginationRequest,
} from '../types/contest.types'

class ContestService {
  private basePath = '/v1/contests'

  // Create a new contest
  async createContest(request: CreateContestRequest): Promise<Contest> {
    const response = await grpcClient.post<CreateContestResponse>(
      this.basePath,
      request
    )
    console.log('Create contest response:', response)
    if (!response.contest) {
      throw new Error('Server returned empty contest data')
    }
    return response.contest
  }

  // Update an existing contest
  async updateContest(request: UpdateContestRequest): Promise<Contest> {
    const response = await grpcClient.put<UpdateContestResponse>(
      `${this.basePath}/${request.id}`,
      request
    )
    return response.contest
  }

  // Get a contest by ID
  async getContest(id: number): Promise<Contest> {
    const response = await grpcClient.get<GetContestResponse>(
      `${this.basePath}/${id}`
    )
    return response.contest
  }

  // Delete a contest
  async deleteContest(id: number): Promise<void> {
    await grpcClient.delete<DeleteContestResponse>(`${this.basePath}/${id}`)
  }

  // List contests with optional filtering and pagination
  async listContests(request: ListContestsRequest = {}): Promise<{
    contests: Contest[]
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
    
    if (request.status) {
      params.append('status', request.status)
    }
    
    if (request.sportType) {
      params.append('sport_type', request.sportType)
    }

    const queryString = params.toString()
    const url = queryString ? `${this.basePath}?${queryString}` : this.basePath

    const response = await grpcClient.get<ListContestsResponse>(url)
    
    return {
      contests: response.contests,
      pagination: response.pagination,
    }
  }

  // Join a contest
  async joinContest(request: JoinContestRequest): Promise<void> {
    await grpcClient.post<JoinContestResponse>(
      `${this.basePath}/${request.contestId}/join`,
      {}
    )
  }

  // Leave a contest
  async leaveContest(request: LeaveContestRequest): Promise<void> {
    await grpcClient.post<LeaveContestResponse>(
      `${this.basePath}/${request.contestId}/leave`,
      {}
    )
  }

  // List participants of a contest
  async listParticipants(request: ListParticipantsRequest): Promise<{
    participants: Participant[]
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
    const url = queryString 
      ? `${this.basePath}/${request.contestId}/participants?${queryString}`
      : `${this.basePath}/${request.contestId}/participants`

    const response = await grpcClient.get<ListParticipantsResponse>(url)
    
    return {
      participants: response.participants,
      pagination: response.pagination,
    }
  }
}

// Singleton instance
export const contestService = new ContestService()
export default contestService

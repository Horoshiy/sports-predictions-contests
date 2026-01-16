import grpcClient from './grpc-client'
import type {
  Score,
  Leaderboard,
  CreateScoreRequest,
  UpdateScoreRequest,
  GetUserScoresRequest,
  GetLeaderboardRequest,
  GetUserRankRequest,
  GetUserStreakRequest,
  UpdateLeaderboardRequest,
  CalculateScoreRequest,
  ListScoresRequest,
  CreateScoreResponse,
  UpdateScoreResponse,
  GetScoreResponse,
  DeleteScoreResponse,
  ListScoresResponse,
  GetUserScoresResponse,
  GetLeaderboardResponse,
  GetUserRankResponse,
  GetUserStreakResponse,
  UpdateLeaderboardResponse,
  CalculateScoreResponse,
  PaginationRequest,
} from '../types/scoring.types'

class ScoringService {
  private basePath = '/v1/scores'
  private leaderboardPath = '/v1/contests'

  // Create a new score
  async createScore(request: CreateScoreRequest): Promise<Score> {
    const response = await grpcClient.post<CreateScoreResponse>(
      this.basePath,
      request
    )
    return response.score
  }

  // Update an existing score
  async updateScore(request: UpdateScoreRequest): Promise<Score> {
    const response = await grpcClient.put<UpdateScoreResponse>(
      `${this.basePath}/${request.id}`,
      request
    )
    return response.score
  }

  // Get a score by ID
  async getScore(id: number): Promise<Score> {
    const response = await grpcClient.get<GetScoreResponse>(
      `${this.basePath}/${id}`
    )
    return response.score
  }

  // Delete a score
  async deleteScore(id: number): Promise<void> {
    await grpcClient.delete<DeleteScoreResponse>(`${this.basePath}/${id}`)
  }

  // List scores with optional filtering and pagination
  async listScores(request: ListScoresRequest = {}): Promise<{
    scores: Score[]
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
    
    if (request.contestId) {
      params.append('contest_id', request.contestId.toString())
    }
    
    if (request.userId) {
      params.append('user_id', request.userId.toString())
    }

    const queryString = params.toString()
    const url = queryString ? `${this.basePath}?${queryString}` : this.basePath

    const response = await grpcClient.get<ListScoresResponse>(url)
    
    return {
      scores: response.scores,
      pagination: response.pagination,
    }
  }

  // Get user scores for a specific contest
  async getUserScores(request: GetUserScoresRequest): Promise<{
    scores: Score[]
    totalPoints: number
  }> {
    const response = await grpcClient.get<GetUserScoresResponse>(
      `/v1/users/${request.userId}/contests/${request.contestId}/scores`
    )
    
    return {
      scores: response.scores,
      totalPoints: response.totalPoints,
    }
  }

  // Get leaderboard for a contest
  async getLeaderboard(request: GetLeaderboardRequest): Promise<Leaderboard> {
    const params = new URLSearchParams()
    
    if (request.limit) {
      params.append('limit', request.limit.toString())
    }

    const queryString = params.toString()
    const url = queryString 
      ? `${this.leaderboardPath}/${request.contestId}/leaderboard?${queryString}`
      : `${this.leaderboardPath}/${request.contestId}/leaderboard`

    const response = await grpcClient.get<GetLeaderboardResponse>(url)
    return response.leaderboard
  }

  // Get user rank in a contest
  async getUserRank(request: GetUserRankRequest): Promise<{
    rank: number
    totalPoints: number
  }> {
    const response = await grpcClient.get<GetUserRankResponse>(
      `${this.leaderboardPath}/${request.contestId}/users/${request.userId}/rank`
    )
    
    return {
      rank: response.rank,
      totalPoints: response.totalPoints,
    }
  }

  // Get user streak in a contest
  async getUserStreak(request: GetUserStreakRequest): Promise<{
    currentStreak: number
    maxStreak: number
    multiplier: number
  }> {
    const response = await grpcClient.get<GetUserStreakResponse>(
      `${this.leaderboardPath}/${request.contestId}/users/${request.userId}/streak`
    )
    
    return {
      currentStreak: response.currentStreak,
      maxStreak: response.maxStreak,
      multiplier: response.multiplier,
    }
  }

  // Update leaderboard for a contest
  async updateLeaderboard(request: UpdateLeaderboardRequest): Promise<Leaderboard> {
    const response = await grpcClient.post<UpdateLeaderboardResponse>(
      `${this.leaderboardPath}/${request.contestId}/leaderboard/update`,
      {}
    )
    return response.leaderboard
  }

  // Calculate score for a prediction
  async calculateScore(request: CalculateScoreRequest): Promise<{
    points: number
    calculationDetails: string
  }> {
    const response = await grpcClient.post<CalculateScoreResponse>(
      `${this.basePath}/calculate`,
      request
    )
    
    return {
      points: response.points,
      calculationDetails: response.calculationDetails,
    }
  }

  // Real-time polling method for leaderboard updates
  async pollLeaderboard(
    contestId: number, 
    limit?: number, 
    intervalMs: number = 5000
  ): Promise<() => void> {
    const poll = async () => {
      try {
        await this.getLeaderboard({ contestId, limit })
      } catch (error) {
        console.error('Failed to poll leaderboard:', error)
      }
    }

    const intervalId = setInterval(poll, intervalMs)
    
    // Return cleanup function
    return () => clearInterval(intervalId)
  }

  // Batch score operations helper
  async batchCreateScores(scores: CreateScoreRequest[]): Promise<Score[]> {
    const results: Score[] = []
    
    // Process in batches to avoid overwhelming the server
    const batchSize = 10
    for (let i = 0; i < scores.length; i += batchSize) {
      const batch = scores.slice(i, i + batchSize)
      const batchPromises = batch.map(score => this.createScore(score))
      const batchResults = await Promise.all(batchPromises)
      results.push(...batchResults)
    }
    
    return results
  }

  // Helper method to get leaderboard with caching
  private leaderboardCache = new Map<string, { data: Leaderboard; timestamp: number }>()
  private readonly CACHE_TTL = 30000 // 30 seconds

  async getCachedLeaderboard(request: GetLeaderboardRequest): Promise<Leaderboard> {
    const cacheKey = `${request.contestId}-${request.limit || 50}`
    const cached = this.leaderboardCache.get(cacheKey)
    
    if (cached && Date.now() - cached.timestamp < this.CACHE_TTL) {
      return cached.data
    }
    
    const leaderboard = await this.getLeaderboard(request)
    this.leaderboardCache.set(cacheKey, {
      data: leaderboard,
      timestamp: Date.now()
    })
    
    return leaderboard
  }

  // Clear leaderboard cache
  clearLeaderboardCache(): void {
    this.leaderboardCache.clear()
  }
}

// Singleton instance
export const scoringService = new ScoringService()
export default scoringService

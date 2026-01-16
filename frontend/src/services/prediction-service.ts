import grpcClient from './grpc-client'
import type {
  Prediction,
  Event,
  SubmitPredictionRequest,
  GetUserPredictionsRequest,
  UpdatePredictionRequest,
  ListEventsRequest,
  SubmitPredictionResponse,
  GetPredictionResponse,
  GetUserPredictionsResponse,
  UpdatePredictionResponse,
  DeletePredictionResponse,
  GetEventResponse,
  ListEventsResponse,
} from '../types/prediction.types'
import type { PaginationResponse } from '../types/common.types'
import type {
  PropType,
  ListPropTypesRequest,
  GetPropTypesResponse,
  ListPropTypesResponse,
} from '../types/props.types'

class PredictionService {
  private predictionsPath = '/v1/predictions'
  private eventsPath = '/v1/events'

  // Submit a new prediction
  async submitPrediction(request: SubmitPredictionRequest): Promise<Prediction> {
    const response = await grpcClient.post<SubmitPredictionResponse>(
      this.predictionsPath,
      request
    )
    return response.prediction
  }

  // Get a prediction by ID
  async getPrediction(id: number): Promise<Prediction> {
    const response = await grpcClient.get<GetPredictionResponse>(
      `${this.predictionsPath}/${id}`
    )
    return response.prediction
  }

  // Get user predictions for a contest
  async getUserPredictions(request: GetUserPredictionsRequest): Promise<{
    predictions: Prediction[]
    pagination: PaginationResponse
  }> {
    const params = new URLSearchParams()
    
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }

    const queryString = params.toString()
    const url = queryString
      ? `${this.predictionsPath}/contest/${request.contestId}?${queryString}`
      : `${this.predictionsPath}/contest/${request.contestId}`

    const response = await grpcClient.get<GetUserPredictionsResponse>(url)
    
    return {
      predictions: response.predictions || [],
      pagination: response.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 },
    }
  }

  // Update an existing prediction
  async updatePrediction(request: UpdatePredictionRequest): Promise<Prediction> {
    const response = await grpcClient.put<UpdatePredictionResponse>(
      `${this.predictionsPath}/${request.id}`,
      { predictionData: request.predictionData }
    )
    return response.prediction
  }

  // Delete a prediction
  async deletePrediction(id: number): Promise<void> {
    await grpcClient.delete<DeletePredictionResponse>(`${this.predictionsPath}/${id}`)
  }

  // Get an event by ID
  async getEvent(id: number): Promise<Event> {
    const response = await grpcClient.get<GetEventResponse>(
      `${this.eventsPath}/${id}`
    )
    return response.event
  }

  // List events with optional filtering
  async listEvents(request: ListEventsRequest = {}): Promise<{
    events: Event[]
    pagination: PaginationResponse
  }> {
    const params = new URLSearchParams()
    
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    
    if (request.sportType) {
      params.append('sport_type', request.sportType)
    }
    
    if (request.status) {
      params.append('status', request.status)
    }

    const queryString = params.toString()
    const url = queryString ? `${this.eventsPath}?${queryString}` : this.eventsPath

    const response = await grpcClient.get<ListEventsResponse>(url)
    
    return {
      events: response.events || [],
      pagination: response.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 },
    }
  }

  // Get prop types for a sport
  async getPropTypes(sportType: string): Promise<PropType[]> {
    const response = await grpcClient.get<GetPropTypesResponse>(
      `/v1/prop-types/${encodeURIComponent(sportType)}`
    )
    if (!response.response?.success) {
      throw new Error(response.response?.message || 'Failed to get prop types')
    }
    return response.propTypes || []
  }

  // List prop types with filtering
  async listPropTypes(request: ListPropTypesRequest = {}): Promise<{
    propTypes: PropType[]
    pagination: PaginationResponse
  }> {
    const params = new URLSearchParams()
    
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.sportType) {
      params.append('sport_type', request.sportType)
    }
    if (request.category) {
      params.append('category', request.category)
    }
    if (request.activeOnly !== undefined) {
      params.append('active_only', request.activeOnly.toString())
    }

    const queryString = params.toString()
    const url = queryString ? `/v1/prop-types?${queryString}` : '/v1/prop-types'

    const response = await grpcClient.get<ListPropTypesResponse>(url)
    
    return {
      propTypes: response.propTypes || [],
      pagination: response.pagination || { page: 1, limit: 20, total: 0, totalPages: 0 },
    }
  }
}

export const predictionService = new PredictionService()
export default predictionService

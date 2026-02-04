import grpcClient from './grpc-client'
import type {
  RiskyEventType,
  ListRiskyEventTypesRequest,
  CreateRiskyEventTypeRequest,
  UpdateRiskyEventTypeRequest,
  GetMatchRiskyEventsRequest,
  SetMatchRiskyEventsRequest,
  SetMatchRiskyEventOutcomeRequest,
  RiskyEventTypeResponse,
  ListRiskyEventTypesResponse,
  GetMatchRiskyEventsResponse,
} from '../types/risky-events.types'

class RiskyEventsService {
  // ==================== Global Event Types ====================

  async listEventTypes(request: ListRiskyEventTypesRequest = {}): Promise<RiskyEventType[]> {
    const params = new URLSearchParams()
    if (request.sportType) params.append('sport_type', request.sportType)
    if (request.activeOnly) params.append('active_only', 'true')
    const url = params.toString() ? `/v1/risky-event-types?${params}` : '/v1/risky-event-types'
    const response = await grpcClient.get<ListRiskyEventTypesResponse>(url)
    return response.eventTypes || []
  }

  async getEventType(id: number): Promise<RiskyEventType> {
    const response = await grpcClient.get<RiskyEventTypeResponse>(`/v1/risky-event-types/${id}`)
    return response.eventType
  }

  async createEventType(request: CreateRiskyEventTypeRequest): Promise<RiskyEventType> {
    // Convert camelCase to snake_case for API
    const payload = {
      slug: request.slug,
      name: request.name,
      name_en: request.nameEn,
      description: request.description,
      default_points: request.defaultPoints,
      sport_type: request.sportType || 'football',
      category: request.category || 'general',
      icon: request.icon,
      sort_order: request.sortOrder || 0,
    }
    const response = await grpcClient.post<RiskyEventTypeResponse>('/v1/risky-event-types', payload)
    return response.eventType
  }

  async updateEventType(request: UpdateRiskyEventTypeRequest): Promise<RiskyEventType> {
    const { id, ...data } = request
    // Convert camelCase to snake_case for API
    const payload: Record<string, unknown> = {}
    if (data.slug !== undefined) payload.slug = data.slug
    if (data.name !== undefined) payload.name = data.name
    if (data.nameEn !== undefined) payload.name_en = data.nameEn
    if (data.description !== undefined) payload.description = data.description
    if (data.defaultPoints !== undefined) payload.default_points = data.defaultPoints
    if (data.sportType !== undefined) payload.sport_type = data.sportType
    if (data.category !== undefined) payload.category = data.category
    if (data.icon !== undefined) payload.icon = data.icon
    if (data.sortOrder !== undefined) payload.sort_order = data.sortOrder
    if (data.isActive !== undefined) payload.is_active = data.isActive

    const response = await grpcClient.put<RiskyEventTypeResponse>(`/v1/risky-event-types/${id}`, payload)
    return response.eventType
  }

  async deleteEventType(id: number): Promise<void> {
    await grpcClient.delete(`/v1/risky-event-types/${id}`)
  }

  // ==================== Match Events ====================

  async getMatchEvents(request: GetMatchRiskyEventsRequest): Promise<GetMatchRiskyEventsResponse> {
    const params = new URLSearchParams()
    params.append('contest_id', request.contestId.toString())
    const response = await grpcClient.get<GetMatchRiskyEventsResponse>(
      `/v1/events/${request.eventId}/risky-events?${params}`
    )
    
    // Validate response structure
    if (!response || typeof response !== 'object') {
      console.warn('[RiskyEventsService] Unexpected response format from getMatchEvents')
      return { events: [], maxSelections: 5 }
    }
    
    return {
      events: response.events || [],
      maxSelections: response.maxSelections || 5,
    }
  }

  async setMatchEvents(request: SetMatchRiskyEventsRequest): Promise<void> {
    const { eventId, events } = request
    // Convert to snake_case
    const payload = {
      events: events.map(e => ({
        risky_event_type_id: e.riskyEventTypeId,
        points: e.points,
        is_enabled: e.isEnabled,
      })),
    }
    await grpcClient.put(`/v1/events/${eventId}/risky-events`, payload)
  }

  async setMatchEventOutcome(request: SetMatchRiskyEventOutcomeRequest): Promise<void> {
    const { eventId, riskyEventTypeId, outcome } = request
    await grpcClient.put(
      `/v1/events/${eventId}/risky-events/${riskyEventTypeId}/outcome`,
      { outcome }
    )
  }
}

export const riskyEventsService = new RiskyEventsService()

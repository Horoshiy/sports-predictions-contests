import type { PaginationRequest, PaginationResponse, ApiResponse } from './common.types'

export interface PropType {
  id: number
  sportType: string
  name: string
  slug: string
  description: string
  category: 'match' | 'player' | 'team'
  valueType: 'over_under' | 'yes_no' | 'team_select' | 'player_select' | 'exact_value'
  defaultLine: number | null
  minValue: number | null
  maxValue: number | null
  pointsCorrect: number
  isActive: boolean
}

export interface PropPrediction {
  propTypeId: number
  propSlug: string
  line?: number
  selection: string
  playerId?: string
  pointsValue: number
}

export interface ListPropTypesRequest {
  sportType?: string
  category?: string
  activeOnly?: boolean
  pagination?: PaginationRequest
}

export interface GetPropTypesResponse {
  response: ApiResponse
  propTypes: PropType[]
}

export interface ListPropTypesResponse {
  response: ApiResponse
  propTypes: PropType[]
  pagination: PaginationResponse
}

export interface PropPredictionFormData {
  propTypeId: number
  propSlug: string
  line?: number
  selection: string
  playerId?: string
  pointsValue: number
}

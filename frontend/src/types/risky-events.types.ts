// Risky Event Type (global)
export interface RiskyEventType {
  id: number
  slug: string
  name: string
  nameEn?: string
  description?: string
  defaultPoints: number
  sportType: string
  category: string
  icon?: string
  isActive: boolean
  sortOrder: number
  createdAt?: string
  updatedAt?: string
}

// Match Risky Event (with overrides applied)
export interface MatchRiskyEvent {
  riskyEventTypeId: number
  slug: string
  name: string
  nameEn?: string
  points: number  // final points (with overrides)
  icon?: string
  isEnabled: boolean
  outcome?: boolean | null  // null=pending, true=happened, false=didn't
}

// Contest Risky Event (selected for contest with custom points)
export interface ContestRiskyEvent {
  riskyEventTypeId: number
  slug: string
  name: string
  nameEn?: string
  points: number
  icon?: string
  enabled?: boolean
}

// Request types
export interface ListRiskyEventTypesRequest {
  sportType?: string
  activeOnly?: boolean
}

export interface CreateRiskyEventTypeRequest {
  slug: string
  name: string
  nameEn?: string
  description?: string
  defaultPoints: number
  sportType?: string
  category?: string
  icon?: string
  sortOrder?: number
}

export interface UpdateRiskyEventTypeRequest {
  id: number
  slug?: string
  name?: string
  nameEn?: string
  description?: string
  defaultPoints?: number
  sportType?: string
  category?: string
  icon?: string
  sortOrder?: number
  isActive?: boolean
}

export interface GetMatchRiskyEventsRequest {
  eventId: number
  contestId: number
}

export interface SetMatchRiskyEventsRequest {
  eventId: number
  events: {
    riskyEventTypeId: number
    points: number
    isEnabled: boolean
  }[]
}

export interface SetMatchRiskyEventOutcomeRequest {
  eventId: number
  riskyEventTypeId: number
  outcome: boolean
}

// Response types
export interface RiskyEventTypeResponse {
  response?: {
    success: boolean
    message: string
  }
  eventType: RiskyEventType
}

export interface ListRiskyEventTypesResponse {
  response?: {
    success: boolean
    message: string
  }
  eventTypes: RiskyEventType[]
}

export interface GetMatchRiskyEventsResponse {
  response?: {
    success: boolean
    message: string
  }
  events: MatchRiskyEvent[]
  maxSelections: number
}

// Categories for risky events
export const RISKY_EVENT_CATEGORIES = [
  { value: 'goals', label: 'Голы', icon: '⚽' },
  { value: 'cards', label: 'Карточки', icon: '🟥' },
  { value: 'defense', label: 'Оборона', icon: '🛡️' },
  { value: 'totals', label: 'Тоталы', icon: '📊' },
  { value: 'halves', label: 'Таймы', icon: '⏱️' },
  { value: 'timing', label: 'Время', icon: '⏰' },
  { value: 'special', label: 'Особые', icon: '✨' },
] as const

export type RiskyEventCategory = typeof RISKY_EVENT_CATEGORIES[number]['value']

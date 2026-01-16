import type { ApiResponse, PaginationRequest, PaginationResponse } from './common.types'

// Entity types matching sports.proto
export interface Sport {
  id: number
  name: string
  slug: string
  description: string
  iconUrl: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface League {
  id: number
  sportId: number
  name: string
  slug: string
  country: string
  season: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Team {
  id: number
  sportId: number
  name: string
  slug: string
  shortName: string
  logoUrl: string
  country: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export type MatchStatus = 'scheduled' | 'live' | 'finished' | 'cancelled' | 'postponed'

export interface Match {
  id: number
  leagueId: number
  homeTeamId: number
  awayTeamId: number
  scheduledAt: string
  status: MatchStatus
  homeScore: number
  awayScore: number
  resultData: string
  createdAt: string
  updatedAt: string
}

// Sport requests
export interface CreateSportRequest {
  name: string
  slug?: string
  description?: string
  iconUrl?: string
}

export interface UpdateSportRequest {
  id: number
  name: string
  slug: string
  description?: string
  iconUrl?: string
  isActive: boolean
}

export interface ListSportsRequest {
  pagination?: PaginationRequest
  activeOnly?: boolean
}

// League requests
export interface CreateLeagueRequest {
  sportId: number
  name: string
  slug?: string
  country?: string
  season?: string
}

export interface UpdateLeagueRequest {
  id: number
  sportId: number
  name: string
  slug: string
  country?: string
  season?: string
  isActive: boolean
}

export interface ListLeaguesRequest {
  pagination?: PaginationRequest
  sportId?: number
  activeOnly?: boolean
}

// Team requests
export interface CreateTeamRequest {
  sportId: number
  name: string
  slug?: string
  shortName?: string
  logoUrl?: string
  country?: string
}

export interface UpdateTeamRequest {
  id: number
  sportId: number
  name: string
  slug: string
  shortName?: string
  logoUrl?: string
  country?: string
  isActive: boolean
}

export interface ListTeamsRequest {
  pagination?: PaginationRequest
  sportId?: number
  activeOnly?: boolean
}

// Match requests
export interface CreateMatchRequest {
  leagueId: number
  homeTeamId: number
  awayTeamId: number
  scheduledAt: string
}

export interface UpdateMatchRequest {
  id: number
  leagueId: number
  homeTeamId: number
  awayTeamId: number
  scheduledAt: string
  status: MatchStatus
  homeScore?: number
  awayScore?: number
  resultData?: string
}

export interface ListMatchesRequest {
  pagination?: PaginationRequest
  leagueId?: number
  teamId?: number
  status?: MatchStatus
}

// Response types
export interface SportResponse {
  response: ApiResponse
  sport: Sport
}

export interface ListSportsResponse {
  response: ApiResponse
  sports: Sport[]
  pagination: PaginationResponse
}

export interface LeagueResponse {
  response: ApiResponse
  league: League
}

export interface ListLeaguesResponse {
  response: ApiResponse
  leagues: League[]
  pagination: PaginationResponse
}

export interface TeamResponse {
  response: ApiResponse
  team: Team
}

export interface ListTeamsResponse {
  response: ApiResponse
  teams: Team[]
  pagination: PaginationResponse
}

export interface MatchResponse {
  response: ApiResponse
  match: Match
}

export interface ListMatchesResponse {
  response: ApiResponse
  matches: Match[]
  pagination: PaginationResponse
}

export interface DeleteResponse {
  response: ApiResponse
}

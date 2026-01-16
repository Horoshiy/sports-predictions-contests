import type { PaginationRequest, PaginationResponse } from './common.types'

export interface Team {
  id: number
  name: string
  description: string
  inviteCode: string
  captainId: number
  maxMembers: number
  currentMembers: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface TeamMember {
  id: number
  teamId: number
  userId: number
  userName: string
  role: 'captain' | 'member'
  status: 'active' | 'inactive'
  joinedAt: string
}

export interface TeamLeaderboardEntry {
  teamId: number
  teamName: string
  totalPoints: number
  rank: number
  memberCount: number
}

export interface CreateTeamRequest {
  name: string
  description: string
  maxMembers: number
}

export interface UpdateTeamRequest {
  id: number
  name: string
  description: string
  maxMembers: number
}

export interface ListTeamsRequest {
  pagination?: PaginationRequest
  myTeamsOnly?: boolean
}

export interface JoinTeamRequest {
  inviteCode: string
}

export interface ListMembersRequest {
  teamId: number
  pagination?: PaginationRequest
}

export interface TeamFormData {
  name: string
  description: string
  maxMembers: number
}

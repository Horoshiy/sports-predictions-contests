import { z } from 'zod'

export const teamSchema = z.object({
  name: z.string().min(1, 'Team name is required').max(100, 'Team name cannot exceed 100 characters'),
  description: z.string().max(500, 'Description cannot exceed 500 characters').optional().default(''),
  maxMembers: z.number().min(2, 'Team must allow at least 2 members').max(50, 'Team cannot exceed 50 members').default(10),
})

export const joinTeamSchema = z.object({
  inviteCode: z.string().min(1, 'Invite code is required').max(20, 'Invalid invite code'),
})

export type TeamSchemaType = z.infer<typeof teamSchema>
export type JoinTeamSchemaType = z.infer<typeof joinTeamSchema>

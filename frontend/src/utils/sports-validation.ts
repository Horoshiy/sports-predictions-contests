import { z } from 'zod'

const slugRegex = /^[a-z0-9]+(?:-[a-z0-9]+)*$/

export const sportSchema = z.object({
  name: z.string().min(1, 'Name is required').max(100, 'Name cannot exceed 100 characters').trim(),
  slug: z.string().optional().or(z.literal('')).refine(
    (val) => !val || slugRegex.test(val),
    'Slug must be lowercase alphanumeric with hyphens'
  ),
  description: z.string().max(500, 'Description cannot exceed 500 characters').optional().or(z.literal('')),
  iconUrl: z.string().optional().or(z.literal('')).refine(
    (val) => !val || /^https?:\/\/.+/.test(val),
    'Must be a valid URL'
  ),
})

export type SportFormData = z.infer<typeof sportSchema>

export const leagueSchema = z.object({
  sportId: z.number().min(1, 'Sport is required'),
  name: z.string().min(1, 'Name is required').max(200, 'Name cannot exceed 200 characters').trim(),
  slug: z.string().optional().or(z.literal('')).refine(
    (val) => !val || slugRegex.test(val),
    'Slug must be lowercase alphanumeric with hyphens'
  ),
  country: z.string().max(100, 'Country cannot exceed 100 characters').optional().or(z.literal('')),
  season: z.string().max(50, 'Season cannot exceed 50 characters').optional().or(z.literal('')),
})

export type LeagueFormData = z.infer<typeof leagueSchema>

export const teamSchema = z.object({
  sportId: z.number().min(1, 'Sport is required'),
  name: z.string().min(1, 'Name is required').max(200, 'Name cannot exceed 200 characters').trim(),
  slug: z.string().optional().or(z.literal('')).refine(
    (val) => !val || slugRegex.test(val),
    'Slug must be lowercase alphanumeric with hyphens'
  ),
  shortName: z.string().max(50, 'Short name cannot exceed 50 characters').optional().or(z.literal('')),
  logoUrl: z.string().optional().or(z.literal('')).refine(
    (val) => !val || /^https?:\/\/.+/.test(val),
    'Must be a valid URL'
  ),
  country: z.string().max(100, 'Country cannot exceed 100 characters').optional().or(z.literal('')),
})

export type TeamFormData = z.infer<typeof teamSchema>

export const matchSchema = z.object({
  leagueId: z.number().min(1, 'League is required'),
  homeTeamId: z.number().min(1, 'Home team is required'),
  awayTeamId: z.number().min(1, 'Away team is required'),
  scheduledAt: z.date(),
  status: z.enum(['scheduled', 'live', 'finished', 'cancelled', 'postponed']).optional(),
  homeScore: z.number().min(0).optional(),
  awayScore: z.number().min(0).optional(),
  resultData: z.string().optional().or(z.literal('')),
}).refine(data => data.homeTeamId !== data.awayTeamId, {
  message: 'Home and away teams must be different',
  path: ['awayTeamId'],
})

export type MatchFormData = z.infer<typeof matchSchema>

export const generateSlug = (name: string): string => {
  return name.toLowerCase().trim().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')
}

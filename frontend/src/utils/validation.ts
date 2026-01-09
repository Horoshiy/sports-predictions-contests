import { z } from 'zod'

// Contest validation schema matching backend constraints
export const contestSchema = z.object({
  title: z
    .string()
    .min(1, 'Title is required')
    .max(200, 'Title cannot exceed 200 characters')
    .trim(),
  description: z
    .string()
    .max(1000, 'Description cannot exceed 1000 characters')
    .optional()
    .or(z.literal('')),
  sportType: z
    .string()
    .min(1, 'Sport type is required')
    .trim(),
  rules: z
    .string()
    .optional()
    .or(z.literal('')),
  startDate: z
    .date()
    .refine(date => date > new Date(), {
      message: 'Start date must be in the future',
    }),
  endDate: z
    .date(),
  maxParticipants: z
    .number()
    .min(0, 'Max participants must be 0 or greater')
    .max(10000, 'Max participants cannot exceed 10,000'),
}).refine(
  (data) => data.endDate > data.startDate,
  {
    message: 'End date must be after start date',
    path: ['endDate'],
  }
)

export type ContestFormData = z.infer<typeof contestSchema>

// Participant validation
export const participantSchema = z.object({
  userId: z.number().min(1, 'User ID is required'),
  role: z.enum(['admin', 'participant'], {
    errorMap: () => ({ message: 'Role must be admin or participant' }),
  }),
})

export type ParticipantFormData = z.infer<typeof participantSchema>

// Search and filter validation
export const contestFiltersSchema = z.object({
  status: z.enum(['draft', 'active', 'completed', 'cancelled']).optional(),
  sportType: z.string().optional(),
  search: z.string().optional(),
  page: z.number().min(1).default(1),
  limit: z.number().min(1).max(100).default(10),
  sortBy: z.string().optional(),
  sortOrder: z.enum(['asc', 'desc']).default('desc'),
})

export type ContestFiltersData = z.infer<typeof contestFiltersSchema>

// Validation helper functions
export const validateRequired = (value: string): boolean => {
  return value.trim().length > 0
}

export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const validateDateRange = (startDate: Date, endDate: Date): boolean => {
  return endDate > startDate
}

export const validateFutureDate = (date: Date): boolean => {
  return date > new Date()
}

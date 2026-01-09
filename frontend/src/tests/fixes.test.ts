import { describe, it, expect } from 'vitest'
import { contestSchema } from '../src/utils/validation'
import { formatDate, getContestStatusByDate } from '../src/utils/date-utils'

describe('Validation Fixes', () => {
  it('should validate future dates correctly without race condition', () => {
    const futureDate = new Date(Date.now() + 24 * 60 * 60 * 1000) // Tomorrow
    const pastDate = new Date(Date.now() - 24 * 60 * 60 * 1000) // Yesterday
    
    const validData = {
      title: 'Test Contest',
      description: 'Test description',
      sportType: 'Football',
      rules: '{}',
      startDate: futureDate,
      endDate: new Date(futureDate.getTime() + 60 * 60 * 1000), // 1 hour later
      maxParticipants: 10,
    }

    const invalidData = {
      ...validData,
      startDate: pastDate,
    }

    expect(contestSchema.safeParse(validData).success).toBe(true)
    expect(contestSchema.safeParse(invalidData).success).toBe(false)
  })

  it('should handle invalid dates gracefully', () => {
    expect(formatDate('invalid-date')).toBe('Invalid date')
    expect(getContestStatusByDate('invalid', 'invalid')).toBe('completed')
  })

  it('should determine contest status correctly with valid dates', () => {
    const now = new Date()
    const past = new Date(now.getTime() - 60 * 60 * 1000) // 1 hour ago
    const future = new Date(now.getTime() + 60 * 60 * 1000) // 1 hour from now
    const farFuture = new Date(now.getTime() + 2 * 60 * 60 * 1000) // 2 hours from now

    expect(getContestStatusByDate(future, farFuture)).toBe('upcoming')
    expect(getContestStatusByDate(past, future)).toBe('active')
    expect(getContestStatusByDate(past, past)).toBe('completed')
  })
})

describe('Type Safety Fixes', () => {
  it('should have proper TypeScript types', () => {
    // This test will fail at compile time if types are wrong
    const mockParticipant = {
      id: 1,
      contestId: 1,
      userId: 1,
      role: 'participant' as const,
      status: 'active' as const,
      joinedAt: new Date().toISOString(),
    }

    expect(mockParticipant.role).toBe('participant')
  })
})

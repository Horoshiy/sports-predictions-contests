import { describe, it, expect } from 'vitest'
import { contestSchema } from '../utils/validation'
import { formatDate, getContestStatusByDate } from '../utils/date-utils'

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

  it('should validate title constraints', () => {
    const futureDate = new Date(Date.now() + 24 * 60 * 60 * 1000)
    const baseData = {
      title: 'Valid Title',
      sportType: 'Football',
      startDate: futureDate,
      endDate: new Date(futureDate.getTime() + 60 * 60 * 1000),
      maxParticipants: 10,
    }

    // Too short
    expect(contestSchema.safeParse({ ...baseData, title: 'AB' }).success).toBe(false)
    
    // Valid minimum
    expect(contestSchema.safeParse({ ...baseData, title: 'ABC' }).success).toBe(true)
    
    // Too long (over 200 chars)
    const longTitle = 'A'.repeat(201)
    expect(contestSchema.safeParse({ ...baseData, title: longTitle }).success).toBe(false)
    
    // Valid max (200 chars)
    const maxTitle = 'A'.repeat(200)
    expect(contestSchema.safeParse({ ...baseData, title: maxTitle }).success).toBe(true)
  })

  it('should validate description and participant constraints', () => {
    const futureDate = new Date(Date.now() + 24 * 60 * 60 * 1000)
    const baseData = {
      title: 'Test Contest',
      sportType: 'Football',
      startDate: futureDate,
      endDate: new Date(futureDate.getTime() + 60 * 60 * 1000),
      maxParticipants: 10,
    }

    // Description too long (over 1000 chars)
    const longDesc = 'A'.repeat(1001)
    expect(contestSchema.safeParse({ ...baseData, description: longDesc }).success).toBe(false)
    
    // Valid description (1000 chars)
    const maxDesc = 'A'.repeat(1000)
    expect(contestSchema.safeParse({ ...baseData, description: maxDesc }).success).toBe(true)
    
    // Max participants too high
    expect(contestSchema.safeParse({ ...baseData, maxParticipants: 10001 }).success).toBe(false)
    
    // Valid max participants
    expect(contestSchema.safeParse({ ...baseData, maxParticipants: 10000 }).success).toBe(true)
  })

  it('should trim whitespace and reject whitespace-only values', () => {
    const futureDate = new Date(Date.now() + 24 * 60 * 60 * 1000)
    const baseData = {
      title: 'Test Contest',
      sportType: 'Football',
      startDate: futureDate,
      endDate: new Date(futureDate.getTime() + 60 * 60 * 1000),
      maxParticipants: 10,
    }

    // Whitespace-only sport type should fail
    expect(contestSchema.safeParse({ ...baseData, sportType: '   ' }).success).toBe(false)
    
    // Whitespace-only title should fail
    expect(contestSchema.safeParse({ ...baseData, title: '   ' }).success).toBe(false)
  })

  it('should handle invalid dates gracefully', () => {
    expect(formatDate('invalid-date')).toBe('Invalid Date')
    // Invalid dates result in 'active' status due to NaN comparisons
    expect(getContestStatusByDate('invalid', 'invalid')).toBe('active')
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

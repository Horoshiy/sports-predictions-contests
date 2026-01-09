import { format, parseISO, isValid, formatDistanceToNow } from 'date-fns'

// Format date for display
export const formatDate = (date: string | Date): string => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date
    if (!isValid(dateObj)) {
      return 'Invalid date'
    }
    return format(dateObj, 'MMM dd, yyyy')
  } catch (error) {
    console.error('Error formatting date:', error)
    return 'Invalid date'
  }
}

// Format datetime for display
export const formatDateTime = (date: string | Date): string => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date
    if (!isValid(dateObj)) {
      return 'Invalid date'
    }
    return format(dateObj, 'MMM dd, yyyy HH:mm')
  } catch (error) {
    console.error('Error formatting datetime:', error)
    return 'Invalid date'
  }
}

// Format date for form inputs (YYYY-MM-DD)
export const formatDateForInput = (date: string | Date): string => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date
    if (!isValid(dateObj)) {
      return ''
    }
    return format(dateObj, 'yyyy-MM-dd')
  } catch (error) {
    console.error('Error formatting date for input:', error)
    return ''
  }
}

// Format datetime for form inputs (YYYY-MM-DDTHH:mm)
export const formatDateTimeForInput = (date: string | Date): string => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date
    if (!isValid(dateObj)) {
      return ''
    }
    return format(dateObj, "yyyy-MM-dd'T'HH:mm")
  } catch (error) {
    console.error('Error formatting datetime for input:', error)
    return ''
  }
}

// Format relative time (e.g., "2 hours ago", "in 3 days")
export const formatRelativeTime = (date: string | Date): string => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date
    if (!isValid(dateObj)) {
      return 'Invalid date'
    }
    return formatDistanceToNow(dateObj, { addSuffix: true })
  } catch (error) {
    console.error('Error formatting relative time:', error)
    return 'Invalid date'
  }
}

// Convert Date to ISO string for API
export const toISOString = (date: Date): string => {
  try {
    if (!isValid(date)) {
      throw new Error('Invalid date')
    }
    return date.toISOString()
  } catch (error) {
    console.error('Error converting to ISO string:', error)
    throw error
  }
}

// Parse ISO string to Date
export const fromISOString = (isoString: string): Date => {
  try {
    const date = parseISO(isoString)
    if (!isValid(date)) {
      throw new Error('Invalid ISO string')
    }
    return date
  } catch (error) {
    console.error('Error parsing ISO string:', error)
    throw error
  }
}

// Check if date is in the past
export const isPastDate = (date: string | Date): boolean => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date
    if (!isValid(dateObj)) {
      return false
    }
    return dateObj < new Date()
  } catch (error) {
    console.error('Error checking if date is past:', error)
    return false
  }
}

// Check if date is in the future
export const isFutureDate = (date: string | Date): boolean => {
  try {
    const dateObj = typeof date === 'string' ? parseISO(date) : date
    if (!isValid(dateObj)) {
      return false
    }
    return dateObj > new Date()
  } catch (error) {
    console.error('Error checking if date is future:', error)
    return false
  }
}

// Get contest status based on dates
export const getContestStatusByDate = (startDate: string | Date, endDate: string | Date): 'upcoming' | 'active' | 'completed' => {
  const now = new Date()
  const start = typeof startDate === 'string' ? parseISO(startDate) : startDate
  const end = typeof endDate === 'string' ? parseISO(endDate) : endDate

  // Validate dates
  if (!isValid(start) || !isValid(end)) {
    return 'completed'
  }

  if (now < start) {
    return 'upcoming'
  } else if (now >= start && now <= end) {
    return 'active'
  } else {
    return 'completed'
  }
}

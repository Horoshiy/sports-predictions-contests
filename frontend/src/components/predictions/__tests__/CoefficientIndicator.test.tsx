import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import '@testing-library/jest-dom'
import CoefficientIndicator from '../CoefficientIndicator'

describe('CoefficientIndicator', () => {
  it('displays coefficient value in full mode', () => {
    render(
      <CoefficientIndicator 
        coefficient={2.0} 
        tier="Early Bird" 
        hoursUntilEvent={200} 
      />
    )
    expect(screen.getByText('2x Points')).toBeInTheDocument()
    expect(screen.getByText('Early Bird')).toBeInTheDocument()
  })

  it('displays coefficient value in compact mode', () => {
    render(
      <CoefficientIndicator 
        coefficient={1.5} 
        tier="Ahead of Time" 
        hoursUntilEvent={100} 
        compact 
      />
    )
    expect(screen.getByText('1.5x')).toBeInTheDocument()
  })

  it('formats time remaining correctly for days', () => {
    render(
      <CoefficientIndicator 
        coefficient={2.0} 
        tier="Early Bird" 
        hoursUntilEvent={200} 
      />
    )
    expect(screen.getByText(/8 days/)).toBeInTheDocument()
  })

  it('formats time remaining correctly for hours', () => {
    render(
      <CoefficientIndicator 
        coefficient={1.1} 
        tier="Last Minute" 
        hoursUntilEvent={18} 
      />
    )
    expect(screen.getByText(/18h/)).toBeInTheDocument()
  })

  it('shows correct color for high coefficient', () => {
    const { container } = render(
      <CoefficientIndicator 
        coefficient={2.0} 
        tier="Early Bird" 
        hoursUntilEvent={200} 
      />
    )
    const chip = container.querySelector('.MuiChip-colorSuccess')
    expect(chip).toBeInTheDocument()
  })
})

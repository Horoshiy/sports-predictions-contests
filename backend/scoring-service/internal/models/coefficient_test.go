package models

import (
	"testing"
	"time"
)

func TestCalculateTimeCoefficient(t *testing.T) {
	// Event 8 days from now
	eventDate := time.Now().Add(8 * 24 * time.Hour)

	tests := []struct {
		name        string
		submittedAt time.Time
		expected    float64
	}{
		{
			name:        "7+ days early should return 2.0x",
			submittedAt: time.Now(),
			expected:    2.0,
		},
		{
			name:        "5 days early should return 1.5x",
			submittedAt: time.Now().Add(3 * 24 * time.Hour),
			expected:    1.5,
		},
		{
			name:        "2 days early should return 1.25x",
			submittedAt: time.Now().Add(6 * 24 * time.Hour),
			expected:    1.25,
		},
		{
			name:        "18 hours early should return 1.1x",
			submittedAt: time.Now().Add(7*24*time.Hour + 6*time.Hour),
			expected:    1.1,
		},
		{
			name:        "6 hours early should return 1.0x",
			submittedAt: time.Now().Add(7*24*time.Hour + 18*time.Hour),
			expected:    1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTimeCoefficient(tt.submittedAt, eventDate)
			if result != tt.expected {
				t.Errorf("CalculateTimeCoefficient() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateTimeCoefficient_EdgeCases(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		submittedAt time.Time
		eventDate   time.Time
		expected    float64
	}{
		{
			name:        "event in past should return 1.0x",
			submittedAt: now,
			eventDate:   now.Add(-24 * time.Hour),
			expected:    1.0,
		},
		{
			name:        "exactly 168 hours should return 2.0x",
			submittedAt: now,
			eventDate:   now.Add(168 * time.Hour),
			expected:    2.0,
		},
		{
			name:        "exactly 72 hours should return 1.5x",
			submittedAt: now,
			eventDate:   now.Add(72 * time.Hour),
			expected:    1.5,
		},
		{
			name:        "exactly 24 hours should return 1.25x",
			submittedAt: now,
			eventDate:   now.Add(24 * time.Hour),
			expected:    1.25,
		},
		{
			name:        "exactly 12 hours should return 1.1x",
			submittedAt: now,
			eventDate:   now.Add(12 * time.Hour),
			expected:    1.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTimeCoefficient(tt.submittedAt, tt.eventDate)
			if result != tt.expected {
				t.Errorf("CalculateTimeCoefficient() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateWithTier(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		hoursAhead    float64
		expectedCoeff float64
		expectedTier  string
	}{
		{"7+ days", 200, 2.0, "Early Bird"},
		{"3-7 days", 100, 1.5, "Ahead of Time"},
		{"1-3 days", 48, 1.25, "Timely"},
		{"12-24 hours", 18, 1.1, "Last Minute"},
		{"<12 hours", 6, 1.0, "Standard"},
		{"past event", -24, 1.0, "Standard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventDate := now.Add(time.Duration(tt.hoursAhead) * time.Hour)
			coeff, tier := CalculateWithTier(now, eventDate)
			if coeff != tt.expectedCoeff {
				t.Errorf("CalculateWithTier() coeff = %v, want %v", coeff, tt.expectedCoeff)
			}
			if tier != tt.expectedTier {
				t.Errorf("CalculateWithTier() tier = %v, want %v", tier, tt.expectedTier)
			}
		})
	}
}

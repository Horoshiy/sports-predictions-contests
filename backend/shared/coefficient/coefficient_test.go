package coefficient

import (
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name             string
		hoursAhead       float64
		expectedCoeff    float64
		expectedTier     string
	}{
		{"7+ days early", 200, 2.0, "Early Bird"},
		{"5 days early", 120, 1.5, "Ahead of Time"},
		{"2 days early", 48, 1.25, "Timely"},
		{"18 hours early", 18, 1.1, "Last Minute"},
		{"6 hours early", 6, 1.0, "Standard"},
		{"event in past", -24, 1.0, "Standard"},
		{"exactly 168 hours", 168, 2.0, "Early Bird"},
		{"exactly 72 hours", 72, 1.5, "Ahead of Time"},
		{"exactly 24 hours", 24, 1.25, "Timely"},
		{"exactly 12 hours", 12, 1.1, "Last Minute"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventDate := now.Add(time.Duration(tt.hoursAhead * float64(time.Hour)))
			result := Calculate(now, eventDate)
			
			if result.Coefficient != tt.expectedCoeff {
				t.Errorf("Calculate().Coefficient = %v, want %v", result.Coefficient, tt.expectedCoeff)
			}
			if result.Tier != tt.expectedTier {
				t.Errorf("Calculate().Tier = %v, want %v", result.Tier, tt.expectedTier)
			}
		})
	}
}

package coefficient

import "time"

// CoefficientResult contains both the multiplier and tier name
type CoefficientResult struct {
	Coefficient float64
	Tier        string
}

// Calculate returns point multiplier and tier based on prediction timing
// Earlier predictions relative to event start earn higher multipliers
func Calculate(submittedAt, eventDate time.Time) CoefficientResult {
	hoursUntilEvent := eventDate.Sub(submittedAt).Hours()

	switch {
	case hoursUntilEvent < 0:
		return CoefficientResult{1.0, "Standard"} // Event already started
	case hoursUntilEvent >= 168: // 7+ days
		return CoefficientResult{2.0, "Early Bird"}
	case hoursUntilEvent >= 72: // 3-7 days
		return CoefficientResult{1.5, "Ahead of Time"}
	case hoursUntilEvent >= 24: // 1-3 days
		return CoefficientResult{1.25, "Timely"}
	case hoursUntilEvent >= 12: // 12-24 hours
		return CoefficientResult{1.1, "Last Minute"}
	default:
		return CoefficientResult{1.0, "Standard"}
	}
}

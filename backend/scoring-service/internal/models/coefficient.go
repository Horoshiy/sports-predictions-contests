package models

import (
	"time"

	"github.com/sports-prediction-contests/shared/coefficient"
)

// CalculateTimeCoefficient returns point multiplier based on prediction timing
func CalculateTimeCoefficient(submittedAt, eventDate time.Time) float64 {
	return coefficient.Calculate(submittedAt, eventDate).Coefficient
}

// CalculateWithTier returns both coefficient and tier in one call
func CalculateWithTier(submittedAt, eventDate time.Time) (float64, string) {
	result := coefficient.Calculate(submittedAt, eventDate)
	return result.Coefficient, result.Tier
}

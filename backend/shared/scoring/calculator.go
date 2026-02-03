package scoring

import (
	"fmt"
)

// ScoreData represents prediction or result score
type ScoreData struct {
	HomeScore int
	AwayScore int
}

// Calculator calculates points based on contest rules
type Calculator struct {
	rules *ContestRules
}

// NewCalculator creates a new calculator with the given rules
func NewCalculator(rules *ContestRules) *Calculator {
	return &Calculator{rules: rules}
}

// CalculationResult contains the result of score calculation
type CalculationResult struct {
	Points  float64
	Details map[string]interface{}
}

// CalculateStandard calculates points for a standard prediction
func (c *Calculator) CalculateStandard(prediction, result ScoreData, isAnyOther bool) CalculationResult {
	details := map[string]interface{}{
		"type":            "standard",
		"predicted_score": fmt.Sprintf("%d:%d", prediction.HomeScore, prediction.AwayScore),
		"actual_score":    fmt.Sprintf("%d:%d", result.HomeScore, result.AwayScore),
		"is_any_other":    isAnyOther,
	}

	// Nil safety check
	if c.rules == nil {
		details["error"] = "nil rules"
		return CalculationResult{Points: 0, Details: details}
	}

	scoring := c.rules.Standard
	if scoring == nil {
		defaultRules := DefaultStandardRules()
		scoring = &defaultRules
	}

	// Handle "any other" prediction
	if isAnyOther {
		// Check if result is NOT one of the common scores (0:0 to 4:4)
		isOther := c.isOtherScore(result.HomeScore, result.AwayScore)
		details["result_is_other"] = isOther
		if isOther {
			details["match_type"] = "any_other_correct"
			return CalculationResult{Points: scoring.AnyOther, Details: details}
		}
		details["match_type"] = "any_other_incorrect"
		return CalculationResult{Points: 0, Details: details}
	}

	// Exact match
	if prediction.HomeScore == result.HomeScore && prediction.AwayScore == result.AwayScore {
		details["match_type"] = "exact_score"
		return CalculationResult{Points: scoring.ExactScore, Details: details}
	}

	// Calculate outcomes
	predictedOutcome := c.determineOutcome(prediction.HomeScore, prediction.AwayScore)
	actualOutcome := c.determineOutcome(result.HomeScore, result.AwayScore)
	details["predicted_outcome"] = predictedOutcome
	details["actual_outcome"] = actualOutcome

	// Goal difference match
	predictedDiff := prediction.HomeScore - prediction.AwayScore
	actualDiff := result.HomeScore - result.AwayScore
	if predictedDiff == actualDiff {
		details["match_type"] = "goal_difference"
		return CalculationResult{Points: scoring.GoalDifference, Details: details}
	}

	// Correct outcome with one team's goals correct
	if predictedOutcome == actualOutcome {
		// Check if home or away goals match
		homeGoalsMatch := prediction.HomeScore == result.HomeScore
		awayGoalsMatch := prediction.AwayScore == result.AwayScore

		if homeGoalsMatch || awayGoalsMatch {
			details["match_type"] = "outcome_plus_team_goals"
			details["home_goals_match"] = homeGoalsMatch
			details["away_goals_match"] = awayGoalsMatch
			points := scoring.CorrectOutcome + scoring.OutcomePlusTeamGoals
			return CalculationResult{Points: points, Details: details}
		}

		// Just correct outcome
		details["match_type"] = "correct_outcome"
		return CalculationResult{Points: scoring.CorrectOutcome, Details: details}
	}

	details["match_type"] = "none"
	return CalculationResult{Points: 0, Details: details}
}

// CalculateRisky calculates points for risky predictions
// selections: list of event slugs the user selected
// outcomes: map of event slug -> whether it occurred (true/false)
func (c *Calculator) CalculateRisky(selections []string, outcomes map[string]bool) CalculationResult {
	details := map[string]interface{}{
		"type":       "risky",
		"selections": selections,
		"outcomes":   outcomes,
	}

	// Nil safety checks
	if c.rules == nil || c.rules.Risky == nil {
		details["error"] = "no risky rules configured"
		return CalculationResult{Points: 0, Details: details}
	}

	var totalPoints float64
	eventResults := make([]map[string]interface{}, 0)

	for _, slug := range selections {
		event := c.rules.Risky.GetEventBySlug(slug)
		if event == nil {
			continue
		}

		occurred, hasOutcome := outcomes[slug]
		if !hasOutcome {
			// Event outcome not yet determined
			continue
		}

		eventResult := map[string]interface{}{
			"slug":     slug,
			"name":     event.Name,
			"points":   event.Points,
			"occurred": occurred,
		}

		if occurred {
			// User guessed correctly: +points
			totalPoints += event.Points
			eventResult["earned"] = event.Points
		} else {
			// User guessed wrong: -points
			totalPoints -= event.Points
			eventResult["earned"] = -event.Points
		}

		eventResults = append(eventResults, eventResult)
	}

	details["event_results"] = eventResults
	details["total_points"] = totalPoints

	return CalculationResult{Points: totalPoints, Details: details}
}

// determineOutcome returns "home", "away", or "draw"
func (c *Calculator) determineOutcome(homeScore, awayScore int) string {
	if homeScore > awayScore {
		return "home"
	}
	if awayScore > homeScore {
		return "away"
	}
	return "draw"
}

// isOtherScore checks if score is outside common range (for "any other" predictions)
// Common scores are typically 0-4 for each team
func (c *Calculator) isOtherScore(homeScore, awayScore int) bool {
	return homeScore > 4 || awayScore > 4
}

// ValidateRiskySelections checks if selections are valid
func (c *Calculator) ValidateRiskySelections(selections []string) error {
	if c.rules.Risky == nil {
		return fmt.Errorf("contest is not risky type")
	}

	if len(selections) > c.rules.Risky.MaxSelections {
		return fmt.Errorf("too many selections: max %d allowed, got %d",
			c.rules.Risky.MaxSelections, len(selections))
	}

	// Check all selections are valid events
	for _, slug := range selections {
		if c.rules.Risky.GetEventBySlug(slug) == nil {
			return fmt.Errorf("unknown event: %s", slug)
		}
	}

	return nil
}

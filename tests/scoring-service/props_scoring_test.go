package scoring_service_test

import (
	"encoding/json"
	"testing"
)

type PropPrediction struct {
	PropTypeID  uint    `json:"prop_type_id"`
	PropSlug    string  `json:"prop_slug"`
	Line        float64 `json:"line"`
	Selection   string  `json:"selection"`
	PointsValue float64 `json:"points_value"`
}

func TestPropsScoring_TotalGoalsOverUnder(t *testing.T) {
	tests := []struct {
		name        string
		line        float64
		selection   string
		totalGoals  int
		wantCorrect bool
	}{
		{"over 2.5 with 3 goals", 2.5, "over", 3, true},
		{"over 2.5 with 2 goals", 2.5, "over", 2, false},
		{"under 2.5 with 2 goals", 2.5, "under", 2, true},
		{"under 2.5 with 3 goals", 2.5, "under", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			correct := evaluateTotalGoalsOU(tt.line, tt.selection, tt.totalGoals)
			if correct != tt.wantCorrect {
				t.Errorf("evaluateTotalGoalsOU() = %v, want %v", correct, tt.wantCorrect)
			}
		})
	}
}

func evaluateTotalGoalsOU(line float64, selection string, totalGoals int) bool {
	if selection == "over" {
		return float64(totalGoals) > line
	}
	return float64(totalGoals) < line
}

func TestPropsScoring_BothTeamsToScore(t *testing.T) {
	tests := []struct {
		name        string
		homeScore   int
		awayScore   int
		selection   string
		wantCorrect bool
	}{
		{"yes with both scoring", 2, 1, "yes", true},
		{"yes with only home scoring", 2, 0, "yes", false},
		{"no with only home scoring", 2, 0, "no", true},
		{"no with both scoring", 1, 1, "no", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			correct := evaluateBTTS(tt.homeScore, tt.awayScore, tt.selection)
			if correct != tt.wantCorrect {
				t.Errorf("evaluateBTTS() = %v, want %v", correct, tt.wantCorrect)
			}
		})
	}
}

func evaluateBTTS(homeScore, awayScore int, selection string) bool {
	btts := homeScore > 0 && awayScore > 0
	if selection == "yes" {
		return btts
	}
	return !btts
}

func TestPropsScoring_JSONParsing(t *testing.T) {
	predictionJSON := `{
		"type": "props",
		"props": [
			{"prop_type_id": 1, "prop_slug": "total-goals-ou", "line": 2.5, "selection": "over", "points_value": 2}
		]
	}`

	var data struct {
		Type  string           `json:"type"`
		Props []PropPrediction `json:"props"`
	}

	err := json.Unmarshal([]byte(predictionJSON), &data)
	if err != nil {
		t.Fatalf("Failed to parse prediction JSON: %v", err)
	}

	if data.Type != "props" {
		t.Errorf("Expected type 'props', got '%s'", data.Type)
	}

	if len(data.Props) != 1 {
		t.Errorf("Expected 1 prop, got %d", len(data.Props))
	}

	if data.Props[0].PropSlug != "total-goals-ou" {
		t.Errorf("Expected prop slug 'total-goals-ou', got '%s'", data.Props[0].PropSlug)
	}
}

package seeder

import (
	"encoding/json"
	"testing"
)

// TestScoreOptionsMatchSchema verifies score options match the schema
func TestScoreOptionsMatchSchema(t *testing.T) {
	factory := &DataFactory{}
	schemaJSON := factory.GenerateDefaultPredictionSchema()
	
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		t.Fatalf("Failed to unmarshal schema: %v", err)
	}
	
	options, ok := schema["options"].([]interface{})
	if !ok {
		t.Fatal("Schema options is not an array")
	}
	
	// Verify "3-3" is included
	found := false
	for _, opt := range options {
		if opt == "3-3" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Score option '3-3' not found in schema")
	}
	
	// Verify we have 16 options
	if len(options) != 16 {
		t.Errorf("Expected 16 score options, got %d", len(options))
	}
}

// TestPredictionDataMarshaling verifies prediction data can be marshaled
func TestPredictionDataMarshaling(t *testing.T) {
	predictionData := map[string]interface{}{
		"home_score":   2,
		"away_score":   1,
		"score_string": "2-1",
	}
	
	jsonData, err := json.Marshal(predictionData)
	if err != nil {
		t.Errorf("Failed to marshal prediction data: %v", err)
	}
	
	// Verify it can be unmarshaled back
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Errorf("Failed to unmarshal prediction data: %v", err)
	}
	
	// Verify values
	if result["score_string"] != "2-1" {
		t.Errorf("Expected score_string '2-1', got %v", result["score_string"])
	}
}

// TestScoreParsingValidation verifies score string parsing
func TestScoreParsingValidation(t *testing.T) {
	testCases := []struct {
		score       string
		expectError bool
		homeScore   int
		awayScore   int
	}{
		{"1-0", false, 1, 0},
		{"2-1", false, 2, 1},
		{"3-3", false, 3, 3},
		{"invalid", true, 0, 0},
		{"1-", true, 0, 0},
		{"-1", true, 0, 0},
		{"", true, 0, 0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.score, func(t *testing.T) {
			parts := splitScore(tc.score)
			if len(parts) != 2 {
				if !tc.expectError {
					t.Errorf("Expected valid score, got invalid format")
				}
				return
			}
			
			home, err1 := parseInt(parts[0])
			away, err2 := parseInt(parts[1])
			
			if tc.expectError {
				if err1 == nil && err2 == nil {
					t.Errorf("Expected error for score %s, but got none", tc.score)
				}
			} else {
				if err1 != nil || err2 != nil {
					t.Errorf("Unexpected error for score %s: %v, %v", tc.score, err1, err2)
				}
				if home != tc.homeScore || away != tc.awayScore {
					t.Errorf("Expected %d-%d, got %d-%d", tc.homeScore, tc.awayScore, home, away)
				}
			}
		})
	}
}

// Helper functions for testing
func splitScore(score string) []string {
	result := []string{}
	parts := ""
	for _, c := range score {
		if c == '-' {
			result = append(result, parts)
			parts = ""
		} else {
			parts += string(c)
		}
	}
	if parts != "" {
		result = append(result, parts)
	}
	return result
}

func parseInt(s string) (int, error) {
	if s == "" {
		return 0, &parseError{"empty string"}
	}
	val := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, &parseError{"invalid character"}
		}
		val = val*10 + int(c-'0')
	}
	return val, nil
}

type parseError struct {
	msg string
}

func (e *parseError) Error() string {
	return e.msg
}

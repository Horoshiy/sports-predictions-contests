package prediction_service_test

import (
	"testing"
)

// PropType test struct matching the model
type PropType struct {
	SportType string
	Name      string
	Category  string
	ValueType string
}

func (p *PropType) ValidateSportType() error {
	if len(p.SportType) == 0 || p.SportType == "   " {
		return errEmpty
	}
	return nil
}

func (p *PropType) ValidateCategory() error {
	validCategories := []string{"match", "player", "team"}
	for _, c := range validCategories {
		if p.Category == c {
			return nil
		}
	}
	return errInvalid
}

func (p *PropType) ValidateValueType() error {
	validTypes := []string{"over_under", "yes_no", "team_select", "player_select", "exact_value"}
	for _, t := range validTypes {
		if p.ValueType == t {
			return nil
		}
	}
	return errInvalid
}

var errEmpty = &testError{"empty"}
var errInvalid = &testError{"invalid"}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func TestPropType_ValidateSportType(t *testing.T) {
	tests := []struct {
		name      string
		sportType string
		wantErr   bool
	}{
		{"valid sport type", "Soccer", false},
		{"empty sport type", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PropType{SportType: tt.sportType}
			err := pt.ValidateSportType()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSportType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPropType_ValidateCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{"valid match", "match", false},
		{"valid player", "player", false},
		{"valid team", "team", false},
		{"invalid category", "invalid", true},
		{"empty category", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PropType{Category: tt.category}
			err := pt.ValidateCategory()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPropType_ValidateValueType(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		wantErr   bool
	}{
		{"valid over_under", "over_under", false},
		{"valid yes_no", "yes_no", false},
		{"valid team_select", "team_select", false},
		{"valid player_select", "player_select", false},
		{"valid exact_value", "exact_value", false},
		{"invalid type", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PropType{ValueType: tt.valueType}
			err := pt.ValidateValueType()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateValueType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

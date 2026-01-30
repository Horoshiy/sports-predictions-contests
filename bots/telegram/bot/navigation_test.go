package bot

import (
	"testing"
)

// TestCalculatePagination tests pagination calculation with various inputs
func TestCalculatePagination(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		itemsPerPage int
		totalItems   int
		wantStart    int
		wantEnd      int
	}{
		{"normal case", 1, 5, 20, 0, 5},
		{"second page", 2, 5, 20, 5, 10},
		{"last page", 4, 5, 18, 15, 18},
		{"invalid page zero", 0, 5, 20, 0, 5},
		{"invalid page negative", -1, 5, 20, 0, 5},
		{"invalid items per page zero", 1, 0, 20, 0, 1},
		{"invalid items per page negative", 1, -5, 20, 0, 1},
		{"negative total items", 1, 5, -1, 0, 0},
		{"overflow protection - large page", 2000000, 5, 20, 0, 0},
		{"overflow protection - large items", 1, 2000, 20, 0, 0},
		{"empty list", 1, 5, 0, 0, 0},
		{"page beyond total", 10, 5, 20, 20, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := CalculatePagination(tt.page, tt.itemsPerPage, tt.totalItems)
			if gotStart != tt.wantStart || gotEnd != tt.wantEnd {
				t.Errorf("CalculatePagination(%d, %d, %d) = (%d, %d), want (%d, %d)",
					tt.page, tt.itemsPerPage, tt.totalItems, gotStart, gotEnd, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

// TestScoreValidationConstants tests that score constants are properly defined
func TestScoreValidationConstants(t *testing.T) {
	if minScore != 0 {
		t.Errorf("minScore = %d, want 0", minScore)
	}
	if maxScore != 20 {
		t.Errorf("maxScore = %d, want 20", maxScore)
	}
	if matchesPerPage != 5 {
		t.Errorf("matchesPerPage = %d, want 5", matchesPerPage)
	}
}

// TestPaginationButtonsDivisionByZero tests that PaginationButtons handles zero itemsPerPage
func TestPaginationButtonsDivisionByZero(t *testing.T) {
	tests := []struct {
		name  string
		state NavigationState
	}{
		{
			name: "zero itemsPerPage",
			state: NavigationState{
				ContestID:    1,
				Page:         1,
				ItemsPerPage: 0,
				TotalItems:   10,
			},
		},
		{
			name: "negative itemsPerPage",
			state: NavigationState{
				ContestID:    1,
				Page:         1,
				ItemsPerPage: -5,
				TotalItems:   10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PaginationButtons panicked with state %+v: %v", tt.state, r)
				}
			}()
			
			buttons := PaginationButtons(tt.state, "test")
			if len(buttons) == 0 {
				t.Error("PaginationButtons returned empty slice")
			}
		})
	}
}

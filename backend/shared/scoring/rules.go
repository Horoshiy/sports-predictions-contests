package scoring

import (
	"encoding/json"
	"errors"
)

// ContestType defines the type of contest
type ContestType string

const (
	ContestTypeStandard    ContestType = "standard"
	ContestTypeRisky       ContestType = "risky"
	ContestTypeTotalizator ContestType = "totalizator"
	ContestTypeRelay       ContestType = "relay"
)

// StandardScoringRules defines points for standard contest
type StandardScoringRules struct {
	ExactScore           float64 `json:"exact_score"`             // точный счёт
	GoalDifference       float64 `json:"goal_difference"`         // разница мячей
	CorrectOutcome       float64 `json:"correct_outcome"`         // верный исход
	OutcomePlusTeamGoals float64 `json:"outcome_plus_team_goals"` // исход + голы одной команды
	AnyOther             float64 `json:"any_other"`               // прогноз "другой"
}

// RiskyEvent defines a risky event type
type RiskyEvent struct {
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	NameEn      string  `json:"name_en,omitempty"`
	Points      float64 `json:"points"`
	Description string  `json:"description,omitempty"`
}

// RiskyScoringRules defines rules for risky contest
type RiskyScoringRules struct {
	MaxSelections int          `json:"max_selections"`
	Events        []RiskyEvent `json:"events,omitempty"`
}

// TotalizatorRules defines rules for totalizator contest
// Admin manually selects matches from different leagues
type TotalizatorRules struct {
	EventCount int                  `json:"event_count"` // number of matches (default 15)
	Scoring    StandardScoringRules `json:"scoring"`     // scoring rules for all matches
}

// RelayRules defines rules for relay (team) contest
// Captain assigns matches to team members, team score is sum of all members
type RelayRules struct {
	TeamSize      int                  `json:"team_size"`       // members per team (2-10)
	EventCount    int                  `json:"event_count"`     // total matches in contest (5-50)
	Scoring       StandardScoringRules `json:"scoring"`         // scoring rules
	AllowReassign bool                 `json:"allow_reassign"`  // can captain reassign after start
}

// ContestRules combines all rule types
type ContestRules struct {
	Type        ContestType           `json:"type"`
	Standard    *StandardScoringRules `json:"scoring,omitempty"`
	Risky       *RiskyScoringRules    `json:"risky,omitempty"`
	Totalizator *TotalizatorRules     `json:"totalizator,omitempty"`
	Relay       *RelayRules           `json:"relay,omitempty"`
}

// DefaultStandardRules returns default scoring for standard contests
func DefaultStandardRules() StandardScoringRules {
	return StandardScoringRules{
		ExactScore:           5,
		GoalDifference:       3,
		CorrectOutcome:       1,
		OutcomePlusTeamGoals: 1,
		AnyOther:             4,
	}
}

// DefaultRiskyRules returns default risky rules
func DefaultRiskyRules() RiskyScoringRules {
	return RiskyScoringRules{
		MaxSelections: 5,
		Events:        DefaultRiskyEvents(),
	}
}

// DefaultRiskyEvents returns default risky events for football
func DefaultRiskyEvents() []RiskyEvent {
	return []RiskyEvent{
		{Slug: "penalty", Name: "Будет пенальти", NameEn: "Penalty awarded", Points: 3},
		{Slug: "red_card", Name: "Будет удаление", NameEn: "Red card shown", Points: 4},
		{Slug: "own_goal", Name: "Будет автогол", NameEn: "Own goal scored", Points: 5},
		{Slug: "hat_trick", Name: "Будет хет-трик", NameEn: "Hat-trick scored", Points: 6},
		{Slug: "clean_sheet_home", Name: "Хозяева на ноль", NameEn: "Home clean sheet", Points: 2},
		{Slug: "clean_sheet_away", Name: "Гости на ноль", NameEn: "Away clean sheet", Points: 3},
		{Slug: "both_teams_score", Name: "Обе забьют", NameEn: "Both teams score", Points: 2},
		{Slug: "over_3_goals", Name: "Больше 3 голов", NameEn: "Over 3.5 goals", Points: 2},
		{Slug: "first_half_draw", Name: "Ничья в 1-м тайме", NameEn: "First half draw", Points: 2},
		{Slug: "comeback", Name: "Камбэк (отыграться)", NameEn: "Comeback from 0:2+", Points: 7},
	}
}

// DefaultTotalizatorRules returns default totalizator rules
func DefaultTotalizatorRules() TotalizatorRules {
	return TotalizatorRules{
		EventCount: 15,
		Scoring:    DefaultStandardRules(),
	}
}

// DefaultRelayRules returns default relay (team) rules
func DefaultRelayRules() RelayRules {
	return RelayRules{
		TeamSize:      5,
		EventCount:    15,
		Scoring:       DefaultStandardRules(),
		AllowReassign: true,
	}
}

// ParseRules parses JSON rules string into ContestRules
func ParseRules(rulesJSON string) (*ContestRules, error) {
	if rulesJSON == "" {
		// Return default standard rules
		defaultRules := DefaultStandardRules()
		return &ContestRules{
			Type:     ContestTypeStandard,
			Standard: &defaultRules,
		}, nil
	}

	var rules ContestRules
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return nil, err
	}

	// Validate and set defaults
	if rules.Type == "" {
		rules.Type = ContestTypeStandard
	}

	if rules.Type == ContestTypeStandard && rules.Standard == nil {
		defaultRules := DefaultStandardRules()
		rules.Standard = &defaultRules
	}

	if rules.Type == ContestTypeRisky && rules.Risky == nil {
		defaultRisky := DefaultRiskyRules()
		rules.Risky = &defaultRisky
	}

	if rules.Type == ContestTypeTotalizator && rules.Totalizator == nil {
		defaultTotalizator := DefaultTotalizatorRules()
		rules.Totalizator = &defaultTotalizator
	}

	if rules.Type == ContestTypeRelay && rules.Relay == nil {
		defaultRelay := DefaultRelayRules()
		rules.Relay = &defaultRelay
	}

	return &rules, nil
}

// ToJSON serializes ContestRules to JSON string
func (r *ContestRules) ToJSON() (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Validate checks if rules are valid
func (r *ContestRules) Validate() error {
	validTypes := map[ContestType]bool{
		ContestTypeStandard:    true,
		ContestTypeRisky:       true,
		ContestTypeTotalizator: true,
		ContestTypeRelay:       true,
	}
	if !validTypes[r.Type] {
		return errors.New("invalid contest type: must be 'standard', 'risky', 'totalizator', or 'relay'")
	}

	if r.Type == ContestTypeStandard {
		if r.Standard == nil {
			return errors.New("standard rules required for standard contest")
		}
		if r.Standard.ExactScore < 0 || r.Standard.GoalDifference < 0 ||
			r.Standard.CorrectOutcome < 0 || r.Standard.AnyOther < 0 {
			return errors.New("scoring points cannot be negative")
		}
	}

	if r.Type == ContestTypeRisky {
		if r.Risky == nil {
			return errors.New("risky rules required for risky contest")
		}
		if r.Risky.MaxSelections < 1 || r.Risky.MaxSelections > 10 {
			return errors.New("max_selections must be between 1 and 10")
		}
		if len(r.Risky.Events) == 0 {
			return errors.New("risky contest must have at least one event")
		}
	}

	if r.Type == ContestTypeTotalizator {
		if r.Totalizator == nil {
			return errors.New("totalizator rules required for totalizator contest")
		}
		if r.Totalizator.EventCount < 5 || r.Totalizator.EventCount > 30 {
			return errors.New("event_count must be between 5 and 30")
		}
		if r.Totalizator.Scoring.ExactScore < 0 || r.Totalizator.Scoring.GoalDifference < 0 ||
			r.Totalizator.Scoring.CorrectOutcome < 0 || r.Totalizator.Scoring.AnyOther < 0 {
			return errors.New("scoring points cannot be negative")
		}
	}

	if r.Type == ContestTypeRelay {
		if r.Relay == nil {
			return errors.New("relay rules required for relay contest")
		}
		if r.Relay.TeamSize < 2 || r.Relay.TeamSize > 10 {
			return errors.New("team_size must be between 2 and 10")
		}
		if r.Relay.EventCount < 5 || r.Relay.EventCount > 50 {
			return errors.New("event_count must be between 5 and 50")
		}
		if r.Relay.Scoring.ExactScore < 0 || r.Relay.Scoring.GoalDifference < 0 ||
			r.Relay.Scoring.CorrectOutcome < 0 || r.Relay.Scoring.AnyOther < 0 {
			return errors.New("scoring points cannot be negative")
		}
	}

	return nil
}

// GetEventBySlug finds a risky event by slug
func (r *RiskyScoringRules) GetEventBySlug(slug string) *RiskyEvent {
	for i := range r.Events {
		if r.Events[i].Slug == slug {
			return &r.Events[i]
		}
	}
	return nil
}

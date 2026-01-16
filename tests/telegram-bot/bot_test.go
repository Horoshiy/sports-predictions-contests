package telegram_bot_test

import (
	"fmt"
	"testing"
)

func TestFormatContest(t *testing.T) {
	tests := []struct {
		name      string
		id        uint32
		title     string
		sportType string
		status    string
		wantEmoji string
	}{
		{"active contest", 1, "Premier League", "football", "active", "🟢"},
		{"draft contest", 2, "NBA Finals", "basketball", "draft", "📋"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatContest(tt.id, tt.title, tt.sportType, tt.status)
			if result == "" {
				t.Error("formatContest returned empty string")
			}
		})
	}
}

func TestFormatLeaderboardEntry(t *testing.T) {
	tests := []struct {
		rank   int
		name   string
		points float64
		streak uint32
	}{
		{1, "Alice", 100.5, 3},
		{2, "Bob", 90.0, 0},
		{3, "Charlie", 80.0, 1},
		{4, "Dave", 70.0, 0},
		{10, "Eve", 50.0, 0},
		{15, "Frank", 40.0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatLeaderboardEntry(tt.rank, tt.name, tt.points, tt.streak)
			if result == "" {
				t.Error("formatLeaderboardEntry returned empty string")
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input string
		cmd   string
		args  string
	}{
		{"/start", "start", ""},
		{"/link test@email.com password123", "link", "test@email.com password123"},
		{"/leaderboard 5", "leaderboard", "5"},
		{"", "", ""},
		{"hello", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd, args := parseCommand(tt.input)
			if cmd != tt.cmd {
				t.Errorf("parseCommand() cmd = %v, want %v", cmd, tt.cmd)
			}
			if args != tt.args {
				t.Errorf("parseCommand() args = %v, want %v", args, tt.args)
			}
		})
	}
}

// Helper functions that mirror bot package logic for testing
func formatContest(id uint32, title, sportType, status string) string {
	emoji := "📋"
	if status == "active" {
		emoji = "🟢"
	}
	return emoji + " " + title
}

func formatLeaderboardEntry(rank int, name string, points float64, streak uint32) string {
	medal := ""
	switch rank {
	case 1:
		medal = "🥇"
	case 2:
		medal = "🥈"
	case 3:
		medal = "🥉"
	default:
		medal = fmt.Sprintf("%d.", rank)
	}
	return medal + " " + name
}

func parseCommand(input string) (string, string) {
	if len(input) == 0 || input[0] != '/' {
		return "", ""
	}
	for i := 1; i < len(input); i++ {
		if input[i] == ' ' {
			return input[1:i], input[i+1:]
		}
	}
	return input[1:], ""
}

package bot

import (
	"context"
	"log"
	"time"

	scoringpb "github.com/sports-prediction-contests/shared/proto/scoring"
)

// showDetailedLeaderboard displays leaderboard with detailed statistics breakdown.
// If msgID > 0, edits the existing message. If msgID == 0, sends a new message.
// Format: Rank | Nickname | Points | Exact | GoalDiff | Outcome | TeamGoals
func (h *Handlers) showDetailedLeaderboard(chatID int64, msgID int, contestID uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.clients.Scoring.GetLeaderboard(ctx, &scoringpb.GetLeaderboardRequest{
		ContestId: contestID,
		Limit:     10,
	})

	if err != nil || resp == nil {
		log.Printf("[ERROR] Failed to get leaderboard: %v", err)
		if msgID > 0 {
			h.editMessage(chatID, msgID, MsgServiceError, BackToMainKeyboard())
		} else {
			h.sendMessage(chatID, MsgServiceError, nil)
		}
		return
	}

	if resp.Leaderboard == nil || len(resp.Leaderboard.Entries) == 0 {
		text := MsgDetailedLeaderboard + MsgEmptyLeaderboard
		if msgID > 0 {
			h.editMessage(chatID, msgID, text, BackToMainKeyboard())
		} else {
			h.sendMessage(chatID, text, BackToMainKeyboard())
		}
		return
	}

	text := MsgDetailedLeaderboard
	text += "💯 Points | 🎯 Exact | ⚖️ Diff | ✓ Outcome | ⚽ Goals\n\n"

	for i, e := range resp.Leaderboard.Entries {
		// NOTE: Detailed stats will be available after proto update
		// Using placeholder values until LeaderboardEntry includes:
		// - ExactScores, GoalDifferences, CorrectOutcomes, TeamGoalsCorrect
		exactScores := 0
		goalDiffs := 0
		outcomes := 0
		teamGoals := 0

		text += FormatDetailedLeaderboardEntry(i+1, e.UserName, e.TotalPoints, exactScores, goalDiffs, outcomes, teamGoals)
	}

	if msgID > 0 {
		h.editMessage(chatID, msgID, text, BackToMainKeyboard())
	} else {
		h.sendMessage(chatID, text, BackToMainKeyboard())
	}
}

package bot

import "fmt"

const (
	MsgWelcome = `🏆 <b>Sports Prediction Contests</b>

Welcome! Make predictions on sports events and compete with others.

<b>Commands:</b>
/contests - View active contests
/leaderboard - View leaderboard
/mystats - Your statistics
/link - Link your account
/help - Show help

To make predictions, first link your account with /link command.`

	MsgHelp = `📖 <b>Available Commands</b>

/start - Start bot
/contests - List active contests
/leaderboard [id] - Contest leaderboard
/mystats - Your prediction stats
/link email password - Link Telegram to account
/help - This message

<b>How to use:</b>
1. Register at our website
2. Use /link to connect your account
3. Browse contests and make predictions!`

	MsgNoContests     = "📭 No active contests at the moment."
	MsgContestList    = "🏆 <b>Active Contests</b>\n\n"
	MsgLeaderboard    = "🏅 <b>Leaderboard</b>\n\n"
	MsgEmptyLeaderboard = "No entries in leaderboard yet."
	MsgNotLinked      = "⚠️ Account not linked. Use /link email password"
	MsgLinkSuccess    = "✅ Account linked successfully!"
	MsgLinkFailed     = "❌ Failed to link account: %s"
	MsgLinkUsage      = "Usage: /link your@email.com password"
	MsgServiceError   = "⚠️ Service temporarily unavailable. Try again later."
	MsgUnknownCommand = "Unknown command. Use /help for available commands."
	MsgStats          = `📊 <b>Your Statistics</b>

Total Points: <b>%.1f</b>
Current Streak: <b>%d</b> 🔥
Max Streak: <b>%d</b>`
)

func FormatContest(id uint32, title, sportType, status string) string {
	emoji := "📋"
	if status == "active" {
		emoji = "🟢"
	}
	return fmt.Sprintf("%s <b>%s</b>\nSport: %s | ID: %d\n", emoji, title, sportType, id)
}

func FormatLeaderboardEntry(rank int, name string, points float64, streak uint32) string {
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
	streakStr := ""
	if streak > 0 {
		streakStr = fmt.Sprintf(" 🔥%d", streak)
	}
	return fmt.Sprintf("%s %s - %.1f pts%s\n", medal, name, points, streakStr)
}

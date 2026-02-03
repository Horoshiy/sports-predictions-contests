package bot

import (
	"fmt"
	"time"

	predictionpb "github.com/sports-prediction-contests/shared/proto/prediction"
)

const (
	MsgWelcome = `🏆 <b>Sports Prediction Contests</b>

Welcome! You're now registered and ready to make predictions on sports events.

<b>Quick Start:</b>
• Use /contests to view active contests
• Select a contest and browse matches
• Make your predictions before matches start
• Check /leaderboard to see rankings
• View /mystats for your performance

<b>Добро пожаловать!</b> Вы зарегистрированы и готовы делать прогнозы на спортивные события.`

	MsgHelp = `📖 <b>Available Commands | Доступные команды</b>

/start - Start bot and register | Начать и зарегистрироваться
/contests - List active contests | Список активных конкурсов
/leaderboard - Contest leaderboard | Таблица лидеров
/mystats - Your prediction stats | Ваша статистика
/link - Link existing web account | Привязать веб-аккаунт
/help - This message | Эта справка

<b>How to use | Как использовать:</b>
1. Browse contests with /contests
2. Select a contest to see matches
3. Make predictions before match starts
4. Compete with others on leaderboard!

<b>Note:</b> Your account is automatically created when you start the bot. If you have an existing web account, use /link to connect it.`

	MsgNoContests         = "📭 No active contests at the moment."
	MsgContestList        = "🏆 <b>Active Contests</b>\n\n"
	MsgLeaderboard        = "🏅 <b>Leaderboard</b>\n\n"
	MsgEmptyLeaderboard   = "No entries in leaderboard yet."
	MsgNotLinked          = "⚠️ Сессия истекла. Нажмите /start для авторизации"
	MsgLinkSuccess        = "✅ Account linked successfully!"
	MsgLinkFailed         = "❌ Failed to link account: %s"
	MsgLinkUsage          = "Usage: /link your@email.com password"
	MsgServiceError       = "⚠️ Service temporarily unavailable. Try again later."
	MsgUnknownCommand     = "Unknown command. Use /help for available commands."
	MsgRegistrationFailed = "❌ Failed to create account. Please try again later."
	MsgStats              = `📊 <b>Your Statistics</b>

Total Points: <b>%.1f</b>
Current Streak: <b>%d</b> 🔥
Max Streak: <b>%d</b>`

	// Match and prediction messages
	MsgMatchList           = "⚽ <b>Matches</b>\n\n"
	MsgNoMatches           = "📭 No matches available."
	MsgMatchDetail         = "⚽ <b>Match Details</b>\n\n"
	MsgMatchNotFound       = "⚠️ Match not found."
	MsgPredictionSuccess   = "✅ Prediction saved!"
	MsgPredictionUpdated   = "✅ Prediction updated!"
	MsgMatchStarted        = "⚠️ Match already started, cannot predict."
	MsgSelectScore         = "Select score prediction:"
	MsgOtherPredictions    = "\n\n👥 <b>Other Predictions:</b>\n"
	MsgDetailedLeaderboard = "🏅 <b>Detailed Leaderboard</b>\n\n"
	MsgSelectContestFirst  = "⚠️ Please select a contest first."
)

// FormatContest formats a contest entry for display in the contest list.
// Returns a formatted string with emoji, title, sport type, and ID.
func FormatContest(id uint32, title, sportType, status string) string {
	emoji := "📋"
	if status == "active" {
		emoji = "🟢"
	}
	// Handle Totalizator (mixed sport type)
	displaySport := sportType
	if sportType == "" || sportType == "mixed" || sportType == "totalizator" {
		displaySport = "🎰 Тотализатор"
	}
	return fmt.Sprintf("%s <b>%s</b>\nТип: %s | ID: %d\n", emoji, title, displaySport, id)
}

// FormatLeaderboardEntry formats a single leaderboard entry with rank, name, points, and streak.
// Ranks 1-3 receive medal emojis (🥇🥈🥉), others show numeric rank.
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

// FormatMatch formats match information with prediction status indicator.
// Shows ✅ if user has made a prediction, ⚪ otherwise.
func FormatMatch(id uint32, homeTeam, awayTeam string, eventDate time.Time, hasPrediction bool) string {
	predIcon := "⚪"
	if hasPrediction {
		predIcon = "✅"
	}
	return fmt.Sprintf("%s <b>%s vs %s</b>\n📅 %s\n\n", predIcon, homeTeam, awayTeam, eventDate.Format("Jan 02, 15:04"))
}

// FormatMatchWithPredictions formats match details including other users' predictions.
// Shows match info, final score if completed, and list of other users' predictions.
func FormatMatchWithPredictions(match *predictionpb.Event, predictions []*predictionpb.Prediction) string {
	text := fmt.Sprintf("⚽ <b>%s vs %s</b>\n\n📅 %s\n", match.HomeTeam, match.AwayTeam, match.EventDate.AsTime().Format("Jan 02, 15:04"))

	if match.Status == "completed" && match.ResultData != "" {
		text += fmt.Sprintf("🏁 Final Score: %s\n", match.ResultData)
	}

	if len(predictions) > 0 {
		text += MsgOtherPredictions
		for _, pred := range predictions {
			text += fmt.Sprintf("• User %d: %s\n", pred.UserId, pred.PredictionData)
		}
	}

	return text
}

// FormatDetailedLeaderboardEntry formats leaderboard entry with detailed statistics breakdown.
// Shows rank, name, total points, and detailed stats (exact scores, goal diffs, outcomes, team goals).
func FormatDetailedLeaderboardEntry(rank int, name string, points float64, exactScores, goalDiffs, outcomes, teamGoals int) string {
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
	return fmt.Sprintf("%s %s\n💯 %.1f pts | 🎯 %d | ⚖️ %d | ✓ %d | ⚽ %d\n\n",
		medal, name, points, exactScores, goalDiffs, outcomes, teamGoals)
}

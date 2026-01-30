# Feature: Complete Telegram Bot Player Experience

## Feature Description

Полная реализация концепции конкурсов прогнозов со стороны игрока через Telegram бота. Бот должен обеспечивать полный цикл взаимодействия: регистрация через Telegram ID (используя firstName и lastName), выбор конкурса, навигация по матчам с пагинацией, подача прогнозов через интуитивные кнопки счета, просмотр таблиц результатов и турнирных таблиц.

Complete implementation of prediction contests from the player's perspective via Telegram bot. The bot should provide a full interaction cycle: registration via Telegram ID (using firstName and lastName), contest selection, match navigation with pagination, prediction submission through intuitive score buttons, viewing results tables and tournament standings.

## User Story

As a **sports prediction contest participant**
I want to **interact with contests entirely through Telegram bot**
So that **I can make predictions, track results, and compete without leaving Telegram**

## Problem Statement

Текущая реализация Telegram бота имеет базовый функционал (просмотр конкурсов, таблицы лидеров, статистика), но не реализует полный игровой цикл:
- Нет регистрации через Telegram (требуется email/password)
- Нет навигации по матчам для прогнозирования
- Нет интерфейса подачи прогнозов
- Нет детальной таблицы результатов с разбивкой по типам угаданных прогнозов
- Нет просмотра прогнозов других участников после начала матча

Current Telegram bot implementation has basic functionality (viewing contests, leaderboards, statistics), but doesn't implement the full player cycle:
- No Telegram-based registration (requires email/password)
- No match navigation for predictions
- No prediction submission interface
- No detailed results table with breakdown by prediction types
- No viewing other participants' predictions after match start

## Solution Statement

Расширить Telegram бота для поддержки полного игрового цикла:
1. **Регистрация через Telegram**: Автоматическое создание пользователя при /start с использованием Telegram ID, firstName, lastName
2. **Навигация по конкурсам**: Выбор конкурса из списка с кнопками
3. **Навигация по матчам**: Пагинированный список матчей конкурса с фильтрацией (непрогнозированные, все)
4. **Интерфейс прогнозирования**: Кнопки счета в 3 колонки (0-0, 1-1, 2-2 / 1-0, 2-0, 2-1 / 3-0, 3-1, 3-2 / 0-1, 0-2, 1-2 / 0-3, 1-3, 2-3 / Любой другой)
5. **Детальная таблица лидеров**: Место, ник, очки, точные счета, разницы мячей, исходы, голы команды
6. **Просмотр прогнозов**: После начала матча показывать прогнозы других участников
7. **Система начисления очков**: 1 за исход, +1 за голы команды, 3 за разницу, 4 за "любой другой", 5 за точный счет

Extend Telegram bot to support full player cycle with all features listed above.

## Feature Metadata

**Feature Type**: Enhancement
**Estimated Complexity**: High
**Primary Systems Affected**: 
- Telegram Bot (bots/telegram/)
- User Service (registration flow)
- Prediction Service (match listing, prediction submission)
- Scoring Service (detailed leaderboard, scoring rules)
- Contest Service (participant management)

**Dependencies**: 
- github.com/go-telegram-bot-api/telegram-bot-api/v5
- Existing gRPC services (user, contest, prediction, scoring)
- PostgreSQL (user storage, predictions)

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Existing Bot Structure:**
- `bots/telegram/main.go` - Entry point, graceful shutdown pattern
- `bots/telegram/bot/bot.go` - Bot lifecycle, update handling
- `bots/telegram/bot/handlers.go` (lines 1-385) - Command and callback handlers, session management
- `bots/telegram/bot/keyboards.go` - Inline keyboard creation patterns
- `bots/telegram/bot/messages.go` - Message constants and formatting functions
- `bots/telegram/clients/clients.go` - gRPC client initialization and management
- `bots/telegram/config/config.go` - Configuration loading

**Proto Definitions:**
- `backend/proto/user.proto` - User service contract (Register, Login)
- `backend/proto/prediction.proto` (lines 1-200) - Prediction and Event messages
- `backend/proto/contest.proto` (lines 1-50) - Contest and Participant messages
- `backend/proto/scoring.proto` - Scoring and Leaderboard messages

**Service Implementations:**
- `backend/user-service/internal/service/user_service.go` - User registration logic
- `backend/prediction-service/internal/service/prediction_service.go` (lines 1-735) - Prediction submission, event listing
- `backend/scoring-service/internal/service/scoring_service.go` (lines 200-400) - Score calculation, leaderboard
- `backend/contest-service/internal/service/contest_service.go` - Contest and participant management

**Models:**
- `backend/user-service/internal/models/user.go` - User model structure
- `backend/prediction-service/internal/models/prediction.go` - Prediction model
- `backend/prediction-service/internal/models/event.go` - Event model
- `backend/scoring-service/internal/models/leaderboard.go` - Leaderboard model

### New Files to Create

- `bots/telegram/bot/registration.go` - Telegram-based user registration handlers
- `bots/telegram/bot/predictions.go` - Match listing and prediction submission handlers
- `bots/telegram/bot/leaderboard_detailed.go` - Enhanced leaderboard with detailed stats
- `bots/telegram/bot/navigation.go` - Navigation state management and pagination
- `bots/telegram/bot/score_buttons.go` - Score prediction button layouts

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Telegram Bot API - Inline Keyboards](https://core.telegram.org/bots/api#inlinekeyboardmarkup)
  - Specific section: InlineKeyboardMarkup structure
  - Why: Required for creating score prediction buttons in 3-column layout
  
- [go-telegram-bot-api Documentation](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5@v5.5.1/tgbotapi)
  - Specific section: NewInlineKeyboardMarkup, NewInlineKeyboardRow
  - Why: Shows proper keyboard construction patterns
  
- [Telegram Bot API - Callback Queries](https://core.telegram.org/bots/api#callbackquery)
  - Specific section: Callback data limitations (64 bytes)
  - Why: Critical for designing callback data format

### Patterns to Follow

**Session Management (from handlers.go:35-50):**
```go
type UserSession struct {
    UserID   uint32
    Email    string
    LinkedAt time.Time
}

// Thread-safe access
func (h *Handlers) getSession(chatID int64) *UserSession {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return h.sessions[chatID]
}

func (h *Handlers) setSession(chatID int64, session *UserSession) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.sessions[chatID] = session
}
```

**Callback Handling (from handlers.go:70-110):**
```go
func (h *Handlers) HandleCallback(cb *tgbotapi.CallbackQuery) {
    data := cb.Data
    chatID := cb.Message.Chat.ID
    msgID := cb.Message.MessageID

    // ALWAYS acknowledge callback first
    h.api.Request(tgbotapi.NewCallback(cb.ID, ""))

    switch {
    case data == "back_main":
        h.editMessage(chatID, msgID, MsgWelcome, MainMenuKeyboard())
    case strings.HasPrefix(data, "contest_"):
        id, _ := strconv.ParseUint(strings.TrimPrefix(data, "contest_"), 10, 32)
        h.handleContestDetail(chatID, msgID, uint32(id))
    }
}
```

**gRPC Client Usage (from handlers.go:120-140):**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := h.clients.Contest.ListContests(ctx, &contestpb.ListContestsRequest{
    Status: "active",
})
if err != nil || resp == nil {
    log.Printf("[ERROR] Failed to list contests: %v", err)
    h.sendMessage(msg.Chat.ID, MsgServiceError, nil)
    return
}
```

**Message Formatting (from messages.go:1-30):**
```go
const (
    MsgWelcome = `🏆 <b>Sports Prediction Contests</b>

Welcome! Make predictions on sports events and compete with others.`
)

func FormatContest(id uint32, title, sportType, status string) string {
    emoji := "📋"
    if status == "active" {
        emoji = "🟢"
    }
    return fmt.Sprintf("%s <b>%s</b>\nSport: %s | ID: %d\n", emoji, title, sportType, id)
}
```

**Keyboard Creation (from keyboards.go:10-35):**
```go
func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
    return tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🏆 Contests", "contests"),
            tgbotapi.NewInlineKeyboardButtonData("🏅 Leaderboard", "leaderboard"),
        ),
    )
}
```

**Error Handling:**
- Always log errors with `[ERROR]` prefix
- Return user-friendly messages (use constants from messages.go)
- Check both `err != nil` and `resp == nil` for gRPC calls

**Naming Conventions:**
- Files: `snake_case.go`
- Functions: `camelCase` (private), `PascalCase` (public)
- Constants: `PascalCase` with `Msg` prefix for messages
- Callback data: `prefix_id` format (e.g., `match_123`, `pred_1_0`)

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - Registration & Session Enhancement

Extend user registration to support Telegram-based signup and enhance session management to track navigation state.

**Tasks:**
- Extend UserSession structure with navigation state
- Implement Telegram-based registration flow
- Add gRPC client for Prediction service
- Create navigation state management utilities

### Phase 2: Core Implementation - Match Navigation & Prediction Interface

Implement match listing with pagination and score prediction button interface.

**Tasks:**
- Create match listing with pagination
- Implement score prediction button layout (3 columns)
- Add prediction submission handlers
- Implement "next unpredicted match" navigation

### Phase 3: Enhanced Leaderboard & Match Details

Implement detailed leaderboard with breakdown and match prediction viewing.

**Tasks:**
- Create detailed leaderboard with stats breakdown
- Implement match details with predictions (after match start)
- Add tournament table view
- Enhance contest detail view with navigation options

### Phase 4: Scoring Rules Implementation

Implement custom scoring rules as specified in requirements.

**Tasks:**
- Update scoring calculation in scoring-service
- Add detailed scoring breakdown
- Implement "any other score" prediction type
- Update leaderboard to show detailed stats

### Phase 5: Testing & Validation

Comprehensive testing of all bot flows.

**Tasks:**
- Test registration flow
- Test prediction submission
- Test navigation and pagination
- Test leaderboard display
- Validate scoring calculations

---

## STEP-BY-STEP TASKS

### Task 1: UPDATE bots/telegram/bot/handlers.go

- **EXTEND**: UserSession structure with navigation state
- **PATTERN**: Existing UserSession structure (handlers.go:25-29)
- **IMPLEMENTATION**:
```go
type UserSession struct {
    UserID      uint32
    Email       string
    LinkedAt    time.Time
    // Navigation state
    CurrentContest uint32
    CurrentPage    int
    ViewMode       string // "matches", "leaderboard", "predictions"
}
```
- **VALIDATE**: `cd bots/telegram && go build`

### Task 2: CREATE bots/telegram/bot/registration.go

- **IMPLEMENT**: Telegram-based user registration
- **PATTERN**: handlers.go handleLink function (lines 320-370)
- **IMPORTS**:
```go
import (
    "context"
    "fmt"
    "log"
    "time"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/sports-prediction-contests/telegram-bot/clients"
    userpb "github.com/sports-prediction-contests/shared/proto/user"
)
```
- **GOTCHA**: Telegram firstName/lastName can be empty, provide defaults
- **VALIDATE**: `cd bots/telegram && go build`


### Task 3: UPDATE bots/telegram/clients/clients.go

- **ADD**: Prediction service client
- **PATTERN**: Existing service client initialization (clients.go:20-70)
- **IMPORTS**: `predictionpb "github.com/sports-prediction-contests/shared/proto/prediction"`
- **IMPLEMENTATION**:
```go
type Clients struct {
    User         userpb.UserServiceClient
    Contest      contestpb.ContestServiceClient
    Scoring      scoringpb.ScoringServiceClient
    Notification notificationpb.NotificationServiceClient
    Prediction   predictionpb.PredictionServiceClient // ADD THIS
    conns        []*grpc.ClientConn
}
```
- **GOTCHA**: Add connection to conns slice for proper cleanup
- **VALIDATE**: `cd bots/telegram && go build`

### Task 4: UPDATE bots/telegram/config/config.go

- **ADD**: Prediction service endpoint configuration
- **PATTERN**: Existing service endpoint fields
- **IMPLEMENTATION**:
```go
type Config struct {
    TelegramBotToken          string
    UserServiceEndpoint       string
    ContestServiceEndpoint    string
    ScoringServiceEndpoint    string
    NotificationServiceEndpoint string
    PredictionServiceEndpoint string // ADD THIS
}
```
- **VALIDATE**: `cd bots/telegram && go build`

### Task 5: CREATE bots/telegram/bot/score_buttons.go

- **IMPLEMENT**: Score prediction button layouts
- **PATTERN**: keyboards.go keyboard creation patterns
- **IMPLEMENTATION**: Create 3-column layout for score buttons
```go
// ScorePredictionKeyboard creates score prediction buttons
// Layout: 3 columns as specified
// Row 1: 0-0, 1-1, 2-2
// Row 2: 1-0, 2-0, 2-1
// Row 3: 3-0, 3-1, 3-2
// Row 4: 0-1, 0-2, 1-2
// Row 5: 0-3, 1-3, 2-3
// Row 6: Any Other (full width)
// Row 7: Back button
func ScorePredictionKeyboard(matchID uint32) tgbotapi.InlineKeyboardMarkup
```
- **GOTCHA**: Callback data limited to 64 bytes, use format `p_matchID_homeScore_awayScore`
- **VALIDATE**: `cd bots/telegram && go build`

### Task 6: CREATE bots/telegram/bot/navigation.go

- **IMPLEMENT**: Navigation state management and pagination utilities
- **PATTERN**: Session management from handlers.go
- **IMPLEMENTATION**:
```go
// NavigationState manages user navigation through contests and matches
type NavigationState struct {
    ContestID   uint32
    Page        int
    ItemsPerPage int
    TotalItems  int
}

// PaginationButtons creates prev/next navigation buttons
func PaginationButtons(state NavigationState, prefix string) []tgbotapi.InlineKeyboardButton

// CalculatePagination calculates page boundaries
func CalculatePagination(page, itemsPerPage, totalItems int) (start, end int)
```
- **VALIDATE**: `cd bots/telegram && go build`

### Task 7: CREATE bots/telegram/bot/predictions.go

- **IMPLEMENT**: Match listing and prediction submission handlers
- **PATTERN**: handlers.go callback handling pattern
- **IMPORTS**:
```go
import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "strconv"
    "strings"
    "time"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/sports-prediction-contests/telegram-bot/clients"
    predictionpb "github.com/sports-prediction-contests/shared/proto/prediction"
    "google.golang.org/protobuf/types/known/timestamppb"
)
```
- **IMPLEMENTATION**:
```go
// handleMatchList shows paginated list of matches for a contest
func (h *Handlers) handleMatchList(chatID int64, msgID int, contestID uint32, page int)

// handleMatchDetail shows match details with prediction buttons
func (h *Handlers) handleMatchDetail(chatID int64, msgID int, matchID uint32)

// handlePredictionSubmit processes score prediction submission
func (h *Handlers) handlePredictionSubmit(chatID int64, msgID int, matchID uint32, homeScore, awayScore int)

// handleAnyOtherScore handles "any other score" prediction
func (h *Handlers) handleAnyOtherScore(chatID int64, msgID int, matchID uint32)

// findNextUnpredictedMatch finds next match without prediction
func (h *Handlers) findNextUnpredictedMatch(ctx context.Context, contestID, userID uint32) (*predictionpb.Event, error)
```
- **GOTCHA**: Check if match already started before allowing prediction
- **VALIDATE**: `cd bots/telegram && go build`

### Task 8: UPDATE bots/telegram/bot/handlers.go - Add Callback Routes

- **ADD**: New callback routes for predictions
- **PATTERN**: Existing HandleCallback switch statement (handlers.go:70-110)
- **IMPLEMENTATION**:
```go
case strings.HasPrefix(data, "matches_"):
    // Format: matches_contestID_page
    parts := strings.Split(strings.TrimPrefix(data, "matches_"), "_")
    contestID, _ := strconv.ParseUint(parts[0], 10, 32)
    page := 1
    if len(parts) > 1 {
        page, _ = strconv.Atoi(parts[1])
    }
    h.handleMatchList(chatID, msgID, uint32(contestID), page)

case strings.HasPrefix(data, "match_"):
    // Format: match_matchID
    id, _ := strconv.ParseUint(strings.TrimPrefix(data, "match_"), 10, 32)
    h.handleMatchDetail(chatID, msgID, uint32(id))

case strings.HasPrefix(data, "p_"):
    // Format: p_matchID_homeScore_awayScore
    parts := strings.Split(strings.TrimPrefix(data, "p_"), "_")
    matchID, _ := strconv.ParseUint(parts[0], 10, 32)
    homeScore, _ := strconv.Atoi(parts[1])
    awayScore, _ := strconv.Atoi(parts[2])
    h.handlePredictionSubmit(chatID, msgID, uint32(matchID), homeScore, awayScore)

case strings.HasPrefix(data, "pany_"):
    // Format: pany_matchID
    id, _ := strconv.ParseUint(strings.TrimPrefix(data, "pany_"), 10, 32)
    h.handleAnyOtherScore(chatID, msgID, uint32(id))
```
- **VALIDATE**: `cd bots/telegram && go build`

### Task 9: UPDATE bots/telegram/bot/keyboards.go - Enhance Contest Detail

- **UPDATE**: ContestDetailKeyboard to include match navigation
- **PATTERN**: Existing ContestDetailKeyboard function
- **IMPLEMENTATION**:
```go
func ContestDetailKeyboard(contestID uint32) tgbotapi.InlineKeyboardMarkup {
    return tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("⚽ Matches", fmt.Sprintf("matches_%d_1", contestID)),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🏅 Leaderboard", fmt.Sprintf("leaderboard_%d", contestID)),
            tgbotapi.NewInlineKeyboardButtonData("🏆 Tournament Table", fmt.Sprintf("tournament_%d", contestID)),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("« Back", "contests"),
        ),
    )
}
```
- **VALIDATE**: `cd bots/telegram && go build`

### Task 10: CREATE bots/telegram/bot/leaderboard_detailed.go

- **IMPLEMENT**: Enhanced leaderboard with detailed statistics
- **PATTERN**: handlers.go showLeaderboard function (lines 220-260)
- **IMPLEMENTATION**:
```go
// showDetailedLeaderboard displays leaderboard with breakdown
// Format: Rank | Nickname | Points | Exact | GoalDiff | Outcome | TeamGoals
func (h *Handlers) showDetailedLeaderboard(chatID int64, msgID int, contestID uint32)

// formatDetailedLeaderboardEntry formats single leaderboard entry
func formatDetailedLeaderboardEntry(rank int, entry *scoringpb.LeaderboardEntry) string
```
- **GOTCHA**: Leaderboard data structure may need extension in scoring service
- **VALIDATE**: `cd bots/telegram && go build`

### Task 11: UPDATE bots/telegram/bot/messages.go - Add New Messages

- **ADD**: Message constants for new features
- **PATTERN**: Existing message constants
- **IMPLEMENTATION**:
```go
const (
    // ... existing messages ...
    
    MsgMatchList = "⚽ <b>Matches</b>\n\n"
    MsgNoMatches = "📭 No matches available."
    MsgMatchDetail = "⚽ <b>Match Details</b>\n\n"
    MsgPredictionSuccess = "✅ Prediction saved!"
    MsgPredictionUpdated = "✅ Prediction updated!"
    MsgMatchStarted = "⚠️ Match already started, cannot predict."
    MsgSelectScore = "Select score prediction:"
    MsgOtherPredictions = "\n\n👥 <b>Other Predictions:</b>\n"
    MsgDetailedLeaderboard = "🏅 <b>Detailed Leaderboard</b>\n\n"
)

// FormatMatch formats match information
func FormatMatch(id uint32, homeTeam, awayTeam string, eventDate time.Time, hasPrediction bool) string

// FormatMatchWithPredictions formats match with other users' predictions
func FormatMatchWithPredictions(match *predictionpb.Event, predictions []*predictionpb.Prediction) string

// FormatDetailedLeaderboardEntry formats leaderboard entry with stats
func FormatDetailedLeaderboardEntry(rank int, name string, points float64, exactScores, goalDiffs, outcomes, teamGoals int) string
```
- **VALIDATE**: `cd bots/telegram && go build`

### Task 12: UPDATE bots/telegram/bot/handlers.go - Implement Registration

- **REFACTOR**: handleStart to include auto-registration
- **PATTERN**: handleLink function (handlers.go:320-370)
- **IMPLEMENTATION**:
```go
func (h *Handlers) handleStart(msg *tgbotapi.Message) {
    session := h.getSession(msg.Chat.ID)
    
    // If already registered, show menu
    if session != nil && session.UserID > 0 {
        h.sendMessage(msg.Chat.ID, MsgWelcome, MainMenuKeyboard())
        return
    }
    
    // Auto-register via Telegram
    h.registerViaTelegram(msg)
}

func (h *Handlers) registerViaTelegram(msg *tgbotapi.Message) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Generate email from Telegram ID
    email := fmt.Sprintf("tg_%d@telegram.bot", msg.From.ID)
    
    // Use firstName + lastName or username as name
    name := strings.TrimSpace(msg.From.FirstName + " " + msg.From.LastName)
    if name == "" {
        name = msg.From.UserName
    }
    if name == "" {
        name = fmt.Sprintf("User%d", msg.From.ID)
    }
    
    // Register user
    resp, err := h.clients.User.Register(ctx, &userpb.RegisterRequest{
        Email:    email,
        Password: fmt.Sprintf("tg_%d_%d", msg.From.ID, time.Now().Unix()),
        Name:     name,
    })
    
    // Handle response and create session
    // ...
}
```
- **GOTCHA**: Check if user already exists, handle duplicate registration gracefully
- **VALIDATE**: `cd bots/telegram && go build`

### Task 13: IMPLEMENT Match Prediction Viewing After Start

- **ADD**: Function to show match predictions after match starts
- **PATTERN**: handlers.go gRPC call patterns
- **IMPLEMENTATION**:
```go
// handleMatchPredictions shows all predictions for a match (after start)
func (h *Handlers) handleMatchPredictions(chatID int64, msgID int, matchID uint32) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Get match details
    matchResp, err := h.clients.Prediction.GetEvent(ctx, &predictionpb.GetEventRequest{
        Id: matchID,
    })
    
    // Check if match started
    if !matchStarted(matchResp.Event) {
        h.editMessage(chatID, msgID, "Match hasn't started yet.", BackToMainKeyboard())
        return
    }
    
    // Get all predictions for this match
    // Format and display with scores if match completed
}
```
- **VALIDATE**: `cd bots/telegram && go build`

### Task 14: UPDATE backend/proto/scoring.proto - Extend Leaderboard Entry

- **ADD**: Detailed statistics fields to LeaderboardEntry message
- **PATTERN**: Existing LeaderboardEntry message structure
- **IMPLEMENTATION**:
```protobuf
message LeaderboardEntry {
    uint32 user_id = 1;
    string user_name = 2;
    double total_points = 3;
    uint32 current_streak = 4;
    uint32 rank = 5;
    // ADD THESE:
    uint32 exact_scores = 6;
    uint32 goal_differences = 7;
    uint32 correct_outcomes = 8;
    uint32 team_goals_correct = 9;
}
```
- **GOTCHA**: After modifying proto, must regenerate Go code
- **VALIDATE**: `cd backend && make proto`

### Task 15: UPDATE backend/scoring-service - Implement Custom Scoring Rules

- **UPDATE**: calculateExactScorePoints function in scoring_service.go
- **PATTERN**: Existing scoring calculation (scoring_service.go:280-340)
- **IMPLEMENTATION**:
```go
// New scoring rules:
// - 5 points: Exact score
// - 4 points: Any other score (not in common predictions)
// - 3 points: Goal difference
// - 1 point: Correct outcome
// - +1 point: Correct goals for one team (if outcome correct)

func (s *ScoringService) calculateExactScorePoints(prediction PredictionData, result ResultData, details map[string]interface{}) (float64, map[string]interface{}) {
    if prediction.HomeScore == nil || prediction.AwayScore == nil {
        return 0, details
    }
    
    predictedHome := *prediction.HomeScore
    predictedAway := *prediction.AwayScore
    
    // Check if "any other" prediction
    isAnyOther := prediction.Type == "any_other"
    
    // Exact match: 5 points
    if predictedHome == result.HomeScore && predictedAway == result.AwayScore {
        details["match_type"] = "exact"
        return 5, details
    }
    
    // Any other score: 4 points (if predicted as "any other" and not common score)
    if isAnyOther && !isCommonScore(result.HomeScore, result.AwayScore) {
        details["match_type"] = "any_other"
        return 4, details
    }
    
    // Goal difference: 3 points
    predictedDiff := predictedHome - predictedAway
    actualDiff := result.HomeScore - result.AwayScore
    if predictedDiff == actualDiff {
        details["match_type"] = "goal_difference"
        return 3, details
    }
    
    // Correct outcome: 1 point + bonus for team goals
    predictedWinner := s.determineWinner(predictedHome, predictedAway)
    if predictedWinner == result.Winner {
        points := 1.0
        
        // +1 for correct home team goals
        if predictedHome == result.HomeScore {
            points += 1.0
            details["home_goals_correct"] = true
        }
        
        // +1 for correct away team goals
        if predictedAway == result.AwayScore {
            points += 1.0
            details["away_goals_correct"] = true
        }
        
        details["match_type"] = "outcome"
        return points, details
    }
    
    return 0, details
}

// isCommonScore checks if score is in common predictions (0-0 to 3-3 range)
func isCommonScore(home, away int) bool {
    return home >= 0 && home <= 3 && away >= 0 && away <= 3
}
```
- **VALIDATE**: `cd backend/scoring-service && go test ./...`

### Task 16: UPDATE backend/scoring-service/internal/repository - Track Detailed Stats

- **UPDATE**: Leaderboard repository to track detailed statistics
- **PATTERN**: Existing leaderboard_repository.go
- **IMPLEMENTATION**: Add fields to track exact_scores, goal_differences, correct_outcomes, team_goals_correct
- **VALIDATE**: `cd backend/scoring-service && go test ./internal/repository/...`

### Task 17: TEST Registration Flow

- **MANUAL**: Test Telegram registration
- **STEPS**:
  1. Start bot with `/start` command
  2. Verify auto-registration creates user
  3. Check session is created
  4. Verify main menu appears
- **VALIDATE**: Check PostgreSQL users table for new Telegram user

### Task 18: TEST Match Navigation and Prediction

- **MANUAL**: Test complete prediction flow
- **STEPS**:
  1. Select contest from list
  2. Navigate to matches
  3. Test pagination (if >5 matches)
  4. Select match
  5. Submit score prediction
  6. Verify prediction saved
  7. Check "next unpredicted match" navigation
- **VALIDATE**: Check PostgreSQL predictions table

### Task 19: TEST Leaderboard Display

- **MANUAL**: Test detailed leaderboard
- **STEPS**:
  1. View contest leaderboard
  2. Verify all columns displayed correctly
  3. Check ranking order
  4. Verify detailed stats (exact scores, goal diffs, etc.)
- **VALIDATE**: Compare with database leaderboard data

### Task 20: TEST Scoring Calculations

- **UNIT**: Test new scoring rules
- **IMPLEMENTATION**: Create test cases in scoring_service_test.go
```go
func TestCustomScoringRules(t *testing.T) {
    // Test exact score: 5 points
    // Test any other: 4 points
    // Test goal difference: 3 points
    // Test outcome: 1 point
    // Test outcome + team goals: 1+1 or 1+2 points
}
```
- **VALIDATE**: `cd backend/scoring-service && go test -v ./internal/service/...`

---

## TESTING STRATEGY

### Unit Tests

**Scoring Service Tests** (`backend/scoring-service/internal/service/scoring_service_test.go`):
- Test exact score prediction (5 points)
- Test "any other" score prediction (4 points)
- Test goal difference prediction (3 points)
- Test outcome prediction (1 point)
- Test outcome + team goals (1+1, 1+2 points)
- Test edge cases (nil scores, invalid data)

**Bot Handler Tests** (`bots/telegram/bot/handlers_test.go`):
- Test session management (get/set)
- Test callback data parsing
- Test pagination calculations
- Test navigation state transitions

### Integration Tests

**End-to-End Bot Flow**:
1. Registration via /start
2. Contest selection
3. Match listing with pagination
4. Prediction submission
5. Leaderboard viewing
6. Match predictions viewing (after start)

### Manual Testing Checklist

- [ ] Registration creates user with Telegram ID
- [ ] Session persists across interactions
- [ ] Contest list displays correctly
- [ ] Match pagination works (prev/next)
- [ ] Score buttons layout correct (3 columns)
- [ ] Prediction submission successful
- [ ] Cannot predict after match start
- [ ] Can view others' predictions after match start
- [ ] Leaderboard shows detailed stats
- [ ] Scoring calculations correct
- [ ] Navigation back buttons work
- [ ] Error messages display properly

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
# Build Telegram bot
cd bots/telegram && go build

# Build scoring service
cd backend/scoring-service && go build

# Regenerate proto files (if modified)
cd backend && make proto
```

### Level 2: Unit Tests

```bash
# Test scoring service
cd backend/scoring-service && go test -v ./...

# Test bot handlers
cd bots/telegram && go test -v ./bot/...
```

### Level 3: Integration Tests

```bash
# Start all services
make docker-services

# Run E2E tests
make e2e-test
```

### Level 4: Manual Validation

```bash
# Start Telegram bot
cd bots/telegram
export TELEGRAM_BOT_TOKEN="your_token"
export USER_SERVICE_ENDPOINT="localhost:8084"
export CONTEST_SERVICE_ENDPOINT="localhost:8085"
export PREDICTION_SERVICE_ENDPOINT="localhost:8086"
export SCORING_SERVICE_ENDPOINT="localhost:8087"
go run main.go

# Test in Telegram:
# 1. Send /start to bot
# 2. Navigate through menus
# 3. Submit predictions
# 4. View leaderboards
```

### Level 5: Database Validation

```bash
# Check registered users
psql -h localhost -U sports_user -d sports_prediction -c "SELECT id, email, name FROM users WHERE email LIKE 'tg_%@telegram.bot';"

# Check predictions
psql -h localhost -U sports_user -d sports_prediction -c "SELECT * FROM predictions WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'tg_%@telegram.bot');"

# Check leaderboard
psql -h localhost -U sports_user -d sports_prediction -c "SELECT * FROM leaderboards ORDER BY total_points DESC LIMIT 10;"
```

---

## ACCEPTANCE CRITERIA

- [ ] User can register via Telegram /start without email/password
- [ ] User can select contest from button list
- [ ] User can navigate matches with pagination (5 per page)
- [ ] User can submit score prediction via 3-column button layout
- [ ] Score buttons include: 0-0 to 3-3 range + "Any Other"
- [ ] After prediction, bot shows next unpredicted match
- [ ] User cannot predict after match starts
- [ ] User can view other participants' predictions after match start
- [ ] Leaderboard shows: Rank, Name, Points, Exact, GoalDiff, Outcome, TeamGoals
- [ ] Scoring rules implemented: 5 (exact), 4 (any other), 3 (diff), 1 (outcome), +1 (team goals)
- [ ] All validation commands pass with zero errors
- [ ] Navigation back buttons work correctly
- [ ] Session persists across interactions
- [ ] Error handling displays user-friendly messages

---

## COMPLETION CHECKLIST

- [ ] All 20 tasks completed in order
- [ ] Each task validation passed
- [ ] Proto files regenerated (if modified)
- [ ] All unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Database validation successful
- [ ] No build errors or warnings
- [ ] Code follows existing patterns
- [ ] All acceptance criteria met

---

## NOTES

### Design Decisions

**Registration Strategy**: Auto-register on /start using Telegram ID as unique identifier. Generate email as `tg_{telegram_id}@telegram.bot` to maintain compatibility with existing user service.

**Callback Data Format**: Use short prefixes to stay within 64-byte limit:
- `p_matchID_home_away` for predictions (e.g., `p_123_2_1`)
- `matches_contestID_page` for match lists
- `match_matchID` for match details

**Pagination**: 5 matches per page to keep messages readable on mobile devices.

**Score Button Layout**: 3-column layout as specified, with "Any Other" as full-width button for emphasis.

**Scoring Rules**: Custom implementation differs from existing scoring service. New rules:
- 5 points: Exact score (was 10)
- 4 points: Any other score (new)
- 3 points: Goal difference (was 5)
- 1 point: Outcome (was 3)
- +1 point: Team goals bonus (new)

### Trade-offs

**Session Storage**: Using in-memory sessions for simplicity. For production, consider Redis for persistence across bot restarts.

**Prediction Data Format**: Using JSON string for flexibility. "Any other" predictions stored as `{"type": "any_other", "home_score": null, "away_score": null}`.

**Leaderboard Stats**: Requires extending proto definition and database schema. Alternative: calculate stats on-the-fly (slower but no schema changes).

### Future Enhancements

- Add filters for match list (unpredicted only, by date)
- Implement contest search
- Add prediction editing (before match start)
- Show prediction history
- Add notifications for match start/results
- Implement tournament bracket view
- Add user statistics dashboard

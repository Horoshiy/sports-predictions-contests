package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	predictionpb "github.com/sports-prediction-contests/shared/proto/prediction"
	"google.golang.org/grpc/metadata"
)

const (
	matchesPerPage = 5
	minScore       = 0
	maxScore       = 20
)

// handleMatchList shows paginated list of matches for a contest
func (h *Handlers) handleMatchList(chatID int64, msgID int, contestID uint32, page int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := h.getSession(chatID)
	if session == nil {
		h.editMessage(chatID, msgID, MsgNotLinked, BackToMainKeyboard())
		return
	}

	// Get events for contest
	resp, err := h.clients.Prediction.ListEvents(ctx, &predictionpb.ListEventsRequest{
		SportType: "",
		Status:    "scheduled",
		ContestId: contestID,
	})

	if err != nil || resp == nil {
		log.Printf("[ERROR] Failed to list events (status=scheduled): %v", err)
		h.editMessage(chatID, msgID, MsgServiceError, BackToMainKeyboard())
		return
	}

	if len(resp.Events) == 0 {
		h.editMessage(chatID, msgID, MsgNoMatches, BackToMainKeyboard())
		return
	}

	// Get user's predictions for this contest
	userPredictions := make(map[uint32]string) // eventId -> "X:Y" or "other"
	predCtx := metadata.AppendToOutgoingContext(ctx, "x-user-id", strconv.FormatUint(uint64(session.UserID), 10))
	predResp, err := h.clients.Prediction.GetUserPredictions(predCtx, &predictionpb.GetUserPredictionsRequest{
		ContestId: contestID,
	})
	if err == nil && predResp != nil && predResp.Predictions != nil {
		for _, pred := range predResp.Predictions {
			var predData struct {
				Type      string `json:"type"`
				HomeScore *int   `json:"home_score"`
				AwayScore *int   `json:"away_score"`
			}
			if json.Unmarshal([]byte(pred.PredictionData), &predData) == nil {
				if predData.Type == "any_other" {
					userPredictions[pred.EventId] = "other"
				} else if predData.HomeScore != nil && predData.AwayScore != nil {
					userPredictions[pred.EventId] = fmt.Sprintf("%d:%d", *predData.HomeScore, *predData.AwayScore)
				}
			}
		}
	}

	// Calculate pagination
	totalMatches := len(resp.Events)
	start, end := CalculatePagination(page, matchesPerPage, totalMatches)

	// Build message
	text := MsgMatchList
	for i := start; i < end && i < len(resp.Events); i++ {
		event := resp.Events[i]
		text += FormatMatch(event.Id, event.HomeTeam, event.AwayTeam, event.EventDate.AsTime(), false)
	}

	// Build keyboard
	var rows [][]tgbotapi.InlineKeyboardButton

	// Match buttons with predictions
	for i := start; i < end && i < len(resp.Events); i++ {
		event := resp.Events[i]
		buttonText := fmt.Sprintf("⚽ %s vs %s", event.HomeTeam, event.AwayTeam)
		if score, ok := userPredictions[event.Id]; ok {
			buttonText = fmt.Sprintf("✅ %s vs %s (%s)", event.HomeTeam, event.AwayTeam, score)
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, fmt.Sprintf("match_%d", event.Id)),
		))
	}

	// Pagination
	navState := NavigationState{
		ContestID:    contestID,
		Page:         page,
		ItemsPerPage: matchesPerPage,
		TotalItems:   totalMatches,
	}
	rows = append(rows, PaginationButtons(navState, "matches"))

	// Back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Back to Contest", fmt.Sprintf("contest_%d", contestID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	h.editMessage(chatID, msgID, text, keyboard)
}

// handleMatchDetail shows match details with prediction buttons
func (h *Handlers) handleMatchDetail(chatID int64, msgID int, matchID uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := h.getSession(chatID)
	if session == nil {
		h.editMessage(chatID, msgID, MsgNotLinked, BackToMainKeyboard())
		return
	}

	// Get match details
	resp, err := h.clients.Prediction.GetEvent(ctx, &predictionpb.GetEventRequest{
		Id: matchID,
	})

	if err != nil || resp == nil || resp.Event == nil {
		log.Printf("[ERROR] Failed to get event %d: %v", matchID, err)
		h.editMessage(chatID, msgID, MsgMatchNotFound, BackToMainKeyboard())
		return
	}

	event := resp.Event
	eventTime := event.EventDate.AsTime()

	// Check if match already started
	if time.Now().After(eventTime) {
		h.editMessage(chatID, msgID, MsgMatchStarted, BackToMainKeyboard())
		return
	}

	// Check for existing prediction
	existingPrediction := ""
	contestID := session.CurrentContest
	log.Printf("[DEBUG] handleMatchDetail: matchID=%d, contestID=%d, userID=%d", matchID, contestID, session.UserID)
	
	if contestID > 0 {
		// Add user_id to gRPC metadata for bot authentication
		predCtx := metadata.AppendToOutgoingContext(ctx, "x-user-id", strconv.FormatUint(uint64(session.UserID), 10))
		
		predResp, err := h.clients.Prediction.GetUserPredictions(predCtx, &predictionpb.GetUserPredictionsRequest{
			ContestId: contestID,
		})
		
		if err != nil {
			log.Printf("[DEBUG] GetUserPredictions error: %v", err)
		} else if predResp != nil && predResp.Predictions != nil {
			log.Printf("[DEBUG] Found %d predictions", len(predResp.Predictions))
			for _, pred := range predResp.Predictions {
				log.Printf("[DEBUG] Prediction: eventId=%d, data=%s", pred.EventId, pred.PredictionData)
				if pred.EventId == matchID {
					// Parse prediction data to extract score
					var predData struct {
						Type      string `json:"type"`
						HomeScore *int   `json:"home_score"`
						AwayScore *int   `json:"away_score"`
					}
					if json.Unmarshal([]byte(pred.PredictionData), &predData) == nil {
						if predData.Type == "any_other" {
							existingPrediction = "\n\n✅ <b>Твой прогноз:</b> Any other"
						} else if predData.HomeScore != nil && predData.AwayScore != nil {
							existingPrediction = fmt.Sprintf("\n\n✅ <b>Твой прогноз:</b> %d : %d", *predData.HomeScore, *predData.AwayScore)
						}
						log.Printf("[DEBUG] Matched! existingPrediction=%s", existingPrediction)
					}
					break
				}
			}
		}
	} else {
		log.Printf("[DEBUG] contestID is 0, skipping prediction lookup")
	}

	// Build message
	selectText := MsgSelectScore
	if existingPrediction != "" {
		selectText = "Выбери новый счёт для изменения прогноза:"
	}
	
	text := fmt.Sprintf("%s<b>%s vs %s</b>\n\n📅 %s%s\n\n%s",
		MsgMatchDetail,
		event.HomeTeam,
		event.AwayTeam,
		eventTime.Format("Jan 02, 15:04"),
		existingPrediction,
		selectText,
	)

	h.editMessage(chatID, msgID, text, ScorePredictionKeyboard(matchID))
}

// handlePredictionSubmit processes score prediction submission
func (h *Handlers) handlePredictionSubmit(chatID int64, msgID int, matchID uint32, homeScore, awayScore int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := h.getSession(chatID)
	if session == nil {
		h.editMessage(chatID, msgID, MsgNotLinked, BackToMainKeyboard())
		return
	}

	// Validate score values
	if homeScore < minScore || awayScore < minScore || homeScore > maxScore || awayScore > maxScore {
		h.editMessage(chatID, msgID,
			fmt.Sprintf("⚠️ Invalid score. Please use values between %d-%d.", minScore, maxScore),
			BackToMainKeyboard())
		return
	}

	// Get match details to verify it hasn't started
	eventResp, err := h.clients.Prediction.GetEvent(ctx, &predictionpb.GetEventRequest{
		Id: matchID,
	})

	if err != nil || eventResp == nil || eventResp.Event == nil {
		log.Printf("[ERROR] Failed to get event %d for validation: %v", matchID, err)
		h.editMessage(chatID, msgID, MsgMatchNotFound, BackToMainKeyboard())
		return
	}

	if time.Now().After(eventResp.Event.EventDate.AsTime()) {
		h.editMessage(chatID, msgID, MsgMatchStarted, BackToMainKeyboard())
		return
	}

	// Create prediction data
	predictionData := map[string]interface{}{
		"type":       "exact_score",
		"home_score": homeScore,
		"away_score": awayScore,
	}

	predictionJSON, err := json.Marshal(predictionData)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal prediction data: %v", err)
		h.editMessage(chatID, msgID, "Failed to save prediction.", BackToMainKeyboard())
		return
	}

	// Submit prediction - require contest to be selected
	contestID := session.CurrentContest
	if contestID == 0 {
		h.editMessage(chatID, msgID, MsgSelectContestFirst, BackToMainKeyboard())
		return
	}

	// Add user_id to gRPC metadata for bot authentication
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", strconv.FormatUint(uint64(session.UserID), 10))

	resp, err := h.clients.Prediction.SubmitPrediction(ctx, &predictionpb.SubmitPredictionRequest{
		ContestId:      contestID,
		EventId:        matchID,
		PredictionData: string(predictionJSON),
	})

	if err != nil || resp == nil || resp.Response == nil || !resp.Response.Success {
		errMsg := "Failed to save prediction"
		if resp != nil && resp.Response != nil {
			errMsg = resp.Response.Message
		}
		log.Printf("[ERROR] Failed to submit prediction (contest=%d, event=%d, user=%d): %v", contestID, matchID, session.UserID, err)
		h.editMessage(chatID, msgID, fmt.Sprintf("❌ %s", errMsg), BackToMainKeyboard())
		return
	}

	// Success message - score prediction
	log.Printf("[INFO] Prediction submitted (user=%d, contest=%d, match=%d, score=%d-%d)", session.UserID, contestID, matchID, homeScore, awayScore)
	successMsg := fmt.Sprintf("%s\n\nPrediction: %d-%d", MsgPredictionSuccess, homeScore, awayScore)

	// Try to find next unpredicted match with fresh context
	nextCtx, nextCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer nextCancel()
	nextEvent, err := h.findNextUnpredictedMatch(nextCtx, contestID, session.UserID)
	if err == nil && nextEvent != nil {
		successMsg += fmt.Sprintf("\n\n⏭️ Next match: %s vs %s", nextEvent.HomeTeam, nextEvent.AwayTeam)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⚽ Predict Next Match", fmt.Sprintf("match_%d", nextEvent.Id)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back to Matches", fmt.Sprintf("matches_%d_1", contestID)),
			),
		)
		h.editMessage(chatID, msgID, successMsg, keyboard)
	} else {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("« Back to Matches", fmt.Sprintf("matches_%d_1", contestID)),
			),
		)
		h.editMessage(chatID, msgID, successMsg, keyboard)
	}
}

// handleAnyOtherScore handles "any other score" prediction
func (h *Handlers) handleAnyOtherScore(chatID int64, msgID int, matchID uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := h.getSession(chatID)
	if session == nil {
		h.editMessage(chatID, msgID, MsgNotLinked, BackToMainKeyboard())
		return
	}

	// Get match details to verify it hasn't started
	eventResp, err := h.clients.Prediction.GetEvent(ctx, &predictionpb.GetEventRequest{
		Id: matchID,
	})

	if err != nil || eventResp == nil || eventResp.Event == nil {
		log.Printf("[ERROR] Failed to get event %d for validation: %v", matchID, err)
		h.editMessage(chatID, msgID, MsgMatchNotFound, BackToMainKeyboard())
		return
	}

	if time.Now().After(eventResp.Event.EventDate.AsTime()) {
		h.editMessage(chatID, msgID, MsgMatchStarted, BackToMainKeyboard())
		return
	}

	// Create "any other" prediction data
	predictionData := map[string]interface{}{
		"type":       "any_other",
		"home_score": nil,
		"away_score": nil,
	}

	predictionJSON, err := json.Marshal(predictionData)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal prediction data: %v", err)
		h.editMessage(chatID, msgID, "Failed to save prediction.", BackToMainKeyboard())
		return
	}

	// Submit prediction - require contest to be selected
	contestID := session.CurrentContest
	if contestID == 0 {
		h.editMessage(chatID, msgID, MsgSelectContestFirst, BackToMainKeyboard())
		return
	}

	// Add user_id to gRPC metadata for bot authentication
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", strconv.FormatUint(uint64(session.UserID), 10))

	resp, err := h.clients.Prediction.SubmitPrediction(ctx, &predictionpb.SubmitPredictionRequest{
		ContestId:      contestID,
		EventId:        matchID,
		PredictionData: string(predictionJSON),
	})

	if err != nil || resp == nil || resp.Response == nil || !resp.Response.Success {
		errMsg := "Failed to save prediction"
		if resp != nil && resp.Response != nil {
			errMsg = resp.Response.Message
		}
		log.Printf("[ERROR] Failed to submit prediction (contest=%d, event=%d, user=%d): %v", contestID, matchID, session.UserID, err)
		h.editMessage(chatID, msgID, fmt.Sprintf("❌ %s", errMsg), BackToMainKeyboard())
		return
	}

	// Success message - any other score
	successMsg := fmt.Sprintf("%s\n\nPrediction: Any Other Score 🎲", MsgPredictionSuccess)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Back to Matches", fmt.Sprintf("matches_%d_1", contestID)),
		),
	)
	h.editMessage(chatID, msgID, successMsg, keyboard)
}

// findNextUnpredictedMatch finds next match without prediction
func (h *Handlers) findNextUnpredictedMatch(ctx context.Context, contestID, userID uint32) (*predictionpb.Event, error) {
	// Get events for this contest
	resp, err := h.clients.Prediction.ListEvents(ctx, &predictionpb.ListEventsRequest{
		SportType: "",
		Status:    "scheduled",
		ContestId: contestID,
	})

	if err != nil || resp == nil || len(resp.Events) == 0 {
		return nil, fmt.Errorf("no events found")
	}

	// Get user's predictions
	predResp, err := h.clients.Prediction.GetUserPredictions(ctx, &predictionpb.GetUserPredictionsRequest{
		ContestId: contestID,
	})

	if err != nil {
		// If error getting predictions, return first event
		return resp.Events[0], nil
	}

	// Build map of predicted event IDs
	predictedEvents := make(map[uint32]bool)
	if predResp != nil && predResp.Predictions != nil {
		for _, pred := range predResp.Predictions {
			predictedEvents[pred.EventId] = true
		}
	}

	// Find first unpredicted event
	for _, event := range resp.Events {
		if !predictedEvents[event.Id] && time.Now().Before(event.EventDate.AsTime()) {
			return event, nil
		}
	}

	return nil, fmt.Errorf("no unpredicted matches found")
}

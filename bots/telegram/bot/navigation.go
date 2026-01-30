package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NavigationState manages user navigation through contests and matches
type NavigationState struct {
	ContestID    uint32
	Page         int
	ItemsPerPage int
	TotalItems   int
}

// PaginationButtons creates prev/next navigation buttons
func PaginationButtons(state NavigationState, prefix string) []tgbotapi.InlineKeyboardButton {
	var buttons []tgbotapi.InlineKeyboardButton

	// Validate itemsPerPage to prevent division by zero
	if state.ItemsPerPage <= 0 {
		state.ItemsPerPage = 1
	}

	totalPages := (state.TotalItems + state.ItemsPerPage - 1) / state.ItemsPerPage
	if totalPages < 1 {
		totalPages = 1
	}

	// Previous button
	if state.Page > 1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			"◀️ Prev",
			fmt.Sprintf("%s_%d_%d", prefix, state.ContestID, state.Page-1),
		))
	}

	// Page indicator
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
		fmt.Sprintf("%d/%d", state.Page, totalPages),
		"page_info",
	))

	// Next button
	if state.Page < totalPages {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			"Next ▶️",
			fmt.Sprintf("%s_%d_%d", prefix, state.ContestID, state.Page+1),
		))
	}

	return buttons
}

// CalculatePagination calculates the start and end indices for a page of items.
// It handles edge cases like invalid page numbers and ensures indices don't exceed totalItems.
// Parameters:
//   - page: 1-based page number (will be clamped to 1 if less)
//   - itemsPerPage: number of items per page (will be clamped to 1 if less or zero)
//   - totalItems: total number of items available (will be clamped to 0 if negative)
// Returns:
//   - start: 0-based start index (inclusive)
//   - end: 0-based end index (exclusive)
func CalculatePagination(page, itemsPerPage, totalItems int) (start, end int) {
	if itemsPerPage <= 0 {
		itemsPerPage = 1
	}
	if totalItems < 0 {
		totalItems = 0
	}
	if page < 1 {
		page = 1
	}

	// Prevent overflow with extreme values
	if page > 1000000 || itemsPerPage > 1000 {
		return 0, 0
	}

	start = (page - 1) * itemsPerPage
	end = start + itemsPerPage

	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	return start, end
}

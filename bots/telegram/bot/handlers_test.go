package bot

import (
	"testing"
	"time"
)

// TestHandlersShutdown tests that Shutdown properly stops the cleanup goroutine
func TestHandlersShutdown(t *testing.T) {
	// Create a minimal handlers instance for testing
	h := &Handlers{
		sessions:   make(map[int64]*UserSession),
		sessionTTL: 1 * time.Second,
		shutdownCh: make(chan struct{}),
	}

	// Start cleanup goroutine
	go h.cleanupSessions()

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	// Shutdown should complete quickly
	done := make(chan struct{})
	go func() {
		h.Shutdown()
		close(done)
	}()

	select {
	case <-done:
		// Success - shutdown completed
	case <-time.After(1 * time.Second):
		t.Error("Shutdown did not complete within 1 second")
	}
}

// TestLastActivityTracking tests that LastActivity is updated on getSession
func TestLastActivityTracking(t *testing.T) {
	h := &Handlers{
		sessions:   make(map[int64]*UserSession),
		sessionTTL: 24 * time.Hour,
		shutdownCh: make(chan struct{}),
	}
	defer h.Shutdown()

	chatID := int64(12345)
	now := time.Now()

	// Create initial session
	h.setSession(chatID, &UserSession{
		UserID:       1,
		Email:        "test@example.com",
		LinkedAt:     now,
		LastActivity: now,
	})

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Get session should update LastActivity
	session := h.getSession(chatID)
	if session == nil {
		t.Fatal("Session not found")
	}

	if !session.LastActivity.After(now) {
		t.Errorf("LastActivity was not updated: got %v, want after %v", session.LastActivity, now)
	}
}

// TestSessionCleanupByLastActivity tests that sessions are cleaned up based on LastActivity
func TestSessionCleanupByLastActivity(t *testing.T) {
	h := &Handlers{
		sessions:   make(map[int64]*UserSession),
		sessionTTL: 100 * time.Millisecond,
		shutdownCh: make(chan struct{}),
	}
	defer h.Shutdown()

	chatID := int64(12345)
	oldTime := time.Now().Add(-200 * time.Millisecond)

	// Create session with old LastActivity
	h.setSession(chatID, &UserSession{
		UserID:       1,
		Email:        "test@example.com",
		LinkedAt:     oldTime,
		LastActivity: oldTime,
	})

	// Manually trigger cleanup
	h.mu.Lock()
	now := time.Now()
	cleaned := 0
	for cid, session := range h.sessions {
		if now.Sub(session.LastActivity) > h.sessionTTL {
			delete(h.sessions, cid)
			cleaned++
		}
	}
	h.mu.Unlock()

	if cleaned != 1 {
		t.Errorf("Expected 1 session to be cleaned, got %d", cleaned)
	}

	// Verify session was removed
	session := h.getSession(chatID)
	if session != nil {
		t.Error("Session should have been cleaned up but still exists")
	}
}

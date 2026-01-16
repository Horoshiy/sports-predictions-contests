package sync

import (
	"context"
	"log"
	"sync"
	"time"
)

// SyncWorker runs periodic sync operations
type SyncWorker struct {
	syncService *SyncService
	interval    time.Duration
	quit        chan bool
	wg          sync.WaitGroup
	running     bool
	mu          sync.Mutex
}

// NewSyncWorker creates a new sync worker
func NewSyncWorker(syncService *SyncService, intervalMins int) *SyncWorker {
	return &SyncWorker{
		syncService: syncService,
		interval:    time.Duration(intervalMins) * time.Minute,
		quit:        make(chan bool),
	}
}

// Start begins the periodic sync worker
func (w *SyncWorker) Start() {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.mu.Unlock()

	w.wg.Add(1)
	go w.run()
	log.Printf("[INFO] Sync worker started with interval %v", w.interval)
}

// Stop gracefully stops the sync worker
func (w *SyncWorker) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	w.mu.Unlock()

	close(w.quit)
	w.wg.Wait()
	log.Println("[INFO] Sync worker stopped")
}

// IsRunning returns whether the worker is running
func (w *SyncWorker) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

func (w *SyncWorker) run() {
	defer w.wg.Done()

	// Run initial sync
	w.runSync()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.runSync()
		case <-w.quit:
			return
		}
	}
}

func (w *SyncWorker) runSync() {
	ctx := context.Background()
	log.Println("[INFO] Starting scheduled sync...")

	// Sync sports first (logging handled by SyncService)
	if _, err := w.syncService.SyncSports(ctx); err != nil {
		log.Printf("[ERROR] Failed to sync sports: %v", err)
	}

	// Sync leagues (logging handled by SyncService)
	if _, err := w.syncService.SyncLeagues(ctx); err != nil {
		log.Printf("[ERROR] Failed to sync leagues: %v", err)
	}

	// Sync match results (logging handled by SyncService)
	if _, err := w.syncService.SyncMatchResults(ctx); err != nil {
		log.Printf("[ERROR] Failed to sync match results: %v", err)
	}

	log.Println("[INFO] Scheduled sync completed")
}

// GetLastSyncAt returns the last sync timestamp from the sync service
func (w *SyncWorker) GetLastSyncAt() *time.Time {
	return w.syncService.GetLastSyncAt()
}

// TriggerSync manually triggers a sync operation
func (w *SyncWorker) TriggerSync(ctx context.Context, entityType string, parentID uint) (int, error) {
	switch entityType {
	case "sports":
		return w.syncService.SyncSports(ctx)
	case "leagues":
		return w.syncService.SyncLeagues(ctx)
	case "teams":
		if parentID == 0 {
			return 0, nil
		}
		return w.syncService.SyncTeamsByLeague(ctx, parentID)
	case "matches":
		if parentID == 0 {
			return 0, nil
		}
		return w.syncService.SyncUpcomingMatches(ctx, parentID)
	case "results":
		return w.syncService.SyncMatchResults(ctx)
	default:
		return 0, nil
	}
}

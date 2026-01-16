package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/sports-prediction-contests/notification-service/internal/channels"
	"github.com/sports-prediction-contests/notification-service/internal/models"
	"github.com/sports-prediction-contests/notification-service/internal/repository"
)

type NotificationJob struct {
	Notification *models.Notification
	Preference   *models.NotificationPreference
}

type WorkerPool struct {
	jobs     chan NotificationJob
	quit     chan bool
	wg       sync.WaitGroup
	telegram *channels.TelegramChannel
	email    *channels.EmailChannel
	repo     repository.NotificationRepositoryInterface
	size     int
}

func NewWorkerPool(size int, telegram *channels.TelegramChannel, email *channels.EmailChannel, repo repository.NotificationRepositoryInterface) *WorkerPool {
	bufferSize := size * 20 // Buffer 20 jobs per worker
	if bufferSize < 100 {
		bufferSize = 100
	}
	return &WorkerPool{
		jobs:     make(chan NotificationJob, bufferSize),
		quit:     make(chan bool),
		telegram: telegram,
		email:    email,
		repo:     repo,
		size:     size,
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.size; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}
	log.Printf("Started %d notification workers", w.size)
}

func (w *WorkerPool) Stop() {
	close(w.quit)
	// Drain remaining jobs
	for {
		select {
		case job := <-w.jobs:
			w.processJob(-1, job) // Process remaining jobs
		default:
			w.wg.Wait()
			log.Println("All notification workers stopped")
			return
		}
	}
}

func (w *WorkerPool) Submit(job NotificationJob) {
	select {
	case w.jobs <- job:
	default:
		log.Println("Warning: notification job queue full, dropping job")
	}
}

func (w *WorkerPool) worker(id int) {
	defer w.wg.Done()
	for {
		select {
		case job := <-w.jobs:
			w.processJob(id, job)
		case <-w.quit:
			return
		}
	}
}

func (w *WorkerPool) processJob(workerID int, job NotificationJob) {
	n := job.Notification
	p := job.Preference
	var sent bool

	switch n.Channel {
	case "telegram":
		if w.telegram.IsEnabled() && p != nil && p.TelegramChatID != 0 {
			if err := w.telegram.Send(p.TelegramChatID, n.Title, n.Message); err != nil {
				log.Printf("Worker %d: telegram send error: %v", workerID, err)
				return
			}
			sent = true
		}
	case "email":
		if w.email.IsEnabled() && p != nil && p.Email != "" {
			if err := w.email.Send(p.Email, n.Title, n.Message); err != nil {
				log.Printf("Worker %d: email send error: %v", workerID, err)
				return
			}
			sent = true
		}
	}

	if sent {
		now := time.Now()
		n.SentAt = &now
		ctx := context.Background()
		if err := w.repo.Update(ctx, n); err != nil {
			log.Printf("Worker %d: failed to update notification sent_at: %v", workerID, err)
		}
	}
}

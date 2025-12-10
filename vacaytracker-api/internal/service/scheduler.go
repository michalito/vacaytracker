package service

import (
	"context"
	"log"
	"sync"
	"time"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/repository/sqlite"
)

// Scheduler handles background scheduled tasks
type Scheduler struct {
	newsletterService *NewsletterService
	settingsRepo      *sqlite.SettingsRepository
	ticker            *time.Ticker
	done              chan bool
	mu                sync.Mutex
	running           bool
}

// NewScheduler creates a new background scheduler
func NewScheduler(
	newsletterService *NewsletterService,
	settingsRepo *sqlite.SettingsRepository,
) *Scheduler {
	return &Scheduler{
		newsletterService: newsletterService,
		settingsRepo:      settingsRepo,
		done:              make(chan bool),
	}
}

// Start begins the scheduler loop
// Checks every hour if newsletter should be sent
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	// Check every hour
	s.ticker = time.NewTicker(1 * time.Hour)

	go func() {
		// Check immediately on startup
		s.checkAndSendNewsletter()

		for {
			select {
			case <-s.ticker.C:
				s.checkAndSendNewsletter()
			case <-s.done:
				s.ticker.Stop()
				return
			}
		}
	}()

	log.Println("[SCHEDULER] Newsletter scheduler started (checking every hour)")
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		s.done <- true
		s.running = false
		log.Println("[SCHEDULER] Newsletter scheduler stopped")
	}
}

// checkAndSendNewsletter determines if newsletter should be sent
func (s *Scheduler) checkAndSendNewsletter() {
	ctx := context.Background()

	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		log.Printf("[SCHEDULER] Failed to get settings: %v", err)
		return
	}

	if !settings.Newsletter.Enabled {
		return
	}

	if !s.shouldSendNewsletter(settings) {
		return
	}

	log.Println("[SCHEDULER] Triggering scheduled newsletter send")

	count, err := s.newsletterService.Send(ctx)
	if err != nil {
		log.Printf("[SCHEDULER] Failed to send newsletter: %v", err)
		return
	}

	log.Printf("[SCHEDULER] Newsletter sent to %d recipients", count)
}

// shouldSendNewsletter checks if it's time to send based on config
func (s *Scheduler) shouldSendNewsletter(settings *domain.Settings) bool {
	return s.shouldSendNewsletterAt(settings, time.Now())
}

// shouldSendNewsletterAt checks if it's time to send based on config at a specific time
// This is separated for testability
func (s *Scheduler) shouldSendNewsletterAt(settings *domain.Settings, now time.Time) bool {
	config := settings.Newsletter

	// Check if newsletter is enabled
	if !config.Enabled {
		return false
	}

	// Check if already sent today
	if config.LastSentAt != nil {
		lastSent := *config.LastSentAt
		if isSameDay(lastSent, now) {
			return false
		}
	}

	switch config.Frequency {
	case "monthly":
		// Send on configured day of month
		return now.Day() == config.DayOfMonth
	case "weekly":
		// Send on Monday (weekday 1)
		return now.Weekday() == time.Monday
	default:
		return false
	}
}

// isSameDay checks if two times are on the same calendar day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

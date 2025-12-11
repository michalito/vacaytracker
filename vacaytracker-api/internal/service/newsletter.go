package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"vacaytracker-api/internal/config"
	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/repository/sqlite"
)

const (
	// LowBalanceThreshold is the number of vacation days below which users get a reminder
	LowBalanceThreshold = 5
)

// NewsletterData holds all content for the newsletter
type NewsletterData struct {
	AppURL            string
	RecipientName     string
	Period            string
	Stats             *sqlite.MonthlyStats
	UpcomingVacations []*domain.TeamVacation
	LowBalanceUsers   []LowBalanceUser
	HasUpcoming       bool
	HasLowBalance     bool
}

// LowBalanceUser represents a user with low vacation balance
type LowBalanceUser struct {
	UserName      string
	RemainingDays int
}

// NewsletterService handles newsletter generation and sending
type NewsletterService struct {
	cfg          *config.Config
	userRepo     *sqlite.UserRepository
	vacationRepo *sqlite.VacationRepository
	settingsRepo *sqlite.SettingsRepository
	emailService *EmailService
}

// NewNewsletterService creates a new NewsletterService
func NewNewsletterService(
	cfg *config.Config,
	userRepo *sqlite.UserRepository,
	vacationRepo *sqlite.VacationRepository,
	settingsRepo *sqlite.SettingsRepository,
	emailService *EmailService,
) *NewsletterService {
	return &NewsletterService{
		cfg:          cfg,
		userRepo:     userRepo,
		vacationRepo: vacationRepo,
		settingsRepo: settingsRepo,
		emailService: emailService,
	}
}

// GetRecipients returns users who have weeklyDigest email preference enabled
func (s *NewsletterService) GetRecipients(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetNewsletterRecipients(ctx)
}

// GetStats returns aggregated statistics for the previous month
func (s *NewsletterService) GetStats(ctx context.Context) (*sqlite.MonthlyStats, string, error) {
	// Get previous month
	now := time.Now()
	prevMonth := now.AddDate(0, -1, 0)
	year := prevMonth.Year()
	month := int(prevMonth.Month())

	stats, err := s.vacationRepo.GetMonthlyStats(ctx, year, month)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get monthly stats: %w", err)
	}

	period := prevMonth.Format("January 2006")
	return stats, period, nil
}

// GetUpcomingVacations returns approved vacations for the next month
func (s *NewsletterService) GetUpcomingVacations(ctx context.Context) ([]*domain.TeamVacation, error) {
	// Get next month
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	year := nextMonth.Year()
	month := int(nextMonth.Month())

	return s.vacationRepo.ListTeam(ctx, month, year)
}

// GetLowBalanceUsers returns users with vacation balance at or below the threshold
func (s *NewsletterService) GetLowBalanceUsers(ctx context.Context) ([]LowBalanceUser, error) {
	users, err := s.userRepo.GetLowBalanceUsers(ctx, LowBalanceThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get low balance users: %w", err)
	}

	result := make([]LowBalanceUser, len(users))
	for i, user := range users {
		result[i] = LowBalanceUser{
			UserName:      user.Name,
			RemainingDays: user.VacationBalance,
		}
	}

	return result, nil
}

// BuildNewsletterData assembles all newsletter content for a specific recipient
func (s *NewsletterService) BuildNewsletterData(ctx context.Context, recipientName string) (*NewsletterData, error) {
	stats, period, err := s.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	upcoming, err := s.GetUpcomingVacations(ctx)
	if err != nil {
		return nil, err
	}

	lowBalance, err := s.GetLowBalanceUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &NewsletterData{
		AppURL:            s.cfg.AppURL,
		RecipientName:     recipientName,
		Period:            period,
		Stats:             stats,
		UpcomingVacations: upcoming,
		LowBalanceUsers:   lowBalance,
		HasUpcoming:       len(upcoming) > 0,
		HasLowBalance:     len(lowBalance) > 0,
	}, nil
}

// GeneratePreview creates a preview of the newsletter without sending
func (s *NewsletterService) GeneratePreview(ctx context.Context) (*dto.NewsletterPreviewResponse, error) {
	recipients, err := s.GetRecipients(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipients: %w", err)
	}

	// Build newsletter data for preview (using "Preview User" as recipient)
	data, err := s.BuildNewsletterData(ctx, "Preview User")
	if err != nil {
		return nil, fmt.Errorf("failed to build newsletter data: %w", err)
	}

	// Render templates using pre-compiled newsletter templates
	htmlBody, err := s.emailService.RenderNewsletterHTML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to render HTML template: %w", err)
	}

	textBody, err := s.emailService.RenderNewsletterText(data)
	if err != nil {
		return nil, fmt.Errorf("failed to render text template: %w", err)
	}

	// Extract recipient emails
	recipientEmails := make([]string, len(recipients))
	for i, r := range recipients {
		recipientEmails[i] = r.Email
	}

	return &dto.NewsletterPreviewResponse{
		Subject:        newsletterSubject,
		HTMLBody:       htmlBody,
		TextBody:       textBody,
		Recipients:     recipientEmails,
		RecipientCount: len(recipients),
	}, nil
}

// Send sends the newsletter to all eligible recipients
func (s *NewsletterService) Send(ctx context.Context) (int, error) {
	recipients, err := s.GetRecipients(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get recipients: %w", err)
	}

	if len(recipients) == 0 {
		log.Println("[NEWSLETTER] No recipients found")
		return 0, nil
	}

	sentCount := 0
	for _, recipient := range recipients {
		// Build personalized newsletter data
		data, err := s.BuildNewsletterData(ctx, recipient.Name)
		if err != nil {
			log.Printf("[NEWSLETTER ERROR] Failed to build data for %s: %v", recipient.Email, err)
			continue
		}

		// Render templates using pre-compiled newsletter templates
		htmlBody, err := s.emailService.RenderNewsletterHTML(data)
		if err != nil {
			log.Printf("[NEWSLETTER ERROR] Failed to render HTML for %s: %v", recipient.Email, err)
			continue
		}

		textBody, err := s.emailService.RenderNewsletterText(data)
		if err != nil {
			log.Printf("[NEWSLETTER ERROR] Failed to render text for %s: %v", recipient.Email, err)
			continue
		}

		// Send email (non-blocking) with newsletter-specific options
		opts := &SendOptions{
			IdempotencyKey: generateIdempotencyKey(recipient.Email, newsletterSubject, data.Period),
			Tags:           []string{"newsletter", "monthly-summary"},
		}
		s.emailService.SendAsync(recipient.Email, newsletterSubject, htmlBody, textBody, opts)
		sentCount++
	}

	// Update last sent timestamp
	if sentCount > 0 {
		if err := s.UpdateLastSent(ctx); err != nil {
			log.Printf("[NEWSLETTER ERROR] Failed to update last sent timestamp: %v", err)
		}
	}

	log.Printf("[NEWSLETTER] Sent to %d recipients", sentCount)
	return sentCount, nil
}

// UpdateLastSent updates the lastSentAt timestamp in settings
func (s *NewsletterService) UpdateLastSent(ctx context.Context) error {
	return s.settingsRepo.UpdateLastNewsletterSent(ctx, time.Now())
}

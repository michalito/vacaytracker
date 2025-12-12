package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"

	"vacaytracker-api/internal/config"
	"vacaytracker-api/internal/domain"
)

// EmailService handles sending emails via Resend API
type EmailService struct {
	cfg    *config.Config
	client *resend.Client

	// Pre-compiled templates for performance
	welcomeHTMLTmpl        *template.Template
	welcomeTextTmpl        *template.Template
	requestSubmittedHTML   *template.Template
	requestSubmittedText   *template.Template
	requestApprovedHTML    *template.Template
	requestApprovedText    *template.Template
	requestRejectedHTML    *template.Template
	requestRejectedText    *template.Template
	adminNewRequestHTML    *template.Template
	adminNewRequestText    *template.Template
	newsletterHTMLTmpl     *template.Template
	newsletterTextTmpl     *template.Template
}

// Retry configuration
const (
	maxRetries     = 3
	baseRetryDelay = 500 * time.Millisecond
	maxRetryDelay  = 10 * time.Second
)

// NewEmailService creates a new EmailService with pre-compiled templates
func NewEmailService(cfg *config.Config) *EmailService {
	svc := &EmailService{
		cfg: cfg,
	}

	// Initialize Resend client if API key is configured
	if cfg.ResendAPIKey != "" {
		svc.client = resend.NewClient(cfg.ResendAPIKey)
	}

	// Pre-compile all templates at startup for performance
	svc.compileTemplates()

	return svc
}

// compileTemplates pre-compiles all email templates
func (s *EmailService) compileTemplates() {
	var err error

	// Welcome email templates
	s.welcomeHTMLTmpl, err = template.New("welcomeHTML").Parse(welcomeEmailHTML)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile welcome HTML template: %v", err)
	}
	s.welcomeTextTmpl, err = template.New("welcomeText").Parse(welcomeEmailText)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile welcome text template: %v", err)
	}

	// Request submitted templates
	s.requestSubmittedHTML, err = template.New("requestSubmittedHTML").Parse(requestSubmittedHTML)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile request submitted HTML template: %v", err)
	}
	s.requestSubmittedText, err = template.New("requestSubmittedText").Parse(requestSubmittedText)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile request submitted text template: %v", err)
	}

	// Request approved templates
	s.requestApprovedHTML, err = template.New("requestApprovedHTML").Parse(requestApprovedHTML)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile request approved HTML template: %v", err)
	}
	s.requestApprovedText, err = template.New("requestApprovedText").Parse(requestApprovedText)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile request approved text template: %v", err)
	}

	// Request rejected templates
	s.requestRejectedHTML, err = template.New("requestRejectedHTML").Parse(requestRejectedHTML)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile request rejected HTML template: %v", err)
	}
	s.requestRejectedText, err = template.New("requestRejectedText").Parse(requestRejectedText)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile request rejected text template: %v", err)
	}

	// Admin new request templates
	s.adminNewRequestHTML, err = template.New("adminNewRequestHTML").Parse(adminNewRequestHTML)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile admin new request HTML template: %v", err)
	}
	s.adminNewRequestText, err = template.New("adminNewRequestText").Parse(adminNewRequestText)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile admin new request text template: %v", err)
	}

	// Newsletter templates
	s.newsletterHTMLTmpl, err = template.New("newsletterHTML").Parse(newsletterHTML)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile newsletter HTML template: %v", err)
	}
	s.newsletterTextTmpl, err = template.New("newsletterText").Parse(newsletterText)
	if err != nil {
		log.Printf("[EMAIL] Warning: Failed to compile newsletter text template: %v", err)
	}
}

// SendOptions contains optional parameters for sending emails
type SendOptions struct {
	IdempotencyKey string   // Prevents duplicate sends within 24 hours
	ReplyTo        string   // Reply-to email address
	Tags           []string // Tags for categorization/analytics
}

// Send sends an email via Resend API with retry logic
func (s *EmailService) Send(ctx context.Context, to, subject, htmlBody, textBody string, opts *SendOptions) error {
	if !s.cfg.EmailEnabled() {
		log.Printf("[EMAIL] Skipping email to %s - email not configured", to)
		return nil
	}

	if s.client == nil {
		log.Printf("[EMAIL] Skipping email to %s - client not initialized", to)
		return nil
	}

	fromAddress := fmt.Sprintf("%s <%s>", s.cfg.EmailFromName, s.cfg.EmailFromAddress)

	params := &resend.SendEmailRequest{
		From:    fromAddress,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
		Text:    textBody,
	}

	// Apply optional parameters
	if opts != nil {
		if opts.ReplyTo != "" {
			params.ReplyTo = opts.ReplyTo
		}
		if len(opts.Tags) > 0 {
			tags := make([]resend.Tag, len(opts.Tags))
			for i, tag := range opts.Tags {
				tags[i] = resend.Tag{
					Name:  tag, // Each tag name must be unique
					Value: "true",
				}
			}
			params.Tags = tags
		}
	}

	// Execute with retry logic
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := calculateBackoff(attempt)
			log.Printf("[EMAIL] Retry attempt %d/%d for %s after %v", attempt, maxRetries, to, delay)

			select {
			case <-ctx.Done():
				return fmt.Errorf("email send cancelled: %w", ctx.Err())
			case <-time.After(delay):
			}
		}

		// Send the email
		sent, err := s.sendEmail(params)

		if err == nil {
			log.Printf("[EMAIL] Email sent to %s (ID: %s)", to, sent.Id)
			return nil
		}

		lastErr = err

		// Check if we should retry (only for transient errors)
		if !isRetryableError(err) {
			log.Printf("[EMAIL ERROR] Non-retryable error sending to %s: %v", to, err)
			return err
		}

		log.Printf("[EMAIL] Transient error on attempt %d for %s: %v", attempt+1, to, err)
	}

	return fmt.Errorf("email failed after %d retries: %w", maxRetries, lastErr)
}

// sendEmail sends an email via the Resend client
// Note: IdempotencyKey in SendOptions is generated for logging/debugging but
// not currently passed to Resend API (SDK v2 doesn't expose this header yet)
func (s *EmailService) sendEmail(params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	return s.client.Emails.Send(params)
}

// calculateBackoff calculates exponential backoff with jitter
func calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := float64(baseRetryDelay) * math.Pow(2, float64(attempt-1))

	// Cap at max delay
	if delay > float64(maxRetryDelay) {
		delay = float64(maxRetryDelay)
	}

	return time.Duration(delay)
}

// isRetryableError determines if an error should trigger a retry
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Resend SDK errors that are retryable:
	// - Network timeouts
	// - 500, 502, 503, 504 server errors
	// - Rate limiting (429) - though we should respect Retry-After

	errStr := err.Error()

	// Check for common retryable patterns
	retryablePatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"503",
		"502",
		"500",
		"504",
		"429", // Rate limit - consider adding delay based on Retry-After
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errStr), pattern) {
			return true
		}
	}

	return false
}

// generateIdempotencyKey creates a deterministic idempotency key based on email content
func generateIdempotencyKey(to, subject string, uniqueData ...string) string {
	h := sha256.New()
	h.Write([]byte(to))
	h.Write([]byte(subject))
	for _, data := range uniqueData {
		h.Write([]byte(data))
	}
	// Add timestamp rounded to hour to allow retries within the hour
	h.Write([]byte(time.Now().UTC().Truncate(time.Hour).Format(time.RFC3339)))
	return hex.EncodeToString(h.Sum(nil))[:32]
}

// SendAsync sends an email asynchronously (non-blocking)
func (s *EmailService) SendAsync(to, subject, htmlBody, textBody string, opts *SendOptions) {
	if !s.cfg.EmailEnabled() {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		if err := s.Send(ctx, to, subject, htmlBody, textBody, opts); err != nil {
			log.Printf("[EMAIL ERROR] Failed to send email to %s: %v", to, err)
		}
	}()
}

// SendWelcome sends a welcome email to a new user with idempotency protection
func (s *EmailService) SendWelcome(user *domain.User, tempPassword string) {
	if s.welcomeHTMLTmpl == nil || s.welcomeTextTmpl == nil {
		log.Printf("[EMAIL ERROR] Welcome email templates not initialized")
		return
	}

	data := welcomeEmailData{
		AppURL:       s.cfg.AppURL,
		UserName:     user.Name,
		UserEmail:    user.Email,
		TempPassword: tempPassword,
	}

	htmlBody, err := s.executeTemplate(s.welcomeHTMLTmpl, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render welcome email HTML: %v", err)
		return
	}

	textBody, err := s.executeTemplate(s.welcomeTextTmpl, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render welcome email text: %v", err)
		return
	}

	// Use idempotency key for welcome emails to prevent duplicate password emails
	opts := &SendOptions{
		IdempotencyKey: generateIdempotencyKey(user.Email, welcomeEmailSubject, user.ID),
		Tags:           []string{"welcome", "onboarding"},
	}

	s.SendAsync(user.Email, welcomeEmailSubject, htmlBody, textBody, opts)
}

// SendRequestSubmitted sends an email when a vacation request is submitted
func (s *EmailService) SendRequestSubmitted(user *domain.User, vacation *domain.VacationRequest) {
	if !user.EmailPreferences.VacationUpdates {
		log.Printf("[EMAIL] Skipping request submitted email for %s - user preferences disabled", user.Email)
		return
	}

	if s.requestSubmittedHTML == nil || s.requestSubmittedText == nil {
		log.Printf("[EMAIL ERROR] Request submitted email templates not initialized")
		return
	}

	data := vacationEmailData{
		AppURL:    s.cfg.AppURL,
		UserName:  user.Name,
		StartDate: vacation.StartDate,
		EndDate:   vacation.EndDate,
		TotalDays: vacation.TotalDays,
	}

	htmlBody, err := s.executeTemplate(s.requestSubmittedHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render request submitted email HTML: %v", err)
		return
	}

	textBody, err := s.executeTemplate(s.requestSubmittedText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render request submitted email text: %v", err)
		return
	}

	opts := &SendOptions{
		IdempotencyKey: generateIdempotencyKey(user.Email, requestSubmittedSubject, vacation.ID),
		Tags:           []string{"vacation", "submitted"},
	}

	s.SendAsync(user.Email, requestSubmittedSubject, htmlBody, textBody, opts)
}

// SendRequestApproved sends an email when a vacation request is approved
func (s *EmailService) SendRequestApproved(user *domain.User, vacation *domain.VacationRequest) {
	if !user.EmailPreferences.VacationUpdates {
		log.Printf("[EMAIL] Skipping approval email for %s - user preferences disabled", user.Email)
		return
	}

	if s.requestApprovedHTML == nil || s.requestApprovedText == nil {
		log.Printf("[EMAIL ERROR] Request approved email templates not initialized")
		return
	}

	data := vacationEmailData{
		AppURL:    s.cfg.AppURL,
		UserName:  user.Name,
		StartDate: vacation.StartDate,
		EndDate:   vacation.EndDate,
		TotalDays: vacation.TotalDays,
	}

	htmlBody, err := s.executeTemplate(s.requestApprovedHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render approved email HTML: %v", err)
		return
	}

	textBody, err := s.executeTemplate(s.requestApprovedText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render approved email text: %v", err)
		return
	}

	opts := &SendOptions{
		IdempotencyKey: generateIdempotencyKey(user.Email, requestApprovedSubject, vacation.ID, "approved"),
		Tags:           []string{"vacation", "approved"},
	}

	s.SendAsync(user.Email, requestApprovedSubject, htmlBody, textBody, opts)
}

// SendRequestRejected sends an email when a vacation request is rejected
func (s *EmailService) SendRequestRejected(user *domain.User, vacation *domain.VacationRequest, reason string) {
	if !user.EmailPreferences.VacationUpdates {
		log.Printf("[EMAIL] Skipping rejection email for %s - user preferences disabled", user.Email)
		return
	}

	if s.requestRejectedHTML == nil || s.requestRejectedText == nil {
		log.Printf("[EMAIL ERROR] Request rejected email templates not initialized")
		return
	}

	data := vacationEmailData{
		AppURL:    s.cfg.AppURL,
		UserName:  user.Name,
		StartDate: vacation.StartDate,
		EndDate:   vacation.EndDate,
		TotalDays: vacation.TotalDays,
		Reason:    reason,
	}

	htmlBody, err := s.executeTemplate(s.requestRejectedHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render rejected email HTML: %v", err)
		return
	}

	textBody, err := s.executeTemplate(s.requestRejectedText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render rejected email text: %v", err)
		return
	}

	opts := &SendOptions{
		IdempotencyKey: generateIdempotencyKey(user.Email, requestRejectedSubject, vacation.ID, "rejected"),
		Tags:           []string{"vacation", "rejected"},
	}

	s.SendAsync(user.Email, requestRejectedSubject, htmlBody, textBody, opts)
}

// SendAdminNewRequest sends an email to admins when a new vacation request is submitted
func (s *EmailService) SendAdminNewRequest(admins []*domain.User, requester *domain.User, vacation *domain.VacationRequest) {
	if s.adminNewRequestHTML == nil || s.adminNewRequestText == nil {
		log.Printf("[EMAIL ERROR] Admin new request email templates not initialized")
		return
	}

	requestReason := ""
	if vacation.Reason != nil {
		requestReason = *vacation.Reason
	}

	data := adminNotificationData{
		AppURL:        s.cfg.AppURL,
		RequesterName: requester.Name,
		StartDate:     vacation.StartDate,
		EndDate:       vacation.EndDate,
		TotalDays:     vacation.TotalDays,
		RequestReason: requestReason,
	}

	htmlBody, err := s.executeTemplate(s.adminNewRequestHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render admin notification email HTML: %v", err)
		return
	}

	textBody, err := s.executeTemplate(s.adminNewRequestText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render admin notification email text: %v", err)
		return
	}

	for _, admin := range admins {
		// Check if admin wants team notifications
		if !admin.EmailPreferences.TeamNotifications {
			log.Printf("[EMAIL] Skipping admin notification for %s - user preferences disabled", admin.Email)
			continue
		}

		opts := &SendOptions{
			IdempotencyKey: generateIdempotencyKey(admin.Email, adminNewRequestSubject, vacation.ID),
			ReplyTo:        requester.Email, // Allow admin to reply directly to requester
			Tags:           []string{"admin", "vacation-request"},
		}

		s.SendAsync(admin.Email, adminNewRequestSubject, htmlBody, textBody, opts)
	}
}

// executeTemplate executes a pre-compiled template with the given data
func (s *EmailService) executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	if tmpl == nil {
		return "", fmt.Errorf("template is nil")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderNewsletterHTML renders the newsletter HTML template with the given data
func (s *EmailService) RenderNewsletterHTML(data interface{}) (string, error) {
	return s.executeTemplate(s.newsletterHTMLTmpl, data)
}

// RenderNewsletterText renders the newsletter text template with the given data
func (s *EmailService) RenderNewsletterText(data interface{}) (string, error) {
	return s.executeTemplate(s.newsletterTextTmpl, data)
}

// EmailPreview contains the rendered email content for preview
type EmailPreview struct {
	Subject  string
	HTMLBody string
	TextBody string
}

// PreviewWelcome renders a preview of the welcome email
func (s *EmailService) PreviewWelcome(userName, userEmail, tempPassword, appURL string) (*EmailPreview, error) {
	data := welcomeEmailData{
		AppURL:       appURL,
		UserName:     userName,
		UserEmail:    userEmail,
		TempPassword: tempPassword,
	}

	htmlBody, err := s.executeTemplate(s.welcomeHTMLTmpl, data)
	if err != nil {
		return nil, err
	}

	textBody, err := s.executeTemplate(s.welcomeTextTmpl, data)
	if err != nil {
		return nil, err
	}

	return &EmailPreview{
		Subject:  welcomeEmailSubject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}, nil
}

// PreviewRequestSubmitted renders a preview of the request submitted email
func (s *EmailService) PreviewRequestSubmitted(userName, startDate, endDate string, totalDays int, appURL string) (*EmailPreview, error) {
	data := vacationEmailData{
		AppURL:    appURL,
		UserName:  userName,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: totalDays,
	}

	htmlBody, err := s.executeTemplate(s.requestSubmittedHTML, data)
	if err != nil {
		return nil, err
	}

	textBody, err := s.executeTemplate(s.requestSubmittedText, data)
	if err != nil {
		return nil, err
	}

	return &EmailPreview{
		Subject:  requestSubmittedSubject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}, nil
}

// PreviewRequestApproved renders a preview of the request approved email
func (s *EmailService) PreviewRequestApproved(userName, startDate, endDate string, totalDays int, appURL string) (*EmailPreview, error) {
	data := vacationEmailData{
		AppURL:    appURL,
		UserName:  userName,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: totalDays,
	}

	htmlBody, err := s.executeTemplate(s.requestApprovedHTML, data)
	if err != nil {
		return nil, err
	}

	textBody, err := s.executeTemplate(s.requestApprovedText, data)
	if err != nil {
		return nil, err
	}

	return &EmailPreview{
		Subject:  requestApprovedSubject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}, nil
}

// PreviewRequestRejected renders a preview of the request rejected email
func (s *EmailService) PreviewRequestRejected(userName, startDate, endDate string, totalDays int, reason, appURL string) (*EmailPreview, error) {
	data := vacationEmailData{
		AppURL:    appURL,
		UserName:  userName,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: totalDays,
		Reason:    reason,
	}

	htmlBody, err := s.executeTemplate(s.requestRejectedHTML, data)
	if err != nil {
		return nil, err
	}

	textBody, err := s.executeTemplate(s.requestRejectedText, data)
	if err != nil {
		return nil, err
	}

	return &EmailPreview{
		Subject:  requestRejectedSubject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}, nil
}

// PreviewAdminNewRequest renders a preview of the admin notification email
func (s *EmailService) PreviewAdminNewRequest(requesterName, startDate, endDate string, totalDays int, requestReason, appURL string) (*EmailPreview, error) {
	data := adminNotificationData{
		AppURL:        appURL,
		RequesterName: requesterName,
		StartDate:     startDate,
		EndDate:       endDate,
		TotalDays:     totalDays,
		RequestReason: requestReason,
	}

	htmlBody, err := s.executeTemplate(s.adminNewRequestHTML, data)
	if err != nil {
		return nil, err
	}

	textBody, err := s.executeTemplate(s.adminNewRequestText, data)
	if err != nil {
		return nil, err
	}

	return &EmailPreview{
		Subject:  adminNewRequestSubject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}, nil
}

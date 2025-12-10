package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"vacaytracker-api/internal/config"
	"vacaytracker-api/internal/domain"
)

const resendAPIURL = "https://api.resend.com/emails"

// EmailService handles sending emails via Resend API
type EmailService struct {
	cfg    *config.Config
	client *http.Client
}

// NewEmailService creates a new EmailService
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// resendRequest represents the Resend API request body
type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

// resendResponse represents the Resend API response
type resendResponse struct {
	ID string `json:"id"`
}

// resendError represents the Resend API error response
type resendError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Name       string `json:"name"`
}

// Send sends an email via Resend API
func (s *EmailService) Send(to, subject, htmlBody, textBody string) error {
	if !s.cfg.EmailEnabled() {
		log.Printf("[EMAIL] Skipping email to %s - email not configured", to)
		return nil
	}

	fromAddress := fmt.Sprintf("%s <%s>", s.cfg.EmailFromName, s.cfg.EmailFromAddress)

	reqBody := resendRequest{
		From:    fromAddress,
		To:      []string{to},
		Subject: subject,
		HTML:    htmlBody,
		Text:    textBody,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	req, err := http.NewRequest("POST", resendAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp resendError
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("email failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("email failed: %s - %s", errResp.Name, errResp.Message)
	}

	var respBody resendResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		// Email was sent even if we can't parse the response
		log.Printf("[EMAIL] Email sent to %s (unable to parse response)", to)
		return nil
	}

	log.Printf("[EMAIL] Email sent to %s (ID: %s)", to, respBody.ID)
	return nil
}

// SendAsync sends an email asynchronously (non-blocking)
func (s *EmailService) SendAsync(to, subject, htmlBody, textBody string) {
	if !s.cfg.EmailEnabled() {
		return
	}
	go func() {
		if err := s.Send(to, subject, htmlBody, textBody); err != nil {
			log.Printf("[EMAIL ERROR] Failed to send email to %s: %v", to, err)
		}
	}()
}

// SendWelcome sends a welcome email to a new user
func (s *EmailService) SendWelcome(user *domain.User, tempPassword string) {
	data := welcomeEmailData{
		AppURL:       s.cfg.AppURL,
		UserName:     user.Name,
		UserEmail:    user.Email,
		TempPassword: tempPassword,
	}

	htmlBody, err := s.renderTemplate(welcomeEmailHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render welcome email HTML: %v", err)
		return
	}

	textBody, err := s.renderTemplate(welcomeEmailText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render welcome email text: %v", err)
		return
	}

	s.SendAsync(user.Email, welcomeEmailSubject, htmlBody, textBody)
}

// SendRequestSubmitted sends an email when a vacation request is submitted
func (s *EmailService) SendRequestSubmitted(user *domain.User, vacation *domain.VacationRequest) {
	if !user.EmailPreferences.VacationUpdates {
		log.Printf("[EMAIL] Skipping request submitted email for %s - user preferences disabled", user.Email)
		return
	}

	data := vacationEmailData{
		AppURL:    s.cfg.AppURL,
		UserName:  user.Name,
		StartDate: vacation.StartDate,
		EndDate:   vacation.EndDate,
		TotalDays: vacation.TotalDays,
	}

	htmlBody, err := s.renderTemplate(requestSubmittedHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render request submitted email HTML: %v", err)
		return
	}

	textBody, err := s.renderTemplate(requestSubmittedText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render request submitted email text: %v", err)
		return
	}

	s.SendAsync(user.Email, requestSubmittedSubject, htmlBody, textBody)
}

// SendRequestApproved sends an email when a vacation request is approved
func (s *EmailService) SendRequestApproved(user *domain.User, vacation *domain.VacationRequest) {
	if !user.EmailPreferences.VacationUpdates {
		log.Printf("[EMAIL] Skipping approval email for %s - user preferences disabled", user.Email)
		return
	}

	data := vacationEmailData{
		AppURL:    s.cfg.AppURL,
		UserName:  user.Name,
		StartDate: vacation.StartDate,
		EndDate:   vacation.EndDate,
		TotalDays: vacation.TotalDays,
	}

	htmlBody, err := s.renderTemplate(requestApprovedHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render approved email HTML: %v", err)
		return
	}

	textBody, err := s.renderTemplate(requestApprovedText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render approved email text: %v", err)
		return
	}

	s.SendAsync(user.Email, requestApprovedSubject, htmlBody, textBody)
}

// SendRequestRejected sends an email when a vacation request is rejected
func (s *EmailService) SendRequestRejected(user *domain.User, vacation *domain.VacationRequest, reason string) {
	if !user.EmailPreferences.VacationUpdates {
		log.Printf("[EMAIL] Skipping rejection email for %s - user preferences disabled", user.Email)
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

	htmlBody, err := s.renderTemplate(requestRejectedHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render rejected email HTML: %v", err)
		return
	}

	textBody, err := s.renderTemplate(requestRejectedText, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render rejected email text: %v", err)
		return
	}

	s.SendAsync(user.Email, requestRejectedSubject, htmlBody, textBody)
}

// SendAdminNewRequest sends an email to admins when a new vacation request is submitted
func (s *EmailService) SendAdminNewRequest(admins []*domain.User, requester *domain.User, vacation *domain.VacationRequest) {
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

	htmlBody, err := s.renderTemplate(adminNewRequestHTML, data)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to render admin notification email HTML: %v", err)
		return
	}

	textBody, err := s.renderTemplate(adminNewRequestText, data)
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
		s.SendAsync(admin.Email, adminNewRequestSubject, htmlBody, textBody)
	}
}

// renderTemplate renders a template string with the given data
func (s *EmailService) renderTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

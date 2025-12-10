package domain

import (
	"time"
)

// VacationStatus represents the status of a vacation request
type VacationStatus string

const (
	StatusPending  VacationStatus = "pending"
	StatusApproved VacationStatus = "approved"
	StatusRejected VacationStatus = "rejected"
)

// VacationRequest represents an employee's vacation request
type VacationRequest struct {
	ID              string         `json:"id"`
	UserID          string         `json:"userId"`
	UserName        string         `json:"userName,omitempty"`  // Populated from JOIN
	UserEmail       string         `json:"userEmail,omitempty"` // Populated from JOIN
	StartDate       string         `json:"startDate"`           // Format: YYYY-MM-DD
	EndDate         string         `json:"endDate"`             // Format: YYYY-MM-DD
	TotalDays       int            `json:"totalDays"`
	Reason          *string        `json:"reason,omitempty"`
	Status          VacationStatus `json:"status"`
	ReviewedBy      *string        `json:"reviewedBy,omitempty"`
	ReviewedAt      *time.Time     `json:"reviewedAt,omitempty"`
	RejectionReason *string        `json:"rejectionReason,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

// IsPending returns true if the request is pending review
func (v *VacationRequest) IsPending() bool {
	return v.Status == StatusPending
}

// IsApproved returns true if the request has been approved
func (v *VacationRequest) IsApproved() bool {
	return v.Status == StatusApproved
}

// IsRejected returns true if the request has been rejected
func (v *VacationRequest) IsRejected() bool {
	return v.Status == StatusRejected
}

// CanBeCancelled returns true if the request can be cancelled
// Only pending requests can be cancelled
func (v *VacationRequest) CanBeCancelled() bool {
	return v.IsPending()
}

// TeamVacation is a simplified view for team calendar display
type TeamVacation struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	TotalDays int    `json:"totalDays"`
}

// ValidStatuses returns all valid vacation status values
func ValidStatuses() []VacationStatus {
	return []VacationStatus{StatusPending, StatusApproved, StatusRejected}
}

// IsValidStatus checks if a status string is valid
func IsValidStatus(status string) bool {
	for _, s := range ValidStatuses() {
		if string(s) == status {
			return true
		}
	}
	return false
}

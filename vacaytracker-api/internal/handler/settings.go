package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/repository/sqlite"
)

// SettingsHandler handles public settings endpoints
type SettingsHandler struct {
	settingsRepo *sqlite.SettingsRepository
}

// NewSettingsHandler creates a new SettingsHandler
func NewSettingsHandler(settingsRepo *sqlite.SettingsRepository) *SettingsHandler {
	return &SettingsHandler{
		settingsRepo: settingsRepo,
	}
}

// PublicSettingsResponse contains only non-sensitive settings
type PublicSettingsResponse struct {
	DefaultVacationDays int `json:"defaultVacationDays"`
	VacationResetMonth  int `json:"vacationResetMonth"`
}

// GetPublic handles GET /api/settings/public
// Returns non-sensitive application settings (authenticated users only)
func (h *SettingsHandler) GetPublic(c *gin.Context) {
	settings, err := h.settingsRepo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
		return
	}

	c.JSON(http.StatusOK, PublicSettingsResponse{
		DefaultVacationDays: settings.DefaultVacationDays,
		VacationResetMonth:  settings.VacationResetMonth,
	})
}

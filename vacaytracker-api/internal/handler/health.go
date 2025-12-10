package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const version = "1.0.0"

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// Check handles GET /health
// Returns the health status of the API
func (h *HealthHandler) Check(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
	}

	c.JSON(http.StatusOK, response)
}

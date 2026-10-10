package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// HealthHandler handles health check requests
type HealthHandler struct {
	Version string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		Version: version,
	}
}

// Health returns the health check response
func (h *HealthHandler) Health(c *gin.Context) {
	response := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   h.Version,
	}

	c.JSON(http.StatusOK, response)
}

// Ready returns the readiness check response
func (h *HealthHandler) Ready(c *gin.Context) {
	response := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   h.Version,
	}

	c.JSON(http.StatusOK, response)
}

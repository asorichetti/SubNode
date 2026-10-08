package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InviteHandler struct {
	DB *gorm.DB
}

type CreateInviteRequest struct {
	MaxUses int             `json:"max_uses" binding:"omitempty,min=1"`
	Expires *time.Time      `json:"expires,omitempty"`
	Role    models.UserRole `json:"role" binding:"omitempty,oneof=member moderator"`
}

type InviteResponse struct {
	ID        uuid.UUID       `json:"id"`
	Code      string          `json:"code"`
	MaxUses   int             `json:"max_uses"`
	Uses      int             `json:"uses"`
	ExpiresAt *time.Time      `json:"expires_at,omitempty"`
	Role      models.UserRole `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
}

// CreateInvite generates a new invite code (admin/mod only)
func (h *InviteHandler) CreateInvite(c *gin.Context) {
	hostID, _ := c.Get("user_id")
	var req CreateInviteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.MaxUses == 0 {
		req.MaxUses = 1
	}
	if req.Role == "" {
		req.Role = models.RoleMember
	}

	// Generate random code
	code := generateInviteCode()

	invite := models.InviteCode{
		Code:      code,
		HostID:    hostID.(uuid.UUID),
		MaxUses:   req.MaxUses,
		Role:      req.Role,
		ExpiresAt: req.Expires,
	}

	if err := h.DB.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invite"})
		return
	}

	c.JSON(http.StatusCreated, InviteResponse{
		ID:        invite.ID,
		Code:      invite.Code,
		MaxUses:   invite.MaxUses,
		Uses:      invite.Uses,
		ExpiresAt: invite.ExpiresAt,
		Role:      invite.Role,
		CreatedAt: invite.CreatedAt,
	})
}

// ListInvites returns all active invite codes (admin/mod only)
func (h *InviteHandler) ListInvites(c *gin.Context) {
	hostID, _ := c.Get("user_id")
	var invites []models.InviteCode

	if err := h.DB.Where("host_id = ? AND deleted_at IS NULL", hostID).
		Order("created_at DESC").
		Find(&invites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invites"})
		return
	}

	response := make([]InviteResponse, len(invites))
	for i, invite := range invites {
		response[i] = InviteResponse{
			ID:        invite.ID,
			Code:      invite.Code,
			MaxUses:   invite.MaxUses,
			Uses:      invite.Uses,
			ExpiresAt: invite.ExpiresAt,
			Role:      invite.Role,
			CreatedAt: invite.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// RevokeInvite deactivates an invite code (admin/mod only)
func (h *InviteHandler) RevokeInvite(c *gin.Context) {
	inviteID := c.Param("id")
	hostID, _ := c.Get("user_id")

	var invite models.InviteCode
	if err := h.DB.Where("id = ? AND host_id = ?", inviteID, hostID).First(&invite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invite not found"})
		return
	}

	if err := h.DB.Model(&invite).Update("deleted_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke invite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invite revoked"})
}

// generateInviteCode creates a random 8-character alphanumeric code
func generateInviteCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:8]
}

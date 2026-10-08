package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NodeHandler struct {
	DB *gorm.DB
}

type CreateNodeRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Rules       string `json:"rules" binding:"omitempty,max=2000"`
}

type NodeResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	Rules       string        `json:"rules"`
	HostID      uuid.UUID     `json:"host_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// CreateNode creates a new node (admin only)
func (h *NodeHandler) CreateNode(c *gin.Context) {
	hostID, _ := c.Get("user_id")
	var req CreateNodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate slug from name
	slug := generateSlug(req.Name)

	// Check if slug is taken
	var count int64
	h.DB.Model(&models.Node{}).Where("slug = ?", slug).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Node name already taken"})
		return
	}

	node := models.Node{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Rules:       req.Rules,
		HostID:      hostID.(uuid.UUID),
	}

	if err := h.DB.Create(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}

	c.JSON(http.StatusCreated, NodeResponse{
		ID:          node.ID,
		Name:        node.Name,
		Slug:        node.Slug,
		Description: node.Description,
		Rules:       node.Rules,
		HostID:      node.HostID,
		CreatedAt:   node.CreatedAt,
		UpdatedAt:   node.UpdatedAt,
	})
}

// ListNodes returns all public nodes
func (h *NodeHandler) ListNodes(c *gin.Context) {
	var nodes []models.Node

	if err := h.DB.Preload("Host").Order("created_at DESC").Find(&nodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nodes"})
		return
	}

	response := make([]NodeResponse, len(nodes))
	for i, node := range nodes {
		response[i] = NodeResponse{
			ID:          node.ID,
			Name:        node.Name,
			Slug:        node.Slug,
			Description: node.Description,
			Rules:       node.Rules,
			HostID:      node.HostID,
			CreatedAt:   node.CreatedAt,
			UpdatedAt:   node.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetNode returns a single node by slug
func (h *NodeHandler) GetNode(c *gin.Context) {
	slug := c.Param("slug")

	var node models.Node
	if err := h.DB.Where("slug = ?", slug).First(&node).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	c.JSON(http.StatusOK, NodeResponse{
		ID:          node.ID,
		Name:        node.Name,
		Slug:        node.Slug,
		Description: node.Description,
		Rules:       node.Rules,
		HostID:      node.HostID,
		CreatedAt:   node.CreatedAt,
		UpdatedAt:   node.UpdatedAt,
	})
}

// DeleteNode deletes a node (admin/host only)
func (h *NodeHandler) DeleteNode(c *gin.Context) {
	slug := c.Param("slug")
	userID, _ := c.Get("user_id")

	var node models.Node
	if err := h.DB.Where("slug = ?", slug).First(&node).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	// Check if user is host or admin
	if node.HostID != userID && !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := h.DB.Delete(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node deleted"})
}

// generateSlug creates a URL-friendly slug from a name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	for _, ch := range slug {
		if !(ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '-') {
			slug = strings.ReplaceAll(slug, string(ch), "")
		}
	}
	return slug
}

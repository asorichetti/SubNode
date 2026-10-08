package handlers

import (
	"net/http"
	"time"

	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

type CreatePostRequest struct {
	Title string   `json:"title" binding:"required,min=1,max=300"`
	Body  string   `json:"body" binding:"omitempty"`
	URL   *string  `json:"url" binding:"omitempty"`
	Tags  []string `json:"tags"`
}

type PostResponse struct {
	ID        uuid.UUID    `json:"id"`
	NodeID    uuid.UUID    `json:"node_id"`
	AuthorID  uuid.UUID    `json:"author_id"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	URL       *string      `json:"url,omitempty"`
	VoteCount int          `json:"vote_count"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type VoteRequest struct {
	Vote models.VoteType `json:"vote" binding:"required,oneof=up down"`
}

// CreatePost creates a new post in a node
func (h *PostHandler) CreatePost(c *gin.Context) {
	slug := c.Param("slug")
	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if node exists
	var node models.Node
	if err := h.DB.Where("slug = ?", slug).First(&node).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		NodeID:   node.ID,
		AuthorID: authorID.(uuid.UUID),
		Title:    req.Title,
		Body:     req.Body,
		URL:      req.URL,
	}

	if err := h.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Handle tags if provided
	if len(req.Tags) > 0 {
		h.addTagsToPost(c, post.ID, req.Tags)
	}

	// Load author
	var author models.User
	h.DB.First(&author, post.AuthorID)

	c.JSON(http.StatusCreated, PostResponse{
		ID:        post.ID,
		NodeID:    post.NodeID,
		AuthorID:  post.AuthorID,
		Title:     post.Title,
		Body:      post.Body,
		URL:       post.URL,
		VoteCount: post.VoteCount,
		Author: UserResponse{
			ID:       author.ID,
			Username: author.Username,
			Email:    author.Email,
			Role:     author.Role,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

// ListPosts returns posts for a node with pagination
func (h *PostHandler) ListPosts(c *gin.Context) {
	slug := c.Param("slug")

	// Get node
	var node models.Node
	if err := h.DB.Where("slug = ?", slug).First(&node).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	// Pagination
	page := 1
	if p := c.Query("page"); p != "" {
		// Parse page (simplified)
		page = 1
	}
	limit := 20
	if l := c.Query("limit"); l != "" {
		// Parse limit (simplified)
		limit = 20
	}

	var posts []models.Post
	if err := h.DB.Where("node_id = ?", node.ID).
		Preload("Author").
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	response := make([]PostResponse, len(posts))
	for i, post := range posts {
		response[i] = PostResponse{
			ID:        post.ID,
			NodeID:    post.NodeID,
			AuthorID:  post.AuthorID,
			Title:     post.Title,
			Body:      post.Body,
			URL:       post.URL,
			VoteCount: post.VoteCount,
			Author: UserResponse{
				ID:       post.Author.ID,
				Username: post.Author.Username,
				Email:    post.Author.Email,
				Role:     post.Author.Role,
			},
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetPost returns a single post with its comments
func (h *PostHandler) GetPost(c *gin.Context) {
	postID := c.Param("id")

	var post models.Post
	if err := h.DB.Where("id = ?", postID).
		Preload("Author").
		Preload("Comments.Author").
		First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, PostResponse{
		ID:        post.ID,
		NodeID:    post.NodeID,
		AuthorID:  post.AuthorID,
		Title:     post.Title,
		Body:      post.Body,
		URL:       post.URL,
		VoteCount: post.VoteCount,
		Author: UserResponse{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Email:    post.Author.Email,
			Role:     post.Author.Role,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

// DeletePost deletes a post (author or admin only)
func (h *PostHandler) DeletePost(c *gin.Context) {
	postID := c.Param("id")
	userID, _ := c.Get("user_id")

	var post models.Post
	if err := h.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user is author or admin
	userRole, _ := c.Get("user_role")
	isAdmin := userRole == "admin"
	if post.AuthorID != userID && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := h.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// VotePost handles voting on a post
func (h *PostHandler) VotePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	var req VoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if post exists
	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user already voted on this post
	var existingVote models.Vote
	if err := h.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingVote).Error; err == nil {
		// Toggle or update vote
		if existingVote.VoteType == req.Vote {
			// Remove vote
			h.DB.Delete(&existingVote)
			h.recalculateVoteCount(postID)
			c.JSON(http.StatusOK, gin.H{"message": "Vote removed"})
			return
		}
		// Update vote
		existingVote.VoteType = req.Vote
		h.DB.Save(&existingVote)
		h.recalculateVoteCount(postID)
		c.JSON(http.StatusOK, gin.H{"message": "Vote updated"})
		return
	}

	// Create new vote
	vote := models.Vote{
		UserID:   userID.(uuid.UUID),
		PostID:   &postID,
		VoteType: req.Vote,
	}
	h.DB.Create(&vote)
	h.recalculateVoteCount(postID)
	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded"})
}

func (h *PostHandler) addTagsToPost(c *gin.Context, postID uuid.UUID, tags []string) {
	for _, tagName := range tags {
		var tag models.Tag
		result := h.DB.Where("name = ?", tagName).First(&tag)
		if result.Error != nil {
			// Create tag if it doesn't exist
			tag = models.Tag{Name: tagName}
			h.DB.Create(&tag)
		}
		// Add post-tag association
		h.DB.Model(&models.Post{}).Where("id = ?", postID).Association("Tags").Append(&tag)
	}
}

func (h *PostHandler) recalculateVoteCount(postID uuid.UUID) {
	var upVotes, downVotes int64
	h.DB.Model(&models.Vote{}).Where("post_id = ? AND vote_type = 'up'", postID).Count(&upVotes)
	h.DB.Model(&models.Vote{}).Where("post_id = ? AND vote_type = 'down'", postID).Count(&downVotes)
	newCount := int(upVotes - downVotes)

	h.DB.Model(&models.Post{}).Where("id = ?", postID).Update("vote_count", newCount)
}

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB     *gorm.DB
	Secret string
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Code     string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
}

// Register handles user registration with invite code
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate invite code
	var invite models.InviteCode
	if err := h.DB.Where("code = ? AND deleted_at IS NULL", req.Code).First(&invite).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or expired invite code"})
		return
	}

	// Check if invite has remaining uses
	if invite.MaxUses > 0 && invite.Uses >= invite.MaxUses {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invite code has reached maximum uses"})
		return
	}

	// Check if invite has expired
	if invite.ExpiresAt != nil && time.Now().After(*invite.ExpiresAt) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invite code has expired"})
		return
	}

	// Check if username is taken
	var existingUser models.User
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Check if email is taken
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     invite.Role,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Increment invite code usage
	invite.Uses++
	if err := h.DB.Save(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invite"})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	u := user.(*models.User)
	c.JSON(http.StatusOK, UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	})
}

// Logout is handled client-side (just invalidate token)
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// generateToken creates a JWT token for the user
func (h *AuthHandler) generateToken(id uuid.UUID, username string, role models.UserRole) (string, error) {
	// Simple token generation - in production, use a proper JWT library
	// This is a simplified version; replace with jwt-go or similar
	token := fmt.Sprintf("%s.%s.%s", id.String(), username, h.Secret)
	return token, nil
}

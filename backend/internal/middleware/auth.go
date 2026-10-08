package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	DB *gorm.DB
}

// JWTAuth validates JWT tokens and sets user context
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Simplified token validation
		// In production, use proper JWT library (jwt-go, etc.)
		userID, err := m.validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Load user from database
		var user models.User
		if err := m.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", &user)
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}

// RequireAdmin checks if the current user is an admin
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireModerator checks if the current user is admin or moderator
func (m *AuthMiddleware) RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || (role != models.RoleAdmin && role != models.RoleModerator) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Moderator access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// validateToken extracts user ID from token
// This is a simplified version - replace with proper JWT validation
func (m *AuthMiddleware) validateToken(token string) (uuid.UUID, error) {
	// Simple token format: userID.username.secret
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	userID, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}

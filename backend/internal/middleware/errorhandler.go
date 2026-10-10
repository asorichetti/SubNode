package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorHandler returns a middleware that handles errors gracefully
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors in the context
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if je, ok := err.Err.(*jsonError); ok {
					c.JSON(je.Code, gin.H{"error": je.Message})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				}
			}
			c.Abort()
			return
		}
	}
}

// jsonError is a custom error type for JSON responses
type jsonError struct {
	Code    int
	Message string
}

func (e *jsonError) Error() string {
	return e.Message
}

// JSONError sets a JSON error in the context
func JSONError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}

// BindError handles binding errors
func BindError(c *gin.Context, err error) {
	msg := "Invalid request"
	if strings.Contains(err.Error(), "Password") {
		msg = "Invalid password"
	} else if strings.Contains(err.Error(), "Username") {
		msg = "Invalid username"
	} else if strings.Contains(err.Error(), "Email") {
		msg = "Invalid email"
	} else if strings.Contains(err.Error(), "Code") {
		msg = "Invalid code"
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

// NotFoundError handles not found errors
func NotFoundError(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, gin.H{"error": resource + " not found"})
}

// AuthError handles authentication errors
func AuthError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
}

// ForbiddenError handles authorization errors
func ForbiddenError(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
}

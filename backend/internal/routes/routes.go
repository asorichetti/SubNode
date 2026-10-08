package routes

import (
	"github.com/asorichetti/subnode/backend/internal/handlers"
	"github.com/asorichetti/subnode/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup configures all routes for the application
func Setup(db *gorm.DB, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	// Initialize handlers
	authHandler := &handlers.AuthHandler{
		DB:     db,
		Secret: "subnode-secret", // TODO: Load from config
	}
	inviteHandler := &handlers.InviteHandler{DB: db}
	nodeHandler := &handlers.NodeHandler{DB: db}
	postHandler := &handlers.PostHandler{DB: db}

	// Public routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public node routes
		nodes := api.Group("/nodes")
		{
			nodes.GET("", nodeHandler.ListNodes)
			nodes.GET("/:slug", nodeHandler.GetNode)
		}

		// Public post listing
		api.GET("/nodes/:slug/posts", postHandler.ListPosts)
		api.GET("/posts/:id", postHandler.GetPost)
	}

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(authMiddleware.JWTAuth())
	{
		// User routes
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Logout)

		// Invite routes (admin/mod only)
		invites := protected.Group("/invites")
		invites.Use(authMiddleware.RequireAdmin())
		{
			invites.POST("", inviteHandler.CreateInvite)
			invites.GET("", inviteHandler.ListInvites)
			invites.DELETE("/:id", inviteHandler.RevokeInvite)
		}

		// Node routes
		nodes := protected.Group("/nodes")
		{
			nodes.POST("", authMiddleware.RequireAdmin(), nodeHandler.CreateNode)
			nodes.DELETE("/:slug", nodeHandler.DeleteNode)
		}

		// Post routes
		posts := protected.Group("/nodes/:slug/posts")
		{
			posts.POST("", postHandler.CreatePost)
		}

		protected.POST("/posts/:id/vote", postHandler.VotePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

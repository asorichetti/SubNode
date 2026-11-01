package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"github.com/asorichetti/subnode/backend/internal/config"
	"github.com/asorichetti/subnode/backend/internal/database"
	"github.com/asorichetti/subnode/backend/internal/handlers"
	"github.com/asorichetti/subnode/backend/internal/middleware"
	"github.com/asorichetti/subnode/backend/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Resolve admin credentials for the initial seed. The default password
	// is only tolerated outside release mode.
	adminUsername, adminPassword := cfg.Admin.Username, cfg.Admin.Password
	if adminUsername == "" {
		log.Fatal("ADMIN_USERNAME must not be empty")
	}
	if adminPassword == "" {
		if cfg.Server.Mode == "release" {
			log.Fatal("ADMIN_PASSWORD must be set in release mode")
		}
		adminPassword = "admin"
		log.Println(`WARNING: using default admin password "admin" (development only)`)
	}
	if adminPassword == "admin" && cfg.Server.Mode == "release" {
		log.Fatal(`refusing to start in release mode with the default admin password "admin"`)
	}

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	// Seed database with default data
	if err := database.Seed(db, adminUsername, adminPassword); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	log.Println("Database seeded")

	// Initialize middleware with JWT
	tokenExpiry := time.Duration(cfg.Auth.TokenExpiry) * time.Minute
	authMiddleware := middleware.NewAuthMiddleware(db, cfg.Auth.JWTSecret, tokenExpiry)

	// Initialize handlers with JWT
	authHandler := handlers.NewAuthHandler(db, cfg.Auth.JWTSecret, tokenExpiry)

	// Setup routes
	router := routes.Setup(db, authMiddleware, authHandler,
		cfg.CORS.AllowedOrigins,
		cfg.RateLimit.Enabled,
		cfg.RateLimit.Requests,
		cfg.RateLimit.Window,
	)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s (mode: %s)", addr, cfg.Server.Mode)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/asorichetti/subnode/backend/internal/config"
	"github.com/asorichetti/subnode/backend/internal/database"
	"github.com/asorichetti/subnode/backend/internal/handlers"
	"github.com/asorichetti/subnode/backend/internal/middleware"
	"github.com/asorichetti/subnode/backend/internal/routes"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	// Seed database with default data
	if err := database.Seed(db, cfg); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	log.Println("Database seeded")

	// Initialize middleware with JWT
	tokenExpiry := time.Duration(cfg.Auth.TokenExpiry) * time.Minute
	authMiddleware := middleware.NewAuthMiddleware(db, cfg.Auth.JWTSecret, tokenExpiry)

	// Initialize handlers with JWT
	authHandler := handlers.NewAuthHandler(db, cfg.Auth.JWTSecret, tokenExpiry)

	// Setup routes
	router := routes.Setup(db, authMiddleware, authHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s (mode: %s)", addr, cfg.Server.Mode)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

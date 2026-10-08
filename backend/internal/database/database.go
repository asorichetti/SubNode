package database

import (
	"fmt"

	"github.com/asorichetti/subnode/backend/internal/config"
	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var defaultModels = []interface{}{
	&models.User{},
	&models.Node{},
	&models.Post{},
	&models.Comment{},
	&models.Vote{},
	&models.InviteCode{},
	&models.Tag{},
	&models.PostTag{},
	&models.UserNode{},
	&models.Moderator{},
}

// New creates a new database connection and runs migrations
func New(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	var err error

	switch cfg.Database.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.Database.GetDSN())
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.SQLitePath)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(defaultModels...)
}

// Seed creates default admin user if none exist
func Seed(db *gorm.DB, cfg *config.Config) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Create default admin user
	admin := models.User{
		Username: "admin",
		Email:    "admin@localhost",
		Role:     models.RoleAdmin,
		Password: "admin", // Will be hashed on actual login setup
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// Create default node
	defaultNode := models.Node{
		Name:        "Welcome",
		Slug:        "welcome",
		Description: "Welcome to your SubNode! This is your first node.",
		Rules:       "1. Be kind\n2. Have fun\n3. Read the rules",
		HostID:      admin.ID,
	}

	if err := db.Create(&defaultNode).Error; err != nil {
		return fmt.Errorf("failed to create default node: %w", err)
	}

	// Create welcome post
	welcomePost := models.Post{
		NodeID:   defaultNode.ID,
		AuthorID: admin.ID,
		Title:    "Welcome to SubNode!",
		Body:     "This is your first post. Edit or delete it as you wish.",
	}

	return db.Create(&welcomePost).Error
}

// UUID generates a new UUID (for use in seed data)
func UUID() uuid.UUID {
	return uuid.New()
}

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asorichetti/subnode/backend/internal/jwt"
	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables without UUID defaults (SQLite doesn't support gen_random_uuid)
	tables := []string{
		`CREATE TABLE users (id TEXT PRIMARY KEY, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL UNIQUE, password TEXT NOT NULL, role TEXT DEFAULT "member", created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE invite_codes (id TEXT PRIMARY KEY, code TEXT NOT NULL UNIQUE, host_id TEXT NOT NULL, max_uses INTEGER DEFAULT 1, uses INTEGER DEFAULT 0, expires_at DATETIME, role TEXT DEFAULT "member", created_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE nodes (id TEXT PRIMARY KEY, host_id TEXT NOT NULL, name TEXT NOT NULL, slug TEXT NOT NULL UNIQUE, description TEXT, rules TEXT, created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE posts (id TEXT PRIMARY KEY, node_id TEXT NOT NULL, author_id TEXT NOT NULL, title TEXT NOT NULL, body TEXT, url TEXT, vote_count INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE votes (id TEXT PRIMARY KEY, user_id TEXT NOT NULL, post_id TEXT, comment_id TEXT, vote_type TEXT NOT NULL, created_at DATETIME)`,
	}

	for _, table := range tables {
		if err := db.Exec(table).Error; err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
	}

	return db
}

func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	t.Helper()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     models.RoleMember,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

func createInviteCode(t *testing.T, db *gorm.DB, hostID uuid.UUID) *models.InviteCode {
	t.Helper()
	// Use raw SQL to avoid UUID issues
	inviteID := uuid.New().String()
	if err := db.Exec("INSERT INTO invite_codes (id, code, host_id, max_uses, uses, role, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		inviteID, "test-invite-code", hostID.String(), 5, 0, "member", time.Now()).Error; err != nil {
		t.Fatalf("Failed to create invite code: %v", err)
	}

	return &models.InviteCode{
		ID:      hostID, // dummy
		Code:    "test-invite-code",
		HostID:  hostID,
		MaxUses: 5,
		Uses:    0,
		Role:    models.RoleMember,
	}
}

func createRouter(db *gorm.DB, handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/register", handler.Register)
	auth.GET("/me", handler.Me)
	return r
}

func TestAuthHandler_Login(t *testing.T) {
	db := setupTestDB(t)
	tokenGen := jwt.NewTokenGenerator("test-secret", 1*time.Hour)
	handler := &AuthHandler{DB: db, tokenGen: tokenGen}
	createTestUser(t, db)

	r := createRouter(db, handler)

	reqBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "testpassword",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	var response AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token in response")
	}

	if response.User.Username != "testuser" {
		t.Errorf("Expected username testuser, got %s", response.User.Username)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	db := setupTestDB(t)
	tokenGen := jwt.NewTokenGenerator("test-secret", 1*time.Hour)
	handler := &AuthHandler{DB: db, tokenGen: tokenGen}
	createTestUser(t, db)

	r := createRouter(db, handler)

	reqBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Register(t *testing.T) {
	db := setupTestDB(t)
	tokenGen := jwt.NewTokenGenerator("test-secret", 1*time.Hour)
	handler := &AuthHandler{DB: db, tokenGen: tokenGen}

	hostID := uuid.New()
	createInviteCode(t, db, hostID)

	r := createRouter(db, handler)

	reqBody, _ := json.Marshal(map[string]string{
		"username": "newuser",
		"email":    "new@example.com",
		"password": "newpassword",
		"code":     "test-invite-code",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	var response AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token in response")
	}

	if response.User.Username != "newuser" {
		t.Errorf("Expected username newuser, got %s", response.User.Username)
	}
}

func TestAuthHandler_Register_InvalidInvite(t *testing.T) {
	db := setupTestDB(t)
	tokenGen := jwt.NewTokenGenerator("test-secret", 1*time.Hour)
	handler := &AuthHandler{DB: db, tokenGen: tokenGen}

	r := createRouter(db, handler)

	reqBody, _ := json.Marshal(map[string]string{
		"username": "newuser",
		"email":    "new@example.com",
		"password": "newpassword",
		"code":     "invalid-code",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAuthHandler_Me(t *testing.T) {
	db := setupTestDB(t)
	tokenGen := jwt.NewTokenGenerator("test-secret", 1*time.Hour)
	handler := &AuthHandler{DB: db, tokenGen: tokenGen}
	testUser := createTestUser(t, db)

	token, _ := tokenGen.GenerateToken(testUser.ID, testUser.Username, testUser.Role)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Manually set the user in context (simulating auth middleware)
	c.Set("user", testUser)
	c.Set("user_id", testUser.ID)
	c.Set("user_role", testUser.Role)

	c.Request, _ = http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler.Me(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Username != "testuser" {
		t.Errorf("Expected username testuser, got %s", response.Username)
	}
}

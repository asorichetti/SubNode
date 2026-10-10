package jwt

import (
	"testing"
	"time"

	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/google/uuid"
)

func TestTokenGenerator_GenerateToken(t *testing.T) {
	secretKey := "test-secret-key"
	expiry := 1 * time.Hour

	gen := NewTokenGenerator(secretKey, expiry)
	userID := uuid.New()
	username := "testuser"
	role := models.RoleMember

	token, err := gen.GenerateToken(userID, username, role)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}
}

func TestTokenGenerator_ValidateToken(t *testing.T) {
	secretKey := "test-secret-key"
	expiry := 1 * time.Hour

	gen := NewTokenGenerator(secretKey, expiry)
	userID := uuid.New()
	username := "testuser"
	role := models.RoleMember

	// Generate token
	token, err := gen.GenerateToken(userID, username, role)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate token
	claims, err := gen.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}

	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}

func TestTokenGenerator_ValidateToken_InvalidSignature(t *testing.T) {
	secretKey := "test-secret-key"
	expiry := 1 * time.Hour

	gen := NewTokenGenerator(secretKey, expiry)

	// Generate token with different secret
	gen2 := NewTokenGenerator("wrong-secret", expiry)
	token, _ := gen2.GenerateToken(uuid.New(), "user", models.RoleMember)

	// Try to validate with original secret
	_, err := gen.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error for invalid signature, got nil")
	}
}

func TestTokenGenerator_ValidateToken_Expired(t *testing.T) {
	secretKey := "test-secret-key"

	// Generate token with past expiry
	gen := NewTokenGenerator(secretKey, -1*time.Hour)
	token, _ := gen.GenerateToken(uuid.New(), "user", models.RoleMember)

	// Validate token (should fail because it's expired)
	_, err := gen.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}
}

func TestTokenGenerator_ValidateToken_InvalidToken(t *testing.T) {
	secretKey := "test-secret-key"
	expiry := 1 * time.Hour

	gen := NewTokenGenerator(secretKey, expiry)

	// Try to validate invalid token
	_, err := gen.ValidateToken("invalid.token.here")
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}

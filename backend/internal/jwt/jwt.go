package jwt

import (
	"errors"
	"time"

	"github.com/asorichetti/subnode/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims with user information
type Claims struct {
	UserID   uuid.UUID       `json:"user_id"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// TokenGenerator handles JWT token creation and validation
type TokenGenerator struct {
	secretKey string
	expiry    time.Duration
}

// NewTokenGenerator creates a new TokenGenerator
func NewTokenGenerator(secretKey string, expiry time.Duration) *TokenGenerator {
	return &TokenGenerator{
		secretKey: secretKey,
		expiry:    expiry,
	}
}

// GenerateToken creates a signed JWT token for the user
func (t *TokenGenerator) GenerateToken(userID uuid.UUID, username string, role models.UserRole) (string, error) {
	now := time.Now()
	expiry := now.Add(t.expiry)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "subnode",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.secretKey))
}

// ValidateToken parses and validates a JWT token, returning the claims
func (t *TokenGenerator) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	cfg := Load()

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Server.Mode != "release" {
		t.Errorf("Expected mode release, got %s", cfg.Server.Mode)
	}

	if cfg.Database.Driver != "sqlite" {
		t.Errorf("Expected driver sqlite, got %s", cfg.Database.Driver)
	}

	if cfg.Auth.JWTSecret != "change-this-in-production" {
		t.Errorf("Expected default JWT secret, got %s", cfg.Auth.JWTSecret)
	}

	if cfg.Auth.TokenExpiry != 60 {
		t.Errorf("Expected token expiry 60, got %d", cfg.Auth.TokenExpiry)
	}

	if cfg.CORS.AllowedOrigins != "*" {
		t.Errorf("Expected CORS allowed origins *, got %s", cfg.CORS.AllowedOrigins)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_MODE", "debug")
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("JWT_SECRET", "custom-secret")
	os.Setenv("JWT_EXPIRY", "120")
	os.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com")
	os.Setenv("RATE_LIMIT_ENABLED", "false")
	os.Setenv("RATE_LIMIT_REQUESTS", "50")
	os.Setenv("RATE_LIMIT_WINDOW", "30")

	cfg := Load()

	if cfg.Server.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Server.Port)
	}

	if cfg.Server.Mode != "debug" {
		t.Errorf("Expected mode debug, got %s", cfg.Server.Mode)
	}

	if cfg.Database.Driver != "postgres" {
		t.Errorf("Expected driver postgres, got %s", cfg.Database.Driver)
	}

	if cfg.Auth.JWTSecret != "custom-secret" {
		t.Errorf("Expected custom JWT secret, got %s", cfg.Auth.JWTSecret)
	}

	if cfg.Auth.TokenExpiry != 120 {
		t.Errorf("Expected token expiry 120, got %d", cfg.Auth.TokenExpiry)
	}

	if cfg.CORS.AllowedOrigins != "https://example.com" {
		t.Errorf("Expected CORS allowed origins https://example.com, got %s", cfg.CORS.AllowedOrigins)
	}

	if cfg.RateLimit.Enabled {
		t.Error("Expected rate limit disabled")
	}

	if cfg.RateLimit.Requests != 50 {
		t.Errorf("Expected rate limit requests 50, got %d", cfg.RateLimit.Requests)
	}

	if cfg.RateLimit.Window != 30*time.Minute {
		t.Errorf("Expected rate limit window 30m, got %v", cfg.RateLimit.Window)
	}

	// Clean up
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SERVER_MODE")
	os.Unsetenv("DB_DRIVER")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRY")
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	os.Unsetenv("RATE_LIMIT_ENABLED")
	os.Unsetenv("RATE_LIMIT_REQUESTS")
	os.Unsetenv("RATE_LIMIT_WINDOW")
}

func TestGetDSN_SQLite(t *testing.T) {
	cfg := &DatabaseConfig{
		Driver:     "sqlite",
		SQLitePath: "test.db",
	}

	dsn := cfg.GetDSN()
	if dsn != "test.db" {
		t.Errorf("Expected DSN test.db, got %s", dsn)
	}
}

func TestGetDSN_Postgres(t *testing.T) {
	cfg := &DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "pass",
		Name:     "db",
		SSLMode:  "disable",
	}

	dsn := cfg.GetDSN()
	expected := "host=localhost port=5432 user=user password=pass dbname=db sslmode=disable"
	if dsn != expected {
		t.Errorf("Expected DSN %s, got %s", expected, dsn)
	}
}

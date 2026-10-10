package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Auth      AuthConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
	Mode string // "debug", "release", "test"
}

type DatabaseConfig struct {
	Driver     string // "postgres" or "sqlite"
	Host       string
	Port       string
	Name       string
	User       string
	Password   string
	SSLMode    string
	SQLitePath string // For SQLite mode
}

type AuthConfig struct {
	JWTSecret   string
	TokenExpiry int // minutes
}

type CORSConfig struct {
	AllowedOrigins string // comma-separated list of allowed origins
}

type RateLimitConfig struct {
	Enabled  bool
	Requests int           // max requests per window
	Window   time.Duration // time window for rate limiting
}

func Load() *Config {
	// Parse rate limit config
	rlEnabled := getEnv("RATE_LIMIT_ENABLED", "true")
	rlRequests := GetEnvInt("RATE_LIMIT_REQUESTS", 100)
	rlWindowMinutes := GetEnvInt("RATE_LIMIT_WINDOW", 15)

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "release"),
		},
		Database: DatabaseConfig{
			Driver:     getEnv("DB_DRIVER", "sqlite"),
			Host:       getEnv("DB_HOST", ""),
			Port:       getEnv("DB_PORT", "5432"),
			Name:       getEnv("DB_NAME", "subnode"),
			User:       getEnv("DB_USER", ""),
			Password:   getEnv("DB_PASSWORD", ""),
			SSLMode:    getEnv("DB_SSLMODE", "disable"),
			SQLitePath: getEnv("DB_SQLITE_PATH", "subnode.db"),
		},
		Auth: AuthConfig{
			JWTSecret:   getEnv("JWT_SECRET", "change-this-in-production"),
			TokenExpiry: GetEnvInt("JWT_EXPIRY", 60), // 60 minutes
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
		RateLimit: RateLimitConfig{
			Enabled:  rlEnabled == "true",
			Requests: rlRequests,
			Window:   time.Duration(rlWindowMinutes) * time.Minute,
		},
	}
}

func (c *DatabaseConfig) GetDSN() string {
	if c.Driver == "sqlite" {
		return c.SQLitePath
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	val := getEnv(key, strconv.Itoa(defaultVal))
	if num, err := strconv.Atoi(val); err == nil {
		return num
	}
	return defaultVal
}

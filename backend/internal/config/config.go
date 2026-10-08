package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
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

func Load() *Config {
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
			TokenExpiry: 60, // 60 minutes
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

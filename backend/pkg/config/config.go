package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't return error if it doesn't exist
		fmt.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "flow"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		},
	}

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

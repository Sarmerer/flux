package common

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	WebOrigin   string
}

func LoadConfig() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/flow?sslmode=disable"),
		WebOrigin:   getEnv("WEB_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

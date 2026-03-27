package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	CORSOrigins []string
}

func Load() Config {
	return Config{
		Port:        envOrDefault("PORT", "8080"),
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://marketplace:marketplace@localhost:5432/marketplace?sslmode=disable"),
		JWTSecret:   envOrDefault("JWT_SECRET", "not-a-real-secret"),
		CORSOrigins: strings.Split(envOrDefault("CORS_ORIGINS", "http://localhost:5173"), ","),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

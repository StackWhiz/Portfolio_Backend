package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	Port        string
	RateLimit   int
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/portfolio_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:        getEnv("PORT", "8080"),
		RateLimit:   getEnvAsInt("RATE_LIMIT", 100),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

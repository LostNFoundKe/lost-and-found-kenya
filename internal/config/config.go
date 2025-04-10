package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// Server settings
	Port        int
	Environment string
	LogLevel    string

	// Database settings
	DatabaseURL string

	// Authentication settings
	JWTSecret     string
	JWTExpiration int // in hours

	// Google Cloud Storage settings
	GCSBucketName      string
	GCSProjectID       string
	GCSCredentialsFile string

	// Optional Redis cache
	RedisURL string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "9080"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	jwtExpiration, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration: %w", err)
	}

	return &Config{
		Port:        port,
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/lostandfound?sslmode=disable"),

		JWTSecret:     getEnv("JWT_SECRET", "your-256-bit-secret"),
		JWTExpiration: jwtExpiration,

		GCSBucketName:      getEnv("GCS_BUCKET_NAME", "lostandfound-kenya"),
		GCSProjectID:       getEnv("GCS_PROJECT_ID", ""),
		GCSCredentialsFile: getEnv("GCS_CREDENTIALS_FILE", ""),

		// Optional Redis cache
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn int // minutes
}

func Load() *Config {
	return &Config{
		Port:         getEnv("USER_SERVICE_PORT", "8081"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://dratbo:P@ssw0rd@localhost:5432/satisfactory?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "Bib233asd18-"),
		JWTExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 60*24),
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
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

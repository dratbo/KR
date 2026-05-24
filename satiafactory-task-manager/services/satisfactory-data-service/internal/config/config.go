package config

import (
	"os"
)

type Config struct {
	Port         string
	DatabaseURL  string
	DataFilePath string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("DATA_SERVICE_PORT", "8083"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://dratbo:P@ssw0rd@localhost:5432/satisfactory?sslmode=disable"),
		DataFilePath: getEnv("DATA_FILE_PATH", "./data/Docs.json"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

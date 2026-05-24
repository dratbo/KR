package config

import (
	"os"
)

type Config struct {
	Port           string
	UserServiceURL string
	TaskServiceURL string
	DataServiceURL string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("GATEWAY_PORT", "8080"),
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
		TaskServiceURL: getEnv("TASK_SERVICE_URL", "http://nginx:8090"),
		DataServiceURL: getEnv("DATA_SERVICE_URL", "http://localhost:8083"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

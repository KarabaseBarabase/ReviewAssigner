package config

import "os"

type Config struct {
	DatabaseURL string
	ServerPort  string
	Environment string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getDatabaseURL(),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getDatabaseURL() string {
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}

	return "postgres://review_user:review_pass@localhost:5432/review_service?sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

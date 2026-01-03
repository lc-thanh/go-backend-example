package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	Port          string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "go_backend_db"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		Port:          getEnv("PORT", "8080"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}
}

// getEnv gets environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

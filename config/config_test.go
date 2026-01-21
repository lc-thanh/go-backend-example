package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalVars := map[string]string{
		"DB_HOST":        os.Getenv("DB_HOST"),
		"DB_PORT":        os.Getenv("DB_PORT"),
		"DB_USER":        os.Getenv("DB_USER"),
		"DB_PASSWORD":    os.Getenv("DB_PASSWORD"),
		"DB_NAME":        os.Getenv("DB_NAME"),
		"JWT_SECRET":     os.Getenv("JWT_SECRET"),
		"PORT":           os.Getenv("PORT"),
		"REDIS_HOST":     os.Getenv("REDIS_HOST"),
		"REDIS_PORT":     os.Getenv("REDIS_PORT"),
		"REDIS_PASSWORD": os.Getenv("REDIS_PASSWORD"),
		"REDIS_DB":       os.Getenv("REDIS_DB"),
	}

	// Restore env vars after test
	defer func() {
		for key, value := range originalVars {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	t.Run("Load config with default values", func(t *testing.T) {
		// Clear all env vars
		for key := range originalVars {
			os.Unsetenv(key)
		}

		cfg := LoadConfig()

		if cfg.DBHost != "localhost" {
			t.Errorf("DBHost = %v, want localhost", cfg.DBHost)
		}
		if cfg.DBPort != "5432" {
			t.Errorf("DBPort = %v, want 5432", cfg.DBPort)
		}
		if cfg.DBUser != "postgres" {
			t.Errorf("DBUser = %v, want postgres", cfg.DBUser)
		}
		if cfg.DBPassword != "postgres" {
			t.Errorf("DBPassword = %v, want postgres", cfg.DBPassword)
		}
		if cfg.DBName != "go_backend_db" {
			t.Errorf("DBName = %v, want go_backend_db", cfg.DBName)
		}
		if cfg.Port != "8080" {
			t.Errorf("Port = %v, want 8080", cfg.Port)
		}
		if cfg.RedisHost != "localhost" {
			t.Errorf("RedisHost = %v, want localhost", cfg.RedisHost)
		}
		if cfg.RedisPort != "6379" {
			t.Errorf("RedisPort = %v, want 6379", cfg.RedisPort)
		}
		if cfg.RedisDB != 0 {
			t.Errorf("RedisDB = %v, want 0", cfg.RedisDB)
		}
	})

	t.Run("Load config with custom values", func(t *testing.T) {
		// Set custom env vars
		os.Setenv("DB_HOST", "custom-host")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("DB_USER", "custom-user")
		os.Setenv("DB_PASSWORD", "custom-password")
		os.Setenv("DB_NAME", "custom-db")
		os.Setenv("JWT_SECRET", "custom-secret")
		os.Setenv("PORT", "9000")
		os.Setenv("REDIS_HOST", "redis-host")
		os.Setenv("REDIS_PORT", "6380")
		os.Setenv("REDIS_PASSWORD", "redis-pass")
		os.Setenv("REDIS_DB", "5")

		cfg := LoadConfig()

		if cfg.DBHost != "custom-host" {
			t.Errorf("DBHost = %v, want custom-host", cfg.DBHost)
		}
		if cfg.DBPort != "3306" {
			t.Errorf("DBPort = %v, want 3306", cfg.DBPort)
		}
		if cfg.DBUser != "custom-user" {
			t.Errorf("DBUser = %v, want custom-user", cfg.DBUser)
		}
		if cfg.DBPassword != "custom-password" {
			t.Errorf("DBPassword = %v, want custom-password", cfg.DBPassword)
		}
		if cfg.DBName != "custom-db" {
			t.Errorf("DBName = %v, want custom-db", cfg.DBName)
		}
		if cfg.JWTSecret != "custom-secret" {
			t.Errorf("JWTSecret = %v, want custom-secret", cfg.JWTSecret)
		}
		if cfg.Port != "9000" {
			t.Errorf("Port = %v, want 9000", cfg.Port)
		}
		if cfg.RedisHost != "redis-host" {
			t.Errorf("RedisHost = %v, want redis-host", cfg.RedisHost)
		}
		if cfg.RedisPort != "6380" {
			t.Errorf("RedisPort = %v, want 6380", cfg.RedisPort)
		}
		if cfg.RedisPassword != "redis-pass" {
			t.Errorf("RedisPassword = %v, want redis-pass", cfg.RedisPassword)
		}
		if cfg.RedisDB != 5 {
			t.Errorf("RedisDB = %v, want 5", cfg.RedisDB)
		}
	})
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "Get existing env var",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "Get non-existing env var with default",
			key:          "NON_EXISTING_KEY",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
		{
			name:         "Get empty env var",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

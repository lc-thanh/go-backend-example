package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name   string
		userID uint
		email  string
		role   string
	}{
		{
			name:   "Generate token for regular user",
			userID: 1,
			email:  "test@example.com",
			role:   "user",
		},
		{
			name:   "Generate token for admin",
			userID: 2,
			email:  "admin@example.com",
			role:   "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.userID, tt.email, tt.role)
			if err != nil {
				t.Errorf("GenerateJWT() error = %v", err)
				return
			}
			if token == "" {
				t.Error("GenerateJWT() returned empty token")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Generate a valid token
	validToken, _ := GenerateJWT(1, "test@example.com", "user")

	// Generate an expired token
	claims := Claims{
		UserID: 1,
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, _ := expiredTokenObj.SignedString([]byte("test-secret-key"))

	tests := []struct {
		name      string
		token     string
		wantError bool
	}{
		{
			name:      "Valid token",
			token:     validToken,
			wantError: false,
		},
		{
			name:      "Invalid token",
			token:     "invalid.token.here",
			wantError: true,
		},
		{
			name:      "Expired token",
			token:     expiredToken,
			wantError: true,
		},
		{
			name:      "Empty token",
			token:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateJWT(tt.token)
			if tt.wantError {
				if err == nil {
					t.Error("ValidateJWT() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateJWT() unexpected error = %v", err)
					return
				}
				if claims == nil {
					t.Error("ValidateJWT() returned nil claims")
					return
				}
				if claims.Email != "test@example.com" {
					t.Errorf("ValidateJWT() email = %v, want test@example.com", claims.Email)
				}
			}
		})
	}
}

func TestGenerateJWTWithoutSecret(t *testing.T) {
	// Ensure JWT_SECRET is not set
	os.Unsetenv("JWT_SECRET")

	token, err := GenerateJWT(1, "test@example.com", "user")
	if err != nil {
		t.Errorf("GenerateJWT() error = %v", err)
		return
	}
	if token == "" {
		t.Error("GenerateJWT() returned empty token")
	}

	// Should use default secret
	claims, err := ValidateJWT(token)
	if err != nil {
		t.Errorf("ValidateJWT() error = %v", err)
		return
	}
	if claims.UserID != 1 {
		t.Errorf("ValidateJWT() userID = %v, want 1", claims.UserID)
	}
}

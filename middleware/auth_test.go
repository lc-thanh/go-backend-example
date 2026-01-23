package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go_backend/models"
	"go_backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Set test mode
	gin.SetMode(gin.TestMode)

	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name           string
		authHeader     string
		setupToken     bool
		expectedStatus int
		expectAbort    bool
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer ",
			setupToken:     true,
			expectedStatus: http.StatusOK,
			expectAbort:    false,
		},
		{
			name:           "Missing authorization header",
			authHeader:     "",
			setupToken:     false,
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "Invalid header format - no Bearer",
			authHeader:     "InvalidToken",
			setupToken:     false,
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "Invalid header format - only Bearer",
			authHeader:     "Bearer",
			setupToken:     false,
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid.token.here",
			setupToken:     false,
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Generate valid token if needed
			var token string
			if tt.setupToken {
				var err error
				token, err = utils.GenerateJWT(1, "test@example.com", "user")
				assert.NoError(t, err)
				tt.authHeader = "Bearer " + token
			}

			// Create request
			c.Request = httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				c.Request.Header.Set("Authorization", tt.authHeader)
			}

			// Execute middleware
			AuthMiddleware()(c)

			// Assertions
			if tt.expectAbort {
				assert.True(t, c.IsAborted())
			} else {
				assert.False(t, c.IsAborted())
				// Check if user context is set (note: code uses "userID" not "user_id")
				userID, exists := c.Get("userID")
				assert.True(t, exists)
				assert.Equal(t, uint(1), userID)
			}
		})
	}
}

func TestAuthMiddlewareWithRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name     string
		userRole string
	}{
		{
			name:     "User role",
			userRole: "user",
		},
		{
			name:     "Admin role",
			userRole: "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate token with specific role
			token, err := utils.GenerateJWT(1, "test@example.com", tt.userRole)
			assert.NoError(t, err)

			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Request.Header.Set("Authorization", "Bearer "+token)

			// Execute middleware
			AuthMiddleware()(c)

			// Assertions
			assert.False(t, c.IsAborted())
			role, exists := c.Get("role")
			assert.True(t, exists)
			assert.Equal(t, tt.userRole, role)
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       string
		expectedStatus int
		expectAbort    bool
	}{
		{
			name:           "Admin user",
			userRole:       "admin",
			expectedStatus: http.StatusOK,
			expectAbort:    false,
		},
		{
			name:           "Regular user",
			userRole:       "user",
			expectedStatus: http.StatusForbidden,
			expectAbort:    true,
		},
		{
			name:           "No role set",
			userRole:       "",
			expectedStatus: http.StatusForbidden,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)

			// Set up test endpoint
			router.GET("/test", func(c *gin.Context) {
				if tt.userRole != "" {
					c.Set("role", tt.userRole)
				}
			}, AdminMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, models.APIResponse{
					Success: true,
					Message: "OK",
				})
			})

			c.Request = httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, c.Request)

			// Assertions
			if tt.expectAbort {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}
		})
	}
}

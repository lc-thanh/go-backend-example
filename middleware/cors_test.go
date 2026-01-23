package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		origin         string
		requestHeaders string
		expectHeaders  map[string]string
	}{
		{
			name:   "GET request with origin",
			method: "GET",
			origin: "http://localhost:3000",
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "POST, OPTIONS, GET, PUT, DELETE, PATCH",
				"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:   "OPTIONS preflight request",
			method: "OPTIONS",
			origin: "http://localhost:3000",
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, OPTIONS, GET, PUT, DELETE, PATCH",
				"Access-Control-Allow-Headers": "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
			},
		},
		{
			name:   "POST request",
			method: "POST",
			origin: "http://example.com",
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)

			// Add CORS middleware and test endpoint
			router.Use(CORSMiddleware())
			router.Handle(tt.method, "/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			// Create request
			c.Request = httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				c.Request.Header.Set("Origin", tt.origin)
			}
			if tt.requestHeaders != "" {
				c.Request.Header.Set("Access-Control-Request-Headers", tt.requestHeaders)
			}

			// Execute request
			router.ServeHTTP(w, c.Request)

			// Assertions
			for header, expectedValue := range tt.expectHeaders {
				actualValue := w.Header().Get(header)
				assert.Equal(t, expectedValue, actualValue, "Header %s mismatch", header)
			}

			// OPTIONS should return 204
			if tt.method == "OPTIONS" {
				assert.Equal(t, http.StatusNoContent, w.Code)
			}
		})
	}
}

func TestCORSMiddlewareWithoutOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest("GET", "/test", nil)
	// Don't set Origin header

	router.ServeHTTP(w, c.Request)

	// Should still set CORS headers
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

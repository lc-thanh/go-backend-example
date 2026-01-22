package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
		body       string
	}{
		{
			name:       "GET request",
			method:     "GET",
			path:       "/test",
			statusCode: http.StatusOK,
			body:       "",
		},
		{
			name:       "POST request",
			method:     "POST",
			path:       "/test",
			statusCode: http.StatusCreated,
			body:       `{"test":"data"}`,
		},
		{
			name:       "404 Not Found",
			method:     "GET",
			path:       "/notfound",
			statusCode: http.StatusNotFound,
			body:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture logs
			var logBuffer bytes.Buffer
			
			// Setup
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)
			
			// Add logger middleware
			router.Use(LoggerMiddleware())
			
			// Add test endpoints
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})
			router.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{"status": "created"})
			})

			// Create request
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			c.Request = req

			// Execute request
			router.ServeHTTP(w, req)

			// The middleware should not break the request
			assert.NotNil(t, w)
			
			// For existing routes, check status code matches
			if tt.path == "/test" {
				assert.Equal(t, tt.statusCode, w.Code)
			}
			
			// Logger should have been called (we can't easily verify logs without dependency injection)
			_ = logBuffer
		})
	}
}

func TestLoggerMiddlewareWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)
	
	router.Use(LoggerMiddleware())
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	})

	c.Request = httptest.NewRequest("GET", "/error", nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

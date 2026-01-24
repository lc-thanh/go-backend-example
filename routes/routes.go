package routes

import (
	"go_backend/controllers"
	"go_backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running...",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)

	// Protected routes (authentication required)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware())

	// User profile routes
	protected.GET("/profile", controllers.GetProfile)

	// Product routes
	products := protected.Group("/products")
	products.GET("", controllers.GetProducts)                             // GET /api/v1/products
	products.GET("/:id", controllers.GetProductByID)                      // GET /api/v1/products/:id
	products.POST("", controllers.CreateProduct)                          // POST /api/v1/products
	products.PUT("/:id", controllers.UpdateProduct)                       // PUT /api/v1/products/:id
	products.DELETE("/:id", controllers.DeleteProduct)                    // DELETE /api/v1/products/:id (soft delete)
	products.POST("/:id/restore", controllers.RestoreProduct)             // POST /api/v1/products/:id/restore
	products.DELETE("/:id/permanent", controllers.PermanentDeleteProduct) // DELETE /api/v1/products/:id/permanent (hard delete)
}

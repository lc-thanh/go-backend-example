package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go_backend/cache"
	"go_backend/database"
	"go_backend/models"

	"github.com/gin-gonic/gin"
)

// GetProducts retrieves all products with pagination
func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("products:page:%d:limit:%d:category:%s", page, limit, category)

	// Try to get from cache
	var cachedResponse models.APIResponse
	if cache.RedisClient != nil {
		if err := cache.Get(cacheKey, &cachedResponse); err == nil {
			c.Header("X-Cache", "HIT")
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	offset := (page - 1) * limit

	var products []models.Product
	var total int64

	query := database.DB.Model(&models.Product{})

	// Filter by category if provided
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Count total records
	query.Count(&total)

	// Get paginated products
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch products",
			Error:   err.Error(),
		})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	response := models.APIResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data: models.PaginationResponse{
			Data:       products,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	// Cache the response for 5 minutes
	if cache.RedisClient != nil {
		if err := cache.Set(cacheKey, response, 5*time.Minute); err != nil {
			log.Printf("Failed to cache products list: %v", err)
		}
	}

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, response)
}

// GetProductByID retrieves a single product by ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	// Generate cache key
	cacheKey := fmt.Sprintf("product:%s", id)

	// Try to get from cache
	var cachedProduct models.Product
	if cache.RedisClient != nil {
		if err := cache.Get(cacheKey, &cachedProduct); err == nil {
			c.Header("X-Cache", "HIT")
			c.JSON(http.StatusOK, models.APIResponse{
				Success: true,
				Message: "Product retrieved successfully",
				Data:    cachedProduct,
			})
			return
		}
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	// Cache the product for 10 minutes
	if cache.RedisClient != nil {
		if err := cache.Set(cacheKey, product, 10*time.Minute); err != nil {
			log.Printf("Failed to cache product %s: %v", id, err)
		}
	}

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		SKU:         req.SKU,
		IsActive:    true,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create product",
			Error:   err.Error(),
		})
		return
	}

	// Invalidate products list cache
	if cache.RedisClient != nil {
		if err := cache.DeletePattern("products:*"); err != nil {
			log.Printf("Failed to invalidate products cache: %v", err)
		}
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Update fields if they are provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update product",
			Error:   err.Error(),
		})
		return
	}

	// Invalidate cache for this product and products list
	if cache.RedisClient != nil {
		if err := cache.Delete(fmt.Sprintf("product:%s", id)); err != nil {
			log.Printf("Failed to delete product cache: %v", err)
		}
		if err := cache.DeletePattern("products:*"); err != nil {
			log.Printf("Failed to invalidate products cache: %v", err)
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct deletes a product (soft delete)
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	// Soft delete
	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete product",
			Error:   err.Error(),
		})
		return
	}

	// Invalidate cache for this product and products list
	if cache.RedisClient != nil {
		if err := cache.Delete(fmt.Sprintf("product:%s", id)); err != nil {
			log.Printf("Failed to delete product cache: %v", err)
		}
		if err := cache.DeletePattern("products:*"); err != nil {
			log.Printf("Failed to invalidate products cache: %v", err)
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product deleted successfully",
		Data:    nil,
	})
}

// RestoreProduct restores a soft-deleted product
func RestoreProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	// Use Unscoped() to query soft-deleted records
	if err := database.DB.Unscoped().First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	// Check if product is already active (not deleted)
	if product.DeletedAt.Time.IsZero() {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Product is not deleted",
			Error:   "Cannot restore a product that is not deleted",
		})
		return
	}

	// Restore by setting DeletedAt to null
	if err := database.DB.Unscoped().Model(&product).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to restore product",
			Error:   err.Error(),
		})
		return
	}

	// Fetch the restored product
	database.DB.First(&product, id)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product restored successfully",
		Data:    product,
	})
}

// PermanentDeleteProduct permanently deletes a product from database
func PermanentDeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	// Use Unscoped() to query all records (including soft-deleted)
	if err := database.DB.Unscoped().First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	// Permanently delete (hard delete)
	if err := database.DB.Unscoped().Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to permanently delete product",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product permanently deleted",
		Data:    nil,
	})
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name" binding:"required,min=3,max=200"`
	Description string         `json:"description"`
	Price       float64        `gorm:"not null" json:"price" binding:"required,gt=0"`
	Stock       int            `gorm:"default:0" json:"stock" binding:"gte=0"`
	Category    string         `json:"category"`
	ImageURL    string         `json:"image_url"`
	SKU         string         `gorm:"unique" json:"sku"` // Stock Keeping Unit
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// CreateProductRequest represents data for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	SKU         string  `json:"sku"`
}

// UpdateProductRequest represents data for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=3,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	SKU         string  `json:"sku"`
	IsActive    *bool   `json:"is_active"` // Pointer to allow null value
}

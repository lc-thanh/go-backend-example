package models

import (
	"testing"
	"time"
)

func TestProductModel(t *testing.T) {
	product := Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Stock:       100,
		Category:    "electronics",
		ImageURL:    "https://example.com/image.jpg",
		SKU:         "TEST-001",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if product.Name != "Test Product" {
		t.Errorf("Name = %v, want Test Product", product.Name)
	}
	if product.Price != 99.99 {
		t.Errorf("Price = %v, want 99.99", product.Price)
	}
	if product.Stock != 100 {
		t.Errorf("Stock = %v, want 100", product.Stock)
	}
	if product.Category != "electronics" {
		t.Errorf("Category = %v, want electronics", product.Category)
	}
	if !product.IsActive {
		t.Error("IsActive should be true")
	}
}

func TestCreateProductRequest(t *testing.T) {
	req := CreateProductRequest{
		Name:        "New Product",
		Description: "New Description",
		Price:       49.99,
		Stock:       50,
		Category:    "books",
		ImageURL:    "https://example.com/book.jpg",
		SKU:         "BOOK-001",
	}

	if req.Name != "New Product" {
		t.Errorf("Name = %v, want New Product", req.Name)
	}
	if req.Price != 49.99 {
		t.Errorf("Price = %v, want 49.99", req.Price)
	}
	if req.Stock != 50 {
		t.Errorf("Stock = %v, want 50", req.Stock)
	}
}

func TestUpdateProductRequest(t *testing.T) {
	isActive := true
	req := UpdateProductRequest{
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       79.99,
		Stock:       75,
		Category:    "clothing",
		ImageURL:    "https://example.com/updated.jpg",
		SKU:         "UPD-001",
		IsActive:    &isActive,
	}

	if req.Name != "Updated Product" {
		t.Errorf("Name = %v, want Updated Product", req.Name)
	}
	if req.Price != 79.99 {
		t.Errorf("Price = %v, want 79.99", req.Price)
	}
	if req.IsActive == nil || *req.IsActive != true {
		t.Error("IsActive should be true")
	}
}

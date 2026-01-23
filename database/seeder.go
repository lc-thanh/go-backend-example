package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"go_backend/models"
)

// SeedData seeds initial data into the database
func SeedData() {
	// Check if admin already exists
	var adminCount int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)

	if adminCount > 0 {
		log.Println("ℹ️  Admin user already exists, skipping seed")
		return
	}

	// Hash the default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Create default admin user
	admin := models.User{
		Email:    "admin@example.com",
		Username: "admin",
		Password: string(hashedPassword),
		FullName: "System Administrator",
		Role:     "admin",
		IsActive: true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Fatal("Failed to seed admin user:", err)
	}

	log.Println("✅ Admin user seeded successfully")
	log.Println("   Email: admin@example.com")
	log.Println("   Password: admin123")
	log.Println("   ⚠️  Please change the password in production!")
}

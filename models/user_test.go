package models

import (
	"testing"
	"time"
)

func TestUserModel(t *testing.T) {
	user := User{
		ID:        1,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FullName:  "Test User",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", user.Email)
	}
	if user.Username != "testuser" {
		t.Errorf("Username = %v, want testuser", user.Username)
	}
	if user.Role != "user" {
		t.Errorf("Role = %v, want user", user.Role)
	}
	if !user.IsActive {
		t.Error("IsActive should be true")
	}
}

func TestUserResponse(t *testing.T) {
	response := UserResponse{
		ID:        1,
		Email:     "test@example.com",
		Username:  "testuser",
		FullName:  "Test User",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if response.Email != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", response.Email)
	}
	if response.Username != "testuser" {
		t.Errorf("Username = %v, want testuser", response.Username)
	}
}

func TestLoginRequest(t *testing.T) {
	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	if loginReq.Email != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", loginReq.Email)
	}
	if loginReq.Password != "password123" {
		t.Errorf("Password = %v, want password123", loginReq.Password)
	}
}

func TestRegisterRequest(t *testing.T) {
	registerReq := RegisterRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
		FullName: "Test User",
	}

	if registerReq.Email != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", registerReq.Email)
	}
	if registerReq.Username != "testuser" {
		t.Errorf("Username = %v, want testuser", registerReq.Username)
	}
	if registerReq.FullName != "Test User" {
		t.Errorf("FullName = %v, want Test User", registerReq.FullName)
	}
}

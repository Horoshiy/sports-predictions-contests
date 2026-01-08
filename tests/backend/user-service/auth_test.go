package tests

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/user-service/internal/models"
)

func TestJWTTokenGeneration(t *testing.T) {
	secret := []byte("test-secret")
	userID := uint(1)
	email := "test@example.com"
	expiration := time.Hour

	token, err := auth.GenerateToken(userID, email, secret, expiration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Validate the token
	claims, err := auth.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestUserPasswordHashing(t *testing.T) {
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	originalPassword := user.Password

	err := user.HashPassword()
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Password should be changed after hashing
	if user.Password == originalPassword {
		t.Error("Password was not hashed")
	}

	// Should be able to verify the original password
	if !user.CheckPassword(originalPassword) {
		t.Error("Password verification failed")
	}

	// Should not verify wrong password
	if user.CheckPassword("wrongpassword") {
		t.Error("Password verification should fail for wrong password")
	}
}

package models

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `gorm:"not null" json:"name"`
	gorm.Model
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return errors.New("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword verifies the provided password against the stored hash
func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// ValidateEmail checks if the email format is valid
func (u *User) ValidateEmail() error {
	if len(u.Email) == 0 {
		return errors.New("email cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateName checks if the name is valid
func (u *User) ValidateName() error {
	if len(strings.TrimSpace(u.Name)) == 0 {
		return errors.New("name cannot be empty")
	}

	if len(u.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}

	return nil
}

// ValidatePassword checks if the password meets requirements
func (u *User) ValidatePassword() error {
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if len(u.Password) > 128 {
		return errors.New("password cannot exceed 128 characters")
	}

	return nil
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Validate fields
	if err := u.ValidateEmail(); err != nil {
		return err
	}

	if err := u.ValidateName(); err != nil {
		return err
	}

	if err := u.ValidatePassword(); err != nil {
		return err
	}

	// Hash password
	return u.HashPassword()
}

package service

import (
	"errors"
	"time"

	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/user-service/internal/models"
	"github.com/sports-prediction-contests/user-service/internal/repository"
)

// AuthServiceInterface defines the contract for authentication service
type AuthServiceInterface interface {
	Register(email, password, name string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
}

// AuthService implements AuthServiceInterface
type AuthService struct {
	userRepo       repository.UserRepositoryInterface
	jwtSecret      []byte
	jwtExpiration  time.Duration
}

// NewAuthService creates a new authentication service instance
func NewAuthService(userRepo repository.UserRepositoryInterface, jwtSecret []byte, jwtExpiration time.Duration) AuthServiceInterface {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// Register creates a new user account
func (s *AuthService) Register(email, password, name string) (*models.User, string, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, "", errors.New("user with this email already exists")
	}

	// Create new user
	user := &models.User{
		Email:    email,
		Password: password,
		Name:     name,
	}

	// Create user in database (validation and password hashing happens in BeforeCreate hook)
	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// Validate input
	if email == "" {
		return nil, "", errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, "", errors.New("password cannot be empty")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

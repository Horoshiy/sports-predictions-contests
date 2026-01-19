package repository

import (
	"errors"

	"github.com/sports-prediction-contests/user-service/internal/models"
	"gorm.io/gorm"
)

// UserRepositoryInterface defines the contract for user repository
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	CreateProfile(profile *models.Profile) error
	UpdateProfile(profile *models.Profile) error
	CreatePreferences(preferences *models.UserPreferences) error
	UpdatePreferences(preferences *models.UserPreferences) error
}

// UserRepository implements UserRepositoryInterface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	result := r.db.Create(user)
	if result.Error != nil {
		// Check for duplicate email error
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("user with this email already exists")
		}
		return result.Error
	}

	return nil
}

// GetByEmail retrieves a user by email address
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByID retrieves a user by ID with profile and preferences
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	result := r.db.Preload("Profile").Preload("Preferences").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.ID == 0 {
		return errors.New("user ID cannot be zero")
	}

	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}

	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// CreateProfile creates a new profile for a user
func (r *UserRepository) CreateProfile(profile *models.Profile) error {
	if profile == nil {
		return errors.New("profile cannot be nil")
	}

	result := r.db.Create(profile)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateProfile updates an existing profile or creates one if it doesn't exist
func (r *UserRepository) UpdateProfile(profile *models.Profile) error {
	if profile == nil {
		return errors.New("profile cannot be nil")
	}

	if profile.UserID == 0 {
		return errors.New("user ID cannot be zero")
	}

	// Try to find existing profile
	var existingProfile models.Profile
	result := r.db.Where("user_id = ?", profile.UserID).First(&existingProfile)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Profile doesn't exist, create it
		return r.CreateProfile(profile)
	} else if result.Error != nil {
		return result.Error
	}

	// Profile exists, use selective update to preserve unmodified fields
	result = r.db.Model(&existingProfile).Updates(profile)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CreatePreferences creates new preferences for a user
func (r *UserRepository) CreatePreferences(preferences *models.UserPreferences) error {
	if preferences == nil {
		return errors.New("preferences cannot be nil")
	}

	result := r.db.Create(preferences)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdatePreferences updates existing preferences or creates them if they don't exist
func (r *UserRepository) UpdatePreferences(preferences *models.UserPreferences) error {
	if preferences == nil {
		return errors.New("preferences cannot be nil")
	}

	if preferences.UserID == 0 {
		return errors.New("user ID cannot be zero")
	}

	// Try to find existing preferences
	var existingPreferences models.UserPreferences
	result := r.db.Where("user_id = ?", preferences.UserID).First(&existingPreferences)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Preferences don't exist, create them
		return r.CreatePreferences(preferences)
	} else if result.Error != nil {
		return result.Error
	}

	// Preferences exist, use selective update to preserve unmodified fields
	result = r.db.Model(&existingPreferences).Updates(preferences)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

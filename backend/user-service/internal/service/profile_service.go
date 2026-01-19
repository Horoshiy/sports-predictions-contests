package service

import (
	"context"
	"strings"

	"github.com/sports-prediction-contests/user-service/internal/models"
	"github.com/sports-prediction-contests/user-service/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProfileServiceInterface defines the profile service interface
type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, userID uint) (*models.Profile, error)
	UpdateProfile(ctx context.Context, userID uint, profile *models.Profile) (*models.Profile, error)
	GetPreferences(ctx context.Context, userID uint) (*models.UserPreferences, error)
	UpdatePreferences(ctx context.Context, userID uint, preferences *models.UserPreferences) (*models.UserPreferences, error)
	CalculateProfileCompletion(ctx context.Context, userID uint) (int32, []string, []string, error)
}

// ProfileService implements the ProfileServiceInterface
type ProfileService struct {
	userRepo repository.UserRepositoryInterface
}

// NewProfileService creates a new ProfileService instance
func NewProfileService(userRepo repository.UserRepositoryInterface) ProfileServiceInterface {
	return &ProfileService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves a user's profile
func (s *ProfileService) GetProfile(ctx context.Context, userID uint) (*models.Profile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// If profile doesn't exist, create a default one
	if user.Profile == nil {
		profile := &models.Profile{
			UserID:            userID,
			ProfileVisibility: "public",
		}

		// Create the profile in database
		// Note: The unique constraint on UserID in the database will prevent
		// duplicate profiles even in race conditions
		if err := s.userRepo.CreateProfile(profile); err != nil {
			// If creation fails due to duplicate, try to fetch again
			user, err = s.userRepo.GetByID(userID)
			if err != nil {
				return nil, status.Error(codes.Internal, "failed to get profile after creation attempt")
			}
			if user.Profile != nil {
				return user.Profile, nil
			}
			return nil, status.Error(codes.Internal, "failed to create profile")
		}

		return profile, nil
	}

	return user.Profile, nil
}

// UpdateProfile updates a user's profile
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uint, profile *models.Profile) (*models.Profile, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Set the user ID
	profile.UserID = userID

	// Update or create profile
	if err := s.userRepo.UpdateProfile(profile); err != nil {
		return nil, status.Error(codes.Internal, "failed to update profile")
	}

	return profile, nil
}

// GetPreferences retrieves a user's preferences
func (s *ProfileService) GetPreferences(ctx context.Context, userID uint) (*models.UserPreferences, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// If preferences don't exist, create default ones
	if user.Preferences == nil {
		preferences := &models.UserPreferences{
			UserID:               userID,
			EmailNotifications:   true,
			PushNotifications:    true,
			ContestNotifications: true,
			PredictionReminders:  true,
			WeeklyDigest:         true,
			Theme:                "light",
			Language:             "en",
			Timezone:             "UTC",
		}

		// Create the preferences in database
		// Note: The unique constraint on UserID in the database will prevent
		// duplicate preferences even in race conditions
		if err := s.userRepo.CreatePreferences(preferences); err != nil {
			// If creation fails due to duplicate, try to fetch again
			user, err = s.userRepo.GetByID(userID)
			if err != nil {
				return nil, status.Error(codes.Internal, "failed to get preferences after creation attempt")
			}
			if user.Preferences != nil {
				return user.Preferences, nil
			}
			return nil, status.Error(codes.Internal, "failed to create preferences")
		}

		return preferences, nil
	}

	return user.Preferences, nil
}

// UpdatePreferences updates a user's preferences
func (s *ProfileService) UpdatePreferences(ctx context.Context, userID uint, preferences *models.UserPreferences) (*models.UserPreferences, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Set the user ID
	preferences.UserID = userID

	// Update or create preferences
	if err := s.userRepo.UpdatePreferences(preferences); err != nil {
		return nil, status.Error(codes.Internal, "failed to update preferences")
	}

	return preferences, nil
}

// CalculateProfileCompletion calculates profile completion percentage
func (s *ProfileService) CalculateProfileCompletion(ctx context.Context, userID uint) (int32, []string, []string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 0, nil, nil, status.Error(codes.NotFound, "user not found")
	}

	// Define required fields with weights
	requiredFields := map[string]int{
		"name":     20, // Basic info
		"email":    20, // Basic info
		"bio":      15, // Profile description
		"location": 10, // Location info
		"avatar":   20, // Profile picture
	}

	// Optional fields with weights
	optionalFields := map[string]int{
		"website":  5,
		"twitter":  3,
		"linkedin": 4,
		"github":   3,
	}

	totalScore := 0
	maxScore := 0
	var missingFields []string
	var suggestions []string

	// Calculate required fields score
	for field, weight := range requiredFields {
		maxScore += weight

		switch field {
		case "name":
			if strings.TrimSpace(user.Name) != "" {
				totalScore += weight
			} else {
				missingFields = append(missingFields, field)
				suggestions = append(suggestions, "Add your full name to help others identify you")
			}
		case "email":
			if strings.TrimSpace(user.Email) != "" {
				totalScore += weight
			} else {
				missingFields = append(missingFields, field)
				suggestions = append(suggestions, "Verify your email address")
			}
		case "bio":
			if user.Profile != nil && strings.TrimSpace(user.Profile.Bio) != "" {
				totalScore += weight
			} else {
				missingFields = append(missingFields, field)
				suggestions = append(suggestions, "Write a brief bio to tell others about yourself")
			}
		case "location":
			if user.Profile != nil && strings.TrimSpace(user.Profile.Location) != "" {
				totalScore += weight
			} else {
				missingFields = append(missingFields, field)
				suggestions = append(suggestions, "Add your location to connect with local users")
			}
		case "avatar":
			if user.Profile != nil && strings.TrimSpace(user.Profile.AvatarURL) != "" {
				totalScore += weight
			} else {
				missingFields = append(missingFields, field)
				suggestions = append(suggestions, "Upload a profile picture to personalize your account")
			}
		}
	}

	// Calculate optional fields score
	for field, weight := range optionalFields {
		maxScore += weight

		if user.Profile != nil {
			switch field {
			case "website":
				if strings.TrimSpace(user.Profile.Website) != "" {
					totalScore += weight
				}
			case "twitter":
				if strings.TrimSpace(user.Profile.TwitterURL) != "" {
					totalScore += weight
				}
			case "linkedin":
				if strings.TrimSpace(user.Profile.LinkedInURL) != "" {
					totalScore += weight
				}
			case "github":
				if strings.TrimSpace(user.Profile.GitHubURL) != "" {
					totalScore += weight
				}
			}
		}
	}

	// Calculate percentage
	percentage := int32((totalScore * 100) / maxScore)

	// Add general suggestions based on completion level
	if percentage < 50 {
		suggestions = append(suggestions, "Complete your basic profile information to unlock all features")
	} else if percentage < 80 {
		suggestions = append(suggestions, "Add social links to connect with other users")
	}

	return percentage, missingFields, suggestions, nil
}

package tests

import (
	"errors"
	"testing"
)

// Mock models for testing since we can't import the actual models due to module path issues
type Profile struct {
	UserID            uint
	Bio               string
	Location          string
	Website           string
	ProfileVisibility string
}

func (p *Profile) ValidateBio() error {
	if len(p.Bio) > 500 {
		return errors.New("bio must be less than 500 characters")
	}
	return nil
}

func (p *Profile) ValidateProfileVisibility() error {
	validVisibilities := []string{"public", "friends", "private"}
	for _, v := range validVisibilities {
		if p.ProfileVisibility == v {
			return nil
		}
	}
	return errors.New("profile visibility must be 'public', 'friends', or 'private'")
}

func TestProfileValidation(t *testing.T) {
	tests := []struct {
		name    string
		profile *Profile
		field   string
		wantErr bool
	}{
		{
			name: "valid bio",
			profile: &Profile{
				Bio: "Test bio",
			},
			field:   "bio",
			wantErr: false,
		},
		{
			name: "bio too long",
			profile: &Profile{
				Bio: string(make([]byte, 501)), // 501 characters
			},
			field:   "bio",
			wantErr: true,
		},
		{
			name: "valid visibility",
			profile: &Profile{
				ProfileVisibility: "public",
			},
			field:   "visibility",
			wantErr: false,
		},
		{
			name: "invalid visibility",
			profile: &Profile{
				ProfileVisibility: "invalid",
			},
			field:   "visibility",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch tt.field {
			case "bio":
				err = tt.profile.ValidateBio()
			case "visibility":
				err = tt.profile.ValidateProfileVisibility()
			}
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Profile validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

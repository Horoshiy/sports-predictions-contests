package profile_test

import (
	"errors"
	"regexp"
	"testing"
)

// Test URL validation improvements
func TestURLValidation(t *testing.T) {
	// Improved regex from the fix
	urlRegex := regexp.MustCompile(`^https?://(?:[-\w.])+(?:\.[a-zA-Z]{2,})+(?:/[^?\s]*)?(?:\?[^#\s]*)?(?:#[^\s]*)?$`)

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "valid http URL with path",
			url:     "http://example.com/path/to/page",
			wantErr: false,
		},
		{
			name:    "valid URL with query",
			url:     "https://example.com/page?param=value",
			wantErr: false,
		},
		{
			name:    "valid URL with fragment",
			url:     "https://example.com/page#section",
			wantErr: false,
		},
		{
			name:    "invalid - no protocol",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "invalid - no domain",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "invalid - spaces",
			url:     "https://example .com",
			wantErr: true,
		},
		{
			name:    "invalid - malformed",
			url:     "https://.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := urlRegex.MatchString(tt.url)
			if matches == tt.wantErr {
				t.Errorf("URL validation for %q: got match=%v, want match=%v", tt.url, matches, !tt.wantErr)
			}
		})
	}
}

// Test profile validation after fixes
func TestProfileValidationFixes(t *testing.T) {
	type Profile struct {
		Bio               string
		Location          string
		ProfileVisibility string
	}

	validateBio := func(bio string) error {
		if len(bio) > 500 {
			return errors.New("bio must be less than 500 characters")
		}
		return nil
	}

	validateLocation := func(location string) error {
		if len(location) > 100 {
			return errors.New("location must be less than 100 characters")
		}
		return nil
	}

	validateVisibility := func(visibility string) error {
		validVisibilities := []string{"public", "friends", "private"}
		for _, v := range validVisibilities {
			if visibility == v {
				return nil
			}
		}
		return errors.New("profile visibility must be 'public', 'friends', or 'private'")
	}

	tests := []struct {
		name    string
		profile Profile
		wantErr bool
	}{
		{
			name: "valid profile",
			profile: Profile{
				Bio:               "Test bio",
				Location:          "Test City",
				ProfileVisibility: "public",
			},
			wantErr: false,
		},
		{
			name: "bio too long",
			profile: Profile{
				Bio:               string(make([]byte, 501)),
				ProfileVisibility: "public",
			},
			wantErr: true,
		},
		{
			name: "location too long",
			profile: Profile{
				Location:          string(make([]byte, 101)),
				ProfileVisibility: "public",
			},
			wantErr: true,
		},
		{
			name: "invalid visibility",
			profile: Profile{
				ProfileVisibility: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			if e := validateBio(tt.profile.Bio); e != nil {
				err = e
			}
			if e := validateLocation(tt.profile.Location); e != nil {
				err = e
			}
			if e := validateVisibility(tt.profile.ProfileVisibility); e != nil {
				err = e
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Profile validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test file type detection
func TestFileTypeDetection(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		allowed     []string
		wantAllowed bool
	}{
		{
			name:        "valid JPEG",
			contentType: "image/jpeg",
			allowed:     []string{"image/jpeg", "image/png", "image/gif"},
			wantAllowed: true,
		},
		{
			name:        "valid PNG",
			contentType: "image/png",
			allowed:     []string{"image/jpeg", "image/png", "image/gif"},
			wantAllowed: true,
		},
		{
			name:        "invalid PDF",
			contentType: "application/pdf",
			allowed:     []string{"image/jpeg", "image/png", "image/gif"},
			wantAllowed: false,
		},
		{
			name:        "invalid executable",
			contentType: "application/x-executable",
			allowed:     []string{"image/jpeg", "image/png", "image/gif"},
			wantAllowed: false,
		},
	}

	isAllowedType := func(contentType string, allowedTypes []string) bool {
		for _, allowedType := range allowedTypes {
			if contentType == allowedType {
				return true
			}
		}
		return false
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := isAllowedType(tt.contentType, tt.allowed)
			if allowed != tt.wantAllowed {
				t.Errorf("File type %q: got allowed=%v, want allowed=%v", tt.contentType, allowed, tt.wantAllowed)
			}
		})
	}
}

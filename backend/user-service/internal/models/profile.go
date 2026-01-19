package models

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// Profile represents extended user profile information
type Profile struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	UserID            uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Bio               string `gorm:"type:text" json:"bio"`
	AvatarURL         string `gorm:"size:500" json:"avatar_url"`
	Location          string `gorm:"size:100" json:"location"`
	Website           string `gorm:"size:200" json:"website"`
	TwitterURL        string `gorm:"size:200" json:"twitter_url"`
	LinkedInURL       string `gorm:"size:200" json:"linkedin_url"`
	GitHubURL         string `gorm:"size:200" json:"github_url"`
	ProfileVisibility string `gorm:"size:20;default:'public'" json:"profile_visibility"`
	gorm.Model
}

// ValidateBio checks if the bio is valid
func (p *Profile) ValidateBio() error {
	if len(p.Bio) > 500 {
		return errors.New("bio must be less than 500 characters")
	}
	return nil
}

// ValidateLocation checks if the location is valid
func (p *Profile) ValidateLocation() error {
	if len(p.Location) > 100 {
		return errors.New("location must be less than 100 characters")
	}
	return nil
}

// ValidateWebsite checks if the website URL is valid
func (p *Profile) ValidateWebsite() error {
	if p.Website == "" {
		return nil
	}

	if len(p.Website) > 200 {
		return errors.New("website URL must be less than 200 characters")
	}

	// More robust URL validation
	urlRegex := regexp.MustCompile(`^https?://(?:[-\w.])+(?:\.[a-zA-Z]{2,})+(?:/[^?\s]*)?(?:\?[^#\s]*)?(?:#[^\s]*)?$`)
	if !urlRegex.MatchString(p.Website) {
		return errors.New("invalid website URL format")
	}

	return nil
}

// ValidateTwitterURL checks if the Twitter URL is valid
func (p *Profile) ValidateTwitterURL() error {
	if p.TwitterURL == "" {
		return nil
	}

	if len(p.TwitterURL) > 200 {
		return errors.New("Twitter URL must be less than 200 characters")
	}

	twitterRegex := regexp.MustCompile(`^https?://(?:www\.)?twitter\.com/[a-zA-Z0-9_]+/?$`)
	if !twitterRegex.MatchString(p.TwitterURL) {
		return errors.New("invalid Twitter URL format")
	}

	return nil
}

// ValidateLinkedInURL checks if the LinkedIn URL is valid
func (p *Profile) ValidateLinkedInURL() error {
	if p.LinkedInURL == "" {
		return nil
	}

	if len(p.LinkedInURL) > 200 {
		return errors.New("LinkedIn URL must be less than 200 characters")
	}

	linkedinRegex := regexp.MustCompile(`^https?://(?:www\.)?linkedin\.com/in/[a-zA-Z0-9-]+/?$`)
	if !linkedinRegex.MatchString(p.LinkedInURL) {
		return errors.New("invalid LinkedIn URL format")
	}

	return nil
}

// ValidateGitHubURL checks if the GitHub URL is valid
func (p *Profile) ValidateGitHubURL() error {
	if p.GitHubURL == "" {
		return nil
	}

	if len(p.GitHubURL) > 200 {
		return errors.New("GitHub URL must be less than 200 characters")
	}

	githubRegex := regexp.MustCompile(`^https?://(?:www\.)?github\.com/[a-zA-Z0-9-]+/?$`)
	if !githubRegex.MatchString(p.GitHubURL) {
		return errors.New("invalid GitHub URL format")
	}

	return nil
}

// ValidateProfileVisibility checks if the profile visibility is valid
func (p *Profile) ValidateProfileVisibility() error {
	validVisibilities := []string{"public", "friends", "private"}
	for _, v := range validVisibilities {
		if p.ProfileVisibility == v {
			return nil
		}
	}
	return errors.New("profile visibility must be 'public', 'friends', or 'private'")
}

// ValidateAll validates all profile fields without setting defaults
func (p *Profile) ValidateAll() error {
	if err := p.ValidateBio(); err != nil {
		return err
	}

	if err := p.ValidateLocation(); err != nil {
		return err
	}

	if err := p.ValidateWebsite(); err != nil {
		return err
	}

	if err := p.ValidateTwitterURL(); err != nil {
		return err
	}

	if err := p.ValidateLinkedInURL(); err != nil {
		return err
	}

	if err := p.ValidateGitHubURL(); err != nil {
		return err
	}

	if err := p.ValidateProfileVisibility(); err != nil {
		return err
	}

	return nil
}

// BeforeCreate is a GORM hook that runs before creating a profile
func (p *Profile) BeforeCreate(tx *gorm.DB) error {
	// Validate all fields
	if err := p.ValidateAll(); err != nil {
		return err
	}

	// Set default visibility if empty (only on create)
	if strings.TrimSpace(p.ProfileVisibility) == "" {
		p.ProfileVisibility = "public"
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a profile
func (p *Profile) BeforeUpdate(tx *gorm.DB) error {
	// Only validate, don't set defaults on update
	return p.ValidateAll()
}

package models

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// UserPreferences represents user preferences and settings
type UserPreferences struct {
	ID                   uint            `gorm:"primaryKey" json:"id"`
	UserID               uint            `gorm:"uniqueIndex;not null" json:"user_id"`
	EmailNotifications   bool            `gorm:"default:true" json:"email_notifications"`
	PushNotifications    bool            `gorm:"default:true" json:"push_notifications"`
	ContestNotifications bool            `gorm:"default:true" json:"contest_notifications"`
	PredictionReminders  bool            `gorm:"default:true" json:"prediction_reminders"`
	WeeklyDigest         bool            `gorm:"default:true" json:"weekly_digest"`
	Theme                string          `gorm:"size:20;default:'light'" json:"theme"`
	Language             string          `gorm:"size:10;default:'en'" json:"language"`
	Timezone             string          `gorm:"size:50;default:'UTC'" json:"timezone"`
	CustomSettings       json.RawMessage `gorm:"type:jsonb" json:"custom_settings"`
	gorm.Model
}

// ValidateTheme checks if the theme is valid
func (up *UserPreferences) ValidateTheme() error {
	validThemes := []string{"light", "dark", "auto"}
	for _, theme := range validThemes {
		if up.Theme == theme {
			return nil
		}
	}
	return errors.New("theme must be 'light', 'dark', or 'auto'")
}

// ValidateLanguage checks if the language code is valid
func (up *UserPreferences) ValidateLanguage() error {
	validLanguages := []string{"en", "ru", "es", "fr", "de"}
	for _, lang := range validLanguages {
		if up.Language == lang {
			return nil
		}
	}
	return errors.New("language must be a valid language code (en, ru, es, fr, de)")
}

// ValidateTimezone checks if the timezone is valid
func (up *UserPreferences) ValidateTimezone() error {
	if len(up.Timezone) == 0 {
		return errors.New("timezone cannot be empty")
	}

	if len(up.Timezone) > 50 {
		return errors.New("timezone must be less than 50 characters")
	}

	return nil
}

// ValidateCustomSettings checks if custom settings JSON is valid
func (up *UserPreferences) ValidateCustomSettings() error {
	if len(up.CustomSettings) == 0 {
		return nil
	}

	var temp interface{}
	if err := json.Unmarshal(up.CustomSettings, &temp); err != nil {
		return errors.New("custom settings must be valid JSON")
	}

	return nil
}

// BeforeCreate is a GORM hook that runs before creating user preferences
func (up *UserPreferences) BeforeCreate(tx *gorm.DB) error {
	// Set defaults if empty
	if up.Theme == "" {
		up.Theme = "light"
	}

	if up.Language == "" {
		up.Language = "en"
	}

	if up.Timezone == "" {
		up.Timezone = "UTC"
	}

	// Validate fields
	if err := up.ValidateTheme(); err != nil {
		return err
	}

	if err := up.ValidateLanguage(); err != nil {
		return err
	}

	if err := up.ValidateTimezone(); err != nil {
		return err
	}

	if err := up.ValidateCustomSettings(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating user preferences
func (up *UserPreferences) BeforeUpdate(tx *gorm.DB) error {
	return up.BeforeCreate(tx)
}

// GetCustomSetting retrieves a custom setting by key
func (up *UserPreferences) GetCustomSetting(key string) (interface{}, error) {
	if len(up.CustomSettings) == 0 {
		return nil, nil
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(up.CustomSettings, &settings); err != nil {
		return nil, err
	}

	return settings[key], nil
}

// SetCustomSetting sets a custom setting by key
func (up *UserPreferences) SetCustomSetting(key string, value interface{}) error {
	var settings map[string]interface{}

	if len(up.CustomSettings) > 0 {
		if err := json.Unmarshal(up.CustomSettings, &settings); err != nil {
			return err
		}
	} else {
		settings = make(map[string]interface{})
	}

	settings[key] = value

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	up.CustomSettings = data
	return nil
}

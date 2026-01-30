package seeder

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// DataFactory provides methods to generate fake data for all entities
type DataFactory struct {
	faker *gofakeit.Faker
	db    *gorm.DB
}

// NewDataFactory creates a new data factory instance
func NewDataFactory(db *gorm.DB, seed int64) *DataFactory {
	faker := gofakeit.New(seed)
	return &DataFactory{
		faker: faker,
		db:    db,
	}
}

// GenerateUsers creates fake users with profiles and preferences
func (f *DataFactory) GenerateUsers(count int) ([]*User, []*Profile, []*UserPreferences, error) {
	log.Printf("Generating %d users with profiles and preferences...", count)

	users := make([]*User, count)
	profiles := make([]*Profile, count)
	preferences := make([]*UserPreferences, count)

	for i := 0; i < count; i++ {
		// Generate secure password for development/testing
		// NOTE: These are development-only passwords, not suitable for production use
		password, err := generateSecurePassword(12)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to generate secure password: %w", err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to hash password: %w", err)
		}

		// Create user (let database auto-assign ID)
		user := &User{
			Email:    f.faker.Email(),
			Password: string(hashedPassword),
			Name:     f.faker.Name(),
		}
		users[i] = user

		// Create profile (will be linked after user is saved)
		profile := &Profile{
			Bio:               f.faker.Sentence(f.faker.Number(1, 4)),
			AvatarURL:         f.faker.ImageURL(200, 200),
			Location:          f.faker.City() + ", " + f.faker.Country(),
			Website:           f.faker.URL(),
			TwitterURL:        "https://twitter.com/" + strings.ToLower(f.faker.Username()),
			LinkedInURL:       "https://linkedin.com/in/" + strings.ToLower(f.faker.Username()),
			GitHubURL:         "https://github.com/" + strings.ToLower(f.faker.Username()),
			ProfileVisibility: f.faker.RandomString([]string{"public", "private", "friends"}),
		}
		profiles[i] = profile

		// Create preferences (will be linked after user is saved)
		preference := &UserPreferences{
			Language:              f.faker.RandomString([]string{"en", "ru", "es", "fr", "de"}),
			Timezone:              f.faker.TimeZone(),
			EmailNotifications:    f.faker.Bool(),
			PushNotifications:     f.faker.Bool(),
			TelegramNotifications: f.faker.Bool(),
			Theme:                 f.faker.RandomString([]string{"light", "dark", "auto"}),
		}
		preferences[i] = preference
	}

	log.Printf("Successfully generated %d users with profiles and preferences", count)
	return users, profiles, preferences, nil
}

// GenerateContests creates fake contests
func (f *DataFactory) GenerateContests(count int, userIDs []uint, sportTypes []string) ([]*Contest, error) {
	log.Printf("Generating %d contests...", count)

	contests := make([]*Contest, count)
	statuses := []string{"draft", "active", "completed"}

	for i := 0; i < count; i++ {
		startDate := f.faker.DateRange(time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 2, 0))
		endDate := startDate.Add(time.Duration(f.faker.Number(7, 30)) * 24 * time.Hour)

		contest := &Contest{
			ID:                  uint(i + 1),
			Title:               f.generateContestTitle(),
			Description:         f.faker.Sentence(f.faker.Number(3, 8)),
			SportType:           f.faker.RandomString(sportTypes),
			Rules:               f.generateContestRules(),
			Status:              f.faker.RandomString(statuses),
			StartDate:           startDate,
			EndDate:             endDate,
			MaxParticipants:     uint(f.faker.Number(0, 1000)),
			CurrentParticipants: uint(f.faker.Number(0, 500)),
			CreatorID:           f.faker.RandomUint(userIDs),
		}
		
		// Add prediction schema to all contests
		schemaJSON := f.GenerateDefaultPredictionSchema()
		contest.PredictionSchema = []byte(schemaJSON)
		
		contests[i] = contest
	}

	log.Printf("Successfully generated %d contests", count)
	return contests, nil
}

// GenerateDefaultPredictionSchema returns the default exact score prediction schema as JSON string
func (f *DataFactory) GenerateDefaultPredictionSchema() string {
	// Return hardcoded schema for exact score predictions
	return `{"type":"exact_score","options":["1-0","0-1","2-0","0-2","2-1","1-2","3-0","0-3","3-1","1-3","3-2","2-3","0-0","1-1","2-2","3-3"],"allow_custom":true}`
}

// GeneratePredictions creates fake predictions
func (f *DataFactory) GeneratePredictions(count int, userIDs []uint, contestIDs []uint, matchIDs []uint) ([]*Prediction, error) {
	log.Printf("Generating %d predictions...", count)

	predictions := make([]*Prediction, 0, count)
	predictionTypes := []string{"match_result", "exact_score", "over_under", "both_teams_score"}
	
	// Track unique combinations to avoid duplicates
	seen := make(map[string]bool)

	for i := 0; i < count*2 && len(predictions) < count; i++ { // Try up to 2x to handle collisions
		submittedAt := f.faker.DateRange(time.Now().AddDate(0, -2, 0), time.Now())
		predType := f.faker.RandomString(predictionTypes)
		
		userID := f.faker.RandomUint(userIDs)
		contestID := f.faker.RandomUint(contestIDs)
		matchID := f.faker.RandomUint(matchIDs)
		
		// Create unique key for this combination
		key := fmt.Sprintf("%d-%d-%d", userID, contestID, matchID)
		if seen[key] {
			continue // Skip duplicate
		}
		seen[key] = true

		prediction := &Prediction{
			ID:             uint(len(predictions) + 1),
			UserID:         userID,
			ContestID:      contestID,
			EventID:        matchID, // EventID is the same as MatchID
			MatchID:        matchID,
			PredictionType: predType,
			PredictionData: f.generatePredictionData(predType),
			SubmittedAt:    submittedAt,
			IsCorrect:      f.generateIsCorrect(),
			Points:         f.faker.Float64Range(0, 10),
		}
		predictions = append(predictions, prediction)
	}

	log.Printf("Successfully generated %d predictions", len(predictions))
	return predictions, nil
}

// Helper methods

func (f *DataFactory) generateContestTitle() string {
	prefixes := []string{"Ultimate", "Championship", "Premier", "Elite", "Grand", "Master"}
	sports := []string{"Football", "Basketball", "Soccer", "Tennis", "Baseball"}
	suffixes := []string{"Challenge", "Tournament", "Cup", "League", "Championship"}

	return fmt.Sprintf("%s %s %s",
		f.faker.RandomString(prefixes),
		f.faker.RandomString(sports),
		f.faker.RandomString(suffixes))
}

func (f *DataFactory) generateContestRules() string {
	rules := map[string]interface{}{
		"scoring_system":   f.faker.RandomString([]string{"standard", "weighted", "progressive"}),
		"points_correct":   f.faker.Number(1, 5),
		"points_incorrect": 0,
		"bonus_multiplier": f.faker.Float64Range(1.0, 2.0),
		"max_predictions":  f.faker.Number(10, 100),
		"allow_late_entry": f.faker.Bool(),
	}

	// Convert to JSON string (simplified)
	return fmt.Sprintf(`{"scoring_system":"%s","points_correct":%d,"points_incorrect":0,"bonus_multiplier":%.1f,"max_predictions":%d,"allow_late_entry":%t}`,
		rules["scoring_system"], rules["points_correct"], rules["bonus_multiplier"], rules["max_predictions"], rules["allow_late_entry"])
}

func (f *DataFactory) generatePredictionData(predType string) string {
	switch predType {
	case "match_result":
		return fmt.Sprintf(`{"result":"%s"}`, f.faker.RandomString([]string{"home", "away", "draw"}))
	case "exact_score":
		return fmt.Sprintf(`{"home_score":%d,"away_score":%d}`, f.faker.Number(0, 5), f.faker.Number(0, 5))
	case "over_under":
		return fmt.Sprintf(`{"line":%.1f,"prediction":"%s"}`, f.faker.Float64Range(1.5, 4.5), f.faker.RandomString([]string{"over", "under"}))
	case "both_teams_score":
		return fmt.Sprintf(`{"prediction":%t}`, f.faker.Bool())
	default:
		return `{"prediction":"unknown"}`
	}
}

func (f *DataFactory) generateIsCorrect() *bool {
	if f.faker.Number(1, 10) <= 3 { // 30% chance of nil (not yet determined)
		return nil
	}
	result := f.faker.Bool()
	return &result
}

// slugify creates a URL-friendly slug from a string
func (f *DataFactory) slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, "'", "")
	return text
}

// generateSecurePassword creates a cryptographically secure password for development use
// Uses a conservative character set to avoid shell escaping and URL encoding issues
func generateSecurePassword(length int) (string, error) {
	// Conservative charset avoiding problematic characters for development environments
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	password := make([]byte, length)

	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

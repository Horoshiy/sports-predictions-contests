package seeder

import (
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Coordinator manages the seeding process across all services
type Coordinator struct {
	db      *gorm.DB
	config  *Config
	factory *DataFactory
	sports  *SportsDataGenerator
}

// NewCoordinator creates a new seeding coordinator
func NewCoordinator(config *Config) (*Coordinator, error) {
	// Connect to database
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Reduce log noise during seeding
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Validate database connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection validation failed: %w", err)
	}

	// Test basic database functionality
	if err := db.Exec("SELECT 1").Error; err != nil {
		return nil, fmt.Errorf("database functionality test failed: %w", err)
	}

	// Auto-migrate all tables
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(
		&User{},
		&Profile{},
		&UserPreferences{},
		&Contest{},
		&Challenge{},
		&ChallengeParticipant{},
		&Sport{},
		&League{},
		&Team{},
		&Match{},
		&Prediction{},
		&Score{},
		&Leaderboard{},
		&UserStreak{},
		&PropType{},
		&UserTeam{},
		&UserTeamMember{},
		&Notification{},
		&NotificationPreference{},
	); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Database migrations completed successfully")

	// Create factory and sports generator
	factory := NewDataFactory(db, config.Seed)
	sports := NewSportsDataGenerator(gofakeit.New(config.Seed))

	return &Coordinator{
		db:      db,
		config:  config,
		factory: factory,
		sports:  sports,
	}, nil
}

// SeedAll performs complete data seeding in the correct dependency order
func (c *Coordinator) SeedAll() (err error) {
	log.Println("Starting comprehensive data seeding...")

	// Clean existing data first
	if err := c.cleanExistingData(); err != nil {
		return fmt.Errorf("failed to clean existing data: %w", err)
	}

	counts := c.config.GetDataCounts()

	// Start transaction for data consistency
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("seeding failed with panic: %v", r)
			log.Printf("Seeding failed with panic: %v", r)
		}
	}()

	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Phase 1: Core entities (no dependencies)
	log.Println("Phase 1: Seeding core entities...")

	// 1. Users (must be first - referenced by everything)
	users, _, _, err := c.seedUsers(tx, counts.Users)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// 2. Sports (referenced by leagues, teams, contests)
	sports, err := c.seedSports(tx, counts.Sports)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed sports: %w", err)
	}

	// Phase 2: Dependent entities
	log.Println("Phase 2: Seeding dependent entities...")

	// 3. Leagues (depends on sports)
	leagues, err := c.seedLeagues(tx, counts.Leagues, sports)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed leagues: %w", err)
	}

	// 4. Teams (depends on sports)
	teams, err := c.seedTeams(tx, counts.Teams, sports)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed teams: %w", err)
	}

	// 5. Matches (depends on leagues and teams)
	matches, err := c.seedMatches(tx, counts.Matches, leagues, teams)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed matches: %w", err)
	}

	// 5.5. Events (create from matches for prediction-service)
	events, err := c.seedEvents(tx, matches, teams, sports)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed events: %w", err)
	}

	// 6. Contests (depends on users and sports)
	contests, err := c.seedContests(tx, counts.Contests, users, sports)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed contests: %w", err)
	}

	// 6.5. Challenges (depends on users and events)
	challenges, err := c.seedChallenges(tx, counts.Challenges, users, events)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed challenges: %w", err)
	}

	// Phase 3: Complex entities
	log.Println("Phase 3: Seeding complex entities...")

	// 7. Predictions (depends on users, contests, events)
	predictions, err := c.seedPredictions(tx, counts.Predictions, users, contests, events)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed predictions: %w", err)
	}

	// 8. Scoring data (depends on predictions)
	err = c.seedScoringData(tx, users, contests, predictions)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed scoring data: %w", err)
	}

	// 9. User teams (depends on users)
	userTeams, err := c.seedUserTeams(tx, counts.UserTeams, users)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed user teams: %w", err)
	}

	// 10. Notifications (depends on users)
	err = c.seedNotifications(tx, users)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed notifications: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit seeding transaction: %w", err)
	}

	log.Printf("Successfully seeded all data:")
	log.Printf("  - %d users with profiles and preferences", len(users))
	log.Printf("  - %d sports", len(sports))
	log.Printf("  - %d leagues", len(leagues))
	log.Printf("  - %d teams", len(teams))
	log.Printf("  - %d matches", len(matches))
	log.Printf("  - %d contests", len(contests))
	log.Printf("  - %d challenges", len(challenges))
	log.Printf("  - %d predictions", len(predictions))
	log.Printf("  - %d user teams", len(userTeams))

	return nil
}

// TestSeed runs the complete seeding process in a transaction and rolls it back for testing
func (c *Coordinator) TestSeed() (err error) {
	log.Println("Starting test seeding (will be rolled back)...")

	counts := c.config.GetDataCounts()

	// Start transaction for testing
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("test seeding failed with panic: %v", r)
			log.Printf("Test seeding failed with panic: %v", r)
			return
		}
		// Always rollback in test mode
		tx.Rollback()
		log.Println("Test transaction rolled back successfully")
	}()

	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to begin test transaction: %w", err)
	}

	// Run all seeding phases in test transaction
	log.Println("Testing Phase 1: Core entities...")
	users, _, _, err := c.seedUsers(tx, minInt(counts.Users, 5)) // Use smaller counts for testing
	if err != nil {
		return fmt.Errorf("test failed to seed users: %w", err)
	}

	sports, err := c.seedSports(tx, minInt(counts.Sports, 2))
	if err != nil {
		return fmt.Errorf("test failed to seed sports: %w", err)
	}

	log.Println("Testing Phase 2: Dependent entities...")
	leagues, err := c.seedLeagues(tx, minInt(counts.Leagues, 3), sports)
	if err != nil {
		return fmt.Errorf("test failed to seed leagues: %w", err)
	}

	teams, err := c.seedTeams(tx, minInt(counts.Teams, 6), sports)
	if err != nil {
		return fmt.Errorf("test failed to seed teams: %w", err)
	}

	matches, err := c.seedMatches(tx, minInt(counts.Matches, 10), leagues, teams)
	if err != nil {
		return fmt.Errorf("test failed to seed matches: %w", err)
	}

	events, err := c.seedEvents(tx, matches, teams, sports)
	if err != nil {
		return fmt.Errorf("test failed to seed events: %w", err)
	}

	contests, err := c.seedContests(tx, minInt(counts.Contests, 2), users, sports)
	if err != nil {
		return fmt.Errorf("test failed to seed contests: %w", err)
	}

	log.Println("Testing Phase 3: Complex entities...")
	predictions, err := c.seedPredictions(tx, minInt(counts.Predictions, 20), users, contests, events)
	if err != nil {
		return fmt.Errorf("test failed to seed predictions: %w", err)
	}

	err = c.seedScoringData(tx, users, contests, predictions)
	if err != nil {
		return fmt.Errorf("test failed to seed scoring data: %w", err)
	}

	_, err = c.seedUserTeams(tx, minInt(counts.UserTeams, 2), users)
	if err != nil {
		return fmt.Errorf("test failed to seed user teams: %w", err)
	}

	err = c.seedNotifications(tx, users)
	if err != nil {
		return fmt.Errorf("test failed to seed notifications: %w", err)
	}

	log.Printf("Test seeding completed successfully:")
	log.Printf("  - Generated %d users, %d sports, %d contests, %d predictions",
		len(users), len(sports), len(contests), len(predictions))
	log.Println("  - All seeding logic validated")

	return nil
}

// Individual seeding methods

func (c *Coordinator) seedUsers(tx *gorm.DB, count int) ([]*User, []*Profile, []*UserPreferences, error) {
	// First, create default admin account
	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to hash admin password: %w", err)
	}

	adminUser := &User{
		Email:    "admin@sportsprediction.com",
		Password: string(hashedPassword),
		Name:     "Admin User",
	}

	// Check if admin already exists
	var existingAdmin User
	result := tx.Where("email = ?", adminUser.Email).First(&existingAdmin)
	if result.Error == nil {
		log.Println("Admin user already exists, skipping creation")
	} else if result.Error == gorm.ErrRecordNotFound {
		// Create admin user
		if err := tx.Create(adminUser).Error; err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create admin user: %w", err)
		}

		// Create admin profile
		adminProfile := &Profile{
			UserID:            adminUser.ID,
			Bio:               "Platform Administrator",
			ProfileVisibility: "public",
		}
		if err := tx.Create(adminProfile).Error; err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create admin profile: %w", err)
		}

		// Create admin preferences
		adminPreferences := &UserPreferences{
			UserID:                adminUser.ID,
			Language:              "en",
			Timezone:              "UTC",
			EmailNotifications:    true,
			PushNotifications:     true,
			TelegramNotifications: false,
			Theme:                 "light",
		}
		if err := tx.Create(adminPreferences).Error; err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create admin preferences: %w", err)
		}

		log.Printf("✅ Created default admin account:")
		log.Printf("   Email: %s", adminUser.Email)
		log.Printf("   Password: %s", adminPassword)
		log.Printf("   ⚠️  Change this password in production!")
	} else {
		return nil, nil, nil, fmt.Errorf("failed to check for existing admin: %w", result.Error)
	}

	// Generate regular users
	users, profiles, preferences, err := c.factory.GenerateUsers(count)
	if err != nil {
		return nil, nil, nil, err
	}

	// Batch insert users first to get auto-assigned IDs
	if err := tx.CreateInBatches(users, c.config.BatchSize).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("failed to insert users: %w", err)
	}

	// Link profiles to users and batch insert
	for i, profile := range profiles {
		profile.UserID = users[i].ID
	}
	if err := tx.CreateInBatches(profiles, c.config.BatchSize).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("failed to insert profiles: %w", err)
	}

	// Link preferences to users and batch insert
	for i, preference := range preferences {
		preference.UserID = users[i].ID
	}
	if err := tx.CreateInBatches(preferences, c.config.BatchSize).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("failed to insert preferences: %w", err)
	}

	// Prepend admin user to the list
	allUsers := append([]*User{adminUser}, users...)

	return allUsers, profiles, preferences, nil
}

func (c *Coordinator) seedSports(tx *gorm.DB, count int) ([]*Sport, error) {
	sports, err := c.sports.GenerateSports(count)
	if err != nil {
		return nil, err
	}

	if err := tx.CreateInBatches(sports, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert sports: %w", err)
	}

	return sports, nil
}

func (c *Coordinator) seedLeagues(tx *gorm.DB, count int, sports []*Sport) ([]*League, error) {
	leagues, err := c.sports.GenerateLeagues(count, sports)
	if err != nil {
		return nil, err
	}

	if err := tx.CreateInBatches(leagues, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert leagues: %w", err)
	}

	return leagues, nil
}

func (c *Coordinator) seedTeams(tx *gorm.DB, count int, sports []*Sport) ([]*Team, error) {
	teams, err := c.sports.GenerateTeams(count, sports)
	if err != nil {
		return nil, err
	}

	if err := tx.CreateInBatches(teams, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert teams: %w", err)
	}

	return teams, nil
}

func (c *Coordinator) seedMatches(tx *gorm.DB, count int, leagues []*League, teams []*Team) ([]*Match, error) {
	matches, err := c.sports.GenerateMatches(count, leagues, teams)
	if err != nil {
		return nil, err
	}

	if err := tx.CreateInBatches(matches, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert matches: %w", err)
	}

	return matches, nil
}

func (c *Coordinator) seedEvents(tx *gorm.DB, matches []*Match, teams []*Team, sports []*Sport) ([]*Event, error) {
	log.Printf("Creating %d events from matches...", len(matches))
	
	// Create team lookup map
	teamMap := make(map[uint]*Team)
	for _, team := range teams {
		teamMap[team.ID] = team
	}
	
	// Create sport lookup map
	sportMap := make(map[uint]*Sport)
	for _, sport := range sports {
		sportMap[sport.ID] = sport
	}
	
	events := make([]*Event, len(matches))
	for i, match := range matches {
		homeTeam := teamMap[match.HomeTeamID]
		awayTeam := teamMap[match.AwayTeamID]
		sport := sportMap[homeTeam.SportID]
		
		// Ensure ResultData is valid JSON or empty
		resultData := match.ResultData
		if resultData == "" {
			resultData = "{}"
		}
		
		events[i] = &Event{
			ID:         match.ID, // Use same ID as match
			Title:      fmt.Sprintf("%s vs %s", homeTeam.Name, awayTeam.Name),
			SportType:  sport.Name,
			HomeTeam:   homeTeam.Name,
			AwayTeam:   awayTeam.Name,
			EventDate:  match.ScheduledAt,
			Status:     match.Status,
			ResultData: resultData,
		}
	}
	
	if err := tx.CreateInBatches(events, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert events: %w", err)
	}
	
	log.Printf("Successfully created %d events", len(events))
	return events, nil
}

func (c *Coordinator) seedContests(tx *gorm.DB, count int, users []*User, sports []*Sport) ([]*Contest, error) {
	// Extract user IDs and sport types
	userIDs := make([]uint, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	sportTypes := make([]string, len(sports))
	for i, sport := range sports {
		sportTypes[i] = sport.Name
	}

	contests, err := c.factory.GenerateContests(count, userIDs, sportTypes)
	if err != nil {
		return nil, err
	}

	if err := tx.CreateInBatches(contests, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert contests: %w", err)
	}

	return contests, nil
}

func (c *Coordinator) seedPredictions(tx *gorm.DB, count int, users []*User, contests []*Contest, events []*Event) ([]*Prediction, error) {
	// Extract IDs
	userIDs := make([]uint, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	contestIDs := make([]uint, len(contests))
	for i, contest := range contests {
		contestIDs[i] = contest.ID
	}

	eventIDs := make([]uint, len(events))
	for i, event := range events {
		eventIDs[i] = event.ID
	}

	predictions, err := c.factory.GeneratePredictions(count, userIDs, contestIDs, eventIDs)
	if err != nil {
		return nil, err
	}

	if err := tx.CreateInBatches(predictions, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert predictions: %w", err)
	}

	return predictions, nil
}

func (c *Coordinator) seedScoringData(tx *gorm.DB, users []*User, contests []*Contest, predictions []*Prediction) error {
	log.Printf("Generating scoring data for %d predictions...", len(predictions))

	// Generate scores
	scores := make([]*Score, 0, len(predictions))
	leaderboards := make(map[string]*Leaderboard) // key: "contestID_userID"
	streaks := make(map[string]*UserStreak)       // key: "contestID_userID"

	for _, prediction := range predictions {
		// Create score entry
		score := &Score{
			UserID:          prediction.UserID,
			ContestID:       prediction.ContestID,
			PredictionID:    prediction.ID,
			Points:          prediction.Points,
			TimeCoefficient: 1.0 + c.factory.faker.Float64Range(-0.2, 0.5), // Random time coefficient
			ScoredAt:        time.Now(),
		}
		scores = append(scores, score)

		// Update leaderboard
		key := fmt.Sprintf("%d_%d", prediction.ContestID, prediction.UserID)
		if leaderboard, exists := leaderboards[key]; exists {
			leaderboard.TotalPoints += score.Points
		} else {
			leaderboards[key] = &Leaderboard{
				ContestID:   prediction.ContestID,
				UserID:      prediction.UserID,
				TotalPoints: score.Points,
				Rank:        0, // Will be calculated later
			}
		}

		// Update streaks
		if streak, exists := streaks[key]; exists {
			if prediction.IsCorrect != nil && *prediction.IsCorrect {
				streak.CurrentStreak++
				if streak.CurrentStreak > streak.MaxStreak {
					streak.MaxStreak = streak.CurrentStreak
				}
			} else if prediction.IsCorrect != nil {
				streak.CurrentStreak = 0
			}
			streak.LastPredictionID = prediction.ID
			streak.LastPredictionCorrect = prediction.IsCorrect
		} else {
			currentStreak := 0
			maxStreak := 0
			if prediction.IsCorrect != nil && *prediction.IsCorrect {
				currentStreak = 1
				maxStreak = 1
			}
			streaks[key] = &UserStreak{
				UserID:                prediction.UserID,
				ContestID:             prediction.ContestID,
				CurrentStreak:         currentStreak,
				MaxStreak:             maxStreak,
				LastPredictionID:      prediction.ID,
				LastPredictionCorrect: prediction.IsCorrect,
			}
		}
	}

	// Insert scores
	if len(scores) > 0 {
		if err := tx.CreateInBatches(scores, c.config.BatchSize).Error; err != nil {
			return fmt.Errorf("failed to insert scores: %w", err)
		}
	}

	// Convert maps to slices and insert
	leaderboardSlice := make([]*Leaderboard, 0, len(leaderboards))
	for _, lb := range leaderboards {
		leaderboardSlice = append(leaderboardSlice, lb)
	}

	if len(leaderboardSlice) > 0 {
		if err := tx.CreateInBatches(leaderboardSlice, c.config.BatchSize).Error; err != nil {
			return fmt.Errorf("failed to insert leaderboards: %w", err)
		}
	}

	streakSlice := make([]*UserStreak, 0, len(streaks))
	for _, streak := range streaks {
		streakSlice = append(streakSlice, streak)
	}

	if len(streakSlice) > 0 {
		if err := tx.CreateInBatches(streakSlice, c.config.BatchSize).Error; err != nil {
			return fmt.Errorf("failed to insert streaks: %w", err)
		}
	}

	log.Printf("Successfully generated scoring data: %d scores, %d leaderboard entries, %d streaks",
		len(scores), len(leaderboardSlice), len(streakSlice))

	return nil
}

func (c *Coordinator) seedUserTeams(tx *gorm.DB, count int, users []*User) ([]*UserTeam, error) {
	log.Printf("Generating %d user teams...", count)

	teams := make([]*UserTeam, count)
	members := make([]*UserTeamMember, 0)

	for i := 0; i < count; i++ {
		captain := users[c.factory.faker.Number(0, len(users)-1)]

		team := &UserTeam{
			ID:             uint(i + 1),
			Name:           c.generateTeamName(),
			Description:    c.factory.faker.Sentence(3),
			InviteCode:     c.factory.faker.LetterN(8),
			CaptainID:      captain.ID,
			MaxMembers:     c.factory.faker.Number(5, 20),
			CurrentMembers: 0, // Will be updated as we add members
			IsActive:       true,
		}
		teams[i] = team

		// Add captain as first member
		captainMember := &UserTeamMember{
			TeamID:   team.ID,
			UserID:   captain.ID,
			Role:     "captain",
			Status:   "active",
			JoinedAt: time.Now().AddDate(0, 0, -c.factory.faker.Number(1, 30)),
		}
		members = append(members, captainMember)
		team.CurrentMembers = 1

		// Add random members
		memberCount := c.factory.faker.Number(1, minInt(team.MaxMembers-1, len(users)-1))
		usedUserIDs := map[uint]bool{captain.ID: true}

		for j := 0; j < memberCount; j++ {
			var member *User
			for {
				member = users[c.factory.faker.Number(0, len(users)-1)]
				if !usedUserIDs[member.ID] {
					usedUserIDs[member.ID] = true
					break
				}
			}

			teamMember := &UserTeamMember{
				TeamID:   team.ID,
				UserID:   member.ID,
				Role:     "member",
				Status:   "active",
				JoinedAt: time.Now().AddDate(0, 0, -c.factory.faker.Number(1, 30)),
			}
			members = append(members, teamMember)
			team.CurrentMembers++
		}
	}

	// Insert teams
	if err := tx.CreateInBatches(teams, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert user teams: %w", err)
	}

	// Insert members
	if len(members) > 0 {
		if err := tx.CreateInBatches(members, c.config.BatchSize).Error; err != nil {
			return nil, fmt.Errorf("failed to insert team members: %w", err)
		}
	}

	log.Printf("Successfully generated %d user teams with %d members", len(teams), len(members))
	return teams, nil
}

func (c *Coordinator) seedNotifications(tx *gorm.DB, users []*User) error {
	log.Printf("Generating notifications for %d users...", len(users))

	notifications := make([]*Notification, 0)
	preferences := make([]*NotificationPreference, 0)

	notificationTypes := []string{"contest_started", "prediction_reminder", "results_available", "team_invitation", "achievement_unlocked"}
	channels := []string{"in_app", "email", "telegram"}

	for _, user := range users {
		// Create notification preferences
		for _, channel := range channels {
			pref := &NotificationPreference{
				UserID:  user.ID,
				Channel: channel,
				Enabled: c.factory.faker.Bool(),
			}
			if channel == "email" {
				pref.Email = user.Email
			}
			preferences = append(preferences, pref)
		}

		// Create some notifications
		notifCount := c.factory.faker.Number(0, 5)
		for i := 0; i < notifCount; i++ {
			notif := &Notification{
				UserID:  user.ID,
				Type:    c.factory.faker.RandomString(notificationTypes),
				Title:   c.factory.faker.Sentence(2),
				Message: c.factory.faker.Sentence(5),
				Channel: c.factory.faker.RandomString(channels),
				IsRead:  c.factory.faker.Bool(),
			}

			if notif.IsRead {
				readTime := time.Now().AddDate(0, 0, -c.factory.faker.Number(1, 7))
				notif.ReadAt = &readTime
			}

			sentTime := time.Now().AddDate(0, 0, -c.factory.faker.Number(1, 14))
			notif.SentAt = &sentTime

			notifications = append(notifications, notif)
		}
	}

	// Insert preferences
	if len(preferences) > 0 {
		if err := tx.CreateInBatches(preferences, c.config.BatchSize).Error; err != nil {
			return fmt.Errorf("failed to insert notification preferences: %w", err)
		}
	}

	// Insert notifications
	if len(notifications) > 0 {
		if err := tx.CreateInBatches(notifications, c.config.BatchSize).Error; err != nil {
			return fmt.Errorf("failed to insert notifications: %w", err)
		}
	}

	log.Printf("Successfully generated %d notification preferences and %d notifications",
		len(preferences), len(notifications))

	return nil
}

// seedChallenges creates realistic challenge data
func (c *Coordinator) seedChallenges(tx *gorm.DB, count int, users []*User, events []*Event) ([]*Challenge, error) {
	log.Printf("Generating %d challenges...", count)

	if len(users) < 2 {
		return nil, fmt.Errorf("need at least 2 users to create challenges, got %d", len(users))
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("need at least 1 event to create challenges")
	}

	challenges := make([]*Challenge, 0, count)
	participants := make([]*ChallengeParticipant, 0, count*2)

	statuses := []string{"pending", "accepted", "declined", "expired", "active", "completed"}
	statusWeights := []int{30, 25, 15, 10, 10, 10} // Weighted distribution

	// Validate that statuses and weights arrays match
	if len(statuses) != len(statusWeights) {
		return nil, fmt.Errorf("status array length (%d) does not match weights array length (%d)", len(statuses), len(statusWeights))
	}

	messages := []string{
		"Let's see who's the better predictor!",
		"Think you can beat me? Prove it!",
		"Challenge accepted - may the best predictor win!",
		"Time to put your prediction skills to the test!",
		"Ready for a friendly competition?",
		"Let's make this match more interesting!",
		"I challenge you to a prediction duel!",
		"Think you know sports better than me?",
		"Game on! Let's see what you've got!",
		"",
	}

	for i := 0; i < count; i++ {
		// Select random challenger and opponent (different users)
		challengerIdx := c.factory.faker.IntRange(0, len(users)-1)
		var opponentIdx int
		
		// Guaranteed different opponent selection
		if len(users) == 2 {
			// With only 2 users, opponent is always the other user
			opponentIdx = 1 - challengerIdx
		} else {
			// With 3+ users, select a different random user
			for attempts := 0; attempts < 100; attempts++ {
				opponentIdx = c.factory.faker.IntRange(0, len(users)-1)
				if opponentIdx != challengerIdx {
					break
				}
				if attempts == 99 {
					// Guaranteed fallback: select next available user
					opponentIdx = (challengerIdx + 1) % len(users)
					if opponentIdx == challengerIdx {
						opponentIdx = (challengerIdx + 2) % len(users)
					}
				}
			}
		}

		challenger := users[challengerIdx]
		opponent := users[opponentIdx]

		// Select random match/event
		event := events[c.factory.faker.IntRange(0, len(events)-1)]

		// Select weighted random status
		totalWeight := 0
		for _, weight := range statusWeights {
			totalWeight += weight
		}
		
		randomValue := c.factory.faker.IntRange(1, totalWeight)
		currentWeight := 0
		status := statuses[0] // fallback
		
		for i, weight := range statusWeights {
			currentWeight += weight
			if randomValue <= currentWeight {
				status = statuses[i]
				break
			}
		}

		// Generate challenge
		challenge := &Challenge{
			ChallengerID: challenger.ID,
			OpponentID:   opponent.ID,
			EventID:      event.ID,
			Message:      c.factory.faker.RandomString(messages),
			Status:       status,
			ExpiresAt:    time.Now().Add(24 * time.Hour), // Default 24h expiration
		}

		// Set timestamps based on status
		now := time.Now()
		baseTime := now.Add(-time.Duration(c.factory.faker.IntRange(1, 168)) * time.Hour) // Up to 1 week ago

		challenge.CreatedAt = baseTime

		switch status {
		case "pending":
			// Still pending, expires in future
			challenge.ExpiresAt = now.Add(time.Duration(c.factory.faker.IntRange(1, 24)) * time.Hour)

		case "accepted", "active":
			// Was accepted
			acceptedTime := baseTime.Add(time.Duration(c.factory.faker.IntRange(1, 12)) * time.Hour)
			challenge.AcceptedAt = &acceptedTime
			challenge.ExpiresAt = baseTime.Add(24 * time.Hour)

		case "declined":
			// Was declined
			challenge.ExpiresAt = baseTime.Add(24 * time.Hour)

		case "expired":
			// Expired without response
			challenge.ExpiresAt = baseTime.Add(time.Duration(c.factory.faker.IntRange(12, 36)) * time.Hour)

		case "completed":
			// Was completed with scores
			acceptedTime := baseTime.Add(time.Duration(c.factory.faker.IntRange(1, 12)) * time.Hour)
			completedTime := acceptedTime.Add(time.Duration(c.factory.faker.IntRange(24, 72)) * time.Hour)
			challenge.AcceptedAt = &acceptedTime
			challenge.CompletedAt = &completedTime
			challenge.ExpiresAt = baseTime.Add(24 * time.Hour)

			// Generate realistic scores
			challengerScore := float64(c.factory.faker.IntRange(0, 20))
			opponentScore := float64(c.factory.faker.IntRange(0, 20))
			challenge.ChallengerScore = challengerScore
			challenge.OpponentScore = opponentScore

			// Determine winner
			if challengerScore > opponentScore {
				challenge.WinnerID = &challenger.ID
			} else if opponentScore > challengerScore {
				challenge.WinnerID = &opponent.ID
			}
			// Tie if scores are equal (WinnerID remains nil)
		}

		challenges = append(challenges, challenge)

		// Create challenge participants
		challengerParticipant := &ChallengeParticipant{
			UserID:   challenger.ID,
			Role:     "challenger",
			Status:   "active",
			JoinedAt: challenge.CreatedAt,
		}

		opponentParticipant := &ChallengeParticipant{
			UserID:   opponent.ID,
			Role:     "opponent",
			Status:   "active",
			JoinedAt: challenge.CreatedAt,
		}

		participants = append(participants, challengerParticipant, opponentParticipant)
	}

	// Insert challenges in batches
	if err := tx.CreateInBatches(challenges, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert challenges: %w", err)
	}

	// Update challenge IDs in participants
	for i, challenge := range challenges {
		participants[i*2].ChallengeID = challenge.ID
		participants[i*2+1].ChallengeID = challenge.ID
	}

	// Insert participants in batches
	if err := tx.CreateInBatches(participants, c.config.BatchSize).Error; err != nil {
		return nil, fmt.Errorf("failed to insert challenge participants: %w", err)
	}

	log.Printf("Successfully generated %d challenges with %d participants", len(challenges), len(participants))

	return challenges, nil
}

// Helper methods

func (c *Coordinator) generateTeamName() string {
	adjectives := []string{"Elite", "Ultimate", "Supreme", "Legendary", "Epic", "Mighty", "Fierce", "Bold", "Swift", "Thunder"}
	nouns := []string{"Predictors", "Champions", "Warriors", "Legends", "Masters", "Titans", "Eagles", "Lions", "Wolves", "Dragons"}

	return fmt.Sprintf("%s %s",
		c.factory.faker.RandomString(adjectives),
		c.factory.faker.RandomString(nouns))
}

// cleanExistingData removes all existing data in reverse dependency order
func (c *Coordinator) cleanExistingData() error {
	log.Println("Cleaning existing data...")

	// Delete in reverse dependency order (children first)
	tables := []string{
		"notifications",
		"user_team_members",
		"challenge_participants",
		"challenges",
		"user_teams",
		"user_streaks",
		"leaderboards",
		"scores",
		"predictions",
		"events",
		"contests",
		"matches",
		"teams",
		"leagues",
		"sports",
		"user_preferences",
		"profiles",
		"users",
	}

	for _, table := range tables {
		if err := c.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			// Ignore errors for tables that don't exist
			log.Printf("Warning: failed to clean table %s: %v", table, err)
		}
	}

	log.Println("✅ Existing data cleaned")
	return nil
}

// Close closes the database connection
func (c *Coordinator) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// SportsDataGenerator provides methods to generate realistic sports-specific data
type SportsDataGenerator struct {
	faker *gofakeit.Faker
}

// NewSportsDataGenerator creates a new sports data generator
func NewSportsDataGenerator(faker *gofakeit.Faker) *SportsDataGenerator {
	return &SportsDataGenerator{
		faker: faker,
	}
}

// GenerateSports creates realistic sports with proper metadata
func (g *SportsDataGenerator) GenerateSports(count int) ([]*Sport, error) {
	log.Printf("Generating %d sports...", count)

	// Predefined realistic sports data
	sportsData := []struct {
		name        string
		description string
		iconURL     string
	}{
		{"Football", "American football with touchdowns and field goals", "https://example.com/icons/football.svg"},
		{"Soccer", "Association football played with feet", "https://example.com/icons/soccer.svg"},
		{"Basketball", "Indoor sport with hoops and dribbling", "https://example.com/icons/basketball.svg"},
		{"Tennis", "Racket sport played on courts", "https://example.com/icons/tennis.svg"},
		{"Baseball", "Bat and ball sport with innings", "https://example.com/icons/baseball.svg"},
		{"Hockey", "Ice hockey with pucks and sticks", "https://example.com/icons/hockey.svg"},
		{"Golf", "Precision sport with clubs and holes", "https://example.com/icons/golf.svg"},
		{"Cricket", "Bat and ball sport with wickets", "https://example.com/icons/cricket.svg"},
	}

	sports := make([]*Sport, 0, count)

	for i := 0; i < count && i < len(sportsData); i++ {
		data := sportsData[i]
		sport := &Sport{
			Name:        data.name,
			Slug:        g.slugify(data.name),
			Description: data.description,
			IconURL:     data.iconURL,
			ExternalID:  fmt.Sprintf("sport_%d", i+1),
			IsActive:    true,
		}
		sports = append(sports, sport)
	}

	// If we need more sports than predefined, generate random ones
	for i := len(sportsData); i < count; i++ {
		sportName := g.faker.Word()
		sport := &Sport{
			Name:        sportName,
			Slug:        g.slugify(sportName),
			Description: fmt.Sprintf("Competitive %s sport", strings.ToLower(sportName)),
			IconURL:     fmt.Sprintf("https://example.com/icons/sport_%d.svg", i+1),
			ExternalID:  fmt.Sprintf("sport_%d", i+1),
			IsActive:    g.faker.Bool(),
		}
		sports = append(sports, sport)
	}

	log.Printf("Successfully generated %d sports", len(sports))
	return sports, nil
}

// GenerateLeagues creates realistic leagues for given sports
func (g *SportsDataGenerator) GenerateLeagues(count int, sports []*Sport) ([]*League, error) {
	log.Printf("Generating %d leagues...", count)

	leagues := make([]*League, count)
	countries := []string{"USA", "UK", "Germany", "Spain", "Italy", "France", "Brazil", "Argentina", "Japan", "Australia"}
	seasons := []string{"2023-24", "2024", "2024-25", "Spring 2024", "Fall 2024"}

	// Predefined league templates by sport
	leagueTemplates := map[string][]string{
		"Football":   {"NFL", "College Football", "Arena Football League"},
		"Soccer":     {"Premier League", "La Liga", "Bundesliga", "Serie A", "Ligue 1", "MLS"},
		"Basketball": {"NBA", "WNBA", "EuroLeague", "College Basketball"},
		"Tennis":     {"ATP Tour", "WTA Tour", "Grand Slam"},
		"Baseball":   {"MLB", "Minor League", "College Baseball"},
		"Hockey":     {"NHL", "KHL", "AHL"},
		"Golf":       {"PGA Tour", "European Tour", "LPGA Tour"},
		"Cricket":    {"IPL", "County Championship", "Big Bash League"},
	}

	for i := 0; i < count; i++ {
		sport := sports[i%len(sports)]

		var leagueName string
		if templates, exists := leagueTemplates[sport.Name]; exists && len(templates) > 0 {
			template := g.faker.RandomString(templates)
			country := g.faker.RandomString(countries)
			if template == "Premier League" || template == "La Liga" {
				leagueName = template
			} else {
				leagueName = fmt.Sprintf("%s %s", country, template)
			}
		} else {
			leagueName = fmt.Sprintf("%s %s League", g.faker.RandomString(countries), sport.Name)
		}

		league := &League{
			SportID:    sport.ID,
			Name:       leagueName,
			Slug:       g.slugify(leagueName),
			Country:    g.faker.RandomString(countries),
			Season:     g.faker.RandomString(seasons),
			ExternalID: fmt.Sprintf("league_%d", i+1),
			IsActive:   g.faker.Bool(),
		}
		leagues[i] = league
	}

	log.Printf("Successfully generated %d leagues", count)
	return leagues, nil
}

// GenerateTeams creates realistic teams for given sports
func (g *SportsDataGenerator) GenerateTeams(count int, sports []*Sport) ([]*Team, error) {
	log.Printf("Generating %d teams...", count)

	teams := make([]*Team, count)

	// Predefined team name components
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose", "Austin", "Jacksonville", "Fort Worth", "Columbus", "Charlotte", "San Francisco", "Indianapolis", "Seattle", "Denver", "Washington", "Boston", "El Paso", "Detroit", "Nashville", "Portland", "Oklahoma City", "Las Vegas", "Louisville", "Baltimore", "Milwaukee", "Albuquerque", "Tucson", "Fresno", "Sacramento", "Kansas City", "Mesa", "Atlanta", "Colorado Springs", "Raleigh", "Omaha", "Miami", "Oakland", "Minneapolis", "Tulsa", "Cleveland", "Wichita", "Arlington"}

	mascots := []string{"Eagles", "Lions", "Tigers", "Bears", "Wolves", "Hawks", "Falcons", "Panthers", "Jaguars", "Bulls", "Rams", "Cowboys", "Giants", "Titans", "Warriors", "Knights", "Kings", "Rangers", "Pirates", "Cardinals", "Rockets", "Thunder", "Lightning", "Storm", "Fire", "Ice", "Blazers", "Flames", "Stars", "Suns", "Moons", "Comets", "Meteors", "Hurricanes", "Tornadoes", "Earthquakes", "Avalanche", "Blizzard", "Cyclones", "Typhoons"}

	countries := []string{"USA", "Canada", "UK", "Germany", "Spain", "Italy", "France", "Brazil", "Argentina", "Mexico", "Japan", "Australia", "Netherlands", "Belgium", "Portugal", "Russia", "Sweden", "Norway", "Denmark", "Finland"}

	for i := 0; i < count; i++ {
		sport := sports[i%len(sports)]
		city := g.faker.RandomString(cities)
		mascot := g.faker.RandomString(mascots)
		teamName := fmt.Sprintf("%s %s", city, mascot)

		// Generate short name (first 3 letters of city + first 3 of mascot)
		shortName := strings.ToUpper(city[:minInt(3, len(city))] + mascot[:minInt(3, len(mascot))])

		team := &Team{
			SportID:    sport.ID,
			Name:       teamName,
			Slug:       g.slugify(teamName),
			ShortName:  shortName,
			LogoURL:    fmt.Sprintf("https://example.com/logos/%s.png", g.slugify(teamName)),
			Country:    g.faker.RandomString(countries),
			ExternalID: fmt.Sprintf("team_%d", i+1),
			IsActive:   g.faker.Bool(),
		}
		teams[i] = team
	}

	log.Printf("Successfully generated %d teams", count)
	return teams, nil
}

// GenerateMatches creates realistic matches with proper scheduling
func (g *SportsDataGenerator) GenerateMatches(count int, leagues []*League, teams []*Team) ([]*Match, error) {
	log.Printf("Generating %d matches...", count)

	matches := make([]*Match, count)

	// Create a map of teams by sport for realistic matchups
	teamsBySport := make(map[uint][]*Team)
	for _, team := range teams {
		teamsBySport[team.SportID] = append(teamsBySport[team.SportID], team)
	}

	for i := 0; i < count; i++ {
		league := leagues[i%len(leagues)]

		// Get teams for this sport
		sportTeams := teamsBySport[league.SportID]
		if len(sportTeams) < 2 {
			// If not enough teams for this sport, use any two teams
			if len(teams) < 2 {
				return nil, fmt.Errorf("not enough teams to generate matches: need at least 2 teams, got %d", len(teams))
			}
			sportTeams = teams
		}

		// Select two different teams with retry limit
		homeTeam := sportTeams[rand.Intn(len(sportTeams))]
		var awayTeam *Team
		maxRetries := 10
		for retry := 0; retry < maxRetries; retry++ {
			awayTeam = sportTeams[rand.Intn(len(sportTeams))]
			if awayTeam.ID != homeTeam.ID {
				break
			}
			if retry == maxRetries-1 {
				return nil, fmt.Errorf("unable to select different teams after %d retries: need at least 2 different teams", maxRetries)
			}
		}

		// Generate realistic match time (within next 3 months)
		scheduledAt := g.faker.DateRange(time.Now(), time.Now().AddDate(0, 3, 0))

		// Determine status based on scheduled time
		status := "scheduled"
		var homeScore, awayScore int
		if scheduledAt.Before(time.Now()) {
			status = g.faker.RandomString([]string{"completed", "live"})
			if status == "completed" {
				homeScore = g.faker.Number(0, 5)
				awayScore = g.faker.Number(0, 5)
			}
		}

		match := &Match{
			LeagueID:    league.ID,
			HomeTeamID:  homeTeam.ID,
			AwayTeamID:  awayTeam.ID,
			ScheduledAt: scheduledAt,
			Status:      status,
			HomeScore:   homeScore,
			AwayScore:   awayScore,
			ResultData:  g.generateMatchResultData(status, homeScore, awayScore),
			ExternalID:  fmt.Sprintf("match_%d", i+1),
		}
		matches[i] = match
	}

	log.Printf("Successfully generated %d matches", count)
	return matches, nil
}

// Helper methods

func (g *SportsDataGenerator) slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, "'", "")
	text = strings.ReplaceAll(text, ".", "")
	return text
}

func (g *SportsDataGenerator) generateMatchResultData(status string, homeScore, awayScore int) string {
	if status != "completed" {
		return ""
	}

	// Generate additional match statistics
	return fmt.Sprintf(`{
		"final_score": {"home": %d, "away": %d},
		"half_time_score": {"home": %d, "away": %d},
		"possession": {"home": %d, "away": %d},
		"shots": {"home": %d, "away": %d},
		"corners": {"home": %d, "away": %d},
		"fouls": {"home": %d, "away": %d},
		"yellow_cards": {"home": %d, "away": %d},
		"red_cards": {"home": %d, "away": %d}
	}`,
		homeScore, awayScore,
		g.faker.Number(0, homeScore), g.faker.Number(0, awayScore),
		g.faker.Number(30, 70), g.faker.Number(30, 70),
		g.faker.Number(5, 20), g.faker.Number(5, 20),
		g.faker.Number(0, 10), g.faker.Number(0, 10),
		g.faker.Number(5, 25), g.faker.Number(5, 25),
		g.faker.Number(0, 5), g.faker.Number(0, 5),
		g.faker.Number(0, 2), g.faker.Number(0, 2))
}

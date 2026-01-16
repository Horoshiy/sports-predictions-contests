package sync

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sports-prediction-contests/sports-service/internal/external"
	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
)

// Rate limit delay between API calls to avoid throttling
const apiRateLimitDelay = 100 * time.Millisecond

// SyncService orchestrates syncing from external API to database
type SyncService struct {
	client     *external.Client
	sportRepo  repository.SportRepositoryInterface
	leagueRepo repository.LeagueRepositoryInterface
	teamRepo   repository.TeamRepositoryInterface
	matchRepo  repository.MatchRepositoryInterface
	lastSyncAt *time.Time
	mu         sync.RWMutex
}

// NewSyncService creates a new sync service
func NewSyncService(
	client *external.Client,
	sportRepo repository.SportRepositoryInterface,
	leagueRepo repository.LeagueRepositoryInterface,
	teamRepo repository.TeamRepositoryInterface,
	matchRepo repository.MatchRepositoryInterface,
) *SyncService {
	return &SyncService{
		client:     client,
		sportRepo:  sportRepo,
		leagueRepo: leagueRepo,
		teamRepo:   teamRepo,
		matchRepo:  matchRepo,
	}
}

// GetLastSyncAt returns the last sync timestamp (thread-safe)
func (s *SyncService) GetLastSyncAt() *time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSyncAt
}

// setLastSyncAt updates the last sync timestamp (thread-safe)
func (s *SyncService) setLastSyncAt() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.lastSyncAt = &now
}

// SyncSports syncs all sports from external API
func (s *SyncService) SyncSports(ctx context.Context) (int, error) {
	sports, err := s.client.GetAllSports()
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, apiSport := range sports {
		sport := s.mapSportFromExternal(apiSport)
		if err := s.sportRepo.Upsert(sport); err != nil {
			log.Printf("[WARN] Failed to upsert sport %s: %v", apiSport.StrSport, err)
			continue
		}
		synced++
	}

	s.setLastSyncAt()
	log.Printf("[INFO] Synced %d sports", synced)
	return synced, nil
}

// SyncLeagues syncs leagues from external API
func (s *SyncService) SyncLeagues(ctx context.Context) (int, error) {
	leagues, err := s.client.GetAllLeagues()
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, apiLeague := range leagues {
		// Find sport by name
		sport, err := s.sportRepo.GetBySlug(slugify(apiLeague.StrSport))
		if err != nil {
			continue // Skip if sport not found
		}

		league := s.mapLeagueFromExternal(apiLeague, sport.ID)
		if err := s.leagueRepo.Upsert(league); err != nil {
			log.Printf("[WARN] Failed to upsert league %s: %v", apiLeague.StrLeague, err)
			continue
		}
		synced++
	}

	s.setLastSyncAt()
	log.Printf("[INFO] Synced %d leagues", synced)
	return synced, nil
}

// SyncTeamsByLeague syncs teams for a specific league
func (s *SyncService) SyncTeamsByLeague(ctx context.Context, leagueID uint) (int, error) {
	league, err := s.leagueRepo.GetByID(leagueID)
	if err != nil {
		return 0, err
	}

	if league.ExternalID == "" {
		return 0, nil // No external ID, can't sync
	}

	teams, err := s.client.GetTeamsByLeague(league.ExternalID)
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, apiTeam := range teams {
		// Find sport by name
		sport, err := s.sportRepo.GetBySlug(slugify(apiTeam.StrSport))
		if err != nil {
			continue
		}

		team := s.mapTeamFromExternal(apiTeam, sport.ID)
		if err := s.teamRepo.Upsert(team); err != nil {
			log.Printf("[WARN] Failed to upsert team %s: %v", apiTeam.StrTeam, err)
			continue
		}
		synced++
	}

	s.setLastSyncAt()
	log.Printf("[INFO] Synced %d teams for league %d", synced, leagueID)
	return synced, nil
}

// SyncUpcomingMatches syncs upcoming matches for teams
func (s *SyncService) SyncUpcomingMatches(ctx context.Context, teamID uint) (int, error) {
	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return 0, err
	}

	if team.ExternalID == "" {
		return 0, nil
	}

	events, err := s.client.GetUpcomingEventsByTeam(team.ExternalID)
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, apiEvent := range events {
		match, err := s.mapMatchFromExternal(apiEvent)
		if err != nil {
			log.Printf("[WARN] Failed to map event %s: %v", apiEvent.IDEvent, err)
			continue
		}

		if err := s.matchRepo.Upsert(match); err != nil {
			log.Printf("[WARN] Failed to upsert match %s: %v", apiEvent.IDEvent, err)
			continue
		}
		synced++
	}

	s.setLastSyncAt()
	log.Printf("[INFO] Synced %d matches for team %d", synced, teamID)
	return synced, nil
}

// SyncMatchResults updates results for completed matches
func (s *SyncService) SyncMatchResults(ctx context.Context) (int, error) {
	// Get matches that are scheduled or live
	matches, _, err := s.matchRepo.List(100, 0, 0, 0, "scheduled")
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, match := range matches {
		if match.ExternalID == "" {
			continue
		}

		// Rate limit API calls to avoid throttling
		time.Sleep(apiRateLimitDelay)

		event, err := s.client.GetEventByID(match.ExternalID)
		if err != nil {
			continue
		}

		// Map API status to internal status and compare
		mappedStatus := s.mapStatus(event.StrStatus)
		if mappedStatus != match.Status {
			match.Status = mappedStatus
			if event.IntHomeScore != "" {
				match.HomeScore, _ = strconv.Atoi(event.IntHomeScore)
			}
			if event.IntAwayScore != "" {
				match.AwayScore, _ = strconv.Atoi(event.IntAwayScore)
			}

			if err := s.matchRepo.Update(match); err != nil {
				log.Printf("[WARN] Failed to update match %d: %v", match.ID, err)
				continue
			}
			synced++
		}
	}

	s.setLastSyncAt()
	log.Printf("[INFO] Updated %d match results", synced)
	return synced, nil
}

// Helper methods

func (s *SyncService) mapSportFromExternal(api external.APISport) *models.Sport {
	return &models.Sport{
		Name:        api.StrSport,
		Slug:        slugify(api.StrSport),
		Description: api.StrFormat,
		IconURL:     api.StrSportThumb,
		ExternalID:  api.IDSport,
		IsActive:    true,
	}
}

func (s *SyncService) mapLeagueFromExternal(api external.APILeague, sportID uint) *models.League {
	return &models.League{
		SportID:    sportID,
		Name:       api.StrLeague,
		Slug:       slugify(api.StrLeague),
		Country:    api.StrCountry,
		ExternalID: api.IDLeague,
		IsActive:   true,
	}
}

func (s *SyncService) mapTeamFromExternal(api external.APITeam, sportID uint) *models.Team {
	return &models.Team{
		SportID:    sportID,
		Name:       api.StrTeam,
		Slug:       slugify(api.StrTeam),
		ShortName:  api.StrTeamShort,
		LogoURL:    api.StrTeamBadge,
		Country:    api.StrCountry,
		ExternalID: api.IDTeam,
		IsActive:   true,
	}
}

func (s *SyncService) mapMatchFromExternal(api external.APIEvent) (*models.Match, error) {
	// Find league by external ID
	league, err := s.leagueRepo.GetByExternalID(api.IDLeague)
	if err != nil {
		return nil, err
	}

	// Find teams by external ID
	homeTeam, err := s.teamRepo.GetByExternalID(api.IDHomeTeam)
	if err != nil {
		return nil, err
	}
	awayTeam, err := s.teamRepo.GetByExternalID(api.IDAwayTeam)
	if err != nil {
		return nil, err
	}

	// Parse scheduled time
	scheduledAt, err := time.Parse("2006-01-02 15:04:05", api.StrTimestamp)
	if err != nil {
		// Try alternative format
		scheduledAt, err = time.Parse("2006-01-02", api.DateEvent)
		if err != nil {
			log.Printf("[WARN] Failed to parse date for event %s, using fallback", api.IDEvent)
			scheduledAt = time.Now().Add(24 * time.Hour) // Default to tomorrow
		}
	}

	homeScore, _ := strconv.Atoi(api.IntHomeScore)
	awayScore, _ := strconv.Atoi(api.IntAwayScore)

	return &models.Match{
		LeagueID:    league.ID,
		HomeTeamID:  homeTeam.ID,
		AwayTeamID:  awayTeam.ID,
		ScheduledAt: scheduledAt,
		Status:      s.mapStatus(api.StrStatus),
		HomeScore:   homeScore,
		AwayScore:   awayScore,
		ExternalID:  api.IDEvent,
	}, nil
}

func (s *SyncService) mapStatus(apiStatus string) string {
	switch strings.ToUpper(apiStatus) {
	case "FT", "AOT", "AET":
		return "completed"
	case "CANC":
		return "cancelled"
	case "PST":
		return "postponed"
	case "1H", "2H", "HT", "ET", "P", "LIVE":
		return "live"
	default:
		return "scheduled"
	}
}

// slugify converts a name to a URL-friendly slug
func slugify(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Use simple character filtering instead of regex for performance
	var result strings.Builder
	lastWasDash := false
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
			lastWasDash = false
		} else if r == '-' && !lastWasDash {
			result.WriteRune('-')
			lastWasDash = true
		}
	}
	return strings.Trim(result.String(), "-")
}

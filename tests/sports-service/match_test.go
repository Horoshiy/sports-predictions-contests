package sports_service_test

import (
	"testing"
	"time"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
)

func TestMatchValidation(t *testing.T) {
	match := &models.Match{
		LeagueID:    1,
		HomeTeamID:  1,
		AwayTeamID:  2,
		ScheduledAt: time.Now().Add(24 * time.Hour),
		Status:      "scheduled",
	}

	if err := match.ValidateTeams(); err != nil {
		t.Errorf("Expected valid teams, got error: %v", err)
	}

	match.HomeTeamID = match.AwayTeamID
	if err := match.ValidateTeams(); err == nil {
		t.Error("Expected error for same home and away team")
	}
}

func TestMatchStatusValidation(t *testing.T) {
	match := &models.Match{Status: "scheduled"}
	if err := match.ValidateStatus(); err != nil {
		t.Errorf("Expected valid status, got error: %v", err)
	}

	match.Status = "invalid"
	if err := match.ValidateStatus(); err == nil {
		t.Error("Expected error for invalid status")
	}
}

func TestMatchRepository(t *testing.T) {
	db := setupTestDB(t)
	sportRepo := repository.NewSportRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football"}
	sportRepo.Create(sport)

	league := &models.League{SportID: sport.ID, Name: "Premier League", Slug: "premier-league"}
	leagueRepo.Create(league)

	homeTeam := &models.Team{SportID: sport.ID, Name: "Team A", Slug: "team-a"}
	awayTeam := &models.Team{SportID: sport.ID, Name: "Team B", Slug: "team-b"}
	teamRepo.Create(homeTeam)
	teamRepo.Create(awayTeam)

	match := &models.Match{
		LeagueID:    league.ID,
		HomeTeamID:  homeTeam.ID,
		AwayTeamID:  awayTeam.ID,
		ScheduledAt: time.Now().Add(24 * time.Hour),
		Status:      "scheduled",
	}

	if err := matchRepo.Create(match); err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}

	retrieved, err := matchRepo.GetByID(match.ID)
	if err != nil {
		t.Fatalf("Failed to get match: %v", err)
	}

	if retrieved.LeagueID != league.ID {
		t.Errorf("Expected league ID %d, got %d", league.ID, retrieved.LeagueID)
	}
}

func TestMatchListByLeague(t *testing.T) {
	db := setupTestDB(t)
	sportRepo := repository.NewSportRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football"}
	sportRepo.Create(sport)

	league1 := &models.League{SportID: sport.ID, Name: "League 1", Slug: "league-1"}
	league2 := &models.League{SportID: sport.ID, Name: "League 2", Slug: "league-2"}
	leagueRepo.Create(league1)
	leagueRepo.Create(league2)

	team1 := &models.Team{SportID: sport.ID, Name: "Team 1", Slug: "team-1"}
	team2 := &models.Team{SportID: sport.ID, Name: "Team 2", Slug: "team-2"}
	teamRepo.Create(team1)
	teamRepo.Create(team2)

	match1 := &models.Match{LeagueID: league1.ID, HomeTeamID: team1.ID, AwayTeamID: team2.ID, ScheduledAt: time.Now(), Status: "scheduled"}
	match2 := &models.Match{LeagueID: league2.ID, HomeTeamID: team1.ID, AwayTeamID: team2.ID, ScheduledAt: time.Now(), Status: "scheduled"}
	matchRepo.Create(match1)
	matchRepo.Create(match2)

	matches, total, err := matchRepo.List(10, 0, league1.ID, 0, "")
	if err != nil {
		t.Fatalf("Failed to list matches: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected 1 match for league1, got %d", total)
	}

	if len(matches) != 1 {
		t.Errorf("Expected 1 match in list, got %d", len(matches))
	}
}

func TestLeagueRepository(t *testing.T) {
	db := setupTestDB(t)
	sportRepo := repository.NewSportRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football"}
	sportRepo.Create(sport)

	league := &models.League{SportID: sport.ID, Name: "Premier League", Slug: "premier-league", Country: "England"}
	if err := leagueRepo.Create(league); err != nil {
		t.Fatalf("Failed to create league: %v", err)
	}

	retrieved, err := leagueRepo.GetByID(league.ID)
	if err != nil {
		t.Fatalf("Failed to get league: %v", err)
	}

	if retrieved.Name != "Premier League" {
		t.Errorf("Expected Premier League, got %s", retrieved.Name)
	}
}

func TestTeamRepository(t *testing.T) {
	db := setupTestDB(t)
	sportRepo := repository.NewSportRepository(db)
	teamRepo := repository.NewTeamRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football"}
	sportRepo.Create(sport)

	team := &models.Team{SportID: sport.ID, Name: "Manchester United", Slug: "manchester-united", ShortName: "MUN"}
	if err := teamRepo.Create(team); err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}

	retrieved, err := teamRepo.GetByID(team.ID)
	if err != nil {
		t.Fatalf("Failed to get team: %v", err)
	}

	if retrieved.ShortName != "MUN" {
		t.Errorf("Expected MUN, got %s", retrieved.ShortName)
	}
}

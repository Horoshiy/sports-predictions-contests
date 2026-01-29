package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
	"github.com/sports-prediction-contests/sports-service/internal/sync"
	"github.com/sports-prediction-contests/shared/proto/common"
	pb "github.com/sports-prediction-contests/shared/proto/sports"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SportsService struct {
	pb.UnimplementedSportsServiceServer
	sportRepo    repository.SportRepositoryInterface
	leagueRepo   repository.LeagueRepositoryInterface
	teamRepo     repository.TeamRepositoryInterface
	matchRepo    repository.MatchRepositoryInterface
	syncWorker   *sync.SyncWorker
	syncEnabled  bool
	syncInterval int
}

func NewSportsService(
	sportRepo repository.SportRepositoryInterface,
	leagueRepo repository.LeagueRepositoryInterface,
	teamRepo repository.TeamRepositoryInterface,
	matchRepo repository.MatchRepositoryInterface,
) *SportsService {
	return &SportsService{
		sportRepo:  sportRepo,
		leagueRepo: leagueRepo,
		teamRepo:   teamRepo,
		matchRepo:  matchRepo,
	}
}

// SetSyncWorker sets the sync worker for the service
func (s *SportsService) SetSyncWorker(worker *sync.SyncWorker, enabled bool, interval int) {
	s.syncWorker = worker
	s.syncEnabled = enabled
	s.syncInterval = interval
}

func (s *SportsService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return &common.Response{
		Success:   true,
		Message:   "Sports Service is healthy",
		Code:      0,
		Timestamp: timestamppb.Now(),
	}, nil
}

// Sport methods
func (s *SportsService) CreateSport(ctx context.Context, req *pb.CreateSportRequest) (*pb.SportResponse, error) {
	sport := &models.Sport{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IconURL:     req.IconUrl,
		IsActive:    true,
	}

	if err := s.sportRepo.Create(sport); err != nil {
		log.Printf("[ERROR] Failed to create sport: %v", err)
		return &pb.SportResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	log.Printf("[INFO] Sport created: %d", sport.ID)
	return &pb.SportResponse{
		Response: &common.Response{Success: true, Message: "Sport created successfully", Code: 0, Timestamp: timestamppb.Now()},
		Sport:    s.sportToProto(sport),
	}, nil
}

func (s *SportsService) GetSport(ctx context.Context, req *pb.GetSportRequest) (*pb.SportResponse, error) {
	sport, err := s.sportRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.SportResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.SportResponse{
		Response: &common.Response{Success: true, Message: "Sport retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Sport:    s.sportToProto(sport),
	}, nil
}

func (s *SportsService) UpdateSport(ctx context.Context, req *pb.UpdateSportRequest) (*pb.SportResponse, error) {
	sport, err := s.sportRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.SportResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}

	sport.Name = req.Name
	sport.Slug = req.Slug
	sport.Description = req.Description
	sport.IconURL = req.IconUrl
	sport.IsActive = req.IsActive

	if err := s.sportRepo.Update(sport); err != nil {
		return &pb.SportResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.SportResponse{
		Response: &common.Response{Success: true, Message: "Sport updated successfully", Code: 0, Timestamp: timestamppb.Now()},
		Sport:    s.sportToProto(sport),
	}, nil
}

func (s *SportsService) DeleteSport(ctx context.Context, req *pb.DeleteSportRequest) (*pb.DeleteResponse, error) {
	// Check if sport has leagues (safe delete)
	leagues, _, err := s.leagueRepo.List(100, 0, uint(req.Id), false)
	if err == nil && len(leagues) > 0 {
		return &pb.DeleteResponse{
			Response: &common.Response{
				Success:   false,
				Message:   fmt.Sprintf("Cannot delete sport with %d leagues. Delete leagues first.", len(leagues)),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	if err := s.sportRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.DeleteResponse{
		Response: &common.Response{Success: true, Message: "Sport deleted successfully", Code: 0, Timestamp: timestamppb.Now()},
	}, nil
}

func (s *SportsService) ListSports(ctx context.Context, req *pb.ListSportsRequest) (*pb.ListSportsResponse, error) {
	if req.Pagination == nil {
		req.Pagination = &common.PaginationRequest{Page: 1, Limit: 10}
	}
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
		req.Pagination.Limit = 10
	}
	page := int(req.Pagination.Page)
	if page <= 0 {
		page = 1
		req.Pagination.Page = 1
	}
	offset := (page - 1) * limit

	sports, total, err := s.sportRepo.List(limit, offset, req.ActiveOnly)
	if err != nil {
		return &pb.ListSportsResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INTERNAL_ERROR), Timestamp: timestamppb.Now()},
		}, nil
	}

	pbSports := make([]*pb.Sport, len(sports))
	for i, sport := range sports {
		pbSports[i] = s.sportToProto(sport)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListSportsResponse{
		Response:   &common.Response{Success: true, Message: "Sports retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Sports:     pbSports,
		Pagination: &common.PaginationResponse{Page: req.Pagination.Page, Limit: req.Pagination.Limit, Total: int32(total), TotalPages: totalPages},
	}, nil
}

// League methods
func (s *SportsService) CreateLeague(ctx context.Context, req *pb.CreateLeagueRequest) (*pb.LeagueResponse, error) {
	// Validate sport exists
	if _, err := s.sportRepo.GetByID(uint(req.SportId)); err != nil {
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: "sport not found", Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}

	league := &models.League{
		SportID:  uint(req.SportId),
		Name:     req.Name,
		Slug:     req.Slug,
		Country:  req.Country,
		Season:   req.Season,
		IsActive: true,
	}

	if err := s.leagueRepo.Create(league); err != nil {
		log.Printf("[ERROR] Failed to create league: %v", err)
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	log.Printf("[INFO] League created: %d", league.ID)
	return &pb.LeagueResponse{
		Response: &common.Response{Success: true, Message: "League created successfully", Code: 0, Timestamp: timestamppb.Now()},
		League:   s.leagueToProto(league),
	}, nil
}

func (s *SportsService) GetLeague(ctx context.Context, req *pb.GetLeagueRequest) (*pb.LeagueResponse, error) {
	league, err := s.leagueRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.LeagueResponse{
		Response: &common.Response{Success: true, Message: "League retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		League:   s.leagueToProto(league),
	}, nil
}

func (s *SportsService) UpdateLeague(ctx context.Context, req *pb.UpdateLeagueRequest) (*pb.LeagueResponse, error) {
	league, err := s.leagueRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	league.SportID = uint(req.SportId)
	league.Name = req.Name
	league.Slug = req.Slug
	league.Country = req.Country
	league.Season = req.Season
	league.IsActive = req.IsActive

	if err := s.leagueRepo.Update(league); err != nil {
		return &pb.LeagueResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.LeagueResponse{
		Response: &common.Response{Success: true, Message: "League updated successfully", Code: 0, Timestamp: timestamppb.Now()},
		League:   s.leagueToProto(league),
	}, nil
}

func (s *SportsService) DeleteLeague(ctx context.Context, req *pb.DeleteLeagueRequest) (*pb.DeleteResponse, error) {
	// Check if league has teams (safe delete)
	teams, _, err := s.teamRepo.List(100, 0, uint(req.Id), false)
	if err == nil && len(teams) > 0 {
		return &pb.DeleteResponse{
			Response: &common.Response{
				Success:   false,
				Message:   fmt.Sprintf("Cannot delete league with %d teams. Delete teams first.", len(teams)),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Check if league has matches
	matches, _, err := s.matchRepo.List(100, 0, uint(req.Id), 0, "")
	if err == nil && len(matches) > 0 {
		return &pb.DeleteResponse{
			Response: &common.Response{
				Success:   false,
				Message:   fmt.Sprintf("Cannot delete league with %d matches. Delete matches first.", len(matches)),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	if err := s.leagueRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.DeleteResponse{
		Response: &common.Response{Success: true, Message: "League deleted successfully", Code: 0, Timestamp: timestamppb.Now()},
	}, nil
}

func (s *SportsService) ListLeagues(ctx context.Context, req *pb.ListLeaguesRequest) (*pb.ListLeaguesResponse, error) {
	if req.Pagination == nil {
		req.Pagination = &common.PaginationRequest{Page: 1, Limit: 10}
	}
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
		req.Pagination.Limit = 10
	}
	page := int(req.Pagination.Page)
	if page <= 0 {
		page = 1
		req.Pagination.Page = 1
	}
	offset := (page - 1) * limit

	leagues, total, err := s.leagueRepo.List(limit, offset, uint(req.SportId), req.ActiveOnly)
	if err != nil {
		return &pb.ListLeaguesResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INTERNAL_ERROR), Timestamp: timestamppb.Now()},
		}, nil
	}

	pbLeagues := make([]*pb.League, len(leagues))
	for i, league := range leagues {
		pbLeagues[i] = s.leagueToProto(league)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListLeaguesResponse{
		Response:   &common.Response{Success: true, Message: "Leagues retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Leagues:    pbLeagues,
		Pagination: &common.PaginationResponse{Page: req.Pagination.Page, Limit: req.Pagination.Limit, Total: int32(total), TotalPages: totalPages},
	}, nil
}

// Team methods
func (s *SportsService) CreateTeam(ctx context.Context, req *pb.CreateTeamRequest) (*pb.TeamResponse, error) {
	// Validate sport exists
	if _, err := s.sportRepo.GetByID(uint(req.SportId)); err != nil {
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: "sport not found", Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}

	team := &models.Team{
		SportID:   uint(req.SportId),
		Name:      req.Name,
		Slug:      req.Slug,
		ShortName: req.ShortName,
		LogoURL:   req.LogoUrl,
		Country:   req.Country,
		IsActive:  true,
	}

	if err := s.teamRepo.Create(team); err != nil {
		log.Printf("[ERROR] Failed to create team: %v", err)
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	log.Printf("[INFO] Team created: %d", team.ID)
	return &pb.TeamResponse{
		Response: &common.Response{Success: true, Message: "Team created successfully", Code: 0, Timestamp: timestamppb.Now()},
		Team:     s.teamToProto(team),
	}, nil
}

func (s *SportsService) GetTeam(ctx context.Context, req *pb.GetTeamRequest) (*pb.TeamResponse, error) {
	team, err := s.teamRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.TeamResponse{
		Response: &common.Response{Success: true, Message: "Team retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Team:     s.teamToProto(team),
	}, nil
}

func (s *SportsService) UpdateTeam(ctx context.Context, req *pb.UpdateTeamRequest) (*pb.TeamResponse, error) {
	team, err := s.teamRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	team.SportID = uint(req.SportId)
	team.Name = req.Name
	team.Slug = req.Slug
	team.ShortName = req.ShortName
	team.LogoURL = req.LogoUrl
	team.Country = req.Country
	team.IsActive = req.IsActive

	if err := s.teamRepo.Update(team); err != nil {
		return &pb.TeamResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.TeamResponse{
		Response: &common.Response{Success: true, Message: "Team updated successfully", Code: 0, Timestamp: timestamppb.Now()},
		Team:     s.teamToProto(team),
	}, nil
}

func (s *SportsService) DeleteTeam(ctx context.Context, req *pb.DeleteTeamRequest) (*pb.DeleteResponse, error) {
	// Check if team has matches (safe delete)
	matches, _, err := s.matchRepo.List(100, 0, 0, uint(req.Id), "")
	if err == nil && len(matches) > 0 {
		return &pb.DeleteResponse{
			Response: &common.Response{
				Success:   false,
				Message:   fmt.Sprintf("Cannot delete team with %d matches. Delete matches first.", len(matches)),
				Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	if err := s.teamRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.DeleteResponse{
		Response: &common.Response{Success: true, Message: "Team deleted successfully", Code: 0, Timestamp: timestamppb.Now()},
	}, nil
}

func (s *SportsService) ListTeams(ctx context.Context, req *pb.ListTeamsRequest) (*pb.ListTeamsResponse, error) {
	if req.Pagination == nil {
		req.Pagination = &common.PaginationRequest{Page: 1, Limit: 10}
	}
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
		req.Pagination.Limit = 10
	}
	page := int(req.Pagination.Page)
	if page <= 0 {
		page = 1
		req.Pagination.Page = 1
	}
	offset := (page - 1) * limit

	teams, total, err := s.teamRepo.List(limit, offset, uint(req.SportId), req.ActiveOnly)
	if err != nil {
		return &pb.ListTeamsResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INTERNAL_ERROR), Timestamp: timestamppb.Now()},
		}, nil
	}

	pbTeams := make([]*pb.Team, len(teams))
	for i, team := range teams {
		pbTeams[i] = s.teamToProto(team)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListTeamsResponse{
		Response:   &common.Response{Success: true, Message: "Teams retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Teams:      pbTeams,
		Pagination: &common.PaginationResponse{Page: req.Pagination.Page, Limit: req.Pagination.Limit, Total: int32(total), TotalPages: totalPages},
	}, nil
}

// Match methods
func (s *SportsService) CreateMatch(ctx context.Context, req *pb.CreateMatchRequest) (*pb.MatchResponse, error) {
	if req.ScheduledAt == nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: "scheduled_at is required", Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	// Validate league exists
	if _, err := s.leagueRepo.GetByID(uint(req.LeagueId)); err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: "league not found", Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}

	// Validate teams exist
	if _, err := s.teamRepo.GetByID(uint(req.HomeTeamId)); err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: "home team not found", Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	if _, err := s.teamRepo.GetByID(uint(req.AwayTeamId)); err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: "away team not found", Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}

	match := &models.Match{
		LeagueID:    uint(req.LeagueId),
		HomeTeamID:  uint(req.HomeTeamId),
		AwayTeamID:  uint(req.AwayTeamId),
		ScheduledAt: req.ScheduledAt.AsTime(),
		Status:      "scheduled",
	}

	if err := s.matchRepo.Create(match); err != nil {
		log.Printf("[ERROR] Failed to create match: %v", err)
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	log.Printf("[INFO] Match created: %d", match.ID)
	return &pb.MatchResponse{
		Response: &common.Response{Success: true, Message: "Match created successfully", Code: 0, Timestamp: timestamppb.Now()},
		Match:    s.matchToProto(match),
	}, nil
}

func (s *SportsService) GetMatch(ctx context.Context, req *pb.GetMatchRequest) (*pb.MatchResponse, error) {
	match, err := s.matchRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.MatchResponse{
		Response: &common.Response{Success: true, Message: "Match retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Match:    s.matchToProto(match),
	}, nil
}

func (s *SportsService) UpdateMatch(ctx context.Context, req *pb.UpdateMatchRequest) (*pb.MatchResponse, error) {
	if req.ScheduledAt == nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: "scheduled_at is required", Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}

	match, err := s.matchRepo.GetByID(uint(req.Id))
	if err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	match.LeagueID = uint(req.LeagueId)
	match.HomeTeamID = uint(req.HomeTeamId)
	match.AwayTeamID = uint(req.AwayTeamId)
	match.ScheduledAt = req.ScheduledAt.AsTime()
	match.Status = req.Status
	match.HomeScore = int(req.HomeScore)
	match.AwayScore = int(req.AwayScore)
	match.ResultData = req.ResultData

	if err := s.matchRepo.Update(match); err != nil {
		return &pb.MatchResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INVALID_ARGUMENT), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.MatchResponse{
		Response: &common.Response{Success: true, Message: "Match updated successfully", Code: 0, Timestamp: timestamppb.Now()},
		Match:    s.matchToProto(match),
	}, nil
}

func (s *SportsService) DeleteMatch(ctx context.Context, req *pb.DeleteMatchRequest) (*pb.DeleteResponse, error) {
	if err := s.matchRepo.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_NOT_FOUND), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.DeleteResponse{
		Response: &common.Response{Success: true, Message: "Match deleted successfully", Code: 0, Timestamp: timestamppb.Now()},
	}, nil
}

func (s *SportsService) ListMatches(ctx context.Context, req *pb.ListMatchesRequest) (*pb.ListMatchesResponse, error) {
	if req.Pagination == nil {
		req.Pagination = &common.PaginationRequest{Page: 1, Limit: 10}
	}
	limit := int(req.Pagination.Limit)
	if limit <= 0 {
		limit = 10
		req.Pagination.Limit = 10
	}
	page := int(req.Pagination.Page)
	if page <= 0 {
		page = 1
		req.Pagination.Page = 1
	}
	offset := (page - 1) * limit

	matches, total, err := s.matchRepo.List(limit, offset, uint(req.LeagueId), uint(req.TeamId), req.Status)
	if err != nil {
		return &pb.ListMatchesResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Code: int32(common.ErrorCode_INTERNAL_ERROR), Timestamp: timestamppb.Now()},
		}, nil
	}

	pbMatches := make([]*pb.Match, len(matches))
	for i, match := range matches {
		pbMatches[i] = s.matchToProto(match)
	}

	totalPages := int32(total) / req.Pagination.Limit
	if int32(total)%req.Pagination.Limit > 0 {
		totalPages++
	}

	return &pb.ListMatchesResponse{
		Response:   &common.Response{Success: true, Message: "Matches retrieved successfully", Code: 0, Timestamp: timestamppb.Now()},
		Matches:    pbMatches,
		Pagination: &common.PaginationResponse{Page: req.Pagination.Page, Limit: req.Pagination.Limit, Total: int32(total), TotalPages: totalPages},
	}, nil
}

// Helper conversion methods
func (s *SportsService) sportToProto(sport *models.Sport) *pb.Sport {
	return &pb.Sport{
		Id:          uint32(sport.ID),
		Name:        sport.Name,
		Slug:        sport.Slug,
		Description: sport.Description,
		IconUrl:     sport.IconURL,
		IsActive:    sport.IsActive,
		CreatedAt:   timestamppb.New(sport.CreatedAt),
		UpdatedAt:   timestamppb.New(sport.UpdatedAt),
	}
}

func (s *SportsService) leagueToProto(league *models.League) *pb.League {
	return &pb.League{
		Id:        uint32(league.ID),
		SportId:   uint32(league.SportID),
		Name:      league.Name,
		Slug:      league.Slug,
		Country:   league.Country,
		Season:    league.Season,
		IsActive:  league.IsActive,
		CreatedAt: timestamppb.New(league.CreatedAt),
		UpdatedAt: timestamppb.New(league.UpdatedAt),
	}
}

func (s *SportsService) teamToProto(team *models.Team) *pb.Team {
	return &pb.Team{
		Id:        uint32(team.ID),
		SportId:   uint32(team.SportID),
		Name:      team.Name,
		Slug:      team.Slug,
		ShortName: team.ShortName,
		LogoUrl:   team.LogoURL,
		Country:   team.Country,
		IsActive:  team.IsActive,
		CreatedAt: timestamppb.New(team.CreatedAt),
		UpdatedAt: timestamppb.New(team.UpdatedAt),
	}
}

func (s *SportsService) matchToProto(match *models.Match) *pb.Match {
	return &pb.Match{
		Id:          uint32(match.ID),
		LeagueId:    uint32(match.LeagueID),
		HomeTeamId:  uint32(match.HomeTeamID),
		AwayTeamId:  uint32(match.AwayTeamID),
		ScheduledAt: timestamppb.New(match.ScheduledAt),
		Status:      match.Status,
		HomeScore:   int32(match.HomeScore),
		AwayScore:   int32(match.AwayScore),
		ResultData:  match.ResultData,
		CreatedAt:   timestamppb.New(match.CreatedAt),
		UpdatedAt:   timestamppb.New(match.UpdatedAt),
	}
}


// TriggerSync manually triggers a sync operation
func (s *SportsService) TriggerSync(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	if s.syncWorker == nil {
		return &pb.SyncResponse{
			Response: &common.Response{
				Success:   false,
				Message:   "Sync is not configured",
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	count, err := s.syncWorker.TriggerSync(ctx, req.EntityType, uint(req.ParentId))
	if err != nil {
		log.Printf("[ERROR] Sync failed for %s: %v", req.EntityType, err)
		return &pb.SyncResponse{
			Response: &common.Response{
				Success:   false,
				Message:   err.Error(),
				Code:      int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.SyncResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Sync completed successfully",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		SyncedCount: int32(count),
		EntityType:  req.EntityType,
	}, nil
}

// GetSyncStatus returns the current sync status
func (s *SportsService) GetSyncStatus(ctx context.Context, req *pb.SyncStatusRequest) (*pb.SyncStatusResponse, error) {
	lastSyncAt := ""
	if s.syncWorker != nil {
		if lastSync := s.syncWorker.GetLastSyncAt(); lastSync != nil {
			lastSyncAt = lastSync.Format(time.RFC3339)
		}
	}

	return &pb.SyncStatusResponse{
		Response: &common.Response{
			Success:   true,
			Message:   "Sync status retrieved",
			Code:      0,
			Timestamp: timestamppb.Now(),
		},
		SyncEnabled:      s.syncEnabled,
		LastSyncAt:       lastSyncAt,
		SyncIntervalMins: int32(s.syncInterval),
	}, nil
}

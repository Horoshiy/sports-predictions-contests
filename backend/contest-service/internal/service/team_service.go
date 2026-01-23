package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sports-prediction-contests/contest-service/internal/models"
	"github.com/sports-prediction-contests/contest-service/internal/repository"
	"github.com/sports-prediction-contests/shared/auth"
	"github.com/sports-prediction-contests/shared/proto/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TeamResponse mirrors proto response structure
type TeamResponse struct {
	Response *common.Response
	Team     *TeamProto
}

type TeamProto struct {
	ID             uint32
	Name           string
	Description    string
	InviteCode     string
	CaptainID      uint32
	MaxMembers     uint32
	CurrentMembers uint32
	IsActive       bool
	CreatedAt      *timestamppb.Timestamp
	UpdatedAt      *timestamppb.Timestamp
}

type TeamMemberProto struct {
	ID       uint32
	TeamID   uint32
	UserID   uint32
	UserName string
	Role     string
	Status   string
	JoinedAt *timestamppb.Timestamp
}

type TeamService struct {
	teamRepo        repository.TeamRepositoryInterface
	memberRepo      repository.TeamMemberRepositoryInterface
	contestEntryRepo repository.TeamContestEntryRepositoryInterface
}

func NewTeamService(
	teamRepo repository.TeamRepositoryInterface,
	memberRepo repository.TeamMemberRepositoryInterface,
	contestEntryRepo repository.TeamContestEntryRepositoryInterface,
) *TeamService {
	return &TeamService{
		teamRepo:        teamRepo,
		memberRepo:      memberRepo,
		contestEntryRepo: contestEntryRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name, description string, maxMembers uint) (*TeamResponse, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &TeamResponse{Response: errorResponse("Authentication required", common.ErrorCode_UNAUTHENTICATED)}, nil
	}

	if maxMembers == 0 {
		maxMembers = 10
	}

	team := &models.Team{
		Name:        name,
		Description: description,
		CaptainID:   userID,
		MaxMembers:  uint(maxMembers),
		IsActive:    true,
	}

	member := &models.TeamMember{
		UserID:   userID,
		Role:     "captain",
		Status:   "active",
		JoinedAt: time.Now(),
	}

	if err := s.teamRepo.CreateWithMember(team, member); err != nil {
		log.Printf("[ERROR] Failed to create team: %v", err)
		return &TeamResponse{Response: errorResponse(err.Error(), common.ErrorCode_INVALID_ARGUMENT)}, nil
	}

	log.Printf("[INFO] Team created: %d by user %d", team.ID, userID)

	return &TeamResponse{
		Response: successResponse("Team created successfully"),
		Team:     s.modelToProto(team),
	}, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, id uint, name, description string, maxMembers uint) (*TeamResponse, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return &TeamResponse{Response: errorResponse("Authentication required", common.ErrorCode_UNAUTHENTICATED)}, nil
	}

	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return &TeamResponse{Response: errorResponse(err.Error(), common.ErrorCode_NOT_FOUND)}, nil
	}

	if !team.IsCaptain(userID) {
		return &TeamResponse{Response: errorResponse("Only captain can update team", common.ErrorCode_PERMISSION_DENIED)}, nil
	}

	team.Name = name
	team.Description = description
	if maxMembers > 0 {
		team.MaxMembers = uint(maxMembers)
	}

	if err := s.teamRepo.Update(team); err != nil {
		return &TeamResponse{Response: errorResponse(err.Error(), common.ErrorCode_INVALID_ARGUMENT)}, nil
	}

	return &TeamResponse{
		Response: successResponse("Team updated successfully"),
		Team:     s.modelToProto(team),
	}, nil
}

func (s *TeamService) GetTeam(ctx context.Context, id uint) (*TeamResponse, error) {
	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return &TeamResponse{Response: errorResponse(err.Error(), common.ErrorCode_NOT_FOUND)}, nil
	}
	return &TeamResponse{
		Response: successResponse("Team retrieved"),
		Team:     s.modelToProto(team),
	}, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, id uint) (*common.Response, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return errorResponse("Authentication required", common.ErrorCode_UNAUTHENTICATED), nil
	}

	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		return errorResponse(err.Error(), common.ErrorCode_NOT_FOUND), nil
	}

	if !team.IsCaptain(userID) {
		return errorResponse("Only captain can delete team", common.ErrorCode_PERMISSION_DENIED), nil
	}

	if err := s.teamRepo.Delete(id); err != nil {
		return errorResponse(err.Error(), common.ErrorCode_INTERNAL_ERROR), nil
	}

	log.Printf("[INFO] Team deleted: %d", id)
	return successResponse("Team deleted successfully"), nil
}

func (s *TeamService) ListTeams(ctx context.Context, page, limit int, myTeamsOnly bool) ([]*TeamProto, int64, error) {
	userID, _ := auth.GetUserIDFromContext(ctx)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	teams, total, err := s.teamRepo.List(limit, offset, userID, myTeamsOnly)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*TeamProto, len(teams))
	for i, t := range teams {
		result[i] = s.modelToProto(t)
	}
	return result, total, nil
}

func (s *TeamService) JoinTeam(ctx context.Context, inviteCode string) (*TeamMemberProto, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("authentication required")
	}

	team, err := s.teamRepo.GetByInviteCode(inviteCode)
	if err != nil {
		return nil, err
	}

	// Check if already a member
	if _, err := s.memberRepo.GetByTeamAndUser(team.ID, userID); err == nil {
		return nil, &joinError{"Already a member of this team"}
	}

	if !team.CanJoin() {
		return nil, &joinError{"Team is full or inactive"}
	}

	member := &models.TeamMember{
		TeamID:   team.ID,
		UserID:   userID,
		Role:     "member",
		Status:   "active",
		JoinedAt: time.Now(),
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, fmt.Errorf("failed to join team: %w", err)
	}

	s.updateMemberCount(team.ID)
	log.Printf("[INFO] User %d joined team %d", userID, team.ID)

	return s.memberToProto(member), nil
}

func (s *TeamService) LeaveTeam(ctx context.Context, teamID uint) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("authentication required")
	}

	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return err
	}

	if team.IsCaptain(userID) {
		return &joinError{"Captain cannot leave team"}
	}

	if err := s.memberRepo.DeleteByTeamAndUser(teamID, userID); err != nil {
		return err
	}

	s.updateMemberCount(teamID)
	log.Printf("[INFO] User %d left team %d", userID, teamID)
	return nil
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, targetUserID uint) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("authentication required")
	}

	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return err
	}

	if !team.IsCaptain(userID) {
		return &joinError{"Only captain can remove members"}
	}

	if team.CaptainID == targetUserID {
		return &joinError{"Cannot remove captain"}
	}

	if err := s.memberRepo.DeleteByTeamAndUser(teamID, targetUserID); err != nil {
		return err
	}

	s.updateMemberCount(teamID)
	log.Printf("[INFO] User %d removed from team %d by captain %d", targetUserID, teamID, userID)
	return nil
}

func (s *TeamService) ListMembers(ctx context.Context, teamID uint, page, limit int) ([]*TeamMemberProto, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	members, total, err := s.memberRepo.ListByTeam(teamID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*TeamMemberProto, len(members))
	for i, m := range members {
		result[i] = s.memberToProto(m)
	}
	return result, total, nil
}

func (s *TeamService) RegenerateInviteCode(ctx context.Context, teamID uint) (string, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("authentication required")
	}

	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return "", err
	}

	if !team.IsCaptain(userID) {
		return "", &joinError{"Only captain can regenerate invite code"}
	}

	newCode, err := models.GenerateInviteCode()
	if err != nil {
		return "", err
	}

	team.InviteCode = newCode
	if err := s.teamRepo.Update(team); err != nil {
		return "", err
	}

	return newCode, nil
}

func (s *TeamService) JoinContestAsTeam(ctx context.Context, teamID, contestID uint) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("authentication required")
	}

	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return err
	}

	if !team.IsCaptain(userID) {
		return &joinError{"Only captain can join contests"}
	}

	entry := &models.TeamContestEntry{
		TeamID:    teamID,
		ContestID: contestID,
		JoinedAt:  time.Now(),
	}

	return s.contestEntryRepo.Create(entry)
}

func (s *TeamService) LeaveContestAsTeam(ctx context.Context, teamID, contestID uint) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("authentication required")
	}

	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		return err
	}

	if !team.IsCaptain(userID) {
		return &joinError{"Only captain can leave contests"}
	}

	return s.contestEntryRepo.Delete(teamID, contestID)
}

func (s *TeamService) GetTeamLeaderboard(ctx context.Context, contestID uint, limit int) ([]*TeamLeaderboardEntryProto, error) {
	if limit <= 0 {
		limit = 10
	}

	entries, err := s.contestEntryRepo.ListByContest(contestID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*TeamLeaderboardEntryProto, len(entries))
	for i, e := range entries {
		result[i] = &TeamLeaderboardEntryProto{
			TeamID:      uint32(e.TeamID),
			TeamName:    e.Team.Name,
			TotalPoints: e.TotalPoints,
			Rank:        uint32(e.Rank),
			MemberCount: uint32(e.Team.CurrentMembers),
		}
	}
	return result, nil
}

func (s *TeamService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return successResponse("Team service is healthy"), nil
}

// Helper methods

func (s *TeamService) updateMemberCount(teamID uint) {
	count, err := s.memberRepo.CountByTeam(teamID)
	if err != nil {
		log.Printf("[ERROR] Failed to count members: %v", err)
		return
	}
	team, err := s.teamRepo.GetByID(teamID)
	if err != nil {
		log.Printf("[ERROR] Failed to get team for member count update: %v", err)
		return
	}
	team.CurrentMembers = uint(count)
	if err := s.teamRepo.Update(team); err != nil {
		log.Printf("[ERROR] Failed to update member count for team %d: %v", teamID, err)
	}
}

func (s *TeamService) modelToProto(t *models.Team) *TeamProto {
	return &TeamProto{
		ID:             uint32(t.ID),
		Name:           t.Name,
		Description:    t.Description,
		InviteCode:     t.InviteCode,
		CaptainID:      uint32(t.CaptainID),
		MaxMembers:     uint32(t.MaxMembers),
		CurrentMembers: uint32(t.CurrentMembers),
		IsActive:       t.IsActive,
		CreatedAt:      timestamppb.New(t.CreatedAt),
		UpdatedAt:      timestamppb.New(t.UpdatedAt),
	}
}

func (s *TeamService) memberToProto(m *models.TeamMember) *TeamMemberProto {
	return &TeamMemberProto{
		ID:       uint32(m.ID),
		TeamID:   uint32(m.TeamID),
		UserID:   uint32(m.UserID),
		Role:     m.Role,
		Status:   m.Status,
		JoinedAt: timestamppb.New(m.JoinedAt),
	}
}

type TeamLeaderboardEntryProto struct {
	TeamID      uint32
	TeamName    string
	TotalPoints float64
	Rank        uint32
	MemberCount uint32
}

type joinError struct{ msg string }

func (e *joinError) Error() string { return e.msg }

func successResponse(msg string) *common.Response {
	return &common.Response{Success: true, Message: msg, Code: 0, Timestamp: timestamppb.Now()}
}

func errorResponse(msg string, code common.ErrorCode) *common.Response {
	return &common.Response{Success: false, Message: msg, Code: int32(code), Timestamp: timestamppb.Now()}
}

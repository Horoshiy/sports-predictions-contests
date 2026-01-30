package service

import (
	"context"

	pb "github.com/sports-prediction-contests/shared/proto/team"
	"github.com/sports-prediction-contests/shared/proto/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TeamServiceGRPC wraps TeamService to implement gRPC interface
type TeamServiceGRPC struct {
	pb.UnimplementedTeamServiceServer
	teamService *TeamService
}

// NewTeamServiceGRPC creates a new gRPC wrapper for TeamService
func NewTeamServiceGRPC(teamService *TeamService) *TeamServiceGRPC {
	return &TeamServiceGRPC{teamService: teamService}
}

func (s *TeamServiceGRPC) CreateTeam(ctx context.Context, req *pb.CreateTeamRequest) (*pb.CreateTeamResponse, error) {
	resp, err := s.teamService.CreateTeam(ctx, req.Name, req.Description, uint(req.MaxMembers))
	if err != nil {
		return &pb.CreateTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	// Defensive nil check
	if resp == nil || resp.Team == nil {
		return &pb.CreateTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: "Internal error: nil response",
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.CreateTeamResponse{
		Response: resp.Response,
		Team: &pb.Team{
			Id: resp.Team.ID,
			Name: resp.Team.Name,
			Description: resp.Team.Description,
			InviteCode: resp.Team.InviteCode,
			CaptainId: resp.Team.CaptainID,
			MaxMembers: resp.Team.MaxMembers,
			CurrentMembers: resp.Team.CurrentMembers,
			IsActive: resp.Team.IsActive,
			CreatedAt: resp.Team.CreatedAt,
			UpdatedAt: resp.Team.UpdatedAt,
		},
	}, nil
}

func (s *TeamServiceGRPC) UpdateTeam(ctx context.Context, req *pb.UpdateTeamRequest) (*pb.UpdateTeamResponse, error) {
	resp, err := s.teamService.UpdateTeam(ctx, uint(req.Id), req.Name, req.Description, uint(req.MaxMembers))
	if err != nil {
		return &pb.UpdateTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	// Defensive nil check
	if resp == nil || resp.Team == nil {
		return &pb.UpdateTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: "Internal error: nil response",
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.UpdateTeamResponse{
		Response: resp.Response,
		Team: &pb.Team{
			Id: resp.Team.ID,
			Name: resp.Team.Name,
			Description: resp.Team.Description,
			InviteCode: resp.Team.InviteCode,
			CaptainId: resp.Team.CaptainID,
			MaxMembers: resp.Team.MaxMembers,
			CurrentMembers: resp.Team.CurrentMembers,
			IsActive: resp.Team.IsActive,
			CreatedAt: resp.Team.CreatedAt,
			UpdatedAt: resp.Team.UpdatedAt,
		},
	}, nil
}

func (s *TeamServiceGRPC) GetTeam(ctx context.Context, req *pb.GetTeamRequest) (*pb.GetTeamResponse, error) {
	resp, err := s.teamService.GetTeam(ctx, uint(req.Id))
	if err != nil {
		return &pb.GetTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	// Defensive nil check
	if resp == nil || resp.Team == nil {
		return &pb.GetTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: "Internal error: nil response",
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.GetTeamResponse{
		Response: resp.Response,
		Team: &pb.Team{
			Id: resp.Team.ID,
			Name: resp.Team.Name,
			Description: resp.Team.Description,
			InviteCode: resp.Team.InviteCode,
			CaptainId: resp.Team.CaptainID,
			MaxMembers: resp.Team.MaxMembers,
			CurrentMembers: resp.Team.CurrentMembers,
			IsActive: resp.Team.IsActive,
			CreatedAt: resp.Team.CreatedAt,
			UpdatedAt: resp.Team.UpdatedAt,
		},
	}, nil
}

func (s *TeamServiceGRPC) DeleteTeam(ctx context.Context, req *pb.DeleteTeamRequest) (*pb.DeleteTeamResponse, error) {
	resp, err := s.teamService.DeleteTeam(ctx, uint(req.Id))
	if err != nil {
		return &pb.DeleteTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.DeleteTeamResponse{Response: resp}, nil
}

func (s *TeamServiceGRPC) ListTeams(ctx context.Context, req *pb.ListTeamsRequest) (*pb.ListTeamsResponse, error) {
	page := 1
	limit := 20
	
	// Handle nil pagination
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		limit = int(req.Pagination.Limit)
	}
	
	// Defensive validation
	if limit <= 0 {
		limit = 20 // default value
	}
	
	teams, total, err := s.teamService.ListTeams(ctx, page, limit, req.MyTeamsOnly)
	if err != nil {
		return &pb.ListTeamsResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	pbTeams := make([]*pb.Team, len(teams))
	for i, t := range teams {
		pbTeams[i] = &pb.Team{
			Id: t.ID,
			Name: t.Name,
			Description: t.Description,
			InviteCode: t.InviteCode,
			CaptainId: t.CaptainID,
			MaxMembers: t.MaxMembers,
			CurrentMembers: t.CurrentMembers,
			IsActive: t.IsActive,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		}
	}
	
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	
	return &pb.ListTeamsResponse{
		Response: &common.Response{
			Success: true,
			Message: "Teams retrieved",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
		Teams: pbTeams,
		Pagination: &common.PaginationResponse{
			Page: int32(page),
			Limit: int32(limit),
			Total: int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *TeamServiceGRPC) JoinTeam(ctx context.Context, req *pb.JoinTeamRequest) (*pb.JoinTeamResponse, error) {
	member, err := s.teamService.JoinTeam(ctx, req.InviteCode)
	if err != nil {
		return &pb.JoinTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.JoinTeamResponse{
		Response: &common.Response{
			Success: true,
			Message: "Joined team successfully",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
		Member: &pb.TeamMember{
			Id: member.ID,
			TeamId: member.TeamID,
			UserId: member.UserID,
			UserName: member.UserName,
			Role: member.Role,
			Status: member.Status,
			JoinedAt: member.JoinedAt,
		},
	}, nil
}

func (s *TeamServiceGRPC) LeaveTeam(ctx context.Context, req *pb.LeaveTeamRequest) (*pb.LeaveTeamResponse, error) {
	err := s.teamService.LeaveTeam(ctx, uint(req.TeamId))
	if err != nil {
		return &pb.LeaveTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.LeaveTeamResponse{
		Response: &common.Response{
			Success: true,
			Message: "Left team successfully",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

func (s *TeamServiceGRPC) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	err := s.teamService.RemoveMember(ctx, uint(req.TeamId), uint(req.UserId))
	if err != nil {
		return &pb.RemoveMemberResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.RemoveMemberResponse{
		Response: &common.Response{
			Success: true,
			Message: "Member removed successfully",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

func (s *TeamServiceGRPC) ListMembers(ctx context.Context, req *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	page := int(req.Pagination.Page)
	limit := int(req.Pagination.Limit)
	
	// Defensive validation
	if limit <= 0 {
		limit = 20 // default value
	}
	
	members, total, err := s.teamService.ListMembers(ctx, uint(req.TeamId), page, limit)
	if err != nil {
		return &pb.ListMembersResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	pbMembers := make([]*pb.TeamMember, len(members))
	for i, m := range members {
		pbMembers[i] = &pb.TeamMember{
			Id: m.ID,
			TeamId: m.TeamID,
			UserId: m.UserID,
			UserName: m.UserName,
			Role: m.Role,
			Status: m.Status,
			JoinedAt: m.JoinedAt,
		}
	}
	
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	
	return &pb.ListMembersResponse{
		Response: &common.Response{
			Success: true,
			Message: "Members retrieved",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
		Members: pbMembers,
		Pagination: &common.PaginationResponse{
			Page: int32(page),
			Limit: int32(limit),
			Total: int32(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *TeamServiceGRPC) RegenerateInviteCode(ctx context.Context, req *pb.RegenerateInviteCodeRequest) (*pb.RegenerateInviteCodeResponse, error) {
	code, err := s.teamService.RegenerateInviteCode(ctx, uint(req.TeamId))
	if err != nil {
		return &pb.RegenerateInviteCodeResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.RegenerateInviteCodeResponse{
		Response: &common.Response{
			Success: true,
			Message: "Invite code regenerated",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
		InviteCode: code,
	}, nil
}

func (s *TeamServiceGRPC) JoinContestAsTeam(ctx context.Context, req *pb.JoinContestAsTeamRequest) (*pb.JoinContestAsTeamResponse, error) {
	err := s.teamService.JoinContestAsTeam(ctx, uint(req.TeamId), uint(req.ContestId))
	if err != nil {
		return &pb.JoinContestAsTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.JoinContestAsTeamResponse{
		Response: &common.Response{
			Success: true,
			Message: "Team joined contest successfully",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

func (s *TeamServiceGRPC) LeaveContestAsTeam(ctx context.Context, req *pb.LeaveContestAsTeamRequest) (*pb.LeaveContestAsTeamResponse, error) {
	err := s.teamService.LeaveContestAsTeam(ctx, uint(req.TeamId), uint(req.ContestId))
	if err != nil {
		return &pb.LeaveContestAsTeamResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	return &pb.LeaveContestAsTeamResponse{
		Response: &common.Response{
			Success: true,
			Message: "Team left contest successfully",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

func (s *TeamServiceGRPC) GetTeamLeaderboard(ctx context.Context, req *pb.GetTeamLeaderboardRequest) (*pb.GetTeamLeaderboardResponse, error) {
	entries, err := s.teamService.GetTeamLeaderboard(ctx, uint(req.ContestId), int(req.Limit))
	if err != nil {
		return &pb.GetTeamLeaderboardResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code: int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}
	
	pbEntries := make([]*pb.TeamLeaderboardEntry, len(entries))
	for i, e := range entries {
		pbEntries[i] = &pb.TeamLeaderboardEntry{
			TeamId: e.TeamID,
			TeamName: e.TeamName,
			TotalPoints: e.TotalPoints,
			Rank: e.Rank,
			MemberCount: e.MemberCount,
		}
	}
	
	return &pb.GetTeamLeaderboardResponse{
		Response: &common.Response{
			Success: true,
			Message: "Leaderboard retrieved",
			Code: 0,
			Timestamp: timestamppb.Now(),
		},
		Entries: pbEntries,
	}, nil
}

func (s *TeamServiceGRPC) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return s.teamService.Check(ctx, req)
}

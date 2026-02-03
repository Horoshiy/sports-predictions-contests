package clients

import (
	"context"
	"fmt"

	teampb "github.com/sports-prediction-contests/shared/proto/team"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TeamClient wraps the gRPC team service client
type TeamClient struct {
	client teampb.TeamServiceClient
	conn   *grpc.ClientConn
}

// NewTeamClient creates a new team service client
func NewTeamClient(endpoint string) (*TeamClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to team service: %w", err)
	}

	client := teampb.NewTeamServiceClient(conn)
	return &TeamClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *TeamClient) Close() error {
	return c.conn.Close()
}

// GetTeam retrieves a team by ID
func (c *TeamClient) GetTeam(ctx context.Context, teamID uint32) (*teampb.Team, error) {
	req := &teampb.GetTeamRequest{
		Id: teamID,
	}

	resp, err := c.client.GetTeam(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	if !resp.Response.Success {
		return nil, fmt.Errorf("team service error: %s", resp.Response.Message)
	}

	return resp.Team, nil
}

// IsTeamCaptain checks if a user is the captain of a team
func (c *TeamClient) IsTeamCaptain(ctx context.Context, teamID uint32, userID uint64) (bool, error) {
	team, err := c.GetTeam(ctx, teamID)
	if err != nil {
		return false, err
	}

	return uint64(team.CaptainId) == userID, nil
}

// GetTeamMembers retrieves all members of a team
func (c *TeamClient) GetTeamMembers(ctx context.Context, teamID uint32) ([]*teampb.TeamMember, error) {
	req := &teampb.ListMembersRequest{
		TeamId: teamID,
	}

	resp, err := c.client.ListMembers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}

	if !resp.Response.Success {
		return nil, fmt.Errorf("team service error: %s", resp.Response.Message)
	}

	return resp.Members, nil
}

// IsUserInTeam checks if a user is a member of a team
func (c *TeamClient) IsUserInTeam(ctx context.Context, teamID uint32, userID uint64) (bool, error) {
	members, err := c.GetTeamMembers(ctx, teamID)
	if err != nil {
		return false, err
	}

	for _, member := range members {
		if uint64(member.UserId) == userID && member.Status == "active" {
			return true, nil
		}
	}

	return false, nil
}

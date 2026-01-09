package clients

import (
	"context"
	"fmt"

	"github.com/sports-prediction-contests/shared/proto/common"
	contestpb "github.com/sports-prediction-contests/shared/proto/contest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ContestClient wraps the gRPC contest service client
type ContestClient struct {
	client contestpb.ContestServiceClient
	conn   *grpc.ClientConn
}

// NewContestClient creates a new contest service client
func NewContestClient(endpoint string) (*ContestClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to contest service: %w", err)
	}

	client := contestpb.NewContestServiceClient(conn)
	return &ContestClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *ContestClient) Close() error {
	return c.conn.Close()
}

// GetContest retrieves a contest by ID
func (c *ContestClient) GetContest(ctx context.Context, contestID uint32) (*contestpb.Contest, error) {
	req := &contestpb.GetContestRequest{
		Id: contestID,
	}

	resp, err := c.client.GetContest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest: %w", err)
	}

	if !resp.Response.Success {
		return nil, fmt.Errorf("contest service error: %s", resp.Response.Message)
	}

	return resp.Contest, nil
}

// ValidateContestParticipation checks if a user is a participant in a contest
func (c *ContestClient) ValidateContestParticipation(ctx context.Context, contestID uint32, userID uint32) error {
	// Get contest details to validate it exists and is active
	contest, err := c.GetContest(ctx, contestID)
	if err != nil {
		return fmt.Errorf("contest validation failed")
	}

	// Check if contest is active
	if contest.Status != "active" {
		return fmt.Errorf("contest is not active")
	}

	// Check if user is a participant
	isParticipant, err := c.IsUserParticipant(ctx, contestID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify participation")
	}

	if !isParticipant {
		return fmt.Errorf("user is not a participant in this contest")
	}
	
	return nil
}

// IsUserParticipant checks if a user is a participant in a contest
func (c *ContestClient) IsUserParticipant(ctx context.Context, contestID uint32, userID uint32) (bool, error) {
	req := &contestpb.ListParticipantsRequest{
		ContestId: contestID,
		Pagination: &common.PaginationRequest{
			Page:     1,
			PageSize: 100, // Check first 100 participants
		},
	}

	resp, err := c.client.ListParticipants(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to list participants: %w", err)
	}

	if !resp.Response.Success {
		return false, fmt.Errorf("contest service error: %s", resp.Response.Message)
	}

	// Check if user is in the participants list
	for _, participant := range resp.Participants {
		if participant.UserId == userID && participant.Status == "active" {
			return true, nil
		}
	}

	return false, nil
}

// IsContestActive checks if a contest is currently active
func (c *ContestClient) IsContestActive(ctx context.Context, contestID uint32) (bool, error) {
	contest, err := c.GetContest(ctx, contestID)
	if err != nil {
		return false, err
	}

	return contest.Status == "active", nil
}

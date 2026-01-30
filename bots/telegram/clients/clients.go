package clients

import (
	"fmt"

	contestpb "github.com/sports-prediction-contests/shared/proto/contest"
	notificationpb "github.com/sports-prediction-contests/shared/proto/notification"
	predictionpb "github.com/sports-prediction-contests/shared/proto/prediction"
	scoringpb "github.com/sports-prediction-contests/shared/proto/scoring"
	userpb "github.com/sports-prediction-contests/shared/proto/user"
	"github.com/sports-prediction-contests/telegram-bot/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	User         userpb.UserServiceClient
	Contest      contestpb.ContestServiceClient
	Scoring      scoringpb.ScoringServiceClient
	Notification notificationpb.NotificationServiceClient
	Prediction   predictionpb.PredictionServiceClient
	conns        []*grpc.ClientConn
}

func New(cfg *config.Config) (*Clients, error) {
	c := &Clients{}
	var err error

	defer func() {
		if err != nil {
			c.Close() // Cleanup on failure
		}
	}()

	// User service
	userConn, err := grpc.Dial(cfg.UserServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	c.conns = append(c.conns, userConn)
	c.User = userpb.NewUserServiceClient(userConn)

	// Contest service
	contestConn, err := grpc.Dial(cfg.ContestServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to contest service: %w", err)
	}
	c.conns = append(c.conns, contestConn)
	c.Contest = contestpb.NewContestServiceClient(contestConn)

	// Scoring service
	scoringConn, err := grpc.Dial(cfg.ScoringServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to scoring service: %w", err)
	}
	c.conns = append(c.conns, scoringConn)
	c.Scoring = scoringpb.NewScoringServiceClient(scoringConn)

	// Notification service
	notificationConn, err := grpc.Dial(cfg.NotificationServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}
	c.conns = append(c.conns, notificationConn)
	c.Notification = notificationpb.NewNotificationServiceClient(notificationConn)

	// Prediction service
	predictionConn, err := grpc.Dial(cfg.PredictionServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to prediction service: %w", err)
	}
	c.conns = append(c.conns, predictionConn)
	c.Prediction = predictionpb.NewPredictionServiceClient(predictionConn)

	return c, nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		if conn != nil {
			conn.Close()
		}
	}
}

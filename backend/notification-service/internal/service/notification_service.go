package service

import (
	"context"
	"log"

	"github.com/sports-prediction-contests/notification-service/internal/models"
	"github.com/sports-prediction-contests/notification-service/internal/repository"
	"github.com/sports-prediction-contests/notification-service/internal/worker"
	"github.com/sports-prediction-contests/shared/proto/common"
	pb "github.com/sports-prediction-contests/shared/proto/notification"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	repo   repository.NotificationRepositoryInterface
	worker *worker.WorkerPool
}

func NewNotificationService(repo repository.NotificationRepositoryInterface, worker *worker.WorkerPool) *NotificationService {
	return &NotificationService{repo: repo, worker: worker}
}

func (s *NotificationService) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	channels := req.Channels
	if len(channels) == 0 {
		channels = []pb.NotificationChannel{pb.NotificationChannel_IN_APP}
	}

	var lastNotification *models.Notification
	for _, ch := range channels {
		notification := &models.Notification{
			UserID:  uint(req.UserId),
			Type:    req.Type.String(),
			Title:   req.Title,
			Message: req.Message,
			Data:    req.Data,
			Channel: channelToString(ch),
		}

		if err := s.repo.Create(ctx, notification); err != nil {
			return &pb.SendNotificationResponse{
				Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
			}, nil
		}

		lastNotification = notification

		if ch != pb.NotificationChannel_IN_APP {
			pref, err := s.repo.GetPreference(ctx, uint(req.UserId), channelToString(ch))
			if err != nil {
				log.Printf("Warning: failed to get preference for user %d channel %s: %v", req.UserId, channelToString(ch), err)
			}
			s.worker.Submit(worker.NotificationJob{Notification: notification, Preference: pref})
		}
	}

	return &pb.SendNotificationResponse{
		Response:     &common.Response{Success: true, Message: "Notification sent", Timestamp: timestamppb.Now()},
		Notification: toProtoNotification(lastNotification),
	}, nil
}

func (s *NotificationService) GetNotification(ctx context.Context, req *pb.GetNotificationRequest) (*pb.GetNotificationResponse, error) {
	notification, err := s.repo.GetByID(ctx, uint(req.Id))
	if err != nil {
		return &pb.GetNotificationResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
		}, nil
	}

	return &pb.GetNotificationResponse{
		Response:     &common.Response{Success: true, Timestamp: timestamppb.Now()},
		Notification: toProtoNotification(notification),
	}, nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	limit, page := 20, 1
	if req.Pagination != nil {
		if req.Pagination.Limit > 0 {
			limit = int(req.Pagination.Limit)
		}
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
	}
	offset := (page - 1) * limit

	notifications, total, err := s.repo.GetByUser(ctx, uint(req.UserId), req.UnreadOnly, limit, offset)
	if err != nil {
		return &pb.GetNotificationsResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
		}, nil
	}

	protoNotifications := make([]*pb.Notification, len(notifications))
	for i, n := range notifications {
		protoNotifications[i] = toProtoNotification(n)
	}

	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	return &pb.GetNotificationsResponse{
		Response:      &common.Response{Success: true, Timestamp: timestamppb.Now()},
		Notifications: protoNotifications,
		Pagination:    &common.PaginationResponse{Total: int32(total), TotalPages: totalPages},
	}, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	if err := s.repo.MarkAsRead(ctx, uint(req.Id)); err != nil {
		return &pb.MarkAsReadResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.MarkAsReadResponse{
		Response: &common.Response{Success: true, Message: "Marked as read", Timestamp: timestamppb.Now()},
	}, nil
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, req *pb.MarkAllAsReadRequest) (*pb.MarkAllAsReadResponse, error) {
	count, err := s.repo.MarkAllAsRead(ctx, uint(req.UserId))
	if err != nil {
		return &pb.MarkAllAsReadResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.MarkAllAsReadResponse{
		Response: &common.Response{Success: true, Timestamp: timestamppb.Now()},
		Count:    uint32(count),
	}, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	count, err := s.repo.GetUnreadCount(ctx, uint(req.UserId))
	if err != nil {
		return &pb.GetUnreadCountResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
		}, nil
	}
	return &pb.GetUnreadCountResponse{
		Response: &common.Response{Success: true, Timestamp: timestamppb.Now()},
		Count:    uint32(count),
	}, nil
}

func (s *NotificationService) GetPreferences(ctx context.Context, req *pb.GetPreferencesRequest) (*pb.GetPreferencesResponse, error) {
	prefs, err := s.repo.GetPreferences(ctx, uint(req.UserId))
	if err != nil {
		return &pb.GetPreferencesResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
		}, nil
	}

	protoPrefs := make([]*pb.NotificationPreference, len(prefs))
	for i, p := range prefs {
		protoPrefs[i] = toProtoPreference(p)
	}

	return &pb.GetPreferencesResponse{
		Response:    &common.Response{Success: true, Timestamp: timestamppb.Now()},
		Preferences: protoPrefs,
	}, nil
}

func (s *NotificationService) UpdatePreference(ctx context.Context, req *pb.UpdatePreferenceRequest) (*pb.UpdatePreferenceResponse, error) {
	pref, err := s.repo.GetPreference(ctx, uint(req.UserId), channelToString(req.Channel))
	if err != nil {
		log.Printf("Warning: failed to get existing preference: %v", err)
	}
	if pref == nil {
		pref = &models.NotificationPreference{
			UserID:  uint(req.UserId),
			Channel: channelToString(req.Channel),
		}
	}

	pref.Enabled = req.Enabled
	pref.TelegramChatID = req.TelegramChatId
	pref.Email = req.Email

	if err := s.repo.UpdatePreference(ctx, pref); err != nil {
		return &pb.UpdatePreferenceResponse{
			Response: &common.Response{Success: false, Message: err.Error(), Timestamp: timestamppb.Now()},
		}, nil
	}

	return &pb.UpdatePreferenceResponse{
		Response:   &common.Response{Success: true, Timestamp: timestamppb.Now()},
		Preference: toProtoPreference(pref),
	}, nil
}

func (s *NotificationService) Check(ctx context.Context, req *emptypb.Empty) (*common.Response, error) {
	return &common.Response{Success: true, Message: "Notification service is healthy", Timestamp: timestamppb.Now()}, nil
}

func channelToString(ch pb.NotificationChannel) string {
	switch ch {
	case pb.NotificationChannel_TELEGRAM:
		return "telegram"
	case pb.NotificationChannel_EMAIL:
		return "email"
	default:
		return "in_app"
	}
}

func toProtoNotification(n *models.Notification) *pb.Notification {
	pn := &pb.Notification{
		Id:        uint32(n.ID),
		UserId:    uint32(n.UserID),
		Type:      pb.NotificationType(pb.NotificationType_value[n.Type]),
		Title:     n.Title,
		Message:   n.Message,
		Data:      n.Data,
		Channel:   pb.NotificationChannel(pb.NotificationChannel_value[n.Channel]),
		IsRead:    n.IsRead,
		CreatedAt: timestamppb.New(n.CreatedAt),
		UpdatedAt: timestamppb.New(n.UpdatedAt),
	}
	if n.SentAt != nil {
		pn.SentAt = timestamppb.New(*n.SentAt)
	}
	if n.ReadAt != nil {
		pn.ReadAt = timestamppb.New(*n.ReadAt)
	}
	return pn
}

func toProtoPreference(p *models.NotificationPreference) *pb.NotificationPreference {
	return &pb.NotificationPreference{
		Id:             uint32(p.ID),
		UserId:         uint32(p.UserID),
		Channel:        pb.NotificationChannel(pb.NotificationChannel_value[p.Channel]),
		Enabled:        p.Enabled,
		TelegramChatId: p.TelegramChatID,
		Email:          p.Email,
		CreatedAt:      timestamppb.New(p.CreatedAt),
		UpdatedAt:      timestamppb.New(p.UpdatedAt),
	}
}

package service

import (
	"context"

	"github.com/sports-prediction-contests/shared/auth"
	pb "github.com/sports-prediction-contests/shared/proto/user"
	"github.com/sports-prediction-contests/shared/proto/common"
	"github.com/sports-prediction-contests/user-service/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserService implements the gRPC UserService
type UserService struct {
	pb.UnimplementedUserServiceServer
	authService AuthServiceInterface
	userRepo    repository.UserRepositoryInterface
}

// NewUserService creates a new UserService instance
func NewUserService(authService AuthServiceInterface, userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		authService: authService,
		userRepo:    userRepo,
	}
}

// Register handles user registration
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, token, err := s.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		return &pb.RegisterResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code:    int32(common.ErrorCode_INVALID_ARGUMENT),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.RegisterResponse{
		Response: &common.Response{
			Success: true,
			Message: "User registered successfully",
			Code:    0,
			Timestamp: timestamppb.Now(),
		},
		User: &pb.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Token: token,
	}, nil
}

// Login handles user authentication
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, token, err := s.authService.Login(req.Email, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code:    int32(common.ErrorCode_UNAUTHENTICATED),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.LoginResponse{
		Response: &common.Response{
			Success: true,
			Message: "Login successful",
			Code:    0,
			Timestamp: timestamppb.Now(),
		},
		User: &pb.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Token: token,
	}, nil
}

// GetProfile retrieves user profile (requires authentication)
func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	// Extract user ID from context (set by JWT interceptor)
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return &pb.GetProfileResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code:    int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.GetProfileResponse{
		Response: &common.Response{
			Success: true,
			Message: "Profile retrieved successfully",
			Code:    0,
			Timestamp: timestamppb.Now(),
		},
		User: &pb.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

// UpdateProfile updates user profile (requires authentication)
func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	// Extract user ID from context (set by JWT interceptor)
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Get existing user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return &pb.UpdateProfileResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code:    int32(common.ErrorCode_NOT_FOUND),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return &pb.UpdateProfileResponse{
			Response: &common.Response{
				Success: false,
				Message: err.Error(),
				Code:    int32(common.ErrorCode_INTERNAL_ERROR),
				Timestamp: timestamppb.Now(),
			},
		}, nil
	}

	return &pb.UpdateProfileResponse{
		Response: &common.Response{
			Success: true,
			Message: "Profile updated successfully",
			Code:    0,
			Timestamp: timestamppb.Now(),
		},
		User: &pb.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

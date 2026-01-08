package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JWTUnaryInterceptor creates a gRPC unary interceptor for JWT authentication
func JWTUnaryInterceptor(secret []byte) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for login/register endpoints
		if strings.Contains(info.FullMethod, "Login") || 
		   strings.Contains(info.FullMethod, "Register") ||
		   strings.Contains(info.FullMethod, "Health") {
			return handler(ctx, req)
		}

		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Get authorization header
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Extract token from Bearer header
		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		// Validate token
		claims, err := ValidateToken(token, secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}

		// Add user info to context
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)

		return handler(ctx, req)
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("user_id").(uint)
	return userID, ok
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("email").(string)
	return email, ok
}

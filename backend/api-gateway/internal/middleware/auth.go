package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sports-prediction-contests/shared/auth"
	"google.golang.org/grpc/metadata"
)

// JWTMiddleware creates HTTP middleware for JWT authentication
func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for login/register endpoints and health checks
			if r.URL.Path == "/health" || strings.HasPrefix(r.URL.Path, "/v1/auth/") {
				next.ServeHTTP(w, r)
				return
			}

			// Extract authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			// Extract token from Bearer header
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			// Validate token
			claims, err := auth.ValidateToken(token, secret)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Add user info to request context for gRPC metadata
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "email", claims.Email)

			// Add authorization header to gRPC metadata
			md := metadata.Pairs("authorization", "Bearer "+token)
			ctx = metadata.NewOutgoingContext(ctx, md)

			// Continue with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

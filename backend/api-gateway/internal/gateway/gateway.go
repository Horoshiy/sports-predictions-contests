package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sports-prediction-contests/api-gateway/internal/config"
	"github.com/sports-prediction-contests/api-gateway/internal/middleware"
	contestpb "github.com/sports-prediction-contests/shared/proto/contest"
	predictionpb "github.com/sports-prediction-contests/shared/proto/prediction"
	scoringpb "github.com/sports-prediction-contests/shared/proto/scoring"
	sportspb "github.com/sports-prediction-contests/shared/proto/sports"
	userpb "github.com/sports-prediction-contests/shared/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Server represents the API Gateway server
type Server struct {
	config *config.Config
	mux    *runtime.ServeMux
}

// NewServer creates a new API Gateway server
func NewServer(cfg *config.Config) (*Server, error) {
	// Create gRPC-Gateway mux with custom error handler
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(customErrorHandler),
	)

	// Register user service
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.UserService, opts)
	if err != nil {
		return nil, err
	}

	// Register contest service
	err = contestpb.RegisterContestServiceHandlerFromEndpoint(ctx, mux, cfg.ContestService, opts)
	if err != nil {
		return nil, err
	}

	// Register prediction service
	err = predictionpb.RegisterPredictionServiceHandlerFromEndpoint(ctx, mux, cfg.PredictionService, opts)
	if err != nil {
		return nil, err
	}

	// Register scoring service
	err = scoringpb.RegisterScoringServiceHandlerFromEndpoint(ctx, mux, cfg.ScoringService, opts)
	if err != nil {
		return nil, err
	}

	// Register sports service
	err = sportspb.RegisterSportsServiceHandlerFromEndpoint(ctx, mux, cfg.SportsService, opts)
	if err != nil {
		return nil, err
	}

	return &Server{
		config: cfg,
		mux:    mux,
	}, nil
}

// Handler returns the HTTP handler with middleware chain
func (s *Server) Handler() http.Handler {
	// Build middleware chain
	handler := http.Handler(s.mux)
	
	// Add health check endpoint
	handler = s.addHealthCheck(handler)
	
	// Add middleware in reverse order (last added = first executed)
	handler = middleware.JWTMiddleware([]byte(s.config.JWTSecret))(handler)
	handler = middleware.CORSMiddleware(s.config.AllowedOrigins)(handler)
	handler = middleware.LoggingMiddleware()(handler)

	return handler
}

// addHealthCheck adds a health check endpoint
func (s *Server) addHealthCheck(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "healthy",
				"service": "api-gateway",
			})
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// customErrorHandler converts gRPC errors to HTTP errors
func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	st := status.Convert(err)
	
	var httpStatus int
	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
	case codes.Internal:
		httpStatus = http.StatusInternalServerError
	case codes.Unavailable:
		httpStatus = http.StatusServiceUnavailable
	default:
		httpStatus = http.StatusInternalServerError
	}

	w.WriteHeader(httpStatus)
	
	errorResp := ErrorResponse{
		Error:   "Request failed",
		Code:    httpStatus,
		Message: st.Message(),
	}
	
	json.NewEncoder(w).Encode(errorResp)
}

package sports

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// RegisterSportsServiceHandlerFromEndpoint is a stub for gRPC-gateway registration
func RegisterSportsServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return nil
}

// UnimplementedSportsServiceServer is a stub for the unimplemented server
type UnimplementedSportsServiceServer struct{}

// RegisterSportsServiceServer is a stub for server registration
func RegisterSportsServiceServer(s *grpc.Server, srv interface{}) {}

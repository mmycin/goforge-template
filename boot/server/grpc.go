package server

import (
	"context"
	"fmt"
	"net"

	"github.com/mmycin/goforge/boot/server/middleware"
	"github.com/mmycin/goforge/core/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCRegistry is a function type for registering gRPC services onto a grpc.Server.
type GRPCRegistry func(server *grpc.Server)

var grpcRegistries []GRPCRegistry

// RegisterGRPC registers a gRPC service registry function to be called when the server starts.
func RegisterGRPC(registry GRPCRegistry) {
	grpcRegistries = append(grpcRegistries, registry)
}

// GetGRPCRegistries returns all globally registered gRPC service registries.
func GetGRPCRegistries() []GRPCRegistry {
	return grpcRegistries
}

// GRPCServer implements the Server interface for gRPC.
type GRPCServer struct {
	server *grpc.Server
	port   int
}

// NewGRPCServer creates a new gRPC server instance.
func NewGRPCServer(cfg *config.AllConfig, registries []GRPCRegistry) *GRPCServer {
	s := grpc.NewServer(middleware.GrpcInterceptor())

	// Register all services (passed-in + globally registered).
	allRegistries := append(registries, grpcRegistries...)
	for _, register := range allRegistries {
		register(s)
	}

	// Enable reflection if configured (useful for tools like grpcurl/Postman).
	if cfg.GRPC.Reflection {
		reflection.Register(s)
	}

	return &GRPCServer{
		server: s,
		port:   cfg.GRPC.Port,
	}
}

// Name returns the server identifier used in log messages.
func (s *GRPCServer) Name() string {
	return "gRPC"
}

// Start begins accepting gRPC connections. Blocks until the server stops.
func (s *GRPCServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("gRPC failed to listen on %s: %w", addr, err)
	}

	log.Info().Str("addr", addr).Msg("Starting gRPC server")
	return s.server.Serve(lis)
}

// Stop shuts the gRPC server down gracefully. If the context expires before
// the graceful drain completes, it forces an immediate stop.
func (s *GRPCServer) Stop(ctx context.Context) error {
	log.Info().Msg("Stopping gRPC server...")

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("gRPC server stopped gracefully")
	case <-ctx.Done():
		log.Warn().Msg("gRPC graceful stop timed out, forcing stop")
		s.server.Stop()
	}
	return nil
}

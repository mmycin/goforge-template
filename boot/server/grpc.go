package server

import (
	"fmt"
	"net"
	"time"

	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/boot/server/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCRegistry is a function type for registering gRPC services
type GRPCRegistry func(server *grpc.Server)

var grpcServices []GRPCRegistry

// RegisterGRPC adds a service registration function to the global registry
func RegisterGRPC(r GRPCRegistry) {
	grpcServices = append(grpcServices, r)
}

// GRPCServer represents the gRPC server
type GRPCServer struct {
	server *grpc.Server
	port   int
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer() *GRPCServer {
	// Apply unified middleware
	s := grpc.NewServer(middleware.GrpcInterceptor())

	// Register all services from the registry
	for _, register := range grpcServices {
		register(s)
	}

	// Enable reflection if configured (useful for tools like grpcurl/Postman)
	if config.GRPC.Reflection {
		reflection.Register(s)
	}

	return &GRPCServer{
		server: s,
		port:   config.GRPC.Port,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Info().Str("addr", addr).Msg("Starting gRPC server")
	return s.server.Serve(lis)
}

// Stop gracefully stops the gRPC server
func (s *GRPCServer) Stop() {
	log.Info().Msg("Stopping gRPC server...")
	s.server.GracefulStop()
}

// StopWithTimeout stops gRPC server gracefully with timeout, or force stops after timeout
func (s *GRPCServer) StopWithTimeout(timeout time.Duration) {
	log.Info().Msg("Stopping gRPC server with timeout...")
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		log.Info().Msg("gRPC server stopped gracefully")
	case <-time.After(timeout):
		log.Warn().Msg("gRPC server graceful stop timed out, forcing stop")
		s.server.Stop()
	}
}

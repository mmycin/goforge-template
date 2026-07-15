package middleware

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/mmycin/goforge/internal/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// GrpcInterceptor combines all unary interceptors
func GrpcInterceptor() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		GrpcRecoveryInterceptor,
		GrpcLoggerInterceptor,
		GrpcAppKeyInterceptor,
		GrpcRateLimitInterceptor,
	)
}

// GrpcAppKeyInterceptor validates X-App-Key from metadata
func GrpcAppKeyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("x-app-key")
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "x-app-key is missing")
	}

	if !ValidateAppKey(values[0]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid x-app-key")
	}

	return handler(ctx, req)
}

// GrpcRateLimitInterceptor limits requests per IP
func GrpcRateLimitInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	limit := config.HTTP.RateLimitPerMinute // Reuse config
	if limit <= 0 {
		return handler(ctx, req)
	}

	var ip string
	if p, ok := peer.FromContext(ctx); ok {
		ip = p.Addr.String()
	} else {
		ip = "unknown"
	}

	limiter := GetLimiter(ip)
	if limiter != nil && !limiter.Allow() {
		return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}

	return handler(ctx, req)
}

// GrpcLoggerInterceptor logs gRPC requests using zerolog
func GrpcLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	latency := time.Since(start)

	code := status.Code(err)
	var ip string
	if p, ok := peer.FromContext(ctx); ok {
		ip = p.Addr.String()
	}

	event := log.Info()
	if err != nil {
		event = log.Error().Err(err)
	}

	event.
		Str("method", info.FullMethod).
		Str("code", code.String()).
		Str("ip", ip).
		Dur("latency", latency).
		Msg("gRPC Request")

	return resp, err
}

// GrpcRecoveryInterceptor recovers from panics in gRPC handlers
func GrpcRecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().
				Interface("panic", r).
				Str("stack", string(debug.Stack())).
				Msg("gRPC Panic Recovered")
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Client defines the interface for our generic gRPC client.
// This allows for mocking in tests and keeps the implementation encapsulated.
type Client interface {
	// Conn returns the underlying grpc.ClientConn.
	// This is used to initialize service-specific proto clients.
	Conn() *grpc.ClientConn

	// Close closes the underlying connection.
	Close() error

	// Target returns the connection target string.
	Target() string
}

// GrpcClient is the concrete implementation of the Client interface.
type GrpcClient struct {
	conn   *grpc.ClientConn
	target string
}

// NewGrpcClient creates and initializes a new GrpcClient.
// It handles the connection logic and provides sensible defaults like insecure credentials
// if no options are provided.
func NewGrpcClient(target string, opts ...grpc.DialOption) (Client, error) {
	if target == "" {
		return nil, fmt.Errorf("grpc client target cannot be empty")
	}

	// Default options if none provided
	if len(opts) == 0 {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Set a default timeout for the dial process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial the target
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc target %s: %w", target, err)
	}

	return &GrpcClient{
		conn:   conn,
		target: target,
	}, nil
}

// Conn returns the *grpc.ClientConn to be used by generated proto clients.
// Example: pb.NewUserServiceClient(client.Conn())
func (c *GrpcClient) Conn() *grpc.ClientConn {
	return c.conn
}

// Close gracefully closes the gRPC connection.
func (c *GrpcClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Target returns the connection target.
func (c *GrpcClient) Target() string {
	return c.target
}

// AppKeyUnaryInterceptor returns a grpc.UnaryClientInterceptor that attaches an 'x-app-key' metadata header.
func AppKeyUnaryInterceptor(appKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-app-key", appKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// AppKeyStreamInterceptor returns a grpc.StreamClientInterceptor that attaches an 'x-app-key' metadata header.
func AppKeyStreamInterceptor(appKey string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-app-key", appKey)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// WithAppKey returns a grpc.DialOption that attaches an 'x-app-key' metadata header to every request.
func WithAppKey(appKey string) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(AppKeyUnaryInterceptor(appKey))
}

// WithStreamAppKey returns a grpc.DialOption that attaches an 'x-app-key' metadata header to every stream.
func WithStreamAppKey(appKey string) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(AppKeyStreamInterceptor(appKey))
}

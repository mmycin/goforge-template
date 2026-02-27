package grpc_test

import (
	"context"
	"testing"

	"github.com/mmycin/goforge/internal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestNewGrpcClient(t *testing.T) {
	target := "localhost:50051"

	// Test basic initialization
	c, err := client.NewGrpcClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create grpc client: %v", err)
	}
	defer c.Close()

	if c.Target() != target {
		t.Errorf("Expected target %s, got %s", target, c.Target())
	}

	if c.Conn() == nil {
		t.Error("Underlying connection is nil")
	}
}

func TestNewGrpcClientEmptyTarget(t *testing.T) {
	_, err := client.NewGrpcClient("")
	if err == nil {
		t.Error("Expected error for empty target, got nil")
	}
}

func TestAppKeyUnaryInterceptor(t *testing.T) {
	appKey := "test-app-key"
	interceptor := client.AppKeyUnaryInterceptor(appKey)

	// Mock unary invoker
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			t.Fatal("metadata not found in outgoing context")
		}

		keys := md.Get("x-app-key")
		if len(keys) == 0 {
			t.Fatal("x-app-key not found in metadata")
		}

		if keys[0] != appKey {
			t.Errorf("expected app key %s, got %s", appKey, keys[0])
		}

		return nil
	}

	// Execution
	err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, mockInvoker)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
}

func TestAppKeyStreamInterceptor(t *testing.T) {
	appKey := "test-stream-app-key"
	interceptor := client.AppKeyStreamInterceptor(appKey)

	// Mock streamer
	mockStreamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			t.Fatal("metadata not found in outgoing context")
		}

		keys := md.Get("x-app-key")
		if len(keys) == 0 {
			t.Fatal("x-app-key not found in metadata")
		}

		if keys[0] != appKey {
			t.Errorf("expected app key %s, got %s", appKey, keys[0])
		}

		return nil, nil
	}

	// Execution
	_, err := interceptor(context.Background(), nil, nil, "/test.Service/StreamMethod", mockStreamer)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
}

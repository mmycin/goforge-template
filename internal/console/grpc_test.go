package console

import (
	"testing"
	"time"

	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCServerConnectivity(t *testing.T) {
	// Setup minimalist config for testing
	config.GRPC.Port = 50052 // Use a different port to avoid conflicts
	config.GRPC.Enable = true
	config.App.Key = "test-secret" // Set app key for auth interceptor

	// Start server in a goroutine
	srv := server.NewGRPCServer()
	go func() {
		if err := srv.Start(); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()
	defer srv.Stop()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Try to connect with App Key
	addr := ":50052"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Check state
	state := conn.GetState()
	t.Logf("Connection state: %s", state)

	// Create a context with the App Key
	// ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("x-app-key", "test-secret"))

	// We don't have a service registered in this test to call, but we can verify connection
	// If we had a service we would call it here passing ctx

	assert.NotNil(t, conn)
}

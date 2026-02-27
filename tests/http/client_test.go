package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmycin/goforge/internal/client"
)

func TestHttpClient_Do_AddsHeader(t *testing.T) {
	apiKey := "test-api-key"

	// Create a test server to verify the header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey := r.Header.Get("X-App-Key")
		if gotKey != apiKey {
			t.Errorf("expected X-App-Key %s, got %s", apiKey, gotKey)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpClient := client.NewHttpClient(client.WithHttpAppKey(apiKey))

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
}

func TestHttpClient_Get_AddsHeader(t *testing.T) {
	apiKey := "test-api-key"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey := r.Header.Get("X-App-Key")
		if gotKey != apiKey {
			t.Errorf("expected X-App-Key %s, got %s", apiKey, gotKey)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpClient := client.NewHttpClient(client.WithHttpAppKey(apiKey))

	resp, err := httpClient.Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
}

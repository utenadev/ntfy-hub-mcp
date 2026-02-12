package ntfy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_Publish(t *testing.T) {
	testMessage := "Test Message"
	testTitle := "Test Title"
	testTopic := "test_topic"

	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/"+testTopic {
			t.Errorf("Expected URL path /%s, got %s", testTopic, r.URL.Path)
		}
		if r.Header.Get("Title") != testTitle {
			t.Errorf("Expected Title header %s, got %s", testTitle, r.Header.Get("Title"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != testMessage {
			t.Errorf("Expected request body %s, got %s", testMessage, string(body))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Publish(testTopic, testMessage, testTitle)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
}

func TestClient_SubscribeOnce(t *testing.T) {
	testMessage := "Hello from ntfy"
	testTopic := "test_sub_topic"

	done := make(chan bool)

	// Mock SSE server - send message and close connection
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+testTopic+"/json" {
			t.Errorf("Expected URL path /%s/json, got %s", testTopic, r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send SSE message with proper format (ntfy.sh returns raw JSON lines)
		jsonData := fmt.Sprintf(`{"id":"1","time":1678886400,"event":"message","topic":"%s","message":"%s"}`, testTopic, testMessage)
		fmt.Fprint(w, jsonData)
		fmt.Fprint(w, "\n")
		w.(http.Flusher).Flush()

		// Wait a bit to ensure client reads the message
		time.Sleep(100 * time.Millisecond)
		close(done)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	msg, err := client.SubscribeOnce(ctx, testTopic)
	if err != nil {
		t.Fatalf("SubscribeOnce failed: %v", err)
	}
	if msg == nil {
		t.Fatalf("Expected a message, got nil")
	}
	if msg.Message != testMessage {
		t.Errorf("Expected message '%s', got '%s'", testMessage, msg.Message)
	}
	if msg.Topic != testTopic {
		t.Errorf("Expected topic '%s', got '%s'", testTopic, msg.Topic)
	}
}

func TestClient_SubscribeOnce_Timeout(t *testing.T) {
	testTopic := "timeout_topic"

	// Mock SSE server that sends no message
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		// Do not send any message
		// Keep handler alive until client disconnects or context cancels
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err := client.SubscribeOnce(ctx, testTopic)
	if err == nil {
		t.Fatalf("Expected timeout error, got nil")
	}
	// Check if the error is due to context timeout
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected context deadline exceeded error, got %v", err)
	}
}

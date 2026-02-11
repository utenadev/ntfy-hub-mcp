package ntfy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Message represents a ntfy message
type Message struct {
	ID      string   `json:"id"`
	Time    int64    `json:"time"`
	Event   string   `json:"event"`
	Topic   string   `json:"topic"`
	Message string   `json:"message"`
	Title   string   `json:"title,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// Client handles communication with ntfy
type Client struct {
	BaseURL string
}

// NewClient creates a new ntfy client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://ntfy.sh"
	}
	return &Client{BaseURL: baseURL}
}

// Publish sends a message to a topic
func (c *Client) Publish(topic string, message string, title string) error {
	url := fmt.Sprintf("%s/%s", c.BaseURL, topic)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(message))
	if err != nil {
		return err
	}

	if title != "" {
		req.Header.Set("Title", title)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to publish message, status: %d", resp.StatusCode)
	}

	return nil
}

// Subscribe listens for messages on a topic using SSE
func (c *Client) Subscribe(topic string, handler func(Message)) error {
	url := fmt.Sprintf("%s/%s/json", c.BaseURL, topic)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var msg Message
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		if err := json.Unmarshal(line, &msg); err != nil {
			// Ignore non-json lines (keep-alive or comments)
			continue
		}

		if msg.Event == "message" {
			handler(msg)
		}
	}

	return scanner.Err()
}

// SubscribeOnce listens for the first message on a topic using SSE with context for timeout
func (c *Client) SubscribeOnce(ctx context.Context, topic string) (*Message, error) {
	url := fmt.Sprintf("%s/%s/json", c.BaseURL, topic)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var msg Message
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}

		if msg.Event == "message" {
			return &msg, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("no message received")
}

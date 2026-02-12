package ntfy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// Default ntfy server URL
	defaultBaseURL = "https://ntfy.sh"

	// HTTP headers
	headerTitle = "Title"

	// SSE event types
	eventMessage = "message"

	// Error messages
	errPublishFailed  = "failed to publish message, status: %d"
	errNoMessage     = "no message received"
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

// HTTPClient interface for testability
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

// Client handles communication with ntfy
type Client struct {
	BaseURL string
	http    HTTPClient
}

// NewClient creates a new ntfy client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		BaseURL: baseURL,
		http:    http.DefaultClient,
	}
}

// setHTTPClient sets a custom HTTP client (used for testing)
func (c *Client) setHTTPClient(client HTTPClient) {
	c.http = client
}

// buildURL builds the full URL for a topic
func (c *Client) buildURL(topic string, withJSON bool) string {
	url := fmt.Sprintf("%s/%s", c.BaseURL, topic)
	if withJSON {
		url += "/json"
	}
	return url
}

// Publish sends a message to a topic
func (c *Client) Publish(topic string, message string, title string) error {
	req, err := http.NewRequest("POST", c.buildURL(topic, false), bytes.NewBufferString(message))
	if err != nil {
		return err
	}

	if title != "" {
		req.Header.Set(headerTitle, title)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(errPublishFailed, resp.StatusCode)
	}

	return nil
}

// Subscribe listens for messages on a topic using SSE
func (c *Client) Subscribe(topic string, handler func(Message)) error {
	resp, err := c.http.Get(c.buildURL(topic, true))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.processMessages(resp.Body, handler)
}

// processMessages processes SSE messages from a reader
func (c *Client) processMessages(r io.Reader, handler func(Message)) error {
	scanner := bufio.NewScanner(r)
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

		if msg.Event == eventMessage {
			handler(msg)
		}
	}

	return scanner.Err()
}

// SubscribeOnce listens for the first message on a topic using SSE with context for timeout
func (c *Client) SubscribeOnce(ctx context.Context, topic string) (*Message, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.buildURL(topic, true), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	msg, err := c.waitForFirstMessage(resp.Body)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf(errNoMessage)
	}
	return msg, nil
}

// waitForFirstMessage waits for the first message event from a reader
func (c *Client) waitForFirstMessage(r io.Reader) (*Message, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var msg Message
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}

		if msg.Event == eventMessage {
			return &msg, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

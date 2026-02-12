package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"ntfy-hub-mcp/ntfy"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	// Server metadata
	serverName    = "ntfy-hub-mcp"
	serverVersion = "1.0.0"

	// Default configuration values
	defaultNtfyURL     = "https://ntfy.sh"
	defaultTopicOut    = "agent-output"
	defaultTopicIn     = "agent-input"
	defaultTimeoutSecs = 60

	// Environment variable keys
	envNtfyURL     = "NTFY_URL"
	envTopicOut    = "NTFY_TOPIC_OUT"
	envTopicIn     = "NTFY_TOPIC_IN"

	// Tool names
	toolPublish     = "ntfy_publish"
	toolWaitForReply = "ntfy_wait_for_reply"

	// Parameter names
	paramMessage      = "message"
	paramTopic        = "topic"
	paramTitle        = "title"
	paramTimeoutSecs  = "timeout_seconds"
	paramPrompt       = "prompt"

	// Error messages
	errMessageRequired  = "message is required"
	errPublishFailed    = "Failed to publish to ntfy: %v"
	errTimeoutWaiting  = "Timed out waiting for reply on topic '%s' after %d seconds"
	errWaitingForReply = "Error waiting for reply: %v"
)

// config holds the server configuration
type config struct {
	ntfyURL      string
	topicOut     string
	topicIn      string
}

func loadConfig() *config {
	return &config{
		ntfyURL:  getEnv(envNtfyURL, defaultNtfyURL),
		topicOut: getEnv(envTopicOut, defaultTopicOut),
		topicIn:  getEnv(envTopicIn, defaultTopicIn),
	}
}

func main() {
	cfg := loadConfig()
	client := ntfy.NewClient(cfg.ntfyURL)

	s := server.NewMCPServer(serverName, serverVersion)
	registerTools(s, client, cfg)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func registerTools(s *server.MCPServer, client *ntfy.Client, cfg *config) {
	s.AddTool(mcp.NewTool(toolPublish,
		mcp.WithDescription("Send a message to a ntfy topic (e.g., for notifications)"),
		mcp.WithString(paramMessage, mcp.Description("The message content to send"), mcp.Required()),
		mcp.WithString(paramTopic, mcp.Description(fmt.Sprintf("The topic to publish to (default: %s)", cfg.topicOut))),
		mcp.WithString(paramTitle, mcp.Description("Optional title for the notification")),
	), makePublishHandler(client, cfg))

	s.AddTool(mcp.NewTool(toolWaitForReply,
		mcp.WithDescription("Wait for a reply from the human on a specific topic. Use this to get user input or approval."),
		mcp.WithString(paramTopic, mcp.Description(fmt.Sprintf("The topic to listen on (default: %s)", cfg.topicIn))),
		mcp.WithNumber(paramTimeoutSecs, mcp.Description("How long to wait for a reply in seconds (default: 60)")),
		mcp.WithString(paramPrompt, mcp.Description("Optional message to send to the human before waiting for a reply")),
	), makeWaitForReplyHandler(client, cfg))
}

func makePublishHandler(client *ntfy.Client, cfg *config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		message := mcp.ParseString(request, paramMessage, "")
		topic := mcp.ParseString(request, paramTopic, cfg.topicOut)
		title := mcp.ParseString(request, paramTitle, "")

		if message == "" {
			return mcp.NewToolResultError(errMessageRequired), nil
		}

		if err := client.Publish(topic, message, title); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(errPublishFailed, err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Message sent to topic '%s'", topic)), nil
	}
}

func makeWaitForReplyHandler(client *ntfy.Client, cfg *config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		topic := mcp.ParseString(request, paramTopic, cfg.topicIn)
		timeoutSeconds := int(mcp.ParseFloat64(request, paramTimeoutSecs, defaultTimeoutSecs))
		prompt := mcp.ParseString(request, paramPrompt, "")

		if prompt != "" {
			_ = client.Publish(cfg.topicOut, prompt, "Instruction Requested")
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		defer cancel()

		msg, err := client.SubscribeOnce(timeoutCtx, topic)
		if err != nil {
			if err == context.DeadlineExceeded {
				return mcp.NewToolResultError(fmt.Sprintf(errTimeoutWaiting, topic, timeoutSeconds)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf(errWaitingForReply, err)), nil
		}

		return mcp.NewToolResultText(msg.Message), nil
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

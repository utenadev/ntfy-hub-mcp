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

func main() {
	// Configuration from environment variables
	ntfyURL := getEnv("NTFY_URL", "https://ntfy.sh")
	defaultTopicOut := getEnv("NTFY_TOPIC_OUT", "agent-output")
	defaultTopicIn := getEnv("NTFY_TOPIC_IN", "agent-input")

	client := ntfy.NewClient(ntfyURL)

	// Create MCP server
	s := server.NewMCPServer(
		"ntfy-hub-mcp",
		"1.0.0",
	)

	// Tool: ntfy_publish
	s.AddTool(mcp.NewTool("ntfy_publish",
		mcp.WithDescription("Send a message to a ntfy topic (e.g., for notifications)"),
		mcp.WithString("message", mcp.Description("The message content to send"), mcp.Required()),
		mcp.WithString("topic", mcp.Description(fmt.Sprintf("The topic to publish to (default: %s)", defaultTopicOut))),
		mcp.WithString("title", mcp.Description("Optional title for the notification")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		message := mcp.ParseString(request, "message", "")
		topic := mcp.ParseString(request, "topic", defaultTopicOut)
		title := mcp.ParseString(request, "title", "")

		if message == "" {
			return mcp.NewToolResultError("message is required"), nil
		}

		err := client.Publish(topic, message, title)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to publish to ntfy: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Message sent to topic '%s'", topic)), nil
	})

	// Tool: ntfy_wait_for_reply
	s.AddTool(mcp.NewTool("ntfy_wait_for_reply",
		mcp.WithDescription("Wait for a reply from the human on a specific topic. Use this to get user input or approval."),
		mcp.WithString("topic", mcp.Description(fmt.Sprintf("The topic to listen on (default: %s)", defaultTopicIn))),
		mcp.WithNumber("timeout_seconds", mcp.Description("How long to wait for a reply in seconds (default: 60)")),
		mcp.WithString("prompt", mcp.Description("Optional message to send to the human before waiting for a reply")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		topic := mcp.ParseString(request, "topic", defaultTopicIn)
		timeoutSeconds := int(mcp.ParseFloat64(request, "timeout_seconds", 60))
		prompt := mcp.ParseString(request, "prompt", "")

		// Optional prompt: send it before waiting
		if prompt != "" {
			_ = client.Publish(defaultTopicOut, prompt, "Instruction Requested")
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		defer cancel()

		msg, err := client.SubscribeOnce(timeoutCtx, topic)
		if err != nil {
			if err == context.DeadlineExceeded {
				return mcp.NewToolResultError(fmt.Sprintf("Timed out waiting for reply on topic '%s' after %d seconds", topic, timeoutSeconds)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Error waiting for reply: %v", err)), nil
		}

		return mcp.NewToolResultText(msg.Message), nil
	})

	// Start the server using stdio transport
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

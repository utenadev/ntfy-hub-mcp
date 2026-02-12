# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`ntfy-hub-mcp` is a lightweight MCP (Model Context Protocol) server that bridges AI agents and humans through `ntfy.sh` for asynchronous communication. The project is written in Go and uses `go-task` for build automation.

### Architecture

```
ntfy-hub-mcp/
├── main.go              # MCP server entry point (stdio transport)
└── ntfy/                # ntfy.sh client package
    ├── client.go        # HTTP/SSE client for publish/subscribe
    └── client_test.go   # Unit tests with httptest mocks
```

**Core Components:**
- **MCP Server** (`main.go`): Uses `github.com/mark3labs/mcp-go/mcp` to expose two tools
  - `ntfy_publish`: Send messages from agent to human
  - `ntfy_wait_for_reply`: Wait for human input/approval with optional prompt
- **ntfy Client** (`ntfy/client.go`): Handles HTTP POST (publish) and SSE streaming (subscribe)

## Development Commands

All commands use `go-task`:

```bash
task build      # Build ntfy-hub-mcp.exe (Windows) or ntfy-hub-mcp (Unix)
task run        # Start the MCP server with configurable env vars
task test       # Run Go tests: go test -v ntfy-hub-mcp/ntfy
task lint       # Run go fmt ./... and go vet ./...
task clean      # Remove build artifacts
task install-tools  # Install golangci-lint
```

## Environment Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `NTFY_URL` | `https://ntfy.sh` | ntfy server base URL |
| `NTFY_TOPIC_OUT` | `agent-output` | Agent→Human notification topic |
| `NTFY_TOPIC_IN` | `agent-input` | Human→Agent instruction topic |

## Key Design Patterns

1. **Environment-based configuration**: Uses `getEnv()` helper with sensible defaults
2. **Context-based timeouts**: `SubscribeOnce` accepts `context.Context` for cancellation
3. **Clean separation**: MCP server logic (main.go) vs ntfy client (ntfy/ package)
4. **SSE message filtering**: Only processes events where `msg.Event == "message"`

## Documentation Structure

The project uses a "cospec"-inspired documentation structure:
- `docs/SPEC.md`: Detailed tool specifications with examples
- `docs/USAGE.md`: Integration guide with Gemini CLI
- `docs/TROUBLESHOOTING.md`: Common issues and solutions
- `docs/BLUEPRINT.md`: Architecture overview and future plans
- `docs/RELATIONSHIP_*.md`: Ecosystem integration context

## MCP Ecosystem

According to `docs/RELATIONSHIP_MCP_HUBS.md`, this server is part of a larger MCP ecosystem:
- Works with `agent-hub-mcp` for project coordination
- Integrates with `gistpad-mcp` for knowledge sharing
- Provides the real-time communication layer for human-agent interaction

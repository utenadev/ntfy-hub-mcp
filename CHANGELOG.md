# Changelog

All notable changes to ntfy-hub-mcp will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Initial Project Setup**: Created `ntfy-hub-mcp` Go project.
  - Renamed module and directory from `ntfy_hub` to `ntfy-hub-mcp`.
  - Implemented MCP server using `github.com/mark3labs/mcp-go`.
  - Provided `ntfy_publish` tool for sending notifications.
  - Provided `ntfy_wait_for_reply` tool for waiting for human input with timeout.
  - Added `SubscribeOnce` method to `ntfy/client.go` for single message reception.
  - Integrated `go-task` with `Taskfile.yml` for build, test, run, lint, and clean tasks.
  - Created `ntfy/client_test.go` with mock HTTP/SSE servers for unit testing.
  - Established `docs` directory structure inspired by `cospec` project.
  - Documented `BLUEPRINT.md`, `PLAN.md`, `SPEC.md`, `USAGE.md`, `TROUBLESHOOTING.md`.
  - Added `RELATIONSHIP_NTFY_SH.md` explaining the connection to `ntfy.sh`.
  - Added `RELATIONSHIP_MCP_HUBS.md` outlining the scope and relationships with `agent-hub-mcp` and `gistpad-mcp`.
  - Created `LICENSE` file (MIT License).
  - Created `README.md` in Japanese.
- **Git Repository Initialization**: Initialized Git repository and pushed to GitHub (`https://github.com/utenadev/ntfy-hub-mcp`).
- **CI/CD**: Added GitHub Actions workflow (`.github/workflows/test.yml`) for automated testing on push and pull requests.
- **Documentation**: Added `CLAUDE.md` for Claude Code guidance with project overview and development commands.
- **Documentation**: Consolidated root `PLAN.md` into `docs/PLAN.md` with architecture overview.

### Changed
- **Refactoring**: Improved code organization in `main.go`:
  - Extracted constants for server metadata, default values, environment variables, tool names, parameters, and error messages.
  - Added `config` struct to hold server configuration.
  - Split tool handlers into separate functions (`makePublishHandler`, `makeWaitForReplyHandler`).
  - Added `registerTools` function for cleaner tool registration.
- **Refactoring**: Enhanced `ntfy/client.go`:
  - Introduced `HTTPClient` interface for better testability.
  - Added `buildURL` helper method for URL construction.
  - Extracted `processMessages` and `waitForFirstMessage` helper methods.
  - Added constants for default values, headers, event types, and error messages.
- **Testing**: Fixed `TestClient_SubscribeOnce` to properly handle SSE message format and connection closure.

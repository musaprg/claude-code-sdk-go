# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go SDK for the Claude Code CLI, providing a clean and idiomatic Go API for interacting with Claude Code. The SDK handles subprocess communication, message parsing, and provides type-safe interfaces for all operations.

## Development Commands

### Testing
- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output  
- `go test -race ./...` - Run tests with race detection
- `go test -race -coverprofile=coverage.out -covermode=atomic ./...` - Run tests with coverage
- `go tool cover -html=coverage.out` - View coverage report
- `go test -run TestSpecificFunction` - Run specific test function

### Building
- `go build -v ./...` - Build all packages
- `go build -o basic_example ./examples/basic_usage` - Build basic usage example
- `go build -o streaming_example ./examples/streaming_mode` - Build streaming mode example

### Code Quality
- `go fmt ./...` - Format all Go code
- `gofmt -s -w .` - Format with simplifications
- `go vet ./...` - Run go vet static analysis
- `golangci-lint run ./...` - Run comprehensive linter (requires golangci-lint)

### Development Helpers
- `go mod tidy` - Clean up go.mod and go.sum
- `go mod download && go mod verify` - Download and verify dependencies

## Architecture Overview

This is a Go SDK for interacting with Claude Code CLI via subprocess communication. The architecture follows a layered design:

### Core Components

1. **Main Package (`./`)**: Public API surface
   - `client.go`: Main Client type and Query convenience function
   - `types.go`: Re-exports all public types from internal packages
   - `errors.go`: Re-exports error types

2. **Internal Transport Layer (`internal/transport/`)**: 
   - `subprocess.go`: Manages Claude CLI subprocess communication via stdin/stdout/stderr
   - Handles JSON message parsing and CLI argument construction
   - Provides cleanup and graceful termination

3. **Internal Types (`internal/types/`)**: Core type definitions
   - Message interfaces (UserMessage, AssistantMessage, SystemMessage, ResultMessage)
   - Content block interfaces (TextBlock, ToolUseBlock, ToolResultBlock)
   - Configuration types (QueryOptions, ClientOptions, McpServerConfig)

4. **Internal Parser (`internal/parser/`)**: 
   - `message.go`: JSON message parsing from CLI output

5. **Internal CLI Utils (`internal/cli/`)**: 
   - `finder.go`: Locates Claude CLI executable

6. **Internal Errors (`internal/errors/`)**: 
   - `errors.go`: Specific error types (CLINotFoundError, CLIConnectionError, ProcessError, MessageParseError)

### Message Flow

The SDK follows this flow:
1. Client creates a SubprocessTransport with CLI path and working directory
2. Transport builds CLI arguments from QueryOptions and executes `claude --output-format stream-json`
3. JSON messages stream from CLI stdout and are parsed by the parser package
4. Parsed messages are sent through Go channels back to the client
5. Transport handles cleanup and graceful process termination

### Key Patterns

- **Interface-based design**: Message and ContentBlock are interfaces allowing polymorphic handling
- **Channel-based streaming**: Uses Go channels instead of iterators for message streaming  
- **Context support**: All operations accept context.Context for cancellation and timeouts
- **Clean separation**: Internal packages handle implementation details, main package provides clean API
- **Type re-exports**: Main package re-exports internal types to provide a single import point

### Prerequisites

- Go 1.24.2 or later
- Node.js and Claude Code CLI: `npm install -g @anthropic-ai/claude-code`
- For development: golangci-lint for comprehensive linting

### Usage Patterns

The SDK supports both simple one-shot queries via `claudecode.Query()` and more advanced usage via `claudecode.NewClient()`. All operations return channels that stream messages as they're received from the CLI subprocess.

Messages are parsed from CLI's `--output-format stream-json` and delivered as typed Go structs implementing the Message interface.

## Important Implementation Details

### Error Handling
The SDK provides specific error types that wrap different failure modes:
- `CLINotFoundError`: Claude CLI binary not found in PATH or common locations
- `CLIConnectionError`: Issues with subprocess pipes or working directory
- `ProcessError`: CLI process execution failures with exit codes and stderr
- `MessageParseError`: JSON parsing errors from CLI output

Always check error types using `errors.As()` for proper error handling.

### Subprocess Management  
The transport layer uses `--output-format stream-json` mode and immediately closes stdin since queries use `--print` mode. The subprocess is terminated gracefully with SIGINT, falling back to SIGKILL after a 5-second timeout.

### Channel-based Streaming
All operations return Go channels that stream messages as they're received. Channels are buffered (size 10) and automatically closed when the CLI process terminates or context is cancelled.

### Context Support
All operations accept `context.Context` for cancellation and timeouts. Context cancellation properly terminates the CLI subprocess and closes all channels.

## Examples

The `examples/` directory contains comprehensive usage examples:

- **`basic_usage/`**: Simple query patterns and client configuration
- **`streaming_mode/`**: Advanced streaming patterns including:
  - Multi-turn conversations
  - Context cancellation and timeouts
  - Custom options and tool configuration
  - Detailed message processing
  
Run examples with:
```bash
go run ./examples/basic_usage
go run ./examples/streaming_mode
```
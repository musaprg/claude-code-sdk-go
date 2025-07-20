# Claude Code SDK for Go

A Go SDK for interacting with Claude Code, providing a clean and idiomatic API for Go developers.

## Overview

This SDK provides both simple one-shot queries and advanced conversation management for Claude Code. It's designed to match the functionality of the official Python SDK while following Go conventions and best practices.

## Installation

```bash
go get github.com/musaprg/claude-code-sdk-go
```

## Prerequisites

- Go 1.24.2 or later
- Node.js (required for Claude Code CLI)
- Claude Code CLI installed: `npm install -g @anthropic-ai/claude-code`

## Quick Start

### Simple Query

```go
package main

import (
    "context"
    "fmt"
    "log"

    claudecode "github.com/musaprg/claude-code-sdk-go"
)

func main() {
    ctx := context.Background()
    
    // Simple query
    messageCh, err := claudecode.Query(ctx, "What is 2 + 2?", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Process messages
    for message := range messageCh {
        switch msg := message.(type) {
        case *claudecode.AssistantMessage:
            for _, block := range msg.Content {
                if textBlock, ok := block.(*claudecode.TextBlock); ok {
                    fmt.Printf("Claude: %s\n", textBlock.Text)
                }
            }
        case *claudecode.ResultMessage:
            fmt.Printf("Cost: $%.4f\n", *msg.TotalCostUSD)
        }
    }
}
```

### Using Client with Options

```go
package main

import (
    "context"
    "fmt"
    "log"

    claudecode "github.com/musaprg/claude-code-sdk-go"
)

func main() {
    ctx := context.Background()
    
    // Create client with options
    client := claudecode.NewClient(&claudecode.ClientOptions{
        CWD: "/path/to/project",
    })
    
    // Configure query options
    options := &claudecode.QueryOptions{
        SystemPrompt:   "You are a helpful coding assistant",
        MaxTurns:       5,
        AllowedTools:   []string{"Read", "Write", "Bash"},
        PermissionMode: claudecode.PermissionModeAcceptEdits,
    }
    
    messageCh, err := client.Query(ctx, "Help me refactor this code", options)
    if err != nil {
        log.Fatal(err)
    }

    for message := range messageCh {
        handleMessage(message)
    }
}

func handleMessage(message claudecode.Message) {
    switch msg := message.(type) {
    case *claudecode.UserMessage:
        fmt.Printf("User: %s\n", msg.Content)
    case *claudecode.AssistantMessage:
        fmt.Println("Claude:")
        for _, block := range msg.Content {
            switch b := block.(type) {
            case *claudecode.TextBlock:
                fmt.Printf("  %s\n", b.Text)
            case *claudecode.ToolUseBlock:
                fmt.Printf("  Using tool: %s\n", b.Name)
            }
        }
    case *claudecode.SystemMessage:
        fmt.Printf("System [%s]: %v\n", msg.Subtype, msg.Data)
    case *claudecode.ResultMessage:
        fmt.Printf("Result: %s (Duration: %dms)\n", msg.Subtype, msg.DurationMs)
        if msg.TotalCostUSD != nil {
            fmt.Printf("Cost: $%.4f\n", *msg.TotalCostUSD)
        }
    }
}
```

## API Reference

### Types

#### Messages

- **UserMessage**: Represents user input
- **AssistantMessage**: Claude's response with content blocks
- **SystemMessage**: System notifications and metadata
- **ResultMessage**: Execution results with timing and cost information

#### Content Blocks

- **TextBlock**: Plain text content
- **ToolUseBlock**: Tool invocation with parameters
- **ToolResultBlock**: Results from tool execution

#### Configuration

```go
type QueryOptions struct {
    AllowedTools                []string
    SystemPrompt                string
    AppendSystemPrompt          string
    MaxTurns                    int
    PermissionMode              PermissionMode
    Model                       string
    CWD                         string
    // ... and more options
}
```

#### Permission Modes

- `PermissionModeDefault`: CLI prompts for dangerous operations
- `PermissionModeAcceptEdits`: Auto-accept file edits
- `PermissionModeBypassPermissions`: Allow all operations (use with caution)

### Functions

#### Query

```go
func Query(ctx context.Context, prompt string, options *QueryOptions) (<-chan Message, error)
```

Simple function for one-shot queries.

#### Client Methods

```go
func NewClient(options *ClientOptions) *Client
func (c *Client) Query(ctx context.Context, prompt string, options *QueryOptions) (<-chan Message, error)
```

Client for more advanced usage with custom configuration.

### Error Handling

The SDK provides specific error types for different failure scenarios:

- **CLINotFoundError**: Claude Code CLI not found
- **CLIConnectionError**: Connection issues with CLI
- **ProcessError**: CLI process execution errors
- **MessageParseError**: Message parsing errors

```go
messageCh, err := claudecode.Query(ctx, prompt, options)
if err != nil {
    var cliNotFound *claudecode.CLINotFoundError
    if errors.As(err, &cliNotFound) {
        fmt.Printf("CLI not found at: %s\n", cliNotFound.CLIPath)
        return
    }
    log.Fatal(err)
}
```

## Architecture

The SDK is organized into several internal packages:

- `internal/types`: Core type definitions
- `internal/errors`: Error types and handling
- `internal/cli`: CLI discovery and utilities
- `internal/parser`: Message parsing from CLI output
- `internal/transport`: Subprocess communication with Claude CLI

The main package re-exports all public types and functions to provide a clean API.

## Differences from Python SDK

While this Go SDK aims for feature parity with the Python SDK, there are some Go-specific adaptations:

1. **Channels instead of AsyncIterators**: Go uses channels for streaming data
2. **Context for cancellation**: Go's context.Context for timeout and cancellation
3. **Interface-based design**: Go interfaces for Message and ContentBlock types
4. **Explicit error handling**: Go's explicit error handling pattern

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests: `go test ./...`
6. Format code: `go fmt ./...`
7. Submit a pull request

## Related Projects

- [Claude Code Python SDK](https://github.com/anthropics/claude-code-sdk-python) - Official Python SDK
- [Claude Code CLI](https://github.com/anthropics/claude-code) - Command-line interface

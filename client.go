package claudecode

import (
	"context"

	"github.com/musaprg/claude-code-sdk-go/internal/transport"
	"github.com/musaprg/claude-code-sdk-go/internal/types"
)

// Client represents a Claude Code SDK client that manages communication with the Claude Code CLI.
// It handles subprocess execution, message parsing, and provides a clean Go API for Claude Code operations.
type Client struct {
	// cliPath is the path to the Claude Code CLI executable.
	cliPath string
	// cwd is the current working directory for Claude Code operations.
	cwd string
}

// NewClient creates a new Claude Code SDK client with the given options.
// If options is nil, default settings will be used (CLI auto-discovery, current working directory).
func NewClient(options *ClientOptions) *Client {
	client := &Client{}

	if options != nil {
		if options.CLIPath != "" {
			client.cliPath = options.CLIPath
		}
		if options.CWD != "" {
			client.cwd = options.CWD
		}
	}

	return client
}

// Query sends a prompt to Claude Code and returns a channel that streams response messages.
// The returned channel will receive messages as they are generated by Claude Code.
// The channel will be closed when the conversation completes or the context is cancelled.
func (c *Client) Query(ctx context.Context, prompt string, options *QueryOptions) (<-chan Message, error) {
	// Create transport
	transport := transport.NewSubprocessTransport(c.cliPath, c.cwd)

	// Convert options to internal type
	var internalOptions *types.QueryOptions
	if options != nil {
		internalOptions = (*types.QueryOptions)(options)
	}

	// Connect and start the query
	if err := transport.Connect(ctx, internalOptions, prompt); err != nil {
		return nil, err
	}

	// Get message channel
	messageCh, err := transport.ReceiveMessages(ctx)
	if err != nil {
		transport.Disconnect()
		return nil, err
	}

	// Wrap the channel to handle cleanup and type conversion
	wrappedCh := make(chan Message, 10)
	go func() {
		defer close(wrappedCh)
		defer transport.Disconnect()

		for message := range messageCh {
			select {
			case wrappedCh <- message:
			case <-ctx.Done():
				return
			}
		}
	}()

	return wrappedCh, nil
}

// Query is a convenience function that creates a default client and executes a query.
// This is equivalent to calling NewClient(nil).Query(ctx, prompt, options).
func Query(ctx context.Context, prompt string, options *QueryOptions) (<-chan Message, error) {
	client := NewClient(nil)
	return client.Query(ctx, prompt, options)
}

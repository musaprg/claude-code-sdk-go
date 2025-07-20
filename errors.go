package claudecode

import "github.com/musaprg/claude-code-sdk-go/internal/errors"

// Re-export error types from internal package.
// These error types provide specific information about different failure modes in the SDK.
type (
	// ClaudeSDKError is the base error type for all SDK-related errors.
	ClaudeSDKError = errors.ClaudeSDKError
	// CLIConnectionError occurs when there are issues with subprocess pipes or working directory setup.
	CLIConnectionError = errors.CLIConnectionError
	// CLINotFoundError occurs when the Claude Code CLI binary cannot be found in PATH or standard locations.
	CLINotFoundError = errors.CLINotFoundError
	// ProcessError occurs when the Claude Code CLI process exits with an error or fails to execute.
	ProcessError = errors.ProcessError
	// CLIJSONDecodeError occurs when the CLI output cannot be parsed as valid JSON.
	CLIJSONDecodeError = errors.CLIJSONDecodeError
	// MessageParseError occurs when JSON from the CLI cannot be parsed into a Message struct.
	MessageParseError = errors.MessageParseError
)

// Re-export error constructor functions from internal package.
// These functions create specific error types with appropriate context and details.
var (
	// NewClaudeSDKError creates a new base SDK error with the given message.
	NewClaudeSDKError = errors.NewClaudeSDKError
	// NewCLIConnectionError creates a new CLI connection error with the given message and underlying error.
	NewCLIConnectionError = errors.NewCLIConnectionError
	// NewCLINotFoundError creates a new CLI not found error with search paths.
	NewCLINotFoundError = errors.NewCLINotFoundError
	// NewProcessError creates a new process error with exit code, stderr, and underlying error.
	NewProcessError = errors.NewProcessError
	// NewCLIJSONDecodeError creates a new JSON decode error with the problematic JSON and underlying error.
	NewCLIJSONDecodeError = errors.NewCLIJSONDecodeError
	// NewMessageParseError creates a new message parse error with the raw JSON and underlying error.
	NewMessageParseError = errors.NewMessageParseError
)

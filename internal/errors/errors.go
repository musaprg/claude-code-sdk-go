package errors

import "fmt"

// ClaudeSDKError is the base error type for all Claude SDK errors.
// It provides error wrapping functionality and serves as the foundation for more specific error types.
type ClaudeSDKError struct {
	// message contains the human-readable error description.
	message string
	// cause contains the underlying error that caused this error, if any.
	cause error
}

func (e *ClaudeSDKError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *ClaudeSDKError) Unwrap() error {
	return e.cause
}

// NewClaudeSDKError creates a new base SDK error with the given message and optional underlying cause.
func NewClaudeSDKError(message string, cause error) *ClaudeSDKError {
	return &ClaudeSDKError{
		message: message,
		cause:   cause,
	}
}

// CLIConnectionError represents errors that occur when establishing or maintaining
// communication with the Claude Code CLI process, such as pipe creation failures
// or working directory access issues.
type CLIConnectionError struct {
	*ClaudeSDKError
}

// NewCLIConnectionError creates a new CLI connection error with the given message and underlying cause.
func NewCLIConnectionError(message string, cause error) *CLIConnectionError {
	return &CLIConnectionError{
		ClaudeSDKError: NewClaudeSDKError(message, cause),
	}
}

// CLINotFoundError represents an error when the Claude Code CLI executable
// cannot be found in the specified path, PATH environment variable, or standard installation locations.
type CLINotFoundError struct {
	*ClaudeSDKError
	// CLIPath contains the path where the CLI was expected to be found.
	CLIPath string
}

// NewCLINotFoundError creates a new CLI not found error with the path that was searched.
func NewCLINotFoundError(cliPath string, cause error) *CLINotFoundError {
	message := fmt.Sprintf("Claude Code CLI not found at path: %s", cliPath)
	return &CLINotFoundError{
		ClaudeSDKError: NewClaudeSDKError(message, cause),
		CLIPath:        cliPath,
	}
}

// ProcessError represents errors that occur during Claude Code CLI process execution,
// such as the process exiting with a non-zero exit code or crashing unexpectedly.
type ProcessError struct {
	*ClaudeSDKError
	// ExitCode contains the exit code returned by the CLI process.
	ExitCode int
	// Stderr contains any error output from the CLI process.
	Stderr string
}

// NewProcessError creates a new process error with the given details about the CLI process failure.
func NewProcessError(message string, exitCode int, stderr string, cause error) *ProcessError {
	return &ProcessError{
		ClaudeSDKError: NewClaudeSDKError(message, cause),
		ExitCode:       exitCode,
		Stderr:         stderr,
	}
}

// CLIJSONDecodeError represents errors that occur when the output from the Claude Code CLI
// cannot be parsed as valid JSON, typically indicating malformed or corrupted output.
type CLIJSONDecodeError struct {
	*ClaudeSDKError
	// RawData contains the raw string data that could not be parsed as JSON.
	RawData string
}

// NewCLIJSONDecodeError creates a new JSON decode error with the problematic raw data.
func NewCLIJSONDecodeError(message string, rawData string, cause error) *CLIJSONDecodeError {
	return &CLIJSONDecodeError{
		ClaudeSDKError: NewClaudeSDKError(message, cause),
		RawData:        rawData,
	}
}

// MessageParseError represents errors that occur when valid JSON from the CLI
// cannot be parsed into the expected Message struct format, indicating unexpected message structure.
type MessageParseError struct {
	*ClaudeSDKError
	// RawData contains the raw data that could not be parsed into a Message.
	RawData any
}

// NewMessageParseError creates a new message parse error with the raw data that could not be parsed.
func NewMessageParseError(message string, rawData any, cause error) *MessageParseError {
	return &MessageParseError{
		ClaudeSDKError: NewClaudeSDKError(message, cause),
		RawData:        rawData,
	}
}

package claudecode

import (
	"errors"
	"strings"
	"testing"
)

func TestClaudeSDKError(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		err := NewClaudeSDKError("test error", nil)
		if err.Error() != "test error" {
			t.Errorf("ClaudeSDKError.Error() = %q, want %q", err.Error(), "test error")
		}
		if err.Unwrap() != nil {
			t.Errorf("ClaudeSDKError.Unwrap() = %v, want nil", err.Unwrap())
		}
	})

	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewClaudeSDKError("test error", cause)
		expected := "test error: underlying error"
		if err.Error() != expected {
			t.Errorf("ClaudeSDKError.Error() = %q, want %q", err.Error(), expected)
		}
		if err.Unwrap() != cause {
			t.Errorf("ClaudeSDKError.Unwrap() = %v, want %v", err.Unwrap(), cause)
		}
	})
}

func TestCLIConnectionError(t *testing.T) {
	cause := errors.New("connection refused")
	err := NewCLIConnectionError("failed to connect", cause)

	if !strings.Contains(err.Error(), "failed to connect") {
		t.Errorf("CLIConnectionError.Error() should contain 'failed to connect', got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("CLIConnectionError.Error() should contain 'connection refused', got %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Errorf("CLIConnectionError should wrap the cause error")
	}
}

func TestCLINotFoundError(t *testing.T) {
	cliPath := "/usr/bin/claude-code"
	cause := errors.New("file not found")
	err := NewCLINotFoundError(cliPath, cause)

	if err.CLIPath != cliPath {
		t.Errorf("CLINotFoundError.CLIPath = %q, want %q", err.CLIPath, cliPath)
	}
	if !strings.Contains(err.Error(), cliPath) {
		t.Errorf("CLINotFoundError.Error() should contain CLI path, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("CLINotFoundError.Error() should contain 'not found', got %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Errorf("CLINotFoundError should wrap the cause error")
	}
}

func TestProcessError(t *testing.T) {
	message := "process failed"
	exitCode := 1
	stderr := "command not found"
	cause := errors.New("exit status 1")
	err := NewProcessError(message, exitCode, stderr, cause)

	if err.ExitCode != exitCode {
		t.Errorf("ProcessError.ExitCode = %d, want %d", err.ExitCode, exitCode)
	}
	if err.Stderr != stderr {
		t.Errorf("ProcessError.Stderr = %q, want %q", err.Stderr, stderr)
	}
	if !strings.Contains(err.Error(), message) {
		t.Errorf("ProcessError.Error() should contain message, got %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Errorf("ProcessError should wrap the cause error")
	}
}

func TestCLIJSONDecodeError(t *testing.T) {
	message := "failed to decode JSON"
	rawData := `{"invalid": json}`
	cause := errors.New("unexpected end of JSON input")
	err := NewCLIJSONDecodeError(message, rawData, cause)

	if err.RawData != rawData {
		t.Errorf("CLIJSONDecodeError.RawData = %q, want %q", err.RawData, rawData)
	}
	if !strings.Contains(err.Error(), message) {
		t.Errorf("CLIJSONDecodeError.Error() should contain message, got %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Errorf("CLIJSONDecodeError should wrap the cause error")
	}
}

func TestErrorChaining(t *testing.T) {
	// Test that error chaining works properly
	rootCause := errors.New("root cause")
	sdkError := NewClaudeSDKError("sdk error", rootCause)
	cliError := NewCLIConnectionError("connection error", sdkError)

	// Test errors.Is works through the chain
	if !errors.Is(cliError, rootCause) {
		t.Errorf("errors.Is should find root cause through error chain")
	}
	if !errors.Is(cliError, sdkError) {
		t.Errorf("errors.Is should find sdk error in chain")
	}

	// Test that error message contains all parts
	errorStr := cliError.Error()
	if !strings.Contains(errorStr, "connection error") {
		t.Errorf("Error string should contain 'connection error', got %q", errorStr)
	}
}

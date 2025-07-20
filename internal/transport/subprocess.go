package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/musaprg/claude-code-sdk-go/internal/cli"
	"github.com/musaprg/claude-code-sdk-go/internal/errors"
	"github.com/musaprg/claude-code-sdk-go/internal/parser"
	"github.com/musaprg/claude-code-sdk-go/internal/types"
)

const (
	maxBufferSize = 1024 * 1024      // 1MB buffer limit
	maxStderrSize = 10 * 1024 * 1024 // 10MB stderr limit
	stderrTimeout = 30 * time.Second
)

// SubprocessTransport handles communication with Claude CLI via subprocess
type SubprocessTransport struct {
	cliPath string
	cwd     string
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
}

// NewSubprocessTransport creates a new subprocess transport
func NewSubprocessTransport(cliPath string, cwd string) *SubprocessTransport {
	return &SubprocessTransport{
		cliPath: cliPath,
		cwd:     cwd,
	}
}

// Connect starts the subprocess
func (t *SubprocessTransport) Connect(ctx context.Context, options *types.QueryOptions, prompt string) error {
	if t.cliPath == "" {
		var err error
		t.cliPath, err = cli.FindCLI()
		if err != nil {
			return err
		}
	}

	// Build command arguments
	args := t.buildCommand(options, prompt)

	t.cmd = exec.CommandContext(ctx, t.cliPath, args...)

	// Set working directory
	if t.cwd != "" {
		if _, err := os.Stat(t.cwd); err != nil {
			return errors.NewCLIConnectionError(
				fmt.Sprintf("Working directory does not exist: %s", t.cwd), err)
		}
		t.cmd.Dir = t.cwd
	}

	// Set environment
	t.cmd.Env = append(os.Environ(),
		"CLAUDE_CODE_ENTRYPOINT=sdk-go",
		"FORCE_COLOR=0",       // Disable color output which might affect buffering
		"NODE_ENV=production") // Ensure consistent node environment

	// Set up pipes
	var err error
	t.stdin, err = t.cmd.StdinPipe()
	if err != nil {
		return errors.NewCLIConnectionError("failed to create stdin pipe", err)
	}

	t.stdout, err = t.cmd.StdoutPipe()
	if err != nil {
		return errors.NewCLIConnectionError("failed to create stdout pipe", err)
	}

	t.stderr, err = t.cmd.StderrPipe()
	if err != nil {
		return errors.NewCLIConnectionError("failed to create stderr pipe", err)
	}

	// Start the process
	if err := t.cmd.Start(); err != nil {
		// Check if the error is due to missing CLI
		if filepath.Base(t.cliPath) == "claude" {
			return errors.NewCLINotFoundError(t.cliPath, err)
		}
		return errors.NewProcessError("failed to start CLI process", 0, "", err)
	}

	// Close stdin immediately since we're using --print mode
	// This prevents the CLI from waiting for interactive input
	if t.stdin != nil {
		t.stdin.Close()
	}

	return nil
}

// ReceiveMessages receives and parses messages from the CLI
func (t *SubprocessTransport) ReceiveMessages(ctx context.Context) (<-chan types.Message, error) {
	if t.cmd == nil {
		return nil, errors.NewCLIConnectionError("not connected", nil)
	}

	messageCh := make(chan types.Message, 10)

	go func() {
		defer close(messageCh)
		defer t.cleanup()

		// Process stdout messages
		scanner := bufio.NewScanner(t.stdout)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			// Check buffer size limit
			if len(line) > maxBufferSize {
				errorMsg := types.NewUserMessage(
					fmt.Sprintf("JSON message exceeded maximum buffer size of %d bytes", maxBufferSize))
				select {
				case messageCh <- errorMsg:
				case <-ctx.Done():
					return
				}
				continue
			}

			// Try to parse complete JSON
			var data map[string]any
			if err := json.Unmarshal([]byte(line), &data); err != nil {
				// Skip invalid JSON lines
				continue
			}

			// Handle control responses (skip for now as they're CLI internal)
			if messageType, ok := data["type"].(string); ok && messageType == "control_response" {
				continue
			}

			// Parse the message
			message, err := parser.ParseMessage(data)
			if err != nil {
				// Send parse error as a user message
				errorMsg := types.NewUserMessage(fmt.Sprintf("Parse error: %v", err))
				select {
				case messageCh <- errorMsg:
				case <-ctx.Done():
					return
				}
				continue
			}

			select {
			case messageCh <- message:
			case <-ctx.Done():
				return
			}
		}

		// Handle scanner error
		if err := scanner.Err(); err != nil && err != io.EOF {
			errorMsg := types.NewUserMessage(fmt.Sprintf("Scanner error: %v", err))
			select {
			case messageCh <- errorMsg:
			case <-ctx.Done():
				return
			}
		}

		// Process stderr and wait for command completion
		t.handleProcessCompletion(ctx, messageCh)
	}()

	return messageCh, nil
}

// Disconnect closes the subprocess
func (t *SubprocessTransport) Disconnect() error {
	return t.cleanup()
}

func (t *SubprocessTransport) buildCommand(options *types.QueryOptions, prompt string) []string {
	args := []string{"--output-format", "stream-json", "--verbose"}

	if options != nil {
		if options.SystemPrompt != "" {
			args = append(args, "--system-prompt", options.SystemPrompt)
		}

		if options.AppendSystemPrompt != "" {
			args = append(args, "--append-system-prompt", options.AppendSystemPrompt)
		}

		if len(options.AllowedTools) > 0 {
			args = append(args, "--allowedTools", strings.Join(options.AllowedTools, ","))
		}

		if options.MaxTurns > 0 {
			args = append(args, "--max-turns", fmt.Sprintf("%d", options.MaxTurns))
		}

		if len(options.DisallowedTools) > 0 {
			args = append(args, "--disallowedTools", strings.Join(options.DisallowedTools, ","))
		}

		if options.Model != "" {
			args = append(args, "--model", options.Model)
		}

		if options.PermissionPromptToolName != "" {
			args = append(args, "--permission-prompt-tool", options.PermissionPromptToolName)
		}

		if options.PermissionMode != "" {
			args = append(args, "--permission-mode", string(options.PermissionMode))
		}

		if options.ContinueConversation {
			args = append(args, "--continue")
		}

		if options.Resume != "" {
			args = append(args, "--resume", options.Resume)
		}

		if len(options.McpServers) > 0 {
			mcpConfig := map[string]any{
				"mcpServers": options.McpServers,
			}
			configJSON, _ := json.Marshal(mcpConfig)
			args = append(args, "--mcp-config", string(configJSON))
		}
	}

	// Add the prompt
	args = append(args, "--print", prompt)

	return args
}

func (t *SubprocessTransport) handleProcessCompletion(ctx context.Context, messageCh chan<- types.Message) {
	// Read stderr with safety limits
	var stderrLines []string
	var stderrSize int

	if t.stderr != nil {
		stderrChan := make(chan string, 100)
		go func() {
			defer close(stderrChan)
			scanner := bufio.NewScanner(t.stderr)
			for scanner.Scan() {
				select {
				case stderrChan <- scanner.Text():
				case <-ctx.Done():
					return
				}
			}
		}()

		// Collect stderr with timeout
		timeout := time.After(stderrTimeout)
		for {
			select {
			case line, ok := <-stderrChan:
				if !ok {
					goto waitForProcess
				}
				lineSize := len(line)
				if stderrSize+lineSize > maxStderrSize {
					stderrLines = append(stderrLines, fmt.Sprintf("[stderr truncated after %d bytes]", stderrSize))
					// Drain remaining lines
					for range stderrChan {
					}
					goto waitForProcess
				}
				stderrLines = append(stderrLines, line)
				stderrSize += lineSize
			case <-timeout:
				stderrLines = append(stderrLines, fmt.Sprintf("[stderr collection timed out after %v]", stderrTimeout))
				goto waitForProcess
			case <-ctx.Done():
				return
			}
		}
	}

waitForProcess:
	// Wait for process to complete
	var exitCode int
	if err := t.cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		} else {
			exitCode = -1
		}
	}

	stderrOutput := strings.Join(stderrLines, "\n")

	// Send error message if process failed
	if exitCode != 0 {
		errorMsg := types.NewUserMessage(
			fmt.Sprintf("Process failed with exit code %d: %s", exitCode, stderrOutput))
		select {
		case messageCh <- errorMsg:
		case <-ctx.Done():
		}
	}
}

func (t *SubprocessTransport) cleanup() error {
	var firstErr error

	if t.stdin != nil {
		if err := t.stdin.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		t.stdin = nil
	}

	if t.stdout != nil {
		if err := t.stdout.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		t.stdout = nil
	}

	if t.stderr != nil {
		if err := t.stderr.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		t.stderr = nil
	}

	if t.cmd != nil && t.cmd.Process != nil {
		// Try graceful termination first
		if err := t.cmd.Process.Signal(os.Interrupt); err == nil {
			// Wait a bit for graceful shutdown
			done := make(chan error, 1)
			go func() {
				done <- t.cmd.Wait()
			}()

			select {
			case <-done:
				// Process terminated gracefully
			case <-time.After(5 * time.Second):
				// Force kill
				t.cmd.Process.Kill()
				<-done
			}
		} else {
			// Force kill immediately
			t.cmd.Process.Kill()
			t.cmd.Wait()
		}
		t.cmd = nil
	}

	return firstErr
}

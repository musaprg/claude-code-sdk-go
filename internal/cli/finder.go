package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/musaprg/claude-code-sdk-go/internal/errors"
)

// FindCLI finds the Claude Code CLI binary
func FindCLI() (string, error) {
	// First check if "claude" is in PATH
	if cli, err := exec.LookPath("claude"); err == nil {
		return cli, nil
	}

	// Common installation locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.NewCLINotFoundError("", err)
	}

	locations := []string{
		filepath.Join(homeDir, ".npm-global", "bin", "claude"),
		"/usr/local/bin/claude",
		filepath.Join(homeDir, ".local", "bin", "claude"),
		filepath.Join(homeDir, "node_modules", ".bin", "claude"),
		filepath.Join(homeDir, ".yarn", "bin", "claude"),
	}

	for _, path := range locations {
		if fileExists(path) {
			return path, nil
		}
	}

	// Check if Node.js is available
	if _, err := exec.LookPath("node"); err != nil {
		return "", errors.NewCLINotFoundError("",
			fmt.Errorf("Claude Code requires Node.js, which is not installed.\n\n"+
				"Install Node.js from: https://nodejs.org/\n\n"+
				"After installing Node.js, install Claude Code:\n"+
				"  npm install -g @anthropic-ai/claude-code"))
	}

	return "", errors.NewCLINotFoundError("",
		fmt.Errorf("Claude Code not found. Install with:\n"+
			"  npm install -g @anthropic-ai/claude-code\n\n"+
			"If already installed locally, try:\n"+
			"  export PATH=\"$HOME/node_modules/.bin:$PATH\"\n\n"+
			"Or specify the path when creating transport:\n"+
			"  NewClient(&ClientOptions{CLIPath: \"/path/to/claude\"})"))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

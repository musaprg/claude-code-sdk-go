package claudecode_test

import (
	"context"
	"fmt"
	"log"

	claudecode "github.com/musaprg/claude-code-sdk-go"
)

func ExampleNewClient() {
	// Create a client with default settings
	client := claudecode.NewClient(nil)
	_ = client
	fmt.Println("Client created with default settings")

	// Create a client with custom options
	client = claudecode.NewClient(&claudecode.ClientOptions{
		CLIPath: "/usr/local/bin/claude",
		CWD:     "/path/to/project",
	})
	_ = client
	fmt.Println("Client created with custom options")

	// Output:
	// Client created with default settings
	// Client created with custom options
}

func ExampleClient_Query() {
	client := claudecode.NewClient(nil)
	ctx := context.Background()

	options := &claudecode.QueryOptions{
		Model:          "sonnet",
		PermissionMode: claudecode.PermissionModeAcceptEdits,
	}

	messages, err := client.Query(ctx, "Explain this code", options)
	if err != nil {
		log.Fatal(err)
	}

	for message := range messages {
		switch msg := message.(type) {
		case *claudecode.AssistantMessage:
			// Handle assistant response
			for _, block := range msg.Content {
				if textBlock, ok := block.(*claudecode.TextBlock); ok {
					fmt.Print(textBlock.Text)
				}
			}
		case *claudecode.ResultMessage:
			// Handle final result with conversation metadata
			fmt.Printf("Conversation completed in %dms\n", msg.DurationMs)
		}
	}
}

func ExampleQuery() {
	ctx := context.Background()

	messages, err := claudecode.Query(ctx, "Write a hello world program", nil)
	if err != nil {
		log.Fatal(err)
	}

	for message := range messages {
		if assistantMsg, ok := message.(*claudecode.AssistantMessage); ok {
			for _, block := range assistantMsg.Content {
				if textBlock, ok := block.(*claudecode.TextBlock); ok {
					fmt.Print(textBlock.Text)
				}
			}
		}
	}
}

func ExampleQueryOptions() {
	ctx := context.Background()

	// Configure query with specific options
	options := &claudecode.QueryOptions{
		Model:                "sonnet",
		PermissionMode:       claudecode.PermissionModeBypassPermissions,
		MaxTurns:             5,
		AllowedTools:         []string{"Read", "Edit", "Bash"},
		DisallowedTools:      []string{"WebFetch"},
		SystemPrompt:         "You are a helpful coding assistant.",
		AppendSystemPrompt:   "Always write clean, well-documented code.",
		ContinueConversation: false,
		MaxThinkingTokens:    1000,
		CWD:                  "/workspace",
		McpTools:             []string{"github", "jira"},
		McpServers: map[string]claudecode.McpServerConfig{
			"github": {
				Type:    "stdio",
				Command: "mcp-github",
				Args:    []string{"--token", "ghp_xxx"},
			},
		},
	}

	messages, err := claudecode.Query(ctx, "Review this pull request", options)
	if err != nil {
		log.Fatal(err)
	}

	for message := range messages {
		// Process messages
		_ = message
	}
}

// streaming_mode.go demonstrates advanced streaming patterns with the Claude Code SDK
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	claudecode "github.com/musaprg/claude-code-sdk-go"
)

func main() {
	// Example 1: Basic streaming
	fmt.Println("=== Basic Streaming Example ===")
	basicStreamingExample()

	// Example 2: Multi-turn conversation
	fmt.Println("\n=== Multi-turn Conversation Example ===")
	multiTurnExample()

	// Example 3: Streaming with custom options
	fmt.Println("\n=== Streaming with Custom Options ===")
	customOptionsExample()

	// Example 4: Context cancellation and timeout
	fmt.Println("\n=== Context Cancellation Example ===")
	timeoutExample()
}

// basicStreamingExample demonstrates the basic streaming pattern
func basicStreamingExample() {
	ctx := context.Background()

	// Simple query with streaming response
	messageCh, err := claudecode.Query(ctx, "What is 2+2?", nil)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	// Process streaming messages
	for message := range messageCh {
		displayMessage(message)
	}
}

// multiTurnExample demonstrates a multi-turn conversation
func multiTurnExample() {
	ctx := context.Background()

	// Create client with custom working directory
	client := claudecode.NewClient(&claudecode.ClientOptions{
		CWD: ".",
	})

	// Configure options for multi-turn conversation
	options := &claudecode.QueryOptions{
		SystemPrompt: "You are a helpful coding assistant. Keep responses concise.",
		MaxTurns:     3,
	}

	// First query
	messageCh, err := client.Query(ctx, "Explain what Go channels are", options)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	fmt.Println("Processing first query...")
	for message := range messageCh {
		displayMessage(message)
	}

	// Continue conversation
	options.ContinueConversation = true
	messageCh2, err := client.Query(ctx, "Can you give me a simple example?", options)
	if err != nil {
		log.Fatal("Follow-up query failed:", err)
	}

	fmt.Println("\nProcessing follow-up query...")
	for message := range messageCh2 {
		displayMessage(message)
	}
}

// customOptionsExample demonstrates advanced configuration
func customOptionsExample() {
	ctx := context.Background()

	// Create client with extensive configuration
	options := &claudecode.QueryOptions{
		SystemPrompt:   "You are an expert Go developer.",
		AllowedTools:   []string{"Read", "Write", "Bash"},
		MaxTurns:       5,
		PermissionMode: claudecode.PermissionModeAcceptEdits,
		Model:          "claude-3-5-sonnet-20241022",
	}

	messageCh, err := claudecode.Query(ctx, "Help me understand Go error handling patterns", options)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	// Process messages with detailed information
	messageCount := 0
	for message := range messageCh {
		messageCount++
		fmt.Printf("Message %d: ", messageCount)
		displayMessageDetailed(message)
	}

	fmt.Printf("Total messages received: %d\n", messageCount)
}

// timeoutExample demonstrates context cancellation and timeout handling
func timeoutExample() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	messageCh, err := claudecode.Query(ctx, "Write a detailed explanation of Go's memory model", nil)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	// Process messages until timeout or completion
	startTime := time.Now()
	for message := range messageCh {
		elapsed := time.Since(startTime)
		fmt.Printf("[%v] ", elapsed.Round(time.Millisecond))
		displayMessage(message)

		// Check if context was cancelled
		select {
		case <-ctx.Done():
			fmt.Printf("Context cancelled: %v\n", ctx.Err())
			return
		default:
		}
	}

	fmt.Printf("Completed in: %v\n", time.Since(startTime).Round(time.Millisecond))
}

// displayMessage provides basic message display
func displayMessage(message claudecode.Message) {
	switch msg := message.(type) {
	case *claudecode.UserMessage:
		fmt.Printf("User: %s\n", msg.Content)

	case *claudecode.AssistantMessage:
		fmt.Print("Claude: ")
		for _, block := range msg.Content {
			switch b := block.(type) {
			case *claudecode.TextBlock:
				fmt.Print(b.Text)
			case *claudecode.ToolUseBlock:
				fmt.Printf("[Using tool: %s]", b.Name)
			case *claudecode.ToolResultBlock:
				fmt.Printf("[Tool result for %s]", b.ToolUseID)
			}
		}
		fmt.Println()

	case *claudecode.SystemMessage:
		fmt.Printf("System [%s]: %v\n", msg.Subtype, msg.Data)

	case *claudecode.ResultMessage:
		fmt.Printf("Result: %s ", msg.Subtype)
		if msg.TotalCostUSD != nil {
			fmt.Printf("(Cost: $%.4f) ", *msg.TotalCostUSD)
		}
		fmt.Printf("(Duration: %dms)\n", msg.DurationMs)
	}
}

// displayMessageDetailed provides detailed message information
func displayMessageDetailed(message claudecode.Message) {
	switch msg := message.(type) {
	case *claudecode.UserMessage:
		fmt.Printf("UserMessage: %s\n", msg.Content)

	case *claudecode.AssistantMessage:
		fmt.Printf("AssistantMessage with %d content blocks:\n", len(msg.Content))
		for i, block := range msg.Content {
			fmt.Printf("  Block %d: ", i+1)
			switch b := block.(type) {
			case *claudecode.TextBlock:
				// Truncate long text for readability
				text := b.Text
				if len(text) > 100 {
					text = text[:100] + "..."
				}
				fmt.Printf("TextBlock: %s\n", text)
			case *claudecode.ToolUseBlock:
				fmt.Printf("ToolUseBlock: %s (ID: %s)\n", b.Name, b.ID)
			case *claudecode.ToolResultBlock:
				fmt.Printf("ToolResultBlock for %s\n", b.ToolUseID)
			}
		}

	case *claudecode.SystemMessage:
		fmt.Printf("SystemMessage [%s]: %v\n", msg.Subtype, msg.Data)

	case *claudecode.ResultMessage:
		fmt.Printf("ResultMessage: %s\n", msg.Subtype)
		fmt.Printf("  Duration: %dms (API: %dms)\n", msg.DurationMs, msg.DurationAPIMs)
		fmt.Printf("  Turns: %d, Error: %v\n", msg.NumTurns, msg.IsError)
		if msg.TotalCostUSD != nil {
			fmt.Printf("  Cost: $%.4f\n", *msg.TotalCostUSD)
		}
		if msg.Usage != nil {
			fmt.Printf("  Usage: %v\n", msg.Usage)
		}
	}
}

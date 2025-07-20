package main

import (
	"context"
	"fmt"
	"log"

	claudecode "github.com/musaprg/claude-code-sdk-go"
)

func main() {
	ctx := context.Background()

	// Example 1: Simple query
	fmt.Println("=== Simple Query ===")
	simpleQuery(ctx)

	fmt.Println("\n=== Client with Options ===")
	// Example 2: Client with options
	clientWithOptions(ctx)
}

func simpleQuery(ctx context.Context) {
	messageCh, err := claudecode.Query(ctx, "What is the capital of France?", nil)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	for message := range messageCh {
		handleMessage(message)
	}
}

func clientWithOptions(ctx context.Context) {
	// Create client with custom working directory
	client := claudecode.NewClient(&claudecode.ClientOptions{
		CWD: ".", // Use current directory
	})

	// Configure query options
	options := &claudecode.QueryOptions{
		SystemPrompt:   "You are a helpful assistant that provides concise answers.",
		MaxTurns:       3,
		PermissionMode: claudecode.PermissionModeDefault,
	}

	messageCh, err := client.Query(ctx, "Explain what Go channels are in 2 sentences.", options)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	for message := range messageCh {
		handleMessage(message)
	}
}

func handleMessage(message claudecode.Message) {
	switch msg := message.(type) {
	case *claudecode.UserMessage:
		fmt.Printf("👤 User: %s\n", msg.Content)

	case *claudecode.AssistantMessage:
		fmt.Print("🤖 Claude: ")
		for i, block := range msg.Content {
			if i > 0 {
				fmt.Print(" ")
			}

			switch b := block.(type) {
			case *claudecode.TextBlock:
				fmt.Print(b.Text)
			case *claudecode.ToolUseBlock:
				fmt.Printf("[Using tool: %s]", b.Name)
			case *claudecode.ToolResultBlock:
				fmt.Printf("[Tool result: %v]", b.Content)
			}
		}
		fmt.Println()

	case *claudecode.SystemMessage:
		fmt.Printf("⚙️  System [%s]: %v\n", msg.Subtype, msg.Data)

	case *claudecode.ResultMessage:
		fmt.Printf("📊 Result: %s (Duration: %dms", msg.Subtype, msg.DurationMs)
		if msg.TotalCostUSD != nil {
			fmt.Printf(", Cost: $%.4f", *msg.TotalCostUSD)
		}
		fmt.Println(")")
	}
}

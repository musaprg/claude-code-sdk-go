package claudecode

import (
	"testing"
)

func TestMessageTypes(t *testing.T) {
	tests := []struct {
		name     string
		message  Message
		expected MessageType
	}{
		{
			name:     "UserMessage",
			message:  NewUserMessage("test content"),
			expected: MessageTypeUser,
		},
		{
			name:     "AssistantMessage",
			message:  NewAssistantMessage([]ContentBlock{NewTextBlock("test")}),
			expected: MessageTypeAssistant,
		},
		{
			name:     "SystemMessage",
			message:  NewSystemMessage("test", map[string]any{"key": "value"}),
			expected: MessageTypeSystem,
		},
		{
			name:     "ResultMessage",
			message:  NewResultMessage("completed", 1000, 800, 1, false, "session1"),
			expected: MessageTypeResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.message.Type(); got != tt.expected {
				t.Errorf("Message.Type() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContentBlockTypes(t *testing.T) {
	tests := []struct {
		name     string
		block    ContentBlock
		expected ContentBlockType
	}{
		{
			name:     "TextBlock",
			block:    NewTextBlock("test"),
			expected: ContentBlockTypeText,
		},
		{
			name:     "ToolUseBlock",
			block:    NewToolUseBlock("id1", "tool1", map[string]any{"param": "value"}),
			expected: ContentBlockTypeToolUse,
		},
		{
			name:     "ToolResultBlock",
			block:    NewToolResultBlock("id1", "result", nil),
			expected: ContentBlockTypeToolResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.block.BlockType(); got != tt.expected {
				t.Errorf("ContentBlock.BlockType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUserMessage(t *testing.T) {
	content := "Hello, world!"
	message := NewUserMessage(content)

	if message.Content != content {
		t.Errorf("UserMessage.Content = %q, want %q", message.Content, content)
	}

	if message.Type() != MessageTypeUser {
		t.Errorf("UserMessage.Type() = %v, want %v", message.Type(), MessageTypeUser)
	}
}

func TestAssistantMessage(t *testing.T) {
	textBlock := NewTextBlock("Hello")
	toolBlock := NewToolUseBlock("id1", "tool1", map[string]any{"param": "value"})
	content := []ContentBlock{textBlock, toolBlock}

	message := NewAssistantMessage(content)

	if len(message.Content) != 2 {
		t.Errorf("AssistantMessage.Content length = %d, want %d", len(message.Content), 2)
	}

	if message.Content[0] != textBlock {
		t.Errorf("AssistantMessage.Content[0] = %v, want %v", message.Content[0], textBlock)
	}

	if message.Content[1] != toolBlock {
		t.Errorf("AssistantMessage.Content[1] = %v, want %v", message.Content[1], toolBlock)
	}

	if message.Type() != MessageTypeAssistant {
		t.Errorf("AssistantMessage.Type() = %v, want %v", message.Type(), MessageTypeAssistant)
	}
}

func TestSystemMessage(t *testing.T) {
	subtype := "info"
	data := map[string]any{"key": "value", "number": 42}

	message := NewSystemMessage(subtype, data)

	if message.Subtype != subtype {
		t.Errorf("SystemMessage.Subtype = %q, want %q", message.Subtype, subtype)
	}

	if message.Data["key"] != "value" {
		t.Errorf("SystemMessage.Data[\"key\"] = %v, want %v", message.Data["key"], "value")
	}

	if message.Type() != MessageTypeSystem {
		t.Errorf("SystemMessage.Type() = %v, want %v", message.Type(), MessageTypeSystem)
	}
}

func TestResultMessage(t *testing.T) {
	subtype := "completed"
	durationMs := 1000
	durationAPIMs := 800
	numTurns := 1
	isError := false
	sessionID := "session123"

	message := NewResultMessage(subtype, durationMs, durationAPIMs, numTurns, isError, sessionID)

	if message.Subtype != subtype {
		t.Errorf("ResultMessage.Subtype = %q, want %q", message.Subtype, subtype)
	}
	if message.DurationMs != durationMs {
		t.Errorf("ResultMessage.DurationMs = %d, want %d", message.DurationMs, durationMs)
	}
	if message.DurationAPIMs != durationAPIMs {
		t.Errorf("ResultMessage.DurationAPIMs = %d, want %d", message.DurationAPIMs, durationAPIMs)
	}
	if message.IsError != isError {
		t.Errorf("ResultMessage.IsError = %v, want %v", message.IsError, isError)
	}
	if message.NumTurns != numTurns {
		t.Errorf("ResultMessage.NumTurns = %d, want %d", message.NumTurns, numTurns)
	}
	if message.SessionID != sessionID {
		t.Errorf("ResultMessage.SessionID = %q, want %q", message.SessionID, sessionID)
	}

	if message.Type() != MessageTypeResult {
		t.Errorf("ResultMessage.Type() = %v, want %v", message.Type(), MessageTypeResult)
	}
}

func TestTextBlock(t *testing.T) {
	text := "Hello, world!"
	block := NewTextBlock(text)

	if block.Text != text {
		t.Errorf("TextBlock.Text = %q, want %q", block.Text, text)
	}

	if block.BlockType() != ContentBlockTypeText {
		t.Errorf("TextBlock.BlockType() = %v, want %v", block.BlockType(), ContentBlockTypeText)
	}
}

func TestToolUseBlock(t *testing.T) {
	id := "tool-123"
	name := "calculator"
	input := map[string]any{"operation": "add", "a": 1, "b": 2}

	block := NewToolUseBlock(id, name, input)

	if block.ID != id {
		t.Errorf("ToolUseBlock.ID = %q, want %q", block.ID, id)
	}
	if block.Name != name {
		t.Errorf("ToolUseBlock.Name = %q, want %q", block.Name, name)
	}
	if block.Input["operation"] != "add" {
		t.Errorf("ToolUseBlock.Input[\"operation\"] = %v, want %v", block.Input["operation"], "add")
	}
	if block.BlockType() != ContentBlockTypeToolUse {
		t.Errorf("ToolUseBlock.BlockType() = %v, want %v", block.BlockType(), ContentBlockTypeToolUse)
	}
}

func TestToolResultBlock(t *testing.T) {
	toolUseID := "tool-123"
	content := "result content"
	isError := false

	block := NewToolResultBlock(toolUseID, content, &isError)

	if block.ToolUseID != toolUseID {
		t.Errorf("ToolResultBlock.ToolUseID = %q, want %q", block.ToolUseID, toolUseID)
	}
	if block.Content != content {
		t.Errorf("ToolResultBlock.Content = %v, want %v", block.Content, content)
	}
	if block.IsError == nil || *block.IsError != isError {
		t.Errorf("ToolResultBlock.IsError = %v, want %v", block.IsError, &isError)
	}
	if block.BlockType() != ContentBlockTypeToolResult {
		t.Errorf("ToolResultBlock.BlockType() = %v, want %v", block.BlockType(), ContentBlockTypeToolResult)
	}
}

func TestQueryOptions(t *testing.T) {
	options := &QueryOptions{
		SystemPrompt:   "You are a helpful assistant",
		MaxTurns:       10,
		AllowedTools:   []string{"calculator", "web_search"},
		CWD:            "/tmp",
		PermissionMode: PermissionModeDefault,
		McpServers: map[string]McpServerConfig{
			"server1": {
				Type:    "stdio",
				Command: "mcp-server",
				Args:    []string{"--port", "8080"},
				Env:     map[string]string{"ENV": "prod"},
			},
		},
	}

	if options.SystemPrompt != "You are a helpful assistant" {
		t.Errorf("QueryOptions.SystemPrompt = %q, want %q", options.SystemPrompt, "You are a helpful assistant")
	}
	if options.MaxTurns != 10 {
		t.Errorf("QueryOptions.MaxTurns = %d, want %d", options.MaxTurns, 10)
	}
	if len(options.AllowedTools) != 2 {
		t.Errorf("QueryOptions.AllowedTools length = %d, want %d", len(options.AllowedTools), 2)
	}
	if options.CWD != "/tmp" {
		t.Errorf("QueryOptions.CWD = %q, want %q", options.CWD, "/tmp")
	}
	if options.PermissionMode != PermissionModeDefault {
		t.Errorf("QueryOptions.PermissionMode = %v, want %v", options.PermissionMode, PermissionModeDefault)
	}
	if len(options.McpServers) != 1 {
		t.Errorf("QueryOptions.McpServers length = %d, want %d", len(options.McpServers), 1)
	}
}

func TestClientOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		client := NewClient(nil)
		if client.cliPath != "" {
			t.Errorf("Client.cliPath = %q, want empty string", client.cliPath)
		}
		if client.cwd != "" {
			t.Errorf("Client.cwd = %q, want empty string", client.cwd)
		}
	})

	t.Run("custom options", func(t *testing.T) {
		customPath := "/usr/local/bin/claude"
		customCWD := "/home/user/project"
		client := NewClient(&ClientOptions{
			CLIPath: customPath,
			CWD:     customCWD,
		})
		if client.cliPath != customPath {
			t.Errorf("Client.cliPath = %q, want %q", client.cliPath, customPath)
		}
		if client.cwd != customCWD {
			t.Errorf("Client.cwd = %q, want %q", client.cwd, customCWD)
		}
	})
}

package claudecode

import "github.com/musaprg/claude-code-sdk-go/internal/types"

// Re-export types from internal package for public API access.
// These type aliases provide a clean interface while keeping implementation details internal.
type (
	// MessageType represents the type of message in a Claude Code conversation.
	MessageType = types.MessageType
	// ContentBlockType represents the type of content block within a message.
	ContentBlockType = types.ContentBlockType
	// PermissionMode defines how tools are permitted to run during a Claude Code session.
	PermissionMode = types.PermissionMode
	// Message represents a message in the Claude Code conversation.
	Message = types.Message
	// ContentBlock represents a content block within a message.
	ContentBlock = types.ContentBlock
	// UserMessage represents a message from the human user to Claude.
	UserMessage = types.UserMessage
	// AssistantMessage represents a message from Claude's AI assistant.
	AssistantMessage = types.AssistantMessage
	// SystemMessage represents internal system messages and metadata from Claude Code.
	SystemMessage = types.SystemMessage
	// ResultMessage represents the final result message containing conversation metadata.
	ResultMessage = types.ResultMessage
	// TextBlock represents a plain text content block within a message.
	TextBlock = types.TextBlock
	// ToolUseBlock represents a tool invocation by the assistant.
	ToolUseBlock = types.ToolUseBlock
	// ToolResultBlock represents the result of a tool execution.
	ToolResultBlock = types.ToolResultBlock
	// McpServerConfig represents configuration for a Model Context Protocol (MCP) server.
	McpServerConfig = types.McpServerConfig
	// QueryOptions contains configuration options for Claude Code queries.
	QueryOptions = types.QueryOptions
	// ClientOptions contains configuration options for creating a new Claude Code SDK client.
	ClientOptions = types.ClientOptions
)

// Re-export constants from internal package.
// These constants define the available message types, content block types, and permission modes.
const (
	// MessageTypeUser represents a message from the human user.
	MessageTypeUser = types.MessageTypeUser
	// MessageTypeAssistant represents a message from Claude's AI assistant.
	MessageTypeAssistant = types.MessageTypeAssistant
	// MessageTypeSystem represents internal system messages and metadata.
	MessageTypeSystem = types.MessageTypeSystem
	// MessageTypeResult represents final result messages with conversation metadata.
	MessageTypeResult = types.MessageTypeResult

	// ContentBlockTypeText represents plain text content.
	ContentBlockTypeText = types.ContentBlockTypeText
	// ContentBlockTypeToolUse represents a tool invocation by the assistant.
	ContentBlockTypeToolUse = types.ContentBlockTypeToolUse
	// ContentBlockTypeToolResult represents the result of a tool execution.
	ContentBlockTypeToolResult = types.ContentBlockTypeToolResult

	// PermissionModeDefault uses the standard permission prompts for tool usage.
	PermissionModeDefault = types.PermissionModeDefault
	// PermissionModeAcceptEdits automatically accepts file edit operations without prompting.
	PermissionModeAcceptEdits = types.PermissionModeAcceptEdits
	// PermissionModeBypassPermissions bypasses all permission checks (recommended only for sandboxes).
	PermissionModeBypassPermissions = types.PermissionModeBypassPermissions
)

// Re-export constructor functions from internal package.
// These functions provide convenient ways to create message and content block instances.
var (
	// NewUserMessage creates a new UserMessage with the given content.
	NewUserMessage = types.NewUserMessage
	// NewAssistantMessage creates a new AssistantMessage with the given content blocks.
	NewAssistantMessage = types.NewAssistantMessage
	// NewSystemMessage creates a new SystemMessage with the given subtype and data.
	NewSystemMessage = types.NewSystemMessage
	// NewResultMessage creates a new ResultMessage with the given parameters.
	NewResultMessage = types.NewResultMessage
	// NewTextBlock creates a new TextBlock with the given text content.
	NewTextBlock = types.NewTextBlock
	// NewToolUseBlock creates a new ToolUseBlock with the given parameters.
	NewToolUseBlock = types.NewToolUseBlock
	// NewToolResultBlock creates a new ToolResultBlock with the given parameters.
	NewToolResultBlock = types.NewToolResultBlock
)

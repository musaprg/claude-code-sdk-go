package types

// MessageType represents the type of message in a Claude Code conversation.
type MessageType string

const (
	// MessageTypeUser represents a message from the human user.
	MessageTypeUser MessageType = "user"
	// MessageTypeAssistant represents a message from Claude's AI assistant.
	MessageTypeAssistant MessageType = "assistant"
	// MessageTypeSystem represents internal system messages and metadata.
	MessageTypeSystem MessageType = "system"
	// MessageTypeResult represents final result messages with conversation metadata.
	MessageTypeResult MessageType = "result"
)

// ContentBlockType represents the type of content block within a message.
type ContentBlockType string

const (
	// ContentBlockTypeText represents plain text content.
	ContentBlockTypeText ContentBlockType = "text"
	// ContentBlockTypeToolUse represents a tool invocation by the assistant.
	ContentBlockTypeToolUse ContentBlockType = "tool_use"
	// ContentBlockTypeToolResult represents the result of a tool execution.
	ContentBlockTypeToolResult ContentBlockType = "tool_result"
)

// PermissionMode defines how tools are permitted to run during a Claude Code session.
type PermissionMode string

const (
	// PermissionModeDefault uses the standard permission prompts for tool usage.
	PermissionModeDefault PermissionMode = "default"
	// PermissionModeAcceptEdits automatically accepts file edit operations without prompting.
	PermissionModeAcceptEdits PermissionMode = "acceptEdits"
	// PermissionModeBypassPermissions bypasses all permission checks (recommended only for sandboxes).
	PermissionModeBypassPermissions PermissionMode = "bypassPermissions"
)

// Message represents a message in the Claude Code conversation.
// All message types implement this interface to provide polymorphic handling.
type Message interface {
	// Type returns the MessageType of this message.
	Type() MessageType
}

// ContentBlock represents a content block within a message.
// Content blocks can contain text, tool uses, or tool results.
type ContentBlock interface {
	// BlockType returns the ContentBlockType of this block.
	BlockType() ContentBlockType
}

// UserMessage represents a message from the human user to Claude.
type UserMessage struct {
	// Content contains the user's prompt or question text.
	Content string `json:"content"`
}

func (m *UserMessage) Type() MessageType {
	return MessageTypeUser
}

func NewUserMessage(content string) *UserMessage {
	return &UserMessage{Content: content}
}

// AssistantMessage represents a message from Claude's AI assistant.
type AssistantMessage struct {
	// Content contains the assistant's response as a sequence of content blocks,
	// which may include text, tool uses, and tool results.
	Content []ContentBlock `json:"content"`
}

func (m *AssistantMessage) Type() MessageType {
	return MessageTypeAssistant
}

func NewAssistantMessage(content []ContentBlock) *AssistantMessage {
	return &AssistantMessage{Content: content}
}

// SystemMessage represents internal system messages and metadata from Claude Code.
type SystemMessage struct {
	// Subtype specifies the kind of system message (e.g., "thinking", "progress").
	Subtype string `json:"subtype"`
	// Data contains arbitrary metadata associated with the system message.
	Data map[string]any `json:"data"`
}

func (m *SystemMessage) Type() MessageType {
	return MessageTypeSystem
}

func NewSystemMessage(subtype string, data map[string]any) *SystemMessage {
	return &SystemMessage{Subtype: subtype, Data: data}
}

// ResultMessage represents the final result message containing conversation metadata and statistics.
type ResultMessage struct {
	// Subtype specifies the kind of result message.
	Subtype string `json:"subtype"`
	// DurationMs is the total conversation duration in milliseconds.
	DurationMs int `json:"duration_ms"`
	// DurationAPIMs is the API processing time in milliseconds.
	DurationAPIMs int `json:"duration_api_ms"`
	// IsError indicates whether the conversation ended with an error.
	IsError bool `json:"is_error"`
	// NumTurns is the number of conversation turns that occurred.
	NumTurns int `json:"num_turns"`
	// SessionID is the unique identifier for this conversation session.
	SessionID string `json:"session_id"`
	// TotalCostUSD is the total cost of the conversation in USD, if available.
	TotalCostUSD *float64 `json:"total_cost_usd,omitempty"`
	// Usage contains token usage statistics and other metrics.
	Usage map[string]any `json:"usage,omitempty"`
	// Result contains the final result text, if any.
	Result *string `json:"result,omitempty"`
}

func (m *ResultMessage) Type() MessageType {
	return MessageTypeResult
}

func NewResultMessage(subtype string, durationMs, durationAPIMs, numTurns int, isError bool, sessionID string) *ResultMessage {
	return &ResultMessage{
		Subtype:       subtype,
		DurationMs:    durationMs,
		DurationAPIMs: durationAPIMs,
		IsError:       isError,
		NumTurns:      numTurns,
		SessionID:     sessionID,
	}
}

// TextBlock represents a plain text content block within a message.
type TextBlock struct {
	// Text contains the actual text content.
	Text string `json:"text"`
}

func (b *TextBlock) BlockType() ContentBlockType {
	return ContentBlockTypeText
}

func NewTextBlock(text string) *TextBlock {
	return &TextBlock{Text: text}
}

// ToolUseBlock represents a tool invocation by the assistant.
type ToolUseBlock struct {
	// ID is the unique identifier for this tool use.
	ID string `json:"id"`
	// Name is the name of the tool being invoked (e.g., "Bash", "Edit", "Read").
	Name string `json:"name"`
	// Input contains the parameters passed to the tool.
	Input map[string]any `json:"input"`
}

func (b *ToolUseBlock) BlockType() ContentBlockType {
	return ContentBlockTypeToolUse
}

func NewToolUseBlock(id, name string, input map[string]any) *ToolUseBlock {
	return &ToolUseBlock{
		ID:    id,
		Name:  name,
		Input: input,
	}
}

// ToolResultBlock represents the result of a tool execution.
type ToolResultBlock struct {
	// ToolUseID is the ID of the corresponding ToolUseBlock.
	ToolUseID string `json:"tool_use_id"`
	// Content contains the tool's output or result data.
	Content any `json:"content,omitempty"`
	// IsError indicates whether the tool execution resulted in an error.
	IsError *bool `json:"is_error,omitempty"`
}

func (b *ToolResultBlock) BlockType() ContentBlockType {
	return ContentBlockTypeToolResult
}

func NewToolResultBlock(toolUseID string, content any, isError *bool) *ToolResultBlock {
	return &ToolResultBlock{
		ToolUseID: toolUseID,
		Content:   content,
		IsError:   isError,
	}
}

// McpServerConfig represents configuration for a Model Context Protocol (MCP) server.
// MCP servers extend Claude Code's capabilities with additional tools and resources.
type McpServerConfig struct {
	// Type specifies the connection type: "stdio", "sse", or "http".
	Type string `json:"type,omitempty"`
	// Command is the executable command for stdio-type servers.
	Command string `json:"command,omitempty"`
	// Args contains command-line arguments for the MCP server executable.
	Args []string `json:"args,omitempty"`
	// Env contains environment variables to set for the MCP server process.
	Env map[string]string `json:"env,omitempty"`
	// URL is the endpoint URL for sse or http-type servers.
	URL string `json:"url,omitempty"`
	// Headers contains HTTP headers for sse or http-type servers.
	Headers map[string]string `json:"headers,omitempty"`
}

// QueryOptions contains configuration options for Claude Code queries.
type QueryOptions struct {
	// AllowedTools is a list of tool names that are explicitly permitted for use.
	// If specified, only these tools will be available to the assistant.
	AllowedTools []string `json:"allowed_tools,omitempty"`
	// MaxThinkingTokens limits the number of tokens Claude can use for internal reasoning.
	MaxThinkingTokens int `json:"max_thinking_tokens,omitempty"`
	// SystemPrompt overrides the default system prompt that defines Claude's behavior.
	SystemPrompt string `json:"system_prompt,omitempty"`
	// AppendSystemPrompt adds additional instructions to the default system prompt.
	AppendSystemPrompt string `json:"append_system_prompt,omitempty"`
	// McpTools is a list of MCP tool names to make available for this query.
	McpTools []string `json:"mcp_tools,omitempty"`
	// McpServers contains MCP server configurations to use for this query.
	McpServers map[string]McpServerConfig `json:"mcp_servers,omitempty"`
	// PermissionMode controls how tool permissions are handled during the session.
	PermissionMode PermissionMode `json:"permission_mode,omitempty"`
	// ContinueConversation continues the most recent conversation if true.
	ContinueConversation bool `json:"continue_conversation,omitempty"`
	// Resume specifies a session ID to resume a previous conversation.
	Resume string `json:"resume,omitempty"`
	// MaxTurns limits the maximum number of conversation turns.
	MaxTurns int `json:"max_turns,omitempty"`
	// DisallowedTools is a list of tool names that are explicitly prohibited.
	DisallowedTools []string `json:"disallowed_tools,omitempty"`
	// Model specifies which Claude model to use (e.g., "sonnet", "opus").
	Model string `json:"model,omitempty"`
	// PermissionPromptToolName customizes the tool name shown in permission prompts.
	PermissionPromptToolName string `json:"permission_prompt_tool_name,omitempty"`
	// CWD sets the current working directory for the Claude Code session.
	CWD string `json:"cwd,omitempty"`
}

// ClientOptions contains configuration options for creating a new Claude Code SDK client.
type ClientOptions struct {
	// CLIPath specifies a custom path to the Claude Code CLI executable.
	// If empty, the SDK will search for the CLI in standard locations.
	CLIPath string
	// CWD sets the current working directory for all operations.
	// If empty, the current process working directory is used.
	CWD string
}

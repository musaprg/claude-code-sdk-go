package parser

import (
	"fmt"

	"github.com/musaprg/claude-code-sdk-go/internal/errors"
	"github.com/musaprg/claude-code-sdk-go/internal/types"
)

// ParseMessage parses message from CLI output into typed Message objects
func ParseMessage(data map[string]any) (types.Message, error) {
	if data == nil {
		return nil, errors.NewMessageParseError(
			"Invalid message data type (expected map, got nil)", data, nil)
	}

	messageType, ok := data["type"].(string)
	if !ok {
		return nil, errors.NewMessageParseError("Message missing 'type' field", data, nil)
	}

	switch messageType {
	case "user":
		return parseUserMessage(data)
	case "assistant":
		return parseAssistantMessage(data)
	case "system":
		return parseSystemMessage(data)
	case "result":
		return parseResultMessage(data)
	default:
		return nil, errors.NewMessageParseError(
			fmt.Sprintf("Unknown message type: %s", messageType), data, nil)
	}
}

func parseUserMessage(data map[string]any) (*types.UserMessage, error) {
	// Try the nested format first: data["message"]["content"]
	if message, ok := data["message"].(map[string]any); ok {
		// Try string content first (typical user message)
		if content, ok := message["content"].(string); ok {
			return types.NewUserMessage(content), nil
		}

		// Try array content (tool result messages)
		if contentArray, ok := message["content"].([]any); ok {
			// For tool result messages, extract the content from the first tool_result block
			for _, blockData := range contentArray {
				if blockMap, ok := blockData.(map[string]any); ok {
					if blockType, ok := blockMap["type"].(string); ok && blockType == "tool_result" {
						if toolContent, ok := blockMap["content"].(string); ok {
							return types.NewUserMessage(toolContent), nil
						}
					}
				}
			}
			// If we can't extract tool result content, use a placeholder
			return types.NewUserMessage("[Tool result message]"), nil
		}
	}

	// Try the direct format: data["content"]
	if content, ok := data["content"].(string); ok {
		return types.NewUserMessage(content), nil
	}

	// If both formats fail, return error
	return nil, errors.NewMessageParseError("Missing required field 'content' in user message", data, nil)
}

func parseAssistantMessage(data map[string]any) (*types.AssistantMessage, error) {
	var contentData []any

	// Try the nested format first: data["message"]["content"]
	if message, ok := data["message"].(map[string]any); ok {
		if content, ok := message["content"].([]any); ok {
			contentData = content
		}
	}

	// Try the direct format: data["content"]
	if contentData == nil {
		if content, ok := data["content"].([]any); ok {
			contentData = content
		}
	}

	// If both formats fail, return error
	if contentData == nil {
		return nil, errors.NewMessageParseError("Missing required field 'content' in assistant message", data, nil)
	}

	var contentBlocks []types.ContentBlock
	for _, blockData := range contentData {
		blockMap, ok := blockData.(map[string]any)
		if !ok {
			return nil, errors.NewMessageParseError("Invalid content block format", data, nil)
		}

		block, err := parseContentBlock(blockMap)
		if err != nil {
			return nil, err
		}
		contentBlocks = append(contentBlocks, block)
	}

	return types.NewAssistantMessage(contentBlocks), nil
}

func parseContentBlock(blockData map[string]any) (types.ContentBlock, error) {
	blockType, ok := blockData["type"].(string)
	if !ok {
		return nil, errors.NewMessageParseError("Content block missing 'type' field", blockData, nil)
	}

	switch blockType {
	case "text":
		text, ok := blockData["text"].(string)
		if !ok {
			return nil, errors.NewMessageParseError("Text block missing 'text' field", blockData, nil)
		}
		return types.NewTextBlock(text), nil

	case "tool_use":
		id, ok := blockData["id"].(string)
		if !ok {
			return nil, errors.NewMessageParseError("Tool use block missing 'id' field", blockData, nil)
		}
		name, ok := blockData["name"].(string)
		if !ok {
			return nil, errors.NewMessageParseError("Tool use block missing 'name' field", blockData, nil)
		}
		input, ok := blockData["input"].(map[string]any)
		if !ok {
			return nil, errors.NewMessageParseError("Tool use block missing 'input' field", blockData, nil)
		}
		return types.NewToolUseBlock(id, name, input), nil

	case "tool_result":
		toolUseID, ok := blockData["tool_use_id"].(string)
		if !ok {
			return nil, errors.NewMessageParseError("Tool result block missing 'tool_use_id' field", blockData, nil)
		}
		content := blockData["content"] // can be nil
		var isError *bool
		if isErrorVal, exists := blockData["is_error"]; exists {
			if isErrorBool, ok := isErrorVal.(bool); ok {
				isError = &isErrorBool
			}
		}
		return types.NewToolResultBlock(toolUseID, content, isError), nil

	default:
		return nil, errors.NewMessageParseError(
			fmt.Sprintf("Unknown content block type: %s", blockType), blockData, nil)
	}
}

func parseSystemMessage(data map[string]any) (*types.SystemMessage, error) {
	subtype, ok := data["subtype"].(string)
	if !ok {
		return nil, errors.NewMessageParseError("Missing required field 'subtype' in system message", data, nil)
	}

	return types.NewSystemMessage(subtype, data), nil
}

func parseResultMessage(data map[string]any) (*types.ResultMessage, error) {
	subtype, ok := data["subtype"].(string)
	if !ok {
		return nil, errors.NewMessageParseError("Missing required field 'subtype' in result message", data, nil)
	}

	durationMs, ok := getIntField(data, "duration_ms")
	if !ok {
		return nil, errors.NewMessageParseError("Missing required field 'duration_ms' in result message", data, nil)
	}

	durationAPIMs, ok := getIntField(data, "duration_api_ms")
	if !ok {
		return nil, errors.NewMessageParseError("Missing required field 'duration_api_ms' in result message", data, nil)
	}

	isError, ok := data["is_error"].(bool)
	if !ok {
		return nil, errors.NewMessageParseError("Missing required field 'is_error' in result message", data, nil)
	}

	numTurns, ok := getIntField(data, "num_turns")
	if !ok {
		return nil, errors.NewMessageParseError("Missing required field 'num_turns' in result message", data, nil)
	}

	sessionID, ok := data["session_id"].(string)
	if !ok {
		return nil, errors.NewMessageParseError("Missing required field 'session_id' in result message", data, nil)
	}

	message := types.NewResultMessage(subtype, durationMs, durationAPIMs, numTurns, isError, sessionID)

	// Optional fields
	if totalCostUSD, exists := data["total_cost_usd"]; exists {
		if cost, ok := totalCostUSD.(float64); ok {
			message.TotalCostUSD = &cost
		}
	}

	if usage, exists := data["usage"]; exists {
		if usageMap, ok := usage.(map[string]any); ok {
			message.Usage = usageMap
		}
	}

	if result, exists := data["result"]; exists {
		if resultStr, ok := result.(string); ok {
			message.Result = &resultStr
		}
	}

	return message, nil
}

func getIntField(data map[string]any, field string) (int, bool) {
	val, exists := data[field]
	if !exists {
		return 0, false
	}

	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

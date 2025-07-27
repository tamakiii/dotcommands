package parser

import (
	"fmt"
	"regexp"
	"strings"

	"llm-mcp-client/mcp"
)

// ToolCallExtractor extracts tool calls from LLM responses
type ToolCallExtractor struct {
	regex *regexp.Regexp
}

// NewToolCallExtractor creates a new tool call extractor
func NewToolCallExtractor() *ToolCallExtractor {
	// Regex to match TOOL_CALL: function_name(param1=value1, param2=value2)
	regex := regexp.MustCompile(`TOOL_CALL:\s*(\w+)\((.*?)\)`)
	return &ToolCallExtractor{regex: regex}
}

// Extract finds and parses tool calls from an LLM response
func (e *ToolCallExtractor) Extract(response string) ([]mcp.ToolCallParams, error) {
	matches := e.regex.FindAllStringSubmatch(response, -1)

	var calls []mcp.ToolCallParams
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		toolName := match[1]
		argsStr := match[2]

		args, err := e.parseArguments(argsStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments for tool %s: %w", toolName, err)
		}

		calls = append(calls, mcp.ToolCallParams{
			Name:      toolName,
			Arguments: args,
		})
	}

	return calls, nil
}

// parseArguments parses function-style arguments into a map
func (e *ToolCallExtractor) parseArguments(argsStr string) (map[string]interface{}, error) {
	args := make(map[string]interface{})
	
	if argsStr == "" {
		return args, nil
	}

	// Split by commas but be careful about arrays
	argPairs, err := e.splitArguments(argsStr)
	if err != nil {
		return nil, err
	}

	for _, pair := range argPairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		parsedValue, err := e.parseValue(value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value for key %s: %w", key, err)
		}

		args[key] = parsedValue
	}

	return args, nil
}

// splitArguments splits arguments by commas, respecting array boundaries
func (e *ToolCallExtractor) splitArguments(argsStr string) ([]string, error) {
	var parts []string
	var current strings.Builder
	bracketDepth := 0
	inQuotes := false
	escapeNext := false

	for _, char := range argsStr {
		if escapeNext {
			current.WriteRune(char)
			escapeNext = false
			continue
		}

		switch char {
		case '\\':
			escapeNext = true
			current.WriteRune(char)
		case '"':
			inQuotes = !inQuotes
			current.WriteRune(char)
		case '[':
			if !inQuotes {
				bracketDepth++
			}
			current.WriteRune(char)
		case ']':
			if !inQuotes {
				bracketDepth--
			}
			current.WriteRune(char)
		case ',':
			if !inQuotes && bracketDepth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	if bracketDepth != 0 {
		return nil, fmt.Errorf("unmatched brackets in arguments")
	}

	if inQuotes {
		return nil, fmt.Errorf("unmatched quotes in arguments")
	}

	return parts, nil
}

// parseValue parses a string value into the appropriate type
func (e *ToolCallExtractor) parseValue(value string) (interface{}, error) {
	// Remove surrounding quotes
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}

	// Handle arrays
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return e.parseArray(value[1 : len(value)-1])
	}

	// Try to parse as number
	if num := e.parseNumber(value); num != nil {
		return num, nil
	}

	// Return as string
	return value, nil
}

// parseArray parses an array string into a slice
func (e *ToolCallExtractor) parseArray(arrayContent string) ([]string, error) {
	if arrayContent == "" {
		return []string{}, nil
	}

	items, err := e.splitArguments(arrayContent)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if strings.HasPrefix(item, "\"") && strings.HasSuffix(item, "\"") {
			item = item[1 : len(item)-1]
		}
		result = append(result, item)
	}

	return result, nil
}

// parseNumber attempts to parse a string as a number
func (e *ToolCallExtractor) parseNumber(s string) interface{} {
	// Try int first
	if i, err := fmt.Sscanf(s, "%d", new(int)); err == nil && i == 1 {
		var num int
		fmt.Sscanf(s, "%d", &num)
		return float64(num) // Convert to float64 for JSON compatibility
	}

	// Try float
	if i, err := fmt.Sscanf(s, "%f", new(float64)); err == nil && i == 1 {
		var num float64
		fmt.Sscanf(s, "%f", &num)
		return num
	}

	return nil
}
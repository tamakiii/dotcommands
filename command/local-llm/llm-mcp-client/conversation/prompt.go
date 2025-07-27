package conversation

import (
	"fmt"
	"strings"

	"llm-mcp-client/mcp"
)

// PromptBuilder creates system prompts for LLM interactions
type PromptBuilder struct{}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// CreateSystemPrompt generates a system prompt that describes available tools
func (pb *PromptBuilder) CreateSystemPrompt(tools []mcp.Tool) string {
	var toolDescriptions []string

	for _, tool := range tools {
		toolDescriptions = append(toolDescriptions, 
			fmt.Sprintf("- %s: %s", tool.Name, tool.Description))
	}

	return fmt.Sprintf(`You are an AI assistant with access to the following tools:

%s

When you want to use a tool, include a tool call in your response using this exact format:
TOOL_CALL: tool_name(parameter1=value1, parameter2=value2)

For example:
- To get a random number: TOOL_CALL: get_random_number(min=1, max=10)
- To get a random string: TOOL_CALL: get_random_string(length=8, charset=alpha)
- To pick from choices: TOOL_CALL: get_random_choice(choices=["apple", "banana", "orange"])

You can use multiple tools in one response if needed. Be helpful and use tools when they would be useful to answer the user's question.`, 
		strings.Join(toolDescriptions, "\n"))
}

// CreateUserPrompt creates a complete prompt including system context and user input
func (pb *PromptBuilder) CreateUserPrompt(systemPrompt, userInput string) string {
	return fmt.Sprintf("%s\n\nUser: %s\n\nAssistant:", systemPrompt, userInput)
}

// CreateToolResultPrompt creates a prompt that includes tool results for final response generation
func (pb *PromptBuilder) CreateToolResultPrompt(originalPrompt string, toolResults []string) string {
	toolResultsText := strings.Join(toolResults, "\n")
	return fmt.Sprintf("%s\n\nTool Results:\n%s\n\nBased on these tool results, provide a final response to the user.\n\nAssistant:", 
		originalPrompt, toolResultsText)
}
package conversation

import (
	"fmt"

	"llm-mcp-client/mcp"
	"llm-mcp-client/ollama"
	"llm-mcp-client/parser"
)

// Handler manages conversation flow between user, LLM, and MCP tools
type Handler struct {
	llmClient       ollama.Client
	mcpClient       mcp.Client
	promptBuilder   *PromptBuilder
	toolExtractor   *parser.ToolCallExtractor
	tools           []mcp.Tool
}

// NewHandler creates a new conversation handler
func NewHandler(llmClient ollama.Client, mcpClient mcp.Client) (*Handler, error) {
	handler := &Handler{
		llmClient:     llmClient,
		mcpClient:     mcpClient,
		promptBuilder: NewPromptBuilder(),
		toolExtractor: parser.NewToolCallExtractor(),
	}

	// Initialize MCP client and get tools
	if err := mcpClient.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize MCP client: %w", err)
	}

	tools, err := mcpClient.ListTools()
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	handler.tools = tools
	return handler, nil
}

// GetToolCount returns the number of available tools
func (h *Handler) GetToolCount() int {
	return len(h.tools)
}

// ProcessUserInput handles a user input and returns the assistant's response
func (h *Handler) ProcessUserInput(userInput string) (string, error) {
	// Create system prompt with available tools
	systemPrompt := h.promptBuilder.CreateSystemPrompt(h.tools)
	fullPrompt := h.promptBuilder.CreateUserPrompt(systemPrompt, userInput)

	// Get LLM response
	response, err := h.llmClient.Generate(fullPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate LLM response: %w", err)
	}

	// Check if LLM wants to use tools
	toolCalls, err := h.toolExtractor.Extract(response)
	if err != nil {
		return "", fmt.Errorf("failed to extract tool calls: %w", err)
	}

	if len(toolCalls) == 0 {
		// No tools requested, return direct response
		return response, nil
	}

	// Execute tool calls
	toolResults, err := h.executeToolCalls(toolCalls)
	if err != nil {
		return "", fmt.Errorf("failed to execute tool calls: %w", err)
	}

	// Get final response with tool results
	finalPrompt := h.promptBuilder.CreateToolResultPrompt(fullPrompt, toolResults)
	finalResponse, err := h.llmClient.Generate(finalPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate final response: %w", err)
	}

	return finalResponse, nil
}

// executeToolCalls executes a list of tool calls and returns formatted results
func (h *Handler) executeToolCalls(toolCalls []mcp.ToolCallParams) ([]string, error) {
	var toolResults []string

	for _, call := range toolCalls {
		result, err := h.mcpClient.CallTool(call.Name, call.Arguments)
		if err != nil {
			toolResults = append(toolResults, 
				fmt.Sprintf("Error calling %s: %v", call.Name, err))
		} else {
			toolResults = append(toolResults, 
				fmt.Sprintf("%s result: %s", call.Name, result))
		}
	}

	return toolResults, nil
}

// HasToolCalls checks if a response contains tool calls without executing them
func (h *Handler) HasToolCalls(response string) (bool, error) {
	toolCalls, err := h.toolExtractor.Extract(response)
	if err != nil {
		return false, err
	}
	return len(toolCalls) > 0, nil
}
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Ollama API structures
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// MCP Protocol structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type MCPClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: llm-mcp-client <mcp-server-path> [model-name]")
		fmt.Println("Example: llm-mcp-client ../bin/mcp-random-server gemma3:12b")
		os.Exit(1)
	}

	mcpServerPath := os.Args[1]
	model := "gemma3:12b"
	if len(os.Args) > 2 {
		model = os.Args[2]
	}

	client := &LLMClient{
		OllamaURL: "http://localhost:11434",
		Model:     model,
	}

	mcpClient, err := NewMCPClient(mcpServerPath)
	if err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
	defer mcpClient.Close()

	if err := mcpClient.Initialize(); err != nil {
		log.Fatalf("Failed to initialize MCP server: %v", err)
	}

	tools, err := mcpClient.ListTools()
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("Connected to MCP server with %d tools available\n", len(tools))
	fmt.Printf("Using model: %s\n\n", model)

	// Interactive loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("User: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}
		if userInput == "quit" || userInput == "exit" {
			break
		}

		// Create context-aware prompt
		systemPrompt := createSystemPrompt(tools)
		fullPrompt := fmt.Sprintf("%s\n\nUser: %s\n\nAssistant:", systemPrompt, userInput)

		// Get LLM response
		response, err := client.Generate(fullPrompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Check if LLM wants to use tools
		toolCalls := extractToolCalls(response)
		if len(toolCalls) > 0 {
			fmt.Printf("Assistant: I'll help you with that. Let me use some tools...\n\n")
			
			// Execute tool calls
			toolResults := make([]string, 0)
			for _, call := range toolCalls {
				result, err := mcpClient.CallTool(call.Name, call.Arguments)
				if err != nil {
					toolResults = append(toolResults, fmt.Sprintf("Error calling %s: %v", call.Name, err))
				} else {
					toolResults = append(toolResults, fmt.Sprintf("%s result: %s", call.Name, result))
				}
			}

			// Get final response with tool results
			toolResultsText := strings.Join(toolResults, "\n")
			finalPrompt := fmt.Sprintf("%s\n\nTool Results:\n%s\n\nBased on these tool results, provide a final response to the user.\n\nAssistant:", fullPrompt, toolResultsText)
			
			finalResponse, err := client.Generate(finalPrompt)
			if err != nil {
				fmt.Printf("Error generating final response: %v\n", err)
				continue
			}
			
			fmt.Printf("Assistant: %s\n\n", finalResponse)
		} else {
			fmt.Printf("Assistant: %s\n\n", response)
		}
	}
}

type LLMClient struct {
	OllamaURL string
	Model     string
}

func (c *LLMClient) Generate(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.OllamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API error: %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", err
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

func NewMCPClient(serverPath string) (*MCPClient, error) {
	cmd := exec.Command(serverPath)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &MCPClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func (c *MCPClient) Close() error {
	c.stdin.Close()
	return c.cmd.Wait()
}

func (c *MCPClient) sendRequest(method string, params interface{}) (*MCPResponse, error) {
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Method:  method,
		Params:  params,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := c.stdin.Write(append(reqBytes, '\n')); err != nil {
		return nil, err
	}

	// Read response
	reader := bufio.NewReader(c.stdout)
	respLine, isPrefix, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, fmt.Errorf("response line too long")
	}

	var resp MCPResponse
	if err := json.Unmarshal(respLine, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}

func (c *MCPClient) Initialize() error {
	_, err := c.sendRequest("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "llm-mcp-client",
			"version": "1.0.0",
		},
	})
	return err
}

func (c *MCPClient) ListTools() ([]map[string]interface{}, error) {
	resp, err := c.sendRequest("tools/list", nil)
	if err != nil {
		return nil, err
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	tools, ok := result["tools"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid tools format")
	}

	var toolList []map[string]interface{}
	for _, tool := range tools {
		if toolMap, ok := tool.(map[string]interface{}); ok {
			toolList = append(toolList, toolMap)
		}
	}

	return toolList, nil
}

func (c *MCPClient) CallTool(name string, arguments map[string]interface{}) (string, error) {
	params := ToolCallParams{
		Name:      name,
		Arguments: arguments,
	}

	resp, err := c.sendRequest("tools/call", params)
	if err != nil {
		return "", err
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("invalid content format")
	}

	if contentItem, ok := content[0].(map[string]interface{}); ok {
		if text, ok := contentItem["text"].(string); ok {
			return text, nil
		}
	}

	return "", fmt.Errorf("could not extract text from response")
}

func createSystemPrompt(tools []map[string]interface{}) string {
	var toolDescriptions []string
	
	for _, tool := range tools {
		name, _ := tool["name"].(string)
		description, _ := tool["description"].(string)
		toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", name, description))
	}

	return fmt.Sprintf(`You are an AI assistant with access to the following tools:

%s

When you want to use a tool, include a tool call in your response using this exact format:
TOOL_CALL: tool_name(parameter1=value1, parameter2=value2)

For example:
- To get a random number: TOOL_CALL: get_random_number(min=1, max=10)
- To get a random string: TOOL_CALL: get_random_string(length=8, charset=alpha)
- To pick from choices: TOOL_CALL: get_random_choice(choices=["apple", "banana", "orange"])

You can use multiple tools in one response if needed. Be helpful and use tools when they would be useful to answer the user's question.`, strings.Join(toolDescriptions, "\n"))
}

func extractToolCalls(response string) []ToolCallParams {
	// Regex to match TOOL_CALL: function_name(param1=value1, param2=value2)
	re := regexp.MustCompile(`TOOL_CALL:\s*(\w+)\((.*?)\)`)
	matches := re.FindAllStringSubmatch(response, -1)

	var calls []ToolCallParams
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		toolName := match[1]
		argsStr := match[2]

		// Parse arguments
		args := make(map[string]interface{})
		if argsStr != "" {
			// Simple argument parsing for key=value pairs
			argPairs := strings.Split(argsStr, ",")
			for _, pair := range argPairs {
				pair = strings.TrimSpace(pair)
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					
					// Remove quotes if present
					if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
						value = value[1 : len(value)-1]
					}
					
					// Handle arrays for choices parameter
					if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
						value = value[1 : len(value)-1]
						choices := strings.Split(value, ",")
						var cleanChoices []string
						for _, choice := range choices {
							choice = strings.TrimSpace(choice)
							if strings.HasPrefix(choice, "\"") && strings.HasSuffix(choice, "\"") {
								choice = choice[1 : len(choice)-1]
							}
							cleanChoices = append(cleanChoices, choice)
						}
						args[key] = cleanChoices
					} else {
						// Try to parse as number
						if num := parseNumber(value); num != nil {
							args[key] = num
						} else {
							args[key] = value
						}
					}
				}
			}
		}

		calls = append(calls, ToolCallParams{
			Name:      toolName,
			Arguments: args,
		})
	}

	return calls
}

func parseNumber(s string) interface{} {
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
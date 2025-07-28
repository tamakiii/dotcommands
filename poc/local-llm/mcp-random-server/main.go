package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

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

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type ToolSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

func main() {
	log.SetOutput(os.Stderr)

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var request MCPRequest
		if err := decoder.Decode(&request); err != nil {
			log.Printf("Error decoding request: %v", err)
			break
		}

		response := handleRequest(request)
		if err := encoder.Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			break
		}
	}
}

func handleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]bool{},
				},
				"serverInfo": ServerInfo{
					Name:    "mcp-random-server",
					Version: "1.0.0",
				},
			},
		}

	case "tools/list":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": []Tool{
					{
						Name:        "get_random_number",
						Description: "Generate a random number within specified range",
						InputSchema: ToolSchema{
							Type: "object",
							Properties: map[string]interface{}{
								"min": map[string]interface{}{
									"type":        "number",
									"description": "Minimum value (inclusive)",
									"default":     0,
								},
								"max": map[string]interface{}{
									"type":        "number",
									"description": "Maximum value (inclusive)",
									"default":     100,
								},
							},
						},
					},
					{
						Name:        "get_random_string",
						Description: "Generate a random string of specified length",
						InputSchema: ToolSchema{
							Type: "object",
							Properties: map[string]interface{}{
								"length": map[string]interface{}{
									"type":        "integer",
									"description": "Length of the random string",
									"default":     10,
									"minimum":     1,
									"maximum":     100,
								},
								"charset": map[string]interface{}{
									"type":        "string",
									"description": "Character set to use (alphanumeric, alpha, numeric)",
									"default":     "alphanumeric",
									"enum":        []string{"alphanumeric", "alpha", "numeric"},
								},
							},
						},
					},
					{
						Name:        "get_random_choice",
						Description: "Pick a random item from a list of choices",
						InputSchema: ToolSchema{
							Type: "object",
							Properties: map[string]interface{}{
								"choices": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{
										"type": "string",
									},
									"description": "List of choices to pick from",
									"minItems":    1,
								},
							},
							Required: []string{"choices"},
						},
					},
				},
			},
		}

	case "tools/call":
		return handleToolCall(req)

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func handleToolCall(req MCPRequest) MCPResponse {
	paramsBytes, _ := json.Marshal(req.Params)
	var params ToolCallParams
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	switch params.Name {
	case "get_random_number":
		return handleRandomNumber(req.ID, params.Arguments)
	case "get_random_string":
		return handleRandomString(req.ID, params.Arguments)
	case "get_random_choice":
		return handleRandomChoice(req.ID, params.Arguments)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Tool not found: %s", params.Name),
			},
		}
	}
}

func handleRandomNumber(id interface{}, args map[string]interface{}) MCPResponse {
	min := 0.0
	max := 100.0

	if v, ok := args["min"]; ok {
		if f, ok := v.(float64); ok {
			min = f
		}
	}
	if v, ok := args["max"]; ok {
		if f, ok := v.(float64); ok {
			max = f
		}
	}

	if min > max {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "min cannot be greater than max",
			},
		}
	}

	result := min + rand.Float64()*(max-min)

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Random number: %.2f", result),
				},
			},
		},
	}
}

func handleRandomString(id interface{}, args map[string]interface{}) MCPResponse {
	length := 10
	charset := "alphanumeric"

	if v, ok := args["length"]; ok {
		if f, ok := v.(float64); ok {
			length = int(f)
		}
	}
	if v, ok := args["charset"]; ok {
		if s, ok := v.(string); ok {
			charset = s
		}
	}

	if length < 1 || length > 100 {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "length must be between 1 and 100",
			},
		}
	}

	var chars string
	switch charset {
	case "alphanumeric":
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	case "alpha":
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "numeric":
		chars = "0123456789"
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "invalid charset, must be one of: alphanumeric, alpha, numeric",
			},
		}
	}

	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Random string: %s", string(result)),
				},
			},
		},
	}
}

func handleRandomChoice(id interface{}, args map[string]interface{}) MCPResponse {
	choices, ok := args["choices"]
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "choices parameter is required",
			},
		}
	}

	choicesArray, ok := choices.([]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "choices must be an array",
			},
		}
	}

	if len(choicesArray) == 0 {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "choices array cannot be empty",
			},
		}
	}

	// Convert to string array
	stringChoices := make([]string, len(choicesArray))
	for i, v := range choicesArray {
		if s, ok := v.(string); ok {
			stringChoices[i] = s
		} else {
			stringChoices[i] = fmt.Sprintf("%v", v)
		}
	}

	choice := stringChoices[rand.Intn(len(stringChoices))]

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Random choice: %s", choice),
				},
			},
		},
	}
}
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/google/uuid"
)

// Client interface for MCP server interactions
type Client interface {
	Initialize() error
	ListTools() ([]Tool, error)
	CallTool(name string, arguments map[string]interface{}) (string, error)
	Close() error
}

// ProcessClient implements the Client interface using subprocess communication
type ProcessClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

// NewProcessClient creates a new MCP client that communicates with a subprocess
func NewProcessClient(serverPath string) (*ProcessClient, error) {
	cmd := exec.Command(serverPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start MCP server: %w", err)
	}

	return &ProcessClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

// Close terminates the MCP server process
func (c *ProcessClient) Close() error {
	c.stdin.Close()
	return c.cmd.Wait()
}

// sendRequest sends a JSON-RPC request to the MCP server
func (c *ProcessClient) sendRequest(method string, params interface{}) (*MCPResponse, error) {
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Method:  method,
		Params:  params,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := c.stdin.Write(append(reqBytes, '\n')); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	// Read response
	reader := bufio.NewReader(c.stdout)
	respLine, isPrefix, err := reader.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if isPrefix {
		return nil, fmt.Errorf("response line too long")
	}

	var resp MCPResponse
	if err := json.Unmarshal(respLine, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}

// Initialize initializes the MCP connection
func (c *ProcessClient) Initialize() error {
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

// ListTools retrieves available tools from the MCP server
func (c *ProcessClient) ListTools() ([]Tool, error) {
	resp, err := c.sendRequest("tools/list", nil)
	if err != nil {
		return nil, err
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	toolsData, ok := result["tools"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid tools format")
	}

	var tools []Tool
	for _, toolData := range toolsData {
		if toolMap, ok := toolData.(map[string]interface{}); ok {
			tool := Tool{}
			if name, ok := toolMap["name"].(string); ok {
				tool.Name = name
			}
			if desc, ok := toolMap["description"].(string); ok {
				tool.Description = desc
			}
			tools = append(tools, tool)
		}
	}

	return tools, nil
}

// CallTool executes a tool on the MCP server
func (c *ProcessClient) CallTool(name string, arguments map[string]interface{}) (string, error) {
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
package mcp

// MCPRequest represents a JSON-RPC request to an MCP server
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC response from an MCP server
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an error in MCP communication
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToolCallParams represents parameters for calling an MCP tool
type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
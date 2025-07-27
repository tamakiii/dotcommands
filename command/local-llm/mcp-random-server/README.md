# MCP Random Server

A simple MCP (Model Context Protocol) server that provides random value generation tools for testing LLM integrations.

## Features

This server provides three tools:

### get_random_number
Generate a random number within a specified range.

**Parameters:**
- `min` (number, optional): Minimum value (inclusive), default: 0
- `max` (number, optional): Maximum value (inclusive), default: 100

### get_random_string  
Generate a random string of specified length and character set.

**Parameters:**
- `length` (integer, optional): Length of the string (1-100), default: 10
- `charset` (string, optional): Character set ("alphanumeric", "alpha", "numeric"), default: "alphanumeric"

### get_random_choice
Pick a random item from a list of choices.

**Parameters:**
- `choices` (array of strings, required): List of choices to pick from

## Usage

### Building
```bash
make build
```

### Running as MCP Server
The server communicates via stdin/stdout using JSON-RPC:

```bash
../bin/mcp-random-server
```

### Testing with curl
You can test the server by sending JSON-RPC requests:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ../bin/mcp-random-server
```

## MCP Protocol Support

This server implements the MCP protocol specification and supports:
- Server initialization
- Tool listing  
- Tool execution
- Proper error handling

## Integration

This server is designed to work with the `llm-mcp-client` to provide random value generation capabilities to local LLMs like gemma3:12b.
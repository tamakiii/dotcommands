# LLM MCP Client

A client that enables local LLMs (like gemma3:12b via Ollama) to interact with MCP (Model Context Protocol) servers.

## Features

- Connects to Ollama for LLM inference
- Communicates with MCP servers via JSON-RPC
- Parses LLM responses for tool calls
- Executes MCP server tools and provides results back to the LLM
- Interactive chat interface

## Prerequisites

1. **Ollama** must be installed and running:
   ```bash
   # Install Ollama (if not already installed)
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Start Ollama service
   ollama serve
   
   # Pull the model (in another terminal)
   ollama pull gemma3:12b
   ```

2. **MCP Server** (e.g., mcp-random-server) must be built:
   ```bash
   make -C ../mcp-random-server build
   ```

## Usage

### Building
```bash
make build
```

### Running
```bash
../bin/llm-mcp-client <mcp-server-path> [model-name]
```

### Examples

Using with the random server:
```bash
../bin/llm-mcp-client ../bin/mcp-random-server gemma3:12b
```

Using with a different model:
```bash
../bin/llm-mcp-client ../bin/mcp-random-server llama3.2:3b
```

## How it Works

1. **Startup**: Connects to specified MCP server and Ollama
2. **Tool Discovery**: Lists available tools from MCP server
3. **Chat Loop**: 
   - Takes user input
   - Sends context-aware prompt to LLM
   - Parses LLM response for tool calls (format: `TOOL_CALL: tool_name(param=value)`)
   - Executes requested tools via MCP server
   - Provides tool results back to LLM for final response

## Tool Call Format

The LLM is instructed to use this format for tool calls:
```
TOOL_CALL: get_random_number(min=1, max=100)
TOOL_CALL: get_random_string(length=10, charset=alpha)
TOOL_CALL: get_random_choice(choices=["option1", "option2", "option3"])
```

## Example Conversation

```
User: Give me a random number between 1 and 10
Assistant: I'll help you with that. Let me use some tools...

get_random_number result: Random number: 7.23
# LLM MCP Client (Python + LangChain)

A modern Python implementation that enables local LLMs (via Ollama) to interact with MCP (Model Context Protocol) servers using LangChain's robust framework for better maintainability and extensibility.

## Features

- **LangChain Integration**: Built on LangChain for robust LLM orchestration
- **Native MCP Support**: Clean integration with MCP servers via standardized protocols
- **Ollama Integration**: Seamless connection to local Ollama models
- **Interactive Chat**: User-friendly command-line chat interface
- **Tool Discovery**: Automatic discovery and integration of MCP server tools
- **Error Handling**: Robust error handling and logging
- **Async Support**: Asynchronous operations for better performance

## Prerequisites

### 1. Python Environment
```bash
# Python 3.8+ required
python3 --version
```

### 2. Ollama Setup
```bash
# Install Ollama (if not already installed)
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama service
ollama serve

# Pull your preferred model (in another terminal)
ollama pull gemma3:12b
# or
ollama pull llama3.2:3b
```

### 3. MCP Server
Build the MCP server you want to use:
```bash
# Example: build the random server
make -C ../mcp-random-server build
```

## Installation & Setup

### Quick Start
```bash
# Setup and build in one command
make setup && make build

# The executable will be created at ../bin/llm-mcp-client
```

### Manual Setup
```bash
# Create virtual environment and install dependencies
make setup

# Or manually:
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Create executable
make build
```

## Usage

### Basic Usage
```bash
../bin/llm-mcp-client <mcp-server-path> [--model <model-name>]
```

### Examples

**Using with the random server:**
```bash
../bin/llm-mcp-client ../bin/mcp-random-server
```

**With a specific model:**
```bash
../bin/llm-mcp-client ../bin/mcp-random-server --model llama3.2:3b
```

**Direct Python execution:**
```bash
# From the llm-mcp-client directory
source venv/bin/activate
python main.py ../bin/mcp-random-server --model gemma3:12b
```

## How It Works

### Architecture
1. **MCP Client**: Manages communication with MCP servers via JSON-RPC
2. **LangChain Integration**: Uses `ChatOllama` for LLM interactions
3. **Tool Wrapper**: Wraps MCP tools as LangChain-compatible tools
4. **Chat Manager**: Handles conversation flow and tool orchestration

### Workflow
1. **Startup**: Connects to MCP server and discovers available tools
2. **Tool Discovery**: Automatically registers MCP tools with LangChain
3. **Chat Loop**: 
   - Takes user input
   - Sends to LLM with tool descriptions
   - Parses LLM response for tool calls (format: `USE_TOOL: tool_name(param=value)`)
   - Executes requested tools via MCP server
   - Returns final response with tool results

### Tool Call Format
The LLM is instructed to use this format for tool calls:
```
USE_TOOL: get_random_number(min=1, max=100)
USE_TOOL: get_random_string(length=10, charset=alpha)
USE_TOOL: get_random_choice(choices=["option1", "option2", "option3"])
```

## Example Conversation

```
🚀 Starting LLM MCP Client (Python + LangChain)
📡 MCP Server: ../bin/mcp-random-server
🤖 Model: gemma3:12b

✓ Connected to MCP server: ../bin/mcp-random-server
✓ Discovered 3 tools: ['get_random_number', 'get_random_string', 'get_random_choice']

💬 Chat started! Type 'quit' or 'exit' to stop.

User: Give me a random number between 1 and 10

🔧 Using tools...
  ✓ get_random_number → Random number: 7
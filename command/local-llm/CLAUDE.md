# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

This is a research and development workspace for implementing local LLM (Large Language Model) solutions on MacBook Pro M3 Max. The primary goal is building reliable local systems for:

- Periodic background tasks (cron jobs) that check files and act on natural language prompts
- Light reasoning tasks like scoring coding agent outputs with numeric evaluations  
- Querying MCP servers to retrieve and analyze conversation histories from coding agents

## Project Structure

```
local-llm/
├── preferred-models.json    # Curated list of local LLM models with ollama.com metadata
├── mcp-random-server/       # MCP server providing random value generation tools
│   ├── main.go             # MCP server implementation
│   ├── go.mod              # Go module definition
│   ├── Makefile            # Build configuration
│   └── README.md           # Server documentation
├── llm-mcp-client/          # LLM client with MCP server integration
│   ├── main.go             # Client implementation with Ollama integration
│   ├── go.mod              # Go module definition
│   ├── Makefile            # Build configuration
│   └── README.md           # Client documentation
├── examples/                # Usage examples and test scripts
│   ├── usage-demo.sh       # Interactive demo script
│   └── test-mcp-server.sh  # MCP server testing script
├── bin/                     # Compiled executables (generated)
│   ├── mcp-random-server   # MCP server binary
│   └── llm-mcp-client      # LLM client binary
├── research/
│   ├── 20250727a/          # Research session exploring model selection and setup
│   │   ├── PROMPT.md       # Original research requirements and use cases
│   │   └── *.md           # Research findings from various LLM assistants
│   └── 20250727b/          # Model performance testing results
│       ├── RESULT.md       # Performance analysis and recommendations
│       └── test_models.py  # Model testing script
└── Makefile                # Main build orchestration
```

## Architecture Overview

This project has evolved from research phase to implementation, featuring a complete MCP integration system:

1. **Model Catalog**: `preferred-models.json` maintains model specifications including size, context length, and input types from ollama.com
2. **MCP Integration**: Full Model Context Protocol implementation enabling local LLMs to use external tools
3. **Research Documentation**: Structured exploration of local LLM options, focusing on model selection, management tools (Ollama, LM Studio), and inference optimization
4. **Target Integration**: Designed to integrate with existing dotcommands workbench for shell workflow automation

### MCP System Components

- **MCP Random Server**: Provides three tools for testing and demonstration:
  - `get_random_number`: Generate random numbers within specified ranges
  - `get_random_string`: Create random strings with configurable length and character sets
  - `get_random_choice`: Pick random items from provided choice arrays

- **LLM MCP Client**: Bridges local LLMs (via Ollama) with MCP servers:
  - Connects to Ollama API for model inference
  - Parses LLM responses for tool call requests
  - Executes MCP server tools and provides results back to the LLM
  - Interactive chat interface with tool integration

## Key Models Under Consideration

The preferred models are optimized for MacBook Pro M3 Max hardware:

- **Multi-modal**: qwen2.5vl:7b, gemma3:12b (support text + images)
- **High Context**: phi3:3.8b, llama3.2:3b, gemma3:12b (128K+ tokens)
- **Lightweight**: llama3.2:3b (2.0GB), phi3.5:3.8b (2.2GB)
- **Capability**: phi4:latest (9.1GB, 14B parameters)

## Development Context

This project operates within the broader dotcommands ecosystem, following Go-based modular architecture principles. The current implementation:

- Integrates with existing shell workflows and JSON-based tool communication
- Follows the established pattern of independent utilities with shared build processes
- Maintains the dotcommands philosophy of lightweight, focused command-line tools
- Provides a foundation for building more sophisticated MCP servers and LLM integrations

### Build and Usage

```bash
# Build all components
make build

# Install dependencies
make install-deps

# Test MCP server functionality
./examples/test-mcp-server.sh

# Run interactive demo (requires Ollama + gemma3:12b)
./examples/usage-demo.sh

# Manual usage
./bin/llm-mcp-client ./bin/mcp-random-server gemma3:12b
```

### Tool Call Format

The LLM client expects tool calls in this format:
```
TOOL_CALL: get_random_number(min=1, max=100)
TOOL_CALL: get_random_string(length=8, charset=alpha)
TOOL_CALL: get_random_choice(choices=["red", "blue", "green"])
```

## Research Methodology

The project uses structured research sessions (e.g., `research/20250727a/`) to explore options systematically. Each session documents findings from different AI assistants to capture diverse perspectives on model selection, tooling, and implementation approaches.
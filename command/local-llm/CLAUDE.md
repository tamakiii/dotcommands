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
├── llm-mcp-client/          # LLM client with MCP server integration (Python + LangChain)
│   ├── main.py             # Client implementation with LangChain and Ollama integration
│   ├── requirements.txt    # Python dependencies
│   ├── venv/               # Python virtual environment
│   ├── Makefile            # Build configuration
│   └── README.md           # Client documentation
├── llm-mcp-go-client/       # Legacy Go implementation (preserved for reference)
│   ├── main.go             # Original Go implementation
│   ├── go.mod              # Go module definition
│   ├── Makefile            # Build configuration
│   └── README.md           # Go client documentation
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

This project has evolved from research phase to production implementation, featuring a complete MCP integration system built on modern Python + LangChain architecture:

1. **Model Catalog**: `preferred-models.json` maintains model specifications including size, context length, and input types from ollama.com
2. **MCP Integration**: Full Model Context Protocol implementation enabling local LLMs to use external tools
3. **LangChain Foundation**: Modern Python implementation using LangChain for robust LLM orchestration and tool integration
4. **Research Documentation**: Structured exploration of local LLM options, focusing on model selection, management tools (Ollama, LM Studio), and inference optimization
5. **Target Integration**: Designed to integrate with existing dotcommands workbench for shell workflow automation

### MCP System Components

- **MCP Random Server**: Provides three tools for testing and demonstration:
  - `get_random_number`: Generate random numbers within specified ranges
  - `get_random_string`: Create random strings with configurable length and character sets
  - `get_random_choice`: Pick random items from provided choice arrays

- **LLM MCP Client** (Python + LangChain): Modern implementation bridging local LLMs with MCP servers:
  - Built on LangChain's ChatOllama for robust LLM integration
  - Native MCP protocol support with clean JSON-RPC communication
  - Async architecture for better performance and reliability
  - Automatic tool discovery and registration
  - Interactive chat interface with visual feedback and error handling
  - Easy extensibility for new model providers and MCP servers

## Key Models Under Consideration

The preferred models are optimized for MacBook Pro M3 Max hardware:

- **Multi-modal**: qwen2.5vl:7b, gemma3:12b (support text + images)
- **High Context**: phi3:3.8b, llama3.2:3b, gemma3:12b (128K+ tokens)
- **Lightweight**: llama3.2:3b (2.0GB), phi3.5:3.8b (2.2GB)
- **Capability**: phi4:latest (9.1GB, 14B parameters)

## Development Context

This project operates within the broader dotcommands ecosystem, following modular architecture principles with a modern Python + LangChain foundation. The current implementation:

- **Hybrid Architecture**: Go MCP servers + Python LangChain client for optimal performance and maintainability
- **Shell Integration**: Seamless integration with existing shell workflows and JSON-based tool communication
- **Modular Design**: Follows the established pattern of independent utilities with shared build processes
- **CLI Philosophy**: Maintains the dotcommands philosophy of lightweight, focused command-line tools
- **Extensible Foundation**: LangChain-based architecture enables easy addition of new model providers, tools, and integrations
- **Migration Strategy**: Preserves original Go implementation as reference while adopting industry-standard Python tooling

### Build and Usage

```bash
# Quick setup and build (recommended)
make install

# Build all components individually
make build

# Test MCP server functionality
./examples/test-mcp-server.sh

# Run interactive demo (requires Ollama + gemma3:12b)
./examples/usage-demo.sh

# Manual usage (absolute path required for MCP server)
./bin/llm-mcp-client $PWD/bin/mcp-random-server --model gemma3:12b

# Alternative models
./bin/llm-mcp-client $PWD/bin/mcp-random-server --model llama3.2:3b
./bin/llm-mcp-client $PWD/bin/mcp-random-server --model phi4:latest
```

### Tool Call Format

The Python LangChain client uses an improved tool call format:
```
USE_TOOL: get_random_number(min=1, max=100)
USE_TOOL: get_random_string(length=8, charset=alpha)  
USE_TOOL: get_random_choice(choices=["red", "blue", "green"])
```

**Legacy Go Client Format** (preserved for reference):
```
TOOL_CALL: get_random_number(min=1, max=100)
TOOL_CALL: get_random_string(length=8, charset=alpha)
TOOL_CALL: get_random_choice(choices=["red", "blue", "green"])
```

## Research Methodology

The project uses structured research sessions (e.g., `research/20250727a/`) to explore options systematically. Each session documents findings from different AI assistants to capture diverse perspectives on model selection, tooling, and implementation approaches.
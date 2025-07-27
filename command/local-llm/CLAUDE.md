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
├── research/
│   └── 20250727a/          # Research session exploring model selection and setup
│       ├── PROMPT.md       # Original research requirements and use cases
│       └── *.md           # Research findings from various LLM assistants
```

## Architecture Overview

This is a research-focused project currently in the exploration phase. The architecture centers around:

1. **Model Catalog**: `preferred-models.json` maintains model specifications including size, context length, and input types from ollama.com
2. **Research Documentation**: Structured exploration of local LLM options, focusing on model selection, management tools (Ollama, LM Studio), and inference optimization
3. **Target Integration**: Designed to integrate with existing dotcommands workbench for shell workflow automation

## Key Models Under Consideration

The preferred models are optimized for MacBook Pro M3 Max hardware:

- **Multi-modal**: qwen2.5vl:7b, gemma3:12b (support text + images)
- **High Context**: phi3:3.8b, llama3.2:3b, gemma3:12b (128K+ tokens)
- **Lightweight**: llama3.2:3b (2.0GB), phi3.5:3.8b (2.2GB)
- **Capability**: phi4:latest (9.1GB, 14B parameters)

## Development Context

This project operates within the broader dotcommands ecosystem, which follows Go-based modular architecture principles. Future implementations will likely:

- Integrate with existing shell workflows and JSON-based tool communication
- Follow the established pattern of independent utilities with shared build processes
- Maintain the dotcommands philosophy of lightweight, focused command-line tools

## Research Methodology

The project uses structured research sessions (e.g., `research/20250727a/`) to explore options systematically. Each session documents findings from different AI assistants to capture diverse perspectives on model selection, tooling, and implementation approaches.
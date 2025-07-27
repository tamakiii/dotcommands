#!/bin/bash

# LLM MCP Integration Demo Script
# This script demonstrates the integration between gemma3:12b and MCP servers

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== LLM MCP Integration Demo ===${NC}"
echo

# Check if Ollama is running
check_ollama() {
    echo -e "${YELLOW}Checking Ollama service...${NC}"
    if ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo -e "${RED}Error: Ollama is not running${NC}"
        echo "Please start Ollama with: ollama serve"
        exit 1
    fi
    echo -e "${GREEN}✓ Ollama is running${NC}"
}

# Check if model is available
check_model() {
    local model="$1"
    echo -e "${YELLOW}Checking if model '$model' is available...${NC}"
    if ! ollama list | grep -q "$model"; then
        echo -e "${RED}Error: Model '$model' not found${NC}"
        echo "Please pull the model with: ollama pull $model"
        exit 1
    fi
    echo -e "${GREEN}✓ Model '$model' is available${NC}"
}

# Build the components
build_components() {
    echo -e "${YELLOW}Building MCP components...${NC}"
    cd "$(dirname "$0")/.."
    make build
    echo -e "${GREEN}✓ Components built successfully${NC}"
}

# Test MCP server
test_mcp_server() {
    echo -e "${YELLOW}Testing MCP server...${NC}"
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./bin/mcp-random-server > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ MCP server is working${NC}"
    else
        echo -e "${RED}Error: MCP server test failed${NC}"
        exit 1
    fi
}

# Interactive demo
run_demo() {
    local model="$1"
    echo -e "${BLUE}Starting interactive demo with $model...${NC}"
    echo -e "${YELLOW}Try these example queries:${NC}"
    echo "- 'Give me a random number between 1 and 100'"
    echo "- 'Generate a random 8-character string'"
    echo "- 'Pick a random color from red, blue, green, yellow'"
    echo "- Type 'quit' to exit"
    echo
    
    ./bin/llm-mcp-client ./bin/mcp-random-server "$model"
}

# Main execution
main() {
    local model="${1:-gemma3:12b}"
    
    echo -e "${BLUE}Using model: $model${NC}"
    echo
    
    check_ollama
    check_model "$model"
    build_components
    test_mcp_server
    
    echo
    echo -e "${GREEN}All checks passed! Starting demo...${NC}"
    echo
    
    run_demo "$model"
}

# Handle script arguments
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Usage: $0 [model-name]"
    echo
    echo "Examples:"
    echo "  $0                    # Use default model (gemma3:12b)"
    echo "  $0 llama3.2:3b       # Use llama3.2:3b model"
    echo "  $0 qwen2.5vl:7b      # Use qwen2.5vl:7b model"
    echo
    exit 0
fi

main "$@"
#!/bin/bash

# Test script for MCP random server
# This script tests all available tools in the MCP server

set -e

echo "=== Testing MCP Random Server ==="
echo

# Build the server first
echo "Building MCP server..."
cd "$(dirname "$0")/.."
make mcp-server

echo "Starting tests..."
echo

# Test 1: Initialize
echo "Test 1: Server initialization"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./bin/mcp-random-server | jq .
echo

# Test 2: List tools  
echo "Test 2: List available tools"
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./bin/mcp-random-server | jq .
echo

# Test 3: Get random number
echo "Test 3: Get random number (1-10)"
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_random_number","arguments":{"min":1,"max":10}}}' | ./bin/mcp-random-server | jq .
echo

# Test 4: Get random string
echo "Test 4: Get random string (length 8, alpha)"
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_random_string","arguments":{"length":8,"charset":"alpha"}}}' | ./bin/mcp-random-server | jq .
echo

# Test 5: Get random choice
echo "Test 5: Get random choice from fruits"
echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_random_choice","arguments":{"choices":["apple","banana","orange","grape"]}}}' | ./bin/mcp-random-server | jq .
echo

echo "=== All tests completed ==="
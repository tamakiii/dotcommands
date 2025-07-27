.PHONY: help setup teardown build clean test mcp-server client demo

SUBDIRS = mcp-random-server llm-mcp-client

help:
	@cat $(firstword $(MAKEFILE_LIST))

setup:
	$(MAKE) -C mcp-random-server setup
	$(MAKE) -C llm-mcp-client setup

build:
	@mkdir -p bin
	$(MAKE) -C mcp-random-server build
	$(MAKE) -C llm-mcp-client build

clean:
	$(MAKE) -C mcp-random-server clean
	$(MAKE) -C llm-mcp-client clean
	@rm -rf bin

test:
	$(MAKE) -C mcp-random-server test
	$(MAKE) -C llm-mcp-client test

# Individual component targets
mcp-server:
	$(MAKE) -C mcp-random-server build

client:
	$(MAKE) -C llm-mcp-client build

# Quick test of the integration
demo: build
	@echo "Starting demo..."
	@echo "Make sure Ollama is running and gemma3:12b is available"
	@echo "Run: ./bin/llm-mcp-client ./bin/mcp-random-server gemma3:12b"

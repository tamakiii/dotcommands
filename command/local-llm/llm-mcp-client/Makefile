#!/usr/bin/make -f

SHELL := /bin/bash
VENV_DIR := venv
PYTHON := $(VENV_DIR)/bin/python
PIP := $(VENV_DIR)/bin/pip
SCRIPT_NAME := llm-mcp-client

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  setup     - Create virtual environment and install dependencies"
	@echo "  build     - Create executable script in ../bin/"
	@echo "  install   - Install dependencies in existing virtual environment"
	@echo "  test      - Run basic tests"
	@echo "  clean     - Remove virtual environment and build artifacts"
	@echo "  run       - Run the client (requires MCP_SERVER_PATH and optional MODEL)"
	@echo ""
	@echo "Usage examples:"
	@echo "  make setup && make build"
	@echo "  make run MCP_SERVER_PATH=../bin/mcp-random-server MODEL=gemma3:12b"

# Create virtual environment and install dependencies
.PHONY: setup
setup: $(VENV_DIR)
	$(PIP) install -r requirements.txt
	@echo "✓ Setup complete. Virtual environment created with dependencies installed."

# Create virtual environment
$(VENV_DIR):
	python3 -m venv $(VENV_DIR)
	$(PIP) install --upgrade pip

# Install dependencies in existing virtual environment
.PHONY: install
install:
	@if [ ! -d "$(VENV_DIR)" ]; then \
		echo "❌ Virtual environment not found. Run 'make setup' first."; \
		exit 1; \
	fi
	$(PIP) install -r requirements.txt

# Create executable script
.PHONY: build
build: setup
	@mkdir -p ../bin
	@echo "#!/bin/bash" > ../bin/$(SCRIPT_NAME)
	@echo "cd \$$(dirname \$$0)/../llm-mcp-client" >> ../bin/$(SCRIPT_NAME)
	@echo "$(PWD)/$(PYTHON) main.py \"\$$@\"" >> ../bin/$(SCRIPT_NAME)
	@chmod +x ../bin/$(SCRIPT_NAME)
	@echo "✓ Executable created: ../bin/$(SCRIPT_NAME)"

# Run the client (for testing)
.PHONY: run
run: setup
	@if [ -z "$(MCP_SERVER_PATH)" ]; then \
		echo "❌ MCP_SERVER_PATH is required. Usage: make run MCP_SERVER_PATH=../bin/mcp-random-server"; \
		exit 1; \
	fi
	$(PYTHON) main.py $(MCP_SERVER_PATH) $(if $(MODEL),--model $(MODEL),)

# Basic tests
.PHONY: test
test: setup
	$(PYTHON) -c "import sys; print('✓ Python executable works')"
	$(PYTHON) -c "import langchain_ollama; print('✓ langchain-ollama imported successfully')"
	$(PYTHON) -c "import langchain_core; print('✓ langchain-core imported successfully')"
	$(PYTHON) main.py --help > /dev/null && echo "✓ main.py help works"
	@echo "✓ All basic tests passed"

# Clean up
.PHONY: clean
clean:
	rm -rf $(VENV_DIR)
	rm -f ../bin/$(SCRIPT_NAME)
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	find . -name "*.pyc" -delete 2>/dev/null || true
	@echo "✓ Cleanup complete"

# Development targets
.PHONY: dev-install
dev-install: setup
	$(PIP) install black isort flake8
	@echo "✓ Development dependencies installed"

.PHONY: format
format: dev-install
	$(VENV_DIR)/bin/black main.py
	$(VENV_DIR)/bin/isort main.py
	@echo "✓ Code formatted"

.PHONY: lint
lint: dev-install
	$(VENV_DIR)/bin/flake8 main.py --max-line-length=100
	@echo "✓ Linting complete"
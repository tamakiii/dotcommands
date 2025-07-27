SUBDIRS = mcp-random-server llm-mcp-client

build:
	@mkdir -p bin
	@for dir in $(SUBDIRS); do \
		echo "Building $$dir..."; \
		$(MAKE) -C $$dir build || exit 1; \
	done

test:
	@for dir in $(SUBDIRS); do \
		echo "Testing $$dir..."; \
		$(MAKE) -C $$dir test || exit 1; \
	done

clean:
	@for dir in $(SUBDIRS); do \
		echo "Cleaning $$dir..."; \
		$(MAKE) -C $$dir clean; \
	done
	@rm -rf bin

install-deps:
	@for dir in $(SUBDIRS); do \
		echo "Installing dependencies for $$dir..."; \
		$(MAKE) -C $$dir install-deps || exit 1; \
	done

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

.PHONY: build test clean install-deps mcp-server client demo
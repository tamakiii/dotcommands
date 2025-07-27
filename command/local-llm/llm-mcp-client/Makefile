build:
	go build -o ../bin/llm-mcp-client .

test:
	go test -v ./...

clean:
	rm -f ../bin/llm-mcp-client

install-deps:
	go mod tidy

.PHONY: build test clean install-deps
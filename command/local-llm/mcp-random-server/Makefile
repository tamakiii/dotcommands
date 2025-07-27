build:
	go build -o ../bin/mcp-random-server .

test:
	go test -v ./...

clean:
	rm -f ../bin/mcp-random-server

install-deps:
	go mod tidy

.PHONY: build test clean install-deps
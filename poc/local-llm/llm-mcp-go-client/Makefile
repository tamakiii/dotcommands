.PHONY: help setup teardown build test clean

help:
	@cat $(firstword $(MAKEFILE_LIST))

setup:
	go mod tidy

teardown:
	# TODO

build:
	go build -o ../bin/llm-mcp-client .

test:
	go test -v ./...

clean:
	rm -f ../bin/llm-mcp-client

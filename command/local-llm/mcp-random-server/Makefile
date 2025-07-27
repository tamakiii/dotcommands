.PHONY: help setup teardown build test clean install-deps

help:
	@cat $(firstword $(MAKEFILE_LIST))

setup:
	go mod tidy

teardown:
	# TODO:

build:
	go build -o ../bin/mcp-random-server .

test:
	go test -v ./...

clean:
	rm -f ../bin/mcp-random-server

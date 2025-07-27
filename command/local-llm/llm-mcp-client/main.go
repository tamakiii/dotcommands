package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"llm-mcp-client/config"
	"llm-mcp-client/conversation"
	"llm-mcp-client/mcp"
	"llm-mcp-client/ollama"
)

func main() {
	// Parse configuration from command line arguments
	cfg, err := config.LoadFromArgs(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Example: llm-mcp-client ../bin/mcp-random-server gemma3:12b")
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize clients
	llmClient := ollama.NewHTTPClient(cfg.OllamaURL, cfg.Model)
	mcpClient, err := mcp.NewProcessClient(cfg.MCPServerPath)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	// Initialize conversation handler
	conversationHandler, err := conversation.NewHandler(llmClient, mcpClient)
	if err != nil {
		log.Fatalf("Failed to initialize conversation handler: %v", err)
	}

	// Display startup information
	fmt.Printf("Connected to MCP server with %d tools available\n", conversationHandler.GetToolCount())
	fmt.Printf("Using model: %s\n\n", cfg.Model)

	// Start interactive loop
	if err := runInteractiveLoop(conversationHandler); err != nil {
		log.Fatalf("Interactive loop error: %v", err)
	}
}

// runInteractiveLoop handles the main conversation loop
func runInteractiveLoop(handler *conversation.Handler) error {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("User: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}
		if userInput == "quit" || userInput == "exit" {
			break
		}

		// Process user input through conversation handler
		response, err := handler.ProcessUserInput(userInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Check if tools were used and provide appropriate feedback
		hasTools, err := handler.HasToolCalls(response)
		if err != nil {
			fmt.Printf("Error checking tool calls: %v\n", err)
		}

		if hasTools {
			fmt.Printf("Assistant: I'll help you with that. Let me use some tools...\n\n")
		}

		fmt.Printf("Assistant: %s\n\n", response)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}
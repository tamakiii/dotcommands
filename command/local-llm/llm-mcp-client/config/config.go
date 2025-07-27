package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	MCPServerPath string
	Model         string
	OllamaURL     string
}

// DefaultConfig returns configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Model:     "gemma3:12b",
		OllamaURL: "http://localhost:11434",
	}
}

// LoadFromArgs parses command line arguments and returns configuration
func LoadFromArgs(args []string) (*Config, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("usage: llm-mcp-client <mcp-server-path> [model-name]")
	}

	config := DefaultConfig()
	config.MCPServerPath = args[1]

	if len(args) > 2 {
		config.Model = args[2]
	}

	// Allow environment variable override
	if url := os.Getenv("OLLAMA_URL"); url != "" {
		config.OllamaURL = url
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.MCPServerPath == "" {
		return fmt.Errorf("MCP server path is required")
	}
	
	if c.Model == "" {
		return fmt.Errorf("model name is required")
	}
	
	if c.OllamaURL == "" {
		return fmt.Errorf("Ollama URL is required")
	}

	return nil
}
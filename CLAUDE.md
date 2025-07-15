# CLAUDE.md

This file provides guidance to Claude Code when working with the dotcommands repository.

## Repository Overview

This is a **Go-based development workbench** for building custom command-line utilities that enhance shell workflows and development productivity. The repository follows a modular architecture where each tool is an independent Go module within a shared workspace.

## Project Architecture

### Structure
```
.dotcommands/
├── README.md              # Project overview
├── command/                # Main workbench directory  
│   ├── Makefile           # Build orchestration
│   ├── README.md          # Comprehensive documentation
│   ├── bin/               # Compiled executables (generated)
│   └── {tool-name}/       # Individual tool modules
│       ├── main.go        # Tool implementation
│       ├── main_test.go   # Test suite
│       ├── Makefile       # Tool-specific build config
│       ├── README.md      # Tool documentation
│       └── go.mod         # Independent Go module
```

### Design Principles
- **Modular Tools**: Each utility is an independent Go module
- **Shared Workspace**: Common `bin/` directory with symlinks to executables
- **Shell Integration**: Tools designed for pipes, JSON output, and script automation
- **Test-Driven**: Comprehensive test coverage for all utilities
- **Documentation-First**: Each tool includes usage examples and integration patterns

## Development Workflow

### Building Tools
```bash
# Build all tools
make -C command build

# Build specific tool
make -C command/git-url build

# Clean artifacts
make -C command clean
```

### Adding New Tools
1. Create new directory under `command/`
2. Initialize Go module: `go mod init`
3. Implement tool following existing patterns
4. Add Makefile with build/test/clean targets
5. Create comprehensive README with usage examples
6. Add test suite with good coverage

### Testing
```bash
# Test all tools
make -C command test

# Test specific tool
make -C command/git-url test
```

## Current Tools

### git-url
**Purpose**: Parse Git URLs and output structured JSON
**Location**: `command/git-url/`
**Usage**: `git remote get-url origin | git-url`
**Integration**: Pipes well with `jq` for JSON processing

## Code Standards

### Go Conventions
- Use standard Go project layout
- Follow effective Go patterns
- Include comprehensive error handling
- Use the `go-git` library for Git operations when needed

### Testing Standards
- Test coverage should exceed implementation code lines
- Include edge cases and error conditions
- Use table-driven tests for multiple scenarios
- Test both positive and negative cases

### Documentation Requirements
- Each tool must have a README with:
  - Purpose and use cases
  - Installation instructions
  - Usage examples
  - Integration patterns
  - Development workflow

### Build Integration
- Each tool has independent `go.mod`
- Shared Makefile patterns for consistency
- Executables placed in shared `bin/` directory
- Use symlinks for easy PATH integration

## Integration Patterns

### Shell Workflow
Tools are designed to integrate naturally with shell workflows:
- Read from stdin when appropriate
- Output structured data (JSON preferred)
- Use exit codes appropriately
- Support piping with other tools

### JSON Output
When outputting structured data:
- Use consistent field naming (snake_case)
- Include all relevant information
- Make output parseable with `jq`
- Consider both human and machine readability

## Development Best Practices

### Security
- Never expose credentials or sensitive data
- Validate all inputs thoroughly
- Use secure defaults
- Follow Go security best practices

### Performance
- Keep tools lightweight and fast
- Avoid unnecessary dependencies
- Use efficient algorithms
- Consider memory usage for large inputs

### Error Handling
- Provide clear error messages
- Use appropriate exit codes
- Log errors to stderr
- Handle edge cases gracefully

## Makefile Standards

### Tool-Level Makefiles
Each tool should include:
```makefile
build:
	go build -o ../bin/tool-name .

test:
	go test -v ./...

clean:
	rm -f ../bin/tool-name

.PHONY: build test clean
```

### Workspace-Level Makefile
The main Makefile should orchestrate all tools:
- Build all tools with single command
- Run all tests
- Clean all artifacts
- Manage shared resources

## Future Direction

This workbench is designed to grow with additional utilities that enhance development workflows. Each new tool should:
- Solve a specific development pain point
- Integrate well with existing shell workflows  
- Follow the established patterns for modularity
- Include comprehensive documentation and tests
- Be lightweight and focused on a single responsibility

## Contributing Guidelines

When adding new tools or modifying existing ones:
1. Follow the established project structure
2. Maintain test coverage standards
3. Update documentation thoroughly
4. Ensure tools work well in shell pipelines
5. Consider cross-platform compatibility
6. Use semantic versioning for releases

This repository serves as a foundation for building a personal collection of development utilities, with emphasis on modularity, testing, and shell integration.
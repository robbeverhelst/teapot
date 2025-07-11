# Claude Assistant Context

## Project Overview

Teapot is a CLI tool for bootstrapping modern full-stack monorepos. It's built with Go and uses the Bubble Tea framework for the terminal UI. The tool helps developers quickly scaffold production-ready projects with pre-configured tooling, infrastructure, and CI/CD pipelines.

### Key Components
- **CLI Interface**: Built with github.com/charmbracelet/bubbletea for interactive prompts
- **Template Engine**: Generates project structures based on user selections
- **Infrastructure Generator**: Creates Docker Compose and Kubernetes configurations
- **CI/CD Setup**: Configures GitHub Actions workflows

## Go Development Rules

### Code Style
- Use `gofmt` and `goimports` for consistent formatting
- Follow standard Go naming conventions (camelCase for private, PascalCase for public)
- Keep functions small and focused on a single responsibility
- Prefer composition over inheritance
- Use meaningful variable and function names

### Error Handling
```go
// Always check and handle errors explicitly
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### Project Structure
```
teapot/
├── cmd/           # Command entry points
├── internal/      # Private application code
│   ├── cli/       # Bubble Tea UI components
│   ├── generator/ # Template and project generators
│   └── config/    # Configuration structures
├── pkg/           # Public packages (if any)
├── templates/     # Project templates
└── main.go        # Main entry point
```

### Testing
- Write table-driven tests for complex logic
- Use `testing.T.Run()` for subtests
- Mock external dependencies
- Aim for >80% test coverage on critical paths

### Dependencies
- Minimize external dependencies
- Use Go modules (`go mod`)
- Prefer standard library when possible
- Current key dependencies:
  - github.com/charmbracelet/bubbletea (TUI framework)
  - github.com/charmbracelet/lipgloss (styling)
  - github.com/charmbracelet/bubbles (UI components)

### Performance Considerations
- Use goroutines for concurrent operations (e.g., generating multiple templates)
- Implement context cancellation for long-running operations
- Buffer channels appropriately
- Profile before optimizing

### CLI Design Principles
- Keep the interface simple and intuitive
- Provide sensible defaults
- Use interactive prompts for complex configurations
- Support both interactive and non-interactive modes
- Implement --help for all commands

### Common Tasks

**Run the CLI:**
```bash
go run main.go
```

**Build:**
```bash
go build -o teapot
```

**Test:**
```bash
go test ./...
```

**Lint:**
```bash
golangci-lint run
```

### Important Notes
- This is a developer tool, so user experience in the terminal is crucial
- The Bubble Tea UI should be responsive and handle terminal resizing
- Error messages should be helpful and suggest solutions
- Keep the binary size reasonable (use `-ldflags="-s -w"` for production builds)
- Support cross-platform compilation (Linux, macOS, Windows)

### Future Considerations
- Plugin system for custom templates
- Configuration file support (.teapot.yml)
- Template marketplace/registry
- Upgrade command for existing projects
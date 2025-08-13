# dns-set Development Context

This file provides context for Claude Code when working on the dns-set project.

## Project Overview

dns-set is a command-line DNS management tool designed for VPS administrators. It automatically updates DNS records by detecting the server's IP address and managing domains from various sources.

## Architecture

### Core Components

```
cmd/dns-set/           # Main application entry point
├── main.go           # CLI setup and command routing

internal/
├── config/           # Configuration management
│   ├── config.go    # Config structure and loading
│   └── env.go       # Environment variable handling
├── domain/          # Domain source management
│   ├── manual.go    # Manual domain input
│   ├── caddyfile.go # Caddyfile parsing
│   └── source.go    # Domain source interface
├── ip/              # IP address detection
│   ├── interface.go # Network interface detection
│   ├── api.go      # External API queries (ip.sb)
│   ├── manual.go   # Manual IP input
│   └── detector.go # IP detection interface
├── dns/             # DNS provider management
│   ├── cloudflare.go # Cloudflare API integration
│   └── provider.go  # DNS provider interface
└── ui/              # User interface
    ├── cli.go       # Interactive CLI
    └── tui.go       # Future TUI implementation
```

### Key Interfaces

```go
// Domain source interface
type DomainSource interface {
    GetDomains() ([]string, error)
    Name() string
}

// IP detector interface
type IPDetector interface {
    GetIPv4() (net.IP, error)
    GetIPv6() (net.IP, error)
    Name() string
}

// DNS provider interface
type DNSProvider interface {
    UpdateRecord(domain string, recordType string, ip net.IP) error
    ListRecords(domain string) ([]Record, error)
    Name() string
}
```

## Development Guidelines

### Testing Strategy

- **Unit tests**: Use Go's standard `testing` package with `testify/assert` for assertions
- **Integration tests**: Use `testify/suite` for complex test scenarios
- **Mocking**: Use `testify/mock` for external dependencies (Cloudflare API, network calls)
- **Test structure**: Follow `_test.go` naming convention, group tests in subtests using `t.Run()`

### Code Organization

- Keep each package focused on a single responsibility
- Use dependency injection for testability
- Implement interfaces before concrete types
- Place shared types and errors in appropriate packages

### Configuration Management

- Use `viper` for configuration management (supports YAML, env vars, flags)
- Config location follows XDG Base Directory specification:
  - Linux/macOS: `~/.config/dns-set/config.yaml`
  - Windows: `%APPDATA%/dns-set/config.yaml`
  - Custom location via `DNS_SET_CONFIG_DIR` environment variable
- Environment variables use `DNS_SET_` prefix
- Support for `.env` files using `github.com/joho/godotenv`
  - Loads from current directory and config directory
  - Supports `.env` and `.env.local` files
- Support config validation and sensible defaults

### External Dependencies

Recommended Go modules:
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/cloudflare/cloudflare-go` - Cloudflare API client
- `github.com/stretchr/testify` - Testing utilities
- `github.com/joho/godotenv` - .env file loading
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/charmbracelet/bubbletea` - Future TUI framework

### Error Handling

- Use wrapped errors with context: `fmt.Errorf("failed to update DNS: %w", err)`
- Define custom error types for different failure scenarios
- Provide user-friendly error messages in CLI output
- Log detailed errors for debugging

### Cloudflare Integration

- Use official `cloudflare-go` SDK
- Support API tokens (recommended) over global API keys
- Handle rate limiting and retry logic
- Support both A and AAAA record types
- Validate zone permissions before attempting updates

### IP Detection Methods

1. **Network Interface**: Parse `ip addr` output or use `net.Interfaces()`
2. **External API**: Query `https://api.ip.sb/ip` with HTTP client
3. **Manual Input**: Interactive prompt with IP validation

### Future TUI Implementation

- Use `bubbletea` framework for reactive TUI
- Implement keyboard shortcuts for common operations
- Show real-time status updates during DNS operations
- Support multiple domain selection with checkboxes

## Build and Test Commands

```bash
# Build
go build -o dns-set ./cmd/dns-set

# Test
go test ./...

# Test with coverage
go test -cover ./...

# Lint (requires golangci-lint)
golangci-lint run

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```

## Common Development Tasks

### Adding a New DNS Provider

1. Implement the `DNSProvider` interface in `internal/dns/`
2. Add provider-specific configuration options
3. Write comprehensive unit tests with mocked API calls
4. Update CLI to include new provider option
5. Add provider setup instructions to README

### Adding a New Domain Source

1. Implement the `DomainSource` interface in `internal/domain/`
2. Add source-specific configuration if needed
3. Write tests covering various input scenarios
4. Update CLI menu to include new source option
5. Document usage in README

### Adding a New IP Detection Method

1. Implement the `IPDetector` interface in `internal/ip/`
2. Handle both IPv4 and IPv6 detection
3. Write tests with network mocking
4. Update CLI selection menu
5. Add method documentation

## Debugging Tips

- Use `-v` flag for verbose output during development
- Log API requests/responses when debugging provider issues
- Test with multiple domains to catch edge cases
- Verify DNS propagation with external tools (`dig`, `nslookup`)
- Use mock HTTP servers for testing external API integrations
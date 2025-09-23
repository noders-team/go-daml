[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE.txt)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Build Status](https://img.shields.io/badge/build-passing-green.svg)]()

# Go DAML SDK

A comprehensive Go SDK for interacting with DAML ledgers, featuring a complete client library, service abstractions, and a powerful code generator for generating type-safe Go code from DAML (.dar) files.

## Overview

The Go DAML SDK provides everything needed to build Go applications that interact with DAML ledgers. It includes a full-featured client library for ledger operations, service layer abstractions for common patterns, and an integrated code generator that creates type-safe Go structs from DAML definitions.

## Features

### SDK Components
- **Complete DAML Client Library** - Full gRPC client for DAML ledger API
- **Service Layer Abstractions** - High-level services for common operations
- **Admin Services** - Package management, user management, party management
- **Ledger Services** - Command submission, event querying, state management
- **Authentication Support** - Bearer token and custom auth implementations
- **Error Handling** - Comprehensive DAML-specific error processing

### Code Generation
- **Type-safe Go code generation** from DAML definitions
- **Custom JSON serialization** for complex DAML types
- **PackageID extraction** and embedding in generated code
- **Cross-platform support** (Linux, macOS, Windows)
- **Clean CLI interface** with comprehensive error handling

## Quick Start

### Installation

#### Option 1: Build from Source
```bash
# Clone the repository
git clone https://github.com/noders-team/go-daml.git
cd go-daml

# Build the CLI tool
make build

# The binary will be available at ./bin/godaml
```

#### Option 2: Install to GOPATH
```bash
make install
# Binary will be installed to $GOPATH/bin/godaml
```

#### Option 3: Use as Go Module
```bash
go get github.com/noders-team/go-daml
```

### SDK Usage

```go
import (
    "context"
    "github.com/noders-team/go-daml/pkg/client"
)

// Create DAML client using builder pattern
bearerToken := "your-auth-token"
grpcAddress := "localhost:8080"
tlsConfig := client.TlsConfig{}

cl, err := client.NewDamlClient(bearerToken, grpcAddress).
    WithTLSConfig(tlsConfig).
    Build(context.Background())
if err != nil {
    log.Fatal().Err(err).Msg("failed to build DAML client")
}
```

### Code Generation

```bash
# Generate Go code from a DAR file
./bin/godaml --dar ./contracts.dar --output ./generated --go_package contracts

# With debug logging
./bin/godaml --dar ./contracts.dar --output ./generated --go_package main --debug
```

### CLI Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `--dar` | ✅ | Path to the DAR file |
| `--output` | ✅ | Output directory where generated Go files will be saved |
| `--go_package` | ✅ | Go package name for generated code |
| `--debug` | ❌ | Enable debug logging (default: false) |

### Help

```bash
./bin/godaml --help
```

## Development

### Prerequisites
- Go 1.23 or later
- Make

### Build Commands
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Development workflow
make dev

# Release workflow
make release

# Show all available commands
make help
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific test
go test ./internal/codegen -v -run TestGetMainDalf
```

## SDK Modules

The Go DAML SDK is organized into the following modules:

### Core SDK (`pkg/`)
- **`pkg/client/`**: DAML ledger gRPC client implementation
  - Connection management and configuration
  - gRPC client bindings for DAML Ledger API
  - Authentication and security handling

- **`pkg/service/`**: High-level service abstractions
  - **`service/ledger/`**: Ledger operations (commands, queries, events)
  - **`service/admin/`**: Administrative operations (packages, users, parties)
  - **`service/testing/`**: Testing utilities and time services

- **`pkg/model/`**: Common data models and type definitions
- **`pkg/auth/`**: Authentication mechanisms (Bearer tokens, etc.)
- **`pkg/errors/`**: DAML-specific error handling and processing

### Code Generation (`internal/codegen/`)
- **`codegen.go`**: DAR file processing and AST generation  
- **`ast.go`**: Type definitions for template structures
- **`template.go`**: Go template processing and binding
- **`source.go.tpl`**: Go code generation template

### CLI Tool (`cmd/`)
- **`cmd/main.go`**: Command-line interface for code generation

### Examples
- **`examples/admin_app/`**: Administrative operations examples
- **`examples/ledger_app/`**: Ledger interaction examples

## SDK Usage Examples

### Example 1: Basic Ledger Client
```go
package main

import (
    "context"
    "os"
    "github.com/noders-team/go-daml/pkg/client"
    "github.com/rs/zerolog/log"
)

func main() {
    // Get configuration from environment
    grpcAddress := os.Getenv("GRPC_ADDRESS")
    if grpcAddress == "" {
        grpcAddress = "localhost:8080"
    }
    
    bearerToken := os.Getenv("BEARER_TOKEN")
    if bearerToken == "" {
        log.Warn().Msg("BEARER_TOKEN environment variable not set")
    }
    
    tlsConfig := client.TlsConfig{}
    
    // Initialize DAML client using builder pattern
    cl, err := client.NewDamlClient(bearerToken, grpcAddress).
        WithTLSConfig(tlsConfig).
        Build(context.Background())
    if err != nil {
        log.Fatal().Err(err).Msg("failed to build DAML client")
    }
    
    // Use the client for ledger operations
    // Example: RunVersionService(cl), RunPackageService(cl), etc.
}
```

### Example 2: Administrative Operations
```go
package main

import (
    "context"
    "os"
    "github.com/noders-team/go-daml/pkg/client"
    "github.com/rs/zerolog/log"
)

func main() {
    grpcAddress := os.Getenv("GRPC_ADDRESS")
    if grpcAddress == "" {
        grpcAddress = "localhost:8080"
    }
    
    bearerToken := os.Getenv("BEARER_TOKEN")
    tlsConfig := client.TlsConfig{}
    
    // Build admin client
    cl, err := client.NewDamlClient(bearerToken, grpcAddress).
        WithTLSConfig(tlsConfig).
        Build(context.Background())
    if err != nil {
        log.Fatal().Err(err).Msg("failed to build DAML client")
    }
    
    // Perform administrative operations
    // Example: RunUsersManagement(cl), RunPackageManagement(cl), etc.
    // See examples/admin_app/ for complete implementations
}
```

## Architecture

### Processing Flow
1. **DAR Extraction**: Unzip and parse DAR archive
2. **Manifest Processing**: Extract metadata and main DALF file
3. **DAML-LF Parsing**: Parse protobuf-encoded DAML definitions
4. **AST Generation**: Convert to internal AST representation
5. **Type Classification**: Identify Records, Variants, Enums
6. **Code Generation**: Apply Go templates to generate type-safe code
7. **File Output**: Write formatted Go source files

## Supported DAML Types

| DAML Type | Go Representation | Notes |
|-----------|------------------|-------|
| `Party` | `PARTY` (string) | Party identifiers |
| `Text` | `TEXT` (string) | Text values |
| `Int` | `INT64` (int64) | 64-bit integers |
| `Bool` | `BOOL` (bool) | Boolean values |
| `Decimal` | `DECIMAL` (*big.Int) | Arbitrary precision decimals |
| `Date` | `DATE` (time.Time) | Date values |
| `Time` | `TIMESTAMP` (time.Time) | Timestamp values |
| `Optional T` | `OPTIONAL` (*interface{}) | Optional values |
| `[T]` | `LIST` ([]string) | Lists |
| `Map T` | `MAP` (map[string]interface{}) | Maps |
| Records | `struct` | Standard Go structs |
| Variants | `struct` with optional fields | Union types with JSON marshaling |
| Enums | `string` | Enumeration values |

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following the existing code style
4. Run tests: `make test`
5. Run linting: `make lint`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

### Development Guidelines
- Follow Go best practices and conventions
- Write comprehensive tests for new functionality
- Update documentation for API changes
- Use meaningful commit messages
- Ensure all tests pass before submitting PR

## DAML Ecosystem

This Go SDK is part of the broader DAML ecosystem:

- **[DAML](https://daml.com/)** - Digital Asset Modeling Language for smart contracts
- **[DAML Ledger API](https://docs.daml.com/app-dev/ledger-api.html)** - gRPC API for interacting with DAML ledgers  
- **[Official DAML SDKs](https://docs.daml.com/app-dev/index.html)** - Java, TypeScript, and other language bindings
- **[Canton](https://github.com/digital-asset/canton)** - High-performance DAML ledger implementation

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE.txt](LICENSE.txt) file for details.

## Support

For questions, issues, or contributions:
- Open an issue on GitHub
- Check existing documentation and examples
- Review the test files for usage patterns
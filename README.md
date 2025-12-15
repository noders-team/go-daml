[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Build Status](https://img.shields.io/badge/build-passing-green.svg)]()

# Go DAML SDK

A comprehensive Go SDK for interacting with DAML ledgers, featuring a complete client library, service abstractions, and a powerful code generator for generating type-safe Go code from DAML (.dar) files.

## Overview

The Go DAML SDK provides a comprehensive toolkit for building Go applications that interact with DAML ledgers. It follows a modular architecture separating concerns between client interaction, service abstraction, and code generation. The SDK includes a full-featured gRPC client for ledger operations, high-level service abstractions for common patterns, and an integrated code generator that creates type-safe Go code from DAML definitions.

The public SDK components are located under `pkg/`, which includes the client, service, model, auth, codec, errors, and types packages intended for direct use by application developers. The `internal/codegen/` directory contains the code generator implementation, which parses DAML-LF (DAML Language Format) binaries extracted from .dar archives and generates corresponding Go structs with proper JSON serialization logic.

## Features

### SDK Components
- **Complete DAML Client Library** - Full gRPC client for DAML Ledger API with connection management, authentication, and TLS support
- **Dual-Connection Support** - Separate connections for ledger and admin endpoints with automatic service routing
- **Service Layer Abstractions** - High-level services for common ledger and administrative operations
- **Ledger Services** - Command submission, command completion, event querying, state management, update service, package service, version service, interactive submission
- **Admin Services** - Package management, user management, party management, participant pruning, command inspection, identity provider configuration
- **Topology Services** - Topology manager read/write operations for namespace delegations, party-to-key mappings, and party-to-participant mappings
- **Authentication Support** - Bearer token authentication with automatic token injection via gRPC interceptors
- **Error Handling** - Comprehensive DAML-specific error processing with categorized error types (authorization, validation, ledger-specific, connection errors)
- **JSON Codec** - Custom JSON serialization/deserialization for complex DAML types including Records, Variants, Enums, and primitive types

### Code Generation
- **Type-safe Go code generation** from DAML definitions with proper type mapping
- **Custom JSON serialization** for complex DAML types (Records, Variants, Enums)
- **PackageID extraction** and embedding in generated code
- **Multi-version DAML-LF support** - Supports both DAML-LF v2 and v3 with automatic version detection
- **Cross-platform support** (Linux, macOS, Windows - amd64 and arm64)
- **Clean CLI interface** with comprehensive error handling and debug logging
- **Template-based generation** using Go's text/template engine for flexible code generation

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
    "github.com/rs/zerolog/log"
)

bearerToken := "your-auth-token"
grpcAddress := "localhost:6865"
adminAddress := "localhost:6866"
tlsConfig := client.TlsConfig{}

cl, err := client.NewDamlClient(bearerToken, grpcAddress).
    WithAdminAddress(adminAddress).
    WithTLSConfig(tlsConfig).
    Build(context.Background())
if err != nil {
    log.Fatal().Err(err).Msg("failed to build DAML client")
}

cl.CommandService.SubmitAndWait(ctx, commands)
cl.EventQueryService.GetActiveContracts(ctx, filter)
cl.StateService.GetLedgerEnd(ctx)
cl.UserMng.CreateUser(ctx, userRequest)
cl.PartyMng.AllocateParty(ctx, partyRequest)
cl.TopologyManagerWrite.GenerateTransactions(ctx, request)
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
  - **Builder pattern** for flexible client configuration
  - Connection management with TLS support
  - Authentication via Bearer tokens with automatic injection
  - Service factory exposing all ledger and admin services
  - gRPC interceptors for authentication and error handling

- **`pkg/service/ledger/`**: Ledger operations
  - **Command Service**: Submit commands synchronously
  - **Command Completion**: Track command completion status
  - **Command Submission**: Asynchronous command submission
  - **Event Query Service**: Query active contracts, transaction trees, flat transactions
  - **Interactive Submission**: Multi-step command submission workflows
  - **State Service**: Query ledger state and configuration
  - **Update Service**: Subscribe to ledger updates
  - **Package Service**: Query and manage DAML packages
  - **Version Service**: Get ledger API version information

- **`pkg/service/admin/`**: Administrative operations
  - **Package Management**: Upload and validate DAR packages
  - **User Management**: Create, update, list, and delete users with rights management
  - **Party Management**: Allocate and manage parties
  - **Participant Pruning**: Prune ledger history
  - **Command Inspection**: Inspect command status
  - **Identity Provider Configuration**: Configure identity providers

- **`pkg/service/topology/`**: Topology management operations
  - **Topology Manager Write**: Generate, authorize, sign, and add topology transactions
  - **Topology Manager Read**: Query namespace delegations, party-to-key mappings, party-to-participant mappings
  - **External Party Support**: Onboarding transactions for external party allocation

- **`pkg/service/testing/`**: Testing utilities
  - **Time Service**: Control ledger time for testing

- **`pkg/model/`**: Common data models and type definitions for ledger and admin operations
- **`pkg/auth/`**: Authentication mechanisms (Bearer token interceptor)
- **`pkg/codec/`**: JSON codec for DAML types with custom marshaling/unmarshaling
- **`pkg/errors/`**: DAML-specific error handling with categorized error types
- **`pkg/types/`**: DAML type system definitions

### Code Generation (`internal/codegen/`)
- **`codegen.go`**: DAR file processing, orchestration, and AST generation
- **`astgen/factory.go`**: AST generator factory with automatic DAML-LF version detection
- **`astgen/v2/`**: DAML-LF v2 AST generation implementation
- **`astgen/v3/`**: DAML-LF v3 AST generation implementation
- **`model/manifest.go`**: DAR manifest parsing and metadata extraction
- **`model/template.go`**: Template data structures for code generation
- **`template.go`**: Go template processing, binding, and file generation
- **`template_test.go`**: Template generation tests

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
    grpcAddress := os.Getenv("GRPC_ADDRESS")
    if grpcAddress == "" {
        grpcAddress = "localhost:8080"
    }
    
    bearerToken := os.Getenv("BEARER_TOKEN")
    if bearerToken == "" {
        log.Warn().Msg("BEARER_TOKEN environment variable not set")
    }
    
    tlsConfig := client.TlsConfig{}
    
    cl, err := client.NewDamlClient(bearerToken, grpcAddress).
        WithTLSConfig(tlsConfig).
        Build(context.Background())
    if err != nil {
        log.Fatal().Err(err).Msg("failed to build DAML client")
    }
    
    ctx := context.Background()
    
    version, err := cl.VersionService.GetLedgerApiVersion(ctx)
    if err != nil {
        log.Error().Err(err).Msg("failed to get ledger version")
    }
    
    ledgerEnd, err := cl.StateService.GetLedgerEnd(ctx)
    if err != nil {
        log.Error().Err(err).Msg("failed to get ledger end")
    }
    
    log.Info().Msgf("Connected to ledger API version: %s", version)
    log.Info().Msgf("Current ledger end offset: %s", ledgerEnd)
}
```

### Example 2: Command Submission and Event Querying
```go
package main

import (
    "context"
    "github.com/noders-team/go-daml/pkg/client"
    "github.com/noders-team/go-daml/pkg/model"
    "github.com/rs/zerolog/log"
)

func main() {
    cl, err := client.NewDamlClient("token", "localhost:8080").
        Build(context.Background())
    if err != nil {
        log.Fatal().Err(err).Msg("failed to build client")
    }
    
    ctx := context.Background()
    
    commands := &model.Commands{
        ApplicationId: "my-app",
        CommandId:     "cmd-001",
        Party:         "Alice",
        Commands:      []interface{}{},
    }
    
    _, err = cl.CommandService.SubmitAndWait(ctx, commands)
    if err != nil {
        log.Error().Err(err).Msg("command submission failed")
        return
    }
    
    filter := &model.TransactionFilter{}
    contracts, err := cl.EventQueryService.GetActiveContracts(ctx, filter)
    if err != nil {
        log.Error().Err(err).Msg("failed to query contracts")
        return
    }
    
    log.Info().Msgf("Retrieved %d active contracts", len(contracts))
}
```

### Example 3: Administrative Operations
```go
package main

import (
    "context"
    "os"
    "github.com/noders-team/go-daml/pkg/client"
    "github.com/noders-team/go-daml/pkg/model"
    "github.com/rs/zerolog/log"
)

func main() {
    grpcAddress := os.Getenv("GRPC_ADDRESS")
    if grpcAddress == "" {
        grpcAddress = "localhost:8080"
    }
    
    bearerToken := os.Getenv("BEARER_TOKEN")
    tlsConfig := client.TlsConfig{}
    
    cl, err := client.NewDamlClient(bearerToken, grpcAddress).
        WithTLSConfig(tlsConfig).
        Build(context.Background())
    if err != nil {
        log.Fatal().Err(err).Msg("failed to build DAML client")
    }
    
    ctx := context.Background()
    
    userReq := &model.CreateUserRequest{
        UserId: "alice@example.com",
        PrimaryParty: "Alice",
    }
    user, err := cl.UserMng.CreateUser(ctx, userReq)
    if err != nil {
        log.Error().Err(err).Msg("failed to create user")
    }
    
    partyReq := &model.AllocatePartyRequest{
        PartyIdHint: "NewParty",
        DisplayName: "New Party Display Name",
    }
    party, err := cl.PartyMng.AllocateParty(ctx, partyReq)
    if err != nil {
        log.Error().Err(err).Msg("failed to allocate party")
    }
    
    darBytes, _ := os.ReadFile("./contracts.dar")
    uploadReq := &model.UploadDarFileRequest{
        DarFile: darBytes,
    }
    _, err = cl.PackageMng.UploadDarFile(ctx, uploadReq)
    if err != nil {
        log.Error().Err(err).Msg("failed to upload DAR")
    }
    
    log.Info().Msgf("Created user: %s", user.UserId)
    log.Info().Msgf("Allocated party: %s", party.PartyId)
}
```

### Example 4: Dual-Connection with Topology Services
```go
package main

import (
    "context"
    "os"
    "github.com/noders-team/go-daml/pkg/client"
    "github.com/noders-team/go-daml/pkg/model"
    "github.com/rs/zerolog/log"
)

func main() {
    ledgerAddress := "localhost:6865"
    adminAddress := "localhost:6866"
    bearerToken := os.Getenv("BEARER_TOKEN")

    cl, err := client.NewDamlClient(bearerToken, ledgerAddress).
        WithAdminAddress(adminAddress).
        Build(context.Background())
    if err != nil {
        log.Fatal().Err(err).Msg("failed to build DAML client")
    }
    defer cl.Close()

    ctx := context.Background()

    proposals := []*model.GenerateTransactionProposal{
        {
            Operation: model.OperationAddReplace,
            Serial:    1,
            Mapping: &model.NamespaceDelegation{
                Namespace:        "namespace-fingerprint",
                TargetKey:        publicKey,
                IsRootDelegation: true,
            },
            Store: &model.StoreID{Value: "authorized"},
        },
    }

    genResp, err := cl.TopologyManagerWrite.GenerateTransactions(ctx, &model.GenerateTransactionsRequest{
        Proposals: proposals,
    })
    if err != nil {
        log.Error().Err(err).Msg("failed to generate topology transactions")
        return
    }

    listResp, err := cl.TopologyManagerRead.ListNamespaceDelegation(ctx, &model.ListNamespaceDelegationRequest{})
    if err != nil {
        log.Error().Err(err).Msg("failed to list namespace delegations")
        return
    }

    log.Info().Msgf("Generated %d topology transactions", len(genResp.GeneratedTransactions))
    log.Info().Msgf("Found %d namespace delegations", len(listResp.Results))
}
```

## Dual-Connection Architecture

The Go DAML SDK supports separate connections for **Ledger API** and **Admin API** endpoints, which is essential for Canton deployments where these services run on different ports.

### When to Use Dual Connections

Use dual connections when:
- Your Canton participant exposes Ledger API and Admin API on different ports (e.g., 6865 and 6866)
- You need to use topology management services alongside ledger operations
- You're working with external party allocation and onboarding transactions

### Connection Routing

The client automatically routes services to the appropriate connection:

**Main Connection (Ledger Address):**
- Command Service, Command Completion, Command Submission
- Event Query Service
- State Service, Update Service
- Package Service, Version Service
- Interactive Submission Service
- User Management, Party Management
- Package Management, Participant Pruning
- Command Inspection, Identity Provider Configuration

**Admin Connection (Admin Address):**
- Topology Manager Write (GenerateTransactions, Authorize, SignTransactions, AddTransactions)
- Topology Manager Read (ListNamespaceDelegation, ListPartyToKeyMapping, ListPartyToParticipant)

### Usage

```go
cl, err := client.NewDamlClient(token, "localhost:6865").
    WithAdminAddress("localhost:6866").  // Optional: if not provided, all services use main address
    Build(ctx)

cl.TopologyManagerWrite.GenerateTransactions(ctx, request)
```

**Backward Compatibility:** If `WithAdminAddress()` is not called, all services (including topology) use the main ledger address, maintaining backward compatibility with existing code.

## Architecture

### Processing Flow
1. **DAR Extraction**: Unzip and parse DAR archive using `UnzipDar()`
2. **Manifest Processing**: Extract metadata (Main-Dalf, Sdk-Version, package names) from `META-INF/MANIFEST.MF`
3. **DAML-LF Version Detection**: Determine DAML-LF version from package metadata
4. **AST Generator Selection**: Choose appropriate AST generator (v2 or v3) via factory pattern
5. **DAML-LF Parsing**: Parse protobuf-encoded DAML definitions from .dalf file
6. **AST Generation**: Convert DAML-LF structures to internal AST representation
7. **Type Classification**: Identify Records, Variants, Enums from DAML type definitions
8. **Package ID Extraction**: Extract and embed package IDs in generated code
9. **Template Application**: Apply Go templates to generate type-safe code with custom JSON marshaling
10. **File Output**: Write formatted Go source files to output directory with proper package declarations

## Supported DAML Types

The SDK's JSON codec (`pkg/codec/json_codec.go`) provides comprehensive support for DAML type serialization:

| DAML Type | Go Representation | JSON Format | Notes |
|-----------|------------------|-------------|-------|
| `Party` | `PARTY` (string) | String | Party identifiers |
| `Text` | `TEXT` (string) | String | Text values |
| `Int64` | `INT64` (int64) | Number | 64-bit integers |
| `Numeric` | `NUMERIC` (string) | String | Arbitrary precision decimals |
| `Bool` | `BOOL` (bool) | Boolean | Boolean values |
| `Date` | `DATE` (time.Time) | String (YYYY-MM-DD) | Date values |
| `Timestamp` | `TIMESTAMP` (time.Time) | String (RFC3339) | Timestamp values |
| `ContractId` | `CONTRACT_ID` (string) | String | Contract identifiers |
| `Optional T` | `OPTIONAL` (*interface{}) | Null or value | Optional values with custom unmarshaling |
| `[T]` | `LIST` ([]interface{}) | Array | Lists with element validation |
| `TextMap T` | `TEXT_MAP` (map[string]interface{}) | Object | Text-keyed maps |
| `GenMap K V` | `GEN_MAP` ([]MapEntry) | Array of {key, value} | Generic maps |
| Records | Go `struct` | Object | Standard Go structs with field mapping |
| Variants | Go `struct` with optional fields | Object with constructor tag | Union types with JSON marshaling |
| Enums | `string` | String | Enumeration values |

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

## Support

For questions, issues, or contributions:
- Open an issue on GitHub
- Check existing documentation and examples
- Review the test files for usage patterns
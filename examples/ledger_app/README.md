# GO-DAML Ledger App (example)

A comprehensive example application demonstrating all DAML Ledger API services using the go-daml SDK.

## Overview

This application showcases the complete DAML ledger functionality in a single executable, including:

- **Version Service** - Get ledger API version and feature information
- **Package Service** - List packages, get package details and status
- **State Service** - Get ledger end, active contracts, and pruned offsets
- **Update Service** - Stream transaction updates from the ledger
- **Command Completion** - Stream command completion events
- **Command Service** - Submit commands and wait for completion
- **Command Submission** - Submit commands asynchronously
- **Event Query** - Query events by contract ID
- **Interactive Submission** - Prepare and execute transactions interactively
- **Time Service** - Get and set ledger time (testing only)

## Structure

```
examples/ledger_app/
├── main.go                    # Entry point and client setup
├── helpers.go                 # Shared helper functions for users and parties
├── version_service.go         # Version service operations
├── package_service.go         # Package service operations
├── state_service.go           # State service operations
├── update_service.go          # Update service operations
├── command_completion.go      # Command completion streaming
├── command_service.go         # Command service operations
├── command_submission.go      # Command submission operations
├── event_query.go             # Event query operations
├── interactive_submission.go  # Interactive submission operations
├── time_service.go           # Time service operations (testing)
└── README.md                 # This file
```

## Running the Application

### Prerequisites

- Go 1.21 or later
- Access to a DAML ledger with ledger API enabled
- Valid bearer token for authentication

### Environment Variables

Set the following environment variables:

```bash
export GRPC_ADDRESS="your-ledger-address:port"
export BEARER_TOKEN="your-jwt-token"
```

### Running

From the project root:

```bash
go run ./examples/ledger_app/*.go
```

Or from the ledger_app directory:

```bash
cd examples/ledger_app
go run *.go
```

### Example with Canton Network

```bash
GRPC_ADDRESS="grpc-ledger.canton-localnet.noders.services:443" \
BEARER_TOKEN="your-jwt-token" \
go run ./examples/ledger_app/*.go
```

## Features Demonstrated

### Version Service
- Get ledger API version information
- Check supported features (user management, party management, offset checkpoint)

### Package Service
- List all known packages
- Get package details including archive payload and hash
- Check package status

### State Service
- Get current ledger end offset
- Get latest pruned offsets
- List connected synchronizers
- Stream active contracts (with timeout)

### Update Service
- Stream transaction updates from ledger
- Handle different update types (transactions, reassignments, checkpoints)
- Process updates with timeout

### Command Completion
- Stream command completion events
- Handle completion responses and offset checkpoints

### Command Service
- Submit commands and wait for completion
- Handle empty command submissions (for testing)

### Command Submission
- Submit commands asynchronously
- Handle submission without waiting for completion

### Event Query
- Query events by contract ID
- Handle created and archived events
- Demonstrate error handling for non-existent contracts

### Interactive Submission
- Get preferred package versions
- Prepare submissions with validation
- Execute prepared transactions
- Handle transaction hashing and signatures

### Time Service
- Get current ledger time (testing ledgers only)
- Set ledger time (testing ledgers only)
- Handle time service operations with proper error handling

## Logging

The application uses structured logging with zerolog:
- **INFO** level for normal operations and results
- **WARN** level for expected errors and warnings
- **FATAL** level for critical errors that stop execution

## Error Handling

The application demonstrates proper error handling patterns:
- Fatal errors stop execution (connection issues, critical API failures)
- Non-fatal errors are logged but allow continuation
- Expected errors (like querying non-existent contracts) are handled gracefully
- Stream timeouts are handled appropriately

## Notes

- Many operations may result in expected errors when running against empty ledgers
- Stream operations include timeouts to prevent hanging
- Command submissions use empty command lists for demonstration purposes
- Time service operations only work with testing ledgers (e.g., Canton sandbox)
- Interactive submission demonstrates the full prepare-execute cycle
- All examples use the same client connection for efficiency
- **User vs Party Distinction**: Examples properly distinguish between UserID (for authentication) and Party (for authorization)
  - `UserID` identifies the authenticated user making the request
  - `ActAs` and filtering use Party identifiers for authorization
  - Helper functions dynamically retrieve both from the ledger system
- **Filter Requirements**: TransactionFilter requires at least one party filter (cannot have empty filtersByParty and filtersForAnyParty simultaneously)

## Development

To add new ledger functionality:

1. Create a new `.go` file with your service functions
2. Add a `Run*Service(cl *client.DamlBindingClient)` function
3. Call your function from `main.go`
4. Follow the existing logging and error handling patterns

## Comparison with Admin App

This ledger_app complements the admin_app example:
- **admin_app** - Covers administrative operations (users, parties, packages, pruning, identity providers)
- **ledger_app** - Covers ledger operations (commands, events, state, updates, interactive submission)

Together, these examples demonstrate the complete DAML Ledger API surface available through the go-daml SDK.
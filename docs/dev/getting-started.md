# Getting Started

This guide takes you from an empty project to a running DAML contract create →
query → exercise flow with the Go DAML SDK.

- [Requirements](#requirements)
- [Install](#install)
- [Build the godaml CLI](#build-the-godaml-cli)
- [Connecting a client](#connecting-a-client)
- [Authentication](#authentication)
- [Generating Go code from a .dar](#generating-go-code-from-a-dar)
- [Your first contract: create → query → exercise](#your-first-contract-create--query--exercise)
- [What to read next](#what-to-read-next)

## Requirements

- **Go 1.25.0+** (the module declares `go 1.25.0`).
- A reachable DAML / Canton participant exposing the gRPC **Ledger API**
  (default Canton port `6865`) and, optionally, an **admin** endpoint
  (default `6866`) for topology operations.
- A **bearer token** accepted by that participant (or Keycloak client
  credentials — see [Authentication](#authentication)).

## Install

Add the SDK to your application module:

```bash
go get github.com/noders-team/go-daml
```

Then import the packages you need:

```go
import (
    "github.com/noders-team/go-daml/pkg/auth"
    "github.com/noders-team/go-daml/pkg/client"
    "github.com/noders-team/go-daml/pkg/model"
    . "github.com/noders-team/go-daml/pkg/types"
)
```

The dot-import of `pkg/types` is the convention used by generated code so that
DAML primitives read naturally (`PARTY("Alice")`, `TEXT("hi")`,
`TEXTMAP{...}`). It is optional — you can import it qualified instead.

## Build the godaml CLI

The code generator is a separate binary. Building from a checkout of this repo:

```bash
make build        # -> ./bin/godaml   (current platform)
make build-all    # cross-compile Linux / macOS / Windows, amd64 + arm64
make install      # -> $GOPATH/bin/godaml
make dev          # fmt + vet + test + build (run before declaring work done)
```

Or pull it in as a Go tool dependency of your own module:

```bash
go get -tool github.com/noders-team/go-daml/cmd/go-daml
go tool go-daml --dar ./contracts.dar --output ./generated --go_package contracts
```

## Connecting a client

The client uses a builder. `NewDamlClient` takes the **ledger gRPC address** and
an **auth token provider**; `Build` opens the connection(s) and returns a
`*client.DamlBindingClient` carrying every service.

```go
package main

import (
    "context"

    "github.com/noders-team/go-daml/pkg/auth"
    "github.com/noders-team/go-daml/pkg/client"
)

func main() {
    ctx := context.Background()

    cl, err := client.NewDamlClient(
        "localhost:6865",                          // ledger gRPC endpoint (required)
        auth.NewBearerTokenProvider("<token>"),    // token provider (required)
    ).
        WithAdminAddress("localhost:6866"). // optional: separate admin endpoint for topology
        Build(ctx)
    if err != nil {
        panic(err)
    }
    defer cl.Close()

    // sanity check — round-trips the version service
    if err := cl.Ping(ctx); err != nil {
        panic(err)
    }
}
```

### Builder reference

| Method | Required | Purpose |
| --- | --- | --- |
| `NewDamlClient(addr, provider)` | yes | Ledger gRPC address + `auth.TokenProvider`. |
| `WithAdminAddress(addr)` | no | Separate endpoint used for **topology** read/write and traffic control. If omitted, those calls fall back to the ledger address. |
| `WithTLSConfig(client.TlsConfig{...})` | no | Switches the connection to TLS. When set, the token is sent via gRPC per-RPC credentials instead of an interceptor. |
| `Build(ctx)` | yes | Opens the gRPC connection(s) and returns `*DamlBindingClient`. |

**Dual connection.** The client holds two gRPC connections: one for the ledger
endpoint and one for the admin endpoint. Ledger / command / state / event /
update / user / party / package calls go over the ledger connection; topology
and traffic-control calls go over the admin connection (or fall back to the
ledger connection when no admin address is configured).

**Lifecycle.** Always `defer cl.Close()` after a successful `Build` — it closes
both connections.

### TLS

```go
cl, err := client.NewDamlClient("ledger.example.com:6865", provider).
    WithTLSConfig(client.TlsConfig{Certificate: "/path/to/ca.pem"}).
    Build(ctx)
```

Without `WithTLSConfig`, the SDK connects with insecure transport credentials —
fine for a local sandbox, **not** for production.

## Authentication

The token provider is the only required auth input. Two providers ship with the
SDK; both satisfy `auth.TokenProvider` and inject `authorization: Bearer <token>`
on every call.

### Static bearer token

```go
provider := auth.NewBearerTokenProvider(os.Getenv("BEARER_TOKEN"))
```

Use this when you already hold a token (CI, manual testing, an upstream service
that mints tokens for you).

### Keycloak (client-credentials, auto-refreshing)

```go
provider, err := auth.NewKeycloakTokenProvider(auth.KeycloakConfig{
    OIDCURL:      os.Getenv("KEYCLOAK_OIDC_URL"),   // realm URL; token endpoint is derived
    TokenURL:     os.Getenv("KEYCLOAK_TOKEN_URL"),  // optional: overrides OIDCURL-derived endpoint
    ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
    ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
    Audience:     os.Getenv("KEYCLOAK_AUDIENCE"),    // optional
})
if err != nil {
    panic(err)
}

cl, err := client.NewDamlClient("localhost:6865", provider).Build(ctx)
```

The Keycloak provider caches the access token and refreshes it before expiry, so
long-running services don't need to re-authenticate manually. You must supply
either `OIDCURL` or `TokenURL`; if both are empty the constructor returns an
error.

## Generating Go code from a .dar

Turn a DAML archive into type-safe Go structs:

```bash
godaml --dar ./contracts.dar --output ./generated --go_package contracts
```

For each DAML template you get a struct implementing the `Template` interface
(`CreateCommand() *model.CreateCommand` and `GetTemplateID() string`), plus a
method per choice; records, variants and enums become typed Go values with
codec-backed JSON marshaling. The full reference — every generated shape, the
two-pass generation flow, DAML-LF version handling, and naming rules — lives in
**[code-generation.md](./code-generation.md)**.

> **Template ID note.** `GetTemplateID()` returns a *package-name* reference
> (`#package-name:Module:Entity`). Some participants require a concrete
> *package-id* reference. See
> [troubleshooting → template/package IDs](./troubleshooting.md#wrong-or-unresolved-template-id).

The examples below use a generated `MappyContract` template:

```go
type MappyContract struct {
    Operator PARTY   `json:"operator"`
    Value    TEXTMAP `json:"value"`
}
```

## Your first contract: create → query → exercise

This is the full happy path using the generated `MappyContract` above. It
uploads the DAR (admin right required), creates a contract, reads it back with
the typed `ContractQuery`, then archives it.

```go
ctx := context.Background()

// 0. Resolve the acting party from the authenticated user.
users, err := cl.UserMng.ListUsers(ctx)
// ... pick the user you authenticated as; party := user.PrimaryParty

party := "Alice::1220..."

// 1. Build a contract value and its create command.
contract := MappyContract{
    Operator: PARTY(party),
    Value:    TEXTMAP{"k1": "v1"},
}
createCmd := contract.CreateCommand()

// 2. Submit and wait for completion.
resp, err := cl.CommandService.SubmitAndWait(ctx, &model.SubmitAndWaitRequest{
    Commands: &model.Commands{
        WorkflowID:   "first-contract",
        UserID:       "alice",
        CommandID:    "create-1",
        ActAs:        []string{party},
        SubmissionID: "create-sub-1",
        DeduplicationPeriod: model.DeduplicationDuration{Duration: 60 * time.Second},
        Commands: []*model.Command{
            {Command: createCmd},
        },
    },
})
if err != nil {
    panic(err)
}
log.Printf("created, updateID=%s", resp.UpdateID)

// 3. Query active contracts of this template, decoded into MappyContract.
query := client.NewContractQuery[MappyContract](cl)
found, err := query.FindContractsByTemplate(ctx, party, contract.GetTemplateID())
if err != nil {
    panic(err)
}

// 4. Archive each one.
for _, c := range found {
    archiveCmd := contract.Archive(c.ContractID)
    _, err := cl.CommandService.SubmitAndWait(ctx, &model.SubmitAndWaitRequest{
        Commands: &model.Commands{
            UserID:       "alice",
            CommandID:    "archive-" + c.ContractID,
            ActAs:        []string{party},
            SubmissionID: "archive-sub-" + c.ContractID,
            DeduplicationPeriod: model.DeduplicationDuration{Duration: 60 * time.Second},
            Commands: []*model.Command{
                {Command: archiveCmd},
            },
        },
    })
    if err != nil {
        panic(err)
    }
}
```

Key shapes to remember:

- A `*model.CreateCommand` / `*model.ExerciseCommand` is wrapped in
  `model.Command{Command: cmd}` before it goes into `Commands.Commands`.
- `SubmitAndWait` takes a `*model.SubmitAndWaitRequest` whose single field is
  `Commands *model.Commands`.
- `ContractQuery[T]` decodes the create arguments of each active contract into
  your generated struct `T`, returning `[]client.Contract[T]` where each item
  carries `ContractID` and `Data`.

## What to read next

- [code-generation.md](./code-generation.md) — the full `godaml` reference:
  every generated type, the two-pass flow, DAML-LF versions, and naming rules.
- [minimal-examples.md](./minimal-examples.md) — focused snippets for every
  service (state, events, updates, completions, admin, topology, interactive
  submission).
- [troubleshooting.md](./troubleshooting.md) — what the common failures look
  like and how to fix them.

# Go DAML SDK — Developer Documentation

`github.com/noders-team/go-daml` is a Go toolkit for building applications that
talk to DAML / Canton ledgers over the gRPC Ledger API. It ships three things in
one module:

1. A gRPC **client library** for the DAML Ledger API (`pkg/client`).
2. **Service abstractions** over ledger, admin, and topology operations (`pkg/service`).
3. A **code generator** (`godaml` CLI) that turns DAML `.dar` files into
   type-safe Go structs.

## Contents

| Document | What it covers |
| --- | --- |
| [getting-started.md](./getting-started.md) | Install, build, connect a client, create your first contract. |
| [code-generation.md](./code-generation.md) | The `godaml` CLI: every generated shape (templates, choices, records, variants, enums), the two-pass flow, DAML-LF versions, naming rules. |
| [minimal-examples.md](./minimal-examples.md) | Copy-paste minimal snippets for every service: commands, state, events, updates, admin, topology. |
| [compatibility.md](./compatibility.md) | Compatibility matrix: Go versions, Ledger API / DAML-LF / Canton ranges, and Python DAZL interop notes (LF/protobuf and JSON/OpenAPI flows). |
| [troubleshooting.md](./troubleshooting.md) | Common integration failures (auth, connection, template IDs, streaming) and how to diagnose them. |

## Quick map of the SDK

```
pkg/client    builder + dual gRPC connection (ledger + admin) + typed ContractQuery
pkg/auth      bearer-token and Keycloak token providers (gRPC interceptors)
pkg/service   ledger / admin / topology service interfaces
pkg/model     request & response structs for every call
pkg/types     DAML primitive types (PARTY, TEXT, INT64, NUMERIC, TEXTMAP, ...)
pkg/codec     custom JSON marshaling for DAML records / variants / enums
pkg/errors    DAML error classification (AsDamlError)
cmd/go-daml   the godaml code-generator CLI
```

All service calls flow through a single `*client.DamlBindingClient` returned by
`Build(ctx)`. Every service is a public field on that struct
(`cl.CommandService`, `cl.StateService`, `cl.UserMng`, `cl.TopologyManagerRead`, …).

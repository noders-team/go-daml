# Compatibility Matrix

What this SDK is built against, which DAML / Canton ledgers it talks to, which
`.dar` files the code generator accepts, and how it relates to the Python
**DAZL** client for cross-language interop.

- [At a glance](#at-a-glance)
- [Go toolchain & dependencies](#go-toolchain--dependencies)
- [Ledger API (runtime / gRPC)](#ledger-api-runtime--grpc)
- [Code generation: DAML-LF & SDK ranges](#code-generation-daml-lf--sdk-ranges)
- [Canton / DAML platform support](#canton--daml-platform-support)
- [Python DAZL interop notes](#python-dazl-interop-notes)
  - [LF / protobuf (gRPC) flow](#lf--protobuf-grpc-flow)
  - [JSON / OpenAPI flow](#json--openapi-flow)
- [How to verify on your stack](#how-to-verify-on-your-stack)

> **Source of truth.** The version numbers below are read directly from this
> repo (`go.mod`, the codegen version selector, and the golden test fixtures).
> Anything about *upstream* DAML/Canton/DAZL releases is a compatibility *note* —
> confirm against the relevant release notes for your deployment. See
> [How to verify on your stack](#how-to-verify-on-your-stack).

## At a glance

Legend: ✓ supported · ✗ not supported · ⚠ conditional (see notes).

| Axis | Value | Status |
| --- | --- | :---: |
| Go toolchain | 1.25.0+ (module declares `go 1.25.0`) | ✓ |
| Ledger API (gRPC) | v2 (`com/daml/ledger/api/v2`) | ✓ |
| Ledger API (gRPC) | v1 (legacy / sandbox) | ✗ |
| DAML-LF — codegen | 1.x (astgen v2) | ✓ |
| DAML-LF — codegen | 2.x (astgen v3) | ✓ |
| DAML SDK of the `.dar` | 1.0 – 1.18, 2.0 – 2.10 → LF1 path (see [full list](#code-generation-daml-lf--sdk-ranges)) | ✓ |
| DAML SDK of the `.dar` | 3.0 – 3.4 → LF2 path (see [full list](#code-generation-daml-lf--sdk-ranges)) | ✓ |
| Verified `.dar` fixtures | SDK 1.18.1, 2.9.1, 3.3.0-snapshot | ✓ |
| Target ledger | Canton (Ledger API v2, incl. Splice / Canton Network) | ✓ |
| Transport | gRPC, plaintext or TLS | ✓ |
| Auth | bearer-token / Keycloak (OIDC) | ✓ |

## Go toolchain & dependencies

From `go.mod`:

| Dependency | Version | Role |
| --- | --- | --- |
| `go` directive | `1.25.0` | Minimum Go toolchain to build the SDK / `godaml`. |
| `github.com/digital-asset/dazl-client/v8` | `v8.9.0` | Generated Go protobuf bindings for the DAML Ledger API and DAML-LF archive format. Shared lineage with the Python DAZL client. |
| `google.golang.org/grpc` | `v1.80.0` | gRPC transport. |
| `google.golang.org/protobuf` | `v1.36.11` | Protobuf runtime. |
| `github.com/shopspring/decimal` | `v1.4.0` | `NUMERIC` / `DECIMAL` handling in `pkg/types`. |

The repo **vendors** its dependencies (`vendor/`). Use the published
`dazl-client v8.9.0` — a local checkout of dazl introduces a proto import cycle.
Run `make deps` to tidy.

**Go version support** (`go` directive is `1.25.0`; Go enforces forward
compatibility, so anything below is rejected):

| Go version | Supported |
| --- | :---: |
| 1.21.x | ✗ |
| 1.22.x | ✗ |
| 1.23.x | ✗ |
| 1.24.x | ✗ |
| **1.25.x** | ✓ |
| 1.26.x and newer | ✓ |

## Ledger API (runtime / gRPC)

Every service in `pkg/service/**` is built on the **Ledger API v2** protobuf
package (`github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2`).
Ledger API v2 is the gRPC API exposed by Canton participant nodes (Canton 2.x and
3.x). This SDK does **not** target the legacy Ledger API v1.

| Ledger API | Supported |
| --- | :---: |
| v1 (`com/daml/ledger/api/v1`, legacy / sandbox) | ✗ |
| **v2** (`com/daml/ledger/api/v2`, Canton) | ✓ |

Practical implications:

- Your participant must expose the **v2** gRPC Ledger API (and, for topology /
  admin operations, the Admin API — typically a separate port; see
  [getting-started → connecting a client](./getting-started.md#connecting-a-client)).
- Features such as interactive submission, synchronizers, reassignments, and
  offset checkpoints are v2 concepts and are surfaced as-is.

## Code generation: DAML-LF & SDK ranges

`godaml` auto-detects the DAML-LF generation from the DAR manifest's
`Sdk-Version` and selects one of two AST generators. The mapping (from the
version selector in `internal/codegen`) is:

| `Sdk-Version` prefix in DAR manifest | astgen | DAML-LF archive proto | DAML-LF generation |
| --- | --- | --- | --- |
| `1.` | v2 | `daml_lf_1` (header proto `1.17`) | **LF 1.x** |
| `2.` | v2 | `daml_lf_1` (header proto `1.17`) | **LF 1.x** |
| `3.` | v3 | `daml_lf_2` (header proto `2.1`) | **LF 2.x** |

**DAML SDK version of the `.dar`.** Selection is by the `Sdk-Version` **major
prefix**, so every minor/patch in a supported line is handled identically. The
full set of released DAML SDK minor versions, each with its status:

| DAML SDK | Generator | DAML-LF | Supported |
| --- | --- | --- | :---: |
| 1.0 | v2 | LF 1.x | ✓ |
| 1.1 | v2 | LF 1.x | ✓ |
| 1.2 | v2 | LF 1.x | ✓ |
| 1.3 | v2 | LF 1.x | ✓ |
| 1.4 | v2 | LF 1.x | ✓ |
| 1.5 | v2 | LF 1.x | ✓ |
| 1.6 | v2 | LF 1.x | ✓ |
| 1.7 | v2 | LF 1.x | ✓ |
| 1.8 | v2 | LF 1.x | ✓ |
| 1.9 | v2 | LF 1.x | ✓ |
| 1.10 | v2 | LF 1.x | ✓ |
| 1.11 | v2 | LF 1.x | ✓ |
| 1.12 | v2 | LF 1.x | ✓ |
| 1.13 | v2 | LF 1.x | ✓ |
| 1.14 | v2 | LF 1.x | ✓ |
| 1.15 | v2 | LF 1.x | ✓ |
| 1.16 | v2 | LF 1.x | ✓ |
| 1.17 | v2 | LF 1.x | ✓ |
| 1.18 ✦ | v2 | LF 1.x | ✓ |
| 2.0 | v2 | LF 1.x | ✓ |
| 2.1 | v2 | LF 1.x | ✓ |
| 2.2 | v2 | LF 1.x | ✓ |
| 2.3 | v2 | LF 1.x | ✓ |
| 2.4 | v2 | LF 1.x | ✓ |
| 2.5 | v2 | LF 1.x | ✓ |
| 2.6 | v2 | LF 1.x | ✓ |
| 2.7 | v2 | LF 1.x | ✓ |
| 2.8 | v2 | LF 1.x | ✓ |
| 2.9 ✦ | v2 | LF 1.x | ✓ |
| 2.10 | v2 | LF 1.x | ✓ |
| 3.0 | v3 | LF 2.x | ✓ |
| 3.1 | v3 | LF 2.x | ✓ |
| 3.2 | v3 | LF 2.x | ✓ |
| 3.3 ✦ | v3 | LF 2.x | ✓ |
| 3.4 | v3 | LF 2.x | ✓ |
| 0.x / any other major | — | — | ✗ (`none supported version`) |

✦ = exercised by a golden fixture in `test-data/` (SDK 1.18.1, 2.9.1,
3.3.0-snapshot) and known to generate. Newer 3.x snapshots (3.5, 3.6) match the
`3.` prefix and use the v3 generator; LF2 decoding still depends on the vendored
`dazl-client v8.9.0` proto (see the DAML-LF table below).

**DAML-LF version** (the archive decoder is chosen by SDK line; the `daml_lf_1`
decoder reads the whole LF 1.x range, `daml_lf_2` reads LF 2.x):

| DAML-LF version | Decoder proto | Supported |
| --- | --- | :---: |
| 1.6 – 1.17 (LF 1.x) | `daml_lf_1` (header `1.17`) | ✓ |
| 2.1 (LF 2.x) | `daml_lf_2` (header `2.1`) | ✓ |
| 2.x newer than the bundled `daml_lf_2` header | `daml_lf_2` | ⚠ decode depends on the vendored `dazl-client v8.9.0` proto — verify with `--debug` |

stdlib / prim DALFs inside the archive are skipped. A DAR whose `Sdk-Version`
doesn't start with `1.`, `2.`, or `3.` is rejected with `none supported version`.
Run with `--debug` to see the detected version and per-DALF steps.

## Canton / DAML platform support

| Capability | Backed by | Supported |
| --- | --- | :---: |
| Ledger / command / state / event / update | Canton Ledger API v2 | ✓ |
| Topology read/write & traffic control | Canton Admin API | ✓ |
| Interactive submission (external signing) | Ledger API v2 | ✓ |
| Splice / Canton Network (Amulet, token standard) | Canton/DAML 3.3 — see [transfer CC](./minimal-examples.md#transfer-canton-coin-cc--amulet) | ✓ |
| OIDC bearer-token auth | `pkg/auth` | ✓ |
| Keycloak auto-refresh auth | `pkg/auth` | ✓ |

**Network / deployment** — all supported as long as the participant exposes the
v2 Ledger API and accepts your token:

| Environment | Supported |
| --- | :---: |
| LocalNet / sandbox | ✓ |
| DevNet | ✓ |
| TestNet | ✓ |
| MainNet | ✓ |

## Python DAZL interop notes

[DAZL](https://github.com/digital-asset/dazl-client) is Digital Asset's **Python**
client for DAML ledgers. This Go SDK depends on the **same `dazl-client/v8`
repository** — specifically its generated **Go** protobuf bindings
(`dazl-client/v8/go/api/...`). That shared lineage is what makes Go ⇄ Python
interop predictable: both sides speak the same wire contract.

There are two interop surfaces. Pick based on what your Python side uses.

### LF / protobuf (gRPC) flow

Both this SDK and Python DAZL talk the **gRPC Ledger API** and exchange
**DAML-LF**-typed payloads encoded as protobuf.

- **Wire contract is shared.** Commands submitted by a Go service and commands
  submitted by a Python DAZL service against the same participant are
  indistinguishable on the ledger — same templates, choices, contract ids.
- **DAML-LF must match the DAR.** A contract created from an LF2 (SDK 3.x) DAR
  must be consumed with code generated from that same DAR on both sides. The Go
  side uses `godaml` (v3 generator for LF2); the Python side uses DAML's Python
  codegen / DAZL typing for the corresponding LF.
- **Template-id discipline.** Both clients must agree on the template id
  reference (package-name vs concrete package-id). See
  [troubleshooting → wrong or unresolved template id](./troubleshooting.md#wrong-or-unresolved-template-id).
- **Ledger API version.** This SDK is **v2-only**. Ensure the Python DAZL release
  you pair with also targets the v2 Ledger API against the same participant.

> **Verify upstream.** DAZL's own DAML-LF / SDK support window is governed by the
> DAZL release you install, not by this repo. Confirm your DAZL version supports
> the LF generation (1.x vs 2.x) of the DARs you share. Treat LF2 (Canton 3.x /
> Splice) support as something to confirm against the DAZL release notes.

### JSON / OpenAPI flow

Canton also exposes the **JSON Ledger API** (HTTP + JSON, described by an
OpenAPI document). This is the most decoupled interop path and is language- and
LF-codegen-agnostic.

- **No shared protobuf needed.** A Python (or any-language) component can drive
  the ledger over HTTP/JSON while this Go SDK drives it over gRPC — they only need
  to agree on template ids and argument shapes.
- **JSON shape parity.** This SDK's `pkg/codec` produces DAML-canonical JSON for
  records, variants, enums, and primitives (e.g. `NUMERIC`/`INT64` as strings for
  precision, variants as `{tag, value}`, timestamps as ISO-8601). That canonical
  encoding is what lines up with the JSON Ledger API / OpenAPI representation, so
  values round-trip cleanly between the gRPC and JSON flows.
- **When to choose it.** Prefer the JSON/OpenAPI flow for light integrations,
  webhooks, or polyglot stacks where you don't want to maintain generated DAML
  bindings on the non-Go side; prefer the gRPC/protobuf flow for high-throughput,
  streaming, or strongly-typed services.

> **Verify upstream.** The JSON Ledger API and its OpenAPI document are provided
> by the Canton participant version you run. Generate the OpenAPI client from
> *your* participant's published spec rather than assuming a fixed shape.

## How to verify on your stack

1. **Ledger version** — call `cl.VersionService.GetLedgerAPIVersion(ctx, ...)` (or
   the `cl.Ping(ctx)` / `cl.ValidateSDKVersion(ctx, sdkVersion)` helpers) to read
   what the participant actually advertises.
2. **DAR generation** — run `godaml --dar <file> ... --debug`; the log shows the
   detected `Sdk-Version` and whether the v2 (LF1) or v3 (LF2) generator was used.
3. **Go toolchain** — `go version` must report **1.25.0+**.
4. **DAZL pairing** — check the installed DAZL release's own compatibility notes
   for the Ledger API version and DAML-LF range it supports, and confirm it
   matches the participant and DARs you share with the Go side.

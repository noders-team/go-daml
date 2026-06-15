# Troubleshooting

Common integration failures with the Go DAML SDK, what they look like, and how to
fix them. Start with [classifying the error](#classifying-any-error) — it tells
you which section below applies.

- [Classifying any error](#classifying-any-error)
- [Build / connection failures](#build--connection-failures)
- [Authentication failures](#authentication-failures)
- [Permission denied](#permission-denied)
- [Wrong or unresolved template id](#wrong-or-unresolved-template-id)
- [Command submission rejected](#command-submission-rejected)
- [Streaming calls hang or return nothing](#streaming-calls-hang-or-return-nothing)
- [TLS problems](#tls-problems)
- [Topology / admin calls fail](#topology--admin-calls-fail)
- [Codegen failures](#codegen-failures)

---

## Classifying any error

Every gRPC error returned by the SDK can be normalized with
`errors.AsDamlError`. It never returns nil and exposes a `CategoryID` you can
branch on.

```go
import "github.com/noders-team/go-daml/pkg/errors"

if err != nil {
    de := errors.AsDamlError(err)
    log.Printf("code=%s category=%d corr=%v msg=%s",
        de.ErrorCode, de.CategoryID, de.CorrelationID, de.Message)
}
```

| `CategoryID` | `ErrorCode` | Meaning | Where to look |
| --- | --- | --- | --- |
| `> 0` | the real DAML code (e.g. `PACKAGE_NAMES_NOT_FOUND`) | A genuine DAML Ledger API error. The participant rejected the request for a domain reason. | The `ErrorCode` and `CorrelationID` — grep the participant logs for the correlation id. |
| `-1` | `DAML_GENERIC_ERROR_CODE` | `nil` was passed to `AsDamlError`. | Nothing — there was no error. |
| `-2` | `DAML_GENERIC_ERROR_CODE` | Not a gRPC error — a plain Go error (often connection/network). | [Build / connection failures](#build--connection-failures). |
| `-3` | `DAML_GENERIC_ERROR_CODE` | A DAML-formatted message whose category id couldn't be parsed. | Treat like a real DAML error; read `Message`. |
| `-5` | `DAML_GENERIC_ERROR_CODE` | A gRPC error **without** DAML's `CODE(cat,corr): msg` envelope — e.g. `Unauthenticated`, `PermissionDenied`, `InvalidArgument`, `DeadlineExceeded`. | The underlying gRPC status — inspect with `status.FromError(err)`. |

For `-5` errors, drop down to the raw gRPC status to get the precise code:

```go
import "google.golang.org/grpc/status"

if st, ok := status.FromError(err); ok {
    log.Printf("grpc code=%s msg=%s", st.Code(), st.Message())
}
```

---

## Build / connection failures

**Symptoms.** `Build(ctx)` returns an error such as
`failed to connect to DAML ledger: ...` or
`failed to connect to DAML admin endpoint: ...`. At call time you get a `-2`
DAML error or a gRPC `Unavailable` status.

**Causes & fixes.**

- **Wrong host/port.** The ledger gRPC port is typically `6865`, the admin port
  `6866`. Confirm what your participant actually exposes — don't assume.
- **Participant not running / not reachable.** Verify with `cl.Ping(ctx)` right
  after `Build`. If `Ping` fails the connection is the problem, not your command.
- **Admin endpoint missing.** If only topology/traffic calls fail but ledger
  calls work, you probably didn't set `WithAdminAddress`, and the fallback
  endpoint doesn't serve the admin API.
- **No deadline.** `grpc.NewClient` connects lazily; a bad address may not
  surface until the first RPC. Always pass a `context` with a timeout to your
  calls so failures are bounded:

  ```go
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  ```

---

## Authentication failures

**Symptoms.** Calls fail immediately with gRPC `Unauthenticated` (a `-5` DAML
error). The message often mentions a missing or invalid `authorization` header,
or an expired/invalid token.

**Causes & fixes.**

- **Empty token.** `NewBearerTokenProvider("")` sends no `authorization` header.
  Check the env var you read the token from is actually set.
- **Expired token.** Static bearer tokens don't refresh. For long-running
  services use `auth.NewKeycloakTokenProvider`, which refreshes automatically.
- **Wrong audience / issuer.** The participant's IDP config must match your
  token's `aud`/`iss`. If you administer the participant, verify with
  `cl.IdentityProviderMng.ListIdentityProviderConfigs(ctx)`; the `Issuer`,
  `JwksURL` and `Audience` must line up with how the token was minted.
- **Keycloak provider constructor error.** `NewKeycloakTokenProvider` returns
  `keycloak token url is empty` when both `OIDCURL` and `TokenURL` are blank —
  supply at least one.
- **Token sent on the wrong channel.** In insecure mode the token is injected
  via gRPC interceptors; in TLS mode via per-RPC credentials. Both are handled
  automatically — but per-RPC credentials over an insecure channel are dropped by
  gRPC. If you set neither TLS nor saw the header arrive, double-check you didn't
  hand-roll dial options that override the SDK's.

---

## Permission denied

**Symptoms.** gRPC `PermissionDenied` (a `-5` DAML error) on admin calls
(`UserMng`, `PartyMng`, `PackageMng`, pruning) or on commands acting as a party.

**Causes & fixes.**

- **User lacks the right.** Admin operations need a `ParticipantAdmin` user;
  acting as a party needs `CanActAs` for that party. Inspect with
  `cl.UserMng.ListUserRights(ctx, userID)`.
- **`ActAs` party mismatch.** The party in `Commands.ActAs` must be one the
  authenticated user is authorized to act as. Resolve the user's primary party
  from `cl.UserMng.GetUser` / `ListUsers` (`user.PrimaryParty`) rather than
  hard-coding it.

---

## Wrong or unresolved template id

This is the most common integration failure with generated code.

**Symptoms.** Submission fails with a DAML error about an unknown template,
package, or `PACKAGE_NAMES_NOT_FOUND` (a real `CategoryID > 0` error).

**Cause.** Generated `GetTemplateID()` returns a **package-name** reference:

```
#package-name:Module:Entity
```

Some participants require a concrete **package-id** reference
(`<64-hex-package-id>:Module:Entity`). The package name isn't resolvable until
the DAR is uploaded and vetted, and even then a participant may demand the id.

**Fix.** After uploading the DAR, look up the real package id and rewrite the
first segment of the template id:

```go
pkgs, _ := cl.PackageMng.ListKnownPackages(ctx)
var packageID string
for _, p := range pkgs {
    if strings.EqualFold(p.Name, "all-kinds-of") {
        packageID = p.PackageID
    }
}

cmd := contract.CreateCommand()
parts := strings.Split(cmd.TemplateID, ":") // ["#all-kinds-of", "Module", "Entity"]
if len(parts) >= 3 {
    parts[0] = packageID
    cmd.TemplateID = strings.Join(parts, ":")
}
```

The same rewrite applies to `*model.ExerciseCommand.TemplateID`. Make sure the
DAR is uploaded **and** the package shows up in `ListKnownPackages` before you
submit — otherwise the id won't resolve.

---

## Command submission rejected

**Symptoms.** `SubmitAndWait` returns a real DAML error (`CategoryID > 0`), or a
`-5` `InvalidArgument`.

**Checklist.**

- **Command not wrapped.** `Commands.Commands` is `[]*model.Command`, and each
  element wraps the command type: `{Command: createCmd}`. A bare
  `*model.CreateCommand` won't compile against the field; a missing wrap is a
  frequent copy-paste slip.
- **Empty `Commands`.** Submitting with an empty `Commands` slice is accepted by
  some participants and rejected by others — make sure you actually populated it.
- **`ActAs` empty.** At least one party is required.
- **Duplicate `CommandID` within the dedup window.** Reusing a `CommandID`
  inside `DeduplicationPeriod` is deduplicated (silently succeeds without a new
  effect) or rejected. Use unique ids per logical command.
- **Bad argument shapes.** Hand-built `Arguments` maps must match the template's
  field names and DAML types. Prefer generated `CreateCommand()` /
  choice methods, which build the maps correctly via `pkg/codec`.

---

## Streaming calls hang or return nothing

Streaming methods (`GetActiveContracts`, `GetUpdates`, `CompletionStream`)
return a **response channel and an error channel**. They don't return data
directly.

**Symptoms.** The program blocks, or you "get nothing back".

**Causes & fixes.**

- **Not draining the error channel.** Always `select` over both channels plus
  `ctx.Done()`. Ignoring `errCh` hides the real failure.
- **No context cancellation.** A long-lived stream runs until the server closes
  it or the context is cancelled. Use `context.WithTimeout` /
  `context.WithCancel` and `cancel()` when you have what you need, otherwise the
  loop never exits.
- **Closed response channel = end of stream.** When `resp, ok := <-respCh` gives
  `ok == false`, the stream finished normally — return, don't treat it as an
  error.
- **Filter excludes everything.** An over-narrow `FiltersByParty` /
  `TemplateFilters` yields an empty stream that closes immediately. Widen the
  filter (empty `InclusiveFilters{}` matches all templates for the party) to
  confirm data exists.
- **Querying at the wrong offset.** `GetActiveContracts` needs a valid
  `ActiveAtOffset`; pass `GetLedgerEnd(...).Offset`. The typed `ContractQuery`
  does this for you.

---

## TLS problems

**Symptoms.** Handshake errors, `connection closed`, or auth headers not
arriving.

**Causes & fixes.**

- **Plaintext vs TLS mismatch.** If the participant requires TLS you must call
  `WithTLSConfig`; if it's plaintext, you must **not**. A TLS client against a
  plaintext port (or vice versa) fails at the handshake.
- **Token dropped.** With TLS the token rides on per-RPC credentials, which gRPC
  only sends over a secure transport. Connecting insecurely while expecting
  per-RPC credentials means the token is silently dropped — set TLS, or rely on
  the interceptor path used in insecure mode.
- **Certificate trust.** Point `client.TlsConfig{Certificate: "/path/to/ca.pem"}`
  at the CA that signed the server certificate.

---

## Topology / admin calls fail but ledger calls work

**Symptom.** `cl.TopologyManagerRead` / `TopologyManagerWrite` calls error while
`CommandService` etc. succeed.

**Cause.** Topology and traffic-control calls use the **admin** connection. If
you didn't set `WithAdminAddress`, they fall back to the ledger endpoint, which
usually doesn't serve the topology API.

**Fix.** Configure the admin endpoint:

```go
cl, err := client.NewDamlClient("localhost:6865", provider).
    WithAdminAddress("localhost:6866").
    Build(ctx)
```

Also check the `BaseQuery.Store` value — reads target a specific store
(`"authorized"`, `"synchronizer:<id>"`, `"temporary:<name>"`); querying the wrong
store returns empty results, not an error.

---

## Codegen failures

**Symptoms.** `godaml` exits non-zero.

**Checklist.**

- **Missing required flags.** `--dar`, `--output`, and `--go_package` are all
  required; omitting any aborts with exit code 1.
- **Bad DAR path.** `--dar` must point at a readable `.dar` archive.
- **Unwritable output dir.** `--output` must be creatable/writable.
- **Opaque error.** Re-run with `--debug` for verbose generation logging; it
  surfaces manifest parsing, DAML-LF version detection (v2 vs v3), and per-DALF
  generation steps.
- **Generated code won't compile.** The generated package imports
  `pkg/codec`, `pkg/model`, and dot-imports `pkg/types`. Ensure your module
  depends on `github.com/noders-team/go-daml` at a compatible version and run
  `go mod tidy`.

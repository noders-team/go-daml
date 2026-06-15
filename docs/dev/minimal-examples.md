# Minimal Examples

Focused, copy-paste snippets for each service. Every example assumes you already
have a connected client:

```go
cl, err := client.NewDamlClient("localhost:6865", auth.NewBearerTokenProvider(token)).
    WithAdminAddress("localhost:6866").
    Build(ctx)
if err != nil {
    panic(err)
}
defer cl.Close()
```

All services hang off `*client.DamlBindingClient` as public fields. Imports used
throughout:

```go
import (
    "context"
    "time"

    "github.com/noders-team/go-daml/pkg/client"
    "github.com/noders-team/go-daml/pkg/model"
)
```

- [Ledger services](#ledger-services)
  - [Submit and wait (create / exercise)](#submit-and-wait-create--exercise)
  - [Transfer Canton Coin (CC / Amulet)](#transfer-canton-coin-cc--amulet)
  - [Fire-and-forget submission](#fire-and-forget-submission)
  - [Read ledger end & active contracts](#read-ledger-end--active-contracts)
  - [Typed contract query](#typed-contract-query)
  - [Stream updates (transactions)](#stream-updates-transactions)
  - [Query events by contract id](#query-events-by-contract-id)
  - [Stream command completions](#stream-command-completions)
  - [Interactive submission (prepare / execute)](#interactive-submission-prepare--execute)
  - [Version & ledger info](#version--ledger-info)
- [Admin services](#admin-services)
  - [Users](#users)
  - [Parties](#parties)
  - [Packages (upload a DAR)](#packages-upload-a-dar)
  - [Identity provider config](#identity-provider-config)
  - [Pruning](#pruning)
- [Topology services](#topology-services)

---

## Ledger services

### Submit and wait (create / exercise)

`SubmitAndWait` blocks until the command completes and returns the update id.

```go
resp, err := cl.CommandService.SubmitAndWait(ctx, &model.SubmitAndWaitRequest{
    Commands: &model.Commands{
        WorkflowID:   "wf-1",
        UserID:       "alice",
        CommandID:    "cmd-1",
        ActAs:        []string{party},
        SubmissionID: "sub-1",
        DeduplicationPeriod: model.DeduplicationDuration{Duration: 60 * time.Second},
        Commands: []*model.Command{
            {Command: &model.CreateCommand{
                TemplateID: "<pkg>:Module:Entity",
                Arguments:  map[string]interface{}{"operator": party},
            }},
        },
    },
})
// resp.UpdateID, resp.CompletionOffset
```

Exercise a choice — swap the command type:

```go
{Command: &model.ExerciseCommand{
    TemplateID: "<pkg>:Module:Entity",
    ContractID: contractID,
    Choice:     "Transfer",
    Arguments:  map[string]interface{}{"newOwner": bob},
}}
```

With generated code you build these via `template.CreateCommand()` and
`template.ChoiceName(contractID, args)` instead of hand-writing the maps.

### Transfer Canton Coin (CC / Amulet)

Canton Coin is just an Amulet holding governed by the Splice Token Standard.
A direct CC transfer is an **`AmuletRules_Transfer`** choice exercised on the
network's `AmuletRules` contract — there is no special SDK call, it's an ordinary
`ExerciseCommand` whose arguments follow the token-standard shape. The pattern
below mirrors the production `transferDirect` in `ambo-canton`
(`internal/service/mana.go`); `go-wallet-daml`'s `TokenStandardController` is an
equivalent wrapper. Both build the same maps over the raw SDK.

**Prerequisites** you must resolve first — in practice these come from a Scan /
app-context lookup that returns both the contract ids **and their disclosed
forms**:

- `amuletRulesTemplateID` / `amuletRulesContractID` — the active `AmuletRules`
  contract (resolve its package id via `cl.PackageMng.ListKnownPackages`, then
  find the contract).
- `openMiningRoundCID` — the current open mining round contract id.
- `disclosed` — `[]*model.DisclosedContract` for `AmuletRules` **and** the open
  mining round. These are contracts your participant doesn't host, so they must
  be disclosed (each needs `TemplateID`, `ContractID`, `CreatedEventBlob`) or the
  exercise fails.
- `inputCIDs` — the sender's Amulet holding (UTXO) contract ids to spend.
- `syncID` — target synchronizer (`cl.StateService.GetConnectedSynchronizers`).

Amulet inputs are a DAML variant; amounts are `NUMERIC`, built with the
`types.NewNumericFromDecimal` helper (CC has 10 implied decimals — the helper
applies the scaling).

```go
import (
    "github.com/shopspring/decimal"
    "github.com/noders-team/go-daml/pkg/types"
)

// Amulet input variant: tag "InputAmulet", value = the holding contract id.
type inputAmulet struct{ ContractID string }

func (v inputAmulet) GetVariantTag() string        { return "InputAmulet" }
func (v inputAmulet) GetVariantValue() interface{} { return types.CONTRACT_ID(v.ContractID) }

func transferCC(ctx context.Context, cl *client.DamlBindingClient,
    sender, receiver, dso, userID, syncID string,
    amuletRulesTemplateID, amuletRulesContractID, openMiningRoundCID string,
    inputCIDs []string,
    amount decimal.Decimal,
    disclosed []*model.DisclosedContract,
) (string, error) {

    inputs := make([]interface{}, 0, len(inputCIDs))
    for _, cid := range inputCIDs {
        inputs = append(inputs, inputAmulet{ContractID: cid})
    }

    emptyGenMap := map[string]interface{}{"_type": "genmap", "value": map[string]interface{}{}}

    exercise := &model.ExerciseCommand{
        TemplateID: amuletRulesTemplateID,
        ContractID: amuletRulesContractID,
        Choice:     "AmuletRules_Transfer",
        Arguments: map[string]interface{}{
            "transfer": map[string]interface{}{
                "sender":   types.PARTY(sender),
                "provider": types.PARTY(sender),
                "inputs":   inputs,
                "outputs": []interface{}{map[string]interface{}{
                    "receiver":         types.PARTY(receiver),
                    "receiverFeeRatio": types.NewNumericFromDecimal(decimal.Zero),
                    "amount":           types.NewNumericFromDecimal(amount),
                }},
            },
            "context": map[string]interface{}{
                "openMiningRound":     types.CONTRACT_ID(openMiningRoundCID),
                "issuingMiningRounds": emptyGenMap,
                "validatorRights":     emptyGenMap,
            },
            "expectedDso": map[string]interface{}{"_type": "optional", "value": types.PARTY(dso)},
        },
    }

    resp, err := cl.CommandService.SubmitAndWait(ctx, &model.SubmitAndWaitRequest{
        Commands: &model.Commands{
            UserID:             userID,
            CommandID:          "cc-transfer-1",
            ActAs:              []string{sender},
            SubmissionID:       "cc-transfer-sub-1",
            SynchronizerID:     syncID,
            DisclosedContracts: disclosed,
            DeduplicationPeriod: model.DeduplicationDuration{Duration: 60 * time.Second},
            Commands:           []*model.Command{{Command: exercise}},
        },
    })
    if err != nil {
        return "", err
    }
    return resp.UpdateID, nil
}
```

Notes:

- **Disclosed contracts are not optional in practice.** `AmuletRules` and the
  open mining round live on other participants; omit them and the transfer is
  rejected. Get them (with their `CreatedEventBlob`) from a Scan lookup.
- **Inputs must cover the amount.** Select holdings whose summed balance ≥ amount;
  the choice returns change as a new Amulet to the sender.
- **`receiverFeeRatio`** splits transfer fees; `0` means the sender pays. The
  `transfer` record and `context` also accept optional `beneficiaries`,
  per-output `lock`, and `context.featuredAppRight` fields, omitted here.
- **One-step preapproved transfers.** If the receiver has a `TransferPreapproval`
  contract, exercise **`TransferPreapproval_Send`** on it instead (`context`
  wrapping `amuletRules` + the same transfer `context`, plus `inputs`, `amount`,
  `sender`, optional `description`) — no receiver acceptance step. Without a
  preapproval, the token-standard **`TransferFactory_Transfer`** produces a
  `TransferInstruction` the receiver must `TransferInstruction_Accept`.
- For production (UTXO selection, mining-round / `AmuletRules` resolution,
  disclosed-contract assembly), use a wrapper like `ambo-canton`'s app-context
  service or `go-wallet-daml`'s `TokenStandardController` rather than building
  these maps by hand.

### Fire-and-forget submission

`CommandSubmission.Submit` returns as soon as the command is accepted for
processing; it does **not** wait for completion. Pair it with the completion
stream if you need the outcome.

```go
_, err := cl.CommandSubmission.Submit(ctx, &model.SubmitRequest{
    Commands: &model.Commands{
        UserID:       "alice",
        CommandID:    "cmd-async-1",
        ActAs:        []string{party},
        SubmissionID: "sub-async-1",
        DeduplicationPeriod: model.DeduplicationDuration{Duration: 60 * time.Second},
        Commands:     []*model.Command{ /* ... */ },
    },
})
```

### Read ledger end & active contracts

```go
end, err := cl.StateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
// end.Offset

req := &model.GetActiveContractsRequest{
    ActiveAtOffset: end.Offset,
    EventFormat: &model.EventFormat{
        Verbose: true,
        FiltersByParty: map[string]*model.Filters{
            party: {Inclusive: &model.InclusiveFilters{
                TemplateFilters: []*model.TemplateFilter{ /* empty = all templates */ },
            }},
        },
    },
}

respCh, errCh := cl.StateService.GetActiveContracts(ctx, req)
for {
    select {
    case resp, ok := <-respCh:
        if !ok {
            return // stream finished
        }
        if e, ok := resp.ContractEntry.(*model.ActiveContractEntry); ok && e.ActiveContract != nil {
            ce := e.ActiveContract.CreatedEvent
            // ce.ContractID, ce.TemplateID, ce.CreateArguments
        }
    case err := <-errCh:
        if err != nil {
            return
        }
    case <-ctx.Done():
        return
    }
}
```

`GetActiveContracts` is a **stream**: it returns a response channel and an error
channel. Always drain both and respect `ctx.Done()`.

### Typed contract query

`ContractQuery[T]` wraps the active-contracts stream and decodes each contract's
create arguments into your generated struct `T`. This is the easiest way to read
contracts.

```go
query := client.NewContractQuery[MappyContract](cl)
contracts, err := query.FindContractsByTemplate(ctx, party, templateID)
for _, c := range contracts {
    // c.ContractID string
    // c.Data       MappyContract
}
```

It snapshots the ledger end internally, so you get a consistent active set at the
time of the call.

### Stream updates (transactions)

```go
req := &model.GetUpdatesRequest{
    BeginExclusive: 0, // tail from the start; use a saved offset to resume
    Filter: &model.TransactionFilter{
        FiltersByParty: map[string]*model.Filters{
            party: {Inclusive: &model.InclusiveFilters{}},
        },
    },
    Verbose: false,
}

respCh, errCh := cl.UpdateService.GetUpdates(ctx, req)
for {
    select {
    case resp, ok := <-respCh:
        if !ok {
            return
        }
        switch {
        case resp.Update.Transaction != nil:
            // resp.Update.Transaction.UpdateID / .Offset / .Events
        case resp.Update.Reassignment != nil:
        case resp.Update.OffsetCheckpoint != nil:
            // checkpoint offset — persist it to resume later
        }
    case err := <-errCh:
        if err != nil {
            return
        }
    case <-ctx.Done():
        return
    }
}
```

Point lookups are available too:
`GetUpdateById`, `GetTransactionByID`, `GetTransactionByOffset`.

### Query events by contract id

Fetch the create (and, if archived, archive) event for one contract:

```go
events, err := cl.EventQuery.GetEventsByContractID(ctx, &model.GetEventsByContractIDRequest{
    ContractID: contractID,
    EventFormat: &model.EventFormat{
        FiltersByParty: map[string]*model.Filters{party: {}},
    },
})
// events.CreateEvent  (*model.CreatedEvent, may be nil)
// events.ArchiveEvent (*model.ArchivedEvent, may be nil)
```

### Stream command completions

Use this with `CommandSubmission.Submit` to learn the fate of async commands.

```go
respCh, errCh := cl.CommandCompletion.CompletionStream(ctx, &model.CompletionStreamRequest{
    UserID:         "alice",
    Parties:        []string{party},
    BeginExclusive: 0,
})
for {
    select {
    case resp, ok := <-respCh:
        if !ok {
            return
        }
        // resp.Response is either a Completion (CommandID, Status, UpdateID, Offset)
        // or an OffsetCheckpoint (Offset)
    case err := <-errCh:
        if err != nil {
            return
        }
    case <-ctx.Done():
        return
    }
}
```

### Interactive submission (prepare / execute)

For externally-signed (multi-party) flows: prepare a transaction, sign the
returned hash off-ledger, then execute.

```go
prep, err := cl.InteractiveSubmissionService.PrepareSubmission(ctx, &model.PrepareSubmissionRequest{
    UserID:         "alice",
    CommandID:      "icmd-1",
    ActAs:          []string{party},
    SynchronizerID: syncID, // from StateService.GetConnectedSynchronizers
    Commands:       []*model.Command{ /* ... */ },
})
// prep.PreparedTransaction []byte, prep.PreparedTransactionHash []byte,
// prep.HashingSchemeVersion, prep.CostEstimation

// ... sign prep.PreparedTransactionHash, then:

_, err = cl.InteractiveSubmissionService.ExecuteSubmission(ctx, &model.ExecuteSubmissionRequest{
    PreparedTransaction:  prep.PreparedTransaction,
    UserID:               "alice",
    SubmissionID:         "iexec-1",
    HashingSchemeVersion: prep.HashingSchemeVersion,
    DeduplicationPeriod:  model.DeduplicationDuration{Duration: 60 * time.Second},
    PartySignatures:      []*model.SinglePartySignatures{ /* your signatures */ },
})
```

`GetPreferredPackageVersion` resolves which package version a set of parties will
accept.

### Version & ledger info

```go
v, err := cl.VersionService.GetLedgerAPIVersion(ctx, &model.GetLedgerAPIVersionRequest{})
// v.Version, v.Features

// convenience wrappers on the client:
err = cl.Ping(ctx)                          // round-trips the version service
err = cl.ValidateSDKVersion(ctx, sdkVersion) // compares against the ledger version
```

---

## Admin services

Admin calls require a user with the appropriate participant rights.

### Users

```go
users, err := cl.UserMng.ListUsers(ctx)
user, err  := cl.UserMng.GetUser(ctx, "alice")

granted, err := cl.UserMng.GrantUserRights(ctx, "alice", "", []*model.Right{
    {Type: model.CanReadAs{Party: party}},
})
```

`model.Right.Type` is one of `CanActAs`, `CanReadAs`, `ParticipantAdmin`,
`IdentityProviderAdmin`.

### Parties

```go
participantID, err := cl.PartyMng.GetParticipantID(ctx)

list, err := cl.PartyMng.ListKnownParties(ctx, "" /*pageToken*/, 100 /*pageSize*/, "" /*idpID*/)
// list.PartyDetails, list.NextPageToken

newParty, err := cl.PartyMng.AllocateParty(ctx, "alice", map[string]string{}, "")
```

### Packages (upload a DAR)

```go
dar, _ := os.ReadFile("./contracts.dar")

// validate first (no state change), then upload
if err := cl.PackageMng.ValidateDarFile(ctx, dar, "validate-1"); err != nil {
    panic(err)
}
if err := cl.PackageMng.UploadDarFile(ctx, dar, "upload-1"); err != nil {
    panic(err)
}

pkgs, err := cl.PackageMng.ListKnownPackages(ctx)
for _, p := range pkgs {
    // p.PackageID, p.Name, p.Version
}
```

The `PackageID` you find here is what you use to pin a template ID to a concrete
package — see
[troubleshooting → template/package IDs](./troubleshooting.md#wrong-or-unresolved-template-id).

### Identity provider config

```go
cfg, err := cl.IdentityProviderMng.CreateIdentityProviderConfig(ctx, &model.IdentityProviderConfig{
    IdentityProviderID: "kc",
    Issuer:             "https://keycloak.example.com/realms/daml",
    JwksURL:            "https://keycloak.example.com/realms/daml/protocol/openid-connect/certs",
    Audience:           "https://daml.example.com",
})
```

### Pruning

```go
err := cl.PruningMng.Prune(ctx, &model.PruneRequest{
    PruneUpTo:                 time.Now().Add(-24 * time.Hour).UnixMicro(),
    SubmissionID:              "prune-1",
    PruneAllDivulgedContracts: false,
})
```

---

## Topology services

Topology read/write use the **admin** connection. Configure `WithAdminAddress`
unless your participant serves both on the same endpoint.

```go
// read: list party -> participant mappings in the authorized store
resp, err := cl.TopologyManagerRead.ListPartyToParticipant(ctx, &model.ListPartyToParticipantRequest{
    BaseQuery:   &model.BaseQuery{Store: &model.StoreID{Value: "authorized"}},
    FilterParty: party,
})
for _, r := range resp.Results {
    // r holds the PartyToParticipantMapping
}
```

Other read calls: `ListNamespaceDelegation`, `ListPartyToKeyMapping`.

Write side (`cl.TopologyManagerWrite`) exposes `Authorize`, `AddTransactions`,
`SignTransactions`, `GenerateTransactions`, and temporary-store management
(`CreateTemporaryTopologyStore` / `DropTemporaryTopologyStore`). These are
advanced operations for onboarding parties and managing namespace delegations;
the `BaseQuery.Store` value selects the store (`"authorized"`,
`"synchronizer:<id>"`, or `"temporary:<name>"`).

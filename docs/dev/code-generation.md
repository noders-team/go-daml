# Code Generation: from `.dar` to type-safe Go

The `godaml` CLI parses a DAML archive (`.dar`) and emits type-safe Go structs
for every template, record, variant and enum it contains. Using generated code
instead of hand-built `map[string]interface{}` arguments is the difference
between compile-time safety and runtime "unknown field" surprises.

- [The CLI](#the-cli)
- [What gets generated](#what-gets-generated)
  - [Package header (MainDalf only)](#package-header-maindalf-only)
  - [Templates](#templates)
  - [Choices](#choices)
  - [Records](#records)
  - [Variants](#variants)
  - [Enums](#enums)
- [How generation works (two-pass)](#how-generation-works-two-pass)
- [DAML-LF version support](#daml-lf-version-support)
- [Naming rules](#naming-rules)
- [Using generated code with the client](#using-generated-code-with-the-client)
- [The template-id caveat](#the-template-id-caveat)

## The CLI

```bash
godaml --dar ./contracts.dar --output ./generated --go_package contracts
```

| Flag | Required | Meaning |
| --- | --- | --- |
| `--dar` | yes | Path to the `.dar` archive. |
| `--output` | yes | Directory to write generated `.go` files into. |
| `--go_package` | yes | Package name placed at the top of every generated file. |
| `--debug` | no | Verbose generation logging (manifest parsing, version detection, per-DALF steps). |

The generator emits **one Go file per DALF** inside the archive. Standard
library / prim DALFs are skipped. Build the CLI with `make build` (→
`./bin/godaml`) or run it as a Go tool — see
[getting-started → build the godaml CLI](./getting-started.md#build-the-godaml-cli).

## What gets generated

Every generated file imports `pkg/codec`, `pkg/model`, and dot-imports
`pkg/types`, so DAML primitives (`PARTY`, `TEXT`, `INT64`, `NUMERIC`, `TEXTMAP`,
…) read naturally.

### Package header (MainDalf only)

The file generated from the archive's main DALF carries package metadata:

```go
const SDKVersion = "3.3.0-snapshot.20250417.0"
const packageName = "all-kinds-of"
const version = "1.0.0"

func GetPackageName() string { return packageName }
func GetVersion() string     { return version }

type Template interface {
    CreateCommand() *model.CreateCommand
    GetTemplateID() string
}
```

### Templates

Each DAML template becomes a struct implementing `Template`:

```go
type MappyContract struct {
    Operator PARTY   `json:"operator"`
    Value    TEXTMAP `json:"value"`
}

func (t MappyContract) GetTemplateID() string {
    return fmt.Sprintf("#%s:%s:%s", packageName, "AllKindsOf", "MappyContract")
}

func (t MappyContract) CreateCommand() *model.CreateCommand {
    args := make(map[string]interface{})
    args["operator"] = t.Operator
    args["value"] = t.Value
    return &model.CreateCommand{
        TemplateID: t.GetTemplateID(),
        Arguments:  args,
    }
}
```

### Choices

Every template gets an `Archive` method, and one method per custom choice. The
method returns a ready-to-submit `*model.ExerciseCommand`.

```go
// Archive — generated for every template
func (t MappyContract) Archive(contractID string) *model.ExerciseCommand {
    return &model.ExerciseCommand{
        TemplateID: t.GetTemplateID(),
        ContractID: contractID,
        Choice:     "Archive",
        Arguments:  map[string]interface{}{},
    }
}
```

A custom choice that takes arguments generates a method whose second parameter is
the generated record type for that choice:

```go
func (t T) MyChoice(contractID string, args MyChoice) *model.ExerciseCommand
```

### Records

DAML records (including choice argument types) become structs with a `ToMap()`
helper and custom JSON marshaling backed by `pkg/codec`:

```go
type OptionalFieldsCleanUp struct{}

func (t OptionalFieldsCleanUp) ToMap() map[string]interface{} { /* ... */ }

func (t OptionalFieldsCleanUp) MarshalJSON() ([]byte, error) {
    return codec.NewJsonCodec().Marshall(t)
}
func (t *OptionalFieldsCleanUp) UnmarshalJSON(data []byte) error {
    return codec.NewJsonCodec().Unmarshall(data, t)
}
```

### Variants

DAML variants (sum types) become structs with one optional pointer field per
constructor, plus `GetVariantTag` / `GetVariantValue` (satisfying
`types.VARIANT`):

```go
type Address struct {
    US *USAddress `json:"US,omitempty"`
    UK *UKAddress `json:"UK,omitempty"`
}

func (v Address) GetVariantTag() string {
    if v.US != nil { return "US" }
    if v.UK != nil { return "UK" }
    return ""
}

func (v Address) GetVariantValue() interface{} {
    if v.US != nil { return v.US }
    if v.UK != nil { return v.UK }
    return nil
}

var _ VARIANT = (*Address)(nil)
```

Set exactly one field; the codec serializes it to DAML's `{tag, value}` shape.

### Enums

DAML enums become a string type with typed constants and `GetEnumConstructor` /
`GetEnumTypeID` (satisfying `types.ENUM`):

```go
type Color string

const (
    ColorRed   Color = "Red"
    ColorGreen Color = "Green"
    ColorBlue  Color = "Blue"
)

func (e Color) GetEnumConstructor() string { return string(e) }
func (e Color) GetEnumTypeID() string {
    return fmt.Sprintf("#%s:%s:%s", packageName, "AllKindsOf", "Color")
}

var _ ENUM = Color("")
```

## How generation works (two-pass)

Cross-DALF interface references mean the generator can't process files in
isolation. It runs two passes:

1. **Pass 1 — interface collection.** Extract the MainDalf, parse every
   dependency DALF, and collect all interface definitions into a global map so
   they're resolvable everywhere.
2. **Pass 2 — generation.** Generate one Go file per DALF using that map, so an
   interface defined in one DALF and implemented in another resolves correctly.
   stdlib/prim DALFs are skipped.

Each file is rendered through `internal/codegen/source.go.tpl` and then run
through `gofmt`.

## DAML-LF version support

The DAML-LF version is auto-detected from the archive manifest and the matching
parser is selected:

- **v2** — DAML SDK 1.x (`daml_lf_1_17` proto), `astgen/v2`.
- **v3** — DAML SDK 2.x+ (`daml_lf_2_1` proto), `astgen/v3`.

You don't choose this; `--debug` shows which version was detected.

## Naming rules

- **Interfaces** are prefixed with `I` to avoid clashing with same-named structs:
  `Transferable` → `ITransferable`.
- **Case is preserved.** PascalCase and acronyms survive
  (`FeaturedAppRight_CreateActivityMarker` → `FeaturedAppRightCreateActivityMarker`;
  `US`/`UK` stay upper-case). Generated identifiers are case-sensitive — match
  them exactly.
- **Underscores are stripped** during primitive/type normalization
  (`Featured_App_Right` → `FeaturedAppRight`).

## Using generated code with the client

```go
contract := MappyContract{
    Operator: PARTY(party),
    Value:    TEXTMAP{"k1": "v1"},
}

// create
createCmd := contract.CreateCommand()

// exercise
archiveCmd := contract.Archive(contractID)

// read back, decoded into the generated struct
query := client.NewContractQuery[MappyContract](cl)
contracts, err := query.FindContractsByTemplate(ctx, party, contract.GetTemplateID())
```

See [getting-started → your first contract](./getting-started.md#your-first-contract-create--query--exercise)
for the full submit flow, and [minimal-examples.md](./minimal-examples.md) for
per-service snippets.

## The template-id caveat

`GetTemplateID()` returns a **package-name** reference
(`#package-name:Module:Entity`). Some participants require a concrete
**package-id** reference (`<64-hex>:Module:Entity`), which only exists after the
DAR is uploaded and vetted. If submission fails with an unknown-template /
`PACKAGE_NAMES_NOT_FOUND` error, rewrite the first segment with the real package
id from `cl.PackageMng.ListKnownPackages`. Full recipe:
[troubleshooting → wrong or unresolved template id](./troubleshooting.md#wrong-or-unresolved-template-id).

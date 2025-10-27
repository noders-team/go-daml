package interfaces

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const (
	PackageID  = "8f919735d2daa1abb780808ad1fed686fc9229a039dc659ccb04e5fd5d071c90"
	SDKVersion = "3.3.0-snapshot.20250507.0"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// ITransferable is a DAML interface
type ITransferable interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// Transfer executes the Transfer choice
	Transfer(contractID string, args Transfer) *model.ExerciseCommand
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	// Check if the type has a toMap method
	type mapper interface {
		toMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.toMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}

// Asset is a Template type
type Asset struct {
	Owner PARTY `json:"owner"`
	Name  TEXT  `json:"name"`
	Value INT64 `json:"value"`
}

// GetTemplateID returns the template ID for this template
func (t Asset) GetTemplateID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "Interfaces", "Asset")
}

// CreateCommand returns a CreateCommand for this template
func (t Asset) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	args["name"] = string(t.Name)

	args["value"] = int64(t.Value)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for Asset using JsonCodec
func (t Asset) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Asset using JsonCodec
func (t *Asset) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for Asset

// Archive exercises the Archive choice on this Asset contract
func (t Asset) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Interfaces", "Asset"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// AssetTransfer exercises the AssetTransfer choice on this Asset contract
func (t Asset) AssetTransfer(contractID string, args AssetTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Interfaces", "Asset"),
		ContractID: contractID,
		Choice:     "AssetTransfer",
		Arguments:  argsToMap(args),
	}
}

// Transfer exercises the Transfer choice on this Asset contract via the ITransferable interface
func (t Asset) Transfer(contractID string, args Transfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Interfaces", "Transferable"),
		ContractID: contractID,
		Choice:     "Transfer",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Asset

var _ ITransferable = (*Asset)(nil)

// AssetTransfer is a Record type
type AssetTransfer struct {
	NewOwner PARTY `json:"newOwner"`
}

// toMap converts AssetTransfer to a map for DAML arguments
func (t AssetTransfer) toMap() map[string]interface{} {
	return map[string]interface{}{
		"newOwner": t.NewOwner.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for AssetTransfer using JsonCodec
func (t AssetTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AssetTransfer using JsonCodec
func (t *AssetTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Token is a Template type
type Token struct {
	Issuer PARTY   `json:"issuer"`
	Owner  PARTY   `json:"owner"`
	Amount NUMERIC `json:"amount"`
}

// GetTemplateID returns the template ID for this template
func (t Token) GetTemplateID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "Interfaces", "Token")
}

// CreateCommand returns a CreateCommand for this template
func (t Token) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["issuer"] = t.Issuer.ToMap()

	args["owner"] = t.Owner.ToMap()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for Token using JsonCodec
func (t Token) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Token using JsonCodec
func (t *Token) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for Token

// Archive exercises the Archive choice on this Token contract
func (t Token) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Interfaces", "Token"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// Transfer exercises the Transfer choice on this Token contract via the ITransferable interface
func (t Token) Transfer(contractID string, args Transfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Interfaces", "Transferable"),
		ContractID: contractID,
		Choice:     "Transfer",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Token

var _ ITransferable = (*Token)(nil)

// Transfer is a Record type
type Transfer struct {
	NewOwner PARTY `json:"newOwner"`
}

// toMap converts Transfer to a map for DAML arguments
func (t Transfer) toMap() map[string]interface{} {
	return map[string]interface{}{
		"newOwner": t.NewOwner.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for Transfer using JsonCodec
func (t Transfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Transfer using JsonCodec
func (t *Transfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferableView is a Record type
type TransferableView struct {
	Owner PARTY `json:"owner"`
}

// toMap converts TransferableView to a map for DAML arguments
func (t TransferableView) toMap() map[string]interface{} {
	return map[string]interface{}{
		"owner": t.Owner.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TransferableView using JsonCodec
func (t TransferableView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferableView using JsonCodec
func (t *TransferableView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ITransferableInterfaceID returns the interface ID for the ITransferable interface
func ITransferableInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("%s:%s:%s", pkgID, "Interfaces", "Transferable")
}

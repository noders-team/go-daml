package model

import (
	"github.com/shopspring/decimal"
)

type TransferInstruction struct {
	ContractID       string
	Sender           string
	Receiver         string
	Amount           decimal.Decimal
	Memo             string
	CreatedEventBlob []byte
}

type AllocationInstruction struct {
	ContractID       string
	Provider         string
	Specification    map[string]interface{}
	CreatedEventBlob []byte
}

type AllocationRequest struct {
	ContractID       string
	Requester        string
	Specification    map[string]interface{}
	CreatedEventBlob []byte
}

type CommandRequest struct {
	Command            *Command
	DisclosedContracts []*DisclosedContract
}

type HoldingUTXO struct {
	ContractID       string
	Amount           decimal.Decimal
	InstrumentID     string
	InstrumentAdmin  string
	Owner            string
	Lock             map[string]interface{}
	CreatedEventBlob []byte
}

type FeaturedAppRight struct {
	TemplateID       string
	ContractID       string
	Payload          map[string]interface{}
	CreatedEventBlob []byte
	CreatedAt        string
}

type MergeUtxosResult struct {
	Commands           []*Command
	DisclosedContracts []*DisclosedContract
}

type Allocation struct {
	ContractID       string
	Provider         string
	Receiver         string
	Amount           decimal.Decimal
	CreatedEventBlob []byte
}

// ChoiceContext mirrors the registry choice-context payload consumed by the
// *FromContext builders. It is normally produced by a token-standard registry
// HTTP call and passed in by the caller.
type ChoiceContext struct {
	ChoiceContextData  map[string]interface{}
	DisclosedContracts []*DisclosedContract
}

// Beneficiary is a featured-app reward split used by the delegate-proxy choices.
type Beneficiary struct {
	Beneficiary string
	Weight      float64
}

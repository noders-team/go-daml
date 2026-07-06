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

type CreateTransferResult struct {
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

package client

import (
	"context"
	"fmt"
	"time"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/service/ledger"
)

type ContractQuery[T any] struct {
	cl *DamlBindingClient
}

func NewContractQuery[T any](client *DamlBindingClient) *ContractQuery[T] {
	return &ContractQuery[T]{
		cl: client,
	}
}

type Contract[T any] struct {
	ContractID string
	Data       T
}

func (c *ContractQuery[T]) FindContractsByTemplate(ctx context.Context, partyID, templateID string) ([]Contract[T], error) {
	var results []Contract[T]
	err := c.scanActiveContractsByTemplate(ctx, contractQuery{
		partyID:    partyID,
		templateID: templateID,
	}, func(evt activeContractEvent) (bool, error) {
		var t T
		if err := ledger.RecordToStruct(evt.arguments, &t); err != nil {
			return false, fmt.Errorf("decode contract %s: %w", evt.contractID, err)
		}
		results = append(results, Contract[T]{ContractID: evt.contractID, Data: t})
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

type contractQuery struct {
	partyID    string
	templateID string
	anyParty   bool
}

type activeContractEvent struct {
	contractID string
	templateID string
	arguments  any
	createdAt  *time.Time
}

func (c *ContractQuery[T]) scanActiveContractsByTemplate(
	ctx context.Context,
	query contractQuery,
	onEvent func(evt activeContractEvent) (stop bool, err error),
) error {
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	req, err := c.newActiveContractsRequest(streamCtx, query)
	if err != nil {
		return err
	}

	respCh, errCh := c.cl.StateService.GetActiveContracts(streamCtx, req)
	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				return nil
			}
			entry, ok := resp.ContractEntry.(*model.ActiveContractEntry)
			if !ok || entry.ActiveContract == nil || entry.ActiveContract.CreatedEvent == nil {
				continue
			}
			evt := entry.ActiveContract.CreatedEvent
			stop, err := onEvent(activeContractEvent{
				contractID: evt.ContractID,
				templateID: evt.TemplateID,
				arguments:  evt.CreateArguments,
				createdAt:  evt.CreatedAt,
			})
			if err != nil {
				return err
			}
			if stop {
				return nil
			}
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("error scanning active contracts: %w", err)
			}
		case <-streamCtx.Done():
			return streamCtx.Err()
		}
	}
}

func (c *ContractQuery[T]) newActiveContractsRequest(ctx context.Context, query contractQuery) (*model.GetActiveContractsRequest, error) {
	ledgerEnd, err := c.cl.StateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger end: %w", err)
	}

	eventFormat := &model.EventFormat{Verbose: true}
	filter := &model.Filters{
		Inclusive: &model.InclusiveFilters{
			TemplateFilters: []*model.TemplateFilter{{TemplateID: query.templateID}},
		},
	}
	if query.anyParty {
		eventFormat.FiltersForAnyParty = filter
	} else {
		eventFormat.FiltersByParty = map[string]*model.Filters{query.partyID: filter}
	}

	return &model.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.Offset,
		EventFormat:    eventFormat,
	}, nil
}

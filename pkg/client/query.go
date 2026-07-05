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
	ContractID       string
	TemplateID       string
	CreatedAt        *time.Time
	CreatedEventBlob []byte
	Data             T
}

func (c *ContractQuery[T]) FindContractsByTemplate(ctx context.Context, partyID, templateID string) ([]Contract[T], error) {
	return c.collect(ctx, contractQuery{
		partyID:    partyID,
		templateID: templateID,
	})
}

func (c *ContractQuery[T]) FindContractsByTemplateAnyParty(ctx context.Context, templateID string) ([]Contract[T], error) {
	return c.collect(ctx, contractQuery{
		templateID: templateID,
		anyParty:   true,
	})
}

func (c *ContractQuery[T]) FindContractsByInterface(ctx context.Context, partyID, interfaceID string) ([]Contract[T], error) {
	return c.collect(ctx, contractQuery{
		partyID:     partyID,
		interfaceID: interfaceID,
	})
}

func (c *ContractQuery[T]) FindContractsByInterfaceAnyParty(ctx context.Context, interfaceID string) ([]Contract[T], error) {
	return c.collect(ctx, contractQuery{
		interfaceID: interfaceID,
		anyParty:    true,
	})
}

func (c *ContractQuery[T]) collect(ctx context.Context, query contractQuery) ([]Contract[T], error) {
	var results []Contract[T]
	err := c.scanActiveContractsByTemplate(ctx, query, func(evt activeContractEvent) (bool, error) {
		var t T
		if err := ledger.RecordToStruct(evt.arguments, &t); err != nil {
			return false, fmt.Errorf("decode contract %s: %w", evt.contractID, err)
		}
		results = append(results, Contract[T]{
			ContractID:       evt.contractID,
			TemplateID:       evt.templateID,
			CreatedAt:        evt.createdAt,
			CreatedEventBlob: evt.createdEventBlob,
			Data:             t,
		})
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

type contractQuery struct {
	partyID     string
	templateID  string
	interfaceID string
	anyParty    bool
}

type activeContractEvent struct {
	contractID       string
	templateID       string
	arguments        any
	createdAt        *time.Time
	createdEventBlob []byte
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
			arguments := evt.CreateArguments
			if query.interfaceID != "" {
				arguments = nil
				for _, iv := range evt.InterfaceViews {
					if iv.ViewValue != nil {
						arguments = iv.ViewValue
						break
					}
				}
				if arguments == nil {
					continue
				}
			}
			stop, err := onEvent(activeContractEvent{
				contractID:       evt.ContractID,
				templateID:       evt.TemplateID,
				arguments:        arguments,
				createdAt:        evt.CreatedAt,
				createdEventBlob: evt.CreatedEventBlob,
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
	inclusive := &model.InclusiveFilters{}
	if query.interfaceID != "" {
		inclusive.InterfaceFilters = []*model.InterfaceFilter{{
			InterfaceID:             query.interfaceID,
			IncludeInterfaceView:    true,
			IncludeCreatedEventBlob: true,
		}}
	} else {
		inclusive.TemplateFilters = []*model.TemplateFilter{{
			TemplateID:              query.templateID,
			IncludeCreatedEventBlob: true,
		}}
	}
	filter := &model.Filters{Inclusive: inclusive}
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

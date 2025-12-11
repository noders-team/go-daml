package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunStateService(cl *client.DamlBindingClient) {
	ledgerEndReq := &model.GetLedgerEndRequest{}
	ledgerEnd, err := cl.StateService.GetLedgerEnd(context.Background(), ledgerEndReq)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get ledger end")
	}

	log.Info().
		Int64("offset", ledgerEnd.Offset).
		Msg("ledger end offset")

	prunedOffsetsReq := &model.GetLatestPrunedOffsetsRequest{}
	prunedOffsets, err := cl.StateService.GetLatestPrunedOffsets(context.Background(), prunedOffsetsReq)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get pruned offsets")
	} else {
		log.Info().
			Int64("participantPrunedUpToInclusive", prunedOffsets.ParticipantPrunedUpToInclusive).
			Int64("allDivulgedContractsPrunedUpToInclusive", prunedOffsets.AllDivulgedContractsPrunedUpToInclusive).
			Msg("latest pruned offsets")
	}

	synchronizersReq := &model.GetConnectedSynchronizersRequest{}
	synchronizers, err := cl.StateService.GetConnectedSynchronizers(context.Background(), synchronizersReq)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get connected synchronizers")
	} else {
		log.Info().
			Interface("synchronizers", synchronizers.ConnectedSynchronizers).
			Msg("connected synchronizers")
	}

	party := getAvailableParty(cl)
	activeContractsReq := &model.GetActiveContractsRequest{
		EventFormat: &model.EventFormat{
			FiltersByParty: map[string]*model.Filters{
				party: {
					Inclusive: &model.InclusiveFilters{
						TemplateFilters: []*model.TemplateFilter{},
					},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	responseCh, errCh := cl.StateService.GetActiveContracts(ctx, activeContractsReq)

	contractCount := 0
	for {
		select {
		case response, ok := <-responseCh:
			if !ok {
				log.Info().Int("totalContracts", contractCount).Msg("active contracts stream completed")
				return
			}
			if response != nil && response.ContractEntry != nil {
				switch entry := response.ContractEntry.(type) {
				case *model.ActiveContractEntry:
					if entry.ActiveContract != nil {
						contractCount++
						log.Info().
							Str("contractID", entry.ActiveContract.CreatedEvent.ContractID).
							Str("templateID", entry.ActiveContract.CreatedEvent.TemplateID).
							Str("synchronizerID", entry.ActiveContract.SynchronizerID).
							Uint64("reassignmentCounter", entry.ActiveContract.ReassignmentCounter).
							Msg("received active contract")
					}
				case *model.IncompleteUnassignedEntry:
					log.Info().Msg("received incomplete unassigned contract")
				case *model.IncompleteAssignedEntry:
					log.Info().Msg("received incomplete assigned contract")
				}
			}
		case err := <-errCh:
			if err != nil {
				log.Warn().Err(err).Msg("active contracts stream error")
				return
			}
		case <-ctx.Done():
			log.Info().Int("totalContracts", contractCount).Msg("active contracts stream timeout")
			return
		}
	}
}

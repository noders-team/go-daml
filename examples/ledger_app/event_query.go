package main

import (
	"context"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunEventQuery(cl *client.DamlBindingClient) {
	contractID := "example-contract-id-123"
	party := getAvailableParty(cl)

	eventReq := &model.GetEventsByContractIDRequest{
		ContractID: contractID,
		EventFormat: &model.EventFormat{
			FiltersByParty: map[string]*model.Filters{
				party: {},
			},
		},
	}

	events, err := cl.EventQuery.GetEventsByContractID(context.Background(), eventReq)
	if err != nil {
		log.Warn().Err(err).
			Str("contractId", contractID).
			Msg("failed to get events by contract ID (expected for non-existent contract)")
	} else {
		eventCount := 0
		if events.CreateEvent != nil {
			eventCount++
		}
		if events.ArchiveEvent != nil {
			eventCount++
		}

		log.Info().
			Str("contractId", contractID).
			Int("eventCount", eventCount).
			Interface("createEvent", events.CreateEvent).
			Interface("archiveEvent", events.ArchiveEvent).
			Msg("got events by contract ID")
	}
}

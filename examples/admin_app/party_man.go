package main

import (
	"context"
	"os"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunPartyManagement(cl *client.DamlBindingClient) {
	participantID, err := cl.PartyMng.GetParticipantID(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get participant ID")
	}
	log.Info().Str("participantID", participantID).Msg("participant ID")

	response, err := cl.PartyMng.ListKnownParties(context.Background(), "", 10, "")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list known parties")
	}
	for _, d := range response.PartyDetails {
		log.Info().Interface("party", d).Msg("received party details")
	}

	allocDetails, err := cl.PartyMng.GetParties(context.Background(), []string{"participant_admin"}, "")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get parties")
	}
	for _, ad := range allocDetails {
		log.Info().Interface("party", ad).Msg("party details")
	}

	newPartyHint := "test_party_" + os.Getenv("USER") + "_" + os.Getenv("HOSTNAME")
	newParty, err := cl.PartyMng.AllocateParty(context.Background(), newPartyHint, map[string]string{
		"description": "New test party",
		"type":        "testing",
	}, "")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to allocate party")
	}
	log.Info().Interface("party", newParty).Msg("new party allocated")

	updatedParty, err := cl.PartyMng.UpdatePartyDetails(context.Background(), &model.PartyDetails{
		Party:   newParty.Party,
		IsLocal: true,
		LocalMetadata: map[string]string{
			"updated":     "true",
			"version":     "v2",
			"description": "Updated test party",
		},
	}, &model.UpdateMask{
		Paths: []string{"local_metadata"},
	})
	if err != nil {
		log.Error().Err(err).Msg("update party details error")
	} else {
		log.Info().Interface("party", updatedParty).Msg("updated party")
	}

	err = cl.PartyMng.UpdatePartyIdentityProviderID(context.Background(), "participant_admin", "", "new-provider-id")
	if err != nil {
		log.Warn().Err(err).Msg("update identity provider error (expected if not supported)")
	}

	updatedResponse, err := cl.PartyMng.ListKnownParties(context.Background(), "", 10, "")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list known parties after update")
	}
	for _, d := range updatedResponse.PartyDetails {
		log.Info().Interface("party", d).Msg("updated party details")
	}
}

package main

import (
	"context"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunVersionService(cl *client.DamlBindingClient) {
	req := &model.GetLedgerAPIVersionRequest{}

	version, err := cl.VersionService.GetLedgerAPIVersion(context.Background(), req)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get ledger API version")
	}

	log.Info().
		Str("version", version.Version).
		Interface("features", version.Features).
		Msg("ledger API version")
}

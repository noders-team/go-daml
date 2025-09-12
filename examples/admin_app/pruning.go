package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunPrunning(cl *client.DamlBindingClient) {
	pruneUpTo := time.Now().Add(-24 * time.Hour).UnixMicro()

	pruneReq := &model.PruneRequest{
		PruneUpTo:                 pruneUpTo,
		SubmissionID:              "prune-" + time.Now().Format("20060102150405"),
		PruneAllDivulgedContracts: false,
	}

	log.Info().
		Time("pruneUpTo", time.UnixMicro(pruneUpTo)).
		Int64("offset", pruneUpTo).
		Msg("attempting to prune ledger")

	err := cl.PruningMng.Prune(context.Background(), pruneReq)
	if err != nil {
		log.Warn().Err(err).Msg("prune operation result")
	} else {
		log.Info().Msg("prune operation completed successfully")
	}
}

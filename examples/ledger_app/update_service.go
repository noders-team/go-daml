package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunUpdateService(cl *client.DamlBindingClient) {
	party := getAvailableParty(cl)
	getUpdatesReq := &model.GetUpdatesRequest{
		BeginExclusive: 0,
		Filter: &model.TransactionFilter{
			FiltersByParty: map[string]*model.Filters{
				party: {
					Inclusive: &model.InclusiveFilters{
						TemplateFilters: []*model.TemplateFilter{},
					},
				},
			},
		},
		Verbose: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	responseCh, errCh := cl.UpdateService.GetUpdates(ctx, getUpdatesReq)

	updateCount := 0
	for {
		select {
		case response, ok := <-responseCh:
			if !ok {
				log.Info().Int("totalUpdates", updateCount).Msg("updates stream completed")
				return
			}
			if response != nil {
				updateCount++
				var updateID, updateType string
				var offset int64

				if response.Update.Transaction != nil {
					updateID = response.Update.Transaction.UpdateID
					offset = response.Update.Transaction.Offset
					updateType = "transaction"
				} else if response.Update.Reassignment != nil {
					updateID = response.Update.Reassignment.UpdateID
					offset = response.Update.Reassignment.Offset
					updateType = "reassignment"
				} else if response.Update.OffsetCheckpoint != nil {
					offset = response.Update.OffsetCheckpoint.Offset
					updateType = "checkpoint"
				}

				log.Info().
					Str("updateId", updateID).
					Int64("offset", offset).
					Str("updateType", updateType).
					Msg("received update")

				if updateCount >= 5 {
					log.Info().Msg("received enough updates, stopping stream")
					cancel()
					return
				}
			}
		case err := <-errCh:
			if err != nil {
				log.Warn().Err(err).Msg("updates stream error")
				return
			}
		case <-ctx.Done():
			log.Info().Int("totalUpdates", updateCount).Msg("updates stream timeout")
			return
		}
	}
}

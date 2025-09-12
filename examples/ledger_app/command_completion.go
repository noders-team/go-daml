package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunCommandCompletion(cl *client.DamlBindingClient) {
	userID, party := getAvailableUserAndParty(cl)
	completionReq := &model.CompletionStreamRequest{
		UserID:         userID,
		Parties:        []string{party},
		BeginExclusive: 0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	responseCh, errCh := cl.CommandCompletion.CompletionStream(ctx, completionReq)

	completionCount := 0
	for {
		select {
		case response, ok := <-responseCh:
			if !ok {
				log.Info().Int("totalCompletions", completionCount).Msg("completion stream completed")
				return
			}
			if response != nil {
				completionCount++
				log.Info().
					Interface("response", response.Response).
					Msg("received completion")

				if completionCount >= 3 {
					log.Info().Msg("received enough completions, stopping stream")
					cancel()
					return
				}
			}
		case err := <-errCh:
			if err != nil {
				log.Warn().Err(err).Msg("completion stream error")
				return
			}
		case <-ctx.Done():
			log.Info().Int("totalCompletions", completionCount).Msg("completion stream timeout")
			return
		}
	}
}

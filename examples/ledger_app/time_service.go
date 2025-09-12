package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunTimeService(cl *client.DamlBindingClient) {
	getTimeReq := &model.GetTimeRequest{}

	timeResp, err := cl.TimeService.GetTime(context.Background(), getTimeReq)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get time (may not be supported in production ledgers)")
	} else {
		log.Info().
			Time("currentTime", timeResp.CurrentTime).
			Msg("current ledger time")

		setTimeReq := &model.SetTimeRequest{
			CurrentTime: timeResp.CurrentTime,
			NewTime:     timeResp.CurrentTime.Add(1 * time.Hour),
		}

		_, err := cl.TimeService.SetTime(context.Background(), setTimeReq)
		if err != nil {
			log.Warn().Err(err).Msg("failed to set time (may not be supported in production ledgers)")
		} else {
			log.Info().
				Time("newTime", setTimeReq.NewTime).
				Msg("set new ledger time")

			newTimeResp, err := cl.TimeService.GetTime(context.Background(), getTimeReq)
			if err != nil {
				log.Warn().Err(err).Msg("failed to get updated time")
			} else {
				log.Info().
					Time("updatedTime", newTimeResp.CurrentTime).
					Msg("confirmed updated ledger time")
			}
		}
	}
}

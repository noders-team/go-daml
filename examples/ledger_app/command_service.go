package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunCommandService(cl *client.DamlBindingClient) {
	now := time.Now()
	userID, party := getAvailableUserAndParty(cl)
	submitReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID:   "example-workflow-" + time.Now().Format("20060102150405"),
			UserID:       userID,
			CommandID:    "cmd-" + time.Now().Format("20060102150405"),
			ActAs:        []string{party},
			SubmissionID: "sub-" + time.Now().Format("20060102150405"),
			DeduplicationPeriod: model.DeduplicationDuration{
				Duration: 60 * time.Second,
			},
			MinLedgerTimeAbs: &now,
			Commands:         []*model.Command{},
		},
	}

	response, err := cl.CommandService.SubmitAndWait(context.Background(), submitReq)
	if err != nil {
		log.Warn().Err(err).Msg("failed to submit and wait (expected with empty commands)")
	} else {
		log.Info().
			Str("updateId", response.UpdateID).
			Int64("completionOffset", response.CompletionOffset).
			Msg("command submitted and completed")
	}
}

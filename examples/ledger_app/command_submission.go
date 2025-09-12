package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunCommandSubmission(cl *client.DamlBindingClient) {
	now := time.Now()
	userID, party := getAvailableUserAndParty(cl)
	submitReq := &model.SubmitRequest{
		Commands: &model.Commands{
			WorkflowID:   "submission-workflow-" + time.Now().Format("20060102150405"),
			UserID:       userID,
			CommandID:    "sub-cmd-" + time.Now().Format("20060102150405"),
			ActAs:        []string{party},
			SubmissionID: "submission-" + time.Now().Format("20060102150405"),
			DeduplicationPeriod: model.DeduplicationDuration{
				Duration: 60 * time.Second,
			},
			MinLedgerTimeAbs: &now,
			Commands:         []*model.Command{},
		},
	}

	response, err := cl.CommandSubmission.Submit(context.Background(), submitReq)
	if err != nil {
		log.Warn().Err(err).Msg("failed to submit command (expected with empty commands)")
	} else {
		log.Info().
			Interface("response", response).
			Msg("command submitted")
	}
}

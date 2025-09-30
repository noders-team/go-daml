package main

import (
	"context"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunUsersManagement(cl *client.DamlBindingClient) {
	users, err := cl.UserMng.ListUsers(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list users")
	}
	for _, u := range users {
		log.Info().Interface("user", u).Msg("received user details")
	}

	user, err := cl.UserMng.GetUser(context.Background(), "participant_admin")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get user")
	}

	log.Info().Interface("user", user).Msg("single user details")

	userRights, err := cl.UserMng.ListUserRights(context.Background(), user.ID)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list user rights")
	}
	for _, r := range userRights {
		log.Info().Interface("right", r).Msg("user rights")
	}

	newRights := make([]*model.Right, 0)
	newRights = append(newRights, &model.Right{Type: model.CanReadAs{Party: "app_provider_localnet-localparty-1"}})

	updatedRights, err := cl.UserMng.GrantUserRights(context.Background(), user.ID, "", newRights)
	if err != nil {
		log.Warn().Err(err).Msg("failed to grant user rights")
	}
	for _, r := range updatedRights {
		log.Info().Interface("right", r).Msg("user rights after grant")
	}

	updatedRights, err = cl.UserMng.RevokeUserRights(context.Background(), user.ID, newRights)
	if err != nil {
		log.Warn().Err(err).Msg("failed to revoke user rights")
	}
	for _, r := range updatedRights {
		log.Info().Interface("right", r).Msg("user rights after revoke")
	}
}

package main

import (
	"context"
	"os"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunIdentityProvider(cl *client.DamlBindingClient) {
	configs, err := cl.IdentityProviderMng.ListIdentityProviderConfigs(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list identity provider configs")
	}
	for _, cfg := range configs {
		log.Info().Interface("config", cfg).Msg("identity provider config")
	}

	newConfig := &model.IdentityProviderConfig{
		IdentityProviderID: "test-provider-" + os.Getenv("USER"),
		IsDeactivated:      false,
		Issuer:             "https://example.com",
		JwksURL:            "https://example.com/.well-known/jwks.json",
		Audience:           "https://daml.network",
	}

	createdConfig, err := cl.IdentityProviderMng.CreateIdentityProviderConfig(context.Background(), newConfig)
	if err != nil {
		log.Error().Err(err).Msg("create identity provider error")
	} else {
		log.Info().Interface("config", createdConfig).Msg("created identity provider")
	}

	if createdConfig != nil {
		retrievedConfig, err := cl.IdentityProviderMng.GetIdentityProviderConfig(context.Background(), createdConfig.IdentityProviderID)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get identity provider config")
		}
		log.Info().Interface("config", retrievedConfig).Msg("retrieved identity provider")

		updatedConfig := &model.IdentityProviderConfig{
			IdentityProviderID: createdConfig.IdentityProviderID,
			IsDeactivated:      false,
			Issuer:             "https://updated.example.com",
			JwksURL:            "https://updated.example.com/.well-known/jwks.json",
			Audience:           "https://daml.network.updated",
		}

		finalConfig, err := cl.IdentityProviderMng.UpdateIdentityProviderConfig(context.Background(), updatedConfig, []string{"issuer", "jwks_url", "audience"})
		if err != nil {
			log.Error().Err(err).Msg("update identity provider error")
		} else {
			log.Info().Interface("config", finalConfig).Msg("updated identity provider")
		}

		err = cl.IdentityProviderMng.DeleteIdentityProviderConfig(context.Background(), createdConfig.IdentityProviderID)
		if err != nil {
			log.Error().Err(err).Msg("delete identity provider error")
		} else {
			log.Info().Msg("identity provider deleted successfully")
		}
	}

	finalConfigs, err := cl.IdentityProviderMng.ListIdentityProviderConfigs(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list final identity provider configs")
	}
	log.Info().Msg("final identity provider configs:")
	for _, cfg := range finalConfigs {
		log.Info().Interface("config", cfg).Msg("identity provider config")
	}
}

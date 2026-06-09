package main

import (
	"context"
	"os"

	"github.com/noders-team/go-daml/pkg/auth"
	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:8080"
	}

	kcCfg := auth.KeycloakConfig{
		OIDCURL:      os.Getenv("KEYCLOAK_OIDC_URL"),
		TokenURL:     os.Getenv("KEYCLOAK_TOKEN_URL"),
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		Audience:     os.Getenv("KEYCLOAK_AUDIENCE"),
	}
	if !kcCfg.Enabled() {
		log.Fatal().Msg("set KEYCLOAK_CLIENT_ID, KEYCLOAK_CLIENT_SECRET and KEYCLOAK_OIDC_URL or KEYCLOAK_TOKEN_URL")
	}

	kc, err := auth.NewKeycloakTokenProvider(kcCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create keycloak token provider")
	}

	builder := client.NewDamlClient(grpcAddress, kc)
	if os.Getenv("GRPC_TLS") == "true" {
		builder = builder.WithTLSConfig(client.TlsConfig{})
	}

	cl, err := builder.Build(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build DAML client")
	}
	defer cl.Close()

	if err := cl.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed to ping ledger")
	}

	version, err := cl.VersionService.GetLedgerAPIVersion(context.Background(), &model.GetLedgerAPIVersionRequest{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get ledger API version")
	}
	log.Info().Str("version", version.Version).Msg("authenticated via keycloak")

	users, err := cl.UserMng.ListUsers(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list users")
	}
	log.Info().Int("users", len(users)).Msg("listed users using keycloak-issued token")
}

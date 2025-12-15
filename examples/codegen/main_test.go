package codegen_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/testutil"
	"github.com/rs/zerolog/log"
)

var cl *client.DamlBindingClient

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel()

	if err := testutil.Setup(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to setup test environment")
	}

	cl = testutil.GetClient()

	code := m.Run()

	testutil.Teardown()

	os.Exit(code)
}

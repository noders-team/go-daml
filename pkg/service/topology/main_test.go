package topology_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/noders-team/go-daml/pkg/testutil"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err := testutil.Setup(ctx)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	testutil.Teardown()
	os.Exit(code)
}

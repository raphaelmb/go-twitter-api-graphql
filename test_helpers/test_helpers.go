package test_helpers

import (
	"context"
	"testing"

	"github.com/raphaelmb/go-twitter-api-graphql/postgres"
	"github.com/stretchr/testify/require"
)

func TearDownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}

package test_helpers

import (
	"context"
	"testing"

	"github.com/raphaelmb/go-twitter-api-graphql"
	"github.com/raphaelmb/go-twitter-api-graphql/faker"
	"github.com/raphaelmb/go-twitter-api-graphql/postgres"
	"github.com/stretchr/testify/require"
)

func TearDownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}

func CreateUser(ctx context.Context, t *testing.T, userRepo twitter.UserRepo) twitter.User {
	t.Helper()

	user, err := userRepo.Create(ctx, twitter.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password,
	})
	require.NoError(t, err)

	return user
}

func LoginUser(ctx context.Context, t *testing.T, user twitter.User) context.Context {
	t.Helper()

	return twitter.PutUserIDIntoContext(ctx, user.ID)
}

//go:build integration

package domain

import (
	"context"
	"testing"

	"github.com/raphaelmb/go-twitter-api-graphql"
	"github.com/raphaelmb/go-twitter-api-graphql/faker"
	"github.com/raphaelmb/go-twitter-api-graphql/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTweetService_Create(t *testing.T) {
	t.Run("non auth user cannot create a tweet", func(t *testing.T) {
		ctx := context.Background()
		_, err := tweetService.Create(ctx, twitter.CreateTweetInput{Body: "hello"})
		require.ErrorIs(t, err, twitter.ErrUnauthenticated)
	})

	t.Run("cana create a tweet", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		input := twitter.CreateTweetInput{
			Body: faker.RandStr(100),
		}

		tweet, err := tweetService.Create(ctx, input)
		require.NoError(t, err)

		require.NotEmpty(t, tweet.ID)
		require.Equal(t, input.Body, tweet.Body)
		require.Equal(t, currentUser.ID, tweet.UserID)
		require.NotEmpty(t, tweet.CreatedAt)
	})
}

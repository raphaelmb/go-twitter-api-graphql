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

	t.Run("cannot create invalid tweet", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)
		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := tweetService.Create(ctx, twitter.CreateTweetInput{Body: "h"})
		require.ErrorIs(t, err, twitter.ErrValidation)
	})

	t.Run("can create a tweet", func(t *testing.T) {
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

func TestIntegrationTweetService_GetByID(t *testing.T) {
	t.Run("can get a tweet by id", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		user := test_helpers.CreateUser(ctx, t, userRepo)
		existingTweet := test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)

		tweet, err := tweetService.GetByID(ctx, existingTweet.ID)
		require.NoError(t, err)

		require.Equal(t, existingTweet.ID, tweet.ID)
		require.Equal(t, existingTweet.Body, tweet.Body)
	})

	t.Run("return error not found if the tweet does not exists", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		_, err := tweetService.GetByID(ctx, faker.UUID())
		require.ErrorIs(t, err, twitter.ErrNotFound)
	})

	t.Run("return error invalid uuid", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		_, err := tweetService.GetByID(ctx, "123")
		require.ErrorIs(t, err, twitter.ErrInvaludUUID)
	})

}

func TestIntegrationTweetService_All(t *testing.T) {
	t.Run("return all tweets", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		user := test_helpers.CreateUser(ctx, t, userRepo)
		test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)

		tweets, err := tweetService.All(ctx)
		require.NoError(t, err)

		require.Len(t, tweets, 3)
	})
}

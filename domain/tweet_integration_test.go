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

func TestIntergrationTweetService_Delete(t *testing.T) {
	t.Run("not auth user cannot delete a tweet", func(t *testing.T) {
		ctx := context.Background()
		err := tweetService.Delete(ctx, faker.UUID())
		require.ErrorIs(t, err, twitter.ErrUnauthenticated)
	})

	t.Run("cannot delete tweet if not the owner", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)
		otherUser := test_helpers.CreateUser(ctx, t, userRepo)
		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		tweet := test_helpers.CreateTweet(ctx, t, tweetRepo, otherUser.ID)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := tweetRepo.GetByID(ctx, tweet.ID)
		require.NoError(t, err)

		err = tweetService.Delete(ctx, tweet.ID)
		require.ErrorIs(t, err, twitter.ErrForbidden)

		_, err = tweetRepo.GetByID(ctx, tweet.ID)
		require.NoError(t, err)
	})

	t.Run("can delete a tweet", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)
		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		tweet := test_helpers.CreateTweet(ctx, t, tweetRepo, currentUser.ID)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := tweetRepo.GetByID(ctx, tweet.ID)
		require.NoError(t, err)

		err = tweetService.Delete(ctx, tweet.ID)
		require.NoError(t, err)

		_, err = tweetRepo.GetByID(ctx, tweet.ID)
		require.ErrorIs(t, err, twitter.ErrNotFound)
	})
}

func TestIntegrationTweetService_CreateReply(t *testing.T) {
	t.Run("not auth user cannot create reply to a tweet", func(t *testing.T) {
		ctx := context.Background()
		_, err := tweetService.CreateReply(ctx, faker.UUID(), twitter.CreateTweetInput{
			Body: faker.RandStr(20),
		})
		require.ErrorIs(t, err, twitter.ErrUnauthenticated)
	})

	t.Run("cannot create a reply to a not found tweet", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := tweetService.CreateReply(ctx, faker.UUID(), twitter.CreateTweetInput{
			Body: faker.RandStr(20),
		})
		require.ErrorIs(t, err, twitter.ErrNotFound)
	})

	t.Run("can create a reply to a tweet", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TearDownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		tweet := test_helpers.CreateTweet(ctx, t, tweetRepo, currentUser.ID)
		input := twitter.CreateTweetInput{
			Body: faker.RandStr(20),
		}

		reply, err := tweetService.CreateReply(ctx, tweet.ID, input)
		require.NoError(t, err)

		require.NotEmpty(t, reply.ID)
		require.Equal(t, input.Body, reply.Body)
		require.Equal(t, currentUser.ID, reply.UserID)
		require.Equal(t, tweet.ID, *reply.ParentID)
		require.NotEmpty(t, reply.CreatedAt)
	})
}

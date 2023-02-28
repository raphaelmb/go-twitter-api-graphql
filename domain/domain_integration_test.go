//go:build integration
// +build integration

package domain

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/raphaelmb/go-twitter-api-graphql"
	"github.com/raphaelmb/go-twitter-api-graphql/config"
	"github.com/raphaelmb/go-twitter-api-graphql/jwt"
	"github.com/raphaelmb/go-twitter-api-graphql/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	conf             *config.Config
	db               *postgres.DB
	authTokenService twitter.AuthTokenService
	tweetService     twitter.TweetService
	authService      twitter.AuthService
	userRepo         twitter.UserRepo
	tweetRepo        twitter.TweetRepo
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	config.LoadEnv(".env.test")

	passwordCost = bcrypt.MinCost

	conf = config.New()
	db = postgres.New(ctx, conf)
	defer db.Close()

	if err := db.Drop(); err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepo = postgres.NewUserRepo(db)
	tweetRepo = postgres.NewTweetRepo(db)

	authTokenService = jwt.NewTokenService(conf)

	authService = NewAuthService(userRepo, authTokenService)
	tweetService = NewTweetService(tweetRepo)

	os.Exit(m.Run())
}

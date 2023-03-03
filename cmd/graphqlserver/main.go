package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raphaelmb/go-twitter-api-graphql/config"
	"github.com/raphaelmb/go-twitter-api-graphql/domain"
	"github.com/raphaelmb/go-twitter-api-graphql/graph"
	"github.com/raphaelmb/go-twitter-api-graphql/jwt"
	"github.com/raphaelmb/go-twitter-api-graphql/postgres"
)

func main() {
	ctx := context.Background()
	config.LoadEnv(".env")
	conf := config.New()
	db := postgres.New(ctx, conf)

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Timeout(time.Second * 60))

	// repo
	userRepo := postgres.NewUserRepo(db)
	tweetRepo := postgres.NewTweetRepo(db)

	// service
	authTokenService := jwt.NewTokenService(conf)
	authService := domain.NewAuthService(userRepo, authTokenService)
	tweetService := domain.NewTweetService(tweetRepo)

	router.Use(authMiddleware(authTokenService))
	router.Handle("/", playground.Handler("twitter clone", "/query"))
	router.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					AuthService:  authService,
					TweetService: tweetService,
				},
			},
		),
	))

	log.Fatal(http.ListenAndServe(":8080", router))
}

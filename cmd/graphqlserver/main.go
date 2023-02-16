package main

import (
	"context"
	"fmt"
	"log"

	"github.com/raphaelmb/go-twitter-api-graphql/config"
	"github.com/raphaelmb/go-twitter-api-graphql/postgres"
)

func main() {
	ctx := context.Background()
	conf := config.New()
	db := postgres.New(ctx, conf)

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("All good")
}

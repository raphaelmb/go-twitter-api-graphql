package postgres

import (
	"context"

	"github.com/raphaelmb/go-twitter-api-graphql"
)

type TweetRepo struct {
	db *DB
}

func NewTweetRepo(db *DB) *TweetRepo {
	return &TweetRepo{
		db: db,
	}
}

func (tr *TweetRepo) All(ctx context.Context) ([]twitter.Tweet, error) {
	panic("not implemented")
}

func (tr *TweetRepo) Create(ctx context.Context, input twitter.Tweet) (twitter.Tweet, error) {
	panic("not implemented")
}

func (tr *TweetRepo) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	panic("not implemented")
}

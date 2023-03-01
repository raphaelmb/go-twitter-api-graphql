package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
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

func (tr *TweetRepo) Create(ctx context.Context, tweet twitter.Tweet) (twitter.Tweet, error) {
	tx, err := tr.db.Pool.Begin(ctx)
	if err != nil {
		return twitter.Tweet{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	t, err := createTweet(ctx, tx, tweet)
	if err != nil {
		return twitter.Tweet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return twitter.Tweet{}, fmt.Errorf("error commiting: %v", err)
	}

	return t, nil
}

func createTweet(ctx context.Context, tx pgx.Tx, tweet twitter.Tweet) (twitter.Tweet, error) {
	query := `INSERT INTO tweets (body, user_id) VALUES ($1, $2) RETURNING *;`

	t := twitter.Tweet{}

	if err := pgxscan.Get(ctx, tx, &t, query, tweet.Body, tweet.UserID); err != nil {
		return twitter.Tweet{}, fmt.Errorf("error insert: %v", err)
	}

	return t, nil
}

func (tr *TweetRepo) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	panic("not implemented")
}

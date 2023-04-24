package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/raphaelmb/go-twitter-api-graphql"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) Create(ctx context.Context, user twitter.User) (twitter.User, error) {
	tx, err := ur.db.Pool.Begin(ctx)
	if err != nil {
		return twitter.User{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	u, err := createUser(ctx, tx, user)
	if err != nil {
		return twitter.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return twitter.User{}, fmt.Errorf("error commiting: %v", err)
	}

	return u, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user twitter.User) (twitter.User, error) {
	query := `INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING *;`

	u := twitter.User{}
	if err := pgxscan.Get(ctx, tx, &u, query, user.Email, user.Username, user.Password); err != nil {
		return twitter.User{}, fmt.Errorf("error insert: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (twitter.User, error) {
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1;`
	u := twitter.User{}
	if err := pgxscan.Get(ctx, ur.db.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (twitter.User, error) {
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1;`
	u := twitter.User{}
	if err := pgxscan.Get(ctx, ur.db.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id string) (twitter.User, error) {
	query := `SELECT * FROM users WHERE id = $1 LIMIT 1;`
	u := twitter.User{}
	if err := pgxscan.Get(ctx, ur.db.Pool, &u, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByIds(ctx context.Context, ids []string) ([]twitter.User, error) {
	return getUsersByIds(ctx, ur.db.Pool, ids)
}

func getUsersByIds(ctx context.Context, q pgxscan.Querier, ids []string) ([]twitter.User, error) {
	query := `SELECT * FROM users WHERE id = ANY($1);`

	var uu []twitter.User

	if err := pgxscan.Select(ctx, q, &uu, query, ids); err != nil {
		return nil, fmt.Errorf("error get users by ids: %+v", err)
	}

	return uu, nil
}

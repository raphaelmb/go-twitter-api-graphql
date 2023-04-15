package twitter

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	TweetMinLength = 2
	TweetMaxLength = 250
)

type CreateTweetInput struct {
	Body string
}

func (ti *CreateTweetInput) Sanitize() {
	ti.Body = strings.TrimSpace(ti.Body)
}

func (ti CreateTweetInput) Validate() error {
	if len(ti.Body) < TweetMinLength {
		return fmt.Errorf("%w: body not long enough, (%d) characters at least", ErrValidation, TweetMinLength)
	}

	if len(ti.Body) > TweetMaxLength {
		return fmt.Errorf("%w: body too long, (%d) characters at most", ErrValidation, TweetMaxLength)
	}

	return nil
}

type Tweet struct {
	ID        string
	Body      string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Tweet) CanDelete(user User) bool {
	return t.UserID == user.ID
}

type TweetService interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, input CreateTweetInput) (Tweet, error)
	GetByID(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}

type TweetRepo interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, tweet Tweet) (Tweet, error)
	GetByID(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}

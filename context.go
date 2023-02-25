package twitter

import (
	"context"
)

type contextKey string

var (
	contextAuthIDKey contextKey = "currentUserId"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	if ctx.Value(contextAuthIDKey) == nil {
		return "", ErrNoUserIdInContext
	}
	userID, ok := ctx.Value(contextAuthIDKey).(string)
	if !ok {
		return "", ErrNoUserIdInContext
	}

	return userID, nil
}

func PutUserIDIntoContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextAuthIDKey, id)
}

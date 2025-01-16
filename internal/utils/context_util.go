package utils

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	userIDKey = contextKey("userID")
)

func SetContextUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetContextUserID(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	userID, ok := v.(string)
	return userID, ok
}

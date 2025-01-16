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

func GetContextUserID(ctx context.Context) (uuid.UUID, bool) {
	val := ctx.Value(userIDKey)

	if val == nil {
		return uuid.UUID{}, false
	}

	userID, ok := val.(uuid.UUID)

	return userID, ok
}

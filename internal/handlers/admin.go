package handlers

import (
	"context"

	"github.com/google/uuid"
)

func IsUserAdmin(userId uuid.UUID, context context.Context, cfg *HandlerSiteConfig) (bool, error) {
	adminCheck, err := cfg.DbQueries.IsAdmin(context, userId.String())

	if err != nil {
		return false, err
	}

	if !adminCheck.IsAdmin {
		return false, nil
	}

	return true, nil

}

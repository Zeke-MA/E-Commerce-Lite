package handlers

import (
	"context"

	"github.com/google/uuid"
)

func (cfg *HandlerSiteConfig) IsUserAdmin(context context.Context, userId uuid.UUID) (bool, error) {
	adminCheck, err := cfg.DbQueries.IsAdmin(context, userId)

	if err != nil {
		return false, err
	}

	if !adminCheck.IsAdmin {
		return false, nil
	}

	return true, nil

}

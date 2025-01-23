package admindb

import (
	"context"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/config"
	"github.com/google/uuid"
)

func IsUserAdmin(userId uuid.UUID, ctx context.Context, cfg *config.SiteConfig) (bool, error) {
	adminCheck, err := cfg.DbQueries.IsAdmin(ctx, userId)

	if err != nil {
		return false, err
	}

	if !adminCheck.IsAdmin {
		return false, nil
	}

	return true, nil

}

package handlers

import (
	"context"
	"net/http"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
	"github.com/google/uuid"
)

func isUserAdmin(userId uuid.UUID, context context.Context, cfg *HandlerSiteConfig) (bool, error) {
	adminCheck, err := cfg.DbQueries.IsAdmin(context, userId.String())

	if err != nil {
		return false, err
	}

	if !adminCheck.IsAdmin {
		return false, nil
	}

	return true, nil

}

func (cfg *HandlerSiteConfig) AddProduct(w http.ResponseWriter, r *http.Request) {
	server.RespondWithJSON(w, http.StatusNoContent, nil)
}

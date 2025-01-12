package handlers

import (
	"net/http"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/auth"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
)

func (cfg *HandlerSiteConfig) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	// if the responding user exists and has a valid access token generate new access token and respond back
	type TokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Access token missing", err)
		return
	}

	_, err = auth.ValidateJWT(bearerToken, cfg.JWTSecret)

	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Unauthorized token request", err)
		return
	}

}

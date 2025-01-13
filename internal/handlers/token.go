package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/auth"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
)

func (cfg *HandlerSiteConfig) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	type TokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	tokenRequest := TokenRequest{}
	err := decoder.Decode(&tokenRequest)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	validRefreshToken, err := cfg.DbQueries.RefreshTokenValid(r.Context(), tokenRequest.RefreshToken)

	if err != nil {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), err)
		return
	}

	if validRefreshToken.Token != tokenRequest.RefreshToken {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	newRefreshToken, err := auth.GenerateJWT(validRefreshToken.UserID, cfg.JWTSecret, cfg.JWTExpiry)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	type TokenResponse struct {
		AccessToken string
	}

	tokenResponse := TokenResponse{AccessToken: newRefreshToken}

	server.RespondWithJSON(w, http.StatusOK, tokenResponse)

}

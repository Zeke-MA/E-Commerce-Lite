package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/auth"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
)

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	AccessToken string
}

func (cfg *HandlerSiteConfig) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	tokenRequest := RefreshToken{}
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

	tokenResponse := AccessToken{AccessToken: newRefreshToken}

	server.RespondWithJSON(w, http.StatusOK, tokenResponse)

}

func (cfg *HandlerSiteConfig) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	bearerAuth, err := auth.GetBearerToken(r.Header)

	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	_, err = auth.ValidateJWT(bearerAuth, cfg.JWTSecret)

	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	tokenRequest := RefreshToken{}
	err = decoder.Decode(&tokenRequest)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	tokenRevokeDb, err := cfg.DbQueries.RevokeRefreshToken(r.Context(), tokenRequest.RefreshToken)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	affected, _ := tokenRevokeDb.RowsAffected()

	if affected == 0 {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), err)
		return
	}

	server.RespondWithJSON(w, http.StatusNoContent, nil)

}

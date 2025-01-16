package middleware

import (
	"context"
	"net/http"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/auth"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/utils"
)

func (cfg *MiddlewareSiteConfig) CheckUserValidated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearerToken, err := auth.GetBearerToken(r.Header)

		if err != nil {
			server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
			return
		}

		requestUserID, err := auth.ValidateJWT(bearerToken, cfg.JWTSecret)

		if err != nil {
			server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
			return
		}

		ctx := context.WithValue(r.Context(), utils.SetContextUserID(r.Context(), requestUserID), requestUserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

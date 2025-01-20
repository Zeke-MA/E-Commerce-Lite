package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
)

func (cfg *MiddlewareSiteConfig) LogIncomingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var i interface{}
		// Acknowledging token is in header logging. Future todo - extract without httputil
		dumpReq, err := httputil.DumpRequest(r, false)

		if err != nil {
			cfg.Logger.Error("Reading Response", slog.String("error", fmt.Sprintf("%v", err)))
			server.RespondWithError(w, http.StatusBadRequest, string(server.MsgBadRequest), err)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)

		if len(bodyBytes) == 0 {
			cfg.Logger.Info("Request has no body", slog.String("details", string(dumpReq)))
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			cfg.Logger.Error("Reading Response Body", slog.String("error", fmt.Sprintf("%v", err)))
			server.RespondWithError(w, http.StatusBadRequest, string(server.MsgBadRequest), err)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.Unmarshal(bodyBytes, &i)

		if err != nil {
			cfg.Logger.Error("Parsing Response Body", slog.String("error", fmt.Sprintf("%v", err)))
			server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
			return
		}

		obj, ok := i.(map[string]interface{})

		if !ok {
			cfg.Logger.Error("Wrong Request Format", slog.String("error", fmt.Sprintf("%v", err)))
			server.RespondWithError(w, http.StatusBadRequest, string(server.MsgBadRequest), err)
			return
		}

		delete(obj, "password")
		delete(obj, "access_token")
		delete(obj, "refresh_token")

		requestBody, err := json.Marshal(obj)

		if err != nil {
			log.Printf("Error marshalling map: %v \n", err)
			server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
			return
		}

		cfg.Logger.Info("Requestor Details", slog.String("details", string(dumpReq)))
		cfg.Logger.Info("Request body", slog.String("request body", string(requestBody)))

		next.ServeHTTP(w, r)
	})
}

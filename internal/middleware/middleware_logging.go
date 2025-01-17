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
)

func (cfg *MiddlewareSiteConfig) LogIncomingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var i interface{}

		req, err := httputil.DumpRequest(r, false)

		if err != nil {
			log.Println(err)
		}

		bodyBytes, err := io.ReadAll(r.Body)

		if err != nil {
			cfg.Logger.Error("Reading Response Body", slog.String("error", fmt.Sprintf("%v", err)))
			next.ServeHTTP(w, r)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.Unmarshal(bodyBytes, &i)

		if err != nil {
			cfg.Logger.Error("Parsing Response Body", slog.String("error", fmt.Sprintf("%v", err)))
		}

		obj, ok := i.(map[string]interface{})

		if !ok {
			cfg.Logger.Error("Wrong Request Format", slog.String("error", fmt.Sprintf("%v", err)))
		}

		delete(obj, "password")
		delete(obj, "access_token")
		delete(obj, "refresh_token")

		jsonBytes, err := json.Marshal(obj)

		if err != nil {
			log.Printf("Error marshalling map: %v \n", err)
		}

		cfg.Logger.Info("Requestor Details", slog.String("details", string(req)))
		cfg.Logger.Info("Request body", slog.String("request body", string(jsonBytes)))

		next.ServeHTTP(w, r)
	})
}

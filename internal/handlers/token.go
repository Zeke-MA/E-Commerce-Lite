package handlers

import "net/http"

func (cfg *HandlerSiteConfig) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	// if the responding user exists and has a valid access token generate new access token and respond back
}

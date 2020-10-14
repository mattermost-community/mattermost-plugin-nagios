package main

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	token := p.getConfiguration().Token

	if token == "" {
		http.Error(w, "This functionality is not configured.", http.StatusNotImplemented)
		return
	}
	if token != r.Header.Get("Authorization") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
}

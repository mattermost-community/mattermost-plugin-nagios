package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/mattermost/mattermost-plugin-starter-template/internal/watcher"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func formatChange(change watcher.Change) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("**%s** has been modified", change.Name))
	b.WriteString(" (-previous +actual):\n\n")
	b.WriteString("```")

	if utf8.RuneCountInString(change.Diff) > 16077 {
		b.WriteString("File has been changed, but the diff is too long.")
	} else {
		b.WriteString(change.Diff)
	}

	b.WriteString("```")

	return b.String()
}

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

	var change watcher.Change

	if err := json.NewDecoder(r.Body).Decode(&change); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: "channelID",
		Message:   formatChange(change),
	}

	if _, err := p.API.CreatePost(post); err != nil {
		p.API.LogError("CreatePost", logErrorKey, err)
	}
}

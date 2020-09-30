package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/ulumuri/go-nagios/nagios"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	client *nagios.Client

	botUserID string
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (p *Plugin) setDefaultKV(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err := p.API.KVSet(key, b); err != nil {
		return fmt.Errorf("p.API.KVSet: %w", err)
	}

	return nil
}

var defaultKVStore = map[string]interface{}{
	logsLimitKey:       defaultLogsLimit,
	logsStartTimeKey:   defaultLogsStartTime,
	reportFrequencyKey: defaultReportFrequency,
}

func (p *Plugin) storeDefaultKV() error {
	for key, val := range defaultKVStore {
		if err := p.setDefaultKV(key, val); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) OnActivate() error {
	botUserID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "nagios",
		DisplayName: "Nagios",
		Description: "Created by the Nagios Plugin.",
	})
	if err != nil {
		return fmt.Errorf("p.Helpers.EnsureBot: %w", err)
	}

	p.botUserID = botUserID

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return fmt.Errorf("p.API.GetBundlePath: %w", err)
	}

	profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "nagios.png"))
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %w", err)
	}

	if err := p.API.SetProfileImage(botUserID, profileImage); err != nil {
		return fmt.Errorf("p.API.SetProfileImage: %w", err)
	}

	if err := p.API.RegisterCommand(nagiosCommand); err != nil {
		return fmt.Errorf("p.API.RegisterCommand: %w", err)
	}

	if err := p.storeDefaultKV(); err != nil {
		return fmt.Errorf("p.storeDefaultKV: %w", err)
	}

	return nil
}

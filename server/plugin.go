package main

import (
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

func (p *Plugin) OnActivate() error {
	config := p.getConfiguration()

	if err := config.isValid(); err != nil {
		return err
	}

	c, err := nagios.NewClient(http.DefaultClient, config.NagiosURL)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	p.client = c

	botUserID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "nagios",
		DisplayName: "Nagios",
		Description: "Created by the Nagios Plugin.",
	})
	if err != nil {
		return fmt.Errorf("EnsureBot: %w", err)
	}

	p.botUserID = botUserID

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return fmt.Errorf("GetBundlePath: %w", err)
	}

	profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "nagios.png"))
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	if err := p.API.SetProfileImage(botUserID, profileImage); err != nil {
		return fmt.Errorf("SetProfileImage: %w", err)
	}

	if err := p.API.RegisterCommand(nagiosCommand); err != nil {
		return fmt.Errorf("RegisterCommand: %w", err)
	}

	return nil
}

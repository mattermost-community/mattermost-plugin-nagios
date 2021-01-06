package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/ulumuri/go-nagios/nagios"
)

// Plugin implements the interface expected by the Mattermost server to
// communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult
	// getConfiguration and setConfiguration for usage.
	configuration *configuration

	client *nagios.Client

	botUserID string

	commandHandlers map[string]commandHandlerFunc
}

func (p *Plugin) storeInitialKV() error {
	initials := map[string]int{
		setLogsLimitKey:       p.getConfiguration().InitialLogsLimit,
		setLogsStartTimeKey:   p.getConfiguration().InitialLogsStartTime,
		setReportFrequencyKey: p.getConfiguration().InitialReportFrequency,
	}

	for key, val := range initials {
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}

		if err := p.API.KVSet(key, b); err != nil {
			return fmt.Errorf("p.API.KVSet: %w", err)
		}
	}

	return nil
}

func (p *Plugin) getProfileImage() ([]byte, error) {
	path, err := p.API.GetBundlePath()
	if err != nil {
		return nil, fmt.Errorf("p.API.GetBundlePath: %w", err)
	}

	// NOTICE: We don't use any of the Nagios logos to avoid legal issues.
	// Instead, we use an image resembling a part of the Nagios Core logo.
	file, err := ioutil.ReadFile(filepath.Join(path, "assets", "orbit-467260.png"))
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %w", err)
	}

	return file, nil
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

	img, err := p.getProfileImage()
	if err != nil {
		return fmt.Errorf("p.getProfileImage: %w", err)
	}

	if err := p.API.SetProfileImage(botUserID, img); err != nil {
		return fmt.Errorf("p.API.SetProfileImage: %w", err)
	}

	ico, err := command.GetIconData(p.API, "assets/orbit-467260.svg")
	if err != nil {
		return fmt.Errorf("command.GetIconData: %w", err)
	}

	if err := p.API.RegisterCommand(p.getCommand(ico)); err != nil {
		return fmt.Errorf("p.API.RegisterCommand: %w", err)
	}

	if err := p.storeInitialKV(); err != nil {
		return fmt.Errorf("p.storeInitialKV: %w", err)
	}

	go p.monitoringReportLoop()

	return nil
}

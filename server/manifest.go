// This file is automatically generated. Do not modify it manually.

package main

import (
	"encoding/json"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

var manifest model.Manifest

const manifestStr = `
{
  "id": "nagios",
  "name": "Nagios",
  "description": "Nagios plugin for Mattermost",
  "homepage_url": "https://github.com/mattermost/mattermost-plugin-nagios",
  "support_url": "https://github.com/mattermost/mattermost-plugin-nagios/issues",
  "release_notes_url": "https://github.com/mattermost/mattermost-plugin-nagios/releases/tag/v1.1.0",
  "icon_path": "assets/orbit-467260.svg",
  "version": "1.1.0",
  "min_server_version": "5.37.0",
  "server": {
    "executables": {
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "settings_schema": {
    "header": "Having problems configuring the plugin? [Check the configuration guide](https://github.com/mattermost/mattermost-plugin-nagios/#configuring-the-plugin).",
    "footer": "To report an issue, make a suggestion or a contribution, or fork your own version of the plugin, [check the repository](https://github.com/mattermost/mattermost-plugin-nagios).",
    "settings": [
      {
        "key": "NagiosURL",
        "display_name": "Nagios URL",
        "type": "text",
        "help_text": "The URL for your Nagios instance. Must start with http:// or https://.",
        "placeholder": "",
        "default": null
      },
      {
        "key": "NagiosUsername",
        "display_name": "Nagios Username",
        "type": "text",
        "help_text": "The Nagios user used to authenticate against your Nagios instance.",
        "placeholder": "",
        "default": null
      },
      {
        "key": "NagiosPassword",
        "display_name": "Nagios Password",
        "type": "text",
        "help_text": "The password of that user.",
        "placeholder": "",
        "default": null
      },
      {
        "key": "Token",
        "display_name": "Token",
        "type": "generated",
        "help_text": "The token for the configuration files watcher.",
        "placeholder": "",
        "default": null
      },
      {
        "key": "InitialLogsLimit",
        "display_name": "Initial logs limit",
        "type": "number",
        "help_text": "Limit the initial number of logs the get-logs action fetches.",
        "placeholder": "Integer",
        "default": 50
      },
      {
        "key": "InitialLogsStartTime",
        "display_name": "Initial logs start time",
        "type": "number",
        "help_text": "Specify the initial age of the oldest log the get-logs action fetches.",
        "placeholder": "Seconds",
        "default": 86400
      },
      {
        "key": "InitialReportFrequency",
        "display_name": "Initial report frequency",
        "type": "number",
        "help_text": "Set the initial frequency of system monitoring reports.",
        "placeholder": "Minutes",
        "default": 1
      }
    ]
  }
}
`

func init() {
	_ = json.NewDecoder(strings.NewReader(manifestStr)).Decode(&manifest)
}

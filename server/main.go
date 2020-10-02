package main

import (
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func main() {
	plugin.ClientMain(&Plugin{
		commandHandlers: map[string]commandHandlerFunc{
			"help":                 nil,
			"set-logs-limit":       setLogsLimit,
			"set-logs-start-time":  setLogsStartTime,
			"get-logs":             getLogs,
			"set-report-frequency": setReportFrequency,
		},
	})
}

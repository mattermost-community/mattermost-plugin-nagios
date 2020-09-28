package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func getAutoCompleteDesc(m map[string]commandHandlerFunc) string {
	var b strings.Builder

	b.WriteString("Available commands: ")

	var i int
	for k := range m {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(k)
	}

	return b.String()
}

var nagiosCommand = &model.Command{
	Trigger:          "nagios",
	AutoComplete:     true,
	AutoCompleteDesc: getAutoCompleteDesc(commandHandlers),
	AutoCompleteHint: "[command]",
	AutocompleteData: getAutocompleteData(),
	DisplayName:      "Nagios",
	Description:      "A Mattermost plugin to interact with Nagios",
}

func getAutocompleteData() *model.AutocompleteData {
	nagios := model.NewAutocompleteData("nagios", "[command]", getAutoCompleteDesc(commandHandlers))

	// GET LOGS
	getlogs := model.NewAutocompleteData("get-logs", "[alerts|notifications]", "Get logs of specific type")

	alerts := model.NewAutocompleteData("alerts", "", "Get alert logs")
	getlogs.AddCommand(alerts)

	notifications := model.NewAutocompleteData("notifications", "", "Get notification logs")
	getlogs.AddCommand(notifications)

	nagios.AddCommand(getlogs)

	// SET LOGS LIMIT
	setlogslimit := model.NewAutocompleteData("set-logs-limit", "[positive integer]", "Set max number of logs to display")
	nagios.AddCommand(setlogslimit)

	// SET LOGS START TIME
	setlogsstarttime := model.NewAutocompleteData("set-logs-start-time", "[seconds]", "Set duration of time(seconds) to display logs from")
	nagios.AddCommand(setlogsstarttime)

	// SET REPORT FREQUENCY
	setreportfrequency := model.NewAutocompleteData("set-report-frequency", "[minutes]", "Set frequency of system monitoring reports")
	nagios.AddCommand(setreportfrequency)

	return nagios
}

func parseCommandArgs(args *model.CommandArgs) (command, action string, parameters []string) {
	fields := strings.Fields(args.Command)

	if len(fields) > 0 {
		command = fields[0]
	}
	if len(fields) > 1 {
		action = fields[1]
	}

	parameters = fields[2:]

	return command, action, parameters
}

func (p *Plugin) getCommandResponse(args *model.CommandArgs, text string) *model.CommandResponse {
	p.API.SendEphemeralPost(args.UserId, &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	})
	return &model.CommandResponse{}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	command, action, parameters := parseCommandArgs(args)

	if command != "/nagios" {
		return &model.CommandResponse{}, nil
	}

	var msg string

	if f, ok := commandHandlers[action]; ok {
		msg = f(p.API, p.client, parameters)
	} else {
		msg = fmt.Sprintf("Unknown action (%s).", action)
	}

	return p.getCommandResponse(args, msg), nil
}

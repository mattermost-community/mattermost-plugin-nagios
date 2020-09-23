package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const helpText = ``

var nagiosCommand = &model.Command{
	Trigger:          "nagios",
	Description:      "A Mattermost plugin to interact with Nagios",
	DisplayName:      "Nagios",
	AutoComplete:     true,
	AutoCompleteDesc: "Available commands: get-log",
	AutoCompleteHint: "[command]",
}

func parseCommandArgs(args *model.CommandArgs) (
	command, action string,
	parameters []string) {

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
	command, action, parameteres := parseCommandArgs(args)

	if command != "/nagios" {
		return &model.CommandResponse{}, nil
	}

	if f, ok := commandHandlers[action]; ok {
		msg := f(p.client, parameteres)
		return p.getCommandResponse(args, msg), nil
	}

	msg := fmt.Sprintf("Unknown action %s", action)
	return p.getCommandResponse(args, msg), nil
}
